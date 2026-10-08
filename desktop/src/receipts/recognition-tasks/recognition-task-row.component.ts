import { CommonModule } from "@angular/common";
import { Component, computed, inject, input } from "@angular/core";
import { MatProgressBarModule } from "@angular/material/progress-bar";
import { ButtonModule } from "../../button";
import { localizeUiTextPair } from "../../i18n/language";
import { PipesModule } from "../../pipes";
import { RecognitionTaskRow } from "./quick-scan-task.models";
import { elapsedText, fileSizeText, taskStatusText, transferBytesText, uploadPercent } from "./quick-scan-task.presentation";
import { QuickScanTaskService } from "./quick-scan-task.service";

@Component({
  selector: "app-recognition-task-row",
  standalone: true,
  imports: [CommonModule, MatProgressBarModule, ButtonModule, PipesModule],
  templateUrl: "./recognition-task-row.component.html",
  styleUrl: "./recognition-task-row.component.scss",
})
export class RecognitionTaskRowComponent {
  public readonly row = input.required<RecognitionTaskRow>();
  public readonly tasks = inject(QuickScanTaskService);
  public readonly text = localizeUiTextPair;
  public readonly statusText = computed(() => taskStatusText(this.row()));
  public readonly percent = computed(() => uploadPercent(this.row()));
  public readonly uploadVisible = computed(() => this.row().local?.sending || this.row().task?.status === "UPLOADING");
  public readonly activeRecognition = computed(() =>
    ["RUNNING", "QUEUED", "DISPATCH_PENDING", "RETRY_WAIT"].includes(this.row().task?.status ?? ""));
  public readonly bytes = fileSizeText;
  public readonly transferBytes = transferBytesText;
  public readonly elapsed = elapsedText;
  public readonly groupId = computed(() => (this.row().task?.groupId ?? this.row().local?.groupId ?? "").toString());

  public chooseFile(event: Event): void {
    const input = event.target as HTMLInputElement;
    const file = input.files?.[0];
    const { task, local } = this.row();
    if (file && task) this.tasks.uploadAgain(task, file);
    else if (file && local) this.tasks.uploadLocalAgain(local, file);
    input.value = "";
  }
}
