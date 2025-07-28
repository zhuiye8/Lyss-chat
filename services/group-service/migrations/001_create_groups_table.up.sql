-- 群组服务数据库迁移 - 创建群组基础表
-- 迁移版本: 001
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 群组基础表
CREATE TABLE groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    avatar_url VARCHAR(500),
    is_public BOOLEAN DEFAULT FALSE,   -- 是否公开可搜索
    join_mode SMALLINT DEFAULT 1,      -- 1:邀请制 2:申请制 3:开放制
    max_members INTEGER DEFAULT 50,
    status SMALLINT DEFAULT 1,         -- 1:正常 2:冻结 3:解散
    owner_user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_groups_owner_user_id ON groups(owner_user_id);
CREATE INDEX idx_groups_status ON groups(status);
CREATE INDEX idx_groups_is_public ON groups(is_public);

-- 添加表注释
COMMENT ON TABLE groups IS '群组基础信息表';
COMMENT ON COLUMN groups.id IS '群组唯一标识符';
COMMENT ON COLUMN groups.name IS '群组名称';
COMMENT ON COLUMN groups.description IS '群组描述';
COMMENT ON COLUMN groups.avatar_url IS '群组头像URL';
COMMENT ON COLUMN groups.is_public IS '是否公开可搜索';
COMMENT ON COLUMN groups.join_mode IS '加入模式：1-邀请制，2-申请制，3-开放制';
COMMENT ON COLUMN groups.max_members IS '最大成员数量';
COMMENT ON COLUMN groups.status IS '群组状态：1-正常，2-冻结，3-解散';
COMMENT ON COLUMN groups.owner_user_id IS '群组所有者用户ID';
COMMENT ON COLUMN groups.created_at IS '创建时间';
COMMENT ON COLUMN groups.updated_at IS '更新时间';

COMMIT;