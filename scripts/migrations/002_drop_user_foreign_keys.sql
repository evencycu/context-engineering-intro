-- Migration: Drop foreign key constraints that reference users, except projects.created_by
-- Reason: users only records who applied for projects

-- audit_logs.user_id -> users.id
ALTER TABLE IF EXISTS audit_logs
    DROP CONSTRAINT IF EXISTS audit_logs_user_id_fkey;

-- destinations.created_by -> users.id
ALTER TABLE IF EXISTS destinations
    DROP CONSTRAINT IF EXISTS destinations_created_by_fkey;

-- notifications.sender_id -> users.id
ALTER TABLE IF EXISTS notifications
    DROP CONSTRAINT IF EXISTS notifications_sender_id_fkey;

-- third_party_bot_api_keys.created_by -> users.id
ALTER TABLE IF EXISTS third_party_bot_api_keys
    DROP CONSTRAINT IF EXISTS third_party_bot_api_keys_created_by_fkey;

-- third_party_bots.created_by -> users.id
ALTER TABLE IF EXISTS third_party_bots
    DROP CONSTRAINT IF EXISTS third_party_bots_created_by_fkey;

-- usage_records.user_id -> users.id
ALTER TABLE IF EXISTS usage_records
    DROP CONSTRAINT IF EXISTS usage_records_user_id_fkey;

-- Note: Keep projects.created_by -> users.id


