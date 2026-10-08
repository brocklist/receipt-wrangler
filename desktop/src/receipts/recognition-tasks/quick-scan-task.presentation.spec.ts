import { elapsedText, taskStatusText, uploadPercent } from "./quick-scan-task.presentation";
import { testTask } from "./recognition-task.test-data";

describe("Quick Scan task presentation", () => {
  beforeEach(() => localStorage.setItem("receipt-wrangler-language", "en-US"));

  it("uses byte progress and waits for server confirmation at 100 percent sent", () => {
    const local = { clientRequestId: "one", fileName: "one", fileSize: 10, groupId: 1,
      registering: false, sending: true, loaded: 10, total: 10, awaitingConfirmation: true };
    const row = { id: "one", task: testTask({ status: "UPLOADING" }), local };
    expect(uploadPercent(row)).toBe(100);
    expect(taskStatusText(row)).toContain("awaiting server confirmation");
    expect(uploadPercent({ ...row, local: { ...local, total: undefined }, task: undefined })).toBeUndefined();
  });

  it("renders stages and terminal states in both supported languages", () => {
    expect(taskStatusText({ id: "1", task: testTask({ status: "RUNNING", stage: "AI" }) })).toBe("Extracting receipt data");
    localStorage.setItem("receipt-wrangler-language", "zh-CN");
    expect(taskStatusText({ id: "1", task: testTask({ status: "RUNNING", stage: "OCR" }) })).toBe("正在识别文字");
    expect(taskStatusText({ id: "1", task: testTask({ status: "UPLOAD_INTERRUPTED" }) })).toBe("上传已中断");
  });

  it("uses actual timestamps for elapsed time, with no estimates", () => {
    expect(elapsedText("2026-10-07T01:00:00Z", "2026-10-07T01:01:05Z", 0)).toBe("1 min 5 sec");
    expect(elapsedText(undefined, undefined, Date.now())).toBe("—");
    expect(elapsedText("invalid", undefined, Date.now())).toBe("—");
  });
});
