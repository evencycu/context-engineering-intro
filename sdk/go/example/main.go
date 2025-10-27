package main

import (
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
            NotifyKey: "cfh-alert-gogo",
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
    
    // Get project destinations
    dests, err := client.ExternalAPI().ApiV1DestinationsNotifyKeyGet(context.Background(), "cfh-alert-gogo")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d destinations\n", len(dests.Destinations))
}
