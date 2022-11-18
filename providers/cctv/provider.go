package cctv

import (
	"fmt"
	"time"
	wheretopark "wheretopark/go"

	"gocv.io/x/gocv"
)

type Provider struct {
	configuration Configuration
	model         *Model
	streams       []*gocv.VideoCapture
}

func (p Provider) GetMetadata() (map[wheretopark.ID]wheretopark.Metadata, error) {
	metadatas := make(map[wheretopark.ID]wheretopark.Metadata, len(p.configuration.ParkingLots))

	for _, parkingLot := range p.configuration.ParkingLots {
		id := wheretopark.GeometryToID(parkingLot.Geometry)
		metadatas[id] = parkingLot.Metadata
	}

	fmt.Printf("obtained %d metadatas\n", len(metadatas))
	return metadatas, nil
}

func (p Provider) GetState() (map[wheretopark.ID]wheretopark.State, error) {
	states := make(map[wheretopark.ID]wheretopark.State)

	img := gocv.NewMat()
	for i, parkingLot := range p.configuration.ParkingLots {
		stream := p.streams[i]
		if ok := stream.Read(&img); !ok {
			return nil, fmt.Errorf("cannot read stream of %s", parkingLot.Name)
		}
		predictions := parkingLot.RunPredictions(p.model, img)
		availableSpots := 0
		for _, prediction := range predictions {
			if prediction > 0.5 {
				availableSpots += 1
			}
		}
		id := wheretopark.GeometryToID(parkingLot.Geometry)
		state := wheretopark.State{
			LastUpdated: time.Now().Format(time.RFC3339),
			AvailableSpots: map[string]uint{
				"CAR": uint(availableSpots),
			},
		}
		states[id] = state
	}
	fmt.Printf("obtained %d states\n", len(states))
	return states, nil
}

func (p Provider) Close() error {
	for _, stream := range p.streams {
		if err := stream.Close(); err != nil {
			return err
		}
	}
	return nil
}

func NewProvider(configurationPath *string, model *Model) (*Provider, error) {
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

	streams := make([]*gocv.VideoCapture, len(configuration.ParkingLots))
	for i, parkingLot := range configuration.ParkingLots {
		fmt.Printf("conencting to %s\n", parkingLot.CameraURL)
		stream, err := gocv.OpenVideoCapture(parkingLot.CameraURL)
		if err != nil {
			return nil, err
		}
		fmt.Printf("connected\n")
		streams[i] = stream
	}

	return &Provider{
		configuration: configuration,
		model:         model,
		streams:       streams,
	}, nil

}
