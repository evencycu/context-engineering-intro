-- Migration: Switch bot_installations unique constraint to conversation_id based
-- Date: 2025-09-23

ALTER TABLE bot_installations
DROP CONSTRAINT IF EXISTS bot_installations_unique_installation;

ALTER TABLE bot_installations
ADD CONSTRAINT bot_installations_bot_tenant_conversation_unique
UNIQUE (bot_id, bot_type, teams_tenant_id, conversation_id);


