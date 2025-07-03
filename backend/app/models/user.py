"""
User model and related entities.

This module defines the user model with role-based access control (RBAC)
and integrates with fastapi-users for authentication.
"""

from enum import Enum
from typing import Optional

from fastapi_users.db import SQLAlchemyBaseUserTableUUID
from sqlalchemy import Boolean, String
from sqlalchemy.orm import Mapped, mapped_column, relationship

from app.db.base import Base


class UserRole(str, Enum):
    """
    User roles for RBAC system.
    
    - USER: Regular user with basic access
    - ADMIN: Administrator with elevated privileges  
    - SUPER_ADMIN: Super administrator with full system access
    """
    USER = "user"
    ADMIN = "admin"
    SUPER_ADMIN = "super_admin"


class User(SQLAlchemyBaseUserTableUUID, Base):
    """
    User model with RBAC support.
    
    Extends fastapi-users base user model with additional fields
    for role-based access control and user management.
    """
    
    # Role-based access control
    role: Mapped[UserRole] = mapped_column(
        String(20),
        default=UserRole.USER,
        nullable=False,
        index=True
    )
    
    # User profile information
    first_name: Mapped[Optional[str]] = mapped_column(String(100), nullable=True)
    last_name: Mapped[Optional[str]] = mapped_column(String(100), nullable=True)
    
    # Account settings
    is_active: Mapped[bool] = mapped_column(Boolean, default=True, nullable=False)
    is_superuser: Mapped[bool] = mapped_column(Boolean, default=False, nullable=False)
    is_verified: Mapped[bool] = mapped_column(Boolean, default=False, nullable=False)
    
    # Relationships
    providers = relationship(
        "Provider",
        back_populates="owner",
        cascade="all, delete-orphan"
    )
    
    model_accesses = relationship(
        "UserModelAccess",
        back_populates="user",
        cascade="all, delete-orphan"
    )
    
    usage_logs = relationship(
        "UsageLog",
        back_populates="user",
        cascade="all, delete-orphan"
    )
    
    uploaded_files = relationship(
        "UploadedFile",
        back_populates="user",
        cascade="all, delete-orphan"
    )
    
    granted_accesses = relationship(
        "UserModelAccess",
        foreign_keys="UserModelAccess.granted_by",
        back_populates="grantor",
        cascade="all, delete-orphan"
    )
    
    @property
    def full_name(self) -> str:
        """Get user's full name."""
        if self.first_name and self.last_name:
            return f"{self.first_name} {self.last_name}"
        return self.email
    
    @property
    def is_admin(self) -> bool:
        """Check if user is admin or super admin."""
        return self.role in [UserRole.ADMIN, UserRole.SUPER_ADMIN]
    
    @property
    def is_super_admin(self) -> bool:
        """Check if user is super admin."""
        return self.role == UserRole.SUPER_ADMIN
    
    def __repr__(self) -> str:
        return f"<User(id={self.id}, email={self.email}, role={self.role})>"