-- 共享服务数据库迁移回滚 - 删除通知消息表  
-- 迁移版本: 002
-- 创建时间: 2025-07-28

BEGIN;

DROP TABLE IF EXISTS notifications CASCADE;

COMMIT;