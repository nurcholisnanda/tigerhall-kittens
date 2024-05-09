package directive

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/middlewares"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/service"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	gc, err := middlewares.GinContextFromContext(ctx)
	if err != nil {
		return nil, err
	}

	tokenData, _ := gc.Value("auth").(*service.JwtCustomClaim)
	if tokenData == nil {
		return nil, &gqlerror.Error{
			Message: "Access Denied",
		}
	}

	return next(ctx)
}
