// pkg/notification/notif.go
package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository/mock"
	"go.uber.org/mock/gomock"
)

func TestNewNotificationService(t *testing.T) {
	ctrl := gomock.NewController(t)
	sightingRepo := mock.NewMockSightingRepository(ctrl)
	userRepo := mock.NewMockUserRepository(ctrl)
	type args struct {
		sr repository.SightingRepository
		ur repository.UserRepository
	}
	tests := []struct {
		name string
		args args
		want *notificationService
	}{
		{
			name: "success",
			args: args{
				sr: sightingRepo,
				ur: userRepo,
			},
			want: NewNotificationService(sightingRepo, userRepo),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotificationService(tt.args.sr, tt.args.ur); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotificationService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notificationService_StartNotificationConsumer(t *testing.T) {
	type fields struct {
		sightingRepo repository.SightingRepository
		userRepo     repository.UserRepository
	}
	tests := []struct {
		name   string
		fields fields
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &notificationService{
				sightingRepo: tt.fields.sightingRepo,
				userRepo:     tt.fields.userRepo,
			}
			s.StartNotificationConsumer()
		})
	}
}

func Test_notificationService_FetchPreviousSighters(t *testing.T) {
	ctrl := gomock.NewController(t)
	sightingRepo := mock.NewMockSightingRepository(ctrl)
	userRepo := mock.NewMockUserRepository(ctrl)
	appendedID := uuid.NewString()
	type fields struct {
		sightingRepo repository.SightingRepository
		userRepo     repository.UserRepository
	}
	type args struct {
		ctx     context.Context
		tigerID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.User
		wantErr bool
		mocks   []*gomock.Call
	}{
		{
			name: "should return error when error fetching prev sighters",
			fields: fields{
				sightingRepo: sightingRepo,
				userRepo:     userRepo,
			},
			args: args{
				ctx:     context.Background(),
				tigerID: uuid.NewString(),
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				sightingRepo.EXPECT().ListUserCreatedSightingByTigerID(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "should return error when error fetching prev sighters",
			fields: fields{
				sightingRepo: sightingRepo,
				userRepo:     userRepo,
			},
			args: args{
				ctx:     context.Background(),
				tigerID: uuid.NewString(),
			},
			want:    nil,
			wantErr: false,
			mocks: []*gomock.Call{
				sightingRepo.EXPECT().ListUserCreatedSightingByTigerID(gomock.Any(), gomock.Any()).Return([]string{uuid.NewString()}, nil),
				userRepo.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "success",
			fields: fields{
				sightingRepo: sightingRepo,
				userRepo:     userRepo,
			},
			args: args{
				ctx:     context.Background(),
				tigerID: uuid.NewString(),
			},
			want:    []*model.User{{ID: appendedID}},
			wantErr: false,
			mocks: []*gomock.Call{
				sightingRepo.EXPECT().ListUserCreatedSightingByTigerID(gomock.Any(), gomock.Any()).Return([]string{uuid.NewString(), appendedID}, nil),
				userRepo.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
				userRepo.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&model.User{ID: appendedID}, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &notificationService{
				sightingRepo: tt.fields.sightingRepo,
				userRepo:     tt.fields.userRepo,
			}
			got, err := s.FetchPreviousSighters(tt.args.ctx, tt.args.tigerID)
			if (err != nil) != tt.wantErr {
				t.Errorf("notificationService.FetchPreviousSighters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("notificationService.FetchPreviousSighters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notificationService_sendEmail(t *testing.T) {
	type fields struct {
		sightingRepo repository.SightingRepository
		userRepo     repository.UserRepository
	}
	type args struct {
		ctx          context.Context
		user         model.User
		notification model.Notification
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "tested",
			args: args{
				ctx: context.Background(),
				user: model.User{
					Name:     "test",
					Email:    "test@example.com",
					Password: "123456",
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &notificationService{
				sightingRepo: tt.fields.sightingRepo,
				userRepo:     tt.fields.userRepo,
			}
			if err := s.sendEmail(tt.args.ctx, tt.args.user, tt.args.notification); (err != nil) != tt.wantErr {
				t.Errorf("notificationService.sendEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
