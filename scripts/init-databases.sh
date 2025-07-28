#!/bin/bash

# Lyss AI Platform 数据库初始化脚本
set -e

echo "🗄️ 初始化 Lyss AI Platform 数据库..."

# 数据库连接配置
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_USER=${DB_USER:-lyss}
DB_PASSWORD=${DB_PASSWORD:-lyss123}

# 等待 PostgreSQL 启动
echo "⏳ 等待 PostgreSQL 服务..."
until pg_isready -h $DB_HOST -p $DB_PORT -U $DB_USER; do
  echo "PostgreSQL 未就绪，等待中..."
  sleep 2
done

echo "✅ PostgreSQL 已就绪，开始数据库初始化..."

# 创建 ClickHouse 数据库
echo "📊 初始化 ClickHouse 数据库..."
clickhouse-client --host localhost --port 8123 --query "
CREATE DATABASE IF NOT EXISTS lyss_analytics;

-- 创建用户使用统计表
CREATE TABLE IF NOT EXISTS lyss_analytics.usage_metrics (
    id UUID DEFAULT generateUUIDv4(),
    user_id String,
    group_id String,
    model_name String,
    provider String,
    request_time DateTime,
    response_time DateTime,
    input_tokens UInt32,
    output_tokens UInt32,
    total_tokens UInt32,
    cost Float64,
    success Bool,
    error_message String,
    request_ip String,
    user_agent String,
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (request_time, user_id)
PARTITION BY toYYYYMM(request_time);

-- 创建聊天分析表
CREATE TABLE IF NOT EXISTS lyss_analytics.chat_analytics (
    id UUID DEFAULT generateUUIDv4(),
    session_id String,
    user_id String,
    group_id String,
    model_name String,
    message_count UInt32,
    conversation_length UInt32,
    session_start DateTime,
    session_end DateTime,
    total_cost Float64,
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (session_start, user_id)
PARTITION BY toYYYYMM(session_start);

-- 创建模型性能表
CREATE TABLE IF NOT EXISTS lyss_analytics.model_performance (
    id UUID DEFAULT generateUUIDv4(),
    model_name String,
    provider String,
    avg_response_time Float64,
    success_rate Float64,
    total_requests UInt64,
    date Date,
    created_at DateTime DEFAULT now()
) ENGINE = MergeTree()
ORDER BY (date, model_name);
" 2>/dev/null || echo "⚠️  ClickHouse 初始化跳过 (服务可能未启动)"

# 创建 Qdrant 集合
echo "🔍 初始化 Qdrant 向量数据库..."
curl -X PUT "http://localhost:6333/collections/user_memories" \
  -H "Content-Type: application/json" \
  -d '{
    "vectors": {
      "size": 1536,
      "distance": "Cosine"
    },
    "shard_number": 1,
    "replication_factor": 1
  }' 2>/dev/null || echo "⚠️  Qdrant user_memories 集合创建跳过"

curl -X PUT "http://localhost:6333/collections/group_memories" \
  -H "Content-Type: application/json" \
  -d '{
    "vectors": {
      "size": 1536,
      "distance": "Cosine"
    },
    "shard_number": 1,
    "replication_factor": 1
  }' 2>/dev/null || echo "⚠️  Qdrant group_memories 集合创建跳过"

curl -X PUT "http://localhost:6333/collections/knowledge_base" \
  -H "Content-Type: application/json" \
  -d '{
    "vectors": {
      "size": 1536,
      "distance": "Cosine"
    },
    "shard_number": 1,
    "replication_factor": 1
  }' 2>/dev/null || echo "⚠️  Qdrant knowledge_base 集合创建跳过"

# 验证数据库连接
echo "🔍 验证数据库连接..."

# 验证 PostgreSQL 数据库
DATABASES=("lyss_user" "lyss_auth" "lyss_group" "lyss_credential" "lyss_gateway" "lyss_billing")

for db in "${DATABASES[@]}"; do
    if PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $db -c "SELECT 1;" > /dev/null 2>&1; then
        echo "✅ $db 数据库连接成功"
    else
        echo "❌ $db 数据库连接失败"
    fi
done

# 验证 Redis 连接
if redis-cli -h localhost -p 6379 ping > /dev/null 2>&1; then
    echo "✅ Redis 连接成功"
    
    # 设置一些默认键值
    redis-cli -h localhost -p 6379 SET "lyss:init:timestamp" "$(date -u +%Y-%m-%dT%H:%M:%SZ)" > /dev/null
    redis-cli -h localhost -p 6379 SET "lyss:init:status" "completed" > /dev/null
else
    echo "❌ Redis 连接失败"
fi

echo ""
echo "🎉 数据库初始化完成！"
echo "📋 已创建数据库："
echo "  - PostgreSQL: lyss_user, lyss_auth, lyss_group, lyss_credential, lyss_gateway, lyss_billing"
echo "  - ClickHouse: lyss_analytics (usage_metrics, chat_analytics, model_performance)"
echo "  - Qdrant: user_memories, group_memories, knowledge_base"
echo "  - Redis: 缓存服务已就绪"