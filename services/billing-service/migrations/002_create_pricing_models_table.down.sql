-- 计费服务数据库迁移回滚 - 删除计费规则表
-- 迁移版本: 002
-- 创建时间: 2025-07-28

BEGIN;

DROP TABLE IF EXISTS pricing_models CASCADE;

COMMIT;