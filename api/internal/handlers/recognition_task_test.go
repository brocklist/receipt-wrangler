package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/utils"
	"receipt-wrangler/api/internal/wranglerasynq"
)

func recognitionHandlerSeed(t *testing.T) (uint, commands.RegisterRecognitionTaskCommand, []byte) {
	t.Helper()
	t.Cleanup(repositories.TruncateTestDb)
	repositories.CreateTestUser()
	repositories.CreateTestGroup()
	grantGroupPerms(t, 1, 1, permissions.GroupReceiptsRead, permissions.GroupReceiptsQuickScan, permissions.GroupActivitiesRerun)
	if _, err := repositories.NewGroupReceiptSettingsRepository(nil).CreateGroupReceiptSettings(1); err != nil {
		t.Fatal(err)
	}
	_, source, _, _ := runtime.Caller(0)
	root := filepath.Clean(filepath.Join(filepath.Dir(source), "..", ".."))
	t.Setenv("BASE_PATH", root)
	file, err := os.ReadFile(filepath.Join(root, "testing", "test.jpg"))
	if err != nil {
		t.Fatal(err)
	}
	return 1, commands.RegisterRecognitionTaskCommand{ClientRequestId: uuid.NewString(), FileName: "receipt.jpg", FileSize: int64(len(file)), GroupId: 1, PaidByUserId: 1, Status: models.OPEN}, file
}

func recognitionHandlerRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Post("/api/recognitionTask", CreateRecognitionTask)
	router.Get("/api/recognitionTask", GetRecognitionTasks)
	router.Get("/api/recognitionTask/{id}", GetRecognitionTask)
	router.Put("/api/recognitionTask/{id}/file", UploadRecognitionTaskFile)
	router.Post("/api/recognitionTask/{id}/retry", RetryRecognitionTask)
	router.Post("/api/systemTask/rerunActivity/{id}", RerunActivity)
	router.Get("/api/systemTask/{id}/sourceFile", GetSystemTaskSourceFile)
	router.Get("/api/systemTask/{id}/sourceFile/download", DownloadSystemTaskSourceFile)
	router.Post("/api/receipt/quickScan", QuickScan)
	return router
}

func recognitionHandlerSetPermissions(t *testing.T, user, group uint, keys ...string) {
	t.Helper()
	role, err := repositories.NewRoleRepository(nil).CreateGroupRole("Recognition changed "+uuid.NewString(), "", keys, nil, nil, nil, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if err = repositories.GetDB().Model(&models.GroupMember{}).Where("user_id = ? AND group_id = ?", user, group).Update("group_role_id", role.ID).Error; err != nil {
		t.Fatal(err)
	}
	services.ClearRolePermissionCacheForTests()
	services.ClearGroupRoleGrantCacheForTests()
}

func recognitionHandlerJSON(t *testing.T, user uint, method, path string, value interface{}) *httptest.ResponseRecorder {
	t.Helper()
	body, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRequest(method, path, bytes.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(user)))
	w := httptest.NewRecorder()
	recognitionHandlerRouter().ServeHTTP(w, r)
	return w
}

func recognitionDecode(t *testing.T, w *httptest.ResponseRecorder) models.RecognitionTask {
	t.Helper()
	var task models.RecognitionTask
	if err := json.Unmarshal(w.Body.Bytes(), &task); err != nil {
		t.Fatalf("decode %s: %v", w.Body.String(), err)
	}
	return task
}

