package providers_test

import (
	"testing"
	wheretopark "wheretopark/go"
	"wheretopark/go/providers"
	"wheretopark/go/tester"

	"github.com/stretchr/testify/assert"
)

func TestCache(t *testing.T) {
	cache, err := providers.NewCache()
	if err != nil {
		t.Fatal(err)
	}

	parkingLot := cache.GetParkingLot(tester.SampleParkingLotID)
	assert.Nil(t, parkingLot)

	err = cache.SetParkingLot(tester.SampleParkingLotID, &tester.SampleParkingLot)
	if err != nil {
		t.Fatal(err)
	}

	parkingLot = cache.GetParkingLot(tester.SampleParkingLotID)
	assert.NotNil(t, tester.SampleParkingLot)
	tester.EqualJson[wheretopark.ParkingLot](t, *parkingLot, tester.SampleParkingLot, "parking lot mismatch")
}
