package wheretopark

import (
	"encoding/json"
	"time"

	"github.com/mmcloughlin/geohash"
	geojson "github.com/paulmach/go.geojson"
	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
)

const (
	SpotTypeCar         = "CAR"
	SpotTypeCarElectric = "CAR_ELECTRIC"
	SpotTypeCarDisabled = "CAR_DISABLED"
	SpotTypeMotorcycle  = "MOTORCYCLE"
	SpotTypeTruck       = "TRUCK"
	SpotTypeBus         = "BUS"
)

const (
	FeatureCovered     = "COVERED"
	FeatureUncovered   = "UNCOVERED"
	FeatureUnderground = "UNDERGROUND"
	FeatureGuarded     = "GUARDED"
	FeatureMonitored   = "MONITORED"
)

const (
	PaymentMethodCash        = "CASH"
	PaymentMethodCard        = "CARD"
	PaymentMethodContactless = "CONTACTLESS"
	PaymentMethodMobile      = "MOBILE"
)

type ID = string
type SpotType = string
type Feature = string
type PaymentMethod = string
type LanguageCode = string

type PricingRule struct {
	Duration  string          `json:"duration"`
	Price     decimal.Decimal `json:"price"`
	Repeating bool            `json:"repeating,omitempty"`
}

func (p PricingRule) Repeated() PricingRule {
	p.Repeating = true
	return p
}

type Rule struct {
	// https://schema.org/openingHours
	// https://wiki.openstreetmap.org/wiki/Key:opening_hours
	Hours string `json:"hours"`
	// If not empty, then applies to only those from this list
	Applies []SpotType    `json:"applies,omitempty"`
	Pricing []PricingRule `json:"pricing"`
}

type Dimensions struct {
	Width  int `json:"width,omitempty"`
	Height int `json:"height,omitempty"`
	Length int `json:"length,omitempty"`
}

func (d *Dimensions) Empty() bool {
	return d.Width == 0 && d.Height == 0 && d.Length == 0
}

type Metadata struct {
	LastUpdated    time.Time               `json:"lastUpdated,omitempty"`
	Name           string                  `json:"name"`
	Address        string                  `json:"address"`
	Geometry       *geojson.Geometry       `json:"geometry"`
	Resources      []string                `json:"resources"`
	TotalSpots     map[SpotType]uint       `json:"totalSpots"`
	MaxDimensions  *Dimensions             `json:"maxDimensions,omitempty"`
	Features       []Feature               `json:"features"`
	PaymentMethods []PaymentMethod         `json:"paymentMethods"`
	Comment        map[LanguageCode]string `json:"comment"`
	Currency       currency.Unit           `json:"currency"`
	Timezone       *time.Location          `json:"timezone"`
	Rules          []Rule                  `json:"rules"`
}

func (m Metadata) MarshalJSON() ([]byte, error) {
	type MetadataAlias Metadata
	return json.Marshal(&struct {
		*MetadataAlias
		LastUpdated string `json:"lastUpdated,omitempty"`
		Currency    string `json:"currency"`
		Timezone    string `json:"timezone"`
	}{
		MetadataAlias: (*MetadataAlias)(&m),
		LastUpdated:   m.LastUpdated.Format(time.DateOnly),
		Currency:      m.Currency.String(),
		Timezone:      m.Timezone.String(),
	})
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
	return geohash.EncodeWithPrecision(latitude, longitude, 10)
}

func GeometryToID(geometry *geojson.Geometry) ID {
	return geohash.EncodeWithPrecision(geometry.Point[1], geometry.Point[0], 10)
}

func (m *Metadata) Validate() error {
	return nil
}
