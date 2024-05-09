package service

import (
	"context"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/domain/model"
)

type TigerService struct {
	// Add repository dependencies here if needed
}

func (s *TigerService) CreateTiger(ctx context.Context, tiger *model.Tiger) error {
	// Add implementation to create a tiger
	return nil
}

func (s *TigerService) GetTigers(ctx context.Context) ([]model.Tiger, error) {
	// Add implementation to list all tigers
	return nil, nil
}
