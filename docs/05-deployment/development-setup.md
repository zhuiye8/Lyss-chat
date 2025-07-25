# Lyss AI Platform 开发环境搭建指南

**版本**: 2.0  
**更新时间**: 2025-01-25  
**技术栈**: Kratos + Docker + 本地混合部署  
**状态**: 已确认

---

## 概述

本指南详细说明如何搭建 Lyss AI Platform 的本地开发环境。基于 **本地服务 + Docker基础设施** 的混合模式，确保开发调试的便利性和环境的一致性。

### 环境架构

```
开发环境架构：
┌─────────────────────────────────────────┐
│  本地Go微服务 (开发调试模式)            │
│  ├── user-service:8001                 │
│  ├── auth-service:8002                 │  
│  ├── group-service:8003                │
│  ├── credential-service:8004           │
│  ├── gateway-service:8000              │
│  └── billing-service:8005              │
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│  Docker基础设施容器                      │
│  ├── PostgreSQL:5432                   │
│  ├── Redis:6379                        │
│  ├── Consul:8500 (服务发现)            │
│  ├── Qdrant:6333 (向量数据库)          │
│  └── Prometheus:9090 (监控,可选)       │
└─────────────────────────────────────────┘
```

---

## 系统要求

### 硬件要求

| 组件 | 最低配置 | 推荐配置 |
|------|----------|----------|
| **CPU** | 4核 | 8核+ |
| **内存** | 8GB | 16GB+ |
| **存储** | 20GB可用空间 | 50GB+ SSD |
| **网络** | 10Mbps | 100Mbps+ |

### 软件要求

| 软件 | 版本要求 | 用途 |
|------|----------|------|
| **Go** | 1.21+ | 后端服务开发 |
| **Node.js** | 18+ | 前端开发和工具链 |
| **Docker** | 20.10+ | 基础设施容器 |
| **Docker Compose** | 2.0+ | 多容器编排 |
| **Git** | 2.30+ | 版本控制 |
| **Make** | 4.0+ | 构建工具 |

### 操作系统支持

- ✅ **macOS** 12+ (Apple Silicon 和 Intel)
- ✅ **Linux** Ubuntu 20.04+, CentOS 8+, Debian 11+
- ✅ **Windows** 10/11 (WSL2 推荐)

---

## 环境安装

### 1. 基础软件安装

#### macOS 安装

```bash
# 安装 Homebrew (如果尚未安装)
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# 安装必需软件
brew install go node docker docker-compose git make

# 启动 Docker Desktop
open -a Docker

# 验证安装
go version          # 应显示 go1.21+ 
node --version      # 应显示 v18+
docker --version    # 应显示 20.10+
docker-compose --version  # 应显示 2.0+
```

#### Ubuntu/Debian 安装

```bash
# 更新包管理器
sudo apt update && sudo apt upgrade -y

# 安装基础工具
sudo apt install -y curl wget git make build-essential

# 安装 Go
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装 Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt install -y nodejs

# 安装 Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER

# 安装 Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 重启终端或重新登录以生效 Docker 组权限
newgrp docker

# 验证安装
go version
node --version
docker --version
docker-compose --version
```

#### Windows (WSL2) 安装

```powershell
# 在 PowerShell (管理员) 中启用 WSL2
wsl --install -d Ubuntu-20.04

# 重启计算机后，在 WSL2 Ubuntu 中执行：
```

```bash
# 在 WSL2 Ubuntu 中执行 Ubuntu 安装步骤
# 然后安装 Docker Desktop for Windows
# 确保在 Docker Desktop 设置中启用 WSL2 集成
```

### 2. IDE 和工具推荐

#### VS Code 插件
```bash
# 安装 VS Code (如果尚未安装)
# macOS
brew install --cask visual-studio-code

# Ubuntu/Debian
wget -qO- https://packages.microsoft.com/keys/microsoft.asc | gpg --dearmor > packages.microsoft.gpg
sudo install -o root -g root -m 644 packages.microsoft.gpg /etc/apt/trusted.gpg.d/
echo "deb [arch=amd64,arm64,armhf signed-by=/etc/apt/trusted.gpg.d/packages.microsoft.gpg] https://packages.microsoft.com/repos/code stable main" | sudo tee /etc/apt/sources.list.d/vscode.list
sudo apt update && sudo apt install code
```

