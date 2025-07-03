/**
 * React Query 配置
 * 统一管理服务端状态，包含查询缓存、重试策略等配置
 */

import { QueryClient } from '@tanstack/react-query';
import { message } from 'antd';

// 创建QueryClient实例
export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      // 数据保持新鲜的时间 (5分钟)
      staleTime: 5 * 60 * 1000,
      
      // 缓存时间 (10分钟)
      cacheTime: 10 * 60 * 1000,
      
      // 重试次数
      retry: (failureCount, error: any) => {
        // 4xx错误不重试
        if (error?.response?.status >= 400 && error?.response?.status < 500) {
          return false;
        }
        // 最多重试2次
        return failureCount < 2;
      },
      
      // 重试延迟 (指数退避)
      retryDelay: (attemptIndex) => Math.min(1000 * 2 ** attemptIndex, 30000),
      
      // 窗口重新获得焦点时重新获取数据
      refetchOnWindowFocus: false,
      
      // 网络重连时重新获取数据  
      refetchOnReconnect: true,
      
      // 全局错误处理
      onError: (error: any) => {
        // React Query会自动处理，这里只做日志记录
        if (process.env.NODE_ENV === 'development') {
          console.error('Query Error:', error);
        }
      },
    },
    
    mutations: {
      // 变更重试次数
      retry: 1,
      
      // 全局成功处理
      onSuccess: () => {
        // 大部分成功情况由具体组件处理
      },
      
      // 全局错误处理
      onError: (error: any) => {
        // 一般错误已经被axios拦截器处理，这里处理特殊情况
        if (process.env.NODE_ENV === 'development') {
          console.error('Mutation Error:', error);
        }
      },
    },
  },
});

// 查询键工厂，统一管理查询键命名
export const queryKeys = {
  // 用户相关
  users: {
    all: ['users'] as const,
    lists: () => [...queryKeys.users.all, 'list'] as const,
    list: (filters: Record<string, any>) => [...queryKeys.users.lists(), filters] as const,
    details: () => [...queryKeys.users.all, 'detail'] as const,
    detail: (id: string) => [...queryKeys.users.details(), id] as const,
    current: () => [...queryKeys.users.all, 'current'] as const,
  },
  
  // 供应商相关
  providers: {
    all: ['providers'] as const,
    lists: () => [...queryKeys.providers.all, 'list'] as const,
    list: (filters: Record<string, any>) => [...queryKeys.providers.lists(), filters] as const,
    details: () => [...queryKeys.providers.all, 'detail'] as const,
    detail: (id: string) => [...queryKeys.providers.details(), id] as const,
    types: () => [...queryKeys.providers.all, 'types'] as const,
    configSchema: (type: string) => [...queryKeys.providers.all, 'config-schema', type] as const,
  },
  
  // 模型相关
  models: {
    all: ['models'] as const,
    lists: () => [...queryKeys.models.all, 'list'] as const,
    list: (filters: Record<string, any>) => [...queryKeys.models.lists(), filters] as const,
    details: () => [...queryKeys.models.all, 'detail'] as const,
    detail: (id: string) => [...queryKeys.models.details(), id] as const,
    available: () => [...queryKeys.models.all, 'available'] as const,
  },
  
  // 文件相关
  files: {
    all: ['files'] as const,
    lists: () => [...queryKeys.files.all, 'list'] as const,
    list: (filters: Record<string, any>) => [...queryKeys.files.lists(), filters] as const,
    details: () => [...queryKeys.files.all, 'detail'] as const,
    detail: (id: string) => [...queryKeys.files.details(), id] as const,
  },
  
  // 使用统计相关
  usage: {
    all: ['usage'] as const,
    stats: (period: string) => [...queryKeys.usage.all, 'stats', period] as const,
    logs: (filters: Record<string, any>) => [...queryKeys.usage.all, 'logs', filters] as const,
  },
  
  // 访问权限相关
  access: {
    all: ['access'] as const,
    lists: () => [...queryKeys.access.all, 'list'] as const,
    list: (filters: Record<string, any>) => [...queryKeys.access.lists(), filters] as const,
    check: (modelId: string) => [...queryKeys.access.all, 'check', modelId] as const,
  },
};