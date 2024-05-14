package service // Change to tigerservice

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/errorhandler"
)

type tigerService struct {
	tigerRepo repository.TigerRepository
}

// NewTigerService creates a new instance of TigerService
func NewTigerService(tigerRepo repository.TigerRepository) TigerService {
	return &tigerService{tigerRepo: tigerRepo}
}

// CreateTiger creates a new tiger in the database
func (s *tigerService) CreateTiger(ctx context.Context, input *model.TigerInput) (*model.Tiger, error) {
	// Input Validation
	if err := validateTigerInput(input); err != nil {
		return nil, err // Return the validation error directly
	}

	// Data Mapping
	tiger := &model.Tiger{
		ID:           uuid.NewString(),
		Name:         input.Name,
		DateOfBirth:  input.DateOfBirth,
		LastSeenTime: input.LastSeenTime,
		Coordinate:   (*model.Coordinate)(input.LastSeenCoordinate), // No need for casting anymore
	}

	// Database Interaction
	err := s.tigerRepo.Create(ctx, tiger)
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") { // Check if it's a unique constraint error
			return nil, errorhandler.NewCustomError(
				fmt.Sprintf("a tiger with the name %s already exists", input.Name),
				http.StatusConflict,
			)
		}
		return nil, errorhandler.NewInternalServerError("unexpected error creating tiger") // Wrap internal errorhandler
	}

	return tiger, nil
}

// validateTigerInput validates the input fields for creating a tiger
func validateTigerInput(input *model.TigerInput) error {
	if !isValidLatitude(input.LastSeenCoordinate.Latitude) ||
		!isValidLongitude(input.LastSeenCoordinate.Longitude) {
		return errorhandler.NewInvalidInputError("Invalid coordinates")
	}
	if input.DateOfBirth.After(time.Now()) {
		return errorhandler.NewInvalidInputError("Date of birth cannot be in the future")
	}
	if input.LastSeenTime.After(time.Now()) {
		return errorhandler.NewInvalidInputError("Last seen time cannot be in the future")
	}
	return nil
}

// ListTigers returns a list of tigers with pagination
func (s *tigerService) ListTigers(ctx context.Context, limit, offset int) ([]*model.Tiger, error) {
	tigers, err := s.tigerRepo.ListTigers(ctx, limit, offset)
	if err != nil {
		return nil, errorhandler.NewInternalServerError("unexpected error listing tigers")
	}
	return tigers, nil
}

// errorhandler Validation Functions
func isValidLatitude(latitude float64) bool {
	return math.Abs(latitude) <= 90
}

func isValidLongitude(longitude float64) bool {
	return math.Abs(longitude) <= 180
}
