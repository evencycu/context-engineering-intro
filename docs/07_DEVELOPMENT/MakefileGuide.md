# 📋 Makefile 指令使用指南

## 1. 文件資訊

- **版本**: v1.0.0
- **撰寫人**: Development Team
- **最後更新**: 2025-10-22
- **目的**: 統一使用 Makefile 管理所有專案指令

---

## 2. 指令統一原則

### 2.1 統一管理

所有專案相關的指令都通過 Makefile 統一管理，不再直接調用腳本：

```bash
# ❌ 舊方式 (不推薦)
./scripts/deploy/docker/deploy.sh restart
./scripts/test/api/simple_*.sh

# ✅ 新方式 (推薦)
make deploy-restart
make test-api
```

### 2.2 指令分類

#### 建構與運行
```bash
make build          # 建構應用程式
make run            # 運行應用程式
make test           # 運行測試
make clean          # 清理建構檔案
make prod-build     # 生產環境建構
```

#### 部署指令 (整合自 deploy.sh)
```bash
make deploy         # 完整部署
make deploy-dev     # 開發部署
make deploy-local   # 本地部署
make deploy-quick   # 快速部署
make deploy-docker  # Docker 部署
make deploy-server  # 僅部署 server
make deploy-openapi # 僅部署 OpenAPI
make deploy-stop    # 停止所有服務
make deploy-restart # 重啟所有服務
make deploy-status  # 檢查服務狀態
make deploy-logs    # 顯示服務日誌
```

#### OpenAPI 管理
```bash
make openapi-start  # 啟動 OpenAPI 服務
make openapi-stop   # 停止 OpenAPI 服務
make openapi-status # 檢查 OpenAPI 狀態
make openapi-open   # 開啟 Swagger UI
```

#### 測試指令
```bash
make test-e2e       # 端到端測試
make test-load      # 負載測試
make test-api       # API 測試
make test-report    # 生成測試報告
```

---

## 3. 常用工作流程

### 3.1 開發環境設置

```bash
# 1. 設置開發環境
make dev-setup

# 2. 啟動 Docker 服務
make docker-up

# 3. 建構應用程式
make build

# 4. 運行應用程式
make run
```

### 3.2 部署流程

```bash
# 1. 完整部署
make deploy

# 2. 檢查服務狀態
make deploy-status

# 3. 查看服務日誌
make deploy-logs
```

### 3.3 測試流程

```bash
# 1. 運行單元測試
make test

# 2. 運行 API 測試
make test-api

# 3. 運行 E2E 測試
make test-e2e

# 4. 生成測試報告
make test-report
```

### 3.4 維護流程

```bash
# 1. 重啟服務
make deploy-restart

# 2. 檢查狀態
make deploy-status

# 3. 查看日誌
make deploy-logs

# 4. 清理建構檔案
make clean
```

---

## 4. 指令對應表

| 舊指令 | 新指令 | 說明 |
|--------|--------|------|
| `./scripts/deploy/docker/deploy.sh restart` | `make deploy-restart` | 重啟所有服務 |
| `./scripts/deploy/docker/deploy.sh status` | `make deploy-status` | 檢查服務狀態 |
| `./scripts/deploy/docker/deploy.sh logs` | `make deploy-logs` | 顯示服務日誌 |
| `./scripts/test/api/simple_*.sh` | `make test-api` | API 測試 |
| `./scripts/openapi.sh start` | `make openapi-start` | 啟動 OpenAPI |
| `./scripts/openapi.sh stop` | `make openapi-stop` | 停止 OpenAPI |

---

## 5. 最佳實踐

### 5.1 指令使用

1. **優先使用 Makefile**: 所有操作都通過 `make` 命令執行
2. **查看可用指令**: 使用 `make help` 查看所有可用指令
3. **保持一致性**: 團隊成員都使用相同的 Makefile 指令

### 5.2 開發流程

1. **開發前**: `make dev-setup` 設置環境
2. **開發中**: `make build` 建構，`make test` 測試
3. **部署前**: `make deploy-status` 檢查狀態
4. **部署後**: `make deploy-logs` 查看日誌

### 5.3 故障排除

1. **服務問題**: `make deploy-status` 檢查狀態
2. **日誌查看**: `make deploy-logs` 查看詳細日誌
3. **重啟服務**: `make deploy-restart` 重啟所有服務
4. **清理環境**: `make clean` 清理建構檔案

---

## 6. 注意事項

### 6.1 向後兼容

- 保留了一些 legacy 命令以確保向後兼容
- 新的開發應該優先使用 Makefile 指令
- 舊的腳本調用方式仍然可用，但不推薦

### 6.2 錯誤處理

- 如果 Makefile 指令失敗，檢查對應的腳本是否存在
- 使用 `make help` 查看所有可用指令
- 檢查 Makefile 中的目標定義是否正確

### 6.3 擴展指令

- 新增指令時，請更新 Makefile 和此文檔
- 保持指令命名的一致性
- 添加適當的說明和幫助信息

---

## 7. 相關資源

- **Makefile**: 專案根目錄的 Makefile
- **部署腳本**: `scripts/deploy/docker/deploy.sh`
- **測試腳本**: `scripts/test/` 目錄下的各種測試腳本
- **OpenAPI 腳本**: `scripts/openapi.sh`

---

*最後更新: 2025-10-22*
