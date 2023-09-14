package wheretopark_test

import (
	"encoding/json"
	"testing"
	"time"
	wheretopark "wheretopark/go"

	"github.com/google/go-cmp/cmp"
	geojson "github.com/paulmach/go.geojson"
	"github.com/shopspring/decimal"
	"golang.org/x/text/currency"
)

var (
	sampleParkingLot = wheretopark.ParkingLot{
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
	sampleParkingLotID = wheretopark.GeometryToID(sampleParkingLot.Metadata.Geometry)
)

func TestEncodeDecodeParkingLot(t *testing.T) {
	data, err := json.Marshal(sampleParkingLot)
	if err != nil {
		t.Fatal(err)
	}
	var decoded wheretopark.ParkingLot
	err = json.Unmarshal(data, &decoded)
	if err != nil {
		t.Fatal(err)
	}
	if diff := cmp.Diff(sampleParkingLot, decoded, cmp.AllowUnexported(currency.Unit{}, time.Location{})); diff != "" {
		t.Errorf("parking lot mismatch (-want +got):\n%s", diff)
	}
}
