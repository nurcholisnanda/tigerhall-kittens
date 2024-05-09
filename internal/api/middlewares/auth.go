package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/service"
)

type ContextKey string

// type authMiddleware struct {
// 	JWT service.JWT
// }

// func New(JWT service.JWT) *authMiddleware {
// 	return &authMiddleware{
// 		JWT: JWT,
// 	}
// }

// func (a *authMiddleware) AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		auth := r.Header.Get("Authorization")

// 		if auth == "" {
// 			next.ServeHTTP(w, r)
// 			return
// 		}

// 		bearer := "Bearer "
// 		auth = auth[len(bearer):]

// 		validate, err := a.JWT.ValidateToken(context.Background(), auth)
// 		if err != nil || !validate.Valid {
// 			http.Error(w, "Invalid token", http.StatusForbidden)
// 			return
// 		}

// 		customClaim, _ := validate.Claims.(*service.JwtCustomClaim)

// 		ctx := context.WithValue(r.Context(), authString("auth"), customClaim)

// 		r = r.WithContext(ctx)
// 		next.ServeHTTP(w, r)
// 	})
// }

// func CtxValue(ctx context.Context) *service.JwtCustomClaim {
// 	raw, _ := ctx.Value(authString("auth")).(*service.JwtCustomClaim)
// 	return raw
// }

// AuthMiddleware is a Gin middleware for JWT-based authentication
func AuthMiddleware(jwtService service.JWT) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Retrieve JWT token from request header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Next()
			return
		}

		tokenString := authHeader[len("Bearer "):]

		// Validate and parse the JWT token
		claims, err := jwtService.ValidateToken(context.Background(), tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}
		customClaim, _ := claims.Claims.(*service.JwtCustomClaim)
		c.Set("auth", customClaim)
		ctx := context.WithValue(c.Request.Context(), ContextKey("ContextKey"), c)
		c.Request = c.Request.WithContext(ctx)
		// Continue to next middleware or handler
		c.Next()
	}
}

func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(ContextKey("ContextKey"))
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
