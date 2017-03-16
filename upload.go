package bqutils

import (
	"cloud.google.com/go/bigquery"
	uuid "github.com/satori/go.uuid"
)

//UploadWrapper wraps a row for uploading through the ValueSaver interface
type UploadWrapper struct {
	Row map[string]bigquery.Value
}

//Save gives the bigquery uploader something to work with, including a
// UUID for insertID to avoid duplicates. could maybe use just an incrementer
func (g UploadWrapper) Save() (map[string]bigquery.Value, string, error) {
	return g.Row, uuid.NewV4().String(), nil
}
