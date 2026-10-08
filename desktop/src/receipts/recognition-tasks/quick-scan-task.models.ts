import { RecognitionTask, ReceiptStatus } from "../../open-api";

export type TaskScope = "own" | "all";
export type TaskBucket = "active" | "history" | "all";

export interface QuickScanSubmission {
  file: File;
  groupId: number;
  paidByUserId?: number;
  status?: ReceiptStatus;
  categoryIds: number[];
  tagIds: number[];
  comment?: string;
}

export interface LocalQuickScanTask {
  clientRequestId: string;
  fileName: string;
  fileSize: number;
  groupId: number;
  paidByUserId?: number;
  status?: QuickScanSubmission["status"];
  categoryIds: number[];
  tagIds: number[];
  comment?: string;
  taskId?: number;
  registering: boolean;
  sending: boolean;
  loaded: number;
  total?: number;
  awaitingConfirmation: boolean;
  reselectRequired?: boolean;
  error?: string;
}

export interface TaskListFilter {
  scope: TaskScope;
  bucket: TaskBucket;
  page: number;
  pageSize: number;
  groupId?: number;
  status?: RecognitionTask["status"];
}

export interface RecognitionTaskRow {
  id: string;
  task?: RecognitionTask;
  local?: LocalQuickScanTask;
}

export const INITIAL_TASK_FILTER: TaskListFilter = {
  scope: "own", bucket: "active", page: 1, pageSize: 25,
};

export function taskIsActive(task: RecognitionTask): boolean {
  return !["SUCCEEDED", "FAILED", "UPLOAD_INTERRUPTED"].includes(task.status);
}
