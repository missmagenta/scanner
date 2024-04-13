package db

import (
	"biocad_internship/internal/model"
	"database/sql"
	"github.com/spf13/viper"
)

type Repository interface {
	Init(config *viper.Viper)
	InsertFile(FileName string, state string, error sql.NullString) error
	InsertMessageData(message model.Message) error
	HandleFileAndData(messages []model.Message, filePath string, parseError error) error
	GetById(UnitGUID string, pageNumber, pageSize int) ([]model.Message, int64)
	AllGetById(UnitGUID string) []model.Message
}
