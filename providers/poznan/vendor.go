package poznan

import (
	"encoding/csv"
	"strconv"
	"strings"
	"time"
)

type Data struct {
	LastUpdated    time.Time
	AvailableSpots uint
}

func parse(data string) (*Data, error) {
	r := csv.NewReader(strings.NewReader(data))
	r.Comma = ';'
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	record := records[1] // take second one, because first one is the list of labels
	lastUpdated, err := time.ParseInLocation("2006-01-02 15:04:05", record[0], defaultTimezone)
	if err != nil {
		return nil, err
	}
	availableSpots, err := strconv.Atoi(record[1])
	if err != nil {
		return nil, err
	}
	return &Data{
		lastUpdated,
		uint(availableSpots),
	}, nil

}
