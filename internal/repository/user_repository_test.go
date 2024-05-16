package repository

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestUserRepoImpl_CreateUser(t *testing.T) {
	ctx := context.Background()
	gormDB, mock := newMockGormDB(t)
	mockUserRepo := NewUserRepoImpl(gormDB)

	user := &model.User{
		ID:       uuid.NewString(),
		Name:     "Test User",
		Email:    "test@mail.com",
		Password: "123456",
		Salt:     "randomsalt",
	}

	// Expected SQL query and behavior
	// Note the specific column names in the INSERT query
	mock.ExpectBegin()
	insertQuery := `INSERT INTO "users" ` +
		`("id","name","email","password","salt") ` +
		`VALUES ($1,$2,$3,$4,$5)`
	mock.ExpectExec(regexp.QuoteMeta(insertQuery)).
		WithArgs(
			user.ID, user.Name, user.Email, user.Password, user.Salt,
		).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	err := mockUserRepo.CreateUser(ctx, user)

	assert.Nil(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserRepoImpl_GetUserByID(t *testing.T) {
	ctx := context.Background()
	gormDB, mock := newMockGormDB(t)
	mockUserRepo := NewUserRepoImpl(gormDB)
	query := `SELECT * FROM "users" WHERE id = $1 ORDER BY "users"."id" LIMIT $2`

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("1", 1). // Use question marks instead of '$1'
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "salt"}).
				AddRow("1", "User One", "user@mail.com", "123456", "randomsalt"))

		user, err := mockUserRepo.GetUserByID(ctx, "1")
		assert.Nil(t, err)

		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "User One", user.Name)

		// We make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("return error not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("1", 1). // Use question marks instead of '$1'
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := mockUserRepo.GetUserByID(ctx, "1")
		assert.Nil(t, user)
		assert.NotNil(t, err)

		// We make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

func TestUserRepoImpl_GetUserByEmail(t *testing.T) {
	ctx := context.Background()
	gormDB, mock := newMockGormDB(t)
	mockUserRepo := NewUserRepoImpl(gormDB)
	query := `SELECT * FROM "users" WHERE email = $1 ORDER BY "users"."id" LIMIT $2`

	t.Run("success", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("user@mail.com", 1). // Use question marks instead of '$1'
			WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "salt"}).
				AddRow("1", "User One", "user@mail.com", "123456", "randomsalt"))

		user, err := mockUserRepo.GetUserByEmail(ctx, "user@mail.com")
		assert.Nil(t, err)

		assert.Equal(t, "1", user.ID)
		assert.Equal(t, "user@mail.com", user.Email)

		// We make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
	t.Run("return error not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs("user@mail.com", 1). // Use question marks instead of '$1'
			WillReturnError(gorm.ErrRecordNotFound)

		user, err := mockUserRepo.GetUserByEmail(ctx, "user@mail.com")
		assert.Nil(t, user)
		assert.NotNil(t, err)

		// We make sure that all expectations were met
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
