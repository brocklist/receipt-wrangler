package services

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/hibiken/asynq"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"net/http"
	"path/filepath"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/utils"
	"time"
)

var ErrRecognitionConflict = errors.New("recognition task changed; refresh its status")
var ErrRecognitionForbidden = errors.New("not permitted to perform this task action")
var ErrRecognitionNotFound = errors.New("recognition task not found")

type RecognitionTaskService struct {
	Repository repositories.RecognitionTaskRepository
}

func NewRecognitionTaskService() RecognitionTaskService {
	return RecognitionTaskService{repositories.NewRecognitionTaskRepository(nil)}
}

func (s RecognitionTaskService) Register(owner uint, c commands.RegisterRecognitionTaskCommand, fingerprint string) (models.RecognitionTask, bool, error) {
	task := models.RecognitionTask{OwnerUserId: owner, ClientRequestId: c.ClientRequestId, FileName: c.FileName, FileSize: c.FileSize, GroupId: c.GroupId, RequestHash: fingerprint, Status: models.RecognitionAwaitingUpload, Stage: models.RecognitionUpload, Version: 1, Generation: 1, MaxAttempts: 4, PaidByUserId: c.PaidByUserId, ReceiptStatus: c.Status, CategoryIds: c.CategoryIds, TagIds: c.TagIds}
	result := s.Repository.GetDB().Clauses(clause.OnConflict{Columns: []clause.Column{{Name: "owner_user_id"}, {Name: "client_request_id"}}, DoNothing: true}).Create(&task)
	if result.Error != nil {
		return task, false, result.Error
	}
	created := result.RowsAffected == 1
	if !created {
		var err error
		task, err = s.Repository.GetByRequest(owner, c.ClientRequestId)
		if err != nil {
			return task, false, err
		}
		if task.RequestHash != fingerprint {
			return task, false, ErrRecognitionConflict
		}
	}
	return task, created, nil
}

func (s RecognitionTaskService) Visible(user uint, task models.RecognitionTask, all bool) (bool, error) {
	p := NewPermissionService(nil)
	if all {
		return p.HasAppPermissions(user, permissions.AppSystemTasksRead)
	}
	if task.OwnerUserId != user {
		return false, nil
	}
	return p.HasGroupPermissions(user, task.GroupId, permissions.GroupReceiptsRead)
}

func (s RecognitionTaskService) Flags(user uint, task *models.RecognitionTask) error {
	p := NewPermissionService(nil)
	quick, err := p.HasGroupPermissions(user, task.GroupId, permissions.GroupReceiptsQuickScan)
	if err != nil {
		return err
	}
	task.CanUpload = quick && user == task.OwnerUserId && (task.Status == models.RecognitionAwaitingUpload || task.Status == models.RecognitionUploadInterrupted)
	rerun, err := p.HasGroupPermissions(user, task.GroupId, permissions.GroupActivitiesRerun)
	if err != nil {
		return err
	}
	task.CanRetry = quick && rerun && task.Status == models.RecognitionFailed && recognitionSourceExists(*task)
	if task.CanRetry {
		if err = s.ValidateSubmission(*task); err != nil {
			task.CanRetry = false
			if !errors.Is(err, ErrRecognitionForbidden) {
				return err
			}
		}
		allowed, grantErr := p.ValidateCategoryTagSelection(user, task.GroupId, task.CategoryIds, task.TagIds)
		if grantErr != nil {
			return grantErr
		}
		task.CanRetry = task.CanRetry && allowed
	}
	return nil
}

func recognitionSourceExists(task models.RecognitionTask) bool {
	return task.SourcePath != "" && utils.AssertWithinDataDir(task.SourcePath) == nil && utils.FileExists(task.SourcePath)
}

func (s RecognitionTaskService) GetVisible(user, id uint) (models.RecognitionTask, error) {
	task, err := s.Repository.Get(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return task, ErrRecognitionNotFound
		}
		return task, err
	}
	all, err := NewPermissionService(nil).HasAppPermissions(user, permissions.AppSystemTasksRead)
	if err != nil {
		return task, err
	}
	visible, err := s.Visible(user, task, all)
	if err != nil {
		return task, err
	}
	if !visible {
		return task, ErrRecognitionNotFound
	}
	err = s.Flags(user, &task)
	return task, err
}

