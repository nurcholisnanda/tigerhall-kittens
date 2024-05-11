// pkg/notification/notif.go
package service

import (
	"context"
	"fmt"
	"net/http"
	"net/smtp"
	"os"
	"time"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/helper"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/logger"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// Channel for sending notifications
var NotificationChan = make(chan model.Notification)

// NotificationService handles notification logic
type notificationService struct {
	sightingRepo repository.SightingRepository
	userRepo     repository.UserRepository
}

// NewNotificationService creates a new NotificationService
func NewNotificationService(sr repository.SightingRepository, ur repository.UserRepository) *notificationService {
	return &notificationService{
		sightingRepo: sr,
		userRepo:     ur,
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
			for _, sighter := range previousSighters {
				s.sendEmail(ctx, *sighter, notification)
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

// sendEmail (Placeholder - you need to implement this)
func (s *notificationService) sendEmail(ctx context.Context, user model.User, notification model.Notification) error {
	// Set up Mailtrap SMTP configuration
	from := mail.NewEmail("TigerHall Kittens", "allendragneel@gmail.com")
	to := mail.NewEmail(user.Name, user.Email)
	plainTextContent := fmt.Sprintf(
		"Hello %s,\n\nA new sighting of tiger %s has been reported on %s.\n\nBest regards,\nTigerHall Kittens",
		user.Name, notification.TigerID, notification.Timestamp.Format(time.RFC822),
	)

	// Set up authentication (replace with your Mailtrap credentials)
	auth := smtp.PlainAuth("", os.Getenv("MAILTRAP_USERNAME"), os.Getenv("MAILTRAP_PASSWORD"), "sandbox.smtp.mailtrap.io")

	// Send the email
	err := smtp.SendMail("sandbox.smtp.mailtrap.io:2525", auth, from.Address, []string{to.Address}, []byte(plainTextContent))
	if err != nil {
		return fmt.Errorf("error sending email: %w", err)
	}

	// Log successful email delivered
	logger.Logger(ctx).Info(ctx, "Sent email notification to user:", user.ID)

	return nil

}
