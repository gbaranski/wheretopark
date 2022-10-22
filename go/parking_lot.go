package wheretopark

import (
	geojson "github.com/paulmach/go.geojson"
)

type ID = string
type SpotType = string
type Feature = string
type PaymentMethod = string
type LanguageCode = string

type PricingRule struct {
	Duration  string  `json:"duration"`
	Price     float32 `json:"price"`
	Repeating bool    `json:"repeating,omitempty"`
}

type Rule struct {
	// https://schema.org/openingHours
	// https://wiki.openstreetmap.org/wiki/Key:opening_hours
	Hours string `json:"hours"`
	// If not empty, then applies to only those from this list
	Applies []SpotType    `json:"applies,omitempty"`
	Pricing []PricingRule `json:"pricing"`
}

type Metadata struct {
	Name           string                  `json:"name"`
	Address        string                  `json:"address"`
	Location       geojson.Feature         `json:"location"`
	Resources      []string                `json:"resources"`
	TotalSpots     map[SpotType]uint       `json:"total-spots"`
	MaxWidth       *int                    `json:"max-width,omitempty"`
	MaxHeight      *int                    `json:"max-height,omitempty"`
	Features       []Feature               `json:"features"`
	PaymentMethods []PaymentMethod         `json:"payment-methods"`
	Comment        map[LanguageCode]string `json:"comment"`
	Currency       string                  `json:"currency"`
	Timezone       string                  `json:"timezone"`
	Rules          []Rule                  `json:"rules"`
}

type State struct {
	LastUpdated    string            `json:"last-updated"`
	AvailableSpots map[SpotType]uint `json:"available-spots"`
}

type ParkingLot struct {
	Metadata Metadata `json:"metadata"`
	State    State    `json:"state"`
}
