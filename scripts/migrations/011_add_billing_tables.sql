-- 添加計費相關表
-- 計費方案表
CREATE TABLE IF NOT EXISTS billing_plans (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price_per_notification DECIMAL(10,4) NOT NULL DEFAULT 0.01,
    price_per_month DECIMAL(10,2) NOT NULL DEFAULT 0.00,
    max_notifications_per_month INTEGER,
    max_projects INTEGER,
    max_users INTEGER,
    features JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 使用記錄表
CREATE TABLE IF NOT EXISTS usage_records (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    record_type VARCHAR(50) NOT NULL, -- 'notification', 'api_call', 'storage'
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_cost DECIMAL(10,4) NOT NULL DEFAULT 0.01,
    total_cost DECIMAL(10,4) NOT NULL DEFAULT 0.01,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 公司計費表
CREATE TABLE IF NOT EXISTS company_billing (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    company_id UUID REFERENCES companies(id) ON DELETE CASCADE UNIQUE,
    billing_plan_id UUID REFERENCES billing_plans(id) ON DELETE SET NULL,
    billing_status VARCHAR(50) DEFAULT 'active', -- 'active', 'suspended', 'cancelled'
    payment_method VARCHAR(50), -- 'credit_card', 'bank_transfer', 'invoice'
    billing_cycle VARCHAR(20) DEFAULT 'monthly', -- 'monthly', 'yearly'
    next_billing_date DATE,
    total_usage_cost DECIMAL(10,2) DEFAULT 0.00,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 創建索引
CREATE INDEX IF NOT EXISTS idx_usage_records_company_id ON usage_records(company_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_project_id ON usage_records(project_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_user_id ON usage_records(user_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_record_type ON usage_records(record_type);
CREATE INDEX IF NOT EXISTS idx_usage_records_created_at ON usage_records(created_at);
CREATE INDEX IF NOT EXISTS idx_company_billing_company_id ON company_billing(company_id);
CREATE INDEX IF NOT EXISTS idx_company_billing_billing_plan_id ON company_billing(billing_plan_id);
