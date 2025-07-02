# LYSS AI 平台 - 数据库设计 (V2.1 优化版)

**版本**: 2.1
**最后更新**: 2025年7月2日

---

## 1. 概述

本文档详细描述了 LYSS AI 平台 V2 架构所使用的数据库结构。我们选用 **PostgreSQL** 作为核心业务数据库，**Qdrant** 作为向量数据库。

## 2. 实体关系图 (ERD) V2.1

下图反映了以 `providers` 表的 `scope` 为核心，并使用 `config_encrypted` 存储灵活配置的最终架构。

```mermaid
erDiagram
    USERS ||--o{ PROVIDERS : "owns (as owner)"
    USERS ||--o{ USER_MODEL_ACCESSES : "receives"
    USERS ||--o{ USAGE_LOGS : "generates"
    
    PROVIDERS {
        UUID id PK
        string name
        string provider_type
        provider_scope scope "Enum: ORGANIZATION, PERSONAL"
        UUID owner_id FK "-> USERS.id"
        text config_encrypted "Encrypted JSONB"
        -- ... other fields
    }
    
    PROVIDERS ||--|{ MODELS : "provides"
    
    MODELS ||--o{ USER_MODEL_ACCESSES : "is accessed via"
    MODELS ||--o{ USAGE_LOGS : "is used in"
```

## 3. 表结构详述 (PostgreSQL)

### 3.1. `users`

*   **描述**: 存储平台的用户信息。
*   **SQL (骨架)**:
    ```sql
    CREATE TYPE user_role AS ENUM ('user', 'admin', 'super_admin');

    CREATE TABLE users (
        id UUID PRIMARY KEY,
        email VARCHAR(255) UNIQUE NOT NULL,
        hashed_password VARCHAR(255) NOT NULL,
        is_active BOOLEAN DEFAULT TRUE NOT NULL,
        is_superuser BOOLEAN DEFAULT FALSE NOT NULL,
        is_verified BOOLEAN DEFAULT FALSE NOT NULL,
        role user_role DEFAULT 'user' NOT NULL,
        created_at TIMESTAMPTZ DEFAULT NOW(),
        updated_at TIMESTAMPTZ DEFAULT NOW()
    );
    ```

### 3.2. `providers` (V2.1 优化)

*   **描述**: 存储供应商配置。`scope` 决定其能力，`config_encrypted` 提供灵活的、面向未来的配置存储。
*   **SQL**:
    ```sql
    CREATE TYPE provider_scope AS ENUM ('ORGANIZATION', 'PERSONAL');

    CREATE TABLE providers (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        name VARCHAR(255) NOT NULL,
        provider_type VARCHAR(100) NOT NULL, -- e.g., 'openai', 'anthropic'
        
        -- V2 核心字段
        scope provider_scope NOT NULL,
        owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        
        -- V2.1 优化: 使用加密的JSONB字段存储灵活的配置
        config_encrypted TEXT NOT NULL, -- 存储加密后的JSON字符串
        
        is_enabled BOOLEAN DEFAULT TRUE,
        created_at TIMESTAMPTZ DEFAULT NOW(),
        updated_at TIMESTAMPTZ DEFAULT NOW()
    );
    CREATE INDEX idx_providers_owner_id ON providers(owner_id);
    CREATE INDEX idx_providers_scope ON providers(scope);
    ```
    **设计思考**: 放弃 `base_url`, `api_key_encrypted` 等固定字段，改为使用单一的 `config_encrypted` 字段。这使得我们可以支持任何复杂的认证方式（如需要`secret`, `project_id`等），而无需修改数据库结构，极大地增强了系统的可扩展性。

### 3.3. `models`

*   **描述**: 存储从供应商同步来的可用模型信息。
*   **SQL**:
    ```sql
    CREATE TABLE models (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        provider_id UUID NOT NULL REFERENCES providers(id) ON DELETE CASCADE,
        model_name VARCHAR(255) NOT NULL,
        -- ... 其他字段
        UNIQUE (provider_id, model_name)
    );
    ```

### 3.4. `user_model_accesses`

*   **描述**: 核心权限表，仅用于分发 `scope=ORGANIZATION` 的模型。
*   **SQL**:
    ```sql
    CREATE TABLE user_model_accesses (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        model_id UUID NOT NULL REFERENCES models(id) ON DELETE CASCADE,
        -- ... 其他权限规则字段
        UNIQUE (user_id, model_id)
    );
    ```

---
(后续表结构无变化)
