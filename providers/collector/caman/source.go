package caman

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/client"

	"github.com/rs/zerolog/log"
)

type Source struct{}

var CAMAN_URL url.URL

func init() {
	strURL := os.Getenv("CAMAN_URL")
	if strURL == "" {
		log.Fatal().Msg("CAMAN_URL is not set")
	}
	url, err := url.Parse(strURL)
	if err != nil {
		log.Fatal().Err(err).Str("strURL", strURL).Msg("failed to parse CAMAN_URL")
	}
	CAMAN_URL = *url
}

func (s Source) UpdateCamanConfiguration() {

}

func (s Source) ParkingLots(ctx context.Context) (map[wheretopark.ID]wheretopark.ParkingLot, error) {
	strURL := fmt.Sprintf("%s/cameras/state", CAMAN_URL.String())
	url, err := url.Parse(strURL)
	if err != nil {
		return nil, err
	}
	vendor, err := client.Get[map[CameraID]CameraState](url, nil)
	if err != nil {
		return nil, err
	}

	// map of parking lot id to map of camera index to camera state
	cameras := make(map[wheretopark.ID][]CameraState)
	for cameraID, state := range *vendor {
		parts := strings.Split(cameraID, "_")
		parkingotID := parts[0]
		cameraIndex, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, err
		}
		cameras[parkingotID][cameraIndex] = state
	}

	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot, len(cameras))
	for id, cameras := range cameras {
		state := wheretopark.State{}
		for _, camera := range cameras {
			if state.LastUpdated.Before(camera.LastUpdated) {
				state.LastUpdated = camera.LastUpdated
			}
			state.AvailableSpots[wheretopark.SpotTypeCar] += camera.AvailableSpots
		}
		parkingLots[id] = wheretopark.ParkingLot{
			Metadata: DefaultConfiguration.ParkingLots[id].Metadata,
			State:    state,
		}
	}

	return parkingLots, nil
}

func New() Source {
	return Source{}
}
