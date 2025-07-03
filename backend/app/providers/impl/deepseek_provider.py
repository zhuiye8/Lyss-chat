"""
DeepSeek provider implementation.

This module implements the DeepSeek API integration following
the OpenAI-compatible interface standard.
"""

from decimal import Decimal
from typing import Any, AsyncGenerator, Dict, List, Type

from pydantic import Field, SecretStr

from app.providers.base import (
    BaseProviderConfig,
    ChatCompletionChunk,
    ChatCompletionRequest,
    ChatCompletionResponse,
    LLMProvider,
    ModelInfo,
)
from app.providers.impl.openai_provider import OpenAIProvider
from app.providers.registry import register_provider


class DeepSeekProviderConfig(BaseProviderConfig):
    """Configuration model for DeepSeek provider."""
    
    api_key: SecretStr = Field(..., description="DeepSeek API key")
    base_url: str = Field(
        "https://api.deepseek.com/v1",
        description="DeepSeek API base URL"
    )
    timeout: int = Field(60, description="Request timeout in seconds")
    max_retries: int = Field(3, description="Maximum retry attempts")


@register_provider("deepseek")
class DeepSeekProvider(LLMProvider):
    """
    DeepSeek API provider implementation.
    
    DeepSeek uses OpenAI-compatible API, so we inherit most functionality
    from OpenAIProvider and customize model-specific details.
    """
    
    def __init__(self, config: Dict[str, Any]):
        super().__init__(config)
        # Create internal OpenAI provider with DeepSeek config
        self._openai_provider = OpenAIProvider({
            "api_key": self.config.api_key.get_secret_value(),
            "base_url": self.config.base_url,
            "timeout": self.config.timeout,
        })
    
    @classmethod
    def get_config_model(cls) -> Type[BaseProviderConfig]:
        return DeepSeekProviderConfig
    
    @classmethod
    def get_provider_info(cls) -> Dict[str, Any]:
        return {
            "name": "DeepSeek",
            "description": "DeepSeek AI models with strong coding and reasoning capabilities",
            "website": "https://www.deepseek.com",
            "documentation": "https://platform.deepseek.com/api-docs",
            "supports_streaming": True,
            "supports_function_calling": True,
            "supports_vision": False,
        }
    
    async def test_connection(self) -> bool:
        """Test connection to DeepSeek API."""
        return await self._openai_provider.test_connection()
    
    async def get_available_models(self) -> List[ModelInfo]:
        """Get available models from DeepSeek API."""
        # DeepSeek specific models with pricing
        models = [
            ModelInfo(
                model_name="deepseek-chat",
                display_name="DeepSeek Chat",
                description="DeepSeek's flagship conversational AI model",
                context_length=32768,
                supports_streaming=True,
                supports_function_calling=True,
                supports_vision=False,
                price_per_1k_prompt_tokens=Decimal("0.00014"),
                price_per_1k_completion_tokens=Decimal("0.00028"),
            ),
            ModelInfo(
                model_name="deepseek-coder",
                display_name="DeepSeek Coder",
                description="Specialized model for code generation and programming tasks",
                context_length=32768,
                supports_streaming=True,
                supports_function_calling=True,
                supports_vision=False,
                price_per_1k_prompt_tokens=Decimal("0.00014"),
                price_per_1k_completion_tokens=Decimal("0.00028"),
            ),
        ]
        
        return models
    
    async def chat_completion(
        self, 
        request: ChatCompletionRequest
    ) -> AsyncGenerator[ChatCompletionChunk, None]:
        """Generate streaming chat completion using DeepSeek API."""
        async for chunk in self._openai_provider.chat_completion(request):
            yield chunk
    
    async def chat_completion_sync(
        self, 
        request: ChatCompletionRequest
    ) -> ChatCompletionResponse:
        """Generate non-streaming chat completion using DeepSeek API."""
        return await self._openai_provider.chat_completion_sync(request)
    
    async def __aenter__(self):
        await self._openai_provider.__aenter__()
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        await self._openai_provider.__aexit__(exc_type, exc_val, exc_tb)