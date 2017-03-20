package bqutils

import (
	"errors"

	"cloud.google.com/go/bigquery"
	"github.com/tealeg/xlsx"
)

func WriteExcelFile(filename string, sheets []SheetConfig) error {

	excelFile := xlsx.NewFile()

	w, err := GetWriter(filename, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

	if err != nil {
		return err
	}

	for _, s := range sheets {
		sheet, err := excelFile.AddSheet(s.SheetName)

		if err != nil {
			return err
		}

		header := []string{}

		for _, f := range s.RowData.Schema {
			header = append(header, f.Name)
		}

		writeExcelRow(header, sheet)

		mapper := func(row map[string]bigquery.Value, schema *bigquery.Schema) (map[string]bigquery.Value, error) {
			if schema == nil {
				return nil, errors.New("Schema is nil")
			}
			writeExcelRow(mapToStringSlice(row, *schema), sheet)
			return nil, nil
		}

		err = MapRows(s.RowData.Rows, &s.RowData.Schema, mapper)

		if err != nil {
			return err
		}
	}

	err = excelFile.Write(w)

	if err != nil {
		return err
	}

	err = w.Close()

	if err != nil {
		return err
	}
	return nil
}

func writeExcelRow(row []string, sheet *xlsx.Sheet) {
	r := sheet.AddRow()
	for _, cell := range row {
		c := r.AddCell()
		c.Value = cell
	}
}

func (e SheetWriterConfig) Execute() (*RowData, error) {
	if e.IsQuery {
		return GetQueryData(e.Project, e.Query)
	}

	dataset, table, err := ParseTableName(e.Table)

	if err != nil {
		return nil, err
	}

	return GetTableData(e.Project, *dataset, *table)
}

func WriteToExcel(ss []SheetWriterConfig, filename string) error {

	sheets, err := sheetConfigToWriterConfig(ss)

	if err != nil {
		return err
	}

	return WriteExcelFile(filename, sheets)
}
