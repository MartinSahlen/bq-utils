package bqutils

import (
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/bigquery"
)

func WriteCsvFile(filename string, rows *bigquery.RowIterator, schema bigquery.Schema) error {

	w, err := GetWriter(filename)

	if err != nil {
		return err
	}

	header := []string{}

	for _, f := range schema {
		header = append(header, f.Name)
	}

	fmt.Fprintln(w, strings.Join(header, ","))

	mapper := func(row map[string]bigquery.Value, schema *bigquery.Schema) (map[string]bigquery.Value, error) {
		if schema == nil {
			return nil, errors.New("Schema is nil")
		}
		_, writeErr := w.Write([]byte(strings.Join(mapToStringSlice(row, *schema), ",") + "\n"))
		return nil, writeErr
	}

	err = MapRows(rows, &schema, mapper)

	if err != nil {
		return err
	}

	return w.Flush()
}

func TableToCsv(project, tablename, filename string) error {

	dataset, table, err := ParseTableName(tablename)

	if err != nil {
		return err
	}

	tableData, err := GetTableData(project, *dataset, *table)

	if err != nil {
		return err
	}

	return WriteCsvFile(filename, tableData.Rows, tableData.Schema)
}

func QueryToCsv(project, query, filename string) error {
	queryData, err := GetQueryData(project, query)

	if err != nil {
		return err
	}

	return WriteCsvFile(filename, queryData.Rows, queryData.Schema)
}
