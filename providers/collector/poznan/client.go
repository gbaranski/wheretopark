package poznan

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type Coordinate struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Metadata struct {
	ID             string     `json:"id"`
	Name           string     `json:"name"`
	ShortName      string     `json:"shortName"`
	Address        string     `json:"address"`
	StreetEntrance string     `json:"streetEntrance"`
	Location       Coordinate `json:"location"`
}

type State struct {
	ID             string `json:"parkingId"`
	AvailableSpots uint   `json:"availableSpots"`
	LastUpdate     string `json:"lastUpdate"`
}

type Row struct {
	LastUpdated    time.Time
	AvailableSpots uint
}

const (
	DATA_URL = "https://www.ztm.poznan.pl/pl/dla-deweloperow/getParkingFile?file=ZTM_ParkAndRide__%s.csv"
)

var client = resty.New()

func getParkingLotData(name string) (*Row, error) {
	url := fmt.Sprintf(DATA_URL, name)
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	body := resp.Body()
	r := csv.NewReader(bytes.NewReader(body))
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	record := records[1] // take second one, because first one is the list of labels
	date, err := time.ParseInLocation("2006-01-02 15:04:05", record[0], defaultLocation)
	if err != nil {
		return nil, err
	}
	availableSpots, err := strconv.Atoi(record[1])
	if err != nil {
		return nil, err
	}
	row := Row{
		LastUpdated:    date,
		AvailableSpots: uint(availableSpots),
	}
	return &row, nil
}

type container struct {
	mu   sync.Mutex
	rows map[string]Row
}

func GetData() (map[string]Row, error) {
	container := container{
		rows: make(map[string]Row, len(configuration.ParkingLots)),
	}
	var wg sync.WaitGroup

	for name := range configuration.ParkingLots {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			row, err := getParkingLotData(name)
			if err != nil {
				log.Err(err).Str("name", name).Msg("failed to get data")
			}
			container.mu.Lock()
			defer container.mu.Unlock()
			container.rows[name] = *row
		}(name)
	}
	wg.Wait()
	return container.rows, nil
}
