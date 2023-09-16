package warsaw_test

import (
	"testing"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/warsaw"
)

func TestWarsaw(t *testing.T) {
	wheretopark.ExamineSource(t, warsaw.New())
}
