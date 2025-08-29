## FEATURE:

1. API server可以送訊息到teams裡面的chatgroup, channel, person
2. 用的是teams framework, 透過azure bot再送到teams
3. 每次notification API call的時候 查詢DB data, 檢查來源跟目標 再送目的地
4. 利用cache減低token跟db 時間
5. 利用queue系統 處理azure bot rate limit issue
6. 需要在db裡面記錄每次call的紀錄 並且提供計費的API 要能依造呼叫者的email 公司別，notification key來分類
7. 支援多種訊息格式：純文字、檔案附件
8. 實作訊息發送狀態追蹤和重試機制
9. 支援訊息排程發送（延遲發送、定期發送）
10. 提供訊息發送歷史查詢和統計分析API
11. 實作權限控制，確保只有授權用戶能發送訊息到指定目標
12. 支援訊息加密和敏感資訊過濾
13. 提供Webhook回調機制，通知發送結果
14. 實作健康檢查和監控端點

## EXAMPLES:

### 基本訊息發送
```bash
POST /api/v1/notifications
{
  "target": {
    "type": "channel",
    "teamId": "team-123",
    "channelId": "channel-456"
  },
  "message": {
    "type": "text",
    "content": "Hello from API!",
    "mentions": ["@user1", "@user2"]
  },
  "sender": "api-user@company.com",
  "notificationKey": "daily-reminder"
}
```

### 檔案附件發送
```bash
POST /api/v1/notifications
{
  "target": {
    "type": "person",
    "userId": "user-789"
  },
  "message": {
    "type": "file",
    "content": "請查看附件報告",
    "attachment": {
      "fileName": "monthly-report.pdf",
      "fileUrl": "https://storage.company.com/reports/monthly-report.pdf",
      "fileSize": 2048576
    }
  },
  "sender": "system@company.com",
  "notificationKey": "monthly-report"
}
```

### 排程發送
```bash
POST /api/v1/notifications/schedule
{
  "target": {"type": "chatgroup", "groupId": "group-123"},
  "message": {"type": "text", "content": "Scheduled message"},
  "schedule": {
    "type": "recurring",
    "pattern": "0 9 * * 1", // 每週一上午9點
    "timezone": "Asia/Taipei"
  }
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
  "notificationKey": "system-maintenance"
}
```

## DOCUMENTATION:

### Microsoft Teams 相關
- [Microsoft Teams Bot Framework Documentation](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/what-are-bots)
- [Teams Bot API Reference](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/how-to/rate-limit)
- [Adaptive Cards Schema](https://adaptivecards.io/explorer/)
- [Teams Graph API](https://docs.microsoft.com/en-us/graph/api/resources/teams-api-overview)

### Azure Bot Service
- [Azure Bot Service Documentation](https://docs.microsoft.com/en-us/azure/bot-service/)
- [Bot Framework SDK](https://github.com/microsoft/botframework-sdk)
- [Bot Framework Rate Limits](https://docs.microsoft.com/en-us/azure/bot-service/bot-service-quotas-and-limits)

### 技術架構參考
- [Redis Caching Patterns](https://redis.io/topics/patterns)
- [Message Queue Best Practices](https://docs.microsoft.com/en-us/azure/service-bus-messaging/service-bus-messaging-overview)
- [Database Design for Notifications](https://docs.microsoft.com/en-us/azure/architecture/patterns/event-sourcing)

## OTHER CONSIDERATIONS:

### 安全性考量
- API Key 輪換機制和過期管理
- 訊息內容的XSS防護和輸入驗證
- 目標地址的權限驗證（防止跨團隊/跨公司發送）
- 敏感資訊的加密存儲和傳輸

### 效能和擴展性
- 訊息發送的異步處理和並發控制
- 大量訊息發送時的批次處理策略
- 資料庫連接池和查詢優化
- 快取策略的失效機制和一致性保證

### 監控和運維
- 完整的日誌記錄（結構化日誌）
- 效能指標收集（響應時間、成功率、錯誤率）
- 告警機制（發送失敗、API延遲、錯誤率超標）
- 健康檢查端點的設計

### 錯誤處理
- 網路異常的重試策略和指數退避
- Teams API 限制的優雅降級
- 無效目標地址的處理和回饋
- 部分失敗的處理（部分用戶收到，部分失敗）

### 合規性
- 訊息發送的審計日誌
- 資料保留政策和清理機制
- 隱私保護（GDPR、CCPA等）
- 訊息內容的合規性檢查

### 測試策略
- 單元測試（API邏輯、驗證規則）
- 整合測試（Teams API模擬、資料庫操作）
- 負載測試（高併發場景）
- 端到端測試（完整發送流程）

### 部署和配置
- 環境配置管理（開發、測試、生產）
- 密鑰和敏感資訊的安全管理
- 容器化部署和服務網格整合
- 藍綠部署和回滾策略
