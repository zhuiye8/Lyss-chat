# Lyss-chat 2.0 API 设计规范

## 目录

1. [设计原则](#设计原则)
2. [API 风格](#api-风格)
3. [URI 设计](#uri-设计)
4. [请求方法](#请求方法)
5. [请求参数](#请求参数)
6. [响应格式](#响应格式)
7. [状态码](#状态码)
8. [错误处理](#错误处理)
9. [版本控制](#版本控制)
10. [认证与授权](#认证与授权)
11. [限流与缓存](#限流与缓存)
12. [API 文档](#api-文档)
13. [示例](#示例)

## 设计原则

Lyss-chat 2.0 的 API 设计遵循以下原则：

1. **一致性**：所有 API 遵循统一的命名规范和交互模式
2. **简单性**：API 设计简单明了，易于理解和使用
3. **可预测性**：相似的资源有相似的 URI 和操作方式
4. **可扩展性**：API 设计支持未来功能扩展，不破坏现有客户端
5. **安全性**：API 设计考虑安全因素，防止常见的安全漏洞
6. **性能**：API 设计考虑性能因素，减少不必要的请求和响应数据

## API 风格

Lyss-chat 2.0 采用 RESTful API 设计风格，具有以下特点：

1. **资源导向**：API 围绕资源设计，而非操作
2. **使用 HTTP 方法**：使用标准 HTTP 方法表示操作
3. **无状态**：服务器不存储客户端状态，每个请求包含所有必要信息
4. **使用 HTTP 状态码**：使用标准 HTTP 状态码表示请求结果
5. **使用 JSON**：请求和响应使用 JSON 格式

## URI 设计

### 基本规则

1. **使用名词表示资源**：URI 使用名词表示资源，而非动词
2. **使用复数形式**：资源集合使用复数形式
3. **使用小写字母**：URI 使用小写字母，单词之间使用连字符（-）分隔
4. **避免文件扩展名**：不使用 .json 等文件扩展名
5. **使用层级结构**：表示资源之间的层级关系

### 基础路径

所有 API 使用以下基础路径：

```
/api/v1
```

### 资源命名

主要资源的 URI 命名如下：

| 资源 | URI |
|------|-----|
| 租户 | `/api/v1/tenants` |
| 用户 | `/api/v1/users` |
| 角色 | `/api/v1/roles` |
| 权限 | `/api/v1/permissions` |
| 工作区 | `/api/v1/workspaces` |
| 画布 | `/api/v1/workspaces/{workspaceId}/canvases` |
| 消息 | `/api/v1/canvases/{canvasId}/messages` |
| 提供商 | `/api/v1/providers` |
| 模型 | `/api/v1/providers/{providerId}/models` |

### 子资源

子资源使用父资源的 URI 作为前缀：

```
/api/v1/workspaces/{workspaceId}/users
/api/v1/users/{userId}/roles
/api/v1/roles/{roleId}/permissions
```

### 查询参数

查询参数用于过滤、排序、分页等操作：

```
/api/v1/users?status=active
/api/v1/messages?canvasId={canvasId}&limit=20&offset=0
/api/v1/models?providerId={providerId}&sort=name
```

## 请求方法

Lyss-chat 2.0 使用标准 HTTP 方法表示操作：

| 方法 | 描述 | 示例 |
|------|------|------|
| GET | 获取资源 | `GET /api/v1/users` |
| POST | 创建资源 | `POST /api/v1/users` |
| PUT | 全量更新资源 | `PUT /api/v1/users/{userId}` |
| PATCH | 部分更新资源 | `PATCH /api/v1/users/{userId}` |
| DELETE | 删除资源 | `DELETE /api/v1/users/{userId}` |

### 方法语义

- **GET**：安全且幂等，不修改资源状态
- **POST**：非幂等，每次调用可能创建新资源
- **PUT**：幂等，多次调用结果相同
- **PATCH**：非幂等，部分更新资源
- **DELETE**：幂等，多次调用结果相同

## 请求参数

### 路径参数

路径参数用于标识特定资源：

```
/api/v1/users/{userId}
/api/v1/workspaces/{workspaceId}/canvases/{canvasId}
```

### 查询参数

查询参数用于过滤、排序、分页等操作：

| 参数 | 描述 | 示例 |
|------|------|------|
| filter | 过滤条件 | `?filter[status]=active` |
| sort | 排序字段 | `?sort=name` 或 `?sort=-created_at` |
| page | 页码 | `?page=2` |
| limit | 每页数量 | `?limit=20` |
| include | 包含关联资源 | `?include=roles,permissions` |
| fields | 指定返回字段 | `?fields=id,name,email` |

### 请求体

POST、PUT、PATCH 请求使用 JSON 格式的请求体：

```json
{
  "name": "用户名",
  "email": "user@example.com",
  "password": "password123",
  "role": "user"
}
```

## 响应格式

### 成功响应

成功响应使用统一的 JSON 格式：

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "用户名",
    "email": "user@example.com",
    "status": "active",
    "created_at": "2023-01-01T12:00:00Z",
    "updated_at": "2023-01-01T12:00:00Z"
  }
}
```

### 集合响应

集合响应包含分页信息：

```json
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "用户1",
      "email": "user1@example.com"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "用户2",
      "email": "user2@example.com"
    }
  ],
  "meta": {
    "total": 100,
    "page": 1,
    "limit": 20,
    "pages": 5
  },
  "links": {
    "self": "/api/v1/users?page=1&limit=20",
    "first": "/api/v1/users?page=1&limit=20",
    "prev": null,
    "next": "/api/v1/users?page=2&limit=20",
    "last": "/api/v1/users?page=5&limit=20"
  }
}
```

### 关联资源

当请求包含 `include` 参数时，响应中包含关联资源：

```json
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "name": "用户名",
    "email": "user@example.com"
  },
  "included": {
    "roles": [
      {
        "id": "550e8400-e29b-41d4-a716-446655440001",
        "name": "管理员"
      }
    ]
  }
}
```

## 状态码

Lyss-chat 2.0 使用标准 HTTP 状态码表示请求结果：

### 成功状态码

| 状态码 | 描述 | 使用场景 |
|--------|------|----------|
| 200 OK | 请求成功 | GET 请求成功返回资源 |
| 201 Created | 资源创建成功 | POST 请求成功创建资源 |
| 204 No Content | 请求成功，无返回内容 | DELETE 请求成功删除资源 |

### 重定向状态码

| 状态码 | 描述 | 使用场景 |
|--------|------|----------|
| 301 Moved Permanently | 资源永久移动到新位置 | 资源 URI 永久变更 |
| 302 Found | 资源临时移动到新位置 | 临时重定向 |
| 304 Not Modified | 资源未修改 | 使用缓存的资源 |

### 客户端错误状态码

| 状态码 | 描述 | 使用场景 |
|--------|------|----------|
| 400 Bad Request | 请求格式错误 | 请求参数格式错误或缺失 |
| 401 Unauthorized | 未认证 | 缺少认证信息或认证失败 |
| 403 Forbidden | 权限不足 | 已认证但权限不足 |
| 404 Not Found | 资源不存在 | 请求的资源不存在 |
| 409 Conflict | 资源冲突 | 资源已存在或版本冲突 |
| 422 Unprocessable Entity | 请求格式正确但语义错误 | 请求参数验证失败 |
| 429 Too Many Requests | 请求过多 | 超过请求频率限制 |

### 服务器错误状态码

| 状态码 | 描述 | 使用场景 |
|--------|------|----------|
| 500 Internal Server Error | 服务器内部错误 | 服务器发生未预期的错误 |
| 502 Bad Gateway | 网关错误 | 上游服务返回无效响应 |
| 503 Service Unavailable | 服务不可用 | 服务器暂时无法处理请求 |
| 504 Gateway Timeout | 网关超时 | 上游服务响应超时 |

## 错误处理

### 错误响应格式

错误响应使用统一的 JSON 格式：

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "请求参数验证失败",
    "details": [
      {
        "field": "email",
        "message": "邮箱格式不正确"
      },
      {
        "field": "password",
        "message": "密码长度必须大于 8 个字符"
      }
    ]
  }
}
```

### 错误码

系统定义了以下错误码：

| 错误码 | 描述 | HTTP 状态码 |
|--------|------|------------|
| AUTHENTICATION_ERROR | 认证失败 | 401 |
| AUTHORIZATION_ERROR | 权限不足 | 403 |
| RESOURCE_NOT_FOUND | 资源不存在 | 404 |
| VALIDATION_ERROR | 请求参数验证失败 | 422 |
| CONFLICT_ERROR | 资源冲突 | 409 |
| RATE_LIMIT_ERROR | 请求频率超限 | 429 |
| INTERNAL_ERROR | 服务器内部错误 | 500 |
| SERVICE_UNAVAILABLE | 服务不可用 | 503 |

### 错误处理最佳实践

1. **提供明确的错误信息**：错误信息应该清晰明确，帮助客户端理解错误原因
2. **包含错误详情**：对于验证错误，提供具体的字段错误信息
3. **使用适当的状态码**：根据错误类型使用适当的 HTTP 状态码
4. **不泄露敏感信息**：错误信息不应包含敏感信息，如密码、内部路径等
5. **记录错误日志**：服务器应记录错误日志，便于问题排查

## 版本控制

为确保 API 的向后兼容性和平滑升级，Lyss-chat 2.0 采用以下版本控制策略：

### URI 版本控制

在 URI 中包含版本号：

```
/api/v1/users
/api/v2/users
```

### 版本升级原则

1. **向后兼容**：新版本应尽量保持向后兼容，不破坏现有客户端
2. **渐进式变更**：重大变更应分阶段进行，给客户端足够的适应时间
3. **版本共存**：新旧版本可以共存一段时间，便于客户端平滑迁移
4. **版本废弃通知**：废弃旧版本前，应提前通知客户端

### 版本生命周期

1. **开发版本**：内部开发使用，不对外发布
2. **测试版本**：提供给测试人员和合作伙伴测试
3. **正式版本**：对所有客户端开放
4. **维护版本**：只修复关键 bug，不添加新功能
5. **废弃版本**：计划废弃，但仍然可用
6. **停用版本**：完全停止服务

## 认证与授权

### 认证机制

Lyss-chat 2.0 使用 JWT（JSON Web Token）进行认证：

1. **获取令牌**：客户端通过登录接口获取 JWT 令牌
2. **使用令牌**：客户端在请求头中携带令牌
3. **验证令牌**：服务器验证令牌的有效性和权限

### 认证流程

```
┌─────────┐                                           ┌─────────┐
│  客户端  │                                           │  服务器  │
└────┬────┘                                           └────┬────┘
     │                                                     │
     │  POST /api/v1/auth/login                           │
     │  {email: "user@example.com", password: "******"}   │
     │ ──────────────────────────────────────────────────>│
     │                                                     │
     │  200 OK                                             │
     │  {access_token: "xxx", refresh_token: "yyy"}        │
     │ <──────────────────────────────────────────────────│
     │                                                     │
     │  GET /api/v1/users                                  │
     │  Authorization: Bearer xxx                          │
     │ ──────────────────────────────────────────────────>│
     │                                                     │
     │  200 OK                                             │
     │  {data: [...]}                                      │
     │ <──────────────────────────────────────────────────│
     │                                                     │
```

### 令牌刷新

为提高安全性，访问令牌（access token）有较短的有效期，刷新令牌（refresh token）有较长的有效期：

```
┌─────────┐                                           ┌─────────┐
│  客户端  │                                           │  服务器  │
└────┬────┘                                           └────┬────┘
     │                                                     │
     │  POST /api/v1/auth/refresh                          │
     │  {refresh_token: "yyy"}                             │
     │ ──────────────────────────────────────────────────>│
     │                                                     │
     │  200 OK                                             │
     │  {access_token: "zzz", refresh_token: "www"}        │
     │ <──────────────────────────────────────────────────│
     │                                                     │
```

### 授权机制

Lyss-chat 2.0 使用基于角色的访问控制（RBAC）进行授权：

1. **角色**：用户被分配一个或多个角色
2. **权限**：每个角色拥有一组权限
3. **资源**：系统中的各种资源（用户、工作区、画布等）
4. **操作**：对资源的各种操作（读取、创建、更新、删除等）

### 权限检查

权限检查在以下层面进行：

1. **API 层**：检查用户是否有权限访问特定 API
2. **服务层**：检查用户是否有权限执行特定操作
3. **数据层**：检查用户是否有权限访问特定数据

## 限流与缓存

### 限流策略

为防止滥用和保护系统资源，Lyss-chat 2.0 实施以下限流策略：

1. **基于 IP 的限流**：限制每个 IP 的请求频率
2. **基于用户的限流**：限制每个用户的请求频率
3. **基于资源的限流**：限制特定资源的访问频率
4. **基于操作的限流**：限制特定操作的执行频率

### 限流响应

当请求超过限制时，返回 429 Too Many Requests 状态码，并在响应头中包含限流信息：

```
HTTP/1.1 429 Too Many Requests
Retry-After: 60
X-RateLimit-Limit: 100
X-RateLimit-Remaining: 0
X-RateLimit-Reset: 1609459200

{
  "error": {
    "code": "RATE_LIMIT_ERROR",
    "message": "请求频率超限，请 60 秒后重试"
  }
}
```

### 缓存策略

为提高性能和减少服务器负载，Lyss-chat 2.0 实施以下缓存策略：

1. **客户端缓存**：使用 HTTP 缓存机制，通过 Cache-Control 和 ETag 头控制
2. **服务器缓存**：使用 Redis 缓存频繁访问的数据
3. **CDN 缓存**：使用 CDN 缓存静态资源

### 缓存控制

通过 HTTP 头控制缓存行为：

```
Cache-Control: max-age=3600, public
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"
Last-Modified: Wed, 21 Oct 2015 07:28:00 GMT
```

### 缓存失效

通过以下机制确保缓存数据的一致性：

1. **基于时间的失效**：设置缓存的过期时间
2. **基于事件的失效**：当数据变更时主动使缓存失效
3. **版本控制**：使用版本号或哈希值标识缓存数据

## API 文档

### 文档工具

Lyss-chat 2.0 使用 OpenAPI（Swagger）规范编写 API 文档，提供以下功能：

1. **API 描述**：详细描述 API 的功能、参数和响应
2. **交互式文档**：提供交互式界面，可以直接测试 API
3. **代码生成**：根据 API 文档生成客户端代码
4. **自动化测试**：根据 API 文档生成测试用例

### 文档结构

API 文档包含以下内容：

1. **API 概述**：API 的基本信息和使用说明
2. **认证方式**：API 的认证机制和使用方法
3. **资源定义**：API 涉及的资源和数据模型
4. **接口定义**：API 的具体接口、参数和响应
5. **错误码**：API 可能返回的错误码和处理方法
6. **示例代码**：API 的使用示例

### 文档示例

```yaml
openapi: 3.0.0
info:
  title: Lyss-chat API
  description: Lyss-chat 2.0 API 文档
  version: 1.0.0
servers:
  - url: https://api.lyss-chat.com/api/v1
    description: 生产环境
  - url: https://staging.lyss-chat.com/api/v1
    description: 测试环境
paths:
  /users:
    get:
      summary: 获取用户列表
      description: 获取符合条件的用户列表
      parameters:
        - name: status
          in: query
          description: 用户状态
          schema:
            type: string
            enum: [active, inactive]
        - name: page
          in: query
          description: 页码
          schema:
            type: integer
            default: 1
        - name: limit
          in: query
          description: 每页数量
          schema:
            type: integer
            default: 20
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserList'
        '401':
          $ref: '#/components/responses/Unauthorized'
        '403':
          $ref: '#/components/responses/Forbidden'
```

## 示例

### 用户认证

#### 登录

请求：

```http
POST /api/v1/auth/login HTTP/1.1
Host: api.lyss-chat.com
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

响应：

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600
  }
}
```

### 用户管理

#### 获取用户列表

请求：

```http
GET /api/v1/users?status=active&page=1&limit=20 HTTP/1.1
Host: api.lyss-chat.com
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

响应：

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "name": "用户1",
      "email": "user1@example.com",
      "status": "active",
      "created_at": "2023-01-01T12:00:00Z"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440001",
      "name": "用户2",
      "email": "user2@example.com",
      "status": "active",
      "created_at": "2023-01-02T12:00:00Z"
    }
  ],
  "meta": {
    "total": 100,
    "page": 1,
    "limit": 20,
    "pages": 5
  },
  "links": {
    "self": "/api/v1/users?page=1&limit=20",
    "first": "/api/v1/users?page=1&limit=20",
    "prev": null,
    "next": "/api/v1/users?page=2&limit=20",
    "last": "/api/v1/users?page=5&limit=20"
  }
}
```

#### 创建用户

请求：

```http
POST /api/v1/users HTTP/1.1
Host: api.lyss-chat.com
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "name": "新用户",
  "email": "newuser@example.com",
  "password": "password123",
  "role_ids": ["550e8400-e29b-41d4-a716-446655440002"]
}
```

响应：

```http
HTTP/1.1 201 Created
Content-Type: application/json
Location: /api/v1/users/550e8400-e29b-41d4-a716-446655440003

{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440003",
    "name": "新用户",
    "email": "newuser@example.com",
    "status": "active",
    "created_at": "2023-01-03T12:00:00Z",
    "updated_at": "2023-01-03T12:00:00Z"
  }
}
```

### 模型管理

#### 获取模型列表

请求：

```http
GET /api/v1/providers/550e8400-e29b-41d4-a716-446655440004/models HTTP/1.1
Host: api.lyss-chat.com
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

响应：

```http
HTTP/1.1 200 OK
Content-Type: application/json

{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440005",
      "provider_id": "550e8400-e29b-41d4-a716-446655440004",
      "model_id": "gpt-4",
      "name": "GPT-4",
      "description": "OpenAI 的 GPT-4 模型",
      "capabilities": {
        "chat": true,
        "image": false,
        "audio": false
      },
      "status": "active"
    },
    {
      "id": "550e8400-e29b-41d4-a716-446655440006",
      "provider_id": "550e8400-e29b-41d4-a716-446655440004",
      "model_id": "gpt-3.5-turbo",
      "name": "GPT-3.5 Turbo",
      "description": "OpenAI 的 GPT-3.5 Turbo 模型",
      "capabilities": {
        "chat": true,
        "image": false,
        "audio": false
      },
      "status": "active"
    }
  ],
  "meta": {
    "total": 2,
    "page": 1,
    "limit": 20,
    "pages": 1
  }
}
```

### 聊天功能

#### 发送消息

请求：

```http
POST /api/v1/canvases/550e8400-e29b-41d4-a716-446655440007/messages HTTP/1.1
Host: api.lyss-chat.com
Content-Type: application/json
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

{
  "role": "user",
  "content": "你好，请介绍一下自己。"
}
```

响应：

```http
HTTP/1.1 201 Created
Content-Type: application/json
Location: /api/v1/messages/550e8400-e29b-41d4-a716-446655440008

{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440008",
    "canvas_id": "550e8400-e29b-41d4-a716-446655440007",
    "role": "user",
    "content": "你好，请介绍一下自己。",
    "created_by": "550e8400-e29b-41d4-a716-446655440000",
    "created_at": "2023-01-03T12:00:00Z"
  }
}
```

#### 获取消息流

请求：

```http
GET /api/v1/canvases/550e8400-e29b-41d4-a716-446655440007/messages/stream HTTP/1.1
Host: api.lyss-chat.com
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

响应：

```http
HTTP/1.1 200 OK
Content-Type: text/event-stream
Transfer-Encoding: chunked

event: message
data: {"id":"550e8400-e29b-41d4-a716-446655440009","role":"assistant","content":"你好","created_at":"2023-01-03T12:00:01Z"}

event: message
data: {"id":"550e8400-e29b-41d4-a716-446655440009","role":"assistant","content":"你好，我是","created_at":"2023-01-03T12:00:01Z"}

event: message
data: {"id":"550e8400-e29b-41d4-a716-446655440009","role":"assistant","content":"你好，我是 Lyss","created_at":"2023-01-03T12:00:01Z"}

event: message
data: {"id":"550e8400-e29b-41d4-a716-446655440009","role":"assistant","content":"你好，我是 Lyss，一个 AI 助手。","created_at":"2023-01-03T12:00:01Z"}

event: done
data: {"id":"550e8400-e29b-41d4-a716-446655440009","token_count":20}
```
```
