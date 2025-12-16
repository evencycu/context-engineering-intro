package services

import (
	"context"
	"errors"
	"sync"

	"github.com/evencycu/TeamsNotifyGoV3/libs/models"
	"github.com/evencycu/TeamsNotifyGoV3/services/teamsnotification/repositories"
	"github.com/sirupsen/logrus"
)

// DirectoryService defines directory operations
type DirectoryService interface {
	SearchUsers(ctx context.Context, query string) ([]models.AzureADUser, error)
	ListGroups(ctx context.Context) ([]models.AzureADGroup, error)
	SyncDirectory(ctx context.Context) error
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
	users, _, err := s.graphService.GetUsers(ctx, "")
	if err != nil {
		s.logger.Errorf("Failed to fetch users from Graph: %v", err)
		// Don't return error immediately, try to sync groups
	} else {
		for _, u := range users {
			if err := s.repo.UpsertUser(ctx, &u); err != nil {
				s.logger.Errorf("Failed to upsert user %s: %v", u.AzureADID, err)
			}
		}
		s.logger.Infof("Synced %d users", len(users))
	}

	// Sync Groups
	groups, _, err := s.graphService.GetGroups(ctx, "")
	if err != nil {
		s.logger.Errorf("Failed to fetch groups from Graph: %v", err)
		return err
	}
	for _, g := range groups {
		if err := s.repo.UpsertGroup(ctx, &g); err != nil {
			s.logger.Errorf("Failed to upsert group %s: %v", g.AzureADID, err)
		}
	}
	s.logger.Infof("Synced %d groups", len(groups))

	return nil
}
