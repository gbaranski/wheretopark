package warsaw

import (
	wheretopark "wheretopark/go"

	_ "embed"

	geojson "github.com/paulmach/go.geojson"
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
		configuration.ParkingLots[k] = wheretopark.Metadata{
			LastUpdated:    defaultLastUpdated,
			Name:           "",
			Address:        v.Address,
			Geometry:       &geojson.Geometry{},
			Resources:      append(v.Resources, ztpBaseResources...),
			TotalSpots:     map[string]uint{},
			MaxDimensions:  &wheretopark.Dimensions{},
			Features:       append(v.Features, ztpBaseFeatures...),
			PaymentMethods: append(v.PaymentMethods, ztpBasePaymentMethods...),
			Comment:        v.Comment,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
			Rules:          v.Rules,
		}
	}
	for k, v := range prParkingLots {
		configuration.ParkingLots[k] = wheretopark.Metadata{
			LastUpdated:    defaultLastUpdated,
			Name:           "",
			Address:        v.Address,
			Geometry:       &geojson.Geometry{},
			Resources:      append(v.Resources, prBaseResources...),
			TotalSpots:     map[string]uint{},
			MaxDimensions:  &wheretopark.Dimensions{},
			Features:       append(v.Features, prBaseFeatures...),
			PaymentMethods: []string{},
			Comment:        prDefaultComment,
			Currency:       defaultCurrency,
			Timezone:       defaultTimezone,
			Rules:          prDefaultRules,
		}
	}
}

var (
	defaultTimezone    = wheretopark.MustLoadLocation("Europe/Warsaw")
	defaultLastUpdated = wheretopark.MustParseDate("2022-12-28")
	defaultCurrency    = currency.PLN
)
