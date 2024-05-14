package service

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	mockRepo "github.com/nurcholisnanda/tigerhall-kittens/internal/repository/mock"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/service/mock"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/imagehandler"
	mockImage "github.com/nurcholisnanda/tigerhall-kittens/pkg/imagehandler/mock"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/storage"
	mockStorage "github.com/nurcholisnanda/tigerhall-kittens/pkg/storage/mock"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestNewSightingService(t *testing.T) {
	ctrl := gomock.NewController(t)
	sightingRepo := mockRepo.NewMockSightingRepository(ctrl)
	tigerRepo := mockRepo.NewMockTigerRepository(ctrl)
	notifSvc := mock.NewMockNotifService(ctrl)
	s3Client := mockStorage.NewMockS3Interface(ctrl)
	imgHandler := mockImage.NewMockImageHandler(ctrl)
	type args struct {
		notifSvc     NotifService
		sightingRepo repository.SightingRepository
		tigerRepo    repository.TigerRepository
		s3Client     storage.S3Interface
		imgHandler   imagehandler.ImageHandler
	}
	tests := []struct {
		name string
		args args
		want SightingService
	}{
		{
			name: "success",
			args: args{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
				notifSvc:     notifSvc,
				s3Client:     s3Client,
				imgHandler:   imgHandler,
			},
			want: NewSightingService(notifSvc, sightingRepo, tigerRepo, s3Client, imgHandler),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSightingService(tt.args.notifSvc, tt.args.sightingRepo, tt.args.tigerRepo, tt.args.s3Client, tt.args.imgHandler); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSightingService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sightingService_CreateSighting(t *testing.T) {
	ctrl := gomock.NewController(t)
	sightingRepo := mockRepo.NewMockSightingRepository(ctrl)
	tigerRepo := mockRepo.NewMockTigerRepository(ctrl)
	imgHandler := mockImage.NewMockImageHandler(ctrl)
	notifSvc := mock.NewMockNotifService(ctrl)
	type fields struct {
		sightingRepo repository.SightingRepository
		tigerRepo    repository.TigerRepository
		imgHandler   imagehandler.ImageHandler
		notifSvc     NotifService
	}
	imageFile := &graphql.Upload{}
	type args struct {
		ctx   context.Context
		input *model.SightingInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Sighting
		wantErr bool
		mocks   []*gomock.Call
	}{
		{
			name: "should return tiger not found error if tiger is not exist",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
			},
			args: args{
				ctx: context.Background(),
				input: &model.SightingInput{
					TigerID: uuid.NewString(),
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound),
			},
		},
		{
			name: "should return unexpected error if probem getting from tiger repo",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
			},
			args: args{
				ctx: context.Background(),
				input: &model.SightingInput{
					TigerID: uuid.NewString(),
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "should return invalid input if inputting invalid coordinator",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
			},
			args: args{
				ctx: context.Background(),
				input: &model.SightingInput{
					TigerID: uuid.NewString(),
					Coordinate: &model.CoordinateInput{
						Latitude:  25,
						Longitude: 200,
					},
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(&model.Tiger{}, nil),
			},
		},
		{
			name: "should return unexpected error if failed to get slighting data",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
			},
			args: args{
				ctx: context.Background(),
				input: &model.SightingInput{
					TigerID: uuid.NewString(),
					Coordinate: &model.CoordinateInput{
						Latitude:  25,
						Longitude: 130,
					},
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(&model.Tiger{}, nil),
				sightingRepo.EXPECT().GetLatestSightingByTigerID(gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
			},
		},
		{
			name: "should return conflict if last coordinate is too close",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
			},
			args: args{
				ctx: context.Background(),
				input: &model.SightingInput{
					TigerID: uuid.NewString(),
					Coordinate: &model.CoordinateInput{
						Latitude:  25,
						Longitude: 130,
					},
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(&model.Tiger{
					Coordinate: &model.Coordinate{Latitude: 25, Longitude: 130},
				}, nil),
				sightingRepo.EXPECT().GetLatestSightingByTigerID(gomock.Any(), gomock.Any()).Return(&model.Sighting{
					Coordinate: &model.Coordinate{Latitude: 25, Longitude: 130},
				}, nil),
			},
		},
		{
			name: "should return error if fail on resizing and upload the image and error create new sighting",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
				imgHandler:   imgHandler,
			},
			args: args{
				ctx: context.Background(),
				input: &model.SightingInput{
					TigerID: uuid.NewString(),
					Coordinate: &model.CoordinateInput{
						Latitude:  70,
						Longitude: -140,
					},
					Image: &graphql.Upload{},
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(&model.Tiger{
					ID:         uuid.NewString(),
					Coordinate: &model.Coordinate{Latitude: 25, Longitude: 130},
				}, nil),
				sightingRepo.EXPECT().GetLatestSightingByTigerID(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound),
				imgHandler.EXPECT().ResizeAndUpload(gomock.Any(), imageFile).Return("", errors.New("any error")),
				sightingRepo.EXPECT().CreateSighting(gomock.Any(), gomock.Any()).Return(errors.New("any error")),
			},
		},
		{
			name: "should return error if fail on creating New Sighting",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
				imgHandler:   imgHandler,
				notifSvc:     notifSvc,
			},
			args: args{
				ctx: context.Background(),
				input: &model.SightingInput{
					TigerID: uuid.NewString(),
					Coordinate: &model.CoordinateInput{
						Latitude:  70,
						Longitude: -140,
					},
					Image: &graphql.Upload{},
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(&model.Tiger{
					Coordinate: &model.Coordinate{Latitude: 25, Longitude: 130},
				}, nil),
				sightingRepo.EXPECT().GetLatestSightingByTigerID(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound),
				imgHandler.EXPECT().ResizeAndUpload(gomock.Any(), imageFile).Return("img url", nil),
				sightingRepo.EXPECT().CreateSighting(gomock.Any(), gomock.Any()).Return(errors.New("any error")),
			},
		},
		{
			name: "success",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
				imgHandler:   imgHandler,
				notifSvc:     notifSvc,
			},
			args: args{
				ctx: context.Background(),
				input: &model.SightingInput{
					TigerID: uuid.NewString(),
					Coordinate: &model.CoordinateInput{
						Latitude:  70,
						Longitude: -140,
					},
					Image: &graphql.Upload{},
				},
			},
			want: &model.Sighting{
				Coordinate: &model.Coordinate{
					Latitude:  70,
					Longitude: -140,
				},
			},
			wantErr: false,
			mocks: []*gomock.Call{
				tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(&model.Tiger{
					Coordinate: &model.Coordinate{Latitude: 25, Longitude: 130},
				}, nil),
				sightingRepo.EXPECT().GetLatestSightingByTigerID(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound),
				imgHandler.EXPECT().ResizeAndUpload(gomock.Any(), imageFile).Return("img url", nil),
				sightingRepo.EXPECT().CreateSighting(gomock.Any(), gomock.Any()).Return(nil),
				notifSvc.EXPECT().SendNotification(gomock.Any()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sightingService{
				sightingRepo: tt.fields.sightingRepo,
				tigerRepo:    tt.fields.tigerRepo,
				imgHandler:   tt.fields.imgHandler,
				notifSvc:     tt.fields.notifSvc,
			}
			got, err := s.CreateSighting(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("sightingService.CreateSighting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.Coordinate, tt.want.Coordinate) {
					t.Errorf("sightingService.CreateSighting() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_calculateDistance(t *testing.T) {
	type args struct {
		coord1 *model.Coordinate
		coord2 *model.Coordinate
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateDistance(tt.args.coord1, tt.args.coord2); got != tt.want {
				t.Errorf("calculateDistance() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_degreesToRadians(t *testing.T) {
	type args struct {
		degrees float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := degreesToRadians(tt.args.degrees); got != tt.want {
				t.Errorf("degreesToRadians() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sightingService_ListSightings(t *testing.T) {
	ctrl := gomock.NewController(t)
	sightingRepo := mockRepo.NewMockSightingRepository(ctrl)
	tigerRepo := mockRepo.NewMockTigerRepository(ctrl)
	sightings := []*model.Sighting{{ID: uuid.NewString(), TigerID: uuid.NewString()}}
	type fields struct {
		sightingRepo repository.SightingRepository
		tigerRepo    repository.TigerRepository
	}
	type args struct {
		ctx     context.Context
		tigerID string
		limit   int
		offset  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*model.Sighting
		wantErr bool
		mocks   []*gomock.Call
	}{
		{
			name: "should return error if database returning record not found",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
			},
			args: args{
				ctx:     context.Background(),
				tigerID: uuid.NewString(),
				limit:   10,
				offset:  0,
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				sightingRepo.EXPECT().GetSightingsByTigerID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, gorm.ErrRecordNotFound),
			},
		},
		{
			name: "should return error if database returning error when get sighting data",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
			},
			args: args{
				ctx:     context.Background(),
				tigerID: uuid.NewString(),
				limit:   10,
				offset:  0,
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				sightingRepo.EXPECT().GetSightingsByTigerID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil, errors.New("any error")),
			},
		},
		{
			name: "success",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
			},
			args: args{
				ctx:     context.Background(),
				tigerID: uuid.NewString(),
				limit:   10,
				offset:  0,
			},
			want:    sightings,
			wantErr: false,
			mocks: []*gomock.Call{
				sightingRepo.EXPECT().GetSightingsByTigerID(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(sightings, nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sightingService{
				sightingRepo: tt.fields.sightingRepo,
				tigerRepo:    tt.fields.tigerRepo,
			}
			got, err := s.ListSightings(tt.args.ctx, tt.args.tigerID, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("sightingService.ListSightings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("sightingService.ListSightings() = %v, want %v", got, tt.want)
			}
		})
	}
}
