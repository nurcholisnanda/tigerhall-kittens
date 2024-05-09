package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaim struct {
	ID string `json:"id"`
	jwt.RegisteredClaims
}

type authClient struct {
	secret string
}

// Authentication client constructor
func NewJWT(secret string) JWT {
	return &authClient{
		secret: secret,
	}
}

// CreateAccessToken will create access token that will be used for user authentication.
// access token will be needed in API that needs user to be authorized
func (c *authClient) GenerateToken(ctx context.Context, userID string) (string, error) {
	if userID == "" {
		return "", errors.New("user id is required")
	}

	currentTime := time.Now().UTC()

	accessTokenClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Hour * 24)),
		IssuedAt:  jwt.NewNumericDate(currentTime),
	}

	claims := JwtCustomClaim{
		ID:               userID,
		RegisteredClaims: accessTokenClaims,
	}

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := at.SignedString([]byte(c.secret))
	if err != nil {
		return "", err
	}

	return accessToken, err
}

// ValidateTokens will validate whether the token is valid
// and will return claims if the user is exist in our database
// otherwise it will return error
func (c *authClient) ValidateToken(ctx context.Context, requestToken string) (*jwt.Token, error) {
	token, _ := jwt.ParseWithClaims(requestToken, &JwtCustomClaim{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("there's a problem with the signing method")
		}
		return []byte(c.secret), nil
	})

	claims, ok := token.Claims.(*JwtCustomClaim)
	if !ok {
		return nil, errors.New("fail to claim token")
	}

	if claims.ExpiresAt.Before(time.Now().UTC()) {
		return nil, errors.New("token is expired please re-login")
	}

	return token, nil
}
