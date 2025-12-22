# 🧑‍💻 Development Guide

## 1. 文件資訊
- **版本**：v1.0  
- **作者**：Tech Lead / Developer  
- **最後更新**：2025-10-08  

---

## 2. 目的
此文件提供開發者進入 Teams Notification API 專案的快速指南，包括環境設定、分支策略、代碼規範與測試準則。

---

## 3. 專案概述

### 3.1 系統架構
Teams Notification API 是一個基於 Go 的微服務架構，提供企業級的 Teams 通知服務。系統採用 Actor 模式處理異步通知，使用 Redis 作為隊列和快取層，PostgreSQL 作為持久化存儲。

### 3.2 核心功能
- **通知發送**：支援 Teams 個人、群組、頻道通知
- **佇列管理**：Redis 驅動的佇列系統與重試機制
- **專案管理**：多租戶專案與權限控制
- **Bot 管理**：Teams Bot 安裝與狀態管理
- **監控告警**：健康檢查與效能監控

---

## 4. 開發環境設定

### 4.1 需求工具
| 類別 | 工具 | 版本要求 |
|------|------|----------|
| 程式語言 | Go | 1.24+ |
| 資料庫 | PostgreSQL | 15+ |
| 快取/佇列 | Redis | 7+ |
| 容器 | Docker / Docker Compose | 最新版 |
| IDE | VSCode | 建議使用 devcontainer |
| 測試工具 | Go test, k6 | 效能測試 |

### 4.2 快速啟動
```bash
# Clone 專案
git clone https://github.com/evencycu/TeamsNotifyGoV3.git
cd TeamsNotifyGoV2

# 啟動開發環境
make docker-up

# 執行資料庫遷移
make db-migrate

# 建置並執行
make build
make run

# 執行測試
make test
```

### 4.3 環境變數設定
```bash
# 必要環境變數
export TEAMS_BOT_APP_ID="your-app-id"
export TEAMS_TENANT_ID="your-tenant-id"
export TEAMS_BOT_APP_PASSWORD="your-password"

# 資料庫連線
export DATABASE_URL="postgresql://teamsnotify:teamsnotify123@localhost:5432/teamsnotify?sslmode=disable"
export REDIS_URL="redis://localhost:6379"

# 日誌等級
export LOG_LEVEL="debug"
```

---

## 5. Git Workflow

```mermaid
flowchart LR
A[feature/login] --> B[Pull Request]
B --> C[develop]
C --> D[release/x.y.z]
D --> E[main]
```

### 5.1 分支策略
- 使用 **GitHub Flow / Trunk-Based Development**
- 分支命名規則：  
  - 功能：`feature/<name>`  
  - 修正：`fix/<name>`  
  - 發佈：`release/<version>`  
- 所有 PR 需通過自動測試與 Code Review

### 5.2 提交規範
```bash
# 提交訊息格式
<type>(<scope>): <description>

# 範例
feat(api): add external notification endpoint
fix(queue): resolve retry mechanism issue
docs(readme): update installation guide
```

---

## 6. 程式風格與規範

### 6.1 Go 程式規範
- 使用 `golangci-lint` 驗證代碼品質
- 命名遵循 Go Conventions（短、簡潔、語意清楚）
- 模組目錄結構：
  ```bash
  /cmd/server/          # 服務入口
  /internal/           # 封裝邏輯
    /api/              # API 層
    /services/         # 業務邏輯
    /database/         # 資料庫層
    /actor/            # Actor 模式
  /api/                # OpenAPI 定義
  /configs/            # 配置檔案
  /scripts/            # 腳本檔案
  ```

### 6.2 代碼品質要求
- 單元測試覆蓋率 ≥ 80%
- 使用 `go fmt` 格式化代碼
- 使用 `go vet` 檢查潛在問題
- 遵循 Go 最佳實踐

---

## 7. 測試與驗證

### 7.1 測試策略
- **單元測試**：使用 `go test` 測試個別函數
- **整合測試**：測試 API 端點與資料庫整合
- **效能測試**：使用 k6 進行負載測試
- **E2E 測試**：完整業務流程測試

### 7.2 測試執行
```bash
# 執行所有測試
make test

# 執行特定測試
go test ./internal/api/handlers/...

# 執行效能測試
k6 run scripts/performance_test.js

# 檢查測試覆蓋率
go test -cover ./...
```

### 7.3 測試環境
- **開發環境**：本地 Docker 容器
- **測試環境**：獨立的測試資料庫和 Redis
- **預生產環境**：模擬生產環境配置

---

## 8. 開發環境參數

### 8.1 必要環境變數
| 變數 | 說明 | 範例 |
|------|------|------|
| `TEAMS_BOT_APP_ID` | Teams Bot 應用程式 ID | `844146d7-4ac9-4e4d-a463-d6e027714e81` |
| `TEAMS_TENANT_ID` | Teams 租戶 ID | `051cece0-e4dc-4aed-b471-bf29824e1ee6` |
| `TEAMS_BOT_APP_PASSWORD` | Teams Bot 密碼 | `your-password` |
| `DATABASE_URL` | PostgreSQL 連線字串 | `postgresql://user:pass@localhost:5432/db` |
| `REDIS_URL` | Redis 伺服器位址 | `redis://localhost:6379` |

