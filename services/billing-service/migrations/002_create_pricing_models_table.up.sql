-- 计费服务数据库迁移 - 创建计费规则表
-- 迁移版本: 002
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 计费规则表
CREATE TABLE pricing_models (
    id BIGSERIAL PRIMARY KEY,
    provider_type VARCHAR(50) NOT NULL,
    model_name VARCHAR(100) NOT NULL,
    input_price_per_1k DECIMAL(10,8) NOT NULL,   -- 每1K input tokens价格
    output_price_per_1k DECIMAL(10,8) NOT NULL,  -- 每1K output tokens价格
    currency VARCHAR(10) DEFAULT 'USD',
    effective_date DATE NOT NULL,
    expiry_date DATE,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT uk_pricing_models_provider_model_date 
        UNIQUE (provider_type, model_name, effective_date)
);

-- 索引设计
CREATE INDEX idx_pricing_models_provider_model ON pricing_models(provider_type, model_name);
CREATE INDEX idx_pricing_models_effective_date ON pricing_models(effective_date);

-- 添加表注释
COMMENT ON TABLE pricing_models IS '计费规则表';
COMMENT ON COLUMN pricing_models.id IS '计费规则唯一标识符';
COMMENT ON COLUMN pricing_models.provider_type IS 'AI供应商类型';
COMMENT ON COLUMN pricing_models.model_name IS '模型名称';
COMMENT ON COLUMN pricing_models.input_price_per_1k IS '每1K输入token价格';
COMMENT ON COLUMN pricing_models.output_price_per_1k IS '每1K输出token价格';
COMMENT ON COLUMN pricing_models.currency IS '货币单位';
COMMENT ON COLUMN pricing_models.effective_date IS '生效日期';
COMMENT ON COLUMN pricing_models.expiry_date IS '失效日期';
COMMENT ON COLUMN pricing_models.is_active IS '是否启用';
COMMENT ON COLUMN pricing_models.created_at IS '创建时间';
COMMENT ON COLUMN pricing_models.updated_at IS '更新时间';

COMMIT;