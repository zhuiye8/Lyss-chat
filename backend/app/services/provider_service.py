"""
Provider service for business logic operations.

This module contains business logic for provider management,
model synchronization, and provider-related operations.
"""

import uuid
from datetime import datetime
from typing import Any, Dict, List, Optional

from sqlalchemy import select
from sqlalchemy.ext.asyncio import AsyncSession
from sqlalchemy.orm import selectinload

from app.models import Model, Provider, ProviderScope, User
from app.providers import create_provider_instance
from app.schemas.provider import ModelSyncResponse


class ProviderService:
    """
    Service class for provider-related business operations.
    
    Encapsulates business logic for provider management, model synchronization,
    and configuration handling.
    """
    
    async def create_provider(
        self,
        db: AsyncSession,
        name: str,
        provider_type: str,
        scope: ProviderScope,
        owner_id: uuid.UUID,
        config_encrypted: str,
        description: Optional[str] = None
    ) -> Provider:
        """
        Create a new provider.
        
        Args:
            db: Database session
            name: Provider display name
            provider_type: Type of provider (e.g., 'openai', 'anthropic')
            scope: Provider scope (ORGANIZATION or PERSONAL)
            owner_id: Owner user ID
            config_encrypted: Encrypted configuration JSON
            description: Optional description
            
        Returns:
            Created provider
        """
        provider = Provider(
            name=name,
            provider_type=provider_type,
            scope=scope,
            owner_id=owner_id,
            config_encrypted=config_encrypted,
            description=description,
            is_enabled=True
        )
        
        db.add(provider)
        await db.commit()
        await db.refresh(provider)
        return provider
    
    async def get_provider_by_id(self, db: AsyncSession, provider_id: str) -> Optional[Provider]:
        """
        Get provider by ID.
        
        Args:
            db: Database session
            provider_id: Provider UUID as string
            
        Returns:
            Provider object or None if not found
        """
        try:
            provider_uuid = uuid.UUID(provider_id)
            result = await db.execute(
                select(Provider)
                .options(selectinload(Provider.models))
                .where(Provider.id == provider_uuid)
            )
            return result.scalar_one_or_none()
        except ValueError:
            return None
    
    async def get_user_providers(
        self, 
        db: AsyncSession, 
        user_id: uuid.UUID
    ) -> List[Provider]:
        """
        Get providers owned by a specific user.
        
        Args:
            db: Database session
            user_id: User UUID
            
        Returns:
            List of providers owned by the user
        """
        result = await db.execute(
            select(Provider)
            .options(selectinload(Provider.models))
            .where(Provider.owner_id == user_id)
            .order_by(Provider.created_at.desc())
        )
        return list(result.scalars().all())
    
    async def get_all_providers(
        self,
        db: AsyncSession,
        skip: int = 0,
        limit: int = 100
    ) -> List[Provider]:
        """
        Get all providers with owner information (admin only).
        
        Args:
            db: Database session
            skip: Number of records to skip
            limit: Maximum number of records
            
        Returns:
            List of providers with owner data
        """
        result = await db.execute(
            select(Provider)
            .options(
                selectinload(Provider.models),
                selectinload(Provider.owner)
            )
            .offset(skip)
            .limit(limit)
            .order_by(Provider.created_at.desc())
        )
        return list(result.scalars().all())
    
    async def update_provider(
        self,
        db: AsyncSession,
        provider: Provider,
        name: Optional[str] = None,
        description: Optional[str] = None,
        is_enabled: Optional[bool] = None,
        config_encrypted: Optional[str] = None
    ) -> Provider:
        """
        Update provider fields.
        
        Args:
            db: Database session
            provider: Provider to update
            name: New name
            description: New description
            is_enabled: Enable/disable flag
            config_encrypted: New encrypted configuration
            
        Returns:
            Updated provider
        """
        if name is not None:
            provider.name = name
        if description is not None:
            provider.description = description
        if is_enabled is not None:
            provider.is_enabled = is_enabled
        if config_encrypted is not None:
            provider.config_encrypted = config_encrypted
        
        provider.updated_at = datetime.utcnow()
        
        await db.commit()
        await db.refresh(provider)
        return provider
    
    async def delete_provider(self, db: AsyncSession, provider: Provider) -> None:
        """
        Delete provider and associated models.
        
        Args:
            db: Database session
            provider: Provider to delete
        """
        await db.delete(provider)
        await db.commit()
    
    async def test_provider_connection(
        self,
        provider_type: str,
        config: Dict[str, Any]
    ) -> bool:
        """
        Test provider connection with given configuration.
        
        Args:
            provider_type: Type of provider
            config: Configuration dictionary
            
        Returns:
            True if connection successful
            
        Raises:
            Exception: If connection fails
        """
        from app.providers.registry import get_provider_class
        
        provider_class = get_provider_class(provider_type)
        provider_instance = provider_class(config)
        
        return await provider_instance.test_connection()
    
    async def sync_provider_models(
        self,
        db: AsyncSession,
        provider: Provider
    ) -> ModelSyncResponse:
        """
        Synchronize models from provider API.
        
        Args:
            db: Database session
            provider: Provider to sync models for
            
        Returns:
            Synchronization results
        """
        errors = []
        models_added = 0
        models_updated = 0
        
        try:
            # Create provider instance
            provider_instance = create_provider_instance(provider)
            
            # Get available models from provider API
            available_models = await provider_instance.get_available_models()
            
            for model_info in available_models:
                try:
                    # Check if model already exists
                    existing_model = await db.execute(
                        select(Model).where(
                            Model.provider_id == provider.id,
                            Model.model_name == model_info.model_name
                        )
                    )
                    existing = existing_model.scalar_one_or_none()
                    
                    if existing:
                        # Update existing model
                        existing.display_name = model_info.display_name
                        existing.description = model_info.description
                        existing.context_length = model_info.context_length
                        existing.max_output_tokens = model_info.max_output_tokens
                        existing.supports_streaming = model_info.supports_streaming
                        existing.supports_function_calling = model_info.supports_function_calling
                        existing.supports_vision = model_info.supports_vision
                        existing.price_per_1k_prompt_tokens = model_info.price_per_1k_prompt_tokens
                        existing.price_per_1k_completion_tokens = model_info.price_per_1k_completion_tokens
                        existing.price_per_1k_tokens = model_info.price_per_1k_tokens
                        existing.updated_at = datetime.utcnow()
                        models_updated += 1
                    else:
                        # Create new model
                        new_model = Model(
                            provider_id=provider.id,
                            model_name=model_info.model_name,
                            display_name=model_info.display_name,
                            description=model_info.description,
                            context_length=model_info.context_length,
                            max_output_tokens=model_info.max_output_tokens,
                            supports_streaming=model_info.supports_streaming,
                            supports_function_calling=model_info.supports_function_calling,
                            supports_vision=model_info.supports_vision,
                            price_per_1k_prompt_tokens=model_info.price_per_1k_prompt_tokens,
                            price_per_1k_completion_tokens=model_info.price_per_1k_completion_tokens,
                            price_per_1k_tokens=model_info.price_per_1k_tokens,
                            is_active=True
                        )
                        db.add(new_model)
                        models_added += 1
                
                except Exception as e:
                    error_msg = f"Failed to sync model {model_info.model_name}: {str(e)}"
                    errors.append(error_msg)
            
            await db.commit()
            
            return ModelSyncResponse(
                provider_id=provider.id,
                models_synced=len(available_models),
                models_added=models_added,
                models_updated=models_updated,
                errors=errors
            )
        
        except Exception as e:
            error_msg = f"Failed to sync models: {str(e)}"
            errors.append(error_msg)
            
            return ModelSyncResponse(
                provider_id=provider.id,
                models_synced=0,
                models_added=0,
                models_updated=0,
                errors=errors
            )
    
    async def get_organization_providers(self, db: AsyncSession) -> List[Provider]:
        """
        Get all ORGANIZATION scoped providers.
        
        Args:
            db: Database session
            
        Returns:
            List of organization providers
        """
        result = await db.execute(
            select(Provider)
            .options(selectinload(Provider.models))
            .where(
                Provider.scope == ProviderScope.ORGANIZATION,
                Provider.is_enabled == True
            )
            .order_by(Provider.created_at.desc())
        )
        return list(result.scalars().all())