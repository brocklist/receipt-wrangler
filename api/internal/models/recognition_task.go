package models

import "time"

type RecognitionTaskStatus string
type RecognitionTaskStage string

const (
	RecognitionAwaitingUpload    RecognitionTaskStatus = "AWAITING_UPLOAD"
	RecognitionUploading         RecognitionTaskStatus = "UPLOADING"
	RecognitionUploadInterrupted RecognitionTaskStatus = "UPLOAD_INTERRUPTED"
	RecognitionDispatchPending   RecognitionTaskStatus = "DISPATCH_PENDING"
	RecognitionQueued            RecognitionTaskStatus = "QUEUED"
	RecognitionRunning           RecognitionTaskStatus = "RUNNING"
	RecognitionRetryWait         RecognitionTaskStatus = "RETRY_WAIT"
	RecognitionSucceeded         RecognitionTaskStatus = "SUCCEEDED"
	RecognitionFailed            RecognitionTaskStatus = "FAILED"
	RecognitionUpload            RecognitionTaskStage  = "UPLOAD"
	RecognitionPreprocessing     RecognitionTaskStage  = "PREPROCESSING"
	RecognitionOCR               RecognitionTaskStage  = "OCR"
	RecognitionAI                RecognitionTaskStage  = "AI"
	RecognitionParsing           RecognitionTaskStage  = "PARSING"
	RecognitionSaving            RecognitionTaskStage  = "SAVING"
	RecognitionDone              RecognitionTaskStage  = "DONE"
)

// RecognitionTask is the durable Quick Scan lifecycle. SystemTask remains its
// terminal processing audit; source paths and worker fencing keys are private.
type RecognitionTask struct {
	BaseModel
	OwnerUserId      uint                  `gorm:"uniqueIndex:recognition_request;index" json:"ownerUserId"`
	ClientRequestId  string                `gorm:"size:64;uniqueIndex:recognition_request" json:"clientRequestId"`
	GroupId          uint                  `gorm:"index" json:"groupId"`
	FileName         string                `json:"fileName"`
	FileSize         int64                 `json:"fileSize"`
	Version          uint                  `gorm:"not null;default:1" json:"version"`
	Status           RecognitionTaskStatus `gorm:"index" json:"status"`
	Stage            RecognitionTaskStage  `json:"stage"`
	QueuedAt         *time.Time            `json:"queuedAt"`
	StartedAt        *time.Time            `json:"startedAt"`
	StageStartedAt   *time.Time            `json:"stageStartedAt"`
	CompletedAt      *time.Time            `json:"completedAt"`
	UploadedBytes    int64                 `json:"uploadedBytes"`
	UploadTotalBytes *int64                `json:"uploadTotalBytes"`
	Attempt          int                   `json:"attempt"`
	MaxAttempts      int                   `gorm:"default:4" json:"maxAttempts"`
	NextRetryAt      *time.Time            `json:"nextRetryAt"`
	FallbackActive   bool                  `json:"fallbackActive"`
	ReceiptId        *uint                 `json:"receiptId"`
	ErrorCode        string                `json:"errorCode"`
	ErrorMessage     string                `json:"errorMessage"`
	CanUpload        bool                  `gorm:"-" json:"canUpload"`
	CanRetry         bool                  `gorm:"-" json:"canRetry"`
	RequestHash      string                `json:"-"`
	SourceHash       string                `json:"-"`
	SourcePath       string                `json:"-"`
	UploadToken      string                `json:"-"`
	AttemptToken     string                `json:"-"`
	Generation       uint                  `gorm:"default:1" json:"-"`
	AsynqTaskId      string                `gorm:"index" json:"-"`
	PaidByUserId     uint                  `json:"-"`
	ReceiptStatus    ReceiptStatus         `json:"-"`
	CategoryIds      []uint                `gorm:"serializer:json" json:"-"`
	TagIds           []uint                `gorm:"serializer:json" json:"-"`
}

func RecognitionActiveStatuses() []RecognitionTaskStatus {
	return []RecognitionTaskStatus{RecognitionAwaitingUpload, RecognitionUploading, RecognitionDispatchPending, RecognitionQueued, RecognitionRunning, RecognitionRetryWait}
}

// RecognitionTaskListActiveStatuses also includes uploads that need user action.
// Keep this separate from RecognitionActiveStatuses, which drives queue recovery.
func RecognitionTaskListActiveStatuses() []RecognitionTaskStatus {
	return append(RecognitionActiveStatuses(), RecognitionUploadInterrupted)
}

// These fields coexist with the legacy Asynq payload during upgrades.
type RecognitionTaskPayload struct {
	RecognitionTaskId uint
	Generation        uint
	AttemptOffset     int
}
