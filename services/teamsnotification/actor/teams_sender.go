package actor

import (
	"context"
	"fmt"
	"log"

	"github.com/evencycu/TeamsNotifyGoV2/libs/models"
)

// teamsSender implements TeamsSender interface for sending messages to Microsoft Teams
type teamsSender struct {
	// This would typically contain Teams API client, credentials, etc.
	// For now, we'll use a mock implementation
}

// TeamsSender defines the interface for sending messages to Teams
type TeamsSender interface {
	Send(ctx context.Context, nd *database.NotificationDestination) (*SendResult, error)
}

// NewTeamsSender creates a new Teams sender
func NewTeamsSender() TeamsSender {
	return &teamsSender{}
}

// Send sends a notification to Teams
func (ts *teamsSender) Send(ctx context.Context, nd *database.NotificationDestination) (*SendResult, error) {
	log.Printf("TeamsSender: Sending notification_dest %s to conversation %s", nd.ID, *nd.ConversationID)

	// TODO: Implement actual Teams API call here
	// This would involve:
	// 1. Getting bot installation details
	// 2. Making HTTP request to Teams API
	// 3. Handling rate limits and errors
	// 4. Parsing Retry-After headers

	// Mock implementation for now
	if nd.RetryCount == 0 {
		// First attempt - simulate rate limit
		return &SendResult{
			Success:    false,
			Error:      fmt.Errorf("simulated rate limit"),
			RetryAfter: 10, // 10 seconds
		}, nil
	}

	// Subsequent attempts - simulate success
	return &SendResult{
		Success:        true,
		TeamsMessageID: "mock-message-id-" + nd.ID.String(),
		RetryAfter:     0,
	}, nil
}

// SendResult contains the outcome of a Teams message send attempt
type SendResult struct {
	Success        bool
	TeamsMessageID string
	RetryAfter     int // From Retry-After header, if any
	Error          error
}
