package poznan_test

import (
	"testing"
	"wheretopark/go/tester"
	"wheretopark/providers/poznan"
)

func TestPoznan(t *testing.T) {
	tester.ExamineProvider(t, poznan.New())
}
