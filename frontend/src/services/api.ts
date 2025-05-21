import axios from 'axios'
import { useAuthStore } from '../store/authStore'

const API_BASE_URL = '/api/v1'

const api = axios.create({
  baseURL: API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
api.interceptors.request.use(
  (config) => {
    const token = useAuthStore.getState().token
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// 响应拦截器
api.interceptors.response.use(
  (response) => {
    return response.data
  },
  async (error) => {
    const originalRequest = error.config
    
    // 如果是 401 错误且没有重试过，尝试刷新 token
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true
      
      try {
        const refreshToken = useAuthStore.getState().refreshToken
        
        if (refreshToken) {
          const response = await axios.post(`${API_BASE_URL}/auth/refresh`, {
            refresh_token: refreshToken
          })
          
          const { access_token } = response.data
          
          // 更新 token
          useAuthStore.getState().updateToken(access_token)
          
          // 更新原始请求的 Authorization 头
          originalRequest.headers.Authorization = `Bearer ${access_token}`
          
          // 重试原始请求
          return axios(originalRequest)
        }
      } catch (refreshError) {
        // 刷新 token 失败，登出用户
        useAuthStore.getState().logout()
        window.location.href = '/login'
        return Promise.reject(refreshError)
      }
    }
    
    return Promise.reject(error)
  }
)

export default api
