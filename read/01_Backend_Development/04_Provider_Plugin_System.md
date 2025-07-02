# LYSS AI 平台 - 供应商插件化架构 (V2.1 优化版)

**版本**: 2.1
**最后更新**: 2025年7月2日

---

## 1. 设计目标与理念

V2.1 版本的插件系统核心在于**灵活性**和**解耦**。我们借鉴了业界优秀的实践（如 `Dify` 和 `LiteLLM`），采用“配置驱动”的模式。

*   **面向接口编程**: 定义统一的 `LLMProvider` 抽象基类 (ABC)，所有插件必须实现其接口。
*   **配置与代码分离**:
    *   **配置 (`config_encrypted`)**: 供应商的认证信息（如 API Key, Secret, URL 等）被加密后，以一个独立的 JSON 对象形式存储在数据库的 `providers.config_encrypted` 字段中。
    *   **代码 (`impl/`)**: 插件代码只负责实现与外部 API 交互的纯粹逻辑。
*   **动态加载与注册**: 系统启动时自动发现所有插件，并通过工厂模式按需实例化。

## 2. 核心组件设计

### 2.1. `base.py`: 抽象基类与配置模型

每个插件都需要定义自己的配置模型，这通过 Pydantic 实现，以获得类型安全和自动验证。

```python
# backend/app/providers/base.py

from abc import ABC, abstractmethod
from typing import AsyncGenerator, List, Dict, Any, Type
from pydantic import BaseModel, Field

# --- 配置模型 ---
class BaseProviderConfig(BaseModel):
    """每个插件都应定义一个继承自此基类的配置模型。"""
    pass

# --- 数据模型 ---
class ModelInfo(BaseModel):
    # ... (无变化)

# --- 抽象基类 ---
class LLMProvider(ABC):
    """
    Abstract Base Class for all LLM providers.
    """
    config: BaseProviderConfig

    def __init__(self, config: Dict[str, Any]):
        """
        Initializes the provider with its configuration.
        The config is first validated against the provider's specific config model.
        """
        ConfigModel = self.get_config_model()
        self.config = ConfigModel(**config)

    @classmethod
    @abstractmethod
    def get_config_model(cls) -> Type[BaseProviderConfig]:
        """
        每个插件必须实现此方法，返回其特定的Pydantic配置模型。
        """
        pass
    
    # ... 其他抽象方法 (get_available_models, chat_completion, etc.) 无变化
```

### 2.2. `factory.py`: 供应商工厂 (V2.1 更新)

工厂现在负责解密、解析 JSON，并使用插件自身的配置模型进行验证。

```python
# backend/app/providers/factory.py

import json
from .registry import get_provider_class
from .base import LLMProvider
from app.models.provider import Provider as ProviderModel
from app.core.security import decrypt # 假设的解密函数

def create_provider_instance(provider_model: ProviderModel) -> LLMProvider:
    """
    Creates a provider instance from a provider SQLAlchemy model.
    """
    provider_cls = get_provider_class(provider_model.provider_type)
    
    # 1. 解密存储在数据库中的配置字符串
    decrypted_config_str = decrypt(provider_model.config_encrypted)
    
    # 2. 将 JSON 字符串解析为字典
    config_dict = json.loads(decrypted_config_str)
    
    # 3. 使用插件类进行实例化 (内部会进行Pydantic验证)
    return provider_cls(config_dict)
```

### 2.3. 具体实现示例: `openai_provider.py` (V2.1 更新)

插件现在需要定义自己的 `ProviderConfig` 模型。

```python
# backend/app/providers/impl/openai_provider.py

from typing import Type
from pydantic import Field, SecretStr
from ..base import LLMProvider, ModelInfo, BaseProviderConfig
from ..registry import register_provider
import openai

# 1. 定义此 Provider 专属的配置模型
class OpenAIProviderConfig(BaseProviderConfig):
    api_key: SecretStr = Field(..., description="OpenAI API Key")
    base_url: str | None = Field(None, description="API Base URL")

@register_provider("openai")
class OpenAIProvider(LLMProvider):
    
    # 2. 实现 get_config_model 方法
    @classmethod
    def get_config_model(cls) -> Type[BaseProviderConfig]:
        return OpenAIProviderConfig

    # 3. 在实现中使用强类型的 self.config
    async def test_connection(self) -> bool:
        client = openai.AsyncOpenAI(
            api_key=self.config.api_key.get_secret_value(),
            base_url=self.config.base_url
        )
        try:
            await client.models.list(timeout=5)
            return True
        except Exception as e:
            raise ConnectionError(f"Failed to connect: {e}")

    async def get_available_models(self) -> List[ModelInfo]:
        client = openai.AsyncOpenAI(
            api_key=self.config.api_key.get_secret_value(),
            base_url=self.config.base_url
        )
        # ...
```

## 4. 优势总结

这种 **“基类 + Pydantic配置模型 + JSONB数据库字段”** 的模式，为平台提供了极致的灵活性和扩展性，是现代插件化系统设计的最佳实践之一。

