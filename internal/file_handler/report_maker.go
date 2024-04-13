package file_handler

import (
	"biocad_internship/db"
	"biocad_internship/internal/pdf"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
)

func makePdfReport(newGuids []string, PDFConfig *viper.Viper, db db.Repository, filePath string) error {
	if len(newGuids) == 0 {
		err := makeReportForNoGuids(filePath, PDFConfig)
		if err != nil {
			return err
		}
		return nil
	}

	for _, newGuid := range newGuids {
		messages := db.AllGetById(newGuid)
		if len(messages) != 0 {
			headers := pdf.ObjHeadToStrArr(messages[0])
			var data [][]string
			for _, m := range messages {
				data = append(data, pdf.ConvertObjectToStrArr(m))
			}
			err := pdf.MakePdfReport(PDFConfig, newGuid, headers, data, filePath)
			if err != nil {
				return fmt.Errorf("ошибка создания отчета: %w", err)
			}
			log.Printf("Отчет по unit_guid %s сформирован\n", newGuid)

		} else {
			log.Println("No messages")
			err := makeReportForNoMessages(filePath, PDFConfig)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func makeReportForNoGuids(filePath string, PDFConfig *viper.Viper) error {
	headers := []string{"Информация"}
	data := [][]string{{"GUIDs не найдено при парсинге файла."}}
	_, filename := filepath.Split(filePath)
	emptyGuidsMessage := "empty_new_guids_" + filename
	err := pdf.MakePdfReport(PDFConfig, emptyGuidsMessage, headers, data, filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания отчета: %w", err)
	}
	log.Printf("Отчет по файлу %s сформирован\n", filename)

	return nil
}

func makeReportForNoMessages(filePath string, PDFConfig *viper.Viper) error {
	headers := []string{"Информация"}
	data := [][]string{{"Новых сообщений по данному guid не обануржено."}}
	_, filename := filepath.Split(filePath)
	informMessage := "no_messages_" + filename
	err := pdf.MakePdfReport(PDFConfig, informMessage, headers, data, filePath)
	if err != nil {
		return fmt.Errorf("ошибка создания отчета: %w", err)
	}
	log.Printf("Отчет по файлу %s сформирован\n", filename)

	return nil
}
