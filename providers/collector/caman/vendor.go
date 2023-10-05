package caman

import (
	"net/url"
	"time"
)

type CameraID = string
type CameraState struct {
	LastUpdated    time.Time `json:"lastUpdated"`
	TotalSpots     uint      `json:"totalSpots"`
	AvailableSpots uint      `json:"availableSpots"`
}

type CameraMetadata struct {
	Url   url.URL `json:"url"`
	Spots []ParkingSpot
}
