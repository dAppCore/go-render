import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ChatSettings, ModelEntry } from './chat.types';

@Component({
  selector: 'chat-settings-panel',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <section class="panel" *ngIf="open">
      <header>
        <strong>Inference settings</strong>
        <div class="actions">
          <button type="button" (click)="reset.emit()">Reset to defaults</button>
          <button type="button" (click)="closed.emit()">Close</button>
        </div>
      </header>
      <label>
        Temperature
        <input type="range" min="0" max="2" step="0.1" [(ngModel)]="draft.temperature" />
        <span>{{ draft.temperature | number: '1.1-1' }}</span>
      </label>
      <label>
        Top P
        <input type="range" min="0" max="1" step="0.05" [(ngModel)]="draft.top_p" />
        <span>{{ draft.top_p | number: '1.2-2' }}</span>
      </label>
      <label>Top K <input type="number" min="0" max="200" [(ngModel)]="draft.top_k" /></label>
      <label>Max tokens <input type="number" min="64" max="32768" [(ngModel)]="draft.max_tokens" /></label>
      <label>
        Context window
        <select [(ngModel)]="draft.context_window">
          <option *ngFor="let option of contextWindows" [ngValue]="option">{{ option }}</option>
        </select>
      </label>
      <label>System prompt <textarea [(ngModel)]="draft.system_prompt"></textarea></label>
      <label>
        Default model
        <select [(ngModel)]="draft.default_model">
          <option *ngFor="let model of models" [ngValue]="model.name">{{ model.name }}</option>
        </select>
      </label>
      <button type="button" class="save" (click)="saved.emit(draft)">Save</button>
    </section>
  `,
  styles: [
    `
      .panel { display: grid; gap: 0.8rem; padding: 1rem; border-radius: 1.2rem; background: rgba(10, 18, 29, 0.95); border: 1px solid rgba(244, 114, 182, 0.18); }
      header, .actions { display: flex; justify-content: space-between; align-items: center; gap: 0.5rem; }
      header { color: #f9a8d4; }
      label { display: grid; gap: 0.35rem; color: #cbd5e1; font-size: 0.86rem; }
      input, textarea, select, button { border-radius: 0.8rem; border: 1px solid rgba(148, 163, 184, 0.2); background: rgba(15, 23, 42, 0.8); color: #f8fafc; padding: 0.72rem 0.85rem; }
      textarea { min-height: 7rem; resize: vertical; }
      .save { background: linear-gradient(135deg, #fb7185, #f59e0b); color: #111827; font-weight: 800; cursor: pointer; }
    `,
  ],
})
export class SettingsPanelComponent {
  readonly contextWindows = [2048, 4096, 8192, 16384, 32768];

  @Input() open = false;
  @Input() models: ModelEntry[] = [];
  @Input() set settings(value: ChatSettings) {
    this.draft = { ...value };
  }
  @Output() saved = new EventEmitter<ChatSettings>();
  @Output() reset = new EventEmitter<void>();
  @Output() closed = new EventEmitter<void>();

  draft: ChatSettings = {
    temperature: 1,
    top_p: 0.95,
    top_k: 64,
    max_tokens: 2048,
    context_window: 8192,
    system_prompt: '',
    default_model: '',
  };
}
