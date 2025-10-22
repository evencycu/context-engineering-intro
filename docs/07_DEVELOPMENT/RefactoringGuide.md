# 🔄 目錄結構重構指南

## 1. 文件資訊

- **版本**：v1.2.0
- **撰寫人**：Development Team
- **最後更新**：2025-10-21
- **目的**：記錄目錄結構重構的變化和影響

---

## 2. 重構概述

本次重構主要目的是提升代碼的可維護性和模組化程度，將共用函式庫和業務邏輯進行分離。

### 2.1 重構目標

- **模組化設計**: 將共用函式庫獨立出來
- **業務邏輯分離**: 將業務服務集中管理
- **提升重用性**: 共用組件可在多個服務中使用
- **改善維護性**: 清晰的目錄結構便於維護

---

## 3. 目錄結構變化

### 3.1 新增目錄

#### `libs/` - 共用函式庫

```
libs/
├── token/                 # Token 管理
│   └── token_manager.go   # Teams Token 快取和管理
├── middleware/            # 中間件
│   ├── auth.go           # 認證中間件
│   ├── error.go          # 錯誤處理中間件
│   ├── logging.go        # 日誌中間件
│   └── business_metrics.go # 業務指標中間件
├── logger/               # 日誌系統
│   └── logger.go         # 結構化日誌記錄
├── models/               # 共用模型
│   └── models.go         # 跨模組共用的資料模型
├── storage/              # 本地存儲
│   └── local_storage.go  # 本地檔案存儲
└── errors/               # 錯誤處理
    └── errors.go         # 統一的錯誤處理機制
```

#### `services/` - 業務服務

```
services/
└── teamsnotification/    # Teams 通知服務
    ├── server.go         # 服務器配置
    ├── handlers/         # HTTP 處理器
    ├── repositories/     # 資料存取層
    ├── actor/           # Actor 模式實作
    └── integration/     # 整合測試
```

### 3.2 保留目錄

#### `internal/` - 內部共用 (Legacy)

```
internal/                 # 內部共用，不對外公開 (Legacy)
├── api/                 # API 層
│   ├── handlers/        # HTTP 處理器
│   ├── middleware/      # 中間件
│   ├── repositories/    # 資料存取層
│   └── services/        # 業務邏輯層
└── database/            # 資料庫模型
    ├── models.go
    ├── schema.sql
    └── init.sql
```

---

## 4. 遷移影響

### 4.1 代碼遷移

#### 已遷移的組件

| 原始位置 | 新位置 | 說明 |
|---------|--------|------|
| `services/teamsnotification/services/token_manager.go` | `libs/token/token_manager.go` | Token 管理 |
| `internal/api/middleware/*` | `libs/middleware/*` | 中間件組件 |
| `internal/api/services/*` | `services/teamsnotification/*` | 業務服務 |

#### 包名變更

| 組件 | 原始包名 | 新包名 |
|------|---------|--------|
| Token Manager | `package services` | `package token` |
| Middleware | `package middleware` | `package middleware` |
| Logger | `package logger` | `package logger` |

### 4.2 導入路徑變更

#### 主要變更

```go
// 舊的導入方式
import "services/teamsnotification/services"

// 新的導入方式
import "libs/token"
import "libs/middleware"
import "services/teamsnotification/handlers"
```

#### 具體範例

```go
// Token Manager 使用
import "libs/token"

// 中間件使用
import "libs/middleware"

// 業務服務使用
import "services/teamsnotification/handlers"
import "services/teamsnotification/repositories"
```

---

## 5. 開發指南

### 5.1 新功能開發

#### 共用組件

- **位置**: `libs/` 目錄
- **用途**: 跨模組共用的功能
- **範例**: 中間件、工具函數、共用模型

#### 業務邏輯

- **位置**: `services/teamsnotification/` 目錄
- **用途**: 特定業務領域的功能
- **範例**: 通知處理、用戶管理、專案管理

### 5.2 代碼組織原則

#### 1. 分層原則

```
libs/           # 基礎設施層
services/       # 業務邏輯層
internal/       # 內部實現層 (Legacy)
```

#### 2. 依賴方向

```
services/ → libs/        # 業務層依賴基礎設施層
internal/ → libs/        # 內部層依賴基礎設施層
services/ → internal/    # 業務層可依賴內部層 (過渡期)
```

#### 3. 模組邊界

