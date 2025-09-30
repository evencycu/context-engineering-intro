# Architecture Overview

## Runtime Composition
- **Entry point**: `cmd/server/main.go` wires configuration, logging, rate limiting, database connection, and registers handlers.
- **HTTP server**: `internal/api/server.go` builds a Gin engine with CORS, request ID, JSON logging, rate limiting, and graceful shutdown.
- **Middleware**: `internal/api/middleware` provides auth (JWT-ready, API key TODO), structured logging, error handling, and security headers.
- **Handlers**: `internal/api/handlers/*` implement thin REST adapters for companies, users, projects, Teams bots, destinations, notifications, message webhooks, and provisioning flows.
- **Services**: `internal/api/services` encapsulate domain logic (CRUD helpers, password hashing, notification broadcasting, Teams activity ingestion, provisioning).
- **Repositories**: `internal/api/repositories` house `sqlx` backed CRUD implementations per aggregate plus helpers for pagination and lookups.
- **Database layer**: `internal/database/models.go` defines entity structs + JSONB helpers; `internal/database/schema.sql` declares PostgreSQL schema and enums.

## Supporting Assets
- **API contracts**: `api/openapi/teams-notification-api.yaml` and `api/proto/notification.proto` describe REST + future gRPC surfaces.
- **Configuration**: `configs/*.yaml` express server/database/cache/auth defaults; `config.example.env` lists environment variables.
- **Docs**: `docs/*.md` cover requirements, ERD, testing, installation, and high-level design notes. `expose_doc/` aggregates curated exports.
- **Tooling**: `scripts/*.sh` provide build, lint, OpenAPI serving, docker helpers, and SQL fixtures. `Makefile` wraps common workflows.
- **Deployments**: `deployments/` contains docker-compose and Kubernetes-style manifests for API, worker, admin, and infra services.

## Execution Flow (Happy Path Notification)
1. Client calls `POST /api/v1/notifications`.
2. Gin middleware attaches request ID, logs payload, enforces rate limit, and captures errors.
3. Handler validates JSON, maps to `services.SendNotificationRequest`.
4. Notification service persists record, then invokes `BroadcastService` to fan out toward stored destinations (current implementation calls stubbed HTTP clients and updates installation metadata).
5. Response returns `202 Accepted` with queued status; follow-up retrieval uses list/get endpoints.

## Key Design Notes
- **Layered modularity**: Handlers ⇄ Services ⇄ Repositories ensures swap-in for alternative transports (e.g., gRPC) without touching HTTP layer.
- **PostgreSQL-first**: Schema enforces enums and JSONB columns for targets/metadata; some foreign keys (e.g., `created_by`, `sender_id`) are tracked but not enforced to keep ingestion flexible.
- **Extensibility**: Queueing, advanced routing, and billing hooks are scaffolded in services and schema but require further implementation work.
- **Observability**: Structured JSON logs, request/response tracing, and health endpoint are built-in; Prometheus/metrics integration referenced in configs but not yet coded.
