package gdansk_test

import (
	"testing"
	"wheretopark/collector/gdansk"
	wheretopark "wheretopark/go"
)

func TestGdansk(t *testing.T) {
	wheretopark.ExamineSource(t, gdansk.New())
}
