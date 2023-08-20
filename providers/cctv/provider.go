package cctv

import (
	"sync"
	"sync/atomic"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/go/provider/sequential"

	"github.com/rs/zerolog/log"
	"gocv.io/x/gocv"
)

type Provider struct {
	configuration Configuration
	model         *Model
	modelMutex    sync.Mutex
	saver         Saver
}

func (p *Provider) Name() string {
	return "cctv"
}

func (p *Provider) Config() sequential.Config {
	return sequential.DEFAULT_CONFIG
}

func (p *Provider) GetMetadatas() (map[wheretopark.ID]wheretopark.Metadata, error) {
	metadatas := make(map[wheretopark.ID]wheretopark.Metadata, len(p.configuration.ParkingLots))

	for _, parkingLot := range p.configuration.ParkingLots {
		id := wheretopark.GeometryToID(parkingLot.Geometry)
		metadatas[id] = parkingLot.Metadata
	}

	return metadatas, nil
}

func (p *Provider) GetConfiguredParkingLots() []ParkingLot {
	return p.configuration.ParkingLots
}

func (p *Provider) ProcessCamera(parkingLot ParkingLot, cameraID int, camera ParkingLotCamera) (uint, error) {
	id := wheretopark.GeometryToID(parkingLot.Geometry)
	log.Info().
		Str("name", parkingLot.Name).
		Int("camera", cameraID).
		Msg("processing parking lot")
	img := gocv.NewMat()
	defer img.Close()
	stream, err := gocv.OpenVideoCapture(camera.URL)
	if err != nil {
		return 0, err
	}
	defer stream.Close()
	captureTime := time.Now()
	if ok := stream.Read(&img); !ok {
		log.Error().
			Str("name", parkingLot.Name).
			Int("camera", cameraID).
			Msg("unable to read stream")
		return 0, nil
	}
	spotImages := ExtractSpots(img, camera.Spots)
	defer func() {
		for _, spotImage := range spotImages {
			spotImage.Close()
		}
	}()
	p.modelMutex.Lock()
	predictions := p.model.PredictMany(spotImages)
	p.modelMutex.Unlock()
	availableSpots := 0
	for _, prediction := range predictions {
		if prediction > 0.5 {
			availableSpots += 1
		}
	}

	err = p.saver.SavePredictions(img, id, cameraID, captureTime, camera.Spots, predictions)
	if err != nil {
		log.Error().
			Str("name", parkingLot.Name).
			Int("camera", cameraID).
			Msg("unable to save predictions")
	}
	return uint(availableSpots), nil
}

func (p *Provider) ProcessParkingLot(parkingLot ParkingLot) wheretopark.State {
	var availableSpots uint32
	var wg sync.WaitGroup
	for k, camera := range parkingLot.Cameras {
		wg.Add(1)
		go func(k int, camera ParkingLotCamera) {
			defer wg.Done()
			onCameraAvailableSpots, err := p.ProcessCamera(parkingLot, k+1, camera)
			if err != nil {
				log.Error().
					Str("name", parkingLot.Name).
					Int("camera", k+1).
					Err(err).
					Send()
			}
			atomic.AddUint32(&availableSpots, uint32(onCameraAvailableSpots))
		}(k, camera)
	}
	wg.Wait()
	return wheretopark.State{
		LastUpdated: time.Now(),
		AvailableSpots: map[string]uint{
			"CAR": uint(availableSpots),
		},
	}
}

// Returns nil if parking lot configuration could not be found
func (p *Provider) ProcessParkingLotByID(id wheretopark.ID) *wheretopark.State {
	for _, parkingLot := range p.configuration.ParkingLots {
		if id == wheretopark.GeometryToID(parkingLot.Geometry) {
			state := p.ProcessParkingLot(parkingLot)
			return &state
		}
	}
	return nil
}

func (p *Provider) GetStates() (map[wheretopark.ID]wheretopark.State, error) {
	states := make(map[wheretopark.ID]wheretopark.State)
	statesMutex := sync.RWMutex{}
	var wg sync.WaitGroup

	for _, parkingLot := range p.configuration.ParkingLots {
		wg.Add(1)
		go func(parkingLot ParkingLot) {
			defer wg.Done()
			id := wheretopark.GeometryToID(parkingLot.Geometry)
			state := p.ProcessParkingLot(parkingLot)
			statesMutex.Lock()
			states[id] = state
			statesMutex.Unlock()
		}(parkingLot)
	}
	wg.Wait()
	return states, nil
}

func NewProvider(configurationPath *string, model *Model, saver Saver) (*Provider, error) {
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

	return &Provider{
		configuration: configuration,
		model:         model,
		modelMutex:    sync.Mutex{},
		saver:         saver,
	}, nil

}
