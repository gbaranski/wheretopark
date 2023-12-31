package forecaster

import (
	"fmt"
	"sort"

	// "sync"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

type ParkingLot struct {
	ID         wheretopark.ID `json:"id"`
	TotalSpots uint           `json:"totalSpots"`
	Interval   time.Duration  `json:"interval"`
	StartDate  time.Time      `json:"startDate"`
	EndDate    time.Time      `json:"endDate"`
	// Sequence of occupied spots
	Sequences []uint `json:"sequences"`
}

type Timeseries struct {
	ParkingLots map[wheretopark.ID]ParkingLot `json:"parkingLots"`
}

func NewTimeseries(meters map[wheretopark.ID]*ParkingMeter) Timeseries {
	allParkingLots := make(map[wheretopark.ID]ParkingLot, len(meters))
	for id, meter := range meters {
		meter.FillEmptyIntervals()
		type ComplexSequence = struct {
			Date          time.Time
			OccupiedSpots uint
		}
		complexSequences := make([]ComplexSequence, 0, len(meter.OccupancyData))
		for date, occupiedSpots := range meter.OccupancyData {
			complexSequences = append(complexSequences, ComplexSequence{
				Date:          date,
				OccupiedSpots: uint(occupiedSpots),
			})
		}
		sort.Slice(complexSequences, func(i, j int) bool {
			return complexSequences[i].Date.Before(complexSequences[j].Date)
		})

		simpleSequences := make([]uint, len(complexSequences))
		for i, complexSequence := range complexSequences {
			simpleSequences[i] = complexSequence.OccupiedSpots
		}

		allParkingLots[id] = ParkingLot{
			ID:         id,
			TotalSpots: meter.TotalSpots(),
			Interval:   Interval,
			StartDate:  complexSequences[0].Date,
			EndDate:    complexSequences[len(complexSequences)-1].Date,
			Sequences:  simpleSequences,
		}
		log.Debug().Msg(fmt.Sprintf("added %d sequences for %s", len(complexSequences), id))
	}
	return Timeseries{
		ParkingLots: allParkingLots,
	}
}

// func (s *Sequence) EncodeCSV() []string {
// 	return []string{
// 		s.ID,
// 		s.Date.Format(time.RFC3339),
// 		strconv.Itoa(int(s.OccupiedSpots)),
// 		strconv.Itoa(int(s.TotalSpots)),
// 	}
// }

// var headers = []string{"id", "date", "occupiedSpots", "totalSpots"}

// func (t *Timeseries) EncodeCSV() [][]string {
// 	data := make([][]string, len(t.sequences)+1)
// 	data[0] = headers
// 	for i, sequence := range t.sequences {
// 		data[i+1] = sequence.EncodeCSV()
// 	}
// 	return data
// }
