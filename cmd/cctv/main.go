package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	wheretopark "wheretopark/go"
	"wheretopark/go/provider/sequential"
	"wheretopark/providers/cctv"

	"github.com/caarlos0/env/v8"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type environment struct {
	Port          int              `env:"PORT" envDefault:"8080"`
	Configuration *string          `env:"CONFIGURATION"`
	Model         string           `env:"MODEL" envDefault:"$HOME/.local/share/wheretopark/cctv/model.onnx" envExpand:"true"`
	SavePath      *string          `env:"SAVE_PATH" envExpand:"true"`
	SaveItems     []cctv.SaveItem  `env:"SAVE_ITEMS" envSeparator:","`
	SaveIDs       []wheretopark.ID `env:"SAVE_IDS" envSeparator:","`
}

func GetParkingLots(provider sequential.Provider, cache *wheretopark.CacheProvider, name string) (map[wheretopark.ID]wheretopark.ParkingLot, error) {
	metadatas, mFound := cache.GetMetadatas(name)
	states, sFound := cache.GetStates(name)

	if mFound && sFound {
		parkingLots, err := wheretopark.JoinMetadatasAndStates(metadatas, states)
		if err != nil {
			return nil, fmt.Errorf("failed to join metadatas and states due to %e", err)
		}
		return parkingLots, nil
	}

	metadatas, err := provider.GetMetadatas()
	if err != nil {
		return nil, fmt.Errorf("failed to get metadatas due to %e", err)
	}
	states, err = provider.GetStates()
	if err != nil {
		return nil, fmt.Errorf("failed to get states due to %e", err)
	}
	parkingLots, err := wheretopark.JoinMetadatasAndStates(metadatas, states)
	if err != nil {
		return nil, fmt.Errorf("failed to join metadatas and states due to %e", err)
	}
	err = cache.SetParkingLots(name, parkingLots)
	if err != nil {
		log.Error().Err(err).Str("provider", name).Msg("failed to update parking lots cache")
	}
	return parkingLots, nil
}

func handleGetParkingLots(
	provider sequential.Provider,
	cache *wheretopark.CacheProvider,
	w http.ResponseWriter,
	r *http.Request,
) {
	providerName := chi.URLParam(r, "provider")
	parkingLots, err := GetParkingLots(provider, cache, providerName)
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
	cache, err := wheretopark.NewCacheProvider()
	if err != nil {
		log.Fatal().Err(err).Msg("create cache fail")
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Get("/parking-lots", func(w http.ResponseWriter, r *http.Request) {
		handleGetParkingLots(provider, cache, w, r)
	})
	log.Info().Msg(fmt.Sprintf("starting server on port %d", environment.Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", environment.Port), r); err != nil {
		log.Fatal().Err(err).Msg("server fail")
	}
}
