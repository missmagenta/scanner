package file_handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func renameToProcessed(filePath string) (newPath string, err error) {
	directory, fileName := filepath.Split(filePath)
	originalFilename := strings.TrimPrefix(fileName, "in_process_")
	newFilename := fmt.Sprintf("processed_%s", originalFilename)
	newPath = filepath.Join(directory, newFilename)

	err = os.Rename(filePath, newPath)
	if err != nil {
		return newPath, err
	}

	return newPath, nil
}
