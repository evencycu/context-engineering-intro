package actor

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/evencycu/TeamsNotifyGoV2/internal/database"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

// NotificationActor handles sending a single notification with retry logic
// Each actor is responsible for one notification_destination
type NotificationActor struct {
	ID                 string
	NotificationDestID uuid.UUID
	NotificationID     uuid.UUID
	ConversationID     string
	Message            string
	BotAppID           string
	BotAppPassword     string
	TenantID           string
	ServiceURL         string
	Installation       *database.BotInstallation
	Redis              *redis.Client
	DB                 ActorDB
	MaxRetries         int
	CurrentRetry       int
	Status             string
	LastError          error
	StopChan           chan struct{}
	DoneChan           chan *ActorResult
}

// ActorDB defines the database operations needed by actors
type ActorDB interface {
	UpdateNotificationDestination(ctx context.Context, id uuid.UUID, update *NotificationDestinationUpdate) error
	GetBotInstallation(ctx context.Context, conversationID string) (*database.BotInstallation, error)
	GetNotificationDestinationByID(ctx context.Context, id uuid.UUID) (*database.NotificationDestination, error)
	GetRetryReadyNotificationDestinations(ctx context.Context, limit int) ([]*database.NotificationDestination, error)
}

// NotificationDestinationUpdate represents fields to update
type NotificationDestinationUpdate struct {
	Status         *string
	ErrorMessage   *string
	TeamsMessageID *string
	SentAt         *time.Time
	RetryCount     *int
	NextRetryAt    *time.Time
	FirstAttemptAt *time.Time
	LastAttemptAt  *time.Time
	FailureReason  *string
	RetryAfter     *int
	ActorID        *string
}

// ActorResult represents the result of an actor's execution
type ActorResult struct {
	ActorID            string
	NotificationDestID uuid.UUID
	Success            bool
	Error              error
	Retrying           bool
	NextRetryAt        *time.Time
	TeamsMessageID     string
}

// NewNotificationActor creates a new notification actor
func NewNotificationActor(
	notificationDestID uuid.UUID,
	notificationID uuid.UUID,
	conversationID string,
	message string,
	installation *database.BotInstallation,
	botAppID string,
	botAppPassword string,
	redis *redis.Client,
	db ActorDB,
) *NotificationActor {
	return &NotificationActor{
		ID:                 fmt.Sprintf("actor-%s", notificationDestID.String()[:8]),
		NotificationDestID: notificationDestID,
		NotificationID:     notificationID,
		ConversationID:     conversationID,
		Message:            message,
		BotAppID:           botAppID,
		BotAppPassword:     botAppPassword,
		TenantID:           installation.TeamsTenantID,
		ServiceURL:         installation.ServiceURL,
		Installation:       installation,
		Redis:              redis,
		DB:                 db,
		MaxRetries:         5,
		CurrentRetry:       0,
		Status:             "pending",
		StopChan:           make(chan struct{}),
		DoneChan:           make(chan *ActorResult, 1),
	}
}

// Start begins the actor's execution
func (a *NotificationActor) Start(ctx context.Context) {
	go a.run(ctx)
}

// Stop stops the actor
func (a *NotificationActor) Stop() {
	close(a.StopChan)
}

// Wait waits for the actor to complete and returns the result
func (a *NotificationActor) Wait() *ActorResult {
	return <-a.DoneChan
}

