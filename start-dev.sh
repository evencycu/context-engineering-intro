#!/bin/bash

# Teams Notification Center - Development Startup Script
# This script helps you start both Backend and Frontend for development

set -e

PROJECT_ROOT=$(dirname "$0")
BACKEND_DIR="$PROJECT_ROOT/Backend"
FRONTEND_DIR="$PROJECT_ROOT/Frontend"

echo "========================================"
echo "  Teams Notification Center Dev Setup"
echo "========================================"

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if docker is running
if ! docker info > /dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  Docker is not running. Please start Docker first.${NC}"
    exit 1
fi

echo ""
echo "Step 1: Starting PostgreSQL & Redis..."
cd "$BACKEND_DIR"
docker-compose -f scripts/docker/docker-compose.yml up -d postgres redis
sleep 5

echo ""
echo "Step 2: Checking database connection..."
docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center -c "SELECT 1" > /dev/null 2>&1 && \
    echo -e "${GREEN}✅ Database connected${NC}" || \
    echo -e "${YELLOW}⚠️  Database not ready, waiting...${NC}"

echo ""
echo -e "${GREEN}========================================"
echo "  Infrastructure is ready!"
echo "========================================${NC}"
echo ""
echo "Now open TWO terminals and run:"
echo ""
echo -e "${YELLOW}Terminal 1 - Backend (port 8080):${NC}"
echo "  cd $BACKEND_DIR"
echo "  source configs/.env && go run ./cmd/server"
echo ""
echo -e "${YELLOW}Terminal 2 - Frontend (port 3000):${NC}"
echo "  cd $FRONTEND_DIR"
echo "  npm run dev"
echo ""
echo "Then open: http://localhost:3000"
echo ""
