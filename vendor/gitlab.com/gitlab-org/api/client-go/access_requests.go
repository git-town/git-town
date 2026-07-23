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
	AccessRequestsServiceInterface interface {
		// ListProjectAccessRequests gets a list of access requests
		// viewable by the authenticated user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/access_requests/#list-access-requests-for-a-group-or-project
		// ListProjectAccessRequests gets a list of access requests
		// viewable by the authenticated user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/access_requests/#list-access-requests-for-a-group-or-project
		ListProjectAccessRequests(pid any, opt *ListAccessRequestsOptions, options ...RequestOptionFunc) ([]*AccessRequest, *Response, error)

		// ListGroupAccessRequests gets a list of access requests
		// viewable by the authenticated user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/access_requests/#list-access-requests-for-a-group-or-project

		// ListGroupAccessRequests gets a list of access requests
		// viewable by the authenticated user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/access_requests/#list-access-requests-for-a-group-or-project
		ListGroupAccessRequests(gid any, opt *ListAccessRequestsOptions, options ...RequestOptionFunc) ([]*AccessRequest, *Response, error)

		// RequestProjectAccess requests access for the authenticated user
		// to a group or project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/access_requests/#request-access-to-a-group-or-project

		// RequestProjectAccess requests access for the authenticated user
		// to a group or project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/access_requests/#request-access-to-a-group-or-project
		RequestProjectAccess(pid any, options ...RequestOptionFunc) (*AccessRequest, *Response, error)

		// RequestGroupAccess requests access for the authenticated user
		// to a group or project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/access_requests/#request-access-to-a-group-or-project
		RequestGroupAccess(gid any, options ...RequestOptionFunc) (*AccessRequest, *Response, error)

		// ApproveProjectAccessRequest approves an access request for the given user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/access_requests/#approve-an-access-request
		ApproveProjectAccessRequest(pid any, user int64, opt *ApproveAccessRequestOptions, options ...RequestOptionFunc) (*AccessRequest, *Response, error)

		// ApproveGroupAccessRequest approves an access request for the given user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/access_requests/#approve-an-access-request
		ApproveGroupAccessRequest(gid any, user int64, opt *ApproveAccessRequestOptions, options ...RequestOptionFunc) (*AccessRequest, *Response, error)

		// DenyProjectAccessRequest denies an access request for the given user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/access_requests/#deny-an-access-request
		DenyProjectAccessRequest(pid any, user int64, options ...RequestOptionFunc) (*Response, error)

		// DenyGroupAccessRequest denies an access request for the given user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/access_requests/#deny-an-access-request
		DenyGroupAccessRequest(gid any, user int64, options ...RequestOptionFunc) (*Response, error)
	}

	// AccessRequestsService handles communication with the project/group
	// access requests related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/access_requests/
	AccessRequestsService struct {
		client *Client
	}
)

var _ AccessRequestsServiceInterface = (*AccessRequestsService)(nil)

// AccessRequest represents a access request for a group or project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/access_requests/
type AccessRequest struct {
	ID          int64            `json:"id"`
	Username    string           `json:"username"`
	Name        string           `json:"name"`
	State       string           `json:"state"`
	CreatedAt   *time.Time       `json:"created_at"`
	RequestedAt *time.Time       `json:"requested_at"`
	AccessLevel AccessLevelValue `json:"access_level"`
}

// ListAccessRequestsOptions represents the available
// ListProjectAccessRequests() or ListGroupAccessRequests() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/access_requests/#list-access-requests-for-a-group-or-project
type ListAccessRequestsOptions struct {
	ListOptions
}

func (s *AccessRequestsService) ListProjectAccessRequests(pid any, opt *ListAccessRequestsOptions, options ...RequestOptionFunc) ([]*AccessRequest, *Response, error) {
	return do[[]*AccessRequest](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/access_requests", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *AccessRequestsService) ListGroupAccessRequests(gid any, opt *ListAccessRequestsOptions, options ...RequestOptionFunc) ([]*AccessRequest, *Response, error) {
	return do[[]*AccessRequest](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/access_requests", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *AccessRequestsService) RequestProjectAccess(pid any, options ...RequestOptionFunc) (*AccessRequest, *Response, error) {
	return do[*AccessRequest](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/access_requests", ProjectID{pid}),
		withRequestOpts(options...),
	)
}

func (s *AccessRequestsService) RequestGroupAccess(gid any, options ...RequestOptionFunc) (*AccessRequest, *Response, error) {
	return do[*AccessRequest](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/access_requests", GroupID{gid}),
		withRequestOpts(options...),
	)
}

// ApproveAccessRequestOptions represents the available
// ApproveProjectAccessRequest() and ApproveGroupAccessRequest() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/access_requests/#approve-an-access-request
type ApproveAccessRequestOptions struct {
	AccessLevel *AccessLevelValue `url:"access_level,omitempty" json:"access_level,omitempty"`
}

func (s *AccessRequestsService) ApproveProjectAccessRequest(pid any, user int64, opt *ApproveAccessRequestOptions, options ...RequestOptionFunc) (*AccessRequest, *Response, error) {
	return do[*AccessRequest](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/access_requests/%d/approve", ProjectID{pid}, user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *AccessRequestsService) ApproveGroupAccessRequest(gid any, user int64, opt *ApproveAccessRequestOptions, options ...RequestOptionFunc) (*AccessRequest, *Response, error) {
	return do[*AccessRequest](s.client,
		withMethod(http.MethodPut),
		withPath("groups/%s/access_requests/%d/approve", GroupID{gid}, user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *AccessRequestsService) DenyProjectAccessRequest(pid any, user int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/access_requests/%d", ProjectID{pid}, user),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *AccessRequestsService) DenyGroupAccessRequest(gid any, user int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/access_requests/%d", GroupID{gid}, user),
		withRequestOpts(options...),
	)
	return resp, err
}
