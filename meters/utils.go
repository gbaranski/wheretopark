package meters

import (
	"os"
	"path/filepath"
	"strings"
)

func listFilesWithExtension(basePath string, extension string) ([]string, error) {
	var filesWithExtension []string

	// Read all files and directories within basePath
	items, err := os.ReadDir(basePath)
	if err != nil {
		return nil, err
	}

	// Iterate through the items
	for _, item := range items {
		if !item.IsDir() { // Ensure it's a file, not a directory
			// Check if the file has the desired extension
			if strings.HasSuffix(item.Name(), "."+extension) {
				fullPath := filepath.Join(basePath, item.Name())
				filesWithExtension = append(filesWithExtension, fullPath)
			}
		}
	}

	return filesWithExtension, nil
}
