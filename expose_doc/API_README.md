# Teams Notification REST API

This guide reflects the current implementation of the Teams notification service (`cmd/server/main.go`). It highlights working endpoints, payload shapes, and areas that are scaffolded but not yet fully wired.

## Base URL
```
http://localhost:8080/api/v1
```

> **Note**: 外部呼叫目前無需 Authorization header。JWT / API Key middleware 僅為 scaffold，尚未連接資料庫或綁定路由。

## Quick Workflow Cheat Sheet
| 任務 | Endpoint | 重點欄位 |
| --- | --- | --- |
| 建立公司 | `POST /companies` | `name`, `contact_email`, `contact_phone`, `address` |
| 建立專案 | `POST /projects` | `company_id`, `notify_key`, `created_by` |
| 建立目的地 | `POST /destinations` | `project_id`, `teams_tenant_id`, `targets` |
| 發送文字通知 | `POST /notifications` | `project_id`, `message_type=text`, `destinations` |
| 發送附件通知 | `POST /notifications` | `message_type=file`, `attachment` |
| 發送 Adaptive Card | `POST /notifications` | `message_type=adaptive_card`, `adaptive_card` |
| 查詢通知 | `GET /notifications/{id}` | 取得狀態 / 錯誤訊息 |
| 重試通知 | `POST /notifications/{id}/retry` | 重新排入佇列 |
| 搜尋目的地 | `GET /destinations/search?team_id=...` | 至少提供一個查詢參數 |
| Bot Webhook | `POST /messages` | 直接接收 Teams 活動 payload |

## Health Checks
```http
GET /health
```
Returns server status, timestamp, version.

## Companies
- `POST /companies` – create company (name, contact info, billing flag).
- `GET /companies` – list with pagination (`limit`, `offset`, `search`, `status`).
- `GET /companies/{id}` – fetch by UUID.
- `PUT /companies/{id}` – update mutable fields.
- `DELETE /companies/{id}` – delete.
- `PATCH /companies/{id}/status` – change status.
- `PATCH /companies/{id}/billing` – enable/disable billing.

### Example
```http
POST /companies
Content-Type: application/json

{
  "name": "Example Corp",
  "contact_email": "ops@example.com",
  "contact_phone": "+886-2-0000-0000",
  "address": "110 台北市信義區信義路五段 7 號",
  "billing_enabled": true
}
```

## Users
- `POST /users` – create user (password hashed inside service layer).
- `GET /users` – pagination + optional `role`/`status` filters.
- `GET /users/{id}` – fetch by UUID.
- `PUT /users/{id}` – update profile/password (re-hashes if provided).
- `DELETE /users/{id}` – delete.
- `PATCH /users/{id}/password` – change password (service validates old password).
- `PATCH /users/{id}/api-key` & `DELETE /users/{id}/api-key` – endpoints exist; persistence methods for API key lifecycle still TODO.
- `GET /users/company/{companyId}`, `GET /users/role/{role}` – list helpers.

## Projects (Notification Keys)
- `POST /projects`
- `GET /projects`
- `GET /projects/{id}`
- `PUT /projects/{id}`
- `DELETE /projects/{id}`
- `PATCH /projects/{id}/limits`
- `GET /projects/company/{companyId}`
- `GET /projects/key/{keyName}`

Fields include `notify_key`, status, daily/monthly limits, priority, and `created_by` user.

## Teams Bots
Only platform bot CRUD is currently wired; third-party routes referenced in older docs were removed in `internal/api/server.go`.

- `POST /bots/platform`
- `GET /bots/platform`
- `GET /bots/platform/{id}`
- `PUT /bots/platform/{id}`
- `DELETE /bots/platform/{id}`
- `PATCH /bots/platform/{id}/status`
- `PATCH /bots/platform/{id}/capabilities`
- `POST /bots/platform/{id}/test`
- `GET /bots/status/{status}` – query by lifecycle state.

