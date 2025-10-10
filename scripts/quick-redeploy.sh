#!/bin/bash

# TeamsNotifyGoV2 快速重構腳本
# 用於代碼更新後的快速重新部署

set -e

# 顏色定義
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

log() {
    echo -e "${BLUE}[$(date +'%H:%M:%S')]${NC} $1"
}

success() {
    echo -e "${GREEN}✅ $1${NC}"
}

echo "🔄 TeamsNotifyGoV2 快速重構"
echo "=========================="
echo ""

# 停止現有服務
log "停止現有服務..."
pkill -f "go run" 2>/dev/null || true
pkill -f "./bin/server" 2>/dev/null || true
docker-compose -f scripts/local-test/docker-compose-local.yml down 2>/dev/null || true

# 重新構建
log "重新構建應用..."
make build

# 重新啟動 Docker 環境
log "重新啟動 Docker 環境..."
./scripts/local-test/start-docker.sh

# 等待服務就緒
log "等待服務就緒..."
sleep 5

# 測試新功能
log "測試新功能..."
if curl -s http://localhost:8080/api/v1/metrics | grep -q "timestamp"; then
    success "系統監控端點正常"
else
    echo "⚠️  系統監控端點可能未完全實現"
fi

success "快速重構完成！"
echo ""
echo "🌐 API Server: http://localhost:8080"
echo "📈 Metrics: http://localhost:8080/api/v1/metrics"
