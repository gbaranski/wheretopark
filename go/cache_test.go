package wheretopark_test

import (
	"encoding/json"
	"testing"
	wheretopark "wheretopark/go"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	cache, err := wheretopark.NewCache()
	if err != nil {
		t.Fatal(err)
	}

	metadatas := cache.GetParkingLots("test")
	assert.Nil(t, metadatas)

	parkingLots := map[wheretopark.ID]wheretopark.ParkingLot{
		sampleParkingLotID: sampleParkingLot,
	}
	err = cache.SetParkingLots("test", parkingLots)
	if err != nil {
		t.Fatal(err)
	}

	parkingLots = cache.GetParkingLots("test")
	assert.NotNil(t, parkingLots)
	equalJson[wheretopark.ParkingLot](t, metadatas[sampleParkingLotID], sampleParkingLot, "parking lot mismatch")
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
