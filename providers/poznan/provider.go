package poznan

import (
	"context"
	"fmt"
	"net/url"
	"sync"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

const URL_FORMAT = "https://www.ztm.poznan.pl/pl/dla-deweloperow/getParkingFile?file=ZTM_ParkAndRide__%s.csv"

type Provider struct{}

func (p Provider) ParkingLots(ctx context.Context) (<-chan map[wheretopark.ID]wheretopark.ParkingLot, error) {
	var wg sync.WaitGroup
	ch := make(chan map[wheretopark.ID]wheretopark.ParkingLot, len(configuration.ParkingLots))
	for name, metadata := range configuration.ParkingLots {
		wg.Add(1)
		go func(name string, metadata wheretopark.Metadata) {
			defer wg.Done()
			url, err := url.Parse(fmt.Sprintf(URL_FORMAT, name))
			if err != nil {
				log.Ctx(ctx).Err(err).Msg("invalid url")
				return
			}
			str, err := wheretopark.GetString(url, nil)
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
			parkingLot := wheretopark.ParkingLot{
				Metadata: metadata,
				State: wheretopark.State{
					LastUpdated: data.LastUpdated,
					AvailableSpots: map[wheretopark.SpotType]uint{
						wheretopark.SpotTypeCar: data.AvailableSpots,
					},
				},
			}
			ch <- map[wheretopark.ID]wheretopark.ParkingLot{
				id: parkingLot,
			}
		}(name, metadata)
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	return ch, nil
}

func New() Provider {
	return Provider{}
}
