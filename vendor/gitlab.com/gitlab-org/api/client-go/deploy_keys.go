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
	DeployKeysServiceInterface interface {
		ListAllDeployKeys(opt *ListInstanceDeployKeysOptions, options ...RequestOptionFunc) ([]*InstanceDeployKey, *Response, error)
		AddInstanceDeployKey(opt *AddInstanceDeployKeyOptions, options ...RequestOptionFunc) (*InstanceDeployKey, *Response, error)
		ListProjectDeployKeys(pid any, opt *ListProjectDeployKeysOptions, options ...RequestOptionFunc) ([]*ProjectDeployKey, *Response, error)
		ListUserProjectDeployKeys(uid any, opt *ListUserProjectDeployKeysOptions, options ...RequestOptionFunc) ([]*ProjectDeployKey, *Response, error)
		GetDeployKey(pid any, deployKey int, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error)
		AddDeployKey(pid any, opt *AddDeployKeyOptions, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error)
		DeleteDeployKey(pid any, deployKey int, options ...RequestOptionFunc) (*Response, error)
		EnableDeployKey(pid any, deployKey int, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error)
		UpdateDeployKey(pid any, deployKey int, opt *UpdateDeployKeyOptions, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error)
	}

	// DeployKeysService handles communication with the keys related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/deploy_keys/
	DeployKeysService struct {
		client *Client
	}
)

var _ DeployKeysServiceInterface = (*DeployKeysService)(nil)

// InstanceDeployKey represents a GitLab deploy key with the associated
// projects it has write access to.
type InstanceDeployKey struct {
	ID                         int                 `json:"id"`
	Title                      string              `json:"title"`
	CreatedAt                  *time.Time          `json:"created_at"`
	ExpiresAt                  *time.Time          `json:"expires_at"`
	Key                        string              `json:"key"`
	Fingerprint                string              `json:"fingerprint"`
	FingerprintSHA256          string              `json:"fingerprint_sha256"`
	ProjectsWithWriteAccess    []*DeployKeyProject `json:"projects_with_write_access"`
	ProjectsWithReadonlyAccess []*DeployKeyProject `json:"projects_with_readonly_access"`
}

func (k InstanceDeployKey) String() string {
	return Stringify(k)
}

// DeployKeyProject refers to a project an InstanceDeployKey has write access to.
type DeployKeyProject struct {
	ID                int        `json:"id"`
	Description       string     `json:"description"`
	Name              string     `json:"name"`
	NameWithNamespace string     `json:"name_with_namespace"`
	Path              string     `json:"path"`
	PathWithNamespace string     `json:"path_with_namespace"`
	CreatedAt         *time.Time `json:"created_at"`
}

func (k DeployKeyProject) String() string {
	return Stringify(k)
}

// ProjectDeployKey represents a GitLab project deploy key.
type ProjectDeployKey struct {
	ID                int        `json:"id"`
	Title             string     `json:"title"`
	Key               string     `json:"key"`
	Fingerprint       string     `json:"fingerprint"`
	FingerprintSHA256 string     `json:"fingerprint_sha256"`
	CreatedAt         *time.Time `json:"created_at"`
	CanPush           bool       `json:"can_push"`
	ExpiresAt         *time.Time `json:"expires_at"`
}

func (k ProjectDeployKey) String() string {
	return Stringify(k)
}

// ListProjectDeployKeysOptions represents the available ListAllDeployKeys()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#list-all-deploy-keys
type ListInstanceDeployKeysOptions struct {
	ListOptions
	Public *bool `url:"public,omitempty" json:"public,omitempty"`
}

// ListAllDeployKeys gets a list of all deploy keys
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#list-all-deploy-keys
func (s *DeployKeysService) ListAllDeployKeys(opt *ListInstanceDeployKeysOptions, options ...RequestOptionFunc) ([]*InstanceDeployKey, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "deploy_keys", opt, options)
	if err != nil {
		return nil, nil, err
	}

	var ks []*InstanceDeployKey
	resp, err := s.client.Do(req, &ks)
	if err != nil {
		return nil, resp, err
	}

	return ks, resp, nil
}

