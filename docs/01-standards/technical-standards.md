# Lyss AI Platform 统一技术规范文档

**版本**: 2.0  
**状态**: 已确认  
**最后更新**: 2025-01-25  
**适用范围**: 全平台所有微服务

## 文档概述

本文档基于2024-2025年最新行业标准制定，是Lyss AI Platform所有开发工作的技术基石。所有微服务、API设计、数据库操作和代码编写必须严格遵循本规范。

---

## 1. RESTful API设计规范

### 1.1 URL设计标准

#### 基础结构
```
https://api.lyss.ai/api/v1/{resource}
```

#### 命名约定
- **集合资源使用复数名词**: `/users`, `/groups`, `/credentials`
- **层级资源关系**: `/groups/{group_id}/members/{user_id}`
- **统一小写格式**: 避免驼峰式命名
- **使用kebab-case**: `/user-profiles`, `/billing-records`
- **禁用动词**: 使用 `/orders` 而非 `/create-order`

#### URL结构示例
```
GET    /api/v1/users                    # 获取用户列表
POST   /api/v1/users                    # 创建新用户
GET    /api/v1/users/{user_id}          # 获取指定用户
PUT    /api/v1/users/{user_id}          # 完整更新用户
PATCH  /api/v1/users/{user_id}          # 部分更新用户
DELETE /api/v1/users/{user_id}          # 删除用户

GET    /api/v1/groups/{group_id}/members     # 获取群组成员
POST   /api/v1/groups/{group_id}/members     # 添加群组成员
```

### 1.2 HTTP方法规范

| 方法 | 用途 | 幂等性 | 请求体 | 响应 |
|------|------|--------|--------|------|
| GET | 查询资源 | ✓ | 无 | 资源数据 |
| POST | 创建资源 | ✗ | 有 | 201 + 新资源 |
| PUT | 完整替换 | ✓ | 有 | 200 + 更新后资源 |
| PATCH | 部分更新 | ✗ | 有 | 200 + 更新后资源 |
| DELETE | 删除资源 | ✓ | 无 | 204 或 200 |

### 1.3 HTTP状态码标准

#### 成功响应 (2xx)
- **200 OK**: GET, PUT, PATCH 成功
- **201 Created**: POST 创建成功，必须返回 Location 头
- **204 No Content**: DELETE 成功，无返回内容

#### 客户端错误 (4xx)
- **400 Bad Request**: 请求参数错误
- **401 Unauthorized**: 未认证或token失效
- **403 Forbidden**: 已认证但无权限
- **404 Not Found**: 资源不存在
- **409 Conflict**: 资源冲突（如重复创建）
- **429 Too Many Requests**: 超出速率限制

#### 服务器错误 (5xx)
- **500 Internal Server Error**: 服务器内部错误
- **502 Bad Gateway**: 上游服务错误
- **503 Service Unavailable**: 服务暂时不可用

### 1.4 统一响应格式

#### 成功响应格式
```json
{
  "success": true,
  "data": {
    // 实际数据内容
  },
  "message": "操作成功",
  "timestamp": "2025-01-24T10:30:00Z",
  "request_id": "req_7d4c2e1f8a9b"
}
```

#### 错误响应格式
```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "用户输入验证失败",
    "details": {
      "field": "email",
      "reason": "邮箱格式不正确"
    }
  },
  "timestamp": "2025-01-24T10:30:00Z",
  "request_id": "req_7d4c2e1f8a9b"
}
```

### 1.5 API版本控制

#### 版本策略
- **URI路径版本**: `/api/v1/`, `/api/v2/`
- **语义化版本**: MAJOR.MINOR.PATCH
- **向下兼容**: 新版本必须支持旧版本客户端
- **弃用策略**: 提前6个月通知版本弃用

---

## 2. 数据库设计规范

### 2.1 命名约定

#### 表名规范
- **使用复数名词**: `users`, `groups`, `credentials`
- **下划线分隔**: `user_groups`, `billing_records`
- **避免缩写**: 使用 `authentication` 而非 `auth`

