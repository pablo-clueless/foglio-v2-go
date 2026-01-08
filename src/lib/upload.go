package lib

import (
	"context"
	"foglio/v2/src/config"
	"log"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

func UploadMultiple(files []*multipart.FileHeader, path string) ([]string, error) {
	ctx := context.Background()
	cld, err := config.UseCloudinary()
	if err != nil {
		return nil, err
	}

	params := uploader.UploadParams{}
	if path != "" {
		params.Folder = path
	}

	urls := make([]string, 0, len(files))

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			return nil, err
		}
		defer func() {
			if err := file.Close(); err != nil {
				log.Printf("Error closing file: %v", err)
			}
		}()

		res, err := cld.Upload.Upload(ctx, file, params)
		if err != nil {
			return nil, err
		}

		urls = append(urls, res.SecureURL)
	}

	return urls, nil
}

func UploadSingle(fileHeader *multipart.FileHeader, path string) (string, error) {
	ctx := context.Background()
	cld, err := config.UseCloudinary()
	if err != nil {
		return "", err
	}

	params := uploader.UploadParams{}
	if path != "" {
		params.Folder = path
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Printf("Error closing file: %v", err)
		}
	}()

	res, err := cld.Upload.Upload(ctx, file, params)
	if err != nil {
		return "", err
	}

	return res.SecureURL, nil
}
