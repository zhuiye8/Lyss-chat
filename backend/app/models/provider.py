"""
Provider model and related entities.

This module defines the AI provider model with scope-based access control
and encrypted configuration storage.
"""

import uuid
from enum import Enum
from typing import Optional

from sqlalchemy import Boolean, ForeignKey, String, Text
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.db.base import Base


class ProviderScope(str, Enum):
    """
    Provider scope determines ownership and distribution capabilities.
    
    - ORGANIZATION: Can be created by admins and distributed to other users
    - PERSONAL: Created by any user for personal use only, cannot be distributed
    """
    ORGANIZATION = "ORGANIZATION"
    PERSONAL = "PERSONAL"


class Provider(Base):
    """
    AI Provider model with scope-based access control.
    
    Stores configuration for AI service providers (OpenAI, Anthropic, etc.)
    with encrypted sensitive data and flexible configuration storage.
    """
    
    __tablename__ = "provider"
    
    # Basic provider information
    name: Mapped[str] = mapped_column(String(255), nullable=False, index=True)
    provider_type: Mapped[str] = mapped_column(String(100), nullable=False, index=True)
    
    # Scope-based access control (V2 core architecture)
    scope: Mapped[ProviderScope] = mapped_column(
        String(20),
        nullable=False,
        index=True
    )
    
    # Owner relationship
    owner_id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("user.id", ondelete="CASCADE"),
        nullable=False,
        index=True
    )
    
    # Encrypted configuration storage (V2.1 optimization)
    config_encrypted: Mapped[str] = mapped_column(
        Text,
        nullable=False,
        comment="Encrypted JSON configuration including API keys, URLs, etc."
    )
    
    # Provider status
    is_enabled: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)
    
    # Optional metadata
    description: Mapped[Optional[str]] = mapped_column(Text, nullable=True)
    
    # Relationships
    owner = relationship("User", back_populates="providers")
    models = relationship(
        "Model",
        back_populates="provider",
        cascade="all, delete-orphan"
    )
    
    @property
    def is_organization_scope(self) -> bool:
        """Check if provider has organization scope."""
        return self.scope == ProviderScope.ORGANIZATION
    
    @property
    def is_personal_scope(self) -> bool:
        """Check if provider has personal scope."""
        return self.scope == ProviderScope.PERSONAL
    
    @property
    def can_be_distributed(self) -> bool:
        """Check if provider models can be distributed to other users."""
        return self.is_organization_scope
    
    def __repr__(self) -> str:
        return f"<Provider(id={self.id}, name={self.name}, type={self.provider_type}, scope={self.scope})>"