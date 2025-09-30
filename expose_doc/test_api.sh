#!/bin/bash

# Teams Notification API smoke test
# Optional env vars:
#   PROJECT_ID, CREATED_BY, TENANT_ID for destination creation
#   SENDER_ID (optional for notifications)

set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"
API_BASE="$BASE_URL/api/v1"
PROJECT_ID="${PROJECT_ID:-}"
CREATED_BY="${CREATED_BY:-}"
TENANT_ID="${TENANT_ID:-}"
SENDER_ID="${SENDER_ID:-}"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

require_tool() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo -e "${RED}缺少工具: $1${NC}" >&2
    exit 1
  fi
}

curl_api() {
  local method="$1"
  local url="$2"
  local data="${3:-}"
  local args=(-sS)
  if [[ "$method" == "GET" ]]; then
    curl "${args[@]}" "$url"
  else
    args+=(-H "Content-Type: application/json" -X "$method" -d "$data")
    curl "${args[@]}" "$url"
  fi
}

print_step() {
  echo -e "${BLUE}==> $1${NC}"
}

print_skip() {
  echo -e "${YELLOW}⚠️  跳過: $1${NC}"
}

print_fail() {
  echo -e "${RED}❌ $1${NC}"
}

print_ok() {
  echo -e "${GREEN}✅ $1${NC}"
}

check_server() {
  print_step "檢查 /health"
  if curl_api GET "$BASE_URL/health" >/dev/null; then
    print_ok "服務器運行中"
  else
    print_fail "服務器未運行，請先執行 go run ./cmd/server"
    exit 1
  fi
}

get_stats() {
  print_step "拉取統計"
  local companies users projects bots destinations notifications
  companies=$(curl_api GET "$API_BASE/companies" | jq '.data | length' 2>/dev/null || echo "0")
  users=$(curl_api GET "$API_BASE/users" | jq '.data | length' 2>/dev/null || echo "0")
  projects=$(curl_api GET "$API_BASE/projects" | jq '.data | length' 2>/dev/null || echo "0")
  bots=$(curl_api GET "$API_BASE/bots/platform" | jq '.data | length' 2>/dev/null || echo "0")
  destinations=$(curl_api GET "$API_BASE/destinations" | jq '.data | length' 2>/dev/null || echo "0")
  notifications=$(curl_api GET "$API_BASE/notifications" | jq '.data | length' 2>/dev/null || echo "0")

  printf "${GREEN}📊 統計${NC}\n"
  printf "  公司: %s\n" "$companies"
  printf "  用戶: %s\n" "$users"
  printf "  專案: %s\n" "$projects"
  printf "  平台 Bot: %s\n" "$bots"
  printf "  目的地: %s\n" "$destinations"
  printf "  通知: %s\n" "$notifications"
}

create_company() {
  print_step "建立測試公司"
  local payload response company_id
  payload='{
    "name": "API 測試公司",
    "contact_email": "test@example.com",
    "contact_phone": "+886-2-1111-2222",
    "address": "台北市測試區測試路 123 號",
    "billing_enabled": true
  }'
  response=$(curl_api POST "$API_BASE/companies" "$payload") || {
    print_fail "建立公司失敗"
    echo "$response"
    return 1
  }
  if echo "$response" | jq -e '.error' >/dev/null 2>&1; then
    print_fail "建立公司失敗"
    echo "$response" | jq
    return 1
  fi
  company_id=$(echo "$response" | jq -er '.data.id' 2>/dev/null || true)
  if [[ -n "$company_id" ]]; then
    printf "${GREEN}公司建立成功: %s${NC}\n" "$company_id"
  else
    print_skip "未取得公司 ID，可能因重複資料"
  fi
}

create_destination_and_notification() {
  if [[ -z "$PROJECT_ID" || -z "$CREATED_BY" || -z "$TENANT_ID" ]]; then
    print_skip "未設定 PROJECT_ID / CREATED_BY / TENANT_ID，跳過目的地與通知測試"
    return
  fi

  print_step "建立目的地"
  local payload response destination_id
  payload=$(jq -n --arg project "$PROJECT_ID" \
                --arg name "API Smoke Destination" \
                --arg desc "Smoke 測試目的地" \
                --arg tenant "$TENANT_ID" \
                --arg created "$CREATED_BY" '
    {
      project_id: $project,
      name: $name,
      description: $desc,
      teams_tenant_id: $tenant,
      targets: [
        {
          type: "channel",
          conversation_id: "19:fake-conversation-id@thread.tacv2",
          display_name: "Smoke 測試",
          tenant_id: $tenant
        }
      ],
      created_by: $created
    }
  ')

  response=$(curl_api POST "$API_BASE/destinations" "$payload") || {
    print_fail "建立目的地失敗"
    echo "$response"
    return
  }

  if echo "$response" | jq -e '.error' >/dev/null 2>&1; then
    print_fail "建立目的地失敗"
    echo "$response" | jq
    return
  fi

  destination_id=$(echo "$response" | jq -er '.data.id' 2>/dev/null || true)
  if [[ -z "$destination_id" ]]; then
    print_skip "無法解析目的地 ID，停止通知測試"
    return
  fi
  print_ok "目的地建立成功: $destination_id"

  print_step "發送通知"
  local notify_payload sender
  sender="$SENDER_ID"
  notify_payload=$(jq -n --arg project "$PROJECT_ID" \
                         --argjson destinations "[\"$destination_id\"]" \
                         --arg content "API smoke 測試通知" \
                         --arg priority "normal" \
                         --arg sender "$sender" '
    {
      project_id: $project,
      sender_id: ($sender | select(length > 0)),
      message_type: "text",
      content: $content,
      priority: $priority,
      destinations: $destinations
    }'
  )

  response=$(curl_api POST "$API_BASE/notifications" "$notify_payload") || {
    print_fail "通知發送失敗"
    echo "$response"
    return
  }

  if echo "$response" | jq -e '.error' >/dev/null 2>&1; then
    print_fail "通知發送失敗"
    echo "$response" | jq
    return
  fi
  print_ok "通知已排入佇列"
}

main() {
  require_tool curl
  require_tool jq

  check_server
  get_stats
  create_company || true
  create_destination_and_notification || true

  print_ok "腳本完成"
}

main "$@"
