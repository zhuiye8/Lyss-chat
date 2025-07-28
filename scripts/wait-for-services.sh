#!/bin/bash

# Lyss AI Platform 等待服务启动脚本
set -e

echo "🚀 等待基础设施服务启动..."

# 进入基础设施目录
cd "$(dirname "$0")/../infrastructure"

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

# 等待 ClickHouse
echo "⏳ 等待 ClickHouse..."
until curl -s http://localhost:8123/ping | grep -q "Ok"; do
  echo "ClickHouse 未就绪，等待中..."
  sleep 2
done
echo "✅ ClickHouse 已就绪"

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
echo ""
echo "📋 服务状态："
echo "- PostgreSQL: http://localhost:5432"
echo "- Redis: http://localhost:6379"  
echo "- ClickHouse: http://localhost:8123"
echo "- Consul: http://localhost:8500"
echo "- Qdrant: http://localhost:6333"

# 检查可选监控服务
if docker-compose ps prometheus | grep -q "Up"; then
  echo "- Prometheus: http://localhost:9090"
fi

if docker-compose ps grafana | grep -q "Up"; then
  echo "- Grafana: http://localhost:3001 (admin/admin123)"
fi

if docker-compose ps jaeger | grep -q "Up"; then
  echo "- Jaeger: http://localhost:16686"
fi

echo ""
echo "🔧 开发环境已准备就绪！"