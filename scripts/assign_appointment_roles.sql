-- ====================================================================
-- 为测试用户分配预约相关角色
-- ====================================================================

-- 注意：运行此脚本前，请根据实际情况调整用户ID

-- 1. 查看当前用户列表
SELECT id, username, name FROM users ORDER BY id LIMIT 20;

-- 2. 为用户分配角色（请根据实际用户ID修改）

-- 示例：为用户分配预约管理员角色（客户）
-- DELETE FROM user_roles WHERE user_id = 2 AND role_id = 10;
-- INSERT INTO user_roles (user_id, role_id, created_at) VALUES (2, 10, NOW());

-- 示例：为用户分配施工员角色
-- DELETE FROM user_roles WHERE user_id = 3 AND role_id = 11;
-- INSERT INTO user_roles (user_id, role_id, created_at) VALUES (3, 11, NOW());

-- 示例：为用户分配作业人员角色
-- DELETE FROM user_roles WHERE user_id = 4 AND role_id = 12;
-- INSERT INTO user_roles (user_id, role_id, created_at) VALUES (4, 12, NOW());

-- 3. 为现有管理员添加预约审批权限
-- 管理员(user_id=1)已经拥有system_config权限，自动拥有审批权限

-- 4. 查看用户的角色分配情况
SELECT
  u.id,
  u.username,
  u.name,
  r.name as role_name
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE r.id IN (10, 11, 12)
ORDER BY u.id;

-- 5. 查看所有用户的权限情况（用于验证）
SELECT
  u.id,
  u.username,
  p.name as permission,
  p.module
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN role_permissions rp ON ur.role_id = rp.role_id
LEFT JOIN permissions p ON rp.permission_id = p.id
WHERE p.module = 'appointment'
ORDER BY u.id, p.name;
