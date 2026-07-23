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
	ContainerRegistryServiceInterface interface {
		// ListProjectRegistryRepositories gets a list of registry repositories in a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_registry/#within-a-project
		ListProjectRegistryRepositories(pid any, opt *ListProjectRegistryRepositoriesOptions, options ...RequestOptionFunc) ([]*RegistryRepository, *Response, error)

		// ListGroupRegistryRepositories gets a list of registry repositories in a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_registry/#within-a-group
		ListGroupRegistryRepositories(gid any, opt *ListGroupRegistryRepositoriesOptions, options ...RequestOptionFunc) ([]*RegistryRepository, *Response, error)

		// GetSingleRegistryRepository gets the details of single registry repository.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_registry/#get-details-of-a-single-repository
		GetSingleRegistryRepository(pid any, opt *GetSingleRegistryRepositoryOptions, options ...RequestOptionFunc) (*RegistryRepository, *Response, error)

		// DeleteRegistryRepository deletes a repository in a registry.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_registry/#delete-registry-repository
		DeleteRegistryRepository(pid any, repository int64, options ...RequestOptionFunc) (*Response, error)

		// ListRegistryRepositoryTags gets a list of tags for given registry repository.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_registry/#list-registry-repository-tags
		ListRegistryRepositoryTags(pid any, repository int64, opt *ListRegistryRepositoryTagsOptions, options ...RequestOptionFunc) ([]*RegistryRepositoryTag, *Response, error)

		// GetRegistryRepositoryTagDetail get details of a registry repository tag.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_registry/#get-details-of-a-registry-repository-tag
		GetRegistryRepositoryTagDetail(pid any, repository int64, tagName string, options ...RequestOptionFunc) (*RegistryRepositoryTag, *Response, error)

		// DeleteRegistryRepositoryTag deletes a registry repository tag.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_registry/#delete-a-registry-repository-tag
		DeleteRegistryRepositoryTag(pid any, repository int64, tagName string, options ...RequestOptionFunc) (*Response, error)

		// DeleteRegistryRepositoryTags deletes repository tags in bulk based on given criteria.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/container_registry/#delete-registry-repository-tags-in-bulk
		DeleteRegistryRepositoryTags(pid any, repository int64, opt *DeleteRegistryRepositoryTagsOptions, options ...RequestOptionFunc) (*Response, error)
	}

	// ContainerRegistryService handles communication with the container registry
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/container_registry/
	ContainerRegistryService struct {
		client *Client
	}
)

var _ ContainerRegistryServiceInterface = (*ContainerRegistryService)(nil)

// RegistryRepository represents a GitLab content registry repository.
//
// GitLab API docs: https://docs.gitlab.com/api/container_registry/
type RegistryRepository struct {
	ID                     int64                    `json:"id"`
	Name                   string                   `json:"name"`
	Path                   string                   `json:"path"`
	ProjectID              int64                    `json:"project_id"`
	Location               string                   `json:"location"`
	CreatedAt              *time.Time               `json:"created_at"`
	CleanupPolicyStartedAt *time.Time               `json:"cleanup_policy_started_at"`
	Status                 *ContainerRegistryStatus `json:"status"`
	TagsCount              int64                    `json:"tags_count"`
	Tags                   []*RegistryRepositoryTag `json:"tags"`
}

func (s RegistryRepository) String() string {
	return Stringify(s)
}

// RegistryRepositoryTag represents a GitLab registry image tag.
//
// GitLab API docs: https://docs.gitlab.com/api/container_registry/
type RegistryRepositoryTag struct {
	Name          string     `json:"name"`
	Path          string     `json:"path"`
	Location      string     `json:"location"`
	Revision      string     `json:"revision"`
	ShortRevision string     `json:"short_revision"`
	Digest        string     `json:"digest"`
	CreatedAt     *time.Time `json:"created_at"`
	TotalSize     int64      `json:"total_size"`
}

func (s RegistryRepositoryTag) String() string {
	return Stringify(s)
}

