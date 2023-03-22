package warsaw_test

import (
	"testing"
	"wheretopark/providers/collector/warsaw"
)

func TestProvider(t *testing.T) {
	provider, err := warsaw.NewProvider()
	if err != nil {
		t.Fatal(err)
	}
	_, err = provider.GetParkingLots()
	if err != nil {
		t.Fatal(err)
	}
}
