package krakow

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
	"wheretopark/forecaster"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"golang.org/x/text/encoding/charmap"
)

type Flowbird struct {
	BasePath string
}

func NewFlowbird(basePath string) Flowbird {
	return Flowbird{
		BasePath: basePath,
	}
}

func (f Flowbird) Name() string {
	return "Flowbird"
}

func (f Flowbird) Load() (map[string]*forecaster.ParkingMeter, error) {
	files, err := ListFilesByExtension(f.BasePath, ".csv")
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
			log.Debug().Str("file", file.Name()).Msg("loading Flowbird file")
			entries, err := f.loadFile(filepath.Join(f.BasePath, file.Name()))
			if err != nil {
				log.Fatal().Err(err).Send()
			}
			mu.Lock()
			for _, entry := range entries {
				if _, ok := meters[entry.Code]; !ok {
					meters[entry.Code] = &forecaster.ParkingMeter{
						Name:          fmt.Sprintf("Parking %s", entry.Code),
						OccupancyData: make(map[time.Time]int),
					}
				}
				meters[entry.Code].AddOccupancy(entry.Date, entry.EndDate)
			}
			mu.Unlock()
		}(file)
	}
	wg.Wait()

	return meters, nil
}

// Godzina serwera;
// Data parkomatu;
// Kod parkomatu;
// Kwota;
// Całkowity czas;
// Czas opłacony;
// Całkowity czas w min.;
// Opłacony czas w min.;
// Opis strefy;
// Opis obwodu;
// Adres;
// Rodzaj;
// Data zakończenia
type flowbirdEntry struct {
	ServerDate       time.Time
	Date             time.Time
	Code             string
	Amount           decimal.Decimal
	TotalTime        time.Duration
	PaidTime         time.Duration
	TotalTimeMinutes int
	PaidTimeMinutes  int
	Zone             string
	Subzone          string
	Address          string
	Type             string
	EndDate          time.Time
}

func (f Flowbird) loadFile(path string) ([]flowbirdEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	reader := csv.NewReader(charmap.ISO8859_2.NewDecoder().Reader(file))
	reader.FieldsPerRecord = 13
	reader.Comma = ';'

	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	entries := make([]flowbirdEntry, 0, len(data))
	skipped := 0
	for i, row := range data {
		if i == 0 {
			continue
		}
		serverDate := row[0]
		date := row[1]
		code := row[2]
		amount := row[3]
		totalTime := row[4]
		paidTime := row[5]
		totalTimeMinutes := row[6]
		paidTimeMinutes := row[7]
		zone := row[8]
		subzone := row[9]
		address := row[10]
		entryType := row[11]
		endDate := row[12]
		if totalTime == "-" || zone == "" || subzone == "" {
			skipped++
			continue
		}

		entry := flowbirdEntry{
			ServerDate:       wheretopark.MustParseDateTimeWith(flowbirdDateFormat, serverDate),
			Date:             wheretopark.MustParseDateTimeWith(flowbirdDateFormat, date),
			Code:             code,
			Amount:           wheretopark.Must(decimal.NewFromString(strings.ReplaceAll(amount, ",", "."))),
			TotalTime:        wheretopark.Must(parseFlowbirdDuration(totalTime)),
			PaidTime:         wheretopark.Must(parseFlowbirdDuration(paidTime)),
			TotalTimeMinutes: wheretopark.Must(strconv.Atoi(totalTimeMinutes)),
			PaidTimeMinutes:  wheretopark.Must(strconv.Atoi(paidTimeMinutes)),
			Zone:             zone,
			Subzone:          subzone,
			Address:          address,
			Type:             entryType,
			EndDate:          wheretopark.MustParseDateTimeWith(flowbirdDateFormat, endDate),
		}
		entries = append(entries, entry)
	}
	log.Debug().Int("entries", len(entries)).Int("skipped", skipped).Str("file", path).Msg("loaded entries")

	return entries, nil
}

const flowbirdDateFormat = "02/01/2006 15:04"

func parseFlowbirdDuration(s string) (time.Duration, error) {
	s = strings.ReplaceAll(s, " ", "")
	return time.ParseDuration(s)
}
