package bitbucketv1

// UpdatePullRequestCommentRequest represents the request body for the Update pull request comment API.
// API doc: https://developer.atlassian.com/server/bitbucket/rest/v805/api-group-pull-requests/#api-api-latest-projects-projectkey-repos-repositoryslug-pull-requests-pullrequestid-comments-commentid-put
type UpdatePullRequestCommentRequest struct {
	// updated comment ID
	ID int `json:"id"`
	// updated comment updated state
	State string `json:"state,omitempty"`
	// existing comment version which need to be fetched from the actual comment
	Version int `json:"version"`
	// updated comment severity
	Severity string `json:"severity,omitempty"`
	// updated comment text
	Text string `json:"text,omitempty"`
	// updated comment properties (should be able to be marshalled into an object according to the documentation
	Properties interface{} `json:"properties,omitempty"`
}
