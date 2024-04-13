package file_handler

import (
	"biocad_internship/db"
	"biocad_internship/internal/pdf/utils"
	"github.com/spf13/viper"
	"log"
	"sync"
)

func FilesProcessing(PDFConfig *viper.Viper, unProcessing chan string, wg *sync.WaitGroup, db db.Repository) {
	defer wg.Done()

	for filePath := range unProcessing {
		log.Println("Принят в обработку " + filePath)
		messages, parseErr := parseFile(filePath)
		processedFilePath, renameErr := renameToProcessed(filePath)

		if renameErr != nil {
			log.Printf("Ошибка переименования файла %s: %s", filePath, parseErr)
			return
		}

		if parseErr != nil {
			log.Printf("В файле %s обнаружена ошибка: %s", filePath, parseErr)
		}

		dbErr := db.HandleFileAndData(messages, processedFilePath, parseErr)
		if dbErr != nil {
			log.Println("Ошибка обработки файла и данных в базе данных:", dbErr)
			continue
		}

		go func() {
			err := makePdfReport(utils.FindUniqueGuids(messages), PDFConfig, db, processedFilePath)
			if err != nil {
				log.Println(err)
			}
		}()
	}
}
