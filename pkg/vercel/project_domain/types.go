package project_domain

type CreateProjectDomain struct {
	Name               string `json:"name"`
	Redirect           string `json:"redirect,omitempty"`
	RedirectStatusCode int    `json:"redirectStatusCode,omitempty"`
	GitBranch          string `json:"gitBranch,omitempty"`
}

type UpdateProjectDomain struct {
	Redirect string `json:"redirect"`
}

type ProjectDomain struct {
	Name      string `json:"name"`
	ProjectID string `json:"projectId"`
	CreatedAt int    `json:"createdAt"`
	UpdatedAt int    `json:"updatedAt"`
	Redirect  string `json:"redirect"`
}
