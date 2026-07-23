import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ModelEntry } from './chat.types';

@Component({
  selector: 'chat-model-selector',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <label class="selector">
      <span>Model</span>
      <div class="selector__row">
        <select [ngModel]="value" (ngModelChange)="valueChange.emit($event)" [disabled]="loading">
          <option *ngFor="let model of models" [ngValue]="model.name">
            {{ model.loaded ? '● ' : '' }}{{ model.name }} · {{ model.architecture }} · {{ model.backend }}
          </option>
        </select>
        <span *ngIf="loading" class="spinner" aria-hidden="true"></span>
      </div>
    </label>
  `,
  styles: [
    `
      .selector { display: grid; gap: 0.35rem; color: #cbd5e1; font-size: 0.82rem; }
      .selector__row { display: flex; align-items: center; gap: 0.65rem; }
      select { min-width: 18rem; border-radius: 0.8rem; border: 1px solid rgba(124, 156, 191, 0.2); background: rgba(8, 21, 35, 0.8); color: #e2e8f0; padding: 0.72rem 0.9rem; }
      .spinner { width: 1rem; height: 1rem; border-radius: 999px; border: 2px solid rgba(148, 163, 184, 0.24); border-top-color: #f59e0b; animation: spin 0.8s linear infinite; }
      @keyframes spin { to { transform: rotate(360deg); } }
    `,
  ],
})
export class ModelSelectorComponent {
  @Input() models: ModelEntry[] = [];
  @Input() value = '';
  @Input() loading = false;
  @Output() valueChange = new EventEmitter<string>();
}
