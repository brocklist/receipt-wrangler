import { HttpContextToken } from "@angular/common/http";

/** Background task synchronization presents its errors without repeated snackbars. */
export const HANDLE_ERROR_LOCALLY = new HttpContextToken<boolean>(() => false);
