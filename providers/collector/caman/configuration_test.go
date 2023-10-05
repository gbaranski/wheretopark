package caman_test

import (
	"log"
	"testing"
	"wheretopark/providers/collector/caman"
)

func TestParse(t *testing.T) {
	log.Println(caman.DefaultConfiguration)
	for _, parkingLot := range caman.DefaultConfiguration.ParkingLots {
		if err := parkingLot.Validate(); err != nil {
			t.Fatalf("%s - invalid parking lot: %s", parkingLot.Name, err)
		}
		if len(parkingLot.Cameras) <= 0 {
			t.Fatalf("%s - no cameras", parkingLot.Name)
		}
	}
}
