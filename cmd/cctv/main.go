package main

import (
	"fmt"
	"sync"
	wheretopark "wheretopark/go"
	"wheretopark/providers/cctv"

	"github.com/caarlos0/env/v8"
	"github.com/rs/zerolog/log"
)

type CctvServer struct {
	provider *cctv.Provider
	cache    *wheretopark.SingularCache
}

func (s CctvServer) GetParkingLotsByIdentifier(id wheretopark.ID) (map[wheretopark.ID]wheretopark.ParkingLot, error) {
	metadatas, err := s.provider.GetMetadatas()
	if err != nil {
		return nil, err
	}
	metadata, found := metadatas[id]
	if !found {
		return nil, fmt.Errorf("metadata for `%s` not found", id)
	}
	state := s.cache.GetState(id)
	log.Debug().Bool("cacheHit", state != nil).Str("id", id).Msg("cache response")
	if state == nil {
		state = s.provider.ProcessParkingLotByID(id)
		if state == nil {
			return nil, fmt.Errorf("failed to process parking lot `%s`", id)
		}
		if err := s.cache.SetState(id, *state); err != nil {
			log.Error().Err(err).Str("id", id).Msg("failed to set cache")
		}
	}
	parkingLot := wheretopark.ParkingLot{
		Metadata: metadata,
		State:    *state,
	}
	return map[wheretopark.ID]wheretopark.ParkingLot{id: parkingLot}, nil
}

func (s CctvServer) GetAllParkingLots() chan map[wheretopark.ID]wheretopark.ParkingLot {
	configurations := s.provider.GetConfiguredParkingLots()
	allParkingLots := make(chan map[wheretopark.ID]wheretopark.ParkingLot, len(configurations))

	var wg sync.WaitGroup
	for _, cfg := range configurations {
		wg.Add(1)
		id := wheretopark.GeometryToID(cfg.Metadata.Geometry)
		go func(cfg cctv.ParkingLot) {
			parkingLots, err := s.GetParkingLotsByIdentifier(id)
			if err != nil {
				log.Error().Err(err).Str("id", id).Msg("failed to get parking lots")
				return
			}
			allParkingLots <- parkingLots
			wg.Done()
		}(cfg)
	}
	go func() {
		wg.Wait()
		close(allParkingLots)

	}()
	return allParkingLots
}

type environment struct {
	Port          int              `env:"PORT" envDefault:"8080"`
	Configuration *string          `env:"CONFIGURATION"`
	Model         string           `env:"MODEL" envDefault:"$HOME/.local/share/wheretopark/cctv/model.onnx" envExpand:"true"`
	SavePath      *string          `env:"SAVE_PATH" envExpand:"true"`
	SaveItems     []cctv.SaveItem  `env:"SAVE_ITEMS" envSeparator:","`
	SaveIDs       []wheretopark.ID `env:"SAVE_IDS" envSeparator:","`
}

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}

	saver := cctv.NewSaver(environment.SavePath, environment.SaveItems, environment.SaveIDs)
	model := cctv.NewModel(environment.Model)
	defer model.Close()

	provider, err := cctv.NewProvider(environment.Configuration, model, saver)
	if err != nil {
		panic(err)
	}

	cache, err := wheretopark.NewSingularCache()
	if err != nil {
		log.Fatal().Err(err).Msg("create cache fail")
	}
	server := CctvServer{
		provider: provider,
		cache:    cache,
	}
	wheretopark.RunServer(server, uint(environment.Port))
}
