package service

import (
	"context"
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/helper"
)

func TestNewJWT(t *testing.T) {
	secret := "secret"
	type args struct {
		secret string
	}
	tests := []struct {
		name string
		args args
		want JWT
	}{
		{
			name: "success",
			args: args{
				secret: secret,
			},
			want: NewJWT(secret),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewJWT(tt.args.secret); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewJWT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authClient_GenerateToken(t *testing.T) {
	secret := "secret"
	type fields struct {
		secret string
	}
	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "should return error if id is nil",
			fields: fields{
				secret: secret,
			},
			args: args{
				ctx:    context.Background(),
				userID: "",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &authClient{
				secret: tt.fields.secret,
			}
			got, err := c.GenerateToken(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("authClient.GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("authClient.GenerateToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authClient_ValidateToken(t *testing.T) {
	secret := "secret"
	userID := uuid.New().String()
	JWT := NewJWT(secret)
	ctx := context.Background()
	token, err := JWT.GenerateToken(ctx, userID)

	expiredClaims := jwt.RegisteredClaims{ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * -24))}
	expiredCustomClaims := helper.JwtCustomClaim{
		ID:               userID,
		RegisteredClaims: expiredClaims,
	}
	expiredJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredCustomClaims)
	expiredToken, _ := expiredJWT.SignedString([]byte(secret))

	if err != nil {
		log.Fatal(err)
	}
	type fields struct {
		secret string
	}
	type args struct {
		ctx          context.Context
		requestToken string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *jwt.Token
		wantErr bool
	}{
		{
			name: "should return error if using an expired token",
			fields: fields{
				secret: secret,
			},
			args: args{
				ctx:          context.Background(),
				requestToken: expiredToken,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "success",
			fields: fields{
				secret: secret,
			},
			args: args{
				ctx:          context.Background(),
				requestToken: token,
			},
			want: &jwt.Token{
				Valid: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &authClient{
				secret: tt.fields.secret,
			}
			got, err := c.ValidateToken(tt.args.ctx, tt.args.requestToken)
			if (err != nil) != tt.wantErr {
				t.Errorf("authClient.ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != nil {
				if !reflect.DeepEqual(got.Valid, tt.want.Valid) {
					t.Errorf("authClient.ValidateToken() = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
