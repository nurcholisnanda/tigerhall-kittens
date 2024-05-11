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
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/helper"
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
func (s *tigerService) CreateTiger(ctx context.Context, input *model.TigerInput) (*model.Tiger, error) {
	// Validations
	if !isValidLatitude(input.LastSeenCoordinate.Latitude) || !isValidLongitude(input.LastSeenCoordinate.Longitude) {
		return nil, &helper.InvalidCoordinatesError{
			Message: "latitude must be between -90 and 90, longitude between -180 and 180",
		}
	}

	if input.DateOfBirth.After(time.Now()) {
		return nil, &helper.InvalidDateOfBirthError{
			Message: "date of birth cannot be in the future",
		}
	}

	if input.LastSeenTime.After(time.Now()) {
		return nil, &helper.InvalidLastSeenTimeError{
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
			return nil, &helper.TigerCreationError{
				Field:   "name",
				Message: fmt.Sprintf("a tiger with the name %s already exists", input.Name),
			}
		}
		logger.Logger(ctx).Error("unexpected error creating tiger", err)
		return nil, helper.ErrInternalServer
	}

	return tiger, nil
}

func (s *tigerService) ListTigers(ctx context.Context, limit int, offset int) ([]*model.Tiger, error) {
	tigers, err := s.tigerRepo.ListTigers(ctx, limit, offset)
	if err != nil {
		logger.Logger(ctx).Error("unexpected error creating tiger", err)
		return nil, fmt.Errorf("unexpected error listing tigers : %v", err.Error())
	}
	return tigers, nil
}

// Helper Validation Functions
func isValidLatitude(latitude float64) bool {
	return math.Abs(latitude) <= 90
}

func isValidLongitude(longitude float64) bool {
	return math.Abs(longitude) <= 180
}
