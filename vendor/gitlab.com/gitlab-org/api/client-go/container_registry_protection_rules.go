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
)

type (
	ContainerRegistryProtectionRulesServiceInterface interface {
		// ListContainerRegistryProtectionRules gets a list of container repository
		// protection rules from a project’s container registry.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_repository_protection_rules/#list-container-repository-protection-rules
		// ListContainerRegistryProtectionRules gets a list of container repository
		// protection rules from a project’s container registry.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_repository_protection_rules/#list-container-repository-protection-rules
		ListContainerRegistryProtectionRules(pid any, options ...RequestOptionFunc) ([]*ContainerRegistryProtectionRule, *Response, error)

		// CreateContainerRegistryProtectionRule creates a container repository
		// protection rule for a project’s container registry.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_repository_protection_rules/#create-a-container-repository-protection-rule

		// CreateContainerRegistryProtectionRule creates a container repository
		// protection rule for a project’s container registry.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_repository_protection_rules/#create-a-container-repository-protection-rule
		CreateContainerRegistryProtectionRule(pid any, opt *CreateContainerRegistryProtectionRuleOptions, options ...RequestOptionFunc) (*ContainerRegistryProtectionRule, *Response, error)

		// UpdateContainerRegistryProtectionRule updates a container repository protection
		// rule for a project’s container registry.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_repository_protection_rules/#update-a-container-repository-protection-rule
		UpdateContainerRegistryProtectionRule(pid any, ruleID int64, opt *UpdateContainerRegistryProtectionRuleOptions, options ...RequestOptionFunc) (*ContainerRegistryProtectionRule, *Response, error)

		// DeleteContainerRegistryProtectionRule deletes a container repository protection
		// rule from a project’s container registry.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_repository_protection_rules/#delete-a-container-repository-protection-rule
		DeleteContainerRegistryProtectionRule(pid any, ruleID int64, options ...RequestOptionFunc) (*Response, error)
	}

	// ContainerRegistryProtectionRulesService handles communication with
	// the container registry protection rules related methods of the
	// GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/container_repository_protection_rules/
	ContainerRegistryProtectionRulesService struct {
		client *Client
	}
)

var _ ContainerRegistryProtectionRulesServiceInterface = (*ContainerRegistryProtectionRulesService)(nil)

// ContainerRegistryProtectionRule represents a GitLab container registry
// protection rule.
//
// GitLab API docs:
// https://docs.gitlab.com/api/container_repository_protection_rules/
type ContainerRegistryProtectionRule struct {
	ID                          int64                     `json:"id"`
	ProjectID                   int64                     `json:"project_id"`
	RepositoryPathPattern       string                    `json:"repository_path_pattern"`
	MinimumAccessLevelForPush   ProtectionRuleAccessLevel `json:"minimum_access_level_for_push"`
	MinimumAccessLevelForDelete ProtectionRuleAccessLevel `json:"minimum_access_level_for_delete"`
}

func (s ContainerRegistryProtectionRule) String() string {
	return Stringify(s)
}

func (s *ContainerRegistryProtectionRulesService) ListContainerRegistryProtectionRules(pid any, options ...RequestOptionFunc) ([]*ContainerRegistryProtectionRule, *Response, error) {
	return do[[]*ContainerRegistryProtectionRule](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/registry/protection/repository/rules", ProjectID{pid}),
		withRequestOpts(options...),
	)
}

// CreateContainerRegistryProtectionRuleOptions represents the available
// CreateContainerRegistryProtectionRule() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/container_repository_protection_rules/#create-a-container-repository-protection-rule
type CreateContainerRegistryProtectionRuleOptions struct {
	RepositoryPathPattern       *string                    `url:"repository_path_pattern,omitempty" json:"repository_path_pattern,omitempty"`
	MinimumAccessLevelForPush   *ProtectionRuleAccessLevel `url:"minimum_access_level_for_push,omitempty" json:"minimum_access_level_for_push,omitempty"`
	MinimumAccessLevelForDelete *ProtectionRuleAccessLevel `url:"minimum_access_level_for_delete,omitempty" json:"minimum_access_level_for_delete,omitempty"`
}

func (s *ContainerRegistryProtectionRulesService) CreateContainerRegistryProtectionRule(pid any, opt *CreateContainerRegistryProtectionRuleOptions, options ...RequestOptionFunc) (*ContainerRegistryProtectionRule, *Response, error) {
	return do[*ContainerRegistryProtectionRule](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/registry/protection/repository/rules", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateContainerRegistryProtectionRuleOptions represents the available
// UpdateContainerRegistryProtectionRule() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/container_repository_protection_rules/#update-a-container-repository-protection-rule
type UpdateContainerRegistryProtectionRuleOptions struct {
	RepositoryPathPattern       *string                    `url:"repository_path_pattern,omitempty" json:"repository_path_pattern,omitempty"`
	MinimumAccessLevelForPush   *ProtectionRuleAccessLevel `url:"minimum_access_level_for_push,omitempty" json:"minimum_access_level_for_push,omitempty"`
	MinimumAccessLevelForDelete *ProtectionRuleAccessLevel `url:"minimum_access_level_for_delete,omitempty" json:"minimum_access_level_for_delete,omitempty"`
}

func (s *ContainerRegistryProtectionRulesService) UpdateContainerRegistryProtectionRule(pid any, ruleID int64, opt *UpdateContainerRegistryProtectionRuleOptions, options ...RequestOptionFunc) (*ContainerRegistryProtectionRule, *Response, error) {
	return do[*ContainerRegistryProtectionRule](s.client,
		withMethod(http.MethodPatch),
		withPath("projects/%s/registry/protection/repository/rules/%d", ProjectID{pid}, ruleID),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ContainerRegistryProtectionRulesService) DeleteContainerRegistryProtectionRule(pid any, ruleID int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/registry/protection/repository/rules/%d", ProjectID{pid}, ruleID),
		withRequestOpts(options...),
	)
	return resp, err
}
