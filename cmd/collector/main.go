package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	wheretopark "wheretopark/go"
	"wheretopark/go/provider"
	"wheretopark/go/provider/sequential"
	"wheretopark/go/provider/simple"
	"wheretopark/providers/collector/gdansk"
	"wheretopark/providers/collector/gdynia"
	"wheretopark/providers/collector/glasgow"
	"wheretopark/providers/collector/poznan"
	"wheretopark/providers/collector/warsaw"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func mustCreateProvider[T provider.Common](createFn func() (T, error)) T {
	provider, err := createFn()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create provider")
	}
	return provider
}

func GetParkingLots(providers []provider.Common, cache *wheretopark.CacheProvider, name string) (map[wheretopark.ID]wheretopark.ParkingLot, error) {
	metadatas, mFound := cache.GetMetadatas(name)
	states, sFound := cache.GetStates(name)

	if mFound && sFound {
		parkingLots, err := wheretopark.JoinMetadatasAndStates(metadatas, states)
		if err != nil {
			return nil, fmt.Errorf("failed to join metadatas and states due to %e", err)
		}
		return parkingLots, nil
	}
	var provider provider.Common
	for _, p := range providers {
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
		return parkingLots, nil
	}

	err = cache.SetParkingLots(name, parkingLots)
	if err != nil {
		log.Error().Err(err).Str("provider", name).Msg("failed to update parking lots cache")
	}
	return parkingLots, nil
}

func handleGetParkingLots(
	providers []provider.Common,
	cache *wheretopark.CacheProvider,
	w http.ResponseWriter,
	r *http.Request,
) {
	providerName := chi.URLParam(r, "provider")
	parkingLots, err := GetParkingLots(providers, cache, providerName)
	if err != nil {
		log.Error().Err(err).Str("provider", providerName).Msg("failed to get parking lots")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	json, err := json.Marshal(parkingLots)
	if err != nil {
		log.Error().Err(err).Str("provider", providerName).Msg("failed to marshal response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(json)
	if err != nil {
		log.Error().Err(err).Str("provider", providerName).Msg("failed to write response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	providers := []provider.Common{
		mustCreateProvider(gdansk.NewProvider),
		mustCreateProvider(gdynia.NewProvider),
		mustCreateProvider(warsaw.NewProvider),
		mustCreateProvider(poznan.NewProvider),
		mustCreateProvider(glasgow.NewProvider),
	}

	cache, err := wheretopark.NewCacheProvider()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create cache")
	}
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/{provider}/parking-lots", func(w http.ResponseWriter, r *http.Request) {
		handleGetParkingLots(providers, cache, w, r)
	})
	port := 8080
	log.Info().Msg(fmt.Sprintf("starting server on port %d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Fatal().Err(err).Msg("server fail")
	}
}
