# LYSS AI 平台 - 核心业务流程实现 (V2.1)

**版本**: 2.1
**最后更新**: 2025年7月2日

---

## 1. 概述

本文档阐述了 V2.1 架构下，平台后端关键业务流程的实现逻辑。

## 2. 业务流程一: 聚合用户可用模型列表

此���程对应 `GET /api/v1/chat/available-models` 端点，是用户开始聊天的第一步。

### 2.1. 流程图

```mermaid
sequenceDiagram
    participant User as 用户
    participant API as 后端API
    participant ChatService as 对话服务
    participant DB as 数据库

    User->>API: 1. 请求 /chat/available-models
    API->>ChatService: 2. get_available_models_for_user(user_id)
    
    par 并行查询
        ChatService->>DB: 3a. 查询 PERSONAL 模型<br/>(providers.owner_id = user_id)
        DB-->>ChatService: 4a. 返回个人模型列表
    and
        ChatService->>DB: 3b. 查询 ORGANIZATION 模型<br/>(JOIN user_model_accesses)
        DB-->>ChatService: 4b. 返回已授权的组织模型列表
    end
    
    ChatService->>ChatService: 5. 合并并去重两个列表
    ChatService-->>API: 6. 返回最终聚合列表
    API-->>User: 7. 返回模型列表JSON
```

### 2.2. 服务层实现

```python
# app/services/chat_service.py
from sqlalchemy.orm import selectinload
from sqlalchemy import select
from app.models import Provider, Model, UserModelAccess

class ChatService:
    async def get_available_models_for_user(self, db: AsyncSession, user_id: uuid.UUID) -> list[Model]:
        # 1. 获取用户自己的 PERSONAL 模型
        personal_stmt = (
            select(Model)
            .join(Provider)
            .where(Provider.owner_id == user_id, Provider.scope == 'PERSONAL', Provider.is_enabled == True)
            .options(selectinload(Model.provider)) # 预加载供应商信息
        )
        personal_models_result = await db.execute(personal_stmt)
        personal_models = personal_models_result.scalars().all()

        # 2. 获取分发给用户的 ORGANIZATION 模型
        org_stmt = (
            select(Model)
            .join(UserModelAccess, Model.id == UserModelAccess.model_id)
            .join(Provider)
            .where(UserModelAccess.user_id == user_id, Provider.scope == 'ORGANIZATION', Provider.is_enabled == True)
            .options(selectinload(Model.provider))
        )
        org_models_result = await db.execute(org_stmt)
        org_models = org_models_result.scalars().all()

        # 3. 合并和去重
        all_models = {model.id: model for model in personal_models}
        all_models.update({model.id: model for model in org_models})

        return list(all_models.values())
```

## 3. 业务流程二: 对话权限验证 (V2.1)

在用户发起对话请求时，验证其对所选模型 `model_id` 的使用权。

### 3.1. 流程图

```mermaid
sequenceDiagram
    participant ChatService as 对话服务
    participant DB as 数据库
    
    ChatService->>DB: 1. 根据 model_id 查询模型及其供应商
    DB-->>ChatService: 2. 返回 Model 和 Provider 信息
    
    alt Provider.scope == 'PERSONAL'
        ChatService->>ChatService: 3a. 检查 Provider.owner_id 是否等于 current_user.id
        alt 检查通过
            ChatService->>ProviderPlugin: 4a. 验证成功，继续
        else 检查失败
            ChatService->>API: 4b. 抛出 403 Forbidden
        end
    else Provider.scope == 'ORGANIZATION'
        ChatService->>DB: 3b. 在 user_model_accesses 表中查询 (user_id, model_id)
        alt 记录存在
             ChatService->>ChatService: 4c. 检查配额和有效期
             ChatService->>ProviderPlugin: 5a. 验证成功，继续
        else 记录不存在
            ChatService->>API: 5b. 抛出 403 Forbidden
        end
    end
```

---
(成本计算等其他流程无核心变化)
