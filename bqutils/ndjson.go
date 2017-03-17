package bqutils

import (
	"bufio"
	"encoding/json"
	"os"

	"cloud.google.com/go/bigquery"
)

func WriteNdJSONFile(fileName string, rows *bigquery.RowIterator) error {

	ndJsonFile, err := os.Create(fileName)

	if err != nil {
		return err
	}

	defer ndJsonFile.Close()

	w := bufio.NewWriter(ndJsonFile)

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
