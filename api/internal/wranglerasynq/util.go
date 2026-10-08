package wranglerasynq

import (
	"encoding/json"
	"errors"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/repositories"
	"receipt-wrangler/api/internal/services"
	"receipt-wrangler/api/internal/structs"
)

// SetActivityFlagsForUser preserves source-file recovery and additionally checks
// the current caller's permissions for durable recognition retries.
func SetActivityFlagsForUser(activities *[]structs.Activity, userId uint) error {
	if len(*activities) == 0 {
		return nil
	}
	inspector, err := GetAsynqInspector()
	if err != nil {
		return nil
	}
	defer inspector.Close()
	lookup := memoizeTaskInfoLookup(inspector.GetTaskInfo)
	applyActivityFlags(activities, lookup)
	for i := range *activities {
		activity := &(*activities)[i]
		if activity.Type != models.QUICK_SCAN || !activity.CanBeRestarted {
			continue
		}
		info, err := lookup(string(models.QuickScanQueue), activity.AsynqTaskId)
		if err != nil || info == nil {
			activity.CanBeRestarted = false
			continue
		}
		var payload models.RecognitionTaskPayload
		if json.Unmarshal(info.Payload, &payload) != nil || payload.RecognitionTaskId == 0 {
			continue
		}
		record, err := repositories.NewRecognitionTaskRepository(nil).Get(payload.RecognitionTaskId)
		if err != nil || record.Generation != payload.Generation || record.Status != models.RecognitionFailed {
			activity.CanBeRestarted = false
			continue
		}
		if err = services.NewRecognitionTaskService().Flags(userId, &record); err != nil {
			return err
		}
		activity.CanBeRestarted = record.CanRetry
	}
	return nil
}

func SystemTaskToQueueName(taskType models.SystemTaskType) (string, error) {
	if string(taskType) == string(models.QUICK_SCAN) {
		return string(models.QuickScanQueue), nil
	}

	if string(taskType) == string(models.EMAIL_UPLOAD) {
		return string(models.EmailReceiptProcessingQueue), nil
	}

	return "", errors.New("unsupported task type")
}
