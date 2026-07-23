import { CommonModule, DatePipe } from '@angular/common';
import {
  AfterViewChecked,
  Component,
  CUSTOM_ELEMENTS_SCHEMA,
  ElementRef,
  Input,
  OnChanges,
  SimpleChanges,
  ViewChild,
} from '@angular/core';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import hljs from 'highlight.js/lib/core';
import bash from 'highlight.js/lib/languages/bash';
import go from 'highlight.js/lib/languages/go';
import javascript from 'highlight.js/lib/languages/javascript';
import json from 'highlight.js/lib/languages/json';
import markdown from 'highlight.js/lib/languages/markdown';
import plaintext from 'highlight.js/lib/languages/plaintext';
import python from 'highlight.js/lib/languages/python';
import typescript from 'highlight.js/lib/languages/typescript';
import xml from 'highlight.js/lib/languages/xml';
import { marked } from 'marked';
import { ChatMessage, ImageAttachment } from './chat.types';

hljs.registerLanguage('bash', bash);
hljs.registerLanguage('go', go);
hljs.registerLanguage('javascript', javascript);
hljs.registerLanguage('json', json);
hljs.registerLanguage('markdown', markdown);
hljs.registerLanguage('plaintext', plaintext);
hljs.registerLanguage('python', python);
hljs.registerLanguage('typescript', typescript);
hljs.registerLanguage('xml', xml);
hljs.registerLanguage('html', xml);

type RenderBlock =
  | { kind: 'text'; html: SafeHtml }
  | { kind: 'code'; code: string; language: string; html: SafeHtml };

