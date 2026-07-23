import { computed, Injectable, effect, inject, signal } from '@angular/core';
import { WebSocketService } from '../services/websocket.service';
import { ChatRoute, ChatRouteMap, CoreGuiChatBindings } from '../generated/core-gui-chat.bindings';
import {
  ChatMessage,
  ChatSettings,
  Conversation,
  ConversationSummary,
  ImageAttachment,
  ModelEntry,
  ThinkingState,
  ToolCall,
  ToolResult,
} from './chat.types';

@Injectable({ providedIn: 'root' })
export class ChatStateService {
  private readonly supportedImageMimeTypes = new Set(['image/png', 'image/jpeg', 'image/webp', 'image/gif']);
  private readonly supportedImageExtensions = ['.png', '.jpg', '.jpeg', '.webp', '.gif'];
  private readonly ws = inject(WebSocketService);
  private readonly bindings = new CoreGuiChatBindings((route, payload) => this.mockInvoke(route, payload));
  private mockConversations: Conversation[] = [];

  readonly conversations = signal<ConversationSummary[]>([]);
  readonly activeConversation = signal<Conversation | null>(null);
  readonly models = signal<ModelEntry[]>([]);
  readonly queuedAttachments = signal<ImageAttachment[]>([]);
  readonly settings = signal<ChatSettings>({
    temperature: 1,
    top_p: 0.95,
    top_k: 64,
    max_tokens: 2048,
    context_window: 8192,
    system_prompt: 'You are a helpful assistant.',
    default_model: '',
  });
  readonly draft = signal('');
  readonly historyQuery = signal('');
  readonly settingsOpen = signal(false);
  readonly sending = signal(false);
  readonly modelSwitching = signal(false);
  readonly selectedModel = signal('');
  readonly nativeDialogAvailable = computed(
    () => typeof window !== 'undefined' && typeof window.__CORE_GUI_INVOKE__ === 'function',
  );
  readonly selectedModelEntry = computed(
    () => this.models().find((model) => model.name === this.selectedModel()) ?? null,
  );
  readonly selectedModelSupportsVision = computed(
    () => this.supportsVision(this.selectedModelEntry()),
  );
  readonly thinkingActive = computed(
    () => this.activeConversation()?.messages.some((message) => message.thinking?.active) ?? false,
  );

  constructor() {
    effect(() => {
      const current = this.activeConversation();
      if (current?.model) {
        this.selectedModel.set(current.model);
      } else if (this.settings().default_model) {
        this.selectedModel.set(this.settings().default_model);
      }
    });
  }

  async init(): Promise<void> {
    this.ws.connect();
    this.bindEvents();
    await this.loadBootstrap();
  }

  async refreshConversation(id: string): Promise<void> {
    const conversation = await this.invoke('gui.chat.conversations.get', { id });
    if (conversation) {
      this.activeConversation.set(conversation);
      this.queuedAttachments.set([]);
    }
  }

  async startConversation(): Promise<void> {
    const conversation = await this.invoke('gui.chat.conversations.new');
    if (conversation) {
      this.activeConversation.set(conversation);
      this.upsertSummary(this.toSummary(conversation));
      this.queuedAttachments.set([]);
    }
  }

  async deleteConversation(id: string): Promise<void> {
    await this.invoke('gui.chat.conversations.delete', { id });
    this.removeSummary(id);
    if (this.activeConversation()?.id === id) {
      this.activeConversation.set(null);
      this.queuedAttachments.set([]);
      const next = this.conversations()[0];
      if (next) {
        await this.refreshConversation(next.id);
      }
    }
  }

  async renameConversation(id: string, title: string): Promise<void> {
    const updated = await this.invoke('gui.chat.conversations.rename', { id, title });
    if (updated) {
      this.mergeConversation(updated);
    }
  }

  async exportConversation(id: string): Promise<void> {
    const markdown = await this.invoke('gui.chat.conversations.export', { id });
    if (!markdown) {
      return;
    }
    const blob = new Blob([markdown], { type: 'text/markdown;charset=utf-8' });
    const url = URL.createObjectURL(blob);
    const anchor = document.createElement('a');
    anchor.href = url;
    anchor.download = `${id}.md`;
    anchor.click();
    URL.revokeObjectURL(url);
  }

