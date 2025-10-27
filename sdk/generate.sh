#!/bin/bash

# Teams Notification SDK Generation Script
# This script generates Go and Java SDKs from OpenAPI specification

set -e

echo "🚀 Generating Teams Notification SDKs..."

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if openapi-generator is installed
if ! command -v openapi-generator &> /dev/null; then
    echo -e "${YELLOW}⚠️  openapi-generator not found. Installing...${NC}"
    
    # Try to install using various methods
    if command -v npm &> /dev/null; then
        echo "Installing via npm..."
        npm install -g @openapitools/openapi-generator-cli
    elif command -v brew &> /dev/null; then
        echo "Installing via homebrew..."
        brew install openapi-generator
    else
        echo "Please install openapi-generator:"
        echo "  npm install -g @openapitools/openapi-generator-cli"
        echo "  or"
        echo "  brew install openapi-generator"
        exit 1
    fi
fi

# Paths
OPENAPI_SPEC="api/openapi/teams-notification-api.yaml"
GO_OUTPUT_DIR="sdk/go"
JAVA_OUTPUT_DIR="sdk/java"

echo -e "${BLUE}📋 OpenAPI Spec: ${OPENAPI_SPEC}${NC}"

# Generate Go SDK
echo -e "\n${BLUE}🔨 Generating Go SDK...${NC}"
openapi-generator-cli generate \
    -i "${OPENAPI_SPEC}" \
    -g go \
    -o "${GO_OUTPUT_DIR}" \
    --additional-properties=packageName=teamsnotify \
    --additional-properties=packageVersion=1.0.0 \
    --skip-validate-spec

echo -e "${GREEN}✅ Go SDK generated successfully${NC}"

# Generate Java SDK
echo -e "\n${BLUE}🔨 Generating Java SDK...${NC}"
openapi-generator-cli generate \
    -i "${OPENAPI_SPEC}" \
    -g java \
    -o "${JAVA_OUTPUT_DIR}" \
    --additional-properties=groupId=com.teamsnotify \
    --additional-properties=artifactId=teams-notify-sdk \
    --additional-properties=artifactVersion=1.0.0 \
    --additional-properties=library=okhttp-gson \
    --skip-validate-spec

echo -e "${GREEN}✅ Java SDK generated successfully${NC}"

# Create Go SDK helper files
echo -e "\n${BLUE}📝 Creating Go SDK helper files...${NC}"
cat > "${GO_OUTPUT_DIR}/client.go" << 'EOF'
package teamsnotify

// Client provides access to all Teams Notification API services
type Client struct {
    baseURL    string
    apiKey     string
    httpClient *HTTPClient
}

// NewClient creates a new Client
func NewClient(baseURL string) *Client {
    return &Client{
        baseURL:    baseURL,
        httpClient: NewHTTPClient(),
    }
}

// NewClientWithAuth creates a new Client with authentication
func NewClientWithAuth(baseURL, apiKey string) *Client {
    return &Client{
        baseURL:    baseURL,
        apiKey:     apiKey,
        httpClient: NewHTTPClient().WithAuth(apiKey),
    }
}

// Notifications returns notifications service
func (c *Client) Notifications() *NotificationService {
    return NewNotificationService(c)
}

// Users returns users service
func (c *Client) Users() *UserService {
    return NewUserService(c)
}

// Projects returns projects service
func (c *Client) Projects() *ProjectService {
    return NewProjectService(c)
}

// Destinations returns destinations service
func (c *Client) Destinations() *DestinationService {
    return NewDestinationService(c)
}

// Bots returns bots service
func (c *Client) Bots() *BotService {
    return NewBotService(c)
}
EOF

# Create Java SDK helper files
echo -e "${BLUE}📝 Creating Java SDK helper files...${NC}"
cat > "${JAVA_OUTPUT_DIR}/Client.java" << 'EOF'
package com.teamsnotify.sdk;

import com.teamsnotify.sdk.services.*;

public class Client {
    private final String baseUrl;
    private final String apiKey;
    private final OkHttpClient httpClient;
    
    private Client(Builder builder) {
        this.baseUrl = builder.baseUrl;
        this.apiKey = builder.apiKey;
        this.httpClient = new OkHttpClient();
    }
    
    public static class Builder {
        private String baseUrl;
        private String apiKey;
        
        public Builder baseUrl(String baseUrl) {
            this.baseUrl = baseUrl;
            return this;
        }
        
        public Builder apiKey(String apiKey) {
            this.apiKey = apiKey;
            return this;
        }
        
        public Client build() {
            return new Client(this);
        }
    }
    
    public NotificationService notifications() {
        return new NotificationService(this);
    }
    
    public UserService users() {
        return new UserService(this);
    }
    
    public ProjectService projects() {
        return new ProjectService(this);
    }
    
    // ... other services
}
EOF

echo -e "\n${GREEN}🎉 SDK generation completed!${NC}"
echo -e "${BLUE}📦 Go SDK: ${GO_OUTPUT_DIR}${NC}"
echo -e "${BLUE}📦 Java SDK: ${JAVA_OUTPUT_DIR}${NC}"

