-- Create unified teams_bots table
CREATE TABLE IF NOT EXISTS teams_bots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,

    -- union discriminator
    type VARCHAR(20) NOT NULL CHECK (type IN ('platform','third_party')),

    -- common/union fields
    company_id UUID, -- only for third_party
    name VARCHAR(255) NOT NULL,
    description TEXT,
    app_id VARCHAR(255) NOT NULL UNIQUE,
    app_password_hash VARCHAR(255) NOT NULL,
    tenant_id VARCHAR(255),
    status VARCHAR(32) NOT NULL DEFAULT 'active',
    webhook_url VARCHAR(500),
    capabilities JSONB DEFAULT '{}'::jsonb,
    rate_limit_per_minute INTEGER DEFAULT 60,
    max_concurrent_requests INTEGER DEFAULT 10,

    -- third-party specific
    api_endpoint VARCHAR(500),
    api_key_hash VARCHAR(255),
    contact_email VARCHAR(255),
    contact_phone VARCHAR(50),
    created_by UUID
);

-- Copy platform_bots into teams_bots (type = platform)
INSERT INTO teams_bots (
    id, created_at, updated_at, type, company_id, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, capabilities, rate_limit_per_minute, max_concurrent_requests, api_endpoint, api_key_hash, contact_email, contact_phone, created_by
)
SELECT id, created_at, updated_at, 'platform', NULL, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, capabilities, rate_limit_per_minute, max_concurrent_requests, NULL, NULL, NULL, NULL, NULL
FROM platform_bots
ON CONFLICT (id) DO NOTHING;

-- Optionally copy third_party_bots if exists (type = third_party)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'third_party_bots') THEN
        INSERT INTO teams_bots (
            id, created_at, updated_at, type, company_id, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, capabilities, rate_limit_per_minute, max_concurrent_requests, api_endpoint, api_key_hash, contact_email, contact_phone, created_by
        )
        SELECT id, created_at, updated_at, 'third_party', company_id, name, description, app_id, app_password_hash, tenant_id, status, webhook_url, capabilities, rate_limit_per_minute, max_concurrent_requests, api_endpoint, api_key_hash, contact_email, contact_phone, created_by
        FROM third_party_bots
        ON CONFLICT (id) DO NOTHING;
    END IF;
END $$;

-- Clear old tables as requested
TRUNCATE TABLE platform_bots RESTART IDENTITY CASCADE;
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'third_party_bots') THEN
        TRUNCATE TABLE third_party_bots RESTART IDENTITY CASCADE;
    END IF;
END $$;


