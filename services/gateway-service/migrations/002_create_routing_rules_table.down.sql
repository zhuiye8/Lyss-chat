-- 网关服务数据库迁移回滚 - 删除路由配置表
-- 迁移版本: 002
-- 创建时间: 2025-07-28

BEGIN;

DROP TABLE IF EXISTS routing_rules CASCADE;

COMMIT;