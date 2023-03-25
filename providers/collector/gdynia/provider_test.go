package gdynia_test

import (
	"testing"
	"wheretopark/go/provider/sequential"
	"wheretopark/providers/collector/gdynia"
)

func TestProvider(t *testing.T) {
	provider, err := gdynia.NewProvider()
	if err != nil {
		t.Fatal(err)
	}
	sequential.ExamineProvider(t, provider)
}
