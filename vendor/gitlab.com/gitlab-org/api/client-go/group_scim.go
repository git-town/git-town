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
	GroupSCIMServiceInterface interface {
		GetSCIMIdentitiesForGroup(gid any, options ...RequestOptionFunc) ([]*GroupSCIMIdentity, *Response, error)
		GetSCIMIdentity(gid any, uid string, options ...RequestOptionFunc) (*GroupSCIMIdentity, *Response, error)
		UpdateSCIMIdentity(gid any, uid string, opt *UpdateSCIMIdentityOptions, options ...RequestOptionFunc) (*Response, error)
		DeleteSCIMIdentity(gid any, uid string, options ...RequestOptionFunc) (*Response, error)
	}

	// GroupSCIMService handles communication with the Group SCIM
	// related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/scim/
	GroupSCIMService struct {
		client *Client
	}
)

// GroupSCIMIdentity represents a GitLab Group SCIM identity.
//
// GitLab API docs:
// https://docs.gitlab.com/api/scim/
type GroupSCIMIdentity struct {
	ExternalUID string `json:"external_uid"`
	UserID      int64  `json:"user_id"`
	Active      bool   `json:"active"`
}

// GetSCIMIdentitiesForGroup gets all SCIM identities for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/scim/#get-scim-identities-for-a-group
func (s *GroupSCIMService) GetSCIMIdentitiesForGroup(gid any, options ...RequestOptionFunc) ([]*GroupSCIMIdentity, *Response, error) {
	return do[[]*GroupSCIMIdentity](s.client,
		withPath("groups/%s/scim/identities", GroupID{gid}),
		withRequestOpts(options...),
	)
}

// GetSCIMIdentity gets a SCIM identity for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/scim/#get-a-single-scim-identity
func (s *GroupSCIMService) GetSCIMIdentity(gid any, uid string, options ...RequestOptionFunc) (*GroupSCIMIdentity, *Response, error) {
	return do[*GroupSCIMIdentity](s.client,
		withPath("groups/%s/scim/%s", GroupID{gid}, uid),
		withRequestOpts(options...),
	)
}

// UpdateSCIMIdentityOptions represent the request options for
// updating a SCIM Identity.
//
// GitLab API docs:
// https://docs.gitlab.com/api/scim/#update-extern_uid-field-for-a-scim-identity
type UpdateSCIMIdentityOptions struct {
	ExternUID *string `url:"extern_uid,omitempty" json:"extern_uid,omitempty"`
}

// UpdateSCIMIdentity updates a SCIM identity.
//
// GitLab API docs:
// https://docs.gitlab.com/api/scim/#update-extern_uid-field-for-a-scim-identity
func (s *GroupSCIMService) UpdateSCIMIdentity(gid any, uid string, opt *UpdateSCIMIdentityOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPatch),
		withPath("groups/%s/scim/%s", GroupID{gid}, uid),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

// DeleteSCIMIdentity deletes a SCIM identity.
//
// GitLab API docs:
// https://docs.gitlab.com/api/scim/#delete-a-single-scim-identity
func (s *GroupSCIMService) DeleteSCIMIdentity(gid any, uid string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/scim/%s", GroupID{gid}, uid),
		withRequestOpts(options...),
	)
	return resp, err
}
