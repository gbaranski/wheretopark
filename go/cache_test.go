package wheretopark_test

import (
	"testing"
	wheretopark "wheretopark/go"

	"github.com/stretchr/testify/assert"
)

func TestPluralCache(t *testing.T) {
	cache, err := wheretopark.NewPluralCache()
	if err != nil {
		t.Fatal(err)
	}
	metadatas := cache.GetMetadatas("test")
	assert.Nil(t, metadatas)
	states := cache.GetStates("test")
	assert.Nil(t, states)

	parkingLots := map[wheretopark.ID]wheretopark.ParkingLot{
		sampleParkingLotID: sampleParkingLot,
	}
	err = cache.SetParkingLots("test", parkingLots)
	if err != nil {
		t.Fatal(err)
	}

	metadatas = cache.GetMetadatas("test")
	assert.NotNil(t, metadatas)
	equalJson[wheretopark.Metadata](t, metadatas[sampleParkingLotID], sampleParkingLot.Metadata, "metadata mismatch")
	states = cache.GetStates("test")
	assert.NotNil(t, states)
	equalJson[wheretopark.State](t, states[sampleParkingLotID], sampleParkingLot.State, "state mismatch")

}

func TestSingularCache(t *testing.T) {
	cache, err := wheretopark.NewSingularCache()
	if err != nil {
		t.Fatal(err)
	}
	state := cache.GetState("abcdefg")
	assert.Nil(t, state)

	err = cache.SetState("abcdefg", sampleParkingLot.State)
	if err != nil {
		t.Fatal(err)
	}

	state = cache.GetState("abcdefg")
	assert.NotNil(t, state)
	equalJson[wheretopark.State](t, *state, sampleParkingLot.State, "state mismatch")

}
