package glasgow_test

import (
	"testing"
	"wheretopark/go/provider/simple"
	"wheretopark/providers/collector/glasgow"
)

func TestProvider(t *testing.T) {
	provider, err := glasgow.NewProvider()
	if err != nil {
		t.Fatal(err)
	}
	simple.ExamineProvider(t, provider)
}
