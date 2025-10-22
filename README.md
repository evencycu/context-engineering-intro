# Teams Notification API Server

A comprehensive Microsoft Teams notification system that provides pre-registered destination management, multi-target messaging, rate limiting compliance, and enterprise-grade features.

## 🚀 Features

- **Pre-registered Destinations**: Secure notification key mapping system
- **Multi-target Support**: Channels, group chats, and individual users
- **Teams API Compliance**: 50 RPS global, 7 RPS per conversation rate limiting
- **Enterprise Ready**: Authentication, authorization, audit trails
- **Usage Tracking**: Comprehensive billing and cost tracking
- **Monitoring**: Health checks, metrics, and observability
- **Scheduled Messaging**: Queue-based message scheduling
- **🆕 Retry Queue & Circuit Breaker**: Automatic retry with exponential backoff and circuit breaker protection
- **🆕 Rate Limit Handling**: Smart handling of Teams API 429 errors with Retry-After support

## 📋 Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Microsoft Teams Bot Framework registration
- PostgreSQL and Redis (or use Docker Compose)

## 🏗️ Architecture

- **Framework**: Go with Gin web framework
- **Database**: PostgreSQL with GORM ORM
- **Cache/Queue**: Redis for caching and message queuing
- **Authentication**: JWT-based with Teams Bot Framework validation
- **Configuration**: Viper with environment variable support
- **Project Structure**: Refactored with `libs/` for shared libraries and `services/` for business logic

## ⚡ Quick Start

### 1. Clone and Setup

```bash
git clone <repository-url>
cd TeamsNotifyGo
```

### 2. Configure Environment

```bash
# Copy environment template
cp .env.example .env

# Edit .env with your Teams Bot credentials:
# - TEAMS_NOTIFY_TEAMS_APP_ID
# - TEAMS_NOTIFY_TEAMS_APP_SECRET  
# - TEAMS_NOTIFY_TEAMS_BOT_ID
# - TEAMS_NOTIFY_AUTH_JWT_SECRET (generated)
# - TEAMS_NOTIFY_DATABASE_PASSWORD
```

### 3. Start with Docker Compose

```bash
# Start all services (from repo root)
docker compose -f deployments/docker/docker-compose.yml up -d

# Or start with admin tools
docker compose -f deployments/docker/docker-compose.yml --profile admin up -d

# Check health
curl http://localhost:8080/health
```

## 🛠️ Development Setup

### Option 1: Full Docker Environment (Recommended)

```bash
# Start only database and Redis services
docker compose -f deployments/docker/docker-compose.yml up -d postgres redis

# Run the Go server locally for debugging
go run ./cmd/server
```

### Option 2: Hot Reload Development

```bash
# Install Air for hot reloading
go install github.com/cosmtrek/air@latest

# Start database services
docker compose -f deployments/docker/docker-compose.yml up -d postgres redis

# Run with hot reload
make dev
```

### Start API Server with Teams Bot credentials and log to file

每次啟動本機伺服器，請先設定 Teams Bot 相關環境變數，並將日誌寫入 `/tmp/teamsnotify_server.log`：

```bash
export TEAMS_BOT_APP_ID=844146d7-4ac9-4e4d-a463-d6e027714e81
export TEAMS_TENANT_ID=051cece0-e4dc-4aed-b471-bf29824e1ee6
export TEAMS_BOT_APP_PASSWORD='HVW8Q~-HmVPi_W0EzIFPbZiL4G1czV8WBMUjQdfy'

nohup env \
  TEAMS_BOT_APP_ID=$TEAMS_BOT_APP_ID \
  TEAMS_TENANT_ID=$TEAMS_TENANT_ID \
  TEAMS_BOT_APP_PASSWORD=$TEAMS_BOT_APP_PASSWORD \
  go run ./cmd/server > /tmp/teamsnotify_server.log 2>&1 & echo $!
```

查看日誌：

```bash
tail -f /tmp/teamsnotify_server.log
```

### Option 3: VS Code Debugging

The project includes VS Code debugging configuration:

1. **F5** to start debugging
2. Automatically starts PostgreSQL and Redis containers
3. Sets development environment variables
4. Enables breakpoints and step-through debugging

### Option 4: Pure Local Development

```bash
# Make sure PostgreSQL and Redis are running locally
# Set local environment variables:
export TEAMS_NOTIFY_DATABASE_HOST=localhost
export TEAMS_NOTIFY_REDIS_HOST=localhost

# Run the server
make run
```

## 🧪 Testing and Validation

### Automated API Testing

```bash
# Run comprehensive API tests
./test/scripts/test_api.sh

# View detailed test documentation
cat test/docs/API_TEST_GUIDE.md
```

