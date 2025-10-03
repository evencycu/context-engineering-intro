-- Add three Teams Bot installation records to bot_installations table
-- Based on the provided JSON data from Teams installation events

-- Record 1: Personal conversation
INSERT INTO bot_installations (
    id,
    bot_id,
    conversation_id,
    conversation_type,
    tenant_id,
    service_url,
    status,
    created_at,
    updated_at
) VALUES (
    uuid_generate_v4(),
    '850e8400-e29b-41d4-a716-446655440001', -- Main Notification Bot ID from seed data
    'a:12mhoHc_sRnffmXHY2H5EvR6MyvmkXiLI5pQ54k3o04gnTMip5k5XPJfrVzA0f8j0mt27QzqCW-Dn5EmRXZa14ckeenzWBArx_V0biX160RcnYMeg5rRzJ6isYrYx-TZR',
    'personal',
    '051cece0-e4dc-4aed-b471-bf29824e1ee6',
    'https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/',
    'active',
    NOW(),
    NOW()
);

-- Record 2: Group chat conversation
INSERT INTO bot_installations (
    id,
    bot_id,
    conversation_id,
    conversation_type,
    tenant_id,
    service_url,
    status,
    created_at,
    updated_at
) VALUES (
    uuid_generate_v4(),
    '850e8400-e29b-41d4-a716-446655440001', -- Main Notification Bot ID from seed data
    '19:f26a8d8a235f430db87a404491cd2ffc@thread.v2',
    'groupChat',
    '051cece0-e4dc-4aed-b471-bf29824e1ee6',
    'https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/',
    'active',
    NOW(),
    NOW()
);

-- Record 3: Channel conversation
INSERT INTO bot_installations (
    id,
    bot_id,
    conversation_id,
    conversation_type,
    tenant_id,
    service_url,
    status,
    created_at,
    updated_at
) VALUES (
    uuid_generate_v4(),
    '850e8400-e29b-41d4-a716-446655440001', -- Main Notification Bot ID from seed data
    '19:lg5lz80dPDcE8OtOolOHKsNZYIZI0IslJnnGDBV2H5A1@thread.tacv2',
    'channel',
    '051cece0-e4dc-4aed-b471-bf29824e1ee6',
    'https://smba.trafficmanager.net/apac/051cece0-e4dc-4aed-b471-bf29824e1ee6/',
    'active',
    NOW(),
    NOW()
);

-- Verify the insertions
SELECT 
    id,
    bot_id,
    conversation_id,
    conversation_type,
    tenant_id,
    status,
    created_at
FROM bot_installations 
WHERE tenant_id = '051cece0-e4dc-4aed-b471-bf29824e1ee6'
ORDER BY created_at DESC
LIMIT 5;
