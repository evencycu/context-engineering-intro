# 🚀 部署模式比較指南

## 📋 概述

TeamsNotifyGoV2 支持三種不同的部署模式，每種模式都有其特定的使用場景和優缺點。

---

## 🎯 部署模式對比

| 特性 | Local Process | Local Docker | Persistent Docker |
|------|---------------|--------------|-------------------|
| **啟動速度** | ⚡ 最快 | 🐌 中等 | 🐌 中等 |
| **資源使用** | 💚 最低 | 🟡 中等 | 🟡 中等 |
| **環境隔離** | ❌ 無 | ✅ 完全隔離 | ✅ 完全隔離 |
| **數據持久化** | ❌ 無 | ❌ 無 | ✅ 完全持久化 |
| **調試便利性** | ✅ 最便利 | 🟡 中等 | 🟡 中等 |
| **生產相似性** | ❌ 低 | ✅ 高 | ✅ 高 |
| **適用場景** | 快速開發 | 功能測試 | 長期開發 |

---

## 🚀 Local Process 模式

### 特點
- **最快啟動**：直接運行 Go 進程，無 Docker 開銷
- **最佳調試體驗**：可以直接使用 IDE 調試器
- **最低資源使用**：不需要 Docker 容器

### 使用場景
- 🔥 **快速開發**：代碼修改後立即測試
- 🐛 **調試問題**：需要詳細的調試信息
- ⚡ **性能測試**：測試真實性能表現

### 命令
```bash
# 基本部署
./scripts/smart-deploy.sh local-process

# 帶遷移和測試
./scripts/smart-deploy.sh local-process --migrate --test

# 使用 Makefile
make deploy-local
```

### 停止方式
```bash
pkill -f "go run"
# 或
pkill -f "./bin/server"
```

---

## 🐳 Local Docker 模式

### 特點
- **環境隔離**：完全隔離的容器環境
- **生產相似**：與生產環境配置一致
- **快速重置**：每次重新部署都是全新環境

### 使用場景
- 🧪 **功能測試**：測試完整的功能流程
- 🔄 **CI/CD 模擬**：模擬持續集成環境
- 🐛 **隔離調試**：避免本地環境干擾

### 命令
```bash
# 基本部署
./scripts/smart-deploy.sh local-docker

# 清理舊映像並測試
./scripts/smart-deploy.sh local-docker --clean --test

# 使用 Makefile
make deploy-docker
```

### 停止方式
```bash
docker-compose -f scripts/local-test/docker-compose-local.yml down
```

---

## 💾 Persistent Docker 模式

### 特點
- **數據持久化**：PostgreSQL 和 Redis 數據不會丟失
- **長期開發**：適合需要保持數據的長期開發
- **生產相似**：最接近生產環境的配置

### 使用場景
- 📊 **數據分析**：需要保持測試數據
- 🔄 **長期開發**：多天開發週期
- 🧪 **集成測試**：測試數據一致性

### 命令
```bash
# 基本部署
./scripts/smart-deploy.sh persistent-docker

# 帶遷移和測試
./scripts/smart-deploy.sh persistent-docker --migrate --test

# 使用 Makefile
make deploy-persistent
```

### 管理命令
```bash
# 啟動持久化服務
make persistent-up

# 停止持久化服務
make persistent-down

# 查看日誌
make persistent-logs

# 清理所有數據
make persistent-clean
```

### 停止方式
```bash
# 停止應用
pkill -f "go run"

# 停止數據庫服務
make persistent-down
```

---

## 🎯 選擇指南

### 快速開發階段
```bash
# 推薦：Local Process
./scripts/smart-deploy.sh local-process
```

### 功能測試階段
```bash
# 推薦：Local Docker
./scripts/smart-deploy.sh local-docker --test
```

### 長期開發階段
```bash
# 推薦：Persistent Docker
./scripts/smart-deploy.sh persistent-docker --migrate
```

### 生產部署前
```bash
# 推薦：Persistent Docker + 完整測試
./scripts/smart-deploy.sh persistent-docker --clean --migrate --test
```

---

## 🔧 高級用法

### 組合使用
```bash
# 開發時使用 Local Process
./scripts/smart-deploy.sh local-process

# 測試時切換到 Docker
./scripts/smart-deploy.sh local-docker --test

# 長期開發時使用 Persistent
./scripts/smart-deploy.sh persistent-docker --migrate
```

### 數據遷移
```bash
# 從 Local Process 到 Persistent Docker
# 1. 導出 Local Process 數據
# 2. 啟動 Persistent Docker
# 3. 導入數據
```

### 性能優化
```bash
# Local Process：最佳性能
./scripts/smart-deploy.sh local-process

# Docker：接近生產性能
./scripts/smart-deploy.sh local-docker --clean
```

---

## 📊 性能對比

| 指標 | Local Process | Local Docker | Persistent Docker |
|------|---------------|--------------|-------------------|
| **啟動時間** | ~2s | ~15s | ~20s |
| **記憶體使用** | ~50MB | ~200MB | ~200MB |
| **CPU 使用** | 最低 | 中等 | 中等 |
| **I/O 性能** | 最佳 | 良好 | 良好 |
| **調試便利性** | 最佳 | 良好 | 良好 |

---

## 🚨 注意事項

### Local Process 模式
- ⚠️ 需要本地安裝 PostgreSQL 和 Redis
- ⚠️ 可能與其他服務衝突
- ⚠️ 環境配置可能與生產不同

### Local Docker 模式
- ⚠️ 每次重啟都會丟失數據
- ⚠️ 需要 Docker 環境
- ⚠️ 調試相對困難

### Persistent Docker 模式
- ⚠️ 數據會一直保留，需要定期清理
- ⚠️ 需要 Docker 環境
- ⚠️ 磁盤空間使用較多

---

## 🔄 模式切換

### 從 Local Process 到 Docker
```bash
# 停止 Local Process
pkill -f "go run"

# 啟動 Docker 模式
./scripts/smart-deploy.sh local-docker
```

### 從 Docker 到 Persistent
```bash
# 停止 Docker
docker-compose -f scripts/local-test/docker-compose-local.yml down

# 啟動 Persistent
./scripts/smart-deploy.sh persistent-docker
```

### 清理切換
```bash
# 完全清理並重新開始
make docker-down
make persistent-clean
./scripts/smart-deploy.sh local-process
```

---

**建立時間**: 2025-10-10  
**維護者**: AI Assistant  
**狀態**: ✅ 完成
