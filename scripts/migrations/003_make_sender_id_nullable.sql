-- Make notifications.sender_id nullable to allow system/automation sends
ALTER TABLE IF EXISTS notifications
    ALTER COLUMN sender_id DROP NOT NULL;

-- Optional: add comment for documentation
COMMENT ON COLUMN notifications.sender_id IS 'Nullable: sender may be system or unknown';