- **libs/**: 無業務邏輯，純技術組件
- **services/**: 包含業務邏輯，可依賴 libs/
- **internal/**: 內部實現，逐步遷移到 services/

---

## 6. 遷移計劃

### 6.1 階段一：基礎設施遷移 ✅

- [x] 創建 `libs/` 目錄結構
- [x] 遷移 Token Manager
- [x] 遷移中間件組件
- [x] 遷移日誌系統
- [x] 更新導入路徑

### 6.2 階段二：業務邏輯遷移 ✅

- [x] 創建 `services/teamsnotification/` 目錄結構
- [x] 遷移業務服務
- [x] 遷移 HTTP 處理器
- [x] 遷移資料存取層
- [x] 更新服務配置

### 6.3 階段三：文檔更新 ✅

- [x] 更新 README.md
- [x] 更新架構文檔
- [x] 更新 OpenAPI 規範
- [x] 創建重構指南

### 6.4 階段四：清理和優化 (進行中)

- [ ] 清理 `internal/` 中的重複代碼
- [ ] 統一錯誤處理機制
- [ ] 優化模組間依賴關係
- [ ] 完善測試覆蓋

---

## 7. 測試策略

### 7.1 單元測試

```bash
# 測試共用函式庫
go test ./libs/...

# 測試業務服務
go test ./services/...

# 測試內部組件
go test ./internal/...
```

### 7.2 整合測試

```bash
# 測試完整服務
go test ./services/teamsnotification/integration/...
```

### 7.3 端到端測試

```bash
# 執行 API 測試
./scripts/test/api/simple_*.sh
```

---

## 8. 常見問題

### 8.1 導入錯誤

**問題**: 找不到模組

```go
// 錯誤
import "services/teamsnotification/services"

// 正確
import "libs/token"
```

**解決方案**: 檢查 `go.mod` 中的模組路徑

### 8.2 包名衝突

**問題**: 包名重複

```go
// 錯誤
package services  // 在多個地方使用

// 正確
package token     // 在 libs/token/
package handlers  // 在 services/teamsnotification/handlers/
```

**解決方案**: 使用具體的包名

### 8.3 循環依賴

**問題**: 模組間循環依賴

**解決方案**: 
- 將共用邏輯移至 `libs/`
- 使用介面解耦
- 重新設計模組邊界

---

## 9. 最佳實踐

### 9.1 代碼組織

1. **單一職責**: 每個模組只負責一個功能領域
2. **依賴倒置**: 依賴抽象而非具體實現
3. **介面隔離**: 定義小而專注的介面
4. **開放封閉**: 對擴展開放，對修改封閉

### 9.2 命名規範

1. **目錄命名**: 使用小寫字母和底線
2. **包命名**: 使用小寫字母，簡潔明瞭
3. **檔案命名**: 使用小寫字母和底線
4. **函數命名**: 使用駝峰命名法

### 9.3 文檔維護

1. **README 更新**: 及時更新專案說明
2. **架構文檔**: 保持架構文檔同步
3. **API 文檔**: 維護 OpenAPI 規範
4. **變更日誌**: 記錄重要變更

---

## 10. 後續計劃

### 10.1 短期目標

- [ ] 完成 `internal/` 目錄的清理
- [ ] 統一錯誤處理機制
- [ ] 完善單元測試覆蓋
- [ ] 優化模組間依賴

### 10.2 中期目標

- [ ] 實現微服務架構
- [ ] 添加服務發現機制
- [ ] 實現配置中心
- [ ] 添加監控和告警

### 10.3 長期目標

- [ ] 支援多租戶架構
- [ ] 實現水平擴展
- [ ] 添加 CI/CD 流水線
- [ ] 實現雲原生部署

---

## 11. 參考資源

### 11.1 相關文檔

- [架構文檔](../02_ARCHITECTURE/Architecture.md)
- [開發指南](DevelopmentGuide.md)
- [API 文檔](../../api/openapi/teams-notification-api.yaml)

### 11.2 外部資源

- [Go 模組管理](https://go.dev/blog/using-go-modules)
- [Go 專案佈局](https://github.com/golang-standards/project-layout)
- [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

---

## 12. 聯絡資訊

如有問題或建議，請聯絡：

- **開發團隊**: development@teamsnotify.local
- **技術負責人**: tech-lead@teamsnotify.local
- **專案經理**: pm@teamsnotify.local

---

*最後更新: 2025-10-21*
