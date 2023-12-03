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

	parkingLot := cache.GetParkingLot(sampleParkingLotID)
	assert.Nil(t, parkingLot)

	err = cache.SetParkingLot(sampleParkingLotID, &sampleParkingLot)
	if err != nil {
		t.Fatal(err)
	}

	parkingLot = cache.GetParkingLot(sampleParkingLotID)
	assert.NotNil(t, sampleParkingLot)
	equalJson[wheretopark.ParkingLot](t, *parkingLot, sampleParkingLot, "parking lot mismatch")
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
