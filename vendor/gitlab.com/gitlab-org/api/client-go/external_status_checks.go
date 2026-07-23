package gitlab

import (
	"net/http"
	"time"
)

type (
	// ExternalStatusChecksServiceInterface defines all the API methods for the ExternalStatusChecksService
	ExternalStatusChecksServiceInterface interface {
		// CreateExternalStatusCheck creates an external status check.
		// Deprecated: to be removed in 1.0; use CreateProjectExternalStatusCheck instead
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#create-external-status-check-service
		CreateExternalStatusCheck(pid any, opt *CreateExternalStatusCheckOptions, options ...RequestOptionFunc) (*Response, error)

		// DeleteExternalStatusCheck deletes an external status check.
		// Deprecated: to be removed in 1.0; use DeleteProjectExternalStatusCheck instead
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#delete-external-status-check-service
		DeleteExternalStatusCheck(pid any, check int64, options ...RequestOptionFunc) (*Response, error)

		// UpdateExternalStatusCheck updates an external status check.
		// Deprecated: to be removed in 1.0; use UpdateProjectExternalStatusCheck instead
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#update-external-status-check-service
		UpdateExternalStatusCheck(pid any, check int64, opt *UpdateExternalStatusCheckOptions, options ...RequestOptionFunc) (*Response, error)

		// ListMergeStatusChecks lists the external status checks that apply to it
		// and their status for a single merge request.
		// Deprecated: to be removed in 1.0; use ListProjectMergeRequestExternalStatusChecks instead
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#list-status-checks-for-a-merge-request
		ListMergeStatusChecks(pid any, mr int64, opt *ListOptions, options ...RequestOptionFunc) ([]*MergeStatusCheck, *Response, error)

		// ListProjectStatusChecks lists the project external status checks.
		// Deprecated: to be removed in 1.0; use ListProjectExternalStatusChecks instead
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#get-project-external-status-check-services
		ListProjectStatusChecks(pid any, opt *ListOptions, options ...RequestOptionFunc) ([]*ProjectStatusCheck, *Response, error)

		// RetryFailedStatusCheckForAMergeRequest retries the specified failed external status check.
		// Deprecated: to be removed in 1.0; use RetryFailedExternalStatusCheckForProjectMergeRequest instead
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#retry-failed-status-check-for-a-merge-request
		RetryFailedStatusCheckForAMergeRequest(pid any, mergeRequest int64, externalStatusCheck int64, options ...RequestOptionFunc) (*Response, error)

		// SetExternalStatusCheckStatus sets the status of an external status check.
		// Deprecated: to be removed in 1.0; use SetProjectMergeRequestExternalStatusCheckStatus instead
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#set-status-of-an-external-status-check
		SetExternalStatusCheckStatus(pid any, mergeRequest int64, opt *SetExternalStatusCheckStatusOptions, options ...RequestOptionFunc) (*Response, error)

		// ListProjectMergeRequestExternalStatusChecks lists the external status checks that apply to it
		// and their status for a single merge request.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#list-status-checks-for-a-merge-request
		ListProjectMergeRequestExternalStatusChecks(pid any, mr int64, opt *ListProjectMergeRequestExternalStatusChecksOptions, options ...RequestOptionFunc) ([]*MergeStatusCheck, *Response, error)

		// ListProjectExternalStatusChecks lists the project external status checks.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#get-project-external-status-check-services
		ListProjectExternalStatusChecks(pid any, opt *ListProjectExternalStatusChecksOptions, options ...RequestOptionFunc) ([]*ProjectStatusCheck, *Response, error)

		// RetryFailedExternalStatusCheckForProjectMergeRequest retries the specified failed external status check.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#retry-failed-status-check-for-a-merge-request
		RetryFailedExternalStatusCheckForProjectMergeRequest(pid any, mergeRequest int64, externalStatusCheck int64, opt *RetryFailedExternalStatusCheckForProjectMergeRequestOptions, options ...RequestOptionFunc) (*Response, error)

		// CreateProjectExternalStatusCheck creates an external status check.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#create-external-status-check-service
		CreateProjectExternalStatusCheck(pid any, opt *CreateProjectExternalStatusCheckOptions, options ...RequestOptionFunc) (*ProjectStatusCheck, *Response, error)

		// UpdateProjectExternalStatusCheck updates an external status check.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#update-external-status-check-service
		UpdateProjectExternalStatusCheck(pid any, check int64, opt *UpdateProjectExternalStatusCheckOptions, options ...RequestOptionFunc) (*ProjectStatusCheck, *Response, error)

		// DeleteProjectExternalStatusCheck deletes an external status check.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#delete-external-status-check-service
		DeleteProjectExternalStatusCheck(pid any, check int64, opt *DeleteProjectExternalStatusCheckOptions, options ...RequestOptionFunc) (*Response, error)

		// SetProjectMergeRequestExternalStatusCheckStatus sets the status of an external status check.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/status_checks/#set-status-of-an-external-status-check
		SetProjectMergeRequestExternalStatusCheckStatus(pid any, mergeRequest int64, opt *SetProjectMergeRequestExternalStatusCheckStatusOptions, options ...RequestOptionFunc) (*Response, error)
	}

	// ExternalStatusChecksService handles communication with the external
	// status check related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/status_checks/
	ExternalStatusChecksService struct {
		client *Client
	}
)