## Destinations
Destinations map a project to one or more Teams targets; payloads use the JSONB helper types defined in `internal/database/models.go`.

- `POST /destinations`
- `GET /destinations`
- `GET /destinations/{id}`
- `PUT /destinations/{id}`
- `DELETE /destinations/{id}`
- `PATCH /destinations/{id}/targets`
- `POST /destinations/{id}/validate`
- `GET /destinations/project/{projectId}`
- `GET /destinations/bot/{botId}`
- `GET /destinations/search`

### Target Shape
```json
{
  "type": "channel",           // "personal", "channel", or "groupchat"
  "conversation_id": "19:abc", // required for channel/groupchat targets
  "email": "user@example.com", // optional for personal targets
  "display_name": "Ops Room",  // optional metadata
  "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
}
```
`validateTeamsTargets` currently enforces non-empty `tenant_id` and `conversation_id` (where applicable).

### Create Example
```http
POST /destinations
Content-Type: application/json

{
  "project_id": "<project-uuid>",
  "name": "Ops Alerts",
  "description": "Primary operations notification targets",
  "teams_tenant_id": "<tenant-uuid>",
  "targets": [
    {
      "type": "channel",
      "conversation_id": "19:team@thread.tacv2",
      "display_name": "Ops Channel",
      "tenant_id": "<tenant-uuid>"
    }
  ],
  "bot_id": "<optional-bot-uuid>",
  "created_by": "<user-uuid>"
}
```

## Notifications
- `POST /notifications` – accepts message metadata, persists, and delegates to `BroadcastService`.
- `GET /notifications`
- `GET /notifications/{id}`
- `POST /notifications/{id}/retry`
- `DELETE /notifications/{id}` – cancel pending.
- `GET /notifications/project/{projectId}`
- `GET /notifications/sender/{senderId}`
- `GET /notifications/status/{status}`
- `GET /notifications/date-range?start_date=...&end_date=...`

### Send Example
```http
POST /notifications
Content-Type: application/json

{
  "project_id": "<project-uuid>",
  "sender_id": "<optional-user-uuid>",
  "message_type": "text",
  "content": "系統維護通知：今晚 22:00 進行更新",
  "mentions": ["@值班"],
  "priority": "high",
  "metadata": {
    "maintenance_window": "2024-06-15T22:00:00+08:00"
  },
  "destinations": ["<destination-uuid>"]
}
```

### Adaptive Card Payload
Provide a JSON Adaptive Card per Teams schema:
```json
{
  "message_type": "adaptive_card",
  "adaptive_card": {
    "type": "AdaptiveCard",
    "version": "1.4",
    "body": [
      {
        "type": "TextBlock",
        "text": "Critical Incident",
        "weight": "Bolder",
        "color": "Attention"
      }
    ]
  }
}
```

## Bot Framework Webhook
- `POST /messages` – receives activity payloads; `messagesService` updates `bot_installations`.
- `POST /messages/proactive/test` – drives a stub proactive send for diagnostics.

## Provisioning Shortcut
`internal/api/handlers/provision` exposes a higher-level provisioning workflow combining companies/projects/destinations. Review handler for payload details before use; endpoints remain experimental.

## Error Format
```json
{
  "error": "Error summary",
  "details": "Optional diagnostic information"
}
```
HTTP status codes follow standard semantics (400 validation failure, 401 missing/invalid JWT, 404 not found, 500 internal errors). Rate limiting returns `429`.

## Known Gaps
- API key authentication, billing/usage endpoints, advanced routing, and queue-backed retries are not yet implemented despite database/schema scaffolding.
- Many routes rely on database state (projects, destinations, bots); seed data via `scripts/test_data_enhanced.sql` or custom inserts before exercising them.

Use this document when integrating against the currently checked-in REST implementation. For full schema and example payloads consult `expose_doc/teams-notification-api.yaml` and the handler/service code.