#### 字段命名规范
- **使用单数名词**: `user_id`, `group_name`, `created_at`
- **主键命名**: `{table_name}_id` (如 `user_id`, `group_id`)
- **外键命名**: 引用主键名称 (如 `user_id` 引用 users.user_id)
- **布尔字段**: `is_active`, `has_permission`
- **时间字段**: `created_at`, `updated_at`, `deleted_at`

### 2.2 数据类型选择标准

#### 数值类型
```sql
-- 主键ID
user_id BIGINT AUTO_INCREMENT PRIMARY KEY

-- 状态枚举
status TINYINT NOT NULL DEFAULT 1

-- 金额 (避免浮点精度问题)
amount DECIMAL(15,4) NOT NULL

-- 百分比
percentage DECIMAL(5,2) NOT NULL
```

#### 字符串类型
```sql
-- 短字符串
email VARCHAR(255) NOT NULL
name VARCHAR(100) NOT NULL

-- 长文本
description TEXT NULL
content LONGTEXT NULL

-- 固定长度
country_code CHAR(2) NOT NULL
```

### 2.3 索引策略

#### 必建索引
- **主键索引**: 自动创建
- **外键索引**: 提升JOIN性能
- **唯一字段索引**: 如email、phone
- **查询频繁字段**: WHERE、ORDER BY使用的字段

---

## 3. JWT身份验证规范

### 3.1 Token结构标准

#### Header
```json
{
  "alg": "RS256",
  "typ": "JWT",
  "kid": "lyss-key-2025-01"
}
```

#### Payload (必需Claims)
```json
{
  "iss": "lyss-auth-service",
  "sub": "user_12345",
  "aud": "lyss-api",
  "exp": 1706097600,
  "iat": 1706096000,
  "jti": "token_abc123def",
  "scope": "user:read user:write group:admin"
}
```

### 3.2 安全要求

#### 密钥管理
- **最小密钥强度**: 256位随机熵
- **算法选择**: RS256 (推荐) 或 ES256
- **密钥轮换**: 每90天轮换一次
- **禁用算法**: 严禁使用 "none" 算法

#### Token生命周期
- **访问Token**: 15-30分钟过期
- **刷新Token**: 24小时过期
- **单设备登录**: 新登录使旧Token失效
- **主动撤销**: 支持Token黑名单机制

---

## 4. 日志记录规范

### 4.1 日志级别定义

| 级别 | 用途 | 示例场景 |
|------|------|----------|
| **DEBUG** | 详细调试信息 | 函数参数、SQL查询、变量值 |
| **INFO** | 一般信息记录 | 用户登录、操作成功、系统启动 |
| **WARN** | 警告信息 | 性能问题、配置警告、重试操作 |
| **ERROR** | 错误信息 | 异常捕获、操作失败、系统错误 |
| **FATAL** | 致命错误 | 系统崩溃、服务不可用 |

### 4.2 结构化日志格式

#### 标准JSON格式
```json
{
  "timestamp": "2025-01-24T10:30:00.123Z",
  "level": "INFO",
  "service": "user-service",
  "request_id": "req_7d4c2e1f8a9b",
  "user_id": "user_12345",
  "message": "用户登录成功",
  "context": {
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0...",
    "endpoint": "/api/v1/auth/login"
  },
  "metrics": {
    "duration_ms": 145,
    "memory_mb": 28.5
  }
}
```

---

## 5. 安全规范

### 5.1 API安全要求

#### 强制HTTPS
- **生产环境**: 必须使用HTTPS (TLS 1.2+)
- **证书验证**: 严格验证SSL证书
- **HSTS**: 启用HTTP严格传输安全
- **重定向**: HTTP自动重定向到HTTPS

#### 输入验证
- 所有用户输入必须验证和清理
- 使用白名单验证方式
- 防止SQL注入、XSS攻击
- 严格的参数类型检查

#### 速率限制
```python
# 示例：速率限制配置
RATE_LIMITS = {
    "login": "5/minute",
    "api_general": "1000/hour",
    "api_premium": "5000/hour",
    "password_reset": "3/hour"
}
```

### 5.2 数据加密要求

#### 静态数据加密
- **敏感字段**: 密码、API Key、个人信息
- **加密算法**: AES-256-GCM
- **密钥管理**: 使用专门的密钥管理服务
- **盐值**: 为每个加密字段使用唯一盐值

