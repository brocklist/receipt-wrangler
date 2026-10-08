# Quick Scan per-file upload and recognition progress

## Goal and scope

Add a durable, per-file task experience for desktop Quick Scan only. Keep Smart Fill, email ingestion, Flutter screens, task cancellation, manual queue scheduling, historical backfill, and deployment out of scope. Preserve the legacy multipart Quick Scan endpoint and response semantics for existing clients. Newly accepted requests to that endpoint also create recognition records after full file receipt/enqueue so they appear in task history, but legacy clients receive no synthetic upload percentage; do not backfill historical records. New desktop Quick Scan submissions use the per-file RecognitionTask API.

## User-visible behavior

- Add a top-bar “识别任务 / Recognition tasks” entry and /receipts/recognition-tasks page. The top bar shows the server-authorized active count.
- After submission, the Quick Scan dialog shows the current batch and each file's group, upload state/real byte percentage, recognition stage, elapsed time, retry timing, result and allowed action. Closing the dialog leaves uploads and polling alive in an application-scoped service.
- Upload files independently with a maximum of two active requests. Each file carries its own existing Quick Scan group, payer, status, categories and tags. A failed file does not block others.
- Show real upload bytes only. If total bytes are unknown, show bytes sent and an indeterminate loader. At 100% sent, show server confirmation until the PUT response says the file is durably accepted/enqueued.
- Render these phases in order: waiting for upload, uploading, waiting for enqueue, queued, preprocessing, OCR when used, AI extraction (show fallback-config indicator when set), parsing, saving, complete. Vision/multimodal processing skips OCR. Do not synthesize recognition percentages or completion estimates.
- Show elapsed time from server stage timestamps. Keep the last known task state during network errors; show a reconnect indicator without marking the task failed or producing repeated error toasts.
- On reload, query server tasks first. Never silently re-send a source file. For an interrupted/not-yet-received upload, ask the user to reselect the original and restart the full upload; byte-offset resume is out of scope.
- Completed tasks open their receipt. Failed tasks show a sanitized summary, attempt count and actual next retry time, plus only server-authorized retry actions.

## API and data contract

Add a RecognitionTask model/table independent of SystemTask so SystemTask retains its existing history semantics. A task stores client/user/group identity, file name/size, upload bytes, status and stage, server timestamps, attempt/max-attempt count, retry generation/version and next retry time, fallback flag, result receipt ID, and safe error code/message. Do not expose source path, upload token, provider prompt/response or credentials.

OpenAPI in api/swagger.yml is the source contract; generate desktop and mobile clients through api/generate-client.sh, never hand-edit generated clients.

- POST /api/recognitionTask: register one file before transfer. Body: clientRequestId UUID, fileName, fileSize, groupId, optional paidByUserId, status, categoryIds, tagIds. Same owner + request ID + same immutable submission fingerprint returns the original task; a mismatched fingerprint returns 409.
- PUT /api/recognitionTask/{id}/file: multipart one file. Record server-received bytes at most once per second, validate the completed file, then persist dispatch-pending state before enqueue. Return the task only after server acceptance; keep dispatch-pending distinct from queued while Redis confirmation is unresolved.
- GET /api/recognitionTask: paged list with scope (own default or all for authorized readers), bucket (active, history, all), task/client-request IDs, status and group filters. Return page data and authorization-scoped counts independent of the current page so the header remains accurate.
- GET /api/recognitionTask/{id}: authorized task refresh.
- POST /api/recognitionTask/{id}/retry: expected task version; atomically create a new retry generation and enqueue, rejecting stale/duplicate clicks with 409.
- Statuses: AWAITING_UPLOAD, UPLOADING, UPLOAD_INTERRUPTED, DISPATCH_PENDING, QUEUED, RUNNING, RETRY_WAIT, SUCCEEDED, FAILED. Recognition stages: UPLOAD, PREPROCESSING, OCR, AI, PARSING, SAVING, DONE. The client presents friendly Chinese/English labels without conflating status and stage.
- Preserve Asynq initial execution plus at most three automatic retries. Persist retry attempts and actual nextRetryAt.

