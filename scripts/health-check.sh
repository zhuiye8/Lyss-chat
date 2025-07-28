#!/bin/bash

# Lyss AI Platform 健康检查脚本
set -e

echo "🏥 Lyss AI Platform 健康检查"
echo "================================"

# 检查基础设施服务
echo ""
echo "📊 基础设施服务状态："

# PostgreSQL 检查
if pg_isready -h localhost -p 5432 -U lyss > /dev/null 2>&1; then
    echo "✅ PostgreSQL: 运行正常"
else
    echo "❌ PostgreSQL: 连接失败"
fi

# Redis 检查
if redis-cli -h localhost -p 6379 ping > /dev/null 2>&1; then
    echo "✅ Redis: 运行正常"
else
    echo "❌ Redis: 连接失败"
fi

# ClickHouse 检查
if curl -s http://localhost:8123/ping | grep -q "Ok"; then
    echo "✅ ClickHouse: 运行正常"
else
    echo "❌ ClickHouse: 连接失败"
fi

# Consul 检查
if curl -s http://localhost:8500/v1/status/leader > /dev/null 2>&1; then
    echo "✅ Consul: 运行正常"
else
    echo "❌ Consul: 连接失败"
fi

# Qdrant 检查
if curl -s http://localhost:6333/health | grep -q "ok"; then
    echo "✅ Qdrant: 运行正常"
else
    echo "❌ Qdrant: 连接失败"
fi

# 检查微服务
echo ""
echo "🔧 微服务状态："

SERVICES=("gateway:8000" "user:8001" "auth:8002" "group:8003" "credential:8004" "billing:8005")

for service in "${SERVICES[@]}"; do
    IFS=':' read -r name port <<< "$service"
    if curl -s "http://localhost:$port/health" > /dev/null 2>&1; then
        echo "✅ $name-service: 运行正常"
    else
        echo "❌ $name-service: 无响应"
    fi
done

# 检查监控服务 (可选)
echo ""
echo "📈 监控服务状态："

if curl -s http://localhost:9090/-/healthy > /dev/null 2>&1; then
    echo "✅ Prometheus: 运行正常"
else
    echo "⚠️  Prometheus: 未启动 (可选)"
fi

if curl -s http://localhost:3001/api/health > /dev/null 2>&1; then
    echo "✅ Grafana: 运行正常"
else
    echo "⚠️  Grafana: 未启动 (可选)"
fi

if curl -s http://localhost:16686 > /dev/null 2>&1; then
    echo "✅ Jaeger: 运行正常"
else
    echo "⚠️  Jaeger: 未启动 (可选)"
fi

# 端口占用检查
echo ""
echo "🔌 端口占用状态："
PORTS=(5432 6379 8123 8500 6333 8000 8001 8002 8003 8004 8005 5173)

for port in "${PORTS[@]}"; do
    if lsof -i :$port > /dev/null 2>&1; then
        echo "✅ Port $port: 已占用"
    else
        echo "❌ Port $port: 空闲"
    fi
done

echo ""
echo "🎯 健康检查完成！"