// run is the main actor loop
func (a *NotificationActor) run(ctx context.Context) {
	defer close(a.DoneChan)

	log.Printf("[%s] Starting actor for notification_dest=%s", a.ID, a.NotificationDestID)

	// Mark as processing
	now := time.Now()
	processing := "processing"
	a.DB.UpdateNotificationDestination(ctx, a.NotificationDestID, &NotificationDestinationUpdate{
		Status:         &processing,
		ActorID:        &a.ID,
		FirstAttemptAt: &now,
		LastAttemptAt:  &now,
	})

	// Main retry loop
	for a.CurrentRetry < a.MaxRetries {
		select {
		case <-a.StopChan:
			log.Printf("[%s] Stopped by signal", a.ID)
			a.DoneChan <- &ActorResult{
				ActorID:            a.ID,
				NotificationDestID: a.NotificationDestID,
				Success:            false,
				Error:              fmt.Errorf("stopped"),
			}
			return
		case <-ctx.Done():
			log.Printf("[%s] Context cancelled", a.ID)
			a.DoneChan <- &ActorResult{
				ActorID:            a.ID,
				NotificationDestID: a.NotificationDestID,
				Success:            false,
				Error:              ctx.Err(),
			}
			return
		default:
		}

		// Check circuit breaker before sending
		if !a.checkCircuitBreaker(ctx) {
			log.Printf("[%s] Circuit breaker is open, scheduling retry", a.ID)
			a.scheduleRetry(ctx, "circuit_breaker_open", nil, nil)
			a.DoneChan <- &ActorResult{
				ActorID:            a.ID,
				NotificationDestID: a.NotificationDestID,
				Success:            false,
				Retrying:           true,
				NextRetryAt:        a.getNextRetryTime(0, nil),
			}
			return
		}

		// Attempt to send
		success, statusCode, retryAfter, teamsMessageID, err := a.attemptSend(ctx)

		// Update last attempt time
		attemptTime := time.Now()
		a.DB.UpdateNotificationDestination(ctx, a.NotificationDestID, &NotificationDestinationUpdate{
			LastAttemptAt: &attemptTime,
		})

		if success {
			// Success!
			log.Printf("[%s] Successfully sent notification", a.ID)
			sentStatus := "sent"
			sentAt := time.Now()
			a.DB.UpdateNotificationDestination(ctx, a.NotificationDestID, &NotificationDestinationUpdate{
				Status:         &sentStatus,
				TeamsMessageID: &teamsMessageID,
				SentAt:         &sentAt,
				ErrorMessage:   nil,
			})

			// Update Redis circuit breaker
			a.Redis.Incr(ctx, "cb:success_count")
			a.Redis.Set(ctx, "cb:failures", 0, 24*time.Hour)

			a.DoneChan <- &ActorResult{
				ActorID:            a.ID,
				NotificationDestID: a.NotificationDestID,
				Success:            true,
				TeamsMessageID:     teamsMessageID,
			}
			return
		}

		// Failed - determine if retryable
		failureReason := a.classifyFailure(statusCode)
		isRetryable := a.isRetryableError(statusCode)

		a.CurrentRetry++

		if !isRetryable || a.CurrentRetry >= a.MaxRetries {
			// Non-retryable or exhausted retries
			log.Printf("[%s] Non-retryable error or exhausted retries: %v", a.ID, err)
			failedStatus := "failed"
			errMsg := fmt.Sprintf("%v", err)
			a.DB.UpdateNotificationDestination(ctx, a.NotificationDestID, &NotificationDestinationUpdate{
				Status:        &failedStatus,
				ErrorMessage:  &errMsg,
				RetryCount:    &a.CurrentRetry,
				FailureReason: &failureReason,
			})

			// Update Redis circuit breaker
			a.Redis.Incr(ctx, "cb:failures")
			a.Redis.Incr(ctx, "cb:failed_count")

			a.DoneChan <- &ActorResult{
				ActorID:            a.ID,
				NotificationDestID: a.NotificationDestID,
				Success:            false,
				Error:              err,
			}
			return
		}

		// Schedule retry
		a.scheduleRetry(ctx, failureReason, retryAfter, err)

		// Update Redis circuit breaker
		a.Redis.Incr(ctx, "cb:failures")
		a.Redis.Incr(ctx, "cb:failed_count")
		if retryAfter != nil {
			a.Redis.Set(ctx, "cb:last_429_retry_after", *retryAfter, 1*time.Hour)
		}

		a.DoneChan <- &ActorResult{
			ActorID:            a.ID,
			NotificationDestID: a.NotificationDestID,
			Success:            false,
			Retrying:           true,
			NextRetryAt:        a.getNextRetryTime(a.CurrentRetry, retryAfter),
		}
		return
	}
}

// checkCircuitBreaker checks Redis circuit breaker state
func (a *NotificationActor) checkCircuitBreaker(ctx context.Context) bool {
	state, err := a.Redis.Get(ctx, "cb:state").Result()
	if err != nil {
		if err == redis.Nil {
			// No state set, assume closed
			return true
		}
		log.Printf("[%s] Failed to get circuit breaker state: %v", a.ID, err)
		return true // Fail open
	}

	if state == "open" {
		// Check if it's time to transition to half-open
		openUntilStr, err := a.Redis.Get(ctx, "cb:open_until").Result()
		if err == nil {
			var openUntil int64
			fmt.Sscanf(openUntilStr, "%d", &openUntil)
			if time.Now().Unix() >= openUntil {
				// Transition to half-open
				a.Redis.Set(ctx, "cb:state", "half-open", 24*time.Hour)
				a.Redis.Set(ctx, "cb:half_open_tests", 0, 24*time.Hour)
				return true
			}
		}
		return false
	}

	if state == "half-open" {
		// Check if we can be a test request
		tests, _ := a.Redis.Incr(ctx, "cb:half_open_tests").Result()
		if tests <= 3 {
			return true
		}
		return false
	}

	return true // closed state
}

