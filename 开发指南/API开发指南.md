# Lyss-chat 2.0 API 开发指南

## 概述

本文档提供 Lyss-chat 2.0 系统 API 开发的详细指南，包括 API 设计原则、请求/响应格式、错误处理、认证授权、限流缓存等内容，确保 API 的一致性、可用性和安全性。

## 目标

1. 建立统一的 API 设计规范
2. 确保 API 的安全性和可靠性
3. 提供清晰的 API 文档和示例
4. 实现高效的 API 性能和可扩展性

## 1. API 设计原则

### 1.1 RESTful 设计

- 使用资源导向设计，而非操作导向
- 使用 HTTP 方法表示操作（GET、POST、PUT、DELETE 等）
- 使用 HTTP 状态码表示请求结果
- 使用 JSON 作为数据交换格式
- 保持 API 的无状态性

### 1.2 URI 设计规范

- 使用名词复数形式表示资源集合
- 使用小写字母和连字符（-）
- 避免使用文件扩展名
- 使用层级结构表示资源关系

```
# 基础路径
/api/v1

# 资源路径示例
/api/v1/users                # 用户集合
/api/v1/users/{id}           # 特定用户
/api/v1/users/{id}/roles     # 用户的角色集合
/api/v1/conversations        # 对话集合
/api/v1/conversations/{id}/messages  # 对话的消息集合
```

### 1.3 HTTP 方法使用

| 方法 | 用途 | 示例 |
|------|------|------|
| GET | 获取资源 | GET /api/v1/users |
| POST | 创建资源 | POST /api/v1/users |
| PUT | 全量更新资源 | PUT /api/v1/users/{id} |
| PATCH | 部分更新资源 | PATCH /api/v1/users/{id} |
| DELETE | 删除资源 | DELETE /api/v1/users/{id} |

## 2. 请求与响应格式

### 2.1 请求格式

#### 请求头

```
Content-Type: application/json
Authorization: Bearer {token}
Accept-Language: zh-CN
```

#### 请求参数

- 路径参数：在 URI 中使用 `{parameter}` 表示
- 查询参数：使用 `?key=value&key2=value2` 格式
- 请求体：使用 JSON 格式

```json
// POST /api/v1/users
{
  "email": "user@example.com",
  "name": "张三",
  "password": "SecurePassword123",
  "role": "user"
}
```

### 2.2 响应格式

#### 成功响应

```json
// 200 OK - 单个资源
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "张三",
    "role": "user",
    "created_at": "2023-06-01T12:00:00Z",
    "updated_at": "2023-06-01T12:00:00Z"
  }
}

// 200 OK - 资源集合
{
  "data": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "user@example.com",
      "name": "张三",
      "role": "user"
    },
    // ...更多资源
  ],
  "meta": {
    "total": 100,
    "page": 1,
    "per_page": 20,
    "total_pages": 5
  },
  "links": {
    "self": "/api/v1/users?page=1&per_page=20",
    "first": "/api/v1/users?page=1&per_page=20",
    "prev": null,
    "next": "/api/v1/users?page=2&per_page=20",
    "last": "/api/v1/users?page=5&per_page=20"
  }
}

// 201 Created
{
  "data": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com",
    "name": "张三",
    "role": "user",
    "created_at": "2023-06-01T12:00:00Z",
    "updated_at": "2023-06-01T12:00:00Z"
  }
}

// 204 No Content
// 无响应体
```

#### 错误响应

```json
// 400 Bad Request, 422 Unprocessable Entity, 等
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "请求参数验证失败",
    "details": [
      {
        "field": "email",
        "message": "无效的邮箱格式"
      },
      {
        "field": "password",
        "message": "密码长度必须至少为8个字符"
      }
    ]
  }
}

// 401 Unauthorized
{
  "error": {
    "code": "UNAUTHORIZED",
    "message": "未授权访问"
  }
}

// 403 Forbidden
{
  "error": {
    "code": "FORBIDDEN",
    "message": "没有权限执行此操作"
  }
}

// 404 Not Found
{
  "error": {
    "code": "RESOURCE_NOT_FOUND",
    "message": "请求的资源不存在"
  }
}

// 500 Internal Server Error
{
  "error": {
    "code": "INTERNAL_ERROR",
    "message": "服务器内部错误"
  }
}
```

