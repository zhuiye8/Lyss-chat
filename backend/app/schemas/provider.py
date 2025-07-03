"""
Provider schemas for API validation.

This module defines Pydantic schemas for provider-related API operations
with proper validation and serialization.
"""

import uuid
from typing import Any, Dict, List, Optional

from pydantic import BaseModel, Field, field_validator

from app.models import ProviderScope


class ProviderBase(BaseModel):
    """Base provider schema with common fields."""
    name: str = Field(..., min_length=1, max_length=255, description="Provider display name")
    provider_type: str = Field(..., description="Provider type (e.g., 'openai', 'anthropic')")
    scope: ProviderScope = Field(..., description="Provider scope: ORGANIZATION or PERSONAL")
    description: Optional[str] = Field(None, max_length=1000, description="Provider description")


class ProviderCreate(ProviderBase):
    """Provider creation schema."""
    config: Dict[str, Any] = Field(..., description="Provider configuration (will be encrypted)")
    
    @field_validator('config')
    @classmethod
    def validate_config_not_empty(cls, v):
        if not v:
            raise ValueError('Configuration cannot be empty')
        return v


class ProviderUpdate(BaseModel):
    """Provider update schema."""
    name: Optional[str] = Field(None, min_length=1, max_length=255)
    description: Optional[str] = Field(None, max_length=1000)
    is_enabled: Optional[bool] = Field(None, description="Enable/disable provider")
    config: Optional[Dict[str, Any]] = Field(None, description="Updated configuration")


class ProviderRead(ProviderBase):
    """Provider response schema."""
    id: uuid.UUID
    owner_id: uuid.UUID
    is_enabled: bool
    created_at: str
    updated_at: str
    
    class Config:
        from_attributes = True


class ProviderWithOwner(ProviderRead):
    """Provider with owner information."""
    owner_email: str
    owner_name: str


class ProviderTestRequest(BaseModel):
    """Provider connection test request."""
    provider_type: str = Field(..., description="Provider type to test")
    config: Dict[str, Any] = Field(..., description="Configuration to test")


class ProviderTestResponse(BaseModel):
    """Provider connection test response."""
    success: bool
    message: str
    error_details: Optional[str] = None


class ProviderConfigSchema(BaseModel):
    """Provider configuration schema response."""
    provider_type: str
    config_schema: Dict[str, Any] = Field(..., description="JSON schema for provider configuration")
    
    class Config:
        from_attributes = True


class ProviderInfo(BaseModel):
    """Provider information and capabilities."""
    provider_type: str
    name: str
    description: str
    website: Optional[str] = None
    documentation: Optional[str] = None
    supports_streaming: bool = False
    supports_function_calling: bool = False
    supports_vision: bool = False
    
    class Config:
        from_attributes = True


class ModelSyncResponse(BaseModel):
    """Model synchronization response."""
    provider_id: uuid.UUID
    models_synced: int
    models_added: int
    models_updated: int
    errors: List[str] = []