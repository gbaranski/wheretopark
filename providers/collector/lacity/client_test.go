package lacity_test

import (
	"log"
	"testing"
	"wheretopark/providers/collector/lacity"
)

func TestSpaceMetadatas(t *testing.T) {
	spaces, err := lacity.GetSpaceMetadatas()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(spaces)
}

func TestSpaceStates(t *testing.T) {
	spaces, err := lacity.GetSpaceStates()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(spaces)
}
