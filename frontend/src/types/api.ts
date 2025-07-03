/**
 * LYSS AI Platform API 类型定义
 * 与后端API接口保持一致的TypeScript类型
 */

// ===== 基础类型 =====
export interface ApiResponse<T = any> {
  success: boolean;
  data?: T;
  message?: string;
  code?: number;
}

export interface PaginationParams {
  page?: number;
  pageSize?: number;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
}

// ===== 用户相关 =====
export interface User {
  id: string;
  email: string;
  first_name?: string;
  last_name?: string;
  role: 'user' | 'admin' | 'super_admin';
  is_active: boolean;
  is_verified: boolean;
  created_at: string;
  updated_at: string;
}

export interface LoginRequest {
  username: string; // 实际使用email
  password: string;
}

export interface RegisterRequest {
  email: string;
  password: string;
  first_name?: string;
  last_name?: string;
}

export interface LoginResponse {
  access_token: string;
  token_type: string;
}

// ===== 供应商相关 =====
export type ProviderScope = 'ORGANIZATION' | 'PERSONAL';

export interface Provider {
  id: string;
  name: string;
  provider_type: string;
  scope: ProviderScope;
  description?: string;
  is_active: boolean;
  owner_id: string;
  created_at: string;
  updated_at: string;
}

export interface ProviderCreateRequest {
  name: string;
  provider_type: string;
  scope: ProviderScope;
  description?: string;
  config: Record<string, any>;
}

export interface ProviderTestRequest {
  provider_type: string;
  config: Record<string, any>;
}

export interface ProviderTestResponse {
  success: boolean;
  message: string;
  error_details?: string;
}

export interface ProviderConfigSchema {
  provider_type: string;
  config_schema: Record<string, any>;
}

// ===== 模型相关 =====
export interface Model {
  id: string;
  provider_id: string;
  name: string;
  display_name: string;
  description?: string;
  model_type: string;
  input_cost_per_token: number;
  output_cost_per_token: number;
  max_tokens?: number;
  supports_streaming: boolean;
  supports_functions: boolean;
  context_window: number;
  is_active: boolean;
  created_at: string;
  updated_at: string;
}

// ===== 对话相关 =====
export interface ChatMessage {
  role: 'user' | 'assistant' | 'system';
  content: string;
}

export interface ChatRequest {
  messages: ChatMessage[];
  model_id: string;
  stream?: boolean;
  temperature?: number;
  max_tokens?: number;
  top_p?: number;
}

export interface ChatResponse {
  content: string;
  finish_reason?: string;
  usage?: {
    prompt_tokens: number;
    completion_tokens: number;
    total_tokens: number;
  };
}

// ===== 文件相关 =====
export interface UploadedFile {
  id: string;
  user_id: string;
  original_filename: string;
  file_path: string;
  file_size: number;
  mime_type: string;
  status: 'processing' | 'completed' | 'failed';
  processing_error?: string;
  total_chunks: number;
  qdrant_collection?: string;
  created_at: string;
  updated_at: string;
}

// ===== 使用统计相关 =====
export interface UsageLog {
  id: string;
  user_id: string;
  model_id: string;
  prompt_tokens: number;
  completion_tokens: number;
  total_tokens: number;
  input_cost: number;
  output_cost: number;
  total_cost: number;
  created_at: string;
}

export interface UsageStats {
  total_requests: number;
  total_tokens: number;
  total_cost: number;
  period: string; // 'today' | 'week' | 'month'
}

// ===== 用户模型访问权限 =====
export interface UserModelAccess {
  id: string;
  user_id: string;
  model_id: string;
  is_active: boolean;
  expires_at?: string;
  daily_quota?: number;
  monthly_quota?: number;
  granted_by: string;
  created_at: string;
  updated_at: string;
}