#!/bin/bash
# Generate and execute curl commands for testing Graph API /chats endpoint
# This script will prompt for credentials if not provided

set -e

# Function to prompt for input if not provided
prompt_if_empty() {
    local var_name=$1
    local prompt_text=$2
    local is_secret=${3:-false}
    
    if [ -z "${!var_name}" ]; then
        if [ "$is_secret" = "true" ]; then
            read -sp "$prompt_text: " value
            echo ""
        else
            read -p "$prompt_text: " value
        fi
        eval "$var_name=\"$value\""
    fi
}

# Get credentials
CLIENT_ID="${1:-${TEAMS_BOT_APP_ID}}"
CLIENT_SECRET="${2:-${TEAMS_BOT_CLIENT_SECRET:-${TEAMS_BOT_APP_PASSWORD}}}"
TENANT_ID="${3:-${TEAMS_TENANT_ID}}"

# Prompt if missing
prompt_if_empty CLIENT_ID "Enter Client ID (TEAMS_BOT_APP_ID)"
prompt_if_empty CLIENT_SECRET "Enter Client Secret (TEAMS_BOT_APP_PASSWORD)" true
prompt_if_empty TENANT_ID "Enter Tenant ID (TEAMS_TENANT_ID)"

echo ""
echo "=========================================="
echo "Microsoft Graph API - Chats Test"
echo "=========================================="
echo "Client ID: $CLIENT_ID"
echo "Tenant ID: $TENANT_ID"
echo "=========================================="
echo ""

# Step 1: Get access token
echo "Step 1: Getting access token..."
TOKEN_URL="https://login.microsoftonline.com/${TENANT_ID}/oauth2/v2.0/token"

TOKEN_RESPONSE=$(curl -s -X POST "$TOKEN_URL" \
    -H "Content-Type: application/x-www-form-urlencoded" \
    -d "grant_type=client_credentials" \
    -d "client_id=${CLIENT_ID}" \
    -d "client_secret=${CLIENT_SECRET}" \
    -d "scope=https://graph.microsoft.com/.default")

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token // empty')

if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
    echo "❌ Failed to get access token"
    echo "Response: $TOKEN_RESPONSE"
    exit 1
fi

echo "✅ Access token obtained: ${ACCESS_TOKEN:0:50}..."
echo ""

