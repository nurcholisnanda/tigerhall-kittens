package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"

	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/bcrypt"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/errorhandler"
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

func generateSalt() (string, error) {
	salt := make([]byte, 16) // 16 bytes is a good size for a salt
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(salt), nil
}

func (s *userService) Register(ctx context.Context, input *model.NewUser) (interface{}, error) {
	if err := input.Validate(); err != nil { // Call the Validate method of your input struct
		return nil, errorhandler.NewInvalidInputError(err.Error())
	}

	// Generate a Random Salt
	salt, err := generateSalt()
	if err != nil {
		return nil, errorhandler.NewInternalServerError("failed to generate salt")
	}

	// Hash the Password with the Salt
	hashedPassword, err := s.bcrypt.HashPassword(input.Password + salt)
	if err != nil {
		return nil, errorhandler.NewInternalServerError("Failed to hash password")
	}

	// Create User with Salt and Hashed Password
	user := &model.User{
		ID:       uuid.NewString(),
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashedPassword), // Store the hash (as string for database)
		Salt:     salt,                   // Store the salt (as string for database)
	}

	// Save to Database
	if err := s.userRepo.CreateUser(ctx, user); err != nil {
		return nil, errorhandler.NewInternalServerError("Failed to create user")
	}

	return user, nil
}

func (s *userService) Login(ctx context.Context, email string, password string) (interface{}, error) {
	// Fetch User by Email
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { // Use errors.Is to check for specific error types
			return nil, errorhandler.NewNotFoundError("User not found")
		}
		return nil, errorhandler.NewInternalServerError("Failed to fetch user")
	}

	// Check Password with Stored Salt
	if err := s.bcrypt.ComparePassword(user.Password, password+user.Salt); err != nil {
		return nil, errorhandler.NewInvalidInputError("Invalid credentials")
	}

	// Generate JWT Tokens
	token, err := s.jwt.GenerateToken(ctx, user.ID)
	if err != nil {
		return nil, errorhandler.NewInternalServerError("Failed to generate access token")
	}

	res := map[string]string{
		"token": token,
	}

	return res, nil
}

func (s *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return s.userRepo.GetUserByID(ctx, id)
}
