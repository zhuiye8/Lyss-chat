import api from './api'
import { User } from '../store/authStore'

interface LoginRequest {
  email: string
  password: string
  tenant_id?: string
}

interface LoginResponse {
  access_token: string
  refresh_token: string
  expires_in: number
  user: User
}

export const authService = {
  login: async (data: LoginRequest): Promise<LoginResponse> => {
    return api.post('/auth/login', data)
  },
  
  register: async (data: {
    email: string
    password: string
    name: string
    tenant_id?: string
  }): Promise<User> => {
    return api.post('/auth/register', data)
  },
  
  refreshToken: async (refreshToken: string): Promise<LoginResponse> => {
    return api.post('/auth/refresh', { refresh_token: refreshToken })
  },
  
  getCurrentUser: async (): Promise<User> => {
    return api.get('/users/me')
  }
}
