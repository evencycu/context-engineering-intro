#!/bin/bash

# Script to create notification_center database and migrate data from codex_teams
# This ensures a clean database setup with the new naming convention

set -e

CONTAINER_NAME="teamsnotify-postgres"
DB_USER="teamsnotify"
SOURCE_DB="codex_teams"
TARGET_DB="notification_center"

echo "=========================================="
echo "  Create notification_center Database"
echo "=========================================="
echo ""

# Check if container is running
echo "🔍 Checking PostgreSQL container..."
if ! docker ps | grep -q "$CONTAINER_NAME"; then
    echo "❌ Container '$CONTAINER_NAME' is not running"
    echo "   Please start it with: docker start $CONTAINER_NAME"
    exit 1
fi
echo "✅ Container is running"
echo ""

# Check if source database exists
echo "🔍 Checking source database ($SOURCE_DB)..."
SOURCE_EXISTS=$(docker exec $CONTAINER_NAME psql -U $DB_USER -lqt | cut -d \| -f 1 | grep -w $SOURCE_DB | wc -l | xargs)
if [ "$SOURCE_EXISTS" -eq "0" ]; then
    echo "⚠️  Source database '$SOURCE_DB' not found"
    echo "   Will create new '$TARGET_DB' database from scratch"
    SOURCE_DB=""
else
    echo "✅ Source database found"
fi
echo ""

# Check if target database already exists
echo "🔍 Checking target database ($TARGET_DB)..."
TARGET_EXISTS=$(docker exec $CONTAINER_NAME psql -U $DB_USER -lqt | cut -d \| -f 1 | grep -w $TARGET_DB | wc -l | xargs)

if [ "$TARGET_EXISTS" -gt "0" ]; then
    echo "⚠️  Database '$TARGET_DB' already exists"
    read -p "   Do you want to drop and recreate it? (yes/no): " confirm
    if [ "$confirm" != "yes" ]; then
        echo "❌ Aborted"
        exit 1
    fi
    
    echo "🗑️  Dropping existing database..."
    docker exec $CONTAINER_NAME psql -U $DB_USER -c "DROP DATABASE IF EXISTS $TARGET_DB;"
    echo "✅ Dropped"
fi

# Create new database
echo ""
echo "🆕 Creating database '$TARGET_DB'..."
docker exec $CONTAINER_NAME psql -U $DB_USER -c "CREATE DATABASE $TARGET_DB WITH OWNER $DB_USER ENCODING 'UTF8';"
echo "✅ Database created"

# If source database exists, dump and restore
if [ -n "$SOURCE_DB" ]; then
    echo ""
    echo "📦 Migrating data from '$SOURCE_DB' to '$TARGET_DB'..."
    
    # Dump source database
    echo "   Dumping source database..."
    docker exec $CONTAINER_NAME pg_dump -U $DB_USER $SOURCE_DB > /tmp/db_dump.sql
    
    # Restore to target database
    echo "   Restoring to target database..."
    docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $TARGET_DB < /tmp/db_dump.sql > /dev/null 2>&1
    
    # Clean up
    rm /tmp/db_dump.sql
    
    echo "✅ Data migrated successfully"
else
    # No source database, run schema initialization
    echo ""
    echo "📋 Initializing schema from scratch..."
    
    # Check if schema.sql exists
    if [ -f "internal/database/schema.sql" ]; then
        echo "   Running schema.sql..."
        docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $TARGET_DB < internal/database/schema.sql
        echo "✅ Schema initialized"
    else
        echo "⚠️  schema.sql not found, skipping"
    fi
    
    # Run all migrations
    echo ""
    echo "   Running migrations..."
    for migration in scripts/migrations/*.sql; do
        if [ -f "$migration" ]; then
            echo "   - $(basename $migration)"
            docker exec -i $CONTAINER_NAME psql -U $DB_USER -d $TARGET_DB < "$migration" 2>&1 | grep -v "already exists\|does not exist, skipping" || true
        fi
    done
    echo "✅ Migrations completed"
fi

# Verify database
echo ""
echo "✅ Verifying database..."
TABLE_COUNT=$(docker exec $CONTAINER_NAME psql -U $DB_USER -d $TARGET_DB -t -c "SELECT COUNT(*) FROM information_schema.tables WHERE table_schema='public';" | xargs)
echo "   Found $TABLE_COUNT tables"

# Show tables
echo ""
echo "📊 Tables in $TARGET_DB:"
docker exec $CONTAINER_NAME psql -U $DB_USER -d $TARGET_DB -c "\dt"

echo ""
echo "=========================================="
echo "✅ Setup Complete!"
echo "=========================================="
echo ""
echo "📋 Database Information:"
echo "   Name: $TARGET_DB"
echo "   User: $DB_USER"
echo "   Tables: $TABLE_COUNT"
echo ""
echo "🔗 Connection String:"
echo "   postgresql://$DB_USER:teamsnotify123@localhost:5432/$TARGET_DB?sslmode=disable"
echo ""
echo "🚀 Next Steps:"
echo "   1. Update your environment: export DATABASE_URL=\"postgresql://$DB_USER:teamsnotify123@localhost:5432/$TARGET_DB?sslmode=disable\""
echo "   2. Start server: ./start_server.sh"
echo "   3. Test queue: ./scripts/test_queue.sh"
echo ""

