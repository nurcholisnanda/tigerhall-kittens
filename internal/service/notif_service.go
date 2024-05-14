// pkg/notification/notif.go
package service

import (
	"context"
	"net/http"
	"time"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/errorhandler"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/logger"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/mailer"
)

// Channel for sending notifications
var notificationChan = make(chan model.Notification)

// NotificationService handles notification logic
type notificationService struct {
	sightingRepo repository.SightingRepository
	userRepo     repository.UserRepository
	mailSvc      mailer.MailService
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(
	sr repository.SightingRepository,
	ur repository.UserRepository,
	mailSvc mailer.MailService,
) NotifService {
	return &notificationService{
		sightingRepo: sr,
		userRepo:     ur,
		mailSvc:      mailSvc,
	}
}

func (s *notificationService) SendNotification(notif model.Notification) {
	notificationChan <- notif
}

// StartNotificationConsumer starts the background goroutine to consume notifications
func (s *notificationService) StartNotificationConsumer(ctx context.Context) {
	go func(parentCtx context.Context) {
		for notification := range notificationChan {
			func(notification model.Notification) {
				ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
				defer cancel()

				// Fetch previous sighters
				previousSighters, err := s.FetchPreviousSighters(ctx, notification.TigerID)
				if err != nil {
					logger.Logger(ctx).Error("Failed to fetch previous sighters:", err)
					return // Skip to the next notification
				}

				// Send email notifications
				for _, sighter := range previousSighters {
					notification.Sighter = sighter.Name
					s.mailSvc.Send(ctx, sighter.Email, "notif_mail.tmpl", notification)
				}
			}(notification)
		}
	}(ctx) // Pass the parent context to the goroutine
}

// fetchPreviousSighters (Placeholder - you need to implement this)
func (s *notificationService) FetchPreviousSighters(ctx context.Context, tigerID string) ([]*model.User, error) {
	userIDs, err := s.sightingRepo.ListUserCreatedSightingByTigerID(ctx, tigerID)
	if err != nil {
		logger.Logger(ctx).Error("Failed to fetch sightings by tiger ID: ", err)
		return nil, errorhandler.NewCustomError("Failed to fetch previous sighters", http.StatusInternalServerError)
	}

	var previousSighters []*model.User
	for _, userID := range userIDs {
		user, err := s.userRepo.GetUserByID(ctx, userID) // Replace with your user repository function
		if err != nil {
			logger.Logger(ctx).Error("Failed to fetch user by ID: ", err)
			continue // Skip this user if there's an error
		}
		previousSighters = append(previousSighters, user)
	}

	return previousSighters, nil
}

func (s *notificationService) CloseNotificationChannel() {
	close(notificationChan)
}
