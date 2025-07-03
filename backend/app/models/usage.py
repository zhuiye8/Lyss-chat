"""
Usage tracking and cost management.

This module defines models for tracking API usage, calculating costs,
and managing user consumption statistics.
"""

import uuid
from decimal import Decimal
from typing import Optional

from sqlalchemy import DECIMAL, CheckConstraint, ForeignKey, Integer
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.db.base import Base


class UsageLog(Base):
    """
    API usage tracking for cost calculation and quota management.
    
    Records every AI API call with token usage and calculated costs
    for billing and monitoring purposes.
    """
    
    __tablename__ = "usage_log"
    
    # User and model relationships
    user_id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("user.id", ondelete="CASCADE"),
        nullable=False,
        index=True
    )
    
    model_id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("model.id", ondelete="CASCADE"),
        nullable=False,
        index=True
    )
    
    # Token usage details
    prompt_tokens: Mapped[int] = mapped_column(Integer, nullable=False)
    completion_tokens: Mapped[int] = mapped_column(Integer, nullable=False)
    total_tokens: Mapped[int] = mapped_column(Integer, nullable=False)
    
    # Cost calculation (in USD)
    cost: Mapped[Decimal] = mapped_column(
        DECIMAL(10, 6),
        nullable=False,
        default=Decimal("0.000000")
    )
    
    # Optional request metadata
    request_id: Mapped[Optional[str]] = mapped_column(
        "request_id",
        nullable=True,
        comment="External request ID for tracing"
    )
    
    # Response metadata
    response_time_ms: Mapped[Optional[int]] = mapped_column(
        Integer,
        nullable=True,
        comment="API response time in milliseconds"
    )
    
    # Relationships
    user = relationship("User", back_populates="usage_logs")
    model = relationship("Model", back_populates="usage_logs")
    
    # Data integrity constraints
    __table_args__ = (
        CheckConstraint(
            "total_tokens = prompt_tokens + completion_tokens",
            name="check_tokens_sum"
        ),
        CheckConstraint("prompt_tokens >= 0", name="check_prompt_tokens_positive"),
        CheckConstraint("completion_tokens >= 0", name="check_completion_tokens_positive"),
        CheckConstraint("cost >= 0", name="check_cost_positive"),
    )
    
    @property
    def cost_per_token(self) -> Decimal:
        """Calculate cost per token for this usage."""
        if self.total_tokens > 0:
            return self.cost / self.total_tokens
        return Decimal("0.000000")
    
    @property
    def efficiency_score(self) -> float:
        """
        Calculate efficiency score based on completion ratio.
        
        Higher completion tokens relative to prompt tokens
        generally indicate more efficient usage.
        """
        if self.total_tokens > 0:
            return float(self.completion_tokens / self.total_tokens)
        return 0.0
    
    def __repr__(self) -> str:
        return (
            f"<UsageLog(user_id={self.user_id}, model_id={self.model_id}, "
            f"total_tokens={self.total_tokens}, cost={self.cost})>"
        )