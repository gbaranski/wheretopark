package warsaw

import (
	wheretopark "wheretopark/go"

	_ "embed"

	"golang.org/x/text/currency"
)

type Configuration struct {
	ParkingLots map[wheretopark.ID]wheretopark.Metadata
}

var configuration = &Configuration{
	ParkingLots: make(map[string]wheretopark.Metadata),
}

func init() {
	for k, v := range ztpParkingLots {
		configuration.ParkingLots[k] = v
	}
	for k, v := range prParkingLots {
		configuration.ParkingLots[k] = v
	}
}

var (
	defaultTimezone    = wheretopark.MustLoadLocation("Europe/Warsaw")
	defaultLastUpdated = wheretopark.MustParseDate("2023-10-07")
	defaultCurrency    = currency.PLN
)
