package main

import (
	"fmt"
	"sync"
	wheretopark "wheretopark/go"
	"wheretopark/go/provider"
	"wheretopark/go/provider/sequential"
	"wheretopark/go/provider/simple"
	"wheretopark/providers/collector/gdansk"
	"wheretopark/providers/collector/gdynia"
	"wheretopark/providers/collector/glasgow"
	"wheretopark/providers/collector/poznan"
	"wheretopark/providers/collector/warsaw"

	"github.com/caarlos0/env/v8"
	"github.com/rs/zerolog/log"
)

func mustCreateProvider[T provider.Common](createFn func() (T, error)) T {
	provider, err := createFn()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create provider")
	}
	return provider
}

type CollectorServer struct {
	providers []provider.Common
	cache     *wheretopark.PluralCache
}

func (s CollectorServer) GetParkingLotsByIdentifier(name string) (map[wheretopark.ID]wheretopark.ParkingLot, error) {
	metadatas := s.cache.GetMetadatas(name)
	states := s.cache.GetStates(name)

	cacheHit := metadatas != nil && states != nil
	log.Debug().Bool("cacheHit", cacheHit).Str("provider", name).Msg("cache response")
	if cacheHit {
		parkingLots, err := wheretopark.JoinMetadatasAndStates(metadatas, states)
		if err != nil {
			return nil, fmt.Errorf("failed to join metadatas and states due to %e", err)
		}
		return parkingLots, nil
	}
	var provider provider.Common
	for _, p := range s.providers {
		if p.Name() == name {
			provider = p
		}
	}

	var parkingLots map[wheretopark.ID]wheretopark.ParkingLot
	var err error

	if simple, ok := provider.(simple.Provider); ok {
		parkingLots, err = simple.GetParkingLots()
		if err != nil {
			return nil, fmt.Errorf("failed to get parking lots due to %e", err)
		}
	} else if sequential, ok := provider.(sequential.Provider); ok {
		metadatas, err := sequential.GetMetadatas()
		if err != nil {
			return nil, fmt.Errorf("failed to get metadatas due to %e", err)
		}
		states, err := sequential.GetStates()
		if err != nil {
			return nil, fmt.Errorf("failed to get states due to %e", err)
		}
		parkingLots, err = wheretopark.JoinMetadatasAndStates(metadatas, states)
		if err != nil {
			return nil, fmt.Errorf("failed to join metadatas and states due to %e", err)
		}
	}

	err = s.cache.SetParkingLots(name, parkingLots)
	if err != nil {
		log.Error().Err(err).Str("provider", name).Msg("failed to update parking lots cache")
	}
	return parkingLots, nil
}

func (s CollectorServer) GetAllParkingLots() chan map[wheretopark.ID]wheretopark.ParkingLot {
	allParkingLots := make(chan map[wheretopark.ID]wheretopark.ParkingLot, len(s.providers))

	var wg sync.WaitGroup
	for _, provider := range s.providers {
		wg.Add(1)
		providerName := provider.Name()
		go func() {
			parkingLots, err := s.GetParkingLotsByIdentifier(providerName)
			if err != nil {
				log.Error().Err(err).Str("provider", providerName).Msg("failed to get parking lots")
				return
			}
			log.Debug().Str("provider", providerName).Msg("got parking lots")
			allParkingLots <- parkingLots
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(allParkingLots)

	}()
	return allParkingLots
}

type environment struct {
	Port uint `env:"PORT" envDefault:"8080"`
}

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}

	providers := []provider.Common{
		mustCreateProvider(gdansk.NewProvider),
		mustCreateProvider(gdynia.NewProvider),
		mustCreateProvider(warsaw.NewProvider),
		mustCreateProvider(poznan.NewProvider),
		mustCreateProvider(glasgow.NewProvider),
	}

	cache, err := wheretopark.NewPluralCache()
	if err != nil {
		log.Fatal().Err(err).Msg("create cache fail")
	}

	server := CollectorServer{
		providers: providers,
		cache:     cache,
	}
	wheretopark.RunServer(server, uint(environment.Port))
}
