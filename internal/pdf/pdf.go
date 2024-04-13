package pdf

import (
	"fmt"
	"github.com/jung-kurt/gofpdf"
	"github.com/spf13/viper"
	"log"
	"os"
	"reflect"
	"unicode/utf8"
)

func ConvertObjectToStrArr(msg interface{}) (result []string) {
	val := reflect.ValueOf(msg)
	typ := reflect.TypeOf(msg)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name
		fieldValue := fmt.Sprintf("%v", field.Interface())
		if fieldName != "UnitGUID" {
			result = append(result, fieldValue)
		}
	}
	return result
}

func ObjHeadToStrArr(msg interface{}) (result []string) {
	val := reflect.ValueOf(msg)
	typ := reflect.TypeOf(msg)

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Name
		if fieldName != "UnitGUID" {
			result = append(result, fieldName)
		}
	}
	return result
}

func MakePdfReport(
	PDFConfig *viper.Viper,
	Header string,
	headers []string,
	data [][]string,
	filePath string) error {

	pdf := gofpdf.New("L", "mm", "A4", "")
	leftMargin := 10.0
	rightMargin := 10.0
	pdf.SetMargins(leftMargin, 10, rightMargin)
	pdf.SetFooterFunc(func() {

		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)

		pdf.CellFormat(10, 10, fmt.Sprint(pdf.PageNo()), "", 0, "L", false, 0, "") // Page number
	})
	fontData, err := os.ReadFile(PDFConfig.GetString("font_path"))
	if err != nil {
		log.Println("Ошибка при чтении файла шрифта:", err)
		return err
	}
	pdf.AddUTF8FontFromBytes("ArialUnicode", "", fontData)

	pdf.AddPage()
	{
		pdf.SetFont(PDFConfig.GetString("headerFont"), "", PDFConfig.GetFloat64("headerSize"))

		l1 := pdf.AddLayer("Layer 1", true)
		pdf.BeginLayer(l1)
		pdf.Write(8, Header+"\n")
		pdf.EndLayer()
	}
	{
		pdf.SetFont("ArialUnicode", "", 12)
		l2 := pdf.AddLayer("Layer 2", true)
		pdf.BeginLayer(l2)
		pdf.Write(7, "Путь к файлу: "+filePath+"\n")
		pdf.EndLayer()
	}

	pdf.SetFont(PDFConfig.GetString("tableFont"), "", PDFConfig.GetFloat64("tableHeaderSize"))

	var sizeArr []float64
	{
		for i := range headers {
			maxLen := utf8.RuneCountInString(headers[i])
			for _, row := range data {
				cellLen := utf8.RuneCountInString(row[i])

				if cellLen > maxLen {
					maxLen = cellLen
				}
			}

			sizeArr = append(sizeArr, float64(maxLen))

		}
	}

	{
		for i, cell := range headers {
			cellWidth := sizeArr[i] * 2.05
			pdf.CellFormat(
				cellWidth,
				7,
				cell,
				PDFConfig.GetString("borderType"),
				0,
				"1",
				false,
				0,
				"")
		}
		pdf.Ln(-1)

		pdf.SetFont(PDFConfig.GetString("tableFont"), "", PDFConfig.GetFloat64("tableSize"))
		for _, row := range data {
			for i, cell := range row {
				cellWidth := sizeArr[i] * 2.05
				pdf.CellFormat(
					cellWidth,
					7,
					cell,
					PDFConfig.GetString("borderType"),
					0,
					"1",
					false,
					0,
					"")
			}

			pdf.Ln(-1)
		}
	}

	err = pdf.OutputFileAndClose(PDFConfig.GetString("output_path") + "/" + Header + ".pdf")
	if err != nil {
		log.Println("Ошибка при сохранении файла:", err)
		return err
	}
	return err
}
