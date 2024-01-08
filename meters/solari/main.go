package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/go/meters"
	"wheretopark/providers/krakow"

	"github.com/jessevdk/go-flags"
	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

var opts struct {
	InputPath  string `short:"i" long:"input" description:"Solari 2000/3000 data input path" required:"true"`
	OutputPath string `short:"o" long:"output" description:"Folder in where to output sequences" required:"false"`
}

func main() {
	wheretopark.InitLogging()
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to parse flags")
	}

	files, err := wheretopark.ListFilesWithExtension(opts.InputPath, "xlsx")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to list files")
	}
	if len(files) == 0 {
		log.Fatal().Msg(fmt.Sprintf("no files found under %s", opts.InputPath))
	}
	allRecords := make([]meters.Record, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup
	for _, path := range files {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			file, err := os.Open(path)
			if err != nil {
				log.Fatal().Err(err).Str("file", path).Msg("failed to open file")
			}
			records, err := loadRecords(file)
			if err != nil {
				log.Fatal().Err(err).Str("file", path).Msg("failed to load records")
			}
			mu.Lock()
			allRecords = append(allRecords, records...)
			mu.Unlock()
		}(path)
	}
	wg.Wait()

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
	spreadsheet, err := excelize.OpenReader(file)
	if err != nil {
		return nil, fmt.Errorf("error reading spreadsheet: %w", err)
	}
	defer spreadsheet.Close()

	rows, err := spreadsheet.GetRows("Transactions")
	if err != nil {
		return nil, err
	}

	records := make([]meters.Record, len(rows))
	for i, row := range rows {
		if i == 0 {
			continue
		}
		strCode := row[1]
		dateStr := row[4]
		durationStr := row[5]
		if durationStr == "" || !strings.Contains(dateStr, "/") {
			// log.Debug().Str("row", fmt.Sprintf("%v", row)).Msg("skipping entry")
			continue
		}
		startDate := wheretopark.MustParseDateTimeWith(dateFormat, dateStr)
		duration := wheretopark.Must(parseDuration(durationStr))

		code, err := strconv.ParseUint(strCode, 10, 32)
		if err != nil {
			log.Warn().Err(err).Str("code", strCode).Msg("failed to parse code")
			continue
		}
		records = append(records, meters.Record{
			Code:      uint(code),
			StartDate: startDate,
			EndDate:   startDate.Add(duration),
		})
	}

	return records, nil
}

const dateFormat = "1/2/06 15:04"

func parseDuration(str string) (time.Duration, error) {
	a := strings.Split(str, ":")
	h := a[0]
	m := a[1]
	s := a[2]
	if s != "00" {
		return 0, fmt.Errorf("invalid duration: %s", s)
	}

	strDuration := fmt.Sprintf("%sh%sm", h, m)
	return time.ParseDuration(strDuration)
}
