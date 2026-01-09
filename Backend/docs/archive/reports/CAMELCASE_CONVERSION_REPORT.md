# API Key Name 轉換報告 - snake_case to camelCase

## 執行日期
2025-10-20

## 轉換範圍

### ✅ 已更新的 Handler 文件 (9個)
1. `internal/api/handlers/external/handler.go` - External API
2. `internal/api/handlers/projects/handler.go` - Projects API
3. `internal/api/handlers/users/handler.go` - Users API
4. `internal/api/handlers/notifications/handler.go` - Notifications API
5. `internal/api/handlers/bots/handler.go` - Bots API
6. `internal/api/handlers/destinations/handler.go` - Destinations API
7. `internal/api/handlers/companies/handler.go` - Companies API
8. `internal/api/handlers/billing/handler.go` - Billing API
9. `internal/api/handlers/files/handler.go` - Files API

### ✅ 已更新的測試腳本 (20個)
- `scripts/test/api/*.sh` - 所有 API 測試腳本
- `scripts/test/integration/*.sh` - 整合測試腳本
- `scripts/test/e2e/*.sh` - E2E 測試腳本
- `scripts/test/security/*.sh` - 安全測試腳本
- `scripts/test/performance/*.sh` - 效能測試腳本

### ✅ 已更新的文檔
- `api/openapi/teams-notification-api.yaml` - OpenAPI 規格文件
- `docs/**/*.md` - 所有文檔中的 JSON 範例
- `README.md`, `API_ROUTING*.md` - 主要文檔

## 轉換對照表

### 常用字段
| snake_case | camelCase |
|------------|-----------|
| `notify_key` | `notifyKey` |
| `message_type` | `messageType` |
| `company_id` | `companyId` |
| `project_id` | `projectId` |
| `created_by` | `createdBy` |
| `daily_limit` | `dailyLimit` |
| `monthly_limit` | `monthlyLimit` |

### 通知相關
| snake_case | camelCase |
|------------|-----------|
| `notification_id` | `notificationId` |
| `sender_id` | `senderId` |
| `adaptive_card` | `adaptiveCard` |
| `destinations_count` | `destinationsCount` |
| `estimated_delivery_time` | `estimatedDeliveryTime` |
| `destinations_sent` | `destinationsSent` |
| `destinations_failed` | `destinationsFailed` |
| `total_destinations` | `totalDestinations` |

### Bot 相關
| snake_case | camelCase |
|------------|-----------|
| `app_id` | `appId` |
| `app_password` | `appPassword` |
| `tenant_id` | `tenantId` |
| `webhook_url` | `webhookUrl` |
| `teams_tenant_id` | `teamsTenantId` |
| `rate_limit_per_minute` | `rateLimitPerMinute` |
| `max_concurrent_requests` | `maxConcurrentRequests` |
| `api_endpoint` | `apiEndpoint` |
| `bot_id` | `botId` |

### 時間戳記
| snake_case | camelCase |
|------------|-----------|
| `created_at` | `createdAt` |
| `updated_at` | `updatedAt` |
| `sent_at` | `sentAt` |

### 錯誤與訊息
| snake_case | camelCase |
|------------|-----------|
| `error_message` | `errorMessage` |

### 檔案相關
| snake_case | camelCase |
|------------|-----------|
| `file_name` | `fileName` |
| `file_size` | `fileSize` |
| `content_type` | `contentType` |

### 計費相關
| snake_case | camelCase |
|------------|-----------|
| `price_per_notification` | `pricePerNotification` |
| `price_per_attachment` | `pricePerAttachment` |
| `price_per_mention` | `pricePerMention` |
| `price_per_adaptive_card` | `pricePerAdaptiveCard` |
| `billing_email` | `billingEmail` |
| `billing_enabled` | `billingEnabled` |
| `billing_plan_id` | `billingPlanId` |
| `payment_method` | `paymentMethod` |

### 公司與用戶
| snake_case | camelCase |
|------------|-----------|
| `contact_email` | `contactEmail` |
| `contact_phone` | `contactPhone` |
| `old_password` | `oldPassword` |
| `new_password` | `newPassword` |

### 目的地相關
| snake_case | camelCase |
|------------|-----------|
| `target_type` | `targetType` |
| `target_id` | `targetId` |
| `team_id` | `teamId` |
| `channel_id` | `channelId` |
| `user_id` | `userId` |
| `group_id` | `groupId` |

## 影響範圍

### ⚠️ Breaking Changes
所有使用舊 API 的客戶端都需要更新：
- 外部 API (`/api/v1/notify`, `/api/v1/destinations/*`)
- 內部 API (`/internal/v1/*`)

### 📝 更新建議
1. 更新所有 API 客戶端使用 camelCase 格式
2. 更新 Postman/Insomnia 集合
3. 更新前端應用的 API 調用
4. 更新自動化腳本

## 測試驗證

### ✅ 測試結果
- camelCase 格式: ✅ 正常工作
- snake_case 格式: ❌ 返回驗證錯誤（預期行為）
- API 功能: ✅ 正常運作
- 服務器啟動: ✅ 成功

### 📊 統計
- 更新的 Go 文件: 9 個
- 更新的測試腳本: 20 個
- 更新的文檔文件: 50+ 個
- 轉換的字段: 48+ 個

## Git 提交建議

```bash
git add internal/api/handlers/ api/openapi/ scripts/test/ docs/
git commit -m "refactor: convert all API field names from snake_case to camelCase

BREAKING CHANGE: All API endpoints now use camelCase for JSON field names

Changes:
- Updated all handler structs to use camelCase JSON tags
- Updated OpenAPI specification with camelCase field names
- Updated all test scripts to use camelCase
- Updated documentation with camelCase examples

Affected APIs:
- External API (/api/v1/notify, /api/v1/destinations)
- Internal API (all /internal/v1/* endpoints)

Migration Guide:
- Replace all snake_case field names with camelCase equivalents
- See CAMELCASE_CONVERSION_REPORT.md for complete mapping

Examples:
- notify_key → notifyKey
- message_type → messageType
- company_id → companyId
"
```

## 備註
- 所有備份文件已保存為 `*.backup_camel`
- 舊的 snake_case 格式將返回驗證錯誤
- 建議客戶端盡快更新以使用新格式
