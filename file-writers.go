package bqutils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"github.com/tealeg/xlsx"
)

type SheetConfig struct {
	Rows      *bigquery.RowIterator
	SheetName string
	Schema    bigquery.Schema
}

func WriteExcelFile(filename string, sheets []SheetConfig) error {

	file := xlsx.NewFile()

	for _, s := range sheets {
		sheet, err := file.AddSheet(s.SheetName)

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
	return file.Save(filename)
}

func WriteCsvFile(fileName string, rows *bigquery.RowIterator, schema bigquery.Schema) error {

	csvFile, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer csvFile.Close()

	w := bufio.NewWriter(csvFile)

	header := []string{}

	for _, f := range schema {
		header = append(header, f.Name)
	}

	fmt.Fprintln(w, strings.Join(header, ","))

	mapper := func(row map[string]bigquery.Value, schema *bigquery.Schema) (map[string]bigquery.Value, error) {
		if schema == nil {
			return nil, errors.New("Schema is nil")
		}
		_, err := fmt.Fprintln(w, strings.Join(mapToStringSlice(row, *schema), ","))
		return nil, err
	}

	MapRows(rows, &schema, mapper)

	return w.Flush()
}

func writeExcelRow(row []string, sheet *xlsx.Sheet) error {
	r := sheet.AddRow()
	for _, cell := range row {
		c := r.AddCell()
		c.Value = cell
	}
	return nil
}

func mapToStringSlice(row map[string]bigquery.Value, schema bigquery.Schema) []string {
	outputRow := []string{}
	for _, f := range schema {
		outputRow = append(outputRow, strings.TrimSpace(fmt.Sprint(row[f.Name])))
	}
	return outputRow
}
