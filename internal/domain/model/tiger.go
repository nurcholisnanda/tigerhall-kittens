package model

import "time"

type Tiger struct {
	ID                uint
	Name              string
	DateOfBirth       time.Time
	LastSeenTimestamp time.Time
	Coordinates
}

type Coordinates struct {
	Lat float64
	Lon float64
}