// ScopeQuery applies ownership and current group access before counts/paging.
func (s RecognitionTaskService) ScopeQuery(user uint, scope string) (*gorm.DB, error) {
	db := s.Repository.GetDB()
	query := db.Model(&models.RecognitionTask{})
	if scope == "all" {
		ok, err := NewPermissionService(nil).HasAppPermissions(user, permissions.AppSystemTasksRead)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, ErrRecognitionForbidden
		}
		return query, nil
	}
	query = query.Where("owner_user_id = ?", user)
	var groups []uint
	if err := db.Model(&models.RecognitionTask{}).Where("owner_user_id = ?", user).Distinct("group_id").Pluck("group_id", &groups).Error; err != nil {
		return nil, err
	}
	allowed := []uint{0}
	p := NewPermissionService(nil)
	for _, group := range groups {
		ok, err := p.HasGroupPermissions(user, group, permissions.GroupReceiptsRead)
		if err != nil {
			return nil, err
		}
		if ok {
			allowed = append(allowed, group)
		}
	}
	return query.Where("group_id IN ?", allowed), nil
}

func (s RecognitionTaskService) ClaimUpload(user, id uint, total *int64) (models.RecognitionTask, string, error) {
	task, err := s.Repository.Get(id)
	if err != nil {
		return task, "", ErrRecognitionNotFound
	}
	if task.OwnerUserId != user {
		return task, "", ErrRecognitionForbidden
	}
	ok, err := NewPermissionService(nil).HasGroupPermissions(user, task.GroupId, permissions.GroupReceiptsQuickScan)
	if err != nil {
		return task, "", err
	}
	if !ok {
		return task, "", ErrRecognitionForbidden
	}
	token := uuid.NewString()
	now := time.Now()
	result := s.Repository.Update(s.Repository.GetDB().Where("id = ? AND status IN ?", id, []models.RecognitionTaskStatus{models.RecognitionAwaitingUpload, models.RecognitionUploadInterrupted}), map[string]interface{}{"status": models.RecognitionUploading, "upload_token": token, "stage": models.RecognitionUpload, "stage_started_at": now, "completed_at": nil, "uploaded_bytes": 0, "upload_total_bytes": total, "error_code": "", "error_message": ""})
	if result.Error != nil {
		return task, "", result.Error
	}
	if result.RowsAffected != 1 {
		return task, "", ErrRecognitionConflict
	}
	task, err = s.Repository.Get(id)
	return task, token, err
}

func (s RecognitionTaskService) UploadProgress(id uint, token string, received int64) error {
	r := s.Repository.Update(s.Repository.GetDB().Where("id = ? AND status = ? AND upload_token = ?", id, models.RecognitionUploading, token), map[string]interface{}{"uploaded_bytes": received, "updated_at": time.Now()})
	if r.Error != nil {
		return r.Error
	}
	if r.RowsAffected != 1 {
		return ErrRecognitionConflict
	}
	return nil
}

func (s RecognitionTaskService) InterruptUpload(id uint, token, code, message string) {
	s.Repository.Update(s.Repository.GetDB().Where("id = ? AND status = ? AND upload_token = ?", id, models.RecognitionUploading, token), map[string]interface{}{"status": models.RecognitionUploadInterrupted, "upload_token": "", "error_code": code, "error_message": message, "completed_at": time.Now()})
}

func (s RecognitionTaskService) AcceptUpload(task models.RecognitionTask, token string, bytes []byte, received int64) error {
	if int64(len(bytes)) != task.FileSize || len(bytes) > commands.RecognitionMaxFileSize {
		return fmt.Errorf("uploaded file size does not match registration")
	}
	if _, err := repositories.NewFileRepository(nil).ValidateFileType(bytes); err != nil {
		return err
	}
	sum := sha256.Sum256(bytes)
	hash := hex.EncodeToString(sum[:])
	root, err := utils.GetDataDir()
	if err != nil {
		return err
	}
	dir := filepath.Join(root, "recognition-uploads", fmt.Sprint(task.ID))
	for _, path := range []string{root, filepath.Dir(dir), dir} {
		if err = utils.EnsureDataDirectory(path); err != nil {
			return err
		}
	}
	// Every upload uses its fencing token in its source name. A stale uploader
	// can never replace a newer upload's immutable file.
	path := filepath.Join(dir, token+".source")
	if err = utils.WriteDataFile(path, bytes); err != nil {
		return err
	}
	now := time.Now()
	result := s.Repository.Update(s.Repository.GetDB().Where("id = ? AND status = ? AND upload_token = ?", task.ID, models.RecognitionUploading, token), map[string]interface{}{"status": models.RecognitionDispatchPending, "source_path": path, "source_hash": hash, "upload_token": "", "uploaded_bytes": received, "queued_at": now, "error_code": "", "error_message": "", "completed_at": nil})
	if result.Error != nil || result.RowsAffected != 1 {
		utils.RemoveDataPath(path)
		if result.Error != nil {
			return result.Error
		}
		return ErrRecognitionConflict
	}
	return nil
}

