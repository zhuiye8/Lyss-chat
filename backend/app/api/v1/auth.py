"""
Authentication API endpoints.

This module provides user authentication, registration,
and account management endpoints using fastapi-users.
"""

from fastapi import APIRouter

from app.core.auth import auth_backend, fastapi_users
from app.schemas.user import UserCreate, UserRead

# Authentication router
auth_router = APIRouter()

# Include fastapi-users authentication routes
auth_router.include_router(
    fastapi_users.get_auth_router(auth_backend),
    prefix="/jwt",
    tags=["auth"]
)

# Include user registration routes
auth_router.include_router(
    fastapi_users.get_register_router(UserRead, UserCreate),
    prefix="/register",
    tags=["auth"]
)

# Include password reset routes
auth_router.include_router(
    fastapi_users.get_reset_password_router(),
    prefix="/reset-password",
    tags=["auth"]
)

# Include email verification routes
auth_router.include_router(
    fastapi_users.get_verify_router(UserRead),
    prefix="/verify",
    tags=["auth"]
)