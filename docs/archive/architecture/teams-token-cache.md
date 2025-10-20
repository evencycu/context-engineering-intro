# Teams Token 快取機制設計

## 概述

本設計文檔描述 Teams Token 快取機制的實現方案，用於優化多 Teams Bot 平台中重複獲取 Microsoft Teams API Token 的性能問題。

## 問題分析

### 當前狀況
1. **重複 Token 獲取**: 每次發送 Teams 訊息都需要重新獲取 OAuth2 Token
2. **多 Bot 支援**: 平台支援多個 Teams Bot，每個 Bot 需要獨立的 Token
3. **Token 有效期**: Microsoft Teams Token 通常有 1 小時的有效期
4. **性能影響**: 每次 Token 獲取需要 1-2 秒的網路請求時間

### 識別到的 Token 獲取場景
1. **Bot Framework Connector Token**: 用於發送訊息到 Teams
2. **Microsoft Graph API Token**: 用於獲取用戶資訊等 Graph API 操作
3. **多租戶支援**: 不同 Teams Bot 可能屬於不同的 Azure AD 租戶

## 設計方案

### 1. Token 快取架構

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   API Server    │    │   Redis Cache    │    │  Microsoft OAuth│
│                 │    │                  │    │                 │
│ ┌─────────────┐ │    │ ┌──────────────┐ │    │ ┌─────────────┐ │
│ │Token Manager│ │◄──►│ │Token Cache   │ │    │ │Token Endpoint│ │
│ │             │ │    │ │              │ │    │ │             │ │
│ └─────────────┘ │    │ └──────────────┘ │    │ └─────────────┘ │
└─────────────────┘    └──────────────────┘    └─────────────────┘
```

### 2. Token 快取策略

#### 2.1 快取鍵設計
```
teams:token:{bot_id}:{tenant_id}:{token_type}
```

其中：
- `bot_id`: Teams Bot 的 App ID
- `tenant_id`: Azure AD 租戶 ID
- `token_type`: `connector` 或 `graph`

#### 2.2 Token 資料結構
```json
{
  "access_token": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiIs...",
  "token_type": "Bearer",
  "expires_in": 3600,
  "expires_at": 1699123456,
  "scope": "https://api.botframework.com/.default",
  "botId": "844146d7-4ac9-4e4d-a463-d6e027714e81",
  "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6"
}
```

### 3. 實現細節

#### 3.1 Token Manager 介面
```go
type TokenManager interface {
    GetConnectorToken(ctx context.Context, botID, tenantID string) (string, error)
    GetGraphToken(ctx context.Context, botID, tenantID string) (string, error)
    InvalidateToken(ctx context.Context, botID, tenantID, tokenType string) error
    RefreshToken(ctx context.Context, botID, tenantID, tokenType string) (string, error)
}
```

#### 3.2 快取邏輯
1. **檢查快取**: 首先檢查 Redis 中是否有有效的 Token
2. **Token 驗證**: 驗證 Token 是否即將過期（提前 5 分鐘刷新）
3. **Token 獲取**: 如果快取無效，從 Microsoft OAuth 端點獲取新 Token
4. **快取更新**: 將新 Token 存入 Redis，設定適當的過期時間

#### 3.3 錯誤處理
1. **Token 失效**: 自動刷新 Token
2. **網路錯誤**: 重試機制（最多 3 次）
3. **認證失敗**: 記錄錯誤並返回適當的錯誤訊息

### 4. 多 Bot 支援

#### 4.1 Bot 配置管理
```go
type BotConfig struct {
    BotID       string `json:"botId"`
    AppPassword string `json:"appPassword"`
    TenantID    string `json:"tenantId"`
    CompanyID   string `json:"companyId"`
}
```

#### 4.2 動態 Bot 發現
- 從 `teams_bots` 表讀取 Bot 配置
- 支援運行時添加/移除 Bot
- 每個 Bot 維護獨立的 Token 快取

### 5. 性能優化

#### 5.1 快取策略
- **TTL 設定**: Token 快取時間設為 50 分鐘（比實際過期時間短 10 分鐘）
- **預刷新**: 在 Token 過期前 5 分鐘自動刷新
- **並發控制**: 使用 Redis SETNX 避免並發刷新同一個 Token

#### 5.2 監控指標
- Token 快取命中率
- Token 獲取延遲
- 快取失效頻率
- 並發請求數量

## 實現計劃

### Phase 1: 基礎 Token 快取
1. 實現 `TokenManager` 介面
2. 實現 Redis 快取邏輯
3. 整合到現有的 `NotificationActor`

### Phase 2: 多 Bot 支援
1. 從資料庫讀取 Bot 配置
2. 實現動態 Bot 管理
3. 支援多租戶 Token 快取

### Phase 3: 性能優化
1. 實現 Token 預刷新
2. 添加監控和指標
3. 優化並發處理

## 配置參數

### 環境變數
```bash
# Redis 配置
REDIS_URL=redis://localhost:6379

