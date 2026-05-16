export interface Chat {
  id: string
  model: string
  title: string
  created_at: string
}

export interface ChatList {
  chats: Chat[]
}

export type MessageRole = 'user' | 'assistant' | 'system' | 'tool'

export type MessageStatus = 'pending' | 'completed' | 'failed'

export interface Message {
  id: string
  role: MessageRole
  content: string
  model_name?: string
  input_tokens?: number
  output_tokens?: number
  tool_call_id?: string | null
  tool_name?: string | null
  tool_args?: Record<string, unknown> | null
  tool_result?: Record<string, unknown> | null
  status: MessageStatus
  created_at: string
}

export interface MessageList {
  messages: Message[]
  has_more: boolean
}

export interface StreamChunk {
  token: string
  done: boolean
}

export interface CreateChatRequest {
  model?: string
  title?: string
}

export interface UpdateChatRequest {
  title?: string
  model?: string
}

export interface SendMessageRequest {
  content: string
}

export interface ErrorResponse {
  error: string
}