**推荐 VS Code 插件**：
```json
{
  "recommendations": [
    "golang.go",                    // Go 语言支持
    "ms-vscode.vscode-typescript-next", // TypeScript 支持
    "Vue.volar",                    // Vue 3 支持
    "ms-azuretools.vscode-docker",  // Docker 支持
    "hashicorp.terraform",          // Terraform 支持
    "ms-kubernetes-tools.vscode-kubernetes-tools", // K8s 支持
    "ms-vscode.rest-client",        // HTTP 客户端
    "bradlc.vscode-tailwindcss",    // Tailwind CSS
    "esbenp.prettier-vscode",       // 代码格式化
    "ms-vscode.live-server",        // 本地服务器
    "GitLens.gitlens"               // Git 增强
  ]
}
```

---

## 项目设置

### 1. 克隆项目

```bash
# 克隆主项目仓库
git clone https://github.com/your-org/lyss-ai-platform.git
cd lyss-ai-platform

# 项目目录结构
lyss-ai-platform/
├── services/                 # 微服务源码
│   ├── user-service/
│   ├── auth-service/
│   ├── group-service/
│   ├── credential-service/
│   ├── gateway-service/
│   └── billing-service/
├── frontend/                 # 前端项目
├── infrastructure/           # 基础设施配置
│   ├── docker-compose.yml
│   ├── docker-compose.override.yml
│   └── configs/
├── docs/                     # 项目文档
├── scripts/                  # 开发脚本
├── Makefile                  # 构建配置
├── .env.example             # 环境变量模板
└── README.md
```

### 2. 环境变量配置

```bash
# 复制环境变量模板
cp .env.example .env

# 编辑环境变量
vim .env
```

**.env 文件配置**：
```bash
# === 基础配置 ===
ENV=development
DEBUG=true
LOG_LEVEL=debug

# === 数据库配置 ===
DB_HOST=localhost
DB_PORT=5432
DB_USER=lyss
DB_PASSWORD=lyss123
DB_NAME=lyss_platform

# === Redis 配置 ===
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# === 服务发现配置 ===
CONSUL_HOST=localhost
CONSUL_PORT=8500

# === 向量数据库配置 ===
QDRANT_HOST=localhost
QDRANT_PORT=6333

# === AI 服务配置 ===
OPENAI_API_KEY=sk-your-openai-key
OPENAI_BASE_URL=https://api.openai.com/v1

AZURE_OPENAI_ENDPOINT=https://your-resource.openai.azure.com/
AZURE_OPENAI_API_KEY=your-azure-key
AZURE_OPENAI_API_VERSION=2023-12-01-preview

ALIYUN_DASH_SCOPE_API_KEY=sk-your-aliyun-key

# === 记忆服务配置 ===
MEM0_API_KEY=your-mem0-key
MEM0_BASE_URL=http://localhost:8080

# === JWT 配置 ===
JWT_SECRET=your-super-secret-jwt-key-change-in-production
JWT_EXPIRES_IN=24h

# === 监控配置 (可选) ===
PROMETHEUS_ENABLED=true
JAEGER_ENDPOINT=http://localhost:14268/api/traces

# === 开发配置 ===
HOT_RELOAD=true
API_CORS_ORIGINS=http://localhost:3000,http://localhost:5173
```

### 3. Go 模块初始化

```bash
# 为每个服务初始化 Go 模块
cd services/user-service
go mod init github.com/your-org/lyss-ai-platform/services/user-service
go mod tidy

cd ../auth-service
go mod init github.com/your-org/lyss-ai-platform/services/auth-service  
go mod tidy

cd ../gateway-service
go mod init github.com/your-org/lyss-ai-platform/services/gateway-service
go mod tidy

# ... 对其他服务重复相同操作

# 返回项目根目录
cd ../../
```

### 4. 前端项目初始化

```bash
# 进入前端目录
cd frontend

# 安装依赖
npm install

# 验证前端项目
npm run dev

# 打开新终端验证前端是否正常启动
curl http://localhost:5173
```

---

## 基础设施启动

### 1. Docker Compose 配置

