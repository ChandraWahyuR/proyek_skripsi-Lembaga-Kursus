package helper

import (
	"context"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"skripsi/config"
	"time"

	"cloud.google.com/go/storage"
)

const (
	UploadPathKategori = "gambar/kategori/"
	UploadPathKursus   = "gambar/kursus/"
	UploadPathUser     = "gambar/users/"
)

type ClientUploader struct {
	cl         *storage.Client
	projectID  string
	BucketName string
	UploadPath string
}

var Uploader *ClientUploader

func InitGCP() {
	config := config.InitConfig()

	os.Setenv("GOOGLE_APPLICATION_CREDENTIALS", config.GoogleCredentials)
	log.Println("GOOGLE_APPLICATION_CREDENTIALS:", os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"))

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	Uploader = &ClientUploader{
		cl:         client,
		BucketName: config.GCP_BucketName,
		projectID:  config.GCP_ProjectID,
	}
}

func (c *ClientUploader) UploadFile(file multipart.File, object string, uploadPath string) error {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	// Set the path dynamically
	fullPath := uploadPath + object

	// Upload an object with storage.Writer
	wc := c.cl.Bucket(c.BucketName).Object(fullPath).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return fmt.Errorf("io.Copy: %v", err)
	}
	if err := wc.Close(); err != nil {
		return fmt.Errorf("Writer.Close: %v", err)
	}

	return nil
}

func (c *ClientUploader) UploadFileGambarKategori(file multipart.File, object string) error {
	return c.UploadFile(file, object, UploadPathKategori)
}

func (c *ClientUploader) UploadFileGambarKursus(file multipart.File, object string) error {
	return c.UploadFile(file, object, UploadPathKursus)
}
