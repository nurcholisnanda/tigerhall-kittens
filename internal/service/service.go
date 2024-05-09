package service

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
)

//go:generate mockgen -source=service.go -destination=mock/service.go -package=mock
type JWT interface {
	GenerateToken(ctx context.Context, userID string) (string, error)
	ValidateToken(ctx context.Context, requestToken string) (*jwt.Token, error)
}

type UserService interface {
	Register(ctx context.Context, input model.NewUser) (interface{}, error)
	Login(ctx context.Context, email string, password string) (interface{}, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
}
