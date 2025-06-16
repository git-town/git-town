//
// Copyright 2021, Sander van Harmelen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package gitlab

import (
	"fmt"
	"net/http"
)

type (
	ValidateServiceInterface interface {
		ProjectNamespaceLint(pid any, opt *ProjectNamespaceLintOptions, options ...RequestOptionFunc) (*ProjectLintResult, *Response, error)
		ProjectLint(pid any, opt *ProjectLintOptions, options ...RequestOptionFunc) (*ProjectLintResult, *Response, error)
	}

	// ValidateService handles communication with the validation related methods of
	// the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/lint/
	ValidateService struct {
		client *Client
	}
)

var _ ValidateServiceInterface = (*ValidateService)(nil)

// LintResult represents the linting results.
//
// GitLab API docs: https://docs.gitlab.com/api/lint/
type LintResult struct {
	Status     string   `json:"status"`
	Errors     []string `json:"errors"`
	Warnings   []string `json:"warnings"`
	MergedYaml string   `json:"merged_yaml"`
}

// ProjectLintResult represents the linting results by project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/lint/
type ProjectLintResult struct {
	Valid      bool      `json:"valid"`
	Errors     []string  `json:"errors"`
	Warnings   []string  `json:"warnings"`
	MergedYaml string    `json:"merged_yaml"`
	Includes   []Include `json:"includes"`
}

// Include contains the details about an include block in the .gitlab-ci.yml file.
// It is used in ProjectLintResult.
//
// Reference can be found at the lint API endpoint in the openapi yaml:
// https://gitlab.com/gitlab-org/gitlab/-/blob/master/doc/api/openapi/openapi_v2.yaml
type Include struct {
	Type           string         `json:"type"`
	Location       string         `json:"location"`
	Blob           string         `json:"blob"`
	Raw            string         `json:"raw"`
	Extra          map[string]any `json:"extra"`
	ContextProject string         `json:"context_project"`
	ContextSHA     string         `json:"context_sha"`
}

// ProjectNamespaceLintOptions represents the available ProjectNamespaceLint() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/lint/#validate-sample-cicd-configuration
type ProjectNamespaceLintOptions struct {
	Content     *string `url:"content,omitempty" json:"content,omitempty"`
	DryRun      *bool   `url:"dry_run,omitempty" json:"dry_run,omitempty"`
	IncludeJobs *bool   `url:"include_jobs,omitempty" json:"include_jobs,omitempty"`
	Ref         *string `url:"ref,omitempty" json:"ref,omitempty"`
}

// ProjectNamespaceLint validates .gitlab-ci.yml content by project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/lint/#validate-sample-cicd-configuration
func (s *ValidateService) ProjectNamespaceLint(pid any, opt *ProjectNamespaceLintOptions, options ...RequestOptionFunc) (*ProjectLintResult, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/ci/lint", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, &opt, options)
	if err != nil {
		return nil, nil, err
	}

	l := new(ProjectLintResult)
	resp, err := s.client.Do(req, l)
	if err != nil {
		return nil, resp, err
	}

	return l, resp, nil
}

// ProjectLintOptions represents the available ProjectLint() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/lint/#validate-a-projects-cicd-configuration
type ProjectLintOptions struct {
	ContentRef  *string `url:"content_ref,omitempty" json:"content_ref,omitempty"`
	DryRunRef   *string `url:"dry_run_ref,omitempty" json:"dry_run_ref,omitempty"`
	DryRun      *bool   `url:"dry_run,omitempty" json:"dry_run,omitempty"`
	IncludeJobs *bool   `url:"include_jobs,omitempty" json:"include_jobs,omitempty"`
	Ref         *string `url:"ref,omitempty" json:"ref,omitempty"`
}

// ProjectLint validates .gitlab-ci.yml content by project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/lint/#validate-a-projects-cicd-configuration
func (s *ValidateService) ProjectLint(pid any, opt *ProjectLintOptions, options ...RequestOptionFunc) (*ProjectLintResult, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/ci/lint", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, &opt, options)
	if err != nil {
		return nil, nil, err
	}

	l := new(ProjectLintResult)
	resp, err := s.client.Do(req, l)
	if err != nil {
		return nil, resp, err
	}

	return l, resp, nil
}
