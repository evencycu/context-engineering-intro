#!/bin/bash
# scripts/local-test/deploy-ut.sh

set -e

echo "🚀 部署 UT 環境到本地 minikube..."

# 檢查 minikube 狀態
if ! minikube status &> /dev/null; then
    echo "❌ minikube 未運行，請先執行 ./scripts/local-test/setup-local-env.sh"
    exit 1
fi

# 建立命名空間
echo "📁 建立命名空間..."
kubectl apply -f deployments/ut/namespace.yaml

# 建立 Secrets
echo "🔐 建立 Secrets..."
kubectl apply -f deployments/ut/secrets.yaml

# 建立 ConfigMap
echo "📋 建立 ConfigMap..."
kubectl apply -f deployments/ut/configmap.yaml

# 部署 PostgreSQL
echo "🐘 部署 PostgreSQL..."
kubectl apply -f deployments/ut/postgres.yaml

# 等待 PostgreSQL 就緒
echo "⏳ 等待 PostgreSQL 就緒..."
kubectl wait --for=condition=ready pod -l app=teams-notification-postgres -n teams-notification-ut --timeout=300s

# 部署 Redis
echo "🔴 部署 Redis..."
kubectl apply -f deployments/ut/redis.yaml

# 等待 Redis 就緒
echo "⏳ 等待 Redis 就緒..."
kubectl wait --for=condition=ready pod -l app=teams-notification-redis -n teams-notification-ut --timeout=300s

# 部署 API Server
echo "🚀 部署 API Server..."
kubectl apply -f deployments/ut/api-server.yaml

# 部署 Worker
echo "⚙️ 部署 Worker..."
kubectl apply -f deployments/ut/worker.yaml

# 部署 Admin
echo "👨‍💼 部署 Admin..."
kubectl apply -f deployments/ut/admin.yaml

# 等待所有 Pod 就緒
echo "⏳ 等待所有 Pod 就緒..."
kubectl wait --for=condition=ready pod -l app=teams-notification-api-server -n teams-notification-ut --timeout=300s
kubectl wait --for=condition=ready pod -l app=teams-notification-worker -n teams-notification-ut --timeout=300s
kubectl wait --for=condition=ready pod -l app=teams-notification-admin -n teams-notification-ut --timeout=300s

echo "✅ UT 環境部署完成！"
echo ""
echo "📊 檢查部署狀態："
kubectl get pods -n teams-notification-ut
kubectl get services -n teams-notification-ut

echo ""
echo "🌐 取得服務 URL："
echo "API Server: $(minikube service teams-notification-api-server -n teams-notification-ut --url)"
echo "Admin: $(minikube service teams-notification-admin -n teams-notification-ut --url)"
