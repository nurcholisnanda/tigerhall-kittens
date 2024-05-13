package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"github.com/nurcholisnanda/tigerhall-kittens/config"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/helper"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/logger"
	"gorm.io/gorm"
)

type sightingService struct {
	sightingRepo repository.SightingRepository
	tigerRepo    repository.TigerRepository
	s3Client     *config.S3Client
}

func NewSightingService(
	sightingRepo repository.SightingRepository,
	tigerRepo repository.TigerRepository,
	s3Client *config.S3Client,
) SightingService {
	return &sightingService{
		sightingRepo: sightingRepo,
		tigerRepo:    tigerRepo,
		s3Client:     s3Client,
	}
}

func (s *sightingService) ListSightings(ctx context.Context, tigerID string, limit int, offset int) ([]*model.Sighting, error) {
	sightings, err := s.sightingRepo.GetSightingsByTigerID(ctx, tigerID, limit, offset)
	if err != nil {
		// Handle errors (e.g., database error)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &helper.TigerNotFound{Message: "Tiger not found"}
		} else {
			logger.Logger(ctx).Error("Failed to list sightings for tiger:", err)
			return nil, helper.NewCustomError("Failed to list sightings", http.StatusInternalServerError)
		}
	}
	return sightings, nil
}

func (s *sightingService) CreateSighting(ctx context.Context, input *model.SightingInput) (*model.Sighting, error) {
	// Check if Tiger exists
	tiger, err := s.tigerRepo.GetTigerByID(ctx, input.TigerID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, &helper.TigerNotFound{Message: "Tiger not found"}
		}
		logger.Logger(ctx).Error("Unexpected error getting tiger by ID: ", err)
		return nil, helper.NewCustomError("Failed to retrieve tiger by ID", http.StatusInternalServerError)
	}

	if !isValidLatitude(input.Coordinate.Latitude) || !isValidLongitude(input.Coordinate.Longitude) {
		return nil, &helper.InvalidCoordinatesError{
			Message: "latitude must be between -90 and 90, longitude between -180 and 180",
		}
	}

	sighting, err := s.sightingRepo.GetLatestSightingByTigerID(ctx, input.TigerID)
	if err != nil {
		if err != gorm.ErrRecordNotFound {
			logger.Logger(ctx).Error("Unexpected error getting tiger by ID: ", err)
			return nil, helper.NewCustomError("Failed to retrieve tiger by ID", http.StatusInternalServerError)
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
		return nil, &helper.SightingTooCloseError{
			Message: fmt.Sprintf("new sighting is too close to the last known location (%.2f meters)", distance),
		}
	}

	imagePath, err := s.GetResizedImage(ctx, input.Image)
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
		return nil, helper.NewCustomError("Failed to create sighting", http.StatusInternalServerError)
	}

	// Create a notification message\
	notification := model.Notification{
		TigerID:   newSighting.TigerID,
		Latitude:  input.Coordinate.Latitude,
		Longitude: input.Coordinate.Longitude,
		Timestamp: time.Now(),
	}

	// Send the notification
	NotificationChan <- notification

	return newSighting, nil
}

func (s *sightingService) GetResizedImage(ctx context.Context, inputImage *graphql.Upload) (string, error) {
	imageData, readErr := io.ReadAll(inputImage.File)
	if readErr != nil {
		logger.Logger(ctx).Error(ctx, "readErr", readErr)

		fmt.Printf("error from file %v", readErr)
	}

	img, format, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		log.Printf("Error decoding image")

		return "", fmt.Errorf("error decoding image: %v got format %v", err, format)
	}
	resizedImage := resize.Resize(250, 200, img, resize.Lanczos3)

	// Encode resized image to base64
	buf := new(bytes.Buffer)
	err = jpeg.Encode(buf, resizedImage, nil)
	if err != nil {
		logger.Logger(ctx).Error(ctx, "Error encoding image", "error", err)
		return "", fmt.Errorf("error encoding image: %v", err)
	}

	// Upload to R2
	bucketName := os.Getenv("R2_BUCKET_NAME") // Get bucket name from env vars

	objectName := fmt.Sprintf("%s.%s", uuid.NewString(), format)
	_, err = s.s3Client.Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
		Body:   bytes.NewReader(buf.Bytes()),
	})
	if err != nil {
		logger.Logger(ctx).Error("Error uploading image to R2", "error", err)
		return "", fmt.Errorf("error uploading image to R2: %w", err)
	}

	objectURL := fmt.Sprintf("https://%s.r2.cloudflarestorage.com/%s/%s", os.Getenv("R2_ACCOUNT_ID"), bucketName, objectName)
	logger.Logger(ctx).Info("Image uploaded to R2", "objectURL", objectURL)
	return objectURL, nil
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

// Helper function to convert degrees to radians
func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}
