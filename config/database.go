package config

import (
	"fmt"
	"log"
	"os"

	"github.com/nurcholisnanda/tigerhall-kittens/internal/api/graph/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// database represents a wrapper around the GORM database instance.
type database struct {
	db *gorm.DB
}

// DBUrl constructs the database connection URL using environment variables.
func DBUrl() string {
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASS"),
	)
}

// NewDatabase initializes a new database connection using GORM.
func NewDatabase() (*database, error) {
	db, err := gorm.Open(postgres.Open(DBUrl()), &gorm.Config{})
	if err != nil {
		log.Panic(err)
	}
	return &database{
		db: db,
	}, nil
}

// GetDB returns the underlying GORM database instance.
func (r *database) GetDB() *gorm.DB {
	return r.db
}

// AutoMigrate performs automatic schema migration for defined models.
func (r *database) AutoMigrate() error {
	return r.db.AutoMigrate(&model.User{})
}
