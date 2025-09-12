-- Teams Notification Service Database Schema
-- PostgreSQL 15+

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable JSONB for better JSON handling
CREATE EXTENSION IF NOT EXISTS "btree_gin";

-- =============================================
-- Core Tables
-- =============================================

-- Companies table
CREATE TABLE companies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    contact_email VARCHAR(255) NOT NULL,
    contact_phone VARCHAR(50),
    address TEXT,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    billing_enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(50) DEFAULT 'user' CHECK (role IN ('admin', 'user', 'viewer')),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    last_login_at TIMESTAMP WITH TIME ZONE,
    password_hash VARCHAR(255),
    api_key_hash VARCHAR(255),
    api_key_expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(company_id, email)
);

-- Projects table
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    key_name VARCHAR(100) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    daily_limit INTEGER DEFAULT 10000,
    monthly_limit INTEGER DEFAULT 300000,
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(company_id, key_name)
);

-- Destinations (Teams targets)
CREATE TABLE destinations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('person', 'channel', 'chatgroup')),
    teams_user_id VARCHAR(255),
    teams_team_id VARCHAR(255),
    teams_channel_id VARCHAR(255),
    teams_group_id VARCHAR(255),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    validation_status VARCHAR(20) DEFAULT 'pending' CHECK (validation_status IN ('pending', 'validated', 'failed')),
    last_validated_at TIMESTAMP WITH TIME ZONE,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_destination_type CHECK (
        (type = 'person' AND teams_user_id IS NOT NULL) OR
        (type = 'channel' AND teams_team_id IS NOT NULL AND teams_channel_id IS NOT NULL) OR
        (type = 'chatgroup' AND teams_group_id IS NOT NULL)
    )
);

-- Notifications table
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id),
    sender_id UUID NOT NULL REFERENCES users(id),
    message_type VARCHAR(20) NOT NULL CHECK (message_type IN ('text', 'file')),
    content TEXT NOT NULL,
    mentions JSONB DEFAULT '[]'::jsonb,
    attachment JSONB,
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high', 'urgent')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'sent', 'failed', 'cancelled')),
    error_message TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    sent_at TIMESTAMP WITH TIME ZONE
);

-- Notification Destinations (Many-to-many relationship)
CREATE TABLE notification_destinations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    destination_id UUID NOT NULL REFERENCES destinations(id),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'processing', 'sent', 'failed', 'cancelled')),
    error_message TEXT,
    teams_message_id VARCHAR(255),
    sent_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- =============================================
-- Billing and Usage Tables
-- =============================================

-- Billing Plans
CREATE TABLE billing_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL,
    description TEXT,
    price_per_notification DECIMAL(10, 6) NOT NULL DEFAULT 0.001,
    price_per_attachment DECIMAL(10, 6) NOT NULL DEFAULT 0.005,
    price_per_mention DECIMAL(10, 6) NOT NULL DEFAULT 0.0005,
    price_per_scheduled DECIMAL(10, 6) NOT NULL DEFAULT 0.002,
    daily_limit INTEGER,
    monthly_limit INTEGER,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Company Billing
CREATE TABLE company_billing (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    billing_plan_id UUID NOT NULL REFERENCES billing_plans(id),
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'suspended', 'cancelled')),
    billing_email VARCHAR(255),
    payment_method VARCHAR(50),
    currency VARCHAR(3) DEFAULT 'USD',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Usage Records
CREATE TABLE usage_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    project_id UUID REFERENCES projects(id),
    user_id UUID REFERENCES users(id),
    notification_id UUID REFERENCES notifications(id),
    record_type VARCHAR(20) NOT NULL CHECK (record_type IN ('notification', 'attachment', 'mention')),
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_price DECIMAL(10, 6) NOT NULL,
    total_cost DECIMAL(10, 6) NOT NULL,
    billing_period DATE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- Audit and Logging Tables
-- =============================================

-- Audit Logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID REFERENCES companies(id),
    user_id UUID REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    resource_type VARCHAR(50) NOT NULL,
    resource_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- System Logs
CREATE TABLE system_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    level VARCHAR(20) NOT NULL CHECK (level IN ('debug', 'info', 'warn', 'error', 'fatal')),
    service VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    context JSONB DEFAULT '{}'::jsonb,
    trace_id VARCHAR(255),
    span_id VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- Configuration Tables
-- =============================================

