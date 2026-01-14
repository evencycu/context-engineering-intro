package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/evencycu/TeamsNotifyGoV3/libs/token"
)

const graphBaseURL = "https://graph.microsoft.com/v1.0"

// GraphService defines operations for Microsoft Graph API
type GraphService interface {
	GetUsers(ctx context.Context, deltaLink string) ([]models.AzureADUser, string, error)
	GetGroups(ctx context.Context, deltaLink string) ([]models.AzureADGroup, string, error)
	GetGroupChannels(ctx context.Context, groupId string) ([]models.AzureADChannel, error)
}


type graphService struct {
	tokenManager token.TokenManager
	tenantID     string
	botID        string
}

// NewGraphService creates a new graph service
func NewGraphService(tokenManager token.TokenManager, tenantID, botID string) GraphService {
	return &graphService{
		tokenManager: tokenManager,
		tenantID:     tenantID,
		botID:        botID,
	}
}

type graphUsersResponse struct {
	Value     []graphUser `json:"value"`
	NextLink  string      `json:"@odata.nextLink"`
	DeltaLink string      `json:"@odata.deltaLink"`
}

type graphUser struct {
	ID                string `json:"id"`
	DisplayName       string `json:"displayName"`
	Mail              string `json:"mail"`
	UserPrincipalName string `json:"userPrincipalName"`
	JobTitle          string `json:"jobTitle"`
	Department        string `json:"department"`
}

func (s *graphService) GetUsers(ctx context.Context, deltaLink string) ([]models.AzureADUser, string, error) {
	// Use delta query if link provided, otherwise initial sync
	requestURL := deltaLink
	if requestURL == "" {
		// Note: Graph API /users does not support sorting by lastModifiedDateTime directly.
		// To get incremental changes, we must rely on the Delta Query mechanism (deltaLink).
		// For the initial sync, we fetch all users.
		requestURL = fmt.Sprintf("%s/users?$select=id,displayName,mail,userPrincipalName,jobTitle,department,createdDateTime&$top=999", graphBaseURL)
	}

	token, err := s.tokenManager.GetGraphToken(ctx, s.botID, s.tenantID) // This gets the Bot token, assuming it has Graph permissions
	if err != nil {
		return nil, "", fmt.Errorf("failed to get token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("graph api call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("graph api error: %s", resp.Status)
	}

	var result graphUsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}

	users := make([]models.AzureADUser, 0, len(result.Value))
	for _, u := range result.Value {
		email := u.Mail
		if email == "" {
			email = u.UserPrincipalName
		}
		users = append(users, models.AzureADUser{
			AzureADID:   u.ID,
			DisplayName: u.DisplayName,
			Email:       email,
			JobTitle:    u.JobTitle,
			Department:  u.Department,
		})
	}

	nextLink := result.NextLink
	if result.DeltaLink != "" {
		nextLink = result.DeltaLink // Return delta link if available for future sync
	}

	return users, nextLink, nil
}

type graphGroupsResponse struct {
	Value     []graphGroup `json:"value"`
	NextLink  string       `json:"@odata.nextLink"`
	DeltaLink string       `json:"@odata.deltaLink"`
}

type graphGroup struct {
	ID          string   `json:"id"`
	DisplayName string   `json:"displayName"`
	Description string   `json:"description"`
	GroupTypes  []string `json:"groupTypes"`
}

func (s *graphService) GetGroups(ctx context.Context, deltaLink string) ([]models.AzureADGroup, string, error) {
	requestURL := deltaLink
	if requestURL == "" {
		requestURL = fmt.Sprintf("%s/groups?$select=id,displayName,description,groupTypes&$top=999", graphBaseURL)
	}

	token, err := s.tokenManager.GetGraphToken(ctx, s.botID, s.tenantID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, "", fmt.Errorf("graph api call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("graph api error: %s", resp.Status)
	}

	var result graphGroupsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, "", fmt.Errorf("failed to decode response: %w", err)
	}

	groups := make([]models.AzureADGroup, 0, len(result.Value))
	for _, g := range result.Value {
		groups = append(groups, models.AzureADGroup{
			AzureADID:   g.ID,
			DisplayName: g.DisplayName,
			Description: g.Description,
			GroupTypes:  models.JSONBStringArray(g.GroupTypes),
		})
	}

	nextLink := result.NextLink
	if result.DeltaLink != "" {
		nextLink = result.DeltaLink
	}

	return groups, nextLink, nil
}

type graphChannelsResponse struct {
	Value []graphChannel `json:"value"`
}

type graphChannel struct {
	ID             string `json:"id"`
	DisplayName    string `json:"displayName"`
	Description    string `json:"description"`
	MembershipType string `json:"membershipType"`
}

func (s *graphService) GetGroupChannels(ctx context.Context, groupId string) ([]models.AzureADChannel, error) {
	requestURL := fmt.Sprintf("%s/teams/%s/channels?$select=id,displayName,description,membershipType", graphBaseURL, groupId)

	token, err := s.tokenManager.GetGraphToken(ctx, s.botID, s.tenantID)
	if err != nil {
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", requestURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("graph api call failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("graph api error: %s", resp.Status)
	}

	var result graphChannelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	channels := make([]models.AzureADChannel, 0, len(result.Value))
	for _, c := range result.Value {
		channels = append(channels, models.AzureADChannel{
			AzureADID:      c.ID,
			TeamID:         groupId, // Parent Team/Group ID
			DisplayName:    c.DisplayName,
			Description:    c.Description,
			MembershipType: c.MembershipType,
			ConversationID: c.ID, // Channel ID is the conversation ID
		})
	}

	return channels, nil
}
