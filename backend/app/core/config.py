"""
Core configuration module for LYSS AI Platform.

This module handles all application settings, environment variables,
and configuration validation using Pydantic Settings.
"""

import secrets
from typing import Any, Dict, List, Optional, Union

from pydantic import AnyHttpUrl, EmailStr, field_validator
from pydantic_core import MultiHostUrl
from pydantic.networks import PostgresDsn
from pydantic_settings import BaseSettings


class Settings(BaseSettings):
    """
    Application settings loaded from environment variables.
    
    All settings follow the principle of explicit configuration
    with secure defaults where possible.
    """
    
    # === Application Settings ===
    PROJECT_NAME: str = "LYSS AI Platform"
    VERSION: str = "0.1.0"
    API_V1_STR: str = "/api/v1"
    SECRET_KEY: str = secrets.token_urlsafe(32)
    AUTH_SECRET: str = secrets.token_urlsafe(32)
    
    # === Environment ===
    ENVIRONMENT: str = "development"  # development, staging, production
    DEBUG: bool = True
    
    # === Security ===
    ACCESS_TOKEN_EXPIRE_MINUTES: int = 60 * 24 * 8  # 8 days
    REFRESH_TOKEN_EXPIRE_MINUTES: int = 60 * 24 * 30  # 30 days
    ENCRYPTION_KEY: str = secrets.token_urlsafe(32)
    
    # === Database ===
    POSTGRES_SERVER: str = "localhost"
    POSTGRES_PORT: int = 5432
    POSTGRES_USER: str = "lyss"
    POSTGRES_PASSWORD: str = ""
    POSTGRES_DB: str = "lyss_db"
    DATABASE_URL: Optional[PostgresDsn] = None
    
    @field_validator("DATABASE_URL", mode="before")
    @classmethod
    def assemble_db_connection(cls, v: Optional[str], values: Dict[str, Any]) -> Any:
        """Assemble database URL from individual components if not provided."""
        if isinstance(v, str):
            return v
        
        # Use the class method parameters to access field values
        return MultiHostUrl.build(
            scheme="postgresql+asyncpg",
            username=values.get("POSTGRES_USER"),
            password=values.get("POSTGRES_PASSWORD"),
            host=values.get("POSTGRES_SERVER"),
            port=values.get("POSTGRES_PORT"),
            path=values.get("POSTGRES_DB") or "",
        )
    
    # === Redis ===
    REDIS_URL: str = "redis://localhost:6379"
    
    # === Qdrant Vector Database ===
    QDRANT_URL: str = "http://localhost:6333"
    
    # === Mem0 Configuration ===
    # Global LLM and Embedding models for memory processing
    MEM0_LLM_PROVIDER: str = "openai"
    MEM0_LLM_MODEL: str = "gpt-4o-mini"
    MEM0_EMBEDDING_MODEL: str = "text-embedding-3-small"
    MEM0_PROVIDER_API_KEY: str = ""
    MEM0_PROVIDER_BASE_URL: Optional[str] = "https://api.openai.com/v1"
    
    # === File Upload Configuration ===
    MAX_FILE_SIZE_MB: int = 50
    ALLOWED_FILE_TYPES: str = "pdf,docx,txt,md"
    UPLOAD_DIR: str = "/tmp/uploads"
    
    def get_allowed_file_types(self) -> List[str]:
        """Get allowed file types as a list."""
        if isinstance(self.ALLOWED_FILE_TYPES, str):
            return [t.strip() for t in self.ALLOWED_FILE_TYPES.split(",")]
        return self.ALLOWED_FILE_TYPES
    
    # === CORS ===
    BACKEND_CORS_ORIGINS: str = "http://localhost:3000,http://localhost:8000,http://localhost:8001"
    
    def get_cors_origins(self) -> List[str]:
        """Get CORS origins as a list."""
        if isinstance(self.BACKEND_CORS_ORIGINS, str):
            return [origin.strip() for origin in self.BACKEND_CORS_ORIGINS.split(",")]
        return self.BACKEND_CORS_ORIGINS
    
    # === Logging ===
    LOG_LEVEL: str = "INFO"
    LOG_FILE: Optional[str] = None
    
    # === Admin User ===
    FIRST_SUPERUSER: EmailStr = "admin@lyss.ai"
    FIRST_SUPERUSER_PASSWORD: str = "changeme"
    
    class Config:
        env_file = ".env"
        case_sensitive = True


# Global settings instance
settings = Settings()


def get_settings() -> Settings:
    """Get application settings instance."""
    return settings