package main

import (
	"context"
	"fmt"
	"log"

	"github.com/teamsnotify/go-sdk"
)

func main() {
	// Initialize the client
	client := sdk.NewAPIClient(sdk.NewConfiguration())
	client.GetConfig().BasePath = "http://localhost:8080"

	// Test external API - Send notification
	fmt.Println("Testing external API - Send notification...")
	externalAPI := client.ExternalAPIAPI

	request := sdk.NewApiV1NotifyPostRequest(context.Background())
	request = request.ExternalNotifyRequest(sdk.ExternalNotifyRequest{
		NotifyKey:   "cfh-alert-gogo",
		Message:     "Hello from Go SDK!",
		MessageType: "text",
		Priority:    "normal",
		Targets:     []string{"all"},
	})

	resp, httpResp, err := externalAPI.ApiV1NotifyPost(*request).Execute()
	if err != nil {
		log.Printf("Error sending notification: %v", err)
		log.Printf("HTTP Response: %v", httpResp)
		return
	}

	fmt.Printf("✅ Notification sent successfully!\n")
	fmt.Printf("   Notification ID: %s\n", resp.GetNotificationId())
	fmt.Printf("   Status: %s\n", resp.GetStatus())
	fmt.Printf("   Destinations Count: %d\n", resp.GetDestinationsCount())

	// Test external API - Get project destinations
	fmt.Println("\nTesting external API - Get project destinations...")
	destsResp, httpResp, err := externalAPI.ApiV1DestinationsNotifyKeyGet(context.Background(), "cfh-alert-gogo").Execute()
	if err != nil {
		log.Printf("Error getting destinations: %v", err)
		log.Printf("HTTP Response: %v", httpResp)
		return
	}

	fmt.Printf("✅ Project destinations retrieved successfully!\n")
	fmt.Printf("   Found %d destinations\n", len(destsResp.GetDestinations()))
	for i, dest := range destsResp.GetDestinations() {
		fmt.Printf("   Destination %d: %s\n", i+1, dest.GetName())
	}

	fmt.Println("\n🎉 All SDK tests completed successfully!")
}
