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
	DeployKeysServiceInterface interface {
		// ListAllDeployKeys gets a list of all deploy keys.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_keys/#list-all-deploy-keys
		// ListAllDeployKeys gets a list of all deploy keys.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_keys/#list-all-deploy-keys
		ListAllDeployKeys(opt *ListInstanceDeployKeysOptions, options ...RequestOptionFunc) ([]*InstanceDeployKey, *Response, error)

		// AddInstanceDeployKey creates a deploy key for the GitLab instance.
		// Requires administrator access.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_keys/#add-deploy-key

		// AddInstanceDeployKey creates a deploy key for the GitLab instance.
		// Requires administrator access.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_keys/#add-deploy-key
		AddInstanceDeployKey(opt *AddInstanceDeployKeyOptions, options ...RequestOptionFunc) (*InstanceDeployKey, *Response, error)

		// ListProjectDeployKeys gets a list of a project's deploy keys.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_keys/#list-deploy-keys-for-project

		// ListProjectDeployKeys gets a list of a project's deploy keys.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_keys/#list-deploy-keys-for-project
		ListProjectDeployKeys(pid any, opt *ListProjectDeployKeysOptions, options ...RequestOptionFunc) ([]*ProjectDeployKey, *Response, error)

		// ListUserProjectDeployKeys gets a list of a user's deploy keys.
		//
		// uid can be either a user ID (int) or a username (string). If a username
		// is provided with a leading "@" (e.g., "@johndoe"), it will be trimmed.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_keys/#list-project-deploy-keys-for-user
		ListUserProjectDeployKeys(uid any, opt *ListUserProjectDeployKeysOptions, options ...RequestOptionFunc) ([]*ProjectDeployKey, *Response, error)

		// GetDeployKey gets a single deploy key.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_keys/#get-a-single-deploy-key
		GetDeployKey(pid any, deployKey int64, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error)

		// AddDeployKey creates a new deploy key for a project. If the deploy key already
		// exists in another project, it will be joined to the project but only if
		// the original one is accessible by the same user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_keys/#add-deploy-key-for-a-project
		AddDeployKey(pid any, opt *AddDeployKeyOptions, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error)

		// DeleteDeployKey deletes a deploy key from a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_keys/#delete-deploy-key
		DeleteDeployKey(pid any, deployKey int64, options ...RequestOptionFunc) (*Response, error)

		// EnableDeployKey enables a deploy key.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_keys/#enable-a-deploy-key
		EnableDeployKey(pid any, deployKey int64, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error)

		// UpdateDeployKey updates a deploy key for a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deploy_keys/#update-deploy-key
		UpdateDeployKey(pid any, deployKey int64, opt *UpdateDeployKeyOptions, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error)
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
	ID                         int64               `json:"id"`
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
	ID                int64      `json:"id"`
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
	ID                int64      `json:"id"`
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

// ListInstanceDeployKeysOptions represents the available ListAllDeployKeys()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#list-all-deploy-keys
type ListInstanceDeployKeysOptions struct {
	ListOptions
	Public *bool `url:"public,omitempty" json:"public,omitempty"`
}

func (s *DeployKeysService) ListAllDeployKeys(opt *ListInstanceDeployKeysOptions, options ...RequestOptionFunc) ([]*InstanceDeployKey, *Response, error) {
	return do[[]*InstanceDeployKey](s.client,
		withPath("deploy_keys"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
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

func (s *DeployKeysService) AddInstanceDeployKey(opt *AddInstanceDeployKeyOptions, options ...RequestOptionFunc) (*InstanceDeployKey, *Response, error) {
	return do[*InstanceDeployKey](s.client,
		withMethod(http.MethodPost),
		withPath("deploy_keys"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListProjectDeployKeysOptions represents the available ListProjectDeployKeys()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#list-deploy-keys-for-project
type ListProjectDeployKeysOptions struct {
	ListOptions
}

func (s *DeployKeysService) ListProjectDeployKeys(pid any, opt *ListProjectDeployKeysOptions, options ...RequestOptionFunc) ([]*ProjectDeployKey, *Response, error) {
	return do[[]*ProjectDeployKey](s.client,
		withPath("projects/%s/deploy_keys", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListUserProjectDeployKeysOptions represents the available ListUserProjectDeployKeys()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#list-project-deploy-keys-for-user
type ListUserProjectDeployKeysOptions struct {
	ListOptions
}

// ListUserProjectDeployKeys gets a list of a user's deploy keys.
//
// uid can be either a user ID (int) or a username (string). If a username
// is provided with a leading "@" (e.g., "@johndoe"), it will be trimmed.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#list-project-deploy-keys-for-user
func (s *DeployKeysService) ListUserProjectDeployKeys(uid any, opt *ListUserProjectDeployKeysOptions, options ...RequestOptionFunc) ([]*ProjectDeployKey, *Response, error) {
	return do[[]*ProjectDeployKey](s.client,
		withPath("users/%s/project_deploy_keys", UserID{uid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DeployKeysService) GetDeployKey(pid any, deployKey int64, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error) {
	return do[*ProjectDeployKey](s.client,
		withPath("projects/%s/deploy_keys/%d", ProjectID{pid}, deployKey),
		withRequestOpts(options...),
	)
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

func (s *DeployKeysService) AddDeployKey(pid any, opt *AddDeployKeyOptions, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error) {
	return do[*ProjectDeployKey](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/deploy_keys", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DeployKeysService) DeleteDeployKey(pid any, deployKey int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/deploy_keys/%d", ProjectID{pid}, deployKey),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *DeployKeysService) EnableDeployKey(pid any, deployKey int64, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error) {
	return do[*ProjectDeployKey](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/deploy_keys/%d/enable", ProjectID{pid}, deployKey),
		withRequestOpts(options...),
	)
}

// UpdateDeployKeyOptions represents the available UpdateDeployKey() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/deploy_keys/#update-deploy-key
type UpdateDeployKeyOptions struct {
	Title   *string `url:"title,omitempty" json:"title,omitempty"`
	CanPush *bool   `url:"can_push,omitempty" json:"can_push,omitempty"`
}

func (s *DeployKeysService) UpdateDeployKey(pid any, deployKey int64, opt *UpdateDeployKeyOptions, options ...RequestOptionFunc) (*ProjectDeployKey, *Response, error) {
	return do[*ProjectDeployKey](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/deploy_keys/%d", ProjectID{pid}, deployKey),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