func recognitionQueueId(id, generation uint) string {
	return fmt.Sprintf("recognition-%d-%d", id, generation)
}

func (s RecognitionTaskService) Dispatch(id uint) error {
	task, err := s.Repository.Get(id)
	if err != nil {
		return err
	}
	if task.Status != models.RecognitionDispatchPending {
		return nil
	}
	if task.Attempt >= task.MaxAttempts {
		s.Repository.Update(s.Repository.GetDB().Where("id = ? AND version = ?", id, task.Version), map[string]interface{}{"status": models.RecognitionFailed, "attempt_token": "", "completed_at": time.Now(), "next_retry_at": nil, "error_code": "PROCESSING_FAILED", "error_message": "Recognition exhausted its attempts. You can retry this task."})
		return nil
	}
	client := repositories.GetAsynqClient()
	payload, _ := json.Marshal(models.RecognitionTaskPayload{RecognitionTaskId: id, Generation: task.Generation, AttemptOffset: task.Attempt})
	queueId := recognitionQueueId(id, task.Generation)
	if client == nil {
		err = fmt.Errorf("task queue unavailable")
	} else {
		_, err = client.Enqueue(asynq.NewTask("receipt:quick_scan", payload), asynq.Queue(string(models.QuickScanQueue)), asynq.MaxRetry(task.MaxAttempts-task.Attempt-1), asynq.TaskID(queueId))
	}
	if err != nil && !errors.Is(err, asynq.ErrTaskIDConflict) {
		s.Repository.Update(s.Repository.GetDB().Where("id = ? AND status = ?", id, models.RecognitionDispatchPending), map[string]interface{}{"error_code": "QUEUE_UNAVAILABLE", "error_message": "Waiting for the task queue to become available."})
		return err
	}
	return s.Repository.Update(s.Repository.GetDB().Where("id = ? AND status = ? AND generation = ?", id, models.RecognitionDispatchPending, task.Generation), map[string]interface{}{"status": models.RecognitionQueued, "asynq_task_id": queueId, "error_code": "", "error_message": ""}).Error
}

func (s RecognitionTaskService) Retry(user, id, version uint) (models.RecognitionTask, error) {
	task, err := s.GetVisible(user, id)
	if err != nil {
		return task, err
	}
	if task.Version != version || task.Status != models.RecognitionFailed {
		return task, ErrRecognitionConflict
	}
	if !task.CanRetry {
		return task, ErrRecognitionForbidden
	}
	if err = s.ValidateSubmission(task); err != nil {
		return task, err
	}
	now := time.Now()
	result := s.Repository.Update(s.Repository.GetDB().Where("id = ? AND version = ? AND status = ?", id, version, models.RecognitionFailed), map[string]interface{}{"status": models.RecognitionDispatchPending, "generation": gorm.Expr("generation + 1"), "attempt_token": "", "asynq_task_id": "", "attempt": 0, "next_retry_at": nil, "completed_at": nil, "started_at": nil, "fallback_active": false, "queued_at": now, "error_code": "", "error_message": ""})
	if result.Error != nil {
		return task, result.Error
	}
	if result.RowsAffected != 1 {
		return task, ErrRecognitionConflict
	}
	// Dispatch gaps are retained in the ledger and repaired by reconciliation.
	_ = s.Dispatch(id)
	return s.GetVisible(user, id)
}

