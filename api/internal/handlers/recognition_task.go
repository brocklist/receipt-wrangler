package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"receipt-wrangler/api/internal/commands"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/permissions"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
)

func recognitionResponse(w http.ResponseWriter, status int, value interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func recognitionError(w http.ResponseWriter, err error) {
	status := services.RecognitionErrorStatus(err)
	message := "Could not update recognition task. Refresh its status and try again."
	if status != http.StatusInternalServerError {
		message = err.Error()
	}
	utils.WriteCustomErrorResponse(w, message, status)
}

func recognitionId(w http.ResponseWriter, r *http.Request) (uint, bool) {
	id, err := utils.StringToUint(chi.URLParam(r, "id"))
	if err != nil || id == 0 {
		utils.WriteCustomErrorResponse(w, "Invalid recognition task ID", http.StatusBadRequest)
		return 0, false
	}
	return id, true
}

func CreateRecognitionTask(w http.ResponseWriter, r *http.Request) {
	var command commands.RegisterRecognitionTaskCommand
	r.Body = http.MaxBytesReader(w, r.Body, 64<<10)
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil {
		utils.WriteCustomErrorResponse(w, "Invalid task registration", http.StatusBadRequest)
		return
	}
	if v := command.Validate(); len(v.Errors) > 0 {
		structs.WriteValidatorErrorResponse(w, v, http.StatusBadRequest)
		return
	}
	user := structs.GetClaims(r).UserId
	service := services.NewRecognitionTaskService()
	ok, err := services.NewPermissionService(nil).HasGroupPermissions(user, command.GroupId, permissions.GroupReceiptsQuickScan)
	if err != nil {
		recognitionError(w, err)
		return
	}
	if !ok {
		recognitionError(w, services.ErrRecognitionForbidden)
		return
	}
	isAllGroup, err := isAllGroupDestination(command.GroupId)
	if err != nil {
		recognitionError(w, err)
		return
	}
	if isAllGroup {
		utils.WriteCustomErrorResponse(w, allGroupCreateMessage, http.StatusBadRequest)
		return
	}
	command.Comment = strings.TrimSpace(command.Comment)
	fingerprint := command.Fingerprint()
	// Replay compares the submitted values before resolving today's defaults.
	// A changed group default must not turn an ambiguous response into a new job.
	existing, err := service.Repository.GetByRequest(user, command.ClientRequestId)
	if err == nil {
		if existing.RequestHash != fingerprint {
			recognitionError(w, services.ErrRecognitionConflict)
			return
		}
		if err = service.Flags(user, &existing); err != nil {
			recognitionError(w, err)
			return
		}
		recognitionResponse(w, http.StatusOK, existing)
		return
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		recognitionError(w, err)
		return
	}
	quick := commands.QuickScanCommand{Files: []multipart.File{nil}, GroupIds: []uint{command.GroupId}, PaidByUserIds: []uint{command.PaidByUserId}, Statuses: []models.ReceiptStatus{command.Status}, CategoryIds: [][]uint{command.CategoryIds}, TagIds: [][]uint{command.TagIds}, Comments: []string{command.Comment}}
	resolved, v, err := services.NewReceiptService(nil).ResolveQuickScanFields(quick, user)
	if err != nil {
		recognitionError(w, err)
		return
	}
	if len(v.Errors) > 0 {
		structs.WriteValidatorErrorResponse(w, v, http.StatusBadRequest)
		return
	}
	allowed, message, err := enforceQuickScanGrantSelection(user, quick)
	if err != nil {
		recognitionError(w, err)
		return
	}
	if !allowed {
		utils.WriteCustomErrorResponse(w, message, http.StatusForbidden)
		return
	}
	command.PaidByUserId, command.Status = resolved[0].PaidByUserId, resolved[0].Status
	command.CategoryIds, command.TagIds = resolved[0].CategoryIds, resolved[0].TagIds
	command.Comment = resolved[0].Comment
	task, created, err := service.Register(user, command, fingerprint)
	if err != nil {
		recognitionError(w, err)
		return
	}
	if err = service.Flags(user, &task); err != nil {
		recognitionError(w, err)
		return
	}
	status := http.StatusOK
	if created {
		status = http.StatusCreated
	}
	recognitionResponse(w, status, task)
}

// Counts the bytes HTTP actually receives, including multipart framing. The
// denominator is Content-Length (the same transfer measured by browser XHR).
type recognitionUploadReader struct {
	io.ReadCloser
	count    int64
	last     time.Time
	progress func(int64) error
}

func (reader *recognitionUploadReader) Read(p []byte) (int, error) {
	n, err := reader.ReadCloser.Read(p)
	reader.count += int64(n)
	if time.Since(reader.last) >= time.Second || err == io.EOF {
		reader.last = time.Now()
		if progressErr := reader.progress(reader.count); progressErr != nil {
			return n, progressErr
		}
	}
	return n, err
}

func UploadRecognitionTaskFile(w http.ResponseWriter, r *http.Request) {
	id, ok := recognitionId(w, r)
	if !ok {
		return
	}
	service := services.NewRecognitionTaskService()
	var total *int64
	if r.ContentLength > 0 {
		length := r.ContentLength
		total = &length
	}
	task, token, err := service.ClaimUpload(structs.GetClaims(r).UserId, id, total)
	if err != nil {
		recognitionError(w, err)
		return
	}
	accepted := false
	defer func() {
		if r.MultipartForm != nil {
			_ = r.MultipartForm.RemoveAll()
		}
		if !accepted {
			service.InterruptUpload(id, token, "UPLOAD_INTERRUPTED", "Upload was interrupted or rejected. Reselect the original file and upload it again.")
		}
	}()
	reader := &recognitionUploadReader{ReadCloser: http.MaxBytesReader(w, r.Body, commands.RecognitionMaxFileSize+(1<<20)), progress: func(received int64) error { return service.UploadProgress(id, token, received) }}
	r.Body = reader
	if err = r.ParseMultipartForm(8 << 20); err != nil {
		var tooLarge *http.MaxBytesError
		status := http.StatusBadRequest
		if errors.As(err, &tooLarge) {
			status = http.StatusRequestEntityTooLarge
		}
		utils.WriteCustomErrorResponse(w, "Upload was interrupted or exceeds the file limit", status)
		return
	}
	files := r.MultipartForm.File["file"]
	if len(files) != 1 || len(r.MultipartForm.File) != 1 {
		utils.WriteCustomErrorResponse(w, "Upload exactly one file", http.StatusBadRequest)
		return
	}
	if files[0].Filename != task.FileName || files[0].Size != task.FileSize {
		utils.WriteCustomErrorResponse(w, "Reselect the registered file with its original name and size", http.StatusBadRequest)
		return
	}
	file, err := files[0].Open()
	if err != nil {
		recognitionError(w, err)
		return
	}
	defer file.Close()
	bytes, err := io.ReadAll(io.LimitReader(file, commands.RecognitionMaxFileSize+1))
	if err != nil {
		recognitionError(w, err)
		return
	}
	if _, err = repositories.NewFileRepository(nil).ValidateFileType(bytes); err != nil {
		utils.WriteCustomErrorResponse(w, "Unsupported receipt file", http.StatusBadRequest)
		return
	}
	if r.Context().Err() != nil {
		return
	}
	if err = service.AcceptUpload(task, token, bytes, reader.count); err != nil {
		recognitionError(w, err)
		return
	}
	accepted = true
	// Durable SQL acceptance succeeds even during a Redis outage. The dispatcher
	// retries this boundary without asking the client to transfer another file.
	_ = service.Dispatch(id)
	task, err = service.GetVisible(structs.GetClaims(r).UserId, id)
	if err != nil {
		recognitionError(w, err)
		return
	}
	recognitionResponse(w, http.StatusAccepted, task)
}

func GetRecognitionTask(w http.ResponseWriter, r *http.Request) {
	id, ok := recognitionId(w, r)
	if !ok {
		return
	}
	task, err := services.NewRecognitionTaskService().GetVisible(structs.GetClaims(r).UserId, id)
	if err != nil {
		recognitionError(w, err)
		return
	}
	recognitionResponse(w, http.StatusOK, task)
}

func RetryRecognitionTask(w http.ResponseWriter, r *http.Request) {
	id, ok := recognitionId(w, r)
	if !ok {
		return
	}
	var command commands.RetryRecognitionTaskCommand
	r.Body = http.MaxBytesReader(w, r.Body, 1024)
	if err := json.NewDecoder(r.Body).Decode(&command); err != nil || command.Version == 0 {
		utils.WriteCustomErrorResponse(w, "Expected task version is required", http.StatusBadRequest)
		return
	}
	task, err := services.NewRecognitionTaskService().Retry(structs.GetClaims(r).UserId, id, command.Version)
	if err != nil {
		recognitionError(w, err)
		return
	}
	recognitionResponse(w, http.StatusAccepted, task)
}

func GetRecognitionTasks(w http.ResponseWriter, r *http.Request) {
	user := structs.GetClaims(r).UserId
	service := services.NewRecognitionTaskService()
	params := r.URL.Query()
	scope, bucket := params.Get("scope"), params.Get("bucket")
	if scope == "" {
		scope = "own"
	}
	if bucket == "" {
		bucket = "all"
	}
	if (scope != "own" && scope != "all") || (bucket != "all" && bucket != "active" && bucket != "history") {
		utils.WriteCustomErrorResponse(w, "Invalid task scope or bucket", http.StatusBadRequest)
		return
	}
	page, pageSize := 1, 25
	var err error
	if params.Get("page") != "" {
		page, err = strconv.Atoi(params.Get("page"))
		if err != nil || page < 1 || page > 100000 {
			utils.WriteCustomErrorResponse(w, "Invalid page", http.StatusBadRequest)
			return
		}
	}
	if params.Get("pageSize") != "" {
		pageSize, err = strconv.Atoi(params.Get("pageSize"))
		if err != nil || pageSize < 1 || pageSize > 100 {
			utils.WriteCustomErrorResponse(w, "Invalid page size", http.StatusBadRequest)
			return
		}
	}
	query, err := service.ScopeQuery(user, scope)
	if err != nil {
		recognitionError(w, err)
		return
	}
	// An empty own-scope lookup can mean either no registration reached the
	// server or the caller lost group read access. Distinguish those cases for
	// recovery without returning any details about a now-inaccessible task.
	if request := params.Get("clientRequestId"); request != "" && scope == "own" {
		var owned models.RecognitionTask
		lookupErr := service.Repository.GetDB().Where("owner_user_id = ? AND client_request_id = ?", user, request).First(&owned).Error
		if lookupErr != nil && !errors.Is(lookupErr, gorm.ErrRecordNotFound) {
			recognitionError(w, lookupErr)
			return
		}
		if lookupErr == nil {
			var visible int64
			if err = query.Session(&gorm.Session{}).Where("id = ?", owned.ID).Count(&visible).Error; err != nil {
				recognitionError(w, err)
				return
			}
			if visible == 0 {
				recognitionError(w, services.ErrRecognitionForbidden)
				return
			}
		}
	}
	response := commands.RecognitionTaskList{Data: []models.RecognitionTask{}}
	count := func(statuses []models.RecognitionTaskStatus, target *int64) error {
		return query.Session(&gorm.Session{}).Where("status IN ?", statuses).Count(target).Error
	}
	if err = count(models.RecognitionActiveStatuses(), &response.ActiveCount); err != nil {
		recognitionError(w, err)
		return
	}
	if err = count([]models.RecognitionTaskStatus{models.RecognitionAwaitingUpload, models.RecognitionUploadInterrupted}, &response.AwaitingUploadCount); err != nil {
		recognitionError(w, err)
		return
	}
	if err = count([]models.RecognitionTaskStatus{models.RecognitionRunning}, &response.RunningCount); err != nil {
		recognitionError(w, err)
		return
	}
	if err = count([]models.RecognitionTaskStatus{models.RecognitionFailed}, &response.FailedCount); err != nil {
		recognitionError(w, err)
		return
	}
	filtered := query.Session(&gorm.Session{})
	if bucket == "active" {
		filtered = filtered.Where("status IN ?", models.RecognitionTaskListActiveStatuses())
	}
	if bucket == "history" {
		filtered = filtered.Where("status NOT IN ?", models.RecognitionTaskListActiveStatuses())
	}
	if request := params.Get("clientRequestId"); request != "" {
		filtered = filtered.Where("client_request_id = ?", request)
	}
	if group := params.Get("groupId"); group != "" {
		id, e := utils.StringToUint(group)
		if e != nil || id == 0 {
			utils.WriteCustomErrorResponse(w, "Invalid group ID", http.StatusBadRequest)
			return
		}
		filtered = filtered.Where("group_id = ?", id)
	}
	if status := params.Get("status"); status != "" {
		valid := false
		for _, s := range append(models.RecognitionActiveStatuses(), models.RecognitionSucceeded, models.RecognitionFailed, models.RecognitionUploadInterrupted) {
			if string(s) == status {
				valid = true
			}
		}
		if !valid {
			utils.WriteCustomErrorResponse(w, "Invalid task status", http.StatusBadRequest)
			return
		}
		filtered = filtered.Where("status = ?", status)
	}
	if ids := params.Get("ids"); ids != "" {
		parts := strings.Split(ids, ",")
		if len(parts) > 100 {
			utils.WriteCustomErrorResponse(w, "Too many task IDs", http.StatusBadRequest)
			return
		}
		values := make([]uint, 0, len(parts))
		for _, part := range parts {
			id, e := utils.StringToUint(strings.TrimSpace(part))
			if e != nil || id == 0 {
				utils.WriteCustomErrorResponse(w, "Invalid task ID", http.StatusBadRequest)
				return
			}
			values = append(values, id)
		}
		filtered = filtered.Where("id IN ?", values)
	}
	if err = filtered.Session(&gorm.Session{}).Count(&response.TotalCount).Error; err != nil {
		recognitionError(w, err)
		return
	}
	if err = filtered.Order("created_at DESC, id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&response.Data).Error; err != nil {
		recognitionError(w, err)
		return
	}
	for i := range response.Data {
		if err = service.Flags(user, &response.Data[i]); err != nil {
			recognitionError(w, err)
			return
		}
	}
	recognitionResponse(w, http.StatusOK, response)
}
