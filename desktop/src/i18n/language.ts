import { computed, Injectable, signal } from "@angular/core";
import chineseMessages from "./zh-CN.json";
import englishOverrides from "./english-overrides.json";

export type AppLanguage = "zh-CN" | "en-US";

const LANGUAGE_STORAGE_KEY = "receipt-wrangler-language";
const chineseToEnglish: Record<string, string> = {};
const chineseMessagesByText: Record<string, string> = chineseMessages;

for (const [english, chinese] of Object.entries(chineseMessages)) {
  chineseToEnglish[chinese] ??= english;
}

Object.assign(chineseToEnglish, englishOverrides.defaults);

export function getAppLanguage(): AppLanguage {
  try {
    return globalThis.localStorage?.getItem(LANGUAGE_STORAGE_KEY) === "en-US"
      ? "en-US"
      : "zh-CN";
  } catch {
    return "zh-CN";
  }
}

export function getAngularLocale(): AppLanguage {
  return getAppLanguage();
}

export function isChineseLanguage(): boolean {
  return getAppLanguage() === "zh-CN";
}

/** Translate known application-owned text; unknown and user-provided values stay intact. */
export function localizeUiText(value: string | null | undefined): string {
  if (value == null || value === "") return value ?? "";

  const leading = value.match(/^\s*/)?.[0] ?? "";
  const trailing = value.match(/\s*$/)?.[0] ?? "";
  const key = value.trim();
  if (!key) return value;

  const translated = getAppLanguage() === "zh-CN"
  ? chineseMessagesByText[key] ?? key
    : chineseToEnglish[key] ?? key;

  return `${leading}${translated}${trailing}`;
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