# Step 2: Test GET /chats (all chats)
echo "Step 2: Testing GET /chats (all chats)..."
echo "----------------------------------------"
CHATS_RESPONSE=$(curl -s -X GET "https://graph.microsoft.com/v1.0/chats" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json")

ERROR=$(echo "$CHATS_RESPONSE" | jq -r '.error.code // empty' 2>/dev/null)
if [ -n "$ERROR" ]; then
    echo "❌ Error: $ERROR"
    ERROR_MESSAGE=$(echo "$CHATS_RESPONSE" | jq -r '.error.message // empty' 2>/dev/null)
    if [ -n "$ERROR_MESSAGE" ]; then
        echo "   Message: $ERROR_MESSAGE"
    fi
    echo ""
    echo "Full response:"
    echo "$CHATS_RESPONSE" | jq '.' 2>/dev/null || echo "$CHATS_RESPONSE"
else
    CHAT_COUNT=$(echo "$CHATS_RESPONSE" | jq '.value | length' 2>/dev/null || echo "0")
    echo "✅ Success! Found $CHAT_COUNT chats"
    if [ "$CHAT_COUNT" -gt 0 ]; then
        echo ""
        echo "First 3 chats:"
        echo "$CHATS_RESPONSE" | jq '.value[0:3] | .[] | {id, chatType, topic, createdDateTime}' 2>/dev/null
    fi
fi
echo ""

# Step 3: Test GET /chats?$filter=chatType eq 'oneOnOne' (personal chats)
echo "Step 3: Testing GET /chats?filter=chatType eq 'oneOnOne' (personal chats)..."
echo "----------------------------------------"
PERSONAL_CHATS_RESPONSE=$(curl -s -X GET "https://graph.microsoft.com/v1.0/chats?\$filter=chatType eq 'oneOnOne'" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json")

ERROR=$(echo "$PERSONAL_CHATS_RESPONSE" | jq -r '.error.code // empty' 2>/dev/null)
if [ -n "$ERROR" ]; then
    echo "❌ Error: $ERROR"
    ERROR_MESSAGE=$(echo "$PERSONAL_CHATS_RESPONSE" | jq -r '.error.message // empty' 2>/dev/null)
    if [ -n "$ERROR_MESSAGE" ]; then
        echo "   Message: $ERROR_MESSAGE"
    fi
else
    CHAT_COUNT=$(echo "$PERSONAL_CHATS_RESPONSE" | jq '.value | length' 2>/dev/null || echo "0")
    echo "✅ Success! Found $CHAT_COUNT personal chats"
    if [ "$CHAT_COUNT" -gt 0 ]; then
        echo ""
        echo "First personal chat:"
        echo "$PERSONAL_CHATS_RESPONSE" | jq '.value[0] | {id, chatType, topic, createdDateTime}' 2>/dev/null
    fi
fi
echo ""

# Step 4: Test GET /chats?$filter=chatType eq 'groupChat' (group chats)
echo "Step 4: Testing GET /chats?filter=chatType eq 'groupChat' (group chats)..."
echo "----------------------------------------"
GROUP_CHATS_RESPONSE=$(curl -s -X GET "https://graph.microsoft.com/v1.0/chats?\$filter=chatType eq 'groupChat'" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json")

ERROR=$(echo "$GROUP_CHATS_RESPONSE" | jq -r '.error.code // empty' 2>/dev/null)
if [ -n "$ERROR" ]; then
    echo "❌ Error: $ERROR"
    ERROR_MESSAGE=$(echo "$GROUP_CHATS_RESPONSE" | jq -r '.error.message // empty' 2>/dev/null)
    if [ -n "$ERROR_MESSAGE" ]; then
        echo "   Message: $ERROR_MESSAGE"
    fi
else
    CHAT_COUNT=$(echo "$GROUP_CHATS_RESPONSE" | jq '.value | length' 2>/dev/null || echo "0")
    echo "✅ Success! Found $CHAT_COUNT group chats"
    if [ "$CHAT_COUNT" -gt 0 ]; then
        echo ""
        echo "First group chat:"
        echo "$GROUP_CHATS_RESPONSE" | jq '.value[0] | {id, chatType, topic, createdDateTime}' 2>/dev/null
    fi
fi
echo ""

# Step 5: Summary
echo "=========================================="
echo "Test Summary"
echo "=========================================="
echo "✅ Access token: Obtained"
echo ""
echo "Test Results:"
ALL_SUCCESS=$(if [ -z "$(echo "$CHATS_RESPONSE" | jq -r '.error.code // empty' 2>/dev/null)"; then echo "✅ Success"; else echo "❌ Failed"; fi)
PERSONAL_SUCCESS=$(if [ -z "$(echo "$PERSONAL_CHATS_RESPONSE" | jq -r '.error.code // empty' 2>/dev/null)"; then echo "✅ Success"; else echo "❌ Failed"; fi)
GROUP_SUCCESS=$(if [ -z "$(echo "$GROUP_CHATS_RESPONSE" | jq -r '.error.code // empty' 2>/dev/null)"; then echo "✅ Success"; else echo "❌ Failed"; fi)

echo "  - GET /chats: $ALL_SUCCESS"
echo "  - GET /chats (personal): $PERSONAL_SUCCESS"
echo "  - GET /chats (group): $GROUP_SUCCESS"
echo ""
echo "Note: If you see 'Insufficient privileges' error, you need to:"
echo "  1. Add 'Chat.Read.All' permission in Azure AD App Registration"
echo "  2. Get tenant admin consent for the permission"
echo "=========================================="
