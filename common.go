package bqutils

import (
	"context"

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
