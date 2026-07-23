package gitlab

import (
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
	ID        int64  `json:"id"`
	ProjectID int64  `json:"project_id"`
	Name      string `json:"name"`
}

// CreateProjectAliasOptions represents the options for creating a project alias.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_aliases/#create-a-project-alias
type CreateProjectAliasOptions struct {
	Name      *string `json:"name" url:"name,omitempty"`
	ProjectID int64   `json:"project_id" url:"project_id,omitempty"`
}

// ListProjectAliases gets a list of all project aliases.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_aliases/#list-all-project-aliases
func (s *ProjectAliasesService) ListProjectAliases(options ...RequestOptionFunc) ([]*ProjectAlias, *Response, error) {
	return do[[]*ProjectAlias](s.client,
		withPath("project_aliases"),
		withRequestOpts(options...),
	)
}

// GetProjectAlias gets details of a project alias.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_aliases/#get-project-alias-details
func (s *ProjectAliasesService) GetProjectAlias(name string, options ...RequestOptionFunc) (*ProjectAlias, *Response, error) {
	return do[*ProjectAlias](s.client,
		withPath("project_aliases/%s", name),
		withRequestOpts(options...),
	)
}

// CreateProjectAlias creates a new project alias.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_aliases/#create-a-project-alias
func (s *ProjectAliasesService) CreateProjectAlias(opt *CreateProjectAliasOptions, options ...RequestOptionFunc) (*ProjectAlias, *Response, error) {
	return do[*ProjectAlias](s.client,
		withMethod(http.MethodPost),
		withPath("project_aliases"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteProjectAlias deletes a project alias.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_aliases/#delete-a-project-alias
func (s *ProjectAliasesService) DeleteProjectAlias(name string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("project_aliases/%s", name),
		withRequestOpts(options...),
	)
	return resp, err
}
