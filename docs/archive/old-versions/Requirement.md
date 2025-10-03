# Teams Notification API Server 需求書

## 專案概述

建立一個企業級的Teams notification API server，能夠安全、高效地發送訊息到Microsoft Teams的各種目標（聊天群組、頻道、個人），並提供完整的監控、計費和運維功能。

**重要說明：** 此服務採用預先申請目的地的方式，使用notification key (project)來做mapping，確保訊息發送的安全性和可控性。

## FEATURE:

1. **核心功能**

   - API server可以送訊息到teams裡面的chatgroup, channel, person
   - 用的是teams framework, 透過azure bot再送到teams
   - 每次notification API call的時候 查詢DB data, 檢查來源跟目標 再送目的地
   - **支援廣播功能：可同時發送訊息到多個預先註冊的目的地**

2. **效能優化**

   - 利用cache減低token跟db 時間
   - 利用queue系統 處理azure bot rate limit issue

3. **記錄與計費**

   - 需要在db裡面記錄每次call的紀錄
   - 提供計費的API 要能依造呼叫者的email 公司別，notification key(project)來分類

4. **訊息格式支援**

   - 支援多種訊息格式：純文字、檔案附件
   - 支援@提及功能

5. **可靠性保障**

   - 實作訊息發送狀態追蹤和重試機制
   - 支援訊息排程發送（延遲發送、定期發送）

6. **管理與分析**

   - 提供訊息發送歷史查詢和統計分析API
   - 實作權限控制，確保只有授權用戶能發送訊息到指定目標

7. **安全性**

   - 支援訊息加密和敏感資訊過濾
   - 提供Webhook回調機制，通知發送結果

8. **運維監控**

   - 實作健康檢查和監控端點
   - 完整的日誌記錄和效能監控

## 目的地管理

### 預先申請流程
1. **申請表單提交**
   - 申請者提供：公司資訊、聯絡人、使用目的
   - 指定Teams目標：teamId, channelId, userId, groupChatId等
   - 設定notification key (project identifier)

2. **管理員審核**
   - 驗證申請者身份和權限
   - 確認Teams目標的有效性
   - 設定發送權限和限制

3. **目的地註冊**
   - 在系統中建立目的地記錄
   - 分配唯一的notification key
   - 設定發送配額和頻率限制

### 目的地類型
- **個人用戶**: 單一Teams用戶 (userId)
- **頻道**: 特定team的特定channel (teamId + channelId)
- **聊天群組**: 私人聊天群組 (groupChatId)
- **廣播群組**: 預定義的多目的地組合

### Teams目標識別符
- **userId**: Teams用戶的唯一識別符
- **teamId**: Teams團隊的唯一識別符
- **channelId**: 團隊內頻道的唯一識別符
- **groupChatId**: 私人聊天群組的唯一識別符

## EXAMPLES:

### 基本訊息發送

```bash
POST /api/v1/notifications
{
  "message": {
    "type": "text",
    "content": "Hello from API!",
    "mentions": ["@user1", "@user2"]
  },
  "sender": "api-user@company.com",
  "notificationKey": "daily-reminder",
  "priority": "normal"
}
```

### 檔案附件發送

```bash
POST /api/v1/notifications
{
  "message": {
    "type": "file",
    "content": "請查看附件報告",
    "attachment": {
      "fileName": "monthly-report.pdf",
      "fileUrl": "https://storage.company.com/reports/monthly-report.pdf",
      "fileSize": 2048576,
      "mimeType": "application/pdf"
    }
  },
  "sender": "system@company.com",
  "notificationKey": "monthly-report",
  "priority": "high"
}
```

### 排程發送

```bash
POST /api/v1/notifications/schedule
{
  "message": {"type": "text", "content": "Scheduled message"},
  "schedule": {
    "type": "recurring",
    "pattern": "0 9 * * 1", // 每週一上午9點
    "timezone": "Asia/Taipei",
    "startDate": "2024-01-01",
    "endDate": "2024-12-31"
  },
  "sender": "scheduler@company.com",
  "notificationKey": "weekly-reminder"
}
```

### 批次發送

```bash
POST /api/v1/notifications/batch
{
  "targets": [
    {"type": "channel", "teamId": "team-123", "channelId": "channel-456"},
    {"type": "person", "userId": "user-789"},
    {"type": "chatgroup", "groupId": "group-123"}
  ],
  "message": {
    "type": "text",
    "content": "重要通知：系統將於今晚進行維護"
  },
  "sender": "admin@company.com",
  "notificationKey": "system-maintenance",
  "priority": "urgent"
}
```

### 查詢發送狀態

```bash
GET /api/v1/notifications/status/{notificationId}
```

### 查詢發送歷史

```bash
GET /api/v1/notifications/history?email=user@company.com&company=CompanyA&key=daily-reminder&startDate=2024-01-01&endDate=2024-01-31
```

### 計費查詢

```bash
GET /api/v1/billing/summary?email=user@company.com&company=CompanyA&period=monthly
```

## 技術架構

### 核心組件

- **API Gateway**: 處理HTTP請求、認證、限流
- **Notification Service**: 核心業務邏輯
- **Teams Bot Service**: 與Microsoft Teams整合
- **Queue Service**: 處理訊息排隊和重試
- **Cache Service**: Redis快取
- **Database**: PostgreSQL/MySQL存儲
- **Monitoring**: 健康檢查和指標收集

### 資料流程

