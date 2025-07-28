-- 用户服务数据库迁移回滚 - 删除用户基础表
-- 迁移版本: 001
-- 创建时间: 2025-07-28

BEGIN;

-- 删除表（会自动删除索引和约束）
DROP TABLE IF EXISTS users CASCADE;

COMMIT;