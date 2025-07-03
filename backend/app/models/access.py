"""
User model access and permission management.

This module defines the access control system for distributing model permissions
to users with quotas and time-based restrictions.
"""

import uuid
from datetime import datetime
from typing import Optional

from sqlalchemy import Boolean, DateTime, ForeignKey, Integer, UniqueConstraint
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.db.base import Base


class UserModelAccess(Base):
    """
    User access permissions for specific models.
    
    This table manages the distribution of ORGANIZATION-scoped models
    to users with configurable quotas and access control.
    """
    
    __tablename__ = "user_model_access"
    
    # Core relationships
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
    
    # Access control
    is_active: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)
    expires_at: Mapped[Optional[datetime]] = mapped_column(
        DateTime(timezone=True),
        nullable=True,
        comment="NULL means permanent access"
    )
    
    # Quota management
    daily_quota: Mapped[Optional[int]] = mapped_column(
        Integer,
        nullable=True,
        comment="Daily API call limit, NULL means unlimited"
    )
    monthly_quota: Mapped[Optional[int]] = mapped_column(
        Integer,
        nullable=True,
        comment="Monthly API call limit, NULL means unlimited"
    )
    
    # Audit trail
    granted_by: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        ForeignKey("user.id"),
        nullable=False,
        comment="Admin who granted this access"
    )
    
    # Relationships
    user = relationship("User", foreign_keys=[user_id], back_populates="model_accesses")
    model = relationship("Model", back_populates="user_accesses")
    grantor = relationship("User", foreign_keys=[granted_by], back_populates="granted_accesses")
    
    # Constraints
    __table_args__ = (
        UniqueConstraint("user_id", "model_id", name="uq_user_model_access"),
    )
    
    @property
    def is_expired(self) -> bool:
        """Check if access has expired."""
        if self.expires_at is None:
            return False
        return datetime.utcnow() > self.expires_at
    
    @property
    def is_valid(self) -> bool:
        """Check if access is valid (active and not expired)."""
        return self.is_active and not self.is_expired
    
    @property
    def has_quota_limits(self) -> bool:
        """Check if access has any quota limitations."""
        return self.daily_quota is not None or self.monthly_quota is not None
    
    def check_quota_remaining(self, current_usage: dict) -> dict:
        """
        Check remaining quota based on current usage.
        
        Args:
            current_usage: Dict with 'daily' and 'monthly' usage counts
            
        Returns:
            Dict with remaining quotas and limits
        """
        result = {
            "daily_limit": self.daily_quota,
            "monthly_limit": self.monthly_quota,
            "daily_used": current_usage.get("daily", 0),
            "monthly_used": current_usage.get("monthly", 0),
            "daily_remaining": None,
            "monthly_remaining": None,
            "can_use": True
        }
        
        if self.daily_quota is not None:
            result["daily_remaining"] = max(0, self.daily_quota - result["daily_used"])
            if result["daily_remaining"] <= 0:
                result["can_use"] = False
        
        if self.monthly_quota is not None:
            result["monthly_remaining"] = max(0, self.monthly_quota - result["monthly_used"])
            if result["monthly_remaining"] <= 0:
                result["can_use"] = False
        
        return result
    
    def __repr__(self) -> str:
        return f"<UserModelAccess(user_id={self.user_id}, model_id={self.model_id}, active={self.is_active})>"