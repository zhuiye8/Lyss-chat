-- 凭证服务数据库迁移 - 创建模型权限表
-- 迁移版本: 002
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 模型权限表
CREATE TABLE model_permissions (
    id BIGSERIAL PRIMARY KEY,
    credential_id BIGINT NOT NULL,
    model_name VARCHAR(100) NOT NULL,
    is_enabled BOOLEAN DEFAULT TRUE,
    priority INTEGER DEFAULT 0,        -- 优先级，数字越大优先级越高
    weight INTEGER DEFAULT 1,          -- 负载均衡权重
    daily_quota BIGINT,                -- 日配额限制
    monthly_quota BIGINT,              -- 月配额限制
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_model_permissions_credential_id 
        FOREIGN KEY (credential_id) REFERENCES provider_credentials(id) ON DELETE CASCADE,
    CONSTRAINT uk_model_permissions_credential_model 
        UNIQUE (credential_id, model_name)
);

-- 索引设计
CREATE INDEX idx_model_permissions_credential_id ON model_permissions(credential_id);
CREATE INDEX idx_model_permissions_model_name ON model_permissions(model_name);

-- 添加表注释
COMMENT ON TABLE model_permissions IS '模型访问权限配置表';
COMMENT ON COLUMN model_permissions.id IS '权限记录唯一标识符';
COMMENT ON COLUMN model_permissions.credential_id IS '关联的凭证ID';
COMMENT ON COLUMN model_permissions.model_name IS '模型名称';
COMMENT ON COLUMN model_permissions.is_enabled IS '是否启用';
COMMENT ON COLUMN model_permissions.priority IS '优先级（数字越大优先级越高）';
COMMENT ON COLUMN model_permissions.weight IS '负载均衡权重';
COMMENT ON COLUMN model_permissions.daily_quota IS '日配额限制';
COMMENT ON COLUMN model_permissions.monthly_quota IS '月配额限制';
COMMENT ON COLUMN model_permissions.created_at IS '创建时间';
COMMENT ON COLUMN model_permissions.updated_at IS '更新时间';

COMMIT;