package wheretopark_test

import (
	"testing"
	wheretopark "wheretopark/go"

	"github.com/stretchr/testify/assert"
)

func TestCacheProvider(t *testing.T) {
	provider, err := wheretopark.NewCacheProvider()
	if err != nil {
		t.Fatal(err)
	}
	metadatas, found := provider.GetMetadatas("test")
	assert.Empty(t, metadatas)
	assert.False(t, found)
	states, found := provider.GetStates("test")
	assert.Empty(t, states)
	assert.False(t, found)

	parkingLots := map[wheretopark.ID]wheretopark.ParkingLot{
		sampleParkingLotID: sampleParkingLot,
	}
	err = provider.SetParkingLots("test", parkingLots)
	if err != nil {
		t.Fatal(err)
	}

	metadatas, found = provider.GetMetadatas("test")
	assert.True(t, found)
	equalJson[wheretopark.Metadata](t, metadatas[sampleParkingLotID], sampleParkingLot.Metadata, "metadata mismatch")
	states, found = provider.GetStates("test")
	assert.True(t, found)
	equalJson[wheretopark.State](t, states[sampleParkingLotID], sampleParkingLot.State, "state mismatch")

}
