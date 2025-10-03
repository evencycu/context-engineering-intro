-- Enhanced test data for Teams Notification Bot Platform
-- Updated to match the new bot_installations schema

-- =============================================
-- Companies
-- =============================================
INSERT INTO companies (id, name, contact_email, contact_phone, address, status, billing_enabled) VALUES
('9589aad2-f5ec-4d25-8ef7-b818b1933c32', '台灣科技股份有限公司', 'info@taiwan-tech.com', '+886-2-2345-6789', '台北市信義區信義路五段7號101大樓', 'active', true),
('a1b2c3d4-e5f6-7890-abcd-ef1234567890', '測試企業有限公司', 'test@test-company.com', '+886-2-1234-5678', '台北市中山區南京東路二段100號', 'active', true),
('b2c3d4e5-f6g7-8901-bcde-f23456789012', '創新科技集團', 'contact@innovate-tech.com', '+886-2-3456-7890', '新北市板橋區文化路一段188號', 'active', true);

-- =============================================
-- Users
-- =============================================
INSERT INTO users (id, company_id, email, name, role, status, password_hash, api_key_hash, api_key_expires_at) VALUES
('643c4d7a-a18f-4caa-aff2-5a0d4439b367', '9589aad2-f5ec-4d25-8ef7-b818b1933c32', 'admin@taiwan-tech.com', '系統管理員', 'admin', 'active', '$2a$10$example_hash_1', '$2a$10$api_key_hash_1', NOW() + INTERVAL '1 year'),
('754d5e8b-b29g-5dbb-bgg3-6b1e5540c478', '9589aad2-f5ec-4d25-8ef7-b818b1933c32', 'user1@taiwan-tech.com', '一般使用者', 'user', 'active', '$2a$10$example_hash_2', '$2a$10$api_key_hash_2', NOW() + INTERVAL '1 year'),
('865e6f9c-c3ah-6ecc-chh4-7c2f6651d589', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'manager@test-company.com', '專案經理', 'admin', 'active', '$2a$10$example_hash_3', '$2a$10$api_key_hash_3', NOW() + INTERVAL '1 year');

-- =============================================
-- Projects
-- =============================================
INSERT INTO projects (id, company_id, notify_key, description, status, daily_limit, monthly_limit, priority, created_by) VALUES
('198f1130-20a9-4c7c-a504-6a055d27e8db', '9589aad2-f5ec-4d25-8ef7-b818b1933c32', 'main-notifications', '主要通知系統', 'active', 10000, 300000, 'high', '643c4d7a-a18f-4caa-aff2-5a0d4439b367'),
('2a9g2241-31ba-5d8d-b615-7b166e38f9ec', '9589aad2-f5ec-4d25-8ef7-b818b1933c32', 'test-project', '測試專案', 'active', 5000, 150000, 'normal', '754d5e8b-b29g-5dbb-bgg3-6b1e5540c478'),
('3bah3352-42cb-6e9e-c726-8c277f49g0fd', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'enterprise-alerts', '企業警報系統', 'active', 20000, 600000, 'high', '865e6f9c-c3ah-6ecc-chh4-7c2f6651d589');

-- =============================================
-- Platform Bots
-- =============================================
INSERT INTO platform_bots (id, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, capabilities, rate_limit_per_minute, max_concurrent_requests) VALUES
('844146d7-4ac9-4e4d-a463-d6e027714e81', '主要通知機器人', '平台統一管理的通知機器人', '844146d7-4ac9-4e4d-a463-d6e027714e81', 'HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy', '051cece0-e4dc-4aed-b471-bf29824e1ee6', 'active', 'https://api.taiwan-tech.com/webhook/bot', '{"receive_message": true, "send_message": true, "file_upload": true, "adaptive_cards": true}', 600, 100),
('955257e8-5bda-5f5e-b574-e7f138825f92', '測試機器人', '用於測試的機器人', '955257e8-5bda-5f5e-b574-e7f138825f92', 'IWX9R~-InWQj_X1FaJQcajM5H2daW9XCNVkRgegz', '051cece0-e4dc-4aed-b471-bf29824e1ee6', 'active', 'https://api.taiwan-tech.com/webhook/test-bot', '{"receive_message": true, "send_message": true, "file_upload": false, "adaptive_cards": true}', 300, 50);

