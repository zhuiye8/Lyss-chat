-- 凭证服务数据库迁移 - 创建供应商凭证表
-- 迁移版本: 001
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 供应商凭证表
CREATE TABLE provider_credentials (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    provider_type VARCHAR(50) NOT NULL, -- openai, claude, gemini
    scope VARCHAR(20) NOT NULL,         -- personal, group
    scope_id BIGINT NOT NULL,           -- user_id 或 group_id
    api_key_encrypted TEXT NOT NULL,    -- 加密存储的 API Key
    api_endpoint VARCHAR(500),
    model_config JSONB,                 -- 模型配置 (temperature, max_tokens等)
    rate_limit_config JSONB,            -- 速率限制配置
    status SMALLINT DEFAULT 1,          -- 1:正常 2:禁用 3:异常
    last_used_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_provider_credentials_scope ON provider_credentials(scope, scope_id);
CREATE INDEX idx_provider_credentials_provider_type ON provider_credentials(provider_type);
CREATE INDEX idx_provider_credentials_status ON provider_credentials(status);

-- 添加表注释
COMMENT ON TABLE provider_credentials IS 'AI供应商凭证管理表';
COMMENT ON COLUMN provider_credentials.id IS '凭证记录唯一标识符';
COMMENT ON COLUMN provider_credentials.name IS '凭证名称';
COMMENT ON COLUMN provider_credentials.provider_type IS 'AI供应商类型（openai, claude, gemini等）';
COMMENT ON COLUMN provider_credentials.scope IS '凭证作用域（personal-个人，group-群组）';
COMMENT ON COLUMN provider_credentials.scope_id IS '作用域ID（用户ID或群组ID）';
COMMENT ON COLUMN provider_credentials.api_key_encrypted IS '加密存储的API密钥';
COMMENT ON COLUMN provider_credentials.api_endpoint IS 'API端点地址';
COMMENT ON COLUMN provider_credentials.model_config IS '模型配置（JSON格式）';
COMMENT ON COLUMN provider_credentials.rate_limit_config IS '速率限制配置（JSON格式）';
COMMENT ON COLUMN provider_credentials.status IS '凭证状态：1-正常，2-禁用，3-异常';
COMMENT ON COLUMN provider_credentials.last_used_at IS '最后使用时间';
COMMENT ON COLUMN provider_credentials.created_at IS '创建时间';
COMMENT ON COLUMN provider_credentials.updated_at IS '更新时间';

COMMIT;