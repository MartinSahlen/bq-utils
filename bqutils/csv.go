package bqutils

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
)

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
		_, writeErr := w.Write([]byte(strings.Join(mapToStringSlice(row, *schema), ",") + "\n"))
		return nil, writeErr
	}

	err = MapRows(rows, &schema, mapper)

	if err != nil {
		return err
	}

	return w.Flush()
}
