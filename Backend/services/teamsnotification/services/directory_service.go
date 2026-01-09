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
	ListGroups(ctx context.Context) ([]models.AzureADGroup, error)
	ListGroupChannels(ctx context.Context, groupId string) ([]models.AzureADChannel, error)
	SyncDirectory(ctx context.Context) error
	// ChatGroup operations
	RegisterChatGroup(ctx context.Context, projectID uuid.UUID, name, chatID string) (*models.ChatGroup, error)
	ListChatGroups(ctx context.Context, projectID uuid.UUID) ([]models.ChatGroup, error)
	RemoveChatGroup(ctx context.Context, id uuid.UUID) error
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

func (s *directoryService) ListGroups(ctx context.Context) ([]models.AzureADGroup, error) {
	return s.repo.ListGroups(ctx)
}

func (s *directoryService) ListGroupChannels(ctx context.Context, groupId string) ([]models.AzureADChannel, error) {
	return s.graphService.GetGroupChannels(ctx, groupId)
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

func (s *directoryService) RemoveChatGroup(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteChatGroup(ctx, id)
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

	return nil
}
