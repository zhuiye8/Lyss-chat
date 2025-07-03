"""
Security utilities for LYSS AI Platform.

This module provides encryption, decryption, password hashing,
and other security-related utilities.
"""

import base64
import secrets
from typing import Any, Union

from cryptography.fernet import Fernet
from passlib.context import CryptContext

from app.core.config import settings

# Password hashing context
pwd_context = CryptContext(schemes=["bcrypt"], deprecated="auto")

# Encryption instance for sensitive data
_encryption_key = base64.urlsafe_b64encode(settings.ENCRYPTION_KEY.encode()[:32].ljust(32, b'\0'))
fernet = Fernet(_encryption_key)


def create_access_token(subject: Union[str, Any]) -> str:
    """
    Create JWT access token.
    
    Args:
        subject: Token subject (usually user ID)
        
    Returns:
        Encoded JWT token
    """
    # This will be implemented with fastapi-users
    # For now, return a placeholder
    return f"placeholder_token_{subject}"


def verify_password(plain_password: str, hashed_password: str) -> bool:
    """
    Verify a password against its hash.
    
    Args:
        plain_password: Plain text password
        hashed_password: Bcrypt hashed password
        
    Returns:
        True if password matches, False otherwise
    """
    return pwd_context.verify(plain_password, hashed_password)


def get_password_hash(password: str) -> str:
    """
    Generate password hash.
    
    Args:
        password: Plain text password
        
    Returns:
        Bcrypt hashed password
    """
    return pwd_context.hash(password)


def encrypt_sensitive_data(data: str) -> str:
    """
    Encrypt sensitive data like API keys.
    
    Args:
        data: Sensitive data to encrypt
        
    Returns:
        Base64 encoded encrypted data
    """
    if not data:
        return ""
    
    encrypted_data = fernet.encrypt(data.encode())
    return base64.b64encode(encrypted_data).decode()


def decrypt_sensitive_data(encrypted_data: str) -> str:
    """
    Decrypt sensitive data like API keys.
    
    Args:
        encrypted_data: Base64 encoded encrypted data
        
    Returns:
        Decrypted plain text data
    """
    if not encrypted_data:
        return ""
    
    try:
        encrypted_bytes = base64.b64decode(encrypted_data.encode())
        decrypted_data = fernet.decrypt(encrypted_bytes)
        return decrypted_data.decode()
    except Exception:
        # Handle decryption errors gracefully
        raise ValueError("Failed to decrypt sensitive data")


def generate_api_key() -> str:
    """
    Generate a secure API key.
    
    Returns:
        Secure random API key
    """
    return secrets.token_urlsafe(32)


def generate_reset_token() -> str:
    """
    Generate a password reset token.
    
    Returns:
        Secure random reset token
    """
    return secrets.token_urlsafe(32)