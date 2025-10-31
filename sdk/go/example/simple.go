package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// Simple external API test without generated SDK
func main() {
	baseURL := "http://localhost:8080"

	// Test 1: Send notification
	fmt.Println("Testing external API - Send notification...")

	requestBody := map[string]interface{}{
		"notifyKey":   "cfh-alert-gogo",
		"message":     "Hello from Go SDK test!",
		"messageType": "text",
		"priority":    "normal",
		"targets":     []string{"all"},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		log.Fatal("Error marshaling JSON:", err)
	}

	resp, err := http.Post(baseURL+"/api/v1/notify", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error sending request:", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		log.Fatal("Error decoding response:", err)
	}

	if resp.StatusCode == 200 {
		fmt.Printf("✅ Notification sent successfully!\n")
		fmt.Printf("   Response: %+v\n", result)
	} else {
		fmt.Printf("❌ Notification failed with status %d\n", resp.StatusCode)
		fmt.Printf("   Response: %+v\n", result)
	}

	// Test 2: Get project destinations
	fmt.Println("\nTesting external API - Get project destinations...")

	resp2, err := http.Get(baseURL + "/api/v1/destinations/cfh-alert-gogo")
	if err != nil {
		log.Fatal("Error getting destinations:", err)
	}
	defer resp2.Body.Close()

	var destsResult map[string]interface{}
	err = json.NewDecoder(resp2.Body).Decode(&destsResult)
	if err != nil {
		log.Fatal("Error decoding destinations response:", err)
	}

	if resp2.StatusCode == 200 {
		fmt.Printf("✅ Project destinations retrieved successfully!\n")
		fmt.Printf("   Response: %+v\n", destsResult)
	} else {
		fmt.Printf("❌ Get destinations failed with status %d\n", resp2.StatusCode)
		fmt.Printf("   Response: %+v\n", destsResult)
	}

	fmt.Println("\n🎉 All external API tests completed!")
}