// attemptSend attempts to send the notification to Teams
func (a *NotificationActor) attemptSend(ctx context.Context) (success bool, statusCode int, retryAfter *int, teamsMessageID string, err error) {
	log.Printf("[%s] Attempting to send (retry %d/%d)", a.ID, a.CurrentRetry, a.MaxRetries)

	// Get token
	token, err := a.getConnectorToken(ctx)
	if err != nil {
		return false, 401, nil, "", fmt.Errorf("failed to get token: %w", err)
	}

	// Prepare payload
	payload := map[string]any{
		"type": "message",
		"text": a.Message,
		"from": map[string]any{
			"id":   a.Installation.RecipientID,
			"name": "Notification Bot",
		},
		"recipient": map[string]any{
			"id": a.Installation.FromID,
		},
		"conversation": map[string]any{
			"id": a.ConversationID,
		},
		"channelId":  "msteams",
		"serviceUrl": a.ServiceURL,
	}

	reqBody, _ := json.Marshal(payload)
	url := a.ServiceURL + "/v3/conversations/" + a.ConversationID + "/activities"

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(string(reqBody)))
	if err != nil {
		return false, 0, nil, "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, 0, nil, "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	statusCode = resp.StatusCode

	// Extract Retry-After if present
	if resp.StatusCode == 429 {
		retryAfterStr := resp.Header.Get("Retry-After")
		if retryAfterStr != "" {
			var ra int
			fmt.Sscanf(retryAfterStr, "%d", &ra)
			retryAfter = &ra
		}
	}

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return false, statusCode, retryAfter, "", fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	// Parse response to get message ID
	var response map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&response); err == nil {
		if id, ok := response["id"].(string); ok {
			teamsMessageID = id
		}
	}

	return true, statusCode, nil, teamsMessageID, nil
}

// getConnectorToken gets a Bot Framework connector token
func (a *NotificationActor) getConnectorToken(ctx context.Context) (string, error) {
	tenantID := a.TenantID
	if tenantID == "" {
		tenantID = "botframework.com"
	}

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Set("client_id", a.BotAppID)
	form.Set("client_secret", a.BotAppPassword)
	form.Set("scope", "https://api.botframework.com/.default")

	tokenURL := "https://login.microsoftonline.com/" + url.PathEscape(tenantID) + "/oauth2/v2.0/token"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenURL, strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := &http.Client{Timeout: 10 * time.Second}
	resp, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("token request failed: %s: %s", resp.Status, string(body))
	}

	var tr struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&tr); err != nil {
		return "", err
	}

	if tr.AccessToken == "" {
		return "", fmt.Errorf("empty access_token")
	}

	return tr.AccessToken, nil
}

// isRetryableError determines if an error is retryable
func (a *NotificationActor) isRetryableError(statusCode int) bool {
	return statusCode == 429 || statusCode == 408 || (statusCode >= 500 && statusCode < 600)
}

// classifyFailure classifies the failure reason
func (a *NotificationActor) classifyFailure(statusCode int) string {
	switch {
	case statusCode == 429:
		return "rate_limit"
	case statusCode == 408:
		return "timeout"
	case statusCode >= 500 && statusCode < 600:
		return "server_error"
	case statusCode == 401 || statusCode == 403:
		return "unauthorized"
	case statusCode >= 400 && statusCode < 500:
		return "bad_request"
	default:
		return "network_error"
	}
}

// scheduleRetry schedules the next retry
func (a *NotificationActor) scheduleRetry(ctx context.Context, failureReason string, retryAfter *int, err error) {
	nextRetryAt := a.getNextRetryTime(a.CurrentRetry, retryAfter)

	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}

	failedStatus := "failed"
	a.DB.UpdateNotificationDestination(ctx, a.NotificationDestID, &NotificationDestinationUpdate{
		Status:        &failedStatus,
		ErrorMessage:  &errMsg,
		RetryCount:    &a.CurrentRetry,
		NextRetryAt:   nextRetryAt,
		FailureReason: &failureReason,
		RetryAfter:    retryAfter,
	})

	log.Printf("[%s] Scheduled retry at %v (retry %d/%d, reason: %s)",
		a.ID, nextRetryAt, a.CurrentRetry, a.MaxRetries, failureReason)
}

// getNextRetryTime calculates the next retry time
func (a *NotificationActor) getNextRetryTime(retryCount int, retryAfter *int) *time.Time {
	var backoff time.Duration

	if retryAfter != nil {
		// Priority: use Retry-After from 429 response
		backoff = time.Duration(*retryAfter) * time.Second
		backoff += 1 * time.Second // Add 1s buffer
		log.Printf("[%s] Using Retry-After: %d seconds (+ 1s buffer)", a.ID, *retryAfter)
	} else {
		// Use exponential backoff
		backoff = time.Duration(math.Pow(2, float64(retryCount))) * time.Second
		if backoff > 5*time.Minute {
			backoff = 5 * time.Minute
		}
		// Add jitter (±10%)
		jitter := backoff / 10
		backoff += time.Duration(rand.Int63n(int64(2*jitter))) - jitter
	}

	nextRetry := time.Now().Add(backoff)
	return &nextRetry
}