-- =============================================
-- Third Party Bots
-- =============================================
INSERT INTO third_party_bots (id, company_id, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, api_endpoint, api_key_hash, capabilities, rate_limit_per_minute, max_concurrent_requests, contact_email, contact_phone, created_by) VALUES
('a1b2c3d4-e5f6-7890-abcd-ef1234567890', '9589aad2-f5ec-4d25-8ef7-b818b1933c32', '第三方通知機器人', '第三方企業的專用機器人', 'third-party-app-123', 'third-party-password-123', '051cece0-e4dc-4aed-b471-bf29824e1ee6', 'active', 'https://third-party.com/webhook', 'https://third-party.com/api', 'third-party-api-key-123', '{"receive_message": true, "send_message": true, "file_upload": true, "adaptive_cards": true}', 200, 20, 'contact@third-party.com', '+886-2-9876-5432', '643c4d7a-a18f-4caa-aff2-5a0d4439b367');

-- =============================================
-- Bot Installations (Enhanced Schema)
-- =============================================
INSERT INTO bot_installations (
    id, bot_id, bot_type, teams_tenant_id, scope, conversation_id, service_url,
    teams_team_id, teams_channel_id, teams_chat_id, teams_user_id,
    recipient_id, from_id, from_aad_object_id,
    installation_status, installed_at, last_activity_at, metadata
) VALUES
-- Personal installations
('11111111-1111-1111-1111-111111111111', '844146d7-4ac9-4e4d-a463-d6e027714e81', 'platform', '051cece0-e4dc-4aed-b471-bf29824e1ee6', 'personal', 
 'a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR',
 'https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/',
 NULL, NULL, NULL, '29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw',
 '28:844146d7-4ac9-4e4d-a463-d6e027714e81', '29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw', '8d40db4b-935f-4e5a-aaae-ba86faad0e12',
 'active', NOW(), NOW(), '{"locale": "zh-TW", "source": "message"}'),

-- Team channel installations
('22222222-2222-2222-2222-222222222222', '844146d7-4ac9-4e4d-a463-d6e027714e81', 'platform', '051cece0-e4dc-4aed-b471-bf29824e1ee6', 'team',
 '19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces',
 'https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/',
 '19:8d40db4b-935f-4e5a-aaae-ba86faad0e12', '19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces', NULL, NULL,
 '28:844146d7-4ac9-4e4d-a463-d6e027714e81', '29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw', '8d40db4b-935f-4e5a-aaae-ba86faad0e12',
 'active', NOW(), NOW(), '{"team": {"id": "19:8d40db4b-935f-4e5a-aaae-ba86faad0e12"}, "channel": {"id": "19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces"}}'),

-- Group chat installations
('33333333-3333-3333-3333-333333333333', '844146d7-4ac9-4e4d-a463-d6e027714e81', 'platform', '051cece0-e4dc-4aed-b471-bf29824e1ee6', 'groupChat',
 '19:meeting_abc123def456ghi789jkl012mno345pqr678stu901vwx234yz@thread.v2',
 'https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/',
 NULL, NULL, '19:meeting_abc123def456ghi789jkl012mno345pqr678stu901vwx234yz@thread.v2', NULL,
 '28:844146d7-4ac9-4e4d-a463-d6e027714e81', '29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw', '8d40db4b-935f-4e5a-aaae-ba86faad0e12',
 'active', NOW(), NOW(), '{"chat": {"id": "19:meeting_abc123def456ghi789jkl012mno345pqr678stu901vwx234yz@thread.v2"}}'),

-- Third party bot installations
('44444444-4444-4444-4444-444444444444', 'a1b2c3d4-e5f6-7890-abcd-ef1234567890', 'third_party', '051cece0-e4dc-4aed-b471-bf29824e1ee6', 'personal',
 'a:thirdparty123456789abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890',
 'https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/',
 NULL, NULL, NULL, '29:thirdparty123456789abcdefghijklmnopqrstuvwxyz1234567890',
 '28:a1b2c3d4-e5f6-7890-abcd-ef1234567890', '29:thirdparty123456789abcdefghijklmnopqrstuvwxyz1234567890', 'thirdparty-aad-object-id-123',
 'active', NOW(), NOW(), '{"locale": "zh-TW", "source": "third_party"}');

