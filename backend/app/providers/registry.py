"""
Provider plugin registry and discovery system.

This module manages the registration and discovery of provider plugins,
enabling dynamic loading and instantiation of different AI providers.
"""

from typing import Dict, List, Type

from app.providers.base import LLMProvider

# Global registry of provider classes
_provider_registry: Dict[str, Type[LLMProvider]] = {}


def register_provider(provider_type: str):
    """
    Decorator to register a provider class in the global registry.
    
    Args:
        provider_type: Unique identifier for the provider (e.g., 'openai', 'anthropic')
        
    Returns:
        Decorator function
        
    Example:
        @register_provider("openai")
        class OpenAIProvider(LLMProvider):
            ...
    """
    def decorator(cls: Type[LLMProvider]) -> Type[LLMProvider]:
        if provider_type in _provider_registry:
            raise ValueError(f"Provider type '{provider_type}' is already registered")
        
        if not issubclass(cls, LLMProvider):
            raise TypeError(f"Provider class must inherit from LLMProvider")
        
        _provider_registry[provider_type] = cls
        return cls
    
    return decorator


def get_provider_class(provider_type: str) -> Type[LLMProvider]:
    """
    Get provider class by type.
    
    Args:
        provider_type: Provider type identifier
        
    Returns:
        Provider class
        
    Raises:
        KeyError: If provider type is not registered
    """
    if provider_type not in _provider_registry:
        raise KeyError(f"Unknown provider type: {provider_type}")
    
    return _provider_registry[provider_type]


def list_provider_types() -> List[str]:
    """
    Get list of all registered provider types.
    
    Returns:
        List of provider type identifiers
    """
    return list(_provider_registry.keys())


def get_provider_info(provider_type: str) -> dict:
    """
    Get provider information and capabilities.
    
    Args:
        provider_type: Provider type identifier
        
    Returns:
        Provider information dictionary
        
    Raises:
        KeyError: If provider type is not registered
    """
    provider_class = get_provider_class(provider_type)
    return provider_class.get_provider_info()


def is_provider_registered(provider_type: str) -> bool:
    """
    Check if a provider type is registered.
    
    Args:
        provider_type: Provider type to check
        
    Returns:
        True if provider is registered
    """
    return provider_type in _provider_registry


# Auto-discovery: Import all provider implementations to trigger registration
def _discover_providers():
    """
    Discover and import all provider implementations.
    
    This function imports all provider modules to ensure they are
    registered in the global registry.
    """
    try:
        # Import all provider implementations
        from app.providers.impl import openai_provider  # noqa: F401
        from app.providers.impl import anthropic_provider  # noqa: F401
        from app.providers.impl import deepseek_provider  # noqa: F401
    except ImportError as e:
        # Some providers might not be available in development
        print(f"Warning: Could not import provider implementation: {e}")


# Auto-discover providers on module import
_discover_providers()