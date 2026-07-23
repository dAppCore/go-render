// SPDX-Licence-Identifier: EUPL-1.2

import { Component, computed, Input, signal } from '@angular/core';
import { RouterLink, RouterLinkActive } from '@angular/router';
import { ProviderDiscoveryService, ProviderInfo } from '../services/provider-discovery.service';

export interface NavItem {
  name: string;
  href: string;
  icon: string;
  element?: { tag: string; source: string };
}

/**
 * ProviderNavComponent renders the sidebar navigation built dynamically
 * from the provider discovery service. Shows icon-only in collapsed mode,
 * expands on click.
 */
@Component({
  selector: 'provider-nav',
  standalone: true,
  imports: [RouterLink, RouterLinkActive],
  template: `
    <nav class="provider-nav">
      <ul role="list" class="nav-list">
        @for (item of navItems(); track item.name) {
          <li>
            <a
              [routerLink]="item.href"
              routerLinkActive="active"
              [routerLinkActiveOptions]="{ exact: true }"
              class="nav-item"
              [title]="item.name"
            >
              <i [class]="item.icon"></i>
              @if (expanded()) {
                <span class="nav-label">{{ item.name }}</span>
              }
            </a>
          </li>
        }
      </ul>
    </nav>
  `,
  styles: [
    `
      .provider-nav {
        display: flex;
        flex-direction: column;
        flex: 1;
      }

      .nav-list {
        list-style: none;
        margin: 0;
        padding: 0;
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 0.25rem;
      }

      .nav-list li {
        width: 100%;
      }

      .nav-item {
        display: flex;
        justify-content: center;
        align-items: center;
        gap: 0.75rem;
        padding: 1rem;
        border-radius: 0.375rem;
        color: #9ca3af;
        text-decoration: none;
        font-size: 0.875rem;
        font-weight: 600;
        line-height: 1.5;
        transition: background 0.15s, color 0.15s;
        height: 4rem;
      }

      .nav-item:hover {
        background: rgba(255, 255, 255, 0.05);
        color: #ffffff;
      }

      .nav-item.active {
        background: rgba(255, 255, 255, 0.05);
        color: #ffffff;
      }

      .nav-item i {
        flex-shrink: 0;
      }

      .nav-label {
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
      }
    `,
  ],
})
export class ProviderNavComponent {
  /** Additional static navigation items to prepend. */
  @Input() staticItems: NavItem[] = [];

  readonly expanded = signal(false);

  constructor(private providerService: ProviderDiscoveryService) {}

  /** Dynamic navigation built from discovered providers and static items. */
  readonly navItems = computed<NavItem[]>(() => {
    const dynamicItems = this.providerService
      .providers()
      .filter((p: ProviderInfo) => p.element)
      .map((p: ProviderInfo) => ({
        name: p.name,
        href: p.name.toLowerCase(),
        icon: 'fa-regular fa-puzzle-piece fa-2xl',
        element: p.element,
      }));

    return [...this.staticItems, ...dynamicItems];
  });

  /** Toggle sidebar expansion. */
  toggle(): void {
    this.expanded.update((v) => !v);
  }
}
