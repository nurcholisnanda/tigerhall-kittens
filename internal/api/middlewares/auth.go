package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/service"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/helper"
)

type ContextKey string

type AuthMiddleware struct {
	userSvc service.UserService
	jwtSvc  service.JWT
}

func NewAuthMiddleware(userSvc service.UserService, jwtSvc service.JWT) *AuthMiddleware {
	return &AuthMiddleware{
		userSvc: userSvc,
		jwtSvc:  jwtSvc,
	}
}

// AuthMiddleware is a Gin middleware for JWT-based authentication
func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		// Retrieve JWT token from request header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Next()
			return
		}

		tokenString := authHeader[len("Bearer "):]

		// Validate and parse the JWT token
		claims, err := m.jwtSvc.ValidateToken(ctx, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		customClaim, _ := claims.Claims.(*helper.JwtCustomClaim)

		if _, err := m.userSvc.GetUserByID(ctx, customClaim.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("auth", customClaim)
		ctx = helper.SetContext(c.Request.Context(), "ContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		// Continue to next middleware or handler
		c.Next()
	}
}
