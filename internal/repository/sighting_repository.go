package repository

import (
	"context"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/helper"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/logger"
	"gorm.io/gorm"
)

type SightingRepositoryImpl struct {
	db *gorm.DB
}

func NewSightingRepositoryImpl(db *gorm.DB) SightingRepository {
	return &SightingRepositoryImpl{db: db}
}

func (r *SightingRepositoryImpl) GetSightingsByTigerID(ctx context.Context, tigerID string, limit int, offset int) ([]*model.Sighting, error) {
	var sightings []*model.Sighting
	if err := r.db.WithContext(ctx).Where("tiger_id = ?", tigerID).Order("last_seen_time desc").
		Offset(offset).Limit(limit).Find(&sightings).Error; err != nil {
		return nil, err
	}
	return sightings, nil
}

func (r *SightingRepositoryImpl) CreateSighting(ctx context.Context, sighting *model.Sighting) error {
	userId, err := helper.GetUserID(ctx)
	if err != nil {
		logger.Logger(ctx).Error("failed to get user id")
	}
	sighting.CreatedBy = userId
	return r.db.WithContext(ctx).Create(sighting).Error
}

func (r *SightingRepositoryImpl) GetLatestSightingByTigerID(ctx context.Context, tigerID string) (*model.Sighting, error) {
	var sighting *model.Sighting
	if err := r.db.WithContext(ctx).Where("tiger_id = ?", tigerID).Order("last_seen_time desc").Take(&sighting).Error; err != nil {
		return nil, err
	}
	return sighting, nil
}

func (r *SightingRepositoryImpl) ListUserCreatedSightingByTigerID(ctx context.Context, tigerID string) ([]string, error) {
	var createdBy []string

	if err := r.db.WithContext(ctx).Table("sightings").Where("tiger_id = ?", tigerID).
		Pluck("created_by", &createdBy).Error; err != nil {
		return nil, err
	}
	return createdBy, nil
}
