package export

import (
	"NULL/knowledgebase/models"
	"NULL/knowledgebase/pkg/setting"
	"NULL/knowledgebase/pkg/upload"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"strconv"
	"strings"
	"time"
)

const EXT = ".xlsx"

// GetExcelFullUrl get the full access path of the Excel file
func GetExcelFullUrl(name string) string {
	return setting.AppSetting.PrefixUrl + "/" + GetExcelPath() + name
}

// GetExcelPath get the relative save path of the Excel file
func GetExcelPath() string {
	return setting.AppSetting.ExportSavePath
}

// GetExcelFullPath Get the full save path of the Excel file
func GetExcelFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetExcelPath()
}

// Write into excel
func WriteIntoExcel(fileName string, records []map[string]string) (string, error) {
	var sheetName = "Sheet1"
	xlsx := excelize.NewFile()
	index := xlsx.NewSheet(sheetName)
	// sort map key
	sorted_keys := models.SortFields(records[0])
	/*
		sorted_keys := make([]string, 0)
		for filed := range records[0] {
			sorted_keys = append(sorted_keys, filed)
		}
	*/
	// Set table head
	var A = 'A'
	for i, field := range sorted_keys {
		var cell string
		if i < 26 {
			cell = fmt.Sprintf("%c", A)
			A++
		} else if i == 26 {
			A = 'A'
			cell = fmt.Sprintf("A%c", A)
			A++
		} else {
			cell = fmt.Sprintf("A%c", A)
			A++
		}
		xlsx.SetCellValue(sheetName, cell+"1", models.ReplaceTableFileds(field))
		xlsx.SetColWidth(sheetName, cell, cell, countWidth(field))
	}
	// Set cell value
	for row, record := range records {
		var A = 'A'
		for i, field := range sorted_keys {
			var cell string
			if i < 26 {
				cell = fmt.Sprintf("%c", A) + strconv.Itoa(row+2)
				A++
			} else if i == 26 {
				A = 'A'
				cell = fmt.Sprintf("A%c", A) + strconv.Itoa(row+2)
				A++
			} else {
				cell = fmt.Sprintf("A%c", A) + strconv.Itoa(row+2)
				A++
			}
			xlsx.SetCellValue(sheetName, cell, record[field])
		}
	}
	// Set active sheet of the workbook
	xlsx.SetActiveSheet(index)
	// Save xlsx file by the given path
	savePath := GetExcelFullPath()
	if err := upload.CheckFile(savePath); err != nil {
		return "", err
	}
	saveName := fileName + strconv.Itoa(int(time.Now().Unix())) + EXT
	scr := savePath + saveName
	if err := xlsx.SaveAs(scr); err != nil {
		return "", err
	}
	return GetExcelFullUrl(saveName), nil
}

// according to cell value, count colwidth
func countWidth(src string) float64 {
	letters := "abcdefghijklmnopqrstuvwxyz"
	letters = letters + strings.ToUpper(letters)
	nums := "0123456789"
	chars := "()$%*@+-=/#"

	numCount := 0
	letterCount := 0
	othersCount := 0
	charsCount := 0

	for _, i := range src {
		switch {
		case strings.ContainsRune(letters, i) == true:
			letterCount += 1
		case strings.ContainsRune(nums, i) == true:
			numCount += 1
		case strings.ContainsRune(chars, i) == true:
			charsCount += 1
		default:
			othersCount += 1
		}
	}

	return float64(numCount+letterCount+charsCount+othersCount*2) * 4
}
