package services

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	config "receipt-wrangler/api/internal/env"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
)

func recognitionFixture(t *testing.T) []byte {
	t.Helper()
	bytes, err := os.ReadFile(filepath.Join(testApiRoot(), "testing", "test.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}

func recognitionSeed(t *testing.T) (RecognitionTaskService, commands.RegisterRecognitionTaskCommand) {
	t.Helper()
	ClearRolePermissionCacheForTests()
	ClearGroupRoleGrantCacheForTests()
	user, group, _ := seedMemberWithGroupRole(t, "recognition-user", []string{permissions.GroupReceiptsRead, permissions.GroupReceiptsQuickScan, permissions.GroupActivitiesRerun})
	if _, err := repositories.NewGroupReceiptSettingsRepository(nil).CreateGroupReceiptSettings(group); err != nil {
		t.Fatal(err)
	}
	command := commands.RegisterRecognitionTaskCommand{ClientRequestId: uuid.NewString(), FileName: "receipt.jpg", FileSize: int64(len(recognitionFixture(t))), GroupId: group, PaidByUserId: user, Status: models.OPEN}
	return NewRecognitionTaskService(), command
}

func recognitionAccept(t *testing.T, service RecognitionTaskService, command commands.RegisterRecognitionTaskCommand) models.RecognitionTask {
	t.Helper()
	task, _, err := service.Register(command.PaidByUserId, command, command.Fingerprint())
	if err != nil {
		t.Fatal(err)
	}
	task, token, err := service.ClaimUpload(task.OwnerUserId, task.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = service.AcceptUpload(task, token, recognitionFixture(t), command.FileSize); err != nil {
		t.Fatal(err)
	}
	task, err = service.Repository.Get(task.ID)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = utils.RemoveAllInDataDir(filepath.Dir(task.SourcePath)) })
	return task
}

func recognitionPipeline(t *testing.T, url string) (RecognitionTaskService, commands.RegisterRecognitionTaskCommand) {
	t.Helper()
	ClearRolePermissionCacheForTests()
	ClearGroupRoleGrantCacheForTests()
	user, group, _ := seedReceiptImagePipeline(t, url)
	role, err := repositories.NewRoleRepository(nil).CreateGroupRole("Recognition role", "", []string{permissions.GroupReceiptsRead, permissions.GroupReceiptsQuickScan, permissions.GroupActivitiesRerun, permissions.GroupCommentsCreate}, nil, nil, nil, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if err = repositories.GetDB().Model(&models.GroupMember{}).Where("group_id = ? AND user_id = ?", group.ID, user.ID).Update("group_role_id", role.ID).Error; err != nil {
		t.Fatal(err)
	}
	if _, err = repositories.NewGroupReceiptSettingsRepository(nil).CreateGroupReceiptSettings(group.ID); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		var files []models.FileData
		repositories.GetDB().Where("receipt_id != ?", 1).Find(&files)
		for _, file := range files {
			if path, e := repositories.NewFileRepository(nil).BuildFilePath(utils.UintToString(file.ReceiptId), utils.UintToString(file.ID), file.Name); e == nil {
				_ = utils.RemoveDataPath(path)
			}
		}
	})
	return NewRecognitionTaskService(), commands.RegisterRecognitionTaskCommand{ClientRequestId: uuid.NewString(), FileName: "receipt.jpg", FileSize: int64(len(recognitionFixture(t))), GroupId: group.ID, PaidByUserId: user.ID, Status: models.OPEN}
}

func recognitionAiResponse() string {
	content, _ := json.Marshal(`{"name":"Tracked receipt","amount":"12.34","date":"2026-10-07T00:00:00Z"}`)
	return `{"message":{"role":"assistant","content":` + string(content) + `},"done":true}`
}

func TestRecognitionRegistrationIdempotency(t *testing.T) {
	t.Cleanup(repositories.TruncateTestDb)
	service, command := recognitionSeed(t)
	first, created, err := service.Register(command.PaidByUserId, command, command.Fingerprint())
	if err != nil || !created {
		t.Fatalf("first registration: %v %v", created, err)
	}
	replayed, created, err := service.Register(command.PaidByUserId, command, command.Fingerprint())
	if err != nil || created || first.ID != replayed.ID {
		t.Fatalf("replay created a task: %+v %v", replayed, err)
	}
	command.FileSize++
	if _, _, err = service.Register(command.PaidByUserId, command, command.Fingerprint()); !errors.Is(err, ErrRecognitionConflict) {
		t.Fatalf("mismatched request must conflict: %v", err)
	}
	var count int64
	service.Repository.GetDB().Model(&models.RecognitionTask{}).Count(&count)
	if count != 1 {
		t.Fatalf("got %d rows", count)
	}
}

func TestRecognitionUploadFencesInterruptedWriter(t *testing.T) {
	t.Cleanup(repositories.TruncateTestDb)
	service, command := recognitionSeed(t)
	task, _, err := service.Register(command.PaidByUserId, command, command.Fingerprint())
	if err != nil {
		t.Fatal(err)
	}
	old, oldToken, err := service.ClaimUpload(task.OwnerUserId, task.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err = service.ClaimUpload(task.OwnerUserId, task.ID, nil); !errors.Is(err, ErrRecognitionConflict) {
		t.Fatalf("concurrent writer allowed: %v", err)
	}
	if err = service.UploadProgress(task.ID, oldToken, 42); err != nil {
		t.Fatal(err)
	}
	progress, _ := service.Repository.Get(task.ID)
	if progress.UploadedBytes != 42 || progress.Version <= old.Version {
		t.Fatal("real byte progress/version missing")
	}
	service.InterruptUpload(task.ID, oldToken, "UPLOAD_INTERRUPTED", "Reselect file")
	fresh, freshToken, err := service.ClaimUpload(task.OwnerUserId, task.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = service.AcceptUpload(old, oldToken, recognitionFixture(t), command.FileSize); !errors.Is(err, ErrRecognitionConflict) {
		t.Fatalf("stale upload accepted: %v", err)
	}
	if err = service.UploadProgress(task.ID, oldToken, 100); !errors.Is(err, ErrRecognitionConflict) {
		t.Fatalf("stale progress accepted: %v", err)
	}
	if err = service.AcceptUpload(fresh, freshToken, recognitionFixture(t), command.FileSize); err != nil {
		t.Fatal(err)
	}
	accepted, _ := service.Repository.Get(task.ID)
	t.Cleanup(func() { _ = utils.RemoveAllInDataDir(filepath.Dir(accepted.SourcePath)) })
	if accepted.Status != models.RecognitionDispatchPending || !recognitionSourceExists(accepted) {
		t.Fatal("accepted file was not durable")
	}
	service.InterruptUpload(task.ID, oldToken, "OLD", "Old writer")
	after, _ := service.Repository.Get(task.ID)
	if after.Status != models.RecognitionDispatchPending {
		t.Fatal("late old writer changed accepted state")
	}
}

func TestRecognitionVisibilityAndRevocation(t *testing.T) {
	t.Cleanup(repositories.TruncateTestDb)
	service, command := recognitionSeed(t)
	task, _, err := service.Register(command.PaidByUserId, command, command.Fingerprint())
	if err != nil {
		t.Fatal(err)
	}
	otherUser := models.User{Username: "other-reader", Password: "p"}
	if err = repositories.GetDB().Create(&otherUser).Error; err != nil {
		t.Fatal(err)
	}
	other := otherUser.ID
	if _, err = service.GetVisible(other, task.ID); !errors.Is(err, ErrRecognitionNotFound) {
		t.Fatalf("other owner can read: %v", err)
	}
	admin, _ := seedUserWithAppRole(t, "global-reader", []string{permissions.AppSystemTasksRead})
	query, err := service.ScopeQuery(admin, "all")
	if err != nil {
		t.Fatal(err)
	}
	var count int64
	query.Count(&count)
	if count != 1 {
		t.Fatal("global reader incorrectly group-filtered")
	}
	visible, err := service.GetVisible(admin, task.ID)
	if err != nil || visible.CanUpload || visible.CanRetry {
		t.Fatalf("global read granted mutation: %+v %v", visible, err)
	}
	repositories.GetDB().Where("user_id = ? AND group_id = ?", task.OwnerUserId, task.GroupId).Delete(&models.GroupMember{})
	if _, err = service.GetVisible(task.OwnerUserId, task.ID); !errors.Is(err, ErrRecognitionNotFound) {
		t.Fatalf("revoked membership still visible: %v", err)
	}
	query, err = service.ScopeQuery(task.OwnerUserId, "own")
	if err != nil {
		t.Fatal(err)
	}
	query.Count(&count)
	if count != 0 {
		t.Fatal("counts leaked revoked group")
	}
	if _, _, err = service.ClaimUpload(task.OwnerUserId, task.ID, nil); !errors.Is(err, ErrRecognitionForbidden) {
		t.Fatalf("revoked upload permitted: %v", err)
	}
}

func TestRecognitionAttemptAndRetryFencing(t *testing.T) {
	t.Cleanup(repositories.TruncateTestDb)
	service, command := recognitionSeed(t)
	task := recognitionAccept(t, service, command)
	old, oldToken, err := service.ClaimAttempt(task.ID, task.Generation, 0)
	if err != nil {
		t.Fatal(err)
	}
	_, currentToken, err := service.ClaimAttempt(task.ID, task.Generation, 1)
	if err != nil {
		t.Fatal(err)
	}
	if err = service.Stage(context.Background(), old, oldToken, models.RecognitionAI, false); !errors.Is(err, ErrRecognitionConflict) {
		t.Fatalf("stale stage permitted: %v", err)
	}
	service.RecordFailure(task.ID, task.Generation, 0, time.Now(), false, oldToken)
	current, _ := service.Repository.Get(task.ID)
	if current.Status != models.RecognitionRunning || current.AttemptToken != currentToken {
		t.Fatal("stale failure changed current attempt")
	}
	service.RecordFailure(task.ID, task.Generation, 3, time.Time{}, true, currentToken)
	failed, _ := service.GetVisible(task.OwnerUserId, task.ID)
	if failed.Status != models.RecognitionFailed || failed.Attempt != 4 || failed.NextRetryAt != nil || !failed.CanRetry {
		t.Fatalf("exhaustion not recorded: %+v", failed)
	}
	retried, err := service.Retry(task.OwnerUserId, task.ID, failed.Version)
	if err != nil {
		t.Fatal(err)
	}
	if retried.Generation != task.Generation+1 || retried.Attempt != 0 {
		t.Fatal("retry did not create fresh generation")
	}
	if _, err = service.Retry(task.OwnerUserId, task.ID, failed.Version); !errors.Is(err, ErrRecognitionConflict) {
		t.Fatalf("double retry permitted: %v", err)
	}
	if err = service.Process(context.Background(), task.ID, task.Generation, 0); err != nil {
		t.Fatal(err)
	}
	current, _ = service.Repository.Get(task.ID)
	if current.Generation != retried.Generation || current.Status != retried.Status {
		t.Fatal("old queue generation touched current task")
	}
}

func TestRecognitionProcessExactlyOnce(t *testing.T) {
	t.Cleanup(repositories.TruncateTestDb)
	server, _ := newMockOllamaServerForService(t, http.StatusOK, recognitionAiResponse())
	service, command := recognitionPipeline(t, server.URL)
	command.Comment = "Team lunch"
	if err := repositories.GetDB().Model(&models.GroupReceiptSettings{}).Where("group_id = ?", command.GroupId).Update("quick_scan_comment_enabled", true).Error; err != nil {
		t.Fatal(err)
	}
	task := recognitionAccept(t, service, command)
	if err := service.Process(context.Background(), task.ID, task.Generation, 0); err != nil {
		t.Fatal(err)
	}
	if err := service.Process(context.Background(), task.ID, task.Generation, 0); err != nil {
		t.Fatal(err)
	}
	saved, _ := service.Repository.Get(task.ID)
	var count int64
	repositories.GetDB().Model(&models.Receipt{}).Where("name = ?", "Tracked receipt").Count(&count)
	if count != 1 || saved.Status != models.RecognitionSucceeded || saved.Stage != models.RecognitionDone || saved.ReceiptId == nil || saved.StartedAt == nil {
		t.Fatalf("not one committed receipt: count=%d %+v", count, saved)
	}
	if recognitionSourceExists(saved) {
		t.Fatal("successful source not removed")
	}
	var comments []models.Comment
	repositories.GetDB().Where("receipt_id = ?", *saved.ReceiptId).Find(&comments)
	if len(comments) != 1 || comments[0].Comment != command.Comment || comments[0].UserId == nil || *comments[0].UserId != task.OwnerUserId {
		t.Fatalf("comment not committed exactly once with receipt: %+v", comments)
	}
	service.RecordFailure(task.ID, task.Generation, 3, time.Time{}, true, "")
	after, _ := service.Repository.Get(task.ID)
	if after.Status != models.RecognitionSucceeded {
		t.Fatal("late framework failure overwrote committed success")
	}
}

func TestRecognitionSaveRollbackRetainsSource(t *testing.T) {
	t.Cleanup(repositories.TruncateTestDb)
	server, _ := newMockOllamaServerForService(t, http.StatusOK, recognitionAiResponse())
	service, command := recognitionPipeline(t, server.URL)
	command.Comment = "Rollback comment"
	if err := repositories.GetDB().Model(&models.GroupReceiptSettings{}).Where("group_id = ?", command.GroupId).Update("quick_scan_comment_enabled", true).Error; err != nil {
		t.Fatal(err)
	}
	task := recognitionAccept(t, service, command)
	db := repositories.GetDB()
	callback := "recognition_test_commit_failure"
	db.Callback().Update().Before("gorm:update").Register(callback, func(tx *gorm.DB) {
		if tx.Statement.Table == "recognition_tasks" {
			if values, ok := tx.Statement.Dest.(map[string]interface{}); ok && values["status"] == models.RecognitionSucceeded {
				tx.AddError(errors.New("injected final save failure"))
			}
		}
	})
	t.Cleanup(func() { db.Callback().Update().Remove(callback) })
	err := service.Process(context.Background(), task.ID, task.Generation, 0)
	if err == nil {
		t.Fatal("injected failure did not rollback")
	}
	var count int64
	db.Model(&models.Receipt{}).Where("name = ?", "Tracked receipt").Count(&count)
	if count != 0 {
		t.Fatal("receipt escaped transaction rollback")
	}
	db.Model(&models.Comment{}).Where("comment = ?", command.Comment).Count(&count)
	if count != 0 {
		t.Fatal("comment escaped receipt transaction rollback")
	}
	after, _ := service.Repository.Get(task.ID)
	if after.ReceiptId != nil || !recognitionSourceExists(after) {
		t.Fatal("failed save lost recoverable source")
	}
	var attempt RecognitionAttemptError
	if !errors.As(err, &attempt) {
		t.Fatal("attempt failure was not fenced")
	}
	service.RecordFailure(task.ID, task.Generation, 3, time.Time{}, true, attempt.Token)
	if _, err = service.Retry(task.OwnerUserId, task.ID, after.Version); !errors.Is(err, ErrRecognitionConflict) {
		t.Fatalf("stale task version accepted: %v", err)
	}
}

func TestRecognitionFallbackAndCancellation(t *testing.T) {
	t.Run("fallback", func(t *testing.T) {
		t.Cleanup(repositories.TruncateTestDb)
		primary, _ := newMockOllamaServerForService(t, http.StatusOK, `{"message":{"content":"not receipt JSON"}}`)
		fallback, _ := newMockOllamaServerForService(t, http.StatusOK, recognitionAiResponse())
		service, command := recognitionPipeline(t, primary.URL)
		var settings models.ReceiptProcessingSettings
		repositories.GetDB().First(&settings)
		settings.ID = 0
		settings.Name = "fallback"
		settings.Url = fallback.URL
		repositories.GetDB().Create(&settings)
		repositories.GetDB().Model(&models.SystemSettings{}).Where("id = ?", 1).Update("fallback_receipt_processing_settings_id", settings.ID)
		task := recognitionAccept(t, service, command)
		if err := service.Process(context.Background(), task.ID, task.Generation, 0); err != nil {
			t.Fatalf("fallback processing: %v (%v)", err, errors.Unwrap(err))
		}
		after, _ := service.Repository.Get(task.ID)
		if !after.FallbackActive || after.Status != models.RecognitionSucceeded {
			t.Fatal("fallback flag/result missing")
		}
	})
	t.Run("cancel", func(t *testing.T) {
		t.Cleanup(repositories.TruncateTestDb)
		reached := make(chan struct{})
		release := make(chan struct{})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(reached)
			select {
			case <-r.Context().Done():
			case <-release:
			}
		}))
		defer server.Close()
		defer close(release)
		service, command := recognitionPipeline(t, server.URL)
		task := recognitionAccept(t, service, command)
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		finished := make(chan error, 1)
		go func() { finished <- service.Process(ctx, task.ID, task.Generation, 0) }()
		select {
		case <-reached:
		case <-time.After(5 * time.Second):
			t.Fatal("provider not called")
		}
		cancel()
		select {
		case err := <-finished:
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("worker context not propagated: %v", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("cancelled AI call did not return")
		}
		var count int64
		repositories.GetDB().Model(&models.Receipt{}).Where("name = ?", "Tracked receipt").Count(&count)
		if count != 0 {
			t.Fatal("cancelled worker created a receipt")
		}
	})
}

func TestRecognitionDispatchAndRestartRecovery(t *testing.T) {
	if os.Getenv("RECOGNITION_TEST_REDIS_HOST") == "" {
		t.Skip("requires isolated test Redis")
	}
	t.Cleanup(repositories.TruncateTestDb)
	t.Setenv("REDIS_HOST", os.Getenv("RECOGNITION_TEST_REDIS_HOST"))
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_USER", "")
	if err := repositories.ConnectToRedis(); err != nil {
		t.Fatal(err)
	}
	defer repositories.ShutdownAsynqClient()
	options, err := config.GetAsynqRedisClientConnectionOptions()
	if err != nil {
		t.Fatal(err)
	}
	inspector := asynq.NewInspector(options)
	defer inspector.Close()
	service, command := recognitionSeed(t)
	task := recognitionAccept(t, service, command)
	if err = service.Dispatch(task.ID); err != nil {
		t.Fatal(err)
	}
	if err = service.Dispatch(task.ID); err != nil {
		t.Fatal(err)
	}
	info, err := inspector.GetTaskInfo(string(models.QuickScanQueue), recognitionQueueId(task.ID, task.Generation))
	if err != nil || info.MaxRetry != 3 {
		t.Fatalf("queue contract: %+v %v", info, err)
	}
	defer inspector.DeleteTask(string(models.QuickScanQueue), info.ID)
	if err = inspector.DeleteTask(string(models.QuickScanQueue), info.ID); err != nil {
		t.Fatal(err)
	}
	if err = service.Reconcile(inspector); err != nil {
		t.Fatal(err)
	}
	if _, err = inspector.GetTaskInfo(string(models.QuickScanQueue), info.ID); err != nil {
		t.Fatal("missing Redis task did not recover")
	}
	_ = repositories.ShutdownAsynqClient()
	command.ClientRequestId = uuid.NewString()
	outage := recognitionAccept(t, service, command)
	if err = service.Dispatch(outage.ID); err == nil {
		t.Fatal("closed Redis client accepted enqueue")
	}
	pending, _ := service.Repository.Get(outage.ID)
	if pending.Status != models.RecognitionDispatchPending || pending.ErrorCode != "QUEUE_UNAVAILABLE" {
		t.Fatal("Redis outage lost durable acceptance")
	}
	if err = repositories.ConnectToRedis(); err != nil {
		t.Fatal(err)
	}
	if err = service.Reconcile(inspector); err != nil {
		t.Fatal(err)
	}
	recovered, _ := service.Repository.Get(outage.ID)
	if recovered.Status != models.RecognitionQueued {
		t.Fatal("outage task did not enqueue after reconnect")
	}
	defer inspector.DeleteTask(string(models.QuickScanQueue), recognitionQueueId(outage.ID, outage.Generation))
	command.ClientRequestId = uuid.NewString()
	upload, _, err := service.Register(command.PaidByUserId, command, command.Fingerprint())
	if err != nil {
		t.Fatal(err)
	}
	upload, _, err = service.ClaimUpload(upload.OwnerUserId, upload.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	oldStageStart := time.Now().Add(-7 * time.Minute)
	service.Repository.GetDB().Model(&models.RecognitionTask{}).Where("id = ?", upload.ID).UpdateColumn("stage_started_at", oldStageStart)
	if err = service.UploadProgress(upload.ID, upload.UploadToken, 1); err != nil {
		t.Fatal(err)
	}
	if err = service.Reconcile(inspector); err != nil {
		t.Fatal(err)
	}
	progressing, _ := service.Repository.Get(upload.ID)
	if progressing.Status != models.RecognitionUploading || progressing.StageStartedAt == nil || !progressing.StageStartedAt.Equal(oldStageStart) {
		t.Fatal("active upload was interrupted or its stage start time was replaced")
	}
	service.Repository.GetDB().Model(&models.RecognitionTask{}).Where("id = ?", upload.ID).UpdateColumn("updated_at", time.Now().Add(-7*time.Minute))
	if err = service.Reconcile(inspector); err != nil {
		t.Fatal(err)
	}
	abandoned, _ := service.Repository.Get(upload.ID)
	if abandoned.Status != models.RecognitionUploadInterrupted || abandoned.UploadToken != "" {
		t.Fatal("abandoned upload writer did not get fenced")
	}
}
