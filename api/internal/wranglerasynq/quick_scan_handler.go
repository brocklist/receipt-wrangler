package wranglerasynq

import (
	"context"
	"encoding/json"
	"github.com/hibiken/asynq"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
)

type QuickScanTaskPayload struct {
	models.RecognitionTaskPayload
	Token            *structs.Claims
	PaidByUserId     uint
	GroupId          uint
	Status           models.ReceiptStatus
	CategoryIds      []uint
	TagIds           []uint
	Comment          string
	TempPath         string
	OriginalFileName string
}

func HandleQuickScanTask(context context.Context, task *asynq.Task) error {
	taskId, err := GetTaskIdFromContext(context)
	if err != nil {
		return HandleError(err)
	}

	var payload QuickScanTaskPayload

	err = json.Unmarshal(task.Payload(), &payload)
	if err != nil {
		return HandleError(err)
	}
	if payload.RecognitionTaskId > 0 {
		retried, _ := asynq.GetRetryCount(context)
		return services.NewRecognitionTaskService().Process(context, payload.RecognitionTaskId, payload.Generation, retried+payload.AttemptOffset)
	}

	receiptService := services.NewReceiptService(nil)
	_, err = receiptService.QuickScan(services.QuickScanParams{
		Token:            payload.Token,
		PaidByUserId:     payload.PaidByUserId,
		GroupId:          payload.GroupId,
		Status:           payload.Status,
		CategoryIds:      payload.CategoryIds,
		TagIds:           payload.TagIds,
		Comment:          payload.Comment,
		TempPath:         payload.TempPath,
		OriginalFileName: payload.OriginalFileName,
		AsynqTaskId:      taskId,
	})
	if err != nil {
		return HandleError(err)
	}

	return nil
}
