// SPDX-Licence-Identifier: EUPL-1.2

import { Injectable } from '@angular/core';
import { ApiConfigService } from './api-config.service';

/**
 * TranslationService provides a simple key-value translation lookup.
 * In production mode it fetches translations from the API; in development
 * it falls back to returning the key as-is.
 */
@Injectable({ providedIn: 'root' })
export class TranslationService {
  private translations = new Map<string, string>();
  private loaded = false;
  private loadingPromise: Promise<void>;

  constructor(private apiConfig: ApiConfigService) {
    this.loadingPromise = this.loadTranslations('en');
  }

  /** Reload translations for a given language. */
  reload(lang: string): Promise<void> {
    this.loaded = false;
    this.loadingPromise = this.loadTranslations(lang);
    return this.loadingPromise;
  }

  /** Translate a key. Returns the key itself if no translation is found. */
  translate(key: string): string {
    if (!this.loaded) {
      return key;
    }
    return this.translations.get(key) ?? key;
  }

  /** Shorthand for translate(). */
  _ = (key: string): string => this.translate(key);

  /** Wait for the initial translation load to complete. */
  onReady(): Promise<void> {
    return this.loadingPromise;
  }

  private async loadTranslations(lang: string): Promise<void> {
    try {
      const res = await fetch(this.apiConfig.url(`/api/v1/i18n/${lang}`));
      if (!res.ok) {
        // API not available — translations will fall through to keys
        this.loaded = true;
        return;
      }

      const messages: Record<string, string> = await res.json();
      this.translations.clear();
      for (const key in messages) {
        if (Object.prototype.hasOwnProperty.call(messages, key)) {
          this.translations.set(key, messages[key]);
        }
      }
      this.loaded = true;
    } catch {
      // Silently fall through — key passthrough is acceptable
      this.loaded = true;
    }
  }
}
