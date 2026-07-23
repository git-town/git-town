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

type (
	NamespacesServiceInterface interface {
		ListNamespaces(opt *ListNamespacesOptions, options ...RequestOptionFunc) ([]*Namespace, *Response, error)
		SearchNamespace(query string, options ...RequestOptionFunc) ([]*Namespace, *Response, error)
		GetNamespace(id any, options ...RequestOptionFunc) (*Namespace, *Response, error)
		NamespaceExists(id any, opt *NamespaceExistsOptions, options ...RequestOptionFunc) (*NamespaceExistance, *Response, error)
	}

	// NamespacesService handles communication with the namespace related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/namespaces/
	NamespacesService struct {
		client *Client
	}
)

var _ NamespacesServiceInterface = (*NamespacesService)(nil)

// Namespace represents a GitLab namespace.
//
// GitLab API docs: https://docs.gitlab.com/api/namespaces/
type Namespace struct {
	ID                          int64    `json:"id"`
	Name                        string   `json:"name"`
	Path                        string   `json:"path"`
	Kind                        string   `json:"kind"`
	FullPath                    string   `json:"full_path"`
	ParentID                    int64    `json:"parent_id"`
	AvatarURL                   *string  `json:"avatar_url"`
	WebURL                      string   `json:"web_url"`
	MembersCountWithDescendants int64    `json:"members_count_with_descendants"`
	BillableMembersCount        int64    `json:"billable_members_count"`
	Plan                        string   `json:"plan"`
	TrialEndsOn                 *ISOTime `json:"trial_ends_on"`
	Trial                       bool     `json:"trial"`
	MaxSeatsUsed                *int64   `json:"max_seats_used"`
	SeatsInUse                  *int64   `json:"seats_in_use"`
}

func (n Namespace) String() string {
	return Stringify(n)
}

// ListNamespacesOptions represents the available ListNamespaces() options.
//
// GitLab API docs: https://docs.gitlab.com/api/namespaces/#list-all-namespaces
type ListNamespacesOptions struct {
	ListOptions
	Search       *string `url:"search,omitempty" json:"search,omitempty"`
	OwnedOnly    *bool   `url:"owned_only,omitempty" json:"owned_only,omitempty"`
	TopLevelOnly *bool   `url:"top_level_only,omitempty" json:"top_level_only,omitempty"`
}

// ListNamespaces gets a list of projects accessible by the authenticated user.
//
// GitLab API docs: https://docs.gitlab.com/api/namespaces/#list-all-namespaces
func (s *NamespacesService) ListNamespaces(opt *ListNamespacesOptions, options ...RequestOptionFunc) ([]*Namespace, *Response, error) {
	return do[[]*Namespace](s.client,
		withPath("namespaces"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// SearchNamespace gets all namespaces that match your string in their name
// or path.
//
// GitLab API docs:
// https://docs.gitlab.com/api/namespaces/#list-all-namespaces
func (s *NamespacesService) SearchNamespace(query string, options ...RequestOptionFunc) ([]*Namespace, *Response, error) {
	var q struct {
		Search string `url:"search,omitempty" json:"search,omitempty"`
	}
	q.Search = query

	return do[[]*Namespace](s.client,
		withPath("namespaces"),
		withAPIOpts(&q),
		withRequestOpts(options...),
	)
}

// GetNamespace gets a namespace by id.
//
// GitLab API docs:
// https://docs.gitlab.com/api/namespaces/#get-details-on-a-namespace
func (s *NamespacesService) GetNamespace(id any, options ...RequestOptionFunc) (*Namespace, *Response, error) {
	namespace, err := parseID(id)
	if err != nil {
		return nil, nil, err
	}
	return do[*Namespace](s.client,
		withPath("namespaces/%s", namespace),
		withRequestOpts(options...),
	)
}

// NamespaceExistance represents a namespace exists result.
//
// GitLab API docs:
// https://docs.gitlab.com/api/namespaces/#verify-namespace-availability
type NamespaceExistance struct {
	Exists   bool     `json:"exists"`
	Suggests []string `json:"suggests"`
}

// NamespaceExistsOptions represents the available NamespaceExists() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/namespaces/#verify-namespace-availability
type NamespaceExistsOptions struct {
	ParentID *int64 `url:"parent_id,omitempty" json:"parent_id,omitempty"`
}

// NamespaceExists checks the existence of a namespace.
//
// GitLab API docs:
// https://docs.gitlab.com/api/namespaces/#verify-namespace-availability
func (s *NamespacesService) NamespaceExists(id any, opt *NamespaceExistsOptions, options ...RequestOptionFunc) (*NamespaceExistance, *Response, error) {
	namespace, err := parseID(id)
	if err != nil {
		return nil, nil, err
	}
	return do[*NamespaceExistance](s.client,
		withPath("namespaces/%s/exists", namespace),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}
