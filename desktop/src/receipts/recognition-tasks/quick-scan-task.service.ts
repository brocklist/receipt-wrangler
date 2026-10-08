import { HttpContext, HttpErrorResponse, HttpEvent, HttpEventType } from "@angular/common/http";
import { DOCUMENT } from "@angular/common";
import { DestroyRef, Injectable, effect, inject, untracked } from "@angular/core";
import { toSignal } from "@angular/core/rxjs-interop";
import { Store } from "@ngxs/store";
import { EMPTY, Observable, Subscription, catchError, concatMap, finalize, from, interval, map, of, tap, throwError } from "rxjs";
import { localizeUiTextPair as text } from "../../i18n/language";
import { HANDLE_ERROR_LOCALLY } from "../../interceptors/local-error.context";
import { Permission, RecognitionTask, RecognitionTaskService } from "../../open-api";
import { SnackbarService } from "../../services";
import { AuthState } from "../../store";
import { LocalQuickScanTask, QuickScanSubmission, TaskListFilter, taskIsActive } from "./quick-scan-task.models";
import { ClearQuickScanTasks, PatchLocalQuickScanTask, PatchQuickScanTasks,
  QuickScanTaskState, SetQuickScanTaskPage, UpsertQuickScanTask } from "./quick-scan-task.state";

@Injectable({ providedIn: "root" })
export class QuickScanTaskService {
  private readonly store = inject(Store);
  private readonly api = inject(RecognitionTaskService);
  private readonly snackbar = inject(SnackbarService);
  private readonly document = inject(DOCUMENT);
  private readonly destroyRef = inject(DestroyRef);
  private readonly localOptions = { context: new HttpContext().set(HANDLE_ERROR_LOCALLY, true) };
  public readonly state = this.store.selectSignal(QuickScanTaskState.model);
  public readonly ownActiveCount = this.store.selectSignal(QuickScanTaskState.ownActiveCount);
  public readonly now = toSignal(interval(1000).pipe(map(() => Date.now())), { initialValue: Date.now() });
  private owner = "";
  private pageUsers = 0;
  private timer?: ReturnType<typeof setTimeout>;
  private poll?: Subscription;
  private pollRunning = false;
  private refreshPending = false;
  private failures = 0;
  private queue: string[] = [];
  private activeTransfers = 0;
  private readonly sources = new Map<string, QuickScanSubmission>();
  private readonly inaccessiblePending = new Set<string>();
  private readonly transfers = new Map<string, Subscription>();
  private readonly retrying = new Set<number>();
  private readonly mutations = new Set<Subscription>();

  constructor() {
    const loggedIn = this.store.selectSignal(AuthState.isLoggedIn);
    const userId = this.store.selectSignal(AuthState.userId);
    const appPermissions = this.store.selectSignal(AuthState.appPermissions);
    const groupPermissions = this.store.selectSignal(AuthState.groupPermissions);
    let previousPermissions = "";
    effect(() => {
      const nextOwner = loggedIn() ? userId() : "";
      const permissions = JSON.stringify([appPermissions(), groupPermissions()]);
      untracked(() => {
        if (this.owner !== nextOwner) {
          this.clearSession();
          this.owner = nextOwner;
          if (this.owner) {
            this.clearLegacySubmissionMetadata();
            this.refresh();
          }
        } else if (this.owner && previousPermissions !== permissions) {
          this.poll?.unsubscribe();
          const filter = this.state().filter;
          this.store.dispatch(new PatchQuickScanTasks({
            tasks: {}, pageIds: [], totalCount: 0, ownActiveCount: 0, readError: undefined,
            permissionBlockedRequestIds: [],
            filter: { ...filter, scope: appPermissions().includes(Permission.AppSystemTasksRead) ? filter.scope : "own" },
          }));
          this.inaccessiblePending.clear();
          this.pruneUnauthorizedReselectRows();
          this.refresh();
        }
        previousPermissions = permissions;
      });
    });
    const reconnect = () => this.refresh();
    this.document.defaultView?.addEventListener("online", reconnect);
    this.document.defaultView?.addEventListener("focus", reconnect);
    this.destroyRef.onDestroy(() => {
      this.document.defaultView?.removeEventListener("online", reconnect);
      this.document.defaultView?.removeEventListener("focus", reconnect);
      this.clearSession();
    });
  }

