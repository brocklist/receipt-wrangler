package commands

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/google/uuid"
	"path/filepath"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/structs"
	"receipt-wrangler/api/internal/utils"
	"strings"
)

const RecognitionMaxFileSize = 50 << 20

type RegisterRecognitionTaskCommand struct {
	ClientRequestId string               `json:"clientRequestId"`
	FileName        string               `json:"fileName"`
	FileSize        int64                `json:"fileSize"`
	GroupId         uint                 `json:"groupId"`
	PaidByUserId    uint                 `json:"paidByUserId"`
	Status          models.ReceiptStatus `json:"status"`
	CategoryIds     []uint               `json:"categoryIds"`
	TagIds          []uint               `json:"tagIds"`
}

func (c RegisterRecognitionTaskCommand) Validate() structs.ValidatorError {
	v := structs.ValidatorError{Errors: map[string]string{}}
	if _, err := uuid.Parse(c.ClientRequestId); err != nil {
		v.Errors["clientRequestId"] = "A UUID is required"
	}
	if strings.TrimSpace(c.FileName) == "" || len(c.FileName) > 255 || !utils.IsSafePathComponent(c.FileName) || filepath.Base(c.FileName) != c.FileName {
		v.Errors["fileName"] = "A file name without path separators is required"
	}
	if c.FileSize <= 0 || c.FileSize > RecognitionMaxFileSize {
		v.Errors["fileSize"] = "File must be between 1 byte and 50 MiB"
	}
	if c.GroupId == 0 {
		v.Errors["groupId"] = "Group is required"
	}
	if c.Status != "" {
		if _, err := c.Status.Value(); err != nil {
			v.Errors["status"] = "Invalid receipt status"
		}
	}
	return v
}

func (c RegisterRecognitionTaskCommand) Fingerprint() string {
	c.ClientRequestId = ""
	data, _ := json.Marshal(c)
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

type RetryRecognitionTaskCommand struct {
	Version uint `json:"version"`
}

type RecognitionTaskList struct {
	Data                []models.RecognitionTask `json:"data"`
	TotalCount          int64                    `json:"totalCount"`
	ActiveCount         int64                    `json:"activeCount"`
	AwaitingUploadCount int64                    `json:"awaitingUploadCount"`
	RunningCount        int64                    `json:"runningCount"`
	FailedCount         int64                    `json:"failedCount"`
}
