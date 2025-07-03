"""
User schemas for API validation.

This module defines Pydantic schemas for user-related API operations
with proper validation and serialization.
"""

import uuid
from typing import Optional

from fastapi_users import schemas
from pydantic import BaseModel, EmailStr

from app.models import UserRole


class UserRead(schemas.BaseUser[uuid.UUID]):
    """
    User response schema for public API.
    
    Includes safe user information without sensitive data.
    """
    role: UserRole
    first_name: Optional[str] = None
    last_name: Optional[str] = None
    full_name: str
    is_admin: bool
    is_super_admin: bool


class UserCreate(schemas.BaseUserCreate):
    """
    User creation schema.
    
    Defines required fields for user registration.
    """
    email: EmailStr
    password: str
    first_name: Optional[str] = None
    last_name: Optional[str] = None
    role: UserRole = UserRole.USER  # Default role


class UserUpdate(schemas.BaseUserUpdate):
    """
    User update schema.
    
    Allows updating user profile information.
    """
    email: Optional[EmailStr] = None
    password: Optional[str] = None
    first_name: Optional[str] = None
    last_name: Optional[str] = None
    is_active: Optional[bool] = None
    is_verified: Optional[bool] = None


class UserRoleUpdate(BaseModel):
    """
    User role update schema (admin only).
    
    Allows admins to change user roles.
    """
    role: UserRole


class UserStats(BaseModel):
    """
    User statistics schema.
    
    Provides usage and cost statistics for users.
    """
    total_requests: int
    total_tokens: int
    total_cost: float
    period_days: int
    
    class Config:
        from_attributes = True