-- =============================================
-- Destinations (Enhanced Targets)
-- =============================================
INSERT INTO destinations (id, project_id, name, description, teams_tenant_id, targets, bot_id, bot_type, status, validation_status, created_by) VALUES
('c5133b8b-6e5c-4359-be17-221483569b97', '198f1130-20a9-4c7c-a504-6a055d27e8db', '主要通知目標', '包含所有類型的通知目標', '051cece0-e4dc-4aed-b471-bf29824e1ee6',
 '[{"type": "person", "user_id": "29:1hr09MmriLZ1ymlHMEG_kFBtId2m8WRHaaxLtgJipvUflqrdoeOatXhzKA1LsQGZPOAMqrSEvoRNThsTlBxUcLw", "aad_object_id": "8d40db4b-935f-4e5a-aaae-ba86faad0e12", "conversation_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR", "display_name": "測試使用者", "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6", "service_url": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/"}, {"type": "channel", "team_id": "19:8d40db4b-935f-4e5a-aaae-ba86faad0e12", "channel_id": "19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces", "conversation_id": "19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces", "display_name": "測試頻道", "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6", "service_url": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/"}, {"type": "chatgroup", "chat_id": "19:meeting_abc123def456ghi789jkl012mno345pqr678stu901vwx234yz@thread.v2", "conversation_id": "19:meeting_abc123def456ghi789jkl012mno345pqr678stu901vwx234yz@thread.v2", "display_name": "測試群組", "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6", "service_url": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/"}]',
 '844146d7-4ac9-4e4d-a463-d6e027714e81', 'platform', 'active', 'validated', '643c4d7a-a18f-4caa-aff2-5a0d4439b367'),

('d6244c9c-7d6d-5460-cf28-332594670a08', '2a9g2241-31ba-5d8d-b615-7b166e38f9ec', '測試目標', '用於測試的目標群組', '051cece0-e4dc-4aed-b471-bf29824e1ee6',
 '[{"type": "person", "user_id": "29:testuser123456789abcdefghijklmnopqrstuvwxyz1234567890", "aad_object_id": "test-user-aad-object-id-456", "conversation_id": "a:testuser123456789abcdefghijklmnopqrstuvwxyz1234567890abcdefghijklmnopqrstuvwxyz1234567890", "display_name": "測試用戶2", "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6", "service_url": "https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/"}]',
 '844146d7-4ac9-4e4d-a463-d6e027714e81', 'platform', 'active', 'validated', '754d5e8b-b29g-5dbb-bgg3-6b1e5540c478');

-- =============================================
-- Notifications (Sample)
-- =============================================
INSERT INTO notifications (id, project_id, sender_id, message_type, content, mentions, priority, status, metadata) VALUES
('e7355dad-8e7e-6571-dg39-4436a581b219', '198f1130-20a9-4c7c-a504-6a055d27e8db', '643c4d7a-a18f-4caa-aff2-5a0d4439b367', 'text', '這是一個測試通知消息 - 包含中文內容', '["@admin", "@user", "@manager"]', 'high', 'sent', '{"broadcast_type": "all", "sent_at": "2025-09-18T10:00:00Z"}'),
('f8466ebe-9f8f-7682-eh4a-5547b692c32a', '2a9g2241-31ba-5d8d-b615-7b166e38f9ec', '754d5e8b-b29g-5dbb-bgg3-6b1e5540c478', 'adaptive_card', '自適應卡片通知', '[]', 'normal', 'sent', '{"card_type": "welcome", "sent_at": "2025-09-18T10:05:00Z"}');

-- =============================================
-- Notification Destinations (Link notifications to destinations)
-- =============================================
INSERT INTO notification_destinations (id, notification_id, destination_id, bot_id, bot_type, status, teams_message_id, sent_at) VALUES
('11111111-1111-1111-1111-111111111111', 'e7355dad-8e7e-6571-dg39-4436a581b219', 'c5133b8b-6e5c-4359-be17-221483569b97', '844146d7-4ac9-4e4d-a463-d6e027714e81', 'platform', 'sent', 'teams_msg_123456789', NOW()),
('22222222-2222-2222-2222-222222222222', 'f8466ebe-9f8f-7682-eh4a-5547b692c32a', 'd6244c9c-7d6d-5460-cf28-332594670a08', '844146d7-4ac9-4e4d-a463-d6e027714e81', 'platform', 'sent', 'teams_msg_987654321', NOW());

-- =============================================
-- Verification Query
-- =============================================
SELECT 
    'Companies' as table_name, COUNT(*) as count 
FROM companies
UNION ALL
SELECT 'Users', COUNT(*) FROM users
UNION ALL
SELECT 'Projects', COUNT(*) FROM projects
UNION ALL
SELECT 'Platform Bots', COUNT(*) FROM platform_bots
UNION ALL
SELECT 'Third Party Bots', COUNT(*) FROM third_party_bots
UNION ALL
SELECT 'Bot Installations', COUNT(*) FROM bot_installations
UNION ALL
SELECT 'Destinations', COUNT(*) FROM destinations
UNION ALL
SELECT 'Notifications', COUNT(*) FROM notifications
UNION ALL
SELECT 'Notification Destinations', COUNT(*) FROM notification_destinations
ORDER BY table_name;