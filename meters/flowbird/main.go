package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/go/meters"
	"wheretopark/providers/krakow"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/charmap"
)

var opts struct {
	InputPath  string `short:"i" long:"input" description:"Flowbird data input path" required:"true"`
	OutputPath string `short:"o" long:"output" description:"Folder in where to output sequences"`
}

func main() {
	wheretopark.InitLogging()
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse flags")
	}

	files, err := wheretopark.ListFilesWithExtension(opts.InputPath, "csv")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list files")
	}
	if len(files) == 0 {
		log.Fatal().Msg(fmt.Sprintf("no files found under %s", opts.InputPath))
	}
	allRecords := make([]meters.Record, 0)
	for _, filePath := range files {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatal().Err(err).Str("file", filePath).Msg("failed to open file")
		}
		records, err := loadRecords(file)
		if err != nil {
			log.Fatal().Err(err).Str("file", filePath).Msg("failed to load records")
		}
		allRecords = append(allRecords, records...)
	}

	placemarks, err := krakow.GetPlacemarks()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get placemarks")
	}
	timeseries := meters.NewTimeseriesFromMeterRecords(allRecords, meters.Options{
		MinimumRecords: meters.DefaultMinimumRecords,
		Interval:       meters.DefaultInterval,
		WorkingScheme:  meters.WorkingScheme{StartHour: 10, EndHour: 20, Weekdays: []time.Weekday{time.Monday, time.Tuesday, time.Wednesday, time.Thursday, time.Friday, time.Saturday}},
		Mapping:        krakow.CodeMapping(placemarks),
	})

	log.
		Info().
		Msg(
			fmt.Sprintf(
				"loaded %d meter records from %d meters from %d files",
				len(allRecords),
				len(timeseries.IDs()),
				len(files),
			),
		)
	for _, id := range timeseries.IDs() {
		fmt.Printf("%s: %d records. %d total spots\n", id, timeseries.CountFor(id), timeseries.MaxOccupancyOf(id))
	}

	if opts.OutputPath != "" {
		err := timeseries.WriteMultipleCSV(opts.OutputPath)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to write csv")
		}
	}
}

func loadRecords(file *os.File) ([]meters.Record, error) {
	reader := csv.NewReader(charmap.ISO8859_2.NewDecoder().Reader(file))
	reader.FieldsPerRecord = 13
	reader.Comma = ';'

	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	records := make([]meters.Record, len(data))
	for i, row := range data {
		if i == 0 {
			continue
		}
		startDate := row[1]
		strCode := row[2]
		totalTime := row[6]
		zone := row[8]
		subzone := row[9]
		endDate := row[12]

		if totalTime == "-" || zone == "" || subzone == "" {
			// log.Debug().Str("row", fmt.Sprintf("%v", row)).Msg("skipping entry")
			continue
		}

		code, err := strconv.ParseUint(strCode, 10, 32)
		if err != nil {
			log.Warn().Err(err).Str("code", row[2]).Msg("failed to parse code")
			continue
		}
		records = append(records, meters.Record{
			StartDate: wheretopark.MustParseDateTimeWith(dateFormat, startDate),
			EndDate:   wheretopark.MustParseDateTimeWith(dateFormat, endDate),
			Code:      uint(code),
		})
	}
	return records, nil
}

const dateFormat = "02/01/2006 15:04"
