package services

import (
	"context"
	"errors"
	"sync"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/repositories"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// DirectoryService defines directory operations
type DirectoryService interface {
	SearchUsers(ctx context.Context, query string) ([]models.AzureADUser, error)
	ListUsers(ctx context.Context, limit, offset int) ([]models.AzureADUser, error)
	ListGroups(ctx context.Context) ([]models.AzureADGroup, error)
	ListTeamsGroups(ctx context.Context) ([]models.AzureADGroup, error)
	ListTeamsGroupsWithChannels(ctx context.Context) ([]models.TeamsGroupWithChannels, error)
	ListGroupChannels(ctx context.Context, groupId string) ([]models.AzureADChannel, error)
	ListChannels(ctx context.Context, teamID string) ([]models.AzureADChannel, error)
	SyncDirectory(ctx context.Context) error
	GetSyncStatus(ctx context.Context) (*models.DirectorySyncStatus, error)
	// ChatGroup operations
	RegisterChatGroup(ctx context.Context, projectID uuid.UUID, name, chatID string) (*models.ChatGroup, error)
	ListChatGroups(ctx context.Context, projectID uuid.UUID) ([]models.ChatGroup, error)
	ListAllChatGroups(ctx context.Context) ([]models.ChatGroup, error)
	RemoveChatGroup(ctx context.Context, id uuid.UUID) error
	// Group chats from bot_installations
	GetGroupChatsFromBotInstallations(ctx context.Context) ([]models.GroupChatFromBotInstallation, error)
}

type directoryService struct {
	repo         repositories.DirectoryRepository
	graphService GraphService
	logger       *logrus.Logger
	syncMutex    sync.Mutex
	isSyncing    bool
}

// NewDirectoryService creates a new directory service
func NewDirectoryService(repo repositories.DirectoryRepository, graphService GraphService) DirectoryService {
	return &directoryService{
		repo:         repo,
		graphService: graphService,
		logger:       logrus.New(),
	}
}

func (s *directoryService) SearchUsers(ctx context.Context, query string) ([]models.AzureADUser, error) {
	return s.repo.SearchUsers(ctx, query)
}

func (s *directoryService) ListUsers(ctx context.Context, limit, offset int) ([]models.AzureADUser, error) {
	return s.repo.ListUsers(ctx, limit, offset)
}

func (s *directoryService) ListGroups(ctx context.Context) ([]models.AzureADGroup, error) {
	return s.repo.ListGroups(ctx)
}

func (s *directoryService) ListTeamsGroups(ctx context.Context) ([]models.AzureADGroup, error) {
	return s.repo.ListTeamsGroups(ctx)
}

func (s *directoryService) ListTeamsGroupsWithChannels(ctx context.Context) ([]models.TeamsGroupWithChannels, error) {
	return s.repo.ListTeamsGroupsWithChannels(ctx)
}

func (s *directoryService) ListGroupChannels(ctx context.Context, groupId string) ([]models.AzureADChannel, error) {
	// Get channels from database (synced channels) instead of Graph API
	return s.repo.ListChannels(ctx, groupId)
}

func (s *directoryService) ListChannels(ctx context.Context, teamID string) ([]models.AzureADChannel, error) {
	return s.repo.ListChannels(ctx, teamID)
}

func (s *directoryService) RegisterChatGroup(ctx context.Context, projectID uuid.UUID, name, chatID string) (*models.ChatGroup, error) {
	cg := &models.ChatGroup{
		ProjectID: projectID,
		Name:      name,
		ChatID:    chatID,
	}
	if err := s.repo.CreateChatGroup(ctx, cg); err != nil {
		return nil, err
	}
	return cg, nil
}

func (s *directoryService) ListChatGroups(ctx context.Context, projectID uuid.UUID) ([]models.ChatGroup, error) {
	return s.repo.ListChatGroups(ctx, projectID)
}

func (s *directoryService) ListAllChatGroups(ctx context.Context) ([]models.ChatGroup, error) {
	return s.repo.ListAllChatGroups(ctx)
}

func (s *directoryService) RemoveChatGroup(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteChatGroup(ctx, id)
}

func (s *directoryService) GetGroupChatsFromBotInstallations(ctx context.Context) ([]models.GroupChatFromBotInstallation, error) {
	return s.repo.GetGroupChatsFromBotInstallations(ctx)
}

func (s *directoryService) SyncDirectory(ctx context.Context) error {
	s.syncMutex.Lock()
	if s.isSyncing {
		s.syncMutex.Unlock()
		return errors.New("sync already in progress")
	}
	s.isSyncing = true
	s.syncMutex.Unlock()

	defer func() {
		s.syncMutex.Lock()
		s.isSyncing = false
		s.syncMutex.Unlock()
	}()

	s.logger.Info("Starting directory sync...")

	// Sync Users
	totalUsers := 0
	nextLink := ""
	for {
		users, newNextLink, err := s.graphService.GetUsers(ctx, nextLink)
		if err != nil {
			s.logger.Errorf("Failed to fetch users from Graph: %v", err)
			break // Stop user sync on error, try groups
		}

		for _, u := range users {
			if err := s.repo.UpsertUser(ctx, &u); err != nil {
				s.logger.Errorf("Failed to upsert user %s: %v", u.AzureADID, err)
			}
		}
		totalUsers += len(users)
		s.logger.Infof("Synced batch of %d users (total: %d)", len(users), totalUsers)

		if newNextLink == "" {
			break
		}
		nextLink = newNextLink
	}
	s.logger.Infof("User sync completed. Total users: %d", totalUsers)

	// Sync Groups
	totalGroups := 0
	nextLink = ""
	for {
		groups, newNextLink, err := s.graphService.GetGroups(ctx, nextLink)
		if err != nil {
			s.logger.Errorf("Failed to fetch groups from Graph: %v", err)
			return err
		}
		for _, g := range groups {
			if err := s.repo.UpsertGroup(ctx, &g); err != nil {
				s.logger.Errorf("Failed to upsert group %s: %v", g.AzureADID, err)
			}
		}
		totalGroups += len(groups)
		s.logger.Infof("Synced batch of %d groups (total: %d)", len(groups), totalGroups)

		if newNextLink == "" {
			break
		}
		nextLink = newNextLink
	}
	s.logger.Infof("Group sync completed. Total groups: %d", totalGroups)

	// Sync Channels from Unified Groups
	s.logger.Info("Starting channels sync from Unified Groups...")
	allGroups, err := s.repo.ListGroups(ctx)
	if err != nil {
		s.logger.Errorf("Failed to list groups for channel sync: %v", err)
		return err
	}

	totalChannels := 0
	unifiedGroupCount := 0
	for _, group := range allGroups {
		// Check if group is a Unified (M365) Group
		isUnified := false
		for _, gt := range group.GroupTypes {
			if gt == "Unified" {
				isUnified = true
				break
			}
		}
		if !isUnified {
			continue
		}

		unifiedGroupCount++
		s.logger.Infof("Syncing channels for Unified Group: %s (%s)", group.DisplayName, group.AzureADID)

		channels, err := s.graphService.GetGroupChannels(ctx, group.AzureADID)
		if err != nil {
			s.logger.Warnf("Failed to fetch channels for group %s: %v (may not be a Teams team)", group.AzureADID, err)
			continue // Not all Unified Groups are Teams teams
		}

		for _, ch := range channels {
			ch.TeamID = group.AzureADID // Ensure team_id is set
			if err := s.repo.UpsertChannel(ctx, &ch); err != nil {
				s.logger.Errorf("Failed to upsert channel %s: %v", ch.AzureADID, err)
			} else {
				totalChannels++
			}
		}
		s.logger.Infof("Synced %d channels for group %s", len(channels), group.DisplayName)
	}
	s.logger.Infof("Channels sync completed. Processed %d Unified Groups, synced %d channels", unifiedGroupCount, totalChannels)

	return nil
}

func (s *directoryService) GetSyncStatus(ctx context.Context) (*models.DirectorySyncStatus, error) {
	s.syncMutex.Lock()
	isSyncing := s.isSyncing
	s.syncMutex.Unlock()

	userCount, err := s.repo.CountUsers(ctx)
	if err != nil {
		s.logger.Errorf("Failed to count users: %v", err)
		userCount = 0
	}

	groupCount, err := s.repo.CountGroups(ctx)
	if err != nil {
		s.logger.Errorf("Failed to count groups: %v", err)
		groupCount = 0
	}

	channelCount, err := s.repo.CountChannels(ctx)
	if err != nil {
		s.logger.Errorf("Failed to count channels: %v", err)
		channelCount = 0
	}

	lastSyncAt, err := s.repo.GetLastSyncTime(ctx)
	if err != nil {
		s.logger.Errorf("Failed to get last sync time: %v", err)
	}

	status := "never"
	if lastSyncAt != nil {
		status = "success"
	}

	return &models.DirectorySyncStatus{
		IsSyncing:      isSyncing,
		LastSyncAt:     lastSyncAt,
		LastSyncStatus: status,
		UserCount:      userCount,
		GroupCount:     groupCount,
		ChannelCount:   channelCount,
	}, nil
}