  async setHistoryQuery(query: string): Promise<void> {
    this.historyQuery.set(query);
    const trimmed = query.trim();
    const route = trimmed ? 'gui.chat.conversations.search' : 'gui.chat.conversations.list';
    const payload = trimmed ? { q: trimmed } : undefined;
    const conversations = await this.invoke(route, payload);
    this.conversations.set(conversations ?? []);
  }

  async saveSettings(settings: ChatSettings): Promise<void> {
    const saved = await this.invoke('gui.chat.settings.save', settings);
    if (saved) {
      this.settings.set(saved);
      if (saved.default_model) {
        this.selectedModel.set(saved.default_model);
      }
    }
  }

  async resetSettings(): Promise<void> {
    const defaults = await this.invoke('gui.chat.settings.defaults');
    if (!defaults) {
      return;
    }
    const saved = await this.invoke('gui.chat.settings.save', defaults);
    if (saved) {
      this.settings.set(saved);
      this.selectedModel.set(saved.default_model || '');
    }
  }

  async changeModel(model: string): Promise<void> {
    if (!model) {
      return;
    }
    this.modelSwitching.set(true);
    try {
      const settings = await this.invoke('gui.chat.selectModel', {
        model,
        conversation_id: this.activeConversation()?.id,
      });
      this.selectedModel.set(model);
      if (settings) {
        this.settings.set(settings);
      }
      const currentId = this.activeConversation()?.id;
      this.activeConversation.update((conversation) => (conversation ? { ...conversation, model } : conversation));
      if (currentId) {
        this.conversations.update((items) =>
          items.map((item) => (item.id === currentId ? { ...item, model } : item)),
        );
      }
    } finally {
      this.modelSwitching.set(false);
    }
  }

  async queueImageFiles(files: FileList | File[]): Promise<void> {
    if (!this.selectedModelSupportsVision()) {
      return;
    }
    const items = Array.from(files);
    for (const file of items) {
      if (!this.isSupportedImageFile(file)) {
        continue;
      }
      const attachment = await this.fileToAttachment(file);
      const queued = await this.invoke('gui.chat.attachImage', {
        conversation_id: this.activeConversation()?.id,
        ...attachment,
      });
      if (queued) {
        this.upsertQueuedAttachment(queued);
      }
    }
  }

  async openImagePicker(): Promise<void> {
    if (!this.selectedModelSupportsVision() || !this.nativeDialogAvailable()) {
      return;
    }
    const paths = await this.invokeGUI<string[]>('gui.dialog.open', {
      title: 'Attach images',
      allowMultiple: true,
      filters: [
        {
          displayName: 'Images',
          pattern: '*.png;*.jpg;*.jpeg;*.webp;*.gif',
          extensions: ['png', 'jpg', 'jpeg', 'webp', 'gif'],
        },
      ],
    });
    if (!paths?.length) {
      return;
    }
    for (const path of paths) {
      const attachment = await this.invoke('gui.chat.attachImageFile', {
        conversation_id: this.activeConversation()?.id,
        path,
      });
      if (attachment) {
        this.upsertQueuedAttachment(attachment);
      }
    }
  }

  async removeQueuedAttachment(index: number): Promise<void> {
    const activeConversationID = this.activeConversation()?.id;
    const removed = await this.invoke('gui.chat.removeImage', {
      conversation_id: activeConversationID,
      index,
    });
    if (!removed) {
      return;
    }
    this.queuedAttachments.update((items) => items.filter((_, itemIndex) => itemIndex !== index));
  }

  async sendMessage(): Promise<void> {
    const content = this.draft().trim();
    if (!content && this.queuedAttachments().length === 0) {
      return;
    }
    this.sending.set(true);
    try {
      const response = await this.invoke('gui.chat.send', {
        conversation_id: this.activeConversation()?.id,
        content,
      });
      if (response) {
        this.activeConversation.set(response);
        this.mergeConversation(response);
        this.draft.set('');
        this.queuedAttachments.set([]);
      }
    } finally {
      this.sending.set(false);
    }
  }

  private async loadBootstrap(): Promise<void> {
    const [models, settings, conversations] = await Promise.all([
      this.invoke('gui.chat.models'),
      this.invoke('gui.chat.settings.load'),
      this.invoke('gui.chat.conversations.list'),
    ]);

    if (models?.length) {
      this.models.set(models);
      const current = models.find((item) => item.loaded) ?? models[0];
      if (current) {
        this.selectedModel.set(current.name);
      }
    }

    if (settings) {
      this.settings.set(settings);
      if (settings.default_model) {
        this.selectedModel.set(settings.default_model);
      }
    }

    if (conversations?.length) {
      this.conversations.set(conversations);
      await this.refreshConversation(conversations[0].id);
      return;
    }
    await this.startConversation();
  }

