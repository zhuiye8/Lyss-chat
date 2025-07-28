-- 用户服务数据库迁移 - 创建用户配额表
-- 迁移版本: 003
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 用户配额表
CREATE TABLE user_quotas (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    total_quota BIGINT DEFAULT 0,      -- 总配额 (tokens)
    used_quota BIGINT DEFAULT 0,       -- 已用配额
    daily_limit BIGINT DEFAULT 0,      -- 日配额限制
    monthly_limit BIGINT DEFAULT 0,    -- 月配额限制
    quota_reset_at DATE,               -- 配额重置日期
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_user_quotas_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uk_user_quotas_user_id UNIQUE (user_id)
);

-- 索引设计
CREATE INDEX idx_user_quotas_user_id ON user_quotas(user_id);

-- 添加表注释
COMMENT ON TABLE user_quotas IS '用户配额管理表';
COMMENT ON COLUMN user_quotas.id IS '配额记录唯一标识符';
COMMENT ON COLUMN user_quotas.user_id IS '关联的用户ID';
COMMENT ON COLUMN user_quotas.total_quota IS '总配额（token数量）';
COMMENT ON COLUMN user_quotas.used_quota IS '已使用配额';
COMMENT ON COLUMN user_quotas.daily_limit IS '日配额限制';
COMMENT ON COLUMN user_quotas.monthly_limit IS '月配额限制';
COMMENT ON COLUMN user_quotas.quota_reset_at IS '配额重置日期';
COMMENT ON COLUMN user_quotas.created_at IS '创建时间';
COMMENT ON COLUMN user_quotas.updated_at IS '更新时间';

COMMIT;