// SPDX-Licence-Identifier: EUPL-1.2

import { Component, Input, OnDestroy, OnInit, signal } from '@angular/core';
import { ProviderDiscoveryService } from '../services/provider-discovery.service';
import { WebSocketService } from '../services/websocket.service';

/**
 * StatusBarComponent renders the footer bar showing time, version,
 * provider count, and connection status.
 */
@Component({
  selector: 'status-bar',
  standalone: true,
  template: `
    <footer class="status-bar">
      <div class="status-left">
        <span class="status-item version">{{ version }}</span>
        <span class="status-item providers">
          <i class="fa-regular fa-puzzle-piece"></i>
          {{ providerCount() }} providers
        </span>
      </div>
      <div class="status-right">
        <span class="status-item connection" [class.connected]="wsConnected()">
          <span class="status-dot"></span>
          {{ wsConnected() ? 'Connected' : 'Disconnected' }}
        </span>
        <span class="status-item time">{{ time() }}</span>
      </div>
    </footer>
  `,
  styles: [
    `
      .status-bar {
        position: fixed;
        inset-inline: 0;
        bottom: 0;
        z-index: 40;
        height: 2.5rem;
        border-top: 1px solid rgb(229 231 235);
        background: #ffffff;
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding-inline: 1rem;
      }

      :host-context(.dark) .status-bar {
        border-color: rgba(255, 255, 255, 0.1);
        background: rgb(17 24 39);
      }

      .status-left,
      .status-right {
        display: flex;
        align-items: center;
        gap: 1rem;
      }

      .status-item {
        font-size: 0.875rem;
        color: rgb(107 114 128);
      }

      :host-context(.dark) .status-item {
        color: rgb(156 163 175);
      }

      .status-item i {
        margin-right: 0.25rem;
      }

      .status-dot {
        display: inline-block;
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: rgb(107 114 128);
        margin-right: 0.375rem;
      }

      .connection.connected .status-dot {
        background: rgb(34 197 94);
        box-shadow: 0 0 4px rgb(34 197 94);
      }

      .time {
        font-family: 'JetBrains Mono', 'Fira Code', monospace;
      }
    `,
  ],
})
export class StatusBarComponent implements OnInit, OnDestroy {
  @Input() version = 'v0.1.0';
  @Input() sidebarWidth = '5rem';

  readonly time = signal('');
  private intervalId: ReturnType<typeof setInterval> | undefined;

  constructor(
    private providerService: ProviderDiscoveryService,
    private wsService: WebSocketService,
  ) {}

  readonly providerCount = () => this.providerService.providers().length;
  readonly wsConnected = () => this.wsService.connected();

  ngOnInit(): void {
    this.updateTime();
    this.intervalId = setInterval(() => this.updateTime(), 1000);
  }

  ngOnDestroy(): void {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }

  private updateTime(): void {
    this.time.set(new Date().toLocaleTimeString());
  }
}
