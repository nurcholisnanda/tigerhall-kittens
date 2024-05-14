package repository

import (
	"context"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/contexthandler"
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

func (r *TigerRepositoryImpl) GetTigerByID(ctx context.Context, id string) (*model.Tiger, error) {
	var tiger *model.Tiger
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&tiger).Error; err != nil {
		return nil, err
	}
	return tiger, nil
}

func (r *TigerRepositoryImpl) Create(ctx context.Context, tiger *model.Tiger) error {
	userId, err := contexthandler.GetUserID(ctx)
	if err != nil {
		logger.Logger(ctx).Error("failed to get user id")
	}
	tiger.CreatedBy = userId
	return r.db.WithContext(ctx).Create(tiger).Error
}

func (r *TigerRepositoryImpl) ListTigers(ctx context.Context, limit int, offset int) ([]*model.Tiger, error) {
	var tigers []*model.Tiger
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("last_seen_time desc").Find(&tigers).Error; err != nil {
		return nil, err
	}
	return tigers, nil
}
