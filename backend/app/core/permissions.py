"""
Role-based access control (RBAC) permissions system.

This module provides decorators and dependencies for implementing
role-based access control throughout the application.
"""

from typing import List

from fastapi import Depends, HTTPException, status

from app.core.auth import current_active_user
from app.models import User, UserRole


class RoleChecker:
    """
    Role-based access control checker.
    
    This class creates dependencies that verify user roles
    for protected API endpoints.
    """
    
    def __init__(self, allowed_roles: List[UserRole]):
        """
        Initialize role checker.
        
        Args:
            allowed_roles: List of roles allowed to access the endpoint
        """
        self.allowed_roles = allowed_roles
    
    def __call__(self, user: User = Depends(current_active_user)) -> User:
        """
        Check if current user has required role.
        
        Args:
            user: Current authenticated user
            
        Returns:
            User if authorized
            
        Raises:
            HTTPException: If user doesn't have required role
        """
        if user.role not in self.allowed_roles:
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail=f"Operation requires one of roles: {[role.value for role in self.allowed_roles]}"
            )
        return user


# Pre-configured role checkers for common use cases
require_admin = RoleChecker([UserRole.ADMIN, UserRole.SUPER_ADMIN])
require_super_admin = RoleChecker([UserRole.SUPER_ADMIN])
require_user = RoleChecker([UserRole.USER, UserRole.ADMIN, UserRole.SUPER_ADMIN])


def check_resource_owner_or_admin(resource_owner_id: str, current_user: User) -> bool:
    """
    Check if user owns resource or is admin.
    
    Args:
        resource_owner_id: ID of the resource owner
        current_user: Current authenticated user
        
    Returns:
        True if user can access resource, False otherwise
    """
    return (
        str(current_user.id) == resource_owner_id 
        or current_user.is_admin
    )


def check_organization_scope_access(current_user: User) -> bool:
    """
    Check if user can create ORGANIZATION scoped resources.
    
    Args:
        current_user: Current authenticated user
        
    Returns:
        True if user can create organization resources
    """
    return current_user.is_admin


class ProviderScopeChecker:
    """
    Provider scope-based access checker.
    
    Validates user permissions based on provider scope
    according to V2 architecture rules.
    """
    
    @staticmethod
    def can_create_organization_provider(user: User) -> bool:
        """Check if user can create ORGANIZATION providers."""
        return user.is_admin
    
    @staticmethod
    def can_create_personal_provider(user: User) -> bool:
        """Check if user can create PERSONAL providers."""
        return True  # Any authenticated user can create personal providers
    
    @staticmethod
    def can_view_provider(provider_owner_id: str, user: User) -> bool:
        """Check if user can view specific provider."""
        return check_resource_owner_or_admin(provider_owner_id, user)
    
    @staticmethod
    def can_modify_provider(provider_owner_id: str, user: User) -> bool:
        """Check if user can modify specific provider."""
        return check_resource_owner_or_admin(provider_owner_id, user)
    
    @staticmethod
    def can_delete_provider(provider_owner_id: str, user: User) -> bool:
        """Check if user can delete specific provider."""
        return (
            check_resource_owner_or_admin(provider_owner_id, user)
            or user.is_super_admin
        )


class ModelAccessChecker:
    """
    Model access permission checker.
    
    Validates user access to models based on scope and distribution rules.
    """
    
    @staticmethod
    def can_distribute_model(user: User) -> bool:
        """Check if user can distribute models to other users."""
        return user.is_admin
    
    @staticmethod
    def can_access_personal_model(provider_owner_id: str, user: User) -> bool:
        """Check if user can access personal provider models."""
        return str(user.id) == provider_owner_id
    
    @staticmethod
    def can_view_usage_stats(target_user_id: str, user: User) -> bool:
        """Check if user can view usage statistics."""
        if user.is_super_admin:
            return True
        if user.is_admin:
            return True  # Admins can view org usage stats
        return str(user.id) == target_user_id  # Users can view own stats