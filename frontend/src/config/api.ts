/**
 * API配置文件
 * 包含后端服务地址、接口路径等配置信息
 */

// 后端服务地址配置
export const API_CONFIG = {
  // 基础URL
  BASE_URL: process.env.NODE_ENV === 'development' 
    ? 'http://localhost:8000' 
    : process.env.REACT_APP_API_URL || 'http://localhost:8000',
  
  // API版本前缀
  API_PREFIX: '/api/v1',
  
  // 请求超时时间 (毫秒)
  TIMEOUT: 30000,
};

// API接口路径常量
export const API_ENDPOINTS = {
  // 认证相关
  AUTH: {
    LOGIN: '/auth/jwt/login',
    LOGOUT: '/auth/jwt/logout', 
    REGISTER: '/auth/register',
    CURRENT_USER: '/auth/users/me',
    REFRESH: '/auth/jwt/refresh',
  },
  
  // 用户管理
  USERS: {
    LIST: '/users',
    DETAIL: (id: string) => `/users/${id}`,
    UPDATE: (id: string) => `/users/${id}`,
    DELETE: (id: string) => `/users/${id}`,
  },
  
  // 供应商管理
  PROVIDERS: {
    LIST: '/providers',
    CREATE: '/providers',
    DETAIL: (id: string) => `/providers/${id}`,
    UPDATE: (id: string) => `/providers/${id}`,
    DELETE: (id: string) => `/providers/${id}`,
    TEST: '/providers/test',
    CONFIG_SCHEMA: (type: string) => `/providers/types/${type}/config-schema`,
    TYPES: '/providers/types',
  },
  
  // 模型管理
  MODELS: {
    LIST: '/models',
    CREATE: '/models',
    DETAIL: (id: string) => `/models/${id}`,
    UPDATE: (id: string) => `/models/${id}`,
    DELETE: (id: string) => `/models/${id}`,
    AVAILABLE: '/models/available',
  },
  
  // 对话功能
  CHAT: {
    COMPLETION: '/chat/completions',
    MODELS: '/chat/models',
  },
  
  // 文件管理
  FILES: {
    LIST: '/files',
    UPLOAD: '/files/upload',
    DETAIL: (id: string) => `/files/${id}`,
    DELETE: (id: string) => `/files/${id}`,
    DOWNLOAD: (id: string) => `/files/${id}/download`,
  },
  
  // 使用统计
  USAGE: {
    STATS: '/usage/stats',
    LOGS: '/usage/logs',
    EXPORT: '/usage/export',
  },
  
  // 用户模型访问权限
  ACCESS: {
    LIST: '/access',
    GRANT: '/access/grant',
    REVOKE: (id: string) => `/access/${id}/revoke`,
    CHECK: (modelId: string) => `/access/check/${modelId}`,
  },
};

// 构建完整的API URL
export const buildApiUrl = (endpoint: string): string => {
  return `${API_CONFIG.BASE_URL}${API_CONFIG.API_PREFIX}${endpoint}`;
};

// 常用的HTTP状态码
export const HTTP_STATUS = {
  OK: 200,
  CREATED: 201,
  NO_CONTENT: 204,
  BAD_REQUEST: 400,
  UNAUTHORIZED: 401,
  FORBIDDEN: 403,
  NOT_FOUND: 404,
  UNPROCESSABLE_ENTITY: 422,
  INTERNAL_SERVER_ERROR: 500,
} as const;

// 请求方法类型
export type HttpMethod = 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';

// 请求配置类型
export interface RequestConfig {
  method?: HttpMethod;
  headers?: Record<string, string>;
  params?: Record<string, any>;
  data?: any;
  timeout?: number;
}