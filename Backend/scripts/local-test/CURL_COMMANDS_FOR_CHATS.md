# CURL Commands for Testing Microsoft Graph API /chats Endpoint

## Prerequisites

You need the following information:
- `CLIENT_ID`: Your Azure AD App Registration Client ID (TEAMS_BOT_APP_ID)
- `CLIENT_SECRET`: Your Azure AD App Registration Client Secret (TEAMS_BOT_APP_PASSWORD)
- `TENANT_ID`: Your Azure AD Tenant ID (TEAMS_TENANT_ID)

## Step 1: Get Access Token

```bash
# Replace the placeholders with your actual values
CLIENT_ID="your-client-id"
CLIENT_SECRET="your-client-secret"
TENANT_ID="your-tenant-id"

# Get access token
TOKEN_RESPONSE=$(curl -s -X POST "https://login.microsoftonline.com/${TENANT_ID}/oauth2/v2.0/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials" \
  -d "client_id=${CLIENT_ID}" \
  -d "client_secret=${CLIENT_SECRET}" \
  -d "scope=https://graph.microsoft.com/.default")

# Extract access token
ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | jq -r '.access_token')

# Verify token
echo "Access Token: ${ACCESS_TOKEN:0:50}..."
```

## Step 2: Test GET /chats (All Chats)

```bash
# Get all chats
curl -X GET "https://graph.microsoft.com/v1.0/chats" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" | jq '.'
```

## Step 3: Test GET /chats (Personal Chats Only)

```bash
# Get personal chats (oneOnOne)
curl -X GET "https://graph.microsoft.com/v1.0/chats?\$filter=chatType eq 'oneOnOne'" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" | jq '.'
```

## Step 4: Test GET /chats (Group Chats Only)

```bash
# Get group chats
curl -X GET "https://graph.microsoft.com/v1.0/chats?\$filter=chatType eq 'groupChat'" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" | jq '.'
```

## Step 5: Test with Additional Query Parameters

```bash
# Get chats with specific fields
curl -X GET "https://graph.microsoft.com/v1.0/chats?\$select=id,chatType,topic,createdDateTime" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" | jq '.'

# Get chats with pagination (top 10)
curl -X GET "https://graph.microsoft.com/v1.0/chats?\$top=10" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" | jq '.'
```

## Complete One-Liner Example

```bash
# Replace these variables
CLIENT_ID="your-client-id"
CLIENT_SECRET="your-client-secret"
TENANT_ID="your-tenant-id"

# Get token and test in one go
TOKEN=$(curl -s -X POST "https://login.microsoftonline.com/${TENANT_ID}/oauth2/v2.0/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=client_credentials" \
  -d "client_id=${CLIENT_ID}" \
  -d "client_secret=${CLIENT_SECRET}" \
  -d "scope=https://graph.microsoft.com/.default" | jq -r '.access_token')

# Test /chats endpoint
curl -X GET "https://graph.microsoft.com/v1.0/chats" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" | jq '.'
```

## Expected Response Format

### Success Response
```json
{
  "@odata.context": "https://graph.microsoft.com/v1.0/$metadata#chats",
  "value": [
    {
      "id": "19:8d40db4b-935f-4e5a-aaae-ba86faad0e12_844146d7-4ac9-4e4d-a463-d6e027714e81@unq.gbl.spaces",
      "chatType": "oneOnOne",
      "topic": null,
      "createdDateTime": "2025-09-18T06:42:39.469Z",
      "webUrl": "https://teams.microsoft.com/l/chat/..."
    },
    {
      "id": "19:f26a8d8a235f430db87a404491cd2ffc@thread.v2",
      "chatType": "groupChat",
      "topic": "Team Discussion",
      "createdDateTime": "2025-09-18T07:00:00.000Z"
    }
  ]
}
```

### Error Response (Insufficient Privileges)
```json
{
  "error": {
    "code": "Forbidden",
    "message": "Insufficient privileges to complete the operation.",
    "innerError": {
      "date": "2025-01-14T12:00:00",
      "request-id": "...",
      "client-request-id": "..."
    }
  }
}
```

## Common Errors and Solutions

### Error: "Insufficient privileges"
**Solution**: 
1. Go to Azure AD App Registration
2. Add `Chat.Read.All` or `Chat.ReadWrite.All` permission (Application permission)
3. Get tenant admin consent

### Error: "Invalid client"
**Solution**: 
- Check if CLIENT_ID and CLIENT_SECRET are correct
- Verify the app registration exists in the tenant

### Error: "Invalid scope"
**Solution**: 
- Ensure the scope is `https://graph.microsoft.com/.default`
- Verify the app has the required permissions

## Using the Test Script

Alternatively, you can use the provided test script:

```bash
cd Backend/scripts/local-test
./test_graph_chats_manual.sh <CLIENT_ID> <CLIENT_SECRET> <TENANT_ID>
```

Example:
```bash
./test_graph_chats_manual.sh 844146d7-4ac9-4e4d-a463-d6e027714e81 'your-secret' 051cece0-e4dc-4aed-b471-bf29824e1ee6
```
