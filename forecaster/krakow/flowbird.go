package krakow

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
	"golang.org/x/text/encoding/charmap"
)

type Flowbird struct {
	BasePath string
}

func NewFlowbird(basePath string) Source {
	return Flowbird{
		BasePath: basePath,
	}
}

func (f Flowbird) Files() ([]string, error) {
	return ListFilesByExtension(f.BasePath, ".csv")
}

func (f Flowbird) LoadFile(file *os.File) (map[meterCode][]meterRecord, error) {
	reader := csv.NewReader(charmap.ISO8859_2.NewDecoder().Reader(file))
	reader.FieldsPerRecord = 13
	reader.Comma = ';'

	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	records := make(map[meterCode][]meterRecord, len(data))
	for i, row := range data {
		if i == 0 {
			continue
		}
		startDate := row[1]
		code := row[2]
		totalTime := row[6]
		zone := row[8]
		subzone := row[9]
		endDate := row[12]

		if totalTime == "-" || zone == "" || subzone == "" {
			log.Debug().Str("row", fmt.Sprintf("%v", row)).Msg("skipping entry")
			continue
		}

		if _, ok := records[code]; !ok {
			records[code] = make([]meterRecord, 0, 8)
		}
		records[code] = append(records[code], meterRecord{
			zone:      zone,
			subzone:   subzone,
			startDate: wheretopark.MustParseDateTimeWith(flowbirdDateFormat, startDate),
			endDate:   wheretopark.MustParseDateTimeWith(flowbirdDateFormat, endDate),
		})
	}
	return records, nil
}

// func (f Flowbird) Load() (map[string]*forecaster.ParkingMeter, error) {
// 	files, err := f.Files()
// 	if err != nil {
// 		return nil, err
// 	}

// 	var wg sync.WaitGroup
// 	var mu sync.Mutex
// 	meters := make(map[string]*forecaster.ParkingMeter)
// 	for _, file := range files {
// 		wg.Add(1)
// 		go func(file os.DirEntry) {
// 			defer wg.Done()
// 			log.Debug().Str("file", file.Name()).Msg("loading Flowbird file")
// 			entries, err := f.loadFile(filepath.Join(f.BasePath, file.Name()))
// 			if err != nil {
// 				log.Fatal().Err(err).Send()
// 			}
// 			mu.Lock()
// 			for _, entry := range entries {
// 				if _, ok := meters[entry.Code]; !ok {
// 					meters[entry.Code] = &forecaster.ParkingMeter{
// 						Name:          fmt.Sprintf("Parking %s", entry.Code),
// 						OccupancyData: make(map[time.Time]uint),
// 					}
// 				}
// 				meters[entry.Code].AddOccupancy(entry.Date, entry.EndDate)
// 			}
// 			mu.Unlock()
// 		}(file)
// 	}
// 	wg.Wait()

// 	return meters, nil
// }

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
//
//
// type flowbirdEntry struct {
// 	ServerDate       time.Time
// 	Date             time.Time
// 	Code             string
// 	Amount           decimal.Decimal
// 	TotalTime        time.Duration
// 	PaidTime         time.Duration
// 	TotalTimeMinutes int
// 	PaidTimeMinutes  int
// 	Zone             string
// 	Subzone          string
// 	Address          string
// 	Type             string
// 	EndDate          time.Time
// }

const flowbirdDateFormat = "02/01/2006 15:04"

func parseFlowbirdDuration(s string) (time.Duration, error) {
	s = strings.ReplaceAll(s, " ", "")
	return time.ParseDuration(s)
}
