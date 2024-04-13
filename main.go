package main

import (
	"biocad_internship/config"
	"biocad_internship/db"
	"biocad_internship/db/impl"
	"biocad_internship/internal/file_handler"
	"biocad_internship/internal/scanner"
	"biocad_internship/server"
	"context"
	"github.com/spf13/viper"
	"sync"
)

func main() {
	ctx := context.Background()
	mainConfig, dbConfig, pdfConfig := config.Init()
	var database db.Repository
	database = dbInit(dbConfig)

	var wg sync.WaitGroup
	notProcessedFiles := make(chan string, 5)

	wg.Add(1)
	go scanner.CheckUnprocessedFiles(mainConfig, &wg, ctx, notProcessedFiles)

	wg.Add(1)
	go file_handler.FilesProcessing(pdfConfig, notProcessedFiles, &wg, database)

	apiServer := server.Server{}
	go apiServer.Init(&wg, database, mainConfig)

	wg.Wait()
}

func dbInit(dbConfig *viper.Viper) (dataBase db.Repository) {
	switch dbConfig.GetString("db_impl") {
	case "mongodb":
		dataBase = &impl.MongoDB{}
	case "clickhouse":
		dataBase = &impl.ClickHouse{}
	default:
		panic("Неверный аргумент типа базы данных")
	}
	dataBase.Init(dbConfig)
	return dataBase
}
