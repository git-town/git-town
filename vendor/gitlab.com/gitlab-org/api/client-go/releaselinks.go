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

import "net/http"

type (
	ReleaseLinksServiceInterface interface {
		ListReleaseLinks(pid any, tagName string, opt *ListReleaseLinksOptions, options ...RequestOptionFunc) ([]*ReleaseLink, *Response, error)
		GetReleaseLink(pid any, tagName string, link int64, options ...RequestOptionFunc) (*ReleaseLink, *Response, error)
		CreateReleaseLink(pid any, tagName string, opt *CreateReleaseLinkOptions, options ...RequestOptionFunc) (*ReleaseLink, *Response, error)
		UpdateReleaseLink(pid any, tagName string, link int64, opt *UpdateReleaseLinkOptions, options ...RequestOptionFunc) (*ReleaseLink, *Response, error)
		DeleteReleaseLink(pid any, tagName string, link int64, options ...RequestOptionFunc) (*ReleaseLink, *Response, error)
	}

	// ReleaseLinksService handles communication with the release link methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/releases/links/
	ReleaseLinksService struct {
		client *Client
	}
)

var _ ReleaseLinksServiceInterface = (*ReleaseLinksService)(nil)

// ReleaseLink represents a release link.
//
// GitLab API docs: https://docs.gitlab.com/api/releases/links/
type ReleaseLink struct {
	ID             int64         `json:"id"`
	Name           string        `json:"name"`
	URL            string        `json:"url"`
	DirectAssetURL string        `json:"direct_asset_url"`
	External       bool          `json:"external"`
	LinkType       LinkTypeValue `json:"link_type"`
}

// ListReleaseLinksOptions represents ListReleaseLinks() options.
//
// GitLab API docs: https://docs.gitlab.com/api/releases/links/#list-links-of-a-release
type ListReleaseLinksOptions struct {
	ListOptions
}

// ListReleaseLinks gets assets as links from a Release.
//
// GitLab API docs: https://docs.gitlab.com/api/releases/links/#list-links-of-a-release
func (s *ReleaseLinksService) ListReleaseLinks(pid any, tagName string, opt *ListReleaseLinksOptions, options ...RequestOptionFunc) ([]*ReleaseLink, *Response, error) {
	return do[[]*ReleaseLink](s.client,
		withPath("projects/%s/releases/%s/assets/links", ProjectID{pid}, tagName),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetReleaseLink returns a link from release assets.
//
// GitLab API docs: https://docs.gitlab.com/api/releases/links/#get-a-release-link
func (s *ReleaseLinksService) GetReleaseLink(pid any, tagName string, link int64, options ...RequestOptionFunc) (*ReleaseLink, *Response, error) {
	return do[*ReleaseLink](s.client,
		withPath("projects/%s/releases/%s/assets/links/%d", ProjectID{pid}, tagName, link),
		withRequestOpts(options...),
	)
}

// CreateReleaseLinkOptions represents CreateReleaseLink() options.
//
// GitLab API docs: https://docs.gitlab.com/api/releases/links/#create-a-release-link
type CreateReleaseLinkOptions struct {
	Name            *string        `url:"name,omitempty" json:"name,omitempty"`
	URL             *string        `url:"url,omitempty" json:"url,omitempty"`
	FilePath        *string        `url:"filepath,omitempty" json:"filepath,omitempty"`
	DirectAssetPath *string        `url:"direct_asset_path,omitempty" json:"direct_asset_path,omitempty"`
	LinkType        *LinkTypeValue `url:"link_type,omitempty" json:"link_type,omitempty"`
}

// CreateReleaseLink creates a link.
//
// GitLab API docs: https://docs.gitlab.com/api/releases/links/#create-a-release-link
func (s *ReleaseLinksService) CreateReleaseLink(pid any, tagName string, opt *CreateReleaseLinkOptions, options ...RequestOptionFunc) (*ReleaseLink, *Response, error) {
	return do[*ReleaseLink](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/releases/%s/assets/links", ProjectID{pid}, tagName),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateReleaseLinkOptions represents UpdateReleaseLink() options.
//
// You have to specify at least one of Name of URL.
//
// GitLab API docs: https://docs.gitlab.com/api/releases/links/#update-a-release-link
type UpdateReleaseLinkOptions struct {
	Name            *string        `url:"name,omitempty" json:"name,omitempty"`
	URL             *string        `url:"url,omitempty" json:"url,omitempty"`
	FilePath        *string        `url:"filepath,omitempty" json:"filepath,omitempty"`
	DirectAssetPath *string        `url:"direct_asset_path,omitempty" json:"direct_asset_path,omitempty"`
	LinkType        *LinkTypeValue `url:"link_type,omitempty" json:"link_type,omitempty"`
}

// UpdateReleaseLink updates an asset link.
//
// GitLab API docs: https://docs.gitlab.com/api/releases/links/#update-a-release-link
func (s *ReleaseLinksService) UpdateReleaseLink(pid any, tagName string, link int64, opt *UpdateReleaseLinkOptions, options ...RequestOptionFunc) (*ReleaseLink, *Response, error) {
	return do[*ReleaseLink](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/releases/%s/assets/links/%d", ProjectID{pid}, tagName, link),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteReleaseLink deletes a link from release.
//
// GitLab API docs: https://docs.gitlab.com/api/releases/links/#delete-a-release-link
func (s *ReleaseLinksService) DeleteReleaseLink(pid any, tagName string, link int64, options ...RequestOptionFunc) (*ReleaseLink, *Response, error) {
	return do[*ReleaseLink](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/releases/%s/assets/links/%d", ProjectID{pid}, tagName, link),
		withRequestOpts(options...),
	)
}
