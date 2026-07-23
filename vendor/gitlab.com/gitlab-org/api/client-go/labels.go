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
	"encoding/json"
	"net/http"
)

type (
	LabelsServiceInterface interface {
		ListLabels(pid any, opt *ListLabelsOptions, options ...RequestOptionFunc) ([]*Label, *Response, error)
		GetLabel(pid any, lid any, options ...RequestOptionFunc) (*Label, *Response, error)
		CreateLabel(pid any, opt *CreateLabelOptions, options ...RequestOptionFunc) (*Label, *Response, error)
		DeleteLabel(pid any, lid any, opt *DeleteLabelOptions, options ...RequestOptionFunc) (*Response, error)
		UpdateLabel(pid any, lid any, opt *UpdateLabelOptions, options ...RequestOptionFunc) (*Label, *Response, error)
		SubscribeToLabel(pid any, lid any, options ...RequestOptionFunc) (*Label, *Response, error)
		UnsubscribeFromLabel(pid any, lid any, options ...RequestOptionFunc) (*Response, error)
		PromoteLabel(pid any, lid any, options ...RequestOptionFunc) (*Response, error)
	}

	// LabelsService handles communication with the label related methods of the
	// GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/labels/
	LabelsService struct {
		client *Client
	}
)

var _ LabelsServiceInterface = (*LabelsService)(nil)

