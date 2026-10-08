import { expect, test } from '@playwright/test';
import type { RecognitionTask } from '../src/open-api';
import { stubTokenRefresh } from './helpers/auth';
import { apiCreateGroup, apiDeleteGroupById, uniqueName, withAdminApi } from './helpers/provisioning';
import { injectQuickScanAppData, mockQuickScanTasks, openQuickScanDialog, uploadQuickScanImages } from './helpers/quick-scan';

test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Recognition tasks', () => {
  let group: { id: number; name: string };
  test.beforeAll(async () => {
    await withAdminApi(async api => {
      group = await apiCreateGroup(api, uniqueName('recognition'));
      const response = await api.put(`/api/group/${group.id}/groupReceiptSettings`, { data: {
        quickScanPaidByEnabled: false, quickScanPaidByRequired: false, quickScanDefaultPaidByType: 'UPLOADER',
        quickScanStatusEnabled: false, quickScanStatusRequired: false, quickScanDefaultStatus: 'OPEN',
        quickScanCategoriesEnabled: false, quickScanCategoriesRequired: false, quickScanTagsEnabled: false, quickScanTagsRequired: false,
      } });
      expect(response.ok()).toBe(true);
    });
  });
  test.afterAll(async () => {
    if (group) await withAdminApi(api => apiDeleteGroupById(api, String(group.id)));
  });
  test.beforeEach(async ({ page }, testInfo) => {
    const language = testInfo.title.includes('in Chinese') ? 'zh-CN' : 'en-US';
    await page.addInitScript(value => localStorage.setItem('receipt-wrangler-language', value), language);
    await stubTokenRefresh(page);
    await injectQuickScanAppData(page, { userPreferences: { quickScanDefaultGroupId: group.id } });
  });

  test('accepts a real file, keeps its task after closing/reload and never replays the upload', async ({ page }, testInfo) => {
    let mutations = 0;
    page.on('request', request => {
      if (/\/api\/recognitionTask(?:\/\d+\/file)?$/.test(new URL(request.url()).pathname) && ['POST', 'PUT'].includes(request.method())) mutations++;
    });
    const dialog = await openQuickScanDialog(page, group.id);
    await uploadQuickScanImages(dialog);
    const registration = page.waitForResponse(response => response.request().method() === 'POST' && new URL(response.url()).pathname === '/api/recognitionTask');
    const upload = page.waitForResponse(response => response.request().method() === 'PUT' && /\/recognitionTask\/\d+\/file$/.test(new URL(response.url()).pathname));
    await dialog.getByTestId('dialog-submit-button').click();
    const registered = await registration;
    expect(registered.status()).toBe(201);
    const task = await registered.json() as RecognitionTask;
    const accepted = await upload;
    expect(accepted.status()).toBe(202);
    const acceptedTask = await accepted.json() as RecognitionTask;
    expect(acceptedTask.uploadTotalBytes).toBe(acceptedTask.uploadedBytes);
    expect(['DISPATCH_PENDING', 'QUEUED', 'RUNNING', 'RETRY_WAIT', 'FAILED', 'SUCCEEDED']).toContain(acceptedTask.status);
    await expect(dialog.getByTestId('quick-scan-current-batch')).toBeVisible();
    await expect(dialog.getByText('Upload complete', { exact: false })).toBeVisible();
    await dialog.getByRole('button', { name: 'View recognition tasks' }).click();
    await expect(page).toHaveURL(/\/receipts\/recognition-tasks$/);
    // Terminal failure is legitimate when this test API has no OCR/AI provider. It remains durable history.
    await page.getByRole('tab', { name: /^(All|全部)$/ }).click();
    const row = page.getByTestId(`recognition-task-${task.id}`);
    await expect(row).toContainText('receipt.png');
    await page.reload();
    await expect(page.getByTestId('recognition-tasks-entry')).toBeVisible();
    await page.getByRole('tab', { name: /^(All|全部)$/ }).click();
    await expect(page.getByTestId(`recognition-task-${task.id}`)).toContainText('receipt.png');
    expect(mutations).toBe(2);
    const snapshot = testInfo.outputPath('recognition-api-reload.png');
    await page.screenshot({ path: snapshot, fullPage: true });
    await testInfo.attach('API-backed task after reload', { path: snapshot, contentType: 'image/png' });
  });

  test('renders stage, retry and interrupted-upload states responsively in Chinese with mocked recognition', async ({ page }, testInfo) => {
    const createdAt = new Date(Date.now() - 65000).toISOString();
    const task = (id: number, patch: Partial<RecognitionTask>): RecognitionTask => ({
      id, clientRequestId: `11111111-1111-4111-8111-${id.toString().padStart(12, '0')}`, fileName: `invoice-${id}.png`,
      fileSize: 1024, groupId: group.id, ownerUserId: 1, version: 1, status: 'RUNNING', stage: 'AI', createdAt,
      updatedAt: createdAt, queuedAt: createdAt, startedAt: createdAt, stageStartedAt: createdAt,
      uploadedBytes: 1024, uploadTotalBytes: 1024, attempt: 1, maxAttempts: 4, fallbackActive: false,
      errorCode: '', errorMessage: '', canRetry: false, canUpload: false, ...patch,
    });
    await mockQuickScanTasks(page, 1, [
      task(1001, { status: 'RUNNING', stage: 'AI', fallbackActive: true }),
      task(1002, { status: 'RETRY_WAIT', stage: 'OCR', attempt: 2, nextRetryAt: new Date(Date.now() + 30000).toISOString() }),
      task(1003, { status: 'UPLOAD_INTERRUPTED', stage: 'UPLOAD', uploadedBytes: 256, canUpload: true }),
      task(1004, { status: 'FAILED', stage: 'AI', errorMessage: 'Recognition failed. Check processing settings.', canRetry: true }),
      task(1005, { status: 'SUCCEEDED', stage: 'DONE', receiptId: 123, completedAt: new Date().toISOString() }),
    ]);
    await page.goto('/receipts/recognition-tasks');
    await expect(page.getByText('识别任务', { exact: true })).toBeVisible();
    await page.getByRole('tab', { name: /^(All|全部)$/ }).click();
    await expect(page.getByTestId('recognition-task-1001')).toContainText('正在提取发票信息');
    await expect(page.getByTestId('recognition-task-1001')).not.toContainText('%');
    await expect(page.getByTestId('recognition-task-1002')).toContainText('下次重试');
    await expect(page.getByTestId('recognition-task-1003')).toContainText('25%');
    await expect(page.getByTestId('recognition-task-1003').getByTestId('recognition-reupload')).toBeVisible();
    await expect(page.getByTestId('recognition-task-1004').getByTestId('recognition-retry')).toBeVisible();
    await expect(page.getByTestId('recognition-task-1005').getByTestId('recognition-open-receipt')).toBeVisible();
    await page.setViewportSize({ width: 390, height: 844 });
    expect(await page.evaluate(() => document.documentElement.scrollWidth)).toBeLessThanOrEqual(390);
    await page.getByTestId('recognition-task-1003').scrollIntoViewIfNeeded();
    const snapshot = testInfo.outputPath('recognition-chinese-mobile.png');
    await page.screenshot({ path: snapshot, fullPage: true });
    await testInfo.attach('Chinese mobile task states (mocked recognition)', { path: snapshot, contentType: 'image/png' });
  });
});
