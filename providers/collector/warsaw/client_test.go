package warsaw_test

import (
	"log"
	"testing"
	"wheretopark/providers/collector/warsaw"

	"github.com/stretchr/testify/assert"
)

func TestData(t *testing.T) {
	data, err := warsaw.GetData()
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, len(data.Result.Parks) > 0, "No data returned")
	log.Println(data)
}
