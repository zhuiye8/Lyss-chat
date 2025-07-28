-- 共享服务数据库迁移回滚 - 删除系统配置表
-- 迁移版本: 001
-- 创建时间: 2025-07-28

BEGIN;

DROP TABLE IF EXISTS system_configs CASCADE;

COMMIT;