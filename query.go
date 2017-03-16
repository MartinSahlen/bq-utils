package bqutils

import (
	"context"
	"errors"

	"cloud.google.com/go/bigquery"
	"google.golang.org/api/iterator"

	bigqueryV2 "google.golang.org/api/bigquery/v2"
)

type MapRow func(row map[string]bigquery.Value, schema *bigquery.Schema) error

func MapRows(rows *bigquery.RowIterator, schema *bigquery.Schema, mapFunc MapRow) error {
	for {
		row := map[string]bigquery.Value{}
		err := rows.Next(&row)
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		err = mapFunc(row, schema)
		if err != nil {
			return err
		}
	}
	return nil
}

func startJob(project, query string) (*bigquery.Job, error) {
	c, err := client(project)

	if err != nil {
		return nil, err
	}

	q := c.Query(query)

	q.QueryConfig = bigquery.QueryConfig{
		Q:                 query,
		UseStandardSQL:    true,
		Priority:          bigquery.InteractivePriority,
		DisableQueryCache: false,
		AllowLargeResults: false,
	}
	return q.Run(context.Background())
}

func GetQuerySchema(project, jobID string) (bigquery.Schema, error) {

	client, err := v2Client()

	if err != nil {
		return nil, err
	}

	jobs := bigqueryV2.NewJobsService(client)

	queryResults, err := jobs.GetQueryResults(project, jobID).Do()

	if err != nil {
		return nil, err
	}

	schema := bigquery.Schema{}

	for _, queryField := range queryResults.Schema.Fields {
		schema = append(schema, parseV2Field(queryField))
	}

	return schema, nil
}

func parseV2Field(v2Field *bigqueryV2.TableFieldSchema) *bigquery.FieldSchema {
	field := &bigquery.FieldSchema{
		Name:        v2Field.Name,
		Description: v2Field.Description,
		Repeated:    isRepeated(v2Field),
		Required:    isRequired(v2Field),
		Type:        fieldType(v2Field),
		Schema:      bigquery.Schema{},
	}

	if field.Type == bigquery.RecordFieldType {
		for _, nestedField := range v2Field.Fields {
			field.Schema = append(field.Schema, parseV2Field(nestedField))
		}
	}

	return field
}

func fieldType(v2Field *bigqueryV2.TableFieldSchema) bigquery.FieldType {
	if v2Field.Type == string(bigquery.StringFieldType) {
		return bigquery.StringFieldType
	} else if v2Field.Type == string(bigquery.BytesFieldType) {
		return bigquery.BytesFieldType
	} else if v2Field.Type == string(bigquery.IntegerFieldType) {
		return bigquery.IntegerFieldType
	} else if v2Field.Type == string(bigquery.FloatFieldType) {
		return bigquery.FloatFieldType
	} else if v2Field.Type == string(bigquery.BooleanFieldType) {
		return bigquery.BooleanFieldType
	} else if v2Field.Type == string(bigquery.TimestampFieldType) {
		return bigquery.TimestampFieldType
	} else if v2Field.Type == string(bigquery.RecordFieldType) {
		return bigquery.RecordFieldType
	} else if v2Field.Type == string(bigquery.DateFieldType) {
		return bigquery.DateFieldType
	} else if v2Field.Type == string(bigquery.TimeFieldType) {
		return bigquery.TimeFieldType
	} else if v2Field.Type == string(bigquery.DateTimeFieldType) {
		return bigquery.DateTimeFieldType
	}
	panic("Unknown fieldtype")
}

func isRepeated(v2Field *bigqueryV2.TableFieldSchema) bool {
	return v2Field.Mode == "REPEATED"
}

func isRequired(v2Field *bigqueryV2.TableFieldSchema) bool {
	return v2Field.Mode == "REQUIRED"
}

func RunQuery(project, query string) (*bigquery.RowIterator, *bigquery.Job, error) {
	job, err := startJob(project, query)

	if err != nil {
		return nil, nil, err
	}

	ctx := context.Background()

	status, err := job.Wait(ctx)

	if err != nil {
		return nil, nil, err
	}

	if status.State != bigquery.Done {
		return nil, nil, errors.New("Query job " + job.ID() + " had errors")
	}
	rows, err := job.Read(ctx)
	return rows, job, err
}
