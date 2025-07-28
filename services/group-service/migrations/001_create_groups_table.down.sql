-- 群组服务数据库迁移回滚 - 删除群组基础表
-- 迁移版本: 001
-- 创建时间: 2025-07-28

BEGIN;

DROP TABLE IF EXISTS groups CASCADE;

COMMIT;