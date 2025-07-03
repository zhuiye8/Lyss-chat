"""
Provider factory for creating instances from database models.

This module provides functionality to create provider instances
from encrypted database configurations.
"""

import json
from typing import Any, Dict

from fastapi import HTTPException, status

from app.core.security import decrypt_sensitive_data
from app.models import Provider as ProviderModel
from app.providers.base import LLMProvider
from app.providers.registry import get_provider_class


def create_provider_instance(provider_model: ProviderModel) -> LLMProvider:
    """
    Create a provider instance from a database model.
    
    This function handles the complete process of:
    1. Getting the appropriate provider class from registry
    2. Decrypting the stored configuration
    3. Validating configuration with Pydantic
    4. Instantiating the provider
    
    Args:
        provider_model: Database provider model with encrypted config
        
    Returns:
        Initialized provider instance
        
    Raises:
        HTTPException: If provider creation fails
    """
    try:
        # Get provider class from registry
        provider_class = get_provider_class(provider_model.provider_type)
        
        # Decrypt stored configuration
        decrypted_config_str = decrypt_sensitive_data(provider_model.config_encrypted)
        
        # Parse JSON configuration
        config_dict = json.loads(decrypted_config_str)
        
        # Create provider instance (Pydantic validation happens inside)
        provider_instance = provider_class(config_dict)
        
        return provider_instance
        
    except KeyError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=f"Unknown provider type: {provider_model.provider_type}"
        ) from e
    
    except json.JSONDecodeError as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to parse provider configuration"
        ) from e
    
    except ValueError as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail="Failed to decrypt provider configuration"
        ) from e
    
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
            detail=f"Failed to create provider instance: {str(e)}"
        ) from e


def validate_provider_config(provider_type: str, config: Dict[str, Any]) -> Dict[str, Any]:
    """
    Validate provider configuration without creating instance.
    
    Args:
        provider_type: Type of provider
        config: Configuration dictionary to validate
        
    Returns:
        Validated configuration dictionary
        
    Raises:
        HTTPException: If validation fails
    """
    try:
        provider_class = get_provider_class(provider_type)
        config_model = provider_class.get_config_model()
        
        # Validate configuration
        validated_config = config_model(**config)
        
        # Return as dictionary
        return validated_config.model_dump()
        
    except KeyError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=f"Unknown provider type: {provider_type}"
        ) from e
    
    except Exception as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=f"Configuration validation failed: {str(e)}"
        ) from e


def get_provider_config_schema(provider_type: str) -> Dict[str, Any]:
    """
    Get JSON schema for provider configuration.
    
    Args:
        provider_type: Type of provider
        
    Returns:
        JSON schema for configuration
        
    Raises:
        HTTPException: If provider type not found
    """
    try:
        provider_class = get_provider_class(provider_type)
        config_model = provider_class.get_config_model()
        
        return config_model.model_json_schema()
        
    except KeyError as e:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=f"Unknown provider type: {provider_type}"
        ) from e