**infrastructure/docker-compose.yml**：
```yaml
version: '3.8'

services:
  # PostgreSQL 主数据库
  postgres:
    image: postgres:15-alpine
    container_name: lyss-postgres
    ports:
      - "5432:5432"
    environment:
      POSTGRES_DB: lyss_platform
      POSTGRES_USER: lyss
      POSTGRES_PASSWORD: lyss123
      POSTGRES_INITDB_ARGS: "--encoding=UTF-8 --locale=C"
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./configs/postgres/init.sql:/docker-entrypoint-initdb.d/init.sql
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U lyss"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # Redis 缓存
  redis:
    image: redis:7-alpine
    container_name: lyss-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
      - ./configs/redis/redis.conf:/usr/local/etc/redis/redis.conf
    command: redis-server /usr/local/etc/redis/redis.conf
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # Consul 服务发现
  consul:
    image: consul:1.17
    container_name: lyss-consul
    ports:
      - "8500:8500"
      - "8600:8600/udp"
    volumes:
      - consul_data:/consul/data
      - ./configs/consul/consul.hcl:/consul/config/consul.hcl
    command: >
      consul agent -server -bootstrap-expect=1 -ui -bind=0.0.0.0 
      -client=0.0.0.0 -config-file=/consul/config/consul.hcl
    healthcheck:
      test: ["CMD", "consul", "members"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # Qdrant 向量数据库
  qdrant:
    image: qdrant/qdrant:v1.7.4
    container_name: lyss-qdrant
    ports:
      - "6333:6333"
      - "6334:6334"
    volumes:
      - qdrant_data:/qdrant/storage
      - ./configs/qdrant/config.yaml:/qdrant/config/production.yaml
    environment:
      QDRANT__SERVICE__HTTP_PORT: 6333
      QDRANT__SERVICE__GRPC_PORT: 6334
    healthcheck:
      test: ["CMD", "curl", "-f", "http://localhost:6333/health"]
      interval: 10s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  # Prometheus 监控 (可选)
  prometheus:
    image: prom/prometheus:v2.48.1
    container_name: lyss-prometheus
    ports:
      - "9090:9090"
    volumes:
      - prometheus_data:/prometheus
      - ./configs/prometheus/prometheus.yml:/etc/prometheus/prometheus.yml
    command:
      - '--config.file=/etc/prometheus/prometheus.yml'
      - '--storage.tsdb.path=/prometheus'
      - '--web.console.libraries=/etc/prometheus/console_libraries'
      - '--web.console.templates=/etc/prometheus/consoles'
      - '--storage.tsdb.retention.time=200h'
      - '--web.enable-lifecycle'
    restart: unless-stopped
    profiles: ["monitoring"]

  # Grafana 可视化 (可选)
  grafana:
    image: grafana/grafana:10.2.3
    container_name: lyss-grafana
    ports:
      - "3001:3000"
    volumes:
      - grafana_data:/var/lib/grafana
      - ./configs/grafana/provisioning:/etc/grafana/provisioning
    environment:
      GF_SECURITY_ADMIN_USER: admin
      GF_SECURITY_ADMIN_PASSWORD: admin123
      GF_USERS_ALLOW_SIGN_UP: false
    restart: unless-stopped
    profiles: ["monitoring"]

volumes:
  postgres_data:
  redis_data:
  consul_data:
  qdrant_data:
  prometheus_data:
  grafana_data:

networks:
  default:
    name: lyss-network
    driver: bridge
```

### 2. 配置文件准备

**infrastructure/configs/postgres/init.sql**：
```sql
-- 创建数据库
CREATE DATABASE lyss_user;
CREATE DATABASE lyss_auth;
CREATE DATABASE lyss_group;
CREATE DATABASE lyss_credential;
CREATE DATABASE lyss_gateway;
CREATE DATABASE lyss_billing;

-- 创建用户并授权
CREATE USER lyss_user_svc WITH ENCRYPTED PASSWORD 'user123';
CREATE USER lyss_auth_svc WITH ENCRYPTED PASSWORD 'auth123';
CREATE USER lyss_group_svc WITH ENCRYPTED PASSWORD 'group123';
CREATE USER lyss_credential_svc WITH ENCRYPTED PASSWORD 'credential123';
CREATE USER lyss_gateway_svc WITH ENCRYPTED PASSWORD 'gateway123';
CREATE USER lyss_billing_svc WITH ENCRYPTED PASSWORD 'billing123';

GRANT ALL PRIVILEGES ON DATABASE lyss_user TO lyss_user_svc;
GRANT ALL PRIVILEGES ON DATABASE lyss_auth TO lyss_auth_svc;
GRANT ALL PRIVILEGES ON DATABASE lyss_group TO lyss_group_svc;
GRANT ALL PRIVILEGES ON DATABASE lyss_credential TO lyss_credential_svc;
GRANT ALL PRIVILEGES ON DATABASE lyss_gateway TO lyss_gateway_svc; 
GRANT ALL PRIVILEGES ON DATABASE lyss_billing TO lyss_billing_svc;

-- 创建扩展
\c lyss_platform;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";
CREATE EXTENSION IF NOT EXISTS "btree_gin";
```

**infrastructure/configs/redis/redis.conf**：
```conf
# Redis 开发环境配置
bind 0.0.0.0
port 6379
timeout 0
tcp-keepalive 300

# 数据持久化
save 900 1
save 300 10
save 60 10000

# 内存管理
maxmemory 512mb
maxmemory-policy allkeys-lru

# 日志配置
loglevel notice
logfile ""

# AOF 持久化 (开发环境关闭以提升性能)
appendonly no

# 客户端连接
maxclients 10000
```

