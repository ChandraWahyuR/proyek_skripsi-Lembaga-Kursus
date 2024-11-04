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
	UploadPathKategori   = "gambar/kategori/"
	UploadPathKursus     = "gambar/kursus/"
	UploadPathUser       = "gambar/users/"
	UploadPathUserKTP    = "gambar/users/ktp"
	UploadPathUserKK     = "gambar/users/kk"
	UploadPathIjazah     = "gambar/users/ijazah"
	UploadPathInstruktur = "gambar/instruktur/"
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

func (c *ClientUploader) UploadFileGambarUser(file multipart.File, object string) error {
	return c.UploadFile(file, object, UploadPathUser)
}

func (c *ClientUploader) UploadFileGambarKategori(file multipart.File, object string) error {
	return c.UploadFile(file, object, UploadPathKategori)
}

func (c *ClientUploader) UploadFileGambarKursus(file multipart.File, object string) error {
	return c.UploadFile(file, object, UploadPathKursus)
}

func (c *ClientUploader) UploadFileGambarInstruktur(file multipart.File, object string) error {
	return c.UploadFile(file, object, UploadPathInstruktur)
}

func (c *ClientUploader) UploadFileGambarKTP(file multipart.File, object string) error {
	return c.UploadFile(file, object, UploadPathUserKTP)
}

func (c *ClientUploader) UploadFileGambarKK(file multipart.File, object string) error {
	return c.UploadFile(file, object, UploadPathUserKK)
}
func (c *ClientUploader) UploadFileIjazah(file multipart.File, object string) error {
	return c.UploadFile(file, object, UploadPathUserKK)
}

// func (c *ClientUploader) DeleteFileGambarKursus(objectUrl string) error {
// 	// Parsing the object name from the full URL
// 	objectName := strings.Replace(objectUrl, fmt.Sprintf("https://storage.googleapis.com/%s/%s", c.BucketName, UploadPathKursus), "", 1)

// 	ctx := context.Background()
// 	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
// 	defer cancel()

// 	// Delete the object in GCS
// 	err := c.cl.Bucket(c.BucketName).Object(UploadPathKursus + objectName).Delete(ctx)
// 	if err != nil {
// 		return fmt.Errorf("Object(%q).Delete: %v", objectName, err)
// 	}
// 	return nil
// }
