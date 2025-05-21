-- 删除索引
DROP INDEX IF EXISTS idx_users_tenant_id;
DROP INDEX IF EXISTS idx_users_email;
DROP INDEX IF EXISTS idx_users_status;

-- 删除表
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;
