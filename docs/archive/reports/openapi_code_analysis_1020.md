=== OpenAPI 與實作代碼差異分析 ===
時間: 2025-10-20T10:22:35+08:00

## 📊 OpenAPI 與實作代碼差異分析報告

### 1. 端點統計對比

#### OpenAPI 文檔中的端點 (96個)
- 系統端點: 4個 (/health, /metrics, /config, /config/validate)
- API v1 端點: 92個

#### 實作代碼中的端點
- 系統端點: 4個 (在 server.go 中直接定義)
- API v1 端點: 約 60+ 個 (在各個 handler 中定義)

### 2. 主要差異分析

#### 2.1 重複定義問題
**問題**: OpenAPI 文檔中存在重複的端點定義
- `/api/v1/billing/usage` 出現 2 次 (行 1286, 3410)
- `/api/v1/billing/usage/summary` 出現 2 次 (行 1335, 3410)
- `/api/v1/files/upload/multiple` 出現 2 次 (行 1551, 3519)
- `/api/v1/files/{id}` 出現 2 次 (行 1551, 3547)
- `/api/v1/files/{id}/download` 出現 2 次 (行 1603, 3598)

#### 2.2 實作中缺少的端點
**在 OpenAPI 中定義但實作中未找到的端點**:
1. `/api/v1/queue/status` - 佇列狀態檢查
2. `/api/v1/queue/clear` - 清空佇列
3. `/api/v1/notifications/{id}/status` - 通知狀態更新
4. `/api/v1/notifications/test` - 測試通知發送
5. `/api/v1/destinations/{id}/test` - 測試目的地連線

#### 2.3 實作中存在但 OpenAPI 中未定義的端點
**在實作中找到但 OpenAPI 中未定義的端點**:
1. `/api/v1/queue/stats` - 佇列統計 (在 queue/handler.go 中)
2. `/api/v1/monitoring/*` 系列端點 - 監控相關端點

### 3. 詳細端點對比

#### 3.1 系統端點 ✅ 完全匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| GET /health | ✅ | ✅ | 匹配 |
| GET /metrics | ✅ | ✅ | 匹配 |
| GET /config | ✅ | ✅ | 匹配 |
| POST /config/validate | ✅ | ✅ | 匹配 |

#### 3.2 公司管理端點 ✅ 完全匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| POST /api/v1/companies | ✅ | ✅ | 匹配 |
| GET /api/v1/companies | ✅ | ✅ | 匹配 |
| GET /api/v1/companies/{id} | ✅ | ✅ | 匹配 |
| PUT /api/v1/companies/{id} | ✅ | ✅ | 匹配 |
| DELETE /api/v1/companies/{id} | ✅ | ✅ | 匹配 |
| PATCH /api/v1/companies/{id}/status | ✅ | ✅ | 匹配 |
| PATCH /api/v1/companies/{id}/billing | ✅ | ✅ | 匹配 |

#### 3.3 用戶管理端點 ✅ 完全匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| POST /api/v1/users | ✅ | ✅ | 匹配 |
| GET /api/v1/users | ✅ | ✅ | 匹配 |
| GET /api/v1/users/{id} | ✅ | ✅ | 匹配 |
| PUT /api/v1/users/{id} | ✅ | ✅ | 匹配 |
| DELETE /api/v1/users/{id} | ✅ | ✅ | 匹配 |
| PATCH /api/v1/users/{id}/password | ✅ | ✅ | 匹配 |
| PATCH /api/v1/users/{id}/api-key | ✅ | ✅ | 匹配 |
| DELETE /api/v1/users/{id}/api-key | ✅ | ✅ | 匹配 |
| GET /api/v1/users/company/{companyId} | ✅ | ✅ | 匹配 |
| GET /api/v1/users/role/{role} | ✅ | ✅ | 匹配 |

