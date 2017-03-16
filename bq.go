package bqutils

import (
	"context"
	"errors"
	"strings"

	"cloud.google.com/go/bigquery"
	bigqueryV2 "google.golang.org/api/bigquery/v2"
)

func GetNumRowsForJob(project, jobID string) (uint64, error) {

	client, err := v2Client()

	if err != nil {
		return 0, err
	}

	jobs := bigqueryV2.NewJobsService(client)

	queryResults, err := jobs.GetQueryResults(project, jobID).Do()

	if err != nil {
		return 0, err
	}

	return queryResults.TotalRows, nil
}

func GetNumRowsForTable(project, dataset, table string) (uint64, error) {

	client, err := v2Client()

	if err != nil {
		return 0, err
	}

	t, err := bigqueryV2.NewTablesService(client).Get(project, dataset, table).Do()

	if err != nil {
		return 0, err
	}

	return t.NumRows, nil
}

func CreateTable(project, dataset, table string, schema bigquery.Schema, force bool) (*bigquery.Table, error) {
	client, err := client(project)

	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = client.Dataset(dataset).Table(table).Create(ctx, schema, bigquery.UseStandardSQL())

	if err != nil {
		errString := err.Error()

		if strings.Contains(errString, "Error 400") {
			return nil, err
		}

		if force && strings.Contains(errString, "Error 409") {
			err = client.Dataset(dataset).Table(table).Delete(ctx)

			if err != nil {
				return nil, err
			}

			err = client.Dataset(dataset).Table(table).Create(ctx, schema, bigquery.UseStandardSQL())

			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return client.Dataset(dataset).Table(table), nil
}

func GetTableRows(project, dataset, table string) (*bigquery.RowIterator, error) {

	client, err := client(project)

	if err != nil {
		return nil, err
	}

	return client.Dataset(dataset).Table(table).Read(context.Background()), nil
}

func ParseTableName(tableName string) (*string, *string, error) {
	s := strings.Split(tableName, ".")
	if len(s) != 2 {
		return nil, nil, errors.New("Malformed table name: " + tableName)
	}
	return &s[0], &s[1], nil
}

func GetTableMeta(project, dataset, table string) (*bigquery.TableMetadata, error) {

	client, err := client(project)

	if err != nil {
		return nil, err
	}

	return client.Dataset(dataset).Table(table).Metadata(context.Background())
}
