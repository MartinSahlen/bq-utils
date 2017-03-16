package bqutils

import (
	"context"
	"log"
	"runtime"

	"google.golang.org/api/iterator"

	"cloud.google.com/go/bigquery"
	"github.com/MartinSahlen/workerpool"
	uuid "github.com/satori/go.uuid"
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

func mapRowsAndUpload(rows *bigquery.RowIterator, schema *bigquery.Schema, mapFunc MapRow, uploader UploaderPool) error {
	for {
		row, done, err := mapRows(rows, schema, mapFunc)
		if done {
			break
		} else if err != nil {
			return err
		}
		uploader.AddRow(row)
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

//UploadWrapper wraps a row for uploading through the ValueSaver interface
type UploadWrapper struct {
	Row map[string]bigquery.Value
}

//Save gives the bigquery uploader something to work with, including a
// UUID for insertID to avoid duplicates. could maybe use just an incrementer
func (u UploadWrapper) Save() (map[string]bigquery.Value, string, error) {
	return u.Row, uuid.NewV4().String(), nil
}

type UploaderPool struct {
	uploader *bigquery.Uploader
	pool     *workerpool.Pool
}

func NewUploaderPool(uploader *bigquery.Uploader, buffer uint64) UploaderPool {
	return UploaderPool{
		uploader: uploader,
		pool:     workerpool.NewPool(runtime.NumCPU()*2*8, buffer),
	}
}

func (u UploaderPool) AddRow(row map[string]bigquery.Value) {
	u.pool.Exec(uploadRowTask{Row: UploadWrapper{Row: row}, Uploader: u.uploader})
}

func (u UploaderPool) Wait() {
	u.pool.Close()
	u.pool.Wait()
}

type uploadRowTask struct {
	Row      UploadWrapper
	Uploader *bigquery.Uploader
}

func (u uploadRowTask) Execute() {
	ctx := context.Background()
	err := u.Uploader.Put(ctx, u.Row)
	if err != nil {
		e, ok := err.(bigquery.PutMultiError)
		if ok {
			for _, m := range e {
				log.Println(m.Error())
				for _, me := range m.Errors {
					log.Println(me.Error())
				}
			}
		}
		log.Println(err.Error())
	}
}

func UploadRows(table *bigquery.Table, schema *bigquery.Schema, rows *bigquery.RowIterator, mapFunc MapRow, buffer uint64) {
	uploader := NewUploaderPool(table.Uploader(), buffer)
	mapRowsAndUpload(rows, schema, mapFunc, uploader)
	uploader.Wait()
}
