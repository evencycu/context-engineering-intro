# 📘 Product Requirement Document (PRD)

## 1. 文件資訊

- **文件版本**：v1.0
- **撰寫人**：Product Manager
- **最後更新**：2025-10-08
- **審核人**：System Analyst, Tech Lead

---

## 2. 專案簡介

### 2.1 背景

建立一個企業級的Teams notification API server，能夠安全、高效地發送訊息到Microsoft Teams的各種目標（聊天群組、頻道、個人），並提供完整的監控、計費和運維功能。

**重要說明：** 此服務採用預先申請目的地的方式，使用notification key (project)來做mapping，確保訊息發送的安全性和可控性。

### 2.2 專案目標

- 統一企業內部通知系統
- 整合 Teams Bot 與內部系統通知流程
- 建立可擴展的通知中台架構
- 提供完整的監控和計費功能

---

## 3. 功能概要


| 功能編號 | 功能名稱     | 說明                                 | 優先級 (P0/P1/P2) |
| -------- | ------------ | ------------------------------------ | ----------------- |
| F-001    | 基本訊息發送 | 支援文字、檔案附件發送到Teams        | P0                |
| F-002    | 目的地管理   | 預先申請和審核Teams目標              | P0                |
| F-003    | 廣播功能     | 同時發送訊息到多個預先註冊的目的地   | P1                |
| F-004    | 排程發送     | 支援延遲發送、定期發送               | P1                |
| F-005    | 批次處理     | 大量訊息的高效批次發送               | P1                |
| F-006    | 狀態追蹤     | 訊息發送狀態追蹤和重試機制           | P1                |
| F-007    | 計費系統     | 依呼叫者、公司、專案分類計費         | P2                |
| F-008    | 監控分析     | 訊息發送歷史查詢和統計分析           | P2                |
| F-009    | 權限控制     | 確保只有授權用戶能發送訊息到指定目標 | P2                |
| F-010    | 安全功能     | 訊息加密、敏感資訊過濾、Webhook回調  | P2                |

---

## 4. 使用者故事 (User Stories)


| 編號   | 故事敘述                                                            | 驗收條件                                 |
| ------ | ------------------------------------------------------------------- | ---------------------------------------- |
| US-001 | 作為系統管理員，我希望能預先申請Teams目的地，以便控制訊息發送範圍   | 成功註冊目的地並獲得notification key     |
| US-002 | 作為開發者，我希望能透過API發送訊息到Teams，以便整合到現有系統      | API調用成功並在Teams收到訊息             |
| US-003 | 作為業務人員，我希望能廣播重要通知到多個Teams群組，以便確保訊息傳達 | 同時發送到多個預註冊的目的地             |
| US-004 | 作為系統管理員，我希望能排程發送定期通知，以便自動化通知流程        | 系統按時自動發送排程訊息                 |
| US-005 | 作為財務人員，我希望能查看API使用計費，以便控制成本                 | 能查詢依使用者、公司、專案分類的計費報表 |
| US-006 | 作為運維人員，我希望能監控訊息發送狀態，以便確保系統正常運作        | 能查看訊息發送歷史和統計分析             |

---

## 5. 成功指標 (KPI)

### 技術指標

- API響應時間 < 200ms
- 訊息發送成功率 > 99.5%
- 系統可用性 > 99.9%
- 支援併發用戶 > 1000

### 商務指標

- 訊息發送量 > 10000/天
- 系統穩定性 > 99.5%
- 成本控制 < 預算的110%

---

## 6. 依賴與限制

### 技術依賴

- Microsoft Teams Bot Framework
- Azure Bot Service
- Microsoft Graph API
- Microsoft Entra Service
- Redis Queue System
- PostgreSQL Database

### 業務限制

- 需要Teams管理員權限進行Bot註冊
- 內部Teams租戶權限受政策影響
- 需要預先申請和審核目的地

---

## 7. 未來展望

- 第三方 Bot 管理基本架構

## 8. 風險評估

### 技術風險

- Teams API限制和變更
- 高併發下的效能問題
- 資料安全合規要求

### 業務風險

- 用戶接受度和採用率
- 與現有系統整合複雜度
- 維護和運營成本

---

## 9. 參考文件

- [Microsoft Teams Bot Framework Documentation](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/what-are-bots)
- [Teams Bot API Reference](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/how-to/rate-limit)
- [Azure Bot Service Documentation](https://docs.microsoft.com/en-us/azure/bot-service/)
