import { HttpEventType, provideHttpClient } from "@angular/common/http";
import { HttpTestingController, TestRequest, provideHttpClientTesting } from "@angular/common/http/testing";
import { TestBed } from "@angular/core/testing";
import { Store, provideStore } from "@ngxs/store";
import { SnackbarService } from "../../services";
import { AuthState } from "../../store/auth.state";
import { Logout, SetPermissions } from "../../store/auth.state.actions";
import { QuickScanTaskService } from "./quick-scan-task.service";
import { PatchQuickScanTasks, QuickScanTaskState, UpsertQuickScanTask } from "./quick-scan-task.state";
import { testPage, testTask } from "./recognition-task.test-data";

describe("QuickScanTaskService", () => {
  let service: QuickScanTaskService;
  let http: HttpTestingController;
  let store: Store;
  const file = () => new File(["data"], "invoice.png", { type: "image/png" });
  const submission = () => ({ file: file(), groupId: 2, categoryIds: [10], tagIds: [20], comment: "Client dinner" });
  const count = (): TestRequest => http.expectOne(request => request.method === "GET" && request.params.get("pageSize") === "1" && request.params.get("bucket") === "active");
  const page = (): TestRequest => http.expectOne(request => request.method === "GET" && request.params.get("pageSize") === "25");

  beforeEach(() => {
    jest.useFakeTimers();
    localStorage.clear();
    localStorage.setItem("receipt-wrangler-language", "en-US");
    TestBed.configureTestingModule({ providers: [
      provideHttpClient(), provideHttpClientTesting(), provideStore([AuthState, QuickScanTaskState]),
      { provide: SnackbarService, useValue: { error: jest.fn() } },
    ] });
    store = TestBed.inject(Store);
    store.reset({ ...store.snapshot(), auth: { userId: "1", expirationDate: (Date.now() / 1000 + 3600).toString(),
      appPermissions: [], groupPermissions: { 2: ["group.receipts.read", "group.receipts.quick-scan"] } } });
    http = TestBed.inject(HttpTestingController);
    service = TestBed.inject(QuickScanTaskService);
    TestBed.tick();
    count().flush(testPage());
  });

  afterEach(() => {
    TestBed.resetTestingModule();
    http.verify({ ignoreCancelled: true });
    jest.useRealTimers();
  });

  it("keeps HTTP Sent/progress events alive, limits transfers to two and starts the next file independently", () => {
    const ids = service.submit([submission(), submission(), submission()]);
    const registrations = http.match(request => request.method === "POST");
    expect(registrations).toHaveLength(2);
    expect(registrations[0].request.body.categoryIds).toEqual([10]);
    expect(registrations[0].request.body.comment).toBe("Client dinner");
    registrations[0].flush(testTask({ clientRequestId: ids[0] }));
    registrations[1].flush(testTask({ id: 2, clientRequestId: ids[1] }));
    const upload = http.expectOne("/api/recognitionTask/1/file");
    expect(upload.request.reportProgress).toBe(true);
    expect((upload.request.body as FormData).get("file")).toEqual(submission().file);
    upload.event({ type: HttpEventType.Sent });
    expect(upload.cancelled).toBe(false);
    upload.event({ type: HttpEventType.UploadProgress, loaded: 50, total: 100 });
    expect(service.state().local[ids[0]].loaded).toBe(50);
    upload.event({ type: HttpEventType.UploadProgress, loaded: 100, total: 100 });
    expect(service.state().local[ids[0]].awaitingConfirmation).toBe(true);
    expect(service.state().tasks[1].status).toBe("AWAITING_UPLOAD");
    upload.flush(testTask({ clientRequestId: ids[0], version: 4, status: "QUEUED", canUpload: false, uploadedBytes: 4 }));
    expect(service.state().tasks[1].status).toBe("QUEUED");
    expect(service.state().local[ids[0]].awaitingConfirmation).toBe(false);
    expect(http.match(request => request.method === "POST")).toHaveLength(1);
    expect(http.expectOne("/api/recognitionTask/2/file").cancelled).toBe(false);
  });

  it("retains task business state during sync errors, backs off and refreshes again on focus", () => {
    service.openPage();
    count().flush(testPage([], 1));
    page().flush(testPage([testTask({ status: "RUNNING", stage: "AI", canUpload: false })], 1));
    jest.advanceTimersByTime(2000);
    count().flush({}, { status: 503, statusText: "Unavailable" });
    expect(service.state().connectionError).toBe(true);
    expect(service.state().tasks[1].status).toBe("RUNNING");
    jest.advanceTimersByTime(1999);
    http.expectNone(request => request.method === "GET");
    jest.advanceTimersByTime(1);
    count().flush({}, { status: 503, statusText: "Unavailable" });
    jest.advanceTimersByTime(3999);
    http.expectNone(request => request.method === "GET");
    window.dispatchEvent(new Event("focus"));
    count().flush(testPage([], 0));
    page().flush(testPage([testTask({ status: "SUCCEEDED", version: 5, canUpload: false, receiptId: 5 })]));
    expect(service.state().connectionError).toBe(false);
    expect(service.state().tasks[1].status).toBe("SUCCEEDED");
  });

  it("starts a waiting file after another upload fails without canceling the independent transfer", () => {
    const ids = service.submit([submission(), submission(), submission()]);
    const registrations = http.match(request => request.method === "POST");
    registrations[0].flush(testTask({ clientRequestId: ids[0] }));
    registrations[1].flush(testTask({ id: 2, clientRequestId: ids[1] }));
    http.expectOne("/api/recognitionTask/1/file").error(new ProgressEvent("error"));
    expect(service.state().tasks[1].status).toBe("AWAITING_UPLOAD");
    expect(service.state().local[ids[0]].error).toContain("unconfirmed");
    expect(http.expectOne("/api/recognitionTask/2/file").cancelled).toBe(false);
    expect(http.expectOne(request => request.method === "POST").request.body.clientRequestId).toBe(ids[2]);
  });

  it("never overlaps polling requests and stops once the page closes and no own task is active", () => {
    service.openPage();
    service.refresh();
    const pending = count();
    jest.advanceTimersByTime(10000);
    http.expectNone(request => request.method === "GET");
    pending.flush(testPage());
    page().flush(testPage());
    jest.advanceTimersByTime(1);
    count().flush(testPage());
    page().flush(testPage());
    service.closePage();
    jest.advanceTimersByTime(2000);
    count().flush(testPage());
    jest.advanceTimersByTime(10000);
    http.expectNone(request => request.method === "GET");
  });

  it("queries ambiguous registrations on recovery without silently replaying a file", () => {
    const id = "11111111-1111-4111-8111-111111111111";
    localStorage.setItem("receipt-wrangler-pending-scans-1", JSON.stringify([id]));
    service.refresh();
    count().flush(testPage());
    http.expectOne(request => request.params.get("clientRequestId") === id)
      .flush(testPage([testTask({ status: "QUEUED", canUpload: false })], 1));
    http.expectNone(request => request.method === "POST" || request.method === "PUT");
    expect(service.state().tasks[1].status).toBe("QUEUED");
    expect(localStorage.getItem("receipt-wrangler-pending-scans-1")).toBeNull();
  });

  it("restores only a generic reselect notice when no task exists after a reload", () => {
    const id = "11111111-1111-4111-8111-111111111111";
    localStorage.setItem("receipt-wrangler-pending-scans-1", JSON.stringify([id]));
    service.refresh();
    count().flush(testPage());
    http.expectOne(request => request.params.get("clientRequestId") === id).flush(testPage());
    expect(service.state().reselectNoticeIds).toEqual([id]);
    expect(service.state().local[id]).toBeUndefined();
    expect(localStorage.getItem("receipt-wrangler-pending-scans-1")).toBeNull();
    expect(localStorage.getItem("receipt-wrangler-quick-scan-submissions-1")).toBeNull();
    http.expectNone(request => request.method === "POST" || request.method === "PUT");
  });

  it("preserves Quick Scan fields in memory after an ambiguous registration is confirmed missing", () => {
    const id = service.submit([submission()])[0];
    http.expectOne(request => request.method === "POST").error(new ProgressEvent("error"));
    count().flush(testPage());
    http.expectOne(request => request.params.get("clientRequestId") === id).flush(testPage());
    const local = service.state().local[id];
    expect(local).toMatchObject({
      reselectRequired: true, fileName: "invoice.png", groupId: 2, categoryIds: [10], tagIds: [20], comment: "Client dinner",
    });
    expect(localStorage.getItem("receipt-wrangler-pending-scans-1")).toBeNull();
    expect(localStorage.getItem("receipt-wrangler-quick-scan-submissions-1")).toBeNull();

    service.uploadLocalAgain(local, file());
    http.expectOne(request => request.params.get("clientRequestId") === id).flush(testPage());
    const registration = http.expectOne(request => request.method === "POST");
    expect(registration.request.body).toMatchObject({
      clientRequestId: id, fileName: "invoice.png", groupId: 2, categoryIds: [10], tagIds: [20], comment: "Client dinner",
    });
    registration.flush(testTask({ clientRequestId: id, canUpload: true }));
    http.expectOne("/api/recognitionTask/1/file").flush(testTask({ clientRequestId: id, status: "QUEUED", canUpload: false }));
  });

  it("removes in-memory reselect details immediately when group permission is revoked", () => {
    const id = "11111111-1111-4111-8111-111111111111";
    store.dispatch(new PatchQuickScanTasks({ local: { [id]: {
      clientRequestId: id, fileName: "invoice.png", fileSize: 4, groupId: 2,
      categoryIds: [], tagIds: [], registering: false, sending: false,
      loaded: 0, awaitingConfirmation: false, reselectRequired: true,
    } } }));
    store.dispatch(new SetPermissions([], { 2: [] }));
    TestBed.tick();
    count().flush(testPage());
    expect(service.state().local[id]).toBeUndefined();

    store.dispatch(new SetPermissions([], { 2: ["group.receipts.read", "group.receipts.quick-scan"] }));
    TestBed.tick();
    count().flush(testPage());
    expect(service.state().local[id]).toBeUndefined();
  });

  it("retains unresolved request IDs when access is denied and reconciles immediately after permission restoration", () => {
    const id = "11111111-1111-4111-8111-111111111111";
    localStorage.setItem("receipt-wrangler-pending-scans-1", JSON.stringify([id]));
    store.dispatch(new SetPermissions([], { 2: [] }));
    TestBed.tick();
    count().flush(testPage());
    http.expectOne(request => request.params.get("clientRequestId") === id)
      .flush({}, { status: 403, statusText: "Forbidden" });
    expect(service.state().permissionBlockedRequestIds).toEqual([id]);
    expect(JSON.parse(localStorage.getItem("receipt-wrangler-pending-scans-1")!)).toEqual([id]);
    jest.advanceTimersByTime(5000);
    http.expectNone(request => request.method === "GET");
    http.expectNone(request => request.method === "POST" || request.method === "PUT");

    store.dispatch(new SetPermissions([], { 2: ["group.receipts.read", "group.receipts.quick-scan"] }));
    TestBed.tick();
    count().flush(testPage());
    http.expectOne(request => request.params.get("clientRequestId") === id)
      .flush(testPage([testTask({ clientRequestId: id })], 1));
    expect(service.state().permissionBlockedRequestIds).toEqual([]);
    expect(service.state().tasks[1].clientRequestId).toBe(id);
    expect(localStorage.getItem("receipt-wrangler-pending-scans-1")).toBeNull();
  });

  it("clears details on access denial without labeling the task failed or retrying the read", () => {
    store.dispatch(new UpsertQuickScanTask(testTask({ status: "RUNNING", stage: "AI" })));
    service.openPage();
    count().flush({}, { status: 403, statusText: "Forbidden" });
    expect(service.state().tasks).toEqual({});
    expect(service.state().connectionError).toBe(false);
    expect(service.state().readError).toContain("no longer accessible");
    jest.advanceTimersByTime(30000);
    http.expectNone(request => request.method === "GET");
  });

  it("invalidates cached details immediately when current group read permission changes", () => {
    store.dispatch(new UpsertQuickScanTask(testTask({ status: "SUCCEEDED", receiptId: 9 })));
    store.dispatch(new SetPermissions([], { 2: [] }));
    TestBed.tick();
    expect(service.state().tasks).toEqual({});
    count().flush(testPage());
  });

  it("requires server upload permission and the original file name/size for whole-file reupload", () => {
    service.uploadAgain(testTask({ canUpload: false }), file());
    service.uploadAgain(testTask(), new File(["x"], "other.png"));
    http.expectNone(request => request.method === "PUT");
    service.uploadAgain(testTask({ status: "UPLOAD_INTERRUPTED" }), file());
    const upload = http.expectOne("/api/recognitionTask/1/file");
    expect(upload.request.method).toBe("PUT");
    expect((upload.request.body as FormData).get("file")).toEqual(file());
  });

  it("does not grant retry from global read access and deduplicates an allowed retry click", () => {
    const task = testTask({ status: "FAILED", canUpload: false });
    service.retry(task);
    http.expectNone(request => request.method === "POST");
    service.retry({ ...task, canRetry: true });
    service.retry({ ...task, canRetry: true });
    const retry = http.expectOne("/api/recognitionTask/1/retry");
    expect(retry.request.body).toEqual({ version: 1 });
    retry.flush({ ...task, status: "DISPATCH_PENDING", canRetry: false, version: 2 });
    expect(service.state().tasks[1].status).toBe("DISPATCH_PENDING");
  });

  it("cancels live transfers and clears files, pending IDs and private state on logout", () => {
    service.submit([submission(), submission(), submission()]);
    const registrations = http.match(request => request.method === "POST");
    expect(registrations).toHaveLength(2);
    store.dispatch(new Logout());
    TestBed.tick();
    expect(registrations.every(request => request.cancelled)).toBe(true);
    expect(service.state().tasks).toEqual({});
    expect(service.state().local).toEqual({});
    expect(service.submit([submission()])).toEqual([]);
    expect(localStorage.getItem("receipt-wrangler-pending-scans-1")).toBeNull();
    jest.advanceTimersByTime(30000);
    http.expectNone(request => request.method === "GET" || request.method === "POST");
  });
});
