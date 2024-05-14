package directive

import (
	"context"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/contexthandler"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

func Auth(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	_, err := contexthandler.GetUserID(ctx)
	if err != nil {
		return nil, &gqlerror.Error{
			Message: "Access Denied",
		}
	}

	return next(ctx)
}
