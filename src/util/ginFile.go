package util

import (
	"bytes"
	"encoding/csv"
	"errors"
	"io"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/tealeg/xlsx"
)

func GinGridFileParser(c *gin.Context, key string) ([]*xlsx.Sheet, error) {
	result := make([]*xlsx.Sheet, 0)

	file, head, err := c.Request.FormFile(key)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 取得檔案格式
	fileSuffix := path.Ext(head.Filename)

	switch fileSuffix {
	case ".xlsx":
		// 將資料轉成byte
		buf := bytes.NewBuffer(nil)
		if _, err = io.Copy(buf, file); err != nil {
			return nil, err
		}

		// 打開資料(儲存成xlsx的結構)
		xlsxFile, err := xlsx.OpenBinary(buf.Bytes())
		if err != nil {
			return nil, err
		}

		result = xlsxFile.Sheets
	case ".csv":
		reader := csv.NewReader(file)
		rows, err := reader.ReadAll()
		if err != nil {
			return nil, err
		}

		sheet := &xlsx.Sheet{
			Name: head.Filename,
			Rows: make([]*xlsx.Row, 0),
		}
		for _, cols := range rows {
			row := &xlsx.Row{
				Cells: make([]*xlsx.Cell, 0),
				Sheet: sheet,
			}
			for _, value := range cols {
				cell := &xlsx.Cell{
					Row:   row,
					Value: value,
				}
				row.Cells = append(row.Cells, cell)
			}
			sheet.Rows = append(sheet.Rows, row)

			sheet.MaxCol = len(cols)
		}

		sheet.MaxRow = len(sheet.Rows)

		result = append(result, sheet)
	default:
		return nil, errors.New("wrong type")
	}

	return result, nil
}
