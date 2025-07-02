# LYSS AI 平台 - API 契约 (V2)

**版本**: 2.0
**格式**: OpenAPI 3.0.3

---

## 1. 概述

本文档定义了 LYSS AI 平台 V2 架构的 RESTful API。

*   **基地址**: `/api/v1`
*   **认证**: Bearer Token (JWT).

## 2. API 端点 (Endpoints)

### 2.1. 认证 (Auth) & 用户 (Users)

*   **管理**: 由 `fastapi-users` 管理，提供标准的用户注册、登录、密码重置、信息获取等接口。
*   **端点**: `/auth/...`, `/users/...`
*   **变更**: 无核心变更。

### 2.2. 供应��管理 (Provider Management) (V2 更新)

*   **标签**: `Providers`
*   **安全**: `BearerAuth`

#### POST /providers/
*   **摘要**: 添加一个新的AI供应商 (ORGANIZATION 或 PERSONAL).
*   **请求体**: `ProviderCreate` schema (包含 `scope` 字段).
*   **核心逻辑**: API 内部必须检查 `scope` 和用户角色。如果 `user.role` 是 `user`，但 `scope` 是 `ORGANIZATION`，则返回 403 Forbidden.
*   **响应 (201)**: `StandardResponse` (data: `ProviderRead`).

#### GET /providers/
*   **摘要**: 获取当前用户可见的供应商列表.
*   **核心逻辑**: 返回 `owner_id` 为当前用户 ID 的所有供应商。
*   **响应 (200)**: `StandardResponse` (data: `list[ProviderRead]`).

### 2.3. 模型与分发 (Models & Access) (V2 更新)

*   **标签**: `Models`, `Admin`
*   **安全**: `BearerAuth`

#### GET /chat/available-models
*   **摘要**: **[V2 新增]** 获取当前用户在聊天界面可用的所有模型。
*   **核心逻辑**: 返回一个聚合列表，包含：
    1.  所有通过 `user_model_accesses` 分发给当前用户的 `ORGANIZATION` 模型。
    2.  当前用户自己创建的所有 `PERSONAL` 供应商下的所有模型。
*   **响应 (200)**: `StandardResponse` (data: `list[ModelRead]`).

#### GET /admin/distributable-models
*   **摘要**: **[V2 新增]** (管��员) 获取可用于分发的模型列表。
*   **核心逻辑**: 仅查询 `scope = 'ORGANIZATION'` 的供应商及其下的所有模型。
*   **权限**: `require_admin`.
*   **响应 (200)**: `StandardResponse` (data: `list[ModelRead]`).

#### POST /admin/accesses/
*   **摘要**: (管理员) 将模型权限分发给用户。
*   **核心逻辑**: `model_id` 必须属于一个 `scope=ORGANIZATION` 的供应商。
*   **权限**: `require_admin`.
*   **响应 (201)**: `StandardResponse`.

### 2.4. 对话 (Chat) & 记忆 (Memory)

*   **变更**: 无核心 API 结构变更。

---

## 3. Schemas (Pydantic Models) (V2 更新)

```yaml
components:
  schemas:
    ProviderScope:
      type: string
      enum: [ORGANIZATION, PERSONAL]

    ProviderBase:
      type: object
      required:
        - name
        - provider_type
        - scope
        - api_key
      properties:
        name:
          type: string
        provider_type:
          type: string
        scope:
          $ref: '#/components/schemas/ProviderScope'
        base_url:
          type: string
          format: uri
        api_key:
          type: string
          description: "未加密的API Key，后端将负责加密存储"

    ProviderCreate:
      $ref: '#/components/schemas/ProviderBase'

    ProviderRead:
      allOf:
        - $ref: '#/components/schemas/ProviderBase'
        - type: object
          properties:
            id:
              type: string
              format: uuid
            owner_id:
              type: string
              format: uuid
            api_key:
              type: string
              description: "为安全起见，读操作不返回key"
              example: "********"

    ModelRead:
      type: object
      properties:
        id:
          type: string
          format: uuid
        model_name:
          type: string
        provider:
          $ref: '#/components/schemas/ProviderRead' # 嵌套供应商信息
        # ... other fields

    ChatCompletionRequest:
      type: object
      required:
        - model_id # 使用模型ID
        - messages
      properties:
        model_id:
          type: string
          format: uuid
        # ... other fields
```
