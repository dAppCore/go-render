// SPDX-Licence-Identifier: EUPL-1.2

import { Injectable, signal } from '@angular/core';
import { ApiConfigService } from './api-config.service';

/**
 * Describes the element specification for a renderable provider.
 */
export interface ElementSpec {
  tag: string;
  source: string;
}

/**
 * Describes a provider as returned by GET /api/v1/providers.
 */
export interface ProviderInfo {
  name: string;
  basePath: string;
  status?: string;
  element?: ElementSpec;
  channels?: string[];
}

/**
 * ProviderDiscoveryService fetches the list of registered providers from
 * the API server and dynamically loads custom element scripts for any
 * Renderable providers.
 */
@Injectable({ providedIn: 'root' })
export class ProviderDiscoveryService {
  private readonly _providers = signal<ProviderInfo[]>([]);
  readonly providers = this._providers.asReadonly();

  private discovered = false;

  constructor(private apiConfig: ApiConfigService) {}

  /** Fetch providers from the API and load custom element scripts. */
  async discover(): Promise<void> {
    if (this.discovered) {
      return;
    }

    try {
      const res = await fetch(this.apiConfig.url('/api/v1/providers'));
      if (!res.ok) {
        console.warn('ProviderDiscoveryService: failed to fetch providers:', res.statusText);
        return;
      }

      const data = await res.json();
      const providers: ProviderInfo[] = data.providers ?? [];
      this._providers.set(providers);
      this.discovered = true;

      // Load custom elements for Renderable providers
      for (const p of providers) {
        if (p.element?.tag && p.element?.source) {
          await this.loadElement(p.element.tag, p.element.source);
        }
      }
    } catch (err) {
      console.warn('ProviderDiscoveryService: discovery failed:', err);
    }
  }

  /** Refresh the provider list (force re-discovery). */
  async refresh(): Promise<void> {
    this.discovered = false;
    await this.discover();
  }

  /** Dynamically load a custom element script if not already registered. */
  private async loadElement(tag: string, source: string): Promise<void> {
    if (customElements.get(tag)) {
      return;
    }

    const script = document.createElement('script');
    script.type = 'module';
    script.src = source.startsWith('http') ? source : this.apiConfig.url(source);
    document.head.appendChild(script);

    try {
      await customElements.whenDefined(tag);
    } catch (err) {
      console.warn(`ProviderDiscoveryService: failed to load element <${tag}>:`, err);
    }
  }
}
