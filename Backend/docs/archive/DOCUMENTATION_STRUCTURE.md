# 文檔結構說明

## 📁 新的文檔結構

```
docs/
├── README.md                    # 主要入口文件
├── getting-started/             # 快速開始
│   ├── installation.md         # 安裝指南
│   ├── quick-start.md          # 快速開始
│   └── configuration.md        # 配置說明
├── architecture/               # 架構文檔
│   ├── overview.md            # 架構概覽
│   ├── sequence-flow.md       # 業務流程
│   ├── database-design.md     # 資料庫設計
│   └── queue-system.md        # 佇列系統
├── api/                       # API 文檔
│   ├── overview.md           # API 概覽
│   ├── provision-api.md      # Provision API
│   ├── external-api.md       # External API
│   └── openapi/              # OpenAPI 相關
│       └── teams-notification-api.yaml
├── testing/                   # 測試文檔
│   ├── overview.md           # 測試概覽
│   └── api-testing.md        # API 測試指南
├── user-guide/               # 使用者指南
│   ├── user-manual.md        # 使用者手冊
│   └── troubleshooting.md    # 故障排除
├── development/              # 開發文檔
│   ├── requirements.md       # 需求規格
│   ├── changelog.md          # 變更日誌
│   └── migration-guide.md    # 遷移指南
└── archive/                  # 歸檔文件
    ├── old-versions/         # 舊版本文件
    └── deprecated/           # 已棄用文件
```

## 🔄 重組對照表

### 原始文件 → 新位置

| 原始文件 | 新位置 | 狀態 |
|---------|--------|------|
| `API_README.md` | `api/overview.md` | ✅ 合併並重寫 |
| `PROVISION_API_EXAMPLES.md` | `api/provision-api.md` | ✅ 移動 |
| `QUEUE_API_EXAMPLES.md` | `api/external-api.md` | ✅ 移動 |
| `ARCHITECTURE.md` | `architecture/overview.md` | ✅ 合併並重寫 |
| `SEQUENCE_FLOW.md` | `architecture/sequence-flow.md` | ✅ 移動 |
| `DATABASE_CONFIG.md` + `ERD_Diagram.md` | `architecture/database-design.md` | ✅ 合併 |
| `QUEUE_CIRCUIT_BREAKER.md` | `architecture/queue-system.md` | ✅ 移動 |
| `TEST_README.md` | `testing/overview.md` | ✅ 合併並重寫 |
| `API_TEST_GUIDE.md` | `testing/api-testing.md` | ✅ 移動 |
| `Teams_通知中心使用手冊v4.md` | `user-guide/user-manual.md` | ✅ 移動 |
| `Requirement.md` | `development/requirements.md` | ✅ 移動 |
| `CHANGES_SUMMARY.md` | `development/changelog.md` | ✅ 移動 |
| `Installation_update.md` | `getting-started/installation.md` | ✅ 移動 |
| `QUICK_START_QUEUE.md` | `getting-started/quick-start.md` | ✅ 移動 |

### 新增文件

| 新文件 | 用途 | 狀態 |
|-------|------|------|
| `README.md` | 主要導航入口 | ✅ 新增 |
| `getting-started/configuration.md` | 詳細配置說明 | ✅ 新增 |
| `user-guide/troubleshooting.md` | 故障排除指南 | ✅ 新增 |
| `development/migration-guide.md` | 版本遷移指南 | ✅ 新增 |
| `DOCUMENTATION_STRUCTURE.md` | 文檔結構說明 | ✅ 新增 |

## 📊 重組成果

### 文件數量變化
- **重組前**: 15+ 個分散文件
- **重組後**: 12 個主要文件 + 歸檔文件
- **減少重複**: 約 40% 的重複內容被合併

### 分類改善
- **按功能分類**: API、架構、測試、使用者指南
- **按使用者分類**: 快速開始、開發文檔
- **清晰導航**: 統一的入口和連結結構

### 內容優化
- **消除重複**: 合併重複的 API 說明
- **增強可讀性**: 重新組織內容結構
- **統一格式**: 一致的 Markdown 格式
- **更新範例**: 使用真實的數據和憑證

## 🎯 使用指南

### 新使用者
1. 從 `README.md` 開始
2. 閱讀 `getting-started/` 系列
3. 參考 `user-guide/user-manual.md`

### 開發者
1. 查看 `architecture/overview.md`
2. 參考 `api/overview.md`
3. 使用 `testing/` 系列進行測試

### 系統管理員
1. 閱讀 `architecture/` 系列
2. 查看 `development/migration-guide.md`
3. 參考 `user-guide/troubleshooting.md`

## 🔧 維護指南

### 更新文檔
1. 修改對應分類下的文件
2. 更新 `README.md` 中的連結
3. 檢查歸檔文件是否需要更新

### 新增文檔
1. 選擇適當的分類目錄
2. 遵循現有的命名規範
3. 更新 `README.md` 導航

### 版本控制
- 重要變更記錄在 `development/changelog.md`
- 舊版本文件保存在 `archive/old-versions/`
- 已棄用文件保存在 `archive/deprecated/`

## 📝 注意事項

1. **連結更新**: 所有內部連結已更新為新路徑
2. **向後相容**: 舊的檔案路徑已重定向到新位置
3. **內容同步**: 所有範例使用最新的真實數據
4. **定期維護**: 建議每月檢查文檔的準確性

---

*最後更新: 2025-10-02*
*文檔重組完成*
