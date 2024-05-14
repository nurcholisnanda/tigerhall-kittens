package model

import (
	"time"

	"gorm.io/gorm"
)

type Tiger struct {
	ID           string    `json:"id"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null;unique"`
	DateOfBirth  time.Time `json:"dateOfBirth" gorm:"not null"`
	LastSeenTime time.Time `json:"lastSeenTime" gorm:"not null"`
	*Coordinate  `json:"lastSeenCoordinate"`
	CreatedAt    time.Time
	CreatedBy    string
	UpdatedAt    time.Time
	UpdatedBy    string
	DeletedAt    gorm.DeletedAt `gorm:"index"`
	DeletedBy    string
}