# Token 快取配置
TOKEN_CACHE_TTL=3000  # 5 分鐘（秒）
TOKEN_REFRESH_AHEAD=300  # 提前 5 分鐘刷新（秒）
TOKEN_MAX_RETRIES=3      # 最大重試次數

# Teams Bot 配置（預設 Bot）
TEAMS_BOT_APP_ID=844146d7-4ac9-4e4d-a463-d6e027714e81
TEAMS_BOT_APP_PASSWORD=HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfh
TEAMS_TENANT_ID=051cece0-e4dc-4aed-b471-bf29824e1ee6
```

### Redis 鍵命名規範
```
teams:token:{bot_id}:{tenant_id}:connector  # Bot Framework Token
teams:token:{bot_id}:{tenant_id}:graph      # Graph API Token
teams:token:refresh:{bot_id}:{tenant_id}   # 刷新鎖定鍵
```

## 測試策略

### 單元測試
- Token 獲取邏輯測試
- 快取命中/未命中測試
- 錯誤處理測試

### 整合測試
- 與 Redis 的整合測試
- 與 Microsoft OAuth 的整合測試
- 多 Bot 場景測試

### 性能測試
- 並發 Token 獲取測試
- 快取性能測試
- 記憶體使用測試

## 監控和日誌

### 關鍵指標
- `teams_token_cache_hits_total`: Token 快取命中次數
- `teams_token_cache_misses_total`: Token 快取未命中次數
- `teams_token_refresh_duration_seconds`: Token 刷新耗時
- `teams_token_active_count`: 活躍 Token 數量

### 日誌格式
```json
{
  "level": "info",
  "msg": "Token refreshed successfully",
  "botId": "844146d7-4ac9-4e4d-a463-d6e027714e81",
  "tenantId": "051cece0-e4dc-4aed-b471-bf29824e1ee6",
  "token_type": "connector",
  "expires_in": 3600,
  "duration_ms": 1250
}
```

## 安全考量

### Token 安全
- Token 在 Redis 中加密存儲
- 設定適當的 Redis ACL
- 定期輪換 Bot 認證資訊

### 存取控制
- 限制 Token Manager 的存取權限
- 實現 Token 存取審計
- 監控異常的 Token 使用模式

## 故障排除

### 常見問題
1. **Token 過期**: 檢查系統時間同步
2. **快取失效**: 檢查 Redis 連接和記憶體使用
3. **認證失敗**: 驗證 Bot 認證資訊
4. **並發問題**: 檢查 Redis SETNX 鎖定機制

### 除錯工具
- Redis CLI 檢查快取狀態
- 日誌分析工具
- 性能監控儀表板

## 結論

這個 Token 快取機制將顯著提升 Teams 通知系統的性能，減少不必要的 Token 獲取請求，並支援多 Bot 和多租戶的複雜場景。通過 Redis 快取和智能刷新策略，可以確保系統的高可用性和性能。

