-- ====================================================================
-- 为用户分配预约相关角色
-- ====================================================================

-- 1. 查看当前用户
SELECT id, username, name FROM users ORDER BY id LIMIT 20;

-- 2. 查看当前角色分配情况
SELECT
  u.id,
  u.username,
  u.name,
  r.name as role_name
FROM users u
LEFT JOIN user_roles ur ON u.id = ur.user_id
LEFT JOIN roles r ON ur.role_id = r.id
WHERE r.id IN (10, 11, 12, 2, 4)
ORDER BY u.id, r.id;

-- 3. 分配角色给用户（请根据实际用户ID修改）

-- 示例：给用户ID=2分配"预约管理员"角色
-- DELETE FROM user_roles WHERE user_id = 2 AND role_id = 10;
-- INSERT INTO user_roles (user_id, role_id, created_at) VALUES (2, 10, CURRENT_TIMESTAMP);

-- 示例：给用户ID=3分配"施工员"角色
-- DELETE FROM user_roles WHERE user_id = 3 AND role_id = 11;
-- INSERT INTO user_roles (user_id, role_id, created_at) VALUES (3, 11, CURRENT_TIMESTAMP);

-- 示例：给用户ID=4分配"作业人员"角色
-- DELETE FROM user_roles WHERE user_id = 4 AND role_id = 12;
-- INSERT INTO user_roles (user_id, role_id, created_at) VALUES (4, 12, CURRENT_TIMESTAMP);

-- 4. 快速设置示例（假设用户ID 2,3,4）
-- 取消上面的注释，根据实际需要执行：

-- 设置用户2为预约管理员（客户）
DELETE FROM user_roles WHERE user_id = 2 AND role_id = 10;
INSERT INTO user_roles (user_id, role_id, created_at) VALUES (2, 10, CURRENT_TIMESTAMP);

-- 设置用户3为施工员
DELETE FROM user_roles WHERE user_id = 3 AND role_id = 11;
INSERT INTO user_roles (user_id, role_id, created_at) VALUES (3, 11, CURRENT_TIMESTAMP);

-- 设置用户4为作业人员
DELETE FROM user_roles WHERE user_id = 4 AND role_id = 12;
INSERT INTO user_roles (user_id, role_id, created_at) VALUES (4, 12, CURRENT_TIMESTAMP);

-- 5. 验证分配结果
SELECT '=== 角色分配结果 ===' AS info;
SELECT
  u.id,
  u.username,
  u.name,
  r.name as role_name,
  r.permissions
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE r.id IN (10, 11, 12, 2, 4)
ORDER BY u.id, r.id;
