#!/bin/bash
# scripts/local-test/start-docker.sh

set -e

echo "🐳 啟動本地 Docker 測試環境..."

# 檢查 Docker 是否運行
if ! docker info &> /dev/null; then
    echo "❌ Docker 未運行，請先啟動 Docker"
    exit 1
fi

# 建置映像
echo "🔨 建置應用程式映像..."
docker build -t teams-notification/api-server:local .
## (Worker/Admin 已移除)

# 啟動服務
echo "🚀 啟動服務..."
docker-compose -f scripts/local-test/docker-compose-local.yml up -d

# 等待服務就緒
echo "⏳ 等待服務就緒..."
sleep 30

# 檢查服務狀態
echo "📊 檢查服務狀態..."
docker-compose -f scripts/local-test/docker-compose-local.yml ps

# 檢查健康狀態
echo "🔍 檢查健康狀態..."
curl -s http://localhost:8080/health | jq '.' || echo "❌ API Server 健康檢查失敗"
## (Admin 已移除)

echo "✅ 本地 Docker 測試環境啟動完成！"
echo ""
echo "🌐 服務 URL："
echo "API Server: http://localhost:8080"
## (Admin/Adminer 已移除)
echo ""
echo "📋 可用命令："
echo "  docker-compose -f scripts/local-test/docker-compose-local.yml logs -f    # 查看日誌"
echo "  docker-compose -f scripts/local-test/docker-compose-local.yml down      # 停止服務"
echo "  docker-compose -f scripts/local-test/docker-compose-local.yml restart   # 重啟服務"
