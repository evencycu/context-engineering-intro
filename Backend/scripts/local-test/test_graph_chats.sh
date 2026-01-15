#!/bin/bash
# Test Microsoft Graph API /chats endpoint
# This script tests if we can get all chats (personal and group chats) using application permissions

set -e

# Load environment variables
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

# Required environment variables
CLIENT_ID="${TEAMS_BOT_APP_ID:-${CLIENT_ID}}"
CLIENT_SECRET="${TEAMS_BOT_CLIENT_SECRET:-${TEAMS_BOT_APP_PASSWORD}}"
TENANT_ID="${TEAMS_TENANT_ID}"

if [ -z "$CLIENT_ID" ] || [ -z "$CLIENT_SECRET" ] || [ -z "$TENANT_ID" ]; then
    echo "Error: Missing required environment variables"
    echo "Required: TEAMS_BOT_APP_ID (or CLIENT_ID), TEAMS_BOT_CLIENT_SECRET (or TEAMS_BOT_APP_PASSWORD), TEAMS_TENANT_ID"
    exit 1
fi

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

echo "✅ Access token obtained"
echo ""

# Step 2: Test GET /chats (all chats)
echo "Step 2: Testing GET /chats (all chats)..."
echo "----------------------------------------"
CHATS_RESPONSE=$(curl -s -X GET "https://graph.microsoft.com/v1.0/chats" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json")

echo "$CHATS_RESPONSE" | jq '.' 2>/dev/null || echo "$CHATS_RESPONSE"

ERROR=$(echo "$CHATS_RESPONSE" | jq -r '.error.code // empty' 2>/dev/null)
if [ -n "$ERROR" ]; then
    echo ""
    echo "❌ Error: $ERROR"
    ERROR_MESSAGE=$(echo "$CHATS_RESPONSE" | jq -r '.error.message // empty' 2>/dev/null)
    if [ -n "$ERROR_MESSAGE" ]; then
        echo "   Message: $ERROR_MESSAGE"
    fi
else
    CHAT_COUNT=$(echo "$CHATS_RESPONSE" | jq '.value | length' 2>/dev/null || echo "0")
    echo ""
    echo "✅ Success! Found $CHAT_COUNT chats"
fi
echo ""

# Step 3: Test GET /chats?$filter=chatType eq 'oneOnOne' (personal chats)
echo "Step 3: Testing GET /chats?filter=chatType eq 'oneOnOne' (personal chats)..."
echo "----------------------------------------"
PERSONAL_CHATS_RESPONSE=$(curl -s -X GET "https://graph.microsoft.com/v1.0/chats?\$filter=chatType eq 'oneOnOne'" \
    -H "Authorization: Bearer ${ACCESS_TOKEN}" \
    -H "Content-Type: application/json")

echo "$PERSONAL_CHATS_RESPONSE" | jq '.' 2>/dev/null || echo "$PERSONAL_CHATS_RESPONSE"

ERROR=$(echo "$PERSONAL_CHATS_RESPONSE" | jq -r '.error.code // empty' 2>/dev/null)
if [ -n "$ERROR" ]; then
    echo ""
    echo "❌ Error: $ERROR"
    ERROR_MESSAGE=$(echo "$PERSONAL_CHATS_RESPONSE" | jq -r '.error.message // empty' 2>/dev/null)
    if [ -n "$ERROR_MESSAGE" ]; then
        echo "   Message: $ERROR_MESSAGE"
    fi
else
    CHAT_COUNT=$(echo "$PERSONAL_CHATS_RESPONSE" | jq '.value | length' 2>/dev/null || echo "0")
    echo ""
    echo "✅ Success! Found $CHAT_COUNT personal chats"
    
    # Show first chat as example
    if [ "$CHAT_COUNT" -gt 0 ]; then
        echo ""
        echo "First personal chat example:"
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

echo "$GROUP_CHATS_RESPONSE" | jq '.' 2>/dev/null || echo "$GROUP_CHATS_RESPONSE"

ERROR=$(echo "$GROUP_CHATS_RESPONSE" | jq -r '.error.code // empty' 2>/dev/null)
if [ -n "$ERROR" ]; then
    echo ""
    echo "❌ Error: $ERROR"
    ERROR_MESSAGE=$(echo "$GROUP_CHATS_RESPONSE" | jq -r '.error.message // empty' 2>/dev/null)
    if [ -n "$ERROR_MESSAGE" ]; then
        echo "   Message: $ERROR_MESSAGE"
    fi
else
    CHAT_COUNT=$(echo "$GROUP_CHATS_RESPONSE" | jq '.value | length' 2>/dev/null || echo "0")
    echo ""
    echo "✅ Success! Found $CHAT_COUNT group chats"
    
    # Show first chat as example
    if [ "$CHAT_COUNT" -gt 0 ]; then
        echo ""
        echo "First group chat example:"
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
echo "  - GET /chats: $(if [ -z "$(echo "$CHATS_RESPONSE" | jq -r '.error.code // empty' 2>/dev/null)"; then echo "✅ Success"; else echo "❌ Failed"; fi)"
echo "  - GET /chats (personal): $(if [ -z "$(echo "$PERSONAL_CHATS_RESPONSE" | jq -r '.error.code // empty' 2>/dev/null)"; then echo "✅ Success"; else echo "❌ Failed"; fi)"
echo "  - GET /chats (group): $(if [ -z "$(echo "$GROUP_CHATS_RESPONSE" | jq -r '.error.code // empty' 2>/dev/null)"; then echo "✅ Success"; else echo "❌ Failed"; fi)"
echo ""
echo "Note: If you see 'Insufficient privileges' error, you need to:"
echo "  1. Add 'Chat.Read.All' permission in Azure AD App Registration"
echo "  2. Get tenant admin consent for the permission"
echo "=========================================="
