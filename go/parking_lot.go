package wheretopark

import (
	"github.com/mmcloughlin/geohash"
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
	Geometry       geojson.Geometry        `json:"geometry"`
	Resources      []string                `json:"resources"`
	TotalSpots     map[SpotType]uint       `json:"totalSpots" yaml:"total-spots"`
	MaxWidth       *int                    `json:"maxWidth,omitempty" yaml:"max-width"`
	MaxHeight      *int                    `json:"maxHeight,omitempty" yaml:"max-height"`
	Features       []Feature               `json:"features"`
	PaymentMethods []PaymentMethod         `json:"paymentMethods" yaml:"payment-methods"`
	Comment        map[LanguageCode]string `json:"comment"`
	Currency       string                  `json:"currency"`
	Timezone       string                  `json:"timezone"`
	Rules          []Rule                  `json:"rules"`
}

type State struct {
	LastUpdated    string            `json:"lastUpdated"`
	AvailableSpots map[SpotType]uint `json:"availableSpots"`
}

type ParkingLot struct {
	Metadata Metadata `json:"metadata"`
	State    State    `json:"state"`
}

func CoordinateToID(latitude, longitude float64) ID {
	return geohash.Encode(latitude, longitude)
}