  public submit(files: QuickScanSubmission[]): string[] {
    if (!this.owner) return [];
    const local = { ...this.state().local };
    const ids = files.map(file => {
      const id = crypto.randomUUID();
      this.sources.set(id, file);
      local[id] = {
        clientRequestId: id, fileName: file.file.name, fileSize: file.file.size,
        groupId: file.groupId, registering: false, sending: false, loaded: 0,
        paidByUserId: file.paidByUserId, status: file.status,
        categoryIds: [...file.categoryIds], tagIds: [...file.tagIds],
        comment: file.comment,
        awaitingConfirmation: false,
      };
      this.queue.push(id);
      return id;
    });
    this.store.dispatch(new PatchQuickScanTasks({ local }));
    this.pump();
    return ids;
  }

  public openPage(): void { this.pageUsers++; this.refresh(); }
  public closePage(): void { this.pageUsers = Math.max(0, this.pageUsers - 1); }
  public dismissReselectNotices(): void { this.store.dispatch(new PatchQuickScanTasks({ reselectNoticeIds: [] })); }

  public setFilter(patch: Partial<TaskListFilter>): void {
    this.store.dispatch(new PatchQuickScanTasks({
      filter: { ...this.state().filter, ...patch },
      pageIds: [], totalCount: 0,
    }));
    this.refresh();
  }

  public refresh(): void {
    if (!this.owner) return;
    clearTimeout(this.timer);
    if (this.pollRunning) { this.refreshPending = true; return; }
    const owner = this.owner;
    let failed = false;
    this.pollRunning = true;
    this.store.dispatch(new PatchQuickScanTasks({ syncing: true }));
    const requests: (() => Observable<unknown>)[] = [
      () => this.api.getRecognitionTasks("own", "active", 1, 1, undefined, undefined, undefined, undefined,
        "body", false, this.localOptions).pipe(tap(page => {
          this.store.dispatch(new PatchQuickScanTasks({ ownActiveCount: page.activeCount }));
        })),
    ];
    if (this.pageUsers) {
      const filter = this.state().filter;
      requests.push(() => this.api.getRecognitionTasks(filter.scope, filter.bucket, filter.page, filter.pageSize,
        undefined, undefined, filter.groupId, filter.status, "body", false, this.localOptions)
        .pipe(tap(page => {
          if (this.state().filter === filter) this.store.dispatch(new SetQuickScanTaskPage(page));
        })));
    }
    const ids = Object.values(this.state().local)
      .map(local => local.taskId).filter((id): id is number => id != null)
      .filter(id => !this.state().tasks[id] || taskIsActive(this.state().tasks[id]));
    for (let start = 0; start < ids.length; start += 100) {
      const chunk = ids.slice(start, start + 100).join(",");
      requests.push(() => this.api.getRecognitionTasks("own", "all", 1, 100, undefined, chunk,
        undefined, undefined, "body", false, this.localOptions).pipe(tap(page => {
          page.data.forEach(task => this.acceptTask(task));
          this.removeInaccessibleTasks(ids.slice(start, start + 100), page.data.map(task => task.id));
        })));
    }
    const pendingScope = this.store.selectSnapshot(AuthState.appPermissions).includes(Permission.AppSystemTasksRead) ? "all" : "own";
    for (const id of this.pendingIds().filter(id => !this.inaccessiblePending.has(id) &&
        !this.transfers.has(id) && !this.state().local[id]?.registering)) {
      requests.push(() => this.api.getRecognitionTasks(pendingScope, "all", 1, 1, id, undefined,
        undefined, undefined, "body", false, this.localOptions).pipe(tap(page => {
          if (page.data.length) {
            this.inaccessiblePending.delete(id);
            this.acceptTask(page.data[0]);
            this.patchLocal(id, { taskId: page.data[0].id, registering: false });
            this.patchPermissionBlockedRequest(id, false);
          } else {
            this.requireReselect(id);
          }
          this.rememberPending(id, false);
        }), catchError((error: HttpErrorResponse) => {
          if (error.status !== 403) return throwError(() => error);
          this.inaccessiblePending.add(id);
          this.patchPermissionBlockedRequest(id, true);
          return of(undefined);
        })));
    }
    const sub = from(requests).pipe(
      concatMap(request => request()),
      catchError((error: HttpErrorResponse) => {
        failed = true;
        if (owner === this.owner) {
          if (error.status === 403 && this.state().filter.scope === "all") {
            this.store.dispatch(new PatchQuickScanTasks({
              tasks: {}, pageIds: [], totalCount: 0, connectionError: false, readError: undefined,
              filter: { ...this.state().filter, scope: "own", page: 1 },
            }));
            this.refreshPending = true;
          } else if (error.status === 0 || error.status >= 500) {
            this.failures++;
            this.store.dispatch(new PatchQuickScanTasks({ connectionError: true, readError: undefined }));
          } else {
            this.failures = 0;
            this.store.dispatch(new PatchQuickScanTasks({
              tasks: {}, local: {}, pageIds: [], totalCount: 0, ownActiveCount: 0,
              activeCount: 0, awaitingUploadCount: 0, runningCount: 0, failedCount: 0,
              connectionError: false,
              readError: [401, 403, 404].includes(error.status)
                ? text("这些任务已不可访问，请确认登录状态和当前权限。", "These tasks are no longer accessible. Check your login and current permissions.")
                : text("暂时无法读取任务，请刷新后重试。", "Cannot read these tasks. Refresh to try again."),
            }));
          }
        }
        return EMPTY;
      }),
      finalize(() => {
        this.pollRunning = false;
        if (owner !== this.owner) return;
        this.store.dispatch(new PatchQuickScanTasks({ syncing: false }));
        if (this.refreshPending) {
          this.refreshPending = false;
          this.timer = setTimeout(() => this.refresh(), 0);
        } else if (!this.state().readError && (this.pageUsers || this.state().ownActiveCount || this.activeTransfers || this.failures ||
            this.pendingIds().some(id => !this.inaccessiblePending.has(id)))) {
          const delay = this.failures ? [2000, 4000, 8000, 15000][Math.min(3, this.failures - 1)] : 2000;
          this.timer = setTimeout(() => this.refresh(), delay);
        }
      }),
    ).subscribe({
      complete: () => {
        if (owner === this.owner && !failed) {
          this.failures = 0;
          this.store.dispatch(new PatchQuickScanTasks({ connectionError: false, readError: undefined, lastSyncedAt: Date.now() }));
        }
      },
    });
    if (!sub.closed) this.poll = sub;
  }

