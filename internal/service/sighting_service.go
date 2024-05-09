package service

import (
	"context"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/domain/model"
)

type SightingService struct {
	// Add repository dependencies here if needed
}

func (s *SightingService) CreateSighting(ctx context.Context, sighting *model.Sighting) error {
	// Add implementation to create a sighting
	return nil
}

func (s *SightingService) ListSightingsByTigerID(ctx context.Context, tigerID uint) ([]model.Sighting, error) {
	// Add implementation to list all sightings of a specific tiger
	return nil, nil
}
