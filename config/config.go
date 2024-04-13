package config

import (
	"github.com/spf13/viper"
	"log"
)

func Init() (mainConfig, dbConfig, PDFConfig *viper.Viper) {
	var err error
	const configFilePath = "app"
	const dbConfigFilePath = "db"
	const pdfConfigFilePath = "pdf"

	mainConfig = viper.New()
	mainConfig.SetConfigName(configFilePath)
	mainConfig.SetConfigType("yml")
	mainConfig.AddConfigPath("config/files")

	err = mainConfig.ReadInConfig()
	if err != nil {
		log.Fatalf("Ошибка при чтении конфигурационного файла: %s", err)
	}

	dbConfig = viper.New()
	dbConfig.SetConfigName(dbConfigFilePath)
	dbConfig.SetConfigType("yml")
	dbConfig.AddConfigPath("config/files")

	err = dbConfig.ReadInConfig()
	if err != nil {
		log.Fatalf("Ошибка при чтении конфигурационного файла: %s", err)
	}

	PDFConfig = viper.New()
	PDFConfig.SetConfigName(pdfConfigFilePath)
	PDFConfig.SetConfigType("yml")
	PDFConfig.AddConfigPath("config/files")
	err = PDFConfig.ReadInConfig()
	if err != nil {
		log.Fatalf("Ошибка при чтении конфигурационного файла: %s", err)
	}

	return mainConfig, dbConfig, PDFConfig
}
