"""
Model schemas for API validation.

This module defines Pydantic schemas for model-related API operations
with proper validation and serialization.
"""

import uuid
from decimal import Decimal
from typing import Optional

from pydantic import BaseModel, Field

from app.schemas.provider import ProviderRead


class ModelBase(BaseModel):
    """Base model schema with common fields."""
    model_name: str = Field(..., description="Model identifier")
    display_name: Optional[str] = Field(None, description="Human-readable model name")
    description: Optional[str] = Field(None, description="Model description")
    context_length: Optional[int] = Field(None, ge=0, description="Maximum context length")
    max_output_tokens: Optional[int] = Field(None, ge=0, description="Maximum output tokens")
    supports_streaming: bool = Field(True, description="Streaming support")
    supports_function_calling: bool = Field(False, description="Function calling support")
    supports_vision: bool = Field(False, description="Vision/image support")


class ModelCreate(ModelBase):
    """Model creation schema."""
    provider_id: uuid.UUID = Field(..., description="Provider ID")
    price_per_1k_prompt_tokens: Optional[Decimal] = Field(None, ge=0, description="Prompt token price")
    price_per_1k_completion_tokens: Optional[Decimal] = Field(None, ge=0, description="Completion token price")
    price_per_1k_tokens: Optional[Decimal] = Field(None, ge=0, description="Unified token price")


class ModelUpdate(BaseModel):
    """Model update schema."""
    display_name: Optional[str] = Field(None, description="Human-readable model name")
    description: Optional[str] = Field(None, description="Model description")
    is_active: Optional[bool] = Field(None, description="Model availability")
    price_per_1k_prompt_tokens: Optional[Decimal] = Field(None, ge=0)
    price_per_1k_completion_tokens: Optional[Decimal] = Field(None, ge=0)
    price_per_1k_tokens: Optional[Decimal] = Field(None, ge=0)


class ModelRead(ModelBase):
    """Model response schema."""
    id: uuid.UUID
    provider_id: uuid.UUID
    is_active: bool
    price_per_1k_prompt_tokens: Optional[Decimal]
    price_per_1k_completion_tokens: Optional[Decimal]
    price_per_1k_tokens: Optional[Decimal]
    created_at: str
    updated_at: str
    
    class Config:
        from_attributes = True


class ModelWithProvider(ModelRead):
    """Model with provider information."""
    provider: ProviderRead


class ModelStats(BaseModel):
    """Model usage statistics."""
    model_id: uuid.UUID
    model_name: str
    total_requests: int
    total_tokens: int
    total_cost: Decimal
    average_tokens_per_request: float
    
    class Config:
        from_attributes = True