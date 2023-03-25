package wheretopark_test

import (
	"log"
	"math/rand"
	"net/url"
	"os"
	"testing"
	"time"
	wheretopark "wheretopark/go"

	"github.com/goccy/go-json"
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

func client(t *testing.T) *wheretopark.Client {
	rawURL := getEnvOr("SURREALDB_URL", "ws://localhost:8000")
	url, err := url.Parse(rawURL)
	if err != nil {
		t.Fatal(err)
	}
	client, err := wheretopark.NewClient(url, "wheretopark", "testing")
	if err != nil {
		t.Fatal(err)
	}
	err = client.SignInWithPassword("root", "password")
	if err != nil {
		t.Fatal(err)
	}
	return client
}

func TestClient(t *testing.T) {
	client := client(t)
	id := RandomID()
	err := client.SetParkingLot(id, sampleParkingLot)
	if err != nil {
		t.Fatal(err)
	}
	if err = client.SetParkingLot(id, sampleParkingLot); err != nil {
		t.Fatal(err)
	}

	obtainedParkingLot, err := client.GetParkingLot(id)
	if err != nil {
		t.Fatal(err)
	}
	equalJson(t, sampleParkingLot, *obtainedParkingLot, "obtained parking lot doesn't match with parking lot that was added")

	err = client.DeleteParkingLot(id)
	if err != nil {
		t.Fatal(err)
	}
	if err = client.DeleteParkingLot(id); err != nil {
		t.Fatal(err)
	}

	exists, err := client.ExistsParkingLot(id)
	if err != nil {
		t.Fatal(err)
	}
	if exists {
		t.Fatalf("client should report that %s does not exist\n", id)
	}
	obtainedParkingLot, err = client.GetParkingLot(id)
	if err != nil {
		t.Fatal(err)
	}
	if obtainedParkingLot != nil {
		t.Fatalf("parkign lot %s should have been deleted", id)
	}

}

func equalJson[T any](t *testing.T, a T, b T, msg string) {
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