  public uploadAgain(task: RecognitionTask, file: File): void {
    if (!task.canUpload || task.ownerUserId.toString() !== this.owner) return;
    if (file.name !== task.fileName || file.size !== task.fileSize) {
      this.snackbar.error(text("请选择原文件，文件名和大小必须与任务一致。", "Select the original file with the same name and size."));
      return;
    }
    const id = task.clientRequestId;
    if (this.queue.includes(id) || this.transfers.has(id)) return;
    this.acceptTask(task);
    this.sources.set(id, { file, groupId: task.groupId, categoryIds: [], tagIds: [] });
    this.store.dispatch(new PatchQuickScanTasks({ local: {
      ...this.state().local, [id]: {
        clientRequestId: id, taskId: task.id, fileName: task.fileName, fileSize: task.fileSize,
        groupId: task.groupId, categoryIds: [], tagIds: [], registering: false, sending: false,
        loaded: 0, awaitingConfirmation: false,
      },
    } }));
    this.queue.push(id);
    this.pump();
  }

  public uploadLocalAgain(local: LocalQuickScanTask, file: File): void {
    if (!local.reselectRequired || local.taskId || local.clientRequestId.length !== 36) return;
    const groupPermissions = this.store.selectSnapshot(AuthState.groupPermissions)[local.groupId] ?? [];
    if (!groupPermissions.includes(Permission.GroupReceiptsRead) || !groupPermissions.includes(Permission.GroupReceiptsQuickScan)) {
      this.snackbar.error(text("当前没有查看或上传此分组收据的权限。", "You no longer have permission to view or upload receipts in this group."));
      return;
    }
    if (file.name !== local.fileName || file.size !== local.fileSize) {
      this.snackbar.error(text("请选择原文件，文件名和大小必须与任务一致。", "Select the original file with the same name and size."));
      return;
    }
    const source: QuickScanSubmission = {
      file, groupId: local.groupId, paidByUserId: local.paidByUserId, status: local.status,
      categoryIds: [...local.categoryIds], tagIds: [...local.tagIds],
      comment: local.comment,
    };
    this.sources.set(local.clientRequestId, source);
    const restored = { ...local, reselectRequired: false, error: undefined, loaded: 0, total: undefined };
    this.store.dispatch(new PatchLocalQuickScanTask(local.clientRequestId, {
      reselectRequired: false, error: undefined, loaded: 0, total: undefined,
    }));
    this.retrySubmission(local.clientRequestId);
  }

