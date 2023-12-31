package krakow

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"wheretopark/forecaster"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

type Solari3000 struct {
	BasePath string
}

func NewSolari3000(basePath string) Solari3000 {
	return Solari3000{
		BasePath: basePath,
	}
}

func (s Solari3000) Name() string {
	return "Solari3000"
}

func (s Solari3000) Load() (map[string]*forecaster.ParkingMeter, error) {
	files, err := ListFilesByExtension(s.BasePath, ".xlsx")
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	var mu sync.Mutex
	meters := make(map[string]*forecaster.ParkingMeter)
	for _, file := range files {
		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()
			log.Debug().Str("file", file.Name()).Msg("loading Solari3000 file")
			entries, err := s.loadFile(filepath.Join(s.BasePath, file.Name()))
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			mu.Lock()
			for _, entry := range entries {
				if _, ok := meters[entry.Code]; !ok {
					meters[entry.Code] = &forecaster.ParkingMeter{
						Name:          fmt.Sprintf("Parking %s", entry.Code),
						OccupancyData: make(map[time.Time]uint),
					}
				}
				meters[entry.Code].AddOccupancy(entry.Date, entry.Date.Add(entry.Duration))
			}
			mu.Unlock()
		}(file)
	}
	wg.Wait()

	return meters, nil
}

type Solari3000Entry struct {
	Result   string
	Code     string
	Zone     string
	Subzone  string
	Date     time.Time
	Duration time.Duration
}

func (s Solari3000) loadFile(path string) ([]Solari3000Entry, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Error().Str("path", path).Err(err).Msg("closing excel file failed")
			panic(err)
		}
	}()
	rows, err := f.GetRows("Transactions")
	if err != nil {
		return nil, err
	}

	entries := make([]Solari3000Entry, 0, len(rows))
	skipped := 0
	for i, row := range rows {
		if i == 0 {
			continue
		}
		result := row[0]
		code := row[1]
		zone := row[2]
		subzone := row[3]
		date := row[4]
		duration := row[5]
		// for some reason this column exists in the spreadsheet, not sure why
		_ = row[6]
		if duration == "" {
			skipped++
			continue
		}

		entry := Solari3000Entry{
			Result:   result,
			Code:     code,
			Zone:     zone,
			Subzone:  subzone,
			Date:     wheretopark.MustParseDateTimeWith(solari3000DateFormat, date),
			Duration: wheretopark.Must(parseSolari3000Duration(duration)),
		}
		entries = append(entries, entry)
	}
	log.Debug().Int("entries", len(entries)).Int("skipped", skipped).Str("file", path).Msg("loaded entries")

	return entries, nil
}

// 15/12/22 9:08
const solari3000DateFormat = "1/2/06 15:04"

func parseSolari3000Duration(str string) (time.Duration, error) {
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
