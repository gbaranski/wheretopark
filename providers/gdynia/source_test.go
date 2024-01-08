package gdynia_test

import (
	"testing"
	wheretopark "wheretopark/go"
	"wheretopark/providers/gdynia"
)

func TestGdynia(t *testing.T) {
	wheretopark.ExamineSource(t, gdynia.New())
}
