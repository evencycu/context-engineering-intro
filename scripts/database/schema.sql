-- Teams Bot Platform Database Schema
-- PostgreSQL 15+
-- 支援多個 Teams Bot 和第三方 Bot 的中台服務

-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Enable JSONB for better JSON handling
CREATE EXTENSION IF NOT EXISTS "btree_gin";

-- =============================================
-- Core Platform Tables
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

-- Projects table (notification keys)
CREATE TABLE projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID NOT NULL REFERENCES companies(id) ON DELETE CASCADE,
    notify_key VARCHAR(100) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    daily_limit INTEGER DEFAULT 10000,
    monthly_limit INTEGER DEFAULT 300000,
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high')),
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(company_id, notify_key)
);

-- =============================================
-- Bot Management Tables
-- =============================================

-- Bot Types
CREATE TYPE bot_type AS ENUM ('platform', 'third_party');

-- Bot Status
CREATE TYPE bot_status AS ENUM ('active', 'inactive', 'suspended', 'maintenance');

-- Unified Teams Bots (platform + third_party)
CREATE TABLE teams_bots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    type VARCHAR(20) NOT NULL CHECK (type IN ('platform','third_party')),
    company_id UUID,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    app_id VARCHAR(255) NOT NULL UNIQUE,
    app_password_hash VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255),
    status bot_status DEFAULT 'active',
    webhook_url VARCHAR(500),
    capabilities JSONB DEFAULT '{}'::jsonb,
    rate_limit_per_minute INTEGER DEFAULT 60,
    max_concurrent_requests INTEGER DEFAULT 10,
    api_endpoint VARCHAR(500),
    api_key_hash VARCHAR(255),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    created_by UUID
);

-- Bot Installations (Bot 安裝記錄)
CREATE TABLE bot_installations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bot_id UUID NOT NULL,
    bot_type bot_type NOT NULL,
    teams_tenant_id VARCHAR(255) NOT NULL,
    -- Conversation type and identification
    conversation_type VARCHAR(20) NOT NULL CHECK (conversation_type IN ('personal', 'channel', 'groupChat')),
    conversation_id VARCHAR(500) NOT NULL, -- Teams conversation ID for sending messages
    service_url VARCHAR(500), -- Bot Framework service URL
    -- Bot and user identification
    recipient_id VARCHAR(255) NOT NULL, -- Bot member ID (usually 28:app_id)
    recipient_name VARCHAR(255), -- Bot display name
    from_id VARCHAR(255), -- Source member ID from Teams event
    from_name VARCHAR(255), -- User display name
    from_aad_object_id VARCHAR(255), -- User's AAD Object ID
    email VARCHAR(255), -- User email address
    description_name VARCHAR(255), -- User description or display name
    -- Status and lifecycle
    installation_status VARCHAR(20) DEFAULT 'active' CHECK (installation_status IN ('active', 'inactive', 'uninstalled', 'stale')),
    installed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    uninstalled_at TIMESTAMP WITH TIME ZONE,
    last_activity_at TIMESTAMP WITH TIME ZONE, -- Last successful message sent/received
    -- Metadata and configuration
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    -- Unique constraint for preventing duplicates (per bot, tenant, conversation)
    UNIQUE(bot_id, bot_type, teams_tenant_id, conversation_id)
);

-- =============================================
-- Destination Management
-- =============================================

-- Destinations (Teams targets) - 支援多個目標的組合
CREATE TABLE destinations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    teams_tenant_id VARCHAR(255) NOT NULL,
    targets JSONB NOT NULL DEFAULT '[]'::jsonb, -- Array of Teams targets
    bot_id UUID,
    status VARCHAR(20) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
    validation_status VARCHAR(20) DEFAULT 'pending' CHECK (validation_status IN ('pending', 'validated', 'failed')),
    last_validated_at TIMESTAMP WITH TIME ZONE,
    created_by UUID NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_targets_array CHECK (jsonb_array_length(targets) > 0)
);

-- =============================================
-- Notification System
-- =============================================

