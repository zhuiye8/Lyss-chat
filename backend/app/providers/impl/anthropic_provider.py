"""
Anthropic Claude provider implementation.

This module implements the Anthropic API integration following
the standardized provider interface.
"""

from decimal import Decimal
from typing import Any, AsyncGenerator, Dict, List, Optional, Type

import httpx
from pydantic import Field, SecretStr

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


class AnthropicProviderConfig(BaseProviderConfig):
    """Configuration model for Anthropic provider."""
    
    api_key: SecretStr = Field(..., description="Anthropic API key")
    base_url: str = Field(
        "https://api.anthropic.com",
        description="Anthropic API base URL"
    )
    timeout: int = Field(60, description="Request timeout in seconds")
    max_retries: int = Field(3, description="Maximum retry attempts")


@register_provider("anthropic")
class AnthropicProvider(LLMProvider):
    """
    Anthropic Claude API provider implementation.
    
    Implements the standardized LLMProvider interface for Anthropic's API,
    handling the different message format and API structure.
    """
    
    def __init__(self, config: Dict[str, Any]):
        super().__init__(config)
        self._client: Optional[httpx.AsyncClient] = None
    
    @classmethod
    def get_config_model(cls) -> Type[BaseProviderConfig]:
        return AnthropicProviderConfig
    
    @classmethod
    def get_provider_info(cls) -> Dict[str, Any]:
        return {
            "name": "Anthropic",
            "description": "Anthropic Claude models for helpful, harmless, and honest AI",
            "website": "https://www.anthropic.com",
            "documentation": "https://docs.anthropic.com",
            "supports_streaming": True,
            "supports_function_calling": True,
            "supports_vision": True,
        }
    
    @property
    def client(self) -> httpx.AsyncClient:
        """Get or create HTTP client."""
        if self._client is None:
            headers = {
                "x-api-key": self.config.api_key.get_secret_value(),
                "Content-Type": "application/json",
                "anthropic-version": "2023-06-01",
            }
            
            self._client = httpx.AsyncClient(
                base_url=self.config.base_url,
                headers=headers,
                timeout=self.config.timeout,
            )
        
        return self._client
    
    async def test_connection(self) -> bool:
        """Test connection to Anthropic API."""
        try:
            # Anthropic doesn't have a models endpoint, so we make a minimal completion request
            payload = {
                "model": "claude-3-haiku-20240307",
                "max_tokens": 1,
                "messages": [{"role": "user", "content": "Hi"}],
            }
            
            response = await self.client.post("/v1/messages", json=payload)
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
        """Get available Claude models."""
        # Anthropic doesn't provide a models API, so we return known models
        models = [
            ModelInfo(
                model_name="claude-3-5-sonnet-20241022",
                display_name="Claude 3.5 Sonnet",
                description="Most intelligent model with enhanced agentic capabilities",
                context_length=200000,
                supports_streaming=True,
                supports_function_calling=True,
                supports_vision=True,
                price_per_1k_prompt_tokens=Decimal("0.0030"),
                price_per_1k_completion_tokens=Decimal("0.0150"),
            ),
            ModelInfo(
                model_name="claude-3-5-haiku-20241022",
                display_name="Claude 3.5 Haiku",
                description="Fastest and most affordable model in the Claude 3.5 family",
                context_length=200000,
                supports_streaming=True,
                supports_function_calling=True,
                supports_vision=True,
                price_per_1k_prompt_tokens=Decimal("0.00080"),
                price_per_1k_completion_tokens=Decimal("0.00400"),
            ),
            ModelInfo(
                model_name="claude-3-opus-20240229",
                display_name="Claude 3 Opus",
                description="Most powerful model for highly complex tasks",
                context_length=200000,
                supports_streaming=True,
                supports_function_calling=True,
                supports_vision=True,
                price_per_1k_prompt_tokens=Decimal("0.0150"),
                price_per_1k_completion_tokens=Decimal("0.0750"),
            ),
            ModelInfo(
                model_name="claude-3-haiku-20240307",
                display_name="Claude 3 Haiku",
                description="Fastest and most compact model for everyday tasks",
                context_length=200000,
                supports_streaming=True,
                supports_function_calling=True,
                supports_vision=True,
                price_per_1k_prompt_tokens=Decimal("0.00025"),
                price_per_1k_completion_tokens=Decimal("0.00125"),
            ),
        ]
        
        return models
    
    def _convert_messages(self, messages: List[Any]) -> tuple[List[Dict], str]:
        """Convert OpenAI format messages to Anthropic format."""
        anthropic_messages = []
        system_message = ""
        
        for msg in messages:
            if msg.role == "system":
                system_message = msg.content
            elif msg.role in ["user", "assistant"]:
                anthropic_messages.append({
                    "role": msg.role,
                    "content": msg.content
                })
        
        return anthropic_messages, system_message
    
    async def chat_completion(
        self, 
        request: ChatCompletionRequest
    ) -> AsyncGenerator[ChatCompletionChunk, None]:
        """Generate streaming chat completion."""
        messages, system_message = self._convert_messages(request.messages)
        
        payload = {
            "model": "claude-3-5-sonnet-20241022",  # Default model
            "max_tokens": request.max_tokens or 4096,
            "messages": messages,
            "stream": True,
        }
        
        if system_message:
            payload["system"] = system_message
        
        if request.temperature is not None:
            payload["temperature"] = request.temperature
        if request.top_p is not None:
            payload["top_p"] = request.top_p
        
        try:
            async with self.client.stream("POST", "/v1/messages", json=payload) as response:
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
                            
                            if data.get("type") == "content_block_delta":
                                delta = data.get("delta", {})
                                content = delta.get("text", "")
                                
                                yield ChatCompletionChunk(content=content)
                            
                            elif data.get("type") == "message_stop":
                                yield ChatCompletionChunk(content="", finish_reason="stop")
                                break
                        
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
        messages, system_message = self._convert_messages(request.messages)
        
        payload = {
            "model": "claude-3-5-sonnet-20241022",  # Default model
            "max_tokens": request.max_tokens or 4096,
            "messages": messages,
        }
        
        if system_message:
            payload["system"] = system_message
        
        if request.temperature is not None:
            payload["temperature"] = request.temperature
        if request.top_p is not None:
            payload["top_p"] = request.top_p
        
        try:
            response = await self.client.post("/v1/messages", json=payload)
            response.raise_for_status()
            data = response.json()
            
            content = ""
            if data.get("content"):
                content = data["content"][0].get("text", "")
            
            usage = data.get("usage", {})
            
            return ChatCompletionResponse(
                content=content,
                finish_reason=data.get("stop_reason", "stop"),
                prompt_tokens=usage.get("input_tokens", 0),
                completion_tokens=usage.get("output_tokens", 0),
                total_tokens=usage.get("input_tokens", 0) + usage.get("output_tokens", 0),
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