// SPDX-Licence-Identifier: EUPL-1.2

import { Component, CUSTOM_ELEMENTS_SCHEMA, Input, OnDestroy, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterLinkActive, RouterOutlet } from '@angular/router';
import { TranslationService } from '../services/translation.service';
import { ProviderDiscoveryService } from '../services/provider-discovery.service';
import { WebSocketService } from '../services/websocket.service';

interface NavItem {
  name: string;
  href: string;
  icon: string;
}

/**
 * ApplicationFrameComponent is the HLCRF (Header, Left nav, Content, Right, Footer)
 * shell for all Core Wails applications. It provides:
 *
 * - Dynamic sidebar navigation populated from ProviderDiscoveryService
 * - Content area rendered via router-outlet for child routes
 * - Footer with time, version, and provider status
 * - Mobile-responsive sidebar with expand/collapse
 * - Dark mode support
 *
 * Ported from core-gui/cmd/lthn-desktop/frontend/src/frame/application.frame.ts
 * with navigation made dynamic via provider discovery.
 */
@Component({
  selector: 'application-frame',
  standalone: true,
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  imports: [CommonModule, RouterOutlet, RouterLink, RouterLinkActive],
  templateUrl: './application-frame.component.html',
  styles: [
    `
      .application-frame {
        min-height: 100vh;
      }

      .frame-main {
        min-height: calc(100vh - 6.5rem);
      }

      .connection-dot {
        display: inline-block;
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: rgb(107 114 128);
        margin-right: 0.375rem;
        vertical-align: middle;
      }

      .connection-dot.connected {
        background: rgb(34 197 94);
        box-shadow: 0 0 4px rgb(34 197 94);
      }
    `,
  ],
})
export class ApplicationFrameComponent implements OnInit, OnDestroy {
  @Input() version = 'v0.1.0';

  sidebarOpen = false;
  userMenuOpen = false;
  time = '';
  private intervalId: ReturnType<typeof setInterval> | undefined;

  /** Static navigation items set by the host application. */
  @Input() staticNavigation: NavItem[] = [];

  /** Combined navigation: static + dynamic from providers. */
  navigation: NavItem[] = [];

  userNavigation: NavItem[] = [];

  constructor(
    public t: TranslationService,
    private providerService: ProviderDiscoveryService,
    private wsService: WebSocketService,
  ) {}

  /** Provider count from discovery service. */
  readonly providerCount = () => this.providerService.providers().length;

  /** WebSocket connection status. */
  readonly wsConnected = () => this.wsService.connected();

  async ngOnInit(): Promise<void> {
    this.updateTime();
    this.intervalId = setInterval(() => this.updateTime(), 1000);

    await this.t.onReady();
    this.initUserNavigation();

    // Discover providers and build navigation
    await this.providerService.discover();
    this.buildNavigation();

    // Connect WebSocket for real-time updates
    this.wsService.connect();
  }

  ngOnDestroy(): void {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }

  private initUserNavigation(): void {
    this.userNavigation = [
      {
        name: this.t._('menu.settings'),
        href: 'settings',
        icon: 'fa-regular fa-gear',
      },
    ];
  }

  private buildNavigation(): void {
    const dynamicItems = this.providerService
      .providers()
      .filter((p) => p.element)
      .map((p) => ({
        name: p.name,
        href: p.name.toLowerCase(),
        icon: 'fa-regular fa-puzzle-piece fa-2xl shrink-0',
      }));

    this.navigation = [...this.staticNavigation, ...dynamicItems];
  }

  private updateTime(): void {
    this.time = new Date().toLocaleTimeString();
  }
}
