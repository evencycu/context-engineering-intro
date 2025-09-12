-- Initialize database with sample data
-- This script runs after schema.sql

-- Insert sample companies
INSERT INTO companies (id, name, contact_email, contact_phone, address, status, billing_enabled, created_at, updated_at) VALUES
('550e8400-e29b-41d4-a716-446655440001', 'Acme Corp', 'admin@acme.com', '+1-555-0101', '123 Business St, City, State', 'active', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440002', 'TechStart Inc', 'info@techstart.com', '+1-555-0102', '456 Innovation Ave, Tech City, CA', 'active', true, NOW(), NOW()),
('550e8400-e29b-41d4-a716-446655440003', 'Global Solutions', 'contact@global.com', '+1-555-0103', '789 Enterprise Blvd, Metro City, NY', 'active', false, NOW(), NOW());

-- Insert sample users
INSERT INTO users (id, company_id, email, name, role, status, password_hash, created_at, updated_at) VALUES
('650e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'admin@acme.com', 'Admin User', 'admin', 'active', '$2a$10$example.hash', NOW(), NOW()),
('650e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001', 'user1@acme.com', 'John Doe', 'user', 'active', '$2a$10$example.hash', NOW(), NOW()),
('650e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440002', 'dev@techstart.com', 'Jane Smith', 'developer', 'active', '$2a$10$example.hash', NOW(), NOW());

-- Insert sample projects
INSERT INTO projects (id, company_id, key_name, description, status, daily_limit, monthly_limit, priority, created_by, created_at, updated_at) VALUES
('750e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', 'acme-alerts', 'Acme Corp notification alerts', 'active', 1000, 30000, 'normal', '650e8400-e29b-41d4-a716-446655440001', NOW(), NOW()),
('750e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440002', 'techstart-monitoring', 'TechStart monitoring notifications', 'active', 500, 15000, 'high', '650e8400-e29b-41d4-a716-446655440003', NOW(), NOW());

-- Insert sample platform bots
INSERT INTO platform_bots (id, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, capabilities, rate_limit_per_minute, max_concurrent_requests, created_at, updated_at) VALUES
('850e8400-e29b-41d4-a716-446655440001', 'Main Notification Bot', 'Primary bot for sending notifications', 'app-12345', '$2a$10$example.hash', 'tenant-12345', 'active', 'https://api.teams.com/webhook/12345', '{"send_message": true, "send_file": true, "send_adaptive_card": true}', 60, 10, NOW(), NOW());

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
