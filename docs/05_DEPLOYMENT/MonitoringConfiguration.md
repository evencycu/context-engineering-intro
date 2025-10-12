# 監控配置指南

## 概述

本文檔描述了 Teams Notify 系統的監控配置，包括系統健康檢查、性能指標、業務指標和警報設置。

## 監控端點

### 1. 系統健康檢查
- **端點**: `GET /api/v1/monitoring/health`
- **描述**: 檢查系統整體健康狀態
- **響應**: 包含數據庫、Redis、API 服務的健康狀態
- **範例**: 
```bash
curl http://localhost:8080/api/v1/monitoring/health | jq '.'
```

### 2. 性能指標
- **端點**: `GET /api/v1/monitoring/performance`
- **描述**: 獲取系統性能指標
- **包含**: 響應時間、吞吐量、錯誤率、資源使用情況
- **範例**:
```bash
curl http://localhost:8080/api/v1/monitoring/performance | jq '.'
```

### 3. 業務健康狀態
- **端點**: `GET /api/v1/monitoring/business`
- **描述**: 獲取業務層面的健康指標
- **包含**: 活躍專案、目的地、通知統計、隊列健康狀態
- **範例**:
```bash
curl http://localhost:8080/api/v1/monitoring/business | jq '.'
```

### 4. 警報狀態
- **端點**: `GET /api/v1/monitoring/alerts`
- **描述**: 獲取當前活躍警報
- **包含**: 警報類型、嚴重程度、組件信息
- **範例**:
```bash
curl http://localhost:8080/api/v1/monitoring/alerts | jq '.'
```

### 5. 監控儀表板
- **端點**: `GET /api/v1/monitoring/dashboard`
- **描述**: 綜合監控儀表板
- **包含**: 所有監控數據的綜合視圖
- **範例**:
```bash
curl http://localhost:8080/api/v1/monitoring/dashboard | jq '.'
```

## 指標說明

### 系統指標
- **Uptime**: 系統運行時間
- **Memory Usage**: 內存使用量 (MB)
- **CPU Usage**: CPU 使用率 (%)
- **Goroutines**: Go 協程數量
- **Database Connections**: 數據庫連接數
- **Redis Connections**: Redis 連接數

### 業務指標
- **Notifications Sent**: 已發送通知數
- **Notifications Failed**: 失敗通知數
- **Notifications Pending**: 待處理通知數
- **Success Rate**: 成功率 (%)
- **Average Response Time**: 平均響應時間 (ms)
- **Teams API Calls**: Teams API 調用次數
- **Teams API Errors**: Teams API 錯誤次數
- **Queue Processing Rate**: 隊列處理速率 (每分鐘)
- **Active Projects**: 活躍專案數
- **Active Destinations**: 活躍目的地數

### 性能指標
- **Response Time**: 響應時間統計 (P50, P95, P99)
- **Throughput**: 吞吐量統計
- **Error Rate**: 錯誤率統計
- **Resource Usage**: 資源使用統計

## 警報配置

### 警報類型
1. **Performance Alerts**
   - 高錯誤率 (成功率 < 90%)
   - 高響應時間 (平均響應時間 > 5000ms)

2. **Integration Alerts**
   - Teams API 錯誤 (錯誤次數 > 10)

3. **System Alerts**
   - 數據庫連接問題
   - Redis 連接問題
   - API 服務問題

### 警報嚴重程度
- **Critical**: 系統關鍵問題
- **Warning**: 性能或業務問題
- **Info**: 信息性警報

## 監控最佳實踐

### 1. 監控頻率
- 健康檢查: 每 30 秒
- 性能指標: 每 5 分鐘
- 業務指標: 每 1 分鐘
- 警報檢查: 每 10 秒

### 2. 閾值設置
- 成功率: < 90% (警告)
- 響應時間: > 5000ms (警告)
- Teams API 錯誤: > 10 次 (嚴重)
- 內存使用: > 80% (警告)
- CPU 使用: > 80% (警告)

### 3. 日誌記錄
- 所有監控端點訪問都會記錄
- 業務事件會記錄詳細信息
- 錯誤和警告會記錄堆棧跟踪

## 集成建議

### 1. 外部監控系統
- **Prometheus**: 用於指標收集
- **Grafana**: 用於可視化
- **AlertManager**: 用於警報管理

### 2. 日誌聚合
- **ELK Stack**: Elasticsearch, Logstash, Kibana
- **Fluentd**: 日誌收集和轉發

### 3. 警報通知
- **Slack**: 即時通知
- **Email**: 郵件通知
- **PagerDuty**: 緊急警報

## 配置示例

### Docker Compose 監控配置
```yaml
version: '3.8'
services:
  prometheus:
    image: prom/prometheus
    ports:
      - "9090:9090"
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
  
  grafana:
    image: grafana/grafana
    ports:
      - "3000:3000"
    environment:
      - GF_SECURITY_ADMIN_PASSWORD=admin
```

### 監控腳本示例
```bash
#!/bin/bash
# 監控健康檢查腳本

API_BASE="http://localhost:8080/api/v1"

# 檢查系統健康
curl -s "$API_BASE/monitoring/health" | jq '.data.status'

# 檢查警報
ALERTS=$(curl -s "$API_BASE/monitoring/alerts" | jq '.data.active_alerts')
if [ "$ALERTS" -gt 0 ]; then
    echo "警告: 發現 $ALERTS 個活躍警報"
fi

# 檢查性能指標
curl -s "$API_BASE/monitoring/performance" | jq '.data.error_rate.success_rate_percent'
```

## 故障排除

### 常見問題
1. **監控端點無響應**: 檢查服務是否運行
2. **指標數據不準確**: 檢查 Redis 連接
3. **警報不觸發**: 檢查閾值設置
4. **日誌不記錄**: 檢查日誌配置

### 調試命令
```bash
# 檢查服務狀態
curl http://localhost:8080/health

# 檢查監控端點
curl http://localhost:8080/api/v1/monitoring/health

# 檢查指標
curl http://localhost:8080/api/v1/metrics

# 檢查日誌
docker logs teamsnotify-api-server
```

## 總結

監控系統提供了全面的系統可視性，包括：
- 實時健康狀態監控
- 性能指標追蹤
- 業務指標分析
- 自動警報機制
- 綜合儀表板視圖

通過正確配置和使用這些監控功能，可以確保系統的穩定運行和快速問題識別。
