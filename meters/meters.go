package meters

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

type Record struct {
	StartDate time.Time
	EndDate   time.Time
}

type DataSource interface {
	Files() ([]string, error)
	LoadRecords(file *os.File) (map[wheretopark.ID][]Record, error)
}

type Meters struct {
	sources map[string]DataSource
}

func NewMeters(sources map[string]DataSource) Meters {
	return Meters{
		sources: sources,
	}
}

const MinimumMeterRecords uint = 50
const MeterInterval time.Duration = time.Minute * 15

func (m Meters) LoadRecordsOf(ctx context.Context, source DataSource) (map[wheretopark.ID][]Record, error) {
	start := time.Now()
	files, err := source.Files()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	allRecords := make(map[wheretopark.ID][]Record)
	for _, file := range files {
		wg.Add(1)
		go func(filepath string) {
			defer wg.Done()
			file, err := os.Open(filepath)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Str("file", filepath).Msg("failed to open file")
			}
			defer file.Close()
			records, err := source.LoadRecords(file)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("failed to load records")
				return
			}
			mu.Lock()
			for id, records := range records {
				allRecords[id] = append(allRecords[id], records...)
			}
			mu.Unlock()

		}(file)
	}
	wg.Wait()

	var totalRecords uint
	for _, records := range allRecords {
		totalRecords += uint(len(records))
	}

	log.
		Ctx(ctx).
		Info().
		Msg(
			fmt.Sprintf(
				"loaded %d meter records from %d files in %d ms",
				totalRecords,
				len(files),
				time.Since(start).Milliseconds()),
		)
	return allRecords, nil
}

func (m Meters) LoadRecords() (map[wheretopark.ID][]Record, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	allRecords := make(map[wheretopark.ID][]Record)
	log.Info().Msg(fmt.Sprintf("loading meter records from %d sources", len(m.sources)))
	for name, source := range m.sources {
		wg.Add(1)
		ctx := log.With().Str("source", name).Logger().WithContext(context.Background())
		go func(source DataSource) {
			defer wg.Done()
			records, err := m.LoadRecordsOf(ctx, source)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("failed to load records")
				return
			}
			mu.Lock()
			for id, records := range records {
				if _, ok := allRecords[id]; !ok {
					allRecords[id] = make([]Record, 0, len(records))
				}
				allRecords[id] = append(allRecords[id], records...)
			}
			mu.Unlock()
		}(source)
	}
	wg.Wait()

	return allRecords, nil
}

func (m Meters) Sequences() (map[wheretopark.ID]map[time.Time]uint, error) {
	allRecords, err := m.LoadRecords()
	if err != nil {
		return nil, fmt.Errorf("error loading records: %w", err)
	}
	allSequences := make(map[wheretopark.ID]map[time.Time]uint, len(allRecords))
	for id, records := range allRecords {
		if len(records) < int(MinimumMeterRecords) {
			log.Debug().Str("id", id).Msg(fmt.Sprintf("not enough records(%d), skipping", len(records)))
			continue

		}
		// create sequences
		sequences := make(map[time.Time]uint)
		for _, record := range records {
			currentInterval := record.StartDate.Truncate(MeterInterval)
			endInterval := record.EndDate.Truncate(MeterInterval)

			for currentInterval.Before(endInterval) {
				sequences[currentInterval]++
				currentInterval = currentInterval.Add(MeterInterval)
			}
		}
		// fill in missing intervals
		{
			earliest := time.Now()
			latest := time.Time{}
			for interval := range sequences {
				if interval.Before(earliest) {
					earliest = interval
				}
				if interval.After(latest) {
					latest = interval
				}
			}
			currentInterval := earliest
			endInterval := latest
			for currentInterval.Before(endInterval) {
				if _, ok := sequences[currentInterval]; !ok {
					sequences[currentInterval] = 0
				}
				currentInterval = currentInterval.Add(MeterInterval)
			}
		}

		allSequences[id] = sequences
	}
	log.Info().Msg(fmt.Sprintf("loaded %d parking lots", len(allSequences)))
	for id, sequences := range allSequences {
		for time, occupiedSpots := range sequences {
			if occupiedSpots > 1000 {
				fmt.Printf("id: %s, time: %s, occupiedSpots: %d\n", id, time, occupiedSpots)
			}
		}
	}
	return allSequences, nil

}

func MaxOccupiedSpots(sequences map[time.Time]uint) uint {
	count := uint(0)
	for _, value := range sequences {
		count = max(count, uint(value))
	}
	return count
}
