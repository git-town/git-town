package gitlab

import (
	"fmt"
	"net/http"
)

type (
	ProjectAliasesServiceInterface interface {
		ListProjectAliases(options ...RequestOptionFunc) ([]*ProjectAlias, *Response, error)
		GetProjectAlias(name string, options ...RequestOptionFunc) (*ProjectAlias, *Response, error)
		CreateProjectAlias(opt *CreateProjectAliasOptions, options ...RequestOptionFunc) (*ProjectAlias, *Response, error)
		DeleteProjectAlias(name string, options ...RequestOptionFunc) (*Response, error)
	}

	// ProjectAliasesService handles communication with the project aliases related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/project_aliases/
	ProjectAliasesService struct {
		client *Client
	}
)

var _ ProjectAliasesServiceInterface = (*ProjectAliasesService)(nil)

// ProjectAlias represents a GitLab project alias.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_aliases/
type ProjectAlias struct {
	ID        int    `json:"id"`
	ProjectID int    `json:"project_id"`
	Name      string `json:"name"`
}

// CreateProjectAliasOptions represents the options for creating a project alias.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_aliases/#create-a-project-alias
type CreateProjectAliasOptions struct {
	Name      *string `json:"name" url:"name,omitempty"`
	ProjectID int     `json:"project_id" url:"project_id,omitempty"`
}

// ListProjectAliases gets a list of all project aliases.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_aliases/#list-all-project-aliases
func (s *ProjectAliasesService) ListProjectAliases(options ...RequestOptionFunc) ([]*ProjectAlias, *Response, error) {
	u := "project_aliases"

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var aliases []*ProjectAlias
	resp, err := s.client.Do(req, &aliases)
	if err != nil {
		return nil, resp, err
	}

	return aliases, resp, nil
}

// GetProjectAlias gets details of a project alias.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_aliases/#get-project-alias-details
func (s *ProjectAliasesService) GetProjectAlias(name string, options ...RequestOptionFunc) (*ProjectAlias, *Response, error) {
	u := fmt.Sprintf("project_aliases/%s", PathEscape(name))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	aliasObj := new(ProjectAlias)
	resp, err := s.client.Do(req, aliasObj)
	if err != nil {
		return nil, resp, err
	}

	return aliasObj, resp, nil
}

// CreateProjectAlias creates a new project alias.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_aliases/#create-a-project-alias
func (s *ProjectAliasesService) CreateProjectAlias(opt *CreateProjectAliasOptions, options ...RequestOptionFunc) (*ProjectAlias, *Response, error) {
	u := "project_aliases"

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	alias := new(ProjectAlias)
	resp, err := s.client.Do(req, alias)
	if err != nil {
		return nil, resp, err
	}

	return alias, resp, nil
}

// DeleteProjectAlias deletes a project alias.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_aliases/#delete-a-project-alias
func (s *ProjectAliasesService) DeleteProjectAlias(name string, options ...RequestOptionFunc) (*Response, error) {
	u := fmt.Sprintf("project_aliases/%s", PathEscape(name))

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
