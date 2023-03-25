package poznan_test

import (
	"testing"
	"wheretopark/go/provider/simple"
	"wheretopark/providers/collector/poznan"
)

func TestProvider(t *testing.T) {
	provider, err := poznan.NewProvider()
	if err != nil {
		t.Fatal(err)
	}
	simple.ExamineProvider(t, provider)
}
