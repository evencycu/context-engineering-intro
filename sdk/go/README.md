# Teams Notification Go SDK

Go SDK for Teams Notification API

## Installation

```bash
go mod init github.com/teamsnotify/go-sdk
go get github.com/teamsnotify/go-sdk
```

## Usage

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "github.com/teamsnotify/go-sdk"
)

func main() {
    // Initialize the client
    client := sdk.NewClient("http://localhost:8080")
    
    // Send a notification using external API
    resp, err := client.ExternalAPI().ApiV1NotifyPost(context.Background(), sdk.ApiV1NotifyPostRequest{
        ExternalNotifyRequest: sdk.ExternalNotifyRequest{
            NotifyKey: "your-notify-key",
            Message: "Hello from Go SDK!",
            MessageType: "text",
            Priority: "normal",
            Targets: []string{"all"},
        },
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Notification sent: %s\n", resp.NotificationId)
}
```

### External API Example

```go
// Send notification via external API
resp, err := client.ExternalAPI().ApiV1NotifyPost(context.Background(), sdk.ApiV1NotifyPostRequest{
    ExternalNotifyRequest: sdk.ExternalNotifyRequest{
        NotifyKey: "your-notify-key",
        Message: "Hello from external API!",
        Targets: []string{"all"},
    },
})

// Get project destinations
dests, err := client.ExternalAPI().ApiV1DestinationsNotifyKeyGet(context.Background(), "your-notify-key")
```

### Internal API Example

```go
// Create authenticated client
client := sdk.NewClientWithAuth("http://localhost:8080", "your-api-key")

// Use protected endpoints
users, err := client.InternalAPIUsers().InternalV1UsersGet(context.Background())
projects, err := client.InternalAPIProjects().InternalV1ProjectsGet(context.Background())
```

## API Reference

### External API
- `ApiV1NotifyPost` - Send notification
- `ApiV1DestinationsNotifyKeyGet` - Get project destinations

### Internal API
- Users: `InternalV1UsersGet`, `InternalV1UsersPost`, etc.
- Projects: `InternalV1ProjectsGet`, `InternalV1ProjectsPost`, etc.
- Notifications: `InternalV1NotificationsGet`, `InternalV1NotificationsPost`, etc.
- Destinations: `InternalV1DestinationsGet`, `InternalV1DestinationsPost`, etc.
- Bots: `InternalV1BotsPlatformGet`, `InternalV1BotsPlatformPost`, etc.

## License

MIT License