"""
Authentication configuration using fastapi-users.

This module sets up JWT-based authentication with role-based access control
and provides user management functionality.
"""

import uuid
from typing import Optional

from fastapi import Depends, Request
from fastapi_users import BaseUserManager, FastAPIUsers, UUIDIDMixin
from fastapi_users.authentication import (
    AuthenticationBackend,
    BearerTransport,
    JWTStrategy,
)
from fastapi_users.db import SQLAlchemyUserDatabase

from app.core.config import settings
from app.db import get_session
from app.models import User


class UserManager(UUIDIDMixin, BaseUserManager[User, uuid.UUID]):
    """
    User manager with custom business logic.
    
    Handles user registration, verification, and role management
    according to LYSS AI Platform requirements.
    """
    
    reset_password_token_secret = settings.AUTH_SECRET
    verification_token_secret = settings.AUTH_SECRET
    
    async def on_after_register(self, user: User, request: Optional[Request] = None):
        """Called after user registration."""
        print(f"User {user.id} has registered.")
    
    async def on_after_login(
        self,
        user: User,
        request: Optional[Request] = None,
    ):
        """Called after user login."""
        print(f"User {user.id} has logged in.")
    
    async def on_after_forgot_password(
        self, user: User, token: str, request: Optional[Request] = None
    ):
        """Called after password reset request."""
        print(f"User {user.id} has requested password reset. Token: {token}")
    
    async def on_after_request_verify(
        self, user: User, token: str, request: Optional[Request] = None
    ):
        """Called after verification request."""
        print(f"Verification requested for user {user.id}. Token: {token}")


async def get_user_db(session=Depends(get_session)):
    """Get user database dependency."""
    yield SQLAlchemyUserDatabase(session, User)


async def get_user_manager(user_db=Depends(get_user_db)):
    """Get user manager dependency."""
    yield UserManager(user_db)


# Authentication backend configuration
bearer_transport = BearerTransport(tokenUrl="auth/jwt/login")

def get_jwt_strategy() -> JWTStrategy:
    """Get JWT authentication strategy."""
    return JWTStrategy(
        secret=settings.SECRET_KEY,
        lifetime_seconds=settings.ACCESS_TOKEN_EXPIRE_MINUTES * 60,
    )

auth_backend = AuthenticationBackend(
    name="jwt",
    transport=bearer_transport,
    get_strategy=get_jwt_strategy,
)

# FastAPI Users instance
fastapi_users = FastAPIUsers[User, uuid.UUID](get_user_manager, [auth_backend])

# Dependencies for protected routes
current_active_user = fastapi_users.current_user(active=True)
current_active_verified_user = fastapi_users.current_user(active=True, verified=True)