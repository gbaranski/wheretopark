package glasgow_test

import (
	"testing"
	"wheretopark/providers/collector/glasgow"
)

func TestProvider(t *testing.T) {
	provider, err := glasgow.NewProvider()
	if err != nil {
		t.Fatal(err)
	}
	_, err = provider.GetParkingLots()
	if err != nil {
		t.Fatal(err)
	}
}
