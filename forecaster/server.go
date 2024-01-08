package forecaster

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	wheretopark "wheretopark/go"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type Server struct {
	pycaster  Pycaster
	sequences map[wheretopark.ID]map[time.Time]uint
}

func NewServer(pycaster Pycaster, sequences map[wheretopark.ID]map[time.Time]uint) Server {
	return Server{
		pycaster,
		sequences,
	}
}

type Forecast struct {
	Predictions []Prediction `json:"predictions"`
}

func (s *Server) handleForecast(w http.ResponseWriter, r *http.Request) {
	parkingID := chi.URLParam(r, "identifier")
	dateOnly := chi.URLParam(r, "dateOnly")
	dateOnlyTime, err := time.Parse(time.DateOnly, dateOnly)
	if err != nil {
		log.Error().Err(err).Str("dateOnly", dateOnly).Msg("invalid dateOnly")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sequences, ok := s.sequences[wheretopark.ID(parkingID)]
	if !ok {
		log.Error().Err(err).Str("parkingID", parkingID).Msg("parkingID not found")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	start := time.Date(dateOnlyTime.Year(), dateOnlyTime.Month(), dateOnlyTime.Day(), 10, 0, 0, 0, time.UTC)
	end := time.Date(dateOnlyTime.Year(), dateOnlyTime.Month(), dateOnlyTime.Day(), 20, 0, 0, 0, time.UTC)
	predictions, err := s.pycaster.Predict(parkingID, sequences, start, end)
	if err != nil {
		log.Error().Err(err).Str("parkingID", parkingID).Msg("error predicting occupancy")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	forecast := Forecast{
		Predictions: predictions,
	}

	json, err := json.Marshal(forecast)
	if err != nil {
		log.Error().Err(err).Str("parkingID", parkingID).Msg("error marshaling predictions")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err = w.Write(json); err != nil {
		log.Error().Err(err).Str("parkingID", parkingID).Msg("error writing response")
	}
}

func (s *Server) Run(r chi.Router, port uint) error {
	log.Info().Msg(fmt.Sprintf("starting server on port %d", port))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func (s *Server) Router() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	// r.Mount("/debug", middleware.Profiler())

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))
	})
	// r.Get("/parking-lots/{identifier}", s.handleParkingLotsByIdentifier)
	r.Get("/predict/{identifier}/{dateOnly}", s.handleForecast)
	// r.Get("/predict/{identifier}/{dateOnly}", s.handleForecast)

	return r

}
