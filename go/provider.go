package wheretopark

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/zerolog/log"
)

type Provider interface {
	Sources() map[string]Source
}

func RunProvider(p Provider, port uint) error {
	cache, err := NewCache()
	if err != nil {
		return fmt.Errorf("failed to create cache: %w", err)
	}

	r := httprouter.New()
	r.GET("/parking-lots/:identifier", func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		identifier := ps.ByName("identifier")
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
	r.GET("/parking-lots", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
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
				_, err = w.Write([]byte(string(json) + "\r\n"))
				if err != nil {
					log.Error().Err(err).Msg("failed to write response")
					return
				}
				mu.Unlock()
				log.Debug().Int("n", len(parkingLots)).Msg("sending parking lots to client")
				if f, ok := w.(http.Flusher); ok {
					f.Flush()
				}
			}(identifier, source)
		}
		wg.Wait()
	})

	log.Info().Msg(fmt.Sprintf("starting server on port %d", port))
	return http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
