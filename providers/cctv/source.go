package cctv

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

type Source struct {
	id       wheretopark.ID
	metadata wheretopark.Metadata
	cameras  []ParkingLotCamera
	model    *Model
	saver    Saver
}

func (s *Source) Metadata(context.Context) (map[wheretopark.ID]wheretopark.Metadata, error) {
	return map[wheretopark.ID]wheretopark.Metadata{
		s.id: s.metadata,
	}, nil
}

func (s *Source) State(ctx context.Context) (map[wheretopark.ID]wheretopark.State, error) {
	var availableSpots uint32
	var wg sync.WaitGroup
	captureTime := time.Now()
	for id, camera := range s.cameras {
		wg.Add(1)
		log.Ctx(ctx).
			Info().
			Int("camera", id).
			Msg("processing parking lot camera")
		go func(id int, camera ParkingLotCamera) {
			defer wg.Done()
			camAvailableSpots, err := s.processCamera(camera)
			if err != nil {
				log.Error().Err(err).Int("id", id).Msg("processing camera fail")
			}
			atomic.AddUint32(&availableSpots, uint32(camAvailableSpots))
		}(id, camera)
	}
	wg.Wait()
	state := wheretopark.State{
		LastUpdated: captureTime,
		AvailableSpots: map[wheretopark.SpotType]uint{
			wheretopark.SpotTypeCar: uint(availableSpots),
		},
	}
	return map[wheretopark.ID]wheretopark.State{
		s.id: state,
	}, nil
}

func (s *Source) processCamera(camera ParkingLotCamera) (uint, error) {
	img, err := GetImageFromCamera(camera.URL)
	if err != nil {
		return 0, fmt.Errorf("unable to get image from camera: %v", err)
	}
	defer img.Close()

	spotImages := ExtractSpots(img, camera.Spots)
	defer func() {
		for _, spotImage := range spotImages {
			spotImage.Close()
		}
	}()
	predictions := s.model.PredictMany(spotImages)
	availableSpots := 0
	for _, prediction := range predictions {
		if prediction > 0.5 {
			availableSpots += 1
		}
	}

	// err = p.saver.SavePredictions(img, id, cameraID, captureTime, camera.Spots, predictions)
	// if err != nil {
	// 	log.Error().
	// 		Str("name", parkingLot.Name).
	// 		Int("camera", cameraID).
	// 		Msg("unable to save predictions")
	// }
	return uint(availableSpots), nil
}
