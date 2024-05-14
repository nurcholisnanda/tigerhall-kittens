// config/config.go

package config

import (
	"fmt"
	"log"
	"os"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/repository"
	"github.com/nurcholisnanda/tigerhall-kittens/internal/service"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/bcrypt"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/mailer"
	"github.com/nurcholisnanda/tigerhall-kittens/pkg/storage"
)

// Dependencies is a struct to handle all dependencies needed for this codebase.
type Dependencies struct {
	SightingRepo        repository.SightingRepository
	TigerRepo           repository.TigerRepository
	UserRepo            repository.UserRepository
	Bcrypt              bcrypt.BcryptInterface
	JWTService          service.JWT
	UserService         service.UserService
	TigerService        service.TigerService
	SightingService     service.SightingService
	NotificationService service.NotifService
}

// InitDependencies is a function to initiate all dependencies
// and will return Dependencies object if there is no error in the initiation process.
func InitDependencies() (Dependencies, error) {
	db, err := NewDatabase()
	if err != nil {
		log.Panic(err)
	}
	db.AutoMigrate() // Automatically migrate database schema
	gormDB := db.GetDB()

	s3Client, err := storage.NewS3Client()
	if err != nil {
		return Dependencies{}, fmt.Errorf("failed to create S3 client: %w", err)
	}
	userRepo := repository.NewUserRepoImpl(gormDB)
	tigerRepo := repository.NewTigerRepositoryImpl(gormDB)
	sightingRepo := repository.NewSightingRepositoryImpl(gormDB)
	mailer := mailer.NewMailService()
	notificationSvc := service.NewNotificationService(sightingRepo, userRepo, mailer)
	JWT := service.NewJWT(os.Getenv("JWT_SECRET"))
	userSvc := service.NewUserService(userRepo, bcrypt.NewBcrypt(), JWT)
	tigerSvc := service.NewTigerService(tigerRepo)
	sightingSvc := service.NewSightingService(notificationSvc, sightingRepo, tigerRepo, s3Client)

	return Dependencies{
		SightingRepo:        sightingRepo,
		TigerRepo:           tigerRepo,
		UserRepo:            userRepo,
		Bcrypt:              bcrypt.NewBcrypt(),
		JWTService:          JWT,
		UserService:         userSvc,
		TigerService:        tigerSvc,
		SightingService:     sightingSvc,
		NotificationService: notificationSvc,
	}, nil
}
