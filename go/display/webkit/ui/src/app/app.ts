import { Component } from '@angular/core';
import { ChatPanelComponent } from '../chat/chat-panel.component';

@Component({
  selector: 'core-display',
  templateUrl: './app.html',
  standalone: true,
  imports: [ChatPanelComponent],
})
export class App {
}
