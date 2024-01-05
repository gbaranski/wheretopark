package meters

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/zerolog/log"
)

type Code = string

type Record struct {
	StartDate time.Time
	EndDate   time.Time
}

type DataSource interface {
	Files() ([]string, error)
	LoadRecords(file *os.File) (map[Code][]Record, error)
}

type Meters struct {
	sources map[string]DataSource
}

func NewMeters(sources map[string]DataSource) Meters {
	return Meters{
		sources: sources,
	}
}

func (m Meters) LoadRecordsOf(ctx context.Context, source DataSource) (map[Code][]Record, error) {
	start := time.Now()
	files, err := source.Files()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	allRecords := make(map[Code][]Record)
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
			for code, records := range records {
				allRecords[code] = append(allRecords[code], records...)
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

func (m Meters) LoadRecords() (map[Code][]Record, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	allRecords := make(map[Code][]Record)
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
			for meterCode, records := range records {
				if _, ok := allRecords[meterCode]; !ok {
					allRecords[meterCode] = make([]Record, 0, len(records))
				}
				allRecords[meterCode] = append(allRecords[meterCode], records...)
			}
			mu.Unlock()
		}(source)
	}
	wg.Wait()

	return allRecords, nil
}
