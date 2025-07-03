"""
Main API router for v1 endpoints.

This module aggregates all API routes and provides the main
API router for the v1 version of the LYSS AI Platform API.
"""

from fastapi import APIRouter

from app.api.v1.auth import auth_router
from app.api.v1.providers import router as providers_router
from app.api.v1.users import router as users_router

# Create main API router
api_router = APIRouter()

# Include all sub-routers
api_router.include_router(auth_router, prefix="/auth", tags=["auth"])
api_router.include_router(users_router, tags=["users"])
api_router.include_router(providers_router, prefix="/providers", tags=["providers"])

# Health check for API
@api_router.get("/health")
async def api_health():
    """API health check endpoint."""
    return {
        "status": "healthy",
        "version": "v1",
        "message": "LYSS AI Platform API is running"
    }