**infrastructure/configs/consul/consul.hcl**：
```hcl
datacenter = "dc1"
data_dir = "/consul/data"
log_level = "INFO"
server = true
bootstrap_expect = 1
ui_config {
  enabled = true
}
bind_addr = "0.0.0.0"
client_addr = "0.0.0.0"
retry_join = ["127.0.0.1"]
connect {
  enabled = true
}
```

### 3. 启动基础设施

```bash
# 进入基础设施目录
cd infrastructure

# 启动核心服务
docker-compose up -d postgres redis consul qdrant

# 检查服务状态
docker-compose ps

# 查看服务日志
docker-compose logs -f

# 启动监控服务 (可选)
docker-compose --profile monitoring up -d prometheus grafana

# 等待所有服务健康检查通过
./scripts/wait-for-services.sh
```

**scripts/wait-for-services.sh**：
```bash
#!/bin/bash

# 等待服务启动脚本
set -e

echo "🚀 等待基础设施服务启动..."

# 等待 PostgreSQL
echo "⏳ 等待 PostgreSQL..."
until docker-compose exec -T postgres pg_isready -U lyss; do
  echo "PostgreSQL 未就绪，等待中..."
  sleep 2
done
echo "✅ PostgreSQL 已就绪"

# 等待 Redis
echo "⏳ 等待 Redis..."
until docker-compose exec -T redis redis-cli ping | grep -q PONG; do
  echo "Redis 未就绪，等待中..."
  sleep 2
done
echo "✅ Redis 已就绪"

# 等待 Consul
echo "⏳ 等待 Consul..."
until curl -s http://localhost:8500/v1/status/leader | grep -q .; do
  echo "Consul 未就绪，等待中..."
  sleep 2
done
echo "✅ Consul 已就绪"

# 等待 Qdrant
echo "⏳ 等待 Qdrant..."
until curl -s http://localhost:6333/health | grep -q "ok"; do
  echo "Qdrant 未就绪，等待中..."
  sleep 2
done
echo "✅ Qdrant 已就绪"

echo "🎉 所有基础设施服务已就绪！"

# 显示服务状态
echo "📋 服务状态："
echo "- PostgreSQL: http://localhost:5432"
echo "- Redis: http://localhost:6379"  
echo "- Consul: http://localhost:8500"
echo "- Qdrant: http://localhost:6333"

if docker-compose ps prometheus | grep -q "Up"; then
  echo "- Prometheus: http://localhost:9090"
fi

if docker-compose ps grafana | grep -q "Up"; then
  echo "- Grafana: http://localhost:3001 (admin/admin123)"
fi
```

---

## 微服务开发

### 1. Makefile 配置

