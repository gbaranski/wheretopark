package wheretopark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

type Server struct {
	cache  *Cache
	source Source
}

func NewServer(cache *Cache, source Source) *Server {
	return &Server{
		cache:  cache,
		source: source,
	}
}

func (s *Server) handleParkingLots(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	var mu sync.Mutex

	send := func(parkingLots map[ID]ParkingLot) {
		log.Info().Int("n", len(parkingLots)).Msg("sending parkings lots")
		json, err := json.Marshal(parkingLots)
		if err != nil {
			log.Error().Err(err).Msg("failed to marshal parking lots")
			return
		}
		mu.Lock()
		defer mu.Unlock()
		_, err = w.Write([]byte(string(json) + "\r\n"))
		if err != nil {
			log.Error().Err(err).Msg("failed to write response")
			return
		}
		log.Debug().Int("n", len(parkingLots)).Msg("sending parking lots to client")
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}

	parkingLots := s.cache.GetAllParkingLots()
	if parkingLots != nil {
		send(parkingLots)
		return
	}

	ctx := log.With().Logger().WithContext(context.TODO())
	ch, err := s.source.ParkingLots(ctx)
	if err != nil {
		log.Error().Err(err).Msg("get parking lots failure")
		return
	}
	for parkingLots := range ch {
		send(parkingLots)
		for id, parkingLot := range parkingLots {
			err := s.cache.SetParkingLot(id, &parkingLot)
			if err != nil {
				log.Error().Err(err).Str("id", id).Msg("set to cache failure")
			}
		}
	}
}

// func (s *Server) handleParkingLotsByIdentifier(w http.ResponseWriter, r *http.Request) {
// 	identifier := chi.URLParam(r, "identifier")
// 	send := func(parkingLots map[wheretopark.ID]wheretopark.ParkingLot) {
// 		log.Info().Int("n", len(parkingLots)).Str("source", identifier).Msg("sending parkings lots")
// 		json, err := json.Marshal(parkingLots)
// 		if err != nil {
// 			log.Error().Err(err).Msg("failed to marshal parking lots")
// 			return
// 		}
// 		_, err = w.Write([]byte(string(json) + "\r\n"))
// 		if err != nil {
// 			log.Error().Err(err).Msg("failed to write response")
// 			return
// 		}
// 		log.Debug().Int("n", len(parkingLots)).Msg("sending parking lots to client")
// 		if f, ok := w.(http.Flusher); ok {
// 			f.Flush()
// 		}
// 	}

// 	parkingLots := s.cache.GetParkingLots(identifier)
// 	if parkingLots != nil {
// 		send(parkingLots)
// 		return
// 	}

// 	source, ok := s.sources[identifier]
// 	if !ok {
// 		w.WriteHeader(http.StatusBadRequest)
// 		w.Write([]byte(fmt.Sprintf("unknown identifier: %s", identifier)))
// 		return
// 	}

// 	ctx := log.With().Str("source", identifier).Logger().WithContext(context.TODO())
// 	ch, err := source.ParkingLots(ctx)
// 	if err != nil {
// 		log.Error().Err(err).Str("identifier", identifier).Msg("get parking lots from source failure")
// 		w.WriteHeader(http.StatusInternalServerError)
// 		return
// 	}

// 	for parkingLots := range ch {
// 		err := s.cache.UpdateParkingLots(identifier, parkingLots)
// 		if err != nil {
// 			log.Error().Err(err).Str("identifier", identifier).Msg("update cache failure")
// 		}
// 		send(parkingLots)
// 	}
// }

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
	r.Get("/parking-lots", s.handleParkingLots)

	return r

}
