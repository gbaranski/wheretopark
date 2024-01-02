package krakow

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"wheretopark/forecaster"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

type Krakow struct {
	metadataPath string
	sources      map[string]Source
}

func NewKrakow(baseSourcePath string) Krakow {
	return Krakow{
		metadataPath: filepath.Join(baseSourcePath, "parkomaty.xml"),
		sources: map[string]Source{
			"Flowbird":    NewFlowbird(filepath.Join(baseSourcePath, "FLOWBIRD")),
			"Solari 2000": NewSolari(filepath.Join(baseSourcePath, "SOLARI 2000")),
			"Solari 3000": NewSolari(filepath.Join(baseSourcePath, "SOLARI 3000")),
		},
	}
}

type meterCode = string

type meterRecord struct {
	zone      string
	subzone   string
	startDate time.Time
	endDate   time.Time
}

type krakowParkingLot struct {
	forecaster.ParkingLot
	Zone string
}

type SourceFile[T any] interface {
	Read(*os.File) T
}

type Source interface {
	Files() ([]string, error)
	LoadFile(*os.File) (map[meterCode][]meterRecord, error)
}

func (k Krakow) LoadSource(ctx context.Context, source Source) (map[meterCode][]meterRecord, error) {
	start := time.Now()
	files, err := source.Files()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	allRecords := make(map[meterCode][]meterRecord)
	for _, file := range files {
		wg.Add(1)
		go func(filepath string) {
			defer wg.Done()
			file, err := os.Open(filepath)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Str("file", filepath).Msg("failed to open file")
			}
			defer file.Close()
			records, err := source.LoadFile(file)
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

func (k Krakow) Load() (map[string]map[wheretopark.ID]forecaster.ParkingLot, error) {
	metadata, err := k.Metadata()
	if err != nil {
		log.Fatal().Err(err).Msg("error loading metadata")
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	allRecords := make(map[meterCode][]meterRecord)
	log.Info().Msg(fmt.Sprintf("loading meter records from %d sources", len(k.sources)))
	for name, source := range k.sources {
		wg.Add(1)
		ctx := log.With().Str("source", name).Logger().WithContext(context.Background())
		go func(source Source) {
			defer wg.Done()
			records, err := k.LoadSource(ctx, source)
			if err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("failed to load records")
				return
			}
			mu.Lock()
			for meterCode, records := range records {
				if _, ok := allRecords[meterCode]; !ok {
					allRecords[meterCode] = make([]meterRecord, 0, len(records))
				}
				allRecords[meterCode] = append(allRecords[meterCode], records...)
			}
			mu.Unlock()
		}(source)
	}
	wg.Wait()

	parkingLots := make(map[wheretopark.ID]krakowParkingLot, len(allRecords))
	for code, records := range allRecords {
		if len(records) < int(forecaster.MinimumRecords) {
			log.Debug().Str("meterCode", code).Msg(fmt.Sprintf("not enough records(%d), skipping", len(records)))
			continue

		}
		placemark, ok := metadata[code]
		if !ok {
			log.Warn().Str("code", code).Msg("no metadata for parking meter, skipping")
			continue
		}

		// create sequences
		sequences := make(map[time.Time]uint)
		for _, record := range records {
			currentInterval := record.startDate.Truncate(forecaster.DefaultInterval)
			endInterval := record.endDate.Truncate(forecaster.DefaultInterval)

			for currentInterval.Before(endInterval) {
				sequences[currentInterval]++
				currentInterval = currentInterval.Add(forecaster.DefaultInterval)
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
				currentInterval = currentInterval.Add(forecaster.DefaultInterval)
			}
		}

		parkingID := wheretopark.CoordinateToID(placemark.Coordinates.Latitude, placemark.Coordinates.Longitude)
		totalSpots := forecaster.MaxOccupiedSpots(sequences)
		parkingLot := krakowParkingLot{
			ParkingLot: forecaster.ParkingLot{
				TotalSpots: totalSpots,
				Sequences:  sequences,
			},
			Zone: string([]rune(placemark.Zone)[7]),
		}
		parkingLots[parkingID] = parkingLot
	}
	log.Info().Msg(fmt.Sprintf("loaded %d parking lots from %d sources", len(parkingLots), len(k.sources)))

	parkingLotsByZone := make(map[string]map[wheretopark.ID]forecaster.ParkingLot)
	for parkingID, parkingLot := range parkingLots {
		if _, ok := parkingLotsByZone[parkingLot.Zone]; !ok {
			parkingLotsByZone[parkingLot.Zone] = make(map[wheretopark.ID]forecaster.ParkingLot)
		}
		parkingLotsByZone[parkingLot.Zone][parkingID] = parkingLot.ParkingLot
	}

	return parkingLotsByZone, nil
}

func (k Krakow) Metadata() (map[string]Placemark, error) {
	placemarks, err := LoadPlacemarks(k.metadataPath)
	if err != nil {
		return nil, err
	}
	placemarksMap := make(map[string]Placemark, len(placemarks))
	for _, placemark := range placemarks {
		placemarksMap[placemark.Code] = placemark
	}
	return placemarksMap, nil

}

func ListFilesByExtension(path string, extension string) ([]string, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	filteredEntries := make([]string, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), extension) {
			filteredEntries = append(filteredEntries, filepath.Join(path, entry.Name()))
		}
	}
	return filteredEntries, nil
}
