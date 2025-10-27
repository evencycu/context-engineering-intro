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
    "fmt"
    "log"
    "github.com/teamsnotify/go-sdk"
)

func main() {
    // Initialize the client
    client := sdk.NewClient("http://localhost:8080")
    
    // Send a notification
    resp, err := client.Notifications().Send(&sdk.SendNotificationRequest{
        NotifyKey: "your-notify-key",
        Message: "Hello from Go SDK!",
        MessageType: "text",
        Priority: "normal",
        Targets: []string{"all"},
    })
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Notification sent: %s\n", resp.NotificationID)
}
```

### External API Example

```go
// Use external API (no authentication required)
external := sdk.NewExternalClient("http://localhost:8080/api/v1")

// Send notification via external API
resp, err := external.SendNotification(&sdk.SendNotificationRequest{
    NotifyKey: "your-notify-key",
    Message: "Hello from external API!",
    Targets: []string{"all"},
})

// Get project destinations
dests, err := external.GetProjectDestinations("your-notify-key")
```

### Authentication

```go
// Create authenticated client
client := sdk.NewClientWithAuth(
    "http://localhost:8080",
    "your-api-key",
)

// Use protected endpoints
users, err := client.Users().List()
projects, err := client.Projects().List()
```

## API Reference

### External API

- `SendNotification(req *SendNotificationRequest) (*SendNotificationResponse, error)`
- `GetProjectDestinations(notifyKey string) (*ProjectDestinationsResponse, error)`

### Internal API

#### Users
- `List() ([]*User, error)`
- `Get(id string) (*User, error)`
- `Create(req *CreateUserRequest) (*User, error)`
- `Update(id string, req *UpdateUserRequest) (*User, error)`
- `Delete(id string) error`

#### Projects
- `List() ([]*Project, error)`
- `Get(id string) (*Project, error)`
- `Create(req *CreateProjectRequest) (*Project, error)`
- `Update(id string, req *UpdateProjectRequest) (*Project, error)`
- `Delete(id string) error`

#### Notifications
- `Send(req *SendNotificationRequest) (*SendNotificationResponse, error)`
- `Get(id string) (*Notification, error)`
- `List(query *NotificationQuery) ([]*Notification, error)`
- `Cancel(id string) error`

## License

MIT License

