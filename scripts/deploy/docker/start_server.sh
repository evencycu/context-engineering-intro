#!/bin/bash

# Start server script with correct database configuration
# This script starts the TeamsNotify server with queue functionality

echo "🚀 Starting TeamsNotify Server with Queue Support..."
echo ""

# Check if server binary exists
if [ ! -f "./server" ]; then
    echo "❌ Server binary not found. Building..."
    go build -o server cmd/server/main.go
    if [ $? -ne 0 ]; then
        echo "❌ Build failed!"
        exit 1
    fi
    echo "✅ Build successful"
fi

# Check if postgres is running
echo "🔍 Checking PostgreSQL connection..."
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "SELECT 1" > /dev/null 2>&1
if [ $? -ne 0 ]; then
    echo "❌ Cannot connect to PostgreSQL (notification_center database)"
    echo "   Please ensure Docker container 'teamsnotify-postgres' is running"
    exit 1
fi
echo "✅ PostgreSQL connection OK"

# Ensure notification_destinations has async columns (next_retry_at, failure_reason, retry_after)
echo "🔍 Checking notification_destinations schema..."
NEED_MIGRATION=$(docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -t -c "
SELECT CASE WHEN 
    EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='notification_destinations' AND column_name='next_retry_at') AND
    EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='notification_destinations' AND column_name='failure_reason') AND
    EXISTS (SELECT 1 FROM information_schema.columns WHERE table_name='notification_destinations' AND column_name='retry_after')
THEN 'no' ELSE 'yes' END;" | xargs)

if [ "$NEED_MIGRATION" = "yes" ]; then
    echo "⚠️  Missing async columns. Applying migration 009..."
    docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < scripts/migrations/009_enhance_notification_destinations_for_async_actor.sql
    if [ $? -eq 0 ]; then
        echo "✅ Migration 009 applied"
    else
        echo "❌ Migration 009 failed"
        exit 1
    fi
fi

# Cleanup legacy failed_notifications artifacts if present
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "DROP VIEW IF EXISTS failed_notifications_queue_status CASCADE;" >/dev/null 2>&1
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "DROP TABLE IF EXISTS failed_notifications CASCADE;" >/dev/null 2>&1

# Kill existing server process
echo "🔍 Checking for existing server process..."
pkill -f './server' 2>/dev/null
if [ $? -eq 0 ]; then
    echo "✅ Stopped existing server"
    sleep 1
fi

# Set environment variables
export DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/notification_center?sslmode=disable"
export TEAMS_BOT_APP_ID="${TEAMS_BOT_APP_ID:-844146d7-4ac9-4e4d-a463-d6e027714e81}"
export TEAMS_TENANT_ID="${TEAMS_TENANT_ID:-051cece0-e4dc-4aed-b471-bf29824e1ee6}"
export TEAMS_BOT_APP_PASSWORD="${TEAMS_BOT_APP_PASSWORD:-HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy}"

echo ""
echo "📋 Configuration:"
echo "   Database: notification_center"
echo "   App ID: $TEAMS_BOT_APP_ID"
echo "   Tenant ID: $TEAMS_TENANT_ID"
echo ""

# Start server
echo "🚀 Starting server..."
echo "   Log file: server.log"
echo "   PID file: server.pid"
echo ""

nohup ./server > server.log 2>&1 &
SERVER_PID=$!
echo $SERVER_PID > server.pid

# Wait a moment and check if server started
sleep 2

if ps -p $SERVER_PID > /dev/null; then
    echo "✅ Server started successfully (PID: $SERVER_PID)"
    echo ""
    echo "📊 Quick Status Check:"
    sleep 1
    
    # Test health endpoint
    HEALTH=$(curl -s http://localhost:8080/health 2>/dev/null)
    if [ $? -eq 0 ]; then
        echo "   Health: $(echo $HEALTH | jq -r '.status' 2>/dev/null || echo 'OK')"
    else
        echo "   Health: Waiting for server..."
    fi
    
    # Test queue endpoint
    QUEUE=$(curl -s http://localhost:8080/api/v1/queue/status 2>/dev/null)
    if [ $? -eq 0 ]; then
        echo "   Queue: $(echo $QUEUE | jq -r '.circuit_state' 2>/dev/null || echo 'OK')"
    else
        echo "   Queue: Initializing..."
    fi
    
    echo ""
    echo "📚 Available Commands:"
    echo "   Stop server:    pkill -f './server' or kill \$(cat server.pid)"
    echo "   View logs:      tail -f server.log"
    echo "   Test queue:     ./scripts/testing/test_queue.sh"
    echo "   Queue status:   curl http://localhost:8080/api/v1/queue/status | jq '.'"
    echo ""
else
    echo "❌ Server failed to start. Check server.log for details:"
    tail -20 server.log
    exit 1
fi

