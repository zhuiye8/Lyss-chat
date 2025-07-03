"""
Provider management API endpoints.

This module provides CRUD operations for AI providers with scope-based
access control and configuration management.
"""

import json
from typing import List

from fastapi import APIRouter, Depends, HTTPException, status
from sqlalchemy.ext.asyncio import AsyncSession

from app.core.auth import current_active_user
from app.core.permissions import ProviderScopeChecker, require_admin
from app.core.security import encrypt_sensitive_data
from app.db import get_session
from app.models import Provider, User
from app.providers import (
    create_provider_instance,
    get_provider_config_schema,
    list_provider_types,
    validate_provider_config,
)
from app.schemas.provider import (
    ModelSyncResponse,
    ProviderConfigSchema,
    ProviderCreate,
    ProviderInfo,
    ProviderRead,
    ProviderTestRequest,
    ProviderTestResponse,
    ProviderUpdate,
    ProviderWithOwner,
)
from app.services.provider_service import ProviderService

router = APIRouter()
provider_service = ProviderService()


@router.get("/types", response_model=List[str], tags=["providers"])
async def list_provider_types_endpoint():
    """
    Get list of available provider types.
    
    Returns all registered provider types that can be configured.
    """
    return list_provider_types()


@router.get("/types/{provider_type}/info", response_model=ProviderInfo, tags=["providers"])
async def get_provider_type_info(provider_type: str):
    """
    Get information about a specific provider type.
    
    Returns capabilities, documentation links, and other metadata
    for the specified provider type.
    """
    try:
        from app.providers.registry import get_provider_info
        info = get_provider_info(provider_type)
        return ProviderInfo(provider_type=provider_type, **info)
    except KeyError:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"Provider type '{provider_type}' not found"
        )


@router.get("/types/{provider_type}/config-schema", response_model=ProviderConfigSchema, tags=["providers"])
async def get_provider_config_schema_endpoint(provider_type: str):
    """
    Get configuration schema for a provider type.
    
    Returns JSON schema that describes the required configuration
    fields for the specified provider type.
    """
    schema = get_provider_config_schema(provider_type)
    return ProviderConfigSchema(provider_type=provider_type, config_schema=schema)


@router.post("/test", response_model=ProviderTestResponse, tags=["providers"])
async def test_provider_connection(
    test_request: ProviderTestRequest,
    current_user: User = Depends(current_active_user)
):
    """
    Test provider connection with given configuration.
    
    Validates configuration and tests connectivity without saving.
    Available to all authenticated users.
    """
    try:
        # Validate configuration first
        validated_config = validate_provider_config(
            test_request.provider_type, 
            test_request.config
        )
        
        # Test connection
        success = await provider_service.test_provider_connection(
            test_request.provider_type,
            validated_config
        )
        
        return ProviderTestResponse(
            success=success,
            message="Connection test successful" if success else "Connection test failed"
        )
    
    except Exception as e:
        return ProviderTestResponse(
            success=False,
            message="Connection test failed",
            error_details=str(e)
        )


@router.post("/", response_model=ProviderRead, tags=["providers"])
async def create_provider(
    provider_data: ProviderCreate,
    db: AsyncSession = Depends(get_session),
    current_user: User = Depends(current_active_user)
):
    """
    Create a new AI provider.
    
    Creates a provider with encrypted configuration. Scope determines
    accessibility: ORGANIZATION providers can be distributed by admins,
    PERSONAL providers are for individual use only.
    """
    # Check scope-based permissions
    if provider_data.scope.value == "ORGANIZATION":
        if not ProviderScopeChecker.can_create_organization_provider(current_user):
            raise HTTPException(
                status_code=status.HTTP_403_FORBIDDEN,
                detail="Only admins can create ORGANIZATION providers"
            )
    
    # Validate and encrypt configuration
    validated_config = validate_provider_config(
        provider_data.provider_type,
        provider_data.config
    )
    encrypted_config = encrypt_sensitive_data(json.dumps(validated_config))
    
    # Create provider
    provider = await provider_service.create_provider(
        db=db,
        name=provider_data.name,
        provider_type=provider_data.provider_type,
        scope=provider_data.scope,
        owner_id=current_user.id,
        config_encrypted=encrypted_config,
        description=provider_data.description
    )
    
    return provider