func (s RecognitionTaskService) ClaimAttempt(id, generation uint, retryCount int) (models.RecognitionTask, string, error) {
	task, err := s.Repository.Get(id)
	if err != nil {
		return task, "", err
	}
	if task.Generation != generation || task.Status == models.RecognitionSucceeded || task.Status == models.RecognitionFailed {
		return task, "", nil
	}
	token := uuid.NewString()
	now := time.Now()
	result := s.Repository.Update(s.Repository.GetDB().Where("id = ? AND generation = ? AND status IN ?", id, generation, []models.RecognitionTaskStatus{models.RecognitionDispatchPending, models.RecognitionQueued, models.RecognitionRunning, models.RecognitionRetryWait}), map[string]interface{}{"status": models.RecognitionRunning, "stage": models.RecognitionPreprocessing, "stage_started_at": now, "started_at": now, "attempt": retryCount + 1, "attempt_token": token, "asynq_task_id": recognitionQueueId(id, generation), "next_retry_at": nil, "error_code": "", "error_message": ""})
	if result.Error != nil {
		return task, "", result.Error
	}
	if result.RowsAffected != 1 {
		return task, "", nil
	}
	task, err = s.Repository.Get(id)
	return task, token, err
}

func (s RecognitionTaskService) Stage(ctx context.Context, task models.RecognitionTask, token string, stage models.RecognitionTaskStage, fallback bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	result := s.Repository.Update(s.Repository.GetDB().Where("id = ? AND generation = ? AND attempt_token = ? AND status = ?", task.ID, task.Generation, token, models.RecognitionRunning), map[string]interface{}{"stage": stage, "stage_started_at": time.Now(), "fallback_active": fallback})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 1 {
		return ErrRecognitionConflict
	}
	return nil
}

type RecognitionAttemptError struct {
	Token string
	Err   error
}

func (e RecognitionAttemptError) Error() string { return "recognition processing failed" }
func (e RecognitionAttemptError) Unwrap() error { return e.Err }

// ValidateSubmission rechecks the saved user's grants and today's required fields.
func (s RecognitionTaskService) ValidateSubmission(task models.RecognitionTask) error {
	p := NewPermissionService(nil)
	ok, err := p.HasGroupPermissions(task.OwnerUserId, task.GroupId, permissions.GroupReceiptsQuickScan)
	if err != nil {
		return err
	}
	if !ok {
		return ErrRecognitionForbidden
	}
	ok, err = p.ValidateCategoryTagSelection(task.OwnerUserId, task.GroupId, task.CategoryIds, task.TagIds)
	if err != nil {
		return err
	}
	if !ok {
		return ErrRecognitionForbidden
	}
	settings, err := repositories.NewGroupReceiptSettingsRepository(nil).GetGroupReceiptSettingsByGroupId(task.GroupId)
	if err != nil {
		return err
	}
	if (settings.QuickScanPaidByEnabled && settings.QuickScanPaidByRequired && task.PaidByUserId == 0) || (settings.QuickScanStatusEnabled && settings.QuickScanStatusRequired && task.ReceiptStatus == "") || (settings.QuickScanCategoriesEnabled && settings.QuickScanCategoriesRequired && len(task.CategoryIds) == 0) || (settings.QuickScanTagsEnabled && settings.QuickScanTagsRequired && len(task.TagIds) == 0) {
		return ErrRecognitionForbidden
	}
	return nil
}

