// pkg/notification/notif.go
package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository/mock"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/mailer"
	mockMailer "github.com/nurcholisnanda/tigerhall-kittens/pkg/mailer/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewNotificationService(t *testing.T) {
	ctrl := gomock.NewController(t)
	sightingRepo := mock.NewMockSightingRepository(ctrl)
	userRepo := mock.NewMockUserRepository(ctrl)
	mockSvc := mockMailer.NewMockMailService(ctrl)
	type args struct {
		sr     repository.SightingRepository
		ur     repository.UserRepository
		mailer mailer.MailService
	}
	tests := []struct {
		name string
		args args
		want NotifService
	}{
		{
			name: "success",
			args: args{
				sr:     sightingRepo,
				ur:     userRepo,
				mailer: mockSvc,
			},
			want: NewNotificationService(sightingRepo, userRepo, mockSvc),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotificationService(tt.args.sr, tt.args.ur, tt.args.mailer); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotificationService() = %v, want %v", got, tt.want)
			}
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

func Test_notificationService_SendNotification(t *testing.T) {
	ctrl := gomock.NewController(t)
	notif := model.Notification{
		Sighter: "nanda",
		TigerID: "2",
	}
	type fields struct {
		sightingRepo repository.SightingRepository
		userRepo     repository.UserRepository
		mailSvc      mailer.MailService
	}
	type args struct {
		notif model.Notification
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "tested",
			fields: fields{
				sightingRepo: mock.NewMockSightingRepository(ctrl),
				userRepo:     mock.NewMockUserRepository(ctrl),
				mailSvc:      mockMailer.NewMockMailService(ctrl),
			},
			args: args{
				notif: notif,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &notificationService{
				sightingRepo: tt.fields.sightingRepo,
				userRepo:     tt.fields.userRepo,
				mailSvc:      tt.fields.mailSvc,
			}
			go func() { // Goroutine to receive from the channel
				time.Sleep(1 * time.Second)
				receivedNotif := <-notificationChan
				assert.Equal(t, tt.args.notif, receivedNotif) // Verify the correct notification was sent
			}()
			s.SendNotification(tt.args.notif)
		})
	}
}

func Test_notificationService_CloseNotificationChannel(t *testing.T) {
	ctrl := gomock.NewController(t)
	type fields struct {
		sightingRepo repository.SightingRepository
		userRepo     repository.UserRepository
		mailSvc      mailer.MailService
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "tested",
			fields: fields{
				sightingRepo: mock.NewMockSightingRepository(ctrl),
				userRepo:     mock.NewMockUserRepository(ctrl),
				mailSvc:      mockMailer.NewMockMailService(ctrl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &notificationService{
				sightingRepo: tt.fields.sightingRepo,
				userRepo:     tt.fields.userRepo,
				mailSvc:      tt.fields.mailSvc,
			}
			// Start a goroutine to attempt sending a message after a delay
			go func() {
				time.Sleep(100 * time.Millisecond) // Wait a short time to ensure channel is created
				defer func() {
					if r := recover(); r != nil {
						t.Errorf("Sending to closed channel should not panic, but got: %v", r)
					}
				}()
				notificationChan <- model.Notification{}
			}()

			// Close the channel
			s.CloseNotificationChannel()

			// Verify that the channel is closed
			select {
			case _, ok := <-notificationChan:
				if ok {
					t.Fatal("Notification channel should be closed after CloseNotificationChannel")
				}
			default:
				// Channel is closed, as expected
			}
		})
	}
}

func Test_notificationService_StartNotificationConsumer(t *testing.T) {
	ctrl := gomock.NewController(t)
	sightingRepo := mock.NewMockSightingRepository(ctrl)
	userRepo := mock.NewMockUserRepository(ctrl)
	mockSvc := mockMailer.NewMockMailService(ctrl)
	type fields struct {
		sightingRepo repository.SightingRepository
		userRepo     repository.UserRepository
		mailSvc      mailer.MailService
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		mocks  []*gomock.Call
	}{
		{
			name: "fail to fetch previous sighters",
			fields: fields{
				sightingRepo: sightingRepo,
				userRepo:     userRepo,
				mailSvc:      mockSvc,
			},
			args: args{
				ctx: context.Background(),
			},
			mocks: []*gomock.Call{
				sightingRepo.EXPECT().ListUserCreatedSightingByTigerID(gomock.Any(), gomock.Any()).
					Return(nil, errors.New("any error")),
			},
		},
		{
			name: "success",
			fields: fields{
				sightingRepo: sightingRepo,
				userRepo:     userRepo,
				mailSvc:      mockSvc,
			},
			args: args{
				ctx: context.Background(),
			},
			mocks: []*gomock.Call{
				sightingRepo.EXPECT().ListUserCreatedSightingByTigerID(gomock.Any(), gomock.Any()).
					Return([]string{"1"}, nil),
				userRepo.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).
					Return(&model.User{ID: "1"}, nil),
				mockSvc.EXPECT().Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &notificationService{
				sightingRepo: tt.fields.sightingRepo,
				userRepo:     tt.fields.userRepo,
				mailSvc:      tt.fields.mailSvc,
			}
			// Create a context with timeout for the test
			ctx, cancel := context.WithTimeout(tt.args.ctx, 5*time.Second)
			defer cancel()

			// Start the consumer
			notificationChan = make(chan model.Notification)
			defer s.CloseNotificationChannel()
			s.StartNotificationConsumer(ctx)

			// Send test notification
			notificationChan <- model.Notification{
				TigerID:   "1",
				Sighter:   "1",
				Timestamp: time.Now(),
			}
		})
	}
}
