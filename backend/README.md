# LYSS AI Platform Backend

A modern, scalable backend for the LYSS AI Platform - a multi-AI provider aggregation and management system.

## Features

- 🔐 **Role-based Access Control (RBAC)** with three user levels
- 🔌 **Pluggable Provider System** supporting OpenAI, Anthropic, DeepSeek and more
- 📊 **Cost Monitoring & Usage Analytics** 
- 🎯 **Scope-based Model Distribution** (Organization vs Personal)
- 🔒 **Encrypted Configuration Storage** for API keys
- 📁 **Document Processing & Q&A** with vector storage
- 🧠 **Smart Memory System** using Mem0
- 🚀 **High Performance** async FastAPI with PostgreSQL

## Tech Stack

- **FastAPI** 0.115.14 - Modern async web framework
- **SQLAlchemy** 2.0.41 - Async ORM with Pydantic integration  
- **PostgreSQL** 16.3 - Primary database
- **Redis** 7.2 - Caching and session storage
- **Qdrant** 1.10.1 - Vector database for embeddings
- **Poetry** 1.8.2 - Dependency management

## Quick Start with Docker

1. **Start all services**:
   ```bash
   cd /root/work/Lyss
   docker-compose up -d
   ```

2. **Check backend logs**:
   ```bash
   docker-compose logs -f backend
   ```

3. **Access API documentation**:
   - OpenAPI Docs: http://localhost:8000/api/v1/docs
   - ReDoc: http://localhost:8000/api/v1/redoc

## Development Setup

### Prerequisites

- Python 3.10+
- PostgreSQL 16+  
- Redis 7+
- Qdrant 1.10+

### Local Installation

1. **Install dependencies**:
   ```bash
   cd backend
   pip install poetry
   poetry install
   ```

2. **Environment setup**:
   ```bash
   cp .env.example .env
   # Edit .env with your configurations
   ```

3. **Start services** (if not using Docker):
   ```bash
   # PostgreSQL
   createdb lyss_db
   
   # Redis
   redis-server
   
   # Qdrant  
   docker run -p 6333:6333 qdrant/qdrant:v1.10.1
   ```

4. **Run development server**:
   ```bash
   python run_dev.py
   ```

## API Overview

### Authentication Endpoints
- `POST /api/v1/auth/jwt/login` - User login
- `POST /api/v1/auth/register` - User registration  
- `POST /api/v1/auth/reset-password` - Password reset

### Provider Management  
- `GET /api/v1/providers/types` - List available provider types
- `POST /api/v1/providers/test` - Test provider configuration
- `POST /api/v1/providers/` - Create new provider
- `GET /api/v1/providers/` - List user's providers
- `POST /api/v1/providers/{id}/sync-models` - Sync models from provider

### User Management
- `GET /api/v1/users/` - List users (admin only)
- `PUT /api/v1/users/{id}/role` - Update user role (super admin only)
- `GET /api/v1/users/{id}/stats` - Get usage statistics

## Architecture Highlights

### Provider Scope System (V2 Core)

```
ORGANIZATION Scope:
├── Created by: Admins only
├── Visibility: Admin dashboard  
├── Distribution: Can be shared with other users
└── Use case: Company-wide AI models

PERSONAL Scope:
├── Created by: Any user
├── Visibility: Owner only
├── Distribution: Cannot be shared
└── Use case: Personal API keys
```

### Plugin Architecture

```python
# Easy to add new providers
@register_provider("custom_ai")
class CustomAIProvider(LLMProvider):
    def get_config_model(cls):
        return CustomAIConfig
    
    async def chat_completion(self, request):
        # Your implementation
        pass
```

### Security Features

- 🔐 **Encrypted Storage**: API keys encrypted with AES-256
- 🛡️ **Role-based Access**: Granular permissions system
- 🔑 **JWT Authentication**: Secure token-based auth
- 📝 **Audit Logs**: Track all provider operations
- 🚫 **Input Validation**: Pydantic schema validation

## Default Admin User

```
Email: admin@lyss.ai
Password: admin123
```

**⚠️ Change this in production!**

## Project Structure

```
backend/
├── app/
│   ├── api/v1/          # API endpoints
│   ├── core/            # Security, config, permissions
│   ├── models/          # SQLAlchemy models  
│   ├── schemas/         # Pydantic schemas
│   ├── services/        # Business logic
│   ├── providers/       # Plugin system
│   │   ├── base/        # Abstract interfaces
│   │   └── impl/        # Provider implementations
│   └── db/              # Database setup
├── alembic/             # Database migrations
└── tests/               # Test suite
```

## Contributing

1. Follow PEP 8 style guidelines
2. Add type hints to all functions
3. Write tests for new features  
4. Update API documentation
5. Use conventional commits

## License

MIT License - see LICENSE file for details