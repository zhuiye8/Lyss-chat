"""
Database base configuration and session management.

This module sets up SQLAlchemy database connections, session management,
and provides the base class for all database models.
"""

import uuid
from datetime import datetime
from typing import Any

from sqlalchemy import DateTime, String
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.ext.asyncio import AsyncSession, async_sessionmaker, create_async_engine
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column

from app.core.config import settings


class Base(DeclarativeBase):
    """
    Base class for all database models.
    
    Provides common fields and functionality for all models:
    - UUID primary key
    - Created/updated timestamps
    - Type annotations for better IDE support
    """
    
    # Common fields for all models
    id: Mapped[uuid.UUID] = mapped_column(
        UUID(as_uuid=True),
        primary_key=True,
        default=uuid.uuid4,
        index=True
    )
    
    created_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        default=datetime.utcnow,
        nullable=False
    )
    
    updated_at: Mapped[datetime] = mapped_column(
        DateTime(timezone=True),
        default=datetime.utcnow,
        onupdate=datetime.utcnow,
        nullable=False
    )


# Database engine configuration
engine = create_async_engine(
    str(settings.DATABASE_URL),
    echo=settings.DEBUG,
    future=True,
    pool_size=10,
    max_overflow=20,
    pool_pre_ping=True,  # Validate connections before use
    pool_recycle=3600,   # Recycle connections every hour
)

# Session factory
AsyncSessionLocal = async_sessionmaker(
    engine,
    class_=AsyncSession,
    expire_on_commit=False,
    autoflush=True,
    autocommit=False,
)


async def get_session() -> AsyncSession:
    """
    Dependency function to get database session.
    
    This function provides database sessions for FastAPI dependency injection.
    Each request gets its own session that is automatically closed after use.
    
    Yields:
        AsyncSession: Database session
    """
    async with AsyncSessionLocal() as session:
        try:
            yield session
        except Exception:
            await session.rollback()
            raise
        finally:
            await session.close()


async def init_db() -> None:
    """
    Initialize database tables.
    
    Creates all tables defined in models.
    Should be called during application startup.
    """
    async with engine.begin() as conn:
        # 导入所有模型以便注册到 Base.metadata 中
        # Import all models to register them with Base.metadata
        from app.models.user import User  # noqa: F401
        from app.models.provider import Provider  # noqa: F401
        from app.models.model import Model  # noqa: F401
        from app.models.access import UserModelAccess  # noqa: F401
        from app.models.usage import UsageLog  # noqa: F401
        from app.models.file import UploadedFile  # noqa: F401
        
        # 创建所有数据表
        await conn.run_sync(Base.metadata.create_all)


async def close_db() -> None:
    """
    Close database connections.
    
    Should be called during application shutdown.
    """
    await engine.dispose()