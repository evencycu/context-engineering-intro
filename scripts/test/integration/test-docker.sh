#!/bin/bash
# scripts/local-test/test-docker.sh

set -e

echo "🧪 測試本地 Docker 環境..."

# 檢查服務狀態
echo "📊 檢查服務狀態..."
docker-compose -f scripts/local-test/docker-compose-local.yml ps

# 測試 API Server
echo "🚀 測試 API Server..."
API_URL="http://localhost:8080"

# 健康檢查
echo "🔍 健康檢查..."
curl -s "$API_URL/health" | jq '.' || echo "❌ 健康檢查失敗"

# 指標端點
echo "📊 指標端點..."
curl -s "$API_URL/metrics" | head -10 || echo "❌ 指標端點失敗"

# 測試外部 API
echo "📡 測試外部 API..."
curl -s -X POST "$API_URL/api/v1/external/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notify_key": "test-key",
    "message": "Hello from local Docker test",
    "message_type": "text",
    "priority": "normal",
    "targets": ["all"]
  }' | jq '.' || echo "❌ 外部 API 測試失敗"

## (Admin 已移除)

# 測試資料庫連線
echo "🐘 測試資料庫連線..."
docker-compose -f scripts/local-test/docker-compose-local.yml exec postgres psql -U teamsnotify -d notification_center -c "SELECT version();" || echo "❌ 資料庫連線失敗"

# 測試 Redis 連線
echo "🔴 測試 Redis 連線..."
docker-compose -f scripts/local-test/docker-compose-local.yml exec redis redis-cli ping || echo "❌ Redis 連線失敗"

echo "✅ 本地 Docker 環境測試完成！"
