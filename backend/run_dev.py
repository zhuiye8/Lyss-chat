"""
Development server runner for LYSS AI Platform.

This script provides a convenient way to run the development server
with proper environment setup and database initialization.
"""

import asyncio
import os
import subprocess
import sys
from pathlib import Path

import uvicorn

# Add the app directory to Python path
app_dir = Path(__file__).parent
sys.path.insert(0, str(app_dir))

def load_env_file():
    """Load environment variables from .env.local file."""
    env_file = app_dir / ".env.local"
    if env_file.exists():
        print(f"Loading environment from {env_file}")
        with open(env_file, 'r') as f:
            for line in f:
                line = line.strip()
                if line and not line.startswith('#') and '=' in line:
                    key, value = line.split('=', 1)
                    os.environ.setdefault(key.strip(), value.strip())
    else:
        print("Warning: .env.local not found, using default values")
        # Set development environment variables
        os.environ.setdefault("ENVIRONMENT", "development")
        os.environ.setdefault("DEBUG", "true")
        os.environ.setdefault("SECRET_KEY", "dev_secret_key_change_in_production")
        os.environ.setdefault("AUTH_SECRET", "dev_auth_secret_change_in_production")
        os.environ.setdefault("ENCRYPTION_KEY", "dev_encryption_key_change_32_chars")
        os.environ.setdefault("DATABASE_URL", "postgresql+asyncpg://lyss:lyss_password_dev@localhost:5432/lyss_db")
        os.environ.setdefault("REDIS_URL", "redis://localhost:6379")
        os.environ.setdefault("QDRANT_URL", "http://localhost:6333")
        os.environ.setdefault("FIRST_SUPERUSER", "admin@lyss.ai")
        os.environ.setdefault("FIRST_SUPERUSER_PASSWORD", "admin123")

def check_and_install_dependencies():
    """Check if dependencies are installed, install if needed."""
    try:
        import fastapi
        import uvicorn
        import sqlalchemy
        print("✓ Dependencies are available")
        return True
    except ImportError:
        print("⚠ Some dependencies are missing. Installing from requirements.txt...")
        try:
            subprocess.check_call([sys.executable, "-m", "pip", "install", "-r", "requirements.txt"])
            print("✓ Dependencies installed successfully")
            return True
        except subprocess.CalledProcessError:
            print("✗ Failed to install dependencies")
            print("Please install manually with: pip install -r requirements.txt")
            return False

def check_database_connection():
    """Check if database services are running."""
    import socket
    
    services = [
        ("PostgreSQL", "localhost", 5432),
        ("Redis", "localhost", 6379),
        ("Qdrant", "localhost", 6333)
    ]
    
    all_running = True
    for name, host, port in services:
        try:
            with socket.create_connection((host, port), timeout=1):
                print(f"✓ {name} is running on {host}:{port}")
        except (socket.timeout, ConnectionRefusedError):
            print(f"✗ {name} is not running on {host}:{port}")
            all_running = False
    
    return all_running


async def init_database():
    """Initialize database and create initial data."""
    try:
        print("Initializing database...")
        from app.initial_data import main
        await main()
        print("Database initialized successfully")
    except Exception as e:
        print(f"Database initialization failed: {e}")
        print("Make sure PostgreSQL, Redis, and Qdrant are running")
        return False
    return True


def main():
    """Main entry point for development server."""
    print("🚀 Starting LYSS AI Platform development server...")
    print("=" * 60)
    
    # Load environment configuration
    load_env_file()
    
    # Check and install dependencies
    if not check_and_install_dependencies():
        print("❌ Cannot proceed without dependencies")
        sys.exit(1)
    
    # Check database services
    if not check_database_connection():
        print("\n⚠️  Some database services are not running!")
        print("Please start them with: docker-compose -f docker-compose.db.yml up -d")
        print("Or use the regular docker-compose up -d if you prefer")
        
        response = input("\nContinue anyway? (y/N): ").lower().strip()
        if response != 'y':
            sys.exit(1)
    
    # Check if we can import the app
    try:
        from app.main import app
        print("✓ Application imported successfully")
    except ImportError as e:
        print(f"✗ Failed to import application: {e}")
        print("Make sure you're in the backend directory and dependencies are installed")
        sys.exit(1)
    
    # Initialize database
    try:
        print("\n📊 Initializing database...")
        if not asyncio.run(init_database()):
            print("⚠️  Database initialization failed. Continuing anyway...")
    except Exception as e:
        print(f"⚠️  Database initialization error: {e}")
        print("Continuing without database initialization...")
    
    # Start the server
    print("\n" + "=" * 60)
    print("🎉 Starting FastAPI development server")
    print("📍 Server: http://localhost:8000")
    print("📚 API Docs: http://localhost:8000/api/v1/docs")
    print("🔧 ReDoc: http://localhost:8000/api/v1/redoc")
    print("👤 Admin: admin@lyss.ai / admin123")
    print("=" * 60)
    
    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=8000,
        reload=True,
        log_level="info"
    )


if __name__ == "__main__":
    main()