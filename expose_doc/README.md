# TeamsNotifyGoV2 Exposure Doc

A condensed reference so another coding agent can understand and extend the Microsoft Teams notification service quickly and safely.

## 1. Product & Scope
- Enterprise-grade Teams notification API with pre-approved destinations.
- Core entities: companies, users, projects (notification keys), Teams bots, destinations, notifications, bot installations.
- Delivery targets: personal chats, channels, group chats; scheduling/analytics planned but not yet implemented.

## 2. Implementation Snapshot
| 功能 | 對應程式 | 現況 | 備註 |
| --- | --- | --- | --- |
| 健康檢查 | `internal/api/server.go` | ✅ | `/health` 回傳 status/version/timestamp |
| 公司 / 使用者 / 專案 CRUD | `internal/api/handlers/{companies,users,projects}` | ✅ | Pagination、狀態與 billing 切換 |
| Teams Bot 管理 | `internal/api/handlers/bots` | ✅ | 僅平台 Bot；第三方 Bot 尚未開放 |
| 目的地管理 | `internal/api/handlers/destinations` | ✅ | Targets 驗證需 `tenant_id` + `conversation_id` |
| 通知發送 | `internal/api/handlers/notifications` | ✅ | 呼叫 `BroadcastService`，回傳 202 Accepted |
| Bot Webhook | `internal/api/handlers/messages` | ✅ | 安裝事件寫入 `bot_installations` |
| 一鍵佈建 | `internal/api/handlers/provision` | ⚠️ Beta | 需求波動大，仍在調整 |
| 佇列 / 排程 / 計費 API | N/A | ⏳ | Schema / docs 有規劃，程式尚未落地 |
| 對外驗證 (JWT / API Key) | `internal/api/middleware/auth.go` | ⏳ | 中介層 scaffold，尚未強制套用；API 目前無需 Authorization |

## 3. Core Requirements (docs/Requirement.md)
- 通过 Azure Bot Framework 向已核准的 Teams 目的地發送通知。
- 以公司 / 專案 / 發送者追蹤用量並保留稽核紀錄。
- 支援廣播、多種訊息格式、@ mention、重試與狀態查詢。
- 透過 Redis / 佇列（規劃中）緩解 Rate Limit；現階段使用記憶體限制器。
- 提供健康檢查、結構化日誌與基本安全標頭。

## 4. Tech Stack & Key Libraries
- Go 1.24、Gin (`github.com/gin-gonic/gin`)、`sqlx`、`logrus`。
- Middle-layer：`gin-contrib/requestid`、`gin-contrib/cors`、`golang.org/x/time/rate`。
- JSON/validator：Sonic、go-json、go-playground validator。
- Tooling：Makefile、shell scripts、Docker Compose、proto/OpenAPI 草稿。

## 5. Repository Map
```
cmd/server          // 服務進入點
internal/api        // server, handlers, middleware, services, repositories
internal/database   // models + schema.sql + usage docs
api/openapi         // REST 規格 (同步於 expose_doc)
api/proto           // gRPC 草案
configs/            // YAML 設定
scripts/            // build/lint/docker/openapi/sql 工具
examples/           // 範例程式與資料
deployments/        // docker-compose / K8s 配置
expose_doc/         // 本文件集
```

## 6. Runtime Components
- `cmd/server/main.go`: 注入 config、建立 middleware、連接 Postgres、註冊所有 handlers。
- `internal/api/server.go`: 建立 Gin router、套用 CORS/RequestID/RateLimiter、處理健康檢查與優雅關閉。
- `internal/api/middleware`: 日誌、錯誤處理、安全標頭、Rate limit 記錄；Auth middleware 僅作為範本，預設不啟用。
- `handlers/*`: 負責 JSON 驗證、呼叫 service、輸出包含 pagination 的回應。
- `services/service.go`: 實作主要領域邏輯（密碼雜湊、API key scaffold、Teams 活動處理、通知傳播等）。
- `services/broadcast_service.go`: 針對每個 Bot Installation 取得 connector token → POST Teams 訊息；尚缺重試/backoff 策略。
- `repositories/repository.go`: 使用 `sqlx` 的 CRUD 實作，涵蓋公司、使用者、專案、Bot、目的地、通知。
- `internal/database/models.go`: Entity 定義與 JSONB Value/Scan helper。

### Broadcast Flow (文字序列圖)
```
Client -> Notifications Handler: POST /api/v1/notifications
Handler -> Notification Service: 驗證 + 寫入 notifications
Notification Service -> Broadcast Service: SendToDestinations(destIDs)
Broadcast Service -> TeamsBot Repo: 取得 Bot 憑證
Broadcast Service -> Installation Repo: 讀取 active installations
Broadcast Service -> Bot Framework: POST /v3/conversations/... (逐一)
Bot Framework -> Broadcast Service: 回傳 HTTP 狀態碼
Broadcast Service -> Installation Repo: UpdateActivity / MarkAsStale
Broadcast Service -> Notification Service: 返回摘要結果
Notification Service -> Client: 202 Accepted + notification_id
```

## 7. Data Model Snapshot
重點欄位與關聯詳見 `expose_doc/ERD_Diagram.md`：公司 → 使用者 / 專案 → 目的地 → 通知，多對多關係透過 `notification_destinations` 表達；Bot 與 installations 提供 Teams 會話資訊。