1. API請求 → 認證驗證
2. 權限檢查 → 目標驗證
3. 訊息入隊 → 快取檢查
4. Teams發送 → 狀態更新
5. 結果記錄 → 計費統計

## DOCUMENTATION:

### Microsoft Teams 相關

- [Microsoft Teams Bot Framework Documentation](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/what-are-bots)
- [Teams Bot API Reference](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/how-to/rate-limit)
- [Adaptive Cards Schema](https://adaptivecards.io/explorer/)
- [Teams Graph API](https://docs.microsoft.com/en-us/graph/api/resources/teams-api-overview)
- [Teams Bot Authentication](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/how-to/authentication/auth-flow-bot)

### Azure Bot Service

- [Azure Bot Service Documentation](https://docs.microsoft.com/en-us/azure/bot-service/)
- [Bot Framework SDK](https://github.com/microsoft/botframework-sdk)
- [Bot Framework Rate Limits](https://docs.microsoft.com/en-us/azure/bot-service/bot-service-quotas-and-limits)
- [Bot Framework Channels](https://docs.microsoft.com/en-us/azure/bot-service/bot-service-channels-reference)

### 技術架構參考

- [Redis Caching Patterns](https://redis.io/topics/patterns)
- [Message Queue Best Practices](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-messaging-overview)
- [Database Design for Notifications](https://docs.microsoft.com/en-us/azure/architecture/patterns/event-sourcing)
- [API Design Best Practices](https://docs.microsoft.com/en-us/azure/architecture/best-practices/api-design)

## OTHER CONSIDERATIONS:

### 安全性考量

- **認證授權**
  - API Key 輪換機制和過期管理
  - JWT Token驗證
  - OAuth 2.0整合
- **資料安全**
  - 訊息內容的XSS防護和輸入驗證
  - 目標地址的權限驗證（防止跨團隊/跨公司發送）
  - 敏感資訊的加密存儲和傳輸
  - HTTPS強制使用

### 效能和擴展性

- **異步處理**
  - 訊息發送的異步處理和並發控制
  - 大量訊息發送時的批次處理策略
- **快取優化**
  - 資料庫連接池和查詢優化
  - 快取策略的失效機制和一致性保證
  - 多層快取架構（記憶體、Redis、CDN）

### 監控和運維

- **日誌系統**
  - 完整的日誌記錄（結構化日誌）
  - 集中式日誌收集和分析
- **監控指標**
  - 效能指標收集（響應時間、成功率、錯誤率）
  - 業務指標監控（發送量、用戶活躍度）
- **告警機制**
  - 發送失敗、API延遲、錯誤率超標
  - 多級告警（郵件、簡訊、Teams通知）
- **健康檢查**
  - 健康檢查端點的設計
  - 依賴服務狀態監控

### 錯誤處理

- **重試策略**
  - 網路異常的重試策略和指數退避
  - 可配置的重試次數和間隔
- **降級處理**
  - Teams API 限制的優雅降級
  - 服務不可用時的備用方案
- **錯誤回饋**
  - 無效目標地址的處理和回饋
  - 部分失敗的處理（部分用戶收到，部分失敗）
  - 詳細的錯誤碼和錯誤訊息

### 合規性

- **審計追蹤**
  - 訊息發送的審計日誌
  - 操作記錄的完整保存
- **資料管理**
  - 資料保留政策和清理機制
  - 資料備份和恢復策略
- **隱私保護**
  - 隱私保護（GDPR、CCPA等）
  - 個人資料的匿名化處理
- **內容合規**
  - 訊息內容的合規性檢查
  - 敏感詞過濾和內容審核

### 測試策略

- **測試覆蓋**
  - 單元測試（API邏輯、驗證規則）
  - 整合測試（Teams API模擬、資料庫操作）
  - 負載測試（高併發場景）
  - 端到端測試（完整發送流程）
- **測試環境**
  - 多環境測試（開發、測試、預生產）
  - 自動化測試CI/CD整合

### 部署和配置

- **環境管理**
  - 環境配置管理（開發、測試、生產）
  - 配置的版本控制和審計
- **安全管理**
  - 密鑰和敏感資訊的安全管理
  - 密鑰輪換和過期管理
- **部署策略**
  - 容器化部署和服務網格整合
  - 藍綠部署和回滾策略
  - 零停機部署

### 災難恢復

- **備份策略**
  - 資料庫定期備份
  - 配置檔案備份
- **恢復計劃**
  - 服務中斷的快速恢復
  - 資料丟失的恢復策略
- **業務連續性**
  - 多區域部署
  - 故障轉移機制

## 開發時程

### Phase 1: MVP (4-6週)

- 基本API架構
- Teams Bot整合
- 簡單的訊息發送
- 基本認證

### Phase 2: 核心功能 (6-8週)

- 快取系統
- 隊列處理
- 資料庫設計
- 權限控制

### Phase 3: 進階功能 (4-6週)

- 排程發送
- 批次處理
- 監控系統
- 計費API

### Phase 4: 優化與測試 (2-4週)

- 效能優化
- 安全加固
- 完整測試
- 部署準備

## 成功指標

### 技術指標

- API響應時間 < 200ms
- 訊息發送成功率 > 99.5%
- 系統可用性 > 99.9%
- 支援併發用戶 > 1000

### 業務指標

- 訊息發送量 > 10000/天
- 用戶滿意度 > 4.5/5
- 系統穩定性 > 99.5%
- 成本控制 < 預算的110%