func recognitionHandlerMultipart(t *testing.T, user uint, path, field string, file []byte, broken bool) *httptest.ResponseRecorder {
	t.Helper()
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(field, "receipt.jpg")
	if err != nil {
		t.Fatal(err)
	}
	part.Write(file)
	if field == "files" {
		writer.WriteField("groupIds", "1")
		writer.WriteField("paidByUserIds", "1")
		writer.WriteField("statuses", string(models.OPEN))
	}
	writer.Close()
	length := int64(body.Len())
	data := body.Bytes()
	if broken {
		data = data[:len(data)/2]
	}
	method := http.MethodPut
	if field == "files" {
		method = http.MethodPost
	}
	r := httptest.NewRequest(method, path, bytes.NewReader(data))
	r.ContentLength = length
	r.Header.Set("Content-Type", writer.FormDataContentType())
	r = r.WithContext(context.WithValue(r.Context(), jwtmiddleware.ContextKey{}, claimsForUser(user)))
	w := httptest.NewRecorder()
	recognitionHandlerRouter().ServeHTTP(w, r)
	return w
}

func TestRecognitionHandlerRegistrationAndConflict(t *testing.T) {
	user, command, _ := recognitionHandlerSeed(t)
	w := recognitionHandlerJSON(t, user, "POST", "/api/recognitionTask", command)
	if w.Code != http.StatusCreated {
		t.Fatalf("register %d: %s", w.Code, w.Body.String())
	}
	first := recognitionDecode(t, w)
	w = recognitionHandlerJSON(t, user, "POST", "/api/recognitionTask", command)
	if w.Code != http.StatusOK || recognitionDecode(t, w).ID != first.ID {
		t.Fatal("same registration did not recover original task")
	}
	command.FileSize++
	w = recognitionHandlerJSON(t, user, "POST", "/api/recognitionTask", command)
	if w.Code != http.StatusConflict {
		t.Fatalf("mismatch %d", w.Code)
	}
	command.ClientRequestId = uuid.NewString()
	command.FileName = "../escape.jpg"
	w = recognitionHandlerJSON(t, user, "POST", "/api/recognitionTask", command)
	if w.Code != http.StatusBadRequest {
		t.Fatal("unsafe filename registered")
	}
	recognitionHandlerSetPermissions(t, user, 1, permissions.GroupReceiptsRead)
	command.FileName = "receipt.jpg"
	w = recognitionHandlerJSON(t, user, "POST", "/api/recognitionTask", command)
	if w.Code != http.StatusForbidden {
		t.Fatal("revoked quick scan permission allowed registration")
	}
}

func TestRecognitionHandlerQuickScanCommentContract(t *testing.T) {
	for _, tc := range []struct {
		name                          string
		enabled, required, permission bool
		comment, saved                string
		status                        int
	}{
		{"required whitespace rejected", true, true, true, "  \n ", "", http.StatusBadRequest},
		{"optional comment preserved", true, false, true, "  Lunch, team\n会议  ", "Lunch, team\n会议", http.StatusCreated},
		{"disabled comment dropped", false, false, true, "incidental", "", http.StatusCreated},
		{"missing comment permission waived", true, true, false, "incidental", "", http.StatusCreated},
		{"oversized shown comment rejected", true, false, true, strings.Repeat("中", models.MaxCommentLength+1), "", http.StatusBadRequest},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Cleanup(repositories.TruncateTestDb)
			user, group := seedQuickScanCommenter(t, tc.permission, models.GroupReceiptSettings{QuickScanCommentEnabled: tc.enabled, QuickScanCommentRequired: tc.required})
			command := commands.RegisterRecognitionTaskCommand{ClientRequestId: uuid.NewString(), FileName: "receipt.jpg", FileSize: 100, GroupId: group, Comment: tc.comment}
			w := recognitionHandlerJSON(t, user, http.MethodPost, "/api/recognitionTask", command)
			if w.Code != tc.status {
				t.Fatalf("registration status %d: %s", w.Code, w.Body.String())
			}
			var tasks []models.RecognitionTask
			repositories.GetDB().Find(&tasks)
			if tc.status == http.StatusCreated {
				if len(tasks) != 1 || tasks[0].Comment != tc.saved {
					t.Fatalf("resolved comment not persisted: %+v", tasks)
				}
				if strings.Contains(w.Body.String(), "Lunch") || strings.Contains(w.Body.String(), "incidental") {
					t.Fatal("private submission comment leaked in task response")
				}
				command.Comment = "changed comment"
				w = recognitionHandlerJSON(t, user, http.MethodPost, "/api/recognitionTask", command)
				if w.Code != http.StatusConflict {
					t.Fatalf("changed comment reused identity: %d", w.Code)
				}
			} else if len(tasks) != 0 {
				t.Fatal("invalid comment registered a task")
			}
		})
	}
}