### 8.2 可選環境變數
| 變數 | 說明 | 預設值 |
|------|------|--------|
| `LOG_LEVEL` | 日誌等級 | `info` |
| `PORT` | 服務端口 | `8080` |
| `ENVIRONMENT` | 環境名稱 | `development` |

---

## 9. Debug 與日誌

### 9.1 日誌設定
```bash
# 啟用 Debug 模式
export LOG_LEVEL=debug

# 查看服務日誌
make docker-logs

# 查看特定服務日誌
docker logs teamsnotify-postgres
docker logs teamsnotify-redis
```

### 9.2 常見錯誤排查
- **RateLimit**：Teams API 速率限制
- **TokenExpired**：Teams Token 過期
- **RedisTimeout**：Redis 連線超時
- **DatabaseError**：資料庫連線錯誤

### 9.3 監控工具
- 使用 `Grafana Loki` 收集集中式日誌
- 內建健康檢查端點：`/health`
- 佇列狀態監控：`/api/v1/queue/status`

---

## 10. CI/CD Pipeline

### 10.1 建置流程
```mermaid
flowchart LR
A["Commit to main"] --> B["GitHub Actions Build"]
B --> C["Run Unit Tests"]
C --> D["Docker Build & Push (ACR)"]
D --> E["Deploy to Azure Container Apps"]
E --> F["Run Smoke Tests"]
F --> G["Notify Teams Channel"]
```

### 10.2 部署階段
| 階段 | 工具 | 動作 |
|------|------|------|
| Build | GitHub Actions | 自動建置 Docker Image |
| Test | GitHub Actions | 單元與整合測試 |
| Deploy | Azure CLI / Bicep | 自動部署 ACA / DB Migration |

### 10.3 環境配置
- **開發環境**：Docker Compose
- **測試環境**：Kubernetes
- **生產環境**：Azure Container Apps

---

## 11. 專案結構詳解

### 11.1 核心模組
```
internal/
├── api/                    # API 層
│   ├── handlers/          # HTTP 處理器
│   ├── middleware/        # 中間件
│   ├── services/          # 業務邏輯
│   └── repositories/      # 資料存取層
├── database/              # 資料庫相關
│   ├── models.go         # 資料模型
│   ├── schema.sql        # 資料庫結構
│   └── init.sql          # 初始化資料
├── actor/                 # Actor 模式
│   └── notification_actor.go
└── queue/                 # 佇列系統
    ├── manager.go
    ├── circuit_breaker.go
    └── retry_policy.go
```

### 11.2 配置檔案
```
configs/
├── common.yaml           # 通用配置
├── apiserver.yaml       # API 服務配置
├── worker.yaml           # Worker 配置
└── admin.yaml            # 管理員配置
```

---

## 12. 開發最佳實踐

### 12.1 代碼品質
- [ ] 命名合理、可讀性佳  
- [ ] 無硬編碼（Hardcoded value）  
- [ ] 單元測試覆蓋率 ≥ 80%  
- [ ] 日誌具備 traceId / correlationId  
- [ ] 錯誤處理明確  

### 12.2 效能優化
- 使用 Redis 快取減少資料庫查詢
- 實作連線池管理資料庫連線
- 使用批次處理提高吞吐量
- 監控記憶體使用和 GC 壓力

### 12.3 安全考量
- 驗證所有輸入參數
- 使用 HTTPS 加密通訊
- 實作速率限制防止濫用
- 定期更新依賴套件

---

## 13. 故障排除

### 13.1 常見問題
1. **服務無法啟動**：檢查環境變數和端口占用
2. **資料庫連線失敗**：確認 PostgreSQL 容器運行狀態
3. **Redis 連線超時**：檢查 Redis 容器和網路配置
4. **Teams API 認證失敗**：驗證 Bot 憑證和權限設定

### 13.2 除錯工具
```bash
# 檢查服務狀態
curl http://localhost:8080/health

# 檢查佇列狀態
curl http://localhost:8080/api/v1/queue/status

# 查看資料庫連線
docker exec teamsnotify-postgres psql -U teamsnotify -d notification_center -c "SELECT 1;"

# 進入資料庫 Shell
make db-shell

# 查看 Redis 連線
docker exec teamsnotify-redis redis-cli ping
```

---

## 14. 參考資源

### 14.1 技術文檔
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Microsoft REST API Guidelines](https://github.com/microsoft/api-guidelines)
- [Twelve-Factor App](https://12factor.net/)

### 14.2 專案文檔
- [系統架構文檔](../02_ARCHITECTURE/Architecture.md)
- [系統設計文檔](../03_DESIGN/SystemDesign.md)
- [測試計劃](../04_TEST/TestPlan.md)
- [部署指南](../05_DEPLOYMENT/DeploymentGuide.md)
- [使用者手冊](../06_USER_GUIDE/UserManual.md)

### 14.3 外部資源
- [Go 官方文檔](https://golang.org/doc/)
- [Gin 框架文檔](https://gin-gonic.com/docs/)
- [PostgreSQL 文檔](https://www.postgresql.org/docs/)
- [Redis 文檔](https://redis.io/documentation)
- [Microsoft Teams Bot Framework](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/what-are-bots)
