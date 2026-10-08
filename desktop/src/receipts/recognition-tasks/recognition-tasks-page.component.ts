import { CommonModule } from "@angular/common";
import { Component, DestroyRef, TemplateRef, computed, inject, viewChild } from "@angular/core";
import { takeUntilDestroyed } from "@angular/core/rxjs-interop";
import { FormControl, ReactiveFormsModule } from "@angular/forms";
import { MatPaginatorModule, PageEvent } from "@angular/material/paginator";
import { MatTableDataSource } from "@angular/material/table";
import { Store } from "@ngxs/store";
import { ButtonModule } from "../../button";
import { LanguageService, localizeUiTextPair } from "../../i18n/language";
import { Permission, RecognitionTaskStatus } from "../../open-api";
import { SelectModule } from "../../select/select.module";
import { SharedUiModule } from "../../shared-ui/shared-ui.module";
import { TableModule } from "../../table/table.module";
import { TableColumn } from "../../table/table-column.interface";
import { AuthState, GroupState } from "../../store";
import { RecognitionTaskRow, TaskBucket, TaskScope } from "./quick-scan-task.models";
import { QuickScanTaskService } from "./quick-scan-task.service";
import { RecognitionTaskRowComponent } from "./recognition-task-row.component";

@Component({
  selector: "app-recognition-tasks-page",
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, MatPaginatorModule, ButtonModule, SelectModule,
    SharedUiModule, TableModule, RecognitionTaskRowComponent],
  templateUrl: "./recognition-tasks-page.component.html",
  styleUrl: "./recognition-tasks-page.component.scss",
})
export class RecognitionTasksPageComponent {
  public readonly tasks = inject(QuickScanTaskService);
  public readonly language = inject(LanguageService);
  private readonly store = inject(Store);
  private readonly destroyRef = inject(DestroyRef);
  public readonly text = localizeUiTextPair;
  private readonly rowTemplate = viewChild<TemplateRef<unknown>>("taskRow");
  private readonly groups = this.store.selectSignal(GroupState.groupsWithoutAll);
  public readonly canReadAll = this.store.selectSignal(AuthState.hasAppPermission(Permission.AppSystemTasksRead));
  public readonly groupFilter = new FormControl<number | undefined>(this.tasks.state().filter.groupId);
  public readonly statusFilter = new FormControl<RecognitionTaskStatus | undefined>(this.tasks.state().filter.status);
  public readonly groupOptions = computed(() => [
    { value: undefined, displayValue: this.text("所有分组", "All groups") },
    ...this.groups().map(group => ({ value: group.id, displayValue: group.name })),
  ]);
  public readonly statusOptions = [
    { value: undefined, displayValue: this.text("所有状态", "All statuses") },
    { value: "AWAITING_UPLOAD", displayValue: this.text("等待上传", "Waiting for upload") },
    { value: "UPLOADING", displayValue: this.text("正在上传", "Uploading") },
    { value: "UPLOAD_INTERRUPTED", displayValue: this.text("上传已中断", "Upload interrupted") },
    { value: "DISPATCH_PENDING", displayValue: this.text("等待入队", "Awaiting queue") },
    { value: "QUEUED", displayValue: this.text("排队中", "Queued") },
    { value: "RUNNING", displayValue: this.text("识别中", "Recognizing") },
    { value: "RETRY_WAIT", displayValue: this.text("等待重试", "Waiting for retry") },
    { value: "SUCCEEDED", displayValue: this.text("完成", "Completed") },
    { value: "FAILED", displayValue: this.text("失败", "Failed") },
  ];
  public readonly tabs = computed(() => [
    { value: "active", label: this.text("进行中", "Active"), count: this.tasks.state().activeCount },
    { value: "history", label: this.text("历史", "History") },
    { value: "all", label: this.text("全部", "All") },
  ]);
  public readonly columns = computed<TableColumn[]>(() => this.rowTemplate() ? [
    { columnHeader: this.text("文件与处理进度", "Files and processing progress"), matColumnDef: "task", template: this.rowTemplate(), sortable: false },
  ] : []);
  public readonly dataSource = computed(() => new MatTableDataSource<RecognitionTaskRow>(
    this.tasks.state().pageIds.map(id => ({
      id: id.toString(), task: this.tasks.state().tasks[id],
      local: Object.values(this.tasks.state().local).find(local => local.taskId === id),
    })),
  ));
  public readonly reselectRows = computed(() => Object.values(this.tasks.state().local)
    .filter(local => local.reselectRequired)
    .map(local => ({ id: local.clientRequestId, local })));

  constructor() {
    this.tasks.openPage();
    this.destroyRef.onDestroy(() => this.tasks.closePage());
    this.groupFilter.valueChanges.pipe(takeUntilDestroyed()).subscribe(value => this.tasks.setFilter({ groupId: value ?? undefined, page: 1 }));
    this.statusFilter.valueChanges.pipe(takeUntilDestroyed()).subscribe(value => this.tasks.setFilter({ status: value ?? undefined, page: 1 }));
  }

  public scopeChanged(scope: TaskScope): void { this.tasks.setFilter({ scope, page: 1 }); }
  public bucketChanged(bucket: string): void { this.tasks.setFilter({ bucket: bucket as TaskBucket, page: 1 }); }
  public pageChanged(event: PageEvent): void { this.tasks.setFilter({ page: event.pageIndex + 1, pageSize: event.pageSize }); }
}
