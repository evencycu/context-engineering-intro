-- 添加批次通知相關表
-- 批次通知表
CREATE TABLE IF NOT EXISTS batch_notifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
    message_type VARCHAR(50) NOT NULL DEFAULT 'text',
    content TEXT NOT NULL,
    priority VARCHAR(20) DEFAULT 'normal', -- 'low', 'normal', 'high'
    mentions TEXT[] DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'processing', 'completed', 'failed', 'cancelled'
    total_targets INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    scheduled_at TIMESTAMP WITH TIME ZONE,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 批次目標表
CREATE TABLE IF NOT EXISTS batch_targets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    batch_id UUID REFERENCES batch_notifications(id) ON DELETE CASCADE,
    target_index INTEGER NOT NULL,
    destination_id UUID REFERENCES destinations(id) ON DELETE SET NULL,
    conversation_id VARCHAR(255),
    user_id VARCHAR(255),
    email VARCHAR(255),
    custom_data JSONB DEFAULT '{}',
    status VARCHAR(50) DEFAULT 'pending', -- 'pending', 'sent', 'failed', 'cancelled'
    error_message TEXT,
    sent_at TIMESTAMP WITH TIME ZONE,
    retry_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 批次模板表
CREATE TABLE IF NOT EXISTS batch_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT,
    project_id UUID REFERENCES projects(id) ON DELETE CASCADE,
    message_type VARCHAR(50) NOT NULL DEFAULT 'text',
    content TEXT NOT NULL,
    priority VARCHAR(20) DEFAULT 'normal',
    mentions TEXT[] DEFAULT '{}',
    metadata JSONB DEFAULT '{}',
    targets JSONB DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 創建索引
CREATE INDEX IF NOT EXISTS idx_batch_notifications_project_id ON batch_notifications(project_id);
CREATE INDEX IF NOT EXISTS idx_batch_notifications_sender_id ON batch_notifications(sender_id);
CREATE INDEX IF NOT EXISTS idx_batch_notifications_status ON batch_notifications(status);
CREATE INDEX IF NOT EXISTS idx_batch_notifications_scheduled_at ON batch_notifications(scheduled_at);
CREATE INDEX IF NOT EXISTS idx_batch_targets_batch_id ON batch_targets(batch_id);
CREATE INDEX IF NOT EXISTS idx_batch_targets_destination_id ON batch_targets(destination_id);
CREATE INDEX IF NOT EXISTS idx_batch_targets_status ON batch_targets(status);
CREATE INDEX IF NOT EXISTS idx_batch_templates_project_id ON batch_templates(project_id);
CREATE INDEX IF NOT EXISTS idx_batch_templates_is_active ON batch_templates(is_active);
