package scanner

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

func isTsv(fileName string) bool {
	return path.Ext(fileName) == ".tsv"
}

func renameToInProcess(filePath string) (newPath string, err error) {
	directory, fileName := filepath.Split(filePath)
	newFilename := fmt.Sprintf("in_process_%s", fileName)
	newPath = filepath.Join(directory, newFilename)

	err = os.Rename(filePath, newPath)
	if err != nil {
		return newPath, err
	}

	return newPath, nil
}

func isFileInProcessing(filename string) bool {
	return len(filename) >= 11 && filename[:11] == "in_process_"
}

func isFileProcessed(filename string) bool {
	return len(filename) >= 10 && filename[:10] == "processed_"
}
