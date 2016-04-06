package fakes

import (
	"bosh-google-cpi/google/image_service"
)

type FakeImageService struct {
	CreateFromURLCalled      bool
	CreateFromURLErr         error
	CreateFromURLID          string
	CreateFromURLSourceURL   string
	CreateFromURLDescription string

	CreateFromTarballCalled      bool
	CreateFromTarballErr         error
	CreateFromTarballID          string
	CreateFromTarballImagePath   string
	CreateFromTarballDescription string

	DeleteCalled bool
	DeleteErr    error

	FindCalled bool
	FindFound  bool
	FindImage  image.Image
	FindErr    error
}

func (i *FakeImageService) CreateFromURL(sourceURL string, description string) (string, error) {
	i.CreateFromURLCalled = true
	i.CreateFromURLSourceURL = sourceURL
	i.CreateFromURLDescription = description
	return i.CreateFromURLID, i.CreateFromURLErr
}

func (i *FakeImageService) CreateFromTarball(imagePath string, description string) (string, error) {
	i.CreateFromTarballCalled = true
	i.CreateFromTarballImagePath = imagePath
	i.CreateFromTarballDescription = description
	return i.CreateFromTarballID, i.CreateFromTarballErr
}

func (i *FakeImageService) Delete(id string) error {
	i.DeleteCalled = true
	return i.DeleteErr
}

func (i *FakeImageService) Find(id string) (image.Image, bool, error) {
	i.FindCalled = true
	return i.FindImage, i.FindFound, i.FindErr
}
