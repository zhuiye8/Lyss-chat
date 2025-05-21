import api from './api'
import { Canvas, Message } from '../store/chatStore'

interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

export const chatService = {
  // 获取画布列表
  getCanvases: async (
    workspace_id: string,
    params?: {
      type?: 'chat' | 'code'
      page?: number
      page_size?: number
    }
  ): Promise<PaginatedResponse<Canvas>> => {
    return api.get('/canvases', {
      params: {
        workspace_id,
        ...params
      }
    })
  },
  
  // 获取画布详情
  getCanvas: async (id: string): Promise<Canvas> => {
    return api.get(`/canvases/${id}`)
  },
  
  // 创建画布
  createCanvas: async (data: {
    title: string
    workspace_id: string
    description?: string
    type?: 'chat' | 'code'
    model_id?: string
  }): Promise<Canvas> => {
    return api.post('/canvases', data)
  },
  
  // 更新画布
  updateCanvas: async (
    id: string,
    data: {
      title?: string
      description?: string
      status?: 'active' | 'archived'
      model_id?: string
    }
  ): Promise<Canvas> => {
    return api.put(`/canvases/${id}`, data)
  },
  
  // 删除画布
  deleteCanvas: async (id: string): Promise<void> => {
    return api.delete(`/canvases/${id}`)
  },
  
  // 获取画布消息列表
  getMessages: async (
    canvas_id: string,
    params?: {
      page?: number
      page_size?: number
    }
  ): Promise<PaginatedResponse<Message>> => {
    return api.get(`/canvases/${canvas_id}/messages`, {
      params
    })
  },
  
  // 发送消息
  sendMessage: async (
    canvas_id: string,
    data: {
      content: string
      parent_id?: string
      metadata?: Record<string, any>
    }
  ): Promise<Message> => {
    return api.post(`/canvases/${canvas_id}/messages`, data)
  }
}
