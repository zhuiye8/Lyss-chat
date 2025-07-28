-- 网关服务数据库迁移 - 创建路由配置表
-- 迁移版本: 002
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 路由配置表
CREATE TABLE routing_rules (
    id BIGSERIAL PRIMARY KEY,
    rule_name VARCHAR(100) NOT NULL,
    model_pattern VARCHAR(100),         -- 模型名称匹配规则
    user_group_pattern VARCHAR(100),    -- 用户组匹配规则
    routing_strategy SMALLINT DEFAULT 1, -- 1:轮询 2:权重 3:最少连接
    fallback_strategy JSONB,           -- 失败回退策略
    is_enabled BOOLEAN DEFAULT TRUE,
    priority INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_routing_rules_is_enabled ON routing_rules(is_enabled);
CREATE INDEX idx_routing_rules_priority ON routing_rules(priority);

-- 添加表注释
COMMENT ON TABLE routing_rules IS '路由配置规则表';
COMMENT ON COLUMN routing_rules.id IS '路由规则唯一标识符';
COMMENT ON COLUMN routing_rules.rule_name IS '规则名称';
COMMENT ON COLUMN routing_rules.model_pattern IS '模型名称匹配规则';
COMMENT ON COLUMN routing_rules.user_group_pattern IS '用户组匹配规则';
COMMENT ON COLUMN routing_rules.routing_strategy IS '路由策略：1-轮询，2-权重，3-最少连接';
COMMENT ON COLUMN routing_rules.fallback_strategy IS '失败回退策略（JSON格式）';
COMMENT ON COLUMN routing_rules.is_enabled IS '是否启用';
COMMENT ON COLUMN routing_rules.priority IS '优先级';
COMMENT ON COLUMN routing_rules.created_at IS '创建时间';
COMMENT ON COLUMN routing_rules.updated_at IS '更新时间';

COMMIT;