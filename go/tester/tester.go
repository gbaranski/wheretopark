package tester

import (
	"context"
	"encoding/json"
	"testing"
	"time"
	wheretopark "wheretopark/go"
	"wheretopark/go/providers"

	geojson "github.com/paulmach/go.geojson"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/currency"
)

var (
	SampleParkingLot = wheretopark.ParkingLot{
		Metadata: wheretopark.Metadata{
			// LastUpdated:    wheretopark.MustParseDate("2021-03-22"),
			Name:           "Galeria Bałtycka",
			Address:        "ul. Dmowskiego",
			Geometry:       geojson.NewPointGeometry([]float64{18.60024, 54.38268}),
			Resources:      []string{"mailto:galeria@galeriabaltycka.pl", "tel:+48-58-521-85-52", "https://www.galeriabaltycka.pl/o-centrum/dojazd-parkingi/parkingi/"},
			TotalSpots:     map[string]uint{"CAR": 1110},
			MaxDimensions:  nil,
			Features:       []string{"COVERED", "UNCOVERED"},
			PaymentMethods: []string{"CASH", "CARD", "CONTACTLESS"},
			Comment:        map[string]string{"pl": "Polski komentarz", "en": "English comment"},
			Currency:       currency.PLN,
			Timezone:       time.UTC,
			Rules:          []wheretopark.Rule{{Hours: "Mo-Sa 08:00-22:00; Su 09:00-21:00", Applies: nil, Pricing: []wheretopark.PricingRule{{Duration: "PT1H", Price: decimal.Zero}, {Duration: "PT2H", Price: decimal.NewFromInt(2)}, {Duration: "PT3H", Price: decimal.NewFromInt(3)}, {Duration: "PT24H", Price: decimal.NewFromInt(25)}, {Duration: "PT1H", Price: decimal.NewFromInt(4), Repeating: true}}}},
		},
		State: wheretopark.State{
			LastUpdated: wheretopark.MustParseDateTime("2022-10-21T23:09:47Z"),
			AvailableSpots: map[string]uint{
				"CAR": 123,
			},
		},
	}
	SampleParkingLotID = wheretopark.GeometryToID(SampleParkingLot.Metadata.Geometry)
)

func ExamineParkingLots(t *testing.T, parkingLots map[wheretopark.ID]wheretopark.ParkingLot) {
	assert.NotEmpty(t, parkingLots)
	for id, parkingLot := range parkingLots {
		t.Logf("%+v", parkingLot)
		if err := parkingLot.Validate(); err != nil {
			t.Fatalf("invalid parking lot %s : %s", id, err)
		}
	}
}

func EqualJson[T any](t *testing.T, a T, b T, msg string) {
	aJson, err := json.Marshal(a)
	if err != nil {
		t.Fatal(err)
	}
	bJson, err := json.Marshal(b)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, aJson, bJson, msg)
}

func ExamineProvider(t *testing.T, provider providers.Provider) {
	ch, err := provider.ParkingLots(context.TODO())
	if err != nil {
		t.Fatalf("provider failed: %s", err)
	}

	parkingLots := make(map[wheretopark.ID]wheretopark.ParkingLot)
	for parkingLot := range ch {
		for id, lot := range parkingLot {
			parkingLots[id] = lot
		}
	}
	ExamineParkingLots(t, parkingLots)
}
