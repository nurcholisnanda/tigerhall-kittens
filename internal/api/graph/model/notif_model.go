package model // Adjust the package name based on your project structure

import "time"

type Notification struct {
	TigerID    string    `json:"tigerID"`
	SightingID string    `json:"sightingID"`
	Timestamp  time.Time `json:"timestamp"`
}
