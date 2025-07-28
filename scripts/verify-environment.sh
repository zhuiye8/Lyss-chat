#!/bin/bash

# Lyss AI Platform 环境验证脚本
set -e

echo "🔍 Lyss AI Platform 环境验证"
echo "===================================="
echo ""

# 检查Docker容器状态
echo "📦 Docker 容器状态："
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" --filter "name=lyss-"

echo ""
echo "🌐 服务连接测试："

# PostgreSQL
if docker exec lyss-postgres pg_isready -U lyss > /dev/null 2>&1; then
    echo "✅ PostgreSQL: 连接正常 (localhost:5432)"
    DB_COUNT=$(docker exec lyss-postgres psql -U lyss -d lyss_platform -t -c "SELECT COUNT(*) FROM pg_database WHERE datname LIKE 'lyss_%';" | tr -d ' ')
    echo "   📊 已创建 $DB_COUNT 个数据库"
else
    echo "❌ PostgreSQL: 连接失败"
fi

# Redis
if docker exec lyss-redis redis-cli ping > /dev/null 2>&1; then
    echo "✅ Redis: 连接正常 (localhost:6379)"
    REDIS_VERSION=$(docker exec lyss-redis redis-cli INFO server | grep redis_version | cut -d: -f2 | tr -d '\r')
    echo "   🔖 版本: $REDIS_VERSION"
else
    echo "❌ Redis: 连接失败"
fi

# ClickHouse
if docker exec lyss-clickhouse clickhouse-client --query "SELECT 1" > /dev/null 2>&1; then
    echo "✅ ClickHouse: 连接正常 (localhost:8123)"
    DB_EXISTS=$(docker exec lyss-clickhouse clickhouse-client --query "EXISTS DATABASE lyss_analytics")
    if [ "$DB_EXISTS" = "1" ]; then
        echo "   📊 lyss_analytics 数据库已创建"
    fi
else
    echo "❌ ClickHouse: 连接失败"
fi

# Consul
if curl -s http://localhost:8500/v1/status/leader > /dev/null 2>&1; then
    echo "✅ Consul: 连接正常 (localhost:8500)"
    echo "   🌐 UI访问: http://localhost:8500/ui"
else
    echo "❌ Consul: 连接失败"
fi

# Qdrant
QDRANT_STATUS=$(curl -s http://localhost:6333/health 2>/dev/null || echo "failed")
if [ "$QDRANT_STATUS" != "failed" ]; then
    echo "✅ Qdrant: 连接正常 (localhost:6333)"
    echo "   🌐 UI访问: http://localhost:6333/dashboard"
else
    echo "❌ Qdrant: 连接失败"
fi

echo ""
echo "🔌 端口监听状态："
PORTS=(5432 6379 8123 8500 6333)
for port in "${PORTS[@]}"; do
    if netstat -tuln 2>/dev/null | grep -q ":$port "; then
        echo "✅ Port $port: 监听中"
    else
        echo "❌ Port $port: 未监听"
    fi
done

echo ""
echo "📋 环境配置验证："

# 检查环境变量文件
if [ -f ".env.example" ]; then
    echo "✅ .env.example: 模板文件存在"
else
    echo "❌ .env.example: 模板文件缺失"
fi

# 检查配置文件
CONFIG_FILES=(
    "infrastructure/docker-compose.yml"
    "infrastructure/configs/postgres/init.sql"
    "infrastructure/configs/redis/redis.conf"
    "infrastructure/configs/clickhouse/config.xml"
    "infrastructure/configs/consul/consul.hcl"
    "infrastructure/configs/qdrant/config.yaml"
)

for config in "${CONFIG_FILES[@]}"; do
    if [ -f "$config" ]; then
        echo "✅ $config: 存在"
    else
        echo "❌ $config: 缺失"
    fi
done

# 检查脚本文件
SCRIPTS=(
    "scripts/wait-for-services.sh"
    "scripts/health-check.sh"
    "scripts/init-databases.sh"
    "scripts/verify-environment.sh"
)

echo ""
echo "🛠️  开发脚本："
for script in "${SCRIPTS[@]}"; do
    if [ -f "$script" ] && [ -x "$script" ]; then
        echo "✅ $script: 可执行"
    else
        echo "❌ $script: 缺失或无执行权限"
    fi
done

echo ""
echo "🎯 验收标准检查："

# 验收标准1: make health-check
echo "1. 健康检查: 基础服务运行状态"
HEALTHY_SERVICES=0
if docker exec lyss-postgres pg_isready -U lyss > /dev/null 2>&1; then ((HEALTHY_SERVICES++)); fi
if docker exec lyss-redis redis-cli ping > /dev/null 2>&1; then ((HEALTHY_SERVICES++)); fi
if docker exec lyss-clickhouse clickhouse-client --query "SELECT 1" > /dev/null 2>&1; then ((HEALTHY_SERVICES++)); fi
if curl -s http://localhost:8500/v1/status/leader > /dev/null 2>&1; then ((HEALTHY_SERVICES++)); fi
if curl -s http://localhost:6333/health > /dev/null 2>&1; then ((HEALTHY_SERVICES++)); fi

if [ $HEALTHY_SERVICES -eq 5 ]; then
    echo "   ✅ 所有基础服务 (5/5) 运行正常"
else
    echo "   ⚠️  部分服务异常 ($HEALTHY_SERVICES/5)"
fi

# 验收标准2: make db-ping  
echo "2. 数据库连接: PostgreSQL + Redis + ClickHouse"
DB_CONNECTIONS=0
if docker exec lyss-postgres pg_isready -U lyss > /dev/null 2>&1; then ((DB_CONNECTIONS++)); fi
if docker exec lyss-redis redis-cli ping > /dev/null 2>&1; then ((DB_CONNECTIONS++)); fi
if docker exec lyss-clickhouse clickhouse-client --query "SELECT 1" > /dev/null 2>&1; then ((DB_CONNECTIONS++)); fi

if [ $DB_CONNECTIONS -eq 3 ]; then
    echo "   ✅ 所有数据库 (3/3) 连接成功"
else
    echo "   ⚠️  部分数据库连接异常 ($DB_CONNECTIONS/3)"
fi

# 验收标准3: make monitoring
echo "3. 监控服务: 可选组件状态"
echo "   ⚠️  监控服务未启动 (需要时可用 docker run 启动 Prometheus + Grafana)"

echo ""
if [ $HEALTHY_SERVICES -eq 5 ] && [ $DB_CONNECTIONS -eq 3 ]; then
    echo "🎉 阶段0 - 开发环境准备: 验证通过!"
    echo "📍 下一步: 可以开始阶段1 - 数据库实施"
    echo ""
    echo "🔗 服务访问地址:"
    echo "   - PostgreSQL: localhost:5432 (用户: lyss, 密码: lyss123)"
    echo "   - Redis: localhost:6379"
    echo "   - ClickHouse: http://localhost:8123"
    echo "   - Consul UI: http://localhost:8500"
    echo "   - Qdrant UI: http://localhost:6333/dashboard"
else
    echo "❌ 阶段0验证失败，请检查服务状态"
    exit 1
fi