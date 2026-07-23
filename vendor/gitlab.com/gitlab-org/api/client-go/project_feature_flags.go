package gitlab

import (
	"net/http"
	"time"
)

type (
	ProjectFeatureFlagServiceInterface interface {
		// ListProjectFeatureFlags returns a list with the feature flags of a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/feature_flags/#list-feature-flags-for-a-project
		ListProjectFeatureFlags(pid any, opt *ListProjectFeatureFlagOptions, options ...RequestOptionFunc) ([]*ProjectFeatureFlag, *Response, error)
		// GetProjectFeatureFlag gets a single feature flag for the specified project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/feature_flags/#get-a-single-feature-flag
		GetProjectFeatureFlag(pid any, name string, options ...RequestOptionFunc) (*ProjectFeatureFlag, *Response, error)
		// CreateProjectFeatureFlag creates a feature flag.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/feature_flags/#create-a-feature-flag
		CreateProjectFeatureFlag(pid any, opt *CreateProjectFeatureFlagOptions, options ...RequestOptionFunc) (*ProjectFeatureFlag, *Response, error)
		// UpdateProjectFeatureFlag updates a feature flag.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/feature_flags/#update-a-feature-flag
		UpdateProjectFeatureFlag(pid any, name string, opt *UpdateProjectFeatureFlagOptions, options ...RequestOptionFunc) (*ProjectFeatureFlag, *Response, error)
		// DeleteProjectFeatureFlag deletes a feature flag.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/feature_flags/#delete-a-feature-flag
		DeleteProjectFeatureFlag(pid any, name string, options ...RequestOptionFunc) (*Response, error)
	}

	// ProjectFeatureFlagService handles operations on gitlab project feature
	// flags using the following api:
	//
	// GitLab API docs: https://docs.gitlab.com/api/feature_flags/
	ProjectFeatureFlagService struct {
		client *Client
	}
)

var _ ProjectFeatureFlagServiceInterface = (*ProjectFeatureFlagService)(nil)

// ProjectFeatureFlag represents a GitLab project iteration.
//
// GitLab API docs: https://docs.gitlab.com/api/feature_flags/
type ProjectFeatureFlag struct {
	Name        string                        `json:"name"`
	Description string                        `json:"description"`
	Active      bool                          `json:"active"`
	Version     string                        `json:"version"`
	CreatedAt   *time.Time                    `json:"created_at"`
	UpdatedAt   *time.Time                    `json:"updated_at"`
	Scopes      []*ProjectFeatureFlagScope    `json:"scopes"`
	Strategies  []*ProjectFeatureFlagStrategy `json:"strategies"`
}

// ProjectFeatureFlagScope defines the scopes of a feature flag
//
// GitLab API docs: https://docs.gitlab.com/api/feature_flags/
type ProjectFeatureFlagScope struct {
	ID               int64  `json:"id"`
	EnvironmentScope string `json:"environment_scope"`
}

// ProjectFeatureFlagStrategy defines the strategy used for a feature flag
//
// GitLab API docs: https://docs.gitlab.com/api/feature_flags/
type ProjectFeatureFlagStrategy struct {
	ID         int64                                `json:"id"`
	Name       string                               `json:"name"`
	Parameters *ProjectFeatureFlagStrategyParameter `json:"parameters"`
	Scopes     []*ProjectFeatureFlagScope           `json:"scopes"`
}

// ProjectFeatureFlagStrategyParameter is used in updating and creating feature flags
//
// GitLab API docs: https://docs.gitlab.com/api/feature_flags/
type ProjectFeatureFlagStrategyParameter struct {
	GroupID    string `json:"groupId,omitempty"`
	UserIDs    string `json:"userIds,omitempty"`
	Percentage string `json:"percentage,omitempty"`

	// Following fields aren't documented in GitLab API docs,
	// but are present in GitLab API since 13.5.
	// Docs: https://docs.getunleash.io/reference/activation-strategies#gradual-rollout
	Rollout    string `json:"rollout,omitempty"`
	Stickiness string `json:"stickiness,omitempty"`
}

