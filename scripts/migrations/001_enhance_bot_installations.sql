-- Migration: Enhance bot_installations table for broadcast functionality
-- Date: 2025-09-18
-- Description: Add new fields to support broadcast messaging to any person/channel/group

-- Add new columns to bot_installations table
ALTER TABLE bot_installations 
ADD COLUMN IF NOT EXISTS scope VARCHAR(20) DEFAULT 'personal' CHECK (scope IN ('personal', 'team', 'groupChat')),
ADD COLUMN IF NOT EXISTS conversation_id VARCHAR(500),
ADD COLUMN IF NOT EXISTS service_url VARCHAR(500),
ADD COLUMN IF NOT EXISTS teams_chat_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS recipient_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS from_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS from_aad_object_id VARCHAR(255),
ADD COLUMN IF NOT EXISTS last_activity_at TIMESTAMP WITH TIME ZONE;

-- Update installation_status constraint to include 'stale'
ALTER TABLE bot_installations 
DROP CONSTRAINT IF EXISTS bot_installations_installation_status_check;

ALTER TABLE bot_installations 
ADD CONSTRAINT bot_installations_installation_status_check 
CHECK (installation_status IN ('active', 'inactive', 'uninstalled', 'stale'));

-- Update unique constraint to use conversation_id instead of multiple fields
ALTER TABLE bot_installations 
DROP CONSTRAINT IF EXISTS bot_installations_bot_id_bot_type_teams_tenant_id_teams_team_id_teams_channel_id_teams_user_id_key;

ALTER TABLE bot_installations 
ADD CONSTRAINT bot_installations_bot_tenant_conversation_unique 
UNIQUE (bot_id, bot_type, teams_tenant_id, conversation_id);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_bot_installations_bot_tenant ON bot_installations(bot_id, bot_type, teams_tenant_id);
CREATE INDEX IF NOT EXISTS idx_bot_installations_scope_status ON bot_installations(scope, installation_status);
CREATE INDEX IF NOT EXISTS idx_bot_installations_conversation ON bot_installations(conversation_id);
CREATE INDEX IF NOT EXISTS idx_bot_installations_team_channel ON bot_installations(teams_team_id, teams_channel_id) WHERE teams_team_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_bot_installations_chat ON bot_installations(teams_chat_id) WHERE teams_chat_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_bot_installations_user ON bot_installations(teams_user_id) WHERE teams_user_id IS NOT NULL;
CREATE INDEX IF NOT EXISTS idx_bot_installations_activity ON bot_installations(last_activity_at) WHERE last_activity_at IS NOT NULL;

-- Update existing records to have required fields
-- Set conversation_id to teams_channel_id for existing records (temporary)
UPDATE bot_installations 
SET conversation_id = COALESCE(teams_channel_id, teams_user_id, 'legacy-' || id::text)
WHERE conversation_id IS NULL;

-- Set recipient_id to a default pattern for existing records
UPDATE bot_installations 
SET recipient_id = '28:legacy-bot'
WHERE recipient_id IS NULL;

-- Set from_id to teams_user_id for existing records
UPDATE bot_installations 
SET from_id = teams_user_id
WHERE from_id IS NULL AND teams_user_id IS NOT NULL;

-- Make conversation_id NOT NULL after populating existing records
ALTER TABLE bot_installations 
ALTER COLUMN conversation_id SET NOT NULL;

-- Make recipient_id NOT NULL after populating existing records
ALTER TABLE bot_installations 
ALTER COLUMN recipient_id SET NOT NULL;

-- Add comment explaining the new structure
COMMENT ON TABLE bot_installations IS 'Enhanced bot installations table supporting broadcast messaging to any person/channel/group';
COMMENT ON COLUMN bot_installations.scope IS 'Installation scope: personal, team, or groupChat';
COMMENT ON COLUMN bot_installations.conversation_id IS 'Teams conversation ID for sending messages';
COMMENT ON COLUMN bot_installations.service_url IS 'Bot Framework service URL for this installation';
COMMENT ON COLUMN bot_installations.recipient_id IS 'Bot member ID (usually 28:app_id)';
COMMENT ON COLUMN bot_installations.from_id IS 'Source member ID from Teams event';
COMMENT ON COLUMN bot_installations.from_aad_object_id IS 'User AAD Object ID for user identification';
COMMENT ON COLUMN bot_installations.last_activity_at IS 'Last successful message sent/received timestamp';