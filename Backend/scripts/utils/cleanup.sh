#!/bin/bash
# scripts/local-test/cleanup.sh

set -e

echo "🧹 清理本地測試環境..."

# 刪除 UT 環境
echo "🗑️ 刪除 UT 環境..."
kubectl delete namespace teams-notification-ut --ignore-not-found=true

# 刪除 UAT 環境
echo "🗑️ 刪除 UAT 環境..."
kubectl delete namespace teams-notification-uat --ignore-not-found=true

# 刪除 PROD 環境
echo "🗑️ 刪除 PROD 環境..."
kubectl delete namespace teams-notification-prod --ignore-not-found=true

# 清理 minikube
echo "🧹 清理 minikube..."
minikube delete --all

echo "✅ 本地測試環境清理完成！"
