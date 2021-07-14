package project_domain

import (
	"encoding/json"
	"fmt"

	"github.com/chronark/terraform-provider-vercel/pkg/vercel/httpApi"
)

type ProjectDomainHandler struct {
	Api httpApi.API
}

func (p *ProjectDomainHandler) Create(projectID string, projectDomain CreateProjectDomain, teamId string) (string, error) {
	url := fmt.Sprintf("/v8/projects/%s/domains", projectID)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	res, err := p.Api.Request("POST", url, projectDomain)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	var createdProjectDomain ProjectDomain
	err = json.NewDecoder(res.Body).Decode(&createdProjectDomain)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s:%s", projectID, createdProjectDomain.Name), nil
}

func (p *ProjectDomainHandler) Read(projectID, name, teamId string) (projectDomain ProjectDomain, err error) {
	url := fmt.Sprintf("/v8/projects/%s/domains/%s", projectID, name)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	res, err := p.Api.Request("GET", url, nil)
	if err != nil {
		return ProjectDomain{}, fmt.Errorf("unable to fetch project domain from vercel: %w", err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&projectDomain)
	if err != nil {
		return ProjectDomain{}, fmt.Errorf("unable to unmarshal project domain: %w", err)
	}

	return projectDomain, nil
}

func (p *ProjectDomainHandler) Update(projectID, name string, updateProjectDomain UpdateProjectDomain, teamId string) (projectDomain ProjectDomain, err error) {
	url := fmt.Sprintf("/v8/projects/%s/domains/%s", projectID, name)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}
	res, err := p.Api.Request("PATCH", url, updateProjectDomain)
	if err != nil {
		return ProjectDomain{}, fmt.Errorf("unable to update project domain: %w", err)
	}
	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(&projectDomain)
	if err != nil {
		return ProjectDomain{}, fmt.Errorf("unable to unmarshal project domain: %w", err)
	}

	return projectDomain, nil
}

func (p *ProjectDomainHandler) Delete(projectID, name, teamId string) error {
	url := fmt.Sprintf("/v8/projects/%s/domains/%s", projectID, name)
	if teamId != "" {
		url = fmt.Sprintf("%s/?teamId=%s", url, teamId)
	}

	res, err := p.Api.Request("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("unable to delete project domain: %w", err)
	}
	defer res.Body.Close()
	return nil
}