// Label represents a GitLab label.
//
// GitLab API docs: https://docs.gitlab.com/api/labels/
type Label struct {
	ID                     int64  `json:"id"`
	Name                   string `json:"name"`
	Color                  string `json:"color"`
	TextColor              string `json:"text_color"`
	Description            string `json:"description"`
	OpenIssuesCount        int64  `json:"open_issues_count"`
	ClosedIssuesCount      int64  `json:"closed_issues_count"`
	OpenMergeRequestsCount int64  `json:"open_merge_requests_count"`
	Subscribed             bool   `json:"subscribed"`
	Priority               int64  `json:"priority"`
	IsProjectLabel         bool   `json:"is_project_label"`
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (l *Label) UnmarshalJSON(data []byte) error {
	type alias Label
	if err := json.Unmarshal(data, (*alias)(l)); err != nil {
		return err
	}

	if l.Name == "" {
		var raw map[string]any
		if err := json.Unmarshal(data, &raw); err != nil {
			return err
		}
		if title, ok := raw["title"].(string); ok {
			l.Name = title
		}
	}

	return nil
}

func (l Label) String() string {
	return Stringify(l)
}

// ListLabelsOptions represents the available ListLabels() options.
//
// GitLab API docs: https://docs.gitlab.com/api/labels/#list-labels
type ListLabelsOptions struct {
	ListOptions
	WithCounts            *bool   `url:"with_counts,omitempty" json:"with_counts,omitempty"`
	IncludeAncestorGroups *bool   `url:"include_ancestor_groups,omitempty" json:"include_ancestor_groups,omitempty"`
	Search                *string `url:"search,omitempty" json:"search,omitempty"`
	Archived              *bool   `url:"archived,omitempty" json:"archived,omitempty"`
}

// ListLabels gets all labels for given project.
//
// GitLab API docs: https://docs.gitlab.com/api/labels/#list-labels
func (s *LabelsService) ListLabels(pid any, opt *ListLabelsOptions, options ...RequestOptionFunc) ([]*Label, *Response, error) {
	return do[[]*Label](s.client,
		withPath("projects/%s/labels", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetLabel get a single label for a given project.
//
// GitLab API docs: https://docs.gitlab.com/api/labels/#get-a-single-project-label
func (s *LabelsService) GetLabel(pid any, lid any, options ...RequestOptionFunc) (*Label, *Response, error) {
	return do[*Label](s.client,
		withPath("projects/%s/labels/%s", ProjectID{pid}, LabelID{lid}),
		withRequestOpts(options...),
	)
}

// CreateLabelOptions represents the available CreateLabel() options.
//
// GitLab API docs: https://docs.gitlab.com/api/labels/#create-a-project-label
type CreateLabelOptions struct {
	Name        *string `url:"name,omitempty" json:"name,omitempty"`
	Color       *string `url:"color,omitempty" json:"color,omitempty"`
	Description *string `url:"description,omitempty" json:"description,omitempty"`
	Priority    *int64  `url:"priority,omitempty" json:"priority,omitempty"`
	Archived    *bool   `url:"archived,omitempty" json:"archived,omitempty"`
}

// CreateLabel creates a new label for given repository with given name and
// color.
//
// GitLab API docs: https://docs.gitlab.com/api/labels/#create-a-new-label
func (s *LabelsService) CreateLabel(pid any, opt *CreateLabelOptions, options ...RequestOptionFunc) (*Label, *Response, error) {
	return do[*Label](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/labels", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteLabelOptions represents the available DeleteLabel() options.
//
// GitLab API docs: https://docs.gitlab.com/api/labels/#delete-a-label
type DeleteLabelOptions struct {
	Name *string `url:"name,omitempty" json:"name,omitempty"`
}

// DeleteLabel deletes a label given by its name or ID.
//
// GitLab API docs: https://docs.gitlab.com/api/labels/#delete-a-label
func (s *LabelsService) DeleteLabel(pid any, lid any, opt *DeleteLabelOptions, options ...RequestOptionFunc) (*Response, error) {
	reqOpts := make([]doOption, 0, 4)
	reqOpts = append(reqOpts,
		withMethod(http.MethodDelete),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)

	if lid != nil {
		reqOpts = append(reqOpts, withPath("projects/%s/labels/%s", ProjectID{pid}, LabelID{lid}))
	} else {
		reqOpts = append(reqOpts, withPath("projects/%s/labels", ProjectID{pid}))
	}

	_, resp, err := do[none](s.client, reqOpts...)
	return resp, err
}

// UpdateLabelOptions represents the available UpdateLabel() options.
//
// GitLab API docs: https://docs.gitlab.com/api/labels/#update-a-project-label
type UpdateLabelOptions struct {
	Name        *string `url:"name,omitempty" json:"name,omitempty"`
	NewName     *string `url:"new_name,omitempty" json:"new_name,omitempty"`
	Color       *string `url:"color,omitempty" json:"color,omitempty"`
	Description *string `url:"description,omitempty" json:"description,omitempty"`
	Priority    *int64  `url:"priority,omitempty" json:"priority,omitempty"`
	Archived    *bool   `url:"archived,omitempty" json:"archived,omitempty"`
}

// UpdateLabel updates an existing label with new name or now color. At least
// one parameter is required, to update the label.
//
// GitLab API docs: https://docs.gitlab.com/api/labels/#edit-an-existing-label
func (s *LabelsService) UpdateLabel(pid any, lid any, opt *UpdateLabelOptions, options ...RequestOptionFunc) (*Label, *Response, error) {
	reqOpts := make([]doOption, 0, 4)
	reqOpts = append(reqOpts,
		withMethod(http.MethodPut),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)

	if lid != nil {
		reqOpts = append(reqOpts, withPath("projects/%s/labels/%s", ProjectID{pid}, LabelID{lid}))
	} else {
		reqOpts = append(reqOpts, withPath("projects/%s/labels", ProjectID{pid}))
	}

	return do[*Label](s.client, reqOpts...)
}

// SubscribeToLabel subscribes the authenticated user to a label to receive
// notifications. If the user is already subscribed to the label, the status
// code 304 is returned.
//
// GitLab API docs:
// https://docs.gitlab.com/api/labels/#subscribe-to-a-label
func (s *LabelsService) SubscribeToLabel(pid any, lid any, options ...RequestOptionFunc) (*Label, *Response, error) {
	return do[*Label](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/labels/%s/subscribe", ProjectID{pid}, LabelID{lid}),
		withRequestOpts(options...),
	)
}

// UnsubscribeFromLabel unsubscribes the authenticated user from a label to not
// receive notifications from it. If the user is not subscribed to the label, the
// status code 304 is returned.
//
// GitLab API docs:
// https://docs.gitlab.com/api/labels/#unsubscribe-from-a-label
func (s *LabelsService) UnsubscribeFromLabel(pid any, lid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/labels/%s/unsubscribe", ProjectID{pid}, LabelID{lid}),
		withRequestOpts(options...),
	)
	return resp, err
}

// PromoteLabel Promotes a project label to a group label.
//
// GitLab API docs:
// https://docs.gitlab.com/api/labels/#promote-a-project-label-to-a-group-label
func (s *LabelsService) PromoteLabel(pid any, lid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/labels/%s/promote", ProjectID{pid}, LabelID{lid}),
		withRequestOpts(options...),
	)
	return resp, err
}
