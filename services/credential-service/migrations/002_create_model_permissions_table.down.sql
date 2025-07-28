-- 凭证服务数据库迁移回滚 - 删除模型权限表
-- 迁移版本: 002
-- 创建时间: 2025-07-28

BEGIN;

DROP TABLE IF EXISTS model_permissions CASCADE;

COMMIT;