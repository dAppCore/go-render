// SPDX-Licence-Identifier: EUPL-1.2

import '@awesome.me/webawesome';
import { platformBrowser } from '@angular/platform-browser';
import { AppModule } from './app/app-module';
import { registerChatHeroElement } from './web-components/chat-hero.element';

registerChatHeroElement();

platformBrowser()
  .bootstrapModule(AppModule, {
    ngZoneEventCoalescing: true,
  })
  .catch((err) => console.error(err));
