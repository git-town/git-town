//
// Copyright 2022, Timo Furrer <tuxtimo@gmail.com>
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
	ClusterAgentsServiceInterface interface {
		// ListAgents returns a list of agents registered for the project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/cluster_agents/#list-the-agents-for-a-project
		// ListAgents returns a list of agents registered for the project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/cluster_agents/#list-the-agents-for-a-project
		ListAgents(pid any, opt *ListAgentsOptions, options ...RequestOptionFunc) ([]*Agent, *Response, error)

		// GetAgent gets a single agent details.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/cluster_agents/#get-details-about-an-agent
		GetAgent(pid any, id int64, options ...RequestOptionFunc) (*Agent, *Response, error)

		// RegisterAgent registers an agent to the project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/cluster_agents/#register-an-agent-with-a-project
		RegisterAgent(pid any, opt *RegisterAgentOptions, options ...RequestOptionFunc) (*Agent, *Response, error)

		// DeleteAgent deletes an existing agent registration.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/cluster_agents/#delete-a-registered-agent
		DeleteAgent(pid any, id int64, options ...RequestOptionFunc) (*Response, error)

		// ListAgentTokens returns a list of tokens for an agent.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/cluster_agents/#list-tokens-for-an-agent
		ListAgentTokens(pid any, aid int64, opt *ListAgentTokensOptions, options ...RequestOptionFunc) ([]*AgentToken, *Response, error)

		// GetAgentToken gets a single agent token.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/cluster_agents/#get-a-single-agent-token
		GetAgentToken(pid any, aid int64, id int64, options ...RequestOptionFunc) (*AgentToken, *Response, error)

		// CreateAgentToken creates a new token for an agent.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/cluster_agents/#create-an-agent-token
		CreateAgentToken(pid any, aid int64, opt *CreateAgentTokenOptions, options ...RequestOptionFunc) (*AgentToken, *Response, error)

		// RevokeAgentToken revokes an agent token.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/cluster_agents/#revoke-an-agent-token
		RevokeAgentToken(pid any, aid int64, id int64, options ...RequestOptionFunc) (*Response, error)
	}

	// ClusterAgentsService handles communication with the cluster agents related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/cluster_agents/
	ClusterAgentsService struct {
		client *Client
	}
)

var _ ClusterAgentsServiceInterface = (*ClusterAgentsService)(nil)

// Agent represents a GitLab agent for Kubernetes.
//
// GitLab API docs: https://docs.gitlab.com/api/cluster_agents/
type Agent struct {
	ID              int64         `json:"id"`
	Name            string        `json:"name"`
	CreatedAt       *time.Time    `json:"created_at"`
	CreatedByUserID int64         `json:"created_by_user_id"`
	ConfigProject   ConfigProject `json:"config_project"`
}

type ConfigProject struct {
	ID                int64      `json:"id"`
	Description       string     `json:"description"`
	Name              string     `json:"name"`
	NameWithNamespace string     `json:"name_with_namespace"`
	Path              string     `json:"path"`
	PathWithNamespace string     `json:"path_with_namespace"`
	CreatedAt         *time.Time `json:"created_at"`
}

func (a Agent) String() string {
	return Stringify(a)
}

// AgentToken represents a GitLab agent token.
//
// GitLab API docs:
// https://docs.gitlab.com/api/cluster_agents/#list-tokens-for-an-agent
type AgentToken struct {
	ID              int64      `json:"id"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	AgentID         int64      `json:"agent_id"`
	Status          string     `json:"status"`
	CreatedAt       *time.Time `json:"created_at"`
	CreatedByUserID int64      `json:"created_by_user_id"`
	LastUsedAt      *time.Time `json:"last_used_at"`
	Token           string     `json:"token"`
}

func (a AgentToken) String() string {
	return Stringify(a)
}

// ListAgentsOptions represents the available ListAgents() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/cluster_agents/#list-the-agents-for-a-project
type ListAgentsOptions struct {
	ListOptions
}

func (s *ClusterAgentsService) ListAgents(pid any, opt *ListAgentsOptions, options ...RequestOptionFunc) ([]*Agent, *Response, error) {
	return do[[]*Agent](s.client,
		withPath("projects/%s/cluster_agents", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ClusterAgentsService) GetAgent(pid any, id int64, options ...RequestOptionFunc) (*Agent, *Response, error) {
	return do[*Agent](s.client,
		withPath("projects/%s/cluster_agents/%d", ProjectID{pid}, id),
		withRequestOpts(options...),
	)
}

// RegisterAgentOptions represents the available RegisterAgent()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/cluster_agents/#register-an-agent-with-a-project
type RegisterAgentOptions struct {
	Name *string `url:"name,omitempty" json:"name,omitempty"`
}

func (s *ClusterAgentsService) RegisterAgent(pid any, opt *RegisterAgentOptions, options ...RequestOptionFunc) (*Agent, *Response, error) {
	return do[*Agent](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/cluster_agents", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ClusterAgentsService) DeleteAgent(pid any, id int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/cluster_agents/%d", ProjectID{pid}, id),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListAgentTokensOptions represents the available ListAgentTokens() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/cluster_agents/#list-tokens-for-an-agent
type ListAgentTokensOptions struct {
	ListOptions
}

func (s *ClusterAgentsService) ListAgentTokens(pid any, aid int64, opt *ListAgentTokensOptions, options ...RequestOptionFunc) ([]*AgentToken, *Response, error) {
	return do[[]*AgentToken](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/cluster_agents/%d/tokens", ProjectID{pid}, aid),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ClusterAgentsService) GetAgentToken(pid any, aid int64, id int64, options ...RequestOptionFunc) (*AgentToken, *Response, error) {
	return do[*AgentToken](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/cluster_agents/%d/tokens/%d", ProjectID{pid}, aid, id),
		withRequestOpts(options...),
	)
}

// CreateAgentTokenOptions represents the available CreateAgentToken() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/cluster_agents/#create-an-agent-token
type CreateAgentTokenOptions struct {
	Name        *string `url:"name,omitempty" json:"name,omitempty"`
	Description *string `url:"description,omitempty" json:"description,omitempty"`
}

func (s *ClusterAgentsService) CreateAgentToken(pid any, aid int64, opt *CreateAgentTokenOptions, options ...RequestOptionFunc) (*AgentToken, *Response, error) {
	return do[*AgentToken](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/cluster_agents/%d/tokens", ProjectID{pid}, aid),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ClusterAgentsService) RevokeAgentToken(pid any, aid int64, id int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/cluster_agents/%d/tokens/%d", ProjectID{pid}, aid, id),
		withRequestOpts(options...),
	)
	return resp, err
}
