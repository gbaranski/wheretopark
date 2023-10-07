package gdynia_test

import (
	"testing"
	"wheretopark/collector/gdynia"
	wheretopark "wheretopark/go"
)

func TestGdynia(t *testing.T) {
	wheretopark.ExamineSource(t, gdynia.New())
}
