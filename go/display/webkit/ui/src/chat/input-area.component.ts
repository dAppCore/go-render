import { CommonModule } from '@angular/common';
import { AfterViewInit, Component, CUSTOM_ELEMENTS_SCHEMA, ElementRef, EventEmitter, Input, Output, ViewChild } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ImageAttachment } from './chat.types';

@Component({
  selector: 'chat-input-area',
  standalone: true,
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  imports: [CommonModule, FormsModule],
  template: `
    <div class="composer" (dragover)="onDragOver($event)" (drop)="onDrop($event)">
      <input #filePicker type="file" accept=".png,.jpg,.jpeg,.webp,.gif" multiple hidden (change)="onFileSelection($event)" />
      <div class="composer__attachments" *ngIf="attachments.length">
        <figure *ngFor="let attachment of attachments; let index = index" class="attachment">
          <img [src]="attachmentSource(attachment)" [alt]="attachment.filename" />
          <figcaption>{{ attachment.filename }}</figcaption>
          <button type="button" (click)="removeAttachment.emit(index)">Remove</button>
        </figure>
      </div>
      <textarea
        #textarea
        [ngModel]="value"
        (ngModelChange)="onValueChange($event)"
        (keydown.enter)="submitOnEnter($event)"
        (paste)="onPaste($event)"
        placeholder="Ask the local model something useful"
      ></textarea>
      <div class="composer__meta">
        <span>{{ value.length }} chars · {{ attachments.length }} image(s)</span>
        <div class="composer__actions">
          <wa-button
            type="button"
            appearance="plain"
            [disabled]="!visionEnabled"
            [attr.title]="!visionEnabled ? visionDisabledReason : null"
            (click)="requestAttachment(filePicker)"
          >
            Attach
          </wa-button>
          <wa-button type="button" variant="brand" [disabled]="disabled" (click)="submit.emit()">Send</wa-button>
        </div>
      </div>
      <p *ngIf="!visionEnabled" class="composer__hint">{{ visionDisabledReason }}</p>
    </div>
  `,
  styles: [
    `
      .composer { display: grid; gap: 0.75rem; padding: 1rem; border-radius: 1.35rem; background: rgba(8, 21, 35, 0.86); border: 1px solid rgba(125, 211, 252, 0.12); }
      .composer__attachments { display: flex; gap: 0.75rem; overflow: auto; }
      .attachment { min-width: 8rem; margin: 0; display: grid; gap: 0.45rem; padding: 0.6rem; border-radius: 1rem; background: rgba(2, 6, 23, 0.76); }
      .attachment img { width: 100%; aspect-ratio: 1.3; object-fit: cover; border-radius: 0.8rem; }
      .attachment figcaption { color: #cbd5e1; font-size: 0.78rem; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
      textarea { width: 100%; min-height: 3.5rem; max-height: 10.5rem; resize: none; border: 0; background: transparent; color: #f8fafc; font: inherit; outline: 0; line-height: 1.6; }
      .composer__meta { display: flex; justify-content: space-between; align-items: center; gap: 1rem; color: #94a3b8; font-size: 0.82rem; }
      .composer__actions { display: flex; gap: 0.6rem; }
      .composer__hint { margin: 0; color: #fbbf24; font-size: 0.8rem; }
      wa-button[disabled] { opacity: 0.4; }
    `,
  ],
})
export class InputAreaComponent implements AfterViewInit {
  @ViewChild('textarea') private readonly textarea?: ElementRef<HTMLTextAreaElement>;

  @Input() value = '';
  @Input() disabled = false;
  @Input() attachments: ImageAttachment[] = [];
  @Input() visionEnabled = true;
  @Input() visionDisabledReason = '';
  @Input() nativePickerEnabled = false;
  @Output() valueChange = new EventEmitter<string>();
  @Output() attachFiles = new EventEmitter<FileList | File[]>();
  @Output() removeAttachment = new EventEmitter<number>();
  @Output() openNativePicker = new EventEmitter<void>();
  @Output() submit = new EventEmitter<void>();

  ngAfterViewInit(): void {
    this.resizeTextarea();
  }

  onValueChange(value: string): void {
    this.valueChange.emit(value);
    queueMicrotask(() => this.resizeTextarea());
  }

  submitOnEnter(event: Event): void {
    const keyboard = event as KeyboardEvent;
    if (!keyboard.shiftKey) {
      keyboard.preventDefault();
      this.submit.emit();
    }
  }

  onDragOver(event: DragEvent): void {
    event.preventDefault();
  }

  onDrop(event: DragEvent): void {
    event.preventDefault();
    if (!this.visionEnabled) {
      return;
    }
    if (event.dataTransfer?.files?.length) {
      this.attachFiles.emit(event.dataTransfer.files);
    }
  }

  onPaste(event: ClipboardEvent): void {
    if (!this.visionEnabled) {
      return;
    }
    const files = Array.from(event.clipboardData?.files ?? []);
    if (files.length > 0) {
      this.attachFiles.emit(files);
    }
  }

  onFileSelection(event: Event): void {
    if (!this.visionEnabled) {
      return;
    }
    const input = event.target as HTMLInputElement;
    if (input.files?.length) {
      this.attachFiles.emit(input.files);
      input.value = '';
    }
  }

  openFilePicker(input: HTMLInputElement): void {
    if (!this.visionEnabled) {
      return;
    }
    input.click();
  }

  requestAttachment(input: HTMLInputElement): void {
    if (!this.visionEnabled) {
      return;
    }
    if (this.nativePickerEnabled) {
      this.openNativePicker.emit();
      return;
    }
    this.openFilePicker(input);
  }

  attachmentSource(attachment: ImageAttachment): string {
    return `data:${attachment.mime_type};base64,${attachment.data}`;
  }

  private resizeTextarea(): void {
    const element = this.textarea?.nativeElement;
    if (!element) {
      return;
    }
    element.style.height = 'auto';
    element.style.height = `${Math.min(element.scrollHeight, 168)}px`;
  }
}
