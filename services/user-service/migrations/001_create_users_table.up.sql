-- 用户服务数据库迁移 - 创建用户基础表
-- 迁移版本: 001
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 用户基础表
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(100),
    avatar_url VARCHAR(500),
    status SMALLINT DEFAULT 1,        -- 1:活跃 2:冻结 3:删除
    role SMALLINT DEFAULT 1,          -- 1:普通用户 2:管理员 3:超级管理员
    email_verified BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);

-- 添加表注释
COMMENT ON TABLE users IS '用户基础信息表';
COMMENT ON COLUMN users.id IS '用户唯一标识符';
COMMENT ON COLUMN users.username IS '用户名，唯一';
COMMENT ON COLUMN users.email IS '邮箱地址，唯一';
COMMENT ON COLUMN users.password_hash IS 'bcrypt加密的密码哈希';
COMMENT ON COLUMN users.display_name IS '显示名称';
COMMENT ON COLUMN users.avatar_url IS '头像URL';
COMMENT ON COLUMN users.status IS '用户状态：1-活跃，2-冻结，3-删除';
COMMENT ON COLUMN users.role IS '用户角色：1-普通用户，2-管理员，3-超级管理员';
COMMENT ON COLUMN users.email_verified IS '邮箱是否已验证';
COMMENT ON COLUMN users.last_login_at IS '最后登录时间';
COMMENT ON COLUMN users.created_at IS '创建时间';
COMMENT ON COLUMN users.updated_at IS '更新时间';

COMMIT;