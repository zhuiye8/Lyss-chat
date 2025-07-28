-- 共享服务数据库迁移 - 创建通知消息表
-- 迁移版本: 002
-- 创建时间: 2025-07-28
-- 数据库: lyss_platform

BEGIN;

-- 通知消息表
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT,
    group_id BIGINT,
    title VARCHAR(200) NOT NULL,
    content TEXT,
    type SMALLINT DEFAULT 1,           -- 1:系统 2:配额 3:安全 4:业务
    priority SMALLINT DEFAULT 1,       -- 1:普通 2:重要 3:紧急
    is_read BOOLEAN DEFAULT FALSE,
    read_at TIMESTAMP,
    expires_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 索引设计
CREATE INDEX idx_notifications_user_id ON notifications(user_id);
CREATE INDEX idx_notifications_is_read ON notifications(is_read);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- 添加表注释
COMMENT ON TABLE notifications IS '系统通知消息表';
COMMENT ON COLUMN notifications.id IS '通知记录唯一标识符';
COMMENT ON COLUMN notifications.user_id IS '接收用户ID';
COMMENT ON COLUMN notifications.group_id IS '关联群组ID';
COMMENT ON COLUMN notifications.title IS '通知标题';
COMMENT ON COLUMN notifications.content IS '通知内容';
COMMENT ON COLUMN notifications.type IS '通知类型：1-系统，2-配额，3-安全，4-业务';
COMMENT ON COLUMN notifications.priority IS '优先级：1-普通，2-重要，3-紧急';
COMMENT ON COLUMN notifications.is_read IS '是否已读';
COMMENT ON COLUMN notifications.read_at IS '阅读时间';
COMMENT ON COLUMN notifications.expires_at IS '过期时间';
COMMENT ON COLUMN notifications.created_at IS '创建时间';

COMMIT;