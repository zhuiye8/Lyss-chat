"""
User service for business logic operations.

This module contains business logic for user management,
statistics, and role-based operations.
"""

import uuid
from datetime import datetime, timedelta
from typing import List, Optional

from sqlalchemy import func, select
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from app.models import User, UserRole, UsageLog


class UserService:
    """
    Service class for user-related business operations.
    
    Encapsulates business logic for user management, statistics,
    and administrative operations.
    """
    
    async def get_user_by_id(self, db: AsyncSession, user_id: str) -> Optional[User]:
        """
        Get user by ID.
        
        Args:
            db: Database session
            user_id: User UUID as string
            
        Returns:
            User object or None if not found
        """
        try:
            user_uuid = uuid.UUID(user_id)
            result = await db.execute(
                select(User).where(User.id == user_uuid)
            )
            return result.scalar_one_or_none()
        except ValueError:
            return None
    
    async def get_users(
        self, 
        db: AsyncSession, 
        skip: int = 0, 
        limit: int = 100
    ) -> List[User]:
        """
        Get paginated list of users.
        
        Args:
            db: Database session
            skip: Number of records to skip
            limit: Maximum number of records to return
            
        Returns:
            List of users
        """
        result = await db.execute(
            select(User)
            .options(selectinload(User.providers))
            .offset(skip)
            .limit(limit)
            .order_by(User.created_at.desc())
        )
        return list(result.scalars().all())
    
    async def update_user_role(
        self, 
        db: AsyncSession, 
        user: User, 
        new_role: UserRole
    ) -> User:
        """
        Update user role.
        
        Args:
            db: Database session
            user: User to update
            new_role: New role to assign
            
        Returns:
            Updated user
        """
        user.role = new_role
        user.updated_at = datetime.utcnow()
        
        # Update superuser flag based on role
        user.is_superuser = new_role == UserRole.SUPER_ADMIN
        
        await db.commit()
        await db.refresh(user)
        return user
    
    async def count_super_admins(self, db: AsyncSession) -> int:
        """
        Count total number of super admins.
        
        Args:
            db: Database session
            
        Returns:
            Number of super admin users
        """
        result = await db.execute(
            select(func.count(User.id)).where(User.role == UserRole.SUPER_ADMIN)
        )
        return result.scalar() or 0
    
    async def get_user_usage_stats(
        self, 
        db: AsyncSession, 
        user_id: uuid.UUID, 
        days: int = 30
    ) -> dict:
        """
        Get user usage statistics for specified period.
        
        Args:
            db: Database session
            user_id: User UUID
            days: Number of days to include in statistics
            
        Returns:
            Dictionary with usage statistics
        """
        start_date = datetime.utcnow() - timedelta(days=days)
        
        result = await db.execute(
            select(
                func.count(UsageLog.id).label('total_requests'),
                func.sum(UsageLog.total_tokens).label('total_tokens'),
                func.sum(UsageLog.cost).label('total_cost')
            )
            .where(
                UsageLog.user_id == user_id,
                UsageLog.created_at >= start_date
            )
        )
        
        stats = result.first()
        
        return {
            'total_requests': stats.total_requests or 0,
            'total_tokens': stats.total_tokens or 0,
            'total_cost': float(stats.total_cost or 0),
            'period_days': days
        }
    
    async def create_user(
        self,
        db: AsyncSession,
        email: str,
        password_hash: str,
        first_name: Optional[str] = None,
        last_name: Optional[str] = None,
        role: UserRole = UserRole.USER
    ) -> User:
        """
        Create a new user.
        
        Args:
            db: Database session
            email: User email
            password_hash: Hashed password
            first_name: User's first name
            last_name: User's last name
            role: User role
            
        Returns:
            Created user
        """
        user = User(
            email=email,
            hashed_password=password_hash,
            first_name=first_name,
            last_name=last_name,
            role=role,
            is_superuser=(role == UserRole.SUPER_ADMIN),
            is_active=True,
            is_verified=True  # Auto-verify for now
        )
        
        db.add(user)
        await db.commit()
        await db.refresh(user)
        return user
    
    async def get_users_by_role(
        self, 
        db: AsyncSession, 
        role: UserRole
    ) -> List[User]:
        """
        Get users by role.
        
        Args:
            db: Database session
            role: User role to filter by
            
        Returns:
            List of users with specified role
        """
        result = await db.execute(
            select(User).where(User.role == role).order_by(User.created_at.desc())
        )
        return list(result.scalars().all())