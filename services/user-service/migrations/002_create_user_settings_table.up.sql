-- 用户服务数据库迁移 - 创建用户配置表
-- 迁移版本: 002
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 用户配置表
CREATE TABLE user_settings (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    setting_key VARCHAR(100) NOT NULL,
    setting_value JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_user_settings_user_id 
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT uk_user_settings_user_key 
        UNIQUE (user_id, setting_key)
);

-- 索引设计
CREATE INDEX idx_user_settings_user_id ON user_settings(user_id);

-- 添加表注释
COMMENT ON TABLE user_settings IS '用户个性化配置表';
COMMENT ON COLUMN user_settings.id IS '配置记录唯一标识符';
COMMENT ON COLUMN user_settings.user_id IS '关联的用户ID';
COMMENT ON COLUMN user_settings.setting_key IS '配置项键名';
COMMENT ON COLUMN user_settings.setting_value IS '配置项值（JSON格式）';
COMMENT ON COLUMN user_settings.created_at IS '创建时间';
COMMENT ON COLUMN user_settings.updated_at IS '更新时间';

COMMIT;