// SPDX-Licence-Identifier: EUPL-1.2

import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ProviderDiscoveryService, ProviderInfo } from '../services/provider-discovery.service';
import { WebSocketService } from '../services/websocket.service';

/**
 * SystemTrayFrameComponent is a 380x480 frameless panel showing:
 * - Provider status cards from the discovery service
 * - Brain connection status
 * - MCP server status
 * - Quick actions
 *
 * Ported from core-gui/cmd/lthn-desktop/frontend/src/frame/system-tray.frame.ts
 * with dynamic provider status cards.
 */
@Component({
  selector: 'system-tray-frame',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="tray-container">
      <!-- Header -->
      <div class="tray-header">
        <div class="tray-logo">
          <i class="fa-regular fa-cube logo-icon"></i>
          <span>Core IDE</span>
        </div>
        <div class="tray-controls">
          <button class="control-btn" (click)="settingsMenuOpen = !settingsMenuOpen" title="Settings">
            <i class="fa-regular fa-gear"></i>
          </button>
        </div>
      </div>

      <!-- Settings dropdown -->
      @if (settingsMenuOpen) {
        <div class="settings-menu">
          @for (item of settingsNavigation; track item.name) {
            <button class="settings-item" (click)="settingsMenuOpen = false">
              {{ item.name }}
            </button>
          }
        </div>
      }

      <!-- Connection Status -->
      <div class="status-section">
        <div class="status-row">
          <span class="status-label">Connection</span>
          <span class="status-value" [class.active]="wsConnected()">
            {{ wsConnected() ? 'Connected' : 'Disconnected' }}
          </span>
        </div>
        <div class="status-row">
          <span class="status-label">Providers</span>
          <span class="status-value">{{ providers().length }}</span>
        </div>
        <div class="status-row">
          <span class="status-label">Time</span>
          <span class="status-value mono">{{ time() }}</span>
        </div>
      </div>

      <!-- Provider Cards -->
      <div class="providers-section">
        <div class="section-header">Providers</div>
        <div class="providers-list">
          @for (provider of providers(); track provider.name) {
            <div class="provider-card">
              <div class="provider-icon">
                <i class="fa-regular fa-puzzle-piece"></i>
              </div>
              <div class="provider-info">
                <span class="provider-name">{{ provider.name }}</span>
                <span class="provider-path">{{ provider.basePath }}</span>
              </div>
              <div class="provider-status">
                <span class="status-indicator active"></span>
              </div>
            </div>
          } @empty {
            <div class="no-providers">No providers registered</div>
          }
        </div>
      </div>

      <!-- Footer -->
      <div class="tray-footer">
        <div class="connection-status" [class.connected]="wsConnected()">
          <div class="footer-dot"></div>
          <span>{{ wsConnected() ? 'Services Running' : 'Ready' }}</span>
        </div>
      </div>
    </div>
  `,
  styles: [
    `
      :host {
        display: block;
        width: 100%;
        height: 100%;
        overflow: hidden;
      }

      .tray-container {
        display: flex;
        flex-direction: column;
        height: 100%;
        background: rgb(17 24 39);
        font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
        color: rgb(156 163 175);
        border-radius: 0.375rem;
      }

      .tray-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0.75rem 1rem;
        background: rgb(31 41 55);
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
      }

      .tray-logo {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.9375rem;
        font-weight: 600;
        color: #ffffff;
      }

      .logo-icon {
        color: rgb(129 140 248);
      }

      .tray-controls {
        display: flex;
        gap: 0.5rem;
      }

      .control-btn {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 32px;
        background: transparent;
        border: 1px solid rgba(255, 255, 255, 0.1);
        border-radius: 6px;
        color: rgb(156 163 175);
        cursor: pointer;
        transition: all 0.15s ease;
      }

      .control-btn:hover {
        background: rgba(255, 255, 255, 0.05);
        border-color: rgba(255, 255, 255, 0.2);
        color: #ffffff;
      }

      .settings-menu {
        background: rgb(31 41 55);
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
        padding: 0.5rem;
      }

      .settings-item {
        display: block;
        width: 100%;
        text-align: left;
        padding: 0.5rem 0.75rem;
        background: transparent;
        border: none;
        border-radius: 0.25rem;
        color: rgb(156 163 175);
        font-size: 0.875rem;
        cursor: pointer;
      }

      .settings-item:hover {
        background: rgba(255, 255, 255, 0.05);
        color: #ffffff;
      }

      .status-section {
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        padding: 0.875rem 1rem;
        background: rgb(31 41 55);
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
      }

      .status-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
      }

      .status-label {
        font-size: 0.8125rem;
        color: rgb(107 114 128);
      }

      .status-value {
        font-size: 0.875rem;
        font-weight: 600;
        color: #ffffff;
      }

      .status-value.active {
        color: rgb(34 197 94);
      }

      .status-value.mono {
        font-family: 'JetBrains Mono', 'Fira Code', monospace;
        font-size: 0.8125rem;
      }

      .providers-section {
        display: flex;
        flex-direction: column;
        flex: 1;
        min-height: 0;
      }

      .section-header {
        padding: 0.625rem 1rem;
        font-size: 0.6875rem;
        font-weight: 600;
        color: rgb(107 114 128);
        text-transform: uppercase;
        letter-spacing: 0.05em;
        border-bottom: 1px solid rgba(255, 255, 255, 0.1);
      }

      .providers-list {
        flex: 1;
        overflow-y: auto;
        padding: 0.5rem;
      }

      .provider-card {
        display: flex;
        align-items: center;
        gap: 0.75rem;
        padding: 0.625rem 0.75rem;
        border-radius: 6px;
        transition: background 0.15s ease;
      }

      .provider-card:hover {
        background: rgba(255, 255, 255, 0.05);
      }

      .provider-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 32px;
        background: rgba(129, 140, 248, 0.15);
        border-radius: 6px;
        color: rgb(129 140 248);
      }

      .provider-info {
        display: flex;
        flex-direction: column;
        flex: 1;
        min-width: 0;
      }

      .provider-name {
        font-size: 0.8125rem;
        font-weight: 500;
        color: #ffffff;
      }

      .provider-path {
        font-size: 0.6875rem;
        color: rgb(107 114 128);
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      .provider-status {
        display: flex;
        align-items: center;
      }

      .status-indicator {
        width: 8px;
        height: 8px;
        border-radius: 50%;
        background: rgb(107 114 128);
      }

      .status-indicator.active {
        background: rgb(34 197 94);
        box-shadow: 0 0 4px rgb(34 197 94);
      }

      .no-providers {
        display: flex;
        align-items: center;
        justify-content: center;
        padding: 2rem;
        color: rgb(107 114 128);
        font-size: 0.8125rem;
      }

      .tray-footer {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0.625rem 1rem;
        background: rgb(31 41 55);
        border-top: 1px solid rgba(255, 255, 255, 0.1);
      }

      .connection-status {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        font-size: 0.75rem;
        color: rgb(107 114 128);
      }

      .footer-dot {
        width: 6px;
        height: 6px;
        border-radius: 50%;
        background: rgb(107 114 128);
      }

      .connection-status.connected .footer-dot {
        background: rgb(34 197 94);
        box-shadow: 0 0 4px rgb(34 197 94);
      }

      .connection-status.connected {
        color: rgb(34 197 94);
      }
    `,
  ],
})
export class SystemTrayFrameComponent implements OnInit, OnDestroy {
  settingsMenuOpen = false;
  readonly time = signal('');
  private intervalId: ReturnType<typeof setInterval> | undefined;

  settingsNavigation = [
    { name: 'Settings', href: '#' },
    { name: 'About', href: '#' },
    { name: 'Check for Updates...', href: '#' },
  ];

  constructor(
    private providerService: ProviderDiscoveryService,
    private wsService: WebSocketService,
  ) {}

  readonly providers = () => this.providerService.providers();
  readonly wsConnected = () => this.wsService.connected();

  async ngOnInit(): Promise<void> {
    this.updateTime();
    this.intervalId = setInterval(() => this.updateTime(), 1000);

    // Discover providers for status cards
    await this.providerService.discover();
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