#### 传输加密
- **TLS版本**: 最低TLS 1.2，推荐TLS 1.3
- **密码套件**: 禁用弱密码套件
- **证书**: 使用可信CA签发的证书

---

## 6. 编码规范

### 6.1 通用编码原则

#### 命名规范
- **变量名**: 驼峰式命名 `userName`, `groupId`
- **常量名**: 全大写下划线分隔 `MAX_RETRY_COUNT`
- **函数名**: 动词开头 `getUserById`, `validateEmail`
- **类名**: 帕斯卡命名 `UserService`, `GroupManager`

#### 代码组织
- **单一职责**: 每个函数只做一件事
- **函数长度**: 建议不超过50行
- **嵌套深度**: 避免超过3层嵌套
- **异常处理**: 明确的异常处理和错误恢复

---

## 7. 性能规范

### 7.1 响应时间目标
- **简单查询**: < 100ms
- **复杂查询**: < 500ms
- **数据写入**: < 200ms
- **文件上传**: < 2s (10MB以内)

### 7.2 数据库性能
- **使用索引**: 确保WHERE、JOIN、ORDER BY字段有索引
- **避免N+1**: 使用JOIN或批量查询替代循环查询
- **分页查询**: 大量数据必须分页，默认20条/页
- **连接池**: 使用数据库连接池，避免频繁连接

---

## 8. 错误处理规范

### 8.1 统一错误码

#### 通用错误码
| 错误码 | HTTP状态码 | 描述 |
|--------|------------|------|
| `INVALID_REQUEST` | 400 | 请求参数无效 |
| `UNAUTHORIZED` | 401 | 未认证或认证失效 |
| `FORBIDDEN` | 403 | 无访问权限 |
| `NOT_FOUND` | 404 | 资源不存在 |
| `CONFLICT` | 409 | 资源冲突 |
| `RATE_LIMITED` | 429 | 请求频率过高 |
| `INTERNAL_ERROR` | 500 | 服务器内部错误 |

#### 业务错误码
- **USER_***：用户相关错误
- **GROUP_***：群组相关错误  
- **API_***：API调用相关错误
- **BILLING_***：计费相关错误

---

## 9. 部署和运维规范

### 9.1 环境配置

#### 环境变量管理
```bash
# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_NAME=lyss_platform
DB_USER=lyss_user
DB_PASSWORD=secure_password

# JWT配置
JWT_SECRET_KEY=your-super-secure-secret-key
JWT_ALGORITHM=RS256
JWT_ACCESS_TOKEN_EXPIRE_MINUTES=30

# 应用配置
DEBUG=false
LOG_LEVEL=INFO
```

### 9.2 健康检查

#### 健康检查端点
- `/health`: 基础健康检查
- `/health/ready`: 就绪检查，验证依赖服务
- `/health/live`: 存活检查，验证应用状态

---

## 10. 监控和告警规范

### 10.1 关键指标监控

#### 应用指标
- **响应时间**: P50, P95, P99响应时间
- **错误率**: HTTP 4xx, 5xx错误比例
- **吞吐量**: QPS (每秒查询数)
- **可用性**: 服务可用时间百分比

#### 业务指标
- **用户活跃度**: DAU, MAU
- **API调用量**: 各服务调用统计
- **计费统计**: Token消费量、成本统计

### 10.2 告警规则

#### 紧急告警
- **服务不可用**: 连续5分钟返回5xx错误
- **响应时间过长**: P95响应时间 > 2秒，持续5分钟
- **数据库连接失败**: 数据库连接失败率 > 10%

#### 警告告警
- **错误率上升**: 5分钟内4xx错误率 > 5%
- **内存使用率高**: 内存使用率 > 80%，持续10分钟
- **CPU使用率高**: CPU使用率 > 80%，持续10分钟

---

## 文档维护

**维护责任人**: 技术架构师  
**审核频率**: 每月一次  
**更新触发条件**: 
- 新技术标准发布
- 重大安全漏洞发现
- 业务需求变更
- 性能问题识别

**版本历史**:
- v1.0 (2025-01-24): 初始版本，基于2024-2025年最新标准制定

---

*本文档是Lyss AI Platform的技术宪法，所有开发活动必须严格遵循此规范。如有疑问或建议，请联系技术架构师。*