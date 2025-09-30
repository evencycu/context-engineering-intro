# Testing Guide (Current Repository State)

The project does not include Go unit tests yet (`rg --files -g '*_test.go'` returns none). A Bash smoke script is provided to exercise running endpoints once the server and database contain seed data.

## Available Assets
- `expose_doc/test_api.sh` – curl-based smoke checks (requires `curl` + `jq`).

## Prerequisites
1. Start dependencies (`docker compose -f deployments/docker/docker-compose.yml up -d postgres redis`).
2. Apply schema (`make db-migrate`) and seed data as needed (`psql -f scripts/test_data_enhanced.sql`).
3. Launch the server (`go run ./cmd/server`) with the necessary Teams Bot environment variables set (see `README.md`).

## Running the Smoke Script
```bash
./expose_doc/test_api.sh
```
The script first checks `/health`, then optionally exercises CRUD/create flows if you supply real UUIDs via environment variables (see comments inside the script). Tests that depend on missing configuration are skipped rather than failing outright.

### Note on Third-Party Bot tests
Third-Party Bot endpoints are currently disabled (routes removed). Any tests referencing `/api/v1/bots/third-party/...` are intentionally skipped. Platform Bot tests remain enabled and pass.

### Failure Quick Reference
| Step | 常見錯誤 | 原因 | 修復 |
| --- | --- | --- | --- |
| `check_server` | connection refused | 服務未啟動或 port 被佔用 | 啟動 `go run ./cmd/server`，確認 8080 空閒 |
| `get_stats` | `jq` 顯示 0 / `null` | DB 尚未有資料 | 匯入 `scripts/test_data_enhanced.sql` 或手動建立 |
| `create_company` | 409 Conflict | 名稱重複 | 調整名稱或清空舊資料 |
| `create_destination` | 400 Invalid request | 缺少 `PROJECT_ID`/`TENANT_ID` 或 targets 欄位 | 設定環境變數並提供合法 UUID / conversation_id |
| `send_notification` | 500 Failed to send notification | Bot 未設定、Teams 憑證錯誤或未安裝 | 檢查 `teams_bots`、`bot_installations`，觀察 broadcast 日誌 |

## Recommendations
- Add Go unit tests under `internal/api/...` to cover services and handlers with `sqlmock` or a containerized Postgres.
- Extend the smoke script or port it to a Go-based integration test once stable fixtures are defined.
- Integrate `make test` into CI after introducing automated test suites.
