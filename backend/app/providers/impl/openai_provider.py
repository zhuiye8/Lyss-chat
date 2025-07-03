"""
OpenAI provider implementation.

This module implements the OpenAI API integration following
the standardized provider interface.
"""

import asyncio
from decimal import Decimal
from typing import Any, AsyncGenerator, Dict, List, Optional, Type

import httpx
from pydantic import BaseModel, Field, SecretStr

from app.providers.base import (
    BaseProviderConfig,
    ChatCompletionChunk,
    ChatCompletionRequest,
    ChatCompletionResponse,
    ConnectionError,
    LLMProvider,
    ModelInfo,
    ProviderError,
    RateLimitError,
)
from app.providers.registry import register_provider


class OpenAIProviderConfig(BaseProviderConfig):
    """Configuration model for OpenAI provider."""
    
    api_key: SecretStr = Field(..., description="OpenAI API key")
    base_url: Optional[str] = Field(
        "https://api.openai.com/v1",
        description="OpenAI API base URL"
    )
    organization: Optional[str] = Field(None, description="Organization ID")
    timeout: int = Field(60, description="Request timeout in seconds")
    max_retries: int = Field(3, description="Maximum retry attempts")


@register_provider("openai")
class OpenAIProvider(LLMProvider):
    """
    OpenAI API provider implementation.
    
    Implements the standardized LLMProvider interface for OpenAI's API,
    supporting both streaming and non-streaming completions.
    """
    
    def __init__(self, config: Dict[str, Any]):
        super().__init__(config)
        self._client: Optional[httpx.AsyncClient] = None
    
    @classmethod
    def get_config_model(cls) -> Type[BaseProviderConfig]:
        return OpenAIProviderConfig
    
    @classmethod
    def get_provider_info(cls) -> Dict[str, Any]:
        return {
            "name": "OpenAI",
            "description": "OpenAI GPT models including GPT-4, GPT-3.5, and others",
            "website": "https://openai.com",
            "documentation": "https://platform.openai.com/docs",
            "supports_streaming": True,
            "supports_function_calling": True,
            "supports_vision": True,
        }
    
    @property
    def client(self) -> httpx.AsyncClient:
        """Get or create HTTP client."""
        if self._client is None:
            headers = {
                "Authorization": f"Bearer {self.config.api_key.get_secret_value()}",
                "Content-Type": "application/json",
            }
            
            if self.config.organization:
                headers["OpenAI-Organization"] = self.config.organization
            
            self._client = httpx.AsyncClient(
                base_url=self.config.base_url,
                headers=headers,
                timeout=self.config.timeout,
            )
        
        return self._client
    
    async def test_connection(self) -> bool:
        """Test connection to OpenAI API."""
        try:
            response = await self.client.get("/models")
            response.raise_for_status()
            return True
        except httpx.HTTPStatusError as e:
            if e.response.status_code == 401:
                raise ConnectionError("Invalid API key") from e
            elif e.response.status_code == 429:
                raise ConnectionError("Rate limit exceeded") from e
            else:
                raise ConnectionError(f"HTTP {e.response.status_code}: {e.response.text}") from e
        except Exception as e:
            raise ConnectionError(f"Connection failed: {str(e)}") from e
    
    async def get_available_models(self) -> List[ModelInfo]:
        """Get available models from OpenAI API."""
        try:
            response = await self.client.get("/models")
            response.raise_for_status()
            data = response.json()
            
            models = []
            # Define known pricing for common models (per 1K tokens in USD)
            pricing = {
                "gpt-4o": {
                    "prompt": Decimal("0.0050"),
                    "completion": Decimal("0.0150")
                },
                "gpt-4o-mini": {
                    "prompt": Decimal("0.000150"),
                    "completion": Decimal("0.000600")
                },
                "gpt-4-turbo": {
                    "prompt": Decimal("0.0100"),
                    "completion": Decimal("0.0300")
                },
                "gpt-3.5-turbo": {
                    "prompt": Decimal("0.0015"),
                    "completion": Decimal("0.0020")
                },
            }
            
            for model_data in data.get("data", []):
                model_id = model_data["id"]
                
                # Only include chat completion models
                if not any(x in model_id for x in ["gpt-4", "gpt-3.5", "gpt-35"]):
                    continue
                
                # Get pricing if available
                model_pricing = pricing.get(model_id, {})
                
                model_info = ModelInfo(
                    model_name=model_id,
                    display_name=model_id.replace("-", " ").title(),
                    description=f"OpenAI {model_id} model",
                    context_length=self._get_context_length(model_id),
                    supports_streaming=True,
                    supports_function_calling=True,
                    supports_vision="vision" in model_id or "gpt-4" in model_id,
                    price_per_1k_prompt_tokens=model_pricing.get("prompt"),
                    price_per_1k_completion_tokens=model_pricing.get("completion"),
                )
                
                models.append(model_info)
            
            return models
            
        except Exception as e:
            raise ProviderError(f"Failed to get models: {str(e)}", "openai") from e
    
    def _get_context_length(self, model_name: str) -> int:
        """Get context length for known models."""
        context_lengths = {
            "gpt-4o": 128000,
            "gpt-4o-mini": 128000,
            "gpt-4-turbo": 128000,
            "gpt-4": 8192,
            "gpt-3.5-turbo": 4096,
        }
        
        for model, length in context_lengths.items():
            if model in model_name:
                return length
        
        return 4096  # Default
    
    async def chat_completion(
        self, 
        request: ChatCompletionRequest
    ) -> AsyncGenerator[ChatCompletionChunk, None]:
        """Generate streaming chat completion."""
        payload = {
            "model": request.messages[0].content if request.messages else "gpt-3.5-turbo",
            "messages": [
                {"role": msg.role, "content": msg.content} 
                for msg in request.messages
            ],
            "stream": True,
        }
        
        # Add optional parameters
        if request.temperature is not None:
            payload["temperature"] = request.temperature
        if request.max_tokens is not None:
            payload["max_tokens"] = request.max_tokens
        if request.top_p is not None:
            payload["top_p"] = request.top_p
        
        try:
            async with self.client.stream("POST", "/chat/completions", json=payload) as response:
                response.raise_for_status()
                
                async for line in response.aiter_lines():
                    if line.startswith("data: "):
                        data_str = line[6:]  # Remove "data: " prefix
                        
                        if data_str.strip() == "[DONE]":
                            yield ChatCompletionChunk(content="", finish_reason="stop")
                            break
                        
                        try:
                            import json
                            data = json.loads(data_str)
                            if "choices" in data and data["choices"]:
                                choice = data["choices"][0]
                                delta = choice.get("delta", {})
                                content = delta.get("content", "")
                                finish_reason = choice.get("finish_reason")
                                
                                yield ChatCompletionChunk(
                                    content=content,
                                    finish_reason=finish_reason
                                )
                        except Exception:
                            # Skip malformed chunks
                            continue
        
        except httpx.HTTPStatusError as e:
            if e.response.status_code == 429:
                raise RateLimitError("Rate limit exceeded") from e
            else:
                raise ProviderError(f"HTTP {e.response.status_code}: {e.response.text}") from e
        except Exception as e:
            raise ProviderError(f"Completion failed: {str(e)}") from e
    
    async def chat_completion_sync(
        self, 
        request: ChatCompletionRequest
    ) -> ChatCompletionResponse:
        """Generate non-streaming chat completion."""
        payload = {
            "model": request.messages[0].content if request.messages else "gpt-3.5-turbo",
            "messages": [
                {"role": msg.role, "content": msg.content} 
                for msg in request.messages
            ],
            "stream": False,
        }
        
        # Add optional parameters
        if request.temperature is not None:
            payload["temperature"] = request.temperature
        if request.max_tokens is not None:
            payload["max_tokens"] = request.max_tokens
        if request.top_p is not None:
            payload["top_p"] = request.top_p
        
        try:
            response = await self.client.post("/chat/completions", json=payload)
            response.raise_for_status()
            data = response.json()
            
            choice = data["choices"][0]
            usage = data["usage"]
            
            return ChatCompletionResponse(
                content=choice["message"]["content"],
                finish_reason=choice["finish_reason"],
                prompt_tokens=usage["prompt_tokens"],
                completion_tokens=usage["completion_tokens"],
                total_tokens=usage["total_tokens"],
            )
        
        except httpx.HTTPStatusError as e:
            if e.response.status_code == 429:
                raise RateLimitError("Rate limit exceeded") from e
            else:
                raise ProviderError(f"HTTP {e.response.status_code}: {e.response.text}") from e
        except Exception as e:
            raise ProviderError(f"Completion failed: {str(e)}") from e
    
    async def __aenter__(self):
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        if self._client:
            await self._client.aclose()