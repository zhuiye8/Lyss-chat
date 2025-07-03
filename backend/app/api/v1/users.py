"""
User management API endpoints.

This module provides user CRUD operations and profile management
with role-based access control.
"""

from typing import List

from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.auth import current_active_user, fastapi_users
from app.core.permissions import require_admin, require_super_admin
from app.db import get_session
from app.models import User
from app.schemas.user import UserRead, UserRoleUpdate, UserStats, UserUpdate
from app.services.user_service import UserService

router = APIRouter()
user_service = UserService()

# Include fastapi-users user management routes
router.include_router(
    fastapi_users.get_users_router(UserRead, UserUpdate),
    prefix="/users",
    tags=["users"]
)


@router.get("/users", response_model=List[UserRead], tags=["users"])
async def list_users(
    skip: int = 0,
    limit: int = 100,
    db: AsyncSession = Depends(get_session),
    current_user: User = Depends(require_admin)
):
    """
    List all users (admin only).
    
    Retrieve a paginated list of all users in the system.
    Only admins and super admins can access this endpoint.
    """
    users = await user_service.get_users(db, skip=skip, limit=limit)
    return users


@router.put("/users/{user_id}/role", response_model=UserRead, tags=["users"])
async def update_user_role(
    user_id: str,
    role_update: UserRoleUpdate,
    db: AsyncSession = Depends(get_session),
    current_user: User = Depends(require_super_admin)
):
    """
    Update user role (super admin only).
    
    Change a user's role. Only super admins can perform this operation
    to prevent privilege escalation.
    """
    user = await user_service.get_user_by_id(db, user_id)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found"
        )
    
    # Prevent self-demotion of the last super admin
    if (current_user.id == user.id and 
        user.is_super_admin and 
        role_update.role != user.role):
        remaining_super_admins = await user_service.count_super_admins(db)
        if remaining_super_admins <= 1:
            raise HTTPException(
                status_code=status.HTTP_400_BAD_REQUEST,
                detail="Cannot demote the last super admin"
            )
    
    updated_user = await user_service.update_user_role(db, user, role_update.role)
    return updated_user


@router.get("/users/{user_id}/stats", response_model=UserStats, tags=["users"])
async def get_user_stats(
    user_id: str,
    days: int = 30,
    db: AsyncSession = Depends(get_session),
    current_user: User = Depends(current_active_user)
):
    """
    Get user usage statistics.
    
    Retrieve usage statistics for a specific user.
    Users can only access their own stats unless they are admins.
    """
    # Check permissions
    if str(current_user.id) != user_id and not current_user.is_admin:
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to view these statistics"
        )
    
    user = await user_service.get_user_by_id(db, user_id)
    if not user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found"
        )
    
    stats = await user_service.get_user_usage_stats(db, user.id, days)
    return stats


@router.get("/me/stats", response_model=UserStats, tags=["users"])
async def get_my_stats(
    days: int = 30,
    db: AsyncSession = Depends(get_session),
    current_user: User = Depends(current_active_user)
):
    """
    Get current user's usage statistics.
    
    Convenient endpoint for users to get their own usage statistics.
    """
    stats = await user_service.get_user_usage_stats(db, current_user.id, days)
    return stats