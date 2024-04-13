package impl

import (
	"biocad_internship/internal/model"
	"database/sql"
	"fmt"
	_ "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"log"
	"os"
	"path/filepath"
	"time"
)

type ClickHouse struct {
	Conn          *sql.DB
	FileTableName string
	GUIDTableName string
}

func (ch *ClickHouse) Init(config *viper.Viper) {
	host := config.GetString("clickhouse_host")
	port := config.GetInt("clickhouse_port")

	envFile := "clickhouse_credentials.env"
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %s", err)
	}
	username := os.Getenv("CLICKHOUSE_USERNAME")
	password := os.Getenv("CLICKHOUSE_PASSWORD")

	conn, err := sql.Open(
		"clickhouse",
		fmt.Sprintf(
			"http://%s:%d?username=%s&password=%s", host, port, username, password))
	if err != nil {
		log.Fatal(err)
	}
	rows, err := conn.Query("SELECT version()")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			log.Fatal(err)
		}

		log.Println("Connected to Clickhouse")
		log.Println("Clickhouse version " + version)
	}

	ch.Conn = conn
	ch.FileTableName = config.GetString("files_table")
	ch.GUIDTableName = config.GetString("guid_table")
}

func (ch *ClickHouse) InsertFile(FileName string, state string, error sql.NullString) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (filename, status, date_processed, error) VALUES (?, ?, ?, ?)",
		ch.FileTableName)

	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println("Ошибка загрузки часового пояса:", err)
		return err
	}

	_, err = ch.Conn.Exec(query, FileName, state, time.Now().In(location).UTC(), error)
	if err != nil {
		log.Println("Ошибка вставки данных:", err)
		return err
	}

	return nil
}

func (ch *ClickHouse) InsertMessageData(msg model.Message) error {
	query := fmt.Sprintf("INSERT INTO %s VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", ch.GUIDTableName)

	_, err := ch.Conn.Exec(query,
		msg.Mqtt,
		msg.Invent,
		msg.UnitGUID,
		msg.MsgID,
		msg.Text,
		msg.Context,
		msg.Class,
		msg.Level,
		msg.Area,
		msg.Addr,
		msg.Block,
		msg.Type,
		msg.Bit,
		msg.InvertBit)

	if err != nil {
		log.Println("Ошибка вставки данных из содержимого файла:", err)
		return err
	}

	return nil
}

func (ch *ClickHouse) HandleFileAndData(messages []model.Message, filePath string, parseError error) error {
	for _, msg := range messages {
		err := ch.InsertMessageData(msg)
		if err != nil {
			return err
		}
	}

	_, fileName := filepath.Split(filePath)

	status := fmt.Sprintf("%d сообщений", len(messages))

	var errorValue sql.NullString
	if parseError != nil {
		errorValue = sql.NullString{String: parseError.Error(), Valid: true}
	} else {
		errorValue = sql.NullString{Valid: false}
	}

	err := ch.InsertFile(fileName, status, errorValue)
	if err != nil {
		return err
	}

	log.Println("Данные помещены в Clickhouse")

	return nil
}

func (ch *ClickHouse) GetById(UnitGUID string, pageNumber, pageSize int) ([]model.Message, int64) {
	var data []model.Message
	var total int64

	offset := (pageNumber - 1) * pageSize

	query := fmt.Sprintf("SELECT * FROM %s WHERE UnitGUID = ? ORDER BY MsgID LIMIT ? OFFSET ?", ch.GUIDTableName)

	rows, err := ch.Conn.Query(query, UnitGUID, pageSize, offset)
	if err != nil {
		log.Println("Ошибка запроса:", err)
		return nil, 0
	}

	defer rows.Close()

	for rows.Next() {
		var msg model.Message
		queryErr := rows.Scan(
			&msg.Mqtt,
			&msg.Invent,
			&msg.UnitGUID,
			&msg.MsgID,
			&msg.Text,
			&msg.Context,
			&msg.Class,
			&msg.Level,
			&msg.Area,
			&msg.Addr,
			&msg.Block,
			&msg.Type,
			&msg.Bit,
			&msg.InvertBit,
		)

		if queryErr != nil {
			log.Println("Ошибка сканирования значений из строки результата запроса:", queryErr)
			continue
		}

		data = append(data, msg)
	}

	countQuery := fmt.Sprintf("SELECT count() FROM %s WHERE UnitGUID = ?", ch.GUIDTableName)
	err = ch.Conn.QueryRow(countQuery, UnitGUID).Scan(&total)
	if err != nil {
		log.Println("Ошибка запроса для получения количества:", err)
		return nil, 0
	}

	return data, total
}

func (ch *ClickHouse) AllGetById(UnitGUID string) []model.Message {
	var data []model.Message

	query := fmt.Sprintf("SELECT * FROM %s WHERE UnitGUID = ?", ch.GUIDTableName)

	rows, err := ch.Conn.Query(query, UnitGUID)
	if err != nil {
		log.Println("Ошибка запроса:", err)
		return nil
	}

	defer rows.Close()

	for rows.Next() {
		var msg model.Message
		err := rows.Scan(
			&msg.Mqtt,
			&msg.Invent,
			&msg.UnitGUID,
			&msg.MsgID,
			&msg.Text,
			&msg.Context,
			&msg.Class,
			&msg.Level,
			&msg.Area,
			&msg.Addr,
			&msg.Block,
			&msg.Type,
			&msg.Bit,
			&msg.InvertBit,
		)

		if err != nil {
			log.Println("Error scanning row:", err)
			continue
		}

		data = append(data, msg)
	}

	return data
}
