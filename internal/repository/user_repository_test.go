package repository

import (
	"context"
	"reflect"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestNewUserRepoImpl(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name string
		args args
		want UserRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepoImpl(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepoImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepoImpl_CreateUser(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx  context.Context
		user *model.User
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
			r := &userRepoImpl{
				db: tt.fields.db,
			}
			if err := r.CreateUser(tt.args.ctx, tt.args.user); (err != nil) != tt.wantErr {
				t.Errorf("userRepoImpl.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepoImpl_GetUserByID(t *testing.T) {
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
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.GetUserByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepoImpl.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepoImpl.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepoImpl_GetUserByEmail(t *testing.T) {
	type fields struct {
		db *gorm.DB
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepoImpl{
				db: tt.fields.db,
			}
			got, err := r.GetUserByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("userRepoImpl.GetUserByEmail() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepoImpl.GetUserByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
