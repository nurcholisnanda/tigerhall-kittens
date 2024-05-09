package main

import (
	"log"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/nurcholisnanda/tigerhall-kittens/config"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/directive"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/middlewares"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/service"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/bcrypt"
)

const defaultPort = "8080"

// Defining the Graphql handler
func graphqlHandler(userSvc service.UserService) gin.HandlerFunc {
	// NewExecutableSchema and Config are in the generated.go file
	// Resolver is in the resolver.go file
	c := graph.Config{Resolvers: &graph.Resolver{
		UserSvc: userSvc,
	}}
	c.Directives.Auth = directive.Auth

	h := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

// Defining the Playground handler
func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func main() {
	//setup database
	db, err := config.NewDatabase()
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate() // Automatically migrate database schema
	gormDB := db.GetDB()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	//initializes dependencies
	userRepo := repository.NewUserRepoImpl(gormDB)
	JWT := service.NewJWT(os.Getenv("SECRET"))
	userSvc := service.NewUserService(userRepo, bcrypt.NewBcrypt(), JWT)
	// middlewares := middlewares.New(JWT)

	// router := mux.NewRouter()
	// router.Use(middlewares.AuthMiddleware)

	// Setting up Gin
	r := gin.Default()
	r.Use(middlewares.AuthMiddleware(JWT))
	r.POST("/query", graphqlHandler(userSvc))
	r.GET("/", playgroundHandler())

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	r.Run()

	// c := graph.Config{Resolvers: &graph.Resolver{
	// 	UserSvc: userSvc,
	// }}
	// c.Directives.Auth = directives.Auth

	// srv := handler.NewDefaultServer(graph.NewExecutableSchema(c))

	// router.Handle("/", playground.Handler("GraphQL playground", "/query"))
	// router.Handle("/query", srv)

	// log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	// log.Fatal(http.ListenAndServe(":"+port, router))
}
