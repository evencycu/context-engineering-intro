# Entity Relationship Snapshot

The schema mirrors `internal/database/schema.sql`. Dashed relationships indicate logical links without an explicit database foreign key constraint.

```
Companies (companies)
- id UUID PK
- name, contact_email, contact_phone, address
- status ENUM(active|inactive|suspended)
- billing_enabled BOOLEAN
- timestamps

Users (users)
- id UUID PK
- company_id UUID FK → companies(id)
- email UNIQUE per company
- name, role ENUM(admin|user|viewer), status ENUM(...)
- password_hash, api_key_hash, api_key_expires_at
- last_login_at, timestamps

Projects (projects)
- id UUID PK
- company_id UUID FK → companies(id)
- notify_key UNIQUE per company
- description, status ENUM(...), daily_limit, monthly_limit, priority ENUM(...)
- created_by UUID (logical link → users.id, constraint NOT enforced)
- timestamps

Teams Bots (teams_bots)
- id UUID PK
- type ENUM(platform|third_party)
- company_id UUID (only for third_party, optional, no FK)
- app_id UNIQUE, app_password_hash, tenant_id
- status ENUM(active|inactive|suspended|maintenance)
- webhook_url, capabilities JSONB, rate_limit fields
- api_endpoint/api_key_hash/contact info (third-party only)
- created_by UUID (no FK)
- timestamps

Bot Installations (bot_installations)
- id UUID PK
- bot_id UUID FK → teams_bots(id)
- bot_type ENUM(platform|third_party)
- teams_tenant_id
- conversation_type ENUM(personal|channel|groupChat)
- conversation_id, service_url
- recipient/from identifiers, email, description_name
- installation_status ENUM(active|inactive|uninstalled|stale)
- installed_at, uninstalled_at, last_activity_at
- metadata JSONB, timestamps
- UNIQUE(bot_id, bot_type, teams_tenant_id, conversation_id)

Destinations (destinations)
- id UUID PK
- project_id UUID FK → projects(id)
- name, description, teams_tenant_id
- targets JSONB (array of TeamsTarget, validated in service layer)
- bot_id UUID (optional, no FK)
- status ENUM(active|inactive|suspended)
- validation_status ENUM(pending|validated|failed)
- last_validated_at
- created_by UUID (logical link → users.id, no FK)
- timestamps

Notifications (notifications)
- id UUID PK
- project_id UUID FK → projects(id)
- sender_id UUID (logical link → users.id, no FK)
- message_type ENUM(text|file|adaptive_card)
- content, mentions JSONB, attachment JSONB, adaptive_card JSONB
- priority ENUM(low|normal|high|urgent)
- status ENUM(pending|processing|sent|failed|cancelled)
- error_message, metadata JSONB
- sent_at, timestamps

Notification Destinations (notification_destinations)
- id UUID PK
- notification_id UUID FK → notifications(id)
- destination_id UUID FK → destinations(id)
- bot_id UUID (optional, no FK)
- bot_type ENUM(platform|third_party)
- status ENUM(pending|processing|sent|failed|cancelled)
- error_message, teams_message_id
- retry_count, max_retries, sent_at, timestamps

Bot Routing Rules (bot_routing_rules)
- id UUID PK
- name, description
- priority INTEGER
- conditions JSONB
- target_bot_id UUID (logical link → teams_bots.id, no FK)
- target_bot_type ENUM(platform|third_party)
- created_by UUID (logical link → users.id, no FK)
- status ENUM(active|inactive|archived)
- timestamps

Queue & Scheduling Tables (queues, scheduled_notifications, notification_retry_queue)
- Present in schema for future queue/worker implementation; currently unused by application code.
```