  private bindEvents(): void {
    this.ws.on('chat.conversation', (payload) => {
      const data = payload as { action?: string; conversation?: Conversation; conversationId?: string; conversation_id?: string };
      const conversationID = data.conversation_id ?? data.conversationId ?? '';
      if (data.conversation) {
        this.mergeConversation(data.conversation);
      }
      if (data.action === 'deleted' && conversationID) {
        this.removeSummary(conversationID);
        if (this.activeConversation()?.id === conversationID) {
          this.activeConversation.set(null);
        }
      }
      if (data.action === 'cleared' && conversationID && this.activeConversation()?.id === conversationID) {
        this.activeConversation.update((conversation) => (conversation ? { ...conversation, messages: [] } : conversation));
      }
    });

    this.ws.on('chat.message', (payload) => {
      const data = payload as { conversationId?: string; conversation_id?: string; message?: ChatMessage; messageId?: string; message_id?: string; state?: string; finishReason?: string };
      const conversationID = data.conversation_id ?? data.conversationId ?? '';
      if (this.activeConversation()?.id !== conversationID) {
        return;
      }
      if (data.message) {
        this.upsertMessage(conversationID, data.message);
      }
      if (data.state === 'started') {
        this.upsertMessage(conversationID, {
          id: data.message_id ?? data.messageId ?? '',
          role: 'assistant',
          content: '',
          created_at: new Date().toISOString(),
          model: this.selectedModel(),
        });
      }
      if (data.state === 'finished') {
        this.patchMessage(conversationID, data.message_id ?? data.messageId ?? '', (message) => ({
          ...message,
          finish_reason: data.finishReason,
        }));
      }
    });

    this.ws.on('chat.token', (payload) => {
      const data = payload as { conversationId?: string; conversation_id?: string; messageId?: string; message_id?: string; content?: string };
      const conversationID = data.conversation_id ?? data.conversationId ?? '';
      const messageID = data.message_id ?? data.messageId ?? '';
      this.patchMessage(conversationID, messageID, (message) => ({
        ...message,
        content: `${message.content ?? ''}${data.content ?? ''}`,
      }));
    });

    this.ws.on('chat.thinking.start', (payload) => {
      const data = payload as { conversationId?: string; conversation_id?: string; messageId?: string; message_id?: string; startedAt?: string };
      this.patchThinking(data, (thinking) => ({
        ...thinking,
        active: true,
        started_at: data.startedAt,
      }));
    });

    this.ws.on('chat.thinking.append', (payload) => {
      const data = payload as { conversationId?: string; conversation_id?: string; messageId?: string; message_id?: string; content?: string };
      this.patchThinking(data, (thinking) => ({
        ...thinking,
        content: `${thinking.content ?? ''}${data.content ?? ''}`,
      }));
    });

    this.ws.on('chat.thinking.end', (payload) => {
      const data = payload as { conversationId?: string; conversation_id?: string; messageId?: string; message_id?: string; durationMs?: number };
      this.patchThinking(data, (thinking) => ({
        ...thinking,
        active: false,
        duration_ms: data.durationMs,
      }));
    });

    this.ws.on('chat.tool.call', (payload) => {
      const data = payload as { conversationId?: string; conversation_id?: string; messageId?: string; message_id?: string; call?: ToolCall };
      const conversationID = data.conversation_id ?? data.conversationId ?? '';
      const messageID = data.message_id ?? data.messageId ?? '';
      if (!data.call) {
        return;
      }
      this.patchMessage(conversationID, messageID, (message) => ({
        ...message,
        tool_calls: [...(message.tool_calls ?? []), data.call!],
      }));
    });

    this.ws.on('chat.tool.result', (payload) => {
      const data = payload as { conversationId?: string; conversation_id?: string; messageId?: string; message_id?: string; result?: ToolResult };
      const conversationID = data.conversation_id ?? data.conversationId ?? '';
      const messageID = data.message_id ?? data.messageId ?? '';
      if (!data.result) {
        return;
      }
      this.patchMessage(conversationID, messageID, (message) => ({
        ...message,
        tool_results: [...(message.tool_results ?? []), data.result!],
      }));
    });

    this.ws.on('chat.image.queued', (payload) => {
      const data = payload as { conversationId?: string; conversation_id?: string; attachment?: ImageAttachment };
      const conversationID = data.conversation_id ?? data.conversationId ?? '';
      const activeID = this.activeConversation()?.id ?? 'draft';
      if (!data.attachment || conversationID !== activeID && !(conversationID === 'draft' && !this.activeConversation()?.messages.length)) {
        return;
      }
      this.upsertQueuedAttachment(data.attachment);
    });
  }

