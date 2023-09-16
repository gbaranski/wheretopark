package wheretopark

import (
	"context"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
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
	r := gin.Default()
	r.ForwardedByClientIP = true
	r.Use(logger.SetLogger())
	r.GET("/parking-lots", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "application/json")
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
				log.Debug().Int("n", len(parkingLots)).Msg("sending parking lots to client")
				c.JSON(http.StatusOK, parkingLots)
			}(identifier, source)
		}
		wg.Wait()
	})
	r.GET("/:identifier/parking-lots", func(c *gin.Context) {
		identifier := c.Param("identifier")
		source := p.Sources()[identifier]
		ctx := log.With().Str("source", identifier).Logger().WithContext(context.TODO())
		parkingLots, err := source.ParkingLots(ctx)
		if err != nil {
			log.Error().Err(err).Str("identifier", identifier).Msg("get parking lots from source failure")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, parkingLots)
	})

	log.Info().Msg(fmt.Sprintf("starting server on port %d", port))
	return r.Run(fmt.Sprintf(":%d", port))
}
