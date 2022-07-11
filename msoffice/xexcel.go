package msoffice

import (
	"github.com/tealeg/xlsx"
)

type XlsSheet struct {
	Name     string
	ColsNum  int
	LinesNum int
	Cells    [][]string
}

//ReadAllExcelData read all cells in string format and output
func ReadAllExcelData(fileName string) ([]*XlsSheet, error) {
	xlsFile, err := xlsx.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	sheets := []*XlsSheet{}

	for _, sheet := range xlsFile.Sheets {
		name := sheet.Name
		colsNum := 0
		linesNum := len(sheet.Rows)
		cells := [][]string{}

		for _, row := range sheet.Rows {
			cols := len(row.Cells)
			if cols > colsNum {
				colsNum = cols
			}
			rowData := []string{}
			for _, cell := range row.Cells {
				val := cell.String()
				//convert special characters
				rowData = append(rowData, val)
			}
			cells = append(cells, rowData)
		}
		sheetData := &XlsSheet{
			Name:     name,
			ColsNum:  colsNum,
			LinesNum: linesNum,
			Cells:    cells,
		}
		sheets = append(sheets, sheetData)
	}
	return sheets, nil
}

func GetColIndexName(idx int) string {
	return string(rune(65 + idx))
}
