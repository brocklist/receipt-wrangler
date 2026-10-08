package wranglerasynq

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/hibiken/asynq"
	"receipt-wrangler/api/internal/logging"
	"receipt-wrangler/api/internal/models"
	"receipt-wrangler/api/internal/services"
)

func trackedRecognitionPayload(task *asynq.Task) (models.RecognitionTaskPayload, bool) {
	var payload models.RecognitionTaskPayload
	if task.Type() != QuickScan || json.Unmarshal(task.Payload(), &payload) != nil || payload.RecognitionTaskId == 0 {
		return payload, false
	}
	return payload, true
}

func recognitionRetryDelay(n int, err error, task *asynq.Task) time.Duration {
	delay := asynq.DefaultRetryDelayFunc(n, err, task)
	if payload, ok := trackedRecognitionPayload(task); ok {
		var attempt services.RecognitionAttemptError
		var token string
		if errors.As(err, &attempt) {
			token = attempt.Token
		}
		services.NewRecognitionTaskService().RecordFailure(payload.RecognitionTaskId, payload.Generation, n+payload.AttemptOffset, time.Now().Add(delay), false, token)
	}
	return delay
}

func recognitionTaskError(ctx context.Context, task *asynq.Task, err error) {
	payload, ok := trackedRecognitionPayload(task)
	if !ok {
		return
	}
	retried, _ := asynq.GetRetryCount(ctx)
	maxRetry, _ := asynq.GetMaxRetry(ctx)
	if retried >= maxRetry || errors.Is(err, asynq.SkipRetry) {
		var attempt services.RecognitionAttemptError
		var token string
		if errors.As(err, &attempt) {
			token = attempt.Token
		}
		services.NewRecognitionTaskService().RecordFailure(payload.RecognitionTaskId, payload.Generation, retried+payload.AttemptOffset, time.Time{}, true, token)
	}
	logging.LogStd(logging.LOG_LEVEL_ERROR, "Quick Scan recognition failed; task state records its next action")
}

// Starts once with the application, independently of worker configuration restarts.
func StartRecognitionTaskReconciler(ctx context.Context) error {
	inspector, err := GetAsynqInspector()
	if err != nil {
		return err
	}
	go func() {
		defer inspector.Close()
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		service := services.NewRecognitionTaskService()
		for {
			if err := service.Reconcile(inspector); err != nil {
				logging.LogStd(logging.LOG_LEVEL_ERROR, "Recognition task reconciliation failed; will retry")
			}
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
			}
		}
	}()
	return nil
}
