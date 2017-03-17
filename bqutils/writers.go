package bqutils

import (
	"context"
	"errors"
	"io"
	"os"
	"strings"

	"cloud.google.com/go/storage"
)

var GCSFilePrefix = "gs://"

func GetWriter(filename, contentType string) (io.WriteCloser, error) {

	if strings.HasPrefix(filename, GCSFilePrefix) {

		if len(filename) == len(GCSFilePrefix) {
			return nil, errors.New("No GCS bucket/file specified")
		}

		bucketPath := filename[5:]

		pathFragments := strings.Split(bucketPath, "/")

		bucketName := pathFragments[0]

		if len(pathFragments) == 1 {
			return nil, errors.New("You have specified a bucket, but not a file or file path")
		}

		filePath := strings.Join(pathFragments[1:], "/")

		return StorageWriter(bucketName, filePath, contentType)

	}
	return FileWriter(filename)
}

func FileWriter(filename string) (io.WriteCloser, error) {
	file, err := os.Create(filename)

	if err != nil {
		return nil, err
	}
	return file, nil
}

func StorageWriter(bucket, filename, contentType string) (io.WriteCloser, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return nil, err
	}
	wc := client.Bucket(bucket).Object(filename).NewWriter(ctx)
	wc.ContentType = contentType

	return wc, nil
}
