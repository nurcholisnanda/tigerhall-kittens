package graph

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/generated"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/service"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/service/mock"
	"go.uber.org/mock/gomock"
)

func Test_authOpsResolver_Login(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := mock.NewMockUserService(ctrl)
	userID := uuid.NewString()
	type fields struct {
		Resolver *Resolver
	}
	type args struct {
		ctx      context.Context
		obj      *model.AuthOps
		email    string
		password string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
		mocks   *gomock.Call
	}{
		{
			name: "should return error if failed to register",
			fields: fields{
				Resolver: &Resolver{
					UserSvc: userSvc,
				},
			},
			args: args{
				ctx:      context.Background(),
				email:    "test@mail.com",
				password: "123456",
			},
			want:    nil,
			wantErr: true,
			mocks: userSvc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.New("any error")),
		},
		{
			name: "success",
			fields: fields{
				Resolver: &Resolver{
					UserSvc: userSvc,
				},
			},
			args: args{
				ctx:      context.Background(),
				email:    "test@mail.com",
				password: "123456",
			},
			want:    &model.User{ID: userID, Name: "Name A"},
			wantErr: false,
			mocks: userSvc.EXPECT().Login(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(&model.User{ID: userID, Name: "Name A"}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authOpsResolver{
				Resolver: tt.fields.Resolver,
			}
			got, err := r.Login(tt.args.ctx, tt.args.obj, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("authOpsResolver.Login() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authOpsResolver.Login() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authOpsResolver_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := mock.NewMockUserService(ctrl)
	userID := uuid.NewString()
	type fields struct {
		Resolver *Resolver
	}
	type args struct {
		ctx   context.Context
		obj   *model.AuthOps
		input model.NewUser
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    interface{}
		wantErr bool
		mocks   *gomock.Call
	}{
		{
			name: "should return error if failed to register",
			fields: fields{
				Resolver: &Resolver{
					UserSvc: userSvc,
				},
			},
			args: args{
				ctx:   context.Background(),
				input: model.NewUser{Name: "Name A"},
			},
			want:    nil,
			wantErr: true,
			mocks: userSvc.EXPECT().Register(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("any error")),
		},
		{
			name: "success",
			fields: fields{
				Resolver: &Resolver{
					UserSvc: userSvc,
				},
			},
			args: args{
				ctx:   context.Background(),
				input: model.NewUser{Name: "Name A"},
			},
			want:    &model.User{ID: userID, Name: "Name A"},
			wantErr: false,
			mocks: userSvc.EXPECT().Register(gomock.Any(), gomock.Any()).
				Return(&model.User{ID: userID, Name: "Name A"}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authOpsResolver{
				Resolver: tt.fields.Resolver,
			}
			got, err := r.Register(tt.args.ctx, tt.args.obj, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("authOpsResolver.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authOpsResolver.Register() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createOpsResolver_CreateSighting(t *testing.T) {
	ctrl := gomock.NewController(t)
	sightSvc := mock.NewMockSightingService(ctrl)
	tigerID := uuid.NewString()
	type fields struct {
		Resolver *Resolver
	}
	type args struct {
		ctx   context.Context
		obj   *model.CreateOps
		input model.SightingInput
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.Sighting
		wantErr bool
		mocks   *gomock.Call
	}{
		{
			name: "should return error if error create sighting",
			fields: fields{
				Resolver: &Resolver{
					SightingSvc: sightSvc,
				},
			},
			args: args{
				ctx: context.Background(),
				input: model.SightingInput{
					TigerID: tigerID,
				},
			},
			want:    nil,
			wantErr: true,
			mocks: sightSvc.EXPECT().CreateSighting(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("any error")),
		},
		{
			name: "success",
			fields: fields{
				Resolver: &Resolver{
					SightingSvc: sightSvc,
				},
			},
			args: args{
				ctx: context.Background(),
				input: model.SightingInput{
					TigerID: tigerID,
				},
			},
			want:    &model.Sighting{TigerID: tigerID},
			wantErr: false,
			mocks: sightSvc.EXPECT().CreateSighting(gomock.Any(), gomock.Any()).
				Return(&model.Sighting{TigerID: tigerID}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &createOpsResolver{
				Resolver: tt.fields.Resolver,
			}
			got, err := r.CreateSighting(tt.args.ctx, tt.args.obj, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("createOpsResolver.CreateSighting() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createOpsResolver.CreateSighting() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createOpsResolver_CreateTiger(t *testing.T) {
	ctrl := gomock.NewController(t)
	tigerSvc := mock.NewMockTigerService(ctrl)
	tigerID := uuid.NewString()
	now := time.Now()
	input := &model.TigerInput{
		Name:         "Tiger A",
		DateOfBirth:  now,
		LastSeenTime: now,
		LastSeenCoordinate: &model.CoordinateInput{
			Latitude:  10,
			Longitude: 50,
		},
	}
	type fields struct {
		Resolver *Resolver
	}
	type args struct {
		ctx   context.Context
		obj   *model.CreateOps
		input model.TigerInput
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
			name: "should return error if error create tiger",
			fields: fields{
				Resolver: &Resolver{
					TigerSvc: tigerSvc,
				},
			},
			args: args{
				ctx:   context.Background(),
				input: *input,
			},
			want:    nil,
			wantErr: true,
			mocks: tigerSvc.EXPECT().CreateTiger(gomock.Any(), gomock.Any()).
				Return(nil, errors.New("any error")),
		},
		{
			name: "success",
			fields: fields{
				Resolver: &Resolver{
					TigerSvc: tigerSvc,
				},
			},
			args: args{
				ctx:   context.Background(),
				input: *input,
			},
			want:    &model.Tiger{ID: tigerID},
			wantErr: false,
			mocks: tigerSvc.EXPECT().CreateTiger(gomock.Any(), gomock.Any()).
				Return(&model.Tiger{ID: tigerID}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &createOpsResolver{
				Resolver: tt.fields.Resolver,
			}
			got, err := r.CreateTiger(tt.args.ctx, tt.args.obj, tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("createOpsResolver.CreateTiger() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createOpsResolver.CreateTiger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listOpsResolver_ListTigers(t *testing.T) {
	ctrl := gomock.NewController(t)
	tigerSvc := mock.NewMockTigerService(ctrl)
	tigerID := uuid.NewString()
	type fields struct {
		Resolver *Resolver
	}
	type args struct {
		ctx    context.Context
		obj    *model.ListOps
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
			name: "should return error if error get tiger list",
			fields: fields{
				Resolver: &Resolver{
					TigerSvc: tigerSvc,
				},
			},
			args: args{
				ctx:    context.Background(),
				limit:  10,
				offset: 0,
			},
			want:    nil,
			wantErr: true,
			mocks: tigerSvc.EXPECT().ListTigers(gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.New("any error")),
		},
		{
			name: "success",
			fields: fields{
				Resolver: &Resolver{
					TigerSvc: tigerSvc,
				},
			},
			args: args{
				ctx:    context.Background(),
				limit:  10,
				offset: 0,
			},
			want:    []*model.Tiger{{ID: tigerID}},
			wantErr: false,
			mocks: tigerSvc.EXPECT().ListTigers(gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]*model.Tiger{{ID: tigerID}}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &listOpsResolver{
				Resolver: tt.fields.Resolver,
			}
			got, err := r.ListTigers(tt.args.ctx, tt.args.obj, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("listOpsResolver.ListTigers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listOpsResolver.ListTigers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_listOpsResolver_ListSightings(t *testing.T) {
	ctrl := gomock.NewController(t)
	sightSvc := mock.NewMockSightingService(ctrl)
	tigerID := uuid.NewString()
	type fields struct {
		Resolver *Resolver
	}
	type args struct {
		ctx     context.Context
		obj     *model.ListOps
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
		mocks   *gomock.Call
	}{
		{
			name: "should return error if error get sighting list",
			fields: fields{
				Resolver: &Resolver{
					SightingSvc: sightSvc,
				},
			},
			args: args{
				ctx:     context.Background(),
				tigerID: tigerID,
				limit:   10,
				offset:  0,
			},
			want:    nil,
			wantErr: true,
			mocks: sightSvc.EXPECT().ListSightings(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil, errors.New("any error")),
		},
		{
			name: "success",
			fields: fields{
				Resolver: &Resolver{
					SightingSvc: sightSvc,
				},
			},
			args: args{
				ctx:     context.Background(),
				tigerID: tigerID,
				limit:   10,
				offset:  0,
			},
			want: []*model.Sighting{
				{
					TigerID: tigerID,
				},
			},
			wantErr: false,
			mocks: sightSvc.EXPECT().ListSightings(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return([]*model.Sighting{{TigerID: tigerID}}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &listOpsResolver{
				Resolver: tt.fields.Resolver,
			}
			got, err := r.ListSightings(tt.args.ctx, tt.args.obj, tt.args.tigerID, tt.args.limit, tt.args.offset)
			if (err != nil) != tt.wantErr {
				t.Errorf("listOpsResolver.ListSightings() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("listOpsResolver.ListSightings() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mutationResolver_Auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := mock.NewMockUserService(ctrl)
	tigerSvc := mock.NewMockTigerService(ctrl)
	sightSvc := mock.NewMockSightingService(ctrl)
	type fields struct {
		Resolver *Resolver
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.AuthOps
		wantErr bool
	}{
		{
			name: "implemented",
			fields: fields{
				Resolver: &Resolver{
					UserSvc:     userSvc,
					TigerSvc:    tigerSvc,
					SightingSvc: sightSvc,
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want:    &model.AuthOps{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mutationResolver{
				Resolver: tt.fields.Resolver,
			}
			got, err := r.Auth(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("mutationResolver.Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mutationResolver.Auth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mutationResolver_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := mock.NewMockUserService(ctrl)
	tigerSvc := mock.NewMockTigerService(ctrl)
	sightSvc := mock.NewMockSightingService(ctrl)
	type fields struct {
		Resolver *Resolver
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.CreateOps
		wantErr bool
	}{
		{
			name: "implemented",
			fields: fields{
				Resolver: &Resolver{
					UserSvc:     userSvc,
					TigerSvc:    tigerSvc,
					SightingSvc: sightSvc,
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want:    &model.CreateOps{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &mutationResolver{
				Resolver: tt.fields.Resolver,
			}
			got, err := r.Create(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("mutationResolver.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mutationResolver.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queryResolver_User(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := mock.NewMockUserService(ctrl)
	tigerSvc := mock.NewMockTigerService(ctrl)
	sightSvc := mock.NewMockSightingService(ctrl)
	id := uuid.NewString()
	type fields struct {
		Resolver *Resolver
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
		mocks   *gomock.Call
	}{
		{
			name: "success get user",
			fields: fields{
				Resolver: &Resolver{
					UserSvc:     userSvc,
					TigerSvc:    tigerSvc,
					SightingSvc: sightSvc,
				},
			},
			args: args{
				ctx: context.Background(),
				id:  id,
			},
			want: &model.User{
				ID: id,
			},
			wantErr: false,
			mocks:   userSvc.EXPECT().GetUserByID(gomock.Any(), gomock.Any()).Return(&model.User{ID: id}, nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &queryResolver{
				Resolver: tt.fields.Resolver,
			}
			got, err := r.User(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("queryResolver.User() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryResolver.User() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_queryResolver_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := mock.NewMockUserService(ctrl)
	tigerSvc := mock.NewMockTigerService(ctrl)
	sightSvc := mock.NewMockSightingService(ctrl)
	type fields struct {
		Resolver *Resolver
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *model.ListOps
		wantErr bool
	}{
		{
			name: "implemented",
			fields: fields{
				Resolver: &Resolver{
					UserSvc:     userSvc,
					TigerSvc:    tigerSvc,
					SightingSvc: sightSvc,
				},
			},
			args: args{
				ctx: context.Background(),
			},
			want:    &model.ListOps{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &queryResolver{
				Resolver: tt.fields.Resolver,
			}
			got, err := r.List(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("queryResolver.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryResolver.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_AuthOps(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := mock.NewMockUserService(ctrl)
	tigerSvc := mock.NewMockTigerService(ctrl)
	sightSvc := mock.NewMockSightingService(ctrl)
	type fields struct {
		UserSvc     service.UserService
		TigerSvc    service.TigerService
		SightingSvc service.SightingService
	}
	tests := []struct {
		name   string
		fields fields
		want   generated.AuthOpsResolver
	}{
		{
			name: "implemented",
			fields: fields{
				UserSvc:     userSvc,
				TigerSvc:    tigerSvc,
				SightingSvc: sightSvc,
			},
			want: &authOpsResolver{
				Resolver: &Resolver{
					UserSvc:     userSvc,
					TigerSvc:    tigerSvc,
					SightingSvc: sightSvc,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resolver{
				UserSvc:     tt.fields.UserSvc,
				TigerSvc:    tt.fields.TigerSvc,
				SightingSvc: tt.fields.SightingSvc,
			}
			if got := r.AuthOps(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolver.AuthOps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_CreateOps(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := mock.NewMockUserService(ctrl)
	tigerSvc := mock.NewMockTigerService(ctrl)
	sightSvc := mock.NewMockSightingService(ctrl)
	type fields struct {
		UserSvc     service.UserService
		TigerSvc    service.TigerService
		SightingSvc service.SightingService
	}
	tests := []struct {
		name   string
		fields fields
		want   generated.CreateOpsResolver
	}{
		{
			name: "implemented",
			fields: fields{
				UserSvc:     userSvc,
				TigerSvc:    tigerSvc,
				SightingSvc: sightSvc,
			},
			want: &createOpsResolver{
				Resolver: &Resolver{
					UserSvc:     userSvc,
					TigerSvc:    tigerSvc,
					SightingSvc: sightSvc,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resolver{
				UserSvc:     tt.fields.UserSvc,
				TigerSvc:    tt.fields.TigerSvc,
				SightingSvc: tt.fields.SightingSvc,
			}
			if got := r.CreateOps(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolver.CreateOps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_ListOps(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := mock.NewMockUserService(ctrl)
	tigerSvc := mock.NewMockTigerService(ctrl)
	sightSvc := mock.NewMockSightingService(ctrl)
	type fields struct {
		UserSvc     service.UserService
		TigerSvc    service.TigerService
		SightingSvc service.SightingService
	}
	tests := []struct {
		name   string
		fields fields
		want   generated.ListOpsResolver
	}{
		{
			name: "implemented",
			fields: fields{
				UserSvc:     userSvc,
				TigerSvc:    tigerSvc,
				SightingSvc: sightSvc,
			},
			want: &listOpsResolver{
				Resolver: &Resolver{
					UserSvc:     userSvc,
					TigerSvc:    tigerSvc,
					SightingSvc: sightSvc,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resolver{
				UserSvc:     tt.fields.UserSvc,
				TigerSvc:    tt.fields.TigerSvc,
				SightingSvc: tt.fields.SightingSvc,
			}
			if got := r.ListOps(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolver.ListOps() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_Mutation(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := mock.NewMockUserService(ctrl)
	tigerSvc := mock.NewMockTigerService(ctrl)
	sightSvc := mock.NewMockSightingService(ctrl)
	type fields struct {
		UserSvc     service.UserService
		TigerSvc    service.TigerService
		SightingSvc service.SightingService
	}
	tests := []struct {
		name   string
		fields fields
		want   generated.MutationResolver
	}{
		{
			name: "implemented",
			fields: fields{
				UserSvc:     userSvc,
				TigerSvc:    tigerSvc,
				SightingSvc: sightSvc,
			},
			want: &mutationResolver{
				Resolver: &Resolver{
					UserSvc:     userSvc,
					TigerSvc:    tigerSvc,
					SightingSvc: sightSvc,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resolver{
				UserSvc:     tt.fields.UserSvc,
				TigerSvc:    tt.fields.TigerSvc,
				SightingSvc: tt.fields.SightingSvc,
			}
			if got := r.Mutation(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolver.Mutation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestResolver_Query(t *testing.T) {
	ctrl := gomock.NewController(t)
	userSvc := mock.NewMockUserService(ctrl)
	tigerSvc := mock.NewMockTigerService(ctrl)
	sightSvc := mock.NewMockSightingService(ctrl)
	type fields struct {
		UserSvc     service.UserService
		TigerSvc    service.TigerService
		SightingSvc service.SightingService
	}
	tests := []struct {
		name   string
		fields fields
		want   generated.QueryResolver
	}{
		{
			name: "implemented",
			fields: fields{
				UserSvc:     userSvc,
				TigerSvc:    tigerSvc,
				SightingSvc: sightSvc,
			},
			want: &queryResolver{
				Resolver: &Resolver{
					UserSvc:     userSvc,
					TigerSvc:    tigerSvc,
					SightingSvc: sightSvc,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Resolver{
				UserSvc:     tt.fields.UserSvc,
				TigerSvc:    tt.fields.TigerSvc,
				SightingSvc: tt.fields.SightingSvc,
			}
			if got := r.Query(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Resolver.Query() = %v, want %v", got, tt.want)
			}
		})
	}
}
