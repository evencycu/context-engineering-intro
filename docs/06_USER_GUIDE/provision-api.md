# Provision API 使用範例

Provision API 提供一次性建立公司、專案和目的地的功能，簡化初始設定流程。

> **注意**: 此文件與 OpenAPI 規範保持同步，最新的 API 定義請參考 `api/openapi/teams-notification-api.yaml`

## API 端點

- `POST /api/v1/provision` - 建立 provision bundle
- `GET /api/v1/provision/{notify_key}` - 取得 provision bundle
- `PUT /api/v1/provision/{notify_key}` - 更新 provision bundle
- `POST /api/v1/provision/{notify_key}/enable` - 啟用專案
- `POST /api/v1/provision/{notify_key}/disable` - 停用專案

## 基本使用範例

### 1. 建立新的 Provision Bundle

```bash
curl -X POST http://localhost:8080/api/v1/provision \
  -H "Content-Type: application/json" \
  -d '{
    "company_id": "e4f160f4-4917-443d-86e9-29071967b768",
    "created_by": "11111111-1111-1111-1111-111111111111",
    "project_name": "測試專案",
    "project_description": "這是一個使用 provision API 建立的測試專案",
    "teams_tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
    "targets": [
      {
        "type": "channel",
        "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
        "conversation_id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2"
      },
      {
        "type": "personal",
        "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
        "conversation_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR"
      }
    ]
  }'
```

**回應範例：**
```json
{
  "data": {
    "notify_key": "5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38",
    "project": {
      "id": "f6a110f8-27f9-477e-a651-6f8cdbc6dca5",
      "company_id": "e4f160f4-4917-443d-86e9-29071967b768",
      "notify_key": "5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38",
      "description": "這是一個使用 provision API 建立的測試專案",
      "status": "active",
      "daily_limit": 10000,
      "monthly_limit": 300000,
      "priority": "normal",
      "created_by": "11111111-1111-1111-1111-111111111111"
    },
    "destination": {
      "id": "cf867660-26ae-4437-828c-6fb2c2e3eb19",
      "project_id": "f6a110f8-27f9-477e-a651-6f8cdbc6dca5",
      "name": "測試專案-default",
      "description": "Auto-provisioned default destination",
      "teams_tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
      "targets": [
        {
          "type": "channel",
          "conversation_id": "19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2",
          "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
        },
        {
          "type": "personal",
          "conversation_id": "a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR",
          "tenant_id": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
        }
      ],
      "status": "active",
      "validation_status": "pending"
    }
  }
}
```

### 2. 取得 Provision Bundle

```bash
curl -s http://localhost:8080/api/v1/provision/5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38 | jq .
```

### 3. 啟用專案

```bash
curl -X POST http://localhost:8080/api/v1/provision/5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38/enable
```

### 4. 停用專案

```bash
curl -X POST http://localhost:8080/api/v1/provision/5984f00fe2c5007ea0edb4e9b6269a4abf304e71070ee05ef05575bfacda5e38/disable
```

## 使用新建立的專案發送通知

```bash
curl -X POST http://localhost:8080/api/v1/notifications \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer test-token" \
  -d '{
    "project_id": "f6a110f8-27f9-477e-a651-6f8cdbc6dca5",
    "message_type": "text",
    "content": "這是一個使用 provision API 建立的專案發送的測試通知",
    "priority": "normal",
    "targets": ["cf867660-26ae-4437-828c-6fb2c2e3eb19"]
  }'
```

## 請求參數說明

### ProvisionCreateRequest

| 欄位 | 類型 | 必填 | 說明 |
|------|------|------|------|
| `company_id` | UUID | 是 | 公司 ID |
| `created_by` | UUID | 是 | 建立者用戶 ID |
| `project_name` | string | 是 | 專案名稱 |
| `project_description` | string | 是 | 專案描述 |
| `teams_tenant_id` | string | 是 | Teams 租戶 ID |
| `targets` | array | 是 | 目標對話列表 |

### Target 物件結構

| 欄位 | 類型 | 必填 | 說明 |
|------|------|------|------|
| `type` | string | 是 | 對話類型：`personal`, `channel`, `groupChat` |
| `tenant_id` | string | 是 | Teams 租戶 ID |
| `conversation_id` | string | 是 | Teams 對話 ID |

## 注意事項

1. **公司資料完整性**：確保 `company_id` 對應的公司記錄有完整的聯絡資訊（`contact_phone`, `address`, `contact_email`）
2. **用戶存在性**：確保 `created_by` 對應的用戶記錄存在
3. **Teams 租戶一致性**：所有 `targets` 中的 `tenant_id` 應該與 `teams_tenant_id` 一致
4. **對話 ID 有效性**：確保 `conversation_id` 對應的 Teams 對話存在且 Bot 已安裝

## 錯誤處理

常見錯誤及解決方法：

- `company not found`: 檢查公司 ID 是否正確，或更新公司資料的必填欄位
- `user not found`: 檢查用戶 ID 是否正確
- `invalid targets`: 檢查 targets 陣列格式是否正確
- `duplicate notify_key`: 專案名稱已存在，請使用不同的專案名稱
