package gitlab

import (
	"net/http"
	"time"
)

type (
	MergeTrainsServiceInterface interface {
		ListProjectMergeTrains(pid any, opt *ListMergeTrainsOptions, options ...RequestOptionFunc) ([]*MergeTrain, *Response, error)
		ListMergeRequestInMergeTrain(pid any, targetBranch string, opts *ListMergeTrainsOptions, options ...RequestOptionFunc) ([]*MergeTrain, *Response, error)
		GetMergeRequestOnAMergeTrain(pid any, mergeRequest int64, options ...RequestOptionFunc) (*MergeTrain, *Response, error)
		AddMergeRequestToMergeTrain(pid any, mergeRequest int64, opts *AddMergeRequestToMergeTrainOptions, options ...RequestOptionFunc) ([]*MergeTrain, *Response, error)
	}

	// MergeTrainsService handles communication with the merge trains related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/merge_trains/
	MergeTrainsService struct {
		client *Client
	}
)

var _ MergeTrainsServiceInterface = (*MergeTrainsService)(nil)

// MergeTrain represents a GitLab merge train.
//
// GitLab API docs: https://docs.gitlab.com/api/merge_trains/
type MergeTrain struct {
	ID           int64                   `json:"id"`
	MergeRequest *MergeTrainMergeRequest `json:"merge_request"`
	User         *BasicUser              `json:"user"`
	Pipeline     *Pipeline               `json:"pipeline"`
	CreatedAt    *time.Time              `json:"created_at"`
	UpdatedAt    *time.Time              `json:"updated_at"`
	TargetBranch string                  `json:"target_branch"`
	Status       string                  `json:"status"`
	MergedAt     *time.Time              `json:"merged_at"`
	Duration     int64                   `json:"duration"`
}

// MergeTrainMergeRequest represents a GitLab merge request inside merge train.
//
// GitLab API docs: https://docs.gitlab.com/api/merge_trains/
type MergeTrainMergeRequest struct {
	ID          int64      `json:"id"`
	IID         int64      `json:"iid"`
	ProjectID   int64      `json:"project_id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	State       string     `json:"state"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	WebURL      string     `json:"web_url"`
}

// ListMergeTrainsOptions represents the available ListMergeTrain() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_trains/#list-merge-trains-for-a-project
type ListMergeTrainsOptions struct {
	ListOptions
	Scope *string `url:"scope,omitempty" json:"scope,omitempty"`
	Sort  *string `url:"sort,omitempty" json:"sort,omitempty"`
}

// ListProjectMergeTrains get a list of merge trains in a project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_trains/#list-merge-trains-for-a-project
func (s *MergeTrainsService) ListProjectMergeTrains(pid any, opt *ListMergeTrainsOptions, options ...RequestOptionFunc) ([]*MergeTrain, *Response, error) {
	return do[[]*MergeTrain](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_trains", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListMergeRequestInMergeTrain gets a list of merge requests added to a merge
// train for the requested target branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_trains/#list-merge-requests-in-a-merge-train
func (s *MergeTrainsService) ListMergeRequestInMergeTrain(pid any, targetBranch string, opts *ListMergeTrainsOptions, options ...RequestOptionFunc) ([]*MergeTrain, *Response, error) {
	return do[[]*MergeTrain](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/merge_trains/%s", ProjectID{pid}, targetBranch),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

// GetMergeRequestOnAMergeTrain Get merge train information for the requested
// merge request.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_trains/#get-the-status-of-a-merge-request-on-a-merge-train
func (s *MergeTrainsService) GetMergeRequestOnAMergeTrain(pid any, mergeRequest int64, options ...RequestOptionFunc) (*MergeTrain, *Response, error) {
	return do[*MergeTrain](s.client,
		withPath("projects/%s/merge_trains/merge_requests/%d", ProjectID{pid}, mergeRequest),
		withRequestOpts(options...),
	)
}

// AddMergeRequestToMergeTrainOptions represents the available
// AddMergeRequestToMergeTrain() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_trains/#add-a-merge-request-to-a-merge-train
type AddMergeRequestToMergeTrainOptions struct {
	AutoMerge *bool   `url:"auto_merge,omitempty" json:"auto_merge,omitempty"`
	SHA       *string `url:"sha,omitempty" json:"sha,omitempty"`
	Squash    *bool   `url:"squash,omitempty" json:"squash,omitempty"`

	// Deprecated: in 17.11, use AutoMerge instead
	WhenPipelineSucceeds *bool `url:"when_pipeline_succeeds,omitempty" json:"when_pipeline_succeeds,omitempty"`
}

// AddMergeRequestToMergeTrain Add a merge request to the merge train targeting
// the merge request’s target branch.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_trains/#add-a-merge-request-to-a-merge-train
func (s *MergeTrainsService) AddMergeRequestToMergeTrain(pid any, mergeRequest int64, opts *AddMergeRequestToMergeTrainOptions, options ...RequestOptionFunc) ([]*MergeTrain, *Response, error) {
	return do[[]*MergeTrain](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/merge_trains/merge_requests/%d", ProjectID{pid}, mergeRequest),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}
