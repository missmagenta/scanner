package file_handler

import (
	"biocad_internship/internal/model"
	"encoding/csv"
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseFile(filePath string) (messages []model.Message, err error) {
	var level, bit, invertBit int

	file, err := os.Open(filePath)
	if err != nil {
		log.Println("Ошибка при открытии файла:", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = '\t'

	lines, err := reader.ReadAll()
	if err != nil {
		log.Println("Ошибка при чтении файла:", err)
	}

	if len(lines) == 0 {
		return messages, errors.New("Файл пуст")
	}

	for numberLine, line := range lines {
		if numberLine == 0 {
			continue
		}

		for i, _ := range line {
			line[i] = strings.Trim(line[i], " ")
		}

		if len(line) > 14 {
			level, bit, invertBit, err = transformSpecificFields(line[8], line[13], line[14])
			if err != nil {
				log.Println("Ошибка при приведении типов n|level|bit файл:"+filePath+": ", err)
			}

			currentMessage := model.Message{
				Mqtt:      line[1],
				Invent:    line[2],
				UnitGUID:  line[3],
				MsgID:     line[4],
				Text:      line[5],
				Context:   line[6],
				Class:     line[7],
				Level:     level,
				Area:      line[9],
				Addr:      line[10],
				Block:     line[11],
				Type:      line[12],
				Bit:       bit,
				InvertBit: invertBit,
			}

			messages = append(messages, currentMessage)

		} else {
			parseErr := errors.New(fmt.Sprintf(
				"Ошибка парсинга файла. Недостаточно данных в строке %d для доступа к индексам 1-14", numberLine))
			return messages, parseErr
		}
	}

	return messages, err
}

func transformSpecificFields(_level, _bit, _invertBit string) (level, bit, invertBit int, err error) {
	level, err = strconv.Atoi(_level)
	if err != nil {
		return level, bit, invertBit, err
	}

	if _bit == "" {
		bit = 1
	} else {
		bit, err = strconv.Atoi(_bit)
	}
	if err != nil {
		return level, bit, invertBit, err
	}

	if _invertBit == "" {
		invertBit = -1
	} else {
		invertBit, err = strconv.Atoi(_invertBit)
	}
	if err != nil {
		return level, bit, invertBit, err
	}
	return level, bit, invertBit, err
}
