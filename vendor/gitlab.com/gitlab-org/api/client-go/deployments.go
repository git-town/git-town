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

package gitlab

import (
	"net/http"
	"time"
)

type (
	// DeploymentsServiceInterface defines all the API methods for the DeploymentsService
	DeploymentsServiceInterface interface {
		// ListProjectDeployments gets a list of deployments in a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deployments/#list-project-deployments
		// ListProjectDeployments gets a list of deployments in a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deployments/#list-project-deployments
		ListProjectDeployments(pid any, opts *ListProjectDeploymentsOptions, options ...RequestOptionFunc) ([]*Deployment, *Response, error)

		// GetProjectDeployment gets a specific deployment for a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deployments/#get-a-specific-deployment
		GetProjectDeployment(pid any, deployment int64, options ...RequestOptionFunc) (*Deployment, *Response, error)

		// CreateProjectDeployment creates a project deployment.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deployments/#create-a-deployment
		CreateProjectDeployment(pid any, opt *CreateProjectDeploymentOptions, options ...RequestOptionFunc) (*Deployment, *Response, error)

		// UpdateProjectDeployment updates a project deployment.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deployments/#update-a-deployment
		UpdateProjectDeployment(pid any, deployment int64, opt *UpdateProjectDeploymentOptions, options ...RequestOptionFunc) (*Deployment, *Response, error)

		// ApproveOrRejectProjectDeployment approves or rejects a blocked deployment.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deployments/#approve-or-reject-a-blocked-deployment
		ApproveOrRejectProjectDeployment(pid any, deployment int64, opt *ApproveOrRejectProjectDeploymentOptions, options ...RequestOptionFunc) (*Response, error)

		// DeleteProjectDeployment deletes a specific deployment.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deployments/#delete-a-specific-deployment
		DeleteProjectDeployment(pid any, deployment int64, options ...RequestOptionFunc) (*Response, error)
	}

	// DeploymentsService handles communication with the deployment related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/deployments/
	DeploymentsService struct {
		client *Client
	}
)

var _ DeploymentsServiceInterface = (*DeploymentsService)(nil)

// Deployment represents the GitLab deployment
type Deployment struct {
	ID          int64                `json:"id"`
	IID         int64                `json:"iid"`
	Ref         string               `json:"ref"`
	SHA         string               `json:"sha"`
	Status      string               `json:"status"`
	CreatedAt   *time.Time           `json:"created_at"`
	UpdatedAt   *time.Time           `json:"updated_at"`
	User        *ProjectUser         `json:"user"`
	Environment *Environment         `json:"environment"`
	Deployable  DeploymentDeployable `json:"deployable"`
}

// DeploymentDeployable represents the Gitlab deployment deployable
type DeploymentDeployable struct {
	ID         int64                        `json:"id"`
	Status     string                       `json:"status"`
	Stage      string                       `json:"stage"`
	Name       string                       `json:"name"`
	Ref        string                       `json:"ref"`
	Tag        bool                         `json:"tag"`
	Coverage   float64                      `json:"coverage"`
	CreatedAt  *time.Time                   `json:"created_at"`
	StartedAt  *time.Time                   `json:"started_at"`
	FinishedAt *time.Time                   `json:"finished_at"`
	Duration   float64                      `json:"duration"`
	User       *User                        `json:"user"`
	Commit     *Commit                      `json:"commit"`
	Pipeline   DeploymentDeployablePipeline `json:"pipeline"`
	Runner     *Runner                      `json:"runner"`
}

// DeploymentDeployablePipeline represents the Gitlab deployment deployable pipeline
type DeploymentDeployablePipeline struct {
	ID        int64      `json:"id"`
	SHA       string     `json:"sha"`
	Ref       string     `json:"ref"`
	Status    string     `json:"status"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// ListProjectDeploymentsOptions represents the available ListProjectDeployments() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deployments/#list-project-deployments
type ListProjectDeploymentsOptions struct {
	ListOptions
	OrderBy     *string `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort        *string `url:"sort,omitempty" json:"sort,omitempty"`
	Environment *string `url:"environment,omitempty" json:"environment,omitempty"`
	Status      *string `url:"status,omitempty" json:"status,omitempty"`

	// Only for GitLab versions less than 14
	UpdatedAfter  *time.Time `url:"updated_after,omitempty" json:"updated_after,omitempty"`
	UpdatedBefore *time.Time `url:"updated_before,omitempty" json:"updated_before,omitempty"`

	// Only for GitLab 14 or higher
	FinishedAfter  *time.Time `url:"finished_after,omitempty" json:"finished_after,omitempty"`
	FinishedBefore *time.Time `url:"finished_before,omitempty" json:"finished_before,omitempty"`
}

func (s *DeploymentsService) ListProjectDeployments(pid any, opts *ListProjectDeploymentsOptions, options ...RequestOptionFunc) ([]*Deployment, *Response, error) {
	return do[[]*Deployment](s.client,
		withPath("projects/%s/deployments", ProjectID{pid}),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

func (s *DeploymentsService) GetProjectDeployment(pid any, deployment int64, options ...RequestOptionFunc) (*Deployment, *Response, error) {
	return do[*Deployment](s.client,
		withPath("projects/%s/deployments/%d", ProjectID{pid}, deployment),
		withRequestOpts(options...),
	)
}

// CreateProjectDeploymentOptions represents the available
// CreateProjectDeployment() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deployments/#create-a-deployment
type CreateProjectDeploymentOptions struct {
	Environment *string                `url:"environment,omitempty" json:"environment,omitempty"`
	Ref         *string                `url:"ref,omitempty" json:"ref,omitempty"`
	SHA         *string                `url:"sha,omitempty" json:"sha,omitempty"`
	Tag         *bool                  `url:"tag,omitempty" json:"tag,omitempty"`
	Status      *DeploymentStatusValue `url:"status,omitempty" json:"status,omitempty"`
}

func (s *DeploymentsService) CreateProjectDeployment(pid any, opt *CreateProjectDeploymentOptions, options ...RequestOptionFunc) (*Deployment, *Response, error) {
	return do[*Deployment](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/deployments", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateProjectDeploymentOptions represents the available
// UpdateProjectDeployment() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deployments/#update-a-deployment
type UpdateProjectDeploymentOptions struct {
	Status *DeploymentStatusValue `url:"status,omitempty" json:"status,omitempty"`
}

func (s *DeploymentsService) UpdateProjectDeployment(pid any, deployment int64, opt *UpdateProjectDeploymentOptions, options ...RequestOptionFunc) (*Deployment, *Response, error) {
	return do[*Deployment](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/deployments/%d", ProjectID{pid}, deployment),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ApproveOrRejectProjectDeploymentOptions represents the available
// ApproveOrRejectProjectDeployment() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deployments/#approve-or-reject-a-blocked-deployment
type ApproveOrRejectProjectDeploymentOptions struct {
	Status        *DeploymentApprovalStatus `url:"status,omitempty" json:"status,omitempty"`
	Comment       *string                   `url:"comment,omitempty" json:"comment,omitempty"`
	RepresentedAs *string                   `url:"represented_as,omitempty" json:"represented_as,omitempty"`
}

func (s *DeploymentsService) ApproveOrRejectProjectDeployment(pid any, deployment int64, opt *ApproveOrRejectProjectDeploymentOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/deployments/%d/approval", ProjectID{pid}, deployment),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *DeploymentsService) DeleteProjectDeployment(pid any, deployment int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/deployments/%d", ProjectID{pid}, deployment),
		withRequestOpts(options...),
	)
	return resp, err
}
