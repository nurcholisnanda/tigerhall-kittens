package handler

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/directive"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/generated"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/service"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/contexthandler"
)

// Defining the Graphql handler
func GraphqlHandler(
	userSvc service.UserService,
	tigerSvc service.TigerService,
	sightingSvc service.SightingService,
) gin.HandlerFunc {
	c := generated.Config{Resolvers: &graph.Resolver{
		UserSvc:     userSvc,
		TigerSvc:    tigerSvc,
		SightingSvc: sightingSvc,
	}}
	c.Directives.Auth = directive.Auth

	h := handler.NewDefaultServer(generated.NewExecutableSchema(c))

	return func(c *gin.Context) {
		// Set up GraphQL context (pass the Gin context)
		ctx := contexthandler.SetContext(c.Request.Context(), "ContextKey", c)
		// Serve the GraphQL request
		h.ServeHTTP(c.Writer, c.Request.WithContext(ctx))
	}
}

// Defining the Playground handler
func PlaygroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
