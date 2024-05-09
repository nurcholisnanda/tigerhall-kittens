package repository

import (
	"context"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"gorm.io/gorm"
)

type userRepoImpl struct {
	db *gorm.DB
}

func NewUserRepoImpl(db *gorm.DB) UserRepository {
	return &userRepoImpl{
		db: db,
	}
}

func (r *userRepoImpl) CreateUser(ctx context.Context, user *model.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *userRepoImpl) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	var user *model.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userRepoImpl) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	var user *model.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
