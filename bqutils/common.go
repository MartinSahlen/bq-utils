package bqutils

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"cloud.google.com/go/bigquery"
	"golang.org/x/oauth2/google"
	bigqueryV2 "google.golang.org/api/bigquery/v2"
)

func client(project string) (*bigquery.Client, error) {
	return bigquery.NewClient(context.Background(), project)
}

func v2Client() (*bigqueryV2.Service, error) {
	httpClient, err := google.DefaultClient(context.Background())
	if err != nil {
		return nil, err
	}
	return bigqueryV2.New(httpClient)
}

type RowData struct {
	Rows    *bigquery.RowIterator
	NumRows uint64
	Schema  bigquery.Schema
}

func mapToStringSlice(row map[string]bigquery.Value, schema bigquery.Schema) []string {
	outputRow := []string{}
	for _, f := range schema {
		outputRow = append(outputRow, strings.TrimSpace(fmt.Sprint(row[f.Name])))
	}
	return outputRow
}

func GetWriter(filename string) (*bufio.Writer, error) {
	file, err := os.Create(filename)

	if err != nil {
		return nil, err
	}

	return bufio.NewWriter(file), nil
}
