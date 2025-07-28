-- 网关服务数据库迁移回滚 - 删除API请求日志表
-- 迁移版本: 001
-- 创建时间: 2025-07-28

BEGIN;

DROP TABLE IF EXISTS api_request_logs CASCADE;

COMMIT;