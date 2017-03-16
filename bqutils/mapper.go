package bqutils

import (
	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"
)

type MapRow func(row map[string]bigquery.Value, schema *bigquery.Schema) (map[string]bigquery.Value, error)

func MapRows(rows *bigquery.RowIterator, schema *bigquery.Schema, mapFunc MapRow) error {
	for {
		_, done, err := mapRows(rows, schema, mapFunc)
		if done {
			break
		} else if err != nil {
			return err
		}
	}
	return nil
}

func mapRows(rows *bigquery.RowIterator, schema *bigquery.Schema, mapFunc MapRow) (map[string]bigquery.Value, bool, error) {
	row := map[string]bigquery.Value{}
	err := rows.Next(&row)

	if err == iterator.Done {
		return nil, true, nil
	}

	if err != nil {
		return nil, false, err
	}

	row, err = mapFunc(row, schema)

	if err != nil {
		return row, false, err
	}
	return row, false, nil
}
