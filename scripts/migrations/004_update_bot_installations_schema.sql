-- Migration: Update bot_installations schema to match Teams Bot Framework events
-- Date: 2025-09-23

-- Add new columns
ALTER TABLE bot_installations 
ADD COLUMN conversation_type VARCHAR(20) CHECK (conversation_type IN ('personal', 'channel', 'groupChat')),
ADD COLUMN from_aad_object_id VARCHAR(255),
ADD COLUMN from_name VARCHAR(255),
ADD COLUMN recipient_name VARCHAR(255);

-- Update existing data to set conversation_type based on existing scope
UPDATE bot_installations 
SET conversation_type = CASE 
    WHEN scope = 'personal' THEN 'personal'
    WHEN scope = 'team' THEN 'channel' 
    WHEN scope = 'groupChat' THEN 'groupChat'
    ELSE 'personal'
END;

-- Make conversation_type NOT NULL after data migration
ALTER TABLE bot_installations ALTER COLUMN conversation_type SET NOT NULL;

-- Drop the old scope column
ALTER TABLE bot_installations DROP COLUMN scope;

-- Update unique constraint to use conversation_type instead of scope
ALTER TABLE bot_installations 
DROP CONSTRAINT IF EXISTS bot_installations_bot_id_bot_type_teams_tenant_id_teams_tea_key;

ALTER TABLE bot_installations 
ADD CONSTRAINT bot_installations_unique_installation 
UNIQUE (bot_id, bot_type, teams_tenant_id, conversation_type, teams_team_id, teams_channel_id, teams_user_id);

-- Add indexes for better performance
CREATE INDEX idx_bot_installations_conversation_type ON bot_installations(conversation_type);
CREATE INDEX idx_bot_installations_from_aad_object_id ON bot_installations(from_aad_object_id);
CREATE INDEX idx_bot_installations_conversation_id ON bot_installations(conversation_id);

-- Add comments
COMMENT ON COLUMN bot_installations.conversation_type IS 'Teams conversation type: personal, channel, or groupChat';
COMMENT ON COLUMN bot_installations.from_aad_object_id IS 'User AAD Object ID from Teams event';
COMMENT ON COLUMN bot_installations.from_name IS 'User display name from Teams event';
COMMENT ON COLUMN bot_installations.recipient_name IS 'Bot display name from Teams event';
COMMENT ON COLUMN bot_installations.conversation_id IS 'Teams conversation ID for sending messages';
COMMENT ON COLUMN bot_installations.service_url IS 'Bot Framework service URL for sending messages';
