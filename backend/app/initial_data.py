"""
Initial data setup for LYSS AI Platform.

This script creates the initial database structure and
sets up the first super admin user.
"""

import asyncio
import logging

from app.core.config import settings
from app.core.security import get_password_hash
from app.db import AsyncSessionLocal, init_db
from app.models import User, UserRole
from app.services.user_service import UserService

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)


async def create_initial_data() -> None:
    """Create initial data including super admin user."""
    
    # Initialize database tables
    logger.info("Creating database tables...")
    await init_db()
    logger.info("Database tables created successfully")
    
    # Create session
    async with AsyncSessionLocal() as session:
        user_service = UserService()
        
        # Check if super admin already exists
        from sqlalchemy import select
        result = await session.execute(
            select(User).where(User.email == settings.FIRST_SUPERUSER)
        )
        existing_admin = result.scalar_one_or_none()
        
        if existing_admin:
            logger.info(f"Super admin user {settings.FIRST_SUPERUSER} already exists")
            return
        
        # Create super admin user
        logger.info(f"Creating super admin user: {settings.FIRST_SUPERUSER}")
        
        password_hash = get_password_hash(settings.FIRST_SUPERUSER_PASSWORD)
        
        super_admin = await user_service.create_user(
            db=session,
            email=settings.FIRST_SUPERUSER,
            password_hash=password_hash,
            first_name="Super",
            last_name="Admin", 
            role=UserRole.SUPER_ADMIN
        )
        
        logger.info(f"Super admin user created with ID: {super_admin.id}")
        logger.info("Initial data setup completed successfully")


async def main() -> None:
    """Main entry point for initial data setup."""
    try:
        await create_initial_data()
    except Exception as e:
        logger.error(f"Failed to create initial data: {e}")
        raise


if __name__ == "__main__":
    asyncio.run(main())