package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	mockRepo "github.com/nurcholisnanda/tigerhall-kittens/internal/repository/mock"
	mockSvc "github.com/nurcholisnanda/tigerhall-kittens/internal/service/mock"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/bcrypt"
	mockBcrypt "github.com/nurcholisnanda/tigerhall-kittens/pkg/bcrypt/mock"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestNewUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mockRepo.NewMockUserRepository(ctrl)
	mockBcrypt := mockBcrypt.NewMockBcryptInterface(ctrl)
	jwt := mockSvc.NewMockJWT(ctrl)

	type args struct {
		userRepo repository.UserRepository
		bcrypt   bcrypt.BcryptInterface
		jwt      JWT
	}
	tests := []struct {
		name string
		args args
		want UserService
	}{
		{
			name: "success",
			args: args{
				userRepo: userRepo,
				bcrypt:   mockBcrypt,
				jwt:      jwt,
			},
			want: NewUserService(userRepo, mockBcrypt, jwt),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.userRepo, tt.args.bcrypt, tt.args.jwt); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mockRepo.NewMockUserRepository(ctrl)
	mockBcrypt := mockBcrypt.NewMockBcryptInterface(ctrl)
	jwt := mockSvc.NewMockJWT(ctrl)

	ctx := context.Background()
	input := &model.NewUser{
		Name:     "name",
		Email:    "email@mail.co.id",
		Password: "password",
	}
	type fields struct {
		userRepo repository.UserRepository
		bcrypt   bcrypt.BcryptInterface
	}
	type args struct {
		ctx   context.Context
		input *model.NewUser
	}
	hashedPassword := []byte("hashed_password")
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		mocks   []*gomock.Call
	}{
		{
			name: "should return error if fail to validate input",
			fields: fields{
				userRepo: userRepo,
				bcrypt:   mockBcrypt,
			},
			args: args{
				ctx:   ctx,
				input: &model.NewUser{},
			},
			wantErr: true,
		},
		{
			name: "should return error if fail to hash password",
			fields: fields{
				userRepo: userRepo,
				bcrypt:   mockBcrypt,
			},
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
			mocks: []*gomock.Call{
				userRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(nil, gorm.ErrRecordNotFound),
				mockBcrypt.EXPECT().HashPassword(gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "should return error if email already exist in database",
			fields: fields{
				userRepo: userRepo,
				bcrypt:   mockBcrypt,
			},
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
			mocks: []*gomock.Call{
				userRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(&model.User{ID: "1"}, nil),
			},
		},
		{
			name: "should return error if fail to createUserCreateUser a new user",
			fields: fields{
				userRepo: userRepo,
				bcrypt:   mockBcrypt,
			},
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: true,
			mocks: []*gomock.Call{
				userRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(nil, gorm.ErrRecordNotFound),
				mockBcrypt.EXPECT().HashPassword(gomock.Any()).Return(hashedPassword, nil),
				userRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(errors.New("any error")),
			},
		},
		{
			name: "success register new user",
			fields: fields{
				userRepo: userRepo,
				bcrypt:   mockBcrypt,
			},
			args: args{
				ctx:   ctx,
				input: input,
			},
			wantErr: false,
			mocks: []*gomock.Call{
				userRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(nil, gorm.ErrRecordNotFound),
				mockBcrypt.EXPECT().HashPassword(gomock.Any()).Return(hashedPassword, nil),
				userRepo.EXPECT().CreateUser(gomock.Any(), gomock.Any()).Return(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				userRepo: tt.fields.userRepo,
				bcrypt:   tt.fields.bcrypt,
				jwt:      jwt,
			}
			_, err := s.Register(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func Test_userService_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mockRepo.NewMockUserRepository(ctrl)
	mockBcrypt := mockBcrypt.NewMockBcryptInterface(ctrl)
	jwt := mockSvc.NewMockJWT(ctrl)

	user := &model.User{
		ID:       uuid.New().String(),
		Name:     "name",
		Email:    "email",
		Password: "password",
	}
	ctx := context.Background()

	type fields struct {
		userRepo repository.UserRepository
		bcrypt   bcrypt.BcryptInterface
	}
	type args struct {
		ctx      context.Context
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
		mocks   []*gomock.Call
	}{
		{
			name: "should return internal error if getting error from database",
			fields: fields{
				userRepo: userRepo,
				bcrypt:   mockBcrypt,
			},
			args: args{
				ctx:      ctx,
				email:    "email",
				password: "password",
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				userRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "should return error if password incorrect",
			fields: fields{
				userRepo: userRepo,
				bcrypt:   mockBcrypt,
			},
			args: args{
				ctx:      ctx,
				email:    "email",
				password: "password",
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				userRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(user, nil),
				mockBcrypt.EXPECT().ComparePassword(gomock.Any(), gomock.Any()).Return(errors.New("any error")),
			},
		},
		{
			name: "should return error if fail to generate token",
			fields: fields{
				userRepo: userRepo,
				bcrypt:   mockBcrypt,
			},
			args: args{
				ctx:      ctx,
				email:    "email",
				password: "password",
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				userRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(user, nil),
				mockBcrypt.EXPECT().ComparePassword(gomock.Any(), gomock.Any()).Return(nil),
				jwt.EXPECT().GenerateToken(ctx, gomock.Any()).Return("", errors.New("any error")),
			},
		},
		{
			name: "success",
			fields: fields{
				userRepo: userRepo,
				bcrypt:   mockBcrypt,
			},
			args: args{
				ctx:      ctx,
				email:    "email",
				password: "password",
			},
			want:    map[string]string{"token": "token"},
			wantErr: false,
			mocks: []*gomock.Call{
				userRepo.EXPECT().GetUserByEmail(ctx, gomock.Any()).Return(user, nil),
				mockBcrypt.EXPECT().ComparePassword(gomock.Any(), gomock.Any()).Return(nil),
				jwt.EXPECT().GenerateToken(ctx, gomock.Any()).Return("token", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				userRepo: tt.fields.userRepo,
				bcrypt:   mockBcrypt,
				jwt:      jwt,
			}
			got, err := s.Login(tt.args.ctx, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_GetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	userRepo := mockRepo.NewMockUserRepository(ctrl)
	type fields struct {
		jwt      JWT
		userRepo repository.UserRepository
		bcrypt   bcrypt.BcryptInterface
	}
	id := uuid.NewString()
	user := &model.User{
		ID: id,
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
		mock    *gomock.Call
	}{
		{
			name: "success",
			fields: fields{
				userRepo: userRepo,
			},
			args: args{
				ctx: context.Background(),
				id:  uuid.NewString(),
			},
			want:    user,
			wantErr: false,
			mock:    userRepo.EXPECT().GetUserByID(context.Background(), gomock.Any()).Return(user, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				jwt:      tt.fields.jwt,
				userRepo: tt.fields.userRepo,
				bcrypt:   tt.fields.bcrypt,
			}
			got, err := s.GetUserByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("userService.GetUserByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.GetUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}