-- Notifications table
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id),
    sender_id UUID,
    message_type VARCHAR(20) NOT NULL CHECK (message_type IN ('text', 'file', 'adaptive_card')),
    content TEXT NOT NULL,
    mentions JSONB DEFAULT '[]'::jsonb,
    attachments JSONB DEFAULT '[]'::jsonb,
    priority VARCHAR(20) DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high')),
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'enqueued', 'processing', 'sent', 'failed', 'cancelled')),
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
    bot_id UUID,
    bot_type bot_type,
    status VARCHAR(20) DEFAULT 'pending' CHECK (status IN ('pending', 'enqueued', 'processing', 'sent', 'failed', 'cancelled')),
    error_message TEXT,
    teams_message_id VARCHAR(255),
    sent_at TIMESTAMP WITH TIME ZONE,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 5,
    -- Async actor fields
    conversation_id VARCHAR(500),
    next_retry_at TIMESTAMP WITH TIME ZONE,
    first_attempt_at TIMESTAMP WITH TIME ZONE,
    last_attempt_at TIMESTAMP WITH TIME ZONE,
    failure_reason TEXT,
    retry_after INT, -- Retry-After from 429 response (seconds)
    actor_id VARCHAR(255), -- ID of the actor handling this notification
    priority VARCHAR(20) NOT NULL DEFAULT 'normal' CHECK (priority IN ('low', 'normal', 'high')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- Bot Routing and Load Balancing
-- =============================================

-- Bot Routing Rules
CREATE TABLE bot_routing_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    priority INTEGER DEFAULT 100,
    conditions JSONB NOT NULL, -- 路由條件
    target_bot_id UUID,
    target_bot_type bot_type,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Bot Health Status
CREATE TABLE bot_health_status (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bot_id UUID NOT NULL,
    bot_type bot_type NOT NULL,
    status VARCHAR(20) NOT NULL CHECK (status IN ('healthy', 'unhealthy', 'degraded')),
    last_check_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    response_time_ms INTEGER,
    error_count INTEGER DEFAULT 0,
    success_count INTEGER DEFAULT 0,
    error_message TEXT,
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- Billing and Usage Tables
-- =============================================

-- Billing Plans
CREATE TABLE billing_plans (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
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

-- Project Billing
CREATE TABLE project_billing (
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

-- Usage Records
CREATE TABLE usage_records (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    record_type VARCHAR(50) NOT NULL, -- 'notification', 'api_call', 'storage'
    quantity INTEGER NOT NULL DEFAULT 1,
    unit_cost DECIMAL(10,4) NOT NULL DEFAULT 0.01,
    total_cost DECIMAL(10,4) NOT NULL DEFAULT 0.01,
    metadata JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- File Management Tables
-- =============================================

-- Files table for file attachments
CREATE TABLE files (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    content_type VARCHAR(100) NOT NULL,
    file_key VARCHAR(255) NOT NULL UNIQUE,
    file_url TEXT,
    description TEXT,
    tags TEXT[] DEFAULT '{}',
    is_public BOOLEAN DEFAULT false,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);


-- =============================================
-- Audit and Logging Tables
-- =============================================

-- Audit Logs
CREATE TABLE audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    company_id UUID REFERENCES companies(id),
    user_id UUID,
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
    bot_id UUID,
    bot_type bot_type,
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
CREATE INDEX idx_companies_status ON companies(status);

-- Users indexes
CREATE INDEX idx_users_company_id ON users(company_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_api_key_hash ON users(api_key_hash);

-- Projects indexes
CREATE INDEX idx_projects_company_id ON projects(company_id);
CREATE INDEX idx_projects_status ON projects(status);

-- Teams Bots indexes
CREATE INDEX idx_teams_bots_type ON teams_bots(type);
CREATE INDEX idx_teams_bots_app_id ON teams_bots(app_id);
CREATE INDEX idx_teams_bots_status ON teams_bots(status);
CREATE INDEX idx_teams_bots_tenant ON teams_bots(tenant_id);
CREATE INDEX idx_teams_bots_company ON teams_bots(company_id);

-- Bot Installations indexes
CREATE INDEX idx_bot_installations_bot_id ON bot_installations(bot_id);
CREATE INDEX idx_bot_installations_bot_type ON bot_installations(bot_type);
CREATE INDEX idx_bot_installations_teams_tenant_id ON bot_installations(teams_tenant_id);
CREATE INDEX idx_bot_installations_conversation_type ON bot_installations(conversation_type);
CREATE INDEX idx_bot_installations_conversation_id ON bot_installations(conversation_id);
CREATE INDEX idx_bot_installations_from_aad_object_id ON bot_installations(from_aad_object_id);
CREATE INDEX idx_bot_installations_status ON bot_installations(installation_status);
CREATE INDEX idx_bot_installations_bot_tenant ON bot_installations(bot_id, bot_type, teams_tenant_id);
-- removed legacy columns

-- Destinations indexes
CREATE INDEX idx_destinations_project_id ON destinations(project_id);
CREATE INDEX idx_destinations_status ON destinations(status);
CREATE INDEX idx_destinations_teams_tenant_id ON destinations(teams_tenant_id);
CREATE INDEX idx_destinations_bot_id ON destinations(bot_id);
-- removed: destination has no bot_type now
-- JSONB indexes for targets array
CREATE INDEX idx_destinations_targets_gin ON destinations USING GIN (targets);
CREATE INDEX idx_destinations_targets_type ON destinations USING GIN ((targets->'type'));
-- New fields per TeamsTarget: email, conversation_id, tenant_id
CREATE INDEX idx_destinations_targets_conversation_id ON destinations USING GIN ((targets->'conversation_id'));
CREATE INDEX idx_destinations_targets_email ON destinations USING GIN ((targets->'email'));
CREATE INDEX idx_destinations_targets_tenant_id ON destinations USING GIN ((targets->'tenant_id'));

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
CREATE INDEX idx_notification_destinations_bot_id ON notification_destinations(bot_id);

-- Bot Routing Rules indexes
CREATE INDEX idx_bot_routing_rules_priority ON bot_routing_rules(priority);
CREATE INDEX idx_bot_routing_rules_active ON bot_routing_rules(is_active);

-- Bot Health Status indexes
CREATE INDEX idx_bot_health_status_bot_id ON bot_health_status(bot_id);
CREATE INDEX idx_bot_health_status_bot_type ON bot_health_status(bot_type);
CREATE INDEX idx_bot_health_status_status ON bot_health_status(status);
CREATE INDEX idx_bot_health_status_last_check ON bot_health_status(last_check_at);

-- Usage Records indexes
CREATE INDEX idx_usage_records_project_id ON usage_records(project_id);
CREATE INDEX idx_usage_records_user_id ON usage_records(user_id);
CREATE INDEX idx_usage_records_record_type ON usage_records(record_type);
CREATE INDEX idx_usage_records_created_at ON usage_records(created_at);

-- Project Billing indexes
CREATE INDEX idx_project_billing_project_id ON project_billing(project_id);
CREATE INDEX idx_project_billing_billing_plan_id ON project_billing(billing_plan_id);

-- Files indexes
CREATE INDEX idx_files_project_id ON files(project_id);
CREATE INDEX idx_files_file_key ON files(file_key);
CREATE INDEX idx_files_is_public ON files(is_public);
CREATE INDEX idx_files_expires_at ON files(expires_at);
CREATE INDEX idx_files_created_at ON files(created_at);


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
CREATE INDEX idx_system_logs_bot_id ON system_logs(bot_id);

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
CREATE TRIGGER update_teams_bots_updated_at BEFORE UPDATE ON teams_bots FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_bot_installations_updated_at BEFORE UPDATE ON bot_installations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_destinations_updated_at BEFORE UPDATE ON destinations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_notifications_updated_at BEFORE UPDATE ON notifications FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_notification_destinations_updated_at BEFORE UPDATE ON notification_destinations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_bot_routing_rules_updated_at BEFORE UPDATE ON bot_routing_rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
-- (removed third party bot api keys trigger)
CREATE TRIGGER update_project_billing_updated_at BEFORE UPDATE ON project_billing FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_system_settings_updated_at BEFORE UPDATE ON system_settings FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_feature_flags_updated_at BEFORE UPDATE ON feature_flags FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_files_updated_at BEFORE UPDATE ON files FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================
-- Views for Common Queries
-- =============================================

-- Notification Summary View
CREATE VIEW notification_summary AS
SELECT 
    n.id,
    n.project_id,
    p.notify_key,
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
LEFT JOIN users u ON n.sender_id = u.id
LEFT JOIN notification_destinations nd ON n.id = nd.notification_id
GROUP BY n.id, p.notify_key, p.company_id, c.name, u.email, n.message_type, n.content, n.priority, n.status, n.created_at, n.sent_at;

-- Bot Status Summary View
CREATE OR REPLACE VIEW bot_status_summary AS
SELECT 
    tb.id AS bot_id,
    tb.name AS bot_name,
    tb.status::text AS status,
    CASE WHEN tb.type='platform' THEN 'platform'::bot_type ELSE 'third_party'::bot_type END AS bot_type,
    tb.app_id,
    COUNT(DISTINCT bi.id) AS installation_count,
    COUNT(DISTINCT CASE WHEN bi.installation_status = 'active' THEN bi.id END) AS active_installations,
    h.status AS health_status,
    h.last_check_at,
    h.response_time_ms
FROM teams_bots tb
LEFT JOIN bot_installations bi 
  ON bi.bot_id = tb.id 
 AND bi.bot_type = CASE WHEN tb.type='platform' THEN 'platform'::bot_type ELSE 'third_party'::bot_type END
LEFT JOIN LATERAL (
    SELECT status, last_check_at, response_time_ms
    FROM bot_health_status bhs
    WHERE bhs.bot_id = tb.id AND bhs.bot_type = CASE WHEN tb.type='platform' THEN 'platform'::bot_type ELSE 'third_party'::bot_type END
    ORDER BY last_check_at DESC
    LIMIT 1
) h ON true
GROUP BY tb.id, tb.name, tb.status, tb.app_id, h.status, h.last_check_at, h.response_time_ms;

-- Usage Summary View
CREATE VIEW usage_summary AS
SELECT 
    ur.project_id,
    p.notify_key,
    p.company_id,
    c.name as company_name,
    ur.record_type,
    SUM(ur.quantity) as total_quantity,
    SUM(ur.total_cost) as total_cost,
    ur.created_at::date as usage_date
FROM usage_records ur
JOIN projects p ON ur.project_id = p.id
JOIN companies c ON p.company_id = c.id
GROUP BY ur.project_id, p.notify_key, p.company_id, c.name, ur.record_type, ur.created_at::date;

-- File Summary View
CREATE VIEW file_summary AS
SELECT 
    f.project_id,
    p.notify_key,
    p.company_id,
    c.name as company_name,
    COUNT(*) as total_files,
    SUM(f.file_size) as total_size,
    COUNT(CASE WHEN f.is_public THEN 1 END) as public_files,
    COUNT(CASE WHEN f.expires_at IS NOT NULL AND f.expires_at > NOW() THEN 1 END) as active_files
FROM files f
JOIN projects p ON f.project_id = p.id
JOIN companies c ON p.company_id = c.id
GROUP BY f.project_id, p.notify_key, p.company_id, c.name;


-- =============================================
-- Indexes for Performance
-- =============================================

-- Bot installations indexes
CREATE INDEX idx_bot_installations_bot_tenant ON bot_installations(bot_id, bot_type, teams_tenant_id);
-- removed invalid legacy indexes (scope, teams_* columns not present)

-- Destinations targets indexes (GIN for JSONB queries)
CREATE INDEX idx_destinations_targets_gin ON destinations USING GIN(targets);

-- Notifications indexes
CREATE INDEX idx_notifications_sender ON notifications(sender_id);
CREATE INDEX idx_notifications_created_at ON notifications(created_at);

-- =============================================
