package service

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
)

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
type JWT interface {
	GenerateToken(ctx context.Context, userID string) (string, error)
	ValidateToken(ctx context.Context, requestToken string) (*jwt.Token, error)
}

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
type UserService interface {
	Register(ctx context.Context, input *model.NewUser) (interface{}, error)
	Login(ctx context.Context, email string, password string) (interface{}, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
type TigerService interface {
	CreateTiger(ctx context.Context, input *model.TigerInput) (*model.Tiger, error)
	ListTigers(ctx context.Context, limit int, offset int) ([]*model.Tiger, error)
}

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
type SightingService interface {
	CreateSighting(ctx context.Context, newSighting *model.SightingInput) (*model.Sighting, error)
	ListSightings(ctx context.Context, tigerID string, limit int, offset int) ([]*model.Sighting, error)
	GetResizedImage(ctx context.Context, inputImage *graphql.Upload) (string, error)
}

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
type MailerInterface interface {
	Send(ctx context.Context, recipient, templateFile string, data interface{}, done chan struct{}) error
}

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
type NotifService interface {
	SendNotification(notif model.Notification)
}
