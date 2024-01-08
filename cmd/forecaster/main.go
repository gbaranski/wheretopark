package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"
	"wheretopark/collector/krakow"
	"wheretopark/forecaster"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"

	"github.com/caarlos0/env/v10"
)

type environment struct {
	Port            uint     `env:"PORT" envDefault:"8080"`
	ForecasterData  string   `env:"FORECASTER_DATA,expand" envDefault:"${HOME}/.local/share/wheretopark/forecaster"`
	ForecasterCache string   `env:"FORECASTER_CACHE,expand" envDefault:"${HOME}/.cache/wheretopark/forecaster"`
	PycasterURL     *url.URL `env:"PYCASTER_URL,required"`
}

func krakowSequences(basePath string) map[wheretopark.ID]map[time.Time]uint {
	placemarks, err := krakow.GetPlacemarks()
	if err != nil {
		log.Fatal().Err(err).Msg("error getting placemarks")
	}
	mapping := krakow.CodeMapping(placemarks)
	meterSources := []meters.DataSource{
		meters.NewFlowbird(filepath.Join(basePath, "FLOWBIRD"), mapping),
		meters.NewSolari(filepath.Join(basePath, "SOLARI 2000"), mapping, meters.SolariVersion2000),
		meters.NewSolari(filepath.Join(basePath, "SOLARI 3000"), mapping, meters.SolariVersion3000),
	}
	meters := meters.NewMeters(meterSources)
	sequences, err := meters.Sequences()
	if err != nil {
		log.Fatal().Err(err).Msg("error getting sequences")
	}
	return sequences
}

// func influxSequences() {
/// ... to be implemented
// }

func main() {
	wheretopark.InitLogging()

	environment := environment{}
	if err := env.Parse(&environment); err != nil {
		log.Fatal().Err(err).Send()
	}

	cachedSequencesPath := filepath.Join(environment.ForecasterCache, "sequences.csv")

	var sequences map[wheretopark.ID]map[time.Time]uint
	if SequencesExist(cachedSequencesPath) {
		sequences = LoadSequences(cachedSequencesPath)
		log.Info().Msg(fmt.Sprintf("loaded %d sequences to %s", len(sequences), cachedSequencesPath))
	} else {
		krakowSequences := krakowSequences(filepath.Join(environment.ForecasterData, "datasets", "krakow"))
		sequences = krakowSequences
		SaveSequences(cachedSequencesPath, sequences)
		log.Info().Msg(fmt.Sprintf("saved %d sequences from %s", len(sequences), cachedSequencesPath))
	}

	for id := range sequences {
		log.Info().Msg(fmt.Sprintf("parkingID: %s", id))
	}

	pycaster := forecaster.NewPycaster(environment.PycasterURL)
	server := forecaster.NewServer(pycaster, sequences)
	router := server.Router()
	if err := server.Run(router, environment.Port); err != nil {
		log.Fatal().Err(err).Msg("run server failure")
	}
}

func SequencesExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		log.Fatal().Err(err).Msg("failed to check if cache file exists")
	}

	return true
}

func SaveSequences(path string, sequences map[wheretopark.ID]map[time.Time]uint) {
	// Ensure the directory exists
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to create directory")
		}
	}

	file, err := os.Create(path)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to create cache file")
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	err = writer.Write([]string{"parkingID", "date", "occupancy"})
	if err != nil {
		log.Fatal().Err(err).Msg("failed to write headers to csv file")
	}

	// Write data
	for id, times := range sequences {
		for t, u := range times {
			err = writer.Write([]string{string(id), t.Format(time.RFC3339), strconv.FormatUint(uint64(u), 10)})
			if err != nil {
				log.Fatal().Err(err).Msg("failed to write record to csv file")
			}
		}
	}
}

func LoadSequences(path string) map[wheretopark.ID]map[time.Time]uint {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to open cache file")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	sequences := make(map[wheretopark.ID]map[time.Time]uint)

	// Read and discard headers
	_, err = reader.Read()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read headers from csv file")
	}

	// Read data
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal().Err(err).Msg("error reading csv file")
		}

		id := wheretopark.ID(record[0])
		t, err := time.Parse(time.RFC3339, record[1])
		if err != nil {
			log.Fatal().Err(err).Msg("error parsing time")
		}
		u, err := strconv.ParseUint(record[2], 10, 32)
		if err != nil {
			log.Fatal().Err(err).Msg("error parsing uint")
		}

		if sequences[id] == nil {
			sequences[id] = make(map[time.Time]uint)
		}
		sequences[id][t] = uint(u)
	}

	return sequences
}
