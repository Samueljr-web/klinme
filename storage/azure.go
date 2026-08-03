package storage

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
)

var Client *azblob.Client

func Connect() {
	connStr := os.Getenv("AZURE_STORAGE_CONNECTION_STRING")

	client, err := azblob.NewClientFromConnectionString(connStr, nil)
	if err != nil {
		panic("Failed to connect to Azure Blob: " + err.Error())
	}

	Client = client
	fmt.Println("Connected to Azure Blob Storage!")
}

func UploadFile(
	ctx context.Context,
	containerName string,
	fileName string,
	file io.Reader,
) (string, error) {
	_, err := Client.UploadStream(
		ctx,
		containerName,
		fileName,
		file,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Return the file URL
	accountName := os.Getenv("AZURE_STORAGE_ACCOUNT")
	url := fmt.Sprintf(
		"https://%s.blob.core.windows.net/%s/%s",
		accountName,
		containerName,
		fileName,
	)

	return url, nil
}
