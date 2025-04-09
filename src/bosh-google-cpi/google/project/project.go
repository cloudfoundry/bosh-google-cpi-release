package project

type Service interface {
	Find(projectId string) string
}
