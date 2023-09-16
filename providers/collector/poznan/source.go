package poznan

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/client"

	"github.com/rs/zerolog/log"
)

const URL_FORMAT = "https://www.ztm.poznan.pl/pl/dla-deweloperow/getParkingFile?file=ZTM_ParkAndRide__%s.csv"

type Source struct{}

type container struct {
	sync.Mutex
	parkingLots map[wheretopark.ID]wheretopark.ParkingLot
}

func (s Source) ParkingLots(ctx context.Context) (map[wheretopark.ID]wheretopark.ParkingLot, error) {
	var wg sync.WaitGroup

	container := container{
		sync.Mutex{},
		make(map[wheretopark.ID]wheretopark.ParkingLot, len(configuration.ParkingLots)),
	}
	for name, metadata := range configuration.ParkingLots {
		wg.Add(1)
		go func(name string, metadata wheretopark.Metadata) {
			defer wg.Done()
			url, err := url.Parse(fmt.Sprintf(URL_FORMAT, name))
			if err != nil {
				log.Ctx(ctx).Err(err).Msg("invalid url")
				return
			}
			str, err := client.GetString(url, nil)
			if err != nil {
				log.Ctx(ctx).Err(err).Str("url", url.String()).Msg("failed to get string")
				return
			}
			data, err := parse(str)
			if err != nil {
				log.Ctx(ctx).Err(err).Str("data", str).Msg("failed to parse data")
				return
			}

			id := wheretopark.GeometryToID(metadata.Geometry)
			container.Lock()
			container.parkingLots[id] = wheretopark.ParkingLot{
				Metadata: metadata,
				State: wheretopark.State{
					LastUpdated: data.LastUpdated,
					AvailableSpots: map[wheretopark.SpotType]uint{
						wheretopark.SpotTypeCar: data.AvailableSpots,
					},
				},
			}
			container.Unlock()
		}(name, metadata)
	}
	wg.Wait()

	return container.parkingLots, nil
}

func New() Source {
	return Source{}
}
