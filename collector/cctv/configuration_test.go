package cctv_test

import (
	"log"
	"testing"
	cctv "wheretopark/collector/cctv"
)

func TestParse(t *testing.T) {
	log.Println(cctv.DefaultConfiguration)
	for _, parkingLot := range cctv.DefaultConfiguration.ParkingLots {
		if err := parkingLot.Validate(); err != nil {
			t.Fatalf("%s - invalid parking lot: %s", parkingLot.Name, err)
		}
		if len(parkingLot.Cameras) <= 0 {
			t.Fatalf("%s - no cameras", parkingLot.Name)
		}
	}
}