  public retry(task: RecognitionTask): void {
    if (!task.canRetry || this.retrying.has(task.id)) return;
    this.retrying.add(task.id);
    this.trackMutation(this.api.retryRecognitionTask(task.id, { version: task.version }).pipe(
      tap(result => { this.acceptTask(result); this.refresh(); }),
      catchError(() => { this.refresh(); return EMPTY; }),
      finalize(() => this.retrying.delete(task.id)),
    ).subscribe());
  }

  public retrySubmission(id: string): void {
    if (!this.sources.has(id) || this.queue.includes(id) || this.transfers.has(id)) return;
    // Reconcile first: an interrupted response must never create a second logical task.
    this.trackMutation(this.api.getRecognitionTasks("own", "all", 1, 1, id, undefined, undefined, undefined,
      "body", false, this.localOptions).pipe(
      tap(page => {
        const task = page.data[0];
        if (task) { this.acceptTask(task); this.patchLocal(id, { taskId: task.id }); }
        if (!task || task.canUpload) { this.queue.push(id); this.pump(); }
      }),
      catchError(() => { this.patchLocal(id, { error: text("暂时无法确认提交状态，请稍后重试。", "Cannot confirm submission yet. Try again later.") }); return EMPTY; }),
    ).subscribe());
  }

  private pump(): void {
    while (this.owner && this.activeTransfers < 2 && this.queue.length) {
      const id = this.queue.shift()!;
      const source = this.sources.get(id);
      const local = this.state().local[id];
      if (!source || !local) continue;
      this.activeTransfers++;
      const owner = this.owner;
      const existing = local.taskId ? this.state().tasks[local.taskId] : undefined;
      this.patchLocal(id, { registering: !existing, error: undefined, loaded: 0, total: undefined, awaitingConfirmation: false });
      if (!existing) this.rememberPending(id, true);
      const registered = existing ? of(existing) : this.api.createRecognitionTask({
        clientRequestId: id, fileName: source.file.name, fileSize: source.file.size,
        groupId: source.groupId, paidByUserId: source.paidByUserId,
        status: source.status, categoryIds: source.categoryIds, tagIds: source.tagIds,
        comment: source.comment,
      }, "body", false, this.localOptions);
      const sub = registered.pipe(
        tap(task => {
          this.acceptTask(task);
          this.rememberPending(id, false);
          this.patchLocal(id, { taskId: task.id, registering: false });
        }),
        concatMap(task => {
          if (!task.canUpload) return of(task);
          this.patchLocal(id, { sending: true });
          return this.api.uploadRecognitionTaskFile(task.id, source.file, "events", true, this.localOptions)
            .pipe(tap(event => this.onUploadEvent(id, event)));
        }),
        catchError((error: HttpErrorResponse) => {
          this.patchLocal(id, {
            registering: false, sending: false,
            error: error.status === 403 ? text("当前没有上传权限。", "You no longer have upload permission.")
              : error.status === 400 ? text("文件无法提交，请检查文件和必填信息。", "Cannot submit this file. Check the file and required fields.")
                : text("提交状态暂未确认，正在同步服务器状态。请确认后重新上传。", "Submission is unconfirmed. Synchronizing server state before reuploading."),
          });
          return EMPTY;
        }),
        finalize(() => {
          this.transfers.delete(id);
          if (owner !== this.owner) return;
          this.activeTransfers--;
          this.patchLocal(id, { registering: false, sending: false });
          this.pump();
          this.refresh();
        }),
      ).subscribe();
      if (!sub.closed) this.transfers.set(id, sub);
    }
  }

  private onUploadEvent(id: string, event: HttpEvent<RecognitionTask>): void {
    if (event.type === HttpEventType.UploadProgress) {
      this.patchLocal(id, { loaded: event.loaded, total: event.total,
        awaitingConfirmation: event.total != null && event.loaded >= event.total });
    }
    if (event.type === HttpEventType.Response && event.body) {
      this.acceptTask(event.body);
      this.patchLocal(id, { sending: false, awaitingConfirmation: false, error: undefined });
      this.sources.delete(id);
    }
  }

  private acceptTask(task: RecognitionTask): void {
    if (this.state().tasks[task.id]?.version > task.version) return;
    this.store.dispatch(new UpsertQuickScanTask(task));
    this.patchLocal(task.clientRequestId, { taskId: task.id, reselectRequired: false, error: undefined });
    if (!["AWAITING_UPLOAD", "UPLOADING", "UPLOAD_INTERRUPTED"].includes(task.status)) {
      this.sources.delete(task.clientRequestId);
      this.patchLocal(task.clientRequestId, { awaitingConfirmation: false, error: undefined });
    }
  }

