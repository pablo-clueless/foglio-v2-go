package config

import (
	"github.com/cloudinary/cloudinary-go/v2"
)

func UseCloudinary() (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(AppConfig.CloudinaryName, AppConfig.CloudinaryKey, AppConfig.CloudinarySecret)
	if err != nil {
		return nil, err
	}

	return cld, err
}