func TestRecognitionHandlerRecoveryDistinguishesRevokedAccessFromMissingRequest(t *testing.T) {
	user, command, _ := recognitionHandlerSeed(t)
	registered := recognitionHandlerJSON(t, user, "POST", "/api/recognitionTask", command)
	if registered.Code != http.StatusCreated {
		t.Fatalf("register %d: %s", registered.Code, registered.Body.String())
	}
	recognitionHandlerSetPermissions(t, user, 1, permissions.GroupReceiptsQuickScan)
	path := "/api/recognitionTask?scope=own&bucket=all&page=1&pageSize=1&clientRequestId=" + command.ClientRequestId
	w := recognitionHandlerJSON(t, user, http.MethodGet, path, nil)
	if w.Code != http.StatusForbidden {
		t.Fatalf("revoked group access was reported as missing: %d %s", w.Code, w.Body.String())
	}
	missing := "/api/recognitionTask?scope=own&bucket=all&page=1&pageSize=1&clientRequestId=" + uuid.NewString()
	w = recognitionHandlerJSON(t, user, http.MethodGet, missing, nil)
	var response commands.RecognitionTaskList
	if w.Code != http.StatusOK || json.Unmarshal(w.Body.Bytes(), &response) != nil || response.TotalCount != 0 {
		t.Fatalf("missing request was not distinguishable: %d %s", w.Code, w.Body.String())
	}
}

func TestRecognitionHandlerReplayDoesNotReapplyDefaults(t *testing.T) {
	user, command, _ := recognitionHandlerSeed(t)
	db := repositories.GetDB()
	db.Model(&models.GroupReceiptSettings{}).Where("group_id = ?", 1).Updates(map[string]interface{}{"quick_scan_paid_by_required": false, "quick_scan_status_required": false, "quick_scan_default_paid_by_type": models.QUICK_SCAN_PAID_BY_UPLOADER, "quick_scan_default_status": models.OPEN})
	command.PaidByUserId = 0
	command.Status = ""
	w := recognitionHandlerJSON(t, user, "POST", "/api/recognitionTask", command)
	if w.Code != http.StatusCreated {
		t.Fatalf("register defaults: %s", w.Body.String())
	}
	first := recognitionDecode(t, w)
	db.Model(&models.GroupReceiptSettings{}).Where("group_id = ?", 1).Update("quick_scan_default_status", models.DRAFT)
	w = recognitionHandlerJSON(t, user, "POST", "/api/recognitionTask", command)
	if w.Code != http.StatusOK || recognitionDecode(t, w).ID != first.ID {
		t.Fatal("configuration change broke ambiguous-response recovery")
	}
	stored, err := repositories.NewRecognitionTaskRepository(nil).Get(first.ID)
	if err != nil || stored.ReceiptStatus != models.OPEN || stored.PaidByUserId != user {
		t.Fatal("replay changed immutable resolved fields")
	}
}

