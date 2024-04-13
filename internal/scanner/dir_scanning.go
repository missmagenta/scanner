package scanner

import (
	"context"
	"github.com/spf13/viper"
	"log"
	"os"
	"sync"
	"time"
)

func CheckUnprocessedFiles(
	mainConfig *viper.Viper,
	wg *sync.WaitGroup,
	ctx context.Context,
	unprocessed chan string) {

	defer wg.Done()

	directory := mainConfig.GetString("path_to_scan")
	frequencyOfChecks := time.Duration(mainConfig.GetInt("frequency_of_dir_checks"))

	for {
		select {
		case <-ctx.Done():
			log.Printf("Завершение проверки директории %s", directory)
			return
		default:
			log.Printf("Происходит сканирование директории %s", directory)
			err := checkUnprocessedFiles(directory, unprocessed)
			if err != nil {
				return
			}
			time.Sleep(frequencyOfChecks * time.Second)
		}
	}
}

func checkUnprocessedFiles(directory string, unprocessed chan string) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Println("Ошибка при чтении директории:", err)
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()
			if isTsv(fileName) && !isFileInProcessing(fileName) && !isFileProcessed(fileName) {
				newPath, _ := renameToInProcess(directory + "/" + fileName)
				log.Printf("%s помечена как in_process", fileName)
				unprocessed <- newPath
			}
		}
	}

	return nil
}
