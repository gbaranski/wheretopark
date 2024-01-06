package forecaster

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"net/url"
	"strconv"
	"strings"
	"time"
	wheretopark "wheretopark/go"

	"github.com/rs/zerolog/log"
)

type Pycaster struct {
	url *url.URL
}

type Prediction struct {
	Date      time.Time `json:"date"`
	Occupancy float32   `json:"occupancy"`
}

func (p Pycaster) Predict(id wheretopark.ID, sequences map[time.Time]uint, dateOnly time.Time) ([]Prediction, error) {
	var buf bytes.Buffer
	err := WriteSequencesCSV(&buf, sequences)
	if err != nil {
		log.Error().Err(err).Str("id", string(id)).Msg("error writing csv")
		return nil, fmt.Errorf("error writing csv: %w", err)
	}

	resp, err := wheretopark.DefaultClient.R().
		SetFileReader("file", fmt.Sprintf("%s.csv", id), &buf).
		Post(p.url.JoinPath("forecast", id, dateOnly.Format(time.DateOnly)).String())
	if err != nil {
		return nil, fmt.Errorf("error posting to server: %w", err)
	}
	fmt.Printf("resp: %+v", resp.String())
	predictions, err := ParsePredictions(resp.String())
	if err != nil {
		return nil, fmt.Errorf("error parsing predictions: %w", err)
	}
	return predictions, nil
}

func ParsePredictions(csvData string) ([]Prediction, error) {
	reader := csv.NewReader(strings.NewReader(csvData))
	var predictions []Prediction

	// Skip the header row.
	if _, err := reader.Read(); err != nil {
		return nil, err
	}

	// Iterate over the CSV rows.
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		predictions = append(predictions, Prediction{
			Date:      wheretopark.MustParseDateTimeWith(time.DateTime, row[0]),
			Occupancy: float32(wheretopark.Must(strconv.ParseFloat(row[1], 32))),
		})
	}

	return predictions, nil

}

func NewPycaster(url *url.URL) Pycaster {
	return Pycaster{url: url}
}
