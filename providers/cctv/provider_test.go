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
	model := cctv.NewModel(modelPath)
	saver := cctv.NewSaver(nil, nil, nil)

	provider, err := cctv.NewProvider(model, saver)
	if err != nil {
		t.Fatal(err)
	}
	defer provider.Close()
	wheretopark.ExamineProvider(t, provider)
}
