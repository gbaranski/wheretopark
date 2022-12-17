package warsaw_test

import (
	"log"
	"testing"
	"wheretopark/providers/collector/warsaw"
)

func TestData(t *testing.T) {
	data, err := warsaw.GetData()
	if err != nil {
		t.Fatal(err)
	}
	log.Println(data)
}
