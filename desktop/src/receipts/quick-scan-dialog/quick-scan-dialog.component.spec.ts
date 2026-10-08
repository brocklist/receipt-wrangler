import { provideHttpClientTesting } from "@angular/common/http/testing";
import { CUSTOM_ELEMENTS_SCHEMA, signal } from "@angular/core";
import { ComponentFixture, TestBed, } from "@angular/core/testing";
import { FormArray, FormControl, ReactiveFormsModule } from "@angular/forms";
import { MatDialog, MatDialogModule, MatDialogRef } from "@angular/material/dialog";
import { MatSnackBarModule } from "@angular/material/snack-bar";
import { NoopAnimationsModule } from "@angular/platform-browser/animations";
import { ActivatedRoute } from "@angular/router";
import { NgxsModule, Store } from "@ngxs/store";
import { CarouselModule } from "ngx-bootstrap/carousel";
import { SharedUiModule } from "src/shared-ui/shared-ui.module";
import { LayoutState } from "src/store/layout.state";
import { ReceiptFileUploadCommand } from "../../interfaces";
import { ApiModule, ReceiptStatus } from "../../open-api";
import { PipesModule } from "../../pipes";
import { SnackbarService } from "../../services";
import { AuthState, GroupState } from "../../store";
import { QuickScanDialogComponent } from "./quick-scan-dialog.component";
import { provideHttpClient, withInterceptorsFromDi } from "@angular/common/http";
import { QuickScanTaskService } from "../recognition-tasks/quick-scan-task.service";
import { initialTaskState } from "../recognition-tasks/quick-scan-task.state";

