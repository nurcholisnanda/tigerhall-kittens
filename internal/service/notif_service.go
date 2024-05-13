// pkg/notification/notif.go
package service

import (
	"context"
	"net/http"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/helper"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/logger"
)

// Channel for sending notifications
var NotificationChan = make(chan model.Notification)

// NotificationService handles notification logic
type notificationService struct {
	sightingRepo repository.SightingRepository
	userRepo     repository.UserRepository
	mailSvc      MailerInterface
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(
	sr repository.SightingRepository,
	ur repository.UserRepository,
	mailSvc MailerInterface,
) *notificationService {
	return &notificationService{
		sightingRepo: sr,
		userRepo:     ur,
		mailSvc:      mailSvc,
	}
}

// StartNotificationConsumer starts the background goroutine to consume notifications
func (s *notificationService) StartNotificationConsumer() {
	ctx := context.Background()
	go func() {
		for notification := range NotificationChan {
			// Fetch users who reported sightings of the same tiger
			previousSighters, err := s.FetchPreviousSighters(ctx, notification.TigerID)
			if err != nil {
				logger.Logger(ctx).Error("Error fetching previous sighters:", err)
				continue // Skip this notification if there's an error
			}

			// Send email notifications
			done := make(chan struct{})
			for _, sighter := range previousSighters {
				s.mailSvc.Send(ctx, sighter.Email, "notif_mail.tmpl", notification, done)
			}
		}
	}()
}

// fetchPreviousSighters (Placeholder - you need to implement this)
func (s *notificationService) FetchPreviousSighters(ctx context.Context, tigerID string) ([]*model.User, error) {
	userIDs, err := s.sightingRepo.ListUserCreatedSightingByTigerID(ctx, tigerID)
	if err != nil {
		logger.Logger(ctx).Error("Failed to fetch sightings by tiger ID: ", err)
		return nil, helper.NewCustomError("Failed to fetch previous sighters", http.StatusInternalServerError)
	}

	// Deduplicate user IDs (in case a user reported multiple sightings)
	uniqueUserIDs := make(map[string]bool)
	for _, userID := range userIDs {
		uniqueUserIDs[userID] = true
	}

	var previousSighters []*model.User
	for userID := range uniqueUserIDs {
		user, err := s.userRepo.GetUserByID(ctx, userID) // Replace with your user repository function
		if err != nil {
			logger.Logger(ctx).Error("Failed to fetch user by ID: ", err)
			continue // Skip this user if there's an error
		}
		previousSighters = append(previousSighters, user)
	}

	return previousSighters, nil
}
