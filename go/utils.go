package wheretopark

import (
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
)

func NewIntPricingRule(duration string, price int32) PricingRule {
	return PricingRule{
		Duration:  duration,
		Price:     decimal.NewFromInt32(price),
		Repeating: false,
	}
}

func NewPricingRule(duration string, price decimal.Decimal) PricingRule {
	return PricingRule{
		Duration:  duration,
		Price:     price,
		Repeating: false,
	}
}

func WithTimeout(fn func() error, timeout time.Duration) error {
	done := make(chan error)
	go func() {
		done <- fn()
	}()
	select {
	case <-time.After(timeout):
		return fmt.Errorf("timeout")
	case err := <-done:
		return err
	}
}

func MustParseDate(date string) time.Time {
	v, err := time.Parse(time.DateOnly, date)
	if err != nil {
		panic(err)
	}
	return v
}

func MustParseDateTime(dateTime string) time.Time {
	v, err := time.Parse(time.RFC3339, dateTime)
	if err != nil {
		panic(err)
	}
	return v
}

func MustLoadLocation(name string) *time.Location {
	location, err := time.LoadLocation(name)
	if err != nil {
		panic(err)
	}
	return location
}

func MustParseURL(v string) *url.URL {
	url, err := url.Parse(v)
	if err != nil {
		panic(err)
	}
	return url
}

func ExtractMetadatas(parkingLots map[ID]ParkingLot) map[ID]Metadata {
	metadatas := make(map[ID]Metadata, len(parkingLots))
	for id, parkingLot := range parkingLots {
		metadatas[id] = parkingLot.Metadata
	}
	return metadatas
}

func ExtractStates(parkingLots map[ID]ParkingLot) map[ID]State {
	states := make(map[ID]State, len(parkingLots))
	for id, parkingLot := range parkingLots {
		states[id] = parkingLot.State
	}
	return states
}

func JoinMetadatasAndStates(metadatas map[ID]Metadata, states map[ID]State) (map[ID]ParkingLot, error) {
	for id := range metadatas {
		_, exists := states[id]
		if !exists {
			delete(metadatas, id)
		}
	}
	for id := range states {
		_, exists := metadatas[id]
		if !exists {
			delete(states, id)
		}
	}
	if len(metadatas) != len(states) {
		return nil, fmt.Errorf("metadatas and states must have the same length")
	}
	parkingLots := make(map[ID]ParkingLot, len(metadatas))
	for id, metadata := range metadatas {
		parkingLots[id] = ParkingLot{
			Metadata: metadata,
			State:    states[id],
		}
	}
	return parkingLots, nil
}

func MergeMaps[K comparable, T any](values ...map[K]T) map[K]T {
	result := make(map[K]T)
	for _, subvalues := range values {
		for id, value := range subvalues {
			result[id] = value
		}
	}
	return result
}

var logLevelMappings = map[string]zerolog.Level{
	"trace": zerolog.TraceLevel,
	"debug": zerolog.DebugLevel,
	"info":  zerolog.InfoLevel,
	"warn":  zerolog.WarnLevel,
	"error": zerolog.ErrorLevel,
	"fatal": zerolog.FatalLevel,
	"panic": zerolog.PanicLevel,
}

func InitLogging() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	logLevelString := os.Getenv("LOG_LEVEL")
	logLevel, exists := logLevelMappings[logLevelString]
	if !exists {
		logLevel = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(logLevel)
}