**项目根目录 Makefile**：
```makefile
# Lyss AI Platform 开发 Makefile
.PHONY: help install start stop test build clean

# 默认目标
help: ## 显示帮助信息
	@echo "Lyss AI Platform 开发环境"
	@echo "可用命令："
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# === 环境管理 ===
install: ## 安装依赖
	@echo "📦 安装 Go 依赖..."
	@for service in services/*/; do \
		if [ -f "$$service/go.mod" ]; then \
			echo "安装 $$service 依赖..."; \
			cd "$$service" && go mod download && go mod tidy && cd ../..; \
		fi; \
	done
	@echo "📦 安装前端依赖..."
	@cd frontend && npm install

setup: ## 初始化开发环境
	@echo "🔧 初始化开发环境..."
	@cp .env.example .env
	@chmod +x scripts/*.sh
	@echo "✅ 环境初始化完成，请编辑 .env 文件"

# === 基础设施管理 ===
infra-start: ## 启动基础设施
	@echo "🚀 启动基础设施容器..."
	@cd infrastructure && docker-compose up -d postgres redis consul qdrant
	@./scripts/wait-for-services.sh

infra-stop: ## 停止基础设施
	@echo "🛑 停止基础设施容器..."
	@cd infrastructure && docker-compose down

infra-logs: ## 查看基础设施日志
	@cd infrastructure && docker-compose logs -f

infra-clean: ## 清理基础设施数据
	@echo "🧹 清理基础设施数据..."
	@cd infrastructure && docker-compose down -v
	@docker system prune -f

# === 服务管理 ===
services-start: ## 启动所有微服务
	@echo "🚀 启动微服务..."
	@make -j6 start-user-service start-auth-service start-group-service start-credential-service start-gateway-service start-billing-service

start-user-service: ## 启动用户服务
	@echo "🚀 启动用户服务..."
	@cd services/user-service && go run cmd/main.go &

start-auth-service: ## 启动认证服务
	@echo "🚀 启动认证服务..."
	@cd services/auth-service && go run cmd/main.go &

start-group-service: ## 启动群组服务
	@echo "🚀 启动群组服务..."
	@cd services/group-service && go run cmd/main.go &

start-credential-service: ## 启动凭证服务
	@echo "🚀 启动凭证服务..."
	@cd services/credential-service && go run cmd/main.go &

start-gateway-service: ## 启动网关服务
	@echo "🚀 启动网关服务..."
	@cd services/gateway-service && go run cmd/main.go &

start-billing-service: ## 启动计费服务
	@echo "🚀 启动计费服务..."
	@cd services/billing-service && go run cmd/main.go &

services-stop: ## 停止所有微服务
	@echo "🛑 停止微服务..."
	@pkill -f "go run cmd/main.go" || true

# === 前端管理 ===
frontend-start: ## 启动前端开发服务器
	@echo "🚀 启动前端..."
	@cd frontend && npm run dev

frontend-build: ## 构建前端
	@echo "🔨 构建前端..."
	@cd frontend && npm run build

# === 完整开发环境 ===
start: infra-start ## 启动完整开发环境
	@echo "⏳ 等待基础设施就绪..."
	@sleep 10
	@make services-start
	@echo "🎉 开发环境已启动！"
	@echo "📋 服务地址："
	@echo "  - 网关服务: http://localhost:8000"
	@echo "  - 前端应用: http://localhost:5173" 
	@echo "  - Consul UI: http://localhost:8500"
	@echo "  - Qdrant UI: http://localhost:6333/dashboard"

stop: services-stop infra-stop ## 停止完整开发环境

restart: stop start ## 重启开发环境

status: ## 检查服务状态
	@echo "📊 服务状态检查..."
	@echo "=== 基础设施容器 ==="
	@cd infrastructure && docker-compose ps
	@echo "=== 微服务进程 ==="
	@ps aux | grep "go run cmd/main.go" | grep -v grep || echo "无运行中的微服务"
	@echo "=== 端口监听状态 ==="
	@lsof -i :5432 -i :6379 -i :8500 -i :6333 -i :8000-8005 -i :5173 2>/dev/null || echo "无相关端口监听"

# === 测试和构建 ===
test: ## 运行所有测试
	@echo "🧪 运行测试..."
	@for service in services/*/; do \
		if [ -f "$$service/go.mod" ]; then \
			echo "测试 $$service..."; \
			cd "$$service" && go test ./... && cd ../..; \
		fi; \
	done
	@cd frontend && npm run test

test-integration: ## 运行集成测试
	@echo "🧪 运行集成测试..."
	@./scripts/integration-tests.sh

build: ## 构建所有服务
	@echo "🔨 构建微服务..."
	@for service in services/*/; do \
		if [ -f "$$service/go.mod" ]; then \
			service_name=$$(basename "$$service"); \
			echo "构建 $$service_name..."; \
			cd "$$service" && go build -o "../../build/$$service_name" cmd/main.go && cd ../..; \
		fi; \
	done
	@make frontend-build

# === 开发工具 ===
migrate: ## 运行数据库迁移
	@echo "🗄️ 运行数据库迁移..."
	@./scripts/migrate.sh

seed: ## 填充测试数据
	@echo "🌱 填充测试数据..."
	@./scripts/seed.sh

logs: ## 查看所有日志
	@echo "📋 查看服务日志..."
	@make infra-logs &
	@tail -f logs/*.log 2>/dev/null || echo "无日志文件"

clean: ## 清理构建产物
	@echo "🧹 清理构建产物..."
	@rm -rf build/
	@rm -rf frontend/dist/
	@go clean -cache
	@cd frontend && npm run clean 2>/dev/null || true

# === 代码质量 ===
lint: ## 代码检查
	@echo "🔍 代码检查..."
	@golangci-lint run ./...
	@cd frontend && npm run lint

format: ## 代码格式化
	@echo "✨ 代码格式化..."
	@gofmt -s -w .
	@cd frontend && npm run format

# === 开发者工具 ===
generate: ## 生成代码
	@echo "⚡ 生成代码..."
	@go generate ./...

mod-tidy: ## 整理 Go 模块
	@echo "📚 整理 Go 模块..."
	@for service in services/*/; do \
		if [ -f "$$service/go.mod" ]; then \
			echo "整理 $$service 模块..."; \
			cd "$$service" && go mod tidy && cd ../..; \
		fi; \
	done

update-deps: ## 更新依赖
	@echo "📦 更新依赖..."
	@for service in services/*/; do \
		if [ -f "$$service/go.mod" ]; then \
			echo "更新 $$service 依赖..."; \
			cd "$$service" && go get -u ./... && go mod tidy && cd ../..; \
		fi; \
	done
	@cd frontend && npm update

# === 监控和调试 ===
monitor: ## 启动监控服务
	@echo "📊 启动监控服务..."
	@cd infrastructure && docker-compose --profile monitoring up -d prometheus grafana
	@echo "监控服务已启动："
	@echo "  - Prometheus: http://localhost:9090"
	@echo "  - Grafana: http://localhost:3001 (admin/admin123)"

debug-gateway: ## 调试网关服务
	@echo "🐛 启动网关服务调试模式..."
	@cd services/gateway-service && dlv debug cmd/main.go

# === 文档和帮助 ===
docs: ## 生成文档
	@echo "📚 生成 API 文档..."
	@swag init -g cmd/main.go -o docs/swagger
	@echo "API 文档已生成到 docs/swagger/"

api-test: ## API 接口测试
	@echo "🧪 运行 API 测试..."
	@./scripts/api-tests.sh
```

