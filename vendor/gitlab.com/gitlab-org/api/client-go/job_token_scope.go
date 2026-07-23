// Copyright 2021, Sander van Harmelen
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitlab

import "net/http"

type (
	JobTokenScopeServiceInterface interface {
		GetProjectJobTokenAccessSettings(pid any, options ...RequestOptionFunc) (*JobTokenAccessSettings, *Response, error)
		PatchProjectJobTokenAccessSettings(pid any, opt *PatchProjectJobTokenAccessSettingsOptions, options ...RequestOptionFunc) (*Response, error)
		GetProjectJobTokenInboundAllowList(pid any, opt *GetJobTokenInboundAllowListOptions, options ...RequestOptionFunc) ([]*Project, *Response, error)
		AddProjectToJobScopeAllowList(pid any, opt *JobTokenInboundAllowOptions, options ...RequestOptionFunc) (*JobTokenInboundAllowItem, *Response, error)
		RemoveProjectFromJobScopeAllowList(pid any, targetProject int64, options ...RequestOptionFunc) (*Response, error)
		GetJobTokenAllowlistGroups(pid any, opt *GetJobTokenAllowlistGroupsOptions, options ...RequestOptionFunc) ([]*Group, *Response, error)
		AddGroupToJobTokenAllowlist(pid any, opt *AddGroupToJobTokenAllowlistOptions, options ...RequestOptionFunc) (*JobTokenAllowlistItem, *Response, error)
		RemoveGroupFromJobTokenAllowlist(pid any, targetGroup int64, options ...RequestOptionFunc) (*Response, error)
	}

	// JobTokenScopeService handles communication with project CI settings
	// such as token permissions.
	//
	// GitLab API docs: https://docs.gitlab.com/api/project_job_token_scopes/
	JobTokenScopeService struct {
		client *Client
	}
)

var _ JobTokenScopeServiceInterface = (*JobTokenScopeService)(nil)

// JobTokenAccessSettings represents job token access attributes for this project.
//
// GitLab API docs: https://docs.gitlab.com/api/project_job_token_scopes/
type JobTokenAccessSettings struct {
	InboundEnabled bool `json:"inbound_enabled"`
}

// GetProjectJobTokenAccessSettings fetch the CI/CD job token access settings (job token scope) of a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#get-a-projects-cicd-job-token-access-settings
func (j *JobTokenScopeService) GetProjectJobTokenAccessSettings(pid any, options ...RequestOptionFunc) (*JobTokenAccessSettings, *Response, error) {
	return do[*JobTokenAccessSettings](j.client,
		withPath("projects/%s/job_token_scope", ProjectID{pid}),
		withRequestOpts(options...),
	)
}

// PatchProjectJobTokenAccessSettingsOptions represents the available
// PatchProjectJobTokenAccessSettings() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#patch-a-projects-cicd-job-token-access-settings
type PatchProjectJobTokenAccessSettingsOptions struct {
	Enabled bool `json:"enabled"`
}

