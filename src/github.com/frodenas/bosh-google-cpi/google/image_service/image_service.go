package image

type Service interface {
	CreateFromURL(sourceURL string, description string) (string, error)
	CreateFromTarball(imagePath string, description string) (string, error)
	Delete(id string) error
	Find(id string) (Image, bool, error)
}
