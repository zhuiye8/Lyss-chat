-- 认证服务数据库迁移回滚 - 删除JWT令牌表
-- 迁移版本: 001
-- 创建时间: 2025-07-28

BEGIN;

-- 删除表（会自动删除索引和约束）
DROP TABLE IF EXISTS jwt_tokens CASCADE;

COMMIT;