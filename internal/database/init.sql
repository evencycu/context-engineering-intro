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
INSERT INTO projects (id, company_id, key_name, description, status, daily_limit, monthly_limit, priority, created_by, created_at, updated_at) VALUES
('750e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'cfh-alert-gogo', 'Holdings notification alerts', 'active', 1000, 30000, 'normal', '650e8400-e29b-41d4-a716-446655440001', NOW(), NOW()),
('750e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 'techstart-monitoring', 'TechStart monitoring notifications', 'active', 500, 15000, 'high', '650e8400-e29b-41d4-a716-446655440003', NOW(), NOW());

-- Insert sample platform bots
INSERT INTO platform_bots (id, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, capabilities, rate_limit_per_minute, max_concurrent_requests, created_at, updated_at) VALUES
('850e8400-e29b-41d4-a716-446655440001', 'Main Notification Bot', 'Primary bot for sending notifications', 'app-12345', '$2a$10$example.hash', 'tenant-12345', 'active', 'https://api.teams.com/webhook/12345', '{"send_message": true, "send_file": true, "send_adaptive_card": true}', 60, 10, NOW(), NOW()),
('850e8400-e29b-41d4-a716-446655440002', 'Lab Test Teams Notify', 'Lab test bot for Teams notifications', '844146d7-4ac9-4e4d-a463-d6e027714e81', '$2a$10$example.hash', '051cece0-e4dc-4aed-b471-bf29824e1ee6', 'active', 'https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/', '{"send_message": true, "send_file": true, "send_adaptive_card": true}', 60, 10, NOW(), NOW());

-- Insert sample bot installations
INSERT INTO bot_installations (id, bot_id, bot_type, teams_tenant_id, conversation_type, conversation_id, service_url, recipient_id, recipient_name, from_id, from_name, from_aad_object_id, installation_status, installed_at, last_activity_at, metadata, created_at, updated_at) VALUES
('a60e8400-e29b-41d4-a716-446655440001', '850e8400-e29b-41d4-a716-446655440002', 'platform', '051cece0-e4dc-4aed-b471-bf29824e1ee6', 'personal', 'a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR', 'https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/', '28:844146d7-4ac9-4e4d-a463-d6e027714e81', 'lab-test-teams-notify', '29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw', '', '8d40db4b-935f-4e5a-aaae-ba86faad0e12', 'active', NOW(), NOW(), '{"action": "add", "channelData": {"settings": {"selectedChannel": {"id": "19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces"}}, "source": {"name": "message"}, "tenant": {"id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"}}, "channelId": "msteams", "conversation": {"conversationType": "personal", "id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR", "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6"}, "entities": [{"locale": "zh-TW", "type": "clientInfo"}], "from": {"aadObjectId": "8d40db4b-935f-4e5a-aaae-ba86faad0e12", "id": "29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw"}, "id": "f:757cd071-d818-ed6a-a463-6d986aa0718d", "locale": "zh-TW", "recipient": {"id": "28:844146d7-4ac9-4e4d-a463-d6e027714e81", "name": "lab-test-teams-notify"}, "serviceUrl": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/", "timestamp": "2025-09-18T06:42:40.192Z", "type": "installationUpdate"}'::jsonb, NOW(), NOW());

-- Insert sample destinations
INSERT INTO destinations (id, project_id, name, description, teams_tenant_id, targets, bot_id, bot_type, status, validation_status, created_by, created_at, updated_at) VALUES
('950e8400-e29b-41d4-a716-446655440001', '750e8400-e29b-41d4-a716-446655440001', 'Development Team', 'Development team notifications', 'tenant-12345', '[{"type": "channel", "team_id": "team-123", "channel_id": "channel-dev", "display_name": "Development Channel"}]', '850e8400-e29b-41d4-a716-446655440001', 'platform', 'active', 'validated', '650e8400-e29b-41d4-a716-446655440001', NOW(), NOW()),
('950e8400-e29b-41d4-a716-446655440002', '750e8400-e29b-41d4-a716-446655440001', 'Management Team', 'Management team notifications', 'tenant-12345', '[{"type": "channel", "team_id": "team-123", "channel_id": "channel-mgmt", "display_name": "Management Channel"}, {"type": "person", "user_id": "user-mgmt-1", "display_name": "Manager 1"}]', '850e8400-e29b-41d4-a716-446655440001', 'platform', 'active', 'validated', '650e8400-e29b-41d4-a716-446655440001', NOW(), NOW());

-- Insert sample billing plans
INSERT INTO billing_plans (id, name, description, price_per_notification, price_per_attachment, price_per_mention, price_per_adaptive_card, daily_limit, monthly_limit, status, created_at, updated_at) VALUES
('a50e8400-e29b-41d4-a716-446655440001', 'Basic Plan', 'Basic notification plan', 0.001, 0.005, 0.0005, 0.01, 1000, 30000, 'active', NOW(), NOW()),
('a50e8400-e29b-41d4-a716-446655440002', 'Pro Plan', 'Professional notification plan', 0.0005, 0.002, 0.0002, 0.005, 5000, 150000, 'active', NOW(), NOW()),
('a50e8400-e29b-41d4-a716-446655440003', 'Enterprise Plan', 'Enterprise notification plan', 0.0001, 0.001, 0.0001, 0.002, 50000, 1500000, 'active', NOW(), NOW());

-- Insert company billing
INSERT INTO company_billing (id, company_id, billing_plan_id, status, billing_email, payment_method, currency, created_at, updated_at) VALUES
('b50e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'a50e8400-e29b-41d4-a716-446655440002', 'active', 'billing@acme.com', 'credit_card', 'USD', NOW(), NOW()),
('b50e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 'a50e8400-e29b-41d4-a716-446655440001', 'active', 'billing@techstart.com', 'paypal', 'USD', NOW(), NOW());

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_notifications_project_id ON notifications(project_id);
CREATE INDEX IF NOT EXISTS idx_notifications_status ON notifications(status);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);
CREATE INDEX IF NOT EXISTS idx_destinations_project_id ON destinations(project_id);
CREATE INDEX IF NOT EXISTS idx_destinations_bot_id ON destinations(bot_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_company_id ON usage_records(company_id);
CREATE INDEX IF NOT EXISTS idx_usage_records_billing_period ON usage_records(billing_period);
