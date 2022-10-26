package wheretopark_test

import (
	"log"
	"math/rand"
	"net/url"
	"os"
	"testing"
	"time"
	wheretopark "wheretopark/go"

	geojson "github.com/paulmach/go.geojson"
	"github.com/stretchr/testify/assert"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyz12345678")

func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandomID() string {
	return RandStringRunes(12)
}

func getEnvOr(name string, or string) string {
	value, exists := os.LookupEnv(name)
	if !exists {
		log.Printf("using default value for %s (%s)", name, or)
		value = or
	}
	return value
}

func client() *wheretopark.Client {
	rawURL := getEnvOr("SURREALDB_URL", "ws://localhost:8000")
	url, err := url.Parse(rawURL)
	if err != nil {
		log.Fatal(err)
	}
	client, err := wheretopark.NewClient(url, "wheretopark", "testing")
	if err != nil {
		log.Fatal(err)
	}
	err = client.SignInWithPassword("root", "root")
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func TestState(t *testing.T) {
	client := client()
	state := wheretopark.State{
		LastUpdated: "2022-10-21T23:09:47+0000",
		AvailableSpots: map[string]uint{
			"CAR": 123,
		},
	}
	id := RandomID()
	err := client.SetState(id, state)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SetState(id, state)
	if err != nil {
		log.Fatal(err)
	}

	obtainedState, err := client.GetState(id)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, state, *obtainedState, "obtained state doesn't match with state that was added")

	err = client.DeleteState(id)
	if err != nil {
		log.Fatal(err)
	}
	err = client.DeleteState(id)
	if err != nil {
		log.Fatal(err)
	}

	exists, err := client.ExistsState(id)
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		log.Fatalf("client should report that %s does not exist\n", id)
	}
	obtainedState, err = client.GetState(id)
	if err != nil {
		log.Fatal(err)
	}
	if obtainedState != nil {
		log.Fatalf("state %s should have been deleted", id)
	}

}

func TestMetadata(t *testing.T) {
	client := client()
	metadata := wheretopark.Metadata{
		Name:     "Galeria Bałtycka",
		Address:  "ul. Dmowskiego",
		Location: *geojson.NewPointFeature([]float64{18.60024, 54.38268}),
		Resources: []string{
			"mailto:galeria@galeriabaltycka.pl",
			"tel:+48-58-521-85-52",
			"https://www.galeriabaltycka.pl/o-centrum/dojazd-parkingi/parkingi/",
		},
		TotalSpots: map[string]uint{
			"CAR": 1110,
		},
		MaxWidth:  nil,
		MaxHeight: nil,
		Features: []string{
			"COVERED",
			"UNCOVERED",
		},
		PaymentMethods: []string{
			"CASH",
			"CARD",
			"CONTACTLESS",
		},
		Comment: map[string]string{
			"pl": "Polski komentarz",
			"en": "English comment",
		},
		Currency: "PLN",
		Timezone: "Europe/Warsaw",
		Rules: []wheretopark.Rule{
			{
				Hours:   "Mo-Sa 08:00-22:00; Su 09:00-21:00",
				Applies: nil,
				Pricing: []wheretopark.PricingRule{
					{
						Duration: "PT1H",
						Price:    0.0,
					},
					{
						Duration: "PT2H",
						Price:    2.0,
					},
					{
						Duration: "PT3H",
						Price:    3.0,
					},
					{
						Duration: "PT24H",
						Price:    25.0,
					},
					{
						Duration:  "PT1H",
						Price:     4.0,
						Repeating: true,
					},
				},
			},
		},
	}
	id := RandomID()
	err := client.SetMetadata(id, metadata)
	if err != nil {
		log.Fatal(err)
	}

	obtainedMetadata, err := client.GetMetadata(id)
	if err != nil {
		log.Fatal(err)
	}
	// weird
	if len(metadata.Location.Properties) == 0 {
		metadata.Location.Properties = nil
	}
	assert.Equal(t, metadata, *obtainedMetadata, "obtained metadata doesn't match with metadata that was added")

	err = client.DeleteMetadata(id)
	if err != nil {
		log.Fatal(err)
	}

	exists, err := client.ExistsMetadata(id)
	if err != nil {
		log.Fatal(err)
	}
	if exists {
		log.Fatalf("client should report that %s does not exist\n", id)
	}
	obtainedMetadata, err = client.GetMetadata(id)
	if err != nil {
		log.Fatal(err)
	}
	if obtainedMetadata != nil {
		log.Fatalf("metadata %s should have been deleted", id)
	}

}
