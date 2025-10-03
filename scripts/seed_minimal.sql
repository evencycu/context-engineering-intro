-- Minimal seed for notification_center compatible with Redis+Actor
-- Safe to run multiple times (uses ON CONFLICT where possible)

-- 1) Company
INSERT INTO companies (id, name, contact_email, status)
VALUES (
    'e4f160f4-4917-443d-86e9-29071967b768',
    'Test Co',
    'admin@test.co',
    'active'
)
ON CONFLICT (id) DO NOTHING;

-- 2) User (created_by)
INSERT INTO users (id, company_id, email, name, role, status)
VALUES (
    '11111111-1111-1111-1111-111111111111',
    'e4f160f4-4917-443d-86e9-29071967b768',
    'admin@test.co',
    'Admin User',
    'admin',
    'active'
)
ON CONFLICT (id) DO NOTHING;

-- 3) Project
INSERT INTO projects (id, company_id, notify_key, description, status, daily_limit, monthly_limit, priority, created_by)
VALUES (
    '5ba0494e-1054-4b8c-a72a-c33b642d1ad2',
    'e4f160f4-4917-443d-86e9-29071967b768',
    'proj-smoke',
    'Smoke test project',
    'active',
    10000,
    300000,
    'normal',
    '11111111-1111-1111-1111-111111111111'
)
ON CONFLICT (id) DO NOTHING;

-- 4) Platform bot (teams_bots)
INSERT INTO teams_bots (id, type, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, rate_limit_per_minute, max_concurrent_requests)
VALUES (
    '850e8400-e29b-41d4-a716-446655440001',
    'platform',
    'Main Notification Bot',
    'Primary bot for sending notifications',
    'app-12345',
    'secret-hash',
    'tenant-12345',
    'active',
    'https://api.teams.example/webhook/12345',
    60,
    10
)
ON CONFLICT (id) DO NOTHING;

-- 5) Destination with targets JSON
INSERT INTO destinations (id, project_id, name, description, teams_tenant_id, targets, bot_id, status, validation_status, created_by)
VALUES (
    '27abfbbf-9352-4ce8-8fdb-16f746bfe245',
    '5ba0494e-1054-4b8c-a72a-c33b642d1ad2',
    'default-destination',
    'Auto seeded default destination',
    'tenant-12345',
    '[{"type":"channel","tenant_id":"tenant-12345","conversation_id":"19:abc123@thread.v2"}]'::jsonb,
    '850e8400-e29b-41d4-a716-446655440001',
    'active',
    'validated',
    '11111111-1111-1111-1111-111111111111'
)
ON CONFLICT (id) DO NOTHING;

-- 6) Bot installation (links bot to tenant/conversation)
INSERT INTO bot_installations (
    id, bot_id, bot_type, teams_tenant_id, conversation_type, conversation_id, service_url,
    recipient_id, recipient_name, from_id, from_name, from_aad_object_id, email, description_name,
    installation_status
)
VALUES (
    'afd4c39e-8fa9-48dd-8f4f-a28f44dfbd23',
    '850e8400-e29b-41d4-a716-446655440001',
    'platform',
    'tenant-12345',
    'channel',
    '19:abc123@thread.v2',
    'https://smba.trafficmanager.net/apac/',
    '28:app-12345',
    'Main Notification Bot',
    'user-1',
    'Demo User',
    '00000000-0000-0000-0000-000000000000',
    'demo.user@test.co',
    'Demo User',
    'active'
)
ON CONFLICT (id) DO NOTHING;

-- 7) Optional: sample notification (no send)
INSERT INTO notifications (id, project_id, message_type, content, priority, status)
VALUES (
    '2c23440d-a8f8-4a5f-b122-9b57ebf26d3f',
    '5ba0494e-1054-4b8c-a72a-c33b642d1ad2',
    'text',
    'Notification DTO',
    'high',
    'pending'
)
ON CONFLICT (id) DO NOTHING;


