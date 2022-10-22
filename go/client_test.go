package wheretopark_test

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
	wheretopark "wheretopark/go"

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
	url := getEnvOr("SURREALDB_URL", "ws://localhost:8000/rpc")
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

	exists, err := client.StateExists(id)
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