  private removeInaccessibleTasks(requested: number[], returned: number[]): void {
    const missing = new Set(requested.filter(id => !returned.includes(id)));
    if (!missing.size) return;
    const tasks = { ...this.state().tasks };
    const local = { ...this.state().local };
    missing.forEach(id => delete tasks[id]);
    Object.values(local).forEach(row => {
      if (row.taskId && missing.has(row.taskId)) {
        this.sources.delete(row.clientRequestId);
        delete local[row.clientRequestId];
      }
    });
    this.store.dispatch(new PatchQuickScanTasks({ tasks, local,
      pageIds: this.state().pageIds.filter(id => !missing.has(id)) }));
  }

  private patchLocal(id: string, patch: Partial<LocalQuickScanTask>): void {
    this.store.dispatch(new PatchLocalQuickScanTask(id, patch));
  }

  private pendingKey(): string { return "receipt-wrangler-pending-scans-" + this.owner; }
  private pendingIds(): string[] {
    try {
      const value: unknown = JSON.parse(localStorage.getItem(this.pendingKey()) || "[]");
      return Array.isArray(value) ? value.filter(id => typeof id === "string" && /^[a-f0-9-]{36}$/i.test(id)) : [];
    } catch { return []; }
  }
  private rememberPending(id: string, add: boolean): void {
    if (!this.owner) return;
    try {
      const ids = this.pendingIds().filter(current => current !== id);
      if (add) ids.push(id);
      if (ids.length) localStorage.setItem(this.pendingKey(), JSON.stringify(ids));
      else localStorage.removeItem(this.pendingKey());
    } catch { /* Tracking remains available from the server if storage is unavailable. */ }
  }

  private clearLegacySubmissionMetadata(): void {
    if (!this.owner) return;
    try { localStorage.removeItem("receipt-wrangler-quick-scan-submissions-" + this.owner); }
    catch { /* Only unresolved request IDs are retained for recovery. */ }
  }

  private requireReselect(id: string): void {
    this.rememberPending(id, false);
    this.patchPermissionBlockedRequest(id, false);
    const local = this.state().local[id];
    if (local) {
      this.patchLocal(id, {
        registering: false, sending: false, awaitingConfirmation: false,
        reselectRequired: true, error: undefined,
      });
      return;
    }
    const notices = new Set(this.state().reselectNoticeIds);
    notices.add(id);
    this.store.dispatch(new PatchQuickScanTasks({ reselectNoticeIds: [...notices] }));
  }

  private patchPermissionBlockedRequest(id: string, blocked: boolean): void {
    const ids = new Set(this.state().permissionBlockedRequestIds);
    if (blocked) ids.add(id); else ids.delete(id);
    this.store.dispatch(new PatchQuickScanTasks({ permissionBlockedRequestIds: [...ids] }));
  }

  private pruneUnauthorizedReselectRows(): void {
    const local = { ...this.state().local };
    const groupPermissions = this.store.selectSnapshot(AuthState.groupPermissions);
    Object.values(local).filter(row => row.reselectRequired).forEach(row => {
      const permissions = groupPermissions?.[row.groupId] ?? [];
      if (!permissions.includes(Permission.GroupReceiptsRead) || !permissions.includes(Permission.GroupReceiptsQuickScan)) {
        this.sources.delete(row.clientRequestId);
        delete local[row.clientRequestId];
      }
    });
    this.store.dispatch(new PatchQuickScanTasks({ local }));
  }
  private clearSession(): void {
    if (this.owner) {
      try { localStorage.removeItem(this.pendingKey()); } catch { /* No file data is persisted. */ }
      this.clearLegacySubmissionMetadata();
    }
    this.owner = "";
    clearTimeout(this.timer);
    this.poll?.unsubscribe();
    this.transfers.forEach(subscription => subscription.unsubscribe());
    this.transfers.clear();
    this.mutations.forEach(subscription => subscription.unsubscribe());
    this.mutations.clear();
    this.retrying.clear();
    this.inaccessiblePending.clear();
    this.sources.clear();
    this.queue = [];
    this.activeTransfers = 0;
    this.failures = 0;
    this.refreshPending = false;
    this.pollRunning = false;
    this.store.dispatch(new ClearQuickScanTasks());
  }

  private trackMutation(subscription: Subscription): void {
    if (subscription.closed) return;
    this.mutations.add(subscription);
    subscription.add(() => this.mutations.delete(subscription));
  }
}
