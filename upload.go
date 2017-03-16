package bqutils

import (
	"context"
	"log"
	"runtime"

	"cloud.google.com/go/bigquery"
	"github.com/MartinSahlen/workerpool"
	uuid "github.com/satori/go.uuid"
)

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

func (u UploaderPool) AddRow(row UploadWrapper) {
	u.pool.Exec(uploadRowTask{Row: row, Uploader: u.uploader})
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
