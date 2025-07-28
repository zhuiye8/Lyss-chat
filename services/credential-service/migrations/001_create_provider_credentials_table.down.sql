-- 凭证服务数据库迁移回滚 - 删除供应商凭证表
-- 迁移版本: 001
-- 创建时间: 2025-07-28

BEGIN;

DROP TABLE IF EXISTS provider_credentials CASCADE;

COMMIT;