## 8. Request Lifecycle（Send Notification）
1. Client 以 `POST /api/v1/notifications` 送出請求。
2. Gin middleware（RequestID、RateLimiter、Logging）包裹請求。
3. Handler 做 JSON 驗證並呼叫 NotificationService。
4. NotificationService 建立紀錄，轉給 BroadcastService。
5. BroadcastService 查詢目的地、Bot 安裝資訊，與 Teams API 互動。
6. Handler 回傳 `202 Accepted` 與 notification ID。
7. 目前無背景 worker，後續狀態更新需透過 Teams 回調或人工。

## 9. API Surfaces & Workflows
- REST：詳見 `expose_doc/API_README.md`，列出所有 `/api/v1` 路由與範例。
- OpenAPI 3：`expose_doc/teams-notification-api.yaml` 精準反映現況（不包含計費/排程）。
- Protobuf（草案）：`api/proto/notification.proto` 描述未來 gRPC 介面。

### Workflow 快速參考
| 任務 | 路徑 | 重點 |
| --- | --- | --- |
| 建立公司 | `POST /api/v1/companies` | 名稱/聯絡資訊必填 |
| 建立專案 | `POST /api/v1/projects` | `company_id`、`notify_key`、`created_by` |
| 建立目的地 | `POST /api/v1/destinations` | Targets 需 `tenant_id`+`conversation_id` |
| 發送文字通知 | `POST /api/v1/notifications` | `message_type=text` |
| 發送附件通知 | 同上 | `message_type=file` + `attachment` 物件 |
| 發送 Adaptive Card | 同上 | `message_type=adaptive_card` + card JSON |
| 查詢通知 | `GET /api/v1/notifications/{id}` | 取得狀態/錯誤訊息 |
| 重試通知 | `POST /api/v1/notifications/{id}/retry` | 重新觸發失敗紀錄 |
| 搜尋目的地 | `GET /api/v1/destinations/search?team_id=...` | 至少提供一個查詢參數 |
| Bot Webhook | `POST /api/v1/messages` | 直接接收 Teams 活動 payload |

## 10. Configuration & Secrets
- `.env` 需要 `DATABASE_URL`、`TEAMS_BOT_APP_ID`、`TEAMS_BOT_APP_PASSWORD`、`TEAMS_TENANT_ID` 等。
- `configs/*.yaml` 提供 server/database/cache/feature flag 預設值。
- 外部 API 無驗證；若要加入需自行啟用 Auth middleware 並儲存金鑰/Token。

## 11. Testing & QA
- 尚無 Go 單元或整合測試。
- `expose_doc/test_api.sh`：環境可變的 smoke script（可設定 `PROJECT_ID`、`CREATED_BY` 等）。
- 詳細步驟與建議參見 `expose_doc/TEST_README.md`。

### Smoke 測試對照表
| Script Step | 依賴 | 常見錯誤 | 修正建議 |
| --- | --- | --- | --- |
| `check_server` | `go run ./cmd/server` | 無法連線 | 確認服務啟動、port 未被占用 |
| `get_stats` | DB 內容 | `jq` 顯示 0 | 匯入測試資料或先建立基礎資料 |
| `create_company` | 無 | 409 Conflict | 更換名稱或清除舊資料 |
| `create_destination` | `PROJECT_ID`、`CREATED_BY`、`TENANT_ID` | 驗證失敗 | 確保帶入合法 UUID + conversation_id |
| `send_notification` | 目的地 + Teams Bot | 500 失敗 | 檢查 Bot 安裝、觀察 `broadcast_service` 日誌 |

## 12. 範例資料與種子
- `scripts/test_data_enhanced.sql` 匯入公司/專案/目的地/通知樣本，方便測試。
- 最小步驟：
  1. `POST /companies` → `POST /projects`（填入公司 ID） → `POST /destinations`。
  2. 從資料庫查詢目的地與 Bot 安裝 UUID，供 `POST /notifications` 使用。
- 建議在 CI 或 README 維護常用 UUID 對照表，以利自動化代理生成範例。

## 13. Build, Tooling & Deployment
- `Makefile`：`docker-up`, `build`, `run`, `test`, `db-migrate`, `openapi-*`。
- `scripts/build.sh`：多平台編譯、Docker build、可選 proto 產生。
- `deployments/`：Docker Compose + K8s/infra 腳本。
- `scripts/openapi.sh`：啟動本地 Swagger UI。

## 14. Observability & Security
- 日誌：Logrus JSON formatter + middleware 輸出請求/回應/錯誤。
- Rate limit：`golang.org/x/time/rate` 全域節流；未依租戶/專案細分。
- 健康檢查：`GET /health`。
- 服務尚未啟用任何外部認證，部署時務必補強。

## 15. Known Gaps
- API Key / JWT 驗證僅止於 scaffold，未實際綁定任何路由。
- Queue、排程、計費、分析尚未實作。
- Broadcast 缺乏批次重試、退避與細緻錯誤分類。
- 沒有自動化測試與 CI；Smoke script 以人工設定為主。

## 16. Quick Start（Local）
1. `docker compose -f deployments/docker/docker-compose.yml up -d postgres redis`
2. `make db-migrate`
3. 設定 `TEAMS_BOT_*`、`DATABASE_URL` 等環境變數。
4. `go run ./cmd/server`
5. `curl http://localhost:8080/health`
6. `./expose_doc/test_api.sh`（視需要傳入 `PROJECT_ID` 等環境變數）。

## 17. Tips for Future Contributors / Coding Agents
- 延伸新資源時沿用 handler → service → repository 分層模式。
- 操作 JSONB 欄位需使用現有 Value/Scan helper，避免自訂編碼錯誤。
- 修改 REST 面時同步更新 `expose_doc/API_README.md` 與 OpenAPI。
- 補上測試與資料種子有助於 CI 與自動化代理。

---
_Last updated: 2025-09-24T14:22:38Z_
