# Backend-Frontend 集成差异分析

## 1. 路由路径差异

### OpenAPI 定义 vs Backend 实际路由

| OpenAPI 定义 | Backend 实际路径 | 状态 | 说明 |
|-------------|-----------------|------|------|
| `POST /messages/send` | `POST /api/v1/notify` | ❌ 不匹配 | 需要适配或创建新端点 |
| `GET /projects` | `GET /internal/v1/projects` | ✅ 匹配 | Frontend 已正确调用 |
| `GET /projects/{projectId}/logs` | `GET /internal/v1/notifications/project/{projectId}` | ❌ 不匹配 | 路径不同 |
| `GET /admin/health` | `GET /internal/v1/monitoring/health` | ❌ 不匹配 | 路径不同 |
| `GET /admin/billing/records` | `GET /internal/v1/billing/usage/summary` | ❌ 不匹配 | 路径和响应格式不同 |
| `GET /projects/{projectId}/templates` | `GET /internal/v1/projects/{id}/templates` | ✅ 匹配 | 参数名不同 (projectId vs id) |
| `GET /directory/users/search` | `GET /internal/v1/directory/users/search` | ✅ 匹配 | 正确 |
| `GET /directory/groups` | `GET /internal/v1/directory/groups` | ✅ 匹配 | 正确 |
| `GET /directory/groups/{groupId}/channels` | `GET /internal/v1/directory/groups/{groupId}/channels` | ✅ 匹配 | 正确 |
| `GET /projects/{projectId}/chat-groups` | `GET /internal/v1/projects/{id}/chat-groups` | ⚠️ 部分匹配 | 参数名不同 |

## 2. 数据模型差异

### 字段命名
- **Backend**: snake_case (`company_id`, `notify_key`, `created_at`)
- **Frontend**: camelCase (`companyId`, `notifyKey`, `createdAt`)
- **状态**: Frontend 已有转换函数，但需要完善

### 响应格式
Backend 使用三种格式：
1. `response.Success()`: `{ code: 200, message: "...", data: {...} }`
2. `gin.H`: `{ success: true, data: {...} }`
3. `gin.H`: `{ data: [...], pagination: {...} }`

Frontend 需要统一处理这些格式。

## 3. 缺失的功能

### Backend 需要实现的端点
- [ ] `POST /messages/send` (OpenAPI 定义但不存在，可用 `/api/v1/notify` 替代)
- [ ] `POST /projects/{projectId}/regenerate-key` (OpenAPI 定义但不存在)
- [ ] `GET /admin/bots` (OpenAPI 定义，backend 有 `/internal/v1/bots/platform`)

### Frontend 需要修复的问题
- [ ] API 路径不匹配（特别是 `/projects/{projectId}/logs`）
- [ ] 参数名不一致（`projectId` vs `id`）
- [ ] 响应格式处理不完整
- [ ] 认证 token 处理（Bearer token, API key）

## 4. 当前集成状态

### ✅ 已正确集成的功能
- Projects CRUD
- Templates CRUD
- Directory (users, groups, channels)
- Chat Groups
- Notifications (部分)

### ⚠️ 需要修复的功能
- Message sending (路径不匹配)
- Project logs (路径不匹配)
- System health (路径不匹配)
- Billing records (路径和格式不匹配)
- API key regeneration (缺失)

### ❌ 未实现的功能
- Audience lists (mock only)
- User tagging (mock only)
- Global templates (mock only)

## 5. 修复优先级

### 高优先级
1. 修复 API 路径映射
2. 统一响应格式处理
3. 修复参数名不一致问题

### 中优先级
4. 实现缺失的端点
5. 添加认证处理
6. 完善错误处理

### 低优先级
7. 实现 mock-only 功能
8. 性能优化
9. 文档完善
