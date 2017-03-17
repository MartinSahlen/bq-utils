package bqutils

import (
	"encoding/json"

	"cloud.google.com/go/bigquery"
)

func WriteNdJSONFile(filename string, rows *bigquery.RowIterator) error {

	w, err := GetWriter(filename)

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

	return w.Flush()
}
