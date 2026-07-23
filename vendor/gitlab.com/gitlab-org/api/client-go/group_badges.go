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
	// GroupBadgesServiceInterface defines all the API methods for the GroupBadgesService
	GroupBadgesServiceInterface interface {
		ListGroupBadges(gid any, opt *ListGroupBadgesOptions, options ...RequestOptionFunc) ([]*GroupBadge, *Response, error)
		GetGroupBadge(gid any, badge int64, options ...RequestOptionFunc) (*GroupBadge, *Response, error)
		AddGroupBadge(gid any, opt *AddGroupBadgeOptions, options ...RequestOptionFunc) (*GroupBadge, *Response, error)
		EditGroupBadge(gid any, badge int64, opt *EditGroupBadgeOptions, options ...RequestOptionFunc) (*GroupBadge, *Response, error)
		DeleteGroupBadge(gid any, badge int64, options ...RequestOptionFunc) (*Response, error)
		PreviewGroupBadge(gid any, opt *GroupBadgePreviewOptions, options ...RequestOptionFunc) (*GroupBadge, *Response, error)
	}

	// GroupBadgesService handles communication with the group badges
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/group_badges/
	GroupBadgesService struct {
		client *Client
	}
)

var _ GroupBadgesServiceInterface = (*GroupBadgesService)(nil)

// BadgeKind represents a GitLab Badge Kind
type BadgeKind string

// all possible values Badge Kind
const (
	ProjectBadgeKind BadgeKind = "project"
	GroupBadgeKind   BadgeKind = "group"
)

// GroupBadge represents a group badge.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_badges/
type GroupBadge struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	LinkURL          string    `json:"link_url"`
	ImageURL         string    `json:"image_url"`
	RenderedLinkURL  string    `json:"rendered_link_url"`
	RenderedImageURL string    `json:"rendered_image_url"`
	Kind             BadgeKind `json:"kind"`
}

// ListGroupBadgesOptions represents the available ListGroupBadges() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_badges/#list-all-badges-of-a-group
type ListGroupBadgesOptions struct {
	ListOptions
	Name *string `url:"name,omitempty" json:"name,omitempty"`
}

// ListGroupBadges gets a list of a group badges.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_badges/#list-all-badges-of-a-group
func (s *GroupBadgesService) ListGroupBadges(gid any, opt *ListGroupBadgesOptions, options ...RequestOptionFunc) ([]*GroupBadge, *Response, error) {
	return do[[]*GroupBadge](s.client,
		withPath("groups/%s/badges", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetGroupBadge gets a group badge.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_badges/#get-a-badge-of-a-group
func (s *GroupBadgesService) GetGroupBadge(gid any, badge int64, options ...RequestOptionFunc) (*GroupBadge, *Response, error) {
	return do[*GroupBadge](s.client,
		withPath("groups/%s/badges/%d", GroupID{gid}, badge),
		withRequestOpts(options...),
	)
}

// AddGroupBadgeOptions represents the available AddGroupBadge() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_badges/#add-a-badge-to-a-group
type AddGroupBadgeOptions struct {
	LinkURL  *string `url:"link_url,omitempty" json:"link_url,omitempty"`
	ImageURL *string `url:"image_url,omitempty" json:"image_url,omitempty"`
	Name     *string `url:"name,omitempty" json:"name,omitempty"`
}

// AddGroupBadge adds a badge to a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_badges/#add-a-badge-to-a-group
func (s *GroupBadgesService) AddGroupBadge(gid any, opt *AddGroupBadgeOptions, options ...RequestOptionFunc) (*GroupBadge, *Response, error) {
	return do[*GroupBadge](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/badges", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// EditGroupBadgeOptions represents the available EditGroupBadge() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_badges/#edit-a-badge-of-a-group
type EditGroupBadgeOptions struct {
	LinkURL  *string `url:"link_url,omitempty" json:"link_url,omitempty"`
	ImageURL *string `url:"image_url,omitempty" json:"image_url,omitempty"`
	Name     *string `url:"name,omitempty" json:"name,omitempty"`
}

// EditGroupBadge updates a badge of a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_badges/#edit-a-badge-of-a-group
func (s *GroupBadgesService) EditGroupBadge(gid any, badge int64, opt *EditGroupBadgeOptions, options ...RequestOptionFunc) (*GroupBadge, *Response, error) {
	return do[*GroupBadge](s.client,
		withMethod(http.MethodPut),
		withPath("groups/%s/badges/%d", GroupID{gid}, badge),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteGroupBadge removes a badge from a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_badges/#remove-a-badge-from-a-group
func (s *GroupBadgesService) DeleteGroupBadge(gid any, badge int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/badges/%d", GroupID{gid}, badge),
		withRequestOpts(options...),
	)
	return resp, err
}

// GroupBadgePreviewOptions represents the available PreviewGroupBadge() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_badges/#preview-a-badge-from-a-group
type GroupBadgePreviewOptions struct {
	LinkURL  *string `url:"link_url,omitempty" json:"link_url,omitempty"`
	ImageURL *string `url:"image_url,omitempty" json:"image_url,omitempty"`
	Name     *string `url:"name,omitempty" json:"name,omitempty"`
}

// PreviewGroupBadge returns how the link_url and image_url final URLs would be after
// resolving the placeholder interpolation.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_badges/#preview-a-badge-from-a-group
func (s *GroupBadgesService) PreviewGroupBadge(gid any, opt *GroupBadgePreviewOptions, options ...RequestOptionFunc) (*GroupBadge, *Response, error) {
	return do[*GroupBadge](s.client,
		withPath("groups/%s/badges/render", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
