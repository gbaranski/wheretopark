package wheretopark

import (
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/mmcloughlin/geohash"
	geojson "github.com/paulmach/go.geojson"
	"github.com/shopspring/decimal"
	"golang.org/x/exp/slices"
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

var (
	SpotTypes      = []SpotType{SpotTypeCar, SpotTypeCarElectric, SpotTypeCarDisabled, SpotTypeMotorcycle, SpotTypeTruck, SpotTypeBus}
	Features       = []Feature{FeatureCovered, FeatureUncovered, FeatureUnderground, FeatureGuarded, FeatureMonitored}
	PaymentMethods = []PaymentMethod{PaymentMethodCash, PaymentMethodCard, PaymentMethodContactless, PaymentMethodMobile}
)

type ID = string
type SpotType = string
type Feature = string
type PaymentMethod = string
type LanguageCode = string

type PricingRule struct {
	Duration  string          `json:"duration"`        // ISO8601 Duration
	Limit     string          `json:"limit,omitempty"` // ISO8601 Duration
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
	LastUpdated    *time.Time              `json:"lastUpdated,omitempty"`
	Name           string                  `json:"name"`
	Address        string                  `json:"address"`
	Geometry       *geojson.Geometry       `json:"geometry"`
	Resources      []string                `json:"resources"`
	TotalSpots     map[SpotType]uint       `json:"totalSpots"`
	MaxDimensions  *Dimensions             `json:"maxDimensions,omitempty"`
	Features       []Feature               `json:"features,omitempty"`
	PaymentMethods []PaymentMethod         `json:"paymentMethods,omitempty"`
	Comment        map[LanguageCode]string `json:"comment,omitempty"`
	Currency       currency.Unit           `json:"currency"`
	Timezone       *time.Location          `json:"timezone"`
	Rules          []Rule                  `json:"rules,omitempty"`
}

type metadataAlias Metadata

type metadataJSON struct {
	*metadataAlias
	LastUpdated string `json:"lastUpdated,omitempty"`
	Currency    string `json:"currency"`
	Timezone    string `json:"timezone"`
}

func (m Metadata) MarshalJSON() ([]byte, error) {
	var lastUpdated string
	if m.LastUpdated != nil {
		lastUpdated = m.LastUpdated.Format(time.DateOnly)
	}
	return json.Marshal(
		&metadataJSON{
			metadataAlias: (*metadataAlias)(&m),
			LastUpdated:   lastUpdated,
			Currency:      m.Currency.String(),
			Timezone:      m.Timezone.String(),
		},
	)
}

func (m *Metadata) UnmarshalJSON(data []byte) error {
	aux := metadataJSON{
		metadataAlias: (*metadataAlias)(m),
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	if aux.LastUpdated > "" {
		if aux.LastUpdated < time.DateOnly {
			return fmt.Errorf("invalid date format: %s", aux.LastUpdated)
		}
		lastUpdated, err := time.Parse(time.DateOnly, aux.LastUpdated[:len(time.DateOnly)])
		if err != nil {
			return err
		}
		m.LastUpdated = &lastUpdated
	}
	m.Currency, err = currency.ParseISO(aux.Currency)
	if err != nil {
		return err
	}
	m.Timezone, err = time.LoadLocation(aux.Timezone)
	if err != nil {
		return err
	}
	return nil
}

type State struct {
	LastUpdated    time.Time         `json:"lastUpdated"`
	AvailableSpots map[SpotType]uint `json:"availableSpots"`
}

type stateAlias State

type stateJSON struct {
	*stateAlias
	LastUpdated string `json:"lastUpdated"`
}

func (s State) MarshalJSON() ([]byte, error) {
	return json.Marshal(&stateJSON{
		stateAlias:  (*stateAlias)(&s),
		LastUpdated: s.LastUpdated.Format(time.RFC3339),
	})
}

func (s *State) UnmarshalJSON(data []byte) error {
	aux := stateJSON{
		stateAlias: (*stateAlias)(s),
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}
	s.LastUpdated, err = time.Parse(time.RFC3339, aux.LastUpdated)
	if err != nil {
		return err
	}
	return nil
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

func (p *ParkingLot) Validate() error {
	if err := p.Metadata.Validate(); err != nil {
		return err
	}
	if err := p.State.Validate(); err != nil {
		return err
	}
	if p.State.AvailableSpots == nil {
		return fmt.Errorf("availableSpots must not be empty")
	}
	for spotType := range p.State.AvailableSpots {
		if _, exists := p.Metadata.TotalSpots[spotType]; !exists {
			return fmt.Errorf("spotType %s from availableSpots is not defined in totalSpots", spotType)
		}
	}
	return nil
}

func (m *Metadata) Validate() error {
	if m.LastUpdated != nil && m.LastUpdated.Unix() == 0 {
		return fmt.Errorf("lastUpdated must not be set to zero")
	}
	if m.Name == "" {
		return fmt.Errorf("name must not be empty")
	}
	if m.Address == "" {
		return fmt.Errorf("address must not be empty")
	}
	defaultGeometry := geojson.Geometry{}
	if m.Geometry == nil || m.Geometry == &defaultGeometry {
		return fmt.Errorf("geometry must not be empty")
	}
	if m.TotalSpots == nil {
		return fmt.Errorf("totalSpots must not be empty")
	}
	for spotType := range m.TotalSpots {
		if !slices.Contains(SpotTypes, spotType) {
			return fmt.Errorf("invalid spotType: %s", spotType)
		}
	}
	for _, resource := range m.Resources {
		if _, err := url.Parse(resource); err != nil {
			return fmt.Errorf("invalid resource URL: %s", resource)
		}
	}
	for _, feature := range m.Features {
		if !slices.Contains(Features, feature) {
			return fmt.Errorf("invalid feature: %s", feature)
		}
	}
	for _, paymentMethod := range m.PaymentMethods {
		if !slices.Contains(PaymentMethods, paymentMethod) {
			return fmt.Errorf("invalid paymentMethod: %s", paymentMethod)
		}
	}
	if m.Timezone == nil {
		return fmt.Errorf("timezone must not be nil")
	}

	// TODO: Validate rules opening hours
	if m.Rules == nil {
		return fmt.Errorf("rules must not be empty")
	}

	return nil
}

func (s *State) Validate() error {
	if s.LastUpdated.Unix() == 0 {
		return fmt.Errorf("lastUpdated must not be set to zero")
	}

	for spotType := range s.AvailableSpots {
		if !slices.Contains(SpotTypes, spotType) {
			return fmt.Errorf("invalid spotType: %s", spotType)
		}
	}
	return nil
}
