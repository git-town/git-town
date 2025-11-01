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
)

type (
	DependencyProxyServiceInterface interface {
		// PurgeGroupDependencyProxy schedules for deletion the cached manifests and blobs
		// for a group. This endpoint requires the Owner role for the group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/dependency_proxy/#purge-the-dependency-proxy-for-a-group
		PurgeGroupDependencyProxy(gid any, options ...RequestOptionFunc) (*Response, error)
	}

	// DependencyProxyService handles communication with the dependency proxy
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/dependency_proxy/
	DependencyProxyService struct {
		client *Client
	}
)

var _ DependencyProxyServiceInterface = (*DependencyProxyService)(nil)

func (s *DependencyProxyService) PurgeGroupDependencyProxy(gid any, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/dependency_proxy/cache", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