  private mergeConversation(conversation: Conversation): void {
    this.upsertSummary(this.toSummary(conversation));
    if (this.activeConversation()?.id === conversation.id || !this.activeConversation()) {
      this.activeConversation.set(conversation);
    }
  }

  private upsertSummary(summary: ConversationSummary): void {
    this.conversations.update((items) => {
      const next = [summary, ...items.filter((item) => item.id !== summary.id)];
      return next.sort((left, right) => Date.parse(right.updated_at) - Date.parse(left.updated_at));
    });
  }

  private removeSummary(id: string): void {
    this.conversations.update((items) => items.filter((item) => item.id !== id));
  }

  private upsertMessage(conversationID: string, message: ChatMessage): void {
    this.activeConversation.update((conversation) => {
      if (!conversation || conversation.id !== conversationID) {
        return conversation;
      }
      const existingIndex = conversation.messages.findIndex((item) => item.id === message.id);
      const messages = [...conversation.messages];
      if (existingIndex >= 0) {
        messages[existingIndex] = { ...messages[existingIndex], ...message };
      } else {
        messages.push(message);
      }
      return { ...conversation, messages, updated_at: message.created_at || conversation.updated_at };
    });
  }

  private patchMessage(conversationID: string, messageID: string, update: (message: ChatMessage) => ChatMessage): void {
    if (!conversationID || !messageID || this.activeConversation()?.id !== conversationID) {
      return;
    }
    this.activeConversation.update((conversation) => {
      if (!conversation) {
        return conversation;
      }
      return {
        ...conversation,
        messages: conversation.messages.map((message) => (message.id === messageID ? update(message) : message)),
      };
    });
  }

  private patchThinking(
    payload: { conversationId?: string; conversation_id?: string; messageId?: string; message_id?: string },
    update: (thinking: ThinkingState) => ThinkingState,
  ): void {
    const conversationID = payload.conversation_id ?? payload.conversationId ?? '';
    const messageID = payload.message_id ?? payload.messageId ?? '';
    this.patchMessage(conversationID, messageID, (message) => ({
      ...message,
      thinking: update(message.thinking ?? { active: false, content: '' }),
    }));
  }

  private toSummary(conversation: Conversation): ConversationSummary {
    return {
      id: conversation.id,
      title: conversation.title,
      model: conversation.model,
      created_at: conversation.created_at,
      updated_at: conversation.updated_at,
      message_count: conversation.messages?.length ?? 0,
    };
  }

  // The UI gates image input from the backend capability contract rather than inferring from names.
  // Example: this.supportsVision({ name: 'lemer', supports_vision: true } as ModelEntry) == true
  private supportsVision(model: ModelEntry | null): boolean {
    return model?.supports_vision ?? false;
  }

  private upsertQueuedAttachment(attachment: ImageAttachment): void {
    this.queuedAttachments.update((items) => {
      if (items.some((item) => this.sameAttachment(item, attachment))) {
        return items;
      }
      return [...items, attachment];
    });
  }

  private sameAttachment(left: ImageAttachment, right: ImageAttachment): boolean {
    return left.filename === right.filename &&
      left.mime_type === right.mime_type &&
      left.data === right.data &&
      left.width === right.width &&
      left.height === right.height;
  }

  private async fileToAttachment(file: File): Promise<ImageAttachment> {
    const data = await this.readFileAsDataURL(file);
    const dimensions = await this.readImageDimensions(data);
    return {
      filename: file.name,
      mime_type: file.type || 'image/png',
      data: data.split(',', 2)[1] ?? data,
      width: dimensions.width,
      height: dimensions.height,
    };
  }

