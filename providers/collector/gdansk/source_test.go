package gdansk_test

import (
	"testing"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/gdansk"
)

func TestGdansk(t *testing.T) {
	wheretopark.ExamineSequentialSource(t, gdansk.New())
}
