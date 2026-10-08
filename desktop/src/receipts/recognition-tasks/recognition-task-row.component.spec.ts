import { signal } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { provideRouter } from "@angular/router";
import { provideStore } from "@ngxs/store";
import { GroupState } from "../../store/group.state";
import { UserState } from "../../store/user.state";
import { QuickScanTaskService } from "./quick-scan-task.service";
import { RecognitionTaskRowComponent } from "./recognition-task-row.component";
import { testTask } from "./recognition-task.test-data";

describe("RecognitionTaskRowComponent", () => {
  let fixture: ComponentFixture<RecognitionTaskRowComponent>;
  const mock = { now: signal(Date.parse("2026-10-07T01:01:00Z")), retry: jest.fn(), uploadAgain: jest.fn(),
    uploadLocalAgain: jest.fn(), retrySubmission: jest.fn() };
  beforeEach(() => {
    localStorage.setItem("receipt-wrangler-language", "en-US");
    jest.clearAllMocks();
    TestBed.configureTestingModule({ imports: [RecognitionTaskRowComponent], providers: [
      provideRouter([]), provideStore([GroupState, UserState]), { provide: QuickScanTaskService, useValue: mock },
    ] });
    fixture = TestBed.createComponent(RecognitionTaskRowComponent);
  });
  const render = (patch: Parameters<typeof testTask>[0]) => {
    fixture.componentRef.setInput("row", { id: "1", task: testTask(patch) });
    fixture.detectChanges();
  };

  it("shows recognition stage and actual elapsed time without a fabricated percentage", () => {
    render({ status: "RUNNING", stage: "AI", startedAt: "2026-10-07T01:00:00Z", stageStartedAt: "2026-10-07T01:00:40Z", canUpload: false });
    expect(fixture.nativeElement.textContent).toContain("Extracting receipt data");
    expect(fixture.nativeElement.textContent).toContain("20 sec");
    expect(fixture.nativeElement.textContent).not.toContain("%");
    expect(fixture.nativeElement.querySelector("mat-progress-bar").getAttribute("mode")).toBe("indeterminate");
  });

  it("shows only server-authorized retry and receipt actions", () => {
    render({ status: "FAILED", canUpload: false, canRetry: false, errorMessage: "Recognition failed", attempt: 4 });
    expect(fixture.nativeElement.textContent).toContain("Recognition failed");
    expect(fixture.nativeElement.querySelector('[data-testid="recognition-retry"]')).toBeNull();
    render({ status: "FAILED", canUpload: false, canRetry: true });
    fixture.nativeElement.querySelector('[data-testid="recognition-retry"] button').click();
    expect(mock.retry).toHaveBeenCalled();
    render({ status: "SUCCEEDED", canUpload: false, receiptId: 12 });
    expect(fixture.nativeElement.querySelector('[data-testid="recognition-open-receipt"]')).toBeTruthy();
  });

  it("allows file selection for an interrupted upload and clears the file input afterward", () => {
    render({ status: "UPLOAD_INTERRUPTED", canUpload: true, uploadedBytes: 64, uploadTotalBytes: 100 });
    expect(fixture.nativeElement.querySelector('[data-testid="recognition-upload-bytes"]').textContent).toContain("64 B");
    expect(fixture.nativeElement.textContent).toContain("64%");
    expect(fixture.nativeElement.textContent).toContain("Received by server; reselect to restart from the beginning");
    expect(fixture.nativeElement.textContent).toContain("restarts from the beginning");
    expect(fixture.nativeElement.textContent).not.toContain("Upload complete");
    expect(fixture.nativeElement.querySelector("mat-progress-bar")).toBeNull();
    const selected = new File(["data"], "invoice.png");
    const input = { files: [selected], value: "C:\\fakepath\\invoice.png" };
    fixture.componentInstance.chooseFile({ target: input } as unknown as Event);
    expect(mock.uploadAgain).toHaveBeenCalledWith(expect.objectContaining({ id: 1 }), selected);
    expect(input.value).toBe("");
  });

  it("does not show transfer progress while the file is waiting to be uploaded", () => {
    render({ status: "AWAITING_UPLOAD", canUpload: true, uploadedBytes: 2 });
    expect(fixture.nativeElement.textContent).toContain("Waiting for upload");
    expect(fixture.nativeElement.querySelector("mat-progress-bar")).toBeNull();
  });

  it("shows a reselect action for a file confirmed absent by the server", () => {
    const local = { clientRequestId: "11111111-1111-4111-8111-111111111111", fileName: "invoice.png",
      fileSize: 4, groupId: 2, categoryIds: [10], tagIds: [20], registering: false, sending: false,
      loaded: 0, awaitingConfirmation: false, reselectRequired: true };
    fixture.componentRef.setInput("row", { id: local.clientRequestId, local });
    fixture.detectChanges();
    expect(fixture.nativeElement.textContent).toContain("server confirmed this file was not received");
    const selected = new File(["data"], "invoice.png");
    const input = { files: [selected], value: "C:\\fakepath\\invoice.png" };
    fixture.componentInstance.chooseFile({ target: input } as unknown as Event);
    expect(mock.uploadLocalAgain).toHaveBeenCalledWith(local, selected);
    expect(input.value).toBe("");
  });
});
