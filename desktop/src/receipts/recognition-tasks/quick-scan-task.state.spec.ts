import { TestBed } from "@angular/core/testing";
import { Store, provideStore } from "@ngxs/store";
import { ClearQuickScanTasks, QuickScanTaskState, SetQuickScanTaskPage, UpsertQuickScanTask } from "./quick-scan-task.state";
import { testPage, testTask } from "./recognition-task.test-data";

describe("QuickScanTaskState", () => {
  let store: Store;
  beforeEach(() => {
    TestBed.configureTestingModule({ providers: [provideStore([QuickScanTaskState])] });
    store = TestBed.inject(Store);
  });

  it("does not replace newer progress with a stale page or upload response", () => {
    store.dispatch(new UpsertQuickScanTask(testTask({ version: 5, status: "SUCCEEDED", receiptId: 10 })));
    store.dispatch(new UpsertQuickScanTask(testTask({ version: 2, status: "RUNNING" })));
    store.dispatch(new SetQuickScanTaskPage(testPage([testTask({ version: 3, status: "QUEUED" })], 1)));
    expect(store.selectSnapshot(QuickScanTaskState.model).tasks[1].status).toBe("SUCCEEDED");
    expect(store.selectSnapshot(QuickScanTaskState.model).pageIds).toEqual([1]);
  });

  it("clears private task state and filters when the session ends", () => {
    store.dispatch(new SetQuickScanTaskPage(testPage([testTask()], 1)));
    store.dispatch(new ClearQuickScanTasks());
    const state = store.selectSnapshot(QuickScanTaskState.model);
    expect(state.tasks).toEqual({});
    expect(state.local).toEqual({});
    expect(state.ownActiveCount).toBe(0);
    expect(state.filter.scope).toBe("own");
  });
});