func (s RecognitionTaskService) Process(ctx context.Context, id, generation uint, retryCount int) (processErr error) {
	task, token, err := s.ClaimAttempt(id, generation, retryCount)
	if err != nil || token == "" {
		return err
	}
	defer func() {
		if processErr != nil {
			processErr = RecognitionAttemptError{Token: token, Err: processErr}
		}
	}()
	if err = s.ValidateSubmission(task); err != nil {
		return err
	}
	bytes, err := utils.ReadDataFile(task.SourcePath)
	if err != nil {
		return fmt.Errorf("source file unavailable: %w", err)
	}
	sum := sha256.Sum256(bytes)
	if hex.EncodeToString(sum[:]) != task.SourceHash {
		return fmt.Errorf("source integrity check failed")
	}
	started := time.Now()
	fallbackActive := false
	command, metadata, err := MagicFillFromImageWithProgress(ctx, commands.MagicFillCommand{ImageData: bytes, Filename: task.FileName}, utils.UintToString(task.GroupId), task.OwnerUserId, func(stage models.RecognitionTaskStage, fallback bool) error {
		fallbackActive = fallback
		return s.Stage(ctx, task, token, stage, fallback)
	})
	if ctx.Err() != nil {
		return ctx.Err()
	}
	tasks, taskErr := NewSystemTaskService(nil).CreateSystemTasksFromMetadata(metadata, started, time.Now(), models.QUICK_SCAN, &task.OwnerUserId, &task.GroupId, recognitionQueueId(id, generation), nil)
	if taskErr != nil {
		return taskErr
	}
	if err != nil {
		return err
	}
	if command.PaidByUserID == 0 {
		command.PaidByUserID = task.PaidByUserId
	}
	if command.Status == "" {
		command.Status = task.ReceiptStatus
	}
	command.GroupId = task.GroupId
	receiptService := NewReceiptService(nil)
	command.Categories, err = receiptService.mergeQuickScanCategories(command.Categories, task.CategoryIds)
	if err != nil {
		return err
	}
	command.Tags, err = receiptService.mergeQuickScanTags(command.Tags, task.TagIds)
	if err != nil {
		return err
	}
	if err = s.ValidateSubmission(task); err != nil {
		return err
	}
	categoryIds := make([]uint, 0, len(command.Categories))
	for _, category := range command.Categories {
		if category.Id != nil {
			categoryIds = append(categoryIds, *category.Id)
		}
	}
	tagIds := make([]uint, 0, len(command.Tags))
	for _, tag := range command.Tags {
		if tag.Id != nil {
			tagIds = append(tagIds, *tag.Id)
		}
	}
	allowed, grantErr := NewPermissionService(nil).ValidateCategoryTagSelection(task.OwnerUserId, task.GroupId, categoryIds, tagIds)
	if grantErr != nil {
		return grantErr
	}
	if !allowed {
		return ErrRecognitionForbidden
	}
	if v := command.Validate(task.OwnerUserId, true); len(v.Errors) > 0 {
		return fmt.Errorf("recognized receipt failed validation")
	}
	if err = s.Stage(ctx, task, token, models.RecognitionSaving, fallbackActive); err != nil {
		return err
	}
	var savedImagePath string
	err = s.Repository.GetDB().Transaction(func(tx *gorm.DB) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		// Conditional UPDATE locks the ledger row on all supported SQL engines.
		// Receipt, image metadata, notifications and success commit together.
		r := repositories.NewRecognitionTaskRepository(tx)
		claim := r.Update(tx.Where("id = ? AND generation = ? AND attempt_token = ? AND status = ?", id, generation, token, models.RecognitionRunning), map[string]interface{}{"stage": models.RecognitionSaving})
		if claim.Error != nil {
			return claim.Error
		}
		if claim.RowsAffected != 1 {
			return ErrRecognitionConflict
		}
		receipt, saveErr := repositories.NewReceiptRepository(tx).CreateReceipt(command, task.OwnerUserId, false)
		if saveErr != nil {
			return saveErr
		}
		file, saveErr := repositories.NewReceiptImageRepository(tx).CreateReceiptImage(models.FileData{Name: task.FileName, Size: uint(len(bytes)), ReceiptId: receipt.ID}, bytes)
		if saveErr != nil {
			return saveErr
		}
		savedImagePath, saveErr = repositories.NewFileRepository(tx).BuildFilePath(utils.UintToString(receipt.ID), utils.UintToString(file.ID), file.Name)
		if saveErr != nil {
			return saveErr
		}
		if ctx.Err() != nil {
			return ctx.Err()
		}
		audit := NewSystemTaskService(tx)
		if _, saveErr = audit.CreateReceiptUploadedSystemTask(nil, receipt, tasks, time.Now()); saveErr != nil {
			return saveErr
		}
		if saveErr = audit.AssociateProcessingSystemTasksToReceipt(tasks, receipt.ID); saveErr != nil {
			return saveErr
		}
		result := r.Update(tx.Where("id = ? AND generation = ? AND attempt_token = ?", id, generation, token), map[string]interface{}{"status": models.RecognitionSucceeded, "stage": models.RecognitionDone, "stage_started_at": time.Now(), "completed_at": time.Now(), "receipt_id": receipt.ID, "attempt_token": "", "error_code": "", "error_message": ""})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected != 1 {
			return ErrRecognitionConflict
		}
		return nil
	})
	if err != nil && savedImagePath != "" {
		_ = utils.RemoveDataPath(savedImagePath)
	}
	if err == nil {
		_ = utils.RemoveDataPath(task.SourcePath)
	}
	return err
}

