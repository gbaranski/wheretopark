package cctv

import (
	"fmt"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
	"gocv.io/x/gocv"
)

type Provider struct {
	configuration Configuration
	model         *Model
	streams       [][]*gocv.VideoCapture
	savePath      *string
}

func (p Provider) GetMetadata() (map[wheretopark.ID]wheretopark.Metadata, error) {
	metadatas := make(map[wheretopark.ID]wheretopark.Metadata, len(p.configuration.ParkingLots))

	for _, parkingLot := range p.configuration.ParkingLots {
		id := wheretopark.GeometryToID(parkingLot.Geometry)
		metadatas[id] = parkingLot.Metadata
	}

	return metadatas, nil
}

func (p Provider) GetState() (map[wheretopark.ID]wheretopark.State, error) {
	states := make(map[wheretopark.ID]wheretopark.State)

	img := gocv.NewMat()
	defer img.Close()
	for i, parkingLot := range p.configuration.ParkingLots {
		availableSpots := 0
		captureTime := time.Now()
		id := wheretopark.GeometryToID(parkingLot.Geometry)
		for k, camera := range parkingLot.Cameras {
			log.Info().
				Str("name", parkingLot.Name).
				Int("camera", k).
				Msg("processing parking lot")
			stream := p.streams[i][k]
			if ok := stream.Read(&img); !ok {
				log.Error().
					Str("name", parkingLot.Name).
					Int("camera", k).
					Msg("unable to read stream")
				continue
			}
			spotImages := ExtractSpots(img, camera.Spots)
			defer func() {
				for _, spotImage := range spotImages {
					spotImage.Close()
				}
			}()
			predictions := p.model.PredictMany(spotImages)
			for _, prediction := range predictions {
				if prediction > 0.5 {
					availableSpots += 1
				}
			}

			if p.savePath != nil {
				basePath := fmt.Sprintf("%s/%s/%02d", *p.savePath, id, k+1)
				err := SavePredictions(img, basePath, captureTime, camera.Spots, predictions)
				if err != nil {
					log.Error().
						Str("name", parkingLot.Name).
						Int("camera", k).
						Msg("unable to save predictions")
				}
			}
		}
		state := wheretopark.State{
			LastUpdated: time.Now().Format(time.RFC3339),
			AvailableSpots: map[string]uint{
				"CAR": uint(availableSpots),
			},
		}
		states[id] = state
	}
	return states, nil
}

func (p Provider) Close() error {
	for _, streams := range p.streams {
		for _, stream := range streams {
			if err := stream.Close(); err != nil {
				return err
			}
		}

	}
	return nil
}

func NewProvider(configurationPath *string, model *Model, savePath *string) (*Provider, error) {
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

	streams := make([][]*gocv.VideoCapture, len(configuration.ParkingLots))
	for i, parkingLot := range configuration.ParkingLots {
		streams[i] = make([]*gocv.VideoCapture, len(parkingLot.Cameras))
		for k, camera := range parkingLot.Cameras {
			log.Info().
				Str("url", camera.URL).
				Str("name", parkingLot.Name).
				Int("camera", k).
				Msg("connecting")
			stream, err := gocv.OpenVideoCapture(camera.URL)
			if err != nil {
				return nil, err
			}
			log.Info().
				Str("url", camera.URL).
				Str("name", parkingLot.Name).
				Int("camera", k).
				Msg("connected")
			streams[i][k] = stream
		}
	}

	return &Provider{
		configuration: configuration,
		model:         model,
		streams:       streams,
		savePath:      savePath,
	}, nil

}