### 2. 热重载开发

**scripts/dev-watch.sh**：
```bash
#!/bin/bash

# 热重载开发脚本
set -e

SERVICE_NAME=${1:-"gateway-service"}
SERVICE_DIR="services/$SERVICE_NAME"

if [ ! -d "$SERVICE_DIR" ]; then
    echo "❌ 服务目录不存在: $SERVICE_DIR"
    exit 1
fi

echo "🔥 启动 $SERVICE_NAME 热重载开发模式..."

# 检查 air 是否安装
if ! command -v air &> /dev/null; then
    echo "📦 安装 air 热重载工具..."
    go install github.com/cosmtrek/air@latest
fi

# 进入服务目录
cd "$SERVICE_DIR"

# 创建 air 配置文件 (如果不存在)
if [ ! -f ".air.toml" ]; then
    cat > .air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/main.go"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_root = false

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
EOF
fi

# 启动热重载
echo "🚀 $SERVICE_NAME 热重载已启动"
echo "💡 修改代码将自动重新编译和重启服务"
air
```

### 3. 调试配置

**VS Code launch.json**：
```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "Debug Gateway Service",
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/services/gateway-service/cmd/main.go",
      "cwd": "${workspaceFolder}/services/gateway-service",
      "env": {
        "ENV": "development",
        "DEBUG": "true"
      },
      "args": []
    },
    {
      "name": "Debug User Service", 
      "type": "go",
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/services/user-service/cmd/main.go",
      "cwd": "${workspaceFolder}/services/user-service",
      "env": {
        "ENV": "development",
        "DEBUG": "true"
      }
    },
    {
      "name": "Debug Auth Service",
      "type": "go", 
      "request": "launch",
      "mode": "auto",
      "program": "${workspaceFolder}/services/auth-service/cmd/main.go",
      "cwd": "${workspaceFolder}/services/auth-service",
      "env": {
        "ENV": "development",
        "DEBUG": "true"
      }
    },
    {
      "name": "Attach to Process",
      "type": "go",
      "request": "attach",
      "mode": "local",
      "processId": "${command:pickProcess}"
    }
  ]
}
```

---

## 完整启动流程

### 1. 首次环境搭建

```bash
# 1. 克隆项目
git clone https://github.com/your-org/lyss-ai-platform.git
cd lyss-ai-platform

# 2. 初始化环境
make setup

# 3. 编辑环境变量 (重要!)
vim .env
# 配置 AI API 密钥等关键信息

# 4. 安装依赖
make install

# 5. 启动基础设施
make infra-start

# 6. 运行数据库迁移
make migrate

# 7. 填充测试数据 (可选)
make seed

# 8. 启动微服务
make services-start

# 9. 启动前端 (新终端)
make frontend-start

# 10. 检查服务状态
make status
```

### 2. 日常开发流程

```bash
# 启动完整开发环境
make start

# 或者分步启动
make infra-start      # 启动基础设施
make services-start   # 启动微服务
make frontend-start   # 启动前端 (新终端)

# 停止开发环境
make stop

# 重启环境
make restart

# 查看日志
make logs

# 运行测试
make test

# 代码检查和格式化
make lint
make format
```

### 3. 热重载开发

```bash
# 启动特定服务的热重载模式
./scripts/dev-watch.sh gateway-service

# 或者使用 VS Code 调试模式
# F5 启动调试，设置断点进行调试
```

---

## 开发工具和脚本

