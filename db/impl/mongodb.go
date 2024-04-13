package impl

import (
	"biocad_internship/internal/model"
	"context"
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"path/filepath"
	"time"
)

type MongoDB struct {
	GUIDCollection *mongo.Collection
	FileCollection *mongo.Collection
	client         *mongo.Client
	ctx            context.Context
}

type FileStatus struct {
	FileName      string    `bson:"filename"`
	Status        string    `bson:"status"`
	DateProcessed time.Time `bson:"date_processed"`
	Error         string    `bson:"error"`
}

func (db *MongoDB) Init(config *viper.Viper) {
	var mongoUri = config.GetString("mongo_host")

	envFile := "mongodb_credentials.env"
	if err := godotenv.Load(envFile); err != nil {
		log.Fatalf("Ошибка загрузки файла .env: %s", err)
	}
	clientOptions := options.Client().ApplyURI(mongoUri)
	clientOptions.Auth = &options.Credential{
		Username: os.Getenv("MONGO_USERNAME"),
		Password: os.Getenv("MONGO_PASSWORD"),
	}

	db.ctx = context.TODO()
	var err error

	db.client, err = mongo.Connect(db.ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = db.client.Ping(db.ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Connected to mongodb")

	db.FileCollection = db.client.Database(config.GetString("mongo_database")).
		Collection(config.GetString("files_collection"))
	db.GUIDCollection = db.client.Database(config.GetString("mongo_database")).
		Collection(config.GetString("guid_collection"))
}

func (db *MongoDB) InsertFile(FileName string, state string, error sql.NullString) error {
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		log.Println("Ошибка загрузки часового пояса:", err)
		return err
	}

	data := FileStatus{
		FileName:      FileName,
		Status:        state,
		DateProcessed: time.Now().In(location).UTC(),
		Error:         error.String,
	}

	log.Printf("File %s FileStatus %s\n", FileName, data)

	messageBSON, err := bson.Marshal(data)
	if err != nil {
		log.Println("Ошибка маршалинга данных message", err)
		return err
	}

	_, err = db.FileCollection.InsertOne(context.Background(), messageBSON)
	if err != nil {
		log.Println("Ошибка вставки данных message", err)
		return err
	}

	return nil
}

func (db *MongoDB) InsertMessageData(message model.Message) error {
	messageBSON, err := bson.Marshal(message)
	if err != nil {
		log.Println("Ошибка маршалинга данных message", err)
		return err
	}

	_, err = db.GUIDCollection.InsertOne(context.Background(), messageBSON)
	if err != nil {
		log.Println("Ошибка вставки данных message", err)
		return err
	}

	return nil
}

func (db *MongoDB) HandleFileAndData(messages []model.Message, filePath string, parseError error) (err error) {
	for _, msg := range messages {
		err = db.InsertMessageData(msg)
		if err != nil {
			return err
		}
	}

	_, fileName := filepath.Split(filePath)
	status := fmt.Sprint(len(messages)) + " сообщений"
	if err != nil {
		status = fmt.Sprint(err)
	}

	var errorValue sql.NullString
	if parseError != nil {
		errorValue = sql.NullString{String: parseError.Error(), Valid: true}
	} else {
		errorValue = sql.NullString{Valid: false}
	}

	errInsert := db.InsertFile(fileName, status, errorValue)
	if errInsert != nil {
		return errInsert
	}
	log.Println("Данные помещены в коллекции MongoDB")

	return nil
}

func (db *MongoDB) GetById(UnitGUID string, pageNumber, pageSize int) (data []model.Message, total int64) {
	total, err := db.GUIDCollection.CountDocuments(context.Background(),
		bson.M{"unitguid": UnitGUID})
	if err != nil {
		log.Println(err)
		return nil, total
	}

	offset := int64((pageNumber - 1) * pageSize)

	cur, err := db.GUIDCollection.Find(
		context.Background(),
		bson.M{"unitguid": UnitGUID},
		options.Find().SetSkip(offset).SetLimit(int64(pageSize)))
	if err != nil {
		return nil, 0
	}

	defer cur.Close(context.Background())

	if err = cur.All(context.Background(), &data); err != nil {
		log.Fatal(err)
	}

	return data, total
}

func (db *MongoDB) AllGetById(UnitGUID string) (data []model.Message) {
	cur, err := db.GUIDCollection.Find(context.Background(),
		bson.M{"unitguid": UnitGUID})

	defer cur.Close(context.Background())

	if err != nil {
		return nil
	}

	if err = cur.All(context.Background(), &data); err != nil {
		log.Fatal(err)
	}
	return data
}
