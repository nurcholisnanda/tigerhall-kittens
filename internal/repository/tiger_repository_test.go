package repository

import (
	"context"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewTigerRepositoryImpl(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want TigerRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTigerRepositoryImpl(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTigerRepositoryImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tigerRepositoryImpl_GetTigerByID(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Tiger
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &tigerRepositoryImpl{
				db: tt.fields.db,
			}
			got, err := r.GetTigerByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("tigerRepositoryImpl.GetTigerByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tigerRepositoryImpl.GetTigerByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tigerRepositoryImpl_Create(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx   context.Context
		tiger *model.Tiger
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &tigerRepositoryImpl{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.tiger); (err != nil) != tt.wantErr {
				t.Errorf("tigerRepositoryImpl.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_tigerRepositoryImpl_ListTigers(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx    context.Context
		limit  int
		offset int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Tiger
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &tigerRepositoryImpl{
				db: tt.fields.db,
			}
			got, err := r.ListTigers(tt.args.ctx, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("tigerRepositoryImpl.ListTigers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tigerRepositoryImpl.ListTigers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTigerRepositoryImpl_Create(t *testing.T) {
	ctx := context.Background()
	gormDB, mock := newMockGormDB(t)
	mockTigerRepo := NewTigerRepositoryImpl(gormDB)
	now := time.Now()

	testTiger := &model.Tiger{
		ID:           uuid.NewString(),
		Name:         "Test Tiger",
		DateOfBirth:  now,
		LastSeenTime: now,
		Coordinate: &model.Coordinate{
			Latitude:  -8.195,
			Longitude: 120.821,
		},
	}

	// Expected SQL query and behavior
	// Note the specific column names in the INSERT query
	mock.ExpectBegin()
	insertQuery := `INSERT INTO "tigers" ` +
		`("id","name","date_of_birth","last_seen_time","latitude","longitude",` +
		`"created_at","created_by","updated_at","updated_by","deleted_at","deleted_by") ` +
		`VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	mock.ExpectExec(regexp.QuoteMeta(insertQuery)).
		WithArgs(
			testTiger.ID, testTiger.Name, testTiger.DateOfBirth, testTiger.LastSeenTime,
			testTiger.Coordinate.Latitude, testTiger.Coordinate.Longitude, sqlmock.AnyArg(),
			"", sqlmock.AnyArg(), "", sqlmock.AnyArg(), "",
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := mockTigerRepo.Create(ctx, testTiger)

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestTigerRepositoryImpl_ListTigers(t *testing.T) {
	ctx := context.Background()
	gormDB, mock := newMockGormDB(t)
	mockTigerRepo := NewTigerRepositoryImpl(gormDB)
	now := time.Now()
	query := `SELECT * FROM "tigers" WHERE "tigers"."deleted_at" IS NULL ORDER BY last_seen_time desc LIMIT $1`
	t.Run("success", func(t *testing.T) {

		// Expected SQL query and behavior
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(10). // Use question marks instead of '$1'
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "date_of_birth", "last_seen_time", "latitude", "longitude"}).
				AddRow("1", "Tiger One", now, now, 10.0, 20.0).
				AddRow("2", "Tiger Two", now, now.Add(-1*time.Hour), 15.0, 25.0))

		tigers, err := mockTigerRepo.ListTigers(ctx, 10, 0) // Remove offset to match the expected query
		assert.Nil(t, err)
		assert.Len(t, tigers, 2) // Check if result length is 2

		assert.Equal(t, "1", tigers[0].ID)
		assert.Equal(t, "Tiger One", tigers[0].Name)

		assert.Equal(t, "2", tigers[1].ID)
		assert.Equal(t, "Tiger Two", tigers[1].Name)

		// We make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("return not found", func(t *testing.T) {
		// Expected SQL query and behavior
		mock.ExpectQuery(regexp.QuoteMeta(query+` OFFSET $2`)).
			WithArgs(10, 1).WillReturnError(gorm.ErrRecordNotFound)

		tigers, err := mockTigerRepo.ListTigers(ctx, 10, 1)
		assert.NotNil(t, err)
		assert.Nil(t, tigers)

		// We make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestTigerRepositoryImpl_GetTigerByID(t *testing.T) {
	ctx := context.Background()
	gormDB, mock := newMockGormDB(t)
	mockTigerRepo := NewTigerRepositoryImpl(gormDB)
	query := `SELECT * FROM "tigers" WHERE id = $1 AND "tigers"."deleted_at" IS NULL ORDER BY "tigers"."id" LIMIT $2`
	now := time.Now()

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("1", 1). // Use question marks instead of '$1'
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "date_of_birth", "last_seen_time", "latitude", "longitude"}).
				AddRow("1", "Tiger One", now, now, 10.0, 20.0))

		tiger, err := mockTigerRepo.GetTigerByID(ctx, "1")
		assert.Nil(t, err)

		assert.Equal(t, "1", tiger.ID)
		assert.Equal(t, "Tiger One", tiger.Name)

		// We make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("return error not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("1", 1). // Use question marks instead of '$1'
			WillReturnError(gorm.ErrRecordNotFound)

		tiger, err := mockTigerRepo.GetTigerByID(ctx, "1")
		assert.Nil(t, tiger)
		assert.NotNil(t, err)

		// We make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