describe("QuickScanDialogComponent", () => {
  let component: QuickScanDialogComponent;
  let fixture: ComponentFixture<QuickScanDialogComponent>;
  let store: Store;

  beforeEach(() => {
    localStorage.setItem("receipt-wrangler-language", "en-US");
    TestBed.configureTestingModule({
    declarations: [QuickScanDialogComponent],
    schemas: [CUSTOM_ELEMENTS_SCHEMA],
    imports: [ApiModule,
        CarouselModule,
        MatDialogModule,
        MatSnackBarModule,
        NgxsModule.forRoot([AuthState, GroupState, LayoutState]),
        NoopAnimationsModule,
        PipesModule,
        ReactiveFormsModule,
        SharedUiModule],
    providers: [
        { provide: QuickScanTaskService, useValue: {
          state: signal(initialTaskState()), submit: jest.fn().mockReturnValue(["batch-1"]),
        } },
        {
            provide: ActivatedRoute,
            useValue: {},
        },
        {
            provide: MatDialog,
            useValue: {}
        },
        {
            provide: MatDialogRef<QuickScanDialogComponent>,
            useValue: {
                close: () => { },
            },
        },
        provideHttpClient(withInterceptorsFromDi()),
        provideHttpClientTesting(),
    ]
});
    fixture = TestBed.createComponent(QuickScanDialogComponent);
    store = TestBed.inject(Store);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it("should create", () => {
    expect(component).toBeTruthy();
  });

  it("should init form correctly", () => {
    component.ngOnInit();

    expect(component.form.value).toEqual({
      paidByUserIds: [],
      statuses: [],
      groupIds: [],
      categories: [],
      tags: [],
    });
  });

  it("should push fileData into images when loaded", () => {
    const originalCreateObjectURL = URL.createObjectURL;
    URL.createObjectURL = jest.fn().mockReturnValue("awesome");
    const fileData = {} as ReceiptFileUploadCommand;
    component.fileLoaded(fileData);

    expect(component.images).toEqual([fileData]);
    URL.createObjectURL = originalCreateObjectURL;
  });

  it("should close the dialog", () => {
    const dialogSpy = jest.spyOn(TestBed.inject(MatDialogRef), "close");

    component.cancelButtonClicked();

    expect(dialogSpy).toHaveBeenCalled();
  });

  it("should show error if no image has been selected", () => {
    const snackbarSpy = jest.spyOn(TestBed.inject(SnackbarService), "error");

    component.submitButtonClicked();

    expect(snackbarSpy).toHaveBeenCalledWith(
      "Please select images to upload"
    );
  });

  it("should push new image when there are no user preferences", () => {
    component.fileLoaded({} as any);

    expect(component.form.value).toEqual({
      paidByUserIds: [""],
      statuses: [""],
      groupIds: [""],
      categories: [[]],
      tags: [[]],
    });
    expect(component.images).toEqual([{} as any]);
  });

  it("should push new image when there user preferences", () => {
    store.reset({
      auth: {
        userPreferences: {
          quickScanDefaultPaidById: 1,
          quickScanDefaultStatus: ReceiptStatus.Open,
          quickScanDefaultGroupId: 1,
        },
      },
      groups: { groups: [], selectedGroupId: "", selectedDashboardId: "" },
    });

    component.fileLoaded({} as any);

    expect(component.form.value).toEqual({
      paidByUserIds: [1],
      statuses: [ReceiptStatus.Open],
      groupIds: [1],
      categories: [[]],
      tags: [[]],
    });
    expect(component.images).toEqual([{} as any]);
  });

  it("should drive field visibility and required-ness from the selected group's settings", () => {
    const group = {
      id: 2,
      groupReceiptSettings: {
        quickScanPaidByEnabled: false,
        quickScanPaidByRequired: false,
        quickScanStatusEnabled: true,
        quickScanStatusRequired: true,
        quickScanCategoriesEnabled: true,
        quickScanCategoriesRequired: true,
        quickScanTagsEnabled: false,
        quickScanTagsRequired: false,
      },
    };
    store.reset({
      auth: {},
      groups: { groups: [group], selectedGroupId: "", selectedDashboardId: "" },
    });

    component.fileLoaded({} as any);
    component.groupIds.at(0).setValue(2);

    expect(component.showPaidBy(0)).toBe(false);
    expect(component.showStatus(0)).toBe(true);
    expect(component.showCategories(0)).toBe(true);
    expect(component.showTags(0)).toBe(false);

    // Required category with no selection is invalid; status is required and empty.
    expect(component.categories.at(0).valid).toBe(false);
    expect(component.statuses.at(0).valid).toBe(false);
    // Hidden paid-by is not required.
    expect(component.paidByUserIds.at(0).valid).toBe(true);
  });

  it("should send category and tag ids per image on submit", () => {
    const originalCreateObjectURL = URL.createObjectURL;
    try {
      URL.createObjectURL = jest.fn().mockReturnValue("blob");

      const group = {
        id: 2,
        groupReceiptSettings: {
          quickScanPaidByEnabled: false,
          quickScanPaidByRequired: false,
          quickScanStatusEnabled: true,
          quickScanStatusRequired: false,
          quickScanCategoriesEnabled: true,
          quickScanCategoriesRequired: true,
          quickScanTagsEnabled: true,
          quickScanTagsRequired: false,
        },
      };
      store.reset({
        auth: {
          groupCategories: { 2: [{ id: 10, name: "c" }] },
          groupTags: { 2: [{ id: 20, name: "t" }] },
        },
        groups: { groups: [group], selectedGroupId: "", selectedDashboardId: "" },
      });

      const serviceSpy = TestBed.inject(QuickScanTaskService).submit;
      const closeSpy = jest.spyOn(TestBed.inject(MatDialogRef), "close");

      const fileData = { file: { name: "a" } } as any;
      component.fileLoaded(fileData);
      component.groupIds.at(0).setValue(2);
      // Categories/tags are per-image FormArrays; selections are pushed onto them
      // (as the autocomplete does), not set wholesale.
      (component.categories.at(0) as FormArray).push(new FormControl({ id: 10, name: "c" }));
      (component.tags.at(0) as FormArray).push(new FormControl({ id: 20, name: "t" }));

      component.submitButtonClicked();

      expect(serviceSpy).toHaveBeenCalledWith([{
        file: fileData.file, groupId: 2, paidByUserId: undefined, status: undefined,
        categoryIds: [10], tagIds: [20],
      }]);
      expect(component.submitted()).toBe(true);
      fixture.detectChanges();
      expect(fixture.nativeElement.querySelector('[data-testid="quick-scan-current-batch"]')).toBeTruthy();
      expect(closeSpy).not.toHaveBeenCalled();
      component.submitButtonClicked();
      expect(serviceSpy).toHaveBeenCalledTimes(1);
    } finally {
      URL.createObjectURL = originalCreateObjectURL;
    }
  });
});
