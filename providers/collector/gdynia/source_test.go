package gdynia_test

import (
	"testing"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/gdynia"
)

func TestGdynia(t *testing.T) {
	wheretopark.ExamineSequentialSource(t, gdynia.New())
}
