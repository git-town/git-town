package gitlab

import (
	"net/http"
	"time"
)

type (
	// RunnerControllersServiceInterface handles communication with the runner
	// controller related methods of the GitLab API. This is an admin-only endpoint.
	//
	// Note: This API is experimental and may change or be removed in future versions.
	//
	// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/
	RunnerControllersServiceInterface interface {
		// ListRunnerControllers gets a list of runner controllers. This is an
		// admin-only endpoint.
		//
		// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#list-all-runner-controllers
		ListRunnerControllers(opt *ListRunnerControllersOptions, options ...RequestOptionFunc) ([]*RunnerController, *Response, error)
		// GetRunnerController retrieves a single runner controller. This is an
		// admin-only endpoint.
		//
		// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#retrieve-a-single-runner-controller
		GetRunnerController(rid int64, options ...RequestOptionFunc) (*RunnerController, *Response, error)
		// CreateRunnerController registers a new runner controller. This is an
		// admin-only endpoint.
		//
		// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#register-a-runner-controller
		CreateRunnerController(opt *CreateRunnerControllerOptions, options ...RequestOptionFunc) (*RunnerController, *Response, error)
		// UpdateRunnerController updates a runner controller. This is an admin-only
		// endpoint.
		//
		// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#update-a-runner-controller
		UpdateRunnerController(rid int64, opt *UpdateRunnerControllerOptions, options ...RequestOptionFunc) (*RunnerController, *Response, error)
		// DeleteRunnerController deletes a runner controller. This is an admin-only
		// endpoint.
		//
		// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#delete-a-runner-controller
		DeleteRunnerController(rid int64, options ...RequestOptionFunc) (*Response, error)
	}

	// RunnerControllersService handles communication with the runner controller
	// related methods of the GitLab API. This is an admin-only endpoint.
	//
	// Note: This API is experimental and may change or be removed in future versions.
	//
	// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/
	RunnerControllersService struct {
		client *Client
	}
)

var _ RunnerControllersServiceInterface = (*RunnerControllersService)(nil)

// RunnerControllerStateValue represents the state of a runner controller.
//
// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/
type RunnerControllerStateValue string

// These constants represent all valid runner controller states.
const (
	RunnerControllerStateDisabled RunnerControllerStateValue = "disabled"
	RunnerControllerStateEnabled  RunnerControllerStateValue = "enabled"
	RunnerControllerStateDryRun   RunnerControllerStateValue = "dry_run"
)

// RunnerController represents a GitLab runner controller.
//
// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/
type RunnerController struct {
	ID          int64                      `json:"id"`
	Description string                     `json:"description"`
	State       RunnerControllerStateValue `json:"state"`
	CreatedAt   *time.Time                 `json:"created_at"`
	UpdatedAt   *time.Time                 `json:"updated_at"`
}

// ListRunnerControllersOptions represents the available
// ListRunnerControllers() options.
//
// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#list-all-runner-controllers
type ListRunnerControllersOptions struct {
	ListOptions
}

func (s *RunnerControllersService) ListRunnerControllers(opt *ListRunnerControllersOptions, options ...RequestOptionFunc) ([]*RunnerController, *Response, error) {
	return do[[]*RunnerController](s.client,
		withPath("runner_controllers"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *RunnerControllersService) GetRunnerController(rid int64, options ...RequestOptionFunc) (*RunnerController, *Response, error) {
	return do[*RunnerController](s.client,
		withPath("runner_controllers/%d", rid),
		withRequestOpts(options...),
	)
}

// CreateRunnerControllerOptions represents the available
// CreateRunnerController() options.
//
// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#register-a-runner-controller
type CreateRunnerControllerOptions struct {
	Description *string                     `url:"description,omitempty" json:"description,omitempty"`
	State       *RunnerControllerStateValue `url:"state,omitempty" json:"state,omitempty"`
}

func (s *RunnerControllersService) CreateRunnerController(opt *CreateRunnerControllerOptions, options ...RequestOptionFunc) (*RunnerController, *Response, error) {
	return do[*RunnerController](s.client,
		withMethod(http.MethodPost),
		withPath("runner_controllers"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateRunnerControllerOptions represents the available
// UpdateRunnerController() options.
//
// GitLab API docs: https://docs.gitlab.com/api/runner_controllers/#update-a-runner-controller
type UpdateRunnerControllerOptions struct {
	Description *string                     `url:"description,omitempty" json:"description,omitempty"`
	State       *RunnerControllerStateValue `url:"state,omitempty" json:"state,omitempty"`
}

func (s *RunnerControllersService) UpdateRunnerController(rid int64, opt *UpdateRunnerControllerOptions, options ...RequestOptionFunc) (*RunnerController, *Response, error) {
	return do[*RunnerController](s.client,
		withMethod(http.MethodPut),
		withPath("runner_controllers/%d", rid),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *RunnerControllersService) DeleteRunnerController(rid int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("runner_controllers/%d", rid),
		withRequestOpts(options...),
	)
	return resp, err
}
