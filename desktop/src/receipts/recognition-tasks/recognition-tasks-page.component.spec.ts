import { provideHttpClient } from "@angular/common/http";
import { signal } from "@angular/core";
import { ComponentFixture, TestBed } from "@angular/core/testing";
import { provideNoopAnimations } from "@angular/platform-browser/animations";
import { provideRouter } from "@angular/router";
import { Store, provideStore } from "@ngxs/store";
import { AuthState } from "../../store/auth.state";
import { SetPermissions } from "../../store/auth.state.actions";
import { GroupState } from "../../store/group.state";
import { UserState } from "../../store/user.state";
import { QuickScanTaskService } from "./quick-scan-task.service";
import { initialTaskState } from "./quick-scan-task.state";
import { RecognitionTasksPageComponent } from "./recognition-tasks-page.component";
import { testTask } from "./recognition-task.test-data";

describe("RecognitionTasksPageComponent", () => {
  let fixture: ComponentFixture<RecognitionTasksPageComponent>;
  let store: Store;
  const mock = { state: signal(initialTaskState()), now: signal(Date.now()), openPage: jest.fn(), closePage: jest.fn(),
    setFilter: jest.fn(), refresh: jest.fn(), retry: jest.fn(), uploadAgain: jest.fn(), dismissReselectNotices: jest.fn() };
  beforeEach(async () => {
    jest.clearAllMocks();
    localStorage.setItem("receipt-wrangler-language", "en-US");
    mock.state.set({ ...initialTaskState(), pageIds: [1], tasks: { 1: testTask({ status: "QUEUED", canUpload: false }) }, totalCount: 1, activeCount: 3 });
    await TestBed.configureTestingModule({ imports: [RecognitionTasksPageComponent], providers: [
      provideHttpClient(), provideRouter([]), provideNoopAnimations(), provideStore([AuthState, GroupState, UserState]),
      { provide: QuickScanTaskService, useValue: mock },
    ] }).compileComponents();
    store = TestBed.inject(Store);
    fixture = TestBed.createComponent(RecognitionTasksPageComponent);
    fixture.detectChanges();
    await fixture.whenStable();
  });

  it("shows server counts, paged task rows and filters for the current user", () => {
    expect(fixture.nativeElement.textContent).toContain("invoice.png");
    expect(fixture.nativeElement.textContent).toContain("Recognition tasks");
    expect(fixture.nativeElement.textContent).not.toContain("All users");
    fixture.componentInstance.groupFilter.setValue(2);
    expect(mock.setFilter).toHaveBeenCalledWith({ groupId: 2, page: 1 });
    fixture.componentInstance.pageChanged({ pageIndex: 1, pageSize: 25, length: 30 });
    expect(mock.setFilter).toHaveBeenCalledWith({ page: 2, pageSize: 25 });
  });

  it("offers the all-users view only with global read permission", async () => {
    store.dispatch(new SetPermissions(["app.system-tasks.read"], {}));
    await fixture.whenStable();
    expect(fixture.nativeElement.textContent).toContain("All users");
    fixture.componentInstance.scopeChanged("all");
    expect(mock.setFilter).toHaveBeenCalledWith({ scope: "all", page: 1 });
  });

  it("keeps known task rows visible during a network reconnect and closes its page subscription", async () => {
    mock.state.update(state => ({ ...state, connectionError: true }));
    await fixture.whenStable();
    expect(fixture.nativeElement.querySelector('[data-testid="recognition-reconnect"]')).toBeTruthy();
    expect(fixture.nativeElement.textContent).toContain("invoice.png");
    fixture.destroy();
    expect(mock.closePage).toHaveBeenCalledTimes(1);
  });

  it("distinguishes reselect guidance from permission-blocked request recovery", async () => {
    mock.state.update(state => ({ ...state,
      reselectNoticeIds: ["11111111-1111-4111-8111-111111111111"],
      permissionBlockedRequestIds: ["22222222-2222-4222-8222-222222222222"],
    }));
    await fixture.whenStable();
    expect(fixture.nativeElement.querySelector('[data-testid="recognition-reselect-notice"]').textContent)
      .toContain("Select the original files");
    expect(fixture.nativeElement.querySelector('[data-testid="recognition-permission-reconnect"]').textContent)
      .toContain("will not replay uploads automatically");
  });

  it("shows files that need to be reselected after the server confirms an empty registration lookup", async () => {
    const id = "11111111-1111-4111-8111-111111111111";
    mock.state.update(state => ({ ...state, local: { [id]: {
      clientRequestId: id, fileName: "invoice.png", fileSize: 4, groupId: 2,
      categoryIds: [], tagIds: [], registering: false, sending: false, loaded: 0,
      awaitingConfirmation: false, reselectRequired: true,
    } } }));
    await fixture.whenStable();
    expect(fixture.nativeElement.querySelector('[data-testid="recognition-reselect-required"]')).toBeTruthy();
    expect(fixture.nativeElement.textContent).toContain("invoice.png");
    expect(fixture.nativeElement.querySelector('[data-testid="recognition-reselect"]')).toBeTruthy();
    expect(fixture.nativeElement.querySelector('[data-testid="recognition-empty"]')).toBeNull();
  });
});
