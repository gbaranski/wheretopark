package poznan_test

import (
	"testing"
	"wheretopark/providers/collector/poznan"
)

func TestProvider(t *testing.T) {
	provider, err := poznan.NewProvider()
	if err != nil {
		t.Fatal(err)
	}
	_, err = provider.GetParkingLots()
	if err != nil {
		t.Fatal(err)
	}
}
