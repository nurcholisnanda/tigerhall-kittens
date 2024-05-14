package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/errorhandler"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/imagehandler"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/logger"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/storage"
	"gorm.io/gorm"
)

type sightingService struct {
	notifSvc     NotifService
	sightingRepo repository.SightingRepository
	tigerRepo    repository.TigerRepository
	s3Client     storage.S3Interface
	imgHandler   imagehandler.ImageHandler
}

func NewSightingService(
	notifSvc NotifService,
	sightingRepo repository.SightingRepository,
	tigerRepo repository.TigerRepository,
	s3Client storage.S3Interface,
	imgHandler imagehandler.ImageHandler,
) SightingService {
	return &sightingService{
		notifSvc:     notifSvc,
		sightingRepo: sightingRepo,
		tigerRepo:    tigerRepo,
		s3Client:     s3Client,
		imgHandler:   imgHandler,
	}
}

func (s *sightingService) ListSightings(ctx context.Context, tigerID string, limit int, offset int) ([]*model.Sighting, error) {
	sightings, err := s.sightingRepo.GetSightingsByTigerID(ctx, tigerID, limit, offset)
	if err != nil {
		// Handle errors (e.g., database error)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &errorhandler.NotFoundError{Message: "Tiger not found"}
		} else {
			logger.Logger(ctx).Error("Failed to list sightings for tiger:", err)
			return nil, errorhandler.NewCustomError("Failed to list sightings", http.StatusInternalServerError)
		}
	}
	return sightings, nil
}

func (s *sightingService) CreateSighting(ctx context.Context, input *model.SightingInput) (*model.Sighting, error) {
	// Check if Tiger exists
	tiger, err := s.tigerRepo.GetTigerByID(ctx, input.TigerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errorhandler.NewNotFoundError("Tiger not found")
		}
		logger.Logger(ctx).Error("Unexpected error getting tiger by ID: ", err)
		return nil, errorhandler.NewCustomError("Failed to retrieve tiger by ID", http.StatusInternalServerError)
	}

	if !isValidLatitude(input.Coordinate.Latitude) || !isValidLongitude(input.Coordinate.Longitude) {
		return nil, errorhandler.NewInvalidInputError("latitude must be between -90 and 90, longitude between -180 and 180")
	}

	sighting, err := s.sightingRepo.GetLatestSightingByTigerID(ctx, input.TigerID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Logger(ctx).Error("Unexpected error getting tiger by ID: ", err)
			return nil, errorhandler.NewCustomError("Failed to retrieve tiger by ID", http.StatusInternalServerError)
		}
	}

	var Coordinate *model.Coordinate
	if sighting == nil {
		Coordinate = tiger.Coordinate
	} else {
		Coordinate = sighting.Coordinate
	}

	distance := calculateDistance((*model.Coordinate)(input.Coordinate), Coordinate)
	if distance < 5000 {
		return nil, &errorhandler.SightingTooCloseError{
			Message: fmt.Sprintf("new sighting is too close to the last known location (%.2f meters)", distance),
		}
	}

	imagePath, err := s.imgHandler.ResizeAndUpload(ctx, input.Image)
	if err != nil {
		logger.Logger(ctx).Error(ctx, "Fail when resizing image", "error", err)
	}

	newSighting := &model.Sighting{
		ID:           uuid.NewString(),
		TigerID:      input.TigerID,
		LastSeenTime: time.Now(),
		Image:        &imagePath,
		Coordinate:   (*model.Coordinate)(input.Coordinate),
	}

	if err := s.sightingRepo.CreateSighting(ctx, newSighting); err != nil {
		logger.Logger(ctx).Error("Unexpected error creating sighting: ", err)
		return nil, errorhandler.NewCustomError("Failed to create sighting", http.StatusInternalServerError)
	}

	// Create a notification message\
	notification := model.Notification{
		TigerID:   newSighting.TigerID,
		Latitude:  input.Coordinate.Latitude,
		Longitude: input.Coordinate.Longitude,
		Timestamp: time.Now(),
	}

	// Send the notification
	s.notifSvc.SendNotification(notification)

	return newSighting, nil
}

func calculateDistance(coord1, coord2 *model.Coordinate) float64 {
	// Earth radius in meters
	const EarthRadius = 6371000
	// Convert latitude and longitude from degrees to radians
	lat1 := degreesToRadians(coord1.Latitude)
	lon1 := degreesToRadians(coord1.Longitude)
	lat2 := degreesToRadians(coord2.Latitude)
	lon2 := degreesToRadians(coord2.Longitude)

	// Haversine formula
	dlon := lon2 - lon1
	dlat := lat2 - lat1
	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	distance := EarthRadius * c
	return distance
}

// errorhandler function to convert degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
