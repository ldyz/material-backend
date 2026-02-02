-- Migration: Create material_plans and material_plan_items tables
-- Description: Creates the core tables for material plan management

-- Create material_plans table
CREATE TABLE IF NOT EXISTS material_plans (
    id BIGSERIAL PRIMARY KEY,
    plan_no VARCHAR(50) UNIQUE NOT NULL,
    plan_name VARCHAR(200) NOT NULL,
    project_id BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    plan_type VARCHAR(20) DEFAULT 'procurement',
    status VARCHAR(20) DEFAULT 'draft',
    priority VARCHAR(20) DEFAULT 'normal',
    planned_start_date TIMESTAMP,
    planned_end_date TIMESTAMP,
    total_budget INTEGER DEFAULT 0,
    actual_cost INTEGER DEFAULT 0,
    description TEXT,
    remark TEXT,
    workflow_instance_id BIGINT,
    creator_id INTEGER NOT NULL,
    creator_name VARCHAR(100) NOT NULL,
    approver_id BIGINT,
    approver_name VARCHAR(100),
    approved_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for material_plans
CREATE INDEX IF NOT EXISTS idx_material_plans_project_id ON material_plans(project_id);
CREATE INDEX IF NOT EXISTS idx_material_plans_status ON material_plans(status);
CREATE INDEX IF NOT EXISTS idx_material_plans_workflow_instance ON material_plans(workflow_instance_id);
CREATE INDEX IF NOT EXISTS idx_material_plans_creator_id ON material_plans(creator_id);

-- Create material_plan_items table
CREATE TABLE IF NOT EXISTS material_plan_items (
    id BIGSERIAL PRIMARY KEY,
    plan_id BIGINT NOT NULL REFERENCES material_plans(id) ON DELETE CASCADE,
    material_id BIGINT REFERENCES materials(id) ON DELETE SET NULL,
    material_name VARCHAR(200) NOT NULL,
    material_code VARCHAR(50),
    specification VARCHAR(200),
    category VARCHAR(100),
    unit VARCHAR(20),
    planned_quantity INTEGER DEFAULT 0,
    arrived_quantity INTEGER DEFAULT 0,
    issued_quantity INTEGER DEFAULT 0,
    unit_price INTEGER DEFAULT 0,
    total_price INTEGER DEFAULT 0,
    required_date TIMESTAMP,
    priority VARCHAR(20) DEFAULT 'normal',
    status VARCHAR(20) DEFAULT 'pending',
    remark TEXT,
    sort_order INTEGER DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for material_plan_items
CREATE INDEX IF NOT EXISTS idx_material_plan_items_plan_id ON material_plan_items(plan_id);
CREATE INDEX IF NOT EXISTS idx_material_plan_items_material_id ON material_plan_items(material_id);

-- Add comments for documentation
COMMENT ON TABLE material_plans IS '物资计划表';
COMMENT ON COLUMN material_plans.plan_no IS '计划编号';
COMMENT ON COLUMN material_plans.plan_name IS '计划名称';
COMMENT ON COLUMN material_plans.project_id IS '关联项目ID';
COMMENT ON COLUMN material_plans.plan_type IS '计划类型: procurement(采购), usage(使用), mixed(混合)';
COMMENT ON COLUMN material_plans.status IS '状态: draft, pending, approved, rejected, active, completed, cancelled';
COMMENT ON COLUMN material_plans.priority IS '优先级: low, normal, high, urgent';
COMMENT ON COLUMN material_plans.total_budget IS '总预算(分)';
COMMENT ON COLUMN material_plans.actual_cost IS '实际成本(分)';

COMMENT ON TABLE material_plan_items IS '物资计划明细表';
COMMENT ON COLUMN material_plan_items.plan_id IS '关联计划ID';
COMMENT ON COLUMN material_plan_items.material_id IS '关联物资ID';
COMMENT ON COLUMN material_plan_items.planned_quantity IS '计划数量';
COMMENT ON COLUMN material_plan_items.arrived_quantity IS '已到货数量';
COMMENT ON COLUMN material_plan_items.issued_quantity IS '已发放数量';
COMMENT ON COLUMN material_plan_items.unit_price IS '单价(分)';
COMMENT ON COLUMN material_plan_items.total_price IS '总价(分)';
COMMENT ON COLUMN material_plan_items.status IS '状态: pending, partial, completed, cancelled';
