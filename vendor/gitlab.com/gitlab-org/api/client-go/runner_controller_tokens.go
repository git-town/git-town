package gitlab

import (
	"net/http"
	"time"
)

type (
	// RunnerControllerTokensServiceInterface handles communication with the runner
	// controller token related methods of the GitLab API. This is an admin-only
	// endpoint.
	//
	// Note: This API is experimental and may change or be removed in future versions.
	//
	// GitLab API docs: https://docs.gitlab.com/api/runner_controller_tokens/
	RunnerControllerTokensServiceInterface interface {
		// ListRunnerControllerTokens lists all runner controller tokens. This is an
		// admin-only endpoint.
		//
		// GitLab API docs: https://docs.gitlab.com/api/runner_controller_tokens/#list-all-runner-controller-tokens
		ListRunnerControllerTokens(rid int64, opt *ListRunnerControllerTokensOptions, options ...RequestOptionFunc) ([]*RunnerControllerToken, *Response, error)
		// GetRunnerControllerToken retrieves a single runner controller token. This
		// is an admin-only endpoint.
		//
		// GitLab API docs: https://docs.gitlab.com/api/runner_controller_tokens/#retrieve-a-single-runner-controller-token
		GetRunnerControllerToken(rid int64, tokenID int64, options ...RequestOptionFunc) (*RunnerControllerToken, *Response, error)
		// CreateRunnerControllerToken creates a new runner controller token. This is
		// an admin-only endpoint.
		//
		// GitLab API docs: https://docs.gitlab.com/api/runner_controller_tokens/#create-a-runner-controller-token
		CreateRunnerControllerToken(rid int64, opt *CreateRunnerControllerTokenOptions, options ...RequestOptionFunc) (*RunnerControllerToken, *Response, error)
		// RotateRunnerControllerToken rotates an existing runner controller token.
		// This is an admin-only endpoint.
		//
		// GitLab API docs: https://docs.gitlab.com/api/runner_controller_tokens/#rotate-a-runner-controller-token
		RotateRunnerControllerToken(rid int64, tokenID int64, options ...RequestOptionFunc) (*RunnerControllerToken, *Response, error)
		// RevokeRunnerControllerToken revokes a runner controller token. This is an
		// admin-only endpoint.
		//
		// GitLab API docs: https://docs.gitlab.com/api/runner_controller_tokens/#revoke-a-runner-controller-token
		RevokeRunnerControllerToken(rid int64, tokenID int64, options ...RequestOptionFunc) (*Response, error)
	}

	// RunnerControllerTokensService handles communication with the runner
	// controller token related methods of the GitLab API. This is an admin-only
	// endpoint.
	//
	// Note: This API is experimental and may change or be removed in future versions.
	//
	// GitLab API docs: https://docs.gitlab.com/api/runner_controller_tokens/
	RunnerControllerTokensService struct {
		client *Client
	}
)

var _ RunnerControllerTokensServiceInterface = (*RunnerControllerTokensService)(nil)

// RunnerControllerToken represents a GitLab runner controller token.
//
// GitLab API docs: https://docs.gitlab.com/api/runner_controller_tokens/
type RunnerControllerToken struct {
	ID                 int64      `json:"id"`
	RunnerControllerID int64      `json:"runner_controller_id"`
	Description        string     `json:"description"`
	Token              string     `json:"token,omitempty"`
	CreatedAt          *time.Time `json:"created_at"`
	UpdatedAt          *time.Time `json:"updated_at"`
}

// ListRunnerControllerTokensOptions represents the available
// ListRunnerControllerTokens() options.
//
// GitLab API docs: https://docs.gitlab.com/api/runner_controller_tokens/#list-all-runner-controller-tokens
type ListRunnerControllerTokensOptions struct {
	ListOptions
}

func (s *RunnerControllerTokensService) ListRunnerControllerTokens(rid int64, opt *ListRunnerControllerTokensOptions, options ...RequestOptionFunc) ([]*RunnerControllerToken, *Response, error) {
	return do[[]*RunnerControllerToken](s.client,
		withPath("runner_controllers/%d/tokens", rid),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *RunnerControllerTokensService) GetRunnerControllerToken(rid int64, tokenID int64, options ...RequestOptionFunc) (*RunnerControllerToken, *Response, error) {
	return do[*RunnerControllerToken](s.client,
		withPath("runner_controllers/%d/tokens/%d", rid, tokenID),
		withRequestOpts(options...),
	)
}

// CreateRunnerControllerTokenOptions represents the available
// CreateRunnerControllerToken() options.
//
// GitLab API docs: https://docs.gitlab.com/api/runner_controller_tokens/#create-a-runner-controller-token
type CreateRunnerControllerTokenOptions struct {
	Description *string `url:"description,omitempty" json:"description,omitempty"`
}

func (s *RunnerControllerTokensService) CreateRunnerControllerToken(rid int64, opt *CreateRunnerControllerTokenOptions, options ...RequestOptionFunc) (*RunnerControllerToken, *Response, error) {
	return do[*RunnerControllerToken](s.client,
		withMethod(http.MethodPost),
		withPath("runner_controllers/%d/tokens", rid),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *RunnerControllerTokensService) RotateRunnerControllerToken(rid int64, tokenID int64, options ...RequestOptionFunc) (*RunnerControllerToken, *Response, error) {
	return do[*RunnerControllerToken](s.client,
		withMethod(http.MethodPost),
		withPath("runner_controllers/%d/tokens/%d/rotate", rid, tokenID),
		withRequestOpts(options...),
	)
}

func (s *RunnerControllerTokensService) RevokeRunnerControllerToken(rid int64, tokenID int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("runner_controllers/%d/tokens/%d", rid, tokenID),
		withRequestOpts(options...),
	)
	return resp, err
}
