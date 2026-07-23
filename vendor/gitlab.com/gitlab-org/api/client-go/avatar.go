//
// Copyright 2021, Pavel Kostohrys
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
	AvatarRequestsServiceInterface interface {
		// GetAvatar gets the avatar URL for a user with the given email address.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/avatar/#get-details-on-an-account-avatar
		GetAvatar(opt *GetAvatarOptions, options ...RequestOptionFunc) (*Avatar, *Response, error)
	}

	// AvatarRequestsService handles communication with the avatar related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/avatar/
	AvatarRequestsService struct {
		client *Client
	}
)

var _ AvatarRequestsServiceInterface = (*AvatarRequestsService)(nil)

// Avatar represents a GitLab avatar.
//
// GitLab API docs: https://docs.gitlab.com/api/avatar/
type Avatar struct {
	AvatarURL string `json:"avatar_url"`
}

// GetAvatarOptions represents the available GetAvatar() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/avatar/#get-details-on-an-account-avatar
type GetAvatarOptions struct {
	Email *string `url:"email,omitempty" json:"email,omitempty"`
	Size  *int64  `url:"size,omitempty" json:"size,omitempty"`
}

func (s *AvatarRequestsService) GetAvatar(opt *GetAvatarOptions, options ...RequestOptionFunc) (*Avatar, *Response, error) {
	return do[*Avatar](s.client,
		withMethod(http.MethodGet),
		withPath("avatar"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
