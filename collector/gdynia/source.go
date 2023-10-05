package gdynia

import (
	"context"
	"time"
	"wheretopark/collector/client"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/currency"
)

var (
	METADATA_URL = wheretopark.MustParseURL("http://api.zdiz.gdynia.pl/ri/rest/parkings")
	STATE_URL    = wheretopark.MustParseURL("http://api.zdiz.gdynia.pl/ri/rest/parking_places")
)

type Source struct{}

func (s Source) ParkingLots(ctx context.Context) (<-chan map[wheretopark.ID]wheretopark.ParkingLot, error) {
	vMetadata, err := client.Get[Metadata](METADATA_URL, nil)
	if err != nil {
		return nil, err
	}

	vState, err := client.Get[State](STATE_URL, nil)
	if err != nil {
		return nil, err
	}

	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot)
	for _, vMetadata := range vMetadata.ParkingLots {
		configuration, exists := configuration.ParkingLots[vMetadata.ID]
		if !exists {
			log.Ctx(ctx).
				Warn().
				Int("id", vMetadata.ID).
				Str("name", vMetadata.Name).
				Msg("missing configuration")
			continue
		}

		vStateIndex := StatePositionByID(*vState, vMetadata.ID)
		vState := (*vState)[vStateIndex]

		metadata := wheretopark.Metadata{
			LastUpdated: configuration.LastUpdated,
			Name:        vMetadata.Name,
			Address:     vMetadata.Address,
			Geometry:    vMetadata.Location,
			Resources:   configuration.Resources,
			TotalSpots: map[string]uint{
				wheretopark.SpotTypeCar: uint(vState.Capacity),
			},
			MaxDimensions:  configuration.MaxDimensions,
			Features:       configuration.Features,
			PaymentMethods: configuration.PaymentMethods,
			Comment:        configuration.Comment,
			Currency:       currency.PLN,
			Timezone:       defaultTimezone,
			Rules:          configuration.Rules,
		}
		lastUpdated, err := time.ParseInLocation("2006-01-02 15:04:05", vState.InsertTime, defaultTimezone)
		if err != nil {
			return nil, err
		}
		state := wheretopark.State{
			LastUpdated: lastUpdated,
			AvailableSpots: map[wheretopark.ID]uint{
				wheretopark.SpotTypeCar: uint(vState.FreePlaces),
			},
		}
		id := wheretopark.GeometryToID(vMetadata.Location)
		parkingLots[id] = wheretopark.ParkingLot{
			Metadata: metadata,
			State:    state,
		}
	}
	ch := make(chan map[wheretopark.ID]wheretopark.ParkingLot, 1)
	ch <- parkingLots
	close(ch)
	return ch, nil
}

func New() Source {
	return Source{}
}
