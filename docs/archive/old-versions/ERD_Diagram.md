# Teams Notification Bot Platform - Entity Relationship Diagram

## Core Platform Entities

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                COMPANIES                                        │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ name (VARCHAR(255), NOT NULL)                                                  │
│ contact_email (VARCHAR(255), NOT NULL)                                         │
│ contact_phone (VARCHAR(50))                                                    │
│ address (TEXT)                                                                  │
│ status (VARCHAR(20), CHECK: active|inactive|suspended)                         │
│ billing_enabled (BOOLEAN, DEFAULT true)                                        │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                  USERS                                          │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ company_id (UUID, FK → companies.id)                                           │
│ email (VARCHAR(255), NOT NULL)                                                 │
│ name (VARCHAR(255), NOT NULL)                                                  │
│ role (VARCHAR(50), CHECK: admin|user|viewer)                                   │
│ status (VARCHAR(20), CHECK: active|inactive|suspended)                         │
│ last_login_at (TIMESTAMP)                                                      │
│ password_hash (VARCHAR(255))                                                   │
│ api_key_hash (VARCHAR(255))                                                    │
│ api_key_expires_at (TIMESTAMP)                                                 │
│ created_at, updated_at (TIMESTAMP)                                             │
│ UNIQUE(company_id, email)                                                      │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                                PROJECTS                                         │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ company_id (UUID, FK → companies.id)                                           │
│ notify_key (VARCHAR(100), NOT NULL)                                            │
│ description (TEXT)                                                              │
│ status (VARCHAR(20), CHECK: active|inactive|suspended)                         │
│ daily_limit (INTEGER, DEFAULT 10000)                                           │
│ monthly_limit (INTEGER, DEFAULT 300000)                                        │
│ priority (VARCHAR(20), CHECK: low|normal|high|urgent)                          │
│ created_by (UUID, FK → users.id)                                               │
│ created_at, updated_at (TIMESTAMP)                                             │
│ UNIQUE(company_id, notify_key)                                                 │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Bot Management Entities

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            PLATFORM_BOTS                                        │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ name (VARCHAR(255), NOT NULL)                                                  │
│ description (TEXT)                                                              │
│ app_id (VARCHAR(255), NOT NULL, UNIQUE)                                        │
│ app_password_hash (VARCHAR(255), NOT NULL)                                     │
│ tenant_id (VARCHAR(255))                                                       │
│ status (bot_status, DEFAULT 'active')                                          │
│ webhook_url (VARCHAR(500))                                                     │
│ capabilities (JSONB, DEFAULT '{}')                                             │
│ rate_limit_per_minute (INTEGER, DEFAULT 600)                                   │
│ max_concurrent_requests (INTEGER, DEFAULT 100)                                 │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            BOT_INSTALLATIONS                                    │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ bot_id (UUID, NOT NULL)                                                        │
│ bot_type (bot_type, NOT NULL)                                                  │
│ teams_tenant_id (VARCHAR(255), NOT NULL)                                       │
│ scope (VARCHAR(20), CHECK: personal|team|groupChat)                            │
│ conversation_id (VARCHAR(500), NOT NULL)                                       │
│ service_url (VARCHAR(500), NOT NULL)                                           │
│ teams_team_id (VARCHAR(255))                                                   │
│ teams_channel_id (VARCHAR(255))                                                │
│ teams_chat_id (VARCHAR(255))                                                   │
│ teams_user_id (VARCHAR(255))                                                   │
│ recipient_id (VARCHAR(255), NOT NULL)                                          │
│ from_id (VARCHAR(255))                                                         │
│ from_aad_object_id (VARCHAR(255))                                              │
│ installation_status (VARCHAR(20), CHECK: active|inactive|uninstalled|stale)   │
│ installed_at (TIMESTAMP)                                                       │
│ uninstalled_at (TIMESTAMP)                                                     │
│ last_activity_at (TIMESTAMP)                                                   │
│ metadata (JSONB, DEFAULT '{}')                                                 │
│ created_at, updated_at (TIMESTAMP)                                             │
│ UNIQUE(bot_id, bot_type, teams_tenant_id, conversation_id)                     │
└─────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────┐
│                          THIRD_PARTY_BOTS                                       │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ company_id (UUID, FK → companies.id)                                           │
│ name (VARCHAR(255), NOT NULL)                                                  │
│ description (TEXT)                                                              │
│ app_id (VARCHAR(255), NOT NULL)                                                │
│ app_password_hash (VARCHAR(255), NOT NULL)                                     │
│ tenant_id (VARCHAR(255))                                                       │
│ status (bot_status, DEFAULT 'active')                                          │
│ webhook_url (VARCHAR(500))                                                     │
│ api_endpoint (VARCHAR(500))                                                    │
│ api_key_hash (VARCHAR(255))                                                    │
│ capabilities (JSONB, DEFAULT '{}')                                             │
│ rate_limit_per_minute (INTEGER, DEFAULT 100)                                   │
│ max_concurrent_requests (INTEGER, DEFAULT 10)                                  │
│ contact_email (VARCHAR(255))                                                   │
│ contact_phone (VARCHAR(50))                                                    │
│ created_by (UUID, FK → users.id)                                               │
│ created_at, updated_at (TIMESTAMP)                                             │
│ UNIQUE(company_id, app_id)                                                     │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            BOT_INSTALLATIONS                                    │
│                              (Same as above)                                   │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Destination & Notification Entities

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              DESTINATIONS                                      │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ project_id (UUID, FK → projects.id)                                            │
│ name (VARCHAR(255), NOT NULL)                                                  │
│ description (TEXT)                                                              │
│ teams_tenant_id (VARCHAR(255), NOT NULL)                                       │
│ targets (JSONB, NOT NULL, DEFAULT '[]')                                        │
│ bot_id (UUID)                                                                   │
│ bot_type (bot_type)                                                             │
│ status (VARCHAR(20), CHECK: active|inactive|suspended)                         │
│ validation_status (VARCHAR(20), CHECK: pending|validated|failed)               │
│ last_validated_at (TIMESTAMP)                                                  │
│ created_by (UUID, FK → users.id)                                               │
│ created_at, updated_at (TIMESTAMP)                                             │
│ CONSTRAINT: jsonb_array_length(targets) > 0                                    │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              NOTIFICATIONS                                     │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ project_id (UUID, FK → projects.id)                                            │
│ sender_id (UUID, FK → users.id)                                                │
│ message_type (VARCHAR(20), CHECK: text|file|adaptive_card)                     │
│ content (TEXT, NOT NULL)                                                       │
│ mentions (JSONB, DEFAULT '[]')                                                 │
│ attachment (JSONB)                                                              │
│ adaptive_card (JSONB)                                                           │
│ priority (VARCHAR(20), CHECK: low|normal|high|urgent)                          │
│ status (VARCHAR(20), CHECK: pending|processing|sent|failed|cancelled)          │
│ error_message (TEXT)                                                            │
│ metadata (JSONB, DEFAULT '{}')                                                 │
│ created_at, updated_at (TIMESTAMP)                                             │
│ sent_at (TIMESTAMP)                                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        NOTIFICATION_DESTINATIONS                               │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ notification_id (UUID, FK → notifications.id)                                  │
│ destination_id (UUID, FK → destinations.id)                                    │
│ bot_id (UUID)                                                                   │
│ bot_type (bot_type)                                                             │
│ status (VARCHAR(20), CHECK: pending|processing|sent|failed|cancelled)          │
│ error_message (TEXT)                                                            │
│ teams_message_id (VARCHAR(255))                                                │
│ sent_at (TIMESTAMP)                                                             │
│ retry_count (INTEGER, DEFAULT 0)                                               │
│ max_retries (INTEGER, DEFAULT 3)                                               │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Bot Management & Health Entities

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            BOT_ROUTING_RULES                                   │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ name (VARCHAR(255), NOT NULL)                                                  │
│ description (TEXT)                                                              │
│ priority (INTEGER, DEFAULT 100)                                                │
│ conditions (JSONB, NOT NULL)                                                   │
│ target_bot_id (UUID)                                                            │
│ target_bot_type (bot_type)                                                      │
│ is_active (BOOLEAN, DEFAULT true)                                              │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────┐
│                            BOT_HEALTH_STATUS                                   │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ bot_id (UUID, NOT NULL)                                                        │
│ bot_type (bot_type, NOT NULL)                                                  │
│ status (VARCHAR(20), CHECK: healthy|unhealthy|degraded)                        │
│ last_check_at (TIMESTAMP)                                                      │
│ response_time_ms (INTEGER)                                                     │
│ error_count (INTEGER, DEFAULT 0)                                               │
│ success_count (INTEGER, DEFAULT 0)                                             │
│ error_message (TEXT)                                                            │
│ metadata (JSONB, DEFAULT '{}')                                                 │
│ created_at (TIMESTAMP)                                                          │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Third Party Bot API Management

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        THIRD_PARTY_BOT_API_KEYS                                │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ bot_id (UUID, FK → third_party_bots.id)                                        │
│ notify_key (VARCHAR(255), NOT NULL)                                            │
│ api_key_hash (VARCHAR(255), NOT NULL)                                          │
│ permissions (JSONB, DEFAULT '{}')                                              │
│ rate_limit_per_minute (INTEGER, DEFAULT 100)                                   │
│ expires_at (TIMESTAMP)                                                          │
│ is_active (BOOLEAN, DEFAULT true)                                              │
│ created_by (UUID, FK → users.id)                                               │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                        THIRD_PARTY_BOT_API_USAGE                                │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ bot_id (UUID, FK → third_party_bots.id)                                        │
│ api_key_id (UUID, FK → third_party_bot_api_keys.id)                            │
│ endpoint (VARCHAR(255), NOT NULL)                                              │
│ method (VARCHAR(10), NOT NULL)                                                 │
│ status_code (INTEGER)                                                           │
│ response_time_ms (INTEGER)                                                     │
│ request_size_bytes (INTEGER)                                                   │
│ response_size_bytes (INTEGER)                                                  │
│ ip_address (INET)                                                               │
│ user_agent (TEXT)                                                               │
│ created_at (TIMESTAMP)                                                          │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Billing & Usage Entities

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              BILLING_PLANS                                     │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ name (VARCHAR(100), NOT NULL)                                                  │
│ description (TEXT)                                                              │
│ price_per_notification (DECIMAL(10,6), DEFAULT 0.001)                          │
│ price_per_attachment (DECIMAL(10,6), DEFAULT 0.005)                            │
│ price_per_mention (DECIMAL(10,6), DEFAULT 0.0005)                              │
│ price_per_adaptive_card (DECIMAL(10,6), DEFAULT 0.002)                         │
│ daily_limit (INTEGER)                                                           │
│ monthly_limit (INTEGER)                                                         │
│ status (VARCHAR(20), CHECK: active|inactive)                                   │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            COMPANY_BILLING                                      │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ company_id (UUID, FK → companies.id)                                           │
│ billing_plan_id (UUID, FK → billing_plans.id)                                  │
│ status (VARCHAR(20), CHECK: active|suspended|cancelled)                         │
│ billing_email (VARCHAR(255))                                                   │
│ payment_method (VARCHAR(50))                                                   │
│ currency (VARCHAR(3), DEFAULT 'USD')                                           │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        │ 1:N
                                        ▼
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              USAGE_RECORDS                                     │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ company_id (UUID, FK → companies.id)                                           │
│ project_id (UUID, FK → projects.id)                                            │
│ user_id (UUID, FK → users.id)                                                  │
│ notification_id (UUID, FK → notifications.id)                                  │
│ bot_id (UUID)                                                                   │
│ bot_type (bot_type)                                                             │
│ record_type (VARCHAR(20), CHECK: notification|attachment|mention|adaptive_card)│
│ quantity (INTEGER, DEFAULT 1)                                                  │
│ unit_price (DECIMAL(10,6), NOT NULL)                                           │
│ total_cost (DECIMAL(10,6), NOT NULL)                                           │
│ billing_period (DATE, NOT NULL)                                                │
│ created_at (TIMESTAMP)                                                          │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Audit & Logging Entities

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                              AUDIT_LOGS                                        │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ company_id (UUID, FK → companies.id)                                           │
│ user_id (UUID, FK → users.id)                                                  │
│ action (VARCHAR(100), NOT NULL)                                                │
│ resource_type (VARCHAR(50), NOT NULL)                                          │
│ resource_id (UUID)                                                              │
│ old_values (JSONB)                                                              │
│ new_values (JSONB)                                                              │
│ ip_address (INET)                                                               │
│ user_agent (TEXT)                                                               │
│ created_at (TIMESTAMP)                                                          │
└─────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────┐
│                              SYSTEM_LOGS                                       │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ level (VARCHAR(20), CHECK: debug|info|warn|error|fatal)                        │
│ service (VARCHAR(50), NOT NULL)                                                │
│ message (TEXT, NOT NULL)                                                       │
│ context (JSONB, DEFAULT '{}')                                                  │
│ trace_id (VARCHAR(255))                                                        │
│ span_id (VARCHAR(255))                                                         │
│ bot_id (UUID)                                                                   │
│ bot_type (bot_type)                                                             │
│ created_at (TIMESTAMP)                                                          │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Configuration Entities

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│                            SYSTEM_SETTINGS                                     │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ key (VARCHAR(100), UNIQUE, NOT NULL)                                           │
│ value (JSONB, NOT NULL)                                                        │
│ description (TEXT)                                                              │
│ is_public (BOOLEAN, DEFAULT false)                                             │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────────────────────┐
│                             FEATURE_FLAGS                                      │
├─────────────────────────────────────────────────────────────────────────────────┤
│ id (UUID, PK)                                                                   │
│ name (VARCHAR(100), UNIQUE, NOT NULL)                                          │
│ description (TEXT)                                                              │
│ enabled (BOOLEAN, DEFAULT false)                                               │
│ company_id (UUID, FK → companies.id)                                           │
│ created_at, updated_at (TIMESTAMP)                                             │
└─────────────────────────────────────────────────────────────────────────────────┘
```

## Key Relationships Summary

### Primary Relationships:
1. **Companies** → **Users** (1:N)
2. **Companies** → **Projects** (1:N)
3. **Users** → **Projects** (1:N, created_by)
4. **Projects** → **Destinations** (1:N)
5. **Projects** → **Notifications** (1:N)
6. **Users** → **Notifications** (1:N, sender_id)
7. **Notifications** → **Notification_Destinations** (1:N)
8. **Destinations** → **Notification_Destinations** (1:N)

### Bot Relationships:
1. **Platform_Bots** → **Bot_Installations** (1:N)
2. **Third_Party_Bots** → **Bot_Installations** (1:N)
3. **Third_Party_Bots** → **Third_Party_Bot_API_Keys** (1:N)
4. **Third_Party_Bot_API_Keys** → **Third_Party_Bot_API_Usage** (1:N)

### Billing Relationships:
1. **Companies** → **Company_Billing** (1:N)
2. **Billing_Plans** → **Company_Billing** (1:N)
3. **Companies** → **Usage_Records** (1:N)
4. **Projects** → **Usage_Records** (1:N)
5. **Users** → **Usage_Records** (1:N)
6. **Notifications** → **Usage_Records** (1:N)

### Audit Relationships:
1. **Companies** → **Audit_Logs** (1:N)
2. **Users** → **Audit_Logs** (1:N)
3. **Companies** → **Feature_Flags** (1:N)

## Enhanced Features for Broadcast:

### Bot Installations Table:
- **Scope-based targeting**: personal, team, groupChat
- **Conversation tracking**: Direct conversation ID for messaging
- **Service URL**: Bot Framework endpoint per installation
- **Activity tracking**: Last successful message timestamp
- **Stale detection**: Mark inactive installations as stale

### Destinations Table:
- **Enhanced targets JSONB**: Support for AAD Object ID, conversation ID, tenant ID
- **Multi-tenant support**: Tenant-specific service URLs
- **Flexible targeting**: Person, channel, group chat combinations

This schema supports the complete Teams notification platform with broadcast capabilities, multi-tenant support, comprehensive audit trails, and billing management.