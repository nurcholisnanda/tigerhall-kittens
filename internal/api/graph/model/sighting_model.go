package model

import (
	"time"

	"gorm.io/gorm"
)

type Sighting struct {
	ID           string    `json:"id"`
	TigerID      string    `json:"tigerID" gorm:"not null;index"`
	LastSeenTime time.Time `json:"lastSeenTime" gorm:"not null"`
	Image        *string   `json:"image" gorm:"type:varchar(255)"`
	*Coordinate  `json:"coordinate"`
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	DeletedBy    string
}
