package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.46

import (
	"context"
	"fmt"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/errorhandler"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/logger"
	"github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

// Login is the resolver for the login field.
func (r *authOpsResolver) Login(ctx context.Context, obj *model.AuthOps, email string, password string) (interface{}, error) {
	return r.UserSvc.Login(ctx, email, password)
}

// Register is the resolver for the register field.
func (r *authOpsResolver) Register(ctx context.Context, obj *model.AuthOps, input model.NewUser) (interface{}, error) {
	return r.UserSvc.Register(ctx, &input)
}

// CreateSighting is the resolver for the createSighting field.
func (r *createOpsResolver) CreateSighting(ctx context.Context, obj *model.CreateOps, input model.SightingInput) (*model.Sighting, error) {
	createdSighting, err := r.SightingSvc.CreateSighting(ctx, &input)
	if err != nil {
		// Error Handling in the Resolver
		switch err.(type) {
		case *errorhandler.InvalidCoordinatesError:
			return nil, &gqlerror.Error{
				Message: "invalid coordinates",
				Extensions: map[string]interface{}{
					"code":    errorhandler.INVALID_INPUT,
					"details": err.Error(),
				},
			}
		case *errorhandler.TigerNotFound:
			return nil, &gqlerror.Error{
				Message: "tiger not found",
				Extensions: map[string]interface{}{
					"code":    errorhandler.NOT_FOUND,
					"details": err.Error(),
				},
			}
		case *errorhandler.SightingTooCloseError:
			return nil, &gqlerror.Error{
				Message: "slighting too close",
				Extensions: map[string]interface{}{
					"code":    errorhandler.CONFLICT,
					"details": err.Error(),
				},
			}
		default:
			// Log the unexpected error for investigation
			logger.Logger(ctx).Error(ctx, "Unexpected error creating sighting", "error", err)
			return nil, fmt.Errorf("internal Server Error")
		}
	}

	return createdSighting, nil
}

// CreateTiger is the resolver for the createTiger field.
func (r *createOpsResolver) CreateTiger(ctx context.Context, obj *model.CreateOps, input model.TigerInput) (*model.Tiger, error) {
	tiger, err := r.TigerSvc.CreateTiger(ctx, &input)
	if err != nil {
		switch err.(type) {
		case *errorhandler.InvalidCoordinatesError:
			return nil, &gqlerror.Error{
				Message: "invalid coordinates",
				Extensions: map[string]interface{}{
					"code":    errorhandler.INVALID_INPUT,
					"details": err.Error(),
				},
			}
		case *errorhandler.InvalidDateOfBirthError:
			return nil, &gqlerror.Error{
				Message: "invalid date of birth",
				Extensions: map[string]interface{}{
					"code":    errorhandler.INVALID_INPUT,
					"details": err.Error(),
				},
			}
		case *errorhandler.InvalidLastSeenTimeError:
			return nil, &gqlerror.Error{
				Message: "invalid last seen time",
				Extensions: map[string]interface{}{
					"code":    errorhandler.INVALID_INPUT,
					"details": err.Error(),
				},
			}
		case *errorhandler.TigerCreationError:
			return nil, &gqlerror.Error{
				Message: "failed to create tiger",
				Extensions: map[string]interface{}{
					"code":    errorhandler.INVALID_INPUT,
					"details": err.Error(),
				},
			}
		default:
			// Log the unexpected error for investigation
			logrus.Error(ctx, "Unexpected error creating tiger", "error:", err.Error())
			return nil, gqlerror.Errorf("Internal Server Error")
		}
	}

	return tiger, nil
}

// ListTigers is the resolver for the ListTigers field.
func (r *listOpsResolver) ListTigers(ctx context.Context, obj *model.ListOps, limit int, offset int) ([]*model.Tiger, error) {
	// Call your tiger service to fetch tigers with pagination
	tigers, err := r.TigerSvc.ListTigers(ctx, limit, offset)
	if err != nil {
		// Log the unexpected error for investigation
		logrus.Error(ctx, "Unexpected error getting tiger list", "error:", err.Error())
		return nil, gqlerror.Errorf("Internal Server Error")
	}
	return tigers, nil
}

// ListSightings is the resolver for the listSightings field.
func (r *listOpsResolver) ListSightings(ctx context.Context, obj *model.ListOps, tigerID string, limit int, offset int) ([]*model.Sighting, error) {
	sightings, err := r.SightingSvc.ListSightings(ctx, tigerID, limit, offset)
	if err != nil {
		// Log the unexpected error for investigation
		logrus.Error(ctx, "Unexpected error getting sighting list", "error:", err.Error())
		return nil, gqlerror.Errorf("Internal Server Error")
	}
	return sightings, nil
}

// Auth is the resolver for the auth field.
func (r *mutationResolver) Auth(ctx context.Context) (*model.AuthOps, error) {
	return &model.AuthOps{}, nil
}

// Create is the resolver for the create field.
func (r *mutationResolver) Create(ctx context.Context) (*model.CreateOps, error) {
	return &model.CreateOps{}, nil
}

// User is the resolver for the user field.
func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return r.UserSvc.GetUserByID(ctx, id)
}

// List is the resolver for the list field.
func (r *queryResolver) List(ctx context.Context) (*model.ListOps, error) {
	return &model.ListOps{}, nil
}

// AuthOps returns AuthOpsResolver implementation.
func (r *Resolver) AuthOps() AuthOpsResolver { return &authOpsResolver{r} }

// CreateOps returns CreateOpsResolver implementation.
func (r *Resolver) CreateOps() CreateOpsResolver { return &createOpsResolver{r} }

// ListOps returns ListOpsResolver implementation.
func (r *Resolver) ListOps() ListOpsResolver { return &listOpsResolver{r} }

// Mutation returns MutationResolver implementation.
func (r *Resolver) Mutation() MutationResolver { return &mutationResolver{r} }

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type authOpsResolver struct{ *Resolver }
type createOpsResolver struct{ *Resolver }
type listOpsResolver struct{ *Resolver }
type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
