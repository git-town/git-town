//
// Copyright 2022, Daniel Steinke
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

import "time"

type (
	GroupIterationsServiceInterface interface {
		ListGroupIterations(gid any, opt *ListGroupIterationsOptions, options ...RequestOptionFunc) ([]*GroupIteration, *Response, error)
	}

	// GroupIterationsService handles communication with the iterations related methods
	// of the GitLab API
	//
	// GitLab API docs: https://docs.gitlab.com/api/group_iterations/
	GroupIterationsService struct {
		client *Client
	}
)

var _ GroupIterationsServiceInterface = (*GroupIterationsService)(nil)

// GroupIteration represents a GitLab iteration.
//
// GitLab API docs: https://docs.gitlab.com/api/group_iterations/
type GroupIteration struct {
	ID          int64      `json:"id"`
	IID         int64      `json:"iid"`
	Sequence    int64      `json:"sequence"`
	GroupID     int64      `json:"group_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	State       int64      `json:"state"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	DueDate     *ISOTime   `json:"due_date"`
	StartDate   *ISOTime   `json:"start_date"`
	WebURL      string     `json:"web_url"`
}

func (i GroupIteration) String() string {
	return Stringify(i)
}

// ListGroupIterationsOptions contains the available ListGroupIterations()
// options
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_iterations/#list-group-iterations
type ListGroupIterationsOptions struct {
	ListOptions
	State            *string `url:"state,omitempty" json:"state,omitempty"`
	Search           *string `url:"search,omitempty" json:"search,omitempty"`
	IncludeAncestors *bool   `url:"include_ancestors,omitempty" json:"include_ancestors,omitempty"`
}

// ListGroupIterations returns a list of group iterations.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_iterations/#list-group-iterations
func (s *GroupIterationsService) ListGroupIterations(gid any, opt *ListGroupIterationsOptions, options ...RequestOptionFunc) ([]*GroupIteration, *Response, error) {
	return do[[]*GroupIteration](s.client,
		withPath("groups/%s/iterations", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