@Component({
  selector: 'chat-message-list',
  standalone: true,
  schemas: [CUSTOM_ELEMENTS_SCHEMA],
  imports: [CommonModule, DatePipe],
  template: `
    <div class="thread-shell">
      <div #viewport class="thread" (scroll)="onScroll()">
        <section *ngIf="thinkingActive" class="thinking-banner">
          <span class="thinking-banner__pulse"></span>
          <span>Thinking...</span>
        </section>

        <article
          *ngFor="let message of messages; trackBy: trackByMessage"
          class="bubble"
          [class.bubble--user]="message.role === 'user'"
        >
          <header>
            <strong>{{ message.role }}</strong>
            <span>{{ message.model || 'local' }}</span>
          </header>

          <ng-container *ngFor="let block of renderBlocks(message.content)">
            <section *ngIf="block.kind === 'text'" class="content markdown" [innerHTML]="block.html"></section>
            <section *ngIf="block.kind === 'code'" class="code-block">
              <div class="code-block__bar">
                <span>{{ block.language || 'code' }}</span>
                <button type="button" class="code-block__copy" (click)="copyCode(block.code)">Copy</button>
              </div>
              <pre><code [innerHTML]="block.html"></code></pre>
            </section>
          </ng-container>

          <div *ngIf="streaming && isLatestAssistantMessage(message)" class="stream-cursor" aria-hidden="true"></div>

          <section *ngIf="message.attachments?.length" class="attachments">
            <figure *ngFor="let attachment of message.attachments" class="attachment">
              <button type="button" class="attachment__button" (click)="openLightbox(attachment)">
                <img [src]="attachmentSource(attachment)" [alt]="attachment.filename" />
              </button>
              <figcaption>{{ attachment.filename }} · {{ attachment.width }}×{{ attachment.height }}</figcaption>
            </figure>
          </section>

          <section *ngIf="message.thinking" class="thinking">
            <div class="thinking__badge" [class.thinking__badge--active]="message.thinking.active">
              Thinking...
            </div>
            <details>
              <summary>Thought for {{ thinkingDuration(message) }}</summary>
              <div class="thinking__content markdown" [innerHTML]="renderMarkdown(message.thinking.content || 'Thinking in progress')"></div>
            </details>
          </section>

          <section *ngIf="message.tool_calls?.length" class="tool">
            <strong>Tool calls</strong>
            <div *ngFor="let call of message.tool_calls" class="tool__block">
              <div class="tool__title">{{ call.name }}</div>
              <pre>{{ call.arguments | json }}</pre>
            </div>
          </section>

          <section *ngIf="message.tool_results?.length" class="tool">
            <strong>Tool results</strong>
            <div *ngFor="let result of message.tool_results" class="tool__block">
              <div class="tool__title">{{ result.tool_call_id }}</div>
              <pre>{{ result.content }}</pre>
            </div>
          </section>

          <footer [attr.title]="message.created_at | date : 'full'">
            {{ message.created_at | date: 'MMM d, HH:mm:ss' }}
          </footer>
        </article>
      </div>

      <button *ngIf="streaming && !pinnedToBottom" type="button" class="scroll-pill" (click)="scrollToBottom()">
        Scroll to bottom
      </button>

      <div *ngIf="lightboxAttachment" class="lightbox" (click)="closeLightbox()">
        <div class="lightbox__dialog" (click)="$event.stopPropagation()">
          <button type="button" class="lightbox__close" (click)="closeLightbox()">Close</button>
          <img
            [src]="attachmentSource(lightboxAttachment)"
            [alt]="lightboxAttachment.filename"
            class="lightbox__image"
          />
          <p class="lightbox__meta">
            {{ lightboxAttachment.filename }} · {{ lightboxAttachment.width }}×{{ lightboxAttachment.height }}
          </p>
        </div>
      </div>
    </div>
  `,
  styles: [
    `
      .thread-shell { position: relative; min-height: 100%; height: 100%; }
      .thread { display: grid; gap: 1rem; align-content: start; min-height: 0; max-height: 100%; overflow: auto; padding-right: 0.25rem; }
      .bubble { max-width: 54rem; padding: 1rem 1.1rem; border-radius: 1.25rem; background: rgba(11, 27, 44, 0.88); border: 1px solid rgba(125, 211, 252, 0.12); box-shadow: 0 12px 40px rgba(0, 0, 0, 0.18); }
      .bubble--user { margin-left: auto; background: linear-gradient(135deg, rgba(8, 47, 73, 0.95), rgba(14, 116, 144, 0.7)); }
      header, footer { display: flex; justify-content: space-between; gap: 1rem; color: #7dd3fc; font-size: 0.72rem; }
      header { margin-bottom: 0.65rem; text-transform: uppercase; letter-spacing: 0.08em; }
      footer { margin-top: 0.9rem; color: #94a3b8; opacity: 0; transition: opacity 0.2s ease; }
      .bubble:hover footer { opacity: 1; }
      .content { color: #ecfeff; line-height: 1.6; }
      .markdown :first-child { margin-top: 0; }
      .markdown :last-child { margin-bottom: 0; }
      .markdown a { color: #7dd3fc; }
      .markdown code { padding: 0.15rem 0.35rem; border-radius: 0.35rem; background: rgba(15, 23, 42, 0.75); }
      .markdown pre code { padding: 0; background: transparent; }
      .attachments { display: flex; gap: 0.8rem; flex-wrap: wrap; margin-top: 0.9rem; }
      .attachment { width: min(16rem, 100%); margin: 0; display: grid; gap: 0.4rem; }
      .attachment__button { border: 0; padding: 0; background: transparent; cursor: zoom-in; }
      .attachment img { width: 100%; border-radius: 1rem; }
      .attachment figcaption { color: #cbd5e1; font-size: 0.76rem; }
      .thinking, .tool { margin-top: 0.85rem; padding-top: 0.85rem; border-top: 1px solid rgba(125, 211, 252, 0.12); color: #cbd5e1; }
      .thinking { display: grid; gap: 0.5rem; }
      .thinking__badge, .thinking-banner { display: inline-flex; width: fit-content; gap: 0.45rem; align-items: center; border-radius: 999px; padding: 0.3rem 0.65rem; background: rgba(15, 23, 42, 0.7); color: #f8fafc; }
      .thinking-banner { position: sticky; top: 0; z-index: 1; margin-bottom: 0.2rem; border: 1px solid rgba(245, 158, 11, 0.2); }
      .thinking-banner__pulse, .thinking__badge--active::before { content: ''; width: 0.55rem; height: 0.55rem; border-radius: 50%; background: #f59e0b; box-shadow: 0 0 0 0 rgba(245, 158, 11, 0.5); animation: pulse 1.4s infinite; }
      details summary { cursor: pointer; color: #e2e8f0; }
      .tool { display: grid; gap: 0.65rem; }
      .tool__block { display: grid; gap: 0.35rem; }
      .tool__title { color: #f8fafc; font-weight: 700; }
      .code-block { margin-top: 0.9rem; overflow: hidden; border-radius: 0.95rem; border: 1px solid rgba(148, 163, 184, 0.16); background: rgba(2, 6, 23, 0.64); }
      .code-block__bar { display: flex; justify-content: space-between; gap: 1rem; align-items: center; padding: 0.55rem 0.8rem; border-bottom: 1px solid rgba(148, 163, 184, 0.12); color: #cbd5e1; font-size: 0.72rem; text-transform: uppercase; letter-spacing: 0.08em; }
      .code-block__copy, .lightbox__close { border: 1px solid rgba(148, 163, 184, 0.2); border-radius: 999px; background: rgba(15, 23, 42, 0.8); color: #e2e8f0; padding: 0.3rem 0.7rem; cursor: pointer; }
      pre { margin: 0; overflow: auto; padding: 0.85rem; color: #f8fafc; }
      .stream-cursor { margin-top: 0.6rem; width: 0.8rem; height: 1.2rem; border-radius: 0.2rem; background: #fde68a; animation: blink 1s steps(2, start) infinite; }
      .scroll-pill { position: absolute; right: 0.75rem; bottom: 0.75rem; border: 1px solid rgba(251, 191, 36, 0.28); border-radius: 999px; background: rgba(124, 45, 18, 0.92); color: #fde68a; padding: 0.75rem 1rem; cursor: pointer; box-shadow: 0 10px 24px rgba(0, 0, 0, 0.25); }
      .lightbox { position: absolute; inset: 0; z-index: 10; display: grid; place-items: center; padding: 1.5rem; background: rgba(2, 6, 23, 0.9); backdrop-filter: blur(10px); }
      .lightbox__dialog { display: grid; gap: 0.9rem; max-width: min(72rem, 100%); max-height: 100%; }
      .lightbox__image { max-width: min(72rem, 100%); max-height: 78vh; border-radius: 1.2rem; }
      .lightbox__meta { margin: 0; color: #cbd5e1; text-align: center; }
      @keyframes pulse {
        0% { transform: scale(0.92); opacity: 0.75; }
        70% { transform: scale(1.08); opacity: 1; }
        100% { transform: scale(0.92); opacity: 0.75; }
      }
      @keyframes blink {
        0%, 49% { opacity: 1; }
        50%, 100% { opacity: 0.2; }
      }
    `,
  ],
})
export class MessageListComponent implements OnChanges, AfterViewChecked {
  @Input() messages: ChatMessage[] = [];
  @Input() streaming = false;
  @Input() thinkingActive = false;

