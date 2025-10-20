# 📚 API 文檔更新指南

## 新增監控 API 端點

### 系統健康監控
```bash
# 檢查系統整體健康狀態
GET /api/v1/monitoring/health

# 響應示例
{
  "data": {
    "status": "healthy",
    "overall": {
      "score": 100,
      "grade": "A",
      "message": "System is healthy"
    }
  }
}
```

### 業務指標監控
```bash
# 獲取業務指標
GET /api/v1/monitoring/business

# 響應示例
{
  "data": {
    "active_projects": 5,
    "active_destinations": 20,
    "notifications_today": 4500,
    "success_rate_percent": 99.8,
    "queue_health": {
      "pending_count": 0,
      "processing_count": 0,
      "failed_count": 0,
      "status": "healthy"
    }
  }
}
```

### 性能指標監控
```bash
# 獲取性能指標
GET /api/v1/monitoring/performance

# 響應示例
{
  "data": {
    "system_metrics": {
      "uptime_seconds": 3600,
      "memory_usage_mb": 107,
      "cpu_usage_percent": 8.34
    },
    "redis_metrics": {
      "connected_clients": 1,
      "used_memory_mb": 2.5,
      "keys_count": 150
    }
  }
}
```

### 告警狀態監控
```bash
# 獲取告警狀態
GET /api/v1/monitoring/alerts

# 響應示例
{
  "data": {
    "active_alerts": 0,
    "alert_history": [],
    "alert_rules": [
      {
        "name": "high_error_rate",
        "threshold": 5.0,
        "status": "normal"
      }
    ]
  }
}
```

### 儀表板數據
```bash
# 獲取儀表板數據
GET /api/v1/monitoring/dashboard

# 響應示例
{
  "data": {
    "summary": {
      "total_notifications": 4500,
      "success_rate": 99.8,
      "active_projects": 5
    },
    "charts": {
      "notifications_timeline": [],
      "error_rate_trend": []
    }
  }
}
```

## 使用示例

### 監控腳本
```bash
#!/bin/bash
# 系統健康檢查腳本

API_URL="http://localhost:8080"

echo "=== 系統健康檢查 ==="
curl -s "$API_URL/api/v1/monitoring/health" | jq '.data.overall'

echo "=== 業務指標 ==="
curl -s "$API_URL/api/v1/monitoring/business" | jq '.data | {notifications_today, success_rate_percent}'

echo "=== 性能指標 ==="
curl -s "$API_URL/api/v1/monitoring/performance" | jq '.data.system_metrics'
```

### 告警配置
```yaml
# alerting.yml
alerts:
  - name: "high_error_rate"
    condition: "error_rate > 5%"
    action: "send_notification"
    
  - name: "actor_pool_full"
    condition: "active_actors >= max_actors"
    action: "scale_up"
    
  - name: "memory_usage_high"
    condition: "memory_usage > 80%"
    action: "restart_service"
```

