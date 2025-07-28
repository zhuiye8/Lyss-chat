-- 群组服务数据库迁移 - 创建群组成员表
-- 迁移版本: 002
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 群组成员表
CREATE TABLE group_members (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL,
    user_id BIGINT NOT NULL,
    role SMALLINT DEFAULT 1,           -- 1:成员 2:管理员 3:群主
    status SMALLINT DEFAULT 1,         -- 1:正常 2:禁言 3:踢出
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT fk_group_members_group_id 
        FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE,
    CONSTRAINT uk_group_members_group_user 
        UNIQUE (group_id, user_id)
);

-- 索引设计
CREATE INDEX idx_group_members_group_id ON group_members(group_id);
CREATE INDEX idx_group_members_user_id ON group_members(user_id);
CREATE INDEX idx_group_members_role ON group_members(role);

-- 添加表注释
COMMENT ON TABLE group_members IS '群组成员表';
COMMENT ON COLUMN group_members.id IS '成员记录唯一标识符';
COMMENT ON COLUMN group_members.group_id IS '关联的群组ID';
COMMENT ON COLUMN group_members.user_id IS '关联的用户ID';
COMMENT ON COLUMN group_members.role IS '成员角色：1-成员，2-管理员，3-群主';
COMMENT ON COLUMN group_members.status IS '成员状态：1-正常，2-禁言，3-踢出';
COMMENT ON COLUMN group_members.joined_at IS '加入时间';
COMMENT ON COLUMN group_members.updated_at IS '更新时间';

COMMIT;