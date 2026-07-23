// Copyright 2022, Daniela Filipe Bento
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitlab

type (
	// DeploymentMergeRequestsServiceInterface defines all the API methods for the DeploymentMergeRequestsService
	DeploymentMergeRequestsServiceInterface interface {
		// ListDeploymentMergeRequests get the merge requests associated with deployment.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/deployments/#list-of-merge-requests-associated-with-a-deployment
		ListDeploymentMergeRequests(pid any, deployment int64, opts *ListMergeRequestsOptions, options ...RequestOptionFunc) ([]*MergeRequest, *Response, error)
	}

	// DeploymentMergeRequestsService handles communication with the deployment's
	// merge requests related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/deployments/#list-of-merge-requests-associated-with-a-deployment
	DeploymentMergeRequestsService struct {
		client *Client
	}
)

var _ DeploymentMergeRequestsServiceInterface = (*DeploymentMergeRequestsService)(nil)

func (s *DeploymentMergeRequestsService) ListDeploymentMergeRequests(pid any, deployment int64, opts *ListMergeRequestsOptions, options ...RequestOptionFunc) ([]*MergeRequest, *Response, error) {
	return do[[]*MergeRequest](s.client,
		withPath("projects/%s/deployments/%d/merge_requests", ProjectID{pid}, deployment),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}
