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
)

type (
	ProtectedTagsServiceInterface interface {
		ListProtectedTags(pid any, opt *ListProtectedTagsOptions, options ...RequestOptionFunc) ([]*ProtectedTag, *Response, error)
		GetProtectedTag(pid any, tag string, options ...RequestOptionFunc) (*ProtectedTag, *Response, error)
		ProtectRepositoryTags(pid any, opt *ProtectRepositoryTagsOptions, options ...RequestOptionFunc) (*ProtectedTag, *Response, error)
		UnprotectRepositoryTags(pid any, tag string, options ...RequestOptionFunc) (*Response, error)
	}

	// ProtectedTagsService handles communication with the protected tag methods
	// of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/protected_tags/
	ProtectedTagsService struct {
		client *Client
	}
)

var _ ProtectedTagsServiceInterface = (*ProtectedTagsService)(nil)

// ProtectedTag represents a protected tag.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_tags/
type ProtectedTag struct {
	Name               string                  `json:"name"`
	CreateAccessLevels []*TagAccessDescription `json:"create_access_levels"`
}

// TagAccessDescription represents the access description for a protected tag.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_tags/
type TagAccessDescription struct {
	ID                     int64            `json:"id"`
	UserID                 int64            `json:"user_id"`
	GroupID                int64            `json:"group_id"`
	DeployKeyID            int64            `json:"deploy_key_id"`
	AccessLevel            AccessLevelValue `json:"access_level"`
	AccessLevelDescription string           `json:"access_level_description"`
}

// ListProtectedTagsOptions represents the available ListProtectedTags()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_tags/#list-protected-tags
type ListProtectedTagsOptions struct {
	ListOptions
}

// ListProtectedTags returns a list of protected tags from a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_tags/#list-protected-tags
func (s *ProtectedTagsService) ListProtectedTags(pid any, opt *ListProtectedTagsOptions, options ...RequestOptionFunc) ([]*ProtectedTag, *Response, error) {
	return do[[]*ProtectedTag](s.client,
		withPath("projects/%s/protected_tags", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetProtectedTag returns a single protected tag or wildcard protected tag.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_tags/#get-a-single-protected-tag-or-wildcard-protected-tag
func (s *ProtectedTagsService) GetProtectedTag(pid any, tag string, options ...RequestOptionFunc) (*ProtectedTag, *Response, error) {
	return do[*ProtectedTag](s.client,
		withPath("projects/%s/protected_tags/%s", ProjectID{pid}, tag),
		withRequestOpts(options...),
	)
}

// ProtectRepositoryTagsOptions represents the available ProtectRepositoryTags()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_tags/#protect-repository-tags
type ProtectRepositoryTagsOptions struct {
	Name              *string                   `url:"name,omitempty" json:"name,omitempty"`
	CreateAccessLevel *AccessLevelValue         `url:"create_access_level,omitempty" json:"create_access_level,omitempty"`
	AllowedToCreate   *[]*TagsPermissionOptions `url:"allowed_to_create,omitempty" json:"allowed_to_create,omitempty"`
}

// TagsPermissionOptions represents a protected tag permission option.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_tags/#protect-repository-tags
type TagsPermissionOptions struct {
	UserID      *int64            `url:"user_id,omitempty" json:"user_id,omitempty"`
	GroupID     *int64            `url:"group_id,omitempty" json:"group_id,omitempty"`
	DeployKeyID *int64            `url:"deploy_key_id,omitempty" json:"deploy_key_id,omitempty"`
	AccessLevel *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
}

// ProtectRepositoryTags protects a single repository tag or several project
// repository tags using a wildcard protected tag.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_tags/#protect-repository-tags
func (s *ProtectedTagsService) ProtectRepositoryTags(pid any, opt *ProtectRepositoryTagsOptions, options ...RequestOptionFunc) (*ProtectedTag, *Response, error) {
	return do[*ProtectedTag](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/protected_tags", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UnprotectRepositoryTags unprotects the given protected tag or wildcard
// protected tag.
//
// GitLab API docs:
// https://docs.gitlab.com/api/protected_tags/#unprotect-repository-tags
func (s *ProtectedTagsService) UnprotectRepositoryTags(pid any, tag string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/protected_tags/%s", ProjectID{pid}, tag),
		withRequestOpts(options...),
	)
	return resp, err
}