### 1. 数据库迁移脚本

**scripts/migrate.sh**：
```bash
#!/bin/bash

# 数据库迁移脚本
set -e

echo "🗄️ 运行数据库迁移..."

# 检查 golang-migrate 是否安装
if ! command -v migrate &> /dev/null; then
    echo "📦 安装 golang-migrate..."
    go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
fi

# 数据库连接配置
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-lyss}
DB_PASSWORD=${DB_PASSWORD:-lyss123}

# 服务列表
SERVICES=("user" "auth" "group" "credential" "gateway" "billing")

for SERVICE in "${SERVICES[@]}"; do
    DB_NAME="lyss_$SERVICE"
    MIGRATION_DIR="services/${SERVICE}-service/migrations"
    
    if [ -d "$MIGRATION_DIR" ]; then
        echo "📋 迁移 $SERVICE 服务数据库..."
        
        DATABASE_URL="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable"
        
        migrate -path "$MIGRATION_DIR" -database "$DATABASE_URL" up
        
        echo "✅ $SERVICE 服务数据库迁移完成"
    else
        echo "⚠️  $SERVICE 服务迁移目录不存在: $MIGRATION_DIR"
    fi
done

echo "🎉 所有数据库迁移完成！"
```

### 2. 测试数据填充

**scripts/seed.sh**：
```bash
#!/bin/bash

# 测试数据填充脚本
set -e

echo "🌱 填充测试数据..."

# 等待服务启动
sleep 5

# API基础URL
API_BASE="http://localhost:8000/api/v1"

# 创建测试用户
echo "👤 创建测试用户..."
curl -X POST "$API_BASE/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "admin",
    "email": "admin@lyss.ai",
    "password": "admin123",
    "full_name": "System Administrator"
  }'

curl -X POST "$API_BASE/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@lyss.ai", 
    "password": "test123",
    "full_name": "Test User"
  }'

# 登录并获取token
echo "🔑 获取认证token..."
ADMIN_TOKEN=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@lyss.ai",
    "password": "admin123"
  }' | jq -r '.token')

# 创建测试群组
echo "👥 创建测试群组..."
curl -X POST "$API_BASE/groups" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "Development Team",
    "description": "开发团队测试群组",
    "type": "team"
  }'

# 添加模型配置
echo "🤖 配置AI模型..."
curl -X POST "$API_BASE/models" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $ADMIN_TOKEN" \
  -d '{
    "name": "gpt-3.5-turbo",
    "provider": "openai",
    "enabled": true,
    "config": {
      "max_tokens": 4096,
      "temperature": 0.7
    }
  }'

echo "✅ 测试数据填充完成！"
echo "📋 测试账号："
echo "  - 管理员: admin@lyss.ai / admin123"
echo "  - 测试用户: test@lyss.ai / test123"
```

### 3. API 测试脚本

**scripts/api-tests.sh**：
```bash
#!/bin/bash

# API 接口测试脚本
set -e

echo "🧪 运行 API 接口测试..."

API_BASE="http://localhost:8000/api/v1"

# 健康检查
echo "🏥 健康检查..."
curl -f "$API_BASE/health" || {
    echo "❌ 健康检查失败"
    exit 1
}

# 用户注册测试
echo "👤 测试用户注册..."
REGISTER_RESPONSE=$(curl -s -X POST "$API_BASE/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "apitest",
    "email": "apitest@lyss.ai",
    "password": "test123",
    "full_name": "API Test User"
  }')

echo "$REGISTER_RESPONSE" | jq .

# 用户登录测试
echo "🔑 测试用户登录..."
LOGIN_RESPONSE=$(curl -s -X POST "$API_BASE/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "apitest@lyss.ai",
    "password": "test123"
  }')

TOKEN=$(echo "$LOGIN_RESPONSE" | jq -r '.token')

if [ "$TOKEN" == "null" ]; then
    echo "❌ 登录失败"
    exit 1
fi

echo "✅ 登录成功，Token: ${TOKEN:0:20}..."

# 获取用户信息测试
echo "📋 测试获取用户信息..."
curl -s "$API_BASE/users/me" \
  -H "Authorization: Bearer $TOKEN" \
  | jq .

# 创建群组测试
echo "👥 测试创建群组..."
GROUP_RESPONSE=$(curl -s -X POST "$API_BASE/groups" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "name": "API Test Group",
    "description": "API测试群组",
    "type": "team"
  }')

GROUP_ID=$(echo "$GROUP_RESPONSE" | jq -r '.id')
echo "✅ 群组创建成功，ID: $GROUP_ID"

# 聊天接口测试
echo "💬 测试聊天接口..."
curl -s -X POST "$API_BASE/chat/completions" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{
    "model": "gpt-3.5-turbo",
    "messages": [
      {
        "role": "user",
        "content": "Hello, this is an API test!"
      }
    ],
    "stream": false
  }' | jq .

echo "🎉 所有 API 测试通过！"
```

