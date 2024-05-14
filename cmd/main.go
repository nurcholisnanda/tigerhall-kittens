package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/nurcholisnanda/tigerhall-kittens/config"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/handler"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/middlewares"
)

const defaultPort = "8080"

func main() {
	//setup env variables
	godotenv.Load()

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	//Set Up Dependencies
	deps, err := config.InitDependencies()
	if err != nil {
		log.Fatalf("Error initializing dependencies: %v", err)
	}

	//Start Notification Consumer
	ctx := context.Background()
	go deps.NotificationService.StartNotificationConsumer(ctx)
	defer deps.NotificationService.CloseNotificationChannel()

	// Setting up Gin
	r := gin.Default()
	r.Use(
		middlewares.Authenticate(deps.UserService, deps.JWTService),
		middlewares.RequestIDMiddleware(),
		middlewares.LoggerMiddleware(),
	)
	r.POST("/query", handler.GraphqlHandler(deps.UserService, deps.TigerService, deps.SightingService))
	r.GET("/", handler.PlaygroundHandler())

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	r.Run()
}
