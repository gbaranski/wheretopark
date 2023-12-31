package forecaster

import (
	"fmt"
	"sort"
	"strconv"

	// "sync"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

type Sequence struct {
	ID            wheretopark.ID
	Date          time.Time
	OccupiedSpots uint
	TotalSpots    uint
}

type Timeseries struct {
	sequences []Sequence
}

func NewTimeseries(meters map[wheretopark.ID]*ParkingMeter) Timeseries {
	allSequences := make([]Sequence, 0)
	for id, meter := range meters {
		log.Debug().Msg(fmt.Sprintf("adding %d sequences for %s", len(meter.OccupancyData), id))
		sequences := make([]Sequence, 0, len(meter.OccupancyData))
		totalSpots := meter.TotalSpots()
		for date, occupiedSpots := range meter.OccupancyData {
			sequences = append(sequences, Sequence{
				ID:            id,
				Date:          date,
				OccupiedSpots: uint(occupiedSpots),
				TotalSpots:    totalSpots,
			})
		}
		sort.Slice(sequences, func(i, j int) bool {
			return sequences[i].Date.Before(sequences[j].Date)
		})
		allSequences = append(allSequences, sequences...)
		log.Debug().Msg(fmt.Sprintf("added %d sequences for %s", len(sequences), id))
	}
	return Timeseries{
		sequences: allSequences,
	}
}

func (s *Sequence) EncodeCSV() []string {
	return []string{
		s.ID,
		s.Date.Format(time.RFC3339),
		strconv.Itoa(int(s.OccupiedSpots)),
		strconv.Itoa(int(s.TotalSpots)),
	}
}

var headers = []string{"id", "date", "occupiedSpots", "totalSpots"}

func (t *Timeseries) EncodeCSV() [][]string {
	data := make([][]string, len(t.sequences)+1)
	data[0] = headers
	for i, sequence := range t.sequences {
		data[i+1] = sequence.EncodeCSV()
	}
	return data
}
