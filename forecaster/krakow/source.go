package krakow

import (
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"wheretopark/forecaster"

	"github.com/rs/zerolog/log"
)

type Krakow struct {
	metadataPath string
	sources      []forecaster.Source
}

func NewKrakow(baseSourcePath string) Krakow {
	sources := []forecaster.Source{
		NewFlowbird(filepath.Join(baseSourcePath, "FLOWBIRD")),
		NewSolari2000(filepath.Join(baseSourcePath, "SOLARI 2000")),
		NewSolari3000(filepath.Join(baseSourcePath, "SOLARI 3000")),
	}
	return Krakow{
		metadataPath: filepath.Join(baseSourcePath, "parkomaty.xml"),
		sources:      sources,
	}
}

func (k Krakow) Load() (map[string]*forecaster.ParkingMeter, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex
	allMeters := make(map[string]*forecaster.ParkingMeter)
	for _, source := range k.sources {
		wg.Add(1)
		go func(source forecaster.Source) {
			defer wg.Done()
			log.Info().Str("source", source.Name()).Msg("loading source")
			start := time.Now()
			meters, err := source.Load()
			if err != nil {
				log.Error().Err(err).Str("source", source.Name()).Msg("failed to load source")
				return
			}
			log.Info().Str("source", source.Name()).Int("meters", len(meters)).Float64("time", time.Since(start).Seconds()).Msg("loaded source")
			mu.Lock()
			for code, meter := range meters {
				allMeters[code] = meter
			}
			mu.Unlock()
		}(source)
	}
	wg.Wait()

	return allMeters, nil
}

func (k Krakow) Metadata() ([]Placemark, error) {
	return LoadPlacemarks(k.metadataPath)

}

func ListFilesByExtension(path string, extension string) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	filteredEntries := make([]os.DirEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		if strings.HasSuffix(entry.Name(), extension) {
			filteredEntries = append(filteredEntries, entry)
		}
	}
	return filteredEntries, nil
}
