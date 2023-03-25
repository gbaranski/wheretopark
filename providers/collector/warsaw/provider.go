package warsaw

import (
	"fmt"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/go/provider/simple"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"

	geojson "github.com/paulmach/go.geojson"
)

type Provider struct{}

func (p Provider) Config() simple.Config {
	return simple.DEFAULT_CONFIG
}

func (p Provider) Name() string {
	return "warsaw"
}

func (p Provider) GetParkingLots() (map[wheretopark.ID]wheretopark.ParkingLot, error) {
	data, err := GetData()
	if err != nil {
		return nil, err
	}
	lastUpdate, err := time.ParseInLocation("2006-01-02T15:04:05", data.Result.Timestamp, defaultTimezone)
	if err != nil {
		return nil, fmt.Errorf("failed to parse last update time: %w", err)
	}

	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot)
	for _, vendor := range data.Result.Parks {
		id := wheretopark.CoordinateToID(vendor.Latitude, vendor.Longitude)
		configuration, exists := configuration.ParkingLots[id]
		if !exists {
			log.Warn().
				Str("name", vendor.Name).
				Str("provider", p.Name()).
				Msg("missing configuration")
			continue
		}
		maxDimensions := &wheretopark.Dimensions{}
		for _, dimension := range vendor.Dimensions {
			width, err := decimal.NewFromString(dimension.Width)
			if err != nil {
				return nil, fmt.Errorf("failed to parse decimal width: `%s`", dimension.Width)
			}
			length, err := decimal.NewFromString(dimension.Width)
			if err != nil {
				return nil, fmt.Errorf("failed to parse decimal width: `%s`", dimension.Width)
			}
			widthInCm := int(width.Mul(decimal.NewFromInt(100)).IntPart())
			lengthInCm := int(length.Mul(decimal.NewFromInt(100)).IntPart())
			if maxDimensions.Width == 0 || widthInCm < maxDimensions.Width {
				maxDimensions.Width = widthInCm
			}
			if maxDimensions.Length == 0 || lengthInCm < maxDimensions.Length {
				maxDimensions.Length = lengthInCm
			}
		}
		if maxDimensions.Empty() {
			maxDimensions = nil
		}
		totalPlaces := TotalPlaces{
			Disabled: 0,
			Standard: 0,
			Electric: 0,
		}
		for _, e := range vendor.TotalPlaces {
			totalPlaces.Disabled += e.Disabled
			totalPlaces.Electric += e.Electric
			totalPlaces.Standard += e.Standard
		}
		metadata := wheretopark.Metadata{
			LastUpdated:    configuration.LastUpdated,
			Name:           vendor.Name,
			Address:        configuration.Address,
			Geometry:       geojson.NewPointGeometry([]float64{vendor.Longitude, vendor.Latitude}),
			Resources:      configuration.Resources,
			TotalSpots:     map[string]uint{wheretopark.SpotTypeCarElectric: totalPlaces.Electric, wheretopark.SpotTypeCar: totalPlaces.Standard, wheretopark.SpotTypeCarDisabled: totalPlaces.Disabled},
			MaxDimensions:  maxDimensions,
			Features:       configuration.Features,
			PaymentMethods: configuration.PaymentMethods,
			Comment:        configuration.Comment,
			Currency:       currency.PLN,
			Timezone:       defaultTimezone,
			Rules:          configuration.Rules,
		}

		state := wheretopark.State{
			LastUpdated: lastUpdate.In(defaultTimezone),
			AvailableSpots: map[string]uint{
				wheretopark.SpotTypeCarElectric: vendor.FreePlacesTotal.Electric,
				wheretopark.SpotTypeCar:         vendor.FreePlacesTotal.Public,
				wheretopark.SpotTypeCarDisabled: vendor.FreePlacesTotal.Disabled,
			},
		}
		parkingLots[id] = wheretopark.ParkingLot{
			Metadata: metadata,
			State:    state,
		}
	}
	return parkingLots, nil
}

func (p Provider) GetState() (map[wheretopark.ID]wheretopark.State, error) {
	data, err := GetData()
	if err != nil {
		return nil, err
	}
	lastUpdate, err := time.ParseInLocation("2006-01-02T15:04:05", data.Result.Timestamp, defaultTimezone)
	if err != nil {
		return nil, fmt.Errorf("failed to parse last update time: %w", err)
	}
	states := make(map[wheretopark.ID]wheretopark.State)

	for _, vendor := range data.Result.Parks {
		id := wheretopark.CoordinateToID(vendor.Latitude, vendor.Longitude)

		state := wheretopark.State{
			LastUpdated: lastUpdate.In(defaultTimezone),
			AvailableSpots: map[string]uint{
				wheretopark.SpotTypeCarElectric: vendor.FreePlacesTotal.Electric,
				wheretopark.SpotTypeCar:         vendor.FreePlacesTotal.Public,
				wheretopark.SpotTypeCarDisabled: vendor.FreePlacesTotal.Disabled,
			},
		}
		states[id] = state
	}
	return states, nil
}

func NewProvider() (simple.Provider, error) {
	return Provider{}, nil
}
