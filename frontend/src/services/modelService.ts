import api from './api'

export interface Model {
  id: string
  provider_id: string
  name: string
  status: 'active' | 'inactive'
  is_public: boolean
}

interface PaginatedResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}

export const modelService = {
  // 获取模型列表
  getModels: async (params?: {
    provider_id?: string
    status?: 'active' | 'inactive'
    page?: number
    page_size?: number
  }): Promise<PaginatedResponse<Model>> => {
    return api.get('/models', { params })
  }
}
