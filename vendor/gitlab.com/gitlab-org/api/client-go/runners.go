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
	"time"
)

type (
	RunnersServiceInterface interface {
		// ListRunners gets a list of runners accessible by the authenticated user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#list-owned-runners
		ListRunners(opt *ListRunnersOptions, options ...RequestOptionFunc) ([]*Runner, *Response, error)
		// ListAllRunners gets a list of all runners in the GitLab instance. Access is
		// restricted to users with admin privileges.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#list-all-runners
		ListAllRunners(opt *ListRunnersOptions, options ...RequestOptionFunc) ([]*Runner, *Response, error)
		// GetRunnerDetails returns details for given runner.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#get-runners-details
		GetRunnerDetails(rid any, options ...RequestOptionFunc) (*RunnerDetails, *Response, error)
		// UpdateRunnerDetails updates details for a given runner.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#update-runners-details
		UpdateRunnerDetails(rid any, opt *UpdateRunnerDetailsOptions, options ...RequestOptionFunc) (*RunnerDetails, *Response, error)
		// RemoveRunner removes a runner.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#delete-a-runner
		RemoveRunner(rid any, options ...RequestOptionFunc) (*Response, error)
		// ListRunnerJobs gets a list of jobs that are being processed or were processed by specified Runner.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#list-jobs-processed-by-a-runner
		ListRunnerJobs(rid any, opt *ListRunnerJobsOptions, options ...RequestOptionFunc) ([]*Job, *Response, error)
		// ListProjectRunners gets a list of runners accessible by the authenticated user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#list-projects-runners
		ListProjectRunners(pid any, opt *ListProjectRunnersOptions, options ...RequestOptionFunc) ([]*Runner, *Response, error)
		// EnableProjectRunner enables an available specific runner in the project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#assign-a-runner-to-project
		EnableProjectRunner(pid any, opt *EnableProjectRunnerOptions, options ...RequestOptionFunc) (*Runner, *Response, error)
		// DisableProjectRunner disables a specific runner from project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#unassign-a-runner-from-project
		DisableProjectRunner(pid any, runner int64, options ...RequestOptionFunc) (*Response, error)
		// ListGroupsRunners lists all runners (specific and shared) available in the
		// group as well it's ancestor groups. Shared runners are listed if at least one
		// shared runner is defined.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#list-groups-runners
		ListGroupsRunners(gid any, opt *ListGroupsRunnersOptions, options ...RequestOptionFunc) ([]*Runner, *Response, error)
		// RegisterNewRunner registers a new Runner for the instance.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#create-a-runner
		RegisterNewRunner(opt *RegisterNewRunnerOptions, options ...RequestOptionFunc) (*Runner, *Response, error)
		// DeleteRegisteredRunner deletes a Runner by Token.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#delete-a-runner-by-authentication-token
		DeleteRegisteredRunner(opt *DeleteRegisteredRunnerOptions, options ...RequestOptionFunc) (*Response, error)
		// DeleteRegisteredRunnerByID deletes a runner by ID.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#delete-a-runner-by-id
		DeleteRegisteredRunnerByID(rid int64, options ...RequestOptionFunc) (*Response, error)
		// VerifyRegisteredRunner registers a new runner for the instance.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#verify-authentication-for-a-registered-runner
		VerifyRegisteredRunner(opt *VerifyRegisteredRunnerOptions, options ...RequestOptionFunc) (*Response, error)
		// ResetRunnerAuthenticationToken resets a runner's authentication token.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runners/#reset-runners-authentication-token-by-using-the-runner-id
		ResetRunnerAuthenticationToken(rid int64, options ...RequestOptionFunc) (*RunnerAuthenticationToken, *Response, error)

		// Deprecated: for removal in GitLab 20.0, see https://docs.gitlab.com/ci/runners/new_creation_workflow/ instead
		ResetInstanceRunnerRegistrationToken(options ...RequestOptionFunc) (*RunnerRegistrationToken, *Response, error)

		// Deprecated: for removal in GitLab 20.0, see https://docs.gitlab.com/ci/runners/new_creation_workflow/ instead
		ResetGroupRunnerRegistrationToken(gid any, options ...RequestOptionFunc) (*RunnerRegistrationToken, *Response, error)

		// Deprecated: for removal in GitLab 20.0, see https://docs.gitlab.com/ci/runners/new_creation_workflow/ instead
		ResetProjectRunnerRegistrationToken(pid any, options ...RequestOptionFunc) (*RunnerRegistrationToken, *Response, error)
	}

	// RunnersService handles communication with the runner related methods of the
	// GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/runners/
	RunnersService struct {
		client *Client
	}
)

