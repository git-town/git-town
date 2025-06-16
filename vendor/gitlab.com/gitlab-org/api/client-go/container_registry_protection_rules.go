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
	ContainerRegistryProtectionRulesServiceInterface interface {
		ListContainerRegistryProtectionRules(pid any, options ...RequestOptionFunc) ([]*ContainerRegistryProtectionRule, *Response, error)
		CreateContainerRegistryProtectionRule(pid any, opt *CreateContainerRegistryProtectionRuleOptions, options ...RequestOptionFunc) (*ContainerRegistryProtectionRule, *Response, error)
		UpdateContainerRegistryProtectionRule(pid any, ruleID int, opt *UpdateContainerRegistryProtectionRuleOptions, options ...RequestOptionFunc) (*ContainerRegistryProtectionRule, *Response, error)
		DeleteContainerRegistryProtectionRule(pid any, ruleID int, options ...RequestOptionFunc) (*Response, error)
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
	ID                          int                       `json:"id"`
	ProjectID                   int                       `json:"project_id"`
	RepositoryPathPattern       string                    `json:"repository_path_pattern"`
	MinimumAccessLevelForPush   ProtectionRuleAccessLevel `json:"minimum_access_level_for_push"`
	MinimumAccessLevelForDelete ProtectionRuleAccessLevel `json:"minimum_access_level_for_delete"`
}

func (s ContainerRegistryProtectionRule) String() string {
	return Stringify(s)
}

// ListContainerRegistryProtectionRules gets a list of container repository
// protection rules from a project’s container registry.
//
// GitLab API docs:
// https://docs.gitlab.com/api/container_repository_protection_rules/#list-container-repository-protection-rules
func (s *ContainerRegistryProtectionRulesService) ListContainerRegistryProtectionRules(pid any, options ...RequestOptionFunc) ([]*ContainerRegistryProtectionRule, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/registry/protection/repository/rules", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var rules []*ContainerRegistryProtectionRule
	resp, err := s.client.Do(req, &rules)
	if err != nil {
		return nil, resp, err
	}

	return rules, resp, nil
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

// CreateContainerRegistryProtectionRule creates a container repository
// protection rule for a project’s container registry.
//
// GitLab API docs:
// https://docs.gitlab.com/api/container_repository_protection_rules/#create-a-container-repository-protection-rule
func (s *ContainerRegistryProtectionRulesService) CreateContainerRegistryProtectionRule(pid any, opt *CreateContainerRegistryProtectionRuleOptions, options ...RequestOptionFunc) (*ContainerRegistryProtectionRule, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/registry/protection/repository/rules", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	rule := new(ContainerRegistryProtectionRule)
	resp, err := s.client.Do(req, rule)
	if err != nil {
		return nil, resp, err
	}

	return rule, resp, nil
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

// UpdateContainerRegistryProtectionRule updates a container repository protection
// rule for a project’s container registry.
//
// GitLab API docs:
// https://docs.gitlab.com/api/container_repository_protection_rules/#update-a-container-repository-protection-rule
func (s *ContainerRegistryProtectionRulesService) UpdateContainerRegistryProtectionRule(pid any, ruleID int, opt *UpdateContainerRegistryProtectionRuleOptions, options ...RequestOptionFunc) (*ContainerRegistryProtectionRule, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/registry/protection/repository/rules/%d", PathEscape(project), ruleID)

	req, err := s.client.NewRequest(http.MethodPatch, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	rule := new(ContainerRegistryProtectionRule)
	resp, err := s.client.Do(req, rule)
	if err != nil {
		return nil, resp, err
	}

	return rule, resp, nil
}

// DeleteContainerRegistryProtectionRule deletes a container repository protection
// rule from a project’s container registry.
//
// GitLab API docs:
// https://docs.gitlab.com/api/container_repository_protection_rules/#delete-a-container-repository-protection-rule
func (s *ContainerRegistryProtectionRulesService) DeleteContainerRegistryProtectionRule(pid any, ruleID int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/registry/protection/repository/rules/%d", PathEscape(project), ruleID)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