func (s RecognitionTaskService) RecordFailure(id, generation uint, retryCount int, next time.Time, exhausted bool, token string) {
	status := models.RecognitionRetryWait
	var completed interface{}
	if exhausted {
		status = models.RecognitionFailed
		completed = time.Now()
	}
	var retryAt interface{}
	if !exhausted && !next.IsZero() {
		retryAt = next
	}
	query := s.Repository.GetDB().Where("id = ? AND generation = ? AND attempt <= ? AND status IN ?", id, generation, retryCount+1, []models.RecognitionTaskStatus{models.RecognitionRunning, models.RecognitionQueued, models.RecognitionRetryWait})
	if token != "" {
		query = query.Where("attempt_token = ?", token)
	}
	s.Repository.Update(query, map[string]interface{}{"status": status, "attempt_token": "", "next_retry_at": retryAt, "attempt": retryCount + 1, "completed_at": completed, "error_code": "PROCESSING_FAILED", "error_message": "Receipt recognition could not finish. Check the processing settings or retry the task."})
}

// Reconcile repairs the SQL/Redis dispatch boundary and process crash states.
// It never promotes a missing Redis task to success or replaces SQL success.
func (s RecognitionTaskService) Reconcile(inspector *asynq.Inspector) error {
	var tasks []models.RecognitionTask
	if err := s.Repository.GetDB().Where("status IN ?", models.RecognitionActiveStatuses()).Find(&tasks).Error; err != nil {
		return err
	}
	for _, task := range tasks {
		if task.Status == models.RecognitionUploading {
			// StageStartedAt measures the upload stage. UpdatedAt is the heartbeat:
			// a long transfer stays valid while the server keeps receiving bytes.
			if time.Since(task.UpdatedAt) > 6*time.Minute {
				s.InterruptUpload(task.ID, task.UploadToken, "UPLOAD_INTERRUPTED", "Upload was interrupted. Reselect the original file to upload it again.")
			}
			continue
		}
		if task.Status == models.RecognitionAwaitingUpload {
			continue
		}
		if task.Status == models.RecognitionDispatchPending {
			_ = s.Dispatch(task.ID)
			continue
		}
		info, err := inspector.GetTaskInfo(string(models.QuickScanQueue), recognitionQueueId(task.ID, task.Generation))
		if errors.Is(err, asynq.ErrTaskNotFound) || errors.Is(err, asynq.ErrQueueNotFound) {
			r := s.Repository.Update(s.Repository.GetDB().Where("id = ? AND version = ?", task.ID, task.Version), map[string]interface{}{"status": models.RecognitionDispatchPending, "attempt_token": ""})
			if r.Error == nil && r.RowsAffected == 1 {
				_ = s.Dispatch(task.ID)
			}
			continue
		}
		if err != nil {
			continue
		}
		var payload models.RecognitionTaskPayload
		_ = json.Unmarshal(info.Payload, &payload)
		if info.State == asynq.TaskStateArchived {
			s.RecordFailure(task.ID, task.Generation, info.Retried+payload.AttemptOffset, time.Time{}, true, task.AttemptToken)
		}
		if info.State == asynq.TaskStateRetry && (task.Status != models.RecognitionRetryWait || task.NextRetryAt == nil || !task.NextRetryAt.Equal(info.NextProcessAt)) {
			s.Repository.Update(s.Repository.GetDB().Where("id = ? AND version = ?", task.ID, task.Version), map[string]interface{}{"status": models.RecognitionRetryWait, "attempt_token": "", "next_retry_at": info.NextProcessAt, "attempt": info.Retried + payload.AttemptOffset, "error_code": "PROCESSING_FAILED", "error_message": "Receipt recognition will retry automatically."})
		}
		if info.State == asynq.TaskStatePending && task.Status == models.RecognitionRunning {
			s.Repository.Update(s.Repository.GetDB().Where("id = ? AND version = ?", task.ID, task.Version), map[string]interface{}{"status": models.RecognitionQueued, "attempt_token": ""})
		}
	}
	return nil
}

func RecognitionErrorStatus(err error) int {
	switch {
	case errors.Is(err, ErrRecognitionForbidden):
		return http.StatusForbidden
	case errors.Is(err, ErrRecognitionNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrRecognitionConflict):
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