## 3. 错误处理

### 3.1 错误码定义

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

### 3.2 错误处理最佳实践

- 提供明确的错误信息
- 包含错误详情（对于验证错误）
- 使用适当的 HTTP 状态码
- 不泄露敏感信息
- 记录错误日志

## 4. 认证与授权

### 4.1 认证机制

使用 JWT（JSON Web Token）进行认证：

1. 客户端通过登录接口获取令牌
2. 客户端在请求头中携带令牌
3. 服务器验证令牌的有效性

```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 4.2 认证 API

```
# 登录
POST /api/v1/auth/login
{
  "email": "user@example.com",
  "password": "SecurePassword123"
}

# 响应
{
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": 3600
  }
}

# 刷新令牌
POST /api/v1/auth/refresh
{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}

# 注销
POST /api/v1/auth/logout
```

### 4.3 授权机制

使用基于角色的访问控制（RBAC）进行授权：

1. 用户被分配一个或多个角色
2. 每个角色拥有一组权限
3. 权限定义为资源和操作的组合

## 5. 分页、排序和过滤

### 5.1 分页

使用 `page` 和 `per_page` 参数进行分页：

```
GET /api/v1/users?page=2&per_page=20
```

### 5.2 排序

使用 `sort` 参数进行排序，前缀 `-` 表示降序：

```
GET /api/v1/users?sort=name         # 按名称升序
GET /api/v1/users?sort=-created_at  # 按创建时间降序
GET /api/v1/users?sort=role,-name   # 先按角色升序，再按名称降序
```

### 5.3 过滤

使用查询参数进行过滤：

```
GET /api/v1/users?role=admin        # 筛选管理员
GET /api/v1/users?status=active     # 筛选活跃用户
GET /api/v1/users?created_at_gte=2023-01-01T00:00:00Z  # 筛选特定日期之后创建的用户
```

## 6. 限流与缓存

### 6.1 限流策略

为防止滥用和保护系统资源，实施以下限流策略：

- 基于 IP 的限流
- 基于用户的限流
- 基于资源的限流
- 基于操作的限流

限流响应：

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

### 6.2 缓存策略

使用 HTTP 缓存机制提高性能：

```
Cache-Control: max-age=3600, public
ETag: "33a64df551425fcc55e4d42a148795d9f25f89d4"
Last-Modified: Wed, 21 Oct 2023 07:28:00 GMT
```

## 7. API 版本控制

使用 URI 中的版本号进行版本控制：

```
/api/v1/users
/api/v2/users
```

版本升级原则：

- 保持向后兼容性
- 渐进式变更
- 版本共存
- 版本废弃通知

## 8. API 文档

使用 OpenAPI（Swagger）规范编写 API 文档：

```yaml
openapi: 3.0.0
info:
  title: Lyss-chat API
  description: Lyss-chat 2.0 API 文档
  version: 1.0.0
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
      responses:
        '200':
          description: 成功
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/UserList'
```

## 9. API 实现示例

### 9.1 用户认证 API

```go
// 登录处理器
func (h *AuthHandler) Login(c *web.Context) {
    // 解析请求
    var req LoginRequest
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, api.ErrorResponse{
            Error: api.Error{
                Code:    "VALIDATION_ERROR",
                Message: "无效的请求格式",
            },
        })
        return
    }
    
    // 验证用户凭据
    tokenPair, err := h.authService.Login(c.Request.Context(), req.TenantID, req.Email, req.Password)
    if err != nil {
        if errors.Is(err, service.ErrInvalidCredentials) {
            c.JSON(http.StatusUnauthorized, api.ErrorResponse{
                Error: api.Error{
                    Code:    "AUTHENTICATION_ERROR",
                    Message: "邮箱或密码不正确",
                },
            })
            return
        }
        
        c.JSON(http.StatusInternalServerError, api.ErrorResponse{
            Error: api.Error{
                Code:    "INTERNAL_ERROR",
                Message: "登录处理失败",
            },
        })
        return
    }
    
    // 返回令牌
    c.JSON(http.StatusOK, api.Response{
        Data: tokenPair,
    })
}
```
