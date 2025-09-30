# Teams Notification API Server – Requirement Alignment

This document keeps the original product intent while clarifying what the current repository delivers versus future work.

## 專案概述

建置一個企業級 Microsoft Teams 通知 API，先行審核目的地（destination），並由通知金鑰（project/notify_key）控管發送權限與度量，確保安全與可追蹤性。

## 功能矩陣

| 功能類別 | 說明 | 現況 |
| --- | --- | --- |
| 通知發送 | 透過預先註冊的目的地發送 Teams 訊息（text/file/adaptive card） | ✅ `internal/api/handlers/notifications` 已實作 CRUD 與派送流程（BroadcastService 尚使用 stub HTTP 呼叫） |
| 多目標廣播 | 單次呼叫可對多個目的地發送 | ✅ `BroadcastService.SendToDestinations` 已實作迴圈與結果記錄 |
| 目的地管理 | 目的地審核、目標驗證、專案綁定 | ✅ CRUD + 驗證端點齊備；`validateTeamsTargets` 強制 tenant_id/會話 ID |
| Teams Bot 管理 | 平台 Bot 資源維運 | ✅ 平台 Bot 路由存在；第三方 Bot 端點已移除（需後續補完） |
| 認證授權 | 對外訪問驗證 | ⏳ Middleware scaffold 未啟用；目前 API 無需 Authorization |
| 使用者 / 公司 / 專案 CRUD | 基礎管理資料 | ✅ Handlers + services/repositories 完整 |
| 計費與統計 | 使用量統計、計費 API | ⏳ Schema 留有欄位，但尚未提供 API |
| 佇列與速率策略 | Redis/queue 處理速率限制 | ⏳ 設計文件與 config 有描述，實作尚未著手 |
| 訊息排程 | 延遲/週期性通知 | ⏳ 需求中記載，尚無程式碼支援 |
| 監控與觀測 | 健康檢查、結構化日誌 | ✅ `/health` + Logrus JSON middleware 已啟用；指標/追蹤待補 |

## 需求摘要（仍有效）
- 預先申請目的地 → 管理員審核 → 建立目的地 → 指派 notification key。
- 支援多種 Teams 目標：個人、頻道、群組（目前以 `type = personal|channel|groupchat` 及 `conversation_id`/`tenant_id` 等欄位儲存）。
- 詳細紀錄：`notifications` 與 `notification_destinations` 表保留送達狀態、錯誤、重試計數。
- 權限控管：Auth middleware 尚未啟用；後續可透過 `RequireRole` / `RequireAPIKey` 加強。

## 範例請求（符合當前 Handler）

```http
POST /api/v1/notifications
Content-Type: application/json

{
  "project_id": "00000000-0000-0000-0000-000000000001",
  "sender_id": "00000000-0000-0000-0000-000000000002",
  "message_type": "text",
  "content": "Hello from API!",
  "mentions": ["@user1"],
  "priority": "normal",
  "destinations": ["00000000-0000-0000-0000-000000000003"],
  "metadata": {
    "ticket": "INC-123"
  }
}
```

檔案附件與 Adaptive Card 只需對應 `message_type` 及 `attachment` / `adaptive_card` 欄位即可。

## 待辦與風險
- 佇列、排程、重試策略尚未落地，需評估 Redis/Kafka 等元件整合。
- API Key 機制需補強（Key 生成、儲存、授權綁定、快取）。
- 計費與報表 API 尚未出現；Schema 僅存放基本資料。
- Graph API / Bot Framework 實際呼叫目前為 stub，需接軌真實 Teams API 並處理 token 交換、錯誤分類。

此文件可作為未來規劃與差距分析的起點，後續實作請同步更新矩陣與範例。