### Development Testing

```bash
# Run all validations (format, vet, lint, test)
make validate

# Run tests only
make test

# Run with verbose output
make test-verbose

# Build the server binary
make build

# Check environment variables
make env-check
```

### Test Coverage

The project includes comprehensive testing for:
- ✅ **All 7 API modules** (Company, User, Project, Bot, Destination, Notification)
- ✅ **CRUD operations** and special functions
- ✅ **JSONB field handling** for complex data structures
- ✅ **Chinese content support** for internationalization
- ✅ **Performance monitoring** with response time tracking
- ✅ **Error handling** and edge cases

## 🐳 Docker Commands

```bash
# Build Docker image
make docker-build

# Run in Docker container
make docker-run

# Start full environment
docker compose -f deployments/docker/docker-compose.yml up -d

# View logs
docker compose -f deployments/docker/docker-compose.yml logs -f teams-notify-server

# Stop services
docker compose -f deployments/docker/docker-compose.yml down
```

## 📊 Monitoring and Admin Tools

### Health Endpoints

- **Health Check**: `GET /health`
- **Database Stats**: Available via health endpoint
- **Service Status**: `make status`

### Admin Interfaces (with --profile admin)

- **pgAdmin**: http://localhost:5050
  - Email: `admin@teamsnotify.local`
  - Password: `admin_password`

- **Redis Commander**: http://localhost:8081
  - Username: `admin`
  - Password: `admin_password`

## 🔧 Configuration

### Required Environment Variables

```bash
# Teams Bot Configuration
TEAMS_NOTIFY_TEAMS_APP_ID=your-bot-app-id
TEAMS_NOTIFY_TEAMS_APP_SECRET=your-bot-app-secret
TEAMS_NOTIFY_TEAMS_BOT_ID=28:your-bot-app-id

# Authentication
TEAMS_NOTIFY_AUTH_JWT_SECRET=your-jwt-secret

# Database
TEAMS_NOTIFY_DATABASE_PASSWORD=your-database-password
```

### Optional Configuration

```bash
# Server
TEAMS_NOTIFY_SERVER_PORT=8080
TEAMS_NOTIFY_SERVER_ENVIRONMENT=development

# Database
TEAMS_NOTIFY_DATABASE_HOST=localhost
TEAMS_NOTIFY_DATABASE_PORT=5432
TEAMS_NOTIFY_DATABASE_USERNAME=teams_notify
TEAMS_NOTIFY_DATABASE_DATABASE_NAME=teams_notify

# Redis
TEAMS_NOTIFY_REDIS_HOST=localhost
TEAMS_NOTIFY_REDIS_PORT=6379
TEAMS_NOTIFY_REDIS_PASSWORD=redis_password

# Logging
TEAMS_NOTIFY_LOGGING_LEVEL=info
TEAMS_NOTIFY_LOGGING_FORMAT=json
```

## 🚦 Rate Limiting

The system complies with Microsoft Teams API rate limits:

- **Global Rate Limit**: 50 requests per second per application
- **Conversation Rate Limit**: 7 requests per second per conversation
- **Automatic Retry**: Exponential backoff with jitter
- **Queue Management**: Redis-based message queuing

## 📡 API Documentation

### Interactive API Documentation

The API includes comprehensive OpenAPI 3.0 documentation with interactive testing:

- **Swagger UI**: http://localhost:8080/docs/
- **OpenAPI Spec**: http://localhost:8080/docs/doc.json
- **ReDoc**: Available via Swagger UI interface

### Generate Documentation

```bash
# Generate API documentation
make docs

# Generate and serve documentation locally
make docs-serve
```

### Core Endpoints

- `GET /health` - Health check and system status
- `GET /ping` - Simple connectivity test
- `GET /docs/` - Interactive API documentation
- `POST /internal/v1/destinations` - Register notification destinations
- `POST /internal/v1/notifications` - Send notifications
- `GET /internal/v1/notifications/{id}` - Get notification status
- `POST /internal/v1/notifications/batch` - Batch notification sending

### Authentication

All API endpoints require JWT authentication:

```bash
Authorization: Bearer <jwt-token>
```

## 📝 Debugging Tips

### Local Development

```bash
# Set debug logging
export TEAMS_NOTIFY_LOGGING_LEVEL=debug
export TEAMS_NOTIFY_SERVER_ENVIRONMENT=development

# Start dependencies
docker-compose up -d postgres redis

# Run locally
go run ./cmd/server

# Check health
curl http://localhost:8080/health
```

### Common Issues

