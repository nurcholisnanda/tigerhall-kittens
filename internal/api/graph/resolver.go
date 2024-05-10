package graph

import "github.com/nurcholisnanda/tigerhall-kittens/internal/service"

//go:generate go get github.com/99designs/gqlgen
//go:generate go run github.com/99designs/gqlgen
// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	UserSvc  service.UserService
	TigerSvc service.TigerService
}
