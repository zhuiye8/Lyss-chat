# LYSS AI 平台 - Mem0 智能记忆模块集成方案 (V2)

**版本**: 2.0
**最后更新**: 2025年7月2日

---

## 1. 概述

智能记忆是提升平台用户体验的关键功能。我们将集成 **Mem0 开源版**，并使用 **Qdrant** 作为其向量存储后端，为用户提供个性化的、具备长期上下文感知能力的对话体验。

## 2. 核心架构：全局唯一的记忆模型

根据 V2 架构，`mem0.ai` 服务自身进行记忆提取和向量化所需的 LLM 和 Embedding 模型是**全局统一配置的**，与用户侧的供应商体系完全分离。

*   **配置方式**: 由超级管理员通过部署环境的环境变量进行设置。
*   **作用**: 这确保了平台记忆生成和检索的一致性和可管理性，避免了因用户使用不同模型而导致记忆质量参差不齐的问题。
*   **用户隔离**: 尽管记忆模型是全局的，但每个用户的记忆数据在存储和检索时，仍通过 `user_id` 进行严格隔离。

## 3. 环境配置 (`.env`)

为了让 Mem0 正确运行，超级管理员需要在部署环境（例如 `.env` 文件）中配置以下变量。

```dotenv
# .env.example

# --- Mem0 Global Configuration ---
# Mem0 服务自身用于提取、查询和总结记忆所使用的LLM和Embedding模型
# 这与用户在聊天时选择的供应商和模型是完全独立的

# 1. Vector Store (Required)
QDRANT_URL="http://qdrant:6333"

# 2. LLM Provider for Memory Processing (Required)
# Options: "openai", "anthropic", "google", "groq", "ollama", etc.
MEM0_LLM_PROVIDER="openai" 

# 3. Models for Memory Processing (Required)
MEM0_LLM_MODEL="gpt-4o"
MEM0_EMBEDDING_MODEL="text-embedding-3-small"

# 4. API Credentials for Memory Processing (Required)
MEM0_PROVIDER_API_KEY="sk-your_openai_api_key_for_mem0"
MEM0_PROVIDER_BASE_URL="https://api.openai.com/v1" # Optional, for proxy or compatible APIs
```

**设计思考**: 这种将“记忆处理模型”与“对话模型”分离的设计，是借鉴了类似系统的先进经验。它允许我们为后台的记忆任务选用性价比最高、最适合文本处理的模型（如 `gpt-4o`），而用户在对话时则可以自由选择其他更专业或更强大的模型，两者互不干扰。

## 4. `MemoryService` 实现

`MemoryService` 的实现保持不变，它会透明地使用 `mem0` 库。`mem0` 库在初始化时，会自动从环境变量中读取上述 `MEM0_*` 和 `QDRANT_*` 配置。

```python
# backend/app/services/memory_service.py

import os
from mem0 import Memory # mem0ai >= 0.1.94
from app.models.user import User
from app.schemas.chat import ChatCompletionRequest

class MemoryService:
    def __init__(self):
        # 在此初始化时，mem0 会自动读取环境变量来配置其内部的
        # LLM Provider, Embedding Model, 和 Qdrant 连接。
        self.mem0 = Memory()

    async def retrieve_memories_as_text(self, user: User, query: str, limit: int = 3) -> str:
        # ... (实现无变化)
        pass

    async def add_memory_from_conversation(self, user: User, request: ChatCompletionRequest, response_text: str):
        # ... (实现无变化)
        pass
```

## 5. 集成流程

集成架构图和在 `ChatService` 中调用的方式保持不变。`ChatService` 只负责调用 `MemoryService` 的方法，而无需关心 `mem0` 内部是如何配置和工作的，这体现了良好的封装性。
