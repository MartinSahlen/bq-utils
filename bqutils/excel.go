package bqutils

import (
	"errors"

	"cloud.google.com/go/bigquery"
	"github.com/tealeg/xlsx"
)

type SheetConfig struct {
	Rows      *bigquery.RowIterator
	SheetName string
	Schema    bigquery.Schema
}

func WriteExcelFile(filename string, sheets []SheetConfig) error {

	excelFile := xlsx.NewFile()

	w, err := GetWriter(filename)

	if err != nil {
		return err
	}

	for _, s := range sheets {
		sheet, err := excelFile.AddSheet(s.SheetName)

		if err != nil {
			return err
		}

		header := []string{}

		for _, f := range s.Schema {
			header = append(header, f.Name)
		}

		writeExcelRow(header, sheet)

		mapper := func(row map[string]bigquery.Value, schema *bigquery.Schema) (map[string]bigquery.Value, error) {
			if schema == nil {
				return nil, errors.New("Schema is nil")
			}
			return nil, writeExcelRow(mapToStringSlice(row, *schema), sheet)
		}

		err = MapRows(s.Rows, &s.Schema, mapper)

		if err != nil {
			return err
		}
	}

	return excelFile.Write(w)
}

func writeExcelRow(row []string, sheet *xlsx.Sheet) error {
	r := sheet.AddRow()
	for _, cell := range row {
		c := r.AddCell()
		c.Value = cell
	}
	return nil
}
