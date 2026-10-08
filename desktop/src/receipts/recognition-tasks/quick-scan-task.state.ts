import { Injectable } from "@angular/core";
import { Action, Selector, State, StateContext } from "@ngxs/store";
import { GetRecognitionTasksResponse, RecognitionTask } from "../../open-api";
import { INITIAL_TASK_FILTER, LocalQuickScanTask, TaskListFilter } from "./quick-scan-task.models";

export interface QuickScanTaskModel {
  tasks: { [id: number]: RecognitionTask };
  local: { [requestId: string]: LocalQuickScanTask };
  pageIds: number[];
  filter: TaskListFilter;
  totalCount: number;
  activeCount: number;
  awaitingUploadCount: number;
  runningCount: number;
  failedCount: number;
  ownActiveCount: number;
  reselectNoticeIds: string[];
  permissionBlockedRequestIds: string[];
  syncing: boolean;
  connectionError: boolean;
  readError?: string;
  lastSyncedAt?: number;
}

export const initialTaskState = (): QuickScanTaskModel => ({
  tasks: {}, local: {}, pageIds: [], filter: { ...INITIAL_TASK_FILTER },
  totalCount: 0, activeCount: 0, awaitingUploadCount: 0, runningCount: 0,
  failedCount: 0, ownActiveCount: 0, reselectNoticeIds: [], permissionBlockedRequestIds: [],
  syncing: false, connectionError: false,
});

export class PatchQuickScanTasks {
  static readonly type = "[Quick Scan Tasks] Patch";
  constructor(public patch: Partial<QuickScanTaskModel>) {}
}

export class UpsertQuickScanTask {
  static readonly type = "[Quick Scan Tasks] Upsert task";
  constructor(public task: RecognitionTask) {}
}

export class PatchLocalQuickScanTask {
  static readonly type = "[Quick Scan Tasks] Patch local transfer";
  constructor(public requestId: string, public patch: Partial<LocalQuickScanTask>) {}
}

export class SetQuickScanTaskPage {
  static readonly type = "[Quick Scan Tasks] Page";
  constructor(public page: GetRecognitionTasksResponse) {}
}

export class ClearQuickScanTasks {
  static readonly type = "[Quick Scan Tasks] Clear";
}

@State<QuickScanTaskModel>({ name: "quickScanTasks", defaults: initialTaskState() })
@Injectable()
export class QuickScanTaskState {
  @Selector()
  static model(state: QuickScanTaskModel): QuickScanTaskModel { return state; }

  @Selector()
  static ownActiveCount(state: QuickScanTaskModel): number { return state.ownActiveCount; }

  @Action(PatchQuickScanTasks)
  patch(ctx: StateContext<QuickScanTaskModel>, action: PatchQuickScanTasks): void {
    ctx.patchState(action.patch);
  }

  @Action(UpsertQuickScanTask)
  upsert(ctx: StateContext<QuickScanTaskModel>, { task }: UpsertQuickScanTask): void {
    const current = ctx.getState().tasks[task.id];
    if (current && current.version > task.version) return;
    ctx.patchState({ tasks: { ...ctx.getState().tasks, [task.id]: task } });
  }

  @Action(PatchLocalQuickScanTask)
  patchLocal(ctx: StateContext<QuickScanTaskModel>, action: PatchLocalQuickScanTask): void {
    const state = ctx.getState();
    const current = state.local[action.requestId];
    if (!current) return;
    ctx.patchState({ local: { ...state.local, [action.requestId]: { ...current, ...action.patch } } });
  }

  @Action(SetQuickScanTaskPage)
  page(ctx: StateContext<QuickScanTaskModel>, { page }: SetQuickScanTaskPage): void {
    const tasks = { ...ctx.getState().tasks };
    for (const task of page.data) {
      if (!tasks[task.id] || tasks[task.id].version <= task.version) tasks[task.id] = task;
    }
    ctx.patchState({
      tasks, pageIds: page.data.map(task => task.id), totalCount: page.totalCount,
      activeCount: page.activeCount, awaitingUploadCount: page.awaitingUploadCount,
      runningCount: page.runningCount, failedCount: page.failedCount,
    });
  }

  @Action(ClearQuickScanTasks)
  clear(ctx: StateContext<QuickScanTaskModel>): void { ctx.setState(initialTaskState()); }
}
