"""
Basic tests for LYSS AI Platform backend.

Tests core functionality including API endpoints, models, and services.
"""

import pytest
from fastapi.testclient import TestClient
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

from app.main import app
from app.db.base import Base, get_session

# Test database URL (in-memory SQLite for testing)
SQLALCHEMY_DATABASE_URL = "sqlite:///./test.db"

# Create test database engine
engine = create_engine(
    SQLALCHEMY_DATABASE_URL, 
    connect_args={"check_same_thread": False}
)

# Create test session
TestingSessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)


def override_get_session():
    """Override database session for testing."""
    try:
        db = TestingSessionLocal()
        yield db
    finally:
        db.close()


# Override the dependency
app.dependency_overrides[get_session] = override_get_session


@pytest.fixture(scope="session")
def client():
    """Create test client."""
    # Create tables
    Base.metadata.create_all(bind=engine)
    
    with TestClient(app) as test_client:
        yield test_client
    
    # Clean up
    Base.metadata.drop_all(bind=engine)


def test_root_endpoint(client):
    """Test root endpoint returns basic information."""
    response = client.get("/")
    assert response.status_code == 200
    data = response.json()
    assert "message" in data
    assert "version" in data
    assert data["message"] == "LYSS AI Platform Backend"


def test_health_endpoint(client):
    """Test health check endpoint."""
    response = client.get("/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "healthy"


def test_api_health_endpoint(client):
    """Test API health check endpoint."""
    response = client.get("/api/v1/health")
    assert response.status_code == 200
    data = response.json()
    assert data["status"] == "healthy"
    assert data["version"] == "v1"


def test_provider_types_endpoint(client):
    """Test provider types endpoint."""
    response = client.get("/api/v1/providers/types")
    assert response.status_code == 200
    data = response.json()
    assert isinstance(data, list)
    # Should have our registered providers
    assert "openai" in data
    assert "anthropic" in data
    assert "deepseek" in data


def test_provider_info_endpoint(client):
    """Test provider info endpoint."""
    response = client.get("/api/v1/providers/types/openai/info")
    assert response.status_code == 200
    data = response.json()
    assert data["provider_type"] == "openai"
    assert data["name"] == "OpenAI"
    assert "supports_streaming" in data


def test_provider_config_schema_endpoint(client):
    """Test provider config schema endpoint."""
    response = client.get("/api/v1/providers/types/openai/config-schema")
    assert response.status_code == 200
    data = response.json()
    assert data["provider_type"] == "openai"
    assert "schema" in data
    assert "properties" in data["schema"]


def test_invalid_provider_type(client):
    """Test invalid provider type returns 404."""
    response = client.get("/api/v1/providers/types/invalid/info")
    assert response.status_code == 404


def test_openapi_docs(client):
    """Test that OpenAPI docs are available."""
    response = client.get("/api/v1/openapi.json")
    assert response.status_code == 200
    data = response.json()
    assert "openapi" in data
    assert "info" in data
    assert data["info"]["title"] == "LYSS AI Platform"