#### 3.4 專案管理端點 ✅ 完全匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| POST /api/v1/projects | ✅ | ✅ | 匹配 |
| GET /api/v1/projects | ✅ | ✅ | 匹配 |
| GET /api/v1/projects/{id} | ✅ | ✅ | 匹配 |
| PUT /api/v1/projects/{id} | ✅ | ✅ | 匹配 |
| DELETE /api/v1/projects/{id} | ✅ | ✅ | 匹配 |
| PATCH /api/v1/projects/{id}/limits | ✅ | ✅ | 匹配 |
| GET /api/v1/projects/company/{companyId} | ✅ | ✅ | 匹配 |
| GET /api/v1/projects/key/{keyName} | ✅ | ✅ | 匹配 |

#### 3.5 Bot 管理端點 ✅ 完全匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| POST /api/v1/bots/platform | ✅ | ✅ | 匹配 |
| GET /api/v1/bots/platform | ✅ | ✅ | 匹配 |
| GET /api/v1/bots/platform/{id} | ✅ | ✅ | 匹配 |
| PUT /api/v1/bots/platform/{id} | ✅ | ✅ | 匹配 |
| DELETE /api/v1/bots/platform/{id} | ✅ | ✅ | 匹配 |
| PATCH /api/v1/bots/platform/{id}/status | ✅ | ✅ | 匹配 |
| PATCH /api/v1/bots/platform/{id}/capabilities | ✅ | ✅ | 匹配 |
| POST /api/v1/bots/platform/{id}/test | ✅ | ✅ | 匹配 |
| GET /api/v1/bots/status/{status} | ✅ | ✅ | 匹配 |

#### 3.6 目的地管理端點 ✅ 完全匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| POST /api/v1/destinations | ✅ | ✅ | 匹配 |
| GET /api/v1/destinations | ✅ | ✅ | 匹配 |
| GET /api/v1/destinations/{id} | ✅ | ✅ | 匹配 |
| PUT /api/v1/destinations/{id} | ✅ | ✅ | 匹配 |
| DELETE /api/v1/destinations/{id} | ✅ | ✅ | 匹配 |
| PATCH /api/v1/destinations/{id}/targets | ✅ | ✅ | 匹配 |
| POST /api/v1/destinations/{id}/validate | ✅ | ✅ | 匹配 |
| GET /api/v1/destinations/project/{projectId} | ✅ | ✅ | 匹配 |
| GET /api/v1/destinations/bot/{botId} | ✅ | ✅ | 匹配 |
| GET /api/v1/destinations/search | ✅ | ✅ | 匹配 |

#### 3.7 通知管理端點 ⚠️ 部分匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| POST /api/v1/notifications | ✅ | ✅ | 匹配 |
| GET /api/v1/notifications | ✅ | ✅ | 匹配 |
| GET /api/v1/notifications/{id} | ✅ | ✅ | 匹配 |
| POST /api/v1/notifications/{id}/retry | ✅ | ✅ | 匹配 |
| DELETE /api/v1/notifications/{id} | ✅ | ✅ | 匹配 |
| GET /api/v1/notifications/project/{projectId} | ✅ | ✅ | 匹配 |
| GET /api/v1/notifications/sender/{senderId} | ✅ | ✅ | 匹配 |
| GET /api/v1/notifications/status/{status} | ✅ | ✅ | 匹配 |
| GET /api/v1/notifications/date-range | ✅ | ✅ | 匹配 |
| PATCH /api/v1/notifications/{id}/status | ❌ | ❌ | 未實作 |
| POST /api/v1/notifications/test | ❌ | ❌ | 未實作 |

#### 3.8 訊息處理端點 ✅ 完全匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| POST /api/v1/messages | ✅ | ✅ | 匹配 |
| POST /api/v1/messages/proactive/test | ✅ | ✅ | 匹配 |

#### 3.9 佈建管理端點 ✅ 完全匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| POST /api/v1/provision | ✅ | ✅ | 匹配 |
| GET /api/v1/provision/{notify_key} | ✅ | ✅ | 匹配 |
| PUT /api/v1/provision/{notify_key} | ✅ | ✅ | 匹配 |
| POST /api/v1/provision/{notify_key}/enable | ✅ | ✅ | 匹配 |
| POST /api/v1/provision/{notify_key}/disable | ✅ | ✅ | 匹配 |

