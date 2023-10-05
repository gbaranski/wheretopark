package poznan_test

import (
	"testing"
	"wheretopark/collector/poznan"
	wheretopark "wheretopark/go"
)

func TestPoznan(t *testing.T) {
	wheretopark.ExamineSource(t, poznan.New())
}
