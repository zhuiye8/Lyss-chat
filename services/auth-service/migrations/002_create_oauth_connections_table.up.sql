-- 认证服务数据库迁移 - 创建OAuth集成表
-- 迁移版本: 002
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- OAuth 集成表
CREATE TABLE oauth_connections (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    provider VARCHAR(50) NOT NULL,     -- github, google, wechat
    provider_user_id VARCHAR(100) NOT NULL,
    access_token_encrypted TEXT,       -- 加密存储
    refresh_token_encrypted TEXT,
    token_expires_at TIMESTAMP,
    user_info JSONB,                   -- 用户在第三方平台的信息
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT uk_oauth_provider_user 
        UNIQUE (provider, provider_user_id)
);

-- 索引设计
CREATE INDEX idx_oauth_connections_user_id ON oauth_connections(user_id);
CREATE INDEX idx_oauth_connections_provider ON oauth_connections(provider);

-- 添加表注释
COMMENT ON TABLE oauth_connections IS 'OAuth第三方平台集成表';
COMMENT ON COLUMN oauth_connections.id IS 'OAuth连接唯一标识符';
COMMENT ON COLUMN oauth_connections.user_id IS '关联的用户ID';
COMMENT ON COLUMN oauth_connections.provider IS '第三方平台名称（github, google, wechat等）';
COMMENT ON COLUMN oauth_connections.provider_user_id IS '第三方平台的用户ID';
COMMENT ON COLUMN oauth_connections.access_token_encrypted IS '加密存储的访问令牌';
COMMENT ON COLUMN oauth_connections.refresh_token_encrypted IS '加密存储的刷新令牌';
COMMENT ON COLUMN oauth_connections.token_expires_at IS '令牌过期时间';
COMMENT ON COLUMN oauth_connections.user_info IS '第三方平台用户信息（JSON格式）';
COMMENT ON COLUMN oauth_connections.is_active IS '连接是否激活';
COMMENT ON COLUMN oauth_connections.created_at IS '创建时间';
COMMENT ON COLUMN oauth_connections.updated_at IS '更新时间';

COMMIT;