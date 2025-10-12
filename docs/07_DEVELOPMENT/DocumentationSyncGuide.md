# 文檔同步機制指南

## 概述

本文檔描述了 Teams Notification API 專案的文檔同步機制，確保文檔與實際實現保持同步。

## 同步檢查工具

### 檢查腳本

**位置**: `scripts/utils/check-doc-sync.sh`

**功能**:
- 檢查 API 服務運行狀態
- 驗證監控端點可用性
- 檢查文檔中的端點路徑正確性
- 生成同步檢查報告

**使用方法**:
```bash
# 運行文檔同步檢查
./scripts/utils/check-doc-sync.sh

# 檢查結果會生成 doc_sync_report.md
```

### 檢查項目

#### 1. API 端點檢查
- `/health` - 基本健康檢查
- `/api/v1/metrics` - 系統指標
- `/api/v1/config` - 配置信息
- `/api/v1/monitoring/health` - 監控健康檢查
- `/api/v1/monitoring/performance` - 性能指標
- `/api/v1/monitoring/business` - 業務指標
- `/api/v1/monitoring/alerts` - 警報狀態
- `/api/v1/monitoring/dashboard` - 監控儀表板
- `/api/v1/queue/stats` - 佇列統計

#### 2. 文檔檢查
- `docs/05_DEPLOYMENT/MonitoringGuide.md`
- `docs/05_DEPLOYMENT/MonitoringConfiguration.md`
- `docs/06_USER_GUIDE/UserManual.md`
- `docs/02_ARCHITECTURE/Architecture.md`

## 同步原則

### 1. 端點路徑同步
- 文檔中的 API 端點路徑必須與實際實現一致
- 新增端點時，必須同步更新相關文檔
- 移除端點時，必須從文檔中移除相關說明

### 2. 功能描述同步
- 文檔中的功能描述必須與實際實現一致
- API 參數、回應格式必須準確
- 使用範例必須可執行

### 3. 架構文檔同步
- 系統架構圖必須反映實際實現
- 技術棧描述必須準確
- 組件關係必須正確

## 同步流程

### 1. 開發階段
- 新增功能時，同步更新相關文檔
- 修改 API 時，更新 API 文檔
- 變更架構時，更新架構文檔

### 2. 測試階段
- 運行文檔同步檢查腳本
- 修復發現的同步問題
- 驗證文檔中的範例可執行

### 3. 發布階段
- 最終檢查文檔同步狀態
- 生成文檔同步報告
- 確認所有文檔與實現一致

## 文檔更新檢查清單

### API 文檔更新
- [ ] 端點路徑正確
- [ ] 請求參數完整
- [ ] 回應格式準確
- [ ] 使用範例可執行
- [ ] 錯誤處理說明完整

### 架構文檔更新
- [ ] 架構圖反映實際實現
- [ ] 技術棧描述準確
- [ ] 組件關係正確
- [ ] 資料流程描述完整

### 使用者指南更新
- [ ] 操作步驟準確
- [ ] 範例代碼可執行
- [ ] 故障排除說明完整
- [ ] 監控使用指南完整

## 自動化檢查

### 定期檢查
建議每週運行一次文檔同步檢查：

```bash
# 每週一運行檢查
./scripts/utils/check-doc-sync.sh
```

### CI/CD 整合
可以將文檔同步檢查整合到 CI/CD 流程中：

```yaml
# .github/workflows/doc-sync.yml
name: Documentation Sync Check
on:
  schedule:
    - cron: '0 9 * * 1'  # 每週一上午 9 點
  workflow_dispatch:

jobs:
  doc-sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Check Documentation Sync
        run: ./scripts/utils/check-doc-sync.sh
```

## 常見問題

### 1. 端點路徑不匹配
**問題**: 文檔中的端點路徑與實際實現不一致
**解決**: 更新文檔中的端點路徑，確保與實際實現一致

### 2. 功能描述過時
**問題**: 文檔中的功能描述與實際實現不符
**解決**: 更新文檔中的功能描述，確保準確反映實際實現

### 3. 範例代碼不可執行
**問題**: 文檔中的範例代碼無法執行
**解決**: 測試並修正範例代碼，確保可執行

## 最佳實踐

### 1. 即時更新
- 新增功能時立即更新文檔
- 修改 API 時同步更新文檔
- 避免累積大量更新

### 2. 版本控制
- 文檔更新與代碼更新同步提交
- 使用有意義的提交訊息
- 標記文檔版本

### 3. 定期檢查
- 每週運行同步檢查
- 每月進行文檔審查
- 季度進行全面檢查

## 工具和資源

### 檢查工具
- `scripts/utils/check-doc-sync.sh` - 文檔同步檢查腳本
- `doc_sync_report.md` - 同步檢查報告

### 相關文檔
- [API 文檔](../06_USER_GUIDE/external-api.md)
- [監控指南](../05_DEPLOYMENT/MonitoringGuide.md)
- [架構文檔](../02_ARCHITECTURE/Architecture.md)

## 總結

文檔同步機制是確保專案文檔品質的重要工具。通過定期檢查和即時更新，可以確保文檔與實際實現保持同步，提供準確的使用指南和技術文檔。