var _ RunnersServiceInterface = (*RunnersService)(nil)

// Runner represents a GitLab CI Runner.
//
// GitLab API docs: https://docs.gitlab.com/api/runners/
type Runner struct {
	ID             int64      `json:"id"`
	Description    string     `json:"description"`
	Paused         bool       `json:"paused"`
	IsShared       bool       `json:"is_shared"`
	RunnerType     string     `json:"runner_type"`
	Name           string     `json:"name"`
	Online         bool       `json:"online"`
	Status         string     `json:"status"`
	Token          string     `json:"token"`
	TokenExpiresAt *time.Time `json:"token_expires_at"`

	// Deprecated: for removal in v5 of the API, use Paused instead
	Active bool `json:"active"`

	// Deprecated: for removal in v5 of the API, returns an empty string from 17.0 onwards, see GraphQL resource CiRunnerManager instead
	IPAddress string `json:"ip_address"`
}

// RunnerDetails represents the GitLab CI runner details.
//
// GitLab API docs: https://docs.gitlab.com/api/runners/
type RunnerDetails struct {
	Paused          bool                   `json:"paused"`
	Description     string                 `json:"description"`
	ID              int64                  `json:"id"`
	IsShared        bool                   `json:"is_shared"`
	RunnerType      string                 `json:"runner_type"`
	ContactedAt     *time.Time             `json:"contacted_at"`
	MaintenanceNote string                 `json:"maintenance_note"`
	Name            string                 `json:"name"`
	Online          bool                   `json:"online"`
	Status          string                 `json:"status"`
	Projects        []RunnerDetailsProject `json:"projects"`
	Token           string                 `json:"token"`
	TagList         []string               `json:"tag_list"`
	RunUntagged     bool                   `json:"run_untagged"`
	Locked          bool                   `json:"locked"`
	AccessLevel     string                 `json:"access_level"`
	MaximumTimeout  int64                  `json:"maximum_timeout"`
	Groups          []RunnerDetailsGroup   `json:"groups"`

	// Deprecated: for removal in v5 of the API, see GraphQL resource CiRunnerManager instead
	Architecture string `json:"architecture"`

	// Deprecated: for removal in v5 of the API, returns an empty string from 17.0 onwards, see GraphQL resource CiRunnerManager instead
	IPAddress string `json:"ip_address"`

	// Deprecated: for removal in v5 of the API, see GraphQL resource CiRunnerManager instead
	Platform string `json:"platform"`

	// Deprecated: for removal in v5 of the API, see GraphQL resource CiRunnerManager instead
	Revision string `json:"revision"`

	// Deprecated: for removal in v5 of the API, see GraphQL resource CiRunnerManager instead
	Version string `json:"version"`

	// Deprecated: for removal in v5 of the API, use Paused instead
	Active bool `json:"active"`
}

// RunnerDetailsProject represents the GitLab CI runner details project.
//
// GitLab API docs: https://docs.gitlab.com/api/runners/
type RunnerDetailsProject struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	NameWithNamespace string `json:"name_with_namespace"`
	Path              string `json:"path"`
	PathWithNamespace string `json:"path_with_namespace"`
}

// RunnerDetailsGroup represents the GitLab CI runner details group.
//
// GitLab API docs: https://docs.gitlab.com/api/runners/
type RunnerDetailsGroup struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	WebURL string `json:"web_url"`
}

// ListRunnersOptions represents the available ListRunners() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#list-owned-runners
type ListRunnersOptions struct {
	ListOptions
	Type    *string   `url:"type,omitempty" json:"type,omitempty"`
	Status  *string   `url:"status,omitempty" json:"status,omitempty"`
	Paused  *bool     `url:"paused,omitempty" json:"paused,omitempty"`
	TagList *[]string `url:"tag_list,comma,omitempty" json:"tag_list,omitempty"`

	// Deprecated: Use Type or Status instead.
	Scope *string `url:"scope,omitempty" json:"scope,omitempty"`
}

