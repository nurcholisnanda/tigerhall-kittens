package contexthandler

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type ContextKey string

type JwtCustomClaim struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

func RetrieveGinContext(ctx context.Context, key string) (*gin.Context, error) {
	ginContext := ctx.Value(ContextKey(key))
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}

	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}

func SetContext(ctx context.Context, key string, val any) context.Context {
	return context.WithValue(ctx, ContextKey(key), val)
}

func GetUserID(ctx context.Context) (string, error) {
	gc, err := RetrieveGinContext(ctx, "ContextKey")
	if err != nil {
		return "", err
	}
	tokenData, ok := gc.Value("auth").(*JwtCustomClaim)
	if !ok {
		return "", errors.New("Access Denied")
	}
	return tokenData.ID, nil
}
