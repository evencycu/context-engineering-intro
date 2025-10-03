#!/usr/bin/env bash
set -euo pipefail

# Load .env if present (non-strict to allow running without it)
if [ -f .env ]; then
  set -a
  # shellcheck disable=SC1091
  source .env
  set +a
fi

TEAMS_APP_ID=${TEAMS_APP_ID:-""}
GOLDEN_CHAT_ID=${GOLDEN_CHAT_ID:-""}
GOLDEN_CHAT_ID_ENC=${GOLDEN_CHAT_ID_ENC:-""}

if [ -z "$TEAMS_APP_ID" ] || [ -z "$GOLDEN_CHAT_ID" ] || [ -z "$GOLDEN_CHAT_ID_ENC" ]; then
  echo "Missing required vars. Ensure .env has TEAMS_APP_ID, GOLDEN_CHAT_ID, GOLDEN_CHAT_ID_ENC" >&2
  exit 1
fi

echo "Using TEAMS_APP_ID=${TEAMS_APP_ID}"
echo "Using GOLDEN_CHAT_ID=${GOLDEN_CHAT_ID}"
echo "Using GOLDEN_CHAT_ID_ENC=${GOLDEN_CHAT_ID_ENC}"

if [ -z "${TOKEN:-}" ]; then
  echo "WARNING: TOKEN not set. Export a Microsoft Graph Bearer token in TOKEN to run curl calls." >&2
fi

cat <<'EOF'
# --- Create group chat (example) ---
# Requires: TOKEN env (Graph access token)
# curl -X POST https://graph.microsoft.com/v1.0/chats \
#   -H "Authorization: Bearer $TOKEN" \
#   -H "Content-Type: application/json" \
#   -d '{
#     "chatType": "group",
#     "topic": "客服協助群組A1",
#     "members": [
#       {"@odata.type":"#microsoft.graph.aadUserConversationMember","roles":["owner"],"user@odata.bind":"https://graph.microsoft.com/v1.0/users/bd5ab633-3f6f-48ad-be04-e144480da48f"},
#       {"@odata.type":"#microsoft.graph.aadUserConversationMember","roles":["owner"],"user@odata.bind":"https://graph.microsoft.com/v1.0/users/76bdf4ee-c942-4c75-8fdf-608909c51abb"}
#     ]
#   }'

# --- Install Teams app into the chat ---
# curl -X POST "https://graph.microsoft.com/v1.0/chats/${GOLDEN_CHAT_ID_ENC}/installedApps" \
#   -H "Authorization: Bearer $TOKEN" \
#   -H "Content-Type: application/json" \
#   -d "{\n  \"teamsApp@odata.bind\": \"https://graph.microsoft.com/v1.0/appCatalogs/teamsApps/${TEAMS_APP_ID}\",\n  \"consentedPermissionSet\": {\n    \"resourceSpecificPermissions\": [\n      {\"permissionValue\":\"ChatMessage.Read.Chat\",\"permissionType\":\"application\"},\n      {\"permissionValue\":\"ChannelMessage.Read.Group\",\"permissionType\":\"application\"}\n    ]\n  }\n}"

# --- Optional: verify installed apps ---
# curl -s -X GET "https://graph.microsoft.com/v1.0/chats/${GOLDEN_CHAT_ID_ENC}/installedApps?$expand=teamsAppDefinition" \
#   -H "Authorization: Bearer $TOKEN" | jq .
EOF

echo "Sample commands printed above. Copy and run after exporting TOKEN."


