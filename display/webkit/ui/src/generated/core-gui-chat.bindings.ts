import {
  ChatSettings,
  Conversation,
  ConversationSummary,
  ImageAttachment,
  ModelEntry,
} from '../chat/chat.types';

declare global {
  interface Window {
    __CORE_GUI_INVOKE__?: (route: string, payload?: unknown) => Promise<unknown> | unknown;
  }
}

export interface ChatRouteMap {
  'gui.chat.clear': {
    request: { id?: string; conversation_id?: string };
    response: Conversation;
  };
  'gui.chat.history': {
    request: { id?: string; conversation_id?: string };
    response: Conversation;
  };
  'gui.chat.models': { request: void; response: ModelEntry[] };
  'gui.chat.settings.defaults': { request: void; response: ChatSettings };
  'gui.chat.settings.load': { request: void; response: ChatSettings };
  'gui.chat.settings.save': { request: ChatSettings; response: ChatSettings };
  'gui.chat.settings.reset': { request: void; response: ChatSettings };
  'gui.chat.selectModel': {
    request: { model: string; conversation_id?: string };
    response: ChatSettings;
  };
  'gui.chat.conversations.list': { request: void; response: ConversationSummary[] };
  'gui.chat.conversations.search': { request: { q: string }; response: ConversationSummary[] };
  'gui.chat.conversations.get': { request: { id?: string; conversation_id?: string }; response: Conversation };
  'gui.chat.conversations.new': { request: void; response: Conversation };
  'gui.chat.conversations.delete': { request: { id: string }; response: void };
  'gui.chat.conversations.rename': { request: { id: string; title: string }; response: Conversation };
  'gui.chat.conversations.export': { request: { id: string }; response: string };
  'gui.chat.conversation.save': { request: Conversation; response: Conversation };
  'gui.chat.attachImage': {
    request: ({ conversation_id?: string } & ImageAttachment);
    response: ImageAttachment;
  };
  'gui.chat.attachImageFile': {
    request: { conversation_id?: string; path: string };
    response: ImageAttachment;
  };
  'gui.chat.removeImage': {
    request: { conversation_id?: string; index: number };
    response: ImageAttachment;
  };
  'gui.chat.send': {
    request: { conversation_id?: string; content: string };
    response: Conversation;
  };
  'gui.chat.thinking.start': {
    request: { conversation_id: string; message_id?: string; started_at?: string };
    response: { conversation_id: string; message_id?: string; started_at?: string };
  };
  'gui.chat.thinking.append': {
    request: { conversation_id: string; message_id?: string; content?: string };
    response: string;
  };
  'gui.chat.thinking.end': {
    request: { conversation_id: string; message_id?: string; started_at?: string; duration_ms?: number };
    response: number;
  };
}

export type ChatRoute = keyof ChatRouteMap;

// Generated binding entry point for the Go chat service.
// chat.invoke('gui.chat.send', { conversation_id: 'conv-1', content: 'Hello' })
export class CoreGuiChatBindings {
  constructor(private readonly fallback: <T>(route: ChatRoute, payload?: unknown) => Promise<T>) {}

  async invoke<RouteName extends ChatRoute>(
    route: RouteName,
    payload?: ChatRouteMap[RouteName]['request'],
  ): Promise<ChatRouteMap[RouteName]['response']> {
    if (typeof window.__CORE_GUI_INVOKE__ === 'function') {
      return (await window.__CORE_GUI_INVOKE__(route, payload)) as ChatRouteMap[RouteName]['response'];
    }
    return this.fallback<ChatRouteMap[RouteName]['response']>(route, payload);
  }
}
