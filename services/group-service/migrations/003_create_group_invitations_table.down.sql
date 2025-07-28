-- 群组服务数据库迁移回滚 - 删除群组邀请表
-- 迁移版本: 003
-- 创建时间: 2025-07-28

BEGIN;

DROP TABLE IF EXISTS group_invitations CASCADE;

COMMIT;