1. **Configuration Errors**: Check environment variables match required prefixes
2. **Database Connection**: Ensure PostgreSQL is running and credentials are correct
3. **Teams Authentication**: Verify Bot Framework credentials and permissions
4. **Rate Limiting**: Check Redis connection for rate limiting functionality

### Logging

- All logs are in JSON format for structured logging
- Use `TEAMS_NOTIFY_LOGGING_LEVEL=debug` for verbose output
- Logs include request IDs for tracing

## 🏭 Production Deployment

### Docker Production

```bash
# Build production image (from repo root)
docker build -f deployments/docker/Dockerfile -t teams-notify-server:latest .

# Run with production environment
docker run -d \
  --name teams-notify-server \
  -p 8080:8080 \
  --env-file .env \
  teams-notify-server:latest
```

### Security Considerations

- Use strong JWT secrets (64+ characters)
- Enable TLS in production
- Secure database and Redis connections
- Implement proper firewall rules
- Regular security updates

## 📚 Development Resources

### Microsoft Teams Bot Framework

- [Bot Framework Documentation](https://docs.microsoft.com/en-us/azure/bot-service/)
- [Teams API Rate Limits](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/how-to/rate-limit)
- [Proactive Messaging](https://docs.microsoft.com/en-us/microsoftteams/platform/bots/how-to/conversations/send-proactive-messages)

### Go Libraries Used

- [Gin Web Framework](https://github.com/gin-gonic/gin)
- [GORM ORM](https://gorm.io/docs/)
- [Viper Configuration](https://github.com/spf13/viper)
- [JWT Golang](https://github.com/golang-jwt/jwt)

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make changes with tests
4. Run validation: `make validate`
5. Submit a pull request

## 📚 Additional Documentation

- **[系統架構文檔](./docs/02_ARCHITECTURE/Architecture.md)** - 完整的系統架構說明
- **[佇列系統詳述](./docs/02_ARCHITECTURE/queue-system-detailed.md)** - 佇列和熔斷器機制說明
- **[Teams Token 快取](./docs/02_ARCHITECTURE/teams-token-cache-detailed.md)** - Teams Token 快取機制
- **[系統設計文檔](./docs/03_DESIGN/SystemDesign.md)** - 系統設計詳述
- **[資料庫設計](./docs/03_DESIGN/ERD.md)** - 資料庫實體關係圖
- **[測試指南](./docs/04_TEST/TestPlan.md)** - 完整的測試計劃
- **[部署指南](./docs/05_DEPLOYMENT/DeploymentGuide.md)** - 部署和配置指南
- **[用戶手冊](./docs/06_USER_GUIDE/UserManual.md)** - 用戶使用指南
- **[開發指南](./docs/07_DEVELOPMENT/DevelopmentGuide.md)** - 開發環境和流程指南
- **[重構指南](./docs/07_DEVELOPMENT/RefactoringGuide.md)** - 目錄結構重構說明

## 🔄 Queue & Circuit Breaker (新功能)

系統現在包含自動重試隊列和熔斷器保護機制，用於處理 Teams API 的速率限制和暫時性錯誤。

### 快速測試

```bash
# 1. 執行 migration 建立隊列表
docker exec -i teamsnotify-postgres psql -U teamsnotify -d notification_center < scripts/migrations/009_enhance_notification_destinations_for_async_actor.sql

# 2. 啟動服務（會自動啟動 Actor Pool）
./server

# 3. 測試 Queue API
./scripts/test_queue.sh
```

### Queue 監控 API

```bash
# 查看隊列狀態
curl http://localhost:8080/internal/v1/queue/status

# 查看熔斷器指標
curl http://localhost:8080/internal/v1/queue/circuit-breaker/metrics

# 重置熔斷器
curl -X POST http://localhost:8080/internal/v1/queue/circuit-breaker/reset
```

### 主要特性

- ✅ **自動重試**: 失敗的通知自動進入重試隊列
- ✅ **指數退避**: 使用指數退避 + 隨機抖動避免雪崩
- ✅ **429 處理**: 智能處理 Teams API 的 Retry-After header
- ✅ **熔斷器保護**: 防止持續向不可用服務發送請求
- ✅ **持久化隊列**: 使用資料庫持久化，防止數據丟失
- ✅ **可觀測性**: 提供監控 API 和指標

詳細說明請參考: [佇列系統詳述](./docs/02_ARCHITECTURE/queue-system-detailed.md)

## 📄 License

[Your License Here]

## 🆘 Support

For issues and questions:

1. Check the troubleshooting section above
2. Review logs for error details
3. Ensure all environment variables are set correctly
4. Verify Teams Bot Framework configuration
5. 查看 Queue 文檔瞭解重試機制