-- System Settings
CREATE TABLE system_settings (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    key VARCHAR(100) UNIQUE NOT NULL,
    value JSONB NOT NULL,
    description TEXT,
    is_public BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Feature Flags
CREATE TABLE feature_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    enabled BOOLEAN DEFAULT false,
    company_id UUID REFERENCES companies(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- Indexes
-- =============================================

-- Companies indexes
CREATE INDEX idx_companies_domain ON companies(domain);
CREATE INDEX idx_companies_status ON companies(status);

-- Users indexes
CREATE INDEX idx_users_company_id ON users(company_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_api_key_hash ON users(api_key_hash);

-- Projects indexes
CREATE INDEX idx_projects_company_id ON projects(company_id);
CREATE INDEX idx_projects_status ON projects(status);

-- Destinations indexes
CREATE INDEX idx_destinations_project_id ON destinations(project_id);
CREATE INDEX idx_destinations_type ON destinations(type);
CREATE INDEX idx_destinations_status ON destinations(status);
CREATE INDEX idx_destinations_teams_user_id ON destinations(teams_user_id);
CREATE INDEX idx_destinations_teams_team_id ON destinations(teams_team_id);
CREATE INDEX idx_destinations_teams_group_id ON destinations(teams_group_id);

-- Notifications indexes
CREATE INDEX idx_notifications_project_id ON notifications(project_id);
CREATE INDEX idx_notifications_sender_id ON notifications(sender_id);
CREATE INDEX idx_notifications_status ON notifications(status);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);
CREATE INDEX idx_notifications_sent_at ON notifications(sent_at);

-- Notification Destinations indexes
CREATE INDEX idx_notification_destinations_notification_id ON notification_destinations(notification_id);
CREATE INDEX idx_notification_destinations_destination_id ON notification_destinations(destination_id);
CREATE INDEX idx_notification_destinations_status ON notification_destinations(status);


-- Usage Records indexes
CREATE INDEX idx_usage_records_company_id ON usage_records(company_id);
CREATE INDEX idx_usage_records_billing_period ON usage_records(billing_period);
CREATE INDEX idx_usage_records_project_id ON usage_records(project_id);
CREATE INDEX idx_usage_records_created_at ON usage_records(created_at);

-- Audit Logs indexes
CREATE INDEX idx_audit_logs_company_id ON audit_logs(company_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- System Logs indexes
CREATE INDEX idx_system_logs_level ON system_logs(level);
CREATE INDEX idx_system_logs_service ON system_logs(service);
CREATE INDEX idx_system_logs_created_at ON system_logs(created_at);
CREATE INDEX idx_system_logs_trace_id ON system_logs(trace_id);

-- =============================================
-- Triggers for Updated At
-- =============================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at triggers to all relevant tables
CREATE TRIGGER update_companies_updated_at BEFORE UPDATE ON companies FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_destinations_updated_at BEFORE UPDATE ON destinations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_notification_destinations_updated_at BEFORE UPDATE ON notification_destinations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_company_billing_updated_at BEFORE UPDATE ON company_billing FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_system_settings_updated_at BEFORE UPDATE ON system_settings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_feature_flags_updated_at BEFORE UPDATE ON feature_flags FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================
-- Views for Common Queries
-- =============================================

-- Notification Summary View
CREATE VIEW notification_summary AS
SELECT 
    n.id,
    n.project_id,
    p.key_name,
    p.company_id,
    c.name as company_name,
    u.email as sender_email,
    n.message_type,
    n.content,
    n.priority,
    n.status,
    n.created_at,
    n.sent_at,
    COUNT(nd.id) as destination_count,
    COUNT(CASE WHEN nd.status = 'sent' THEN 1 END) as sent_count,
    COUNT(CASE WHEN nd.status = 'failed' THEN 1 END) as failed_count
FROM notifications n
JOIN projects p ON n.project_id = p.id
JOIN companies c ON p.company_id = c.id
JOIN users u ON n.sender_id = u.id
LEFT JOIN notification_destinations nd ON n.id = nd.notification_id
GROUP BY n.id, p.key_name, p.company_id, c.name, u.email, n.message_type, n.content, n.priority, n.status, n.created_at, n.sent_at;

-- Usage Summary View
CREATE VIEW usage_summary AS
SELECT 
    ur.company_id,
    c.name as company_name,
    ur.billing_period,
    ur.record_type,
    SUM(ur.quantity) as total_quantity,
    SUM(ur.total_cost) as total_cost,
    COUNT(DISTINCT ur.project_id) as projects_used
FROM usage_records ur
JOIN companies c ON ur.company_id = c.id
GROUP BY ur.company_id, c.name, ur.billing_period, ur.record_type;

-- =============================================
-- Initial Data
-- =============================================

-- Insert default billing plan
INSERT INTO billing_plans (name, description, price_per_notification, price_per_attachment, price_per_mention, price_per_scheduled, daily_limit, monthly_limit) VALUES
('Standard', 'Standard billing plan for Teams notifications', 0.001, 0.005, 0.0005, 0.002, 10000, 300000);

-- Insert default system settings
INSERT INTO system_settings (key, value, description, is_public) VALUES
('max_message_length', '4000', 'Maximum message length in characters', true),
('max_attachment_size', '10485760', 'Maximum attachment size in bytes (10MB)', true),
('allowed_file_types', '["pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "jpg", "jpeg", "png", "gif"]', 'Allowed file types for attachments', true),
('rate_limit_per_minute', '100', 'Rate limit per minute per user', true),
('retry_attempts', '3', 'Maximum retry attempts for failed notifications', false),
('cleanup_retention_days', '90', 'Number of days to retain old notifications', false);

-- Insert default feature flags
INSERT INTO feature_flags (name, description, enabled) VALUES
('file_attachments', 'Enable file attachment feature', true),
('mentions', 'Enable @mention feature', true),
('webhooks', 'Enable webhook notifications', false),
('analytics', 'Enable analytics dashboard', true),
('billing', 'Enable billing and usage tracking', true),
('admin_panel', 'Enable admin panel', true),
('audit_logs', 'Enable audit logging', true),
('maintenance_mode', 'Enable maintenance mode', false);
