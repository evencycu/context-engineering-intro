-- 010_add_priority_to_notification_destinations.sql
-- Add priority column to notification_destinations (low|normal|high)

ALTER TABLE notification_destinations
ADD COLUMN IF NOT EXISTS priority VARCHAR(20) NOT NULL DEFAULT 'normal'
    CHECK (priority IN ('low','normal','high'));

-- Optional index to help querying by priority
CREATE INDEX IF NOT EXISTS idx_notification_destinations_priority
    ON notification_destinations (priority);


