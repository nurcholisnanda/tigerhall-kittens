package model

import (
	"time"

	"gorm.io/gorm"
)

type Sighting struct {
	ID                  string    `json:"id"`
	TigerID             string    `json:"tigerID" gorm:"not null"`
	LastSeenTime        time.Time `json:"lastSeenTime" gorm:"not null"`
	Image               *string   `json:"image" gorm:"type:text"`
	*LastSeenCoordinate `json:"lastSeenCoordinate"`
	CreatedAt           time.Time
	CreatedBy           string `gorm:"index"`
	UpdatedAt           time.Time
	UpdatedBy           string
	DeletedAt           gorm.DeletedAt `gorm:"index"`
	DeletedBy           string
}
