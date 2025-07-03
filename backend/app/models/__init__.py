"""
Database models for LYSS AI Platform.

This module exports all database models for easy importing
and ensures proper model registration with SQLAlchemy.
"""

from .access import UserModelAccess
from .file import UploadedFile
# Alias for convenience in database imports
File = UploadedFile
from .model import Model
from .provider import Provider, ProviderScope
from .usage import UsageLog
from .user import User, UserRole

# Export all models for easy importing
__all__ = [
    "User",
    "UserRole",
    "Provider",
    "ProviderScope",
    "Model",
    "UserModelAccess",
    "UsageLog",
    "UploadedFile",
    "File",
]