// AddInstanceDeployKeyOptions represents the available AddInstanceDeployKey()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#add-deploy-key
type AddInstanceDeployKeyOptions struct {
	Key       *string    `url:"key,omitempty" json:"key,omitempty"`
	Title     *string    `url:"title,omitempty" json:"title,omitempty"`
	ExpiresAt *time.Time `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// AddInstanceDeployKey creates a deploy key for the GitLab instance.
// Requires administrator access.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#add-deploy-key
func (s *DeployKeysService) AddInstanceDeployKey(opt *AddInstanceDeployKeyOptions, options ...RequestOptionFunc) (*InstanceDeployKey, *Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "deploy_keys", opt, options)
	if err != nil {
		return nil, nil, err
	}

	key := new(InstanceDeployKey)
	resp, err := s.client.Do(req, &key)
	if err != nil {
		return nil, resp, err
	}

	return key, resp, nil
}

// ListProjectDeployKeysOptions represents the available ListProjectDeployKeys()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#list-deploy-keys-for-project
type ListProjectDeployKeysOptions ListOptions

// ListProjectDeployKeys gets a list of a project's deploy keys
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#list-deploy-keys-for-project
func (s *DeployKeysService) ListProjectDeployKeys(pid any, opt *ListProjectDeployKeysOptions, options ...RequestOptionFunc) ([]*ProjectDeployKey, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/deploy_keys", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var ks []*ProjectDeployKey
	resp, err := s.client.Do(req, &ks)
	if err != nil {
		return nil, resp, err
	}

	return ks, resp, nil
}

// ListUserProjectDeployKeysOptions represents the available ListUserProjectDeployKeys()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#list-project-deploy-keys-for-user
type ListUserProjectDeployKeysOptions ListOptions

// ListUserProjectDeployKeys gets a list of a user's deploy keys
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#list-project-deploy-keys-for-user
func (s *DeployKeysService) ListUserProjectDeployKeys(uid any, opt *ListUserProjectDeployKeysOptions, options ...RequestOptionFunc) ([]*ProjectDeployKey, *Response, error) {
	user, err := parseID(uid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("users/%s/project_deploy_keys", PathEscape(user))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var ks []*ProjectDeployKey
	resp, err := s.client.Do(req, &ks)
	if err != nil {
		return nil, resp, err
	}

	return ks, resp, nil
}

// GetDeployKey gets a single deploy key.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#get-a-single-deploy-key
func (s *DeployKeysService) GetDeployKey(pid any, deployKey int, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/deploy_keys/%d", PathEscape(project), deployKey)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	k := new(ProjectDeployKey)
	resp, err := s.client.Do(req, k)
	if err != nil {
		return nil, resp, err
	}

	return k, resp, nil
}

// AddDeployKeyOptions represents the available ADDDeployKey() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#add-deploy-key-for-a-project
type AddDeployKeyOptions struct {
	Key       *string    `url:"key,omitempty" json:"key,omitempty"`
	Title     *string    `url:"title,omitempty" json:"title,omitempty"`
	CanPush   *bool      `url:"can_push,omitempty" json:"can_push,omitempty"`
	ExpiresAt *time.Time `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

// AddDeployKey creates a new deploy key for a project. If deploy key already
// exists in another project - it will be joined to project but only if
// original one is accessible by the same user.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#add-deploy-key-for-a-project
func (s *DeployKeysService) AddDeployKey(pid any, opt *AddDeployKeyOptions, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/deploy_keys", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	k := new(ProjectDeployKey)
	resp, err := s.client.Do(req, k)
	if err != nil {
		return nil, resp, err
	}

	return k, resp, nil
}

// DeleteDeployKey deletes a deploy key from a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#delete-deploy-key
func (s *DeployKeysService) DeleteDeployKey(pid any, deployKey int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/deploy_keys/%d", PathEscape(project), deployKey)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// EnableDeployKey enables a deploy key.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#enable-a-deploy-key
func (s *DeployKeysService) EnableDeployKey(pid any, deployKey int, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/deploy_keys/%d/enable", PathEscape(project), deployKey)

	req, err := s.client.NewRequest(http.MethodPost, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	k := new(ProjectDeployKey)
	resp, err := s.client.Do(req, k)
	if err != nil {
		return nil, resp, err
	}

	return k, resp, nil
}

// UpdateDeployKeyOptions represents the available UpdateDeployKey() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#update-deploy-key
type UpdateDeployKeyOptions struct {
	Title   *string `url:"title,omitempty" json:"title,omitempty"`
	CanPush *bool   `url:"can_push,omitempty" json:"can_push,omitempty"`
}

// UpdateDeployKey updates a deploy key for a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#update-deploy-key
func (s *DeployKeysService) UpdateDeployKey(pid any, deployKey int, opt *UpdateDeployKeyOptions, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/deploy_keys/%d", PathEscape(project), deployKey)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	k := new(ProjectDeployKey)
	resp, err := s.client.Do(req, k)
	if err != nil {
		return nil, resp, err
	}

	return k, resp, nil
}
