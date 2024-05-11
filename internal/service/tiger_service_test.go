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
	mockRepo "github.com/nurcholisnanda/tigerhall-kittens/internal/repository/mock"
	"go.uber.org/mock/gomock"
)

func TestNewTigerService(t *testing.T) {
	ctrl := gomock.NewController(t)
	tigerRepo := mockRepo.NewMockTigerRepository(ctrl)
	type args struct {
		tigerRepo repository.TigerRepository
	}
	tests := []struct {
		name string
		args args
		want TigerService
	}{
		{
			name: "success",
			args: args{
				tigerRepo: tigerRepo,
			},
			want: NewTigerService(tigerRepo),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewTigerService(tt.args.tigerRepo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewTigerService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_tigerService_CreateTiger(t *testing.T) {
	ctrl := gomock.NewController(t)
	tigerRepo := mockRepo.NewMockTigerRepository(ctrl)
	invalidCoordinate := model.LastSeenCoordinateInput{
		Latitude:  -25,
		Longitude: 255,
	}
	coordinate := model.LastSeenCoordinateInput{
		Latitude:  26,
		Longitude: 150,
	}
	input := model.TigerInput{
		Name:         "tiger a",
		DateOfBirth:  time.Now().AddDate(-14, -2, 0),
		LastSeenTime: time.Now().Add(time.Hour * -1),
	}
	invalidInput, validInput := input, input
	invalidInput.LastSeenCoordinate = &invalidCoordinate
	validInput.LastSeenCoordinate = &coordinate

	tiger := &model.Tiger{
		ID:                 uuid.NewString(),
		Name:               "tiger a",
		DateOfBirth:        input.DateOfBirth,
		LastSeenTime:       input.LastSeenTime,
		LastSeenCoordinate: (*model.LastSeenCoordinate)(&coordinate),
	}

	type fields struct {
		tigerRepo repository.TigerRepository
	}
	type args struct {
		ctx   context.Context
		input *model.TigerInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Tiger
		wantErr bool
		mocks   *gomock.Call
	}{
		{
			name: "should return error if inserting invalid coordinate",
			fields: fields{
				tigerRepo: tigerRepo,
			},
			args: args{
				ctx:   context.Background(),
				input: &invalidInput,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return error if inserting invalid date of birth",
			fields: fields{
				tigerRepo: tigerRepo,
			},
			args: args{
				ctx: context.Background(),
				input: &model.TigerInput{
					Name:               "tiger a",
					DateOfBirth:        time.Now().Add(time.Hour * 1),
					LastSeenTime:       time.Now().Add((time.Hour * -1)),
					LastSeenCoordinate: &coordinate,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return error if inserting invalid last seen time",
			fields: fields{
				tigerRepo: tigerRepo,
			},
			args: args{
				ctx: context.Background(),
				input: &model.TigerInput{
					Name:               "tiger a",
					DateOfBirth:        time.Now().Add(time.Hour * -11),
					LastSeenTime:       time.Now().Add((time.Hour * 1)),
					LastSeenCoordinate: &coordinate,
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "should return error if fail to create new tiger to the database",
			fields: fields{
				tigerRepo: tigerRepo,
			},
			args: args{
				ctx:   context.Background(),
				input: &validInput,
			},
			want:    nil,
			wantErr: true,
			mocks:   tigerRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("any error")),
		},
		{
			name: "should return error if input the same name",
			fields: fields{
				tigerRepo: tigerRepo,
			},
			args: args{
				ctx:   context.Background(),
				input: &validInput,
			},
			want:    nil,
			wantErr: true,
			mocks:   tigerRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(errors.New("unique constraint")),
		},
		{
			name: "success create new tiger",
			fields: fields{
				tigerRepo: tigerRepo,
			},
			args: args{
				ctx:   context.Background(),
				input: &validInput,
			},
			want:    tiger,
			wantErr: false,
			mocks:   tigerRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &tigerService{
				tigerRepo: tt.fields.tigerRepo,
			}
			got, err := s.CreateTiger(tt.args.ctx, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("tigerService.CreateTiger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.Name, tt.want.Name) {
					t.Errorf("tigerService.CreateTiger() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func Test_tigerService_ListTigers(t *testing.T) {
	ctrl := gomock.NewController(t)
	tigerRepo := mockRepo.NewMockTigerRepository(ctrl)
	tigers := []*model.Tiger{
		{
			ID:   uuid.NewString(),
			Name: "tiger a",
		},
	}
	type fields struct {
		tigerRepo repository.TigerRepository
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
		mocks   *gomock.Call
	}{
		{
			name: "should return error if failed to get data from db",
			fields: fields{
				tigerRepo: tigerRepo,
			},
			args: args{
				ctx:    context.Background(),
				limit:  1,
				offset: 0,
			},
			want:    nil,
			wantErr: true,
			mocks:   tigerRepo.EXPECT().ListTigers(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("any error")),
		},
		{
			name: "success",
			fields: fields{
				tigerRepo: tigerRepo,
			},
			args: args{
				ctx:    context.Background(),
				limit:  1,
				offset: 0,
			},
			want:    tigers,
			wantErr: false,
			mocks:   tigerRepo.EXPECT().ListTigers(gomock.Any(), gomock.Any(), gomock.Any()).Return(tigers, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &tigerService{
				tigerRepo: tt.fields.tigerRepo,
			}
			got, err := s.ListTigers(tt.args.ctx, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("tigerService.ListTigers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("tigerService.ListTigers() = %v, want %v", got, tt.want)
			}
		})
	}
}
