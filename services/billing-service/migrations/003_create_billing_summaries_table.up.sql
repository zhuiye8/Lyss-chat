-- 计费服务数据库迁移 - 创建账单汇总表
-- 迁移版本: 003
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 账单汇总表
CREATE TABLE billing_summaries (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    group_id BIGINT,
    billing_period DATE NOT NULL,       -- 账单周期 (YYYY-MM-01)
    
    -- 统计数据
    total_requests INTEGER DEFAULT 0,
    total_tokens BIGINT DEFAULT 0,
    total_cost DECIMAL(15,2) DEFAULT 0,
    
    -- 按供应商分组统计
    provider_breakdown JSONB,           -- {provider: {requests, tokens, cost}}
    
    status SMALLINT DEFAULT 1,          -- 1:计算中 2:已确认 3:已支付
    generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    confirmed_at TIMESTAMP,
    
    CONSTRAINT uk_billing_summaries_user_period 
        UNIQUE (user_id, group_id, billing_period)
);

-- 索引设计
CREATE INDEX idx_billing_summaries_user_id ON billing_summaries(user_id);
CREATE INDEX idx_billing_summaries_billing_period ON billing_summaries(billing_period);

-- 添加表注释
COMMENT ON TABLE billing_summaries IS '账单汇总表';
COMMENT ON COLUMN billing_summaries.id IS '账单汇总唯一标识符';
COMMENT ON COLUMN billing_summaries.user_id IS '用户ID';
COMMENT ON COLUMN billing_summaries.group_id IS '群组ID';
COMMENT ON COLUMN billing_summaries.billing_period IS '账单周期（年月-01格式）';
COMMENT ON COLUMN billing_summaries.total_requests IS '总请求数';
COMMENT ON COLUMN billing_summaries.total_tokens IS '总token数';
COMMENT ON COLUMN billing_summaries.total_cost IS '总成本';
COMMENT ON COLUMN billing_summaries.provider_breakdown IS '按供应商分组统计（JSON格式）';
COMMENT ON COLUMN billing_summaries.status IS '账单状态：1-计算中，2-已确认，3-已支付';
COMMENT ON COLUMN billing_summaries.generated_at IS '生成时间';
COMMENT ON COLUMN billing_summaries.confirmed_at IS '确认时间';

COMMIT;