func TestRecognitionHandlerUploadInterruptionAndRecovery(t *testing.T) {
	user, command, file := recognitionHandlerSeed(t)
	w := recognitionHandlerJSON(t, user, "POST", "/api/recognitionTask", command)
	task := recognitionDecode(t, w)
	path := fmt.Sprintf("/api/recognitionTask/%d/file", task.ID)
	w = recognitionHandlerMultipart(t, user, path, "file", file, true)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("broken upload: %d %s", w.Code, w.Body.String())
	}
	w = recognitionHandlerJSON(t, user, "GET", fmt.Sprintf("/api/recognitionTask/%d", task.ID), nil)
	interrupted := recognitionDecode(t, w)
	if interrupted.Status != models.RecognitionUploadInterrupted || !interrupted.CanUpload || interrupted.UploadedBytes == 0 {
		t.Fatalf("interrupted state: %+v", interrupted)
	}
	w = recognitionHandlerMultipart(t, user, path, "file", file, false)
	if w.Code != http.StatusAccepted {
		t.Fatalf("complete reupload: %d %s", w.Code, w.Body.String())
	}
	accepted := recognitionDecode(t, w)
	if accepted.Status != models.RecognitionDispatchPending && accepted.Status != models.RecognitionQueued {
		t.Fatal("server did not confirm durable receipt")
	}
	if accepted.UploadTotalBytes == nil || accepted.UploadedBytes != *accepted.UploadTotalBytes || accepted.Version <= interrupted.Version {
		t.Fatal("server byte meter or version was missing")
	}
	stored, _ := repositories.NewRecognitionTaskRepository(nil).Get(task.ID)
	t.Cleanup(func() { utils.RemoveAllInDataDir(filepath.Dir(stored.SourcePath)) })
	w = recognitionHandlerMultipart(t, user, path, "file", file, false)
	if w.Code != http.StatusConflict {
		t.Fatal("accepted upload was replayed")
	}
	var count int64
	repositories.GetDB().Model(&models.RecognitionTask{}).Count(&count)
	if count != 1 {
		t.Fatal("reupload duplicated task")
	}
	body := w.Body.String()
	if bytes.Contains([]byte(body), []byte("source_path")) {
		t.Fatal("private source path exposed")
	}
}

