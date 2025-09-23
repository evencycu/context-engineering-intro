-- Clear all data from the database
-- WARNING: This will delete ALL data in the database

-- Disable foreign key checks temporarily
SET session_replication_role = replica;

-- Clear all tables in reverse dependency order
TRUNCATE TABLE notification_destinations CASCADE;
TRUNCATE TABLE notifications CASCADE;
TRUNCATE TABLE destinations CASCADE;
TRUNCATE TABLE bot_installations CASCADE;
TRUNCATE TABLE third_party_bot_api_usage CASCADE;
TRUNCATE TABLE third_party_bot_api_keys CASCADE;
TRUNCATE TABLE third_party_bots CASCADE;
TRUNCATE TABLE platform_bots CASCADE;
TRUNCATE TABLE projects CASCADE;
TRUNCATE TABLE users CASCADE;
TRUNCATE TABLE companies CASCADE;
TRUNCATE TABLE bot_routing_rules CASCADE;
TRUNCATE TABLE bot_health_status CASCADE;
TRUNCATE TABLE company_billing CASCADE;
TRUNCATE TABLE usage_records CASCADE;

-- Re-enable foreign key checks
SET session_replication_role = DEFAULT;

-- Reset sequences (if any)
-- Note: PostgreSQL doesn't have sequences for UUID primary keys

-- Verify tables are empty
SELECT 'Companies' as table_name, COUNT(*) as count FROM companies
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
SELECT 'Notifications', COUNT(*) FROM notifications;