// ListProjectRegistryRepositoriesOptions represents the available
// ListProjectRegistryRepositories() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/container_registry/#list-registry-repositories
type ListProjectRegistryRepositoriesOptions struct {
	ListOptions

	Tags      *bool `url:"tags,omitempty" json:"tags,omitempty"`
	TagsCount *bool `url:"tags_count,omitempty" json:"tags_count,omitempty"`
}

// ListGroupRegistryRepositoriesOptions represents the available
// ListGroupRegistryRepositories() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/container_registry/#within-a-group
type ListGroupRegistryRepositoriesOptions struct {
	ListOptions
}

func (s *ContainerRegistryService) ListProjectRegistryRepositories(pid any, opt *ListProjectRegistryRepositoriesOptions, options ...RequestOptionFunc) ([]*RegistryRepository, *Response, error) {
	return do[[]*RegistryRepository](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/registry/repositories", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ContainerRegistryService) ListGroupRegistryRepositories(gid any, opt *ListGroupRegistryRepositoriesOptions, options ...RequestOptionFunc) ([]*RegistryRepository, *Response, error) {
	return do[[]*RegistryRepository](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/registry/repositories", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetSingleRegistryRepositoryOptions represents the available
// GetSingleRegistryRepository() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/container_registry/#get-details-of-a-single-repository
type GetSingleRegistryRepositoryOptions struct {
	Tags      *bool `url:"tags,omitempty" json:"tags,omitempty"`
	TagsCount *bool `url:"tags_count,omitempty" json:"tags_count,omitempty"`
}

func (s *ContainerRegistryService) GetSingleRegistryRepository(pid any, opt *GetSingleRegistryRepositoryOptions, options ...RequestOptionFunc) (*RegistryRepository, *Response, error) {
	return do[*RegistryRepository](s.client,
		withMethod(http.MethodGet),
		withPath("registry/repositories/%s", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ContainerRegistryService) DeleteRegistryRepository(pid any, repository int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/registry/repositories/%d", ProjectID{pid}, repository),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListRegistryRepositoryTagsOptions represents the available
// ListRegistryRepositoryTags() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/container_registry/#list-registry-repository-tags
type ListRegistryRepositoryTagsOptions struct {
	ListOptions
}

func (s *ContainerRegistryService) ListRegistryRepositoryTags(pid any, repository int64, opt *ListRegistryRepositoryTagsOptions, options ...RequestOptionFunc) ([]*RegistryRepositoryTag, *Response, error) {
	return do[[]*RegistryRepositoryTag](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/registry/repositories/%d/tags", ProjectID{pid}, repository),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ContainerRegistryService) GetRegistryRepositoryTagDetail(pid any, repository int64, tagName string, options ...RequestOptionFunc) (*RegistryRepositoryTag, *Response, error) {
	return do[*RegistryRepositoryTag](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/registry/repositories/%d/tags/%s", ProjectID{pid}, repository, tagName),
		withRequestOpts(options...),
	)
}

func (s *ContainerRegistryService) DeleteRegistryRepositoryTag(pid any, repository int64, tagName string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/registry/repositories/%d/tags/%s", ProjectID{pid}, repository, tagName),
		withRequestOpts(options...),
	)
	return resp, err
}

// DeleteRegistryRepositoryTagsOptions represents the available
// DeleteRegistryRepositoryTags() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/container_registry/#delete-registry-repository-tags-in-bulk
type DeleteRegistryRepositoryTagsOptions struct {
	NameRegexpDelete *string `url:"name_regex_delete,omitempty" json:"name_regex_delete,omitempty"`
	NameRegexpKeep   *string `url:"name_regex_keep,omitempty" json:"name_regex_keep,omitempty"`
	KeepN            *int64  `url:"keep_n,omitempty" json:"keep_n,omitempty"`
	OlderThan        *string `url:"older_than,omitempty" json:"older_than,omitempty"`

	// Deprecated: NameRegexp is deprecated in favor of NameRegexpDelete.
	NameRegexp *string `url:"name_regex,omitempty" json:"name_regex,omitempty"`
}

func (s *ContainerRegistryService) DeleteRegistryRepositoryTags(pid any, repository int64, opt *DeleteRegistryRepositoryTagsOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/registry/repositories/%d/tags", ProjectID{pid}, repository),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}