@router.get("/", response_model=List[ProviderRead], tags=["providers"])
async def list_providers(
    db: AsyncSession = Depends(get_session),
    current_user: User = Depends(current_active_user)
):
    """
    Get list of providers accessible to current user.
    
    Returns providers owned by the current user.
    """
    providers = await provider_service.get_user_providers(db, current_user.id)
    return providers


@router.get("/all", response_model=List[ProviderWithOwner], tags=["providers"])
async def list_all_providers(
    skip: int = 0,
    limit: int = 100,
    db: AsyncSession = Depends(get_session),
    current_user: User = Depends(require_admin)
):
    """
    Get list of all providers (admin only).
    
    Returns all providers in the system with owner information.
    Only available to administrators.
    """
    providers = await provider_service.get_all_providers(db, skip=skip, limit=limit)
    return providers


@router.get("/{provider_id}", response_model=ProviderRead, tags=["providers"])
async def get_provider(
    provider_id: str,
    db: AsyncSession = Depends(get_session),
    current_user: User = Depends(current_active_user)
):
    """
    Get specific provider by ID.
    
    Users can only access their own providers unless they are admins.
    """
    provider = await provider_service.get_provider_by_id(db, provider_id)
    
    if not provider:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Provider not found"
        )
    
    # Check access permissions
    if not ProviderScopeChecker.can_view_provider(str(provider.owner_id), current_user):
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to access this provider"
        )
    
    return provider


@router.put("/{provider_id}", response_model=ProviderRead, tags=["providers"])
async def update_provider(
    provider_id: str,
    provider_data: ProviderUpdate,
    db: AsyncSession = Depends(get_session),
    current_user: User = Depends(current_active_user)
):
    """
    Update provider configuration.
    
    Allows updating provider settings and configuration.
    Only the owner or admin can modify providers.
    """
    provider = await provider_service.get_provider_by_id(db, provider_id)
    
    if not provider:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Provider not found"
        )
    
    # Check modify permissions
    if not ProviderScopeChecker.can_modify_provider(str(provider.owner_id), current_user):
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to modify this provider"
        )
    
    # Handle configuration update if provided
    encrypted_config = None
    if provider_data.config is not None:
        validated_config = validate_provider_config(
            provider.provider_type,
            provider_data.config
        )
        encrypted_config = encrypt_sensitive_data(json.dumps(validated_config))
    
    updated_provider = await provider_service.update_provider(
        db=db,
        provider=provider,
        name=provider_data.name,
        description=provider_data.description,
        is_enabled=provider_data.is_enabled,
        config_encrypted=encrypted_config
    )
    
    return updated_provider


@router.delete("/{provider_id}", tags=["providers"])
async def delete_provider(
    provider_id: str,
    db: AsyncSession = Depends(get_session),
    current_user: User = Depends(current_active_user)
):
    """
    Delete provider.
    
    Permanently removes provider and all associated models.
    Only the owner or super admin can delete providers.
    """
    provider = await provider_service.get_provider_by_id(db, provider_id)
    
    if not provider:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Provider not found"
        )
    
    # Check delete permissions
    if not ProviderScopeChecker.can_delete_provider(str(provider.owner_id), current_user):
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to delete this provider"
        )
    
    await provider_service.delete_provider(db, provider)
    
    return {"message": "Provider deleted successfully"}


@router.post("/{provider_id}/sync-models", response_model=ModelSyncResponse, tags=["providers"])
async def sync_provider_models(
    provider_id: str,
    db: AsyncSession = Depends(get_session),
    current_user: User = Depends(current_active_user)
):
    """
    Synchronize models from provider API.
    
    Fetches available models from the provider's API and updates
    the local database. Only the owner or admin can sync models.
    """
    provider = await provider_service.get_provider_by_id(db, provider_id)
    
    if not provider:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="Provider not found"
        )
    
    # Check modify permissions
    if not ProviderScopeChecker.can_modify_provider(str(provider.owner_id), current_user):
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="Not authorized to sync models for this provider"
        )
    
    if not provider.is_enabled:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="Cannot sync models for disabled provider"
        )
    
    sync_result = await provider_service.sync_provider_models(db, provider)
    return sync_result