package glasgow_test

import (
	"testing"
	wheretopark "wheretopark/go"
	"wheretopark/providers/collector/glasgow"
)

func TestGlasgow(t *testing.T) {
	wheretopark.ExamineSource(t, glasgow.New())
}
