package gitlab

import (
	"net/http"
	"time"
)

type (
	// RunnerControllerScopesServiceInterface handles communication with the
	// runner controller scopes related methods of the GitLab API. This is an
	// admin-only endpoint.
	//
	// Note: This API is experimental and may change or be removed in future versions.
	//
	// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#runner-controller-scopes
	RunnerControllerScopesServiceInterface interface {
		// ListRunnerControllerScopes lists all scopes for a specific runner
		// controller. This is an admin-only endpoint.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runner_controllers/#list-all-scopes-for-a-runner-controller
		ListRunnerControllerScopes(rid int64, options ...RequestOptionFunc) (*RunnerControllerScopes, *Response, error)
		// AddRunnerControllerInstanceScope adds an instance-level scope to a
		// runner controller. This is an admin-only endpoint.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runner_controllers/#add-instance-level-scope
		AddRunnerControllerInstanceScope(rid int64, options ...RequestOptionFunc) (*RunnerControllerInstanceLevelScoping, *Response, error)
		// RemoveRunnerControllerInstanceScope removes an instance-level scope
		// from a runner controller. This is an admin-only endpoint.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runner_controllers/#remove-instance-level-scope
		RemoveRunnerControllerInstanceScope(rid int64, options ...RequestOptionFunc) (*Response, error)
		// AddRunnerControllerRunnerScope adds a runner scope to a runner
		// controller. This is an admin-only endpoint. The runner must be an
		// instance-level runner.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runner_controllers/#add-runner-scope
		AddRunnerControllerRunnerScope(rid, runnerID int64, options ...RequestOptionFunc) (*RunnerControllerRunnerLevelScoping, *Response, error)
		// RemoveRunnerControllerRunnerScope removes a runner scope from a runner
		// controller. This is an admin-only endpoint.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/runner_controllers/#remove-runner-scope
		RemoveRunnerControllerRunnerScope(rid, runnerID int64, options ...RequestOptionFunc) (*Response, error)
	}

	// RunnerControllerScopesService handles communication with the runner
	// controller scopes related methods of the GitLab API. This is an admin-only
	// endpoint.
	//
	// Note: This API is experimental and may change or be removed in future versions.
	//
	// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#runner-controller-scopes
	RunnerControllerScopesService struct {
		client *Client
	}
)

var _ RunnerControllerScopesServiceInterface = (*RunnerControllerScopesService)(nil)

// RunnerControllerInstanceLevelScoping represents an instance-level scoping
// for a GitLab runner controller.
//
// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#runner-controller-scopes
type RunnerControllerInstanceLevelScoping struct {
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// RunnerControllerRunnerLevelScoping represents a runner-level scoping for a
// GitLab runner controller.
//
// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#runner-controller-scopes
type RunnerControllerRunnerLevelScoping struct {
	RunnerID  int64      `json:"runner_id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

// RunnerControllerScopes represents all scopes configured for a GitLab runner
// controller.
//
// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#runner-controller-scopes
type RunnerControllerScopes struct {
	InstanceLevelScopings []*RunnerControllerInstanceLevelScoping `json:"instance_level_scopings"`
	RunnerLevelScopings   []*RunnerControllerRunnerLevelScoping   `json:"runner_level_scopings"`
}

func (s *RunnerControllerScopesService) ListRunnerControllerScopes(rid int64, options ...RequestOptionFunc) (*RunnerControllerScopes, *Response, error) {
	return do[*RunnerControllerScopes](s.client,
		withPath("runner_controllers/%d/scopes", rid),
		withRequestOpts(options...),
	)
}

func (s *RunnerControllerScopesService) AddRunnerControllerInstanceScope(rid int64, options ...RequestOptionFunc) (*RunnerControllerInstanceLevelScoping, *Response, error) {
	return do[*RunnerControllerInstanceLevelScoping](s.client,
		withMethod(http.MethodPost),
		withPath("runner_controllers/%d/scopes/instance", rid),
		withRequestOpts(options...),
	)
}

func (s *RunnerControllerScopesService) RemoveRunnerControllerInstanceScope(rid int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("runner_controllers/%d/scopes/instance", rid),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *RunnerControllerScopesService) AddRunnerControllerRunnerScope(rid, runnerID int64, options ...RequestOptionFunc) (*RunnerControllerRunnerLevelScoping, *Response, error) {
	return do[*RunnerControllerRunnerLevelScoping](s.client,
		withMethod(http.MethodPost),
		withPath("runner_controllers/%d/scopes/runners/%d", rid, runnerID),
		withRequestOpts(options...),
	)
}

func (s *RunnerControllerScopesService) RemoveRunnerControllerRunnerScope(rid, runnerID int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("runner_controllers/%d/scopes/runners/%d", rid, runnerID),
		withRequestOpts(options...),
	)
	return resp, err
}
