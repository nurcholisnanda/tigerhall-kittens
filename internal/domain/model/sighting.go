package model

import "time"

type Sighting struct {
	ID        uint
	TigerID   uint
	Timestamp time.Time
	ImageURL  string
	Coordinates
}