  private readFileAsDataURL(file: File): Promise<string> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();
      reader.onerror = () => reject(reader.error);
      reader.onload = () => resolve(String(reader.result ?? ''));
      reader.readAsDataURL(file);
    });
  }

  private isSupportedImageFile(file: File): boolean {
    const mimeType = (file.type ?? '').toLowerCase();
    if (this.supportedImageMimeTypes.has(mimeType)) {
      return true;
    }
    const lowerName = file.name.toLowerCase();
    return this.supportedImageExtensions.some((suffix) => lowerName.endsWith(suffix));
  }

  private readImageDimensions(source: string): Promise<{ width: number; height: number }> {
    return new Promise((resolve) => {
      const image = new Image();
      image.onload = () => resolve({ width: image.width, height: image.height });
      image.onerror = () => resolve({ width: 0, height: 0 });
      image.src = source;
    });
  }

  private invoke<RouteName extends ChatRoute>(
    route: RouteName,
    payload?: ChatRouteMap[RouteName]['request'],
  ): Promise<ChatRouteMap[RouteName]['response']> {
    return this.bindings.invoke(route, payload);
  }

  private async invokeGUI<T>(route: string, payload?: unknown): Promise<T | null> {
    if (typeof window === 'undefined' || typeof window.__CORE_GUI_INVOKE__ !== 'function') {
      return null;
    }
    return (await window.__CORE_GUI_INVOKE__(route, payload)) as T;
  }

  private async mockInvoke<T>(route: ChatRoute, payload?: unknown): Promise<T> {
    if (route === 'gui.chat.models') {
      return [
        { name: 'lemer', architecture: 'gemma3', quant_bits: 4, size_bytes: 1500000000, loaded: true, backend: 'metal', supports_vision: true },
        { name: 'lemma', architecture: 'qwen3', quant_bits: 8, size_bytes: 3200000000, loaded: false, backend: 'ollama', supports_vision: false },
      ] as T;
    }
    if (route === 'gui.chat.settings.load') {
      return this.settings() as T;
    }
    if (route === 'gui.chat.settings.defaults') {
      return {
        temperature: 1,
        top_p: 0.95,
        top_k: 64,
        max_tokens: 2048,
        context_window: 8192,
        system_prompt: 'You are a helpful assistant.',
        default_model: '',
      } as T;
    }
    if (route === 'gui.chat.settings.save') {
      const settings = payload as ChatSettings;
      this.settings.set(settings);
      if (settings.default_model) {
        this.selectedModel.set(settings.default_model);
      }
      return settings as T;
    }
    if (route === 'gui.chat.settings.reset') {
      const defaults = (await this.mockInvoke('gui.chat.settings.defaults')) as ChatSettings;
      this.settings.set(defaults);
      return defaults as T;
    }
    if (route === 'gui.chat.selectModel') {
      const model = (payload as { model?: string })?.model ?? this.settings().default_model;
      const updated = {
        ...this.settings(),
        default_model: model,
      };
      if (this.activeConversation()) {
        const currentId = this.activeConversation()?.id;
        this.activeConversation.update((conversation) => (conversation ? { ...conversation, model } : conversation));
        if (currentId) {
          this.conversations.update((items) =>
            items.map((item) => (item.id === currentId ? { ...item, model } : item)),
          );
          this.mockConversations = this.mockConversations.map((item) => (item.id === currentId ? { ...item, model } : item));
        }
      }
      return updated as T;
    }
    if (route === 'gui.chat.conversations.get') {
      const id = (payload as { id?: string; conversation_id?: string })?.id ?? (payload as { id?: string; conversation_id?: string })?.conversation_id;
      const found = this.mockConversations.find((item) => item.id === id);
      if (found) {
        return found as T;
      }
      if (this.activeConversation()?.id === id) {
        return this.activeConversation() as T;
      }
      return null as T;
    }
    if (route === 'gui.chat.conversations.rename') {
      const { id, title } = payload as { id?: string; title?: string };
      const titleValue = title?.trim() || 'New Chat';
      const currentId = this.activeConversation()?.id;
      if (currentId === id) {
        this.activeConversation.update((conversation) => (conversation ? { ...conversation, title: titleValue } : conversation));
      }
      this.conversations.update((items) =>
        items.map((item) => (item.id === id ? { ...item, title: titleValue } : item)),
      );
      this.mockConversations = this.mockConversations.map((item) => (item.id === id ? { ...item, title: titleValue } : item));
      const found = this.mockConversations.find((item) => item.id === id);
      if (found) {
        return found as T;
      }
      if (this.activeConversation()?.id === id) {
        return this.activeConversation() as T;
      }
      return { id, title: titleValue } as T;
    }
    if (route === 'gui.chat.conversations.delete') {
      const id = (payload as { id?: string; conversation_id?: string })?.id ?? (payload as { id?: string; conversation_id?: string })?.conversation_id;
      this.conversations.update((items) => items.filter((item) => item.id !== id));
      this.mockConversations = this.mockConversations.filter((item) => item.id !== id);
      if (this.activeConversation()?.id === id) {
        this.activeConversation.set(null);
      }
      return undefined as T;
    }
    if (route === 'gui.chat.conversations.list' || route === 'gui.chat.conversations.search') {
      const query = ((payload as { q?: string })?.q ?? '').trim().toLowerCase();
      if (!query) {
        return this.mockConversations.map((item) => this.toSummary(item)) as T;
      }
      return this.mockConversations
        .filter((item) =>
          item.title.toLowerCase().includes(query) ||
          item.model.toLowerCase().includes(query) ||
          item.messages.some((message) => message.content.toLowerCase().includes(query)),
        )
        .map((item) => this.toSummary(item)) as T;
    }
    if (route === 'gui.chat.conversations.export') {
      const id = (payload as { id?: string; conversation_id?: string })?.id ?? (payload as { id?: string; conversation_id?: string })?.conversation_id;
      const found = this.mockConversations.find((item) => item.id === id);
      const conversation = found ?? (this.activeConversation()?.id === id ? this.activeConversation() : null);
      if (!conversation) {
        return '# Exported Conversation\n' as T;
      }
      return [
        `# ${conversation.title}`,
        '',
        ...conversation.messages.map((message) => {
          const heading = `## ${message.role.charAt(0).toUpperCase() + message.role.slice(1)}`;
          const body = message.content ? `${message.content}\n` : '';
          return `${heading}\n\n${body}`.trimEnd();
        }),
        '',
      ].join('\n') as T;
    }
    if (route === 'gui.chat.attachImage') {
      const attachment = payload as ImageAttachment;
      this.queuedAttachments.update((items) => [...items, attachment]);
      return payload as T;
    }
    if (route === 'gui.chat.removeImage') {
      const index = (payload as { index?: number })?.index ?? -1;
      const current = this.queuedAttachments();
      if (index < 0 || index >= current.length) {
        return null as T;
      }
      const [removed] = current.slice(index, index + 1);
      this.queuedAttachments.update((items) => items.filter((_, itemIndex) => itemIndex !== index));
      return removed as T;
    }
    if (route === 'gui.chat.conversations.new') {
      const now = new Date().toISOString();
      const conversation = {
        id: `conv-${Date.now().toString(36)}`,
        title: 'New Chat',
        model: this.selectedModel() || 'lemer',
        created_at: now,
        updated_at: now,
        message_count: 0,
        messages: [],
      };
      this.activeConversation.set(conversation);
      this.mockConversations = [conversation, ...this.mockConversations.filter((item) => item.id !== conversation.id)];
      this.upsertSummary(conversation);
      return conversation as T;
    }
    if (route === 'gui.chat.send') {
      const conversation = this.activeConversation() ?? ((await this.mockInvoke('gui.chat.conversations.new')) as Conversation);
      const content = (payload as { content?: string })?.content ?? '';
      const now = new Date().toISOString();
      const attachments = this.queuedAttachments();
      const updated = {
        ...conversation,
        updated_at: now,
        title: conversation.title === 'New Chat' ? content.slice(0, 48) || 'New Chat' : conversation.title,
        messages: [
          ...conversation.messages,
          { id: `user-${Date.now()}`, role: 'user', content, created_at: now, model: this.selectedModel(), attachments },
          {
            id: `assistant-${Date.now() + 1}`,
            role: 'assistant',
            content: `Local mock response for: ${content || 'image input'}`,
            created_at: now,
            model: this.selectedModel(),
          },
        ],
      };
      this.draft.set('');
      this.queuedAttachments.set([]);
      this.activeConversation.set(updated);
      this.mockConversations = [updated, ...this.mockConversations.filter((item) => item.id !== updated.id)];
      this.upsertSummary(updated);
      return updated as T;
    }
    return {} as T;
  }

}
