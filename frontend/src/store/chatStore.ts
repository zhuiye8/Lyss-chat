import { create } from 'zustand'

export interface Canvas {
  id: string
  workspace_id: string
  title: string
  description?: string
  type: 'chat' | 'code'
  status: 'active' | 'archived'
  model_id?: string
  created_at: string
}

export interface Message {
  id: string
  canvas_id: string
  parent_id?: string
  role: 'user' | 'assistant' | 'system'
  content: string
  metadata?: Record<string, any>
  created_at: string
}

interface ChatState {
  canvases: Canvas[]
  currentCanvas: Canvas | null
  messages: Message[]
  loading: boolean
  error: string | null
  
  setCanvases: (canvases: Canvas[]) => void
  setCurrentCanvas: (canvas: Canvas | null) => void
  setMessages: (messages: Message[]) => void
  addMessage: (message: Message) => void
  updateMessage: (id: string, content: string) => void
  setLoading: (loading: boolean) => void
  setError: (error: string | null) => void
  clearMessages: () => void
}

export const useChatStore = create<ChatState>((set) => ({
  canvases: [],
  currentCanvas: null,
  messages: [],
  loading: false,
  error: null,
  
  setCanvases: (canvases) => set({ canvases }),
  setCurrentCanvas: (canvas) => set({ currentCanvas: canvas }),
  setMessages: (messages) => set({ messages }),
  addMessage: (message) => set((state) => ({ 
    messages: [...state.messages, message] 
  })),
  updateMessage: (id, content) => set((state) => ({
    messages: state.messages.map(msg => 
      msg.id === id ? { ...msg, content } : msg
    )
  })),
  setLoading: (loading) => set({ loading }),
  setError: (error) => set({ error }),
  clearMessages: () => set({ messages: [] })
}))
