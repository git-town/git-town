//
// Copyright 2021, Igor Varavko
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

import "net/http"

type (
	PlanLimitsServiceInterface interface {
		GetCurrentPlanLimits(opt *GetCurrentPlanLimitsOptions, options ...RequestOptionFunc) (*PlanLimit, *Response, error)
		ChangePlanLimits(opt *ChangePlanLimitOptions, options ...RequestOptionFunc) (*PlanLimit, *Response, error)
	}

	// PlanLimitsService handles communication with the repositories related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/plan_limits/
	PlanLimitsService struct {
		client *Client
	}
)

var _ PlanLimitsServiceInterface = (*PlanLimitsService)(nil)

// PlanLimit represents a GitLab pipeline.
//
// GitLab API docs: https://docs.gitlab.com/api/plan_limits/
type PlanLimit struct {
	ConanMaxFileSize           int64 `json:"conan_max_file_size,omitempty"`
	GenericPackagesMaxFileSize int64 `json:"generic_packages_max_file_size,omitempty"`
	HelmMaxFileSize            int64 `json:"helm_max_file_size,omitempty"`
	MavenMaxFileSize           int64 `json:"maven_max_file_size,omitempty"`
	NPMMaxFileSize             int64 `json:"npm_max_file_size,omitempty"`
	NugetMaxFileSize           int64 `json:"nuget_max_file_size,omitempty"`
	PyPiMaxFileSize            int64 `json:"pypi_max_file_size,omitempty"`
	TerraformModuleMaxFileSize int64 `json:"terraform_module_max_file_size,omitempty"`
}

// GetCurrentPlanLimitsOptions represents the available GetCurrentPlanLimits()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/plan_limits/#get-current-plan-limits
type GetCurrentPlanLimitsOptions struct {
	PlanName *string `url:"plan_name,omitempty" json:"plan_name,omitempty"`
}

// GetCurrentPlanLimits lists the current limits of a plan on the GitLab instance.
//
// GitLab API docs:
// https://docs.gitlab.com/api/plan_limits/#get-current-plan-limits
func (s *PlanLimitsService) GetCurrentPlanLimits(opt *GetCurrentPlanLimitsOptions, options ...RequestOptionFunc) (*PlanLimit, *Response, error) {
	return do[*PlanLimit](s.client,
		withPath("application/plan_limits"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ChangePlanLimitOptions represents the available ChangePlanLimits() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/plan_limits/#change-plan-limits
type ChangePlanLimitOptions struct {
	PlanName                   *string `url:"plan_name,omitempty" json:"plan_name,omitempty"`
	ConanMaxFileSize           *int64  `url:"conan_max_file_size,omitempty" json:"conan_max_file_size,omitempty"`
	GenericPackagesMaxFileSize *int64  `url:"generic_packages_max_file_size,omitempty" json:"generic_packages_max_file_size,omitempty"`
	HelmMaxFileSize            *int64  `url:"helm_max_file_size,omitempty" json:"helm_max_file_size,omitempty"`
	MavenMaxFileSize           *int64  `url:"maven_max_file_size,omitempty" json:"maven_max_file_size,omitempty"`
	NPMMaxFileSize             *int64  `url:"npm_max_file_size,omitempty" json:"npm_max_file_size,omitempty"`
	NugetMaxFileSize           *int64  `url:"nuget_max_file_size,omitempty" json:"nuget_max_file_size,omitempty"`
	PyPiMaxFileSize            *int64  `url:"pypi_max_file_size,omitempty" json:"pypi_max_file_size,omitempty"`
	TerraformModuleMaxFileSize *int64  `url:"terraform_module_max_file_size,omitempty" json:"terraform_module_max_file_size,omitempty"`
}

// ChangePlanLimits modifies the limits of a plan on the GitLab instance.
//
// GitLab API docs:
// https://docs.gitlab.com/api/plan_limits/#change-plan-limits
func (s *PlanLimitsService) ChangePlanLimits(opt *ChangePlanLimitOptions, options ...RequestOptionFunc) (*PlanLimit, *Response, error) {
	return do[*PlanLimit](s.client,
		withMethod(http.MethodPut),
		withPath("application/plan_limits"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
