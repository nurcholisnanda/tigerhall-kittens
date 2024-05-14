package middlewares

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/service"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/contexthandler"
)

type ContextKey string

// Authenticate is a Gin middleware for JWT-based authentication
func Authenticate(userSvc service.UserService, jwtSvc service.JWT) gin.HandlerFunc {
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
		claims, err := jwtSvc.ValidateToken(ctx, tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		customClaim, _ := claims.Claims.(*contexthandler.JwtCustomClaim)

		if _, err := userSvc.GetUserByID(ctx, customClaim.ID); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("auth", customClaim)
		ctx = contexthandler.SetContext(c.Request.Context(), "ContextKey", c)
		c.Request = c.Request.WithContext(ctx)
		// Continue to next middleware or handler
		c.Next()
	}
}
