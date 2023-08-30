package wheretopark

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type Server interface {
	GetAllParkingLots() chan map[ID]ParkingLot
	GetParkingLotsByIdentifier(identifier string) (map[ID]ParkingLot, error)
}

func RunServer(s Server, port uint) {
	r := gin.Default()
	r.Use(logger.SetLogger())
	r.GET("/parking-lots", func(c *gin.Context) {
		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Header().Set("Content-Type", "application/json")
		stream := s.GetAllParkingLots()
		for parkingLots := range stream {
			json, err := json.Marshal(parkingLots)
			if err != nil {
				log.Error().Err(err).Msg("failed to marshal parking lots")
				break
			}
			_, err = c.Writer.WriteString(string(json) + "\r\n")
			if err != nil {
				log.Error().Err(err).Msg("failed to write response")
				break
			}
			log.Debug().Int("n", len(parkingLots)).Msg("sending parking lots to client")
			c.Writer.Flush()

			time.Sleep(1 * time.Second)
		}
	})
	r.GET("/:identifier/parking-lots", func(c *gin.Context) {
		identifier := c.Param("identifier")
		parkingLots, err := s.GetParkingLotsByIdentifier(identifier)
		if err != nil {
			log.Error().Err(err).Str("identifier", identifier).Msg("getParkingLots failure")
			c.Status(http.StatusInternalServerError)
			return
		}

		c.JSON(http.StatusOK, parkingLots)
	})

	log.Info().Msg(fmt.Sprintf("starting server on port %d", port))
	if err := r.Run(fmt.Sprintf(":%d", port)); err != nil {
		log.Fatal().Err(err).Msg("server fail")
	}
}
