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

type Provider interface {
	Sources() map[string]Source
}

func RunProvider(p Provider, port uint) error {
	router, err := GetProviderRouter(p)
	if err != nil {
		return err
	}
	return RunRouter(router, port)
}

func RunRouter(r chi.Router, port uint) error {
	log.Info().Msg(fmt.Sprintf("starting server on port %d", port))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}

func GetProviderRouter(p Provider) (chi.Router, error) {
	cache, err := NewCache()
	if err != nil {
		return nil, fmt.Errorf("failed to create cache: %w", err)
	}

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
	r.Get("/parking-lots/{identifier}", func(w http.ResponseWriter, r *http.Request) {
		identifier := chi.URLParam(r, "identifier")
		source := p.Sources()[identifier]
		ctx := log.With().Str("source", identifier).Logger().WithContext(context.TODO())
		parkingLots, err := source.ParkingLots(ctx)
		if err != nil {
			log.Error().Err(err).Str("identifier", identifier).Msg("get parking lots from source failure")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json, err := json.Marshal(parkingLots)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Error().Err(err).Str("identifier", identifier).Msg("failed to marshal parking lots")
			return
		}
		_, err = w.Write(json)
		if err != nil {
			log.Error().Err(err).Str("identifier", identifier).Msg("failed to write response")
			return
		}
	})
	r.Get("/parking-lots", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		var mu sync.Mutex
		var wg sync.WaitGroup
		for identifier, source := range p.Sources() {
			wg.Add(1)
			go func(identifier string, source Source) {
				defer wg.Done()
				ctx := log.With().Str("source", identifier).Logger().WithContext(context.TODO())
				parkingLots, err := cache.GetParkingLotsOrUpdate(identifier, func() (map[ID]ParkingLot, error) {
					return source.ParkingLots(ctx)
				})
				if err != nil {
					log.Error().Err(err).Str("identifier", identifier).Msg("get parking lots faliure")
					return
				}
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
			}(identifier, source)
		}
		wg.Wait()
	})

	return r, nil
}