var _ ExternalStatusChecksServiceInterface = (*ExternalStatusChecksService)(nil)

type MergeStatusCheck struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	ExternalURL string `json:"external_url"`
	Status      string `json:"status"`
}

type ProjectStatusCheck struct {
	ID                int64                        `json:"id"`
	Name              string                       `json:"name"`
	ProjectID         int64                        `json:"project_id"`
	ExternalURL       string                       `json:"external_url"`
	HMAC              bool                         `json:"hmac"`
	ProtectedBranches []StatusCheckProtectedBranch `json:"protected_branches"`
}

type StatusCheckProtectedBranch struct {
	ID                        int64      `json:"id"`
	ProjectID                 int64      `json:"project_id"`
	Name                      string     `json:"name"`
	CreatedAt                 *time.Time `json:"created_at"`
	UpdatedAt                 *time.Time `json:"updated_at"`
	CodeOwnerApprovalRequired bool       `json:"code_owner_approval_required"`
}

// ListMergeStatusChecks lists the external status checks that apply to it
// and their status for a single merge request.
// Deprecated: to be removed in 1.0; use ListProjectMergeRequestExternalStatusChecks instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/status_checks/#list-status-checks-for-a-merge-request
func (s *ExternalStatusChecksService) ListMergeStatusChecks(pid any, mr int64, opt *ListOptions, options ...RequestOptionFunc) ([]*MergeStatusCheck, *Response, error) {
	return do[[]*MergeStatusCheck](s.client,
		withPath("projects/%s/merge_requests/%d/status_checks", ProjectID{pid}, mr),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// SetExternalStatusCheckStatusOptions represents the available
// SetExternalStatusCheckStatus() options.
// Deprecated: to be removed in 1.0; use SetProjectMergeRequestExternalStatusCheckStatusOptions instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/status_checks/#set-status-of-an-external-status-check
type SetExternalStatusCheckStatusOptions struct {
	SHA                   *string `url:"sha,omitempty" json:"sha,omitempty"`
	ExternalStatusCheckID *int64  `url:"external_status_check_id,omitempty" json:"external_status_check_id,omitempty"`
	Status                *string `url:"status,omitempty" json:"status,omitempty"`
}

// SetExternalStatusCheckStatus sets the status of an external status check.
// Deprecated: to be removed in 1.0; use SetProjectMergeRequestExternalStatusCheckStatus instead
//
// Gitlab API docs:
// https://docs.gitlab.com/api/status_checks/#set-status-of-an-external-status-check
func (s *ExternalStatusChecksService) SetExternalStatusCheckStatus(pid any, mergeRequest int64, opt *SetExternalStatusCheckStatusOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/status_check_responses", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *ExternalStatusChecksService) ListProjectStatusChecks(pid any, opt *ListOptions, options ...RequestOptionFunc) ([]*ProjectStatusCheck, *Response, error) {
	return do[[]*ProjectStatusCheck](s.client,
		withPath("projects/%s/external_status_checks", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreateExternalStatusCheckOptions represents the available
// CreateExternalStatusCheck() options.
// Deprecated: to be removed in 1.0; use CreateProjectExternalStatusCheckOptions instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/status_checks/#create-external-status-check-service
type CreateExternalStatusCheckOptions struct {
	Name               *string  `url:"name,omitempty" json:"name,omitempty"`
	ExternalURL        *string  `url:"external_url,omitempty" json:"external_url,omitempty"`
	ProtectedBranchIDs *[]int64 `url:"protected_branch_ids,omitempty" json:"protected_branch_ids,omitempty"`
}

func (s *ExternalStatusChecksService) CreateExternalStatusCheck(pid any, opt *CreateExternalStatusCheckOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/external_status_checks", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

// DeleteExternalStatusCheck deletes an external status check.
// Deprecated: to be removed in 1.0; use DeleteProjectExternalStatusCheck instead
//
// Gitlab API docs:
// https://docs.gitlab.com/api/status_checks/#delete-external-status-check-service
func (s *ExternalStatusChecksService) DeleteExternalStatusCheck(pid any, check int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/external_status_checks/%d", ProjectID{pid}, check),
		withRequestOpts(options...),
	)
	return resp, err
}

// UpdateExternalStatusCheckOptions represents the available
// UpdateExternalStatusCheck() options.
// Deprecated: to be removed in 1.0; use UpdateProjectExternalStatusCheckOptions instead
//
// GitLab API docs:
// https://docs.gitlab.com/api/status_checks/#update-external-status-check-service
type UpdateExternalStatusCheckOptions struct {
	Name               *string  `url:"name,omitempty" json:"name,omitempty"`
	ExternalURL        *string  `url:"external_url,omitempty" json:"external_url,omitempty"`
	ProtectedBranchIDs *[]int64 `url:"protected_branch_ids,omitempty" json:"protected_branch_ids,omitempty"`
}

// UpdateExternalStatusCheck updates an external status check.
// Deprecated: to be removed in 1.0; use UpdateProjectExternalStatusCheck instead
//
// Gitlab API docs:
// https://docs.gitlab.com/api/status_checks/#update-external-status-check-service
func (s *ExternalStatusChecksService) UpdateExternalStatusCheck(pid any, check int64, opt *UpdateExternalStatusCheckOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/external_status_checks/%d", ProjectID{pid}, check),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

// RetryFailedStatusCheckForAMergeRequest retries the specified failed external status check.
// Deprecated: to be removed in 1.0; use RetryFailedExternalStatusCheckForProjectMergeRequest instead
//
// Gitlab API docs:
// https://docs.gitlab.com/api/status_checks/#retry-failed-status-check-for-a-merge-request
func (s *ExternalStatusChecksService) RetryFailedStatusCheckForAMergeRequest(pid any, mergeRequest int64, externalStatusCheck int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/status_checks/%d/retry", ProjectID{pid}, mergeRequest, externalStatusCheck),
		withRequestOpts(options...),
	)
	return resp, err
}

// ListProjectMergeRequestExternalStatusChecksOptions represents the available
// ListProjectMergeRequestExternalStatusChecks() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/status_checks/#list-status-checks-for-a-merge-request
type ListProjectMergeRequestExternalStatusChecksOptions struct {
	ListOptions
}

// ListProjectMergeRequestExternalStatusChecks lists the external status checks that apply to it
// and their status for a single merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/status_checks/#list-status-checks-for-a-merge-request
func (s *ExternalStatusChecksService) ListProjectMergeRequestExternalStatusChecks(pid any, mr int64, opt *ListProjectMergeRequestExternalStatusChecksOptions, options ...RequestOptionFunc) ([]*MergeStatusCheck, *Response, error) {
	return do[[]*MergeStatusCheck](s.client,
		withPath("projects/%s/merge_requests/%d/status_checks", ProjectID{pid}, mr),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListProjectExternalStatusChecksOptions represents the available
// ListProjectExternalStatusChecks() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/status_checks/#get-project-external-status-check-services
type ListProjectExternalStatusChecksOptions struct {
	ListOptions
}

func (s *ExternalStatusChecksService) ListProjectExternalStatusChecks(pid any, opt *ListProjectExternalStatusChecksOptions, options ...RequestOptionFunc) ([]*ProjectStatusCheck, *Response, error) {
	return do[[]*ProjectStatusCheck](s.client,
		withPath("projects/%s/external_status_checks", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreateProjectExternalStatusCheckOptions represents the available
// CreateProjectExternalStatusCheck() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/status_checks/#create-external-status-check-service
type CreateProjectExternalStatusCheckOptions struct {
	Name               *string  `url:"name,omitempty" json:"name,omitempty"`
	ExternalURL        *string  `url:"external_url,omitempty" json:"external_url,omitempty"`
	SharedSecret       *string  `url:"shared_secret,omitempty" json:"shared_secret,omitempty"`
	ProtectedBranchIDs *[]int64 `url:"protected_branch_ids,omitempty" json:"protected_branch_ids,omitempty"`
}

func (s *ExternalStatusChecksService) CreateProjectExternalStatusCheck(pid any, opt *CreateProjectExternalStatusCheckOptions, options ...RequestOptionFunc) (*ProjectStatusCheck, *Response, error) {
	return do[*ProjectStatusCheck](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/external_status_checks", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteProjectExternalStatusCheckOptions represents the available
// DeleteProjectExternalStatusCheck() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/status_checks/#delete-external-status-check-service
type DeleteProjectExternalStatusCheckOptions struct{}

// DeleteProjectExternalStatusCheck deletes an external status check.
//
// Gitlab API docs:
// https://docs.gitlab.com/api/status_checks/#delete-external-status-check-service
func (s *ExternalStatusChecksService) DeleteProjectExternalStatusCheck(pid any, check int64, opt *DeleteProjectExternalStatusCheckOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/external_status_checks/%d", ProjectID{pid}, check),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

// UpdateProjectExternalStatusCheckOptions represents the available
// UpdateProjectExternalStatusCheck() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/status_checks/#update-external-status-check-service
type UpdateProjectExternalStatusCheckOptions struct {
	Name               *string  `url:"name,omitempty" json:"name,omitempty"`
	ExternalURL        *string  `url:"external_url,omitempty" json:"external_url,omitempty"`
	SharedSecret       *string  `url:"shared_secret,omitempty" json:"shared_secret,omitempty"`
	ProtectedBranchIDs *[]int64 `url:"protected_branch_ids,omitempty" json:"protected_branch_ids,omitempty"`
}

// UpdateProjectExternalStatusCheck updates an external status check.
//
// Gitlab API docs:
// https://docs.gitlab.com/api/status_checks/#update-external-status-check-service
func (s *ExternalStatusChecksService) UpdateProjectExternalStatusCheck(pid any, check int64, opt *UpdateProjectExternalStatusCheckOptions, options ...RequestOptionFunc) (*ProjectStatusCheck, *Response, error) {
	return do[*ProjectStatusCheck](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/external_status_checks/%d", ProjectID{pid}, check),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// RetryFailedExternalStatusCheckForProjectMergeRequestOptions represents the available
// RetryFailedExternalStatusCheckForProjectMergeRequest() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/status_checks/#retry-failed-status-check-for-a-merge-request
type RetryFailedExternalStatusCheckForProjectMergeRequestOptions struct{}

// RetryFailedExternalStatusCheckForProjectMergeRequest retries the specified failed external status check.
//
// Gitlab API docs:
// https://docs.gitlab.com/api/status_checks/#retry-failed-status-check-for-a-merge-request
func (s *ExternalStatusChecksService) RetryFailedExternalStatusCheckForProjectMergeRequest(pid any, mergeRequest int64, externalStatusCheck int64, opt *RetryFailedExternalStatusCheckForProjectMergeRequestOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/status_checks/%d/retry", ProjectID{pid}, mergeRequest, externalStatusCheck),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

// SetProjectMergeRequestExternalStatusCheckStatusOptions represents the available
// SetProjectMergeRequestExternalStatusCheckStatus() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/status_checks/#set-status-of-an-external-status-check
type SetProjectMergeRequestExternalStatusCheckStatusOptions struct {
	SHA                   *string `url:"sha,omitempty" json:"sha,omitempty"`
	ExternalStatusCheckID *int64  `url:"external_status_check_id,omitempty" json:"external_status_check_id,omitempty"`
	Status                *string `url:"status,omitempty" json:"status,omitempty"`
}

// SetProjectMergeRequestExternalStatusCheckStatus sets the status of an external status check.
//
// Gitlab API docs:
// https://docs.gitlab.com/api/status_checks/#set-status-of-an-external-status-check
func (s *ExternalStatusChecksService) SetProjectMergeRequestExternalStatusCheckStatus(pid any, mergeRequest int64, opt *SetProjectMergeRequestExternalStatusCheckStatusOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_requests/%d/status_check_responses", ProjectID{pid}, mergeRequest),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}
