import { api } from '@/api/client'
import type {
  Chat,
  ChatList,
  Message,
  MessageList,
  CreateChatRequest,
  UpdateChatRequest,
  SendMessageRequest,
} from '@/types/chat'

export const createChat = (payload: CreateChatRequest) => {
  return api.post<Chat>('/chats', payload)
}

export const getChats = (limit = 50) => {
  return api.get<ChatList>('/chats', {
    params: { limit },
  })
}

export const getChatById = (chatId: string) => {
  return api.get<Chat>(`/chats/${chatId}`)
}

export const updateChat = (chatId: string, payload: UpdateChatRequest) => {
  return api.put<Chat>(`/chats/${chatId}`, payload)
}

export const deleteChat = (chatId: string) => {
  return api.delete(`/chats/${chatId}`)
}

export const generateTitle = (chatId: string) => {
  return api.post<{ title: string }>(`/chats/${chatId}/generate-title`)
}

export const getHistory = (chatId: string, limit = 50, before?: string) => {
  return api.get<MessageList>(`/chats/${chatId}/messages`, {
    params: { limit, before },
  })
}

export const sendMessage = (chatId: string, payload: SendMessageRequest, stream = false) => {
  return api.post<Message>(`/chats/${chatId}/messages`, payload, {
    params: { stream },
  })
}

export const streamMessage = (
  chatId: string,
  content: string,
  onToken: (token: string) => void,
  onDone: () => void,
  onError: (error: string) => void,
): (() => void) => {
  const controller = new AbortController()

  fetch(`/api/chats/${chatId}/messages?stream=true`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    body: JSON.stringify({ content }),
    signal: controller.signal,
  })
    .then(async (res) => {
      if (!res.ok || !res.body) {
        onError('Request failed')
        return
      }

      const reader = res.body.getReader()
      const decoder = new TextDecoder()
      let buffer = ''

      while (true) {
        const { done, value } = await reader.read()
        if (done) break

        buffer += decoder.decode(value, { stream: true })
        const lines = buffer.split('\n')
        buffer = lines.pop() ?? ''

        for (const line of lines) {
          if (!line.startsWith('data:')) continue
          try {
            const json = JSON.parse(line.slice(5).trim())
            if (json.token) onToken(json.token)
            if (json.done) {
              onDone()
              return
            }
          } catch {
            // skip malformed lines
          }
        }
      }
      onDone()
    })
    .catch((err) => {
      if (err.name !== 'AbortError') {
        onError(err.message || 'Stream failed')
      }
    })

  return () => controller.abort()
}
