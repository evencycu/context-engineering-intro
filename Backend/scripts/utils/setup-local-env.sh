#!/bin/bash
# scripts/local-test/setup-local-env.sh

set -e

echo "🚀 建立本地測試環境..."

# 檢查必要工具
echo "🔍 檢查必要工具..."

# 檢查 Docker
if ! command -v docker &> /dev/null; then
    echo "❌ Docker 未安裝，請先安裝 Docker"
    exit 1
fi

# 檢查 kubectl
if ! command -v kubectl &> /dev/null; then
    echo "❌ kubectl 未安裝，請先安裝 kubectl"
    exit 1
fi

# 檢查 minikube
if ! command -v minikube &> /dev/null; then
    echo "❌ minikube 未安裝，請先安裝 minikube"
    exit 1
fi

echo "✅ 所有必要工具已安裝"

# 啟動 minikube
echo "🚀 啟動 minikube..."
minikube start --driver=docker --memory=4096 --cpus=2

# 等待 minikube 就緒
echo "⏳ 等待 minikube 就緒..."
kubectl wait --for=condition=ready node minikube --timeout=300s

# 啟用必要的 addons
echo "🔧 啟用 minikube addons..."
minikube addons enable ingress
minikube addons enable metrics-server

# 檢查 minikube 狀態
echo "📊 檢查 minikube 狀態..."
minikube status

echo "✅ 本地測試環境建立完成！"
echo ""
echo "📋 可用命令："
echo "  minikube dashboard    # 開啟 Kubernetes 儀表板"
echo "  minikube tunnel       # 開啟 LoadBalancer 支援"
echo "  kubectl get nodes     # 檢查節點狀態"
echo "  kubectl get pods -A   # 檢查所有 Pod"