## Authorization and privacy

- A normal user reads only their own task while their current group membership still grants group.receipts.read.
- A caller with app.system-tasks.read may request the all-users view across owners and groups, matching the existing system-wide activity log. This global read permission is the explicit visibility grant for that view and does not depend on group.receipts.read.
- Global read permission does not grant upload or retry permission. Upload is owner-only. Retry rechecks current group Quick Scan and activity-rerun permission plus field/category/tag grants.
- Return only sanitized safe error summaries. Do not include internal filesystem locations, raw provider errors containing sensitive data, prompts, tokens or source bytes in list responses/logs.

## Persistence, retries and recovery

- Store the single source file below the existing durable data directory using the repository's safe data-file helpers. Keep it on recognition failure for retry and remove it only after successful receipt persistence.
- Create the task before transfer. Use a separate upload attempt identity/temp path so a stale or concurrent upload cannot overwrite a newer accepted file. On interrupted transfer, immediately persist UPLOAD_INTERRUPTED; after six minutes reconcile abandoned UPLOADING rows to interrupted and invalidate that writer before allowing a complete re-upload.
- Persist DISPATCH_PENDING before publishing to Redis. Reconcile pending database/queue states on startup and every 30 seconds; recover enqueue failures and process crashes without duplicate receipts.
- Give every worker execution a fenced attempt identity. A stale execution cannot update current task state. Recognition success, receipt, image and task result commit in one database transaction; duplicate delivery cannot create a second receipt.
- Keep all attempts for one logical task, including primary-to-fallback AI configuration, under the same task ID. Record stage transitions using server time; visual extraction stages must not claim a numeric AI completion percentage.

## Desktop state and synchronization

- Reuse NGXS for task state. An application-scoped task service owns selected File objects, upload requests and state synchronization; closing the dialog or navigating away does not cancel them.
- Persist no credentials or raw file data in local storage. Keep only narrowly scoped unresolved client-request IDs if needed to reconcile an ambiguous registration response; clear them on user change/logout.
- If an own-scope request-ID lookup returns 403 because group read access was revoked, keep that ID unresolved and do not report the file missing or upload again. Reconcile after permissions change; only an accessible empty lookup becomes a reselect prompt.
- Poll every two seconds, at most one request at a time, while the task page is visible or the app has active tasks. On transport failure, back off at 2, 4, 8 then 15 seconds. Resume immediately on network recovery, window focus or task-page open.
- Abort state polling on user change/logout. Treat UPLOAD_INTERRUPTED as requiring user reselection and full-file upload, never as proof that a partial file is resumable.
- Use the app's existing Chinese/English language preference. Support narrow screens and accessible status/progress announcements.

## Acceptance and validation

- Backend: focused tests for idempotent registration/conflicting reuse, owner/group/admin visibility, revoked permissions, upload byte updates and interruption, dispatch recovery, enqueue/Redis outage, attempt fencing, retry version/generation, fallback, transaction rollback and duplicate delivery. Run Go build and the complete Go test suite.
- Desktop: Jest tests for XHR byte events vs final response, max-two uploads, independent file failures, configuration mapping, 100%-sent/server-confirmation distinction, NGXS/service lifecycle, polling non-overlap/backoff/reconnect, permission-aware list/retry/result navigation, bilingual strings and narrow layout. Run Jest and production build.
- UI acceptance: use a real browser against the API-backed app where supported, exercise the key upload → task progress → complete/failure/reconnect flow, and retain a screenshot. Record URL, API mode, role/data state, interactions, result and untested states. Keep existing Quick Scan, Smart Fill and email regression checks.
- Review the final diff explicitly. Once source is stable, perform one affected AOCI incremental maintenance check; report separately if the governance tool blocks that check.
- Do not claim an API-backed or live acceptance path from mocks, build, health checks or page load alone.