func (s *RunnersService) ListRunners(opt *ListRunnersOptions, options ...RequestOptionFunc) ([]*Runner, *Response, error) {
	res, resp, err := do[[]*Runner](s.client,
		withMethod(http.MethodGet),
		withPath("runners"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

func (s *RunnersService) ListAllRunners(opt *ListRunnersOptions, options ...RequestOptionFunc) ([]*Runner, *Response, error) {
	res, resp, err := do[[]*Runner](s.client,
		withMethod(http.MethodGet),
		withPath("runners/all"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

func (s *RunnersService) GetRunnerDetails(rid any, options ...RequestOptionFunc) (*RunnerDetails, *Response, error) {
	res, resp, err := do[*RunnerDetails](s.client,
		withMethod(http.MethodGet),
		withPath("runners/%s", RunnerID{rid}),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

// UpdateRunnerDetailsOptions represents the available UpdateRunnerDetails() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#update-runners-details
type UpdateRunnerDetailsOptions struct {
	Description     *string   `url:"description,omitempty" json:"description,omitempty"`
	Paused          *bool     `url:"paused,omitempty" json:"paused,omitempty"`
	TagList         *[]string `url:"tag_list[],omitempty" json:"tag_list,omitempty"`
	RunUntagged     *bool     `url:"run_untagged,omitempty" json:"run_untagged,omitempty"`
	Locked          *bool     `url:"locked,omitempty" json:"locked,omitempty"`
	AccessLevel     *string   `url:"access_level,omitempty" json:"access_level,omitempty"`
	MaximumTimeout  *int64    `url:"maximum_timeout,omitempty" json:"maximum_timeout,omitempty"`
	MaintenanceNote *string   `url:"maintenance_note,omitempty" json:"maintenance_note,omitempty"`

	// Deprecated: for removal in v5 of the API, use Paused instead
	Active *bool `url:"active,omitempty" json:"active,omitempty"`
}

func (s *RunnersService) UpdateRunnerDetails(rid any, opt *UpdateRunnerDetailsOptions, options ...RequestOptionFunc) (*RunnerDetails, *Response, error) {
	res, resp, err := do[*RunnerDetails](s.client,
		withMethod(http.MethodPut),
		withPath("runners/%s", RunnerID{rid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

func (s *RunnersService) RemoveRunner(rid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("runners/%s", RunnerID{rid}),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListRunnerJobsOptions represents the available ListRunnerJobs()
// options. Status can be one of: running, success, failed, canceled.
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#list-jobs-processed-by-a-runner
type ListRunnerJobsOptions struct {
	ListOptions
	Status  *string `url:"status,omitempty" json:"status,omitempty"`
	OrderBy *string `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort    *string `url:"sort,omitempty" json:"sort,omitempty"`
}

func (s *RunnersService) ListRunnerJobs(rid any, opt *ListRunnerJobsOptions, options ...RequestOptionFunc) ([]*Job, *Response, error) {
	res, resp, err := do[[]*Job](s.client,
		withMethod(http.MethodGet),
		withPath("runners/%s/jobs", RunnerID{rid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

// ListProjectRunnersOptions represents the available ListProjectRunners()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#list-projects-runners
type ListProjectRunnersOptions ListRunnersOptions

func (s *RunnersService) ListProjectRunners(pid any, opt *ListProjectRunnersOptions, options ...RequestOptionFunc) ([]*Runner, *Response, error) {
	res, resp, err := do[[]*Runner](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/runners", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

// EnableProjectRunnerOptions represents the available EnableProjectRunner()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#assign-a-runner-to-project
type EnableProjectRunnerOptions struct {
	RunnerID int64 `json:"runner_id"`
}

func (s *RunnersService) EnableProjectRunner(pid any, opt *EnableProjectRunnerOptions, options ...RequestOptionFunc) (*Runner, *Response, error) {
	res, resp, err := do[*Runner](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/runners", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

func (s *RunnersService) DisableProjectRunner(pid any, runner int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/runners/%d", ProjectID{pid}, runner),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListGroupsRunnersOptions represents the available ListGroupsRunners() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#list-groups-runners
type ListGroupsRunnersOptions struct {
	ListOptions
	Type    *string   `url:"type,omitempty" json:"type,omitempty"`
	Status  *string   `url:"status,omitempty" json:"status,omitempty"`
	TagList *[]string `url:"tag_list,comma,omitempty" json:"tag_list,omitempty"`
}

func (s *RunnersService) ListGroupsRunners(gid any, opt *ListGroupsRunnersOptions, options ...RequestOptionFunc) ([]*Runner, *Response, error) {
	res, resp, err := do[[]*Runner](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/runners", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

// RegisterNewRunnerOptions represents the available RegisterNewRunner()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#create-a-runner
type RegisterNewRunnerOptions struct {
	Token           *string                       `url:"token" json:"token"`
	Description     *string                       `url:"description,omitempty" json:"description,omitempty"`
	Info            *RegisterNewRunnerInfoOptions `url:"info,omitempty" json:"info,omitempty"`
	Paused          *bool                         `url:"paused,omitempty" json:"paused,omitempty"`
	Locked          *bool                         `url:"locked,omitempty" json:"locked,omitempty"`
	RunUntagged     *bool                         `url:"run_untagged,omitempty" json:"run_untagged,omitempty"`
	TagList         *[]string                     `url:"tag_list[],omitempty" json:"tag_list,omitempty"`
	AccessLevel     *string                       `url:"access_level,omitempty" json:"access_level,omitempty"`
	MaximumTimeout  *int64                        `url:"maximum_timeout,omitempty" json:"maximum_timeout,omitempty"`
	MaintenanceNote *string                       `url:"maintenance_note,omitempty" json:"maintenance_note,omitempty"`

	// Deprecated: for removal in v5 of the API, use Paused instead
	Active *bool `url:"active,omitempty" json:"active,omitempty"`
}

// RegisterNewRunnerInfoOptions represents the info hashmap parameter in
// RegisterNewRunnerOptions.
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#create-a-runner
type RegisterNewRunnerInfoOptions struct {
	Name         *string `url:"name,omitempty" json:"name,omitempty"`
	Version      *string `url:"version,omitempty" json:"version,omitempty"`
	Revision     *string `url:"revision,omitempty" json:"revision,omitempty"`
	Platform     *string `url:"platform,omitempty" json:"platform,omitempty"`
	Architecture *string `url:"architecture,omitempty" json:"architecture,omitempty"`
}

func (s *RunnersService) RegisterNewRunner(opt *RegisterNewRunnerOptions, options ...RequestOptionFunc) (*Runner, *Response, error) {
	res, resp, err := do[*Runner](s.client,
		withMethod(http.MethodPost),
		withPath("runners"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

// DeleteRegisteredRunnerOptions represents the available
// DeleteRegisteredRunner() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#delete-a-runner-by-authentication-token
type DeleteRegisteredRunnerOptions struct {
	Token *string `url:"token" json:"token"`
}

func (s *RunnersService) DeleteRegisteredRunner(opt *DeleteRegisteredRunnerOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("runners"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *RunnersService) DeleteRegisteredRunnerByID(rid int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath(fmt.Sprintf("runners/%d", rid)),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}

// VerifyRegisteredRunnerOptions represents the available
// VerifyRegisteredRunner() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#verify-authentication-for-a-registered-runner
type VerifyRegisteredRunnerOptions struct {
	Token *string `url:"token" json:"token"`
}

func (s *RunnersService) VerifyRegisteredRunner(opt *VerifyRegisteredRunnerOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("runners/verify"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

type RunnerRegistrationToken struct {
	Token          *string    `url:"token" json:"token"`
	TokenExpiresAt *time.Time `url:"token_expires_at" json:"token_expires_at"`
}

// ResetInstanceRunnerRegistrationToken resets the instance runner registration
// token.
// Deprecated: for removal in GitLab 20.0, see https://docs.gitlab.com/ci/runners/new_creation_workflow/ instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#reset-instances-runner-registration-token
func (s *RunnersService) ResetInstanceRunnerRegistrationToken(options ...RequestOptionFunc) (*RunnerRegistrationToken, *Response, error) {
	res, resp, err := do[*RunnerRegistrationToken](s.client,
		withMethod(http.MethodPost),
		withPath("runners/reset_registration_token"),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

// ResetGroupRunnerRegistrationToken resets a group's runner registration token.
// Deprecated: for removal in GitLab 20.0, see https://docs.gitlab.com/ci/runners/new_creation_workflow/ instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#reset-groups-runner-registration-token
func (s *RunnersService) ResetGroupRunnerRegistrationToken(gid any, options ...RequestOptionFunc) (*RunnerRegistrationToken, *Response, error) {
	res, resp, err := do[*RunnerRegistrationToken](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/runners/reset_registration_token", GroupID{gid}),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

// ResetProjectRunnerRegistrationToken resets a projects's runner registration token.
// Deprecated: for removal in GitLab 20.0, see https://docs.gitlab.com/ci/runners/new_creation_workflow/ instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/runners/#reset-projects-runner-registration-token
func (s *RunnersService) ResetProjectRunnerRegistrationToken(pid any, options ...RequestOptionFunc) (*RunnerRegistrationToken, *Response, error) {
	res, resp, err := do[*RunnerRegistrationToken](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/runners/reset_registration_token", ProjectID{pid}),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}

type RunnerAuthenticationToken struct {
	Token          *string    `url:"token" json:"token"`
	TokenExpiresAt *time.Time `url:"token_expires_at" json:"token_expires_at"`
}

func (s *RunnersService) ResetRunnerAuthenticationToken(rid int64, options ...RequestOptionFunc) (*RunnerAuthenticationToken, *Response, error) {
	res, resp, err := do[*RunnerAuthenticationToken](s.client,
		withMethod(http.MethodPost),
		withPath("runners/%d/reset_authentication_token", rid),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return res, resp, nil
}
