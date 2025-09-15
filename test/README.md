# 測試資料夾

本資料夾包含 Teams Notification API 的所有測試相關文件。

## 資料夾結構

```
test/
├── README.md                 # 本文件
├── api/                      # API 測試相關文件
│   └── (未來可放置 API 測試代碼)
├── scripts/                  # 測試腳本
│   └── test_api.sh          # 自動化 API 測試腳本
└── docs/                     # 測試文檔
    ├── API_TEST_GUIDE.md    # 完整的 API 測試指南
    └── TEST_README.md       # 測試文件說明
```

## 快速開始

### 運行自動化測試
```bash
# 從項目根目錄運行
./test/scripts/test_api.sh

# 或進入測試腳本目錄
cd test/scripts
./test_api.sh
```

### 查看測試文檔
```bash
# 查看完整測試指南
cat test/docs/API_TEST_GUIDE.md

# 查看測試說明
cat test/docs/TEST_README.md
```

## 測試腳本說明

### test_api.sh
自動化測試腳本，提供：
- ✅ 健康檢查測試
- ✅ 所有 API 端點測試
- ✅ 創建操作測試
- ✅ JSONB 字段測試
- ✅ 中文內容測試
- ✅ 性能測試
- ✅ 彩色輸出和錯誤處理

### 使用方法
```bash
# 確保服務器運行在 8080 端口
go run ./cmd/server

# 在另一個終端運行測試
./test/scripts/test_api.sh
```

## 測試覆蓋範圍

- **Company API**: 公司管理功能
- **User API**: 用戶管理功能
- **Project API**: 項目管理功能
- **Bot API**: 機器人管理功能
  - Platform Bot (平台機器人)
  - Third Party Bot (第三方機器人)
- **Destination API**: 目的地管理功能
- **Notification API**: 通知管理功能

## 測試數據

測試腳本會創建以下測試數據：
- 測試公司
- 測試目的地
- 測試通知
- 複雜的 JSONB 配置
- 中文內容測試

## 性能指標

正常響應時間應該在：
- 健康檢查: < 1ms
- 基本 API: < 5ms
- 複雜查詢: < 10ms

## 故障排除

### 常見問題
1. **服務器未運行**: 先啟動 `go run ./cmd/server`
2. **端口被占用**: 檢查 8080 端口
3. **權限錯誤**: 運行 `chmod +x test/scripts/test_api.sh`
4. **依賴缺失**: 確保安裝 `curl` 和 `jq`

### 日誌查看
```bash
# 查看服務器日誌
go run ./cmd/server 2>&1 | tee server.log

# 查看測試輸出
./test/scripts/test_api.sh 2>&1 | tee test.log
```

## 自定義測試

### 修改測試數據
編輯 `test/scripts/test_api.sh` 中的測試數據

### 添加新測試
在 `test/scripts/test_api.sh` 中添加新的測試函數

### 創建新的測試腳本
在 `test/scripts/` 目錄下創建新的測試腳本

## 更新日誌

- **v1.0.0** (2025-09-15): 初始版本
  - 完整的測試資料夾結構
  - 自動化測試腳本
  - 詳細的測試文檔
  - 性能監控功能
