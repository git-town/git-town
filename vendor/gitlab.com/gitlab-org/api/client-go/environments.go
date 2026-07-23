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
	"net/http"
	"time"
)

type (
	// EnvironmentsServiceInterface defines all the API methods for the EnvironmentsService
	EnvironmentsServiceInterface interface {
		// ListEnvironments gets a list of environments from a project, sorted by name
		// alphabetically.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/environments/#list-environments
		ListEnvironments(pid any, opts *ListEnvironmentsOptions, options ...RequestOptionFunc) ([]*Environment, *Response, error)
		GetEnvironment(pid any, environment int64, options ...RequestOptionFunc) (*Environment, *Response, error)
		CreateEnvironment(pid any, opt *CreateEnvironmentOptions, options ...RequestOptionFunc) (*Environment, *Response, error)
		EditEnvironment(pid any, environment int64, opt *EditEnvironmentOptions, options ...RequestOptionFunc) (*Environment, *Response, error)
		DeleteEnvironment(pid any, environment int64, options ...RequestOptionFunc) (*Response, error)
		StopEnvironment(pid any, environmentID int64, opt *StopEnvironmentOptions, options ...RequestOptionFunc) (*Environment, *Response, error)
	}

	// EnvironmentsService handles communication with the environment related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/environments/
	EnvironmentsService struct {
		client *Client
	}
)

var _ EnvironmentsServiceInterface = (*EnvironmentsService)(nil)

// Environment represents a GitLab environment.
//
// GitLab API docs: https://docs.gitlab.com/api/environments/
type Environment struct {
	ID                  int64       `json:"id"`
	Name                string      `json:"name"`
	Slug                string      `json:"slug"`
	Description         string      `json:"description"`
	State               string      `json:"state"`
	Tier                string      `json:"tier"`
	ExternalURL         string      `json:"external_url"`
	Project             *Project    `json:"project"`
	CreatedAt           *time.Time  `json:"created_at"`
	UpdatedAt           *time.Time  `json:"updated_at"`
	LastDeployment      *Deployment `json:"last_deployment"`
	ClusterAgent        *Agent      `json:"cluster_agent"`
	KubernetesNamespace string      `json:"kubernetes_namespace"`
	FluxResourcePath    string      `json:"flux_resource_path"`
	AutoStopAt          *time.Time  `json:"auto_stop_at"`
	AutoStopSetting     string      `json:"auto_stop_setting"`
}

func (env Environment) String() string {
	return Stringify(env)
}

// ListEnvironmentsOptions represents the available ListEnvironments() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/environments/#list-environments
type ListEnvironmentsOptions struct {
	ListOptions
	Name   *string `url:"name,omitempty" json:"name,omitempty"`
	Search *string `url:"search,omitempty" json:"search,omitempty"`
	States *string `url:"states,omitempty" json:"states,omitempty"`
}

func (s *EnvironmentsService) ListEnvironments(pid any, opts *ListEnvironmentsOptions, options ...RequestOptionFunc) ([]*Environment, *Response, error) {
	return do[[]*Environment](s.client,
		withPath("projects/%s/environments", ProjectID{pid}),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

func (s *EnvironmentsService) GetEnvironment(pid any, environment int64, options ...RequestOptionFunc) (*Environment, *Response, error) {
	return do[*Environment](s.client,
		withPath("projects/%s/environments/%d", ProjectID{pid}, environment),
		withRequestOpts(options...),
	)
}

// CreateEnvironmentOptions represents the available CreateEnvironment() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/environments/#create-a-new-environment
type CreateEnvironmentOptions struct {
	Name                *string `url:"name,omitempty" json:"name,omitempty"`
	Description         *string `url:"description,omitempty" json:"description,omitempty"`
	ExternalURL         *string `url:"external_url,omitempty" json:"external_url,omitempty"`
	Tier                *string `url:"tier,omitempty" json:"tier,omitempty"`
	ClusterAgentID      *int64  `url:"cluster_agent_id,omitempty" json:"cluster_agent_id,omitempty"`
	KubernetesNamespace *string `url:"kubernetes_namespace,omitempty" json:"kubernetes_namespace,omitempty"`
	FluxResourcePath    *string `url:"flux_resource_path,omitempty" json:"flux_resource_path,omitempty"`
	AutoStopSetting     *string `url:"auto_stop_setting,omitempty" json:"auto_stop_setting,omitempty"`
}

func (s *EnvironmentsService) CreateEnvironment(pid any, opt *CreateEnvironmentOptions, options ...RequestOptionFunc) (*Environment, *Response, error) {
	return do[*Environment](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/environments", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// EditEnvironmentOptions represents the available EditEnvironment() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/environments/#update-an-existing-environment
type EditEnvironmentOptions struct {
	Name                *string `url:"name,omitempty" json:"name,omitempty"`
	Description         *string `url:"description,omitempty" json:"description,omitempty"`
	ExternalURL         *string `url:"external_url,omitempty" json:"external_url,omitempty"`
	Tier                *string `url:"tier,omitempty" json:"tier,omitempty"`
	ClusterAgentID      *int64  `url:"cluster_agent_id,omitempty" json:"cluster_agent_id,omitempty"`
	KubernetesNamespace *string `url:"kubernetes_namespace,omitempty" json:"kubernetes_namespace,omitempty"`
	FluxResourcePath    *string `url:"flux_resource_path,omitempty" json:"flux_resource_path,omitempty"`
	AutoStopSetting     *string `url:"auto_stop_setting,omitempty" json:"auto_stop_setting,omitempty"`
}

func (s *EnvironmentsService) EditEnvironment(pid any, environment int64, opt *EditEnvironmentOptions, options ...RequestOptionFunc) (*Environment, *Response, error) {
	return do[*Environment](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/environments/%d", ProjectID{pid}, environment),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *EnvironmentsService) DeleteEnvironment(pid any, environment int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/environments/%d", ProjectID{pid}, environment),
		withRequestOpts(options...),
	)
	return resp, err
}

// StopEnvironmentOptions represents the available StopEnvironment() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/environments/#stop-an-environment
type StopEnvironmentOptions struct {
	Force *bool `url:"force,omitempty" json:"force,omitempty"`
}

func (s *EnvironmentsService) StopEnvironment(pid any, environmentID int64, opt *StopEnvironmentOptions, options ...RequestOptionFunc) (*Environment, *Response, error) {
	return do[*Environment](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/environments/%d/stop", ProjectID{pid}, environmentID),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
