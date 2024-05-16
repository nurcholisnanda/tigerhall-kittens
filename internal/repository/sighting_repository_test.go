package repository

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/contexthandler"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestSightingRepositoryImplGetSightingsByTigerID(t *testing.T) {
	gormDB, sqlMock := newMockGormDB(t)
	sightingRepo := NewSightingRepositoryImpl(gormDB)
	query := `SELECT * FROM "sightings" WHERE tiger_id = $1 AND "sightings"."deleted_at" IS NULL ` +
		`ORDER BY last_seen_time desc LIMIT $2`
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		sqlMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("tiger_1", 10).
			WillReturnRows(
				sqlmock.NewRows([]string{
					"id", "tiger_id", "last_seen_time", "image", "latitude", "longitude", "created_at", "created_by",
				}).
					AddRow("1", "tiger_1", now, "image_url", 10.0, 20.0, now, "user_1").
					AddRow("2", "tiger_1", now, "image_url", 15.0, 24.0, now, "user_2"),
			)
		users, err := sightingRepo.GetSightersByTigerID(context.Background(), "tiger_1", 10, 0)
		assert.Nil(t, err)
		assert.NotNil(t, users)

		assert.Equal(t, "1", users[0].ID)
		assert.Equal(t, "2", users[1].ID)

		assert.Equal(t, "tiger_1", users[0].TigerID)
		assert.Equal(t, "tiger_1", users[1].TigerID)

		// We make sure that all expectations were met
		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("return not found", func(t *testing.T) {
		sqlMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("tiger_2", 10).WillReturnError(gorm.ErrRecordNotFound)

		users, err := sightingRepo.GetSightersByTigerID(context.Background(), "tiger_2", 10, 0)
		assert.Nil(t, users)
		assert.NotNil(t, err)

		assert.Equal(t, gorm.ErrRecordNotFound, err)

		// We make sure that all expectations were met
		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestSightingRepositoryImplCreateSighting(t *testing.T) {
	ctx := context.Background()
	ctx = contexthandler.SetContext(ctx, "ContextKey", uuid.NewString())
	gormDB, mock := newMockGormDB(t)
	sightingRepo := NewSightingRepositoryImpl(gormDB)
	now := time.Now()
	imageUrl := "image_url"

	sighting := &model.Sighting{
		ID:           uuid.NewString(),
		TigerID:      uuid.NewString(),
		Image:        &imageUrl,
		LastSeenTime: now,
		Coordinate: &model.Coordinate{
			Latitude:  -8.195,
			Longitude: 120.821,
		},
	}

	// Expected SQL query and behavior
	// Note the specific column names in the INSERT query
	mock.ExpectBegin()
	insertQuery := `INSERT INTO "sightings" ` +
		`("id","tiger_id","last_seen_time","image","latitude","longitude",` +
		`"created_at","created_by","updated_at","updated_by","deleted_at","deleted_by") ` +
		`VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	mock.ExpectExec(regexp.QuoteMeta(insertQuery)).
		WithArgs(
			sighting.ID, sighting.TigerID, sighting.LastSeenTime, sighting.Image,
			sighting.Coordinate.Latitude, sighting.Coordinate.Longitude, sqlmock.AnyArg(),
			sighting.CreatedBy, sqlmock.AnyArg(), "", sqlmock.AnyArg(), "",
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := sightingRepo.CreateSighting(ctx, sighting)

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestSightingRepositoryImplGetLatestSightingByTigerID(t *testing.T) {
	gormDB, sqlMock := newMockGormDB(t)
	sightingRepo := NewSightingRepositoryImpl(gormDB)
	query := `SELECT * FROM "sightings" WHERE tiger_id = $1 AND "sightings"."deleted_at" IS NULL ` +
		`ORDER BY last_seen_time desc LIMIT $2`
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		sqlMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("tiger_1", 1).
			WillReturnRows(
				sqlmock.NewRows([]string{
					"id", "tiger_id", "last_seen_time", "image", "latitude", "longitude", "created_at", "created_by",
				}).
					AddRow("1", "tiger_1", now, "image_url", 10.0, 20.0, now, "user_1"),
			)
		user, err := sightingRepo.GetLatestSightingByTigerID(context.Background(), "tiger_1")
		assert.Nil(t, err)
		assert.NotNil(t, user)

		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "tiger_1", user.TigerID)

		// We make sure that all expectations were met
		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("return not found", func(t *testing.T) {
		sqlMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("tiger_2", 1).WillReturnError(gorm.ErrRecordNotFound)

		users, err := sightingRepo.GetLatestSightingByTigerID(context.Background(), "tiger_2")
		assert.Nil(t, users)
		assert.NotNil(t, err)

		assert.Equal(t, gorm.ErrRecordNotFound, err)

		// We make sure that all expectations were met
		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestSightingRepositoryImplListUserCreatedSightingByTigerID(t *testing.T) {
	gormDB, sqlMock := newMockGormDB(t)
	sightingRepo := NewSightingRepositoryImpl(gormDB)
	query := `SELECT DISTINCT created_by FROM "sightings" WHERE tiger_id = $1`

	t.Run("success", func(t *testing.T) {
		sqlMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("tiger_1").
			WillReturnRows(
				sqlmock.NewRows([]string{"created_by"}).AddRow("user_1").AddRow("user_2"),
			)
		userIDs, err := sightingRepo.ListUserCreatedSightingByTigerID(context.Background(), "tiger_1")
		assert.Nil(t, err)
		assert.NotNil(t, userIDs)

		assert.Equal(t, "user_1", userIDs[0])
		assert.Equal(t, "user_2", userIDs[1])

		// We make sure that all expectations were met
		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("return not found", func(t *testing.T) {
		sqlMock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("tiger_2").WillReturnError(gorm.ErrRecordNotFound)

		users, err := sightingRepo.ListUserCreatedSightingByTigerID(context.Background(), "tiger_2")
		assert.Nil(t, users)
		assert.NotNil(t, err)

		assert.Equal(t, gorm.ErrRecordNotFound, err)

		// We make sure that all expectations were met
		if err := sqlMock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
