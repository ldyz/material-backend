-- 为用户分配预约相关角色
-- 用户ID 2=wqs, 3=julei, 4=libo

-- 分配预约管理员角色给用户2
DELETE FROM user_roles WHERE user_id = 2 AND role_id = 10;
INSERT INTO user_roles (user_id, role_id) VALUES (2, 10);

-- 分配施工员角色给用户3
DELETE FROM user_roles WHERE user_id = 3 AND role_id = 11;
INSERT INTO user_roles (user_id, role_id) VALUES (3, 11);

-- 分配作业人员角色给用户4
DELETE FROM user_roles WHERE user_id = 4 AND role_id = 12;
INSERT INTO user_roles (user_id, role_id) VALUES (4, 12);

-- 查看分配结果
SELECT
  u.id,
  u.username,
  r.name as role_name
FROM users u
JOIN user_roles ur ON u.id = ur.user_id
JOIN roles r ON ur.role_id = r.id
WHERE r.id IN (10, 11, 12, 2, 4)
ORDER BY u.id, r.id;
