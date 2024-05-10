package repository

import (
	"context"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/helper"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/logger"
	"gorm.io/gorm"
)

type TigerRepositoryImpl struct {
	db *gorm.DB
}

func NewTigerRepositoryImpl(db *gorm.DB) TigerRepository {
	return &TigerRepositoryImpl{
		db: db,
	}
}

func (r *TigerRepositoryImpl) Create(ctx context.Context, tiger *model.Tiger) error {
	userId, err := helper.GetUserID(ctx)
	if err != nil {
		logger.Logger(ctx).Error("failed to get user id")
	}
	tiger.CreatedBy = userId
	return r.db.WithContext(ctx).Create(tiger).Error
}

// func (r *TigerRepositoryImpl) FindAll() ([]model.Tiger, error) {
// 	// Add implementation to retrieve all tigers from the database
// 	return nil, nil
// }
