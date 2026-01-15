#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# 載入 TEAMS_* 環境變數（優先用 Backend/configs/.env，其次 .env）
if [ -f "$ROOT_DIR/Backend/configs/.env" ]; then
  # shellcheck disable=SC2046
  export $(grep -E '^(TEAMS_TENANT_ID|TEAMS_BOT_APP_ID|TEAMS_BOT_APP_PASSWORD)=' "$ROOT_DIR/Backend/configs/.env" | xargs)
elif [ -f "$ROOT_DIR/.env" ]; then
  # shellcheck disable=SC2046
  export $(grep -E '^(TEAMS_TENANT_ID|TEAMS_BOT_APP_ID|TEAMS_BOT_APP_PASSWORD)=' "$ROOT_DIR/.env" | xargs)
fi

if [[ -z "${TEAMS_TENANT_ID:-}" || -z "${TEAMS_BOT_APP_ID:-}" || -z "${TEAMS_BOT_APP_PASSWORD:-}" ]]; then
  echo "TEAMS_TENANT_ID / TEAMS_BOT_APP_ID / TEAMS_BOT_APP_PASSWORD 沒設定，先在 .env 或 Backend/configs/.env 補好再跑。"
  exit 1
fi

TENANT_ID="$TEAMS_TENANT_ID"
CLIENT_ID="$TEAMS_BOT_APP_ID"
CLIENT_SECRET="$TEAMS_BOT_APP_PASSWORD"

echo "=== 1. 取得 Graph Token ==="
TOKEN_RESPONSE=$(curl -s -X POST \
  "https://login.microsoftonline.com/${TENANT_ID}/oauth2/v2.0/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials" \
  -d "client_id=${CLIENT_ID}" \
  -d "client_secret=${CLIENT_SECRET}" \
  -d "scope=https%3A%2F%2Fgraph.microsoft.com%2F.default")

echo "$TOKEN_RESPONSE" | jq . >/dev/null 2>&1 || {
  echo "回傳不是 JSON："
  echo "$TOKEN_RESPONSE"
  exit 1
}

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token // empty')

if [ -z "$ACCESS_TOKEN" ] || [ "$ACCESS_TOKEN" = "null" ]; then
  echo "access_token 取得失敗："
  echo "$TOKEN_RESPONSE" | jq .
  exit 1
fi

echo "取得 token 成功"

echo
echo "=== 2. 打 /users 取樣 ==="
curl -s \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  "https://graph.microsoft.com/v1.0/users?\$select=id,displayName,mail,userPrincipalName,jobTitle,department&\$top=5" \
  | jq .

echo
echo "=== 3. 打 /groups 取樣 ==="
curl -s \
  -H "Authorization: Bearer $ACCESS_TOKEN" \
  "https://graph.microsoft.com/v1.0/groups?\$select=id,displayName,description,groupTypes&\$top=5" \
  | jq .

