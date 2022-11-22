package cctv

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	wheretopark "wheretopark/go"

	"gocv.io/x/gocv"
)

type Provider struct {
	configuration Configuration
	model         *Model
	streams       []*gocv.VideoCapture
	savePath      *string
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
		captureTime := time.Now()
		spotImages := ExtractSpots(img, parkingLot.Spots)
		predictions := p.model.PredictMany(spotImages)
		availableSpots := 0
		for _, prediction := range predictions {
			if prediction > 0.5 {
				availableSpots += 1
			}
		}

		if p.savePath != nil {
			// Save files for history and continious training
			{
				basePath := fmt.Sprintf("%s/%s/%s", *p.savePath, parkingLot.Name, captureTime.Format(time.RFC3339))
				err := os.MkdirAll(basePath, os.ModePerm)
				if err != nil {
					log.Println(err)
				}
				// Write raw image, without any tweaks
				{
					gocv.IMWrite(fmt.Sprintf("%s/raw.jpg", basePath), img)
				}
				// Write file with results
				{
					spotResults := make([]SpotResult, len(predictions))
					for i, prediction := range predictions {
						spot := parkingLot.Spots[i]
						spotResults[i] = SpotResult{
							Prediction: prediction,
							Points:     spot.Points,
						}
					}
					result := Result{
						Spots: spotResults,
					}
					resultJSON, err := json.Marshal(result)
					if err != nil {
						log.Fatalf("cannot marshall result: %v", err)
					}
					err = os.WriteFile(fmt.Sprintf("%s/result.json", basePath), resultJSON, 0644)
					if err != nil {
						log.Fatalf("cannot write result.json: %v", err)
					}
				}
				// Write image with drawn spots and predicted values
				{
					for i, prediction := range predictions {
						spot := parkingLot.Spots[i]
						VisualizeSpotPrediction(&img, spot, prediction)
					}
					gocv.IMWrite(fmt.Sprintf("%s/visualization.jpg", basePath), img)
				}
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
		savePath:      savePath,
	}, nil

}
