export interface ThinkingState {
  active: boolean;
  content: string;
  started_at?: string;
  ended_at?: string;
  duration_ms?: number;
}

export interface ToolCall {
  id: string;
  name: string;
  arguments: Record<string, unknown>;
}

export interface ToolResult {
  tool_call_id: string;
  content: string;
}

export interface ImageAttachment {
  filename: string;
  mime_type: string;
  data: string;
  width: number;
  height: number;
}

export interface ChatMessage {
  id: string;
  role: string;
  content: string;
  created_at: string;
  model?: string;
  finish_reason?: string;
  thinking?: ThinkingState;
  tool_calls?: ToolCall[];
  tool_results?: ToolResult[];
  attachments?: ImageAttachment[];
}

export interface ConversationSummary {
  id: string;
  title: string;
  model: string;
  created_at: string;
  updated_at: string;
  message_count: number;
}

export interface Conversation extends ConversationSummary {
  messages: ChatMessage[];
  settings?: ChatSettings | null;
}

export interface ModelEntry {
  name: string;
  architecture: string;
  quant_bits: number;
  size_bytes: number;
  loaded: boolean;
  backend: string;
  supports_vision: boolean;
}

export interface ChatSettings {
  temperature: number;
  top_p: number;
  top_k: number;
  max_tokens: number;
  context_window: number;
  system_prompt: string;
  default_model: string;
}
