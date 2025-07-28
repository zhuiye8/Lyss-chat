-- 计费服务数据库迁移回滚 - 删除账单汇总表
-- 迁移版本: 003
-- 创建时间: 2025-07-28

BEGIN;

DROP TABLE IF EXISTS billing_summaries CASCADE;

COMMIT;