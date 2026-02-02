-- Verification Script for Material to Plan Migration
-- Run this after migration to verify data integrity

-- 1. Check created plans
SELECT
    p.id,
    p.plan_no,
    p.plan_name,
    p.project_id,
    pr.name as project_name,
    p.status,
    p.items_count,
    p.total_budget / 100.0 as total_budget,
    p.creator_name,
    p.created_at
FROM material_plans p
LEFT JOIN projects pr ON p.project_id = pr.id
ORDER BY p.id;

-- 2. Check plan items
SELECT
    pi.id,
    pi.plan_id,
    p.plan_no,
    pi.material_id,
    m.name as material_name,
    pi.material_name,
    pi.specification,
    pi.unit,
    pi.planned_quantity,
    pi.arrived_quantity,
    pi.issued_quantity,
    pi.unit_price / 100.0 as unit_price,
    pi.total_price / 100.0 as total_price,
    pi.status
FROM material_plan_items pi
LEFT JOIN material_plans p ON pi.plan_id = p.id
LEFT JOIN materials m ON pi.material_id = m.id
ORDER BY pi.plan_id, pi.sort_order;

-- 3. Verify plan-item integrity
SELECT
    p.id,
    p.plan_no,
    p.plan_name,
    COUNT(pi.id) as item_count
FROM material_plans p
LEFT JOIN material_plan_items pi ON p.id = pi.plan_id
GROUP BY p.id, p.plan_no, p.plan_name
HAVING COUNT(pi.id) = 0; -- Plans without items

-- 4. Check items without materials
SELECT
    pi.id,
    pi.plan_id,
    pi.material_id,
    pi.material_name
FROM material_plan_items pi
LEFT JOIN materials m ON pi.material_id = m.id
WHERE pi.material_id IS NOT NULL AND m.id IS NULL;

-- 5. Summary statistics
SELECT
    COUNT(DISTINCT p.id) as total_plans,
    COUNT(DISTINCT pi.id) as total_items,
    COUNT(DISTINCT p.project_id) as projects_with_plans,
    SUM(p.total_budget) / 100.0 as total_budget_all_plans
FROM material_plans p
LEFT JOIN material_plan_items pi ON p.id = pi.plan_id;

-- 6. Status distribution
SELECT
    status,
    COUNT(*) as count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
FROM material_plans
GROUP BY status
ORDER BY count DESC;

-- 7. Items status distribution
SELECT
    status,
    COUNT(*) as count,
    ROUND(COUNT(*) * 100.0 / SUM(COUNT(*)) OVER(), 2) as percentage
FROM material_plan_items
GROUP BY status
ORDER BY count DESC;
