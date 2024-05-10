package service

import (
	"context"
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/errors"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/logger"
)

type tigerService struct {
	tigerRepo repository.TigerRepository
}

func NewTigerService(tigerRepo repository.TigerRepository) TigerService {
	return &tigerService{
		tigerRepo: tigerRepo,
	}
}

// CreateTiger Function
func (s *tigerService) CreateTiger(ctx context.Context, input model.TigerInput) (*model.Tiger, error) {
	// Validations
	if !isValidLatitude(input.LastSeenCoordinate.Latitude) || !isValidLongitude(input.LastSeenCoordinate.Longitude) {
		return nil, &errors.InvalidCoordinatesError{
			Message: "latitude must be between -90 and 90, longitude between -180 and 180",
		}
	}

	if input.DateOfBirth.After(time.Now()) {
		return nil, &errors.InvalidDateOfBirthError{
			Message: "date of birth cannot be in the future",
		}
	}

	if input.LastSeenTime.After(time.Now()) {
		return nil, &errors.InvalidLastSeenTimeError{
			Message: "last seen time cannot be in the future",
		}
	}

	// Data Mapping
	tiger := &model.Tiger{
		ID:                 uuid.NewString(),
		Name:               input.Name,
		DateOfBirth:        input.DateOfBirth,
		LastSeenTime:       input.LastSeenTime,
		LastSeenCoordinate: (*model.LastSeenCoordinate)(input.LastSeenCoordinate),
	}

	// Database Interaction
	if err := s.tigerRepo.Create(ctx, tiger); err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			logger.Logger(ctx).Error("a tiger with the name already exists", err)
			return nil, &errors.TigerCreationError{
				Field:   "name",
				Message: fmt.Sprintf("a tiger with the name %s already exists", input.Name),
			}
		}
		logger.Logger(ctx).Error("Unexpected error creating tiger", err)
		return nil, errors.ErrInternalServer
	}

	return tiger, nil
}

// Helper Validation Functions
func isValidLatitude(latitude float64) bool {
	return math.Abs(latitude) <= 90
}

func isValidLongitude(longitude float64) bool {
	return math.Abs(longitude) <= 180
}
