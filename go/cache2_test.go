package wheretopark_test

import (
	"testing"
	wheretopark "wheretopark/go"

	"github.com/stretchr/testify/assert"
)

func TestCache2(t *testing.T) {
	cache, err := wheretopark.NewCache2()
	if err != nil {
		t.Fatal(err)
	}

	parkingLots := cache.GetParkingLots("test")
	assert.Nil(t, parkingLots)

	sampleParkingLots := map[wheretopark.ID]wheretopark.ParkingLot{
		sampleParkingLotID: sampleParkingLot,
	}
	err = cache.SetParkingLots("test", sampleParkingLots)
	if err != nil {
		t.Fatal(err)
	}

	parkingLots = cache.GetParkingLots("test")
	assert.NotNil(t, sampleParkingLots)
	equalJson[map[wheretopark.ID]wheretopark.ParkingLot](t, parkingLots, sampleParkingLots, "parking lot mismatch")
}
