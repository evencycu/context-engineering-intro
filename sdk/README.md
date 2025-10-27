# Teams Notification SDKs

This directory contains automatically generated SDKs for the Teams Notification API.

## Available SDKs

### Go SDK
- 📦 **Package**: `github.com/teamsnotify/go-sdk`
- 📁 **Location**: `sdk/go/`
- 📖 **Documentation**: [Go SDK README](./go/README.md)

### Java SDK
- 📦 **Package**: `com.teamsnotify:teams-notify-sdk`
- 📁 **Location**: `sdk/java/`
- 📖 **Documentation**: [Java SDK README](./java/README.md)

## Generating SDKs

SDKs are automatically generated from the OpenAPI specification using `openapi-generator`.

### Prerequisites

Install `openapi-generator`:

```bash
# Using npm
npm install -g @openapitools/openapi-generator-cli

# Using homebrew (macOS)
brew install openapi-generator
```

### Generate SDKs

```bash
# Make the script executable
chmod +x sdk/generate.sh

# Run the generation script
./sdk/generate.sh
```

This will generate both Go and Java SDKs based on `api/openapi/teams-notification-api.yaml`.

## Manual Generation

You can also generate SDKs manually:

### Go SDK

```bash
openapi-generator-cli generate \
    -i api/openapi/teams-notification-api.yaml \
    -g go \
    -o sdk/go \
    --additional-properties=packageName=teamsnotify
```

### Java SDK

```bash
openapi-generator-cli generate \
    -i api/openapi/teams-notification-api.yaml \
    -g java \
    -o sdk/java \
    --additional-properties=groupId=com.teamsnotify \
    --additional-properties=artifactId=teams-notify-sdk \
    --additional-properties=library=okhttp-gson
```

## Integration

### Go Project

```go
import "github.com/teamsnotify/go-sdk"

client := sdk.NewClient("http://localhost:8080")
// Use client...
```

### Java Project

```xml
<dependency>
    <groupId>com.teamsnotify</groupId>
    <artifactId>teams-notify-sdk</artifactId>
    <version>1.0.0</version>
</dependency>
```

## Updating SDKs

When the API changes, regenerate the SDKs:

```bash
# Update the OpenAPI spec
# ...

# Regenerate SDKs
./sdk/generate.sh

# Commit changes
git add sdk/
git commit -m "chore: update SDKs"
```

## Building and Publishing

### Go SDK

```bash
cd sdk/go
go mod tidy
go build ./...

# Publish to GitHub
git tag -a v1.0.0 -m "Release v1.0.0"
git push origin v1.0.0
```

### Java SDK

```bash
cd sdk/java
mvn clean install

# Publish to Maven Central
mvn deploy
```

## Contributing

When adding new API endpoints:

1. Update `api/openapi/teams-notification-api.yaml`
2. Run `./sdk/generate.sh` to regenerate SDKs
3. Update SDK examples in their respective README files
4. Test the SDKs with the updated API
5. Commit changes

## Support

For issues and questions about the SDKs:

- **Go SDK**: See [Go SDK README](./go/README.md)
- **Java SDK**: See [Java SDK README](./java/README.md)
- **General**: Open an issue in the main repository

