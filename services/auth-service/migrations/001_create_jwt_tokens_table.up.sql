-- 认证服务数据库迁移 - 创建JWT令牌表
-- 迁移版本: 001
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- JWT 令牌表
CREATE TABLE jwt_tokens (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token_hash VARCHAR(255) NOT NULL,  -- JWT token 的 hash
    refresh_token_hash VARCHAR(255),   -- Refresh token 的 hash
    expires_at TIMESTAMP NOT NULL,
    refresh_expires_at TIMESTAMP,
    is_revoked BOOLEAN DEFAULT FALSE,
    device_info JSONB,                 -- 设备信息
    ip_address INET,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT uk_jwt_tokens_token_hash UNIQUE (token_hash)
);

-- 索引设计
CREATE INDEX idx_jwt_tokens_user_id ON jwt_tokens(user_id);
CREATE INDEX idx_jwt_tokens_expires_at ON jwt_tokens(expires_at);
CREATE INDEX idx_jwt_tokens_is_revoked ON jwt_tokens(is_revoked);

-- 添加表注释
COMMENT ON TABLE jwt_tokens IS 'JWT访问令牌管理表';
COMMENT ON COLUMN jwt_tokens.id IS '令牌记录唯一标识符';
COMMENT ON COLUMN jwt_tokens.user_id IS '关联的用户ID';
COMMENT ON COLUMN jwt_tokens.token_hash IS 'JWT访问令牌的SHA-256哈希值';
COMMENT ON COLUMN jwt_tokens.refresh_token_hash IS '刷新令牌的SHA-256哈希值';
COMMENT ON COLUMN jwt_tokens.expires_at IS '访问令牌过期时间';
COMMENT ON COLUMN jwt_tokens.refresh_expires_at IS '刷新令牌过期时间';
COMMENT ON COLUMN jwt_tokens.is_revoked IS '令牌是否已被撤销';
COMMENT ON COLUMN jwt_tokens.device_info IS '设备信息（JSON格式）';
COMMENT ON COLUMN jwt_tokens.ip_address IS '创建令牌时的IP地址';
COMMENT ON COLUMN jwt_tokens.created_at IS '创建时间';

COMMIT;