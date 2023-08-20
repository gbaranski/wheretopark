package wheretopark

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
)

type Server interface {
	GetAllParkingLots() chan map[ID]ParkingLot
	GetParkingLotsByIdentifier(identifier string) (map[ID]ParkingLot, error)
}

func RunServer(s Server, port uint) {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/parking-lots", func(w http.ResponseWriter, r *http.Request) {
		stream := s.GetAllParkingLots()
		for parkingLots := range stream {
			render.JSON(w, r, parkingLots)
			log.Debug().Int("n", len(parkingLots)).Msg("sending parking lots to client")
		}
	})
	r.Get("/{identifier}/parking-lots", func(w http.ResponseWriter, r *http.Request) {
		identifier := chi.URLParam(r, "identifier")
		parkingLots, err := s.GetParkingLotsByIdentifier(identifier)
		if err != nil {
			log.Error().Err(err).Str("identifier", identifier).Msg("getParkingLots failure")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		json, err := json.Marshal(parkingLots)
		if err != nil {
			log.Error().Err(err).Str("identifier", identifier).Msg("JSON Marshall failure")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(json)
		if err != nil {
			log.Error().Err(err).Str("identifier", identifier).Msg("Write Response failure")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	})

	log.Info().Msg(fmt.Sprintf("starting server on port %d", port))
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), r); err != nil {
		log.Fatal().Err(err).Msg("server fail")
	}
}
