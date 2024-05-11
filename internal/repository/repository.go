package repository

import (
	"context"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
)

//go:generate mockgen -source=repository.go -destination=mock/repository.go -package=mock
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

//go:generate mockgen -source=repository.go -destination=mock/repository.go -package=mock
type TigerRepository interface {
	Create(ctx context.Context, tiger *model.Tiger) error
	GetTigerByID(ctx context.Context, id string) (*model.Tiger, error)
	ListTigers(ctx context.Context, limit int, offset int) ([]*model.Tiger, error)
}

//go:generate mockgen -source=repository.go -destination=mock/repository.go -package=mock
type SightingRepository interface {
	GetSightingsByTigerID(ctx context.Context, tigerID string, limit int, offset int) ([]*model.Sighting, error)
	CreateSighting(ctx context.Context, sighting *model.Sighting) error
	GetLatestSightingByTigerID(ctx context.Context, tigerID string) (*model.Sighting, error)
	ListUserCreatedSightingByTigerID(ctx context.Context, tigerID string) ([]string, error)
}
