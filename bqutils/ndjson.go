package bqutils

import (
	"encoding/json"

	"cloud.google.com/go/bigquery"
)

func WriteNdJSONFile(filename string, rows *bigquery.RowIterator) error {

	w, err := GetWriter(filename, "application/ndjson")

	if err != nil {
		return err
	}

	mapper := func(row map[string]bigquery.Value, schema *bigquery.Schema) (map[string]bigquery.Value, error) {
		line, writeErr := json.Marshal(row)
		if writeErr != nil {
			return nil, writeErr
		}
		_, writeErr = w.Write(append(line, []byte("\n")...))
		return nil, writeErr
	}

	err = MapRows(rows, nil, mapper)

	if err != nil {
		return err
	}

	err = w.Close()

	if err != nil {
		return err
	}
	return nil
}

func QueryToNdJSON(project, query, filename string) error {
	queryData, err := GetQueryData(project, query)

	if err != nil {
		return err
	}

	return WriteNdJSONFile(filename, queryData.Rows)
}

func TableToNdJSON(project, tablename, filename string) error {
	dataset, table, err := ParseTableName(tablename)

	if err != nil {
		return err
	}

	queryData, err := GetTableData(project, *dataset, *table)

	if err != nil {
		return err
	}

	return WriteNdJSONFile(filename, queryData.Rows)
}
