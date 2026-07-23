// SPDX-Licence-Identifier: EUPL-1.2

import { Injectable } from '@angular/core';

/**
 * ApiConfigService provides a configurable base URL for all API calls.
 * Defaults to the current origin (Wails embedded) but can be overridden
 * for development or remote connections.
 */
@Injectable({ providedIn: 'root' })
export class ApiConfigService {
  private _baseUrl = '';

  /** The API base URL without a trailing slash. */
  get baseUrl(): string {
    return this._baseUrl;
  }

  /** Override the base URL. Strips trailing slash if present. */
  set baseUrl(url: string) {
    this._baseUrl = url.replace(/\/+$/, '');
  }

  /** Build a full URL for the given path. */
  url(path: string): string {
    const cleanPath = path.startsWith('/') ? path : `/${path}`;
    return `${this._baseUrl}${cleanPath}`;
  }
}
