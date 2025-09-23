#!/bin/bash

# Reset Database Script
# This script clears the database, runs migrations, and loads enhanced test data

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Database connection parameters
DB_HOST=${DB_HOST:-localhost}
DB_PORT=${DB_PORT:-5432}
DB_NAME=${DB_NAME:-teams_notification_bot}
DB_USER=${DB_USER:-postgres}
DB_PASSWORD=${DB_PASSWORD:-password}

echo -e "${YELLOW}🔄 Resetting Teams Notification Bot Database...${NC}"

# Check if psql is available
if ! command -v psql &> /dev/null; then
    echo -e "${RED}❌ psql command not found. Please install PostgreSQL client tools.${NC}"
    exit 1
fi

# Set PGPASSWORD environment variable
export PGPASSWORD="$DB_PASSWORD"

echo -e "${YELLOW}📋 Step 1: Clearing existing data...${NC}"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$(dirname "$0")/clear_database.sql"

echo -e "${YELLOW}📋 Step 2: Running migration to enhance bot_installations table...${NC}"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$(dirname "$0")/migrations/001_enhance_bot_installations.sql"

echo -e "${YELLOW}📋 Step 3: Loading enhanced test data...${NC}"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -f "$(dirname "$0")/test_data_enhanced.sql"

echo -e "${GREEN}✅ Database reset completed successfully!${NC}"
echo ""
echo -e "${GREEN}📊 Database Summary:${NC}"
psql -h "$DB_HOST" -p "$DB_PORT" -U "$DB_USER" -d "$DB_NAME" -c "
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
"

echo ""
echo -e "${GREEN}🎯 Enhanced Features Available:${NC}"
echo "  • Bot installations with scope (personal/team/groupChat)"
echo "  • Conversation ID tracking for proactive messaging"
echo "  • Enhanced destination targets with AAD Object ID"
echo "  • Service URL and tenant-specific configurations"
echo "  • Activity tracking and stale detection"
echo ""
echo -e "${GREEN}🚀 Ready to test broadcast functionality!${NC}"