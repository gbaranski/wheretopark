package poznan

import (
	wheretopark "wheretopark/go"
	"wheretopark/go/provider/simple"

	"github.com/rs/zerolog/log"
)

type Provider struct {
	mapping map[string]wheretopark.ID
}

func (p Provider) Name() string {
	return "poznan"
}

func (p Provider) Config() simple.Config {
	return simple.DEFAULT_CONFIG
}

func (p Provider) GetParkingLots() (map[wheretopark.ID]wheretopark.ParkingLot, error) {
	vendorData, err := GetData()
	if err != nil {
		return nil, err
	}
	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot)
	for name, data := range vendorData {
		metadata, exists := configuration.ParkingLots[name]
		if !exists {
			log.Warn().
				Str("name", name).
				Msg("missing configuration")
			continue
		}
		id := wheretopark.GeometryToID(metadata.Geometry)
		state := wheretopark.State{
			LastUpdated: data.LastUpdated.UTC(),
			AvailableSpots: map[string]uint{
				wheretopark.SpotTypeCar: data.AvailableSpots,
			},
		}
		parkingLots[id] = wheretopark.ParkingLot{
			Metadata: metadata,
			State:    state,
		}
	}
	return parkingLots, nil
}

func NewProvider() (simple.Provider, error) {
	return Provider{
		mapping: make(map[string]wheretopark.ID),
	}, nil
}
