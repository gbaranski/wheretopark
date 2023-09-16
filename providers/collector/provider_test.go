package collector_test

import (
	"testing"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector"
)

func TestProvider(t *testing.T) {
	provider, err := collector.NewProvider()
	if err != nil {
		t.Fatal(err)
	}
	wheretopark.ExamineProvider(t, provider)
}
