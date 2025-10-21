#!/bin/bash
# scripts/local-test/test-api.sh

set -e

echo "🧪 測試 UT 環境 API..."

# 取得 API Server URL
API_URL=$(minikube service teams-notification-apiserver -n teams-notification-ut --url)
echo "📡 API Server URL: $API_URL"

# 等待 API Server 就緒
echo "⏳ 等待 API Server 就緒..."
for i in {1..30}; do
    if curl -s "$API_URL/health" > /dev/null; then
        echo "✅ API Server 已就緒"
        break
    fi
    echo "⏳ 等待中... ($i/30)"
    sleep 10
done

# 測試健康檢查
echo "🔍 測試健康檢查..."
curl -s "$API_URL/health" | jq '.' || echo "❌ 健康檢查失敗"

# 測試指標端點
echo "📊 測試指標端點..."
curl -s "$API_URL/metrics" | head -10 || echo "❌ 指標端點失敗"

# 測試 API 端點
echo "🚀 測試 API 端點..."

# 測試外部 API
echo "📡 測試外部 API..."
curl -s -X POST "$API_URL/api/v1/notify" \
  -H "Content-Type: application/json" \
  -d '{
    "notifyKey": "test-key",
    "message": "Hello from local test",
    "messageType": "text",
    "priority": "normal",
    "targets": ["all"]
  }' | jq '.' || echo "❌ 外部 API 測試失敗"

# 測試配置端點
echo "⚙️ 測試配置端點..."
curl -s "$API_URL/config" | jq '.' || echo "❌ 配置端點失敗"

echo "✅ API 測試完成！"
