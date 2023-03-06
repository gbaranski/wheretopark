package poznan

import (
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
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
	mapping map[string]wheretopark.ID
}

func (p Provider) Name() string {
	return "poznan"
}

func (p Provider) GetMetadata() (map[wheretopark.ID]wheretopark.Metadata, error) {
	vendorData, err := GetData()
	if err != nil {
		return nil, err
	}
	metadatas := make(map[wheretopark.ID]wheretopark.Metadata)
	for name := range vendorData {
		metadata, exists := configuration.ParkingLots[name]
		if !exists {
			log.Warn().
				Str("name", name).
				Msg("missing configuration")
			continue
		}
		id := wheretopark.GeometryToID(metadata.Geometry)
		metadatas[id] = metadata
	}
	return metadatas, nil
}

func (p Provider) GetState() (map[wheretopark.ID]wheretopark.State, error) {
	vendorData, err := GetData()
	if err != nil {
		return nil, err
	}
	states := make(map[wheretopark.ID]wheretopark.State)
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
			LastUpdated: data.LastUpdated.UTC().Format(time.RFC3339),
			AvailableSpots: map[string]uint{
				wheretopark.SpotTypeCar: data.AvailableSpots,
			},
		}
		states[id] = state
	}
	return states, nil
}

func NewProvider() (wheretopark.Provider, error) {
	return Provider{
		mapping: make(map[string]wheretopark.ID),
	}, nil
}
