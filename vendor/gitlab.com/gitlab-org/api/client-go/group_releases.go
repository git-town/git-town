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
)

type (
	GroupReleasesServiceInterface interface {
		ListGroupReleases(gid interface{}, opts *ListGroupReleasesOptions, options ...RequestOptionFunc) ([]*Release, *Response, error)
	}

	// GroupReleasesService handles communication with the group
	// releases related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/group_releases.html
	GroupReleasesService struct {
		client *Client
	}
)

var _ GroupReleasesServiceInterface = (*GroupReleasesService)(nil)

// ListGroupReleasesOptions represents the available ListGroupReleases() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_releases.html#list-group-releases
type ListGroupReleasesOptions struct {
	ListOptions
	Simple *bool `url:"simple,omitempty" json:"simple,omitempty"`
}

// ListGroupReleases gets a list of releases for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_releases.html#list-group-releases
func (s *GroupReleasesService) ListGroupReleases(gid interface{}, opts *ListGroupReleasesOptions, options ...RequestOptionFunc) ([]*Release, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/releases", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, opts, options)
	if err != nil {
		return nil, nil, err
	}

	var releases []*Release
	resp, err := s.client.Do(req, &releases)
	if err != nil {
		return nil, resp, err
	}
	return releases, resp, nil
}
