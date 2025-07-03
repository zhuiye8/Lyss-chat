/**
 * 统一的HTTP请求工具
 * 基于axios，包含请求/响应拦截器、错误处理、认证等功能
 */

import axios, { AxiosInstance, AxiosRequestConfig, AxiosResponse } from 'axios';
import { message } from 'antd';
import { API_CONFIG, HTTP_STATUS } from '@/config/api';

// 创建axios实例
const request: AxiosInstance = axios.create({
  baseURL: `${API_CONFIG.BASE_URL}${API_CONFIG.API_PREFIX}`,
  timeout: API_CONFIG.TIMEOUT,
  headers: {
    'Content-Type': 'application/json',
  },
});

// 获取存储的token
const getToken = (): string | null => {
  return localStorage.getItem('access_token');
};

// 存储token
export const setToken = (token: string): void => {
  localStorage.setItem('access_token', token);
};

// 清除token
export const removeToken = (): void => {
  localStorage.removeItem('access_token');
};

// 请求拦截器
request.interceptors.request.use(
  (config: AxiosRequestConfig) => {
    // 添加认证头
    const token = getToken();
    if (token && config.headers) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    
    // 打印请求信息 (开发环境)
    if (process.env.NODE_ENV === 'development') {
      console.log('🚀 API Request:', {
        method: config.method?.toUpperCase(),
        url: config.url,
        data: config.data,
        params: config.params,
      });
    }
    
    return config;
  },
  (error) => {
    console.error('请求配置错误:', error);
    return Promise.reject(error);
  },
);

// 响应拦截器
request.interceptors.response.use(
  (response: AxiosResponse) => {
    // 打印响应信息 (开发环境)
    if (process.env.NODE_ENV === 'development') {
      console.log('✅ API Response:', {
        status: response.status,
        url: response.config.url,
        data: response.data,
      });
    }
    
    return response;
  },
  (error) => {
    // 处理不同的错误状态
    if (error.response) {
      const { status, data } = error.response;
      
      switch (status) {
        case HTTP_STATUS.UNAUTHORIZED:
          message.error('登录已过期，请重新登录');
          removeToken();
          // 跳转到登录页面
          window.location.href = '/user/login';
          break;
          
        case HTTP_STATUS.FORBIDDEN:
          message.error('没有权限访问该资源');
          break;
          
        case HTTP_STATUS.NOT_FOUND:
          message.error('请求的资源不存在');
          break;
          
        case HTTP_STATUS.UNPROCESSABLE_ENTITY:
          // 处理表单验证错误
          if (data?.detail && Array.isArray(data.detail)) {
            const errorMessages = data.detail.map((err: any) => err.msg).join(', ');
            message.error(`输入错误: ${errorMessages}`);
          } else {
            message.error(data?.detail || '请求参数错误');
          }
          break;
          
        case HTTP_STATUS.INTERNAL_SERVER_ERROR:
          message.error('服务器内部错误，请稍后重试');
          break;
          
        default:
          message.error(data?.detail || data?.message || '请求失败，请稍后重试');
      }
    } else if (error.request) {
      // 网络错误
      message.error('网络连接失败，请检查网络设置');
    } else {
      // 其他错误
      message.error('请求配置错误');
    }
    
    // 打印错误信息 (开发环境)
    if (process.env.NODE_ENV === 'development') {
      console.error('❌ API Error:', {
        message: error.message,
        status: error.response?.status,
        data: error.response?.data,
        config: error.config,
      });
    }
    
    return Promise.reject(error);
  },
);

// 封装常用的请求方法
export const api = {
  get: <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> => {
    return request.get(url, config).then(res => res.data);
  },
  
  post: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
    return request.post(url, data, config).then(res => res.data);
  },
  
  put: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
    return request.put(url, data, config).then(res => res.data);
  },
  
  delete: <T = any>(url: string, config?: AxiosRequestConfig): Promise<T> => {
    return request.delete(url, config).then(res => res.data);
  },
  
  patch: <T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<T> => {
    return request.patch(url, data, config).then(res => res.data);
  },
};

export default request;