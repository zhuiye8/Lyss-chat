"""
Abstract base classes for AI provider plugins.

This module defines the interfaces that all provider plugins must implement
to ensure consistent behavior across different AI services.
"""

from abc import ABC, abstractmethod
from decimal import Decimal
from typing import Any, AsyncGenerator, Dict, List, Optional, Type

from pydantic import BaseModel, Field


class BaseProviderConfig(BaseModel):
    """
    Base configuration model for all providers.
    
    Each provider should inherit from this and define their specific
    configuration fields with proper validation.
    """
    pass


class ModelInfo(BaseModel):
    """
    Information about an AI model provided by a service.
    
    This standardized format allows the platform to work with
    models from different providers uniformly.
    """
    
    model_name: str = Field(..., description="Unique model identifier")
    display_name: Optional[str] = Field(None, description="Human-readable model name")
    description: Optional[str] = Field(None, description="Model description")
    
    # Model capabilities
    context_length: Optional[int] = Field(None, description="Maximum context length in tokens")
    max_output_tokens: Optional[int] = Field(None, description="Maximum output tokens")
    supports_streaming: bool = Field(True, description="Whether model supports streaming")
    supports_function_calling: bool = Field(False, description="Whether model supports function calling")
    supports_vision: bool = Field(False, description="Whether model supports vision/image input")
    
    # Pricing (per 1K tokens in USD)
    price_per_1k_prompt_tokens: Optional[Decimal] = Field(None, description="Cost per 1K prompt tokens")
    price_per_1k_completion_tokens: Optional[Decimal] = Field(None, description="Cost per 1K completion tokens")
    price_per_1k_tokens: Optional[Decimal] = Field(None, description="Unified cost per 1K tokens")
    
    class Config:
        from_attributes = True


class ChatMessage(BaseModel):
    """Standardized chat message format."""
    role: str = Field(..., description="Message role: user, assistant, system")
    content: str = Field(..., description="Message content")
    
    class Config:
        from_attributes = True


class ChatCompletionRequest(BaseModel):
    """Standardized chat completion request format."""
    messages: List[ChatMessage] = Field(..., description="List of chat messages")
    stream: bool = Field(False, description="Whether to stream the response")
    temperature: Optional[float] = Field(None, ge=0.0, le=2.0, description="Sampling temperature")
    max_tokens: Optional[int] = Field(None, gt=0, description="Maximum tokens to generate")
    top_p: Optional[float] = Field(None, ge=0.0, le=1.0, description="Top-p sampling parameter")
    
    class Config:
        from_attributes = True


class ChatCompletionChunk(BaseModel):
    """Streaming chat completion chunk."""
    content: str = Field("", description="Incremental content")
    finish_reason: Optional[str] = Field(None, description="Reason for completion")
    
    class Config:
        from_attributes = True


class ChatCompletionResponse(BaseModel):
    """Complete chat completion response."""
    content: str = Field(..., description="Generated content")
    finish_reason: str = Field(..., description="Reason for completion")
    prompt_tokens: int = Field(..., description="Number of prompt tokens used")
    completion_tokens: int = Field(..., description="Number of completion tokens generated")
    total_tokens: int = Field(..., description="Total tokens used")
    
    class Config:
        from_attributes = True


class LLMProvider(ABC):
    """
    Abstract base class for all LLM provider implementations.
    
    This class defines the interface that all provider plugins must implement
    to ensure consistent behavior and enable the factory pattern.
    """
    
    config: BaseProviderConfig
    
    def __init__(self, config: Dict[str, Any]):
        """
        Initialize the provider with configuration.
        
        Args:
            config: Configuration dictionary from encrypted storage
        """
        ConfigModel = self.get_config_model()
        self.config = ConfigModel(**config)
    
    @classmethod
    @abstractmethod
    def get_config_model(cls) -> Type[BaseProviderConfig]:
        """
        Return the Pydantic model for this provider's configuration.
        
        Returns:
            Pydantic model class for configuration validation
        """
        pass
    
    @classmethod
    @abstractmethod
    def get_provider_info(cls) -> Dict[str, Any]:
        """
        Return provider metadata and capabilities.
        
        Returns:
            Dictionary with provider information
        """
        pass
    
    @abstractmethod
    async def test_connection(self) -> bool:
        """
        Test the connection to the provider's API.
        
        Returns:
            True if connection is successful
            
        Raises:
            ConnectionError: If connection fails with details
        """
        pass
    
    @abstractmethod
    async def get_available_models(self) -> List[ModelInfo]:
        """
        Retrieve list of available models from the provider.
        
        Returns:
            List of ModelInfo objects
            
        Raises:
            ProviderError: If unable to retrieve models
        """
        pass
    
    @abstractmethod
    async def chat_completion(
        self, 
        request: ChatCompletionRequest
    ) -> AsyncGenerator[ChatCompletionChunk, None]:
        """
        Generate chat completion with streaming support.
        
        Args:
            request: Chat completion request
            
        Yields:
            ChatCompletionChunk objects for streaming response
            
        Raises:
            ProviderError: If completion fails
        """
        pass
    
    @abstractmethod
    async def chat_completion_sync(
        self, 
        request: ChatCompletionRequest
    ) -> ChatCompletionResponse:
        """
        Generate non-streaming chat completion.
        
        Args:
            request: Chat completion request
            
        Returns:
            Complete ChatCompletionResponse
            
        Raises:
            ProviderError: If completion fails
        """
        pass


class ProviderError(Exception):
    """Base exception for provider-related errors."""
    
    def __init__(self, message: str, provider_type: Optional[str] = None):
        self.message = message
        self.provider_type = provider_type
        super().__init__(self.message)


class ConnectionError(ProviderError):
    """Exception raised when provider connection fails."""
    pass


class ModelNotFoundError(ProviderError):
    """Exception raised when requested model is not available."""
    pass


class RateLimitError(ProviderError):
    """Exception raised when rate limit is exceeded."""
    pass