// PatchProjectJobTokenAccessSettings patch the Limit access to this project setting (job token scope) of a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#patch-a-projects-cicd-job-token-access-settings
func (j *JobTokenScopeService) PatchProjectJobTokenAccessSettings(pid any, opt *PatchProjectJobTokenAccessSettingsOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](j.client,
		withMethod(http.MethodPatch),
		withPath("projects/%s/job_token_scope", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

// JobTokenInboundAllowItem represents a single job token inbound allowlist item.
//
// GitLab API docs: https://docs.gitlab.com/api/project_job_token_scopes/
type JobTokenInboundAllowItem struct {
	SourceProjectID int64 `json:"source_project_id"`
	TargetProjectID int64 `json:"target_project_id"`
}

// GetJobTokenInboundAllowListOptions represents the available
// GetJobTokenInboundAllowList() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#get-a-projects-cicd-job-token-inbound-allowlist
type GetJobTokenInboundAllowListOptions struct {
	ListOptions
}

// GetProjectJobTokenInboundAllowList fetches the CI/CD job token inbound
// allowlist (job token scope) of a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#get-a-projects-cicd-job-token-inbound-allowlist
func (j *JobTokenScopeService) GetProjectJobTokenInboundAllowList(pid any, opt *GetJobTokenInboundAllowListOptions, options ...RequestOptionFunc) ([]*Project, *Response, error) {
	return do[[]*Project](j.client,
		withPath("projects/%s/job_token_scope/allowlist", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// JobTokenInboundAllowOptions represents the available
// AddProjectToJobScopeAllowList() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#add-a-project-to-a-cicd-job-token-inbound-allowlist
type JobTokenInboundAllowOptions struct {
	TargetProjectID *int64 `url:"target_project_id,omitempty" json:"target_project_id,omitempty"`
}

// AddProjectToJobScopeAllowList adds a new project to a project's job token
// inbound allow list.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#add-a-project-to-a-cicd-job-token-inbound-allowlist
func (j *JobTokenScopeService) AddProjectToJobScopeAllowList(pid any, opt *JobTokenInboundAllowOptions, options ...RequestOptionFunc) (*JobTokenInboundAllowItem, *Response, error) {
	return do[*JobTokenInboundAllowItem](j.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/job_token_scope/allowlist", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RemoveProjectFromJobScopeAllowList removes a project from a project's job
// token inbound allow list.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#remove-a-project-from-a-cicd-job-token-inbound-allowlist
func (j *JobTokenScopeService) RemoveProjectFromJobScopeAllowList(pid any, targetProject int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](j.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/job_token_scope/allowlist/%d", ProjectID{pid}, targetProject),
		withRequestOpts(options...),
	)
	return resp, err
}

// JobTokenAllowlistItem represents a single job token allowlist item.
//
// GitLab API docs: https://docs.gitlab.com/api/project_job_token_scopes/
type JobTokenAllowlistItem struct {
	SourceProjectID int64 `json:"source_project_id"`
	TargetGroupID   int64 `json:"target_group_id"`
}

// GetJobTokenAllowlistGroupsOptions represents the available
// GetJobTokenAllowlistGroups() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#get-a-projects-cicd-job-token-allowlist-of-groups
type GetJobTokenAllowlistGroupsOptions struct {
	ListOptions
}

// GetJobTokenAllowlistGroups fetches the CI/CD job token allowlist groups
// (job token scopes) of a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#get-a-projects-cicd-job-token-allowlist-of-groups
func (j *JobTokenScopeService) GetJobTokenAllowlistGroups(pid any, opt *GetJobTokenAllowlistGroupsOptions, options ...RequestOptionFunc) ([]*Group, *Response, error) {
	return do[[]*Group](j.client,
		withPath("projects/%s/job_token_scope/groups_allowlist", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// AddGroupToJobTokenAllowlistOptions represents the available
// AddGroupToJobTokenAllowlist() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#add-a-group-to-a-cicd-job-token-allowlist
type AddGroupToJobTokenAllowlistOptions struct {
	TargetGroupID *int64 `url:"target_group_id,omitempty" json:"target_group_id,omitempty"`
}

// AddGroupToJobTokenAllowlist adds a new group to a project's job token
// inbound groups allow list.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#add-a-group-to-a-cicd-job-token-allowlist
func (j *JobTokenScopeService) AddGroupToJobTokenAllowlist(pid any, opt *AddGroupToJobTokenAllowlistOptions, options ...RequestOptionFunc) (*JobTokenAllowlistItem, *Response, error) {
	return do[*JobTokenAllowlistItem](j.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/job_token_scope/groups_allowlist", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RemoveGroupFromJobTokenAllowlist removes a group from a project's job
// token inbound groups allow list.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_job_token_scopes/#remove-a-group-from-a-cicd-job-token-allowlist
func (j *JobTokenScopeService) RemoveGroupFromJobTokenAllowlist(pid any, targetGroup int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](j.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/job_token_scope/groups_allowlist/%d", ProjectID{pid}, targetGroup),
		withRequestOpts(options...),
	)
	return resp, err
}
