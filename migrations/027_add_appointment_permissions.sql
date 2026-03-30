-- 添加施工预约相关权限到角色
-- 运行时间: 2026-03-29

-- 更新admin角色，添加预约相关权限
UPDATE roles
SET permissions = permissions || ',appointment_view,appointment_create,appointment_edit,appointment_delete,appointment_submit,appointment_approve,appointment_assign,appointment_cancel'
WHERE name = 'admin'
AND permissions NOT LIKE '%appointment_view%';

-- 更新项目经理角色
UPDATE roles
SET permissions = permissions || ',appointment_view,appointment_create,appointment_edit,appointment_submit,appointment_approve,appointment_assign'
WHERE name = '项目经理'
AND permissions NOT LIKE '%appointment_view%';

-- 更新施工员角色
UPDATE roles
SET permissions = permissions || ',appointment_view,appointment_approve,appointment_assign'
WHERE name = '施工员'
AND permissions NOT LIKE '%appointment_view%';

-- 更新预约管理员角色
UPDATE roles
SET permissions = permissions || ',appointment_view,appointment_create,appointment_edit,appointment_delete,appointment_submit,appointment_approve,appointment_assign,appointment_cancel'
WHERE name = '预约管理员'
AND permissions NOT LIKE '%appointment_view%';

-- 更新材料员角色
UPDATE roles
SET permissions = permissions || ',appointment_view,appointment_create'
WHERE name = '材料员'
AND permissions NOT LIKE '%appointment_view%';