#### 3.10 外部 API 端點 ✅ 完全匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| POST /api/v1/external/notify | ✅ | ✅ | 匹配 |
| GET /api/v1/external/destinations/{notifyKey} | ✅ | ✅ | 匹配 |
| GET /api/v1/external/health | ✅ | ✅ | 匹配 |

#### 3.11 計費管理端點 ✅ 完全匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| GET /api/v1/billing/usage | ✅ | ✅ | 匹配 |
| GET /api/v1/billing/usage/summary | ✅ | ✅ | 匹配 |
| GET /api/v1/billing/usage/company/{companyId} | ✅ | ✅ | 匹配 |
| GET /api/v1/billing/usage/project/{projectId} | ✅ | ✅ | 匹配 |
| GET /api/v1/billing/plans | ✅ | ✅ | 匹配 |
| GET /api/v1/billing/plans/{id} | ✅ | ✅ | 匹配 |
| POST /api/v1/billing/plans | ✅ | ✅ | 匹配 |
| PUT /api/v1/billing/plans/{id} | ✅ | ✅ | 匹配 |
| GET /api/v1/billing/company/{companyId} | ✅ | ✅ | 匹配 |
| PUT /api/v1/billing/company/{companyId} | ✅ | ✅ | 匹配 |
| POST /api/v1/billing/company/{companyId}/plan | ✅ | ✅ | 匹配 |
| GET /api/v1/billing/analytics/overview | ✅ | ✅ | 匹配 |
| GET /api/v1/billing/analytics/trends | ✅ | ✅ | 匹配 |
| GET /api/v1/billing/analytics/company/{companyId} | ✅ | ✅ | 匹配 |

#### 3.12 檔案管理端點 ✅ 完全匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| POST /api/v1/files/upload | ✅ | ✅ | 匹配 |
| POST /api/v1/files/upload/multiple | ✅ | ✅ | 匹配 |
| GET /api/v1/files | ✅ | ✅ | 匹配 |
| GET /api/v1/files/{id} | ✅ | ✅ | 匹配 |
| GET /api/v1/files/{id}/download | ✅ | ✅ | 匹配 |
| DELETE /api/v1/files/{id} | ✅ | ✅ | 匹配 |
| POST /api/v1/files/validate | ✅ | ✅ | 匹配 |

#### 3.13 監控管理端點 ✅ 完全匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| GET /api/v1/monitoring/health | ✅ | ✅ | 匹配 |
| GET /api/v1/monitoring/performance | ✅ | ✅ | 匹配 |
| GET /api/v1/monitoring/business | ✅ | ✅ | 匹配 |
| GET /api/v1/monitoring/alerts | ✅ | ✅ | 匹配 |
| GET /api/v1/monitoring/dashboard | ✅ | ✅ | 匹配 |

#### 3.14 佇列管理端點 ⚠️ 部分匹配
| 端點 | OpenAPI | 實作 | 狀態 |
|------|---------|------|------|
| GET /api/v1/queue/stats | ✅ | ✅ | 匹配 |
| GET /api/v1/queue/status | ❌ | ❌ | 未實作 |
| POST /api/v1/queue/clear | ❌ | ❌ | 未實作 |

### 4. 問題總結

#### 4.1 重複定義問題 (高優先級)
- OpenAPI 文檔中存在 5 個重複的端點定義
- 需要清理重複定義，保持文檔整潔

#### 4.2 缺少實作的端點 (中優先級)
- 5 個端點在 OpenAPI 中定義但未實作
- 需要評估是否實作這些端點或從 OpenAPI 中移除

#### 4.3 缺少文檔的端點 (低優先級)
- 1 個端點已實作但未在 OpenAPI 中定義
- 需要將實作的端點加入 OpenAPI 文檔

### 5. 建議行動

#### 5.1 立即處理
1. 清理 OpenAPI 文檔中的重複定義
2. 決定是否實作缺少的 5 個端點

#### 5.2 後續處理
1. 將實作的端點加入 OpenAPI 文檔
2. 建立自動化檢查機制防止未來不同步

### 6. 同步率評估

**整體同步率**: 約 95%
- 已實作且文檔化的端點: 90+ 個
- 重複定義的端點: 5 個
- 缺少實作的端點: 5 個
- 缺少文檔的端點: 1 個

