-- Initialize database with sample data
-- This script runs after schema.sql

-- Insert sample companies
INSERT INTO companies (id, name, contact_email, contact_phone, address, status, billing_enabled, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440001', '國泰金控', 'admin@cathayholdings.com.tw', '+1-555-0101', '123 Business St, City, State', 'active', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440002', '國泰人壽', 'admin@cathaylife.com.tw', '+1-555-0102', '456 Innovation Ave, Tech City, CA', 'active', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440003', '國泰世華銀行', 'admin@cathaybank.com.tw', '+1-555-0103', '789 Enterprise Blvd, Metro City, NY', 'active', false, NOW(), NOW());

-- Insert sample users
INSERT INTO users (id, company_id, email, name, role, status, password_hash, created_at, updated_at) VALUES
('650e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'admin@cathayholdings.com.tw', 'Admin User', 'admin', 'active', '$2a$10$example.hash', NOW(), NOW()),
('650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001', 'admin@cathaylife.com.tw', 'John Doe', 'user', 'active', '$2a$10$example.hash', NOW(), NOW()),
('650e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440002', 'admin@cathaybank.com.tw', 'Jane Smith', 'user', 'active', '$2a$10$example.hash', NOW(), NOW());

-- Insert sample projects
INSERT INTO projects (id, company_id, notify_key, description, status, daily_limit, monthly_limit, priority, created_by, created_at, updated_at) VALUES
('750e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'cfh-alert-gogo', 'Holdings notification alerts', 'active', 1000, 30000, 'normal', '650e8400-e29b-41d4-a716-446655440001', NOW(), NOW()),
('750e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 'techstart-monitoring', 'TechStart monitoring notifications', 'active', 500, 15000, 'high', '650e8400-e29b-41d4-a716-446655440003', NOW(), NOW());

-- Insert sample teams bots (type = platform)
INSERT INTO teams_bots (id, type, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, capabilities, rate_limit_per_minute, max_concurrent_requests, created_at, updated_at) VALUES
('850e8400-e29b-41d4-a716-446655440001', 'platform', 'Main Notification Bot', 'Primary bot for sending notifications', 'app-12345', '$2a$10$example.hash', 'tenant-12345', 'active', 'https://api.teams.com/webhook/12345', '{"send_message": true, "send_file": true, "send_adaptive_card": true}', 60, 10, NOW(), NOW()),
('850e8400-e29b-41d4-a716-446655440002', 'platform', 'Lab Test Teams Notify', 'Lab test bot for Teams notifications', '844146d7-4ac9-4e4d-a463-d6e027714e81', 'HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy', '051cece0-e4dc-4aed-b471-bf29824e1ee6', 'active', 'https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/', '{"send_message": true, "send_file": true, "send_adaptive_card": true}', 60, 10, NOW(), NOW());

-- Insert sample bot installations
INSERT INTO bot_installations (id, bot_id, bot_type, teams_tenant_id, conversation_type, conversation_id, service_url, recipient_id, recipient_name, from_id, from_name, from_aad_object_id, installation_status, installed_at, last_activity_at, metadata, created_at, updated_at) VALUES
('a60e8400-e29b-41d4-a716-446655440001', '850e8400-e29b-41d4-a716-446655440002', 'platform', '051cece0-e4dc-4aed-b471-bf29824e1ee6', 'personal', 'a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR', 'https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/', '28:844146d7-4ac9-4e4d-a463-d6e027714e81', 'lab-test-teams-notify', '29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw', '', '8d40db4b-935f-4e5a-aaae-ba86faad0e12', 'active', NOW(), NOW(), '{"action": "add", "channelData": {"settings": {"selectedChannel": {"id": "19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces"}}, "source": {"name": "message"}, "tenant": {"id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"}}, "channelId": "msteams", "conversation": {"conversationType": "personal", "id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR", "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6"}, "entities": [{"locale": "zh-TW", "type": "clientInfo"}], "from": {"aadObjectId": "8d40db4b-935f-4e5a-aaae-ba86faad0e12", "id": "29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw"}, "id": "f:757cd071-d818-ed6a-a463-6d986aa0718d", "locale": "zh-TW", "recipient": {"id": "28:844146d7-4ac9-4e4d-a463-d6e027714e81", "name": "lab-test-teams-notify"}, "serviceUrl": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/", "timestamp": "2025-09-18T06:42:40.192Z", "type": "installationUpdate"}'::jsonb, NOW(), NOW());

-- Insert sample destinations
INSERT INTO destinations (id, project_id, name, description, teams_tenant_id, targets, bot_id, bot_type, status, validation_status, created_by, created_at, updated_at) VALUES
('950e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440001', 'Development Team', 'Development team notifications', 'tenant-12345', '[{"type": "channel", "team_id": "team-123", "channel_id": "channel-dev", "display_name": "Development Channel"}]', '850e8400-e29b-41d4-a716-446655440001', 'platform', 'active', 'validated', '650e8400-e29b-41d4-a716-446655440001', NOW(), NOW()),
('950e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440001', 'Management Team', 'Management team notifications', 'tenant-12345', '[{"type": "channel", "team_id": "team-123", "channel_id": "channel-mgmt", "display_name": "Management Channel"}, {"type": "person", "user_id": "user-mgmt-1", "display_name": "Manager 1"}]', '850e8400-e29b-41d4-a716-446655440001', 'platform', 'active', 'validated', '650e8400-e29b-41d4-a716-446655440001', NOW(), NOW());

-- Insert sample billing plans
INSERT INTO billing_plans (id, name, description, price_per_notification, price_per_month, max_notifications_per_month, max_projects, max_users, features, is_active, created_at, updated_at) VALUES
('a50e8400-e29b-41d4-a716-446655440001', 'Basic Plan', 'Basic notification plan', 0.01, 0.00, 1000, 2, 5, '{"notifications": true}', true, NOW(), NOW()),
('a50e8400-e29b-41d4-a716-446655440002', 'Pro Plan', 'Professional notification plan', 0.005, 25.00, 5000, 10, 25, '{"notifications": true, "files": true}', true, NOW(), NOW()),
('a50e8400-e29b-41d4-a716-446655440003', 'Enterprise Plan', 'Enterprise notification plan', 0.001, 100.00, 50000, 50, 100, '{"notifications": true, "files": true, "analytics": true, "priority_support": true}', true, NOW(), NOW());

-- Insert company billing
INSERT INTO company_billing (id, company_id, billing_plan_id, billing_status, payment_method, billing_cycle, next_billing_date, total_usage_cost, created_at, updated_at) VALUES
('b50e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'a50e8400-e29b-41d4-a716-446655440002', 'active', 'credit_card', 'monthly', '2025-02-01', 15.50, NOW(), NOW()),
('b50e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 'a50e8400-e29b-41d4-a716-446655440001', 'active', 'bank_transfer', 'monthly', '2025-02-01', 8.75, NOW(), NOW());

-- Insert sample usage records
INSERT INTO usage_records (id, company_id, project_id, user_id, record_type, quantity, unit_cost, total_cost, metadata, created_at) VALUES
('c50e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440001', '650e8400-e29b-41d4-a716-446655440001', 'notification', 100, 0.01, 1.00, '{"message_type": "text", "priority": "normal"}', NOW()),
('c50e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440001', '650e8400-e29b-41d4-a716-446655440001', 'notification', 50, 0.01, 0.50, '{"message_type": "file", "priority": "high"}', NOW()),
('c50e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440002', '650e8400-e29b-41d4-a716-446655440003', 'notification', 200, 0.005, 1.00, '{"message_type": "text", "priority": "normal"}', NOW());

-- Insert sample files
INSERT INTO files (id, project_id, file_name, file_size, content_type, file_key, file_url, description, tags, is_public, expires_at, created_at, updated_at) VALUES
('d50e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440001', 'report.pdf', 1024000, 'application/pdf', 'files/2025/01/05/report_20250105.pdf', 'https://storage.example.com/files/2025/01/05/report_20250105.pdf', 'Monthly financial report', '{"report", "financial", "monthly"}', false, '2025-02-05 23:59:59', NOW(), NOW()),
('d50e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440001', 'presentation.pptx', 2048000, 'application/vnd.openxmlformats-officedocument.presentationml.presentation', 'files/2025/01/05/presentation_20250105.pptx', 'https://storage.example.com/files/2025/01/05/presentation_20250105.pptx', 'Q1 2025 presentation', '{"presentation", "q1", "2025"}', true, NULL, NOW(), NOW()),
('d50e8400-e29b-41d4-a716-446655440003', '750e8400-e29b-41d4-a716-446655440002', 'data.xlsx', 512000, 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet', 'files/2025/01/05/data_20250105.xlsx', 'https://storage.example.com/files/2025/01/05/data_20250105.xlsx', 'Analytics data export', '{"data", "analytics", "export"}', false, '2025-01-12 23:59:59', NOW(), NOW());


-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_notifications_project_id ON notifications(project_id);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
CREATE INDEX IF NOT EXISTS idx_destinations_project_id ON destinations(project_id);
CREATE INDEX IF NOT EXISTS idx_destinations_bot_id ON destinations(bot_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_company_id ON usage_records(company_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_project_id ON usage_records(project_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_user_id ON usage_records(user_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_record_type ON usage_records(record_type);
CREATE INDEX IF NOT EXISTS idx_usage_records_created_at ON usage_records(created_at);
CREATE INDEX IF NOT EXISTS idx_company_billing_company_id ON company_billing(company_id);
CREATE INDEX IF NOT EXISTS idx_company_billing_billing_plan_id ON company_billing(billing_plan_id);
CREATE INDEX IF NOT EXISTS idx_files_project_id ON files(project_id);
CREATE INDEX IF NOT EXISTS idx_files_file_key ON files(file_key);
CREATE INDEX IF NOT EXISTS idx_files_is_public ON files(is_public);
CREATE INDEX IF NOT EXISTS idx_files_expires_at ON files(expires_at);
CREATE INDEX IF NOT EXISTS idx_files_created_at ON files(created_at);
