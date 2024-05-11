package service

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	mockRepo "github.com/nurcholisnanda/tigerhall-kittens/internal/repository/mock"
	"go.uber.org/mock/gomock"
	"gorm.io/gorm"
)

func TestNewSightingService(t *testing.T) {
	ctrl := gomock.NewController(t)
	sightingRepo := mockRepo.NewMockSightingRepository(ctrl)
	tigerRepo := mockRepo.NewMockTigerRepository(ctrl)
	type args struct {
		sightingRepo repository.SightingRepository
		tigerRepo    repository.TigerRepository
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
			},
			want: NewSightingService(sightingRepo, tigerRepo),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSightingService(tt.args.sightingRepo, tt.args.tigerRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSightingService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_sightingService_CreateSighting(t *testing.T) {
	ctrl := gomock.NewController(t)
	sightingRepo := mockRepo.NewMockSightingRepository(ctrl)
	tigerRepo := mockRepo.NewMockTigerRepository(ctrl)
	type fields struct {
		sightingRepo repository.SightingRepository
		tigerRepo    repository.TigerRepository
	}
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
					LastSeenCoordinate: &model.LastSeenCoordinateInput{
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
					LastSeenCoordinate: &model.LastSeenCoordinateInput{
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
			name: "should invalid input if last seen time is in the future",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
			},
			args: args{
				ctx: context.Background(),
				input: &model.SightingInput{
					TigerID: uuid.NewString(),
					LastSeenCoordinate: &model.LastSeenCoordinateInput{
						Latitude:  25,
						Longitude: 130,
					},
					LastSeenTime: time.Now().Add(1 * time.Hour),
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(&model.Tiger{}, nil),
				sightingRepo.EXPECT().GetLatestSightingByTigerID(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound),
			},
		},
		{
			name: "should invalid input if last seen time is before recorded last seen",
			fields: fields{
				sightingRepo: sightingRepo,
				tigerRepo:    tigerRepo,
			},
			args: args{
				ctx: context.Background(),
				input: &model.SightingInput{
					TigerID: uuid.NewString(),
					LastSeenCoordinate: &model.LastSeenCoordinateInput{
						Latitude:  25,
						Longitude: 130,
					},
					LastSeenTime: time.Now().Add(-5 * time.Hour),
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(&model.Tiger{}, nil),
				sightingRepo.EXPECT().GetLatestSightingByTigerID(gomock.Any(), gomock.Any()).Return(&model.Sighting{LastSeenTime: time.Now().Add(-3 * time.Hour)}, nil),
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
					LastSeenCoordinate: &model.LastSeenCoordinateInput{
						Latitude:  25,
						Longitude: 130,
					},
					LastSeenTime: time.Now().Add(-5 * time.Hour),
				},
			},
			want:    nil,
			wantErr: true,
			mocks: []*gomock.Call{
				tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(&model.Tiger{
					LastSeenCoordinate: &model.LastSeenCoordinate{Latitude: 25, Longitude: 130},
				}, nil),
				sightingRepo.EXPECT().GetLatestSightingByTigerID(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound),
			},
		},
		// {
		// 	name: "should return error if fail on creating New Sighting",
		// 	fields: fields{
		// 		sightingRepo: sightingRepo,
		// 		tigerRepo:    tigerRepo,
		// 	},
		// 	args: args{
		// 		ctx: context.Background(),
		// 		input: &model.SightingInput{
		// 			TigerID: uuid.NewString(),
		// 			LastSeenCoordinate: &model.LastSeenCoordinateInput{
		// 				Latitude:  70,
		// 				Longitude: -140,
		// 			},
		// 			LastSeenTime: time.Now().Add(-5 * time.Hour),
		// 		},
		// 	},
		// 	want:    nil,
		// 	wantErr: true,
		// 	mocks: []*gomock.Call{
		// 		tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(&model.Tiger{
		// 			LastSeenCoordinate: &model.LastSeenCoordinate{Latitude: 25, Longitude: 130},
		// 		}, nil),
		// 		sightingRepo.EXPECT().GetLatestSightingByTigerID(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound),
		// 		sightingRepo.EXPECT().CreateSighting(gomock.Any(), gomock.Any()).Return(errors.New("any error")),
		// 	},
		// },
		// {
		// 	name: "success",
		// 	fields: fields{
		// 		sightingRepo: sightingRepo,
		// 		tigerRepo:    tigerRepo,
		// 	},
		// 	args: args{
		// 		ctx: context.Background(),
		// 		input: &model.SightingInput{
		// 			TigerID: uuid.NewString(),
		// 			LastSeenCoordinate: &model.LastSeenCoordinateInput{
		// 				Latitude:  70,
		// 				Longitude: -140,
		// 			},
		// 			LastSeenTime: time.Now().Add(-5 * time.Hour),
		// 		},
		// 	},
		// 	want: &model.Sighting{
		// 		LastSeenCoordinate: &model.LastSeenCoordinate{
		// 			Latitude:  70,
		// 			Longitude: -140,
		// 		},
		// 	},
		// 	wantErr: false,
		// 	mocks: []*gomock.Call{
		// 		tigerRepo.EXPECT().GetTigerByID(gomock.Any(), gomock.Any()).Return(&model.Tiger{
		// 			LastSeenCoordinate: &model.LastSeenCoordinate{Latitude: 25, Longitude: 130},
		// 		}, nil),
		// 		sightingRepo.EXPECT().GetLatestSightingByTigerID(gomock.Any(), gomock.Any()).Return(nil, gorm.ErrRecordNotFound),
		// 		sightingRepo.EXPECT().CreateSighting(gomock.Any(), gomock.Any()).Return(nil),
		// 	},
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sightingService{
				sightingRepo: tt.fields.sightingRepo,
				tigerRepo:    tt.fields.tigerRepo,
			}
			got, err := s.CreateSighting(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("sightingService.CreateSighting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.LastSeenCoordinate, tt.want.LastSeenCoordinate) {
					t.Errorf("sightingService.CreateSighting() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_sightingService_GetResizedImage(t *testing.T) {
	type fields struct {
		sightingRepo repository.SightingRepository
		tigerRepo    repository.TigerRepository
	}
	type args struct {
		ctx       context.Context
		imageData *graphql.Upload
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &sightingService{
				sightingRepo: tt.fields.sightingRepo,
				tigerRepo:    tt.fields.tigerRepo,
			}
			got, err := s.GetResizedImage(tt.args.ctx, tt.args.imageData)
			if (err != nil) != tt.wantErr {
				t.Errorf("sightingService.GetResizedImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("sightingService.GetResizedImage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_calculateDistance(t *testing.T) {
	type args struct {
		coord1 *model.LastSeenCoordinate
		coord2 *model.LastSeenCoordinate
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
