package repository

import (
	"context"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
)

//go:generate mockgen -source=repository.go -destination=mock/repository.go -package=mock
type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindUserByID(ctx context.Context, id string) (*model.User, error)
	FindUserByEmail(ctx context.Context, email string) (*model.User, error)
}

// type TigerRepository interface {
// 	Create(tiger *model.Tiger) error
// 	FindAll() ([]model.Tiger, error)
// }

// type SightingRepository interface {
// 	Create(tiger *model.Sighting) error
// 	FindAll() ([]model.Sighting, error)
// }
