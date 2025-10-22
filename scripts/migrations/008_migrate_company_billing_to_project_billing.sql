-- Migration: Migrate company_billing to project_billing
-- Date: 2025-10-21
-- Description: Remove company_billing table and create project_billing table

-- Step 1: Create project_billing table if it doesn't exist
CREATE TABLE IF NOT EXISTS project_billing (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE UNIQUE,
    billing_plan_id UUID REFERENCES billing_plans(id) ON DELETE SET NULL,
    billing_status VARCHAR(50) DEFAULT 'active', -- 'active', 'suspended', 'cancelled'
    payment_method VARCHAR(50), -- 'credit_card', 'bank_transfer', 'invoice'
    billing_cycle VARCHAR(20) DEFAULT 'monthly', -- 'monthly', 'yearly'
    next_billing_date DATE,
    total_usage_cost DECIMAL(10,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Step 2: Create indexes for project_billing
CREATE INDEX IF NOT EXISTS idx_project_billing_project_id ON project_billing(project_id);
CREATE INDEX IF NOT EXISTS idx_project_billing_billing_plan_id ON project_billing(billing_plan_id);

-- Step 3: Create trigger for project_billing
CREATE TRIGGER update_project_billing_updated_at 
    BEFORE UPDATE ON project_billing 
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Step 4: Migrate data from company_billing to project_billing
-- For each company, find its first project and migrate billing data
INSERT INTO project_billing (
    id,
    project_id,
    billing_plan_id,
    billing_status,
    payment_method,
    billing_cycle,
    next_billing_date,
    total_usage_cost,
    created_at,
    updated_at
)
SELECT 
    cb.id,
    p.id as project_id,
    cb.billing_plan_id,
    cb.billing_status,
    cb.payment_method,
    cb.billing_cycle,
    cb.next_billing_date,
    cb.total_usage_cost,
    cb.created_at,
    cb.updated_at
FROM company_billing cb
JOIN companies c ON cb.company_id = c.id
JOIN projects p ON c.id = p.company_id
WHERE p.id = (
    SELECT id FROM projects 
    WHERE company_id = cb.company_id 
    ORDER BY created_at ASC 
    LIMIT 1
)
ON CONFLICT (project_id) DO NOTHING;

-- Step 5: Drop company_billing table and its dependencies
DROP TRIGGER IF EXISTS update_company_billing_updated_at ON company_billing;
DROP INDEX IF EXISTS idx_company_billing_company_id;
DROP INDEX IF EXISTS idx_company_billing_billing_plan_id;
DROP TABLE IF EXISTS company_billing;

-- Step 6: Update usage_summary view to use project_billing
DROP VIEW IF EXISTS usage_summary;
CREATE VIEW usage_summary AS
SELECT 
    ur.project_id,
    p.notify_key,
    p.company_id,
    c.name as company_name,
    ur.record_type,
    SUM(ur.quantity) as total_quantity,
    SUM(ur.total_cost) as total_cost,
    DATE(ur.created_at) as usage_date
FROM usage_records ur
JOIN projects p ON ur.project_id = p.id
JOIN companies c ON p.company_id = c.id
GROUP BY ur.project_id, p.notify_key, p.company_id, c.name, ur.record_type, DATE(ur.created_at);
