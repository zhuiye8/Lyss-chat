-- 计费服务数据库迁移 - 创建使用记录表
-- 迁移版本: 001
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 使用记录表
CREATE TABLE usage_records (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    group_id BIGINT,
    credential_id BIGINT NOT NULL,
    request_id VARCHAR(100) NOT NULL,
    model_name VARCHAR(100) NOT NULL,
    provider_type VARCHAR(50) NOT NULL,
    
    -- Token 统计
    prompt_tokens INTEGER DEFAULT 0,
    completion_tokens INTEGER DEFAULT 0,
    total_tokens INTEGER DEFAULT 0,
    
    -- 成本计算
    input_cost DECIMAL(15,8) DEFAULT 0,    -- 输入成本 (USD)
    output_cost DECIMAL(15,8) DEFAULT 0,   -- 输出成本 (USD)  
    total_cost DECIMAL(15,8) DEFAULT 0,    -- 总成本 (USD)
    
    -- 时间信息
    usage_date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_usage_records_user_id ON usage_records(user_id);
CREATE INDEX idx_usage_records_usage_date ON usage_records(usage_date);
CREATE INDEX idx_usage_records_group_id ON usage_records(group_id);

-- 添加表注释
COMMENT ON TABLE usage_records IS '用户使用记录表';
COMMENT ON COLUMN usage_records.id IS '使用记录唯一标识符';
COMMENT ON COLUMN usage_records.user_id IS '使用者用户ID';
COMMENT ON COLUMN usage_records.group_id IS '关联的群组ID';
COMMENT ON COLUMN usage_records.credential_id IS '使用的凭证ID';
COMMENT ON COLUMN usage_records.request_id IS '关联的请求ID';
COMMENT ON COLUMN usage_records.model_name IS '使用的模型名称';
COMMENT ON COLUMN usage_records.provider_type IS 'AI供应商类型';
COMMENT ON COLUMN usage_records.prompt_tokens IS '输入token数量';
COMMENT ON COLUMN usage_records.completion_tokens IS '输出token数量';
COMMENT ON COLUMN usage_records.total_tokens IS '总token数量';
COMMENT ON COLUMN usage_records.input_cost IS '输入成本（美元）';
COMMENT ON COLUMN usage_records.output_cost IS '输出成本（美元）';
COMMENT ON COLUMN usage_records.total_cost IS '总成本（美元）';
COMMENT ON COLUMN usage_records.usage_date IS '使用日期';
COMMENT ON COLUMN usage_records.created_at IS '创建时间';

COMMIT;