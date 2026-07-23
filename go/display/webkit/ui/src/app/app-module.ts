// SPDX-Licence-Identifier: EUPL-1.2

import { DoBootstrap, Injector, NgModule, provideBrowserGlobalErrorListeners } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';
import { createCustomElement } from '@angular/elements';

import { App } from './app';
import { ChatPanelComponent } from '../chat/chat-panel.component';

@NgModule({
  imports: [BrowserModule, App, ChatPanelComponent],
  providers: [provideBrowserGlobalErrorListeners()],
})
export class AppModule implements DoBootstrap {
  constructor(private injector: Injector) {
    const el = createCustomElement(App, { injector });
    customElements.define('core-display', el);
    if (!customElements.get('core-chat-panel')) {
      customElements.define('core-chat-panel', createCustomElement(ChatPanelComponent, { injector }));
    }
  }

  ngDoBootstrap() {}
}