  @ViewChild('viewport') private readonly viewport?: ElementRef<HTMLDivElement>;

  lightboxAttachment: ImageAttachment | null = null;
  private pendingScroll = false;
  pinnedToBottom = true;

  constructor(private readonly sanitizer: DomSanitizer) {}

  ngOnChanges(changes: SimpleChanges): void {
    if (changes['messages'] || changes['streaming'] || changes['thinkingActive']) {
      this.pendingScroll = true;
    }
  }

  ngAfterViewChecked(): void {
    if (this.pendingScroll) {
      this.pendingScroll = false;
      if (this.pinnedToBottom) {
        queueMicrotask(() => this.scrollToBottom());
      }
    }
  }

  trackByMessage(_index: number, message: ChatMessage): string {
    return message.id;
  }

  attachmentSource(attachment: ImageAttachment): string {
    return `data:${attachment.mime_type};base64,${attachment.data}`;
  }

  openLightbox(attachment: ImageAttachment): void {
    this.lightboxAttachment = attachment;
  }

  closeLightbox(): void {
    this.lightboxAttachment = null;
  }

  isLatestAssistantMessage(message: ChatMessage): boolean {
    const assistantMessages = this.messages.filter((entry) => entry.role === 'assistant');
    return assistantMessages.at(-1)?.id === message.id;
  }

  thinkingDuration(message: ChatMessage): string {
    const duration = message.thinking?.duration_ms ?? 0;
    return duration > 0 ? `${(duration / 1000).toFixed(1)}s` : 'in progress';
  }

  renderBlocks(content: string): RenderBlock[] {
    const source = content ?? '';
    if (!source.trim()) {
      return [];
    }

    const blocks: RenderBlock[] = [];
    const pattern = /```([\w-]+)?\n([\s\S]*?)```/g;
    let lastIndex = 0;
    let match: RegExpExecArray | null;

    while ((match = pattern.exec(source)) !== null) {
      const before = source.slice(lastIndex, match.index).trim();
      if (before) {
        blocks.push({ kind: 'text', html: this.renderMarkdown(before) });
      }
      const code = match[2].replace(/\n+$/, '');
      blocks.push({
        kind: 'code',
        language: match[1] ?? '',
        code,
        html: this.highlightCode(code, match[1] ?? ''),
      });
      lastIndex = pattern.lastIndex;
    }

    const after = source.slice(lastIndex).trim();
    if (after) {
      blocks.push({ kind: 'text', html: this.renderMarkdown(after) });
    }

    if (blocks.length === 0) {
      blocks.push({ kind: 'text', html: this.renderMarkdown(source) });
    }
    return blocks;
  }

  renderMarkdown(content: string): SafeHtml {
    const sanitized = (content ?? '').replace(/<[^>]*>/g, '');
    const html = marked.parse(sanitized, { async: false, gfm: true, breaks: true });
    return this.sanitizer.bypassSecurityTrustHtml(String(html));
  }

  copyCode(code: string): void {
    if (!code) {
      return;
    }
    if (navigator.clipboard?.writeText) {
      void navigator.clipboard.writeText(code);
      return;
    }
    const textarea = document.createElement('textarea');
    textarea.value = code;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.focus();
    textarea.select();
    document.execCommand('copy');
    document.body.removeChild(textarea);
  }

  onScroll(): void {
    const element = this.viewport?.nativeElement;
    if (!element) {
      return;
    }
    const distanceFromBottom = element.scrollHeight - element.scrollTop - element.clientHeight;
    this.pinnedToBottom = distanceFromBottom < 48;
  }

  scrollToBottom(): void {
    const element = this.viewport?.nativeElement;
    if (!element) {
      return;
    }
    element.scrollTop = element.scrollHeight;
    this.pinnedToBottom = true;
  }

  private highlightCode(code: string, language: string): SafeHtml {
    const normalized = language.trim().toLowerCase();
    const value =
      normalized && hljs.getLanguage(normalized)
        ? hljs.highlight(code, { language: normalized }).value
        : hljs.highlightAuto(code).value;
    return this.sanitizer.bypassSecurityTrustHtml(value);
  }
}
