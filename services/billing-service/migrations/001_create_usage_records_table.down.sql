-- 计费服务数据库迁移回滚 - 删除使用记录表
-- 迁版本: 001
-- 创建时间: 2025-07-28

BEGIN;

DROP TABLE IF EXISTS usage_records CASCADE;

COMMIT;