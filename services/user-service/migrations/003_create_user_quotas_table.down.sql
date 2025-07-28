-- 用户服务数据库迁移回滚 - 删除用户配额表
-- 迁移版本: 003
-- 创建时间: 2025-07-28

BEGIN;

-- 删除表（会自动删除索引和约束）
DROP TABLE IF EXISTS user_quotas CASCADE;

COMMIT;