func (i ProjectFeatureFlag) String() string {
	return Stringify(i)
}

// ListProjectFeatureFlagOptions contains the options for ListProjectFeatureFlags
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flags/#list-feature-flags-for-a-project
type ListProjectFeatureFlagOptions struct {
	ListOptions
	Scope *string `url:"scope,omitempty" json:"scope,omitempty"`
}

func (s *ProjectFeatureFlagService) ListProjectFeatureFlags(pid any, opt *ListProjectFeatureFlagOptions, options ...RequestOptionFunc) ([]*ProjectFeatureFlag, *Response, error) {
	return do[[]*ProjectFeatureFlag](s.client,
		withPath("projects/%s/feature_flags", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ProjectFeatureFlagService) GetProjectFeatureFlag(pid any, name string, options ...RequestOptionFunc) (*ProjectFeatureFlag, *Response, error) {
	return do[*ProjectFeatureFlag](s.client,
		withPath("projects/%s/feature_flags/%s", ProjectID{pid}, name),
		withRequestOpts(options...),
	)
}

// CreateProjectFeatureFlagOptions represents the available
// CreateProjectFeatureFlag() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flags/#create-a-feature-flag
type CreateProjectFeatureFlagOptions struct {
	Name        *string                        `url:"name,omitempty" json:"name,omitempty"`
	Description *string                        `url:"description,omitempty" json:"description,omitempty"`
	Version     *string                        `url:"version,omitempty" json:"version,omitempty"`
	Active      *bool                          `url:"active,omitempty" json:"active,omitempty"`
	Strategies  *[]*FeatureFlagStrategyOptions `url:"strategies,omitempty" json:"strategies,omitempty"`
}

// FeatureFlagStrategyOptions represents the available feature flag strategy
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flags/#create-a-feature-flag
type FeatureFlagStrategyOptions struct {
	ID         *int64                               `url:"id,omitempty" json:"id,omitempty"`
	Name       *string                              `url:"name,omitempty" json:"name,omitempty"`
	Parameters *ProjectFeatureFlagStrategyParameter `url:"parameters,omitempty" json:"parameters,omitempty"`
	Scopes     *[]*ProjectFeatureFlagScope          `url:"scopes,omitempty" json:"scopes,omitempty"`
}

// ProjectFeatureFlagScopeOptions represents the available feature flag scope
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flags/#create-a-feature-flag
type ProjectFeatureFlagScopeOptions struct {
	ID               *int64  `url:"id,omitempty" json:"id,omitempty"`
	EnvironmentScope *string `url:"id,omitempty" json:"environment_scope,omitempty"`
}

func (s *ProjectFeatureFlagService) CreateProjectFeatureFlag(pid any, opt *CreateProjectFeatureFlagOptions, options ...RequestOptionFunc) (*ProjectFeatureFlag, *Response, error) {
	return do[*ProjectFeatureFlag](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/feature_flags", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateProjectFeatureFlagOptions represents the available
// UpdateProjectFeatureFlag() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/feature_flags/#update-a-feature-flag
type UpdateProjectFeatureFlagOptions struct {
	Name        *string                        `url:"name,omitempty" json:"name,omitempty"`
	Description *string                        `url:"description,omitempty" json:"description,omitempty"`
	Active      *bool                          `url:"active,omitempty" json:"active,omitempty"`
	Strategies  *[]*FeatureFlagStrategyOptions `url:"strategies,omitempty" json:"strategies,omitempty"`
}

func (s *ProjectFeatureFlagService) UpdateProjectFeatureFlag(pid any, name string, opt *UpdateProjectFeatureFlagOptions, options ...RequestOptionFunc) (*ProjectFeatureFlag, *Response, error) {
	return do[*ProjectFeatureFlag](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/feature_flags/%s", ProjectID{pid}, name),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *ProjectFeatureFlagService) DeleteProjectFeatureFlag(pid any, name string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/feature_flags/%s", ProjectID{pid}, name),
		withRequestOpts(options...),
	)
	return resp, err
}
