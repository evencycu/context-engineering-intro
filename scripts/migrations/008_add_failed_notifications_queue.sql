-- Migration: Add failed notifications queue table
-- Purpose: Store failed notifications for retry with circuit breaker support
-- Date: 2025-10-01

-- Create failed_notifications table for retry queue
CREATE TABLE IF NOT EXISTS failed_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    notification_id UUID NOT NULL REFERENCES notifications(id) ON DELETE CASCADE,
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    target_id VARCHAR(500) NOT NULL, -- conversation_id or email
    message TEXT NOT NULL,
    reason VARCHAR(50) NOT NULL, -- rate_limit, timeout, unauthorized, server_error, network_error, unknown
    error_message TEXT,
    retry_count INTEGER NOT NULL DEFAULT 0,
    max_retries INTEGER NOT NULL DEFAULT 5,
    next_retry_at TIMESTAMP WITH TIME ZONE NOT NULL,
    retry_after INTEGER, -- From Retry-After header (seconds)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    last_attempt_at TIMESTAMP WITH TIME ZONE,
    
    -- Additional context
    installation_id UUID REFERENCES bot_installations(id) ON DELETE SET NULL,
    bot_id UUID REFERENCES teams_bots(id) ON DELETE SET NULL,
    metadata JSONB DEFAULT '{}'::jsonb,
    
    -- Indexes for efficient querying
    CONSTRAINT check_retry_count CHECK (retry_count <= max_retries)
);

-- Create indexes for efficient queue processing
CREATE INDEX IF NOT EXISTS idx_failed_notifications_next_retry 
    ON failed_notifications(next_retry_at) 
    WHERE retry_count < max_retries;

CREATE INDEX IF NOT EXISTS idx_failed_notifications_notification_id 
    ON failed_notifications(notification_id);

CREATE INDEX IF NOT EXISTS idx_failed_notifications_project_id 
    ON failed_notifications(project_id);

CREATE INDEX IF NOT EXISTS idx_failed_notifications_reason 
    ON failed_notifications(reason);

CREATE INDEX IF NOT EXISTS idx_failed_notifications_created_at 
    ON failed_notifications(created_at DESC);

-- Create function to automatically update updated_at timestamp
CREATE OR REPLACE FUNCTION update_failed_notifications_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for auto-updating updated_at
DROP TRIGGER IF EXISTS trigger_update_failed_notifications_updated_at ON failed_notifications;
CREATE TRIGGER trigger_update_failed_notifications_updated_at
    BEFORE UPDATE ON failed_notifications
    FOR EACH ROW
    EXECUTE FUNCTION update_failed_notifications_updated_at();

-- Create view for queue status monitoring
CREATE OR REPLACE VIEW failed_notifications_queue_status AS
SELECT 
    COUNT(*) FILTER (WHERE retry_count < max_retries AND next_retry_at > NOW()) as pending_count,
    COUNT(*) FILTER (WHERE retry_count < max_retries AND next_retry_at <= NOW()) as ready_for_retry_count,
    COUNT(*) FILTER (WHERE retry_count >= max_retries) as exhausted_count,
    MIN(next_retry_at) FILTER (WHERE retry_count < max_retries) as next_retry_due,
    MIN(created_at) FILTER (WHERE retry_count < max_retries) as oldest_pending,
    MAX(updated_at) as last_activity
FROM failed_notifications;

COMMENT ON TABLE failed_notifications IS 'Queue for storing failed notification attempts with retry mechanism';
COMMENT ON COLUMN failed_notifications.reason IS 'Failure reason: rate_limit, timeout, unauthorized, server_error, network_error, unknown';
COMMENT ON COLUMN failed_notifications.retry_after IS 'Seconds to wait before retry from Retry-After header';
COMMENT ON COLUMN failed_notifications.next_retry_at IS 'Timestamp when this notification should be retried';
COMMENT ON VIEW failed_notifications_queue_status IS 'Real-time status of the failed notifications queue';

