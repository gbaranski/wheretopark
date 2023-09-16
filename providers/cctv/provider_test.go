package cctv_test

import (
	"fmt"
	"os"
	"testing"
	wheretopark "wheretopark/go"
	cctv "wheretopark/providers/cctv"
)

func TestProvider(t *testing.T) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}
	modelPath := fmt.Sprintf("%s/.local/share/wheretopark/cctv/model.onnx", homeDirectory)

	provider, err := cctv.NewProvider(modelPath, nil, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer provider.Close()
	wheretopark.ExamineProvider(t, provider)
}