func TestRecognitionHandlerScopedCountsAndFilters(t *testing.T) {
	user, command, _ := recognitionHandlerSeed(t)
	service := services.NewRecognitionTaskService()
	for i := 0; i < 3; i++ {
		command.ClientRequestId = uuid.NewString()
		task, _, err := service.Register(user, command, command.Fingerprint())
		if err != nil {
			t.Fatal(err)
		}
		if i == 2 {
			service.Repository.Update(repositories.GetDB().Where("id = ?", task.ID), map[string]interface{}{"status": models.RecognitionFailed})
		}
	}
	command.ClientRequestId = uuid.NewString()
	command.GroupId = 2
	_, _, err := service.Register(99, command, command.Fingerprint())
	if err != nil {
		t.Fatal(err)
	}
	w := recognitionHandlerJSON(t, user, "GET", "/api/recognitionTask?scope=own&bucket=history&pageSize=1&status=FAILED&groupId=1", nil)
	var page commands.RecognitionTaskList
	if err = json.Unmarshal(w.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 || page.TotalCount != 1 || len(page.Data) != 1 || page.ActiveCount != 2 || page.FailedCount != 1 {
		t.Fatalf("counts affected by page/filter: %+v %s", page, w.Body.String())
	}
	w = recognitionHandlerJSON(t, user, "GET", "/api/recognitionTask?scope=all", nil)
	if w.Code != 403 {
		t.Fatal("normal user read all tasks")
	}
	grantAppPerms(t, 2, permissions.AppSystemTasksRead)
	w = recognitionHandlerJSON(t, 2, "GET", "/api/recognitionTask?scope=all&pageSize=1", nil)
	if err = json.Unmarshal(w.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	if w.Code != 200 || page.TotalCount != 4 || page.ActiveCount != 3 {
		t.Fatalf("global reader incorrectly scoped: %s", w.Body.String())
	}
	for _, task := range page.Data {
		if task.CanUpload || task.CanRetry {
			t.Fatal("global read granted mutations")
		}
	}
	recognitionHandlerSetPermissions(t, user, 1, permissions.GroupReceiptsQuickScan)
	w = recognitionHandlerJSON(t, user, "GET", "/api/recognitionTask?scope=own", nil)
	if err = json.Unmarshal(w.Body.Bytes(), &page); err != nil {
		t.Fatal(err)
	}
	if page.TotalCount != 0 || page.ActiveCount != 0 {
		t.Fatal("revoked read permission leaked counts")
	}
}

func TestRecognitionInterruptedUploadAppearsInActionableActiveView(t *testing.T) {
	user, command, _ := recognitionHandlerSeed(t)
	service := services.NewRecognitionTaskService()
	task, _, err := service.Register(user, command, command.Fingerprint())
	if err != nil {
		t.Fatal(err)
	}
	service.Repository.Update(repositories.GetDB().Where("id = ?", task.ID), map[string]interface{}{"status": models.RecognitionUploadInterrupted})

	w := recognitionHandlerJSON(t, user, "GET", "/api/recognitionTask?scope=own&bucket=active", nil)
	var active commands.RecognitionTaskList
	if err = json.Unmarshal(w.Body.Bytes(), &active); err != nil {
		t.Fatal(err)
	}
	if w.Code != http.StatusOK || active.TotalCount != 1 || len(active.Data) != 1 || !active.Data[0].CanUpload {
		t.Fatalf("interrupted upload missing from active view: %s", w.Body.String())
	}

	w = recognitionHandlerJSON(t, user, "GET", "/api/recognitionTask?scope=own&bucket=history", nil)
	var history commands.RecognitionTaskList
	if err = json.Unmarshal(w.Body.Bytes(), &history); err != nil {
		t.Fatal(err)
	}
	if history.TotalCount != 0 {
		t.Fatalf("actionable upload leaked into history view: %s", w.Body.String())
	}
}

func recognitionHandlerRedis(t *testing.T) {
	t.Helper()
	host := os.Getenv("RECOGNITION_TEST_REDIS_HOST")
	if host == "" {
		t.Skip("requires isolated test Redis")
	}
	t.Setenv("REDIS_HOST", host)
	t.Setenv("REDIS_PORT", "6379")
	t.Setenv("REDIS_PASSWORD", "")
	t.Setenv("REDIS_USER", "")
	if err := repositories.ConnectToRedis(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { repositories.ShutdownAsynqClient() })
}

func TestRecognitionLegacyQuickScanEmptyResponseAndTracking(t *testing.T) {
	user, _, file := recognitionHandlerSeed(t)
	recognitionHandlerRedis(t)
	w := recognitionHandlerMultipart(t, user, "/api/receipt/quickScan", "files", file, false)
	if w.Code != http.StatusOK || w.Body.Len() != 0 {
		t.Fatalf("legacy response changed: %d %q", w.Code, w.Body.String())
	}
	var tasks []models.RecognitionTask
	repositories.GetDB().Find(&tasks)
	if len(tasks) != 1 || tasks[0].Status != models.RecognitionQueued || tasks[0].UploadTotalBytes != nil {
		t.Fatalf("legacy acceptance not tracked without fake percentage: %+v", tasks)
	}
	task := tasks[0]
	t.Cleanup(func() { utils.RemoveAllInDataDir(filepath.Dir(task.SourcePath)) })
	inspector, err := wranglerasynq.GetAsynqInspector()
	if err != nil {
		t.Fatal(err)
	}
	defer inspector.Close()
	defer inspector.DeleteTask(string(models.QuickScanQueue), task.AsynqTaskId)
	info, err := inspector.GetTaskInfo(string(models.QuickScanQueue), task.AsynqTaskId)
	if err != nil {
		t.Fatal(err)
	}
	var payload models.RecognitionTaskPayload
	if err = json.Unmarshal(info.Payload, &payload); err != nil || payload.RecognitionTaskId != task.ID {
		t.Fatal("legacy accepted file used untracked queue payload")
	}
}

func TestRecognitionTrackedActivityRerunUsesGenerationFence(t *testing.T) {
	user, command, file := recognitionHandlerSeed(t)
	recognitionHandlerSetPermissions(t, user, command.GroupId, permissions.GroupReceiptsRead, permissions.GroupReceiptsQuickScan, permissions.GroupActivitiesRerun, permissions.GroupActivitiesRead)
	recognitionHandlerRedis(t)
	service := services.NewRecognitionTaskService()
	task, _, err := service.Register(user, command, command.Fingerprint())
	if err != nil {
		t.Fatal(err)
	}
	task, token, err := service.ClaimUpload(user, task.ID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if err = service.AcceptUpload(task, token, file, int64(len(file))); err != nil {
		t.Fatal(err)
	}
	if err = service.Dispatch(task.ID); err != nil {
		t.Fatal(err)
	}
	task, _ = service.Repository.Get(task.ID)
	t.Cleanup(func() { utils.RemoveAllInDataDir(filepath.Dir(task.SourcePath)) })
	inspector, err := wranglerasynq.GetAsynqInspector()
	if err != nil {
		t.Fatal(err)
	}
	defer inspector.Close()
	defer inspector.DeleteTask(string(models.QuickScanQueue), task.AsynqTaskId)
	if err = inspector.ArchiveTask(string(models.QuickScanQueue), task.AsynqTaskId); err != nil {
		t.Fatal(err)
	}
	service.Repository.Update(repositories.GetDB().Where("id = ?", task.ID), map[string]interface{}{"status": models.RecognitionFailed, "attempt": 4, "completed_at": time.Now()})
	audit := models.SystemTask{Type: models.QUICK_SCAN, Status: models.SYSTEM_TASK_FAILED, AssociatedEntityType: models.RECEIPT_PROCESSING_SETTINGS, AsynqTaskId: task.AsynqTaskId, GroupId: &task.GroupId}
	if err = repositories.GetDB().Create(&audit).Error; err != nil {
		t.Fatal(err)
	}
	sourcePath := fmt.Sprintf("/api/systemTask/%d/sourceFile", audit.ID)
	source := recognitionHandlerJSON(t, user, http.MethodGet, sourcePath, nil)
	if source.Code != http.StatusOK || !strings.Contains(source.Body.String(), "receipt.jpg") || !strings.Contains(source.Body.String(), "data:image/") {
		t.Fatalf("durable activity preview failed: %d %s", source.Code, source.Body.String())
	}
	source = recognitionHandlerJSON(t, user, http.MethodGet, sourcePath+"/download", nil)
	if source.Code != http.StatusOK || !bytes.Equal(source.Body.Bytes(), file) {
		t.Fatalf("durable activity download changed bytes: %d", source.Code)
	}
	path := fmt.Sprintf("/api/systemTask/rerunActivity/%d", audit.ID)
	w := recognitionHandlerJSON(t, user, "POST", path, nil)
	if w.Code != 200 {
		t.Fatalf("tracked activity rerun: %d %s", w.Code, w.Body.String())
	}
	after, _ := service.Repository.Get(task.ID)
	defer inspector.DeleteTask(string(models.QuickScanQueue), after.AsynqTaskId)
	if after.Generation != task.Generation+1 || after.Status != models.RecognitionQueued {
		t.Fatal("activity did not use recognition retry")
	}
	source = recognitionHandlerJSON(t, user, http.MethodGet, sourcePath+"/download", nil)
	if source.Code != http.StatusNotFound {
		t.Fatalf("old activity generation exposed current upload: %d", source.Code)
	}
	w = recognitionHandlerJSON(t, user, "POST", path, nil)
	if w.Code != 409 {
		t.Fatalf("double activity rerun accepted: %d", w.Code)
	}
	current, _ := service.Repository.Get(task.ID)
	if current.Generation != after.Generation {
		t.Fatal("double click created an extra generation")
	}
}
