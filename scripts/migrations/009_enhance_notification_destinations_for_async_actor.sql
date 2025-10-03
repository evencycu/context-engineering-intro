-- Migration: Enhance notification_destinations for async actor-based sending
-- This migration adds fields needed for actor-based async notification sending with queue management

-- 1. Add new columns to notification_destinations
ALTER TABLE notification_destinations 
ADD COLUMN IF NOT EXISTS conversation_id VARCHAR(500),
ADD COLUMN IF NOT EXISTS next_retry_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS first_attempt_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS last_attempt_at TIMESTAMP WITH TIME ZONE,
ADD COLUMN IF NOT EXISTS failure_reason TEXT,
ADD COLUMN IF NOT EXISTS retry_after INT, -- Retry-After from 429 response (seconds)
ADD COLUMN IF NOT EXISTS actor_id VARCHAR(255); -- ID of the actor handling this notification

-- 2. Update max_retries default to 5 (was 3)
ALTER TABLE notification_destinations 
ALTER COLUMN max_retries SET DEFAULT 5;

-- 3. Add new indexes for actor-based processing
CREATE INDEX IF NOT EXISTS idx_notification_destinations_next_retry 
ON notification_destinations(next_retry_at) 
WHERE status IN ('pending', 'failed') AND retry_count < max_retries;

CREATE INDEX IF NOT EXISTS idx_notification_destinations_conversation 
ON notification_destinations(conversation_id) 
WHERE conversation_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_notification_destinations_actor 
ON notification_destinations(actor_id) 
WHERE actor_id IS NOT NULL;

CREATE INDEX IF NOT EXISTS idx_notification_destinations_retry_processing 
ON notification_destinations(status, next_retry_at, retry_count);

-- 4. Create a view for retry-ready notifications
CREATE OR REPLACE VIEW notification_destinations_retry_ready AS
SELECT 
    id,
    notification_id,
    destination_id,
    conversation_id,
    bot_id,
    bot_type,
    status,
    retry_count,
    max_retries,
    next_retry_at,
    failure_reason,
    retry_after,
    error_message,
    created_at,
    updated_at
FROM notification_destinations
WHERE status IN ('pending', 'failed')
  AND retry_count < max_retries
  AND (next_retry_at IS NULL OR next_retry_at <= NOW());

-- 5. Create a view for notification destination statistics
CREATE OR REPLACE VIEW notification_destination_stats AS
SELECT 
    status,
    COUNT(*) as count,
    COUNT(CASE WHEN retry_count >= max_retries THEN 1 END) as exhausted_count,
    COUNT(CASE WHEN next_retry_at IS NOT NULL AND next_retry_at > NOW() THEN 1 END) as scheduled_count,
    AVG(retry_count) as avg_retry_count,
    MAX(retry_count) as max_retry_count_seen,
    MIN(next_retry_at) as next_retry_due,
    MAX(last_attempt_at) as last_attempt_time
FROM notification_destinations
GROUP BY status;

-- 6. Add comment to explain the schema
COMMENT ON COLUMN notification_destinations.conversation_id IS 'Teams conversation ID for this specific target';
COMMENT ON COLUMN notification_destinations.next_retry_at IS 'Timestamp when this notification should be retried (NULL = immediate)';
COMMENT ON COLUMN notification_destinations.first_attempt_at IS 'Timestamp of the first send attempt';
COMMENT ON COLUMN notification_destinations.last_attempt_at IS 'Timestamp of the last send attempt';
COMMENT ON COLUMN notification_destinations.failure_reason IS 'Reason for failure: rate_limit, timeout, server_error, network_error, unauthorized';
COMMENT ON COLUMN notification_destinations.retry_after IS 'Retry-After value from 429 response (seconds)';
COMMENT ON COLUMN notification_destinations.actor_id IS 'ID of the actor currently handling this notification';

-- 7. Drop old failed_notifications table (if exists)
-- We'll use notification_destinations for queue management instead
DROP TABLE IF EXISTS failed_notifications CASCADE;

-- 8. Remove the old failed_notifications_queue_status view (if exists)
DROP VIEW IF EXISTS failed_notifications_queue_status CASCADE;

COMMENT ON TABLE notification_destinations IS 'Manages notification delivery to each destination with actor-based async sending and retry logic';

