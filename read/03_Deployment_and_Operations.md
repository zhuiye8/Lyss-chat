# LYSS AI 平台 - 部署与运维

**版本**: 1.0
**最后更新**: 2025年7月2日

---

## 1. 概述

本文档为 LYSS AI 平台的部署、配置和日常运维提供了指导。我们的目标是提供一个标准化的、基于容器的部署方案，以简化开发、测试和生产环境的管理。

## 2. 部署架构

我们采用 **Docker** 和 **Docker Compose** 作为核心的容器化和编排工具。这种方法将整个应用（前端、后端、数据库）及其依赖项打包成一组可移植的、自包含的服务。

### 2.1. `docker-compose.yml`

项目根目录下的 `docker-compose.yml` 文件定义了整个应用的服务栈。

```yaml
version: '3.8'

services:
  # 1. 后端 FastAPI 服务
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    ports:
      - "8000:8000"
    volumes:
      - ./backend:/app # 开发模式下挂载代码，实现热重载
    env_file:
      - ./backend/.env
    depends_on:
      - db
      - redis
      - qdrant
    command: uvicorn app.main:app --host 0.0.0.0 --port 8000 --reload

  # 2. 前端 React 服务
  frontend:
    build:
      context: ./frontend
      dockerfile: Dockerfile
    ports:
      - "8001:80" # Nginx 默认监听80端口
    depends_on:
      - backend

  # 3. PostgreSQL 数据库
  db:
    image: postgres:16.3-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data/
    env_file:
      - ./backend/.env # 复用后端环境变量配置
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${POSTGRES_USER} -d ${POSTGRES_DB}"]
      interval: 5s
      timeout: 5s
      retries: 5

  # 4. Redis 缓存
  redis:
    image: redis:7.2-alpine
    ports:
      - "6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5

  # 5. Qdrant 向量数据库
  qdrant:
    image: qdrant/qdrant:v1.10.1
    ports:
      - "6333:6333"
      - "6334:6334"
    volumes:
      - qdrant_data:/qdrant/storage

volumes:
  postgres_data:
  qdrant_data:
```

### 2.2. 后端 `Dockerfile`

```dockerfile
# backend/Dockerfile
FROM python:3.10-slim

WORKDIR /app

# 使用 Poetry 进行依赖管理
RUN pip install poetry==1.8.2
COPY poetry.lock pyproject.toml ./
RUN poetry config virtualenvs.create false && poetry install --no-dev --no-root

# 拷贝应用代码
COPY ./app ./app

# 暴露端口
EXPOSE 8000
```

### 2.3. 前端 `Dockerfile` (生产环境)

此 Dockerfile 用于构建生��环境下的静态文件并使用 Nginx 提供服务。

```dockerfile
# frontend/Dockerfile

# --- Build Stage ---
FROM node:20-alpine AS builder
WORKDIR /app
COPY package.json pnpm-lock.yaml ./
RUN npm install -g pnpm@9.5.0 && pnpm install
COPY . .
RUN pnpm build

# --- Production Stage ---
FROM nginx:1.25-alpine
COPY --from=builder /app/dist /usr/share/nginx/html
# 可选：拷贝自定义的 nginx 配置
# COPY nginx.conf /etc/nginx/conf.d/default.conf
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
```

## 3. 环境变量配置

所有敏感信息和环境特定的配置都必须通过环境变量进行管理，严禁硬编码在代码中。

### 3.1. 后端 (`backend/.env`)

需要创建一个 `.env` 文件，可以从 `.env.example` 复制。

```dotenv
# backend/.env.example

# --- Application ---
ENVIRONMENT=development # development or production
SECRET_KEY=your_strong_secret_key # 用于fastapi-users
AUTH_SECRET=another_strong_secret_key # 同样用于fastapi-users

# --- Database (PostgreSQL) ---
POSTGRES_SERVER=db
POSTGRES_PORT=5432
POSTGRES_USER=lyss
POSTGRES_PASSWORD=your_db_password
POSTGRES_DB=lyss_db
DATABASE_URL=postgresql+asyncpg://${POSTGRES_USER}:${POSTGRES_PASSWORD}@${POSTGRES_SERVER}:${POSTGRES_PORT}/${POSTGRES_DB}

# --- Cache (Redis) ---
REDIS_URL=redis://redis:6379

# --- Vector Store (Qdrant) ---
QDRANT_URL=http://qdrant:6333

# --- Mem0 Configuration ---
# Mem0 会自动使用 QDRANT_URL，但也可以显式配置
MEM0_VECTOR_STORE_URI=${QDRANT_URL}
```

## 4. 启动与运维

### 4.1. 首次启动

1.  **克隆代码**: `git clone ...`
2.  **配置环境**:
    *   `cp backend/.env.example backend/.env`
    *   编辑 `backend/.env` 文件，填入所有必需的值，特别是 `SECRET_KEY`, `AUTH_SECRET` 和 `POSTGRES_PASSWORD`。
3.  **构建并启动服务**:
    ```bash
    docker-compose up --build -d
    ```
4.  **初始化数据库 (首次)**:
    *   进入后端容器: `docker-compose exec backend bash`
    *   在容器内运行初始化脚本: `python -m app.initial_data` (此脚本需手动创建，用于创建超级管理员等初始数据)。

### 4.2. 日常运维

*   **查看日志**: `docker-compose logs -f <service_name>` (例如: `docker-compose logs -f backend`)
*   **停止服务**: `docker-compose down`
*   **更新镜像**: `docker-compose pull`
*   **重建服务**: `docker-compose up --build -d <service_name>`

## 5. 备份与恢复

*   **PostgreSQL**: 数据持久化在 `postgres_data` Docker 卷中。应定期使用 `pg_dump` 对数据库进行逻辑备份。
*   **Qdrant**: 数据持久化在 `qdrant_data` 卷中。Qdrant 提供了快照 (snapshot) API，可以用于创建和恢复数据备份。

通过标准化的容器部署方案，我们可以确保开发、测试和生产环境的一致性，极大地降低了部署和运维的复杂性。
