package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/bcrypt"
	"gorm.io/gorm"
)

type userService struct {
	jwt      JWT
	userRepo repository.UserRepository
	bcrypt   bcrypt.BcryptInterface
}

func NewUserService(userRepo repository.UserRepository, bcrypt bcrypt.BcryptInterface, jwt JWT) UserService {
	return &userService{
		userRepo: userRepo,
		bcrypt:   bcrypt,
		jwt:      jwt,
	}
}

func (s *userService) Register(ctx context.Context, input model.NewUser) (interface{}, error) {
	//hashing password
	hashedPassword, err := s.bcrypt.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user, err := s.userRepo.FindUserByEmail(ctx, input.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if user != nil {
		return nil, errors.New("email already registered in our database")
	}

	user = &model.User{
		ID:       uuid.New().String(),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	res := map[string]string{"message": "register success"}

	return res, nil
}

func (s *userService) Login(ctx context.Context, email string, password string) (interface{}, error) {
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := s.bcrypt.ComparePassword(user.Password, password); err != nil {
		return nil, err
	}

	token, err := s.jwt.GenerateToken(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	res := map[string]string{
		"token": token,
	}

	return res, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.FindUserByID(ctx, id)
}