---

## 常见问题和解决方案

### 1. 端口冲突问题

```bash
# 检查端口占用
lsof -i :5432  # PostgreSQL
lsof -i :6379  # Redis
lsof -i :8500  # Consul
lsof -i :8000  # Gateway

# 停止占用端口的进程
sudo kill -9 $(lsof -t -i :8000)

# 修改端口配置
# 编辑 .env 文件或 docker-compose.yml 文件修改端口
```

### 2. Docker 权限问题

```bash
# Linux 下添加用户到 docker 组
sudo usermod -aG docker $USER
newgrp docker

# 重启 Docker 服务
sudo systemctl restart docker
```

### 3. Go 模块下载问题

```bash
# 设置 Go 代理 (中国地区)
go env -w GOPROXY=https://goproxy.cn,direct
go env -w GOSUMDB=sum.golang.google.cn

# 清理模块缓存
go clean -modcache

# 重新下载依赖
go mod download
```

### 4. 数据库连接问题

```bash
# 检查数据库是否正常启动
docker-compose logs postgres

# 手动连接测试
psql -h localhost -p 5432 -U lyss -d lyss_platform

# 重置数据库
make infra-clean
make infra-start
make migrate
```

### 5. 前端启动问题

```bash
# 清理 node_modules
cd frontend
rm -rf node_modules package-lock.json
npm install

# 检查 Node.js 版本
node --version  # 应该是 18+

# 检查端口冲突
lsof -i :5173
```

---

## 开发最佳实践

### 1. 代码规范

```bash
# 安装 pre-commit hooks
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# 创建 .golangci.yml 配置文件
cat > .golangci.yml << 'EOF'
run:
  timeout: 5m
  issues-exit-code: 1
  tests: true

linters:
  enable:
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - typecheck
    - unused
    - varcheck
    - deadcode
    - structcheck
    - misspell
    - unparam
    - unconvert
    - gofmt
    - goimports
    - golint
    - gocyclo
    - dupl

linters-settings:
  gocyclo:
    min-complexity: 15
  dupl:
    threshold: 100
  misspell:
    locale: US
EOF

# 运行代码检查
make lint
```

### 2. Git 工作流

```bash
# 功能分支开发
git checkout -b feature/user-auth
git add .
git commit -m "feat: implement user authentication"
git push origin feature/user-auth

# 提交信息规范
# feat: 新功能
# fix: 修复bug  
# docs: 文档更新
# style: 代码格式调整
# refactor: 代码重构
# test: 测试相关
# chore: 构建工具、依赖更新等
```

### 3. 调试技巧

```bash
# 查看服务日志
docker-compose logs -f postgres
tail -f logs/gateway-service.log

# 数据库查询调试
psql -h localhost -p 5432 -U lyss -d lyss_platform
\dt  # 查看表
\d users  # 查看表结构

# Redis 调试
redis-cli -h localhost -p 6379
KEYS *
GET user:123

# 网络调试
curl -v http://localhost:8000/api/v1/health
netstat -tulpn | grep :8000
```

### 4. 性能监控

```bash
# 启动监控服务
make monitor

# 查看系统资源使用
htop
docker stats

# Go 性能分析
go tool pprof http://localhost:8001/debug/pprof/profile
go tool pprof http://localhost:8001/debug/pprof/heap
```

---

## 总结

本开发环境搭建指南提供了完整的 Lyss AI Platform 本地开发环境配置流程。通过本地服务 + Docker 基础设施的混合模式，实现了开发调试的便利性和环境的一致性。

### 核心特点

1. **混合架构**: 微服务本地运行便于调试，基础设施Docker化保证一致性
2. **自动化工具**: 完整的 Makefile 和脚本支持一键启动
3. **热重载支持**: 代码修改自动重新编译和重启
4. **调试友好**: VS Code 调试配置和详细的日志输出
5. **测试完备**: 单元测试、集成测试和API测试支持

### 快速开始

```bash
# 一键启动开发环境
git clone https://github.com/your-org/lyss-ai-platform.git
cd lyss-ai-platform
make setup
make install
make start
```

现在您可以开始 Lyss AI Platform 的开发之旅了！

---

*本文档将随着项目发展持续更新，确保开发环境的最佳体验。*

**最后更新**: 2025-01-25  
**下次检查**: 2025-02-01