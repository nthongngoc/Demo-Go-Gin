package cloudstorage

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"time"

	cloudStorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/khoa5773/go-server/src/configs"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

var storageBucket *cloudStorage.BucketHandle

func init() {
	ctx := context.Background()
	config := &firebase.Config{
		StorageBucket: fmt.Sprintf(configs.ConfigsService.GGCloudStorageBucket),
	}
	opt := option.WithCredentialsFile(configs.ConfigsService.GGServiceAccountPath)
	app, err := firebase.NewApp(ctx, config, opt)
	if err != nil {
		log.Fatalln(err)
	}

	client, err := app.Storage(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	storageBucket, err = client.DefaultBucket()
	if err != nil {
		log.Fatalln(err)
	}
}

func UploadFileToCloudStorage(
	ctx context.Context, file *multipart.FileHeader, dataID string, projectID string, mimeType string,
) (string, error) {
	fileName := fmt.Sprintf("%s-%s.%s", dataID, projectID, mimeType)
	object := storageBucket.Object(fileName)
	w := object.NewWriter(ctx)

	content, err := file.Open()
	if err != nil {
		return "", err
	}

	defer content.Close()
	byteData, err := ioutil.ReadAll(content)
	if err != nil {
		return "", err
	}

	_, err = w.Write(byteData)

	if err != nil {
		return "", err
	}

	if err := w.Close(); err != nil {
		return "", err
	}

	return fmt.Sprintf("http://storage.googleapis.com/%s/%s", configs.ConfigsService.GGCloudStorageBucket, fileName), nil
}

func GenerateSignedUrl(dataID string, projectID string, mimeType string) (string, error) {
	fileName := fmt.Sprintf("%s-%s.%s", dataID, projectID, mimeType)

	saKey, err := ioutil.ReadFile(configs.ConfigsService.GGServiceAccountPath)
	if err != nil {
		log.Fatalln(err)
	}

	cfg, err := google.JWTConfigFromJSON(saKey)
	if err != nil {
		log.Fatalln(err)
	}
	url, err := cloudStorage.SignedURL(configs.ConfigsService.GGCloudStorageBucket, fileName, &cloudStorage.SignedURLOptions{
		GoogleAccessID: cfg.Email,
		PrivateKey:     cfg.PrivateKey,
		Method:         "GET",
		Expires:        time.Now().Add(3 * time.Hour),
	})

	if err != nil {
		return "", err
	}

	return url, nil
}

func RemoveFileFromCloudStorage(
	ctx context.Context, dataID string, projectID string, mimeType string,
) (bool, error) {
	object := storageBucket.Object(fmt.Sprintf("%s-%s.%s", dataID, projectID, mimeType))
	if err := object.Delete(ctx); err != nil {
		return false, err
	}

	return true, nil
}
