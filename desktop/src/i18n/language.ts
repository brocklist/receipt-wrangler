import { computed, Injectable, signal } from "@angular/core";

export type AppLanguage = "zh-CN" | "en-US";

const LANGUAGE_STORAGE_KEY = "receipt-wrangler-language";

export function getAppLanguage(): AppLanguage {
  try {
    return globalThis.localStorage?.getItem(LANGUAGE_STORAGE_KEY) === "en-US"
      ? "en-US"
      : "zh-CN";
  } catch {
    return "zh-CN";
  }
}

/** Select an explicit translation for app-owned messages with interpolated user data. */
export function localizeUiTextPair(chinese: string, english: string): string {
  return getAppLanguage() === "zh-CN" ? chinese : english;
}

@Injectable({ providedIn: "root" })
export class LanguageService {
  public readonly language = signal<AppLanguage>(getAppLanguage());
  public readonly isChinese = computed(() => this.language() === "zh-CN");

  public toggle(): void {
    const nextLanguage: AppLanguage = this.isChinese() ? "en-US" : "zh-CN";
    try {
      globalThis.localStorage?.setItem(LANGUAGE_STORAGE_KEY, nextLanguage);
    } catch {
      return;
    }
    this.language.set(nextLanguage);
    if (typeof document !== "undefined") document.documentElement.lang = nextLanguage;
    if (typeof window !== "undefined") window.location.reload();
  }
}
