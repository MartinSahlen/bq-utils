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

		for _, f := range s.Schema {
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

		err = MapRows(s.Rows, &s.Schema, mapper)

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

//StitchSheetNames : Since we get two arrays for queries and their corresponding sheet names,
//We need to "stitch" them together. Docopt guarantees that this will not blast
//Because the slices will have the same length
func StitchSheetNames(queriesOrTables, sheetNames []string, project string, isQuery bool) []SheetWriterConfig {
	writeConfigs := []SheetWriterConfig{}

	for i, queryOrTable := range queriesOrTables {

		writeConfig := SheetWriterConfig{
			SheetName: sheetNames[i],
			Project:   project,
		}

		if isQuery {
			writeConfig.Query = queryOrTable
			writeConfig.IsQuery = true
		} else {
			writeConfig.Table = queryOrTable
			writeConfig.IsQuery = false
		}
		writeConfigs = append(writeConfigs, writeConfig)
	}
	return writeConfigs
}

type SheetWriterConfig struct {
	IsQuery   bool
	Project   string
	Query     string
	Table     string
	SheetName string
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

	sheets := []SheetConfig{}

	for _, s := range ss {
		rowData, err := s.Execute()

		if err != nil {
			return err
		}

		sheets = append(sheets, SheetConfig{
			SheetName: s.SheetName,
			Schema:    rowData.Schema,
			Rows:      rowData.Rows,
		})
	}
	return WriteExcelFile(filename, sheets)
}
