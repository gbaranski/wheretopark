package krakow_test

import (
	"testing"
	"wheretopark/collector/krakow"
)

func TestGetPlacemarks(t *testing.T) {
	_, err := krakow.GetPlacemarks()
	if err != nil {
		t.Fatal(err)
	}
}

func TestMetadata(t *testing.T) {
	placemarks, err := krakow.GetPlacemarks()
	if err != nil {
		t.Fatal(err)
	}
	for _, placemark := range placemarks {
		metadata := placemark.Metadata(0)
		err := metadata.Validate()
		if err != nil {
			t.Fatal(err)
		}
	}
}
