package warsaw

import (
	"fmt"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"

	geojson "github.com/paulmach/go.geojson"
)

var defaultLocation *time.Location

func init() {
	location, err := time.LoadLocation("Europe/Warsaw")
	if err != nil {
		panic(err)
	}
	defaultLocation = location
}

type Provider struct {
	configuration Configuration
}

func (p Provider) GetMetadata() (map[wheretopark.ID]wheretopark.Metadata, error) {
	data, err := GetData()
	if err != nil {
		return nil, err
	}
	metadatas := make(map[wheretopark.ID]wheretopark.Metadata)
	for _, vendor := range data.Result.Parks {
		id := wheretopark.CoordinateToID(vendor.Latitude, vendor.Longitude)
		configuration, exists := p.configuration.ParkingLots[id]
		if !exists {
			log.Warn().Str("name", vendor.Name).Msg("missing configuration")
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
			if maxDimensions.Width == nil || widthInCm < *maxDimensions.Width {
				maxDimensions.Width = &widthInCm
			}
			if maxDimensions.Length == nil || lengthInCm < *maxDimensions.Length {
				maxDimensions.Length = &lengthInCm
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
			Name:      vendor.Name,
			Address:   configuration.Address,
			Geometry:  *geojson.NewPointGeometry([]float64{vendor.Longitude, vendor.Latitude}),
			Resources: configuration.Resources,
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCarElectric: totalPlaces.Electric,
				wheretopark.SpotTypeCar:         totalPlaces.Standard,
				wheretopark.SpotTypeCarDisabled: totalPlaces.Disabled,
			},
			MaxDimensions:  maxDimensions,
			Features:       configuration.Features,
			PaymentMethods: configuration.PaymentMethods,
			Comment:        configuration.Comment,
			Currency:       "PLN",
			Timezone:       "Europe/Warsaw",
			Rules:          configuration.Rules,
		}
		metadatas[id] = metadata
	}
	return metadatas, nil
}

func (p Provider) GetState() (map[wheretopark.ID]wheretopark.State, error) {
	data, err := GetData()
	if err != nil {
		return nil, err
	}
	lastUpdate, err := time.ParseInLocation("2006-01-02T15:04:05", data.Result.Timestamp, defaultLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to parse last update time: %w", err)
	}
	states := make(map[wheretopark.ID]wheretopark.State)

	for _, vendor := range data.Result.Parks {
		id := wheretopark.CoordinateToID(vendor.Latitude, vendor.Longitude)

		state := wheretopark.State{
			LastUpdated: lastUpdate.In(defaultLocation).Format(time.RFC3339),
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

func NewProvider(configurationPath *string) (wheretopark.Provider, error) {
	var configuration Configuration
	if configurationPath == nil {
		configuration = DefaultConfiguration
	} else {
		newConfiguration, err := LoadConfiguration(*configurationPath)
		if err != nil {
			return nil, err
		}
		configuration = *newConfiguration
	}
	return Provider{
		configuration: configuration,
	}, nil

}
