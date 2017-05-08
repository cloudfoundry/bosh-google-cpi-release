package project

type GoogleProjectService struct {
	ProjectId string
}

func (s GoogleProjectService) Find(projectId string) string {
	// Prefer caller provided projectId over default projectId
	if projectId != "" {
		return projectId
	}

	return s.ProjectId
}

func NewGoogleProjectService(projectId string) GoogleProjectService {
	return GoogleProjectService{ProjectId: projectId}
}
