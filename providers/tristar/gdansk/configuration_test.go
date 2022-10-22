package gdansk_test

import (
	"log"
	"testing"
	"wheretopark/providers/tristar/gdansk"
)

func TestParse(t *testing.T) {
	configuration, err := gdansk.LoadConfiguration("configuration.yaml")
	if err != nil {
		t.Fatal(err)
	}
	log.Println(configuration)
}
