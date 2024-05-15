package repository

import (
	"context"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/contexthandler"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/logger"
	"gorm.io/gorm"
)

// tigerRepositoryImpl implements the TigerRepository interface
type tigerRepositoryImpl struct {
	db *gorm.DB
}

// NewTigerRepositoryImpl creates a new TigerRepository instance
func NewTigerRepositoryImpl(db *gorm.DB) TigerRepository {
	return &tigerRepositoryImpl{
		db: db,
	}
}

// GetTigerByID retrieves a tiger by its ID
func (r *tigerRepositoryImpl) GetTigerByID(ctx context.Context, id string) (*model.Tiger, error) {
	var tiger *model.Tiger
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&tiger).Error; err != nil {
		return nil, err
	}
	return tiger, nil
}

// Create saves a new tiger to the database
func (r *tigerRepositoryImpl) Create(ctx context.Context, tiger *model.Tiger) error {
	userId, err := contexthandler.GetUserID(ctx)
	if err != nil {
		logger.Logger(ctx).Error("failed to get user id")
	}
	tiger.CreatedBy = userId
	return r.db.WithContext(ctx).Create(tiger).Error
}

// ListTigers retrieves a paginated list of tigers
func (r *tigerRepositoryImpl) ListTigers(ctx context.Context, limit int, offset int) ([]*model.Tiger, error) {
	var tigers []*model.Tiger
	if err := r.db.WithContext(ctx).Offset(offset).Limit(limit).Order("last_seen_time desc").Find(&tigers).Error; err != nil {
		return nil, err
	}
	return tigers, nil
}
