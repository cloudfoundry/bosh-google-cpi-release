package gimage

import (
	"google.golang.org/api/compute/v1"
)

type ImageService interface {
	CreateFromURL(sourceURL string, description string) (string, error)
	CreateFromTarball(imagePath string, description string) (string, error)
	Delete(id string) error
	Find(id string) (*compute.Image, bool, error)
}
