package model // Adjust the package name based on your project structure

import (
	"time"
)

type Notification struct {
	Sighter   string
	TigerID   string    `json:"tigerID"`
	Timestamp time.Time `json:"timestamp"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}
