import { GetRecognitionTasksResponse, RecognitionTask } from "../../open-api";

export function testTask(patch: Partial<RecognitionTask> = {}): RecognitionTask {
  return {
    id: 1, clientRequestId: "11111111-1111-4111-8111-111111111111", fileName: "invoice.png", fileSize: 4,
    groupId: 2, ownerUserId: 1, version: 1, status: "AWAITING_UPLOAD", stage: "UPLOAD",
    createdAt: "2026-10-07T01:00:00Z", updatedAt: "2026-10-07T01:00:00Z", uploadedBytes: 0,
    uploadTotalBytes: 4, attempt: 0, maxAttempts: 4, fallbackActive: false,
    errorCode: "", errorMessage: "", canUpload: true, canRetry: false, ...patch,
  };
}

export function testPage(data: RecognitionTask[] = [], activeCount = 0): GetRecognitionTasksResponse {
  return { data, totalCount: data.length, activeCount, awaitingUploadCount: 0, runningCount: 0, failedCount: 0 };
}
