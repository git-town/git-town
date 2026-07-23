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
	// FeaturesServiceInterface defines all the API methods for the FeaturesService
	FeaturesServiceInterface interface {
		ListFeatures(options ...RequestOptionFunc) ([]*Feature, *Response, error)
		ListFeatureDefinitions(options ...RequestOptionFunc) ([]*FeatureDefinition, *Response, error)
		SetFeatureFlag(name string, opt *SetFeatureFlagOptions, options ...RequestOptionFunc) (*Feature, *Response, error)
		DeleteFeatureFlag(name string, options ...RequestOptionFunc) (*Response, error)
	}

	// FeaturesService handles the communication with the application FeaturesService
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/features/
	FeaturesService struct {
		client *Client
	}
)

var _ FeaturesServiceInterface = (*FeaturesService)(nil)

// Feature represents a GitLab feature flag.
//
// GitLab API docs: https://docs.gitlab.com/api/features/
type Feature struct {
	Name       string `json:"name"`
	State      string `json:"state"`
	Gates      []Gate
	Definition *FeatureDefinition `json:"definition"`
}

// Gate represents a gate of a GitLab feature flag.
//
// GitLab API docs: https://docs.gitlab.com/api/features/
type Gate struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

func (f Feature) String() string {
	return Stringify(f)
}

// ListFeatures gets a list of feature flags
//
// GitLab API docs:
// https://docs.gitlab.com/api/features/#list-all-features
func (s *FeaturesService) ListFeatures(options ...RequestOptionFunc) ([]*Feature, *Response, error) {
	return do[[]*Feature](s.client,
		withPath("features"),
		withRequestOpts(options...),
	)
}

// FeatureDefinition represents a Feature Definition.
//
// GitLab API docs:
// https://docs.gitlab.com/api/features/#list-all-feature-definitions
type FeatureDefinition struct {
	Name            string `json:"name"`
	IntroducedByURL string `json:"introduced_by_url"`
	RolloutIssueURL string `json:"rollout_issue_url"`
	Milestone       string `json:"milestone"`
	LogStateChanges bool   `json:"log_state_changes"`
	Type            string `json:"type"`
	Group           string `json:"group"`
	DefaultEnabled  bool   `json:"default_enabled"`
}

func (fd FeatureDefinition) String() string {
	return Stringify(fd)
}

// ListFeatureDefinitions gets a lists of all feature definitions.
//
// GitLab API docs:
// https://docs.gitlab.com/api/features/#list-all-feature-definitions
func (s *FeaturesService) ListFeatureDefinitions(options ...RequestOptionFunc) ([]*FeatureDefinition, *Response, error) {
	return do[[]*FeatureDefinition](s.client,
		withPath("features/definitions"),
		withRequestOpts(options...),
	)
}

// SetFeatureFlagOptions represents the available options for
// SetFeatureFlag().
//
// GitLab API docs:
// https://docs.gitlab.com/api/features/#set-or-create-a-feature
type SetFeatureFlagOptions struct {
	Value        any    `url:"value" json:"value"`
	Key          string `url:"key" json:"key"`
	FeatureGroup string `url:"feature_group" json:"feature_group"`
	User         string `url:"user" json:"user"`
	Group        string `url:"group" json:"group"`
	Namespace    string `url:"namespace" json:"namespace"`
	Project      string `url:"project" json:"project"`
	Repository   string `url:"repository" json:"repository"`
	Force        bool   `url:"force" json:"force"`
}

// SetFeatureFlag sets or creates a feature flag gate
//
// GitLab API docs:
// https://docs.gitlab.com/api/features/#set-or-create-a-feature
func (s *FeaturesService) SetFeatureFlag(name string, opt *SetFeatureFlagOptions, options ...RequestOptionFunc) (*Feature, *Response, error) {
	return do[*Feature](s.client,
		withMethod(http.MethodPost),
		withPath("features/%s", name),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteFeatureFlag deletes a feature flag.
//
// GitLab API docs:
// https://docs.gitlab.com/api/features/#delete-a-feature
func (s *FeaturesService) DeleteFeatureFlag(name string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("features/%s", name),
		withRequestOpts(options...),
	)
	return resp, err
}
