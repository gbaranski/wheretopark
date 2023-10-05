package cctv

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
	wheretopark "wheretopark/go"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"gocv.io/x/gocv"
)

type Source struct {
	configuration Configuration
	model         Model
	saver         Saver
}

func (s Source) stateOf(ctx context.Context, parkingLot ConfiguredParkingLot) (wheretopark.State, error) {
	availableSpotsPtr := make(map[wheretopark.SpotType]*uint32, len(wheretopark.SpotTypes))
	for _, spotType := range wheretopark.SpotTypes {
		availableSpotsPtr[spotType] = new(uint32)
	}

	var wg sync.WaitGroup
	captureTime := time.Now()
	for id, camera := range parkingLot.Cameras {
		wg.Add(1)
		log.Ctx(ctx).
			Info().
			Int("camera", id).
			Msg("processing parking lot camera")
		go func(id int, camera Camera) {
			defer wg.Done()
			camAvailableSpots, err := s.processCamera(camera)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Int("id", id).Msg("processing camera fail")
				return
			}
			for spotType, count := range camAvailableSpots {
				atomic.AddUint32(availableSpotsPtr[spotType], uint32(count))
			}
		}(id, camera)
	}
	wg.Wait()
	log.Ctx(ctx).
		Info().
		Str("duration", time.Since(captureTime).String()).
		Msg("finished processing cameras")

	availableSpots := make(map[wheretopark.SpotType]uint)
	// TODO: Make clients to comply with missing spot type car
	availableSpots[wheretopark.SpotTypeCar] = 0
	for spotType, countPtr := range availableSpotsPtr {
		if *countPtr > 0 {
			availableSpots[spotType] = uint(*countPtr)
		}
	}
	state := wheretopark.State{
		LastUpdated:    captureTime,
		AvailableSpots: availableSpots,
	}
	return state, nil

}

func (s Source) processCamera(camera Camera) (map[wheretopark.SpotType]uint, error) {
	img, err := GetImageFromCamera(camera.URL)
	if err != nil {
		return nil, fmt.Errorf("unable to get image from camera: %v", err)
	}
	defer img.Close()

	spotImages := ExtractSpots(img, camera.Spots)
	defer func() {
		for _, spotImage := range spotImages {
			spotImage.Close()
		}
	}()
	predictions := s.model.PredictMany(spotImages)

	availableSpots := map[wheretopark.SpotType]uint{}
	for i, prediction := range predictions {
		if IsVacant(prediction) {
			spotType := camera.Spots[i].Type
			availableSpots[spotType]++
		}
	}

	// err = p.saver.SavePredictions(img, id, cameraID, captureTime, camera.Spots, predictions)
	// if err != nil {
	// 	log.Error().
	// 		Str("name", parkingLot.Name).
	// 		Int("camera", cameraID).
	// 		Msg("unable to save predictions")
	// }
	return availableSpots, nil
}

func (s Source) ParkingLots(ctx context.Context) (<-chan map[wheretopark.ID]wheretopark.ParkingLot, error) {
	var wg sync.WaitGroup

	ch := make(chan map[wheretopark.ID]wheretopark.ParkingLot, len(s.configuration.ParkingLots))
	for id, cfg := range s.configuration.ParkingLots {
		wg.Add(1)
		go func(id wheretopark.ID, cfg ConfiguredParkingLot) {
			defer wg.Done()
			ctx := log.With().Str("id", id).Logger().WithContext(ctx)
			state, err := s.stateOf(ctx, cfg)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("stateOf fail")
				return
			}
			ch <- map[wheretopark.ID]wheretopark.ParkingLot{
				id: {
					Metadata: cfg.Metadata,
					State:    state,
				},
			}
		}(id, cfg)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch, nil
}

func (s Source) HandleVisualizeCamera(w http.ResponseWriter, r *http.Request) {
	id := wheretopark.ID(chi.URLParam(r, "id"))
	parkingLot, exists := s.configuration.ParkingLots[id]
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("parking lot not found"))
		return
	}

	cameraID, err := strconv.Atoi(chi.URLParam(r, "camera"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid camera ID"))
		return
	}
	if cameraID > len(parkingLot.Cameras)-1 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("camera not found"))
		return
	}
	camera := parkingLot.Cameras[cameraID]
	img, err := generateVisualizationFor(s.model, &camera)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("generate visualization failure: %s", err)))
		return
	}
	defer img.Close()

	buf, err := gocv.IMEncode(".jpg", img)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Err(err).Msg("encode image failure")
		return
	}
	defer buf.Close()

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(buf.GetBytes())
}

func (s Source) ConfigureRoutes(r *chi.Mux) error {
	r.Get("/visualize/{id}/{camera}", s.HandleVisualizeCamera)
	return nil
}

func New(environment Environment) Source {
	model := NewModel(environment.ModelPath)
	saver := NewSaver(environment.SavePath, environment.SaveItems, environment.SaveIDs)
	return Source{
		configuration: DefaultConfiguration,
		model:         model,
		saver:         saver,
	}
}
