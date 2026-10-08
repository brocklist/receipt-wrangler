import { localizeUiTextPair as text } from "../../i18n/language";
import { RecognitionTask } from "../../open-api";
import { RecognitionTaskRow } from "./quick-scan-task.models";

export function taskStatusText(row: RecognitionTaskRow): string {
  if (row.local?.reselectRequired) return text("需要重新选择原文件", "Original file needs to be selected again");
  if (row.local?.registering) return text("正在创建任务", "Registering task");
  if (row.local?.awaitingConfirmation && (!row.task || ["AWAITING_UPLOAD", "UPLOADING"].includes(row.task.status))) {
    return text("上传已发送，等待服务器确认", "Upload sent, awaiting server confirmation");
  }
  if (row.local?.sending) return row.local.awaitingConfirmation
    ? text("上传完成，等待服务器确认", "Upload sent, awaiting server confirmation")
    : text("正在上传", "Uploading");
  if (!row.task) return row.local?.error
    ? text("提交状态待确认", "Submission needs confirmation")
    : text("等待上传", "Waiting for upload");
  switch (row.task.status) {
    case "AWAITING_UPLOAD": return text("等待上传", "Waiting for upload");
    case "UPLOADING": return text("正在上传", "Uploading");
    case "UPLOAD_INTERRUPTED": return text("上传已中断", "Upload interrupted");
    case "DISPATCH_PENDING": return text("文件已接收，等待入队", "File received, awaiting queue");
    case "QUEUED": return text("排队中", "Queued");
    case "RETRY_WAIT": return text("等待自动重试", "Waiting for automatic retry");
    case "SUCCEEDED": return text("识别完成", "Completed");
    case "FAILED": return text("识别失败", "Failed");
    default: return taskStageText(row.task);
  }
}

export function taskStageText(task: RecognitionTask): string {
  switch (task.stage) {
    case "PREPROCESSING": return text("正在预处理", "Preprocessing");
    case "OCR": return text("正在识别文字", "Recognizing text");
    case "AI": return text("正在提取发票信息", "Extracting receipt data");
    case "PARSING": return text("正在解析结果", "Parsing result");
    case "SAVING": return text("正在保存收据", "Saving receipt");
    case "DONE": return text("完成", "Complete");
    default: return text("处理中", "Processing");
  }
}

export function uploadPercent(row: RecognitionTaskRow): number | undefined {
  const total = row.local?.total ?? row.task?.uploadTotalBytes;
  const loaded = row.local?.loaded ?? row.task?.uploadedBytes;
  if (!total || loaded == null) return undefined;
  return Math.min(100, Math.max(0, Math.floor(loaded * 100 / total)));
}

export function elapsedText(start: string | undefined, end: string | undefined, now: number): string {
  if (!start) return "—";
  const startAt = Date.parse(start);
  const endAt = end ? Date.parse(end) : now;
  if (!Number.isFinite(startAt) || !Number.isFinite(endAt)) return "—";
  const seconds = Math.max(0, Math.floor((endAt - startAt) / 1000));
  return seconds < 60 ? seconds + text(" 秒", " sec")
    : Math.floor(seconds / 60) + text(" 分 ", " min ") + seconds % 60 + text(" 秒", " sec");
}

export function fileSizeText(bytes: number | undefined): string {
  if (bytes == null) return "—";
  return bytes < 1024 * 1024
    ? (bytes / 1024).toFixed(1) + " KB"
    : (bytes / (1024 * 1024)).toFixed(1) + " MB";
}

export function transferBytesText(bytes: number | undefined): string {
  return bytes == null ? "—" : bytes.toLocaleString() + " B";
}
