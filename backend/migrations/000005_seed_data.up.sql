-- 插入默认租户
INSERT INTO tenants (id, name, domain, status, created_at, updated_at)
VALUES 
    ('11111111-1111-1111-1111-111111111111', 'Default Tenant', 'default.lyss-chat.com', 'active', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 插入默认权限
INSERT INTO permissions (id, code, name, description, resource, action, created_at, updated_at)
VALUES
    ('20000000-0000-0000-0000-000000000001', 'users:read', '查看用户', '允许查看用户信息', 'users', 'read', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000002', 'users:create', '创建用户', '允许创建新用户', 'users', 'create', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000003', 'users:update', '更新用户', '允许更新用户信息', 'users', 'update', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000004', 'users:delete', '删除用户', '允许删除用户', 'users', 'delete', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000005', 'roles:read', '查看角色', '允许查看角色信息', 'roles', 'read', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000006', 'roles:create', '创建角色', '允许创建新角色', 'roles', 'create', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000007', 'roles:update', '更新角色', '允许更新角色信息', 'roles', 'update', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000008', 'roles:delete', '删除角色', '允许删除角色', 'roles', 'delete', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000009', 'models:read', '查看模型', '允许查看模型信息', 'models', 'read', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000010', 'models:create', '创建模型', '允许创建新模型', 'models', 'create', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000011', 'models:update', '更新模型', '允许更新模型信息', 'models', 'update', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000012', 'models:delete', '删除模型', '允许删除模型', 'models', 'delete', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000013', 'canvases:read', '查看画布', '允许查看画布信息', 'canvases', 'read', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000014', 'canvases:create', '创建画布', '允许创建新画布', 'canvases', 'create', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000015', 'canvases:update', '更新画布', '允许更新画布信息', 'canvases', 'update', NOW(), NOW()),
    ('20000000-0000-0000-0000-000000000016', 'canvases:delete', '删除画布', '允许删除画布', 'canvases', 'delete', NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 插入默认角色
INSERT INTO roles (id, tenant_id, name, description, is_system, created_at, updated_at)
VALUES
    ('30000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', 'Admin', '系统管理员，拥有所有权限', TRUE, NOW(), NOW()),
    ('30000000-0000-0000-0000-000000000002', '11111111-1111-1111-1111-111111111111', 'User', '普通用户，拥有基本权限', TRUE, NOW(), NOW())
ON CONFLICT DO NOTHING;

-- 为管理员角色分配所有权限
INSERT INTO role_permissions (id, role_id, permission_id, created_at)
SELECT 
    md5(random()::text || clock_timestamp()::text)::uuid, 
    '30000000-0000-0000-0000-000000000001', 
    id, 
    NOW()
FROM permissions
ON CONFLICT DO NOTHING;

-- 为普通用户角色分配基本权限
INSERT INTO role_permissions (id, role_id, permission_id, created_at)
VALUES
    (md5(random()::text || clock_timestamp()::text)::uuid, '30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000001', NOW()),
    (md5(random()::text || clock_timestamp()::text)::uuid, '30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000009', NOW()),
    (md5(random()::text || clock_timestamp()::text)::uuid, '30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000013', NOW()),
    (md5(random()::text || clock_timestamp()::text)::uuid, '30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000014', NOW()),
    (md5(random()::text || clock_timestamp()::text)::uuid, '30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000015', NOW()),
    (md5(random()::text || clock_timestamp()::text)::uuid, '30000000-0000-0000-0000-000000000002', '20000000-0000-0000-0000-000000000016', NOW())
ON CONFLICT DO NOTHING;
