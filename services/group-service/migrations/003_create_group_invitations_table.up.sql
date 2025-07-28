-- 群组服务数据库迁移 - 创建群组邀请表
-- 迁移版本: 003
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 群组邀请表
CREATE TABLE group_invitations (
    id BIGSERIAL PRIMARY KEY,
    group_id BIGINT NOT NULL,
    inviter_user_id BIGINT NOT NULL,
    invitee_email VARCHAR(255),
    invitee_user_id BIGINT,
    invitation_code VARCHAR(50) UNIQUE NOT NULL,
    status SMALLINT DEFAULT 1,         -- 1:待处理 2:已接受 3:已拒绝 4:已过期
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    processed_at TIMESTAMP,
    
    CONSTRAINT fk_group_invitations_group_id 
        FOREIGN KEY (group_id) REFERENCES groups(id) ON DELETE CASCADE
);

-- 索引设计
CREATE INDEX idx_group_invitations_invitation_code ON group_invitations(invitation_code);
CREATE INDEX idx_group_invitations_group_id ON group_invitations(group_id);
CREATE INDEX idx_group_invitations_status ON group_invitations(status);

-- 添加表注释
COMMENT ON TABLE group_invitations IS '群组邀请表';
COMMENT ON COLUMN group_invitations.id IS '邀请记录唯一标识符';
COMMENT ON COLUMN group_invitations.group_id IS '关联的群组ID';
COMMENT ON COLUMN group_invitations.inviter_user_id IS '邀请人用户ID';
COMMENT ON COLUMN group_invitations.invitee_email IS '被邀请人邮箱';
COMMENT ON COLUMN group_invitations.invitee_user_id IS '被邀请人用户ID';
COMMENT ON COLUMN group_invitations.invitation_code IS '邀请码';
COMMENT ON COLUMN group_invitations.status IS '邀请状态：1-待处理，2-已接受，3-已拒绝，4-已过期';
COMMENT ON COLUMN group_invitations.expires_at IS '过期时间';
COMMENT ON COLUMN group_invitations.created_at IS '创建时间';
COMMENT ON COLUMN group_invitations.processed_at IS '处理时间';

COMMIT;