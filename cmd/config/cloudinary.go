package config

import (
	"log"
	"sync"

	"github.com/cloudinary/cloudinary-go/v2"
)

var ClientCloudinary *cloudinary.Cloudinary
var once sync.Once

func InitCloudinary(configuration *Configuration) {
	once.Do(func() {
		var cloudinary, err = cloudinary.NewFromURL(configuration.Cloudinary.Url)
		if err != nil {
			log.Fatalf("Failed to intialize Cloudinary, %v", err)
		}

		ClientCloudinary = cloudinary
	})
}

func GetCloudinaryClient() *cloudinary.Cloudinary {
	return ClientCloudinary
}
