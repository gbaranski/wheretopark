package poznan_test

import (
	"testing"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/poznan"
)

func TestPoznan(t *testing.T) {
	wheretopark.ExamineSource(t, poznan.New())
}