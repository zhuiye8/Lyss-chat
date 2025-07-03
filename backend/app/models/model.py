"""
Model and related entities.

This module defines AI models provided by various providers and their capabilities.
"""

import uuid
from decimal import Decimal
from typing import Optional

from sqlalchemy import DECIMAL, Boolean, ForeignKey, Integer, String, Text, UniqueConstraint
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.db.base import Base


class Model(Base):
    """
    AI Model provided by a specific provider.
    
    Represents individual AI models (like gpt-4, claude-3.5-sonnet)
    with their capabilities and pricing information.
    """
    
    __tablename__ = "model"
    
    # Provider relationship
    provider_id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("provider.id", ondelete="CASCADE"),
        nullable=False,
        index=True
    )
    
    # Model identification
    model_name: Mapped[str] = mapped_column(String(255), nullable=False, index=True)
    display_name: Mapped[Optional[str]] = mapped_column(String(255), nullable=True)
    
    # Model capabilities
    context_length: Mapped[Optional[int]] = mapped_column(Integer, nullable=True)
    max_output_tokens: Mapped[Optional[int]] = mapped_column(Integer, nullable=True)
    supports_streaming: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)
    supports_function_calling: Mapped[bool] = mapped_column(Boolean, default=False, nullable=False)
    supports_vision: Mapped[bool] = mapped_column(Boolean, default=False, nullable=False)
    
    # Pricing information (per 1K tokens in USD)
    price_per_1k_prompt_tokens: Mapped[Optional[Decimal]] = mapped_column(
        DECIMAL(10, 6), nullable=True
    )
    price_per_1k_completion_tokens: Mapped[Optional[Decimal]] = mapped_column(
        DECIMAL(10, 6), nullable=True
    )
    price_per_1k_tokens: Mapped[Optional[Decimal]] = mapped_column(
        DECIMAL(10, 6), nullable=True, comment="Unified pricing for models with same input/output cost"
    )
    
    # Model metadata
    description: Mapped[Optional[str]] = mapped_column(Text, nullable=True)
    is_active: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)
    
    # Relationships
    provider = relationship("Provider", back_populates="models")
    user_accesses = relationship(
        "UserModelAccess",
        back_populates="model",
        cascade="all, delete-orphan"
    )
    usage_logs = relationship(
        "UsageLog",
        back_populates="model",
        cascade="all, delete-orphan"
    )
    
    # Constraints
    __table_args__ = (
        UniqueConstraint("provider_id", "model_name", name="uq_provider_model"),
    )
    
    @property
    def full_name(self) -> str:
        """Get model's full display name."""
        if self.display_name:
            return self.display_name
        return self.model_name
    
    @property
    def has_differential_pricing(self) -> bool:
        """Check if model has different pricing for prompt and completion tokens."""
        return (
            self.price_per_1k_prompt_tokens is not None 
            and self.price_per_1k_completion_tokens is not None
        )
    
    @property
    def has_unified_pricing(self) -> bool:
        """Check if model has unified pricing for all tokens."""
        return self.price_per_1k_tokens is not None
    
    def calculate_cost(self, prompt_tokens: int, completion_tokens: int) -> Decimal:
        """
        Calculate cost for token usage.
        
        Args:
            prompt_tokens: Number of prompt tokens
            completion_tokens: Number of completion tokens
            
        Returns:
            Cost in USD
        """
        if self.has_differential_pricing:
            prompt_cost = Decimal(prompt_tokens) * self.price_per_1k_prompt_tokens / 1000
            completion_cost = Decimal(completion_tokens) * self.price_per_1k_completion_tokens / 1000
            return prompt_cost + completion_cost
        elif self.has_unified_pricing:
            total_tokens = prompt_tokens + completion_tokens
            return Decimal(total_tokens) * self.price_per_1k_tokens / 1000
        else:
            return Decimal("0.00")
    
    def __repr__(self) -> str:
        return f"<Model(id={self.id}, name={self.model_name}, provider={self.provider_id})>"