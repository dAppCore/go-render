// SPDX-Licence-Identifier: EUPL-1.2

// Frame components
export { ApplicationFrameComponent } from './frame/application-frame.component';
export { SystemTrayFrameComponent } from './frame/system-tray-frame.component';

// Services
export { ApiConfigService } from './services/api-config.service';
export { ProviderDiscoveryService, type ProviderInfo, type ElementSpec } from './services/provider-discovery.service';
export { WebSocketService, type WSMessage } from './services/websocket.service';
export { TranslationService } from './services/translation.service';

// Components
export { ProviderHostComponent } from './components/provider-host.component';
export { ProviderNavComponent, type NavItem } from './components/provider-nav.component';
export { StatusBarComponent } from './components/status-bar.component';
export { ChatPanelComponent } from './chat/chat-panel.component';
