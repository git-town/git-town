package gitlab

import (
	"net/http"
	"time"
)

type (
	GroupCredentialsServiceInterface interface {
		// ListGroupPersonalAccessTokens lists all personal access tokens
		// associated with enterprise users in a top-level group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/groups/#list-all-personal-access-tokens-for-a-group
		ListGroupPersonalAccessTokens(gid any, opt *ListGroupPersonalAccessTokensOptions, options ...RequestOptionFunc) ([]*GroupPersonalAccessToken, *Response, error)
		// ListGroupSSHKeys lists all SSH public keys associated with
		// enterprise users in a top-level group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/groups/#list-all-ssh-keys-for-a-group
		ListGroupSSHKeys(gid any, opt *ListGroupSSHKeysOptions, options ...RequestOptionFunc) ([]*GroupSSHKey, *Response, error)
		// RevokeGroupPersonalAccessToken revokes a specified personal access token
		// for an enterprise user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/groups/#revoke-a-personal-access-token-for-an-enterprise-user
		RevokeGroupPersonalAccessToken(gid any, tokenID int64, options ...RequestOptionFunc) (*Response, error)
		// DeleteGroupSSHKey deletes a specified SSH public key for an
		// enterprise user associated with the top-level group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/groups/#delete-an-ssh-key-for-an-enterprise-user
		DeleteGroupSSHKey(gid any, keyID int64, options ...RequestOptionFunc) (*Response, error)
	}

	// GroupCredentialsService handles communication with the top-level group
	// credentials inventory management endpoints of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/groups/#credentials-inventory-management
	GroupCredentialsService struct {
		client *Client
	}
)

// GroupPersonalAccessToken represents a group enterprise users personal access token.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-all-personal-access-tokens-for-a-group
type GroupPersonalAccessToken struct {
	ID          int64      `json:"id"`
	Name        string     `json:"name"`
	Revoked     bool       `json:"revoked"`
	CreatedAt   *time.Time `json:"created_at"`
	Description string     `json:"description"`
	Scopes      []string   `json:"scopes"`
	UserID      int64      `json:"user_id"`
	LastUsedAt  *time.Time `json:"last_used_at,omitempty"`
	Active      bool       `json:"active"`
	ExpiresAt   *ISOTime   `json:"expires_at"`
}

// ListGroupPersonalAccessTokensOptions represents the available
// ListGroupPersonalAccessTokens() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-all-personal-access-tokens-for-a-group
type ListGroupPersonalAccessTokensOptions struct {
	ListOptions
	CreatedAfter   *ISOTime `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore  *ISOTime `url:"created_before,omitempty" json:"created_before,omitempty"`
	LastUsedAfter  *ISOTime `url:"last_used_after,omitempty" json:"last_used_after,omitempty"`
	LastUsedBefore *ISOTime `url:"last_used_before,omitempty" json:"last_used_before,omitempty"`
	Revoked        *bool    `url:"revoked,omitempty" json:"revoked,omitempty"`
	Search         *string  `url:"search,omitempty" json:"search,omitempty"`
	State          *string  `url:"state,omitempty" json:"state,omitempty"`
}

func (g *GroupCredentialsService) ListGroupPersonalAccessTokens(gid any, opt *ListGroupPersonalAccessTokensOptions, options ...RequestOptionFunc) ([]*GroupPersonalAccessToken, *Response, error) {
	return do[[]*GroupPersonalAccessToken](g.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/manage/personal_access_tokens", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GroupSSHKey represents a group enterprise users public SSH key.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-all-ssh-keys-for-a-group
type GroupSSHKey struct {
	ID         int64      `json:"id"`
	Title      string     `json:"title"`
	CreatedAt  *time.Time `json:"created_at"`
	ExpiresAt  *time.Time `json:"expires_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
	UsageType  string     `json:"usage_type"`
	UserID     int64      `json:"user_id"`
}

// ListGroupSSHKeysOptions represents the available
// ListGroupSSHKeys() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/groups/#list-all-ssh-keys-for-a-group
type ListGroupSSHKeysOptions struct {
	ListOptions
	CreatedAfter  *ISOTime `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore *ISOTime `url:"created_before,omitempty" json:"created_before,omitempty"`
	ExpiresBefore *ISOTime `url:"expires_before,omitempty" json:"expires_before,omitempty"`
	ExpiresAfter  *ISOTime `url:"expires_after,omitempty" json:"expires_after,omitempty"`
}

func (g *GroupCredentialsService) ListGroupSSHKeys(gid any, opt *ListGroupSSHKeysOptions, options ...RequestOptionFunc) ([]*GroupSSHKey, *Response, error) {
	return do[[]*GroupSSHKey](g.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/manage/ssh_keys", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (g *GroupCredentialsService) RevokeGroupPersonalAccessToken(gid any, tokenID int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](g.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/manage/personal_access_tokens/%d", GroupID{gid}, tokenID),
		withRequestOpts(options...),
	)
	return resp, err
}

func (g *GroupCredentialsService) DeleteGroupSSHKey(gid any, keyID int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](g.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/manage/ssh_keys/%d", GroupID{gid}, keyID),
		withRequestOpts(options...),
	)
	return resp, err
}
