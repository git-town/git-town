package gitlab

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	TerraformStatesServiceInterface interface {
		List(projectFullPath string, options ...RequestOptionFunc) ([]TerraformState, *Response, error)
		Get(projectFullPath string, name string, options ...RequestOptionFunc) (*TerraformState, *Response, error)
		Download(pid any, name string, serial uint64, options ...RequestOptionFunc) (io.Reader, *Response, error)
		DownloadLatest(pid any, name string, options ...RequestOptionFunc) (io.Reader, *Response, error)
		Delete(pid any, name string, options ...RequestOptionFunc) (*Response, error)
		DeleteVersion(pid any, name string, serial uint64, options ...RequestOptionFunc) (*Response, error)
		Lock(pid any, name string, options ...RequestOptionFunc) (*Response, error)
		Unlock(pid any, name string, options ...RequestOptionFunc) (*Response, error)
	}

	// TerraformStatesService handles communication with the GitLab-managed Terraform state API
	//
	// GitLab API docs: https://docs.gitlab.com/user/infrastructure/iac/terraform_state/
	TerraformStatesService struct {
		client *Client
	}
)

var _ TerraformStatesServiceInterface = (*TerraformStatesService)(nil)

// TerraformState represents a Terraform state.
//
// GitLab API docs: https://docs.gitlab.com/api/graphql/reference/#terraformstate
type TerraformState struct {
	Name          string                `json:"name"`
	LatestVersion TerraformStateVersion `json:"latestVersion"`
	CreatedAt     time.Time             `json:"createdAt"`
	UpdatedAt     time.Time             `json:"updatedAt"`
	DeletedAt     time.Time             `json:"deletedAt"`
	LockedAt      time.Time             `json:"lockedAt"`
}

// TerraformStateVersion represents a Terraform state version.
//
// GitLab API docs: https://docs.gitlab.com/api/graphql/reference/#terraformstateversion
type TerraformStateVersion struct {
	Serial       uint64    `json:"serial"`
	DownloadPath string    `json:"downloadPath"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

// List returns all Terraform states
func (s *TerraformStatesService) List(projectFullPath string, options ...RequestOptionFunc) ([]TerraformState, *Response, error) {
	query := GraphQLQuery{
		Query: fmt.Sprintf(`
			query {
				project(fullPath: %q) {
					terraformStates {
						nodes {
							name
							createdAt
							deletedAt
							latestVersion {
								createdAt
								updatedAt
								downloadPath
								serial
							}
							updatedAt
							lockedAt
						}
					}
				}
			}
		`, projectFullPath),
	}

	var response struct {
		Data struct {
			Project *struct {
				TerraformStates struct {
					Nodes []TerraformState `json:"nodes"`
				} `json:"terraformStates"`
			} `json:"project"`
		} `json:"data"`
	}
	resp, err := s.client.GraphQL.Do(query, &response, options...)
	if err != nil {
		return nil, resp, err
	}
	if response.Data.Project == nil {
		return nil, resp, ErrNotFound
	}

	return response.Data.Project.TerraformStates.Nodes, resp, nil
}

// Get returns a single Terraform state
func (s *TerraformStatesService) Get(projectFullPath string, name string, options ...RequestOptionFunc) (*TerraformState, *Response, error) {
	query := GraphQLQuery{
		Query: fmt.Sprintf(`
			query {
				project(fullPath: %q) {
					terraformState(name: %q) {
						name
						createdAt
						deletedAt
						latestVersion {
							createdAt
							updatedAt
							downloadPath
							serial
						}
						updatedAt
						lockedAt
					}
				}
			}
		`, projectFullPath, name),
	}

	var response struct {
		Data struct {
			Project *struct {
				TerraformState *TerraformState `json:"terraformState"`
			} `json:"project"`
		} `json:"data"`
	}
	resp, err := s.client.GraphQL.Do(query, &response, options...)
	if err != nil {
		return nil, resp, err
	}
	if response.Data.Project == nil || response.Data.Project.TerraformState == nil {
		return nil, resp, ErrNotFound
	}

	return response.Data.Project.TerraformState, resp, nil
}

func (s *TerraformStatesService) DownloadLatest(pid any, name string, options ...RequestOptionFunc) (io.Reader, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	uri := fmt.Sprintf("projects/%s/terraform/state/%s", PathEscape(project), PathEscape(name))

	req, err := s.client.NewRequest(http.MethodGet, uri, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var b bytes.Buffer
	resp, err := s.client.Do(req, &b)
	if err != nil {
		return nil, resp, err
	}

	return &b, resp, nil
}

func (s *TerraformStatesService) Download(pid any, name string, serial uint64, options ...RequestOptionFunc) (io.Reader, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	uri := fmt.Sprintf("projects/%s/terraform/state/%s/versions/%d", PathEscape(project), PathEscape(name), serial)

	req, err := s.client.NewRequest(http.MethodGet, uri, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var b bytes.Buffer
	resp, err := s.client.Do(req, &b)
	if err != nil {
		return nil, resp, err
	}

	return &b, resp, nil
}

// Delete deletes a single Terraform state
//
// GitLab API docs: https://docs.gitlab.com/user/infrastructure/iac/terraform_state/
func (s *TerraformStatesService) Delete(pid any, name string, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("projects/%s/terraform/state/%s", PathEscape(project), PathEscape(name))

	req, err := s.client.NewRequest(http.MethodDelete, uri, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// DeleteVersion deletes a single Terraform state version
//
// GitLab API docs: https://docs.gitlab.com/user/infrastructure/iac/terraform_state/
func (s *TerraformStatesService) DeleteVersion(pid any, name string, serial uint64, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("projects/%s/terraform/state/%s/versions/%d", PathEscape(project), PathEscape(name), serial)

	req, err := s.client.NewRequest(http.MethodDelete, uri, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Lock locks a single Terraform state
//
// GitLab API docs: https://docs.gitlab.com/user/infrastructure/iac/terraform_state/
func (s *TerraformStatesService) Lock(pid any, name string, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("projects/%s/terraform/state/%s/lock", PathEscape(project), PathEscape(name))

	req, err := s.client.NewRequest(http.MethodPost, uri, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// Unlock unlocks a single Terraform state
//
// GitLab API docs: https://docs.gitlab.com/user/infrastructure/iac/terraform_state/
func (s *TerraformStatesService) Unlock(pid any, name string, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	uri := fmt.Sprintf("projects/%s/terraform/state/%s/lock", PathEscape(project), PathEscape(name))

	req, err := s.client.NewRequest(http.MethodDelete, uri, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
