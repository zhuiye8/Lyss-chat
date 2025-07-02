# LYSS AI 平台 - 前端与后端API交互规范

**版本**: 1.0
**最后更新**: 2025年7月2日

---

## 1. 概述

为了确保前端与后端之间通信的稳定性、一致性和可维护性，我们必须建立一套标准的API交互规范。本文档定义了前端应用中API请求的完整生命周期管理，包括服务层结构、请求封装、响应处理和错误管理。

## 2. 服务层 (Service Layer) 结构

所有与后端API的交互逻辑都必须封装在 **服务层** 中。这有助于将数据获取逻辑与UI组件分离，提高代码的可重用性和可测试性。

*   **目录结构**: `src/services/lyss/`
*   **文件组织**: 按业务领域（或API的Controller/Router）组织文件。

```
src/services/lyss/
├── userAPI.ts         # 用户相关接口 (获取当前用户, 列表等)
├── providerAPI.ts     # 供应商管理接口
├── modelAPI.ts        # 模型管理与分发接口
├── chatAPI.ts         # 对话接口
└── typings.d.ts       # API相关的TypeScript类型定义
```

### 2.1. 类型定义 (`typings.d.ts`)

所有API请求参数和响应数据的TypeScript类型都应集中定义在此文件中，以便于管理和复用。

```typescript
// src/services/lyss/typings.d.ts

declare namespace API {
  // --- Provider ---
  type Provider = {
    id: string;
    name: string;
    provider_type: 'openai' | 'anthropic';
    // ... other fields
  };

  type ProviderCreateParams = {
    name: string;
    // ...
  };

  // --- Standard Response ---
  type StandardResponse<T> = {
    code: number;
    message: string;
    data: T;
  };
}
```

### 2.2. API函数示例 (`providerAPI.ts`)

每个API函数都应是导出的异步函数，负责处理一个具体的API端点调用。

```typescript
// src/services/lyss/providerAPI.ts
import { request } from '@umijs/max'; // Ant Design Pro 封装的请求库

/** 获取供应商列表 GET /api/v1/providers */
export async function getProviders(options?: { [key: string]: any }) {
  return request<API.StandardResponse<API.Provider[]>>('/api/v1/providers', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 创建供应商 POST /api/v1/providers */
export async function createProvider(body: API.ProviderCreateParams, options?: { [key: string]: any }) {
  return request<API.StandardResponse<API.Provider>>('/api/v1/providers', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}
```

## 3. 请求封装与配置

我们使用 Ant Design Pro 内置的 `request` 函数（基于 `umi-request` 或 `axios`），它提供了统一的请求和响应拦截器、错误处理等功能。

### 3.1. 全局请求配置 (`src/requestErrorConfig.ts`)

这是配置全局请求行为的地方。

*   **请求拦截器**:
    *   自动为每个请求（除了登录/注册）添加 `Authorization` 头，值为存储在 `localStorage` 中的 JWT Token。
*   **响应拦截器**:
    *   **业务成功判断**: 检查后端返回的 `StandardResponse` 中的 `code` 字段。如果 `code !== 0`，即使HTTP状态码是200，也应将其视为一个业务错误，并 `Promise.reject`。
    *   **数据提取**: 如果 `code === 0`，直接返回 `response.data` 部分，这样业务代码（如React Query的 `queryFn`）就可以直接拿到所需的数据，无需每次都解包。
    *   **未登录处理**: 如果API返回特定的HTTP状态码（如401 Unauthorized），自动重定向到登录页面。

### 3.2. 配置示例

```typescript
// src/requestErrorConfig.ts
import type { RequestConfig } from '@umijs/max';
import { message, notification } from 'antd';

// 错误处理方案： 错误类型
enum ErrorShowType {
  SILENT = 0,
  WARN_MESSAGE = 1,
  ERROR_MESSAGE = 2,
  NOTIFICATION = 3,
  REDIRECT = 9,
}

export const errorConfig: RequestConfig = {
  // 请求拦截器
  requestInterceptors: [
    (config: RequestOptions) => {
      const token = localStorage.getItem('auth_token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    },
  ],

  // 响应拦截器
  responseInterceptors: [
    (response) => {
      const { data } = response as unknown as { data: API.StandardResponse<any> };
      
      // 后端自定义业务错误
      if (data?.code !== 0) {
        message.error(data.message || 'Request error, please try again.');
        // 抛出错误，以便 useQuery/useMutation 的 onError 回调可以捕获
        throw new Error(data.message);
      }
      
      // 直接返回核心数据，简化业务代码
      response.data = data.data;
      return response;
    },
  ],
  
  // ... 其他错误处理配置
};
```

## 4. 流式API请求 (`chatAPI.ts`)

对于流式API（如聊天），我们不能使用默认的 `request` 配置，因为它会等待完整的JSON响应。我们需要使用原生的 `fetch` API 来处理 `ReadableStream`。

```typescript
// src/services/lyss/chatAPI.ts

export async function chatCompletion(body: API.ChatCompletionParams) {
  const token = localStorage.getItem('auth_token');
  
  // 直接使用 fetch API
  const response = await fetch('/api/v1/chat/completions', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${token}`,
    },
    body: JSON.stringify({ ...body, stream: true }),
  });

  if (!response.ok) {
    // 处理错误情况
    const errorData = await response.json();
    throw new Error(errorData.message || 'Chat API request failed');
  }

  // 将 Response 对象直接返回给 @ant-design/pro-chat
  // ProChat 内部会处理 ReadableStream
  return response;
}
```

通过这套规范，我们确保了前端API交互的健壮性和一致性，为构建稳定可靠的用户界面打下了坚实的基础。
