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

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"text/template"
	"time"
)

type (
	UsersServiceInterface interface {
		// ListUsers gets a list of users.
		//
		// GitLab API docs: https://docs.gitlab.com/api/users/#list-users
		ListUsers(opt *ListUsersOptions, options ...RequestOptionFunc) ([]*User, *Response, error)
		// GetUser gets a single user.
		//
		// GitLab API docs: https://docs.gitlab.com/api/users/#get-a-single-user
		GetUser(user int64, opt GetUsersOptions, options ...RequestOptionFunc) (*User, *Response, error)
		// CreateUser creates a new user. Note only administrators can create new users.
		//
		// GitLab API docs: https://docs.gitlab.com/api/users/#create-a-user
		CreateUser(opt *CreateUserOptions, options ...RequestOptionFunc) (*User, *Response, error)
		// ModifyUser modifies an existing user. Only administrators can change attributes
		// of a user.
		//
		// GitLab API docs: https://docs.gitlab.com/api/users/#modify-a-user
		ModifyUser(user int64, opt *ModifyUserOptions, options ...RequestOptionFunc) (*User, *Response, error)
		// DeleteUser deletes a user. Available only for administrators. This is an
		// idempotent function, calling this function for a non-existent user id still
		// returns a status code 200 OK. The JSON response differs if the user was
		// actually deleted or not. In the former the user is returned and in the
		// latter not.
		//
		// GitLab API docs: https://docs.gitlab.com/api/users/#delete-a-user
		DeleteUser(user int64, options ...RequestOptionFunc) (*Response, error)
		// CurrentUser gets currently authenticated user.
		//
		// GitLab API docs: https://docs.gitlab.com/api/users/#get-the-current-user
		CurrentUser(options ...RequestOptionFunc) (*User, *Response, error)
		// CurrentUserStatus retrieves the user status
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/users/#get-your-user-status
		CurrentUserStatus(options ...RequestOptionFunc) (*UserStatus, *Response, error)
		// GetUserStatus retrieves a user's status.
		//
		// uid can be either a user ID (int) or a username (string); will trim one "@" character off the username, if present.
		// Other types will cause an error to be returned.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/users/#get-the-status-of-a-user
		GetUserStatus(uid any, options ...RequestOptionFunc) (*UserStatus, *Response, error)
		// SetUserStatus sets the user's status
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/users/#set-your-user-status
		SetUserStatus(opt *UserStatusOptions, options ...RequestOptionFunc) (*UserStatus, *Response, error)
		// GetUserAssociationsCount gets a list of a specified user's associations.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/users/#get-a-count-of-a-users-projects-groups-issues-and-merge-requests
		GetUserAssociationsCount(user int64, options ...RequestOptionFunc) (*UserAssociationsCount, *Response, error)
		// ListSSHKeys gets a list of currently authenticated user's SSH keys.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_keys/#list-all-ssh-keys
		ListSSHKeys(opt *ListSSHKeysOptions, options ...RequestOptionFunc) ([]*SSHKey, *Response, error)
		// ListSSHKeysForUser gets a list of a specified user's SSH keys.
		//
		// uid can be either a user ID (int) or a username (string). If a username
		// is provided with a leading "@" (e.g., "@johndoe"), it will be trimmed.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_keys/#list-all-ssh-keys-for-a-user
		ListSSHKeysForUser(uid any, opt *ListSSHKeysForUserOptions, options ...RequestOptionFunc) ([]*SSHKey, *Response, error)
		// GetSSHKey gets a single key.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_keys/#get-an-ssh-key
		GetSSHKey(key int64, options ...RequestOptionFunc) (*SSHKey, *Response, error)
		// GetSSHKeyForUser gets a single key for a given user.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_keys/#get-an-ssh-key-for-a-user
		GetSSHKeyForUser(user int64, key int64, options ...RequestOptionFunc) (*SSHKey, *Response, error)
		// AddSSHKey creates a new key owned by the currently authenticated user.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_keys/#add-an-ssh-key
		AddSSHKey(opt *AddSSHKeyOptions, options ...RequestOptionFunc) (*SSHKey, *Response, error)
		// AddSSHKeyForUser creates new key owned by specified user. Available only for
		// admin.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_keys/#add-an-ssh-key-for-a-user
		AddSSHKeyForUser(user int64, opt *AddSSHKeyOptions, options ...RequestOptionFunc) (*SSHKey, *Response, error)
		// DeleteSSHKey deletes key owned by currently authenticated user. This is an
		// idempotent function and calling it on a key that is already deleted or not
		// available results in 200 OK.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_keys/#delete-an-ssh-key
		DeleteSSHKey(key int64, options ...RequestOptionFunc) (*Response, error)
		// DeleteSSHKeyForUser deletes key owned by a specified user. Available only
		// for admin.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_keys/#delete-an-ssh-key-for-a-user
		DeleteSSHKeyForUser(user, key int64, options ...RequestOptionFunc) (*Response, error)
		// ListGPGKeys gets a list of currently authenticated user’s GPG keys.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_keys/#list-all-gpg-keys
		ListGPGKeys(options ...RequestOptionFunc) ([]*GPGKey, *Response, error)
		// GetGPGKey gets a specific GPG key of currently authenticated user.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_keys/#get-a-gpg-key
		GetGPGKey(key int64, options ...RequestOptionFunc) (*GPGKey, *Response, error)
		// AddGPGKey creates a new GPG key owned by the currently authenticated user.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_keys/#add-a-gpg-key
		AddGPGKey(opt *AddGPGKeyOptions, options ...RequestOptionFunc) (*GPGKey, *Response, error)
		// DeleteGPGKey deletes a GPG key owned by currently authenticated user.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_keys/#delete-a-gpg-key
		DeleteGPGKey(key int64, options ...RequestOptionFunc) (*Response, error)
		// ListGPGKeysForUser gets a list of a specified user’s GPG keys.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_keys/#list-all-gpg-keys-for-a-user
		ListGPGKeysForUser(user int64, options ...RequestOptionFunc) ([]*GPGKey, *Response, error)
		// GetGPGKeyForUser gets a specific GPG key for a given user.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_keys/#get-a-gpg-key-for-a-user
		GetGPGKeyForUser(user, key int64, options ...RequestOptionFunc) (*GPGKey, *Response, error)
		// AddGPGKeyForUser creates new GPG key owned by the specified user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_keys/#add-a-gpg-key-for-a-user
		AddGPGKeyForUser(user int64, opt *AddGPGKeyOptions, options ...RequestOptionFunc) (*GPGKey, *Response, error)
		// DeleteGPGKeyForUser deletes a GPG key owned by a specified user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_keys/#delete-a-gpg-key-for-a-user
		DeleteGPGKeyForUser(user, key int64, options ...RequestOptionFunc) (*Response, error)
		// ListEmails gets a list of currently authenticated user's Emails.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_email_addresses/#list-all-email-addresses
		ListEmails(options ...RequestOptionFunc) ([]*Email, *Response, error)
		// ListEmailsForUser gets a list of a specified user's Emails. Available
		// only for admin
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_email_addresses/#list-all-email-addresses-for-a-user
		ListEmailsForUser(user int64, opt *ListEmailsForUserOptions, options ...RequestOptionFunc) ([]*Email, *Response, error)
		// GetEmail gets a single email.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_email_addresses/#get-details-on-an-email-address
		GetEmail(email int64, options ...RequestOptionFunc) (*Email, *Response, error)
		// AddEmail creates a new email owned by the currently authenticated user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_email_addresses/#add-an-email-address
		AddEmail(opt *AddEmailOptions, options ...RequestOptionFunc) (*Email, *Response, error)
		// AddEmailForUser creates new email owned by specified user. Available only for
		// admin.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_email_addresses/#add-an-email-address-for-a-user
		AddEmailForUser(user int64, opt *AddEmailOptions, options ...RequestOptionFunc) (*Email, *Response, error)
		// DeleteEmail deletes email owned by currently authenticated user. This is an
		// idempotent function and calling it on a key that is already deleted or not
		// available results in 200 OK.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_email_addresses/#delete-an-email-address
		DeleteEmail(email int64, options ...RequestOptionFunc) (*Response, error)
		// DeleteEmailForUser deletes email owned by a specified user. Available only
		// for admin.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_email_addresses/#delete-an-email-address-for-a-user
		DeleteEmailForUser(user, email int64, options ...RequestOptionFunc) (*Response, error)
		// BlockUser blocks the specified user. Available only for admin.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_moderation/#block-access-to-a-user
		BlockUser(user int64, options ...RequestOptionFunc) error
		// UnblockUser unblocks the specified user. Available only for admin.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_moderation/#unblock-access-to-a-user
		UnblockUser(user int64, options ...RequestOptionFunc) error
		// BanUser bans the specified user. Available only for admin.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_moderation/#ban-a-user
		BanUser(user int64, options ...RequestOptionFunc) error
		// UnbanUser unbans the specified user. Available only for admin.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_moderation/#unban-a-user
		UnbanUser(user int64, options ...RequestOptionFunc) error
		// DeactivateUser deactivate the specified user. Available only for admin.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_moderation/#deactivate-a-user
		DeactivateUser(user int64, options ...RequestOptionFunc) error
		// ActivateUser activate the specified user. Available only for admin.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_moderation/#reactivate-a-user
		ActivateUser(user int64, options ...RequestOptionFunc) error
		// ApproveUser approve the specified user. Available only for admin.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_moderation/#approve-access-to-a-user
		ApproveUser(user int64, options ...RequestOptionFunc) error
		// RejectUser reject the specified user. Available only for admin.
		//
		// GitLab API docs: https://docs.gitlab.com/api/user_moderation/#reject-access-to-a-user
		RejectUser(user int64, options ...RequestOptionFunc) error
		// GetAllImpersonationTokens retrieves all impersonation tokens of a user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_tokens/#list-all-impersonation-tokens-for-a-user
		GetAllImpersonationTokens(user int64, opt *GetAllImpersonationTokensOptions, options ...RequestOptionFunc) ([]*ImpersonationToken, *Response, error)
		// GetImpersonationToken retrieves an impersonation token of a user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_tokens/#get-an-impersonation-token-for-a-user
		GetImpersonationToken(user, token int64, options ...RequestOptionFunc) (*ImpersonationToken, *Response, error)
		// CreateImpersonationToken creates an impersonation token.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_tokens/#create-an-impersonation-token
		CreateImpersonationToken(user int64, opt *CreateImpersonationTokenOptions, options ...RequestOptionFunc) (*ImpersonationToken, *Response, error)
		// RevokeImpersonationToken revokes an impersonation token.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_tokens/#revoke-an-impersonation-token
		RevokeImpersonationToken(user, token int64, options ...RequestOptionFunc) (*Response, error)
		// CreatePersonalAccessToken creates a personal access token.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_tokens/#create-a-personal-access-token-for-a-user
		CreatePersonalAccessToken(user int64, opt *CreatePersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error)
		// CreatePersonalAccessTokenForCurrentUser creates a personal access token with limited scopes for the currently authenticated user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_tokens/#create-a-personal-access-token
		CreatePersonalAccessTokenForCurrentUser(opt *CreatePersonalAccessTokenForCurrentUserOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error)
		// GetUserActivities retrieves user activities (admin only)
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/users/#list-a-users-activity
		GetUserActivities(opt *GetUserActivitiesOptions, options ...RequestOptionFunc) ([]*UserActivity, *Response, error)
		// GetUserMemberships retrieves a list of the user's memberships.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/users/#list-projects-and-groups-that-a-user-is-a-member-of
		GetUserMemberships(user int64, opt *GetUserMembershipOptions, options ...RequestOptionFunc) ([]*UserMembership, *Response, error)
		// DisableTwoFactor disables two factor authentication for the specified user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/users/#disable-two-factor-authentication-for-a-user
		DisableTwoFactor(user int64, options ...RequestOptionFunc) error
		// CreateUserRunner creates a runner linked to the current user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/users/#create-a-runner-linked-to-a-user
		CreateUserRunner(opts *CreateUserRunnerOptions, options ...RequestOptionFunc) (*UserRunner, *Response, error)
		// CreateServiceAccountUser creates a new service account user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_service_accounts/#create-a-service-account-user
		CreateServiceAccountUser(opts *CreateServiceAccountUserOptions, options ...RequestOptionFunc) (*User, *Response, error)
		// ListServiceAccounts lists all service accounts.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/user_service_accounts/#list-all-service-account-users
		ListServiceAccounts(opt *ListServiceAccountsOptions, options ...RequestOptionFunc) ([]*ServiceAccount, *Response, error)
		// UploadAvatar uploads an avatar to the current user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/users/#upload-an-avatar-for-yourself
		UploadAvatar(avatar io.Reader, filename string, options ...RequestOptionFunc) (*User, *Response, error)
		// DeleteUserIdentity deletes a user's authentication identity using the provider
		// name associated with that identity. Only available for administrators.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/users/#delete-authentication-identity-from-a-user
		DeleteUserIdentity(user int64, provider string, options ...RequestOptionFunc) (*Response, error)

		// events.go
		ListUserContributionEvents(uid any, opt *ListContributionEventsOptions, options ...RequestOptionFunc) ([]*ContributionEvent, *Response, error)
	}

	// UsersService handles communication with the user related methods of
	// the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/users/
	UsersService struct {
		client *Client
	}
)

var _ UsersServiceInterface = (*UsersService)(nil)

// List a couple of standard errors.
var (
	ErrUserActivatePrevented         = errors.New("cannot activate a user that is blocked by admin or by LDAP synchronization")
	ErrUserApprovePrevented          = errors.New("cannot approve a user that is blocked by admin or by LDAP synchronization")
	ErrUserBlockPrevented            = errors.New("cannot block a user that is already blocked by LDAP synchronization")
	ErrUserConflict                  = errors.New("user does not have a pending request")
	ErrUserDeactivatePrevented       = errors.New("cannot deactivate a user that is blocked by admin or by LDAP synchronization")
	ErrUserDisableTwoFactorPrevented = errors.New("cannot disable two factor authentication if not authenticated as administrator")
	ErrUserNotFound                  = errors.New("user does not exist")
	ErrUserRejectPrevented           = errors.New("cannot reject a user if not authenticated as administrator")
	ErrUserTwoFactorNotEnabled       = errors.New("cannot disable two factor authentication if not enabled")
	ErrUserUnblockPrevented          = errors.New("cannot unblock a user that is blocked by LDAP synchronization")

	errUnexpectedResultCode = errors.New("received unexpected result code")
)

// BasicUser included in other service responses (such as merge requests, pipelines, etc).
type BasicUser struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`

	// State represents the administrative status of the user account.
	// Common values: "active", "blocked", "deactivated", "banned",
	// "ldap_blocked", "blocked_pending_approval".
	//
	// This is independent from the Locked field: State tracks permanent
	// administrative actions, while Locked handles temporary login failures.
	State string `json:"state"`

	// Locked indicates whether the user account is temporarily locked due to
	// excessive failed login attempts. This is separate from administrative
	// blocking (the State field). Locks automatically expire after a configured
	// time period (default: 10 minutes).
	Locked bool `json:"locked"`

	CreatedAt *time.Time `json:"created_at"`
	AvatarURL string     `json:"avatar_url"`
	WebURL    string     `json:"web_url"`
}

// ServiceAccount represents a GitLab service account.
//
// GitLab API docs:
// https://docs.gitlab.com/api/user_service_accounts/
type ServiceAccount struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Name     string `json:"name"`
}

// User represents a GitLab user.
//
// GitLab API docs: https://docs.gitlab.com/api/users/
type User struct {
	ID                             int64              `json:"id"`
	Username                       string             `json:"username"`
	Email                          string             `json:"email"`
	Name                           string             `json:"name"`
	State                          string             `json:"state"`
	WebURL                         string             `json:"web_url"`
	CreatedAt                      *time.Time         `json:"created_at"`
	Bio                            string             `json:"bio"`
	Bot                            bool               `json:"bot"`
	Location                       string             `json:"location"`
	PublicEmail                    string             `json:"public_email"`
	Skype                          string             `json:"skype"`
	Linkedin                       string             `json:"linkedin"`
	Twitter                        string             `json:"twitter"`
	WebsiteURL                     string             `json:"website_url"`
	Organization                   string             `json:"organization"`
	JobTitle                       string             `json:"job_title"`
	ExternUID                      string             `json:"extern_uid"`
	Provider                       string             `json:"provider"`
	ThemeID                        int64              `json:"theme_id"`
	LastActivityOn                 *ISOTime           `json:"last_activity_on"`
	ColorSchemeID                  int64              `json:"color_scheme_id"`
	IsAdmin                        bool               `json:"is_admin"`
	IsAuditor                      bool               `json:"is_auditor"`
	AvatarURL                      string             `json:"avatar_url"`
	CanCreateGroup                 bool               `json:"can_create_group"`
	CanCreateProject               bool               `json:"can_create_project"`
	CanCreateOrganization          bool               `json:"can_create_organization"`
	ProjectsLimit                  int64              `json:"projects_limit"`
	CurrentSignInAt                *time.Time         `json:"current_sign_in_at"`
	CurrentSignInIP                *net.IP            `json:"current_sign_in_ip"`
	LastSignInAt                   *time.Time         `json:"last_sign_in_at"`
	LastSignInIP                   *net.IP            `json:"last_sign_in_ip"`
	ConfirmedAt                    *time.Time         `json:"confirmed_at"`
	TwoFactorEnabled               bool               `json:"two_factor_enabled"`
	Note                           string             `json:"note"`
	Identities                     []*UserIdentity    `json:"identities"`
	External                       bool               `json:"external"`
	PrivateProfile                 bool               `json:"private_profile"`
	SharedRunnersMinutesLimit      int64              `json:"shared_runners_minutes_limit"`
	ExtraSharedRunnersMinutesLimit int64              `json:"extra_shared_runners_minutes_limit"`
	UsingLicenseSeat               bool               `json:"using_license_seat"`
	CustomAttributes               []*CustomAttribute `json:"custom_attributes"`
	NamespaceID                    int64              `json:"namespace_id"`
	Locked                         bool               `json:"locked"`
	CreatedBy                      *BasicUser         `json:"created_by"`
}

// UserIdentity represents a user identity.
type UserIdentity struct {
	Provider  string `json:"provider"`
	ExternUID string `json:"extern_uid"`
}

// UserAvatar represents a GitLab user avatar.
//
// GitLab API docs: https://docs.gitlab.com/api/users/
type UserAvatar struct {
	Filename string
	Image    io.Reader
}

// MarshalJSON implements the json.Marshaler interface.
func (a *UserAvatar) MarshalJSON() ([]byte, error) {
	if a.Filename == "" && a.Image == nil {
		return []byte(`""`), nil
	}
	type alias UserAvatar
	return json.Marshal((*alias)(a))
}

// ListUsersOptions represents the available ListUsers() options.
//
// GitLab API docs: https://docs.gitlab.com/api/users/#list-users
type ListUsersOptions struct {
	ListOptions
	Active          *bool   `url:"active,omitempty" json:"active,omitempty"`
	Blocked         *bool   `url:"blocked,omitempty" json:"blocked,omitempty"`
	Humans          *bool   `url:"humans,omitempty" json:"humans,omitempty"`
	ExcludeInternal *bool   `url:"exclude_internal,omitempty" json:"exclude_internal,omitempty"`
	ExcludeActive   *bool   `url:"exclude_active,omitempty" json:"exclude_active,omitempty"`
	ExcludeExternal *bool   `url:"exclude_external,omitempty" json:"exclude_external,omitempty"`
	ExcludeHumans   *bool   `url:"exclude_humans,omitempty" json:"exclude_humans,omitempty"`
	PublicEmail     *string `url:"public_email,omitempty" json:"public_email,omitempty"`

	// The options below are only available for admins.
	Search               *string    `url:"search,omitempty" json:"search,omitempty"`
	Username             *string    `url:"username,omitempty" json:"username,omitempty"`
	ExternalUID          *string    `url:"extern_uid,omitempty" json:"extern_uid,omitempty"`
	Provider             *string    `url:"provider,omitempty" json:"provider,omitempty"`
	CreatedBefore        *time.Time `url:"created_before,omitempty" json:"created_before,omitempty"`
	CreatedAfter         *time.Time `url:"created_after,omitempty" json:"created_after,omitempty"`
	OrderBy              *string    `url:"order_by,omitempty" json:"order_by,omitempty"`
	Sort                 *string    `url:"sort,omitempty" json:"sort,omitempty"`
	TwoFactor            *string    `url:"two_factor,omitempty" json:"two_factor,omitempty"`
	Admins               *bool      `url:"admins,omitempty" json:"admins,omitempty"`
	External             *bool      `url:"external,omitempty" json:"external,omitempty"`
	WithoutProjects      *bool      `url:"without_projects,omitempty" json:"without_projects,omitempty"`
	WithCustomAttributes *bool      `url:"with_custom_attributes,omitempty" json:"with_custom_attributes,omitempty"`
	WithoutProjectBots   *bool      `url:"without_project_bots,omitempty" json:"without_project_bots,omitempty"`
}

func (s *UsersService) ListUsers(opt *ListUsersOptions, options ...RequestOptionFunc) ([]*User, *Response, error) {
	return do[[]*User](s.client,
		withPath("users"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// GetUsersOptions represents the available GetUser() options.
//
// GitLab API docs: https://docs.gitlab.com/api/users/#get-a-single-user
type GetUsersOptions struct {
	WithCustomAttributes *bool `url:"with_custom_attributes,omitempty" json:"with_custom_attributes,omitempty"`
}

func (s *UsersService) GetUser(user int64, opt GetUsersOptions, options ...RequestOptionFunc) (*User, *Response, error) {
	return do[*User](s.client,
		withPath("users/%d", user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreateUserOptions represents the available CreateUser() options.
//
// GitLab API docs: https://docs.gitlab.com/api/users/#create-a-user
type CreateUserOptions struct {
	Admin               *bool       `url:"admin,omitempty" json:"admin,omitempty"`
	Avatar              *UserAvatar `url:"-" json:"-"`
	Bio                 *string     `url:"bio,omitempty" json:"bio,omitempty"`
	CanCreateGroup      *bool       `url:"can_create_group,omitempty" json:"can_create_group,omitempty"`
	Email               *string     `url:"email,omitempty" json:"email,omitempty"`
	External            *bool       `url:"external,omitempty" json:"external,omitempty"`
	ExternUID           *string     `url:"extern_uid,omitempty" json:"extern_uid,omitempty"`
	ForceRandomPassword *bool       `url:"force_random_password,omitempty" json:"force_random_password,omitempty"`
	JobTitle            *string     `url:"job_title,omitempty" json:"job_title,omitempty"`
	Linkedin            *string     `url:"linkedin,omitempty" json:"linkedin,omitempty"`
	Location            *string     `url:"location,omitempty" json:"location,omitempty"`
	Name                *string     `url:"name,omitempty" json:"name,omitempty"`
	Note                *string     `url:"note,omitempty" json:"note,omitempty"`
	Organization        *string     `url:"organization,omitempty" json:"organization,omitempty"`
	Password            *string     `url:"password,omitempty" json:"password,omitempty"`
	PrivateProfile      *bool       `url:"private_profile,omitempty" json:"private_profile,omitempty"`
	ProjectsLimit       *int64      `url:"projects_limit,omitempty" json:"projects_limit,omitempty"`
	Provider            *string     `url:"provider,omitempty" json:"provider,omitempty"`
	ResetPassword       *bool       `url:"reset_password,omitempty" json:"reset_password,omitempty"`
	SkipConfirmation    *bool       `url:"skip_confirmation,omitempty" json:"skip_confirmation,omitempty"`
	Skype               *string     `url:"skype,omitempty" json:"skype,omitempty"`
	ThemeID             *int64      `url:"theme_id,omitempty" json:"theme_id,omitempty"`
	Twitter             *string     `url:"twitter,omitempty" json:"twitter,omitempty"`
	Username            *string     `url:"username,omitempty" json:"username,omitempty"`
	WebsiteURL          *string     `url:"website_url,omitempty" json:"website_url,omitempty"`
	ViewDiffsFileByFile *bool       `url:"view_diffs_file_by_file,omitempty" json:"view_diffs_file_by_file,omitempty"`
	PublicEmail         *string     `url:"public_email,omitempty" json:"public_email,omitempty"`
}

func (s *UsersService) CreateUser(opt *CreateUserOptions, options ...RequestOptionFunc) (*User, *Response, error) {
	reqOpts := []doOption{
		withMethod(http.MethodPost),
		withPath("users"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	}
	if opt.Avatar != nil {
		reqOpts = append(reqOpts, withUpload(opt.Avatar.Image, opt.Avatar.Filename, UploadAvatar))
	}
	return do[*User](s.client, reqOpts...)
}

// ModifyUserOptions represents the available ModifyUser() options.
//
// GitLab API docs: https://docs.gitlab.com/api/users/#modify-a-user
type ModifyUserOptions struct {
	Admin               *bool       `url:"admin,omitempty" json:"admin,omitempty"`
	Avatar              *UserAvatar `url:"-" json:"avatar,omitempty"`
	Bio                 *string     `url:"bio,omitempty" json:"bio,omitempty"`
	CanCreateGroup      *bool       `url:"can_create_group,omitempty" json:"can_create_group,omitempty"`
	CommitEmail         *string     `url:"commit_email,omitempty" json:"commit_email,omitempty"`
	Email               *string     `url:"email,omitempty" json:"email,omitempty"`
	External            *bool       `url:"external,omitempty" json:"external,omitempty"`
	ExternUID           *string     `url:"extern_uid,omitempty" json:"extern_uid,omitempty"`
	JobTitle            *string     `url:"job_title,omitempty" json:"job_title,omitempty"`
	Linkedin            *string     `url:"linkedin,omitempty" json:"linkedin,omitempty"`
	Location            *string     `url:"location,omitempty" json:"location,omitempty"`
	Name                *string     `url:"name,omitempty" json:"name,omitempty"`
	Note                *string     `url:"note,omitempty" json:"note,omitempty"`
	Organization        *string     `url:"organization,omitempty" json:"organization,omitempty"`
	Password            *string     `url:"password,omitempty" json:"password,omitempty"`
	PrivateProfile      *bool       `url:"private_profile,omitempty" json:"private_profile,omitempty"`
	ProjectsLimit       *int64      `url:"projects_limit,omitempty" json:"projects_limit,omitempty"`
	Provider            *string     `url:"provider,omitempty" json:"provider,omitempty"`
	PublicEmail         *string     `url:"public_email,omitempty" json:"public_email,omitempty"`
	SkipReconfirmation  *bool       `url:"skip_reconfirmation,omitempty" json:"skip_reconfirmation,omitempty"`
	Skype               *string     `url:"skype,omitempty" json:"skype,omitempty"`
	ThemeID             *int64      `url:"theme_id,omitempty" json:"theme_id,omitempty"`
	Twitter             *string     `url:"twitter,omitempty" json:"twitter,omitempty"`
	Username            *string     `url:"username,omitempty" json:"username,omitempty"`
	WebsiteURL          *string     `url:"website_url,omitempty" json:"website_url,omitempty"`
	ViewDiffsFileByFile *bool       `url:"view_diffs_file_by_file,omitempty" json:"view_diffs_file_by_file,omitempty"`
}

func (s *UsersService) ModifyUser(user int64, opt *ModifyUserOptions, options ...RequestOptionFunc) (*User, *Response, error) {
	reqOpts := []doOption{
		withMethod(http.MethodPut),
		withPath("users/%d", user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	}
	if opt.Avatar != nil && (opt.Avatar.Filename != "" || opt.Avatar.Image != nil) {
		reqOpts = append(reqOpts, withUpload(opt.Avatar.Image, opt.Avatar.Filename, UploadAvatar))
	}
	return do[*User](s.client, reqOpts...)
}

func (s *UsersService) DeleteUser(user int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("users/%d", user),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *UsersService) CurrentUser(options ...RequestOptionFunc) (*User, *Response, error) {
	return do[*User](s.client,
		withPath("user"),
		withRequestOpts(options...),
	)
}

// UserStatus represents the current status of a user
//
// GitLab API docs:
// https://docs.gitlab.com/api/users/#get-your-user-status
type UserStatus struct {
	Emoji         string            `json:"emoji"`
	Availability  AvailabilityValue `json:"availability"`
	Message       string            `json:"message"`
	MessageHTML   string            `json:"message_html"`
	ClearStatusAt *time.Time        `json:"clear_status_at"`
}

func (s *UsersService) CurrentUserStatus(options ...RequestOptionFunc) (*UserStatus, *Response, error) {
	return do[*UserStatus](s.client,
		withPath("user/status"),
		withRequestOpts(options...),
	)
}

// GetUserStatus retrieves a user's status.
//
// uid can be either a user ID (int) or a username (string). If a username
// is provided with a leading "@" (e.g., "@johndoe"), it will be trimmed.
//
// GitLab API docs:
// https://docs.gitlab.com/api/users/#get-the-status-of-a-user
func (s *UsersService) GetUserStatus(uid any, options ...RequestOptionFunc) (*UserStatus, *Response, error) {
	return do[*UserStatus](s.client,
		withPath("users/%s/status", UserID{uid}),
		withRequestOpts(options...),
	)
}

// UserStatusOptions represents the options required to set the status
//
// GitLab API docs:
// https://docs.gitlab.com/api/users/#set-your-user-status
type UserStatusOptions struct {
	Emoji            *string                `url:"emoji,omitempty" json:"emoji,omitempty"`
	Availability     *AvailabilityValue     `url:"availability,omitempty" json:"availability,omitempty"`
	Message          *string                `url:"message,omitempty" json:"message,omitempty"`
	ClearStatusAfter *ClearStatusAfterValue `url:"clear_status_after,omitempty" json:"clear_status_after,omitempty"`
}

func (s *UsersService) SetUserStatus(opt *UserStatusOptions, options ...RequestOptionFunc) (*UserStatus, *Response, error) {
	return do[*UserStatus](s.client,
		withMethod(http.MethodPut),
		withPath("user/status"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UserAssociationsCount represents the user associations count.
//
// GitLab API docs:
// https://docs.gitlab.com/api/users/#get-a-count-of-a-users-projects-groups-issues-and-merge-requests
type UserAssociationsCount struct {
	GroupsCount        int64 `json:"groups_count"`
	ProjectsCount      int64 `json:"projects_count"`
	IssuesCount        int64 `json:"issues_count"`
	MergeRequestsCount int64 `json:"merge_requests_count"`
}

// GetUserAssociationsCount gets a list of a specified user associations.
//
// GitLab API docs:
// https://docs.gitlab.com/api/users/#get-a-count-of-a-users-projects-groups-issues-and-merge-requests
func (s *UsersService) GetUserAssociationsCount(user int64, options ...RequestOptionFunc) (*UserAssociationsCount, *Response, error) {
	return do[*UserAssociationsCount](s.client,
		withPath("users/%d/associations_count", user),
		withRequestOpts(options...),
	)
}

// SSHKey represents a SSH key.
//
// GitLab API docs: https://docs.gitlab.com/api/user_keys/#list-all-ssh-keys
type SSHKey struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Key       string     `json:"key"`
	CreatedAt *time.Time `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at"`
	UsageType string     `json:"usage_type"`
}

// ListSSHKeysOptions represents the available ListSSHKeys options.
//
// GitLab API docs: https://docs.gitlab.com/api/user_keys/#list-all-ssh-keys
type ListSSHKeysOptions struct {
	ListOptions
}

func (s *UsersService) ListSSHKeys(opt *ListSSHKeysOptions, options ...RequestOptionFunc) ([]*SSHKey, *Response, error) {
	return do[[]*SSHKey](s.client,
		withPath("user/keys"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListSSHKeysForUserOptions represents the available ListSSHKeysForUser() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/user_keys/#list-all-ssh-keys-for-a-user
type ListSSHKeysForUserOptions struct {
	ListOptions
}

// ListSSHKeysForUser gets a list of a specified user's SSH keys.
//
// uid can be either a user ID (int) or a username (string). If a username
// is provided with a leading "@" (e.g., "@johndoe"), it will be trimmed.
//
// GitLab API docs:
// https://docs.gitlab.com/api/user_keys/#list-all-ssh-keys-for-a-user
func (s *UsersService) ListSSHKeysForUser(uid any, opt *ListSSHKeysForUserOptions, options ...RequestOptionFunc) ([]*SSHKey, *Response, error) {
	return do[[]*SSHKey](s.client,
		withPath("users/%s/keys", UserID{uid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *UsersService) GetSSHKey(key int64, options ...RequestOptionFunc) (*SSHKey, *Response, error) {
	return do[*SSHKey](s.client,
		withPath("user/keys/%d", key),
		withRequestOpts(options...),
	)
}

func (s *UsersService) GetSSHKeyForUser(user int64, key int64, options ...RequestOptionFunc) (*SSHKey, *Response, error) {
	return do[*SSHKey](s.client,
		withPath("users/%d/keys/%d", user, key),
		withRequestOpts(options...),
	)
}

// AddSSHKeyOptions represents the available AddSSHKey() options.
//
// GitLab API docs: https://docs.gitlab.com/api/user_keys/#add-an-ssh-key
type AddSSHKeyOptions struct {
	Title     *string  `url:"title,omitempty" json:"title,omitempty"`
	Key       *string  `url:"key,omitempty" json:"key,omitempty"`
	ExpiresAt *ISOTime `url:"expires_at,omitempty" json:"expires_at,omitempty"`
	UsageType *string  `url:"usage_type,omitempty" json:"usage_type,omitempty"`
}

func (s *UsersService) AddSSHKey(opt *AddSSHKeyOptions, options ...RequestOptionFunc) (*SSHKey, *Response, error) {
	return do[*SSHKey](s.client,
		withMethod(http.MethodPost),
		withPath("user/keys"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *UsersService) AddSSHKeyForUser(user int64, opt *AddSSHKeyOptions, options ...RequestOptionFunc) (*SSHKey, *Response, error) {
	return do[*SSHKey](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/keys", user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *UsersService) DeleteSSHKey(key int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("user/keys/%d", key),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *UsersService) DeleteSSHKeyForUser(user, key int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("users/%d/keys/%d", user, key),
		withRequestOpts(options...),
	)
	return resp, err
}

// GPGKey represents a GPG key.
//
// GitLab API docs: https://docs.gitlab.com/api/user_keys/#list-all-gpg-keys
type GPGKey struct {
	ID        int64      `json:"id"`
	Key       string     `json:"key"`
	CreatedAt *time.Time `json:"created_at"`
}

func (s *UsersService) ListGPGKeys(options ...RequestOptionFunc) ([]*GPGKey, *Response, error) {
	return do[[]*GPGKey](s.client,
		withPath("user/gpg_keys"),
		withRequestOpts(options...),
	)
}

func (s *UsersService) GetGPGKey(key int64, options ...RequestOptionFunc) (*GPGKey, *Response, error) {
	return do[*GPGKey](s.client,
		withPath("user/gpg_keys/%d", key),
		withRequestOpts(options...),
	)
}

// AddGPGKeyOptions represents the available AddGPGKey() options.
//
// GitLab API docs: https://docs.gitlab.com/api/user_keys/#add-a-gpg-key
type AddGPGKeyOptions struct {
	Key *string `url:"key,omitempty" json:"key,omitempty"`
}

func (s *UsersService) AddGPGKey(opt *AddGPGKeyOptions, options ...RequestOptionFunc) (*GPGKey, *Response, error) {
	return do[*GPGKey](s.client,
		withMethod(http.MethodPost),
		withPath("user/gpg_keys"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *UsersService) DeleteGPGKey(key int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("user/gpg_keys/%d", key),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *UsersService) ListGPGKeysForUser(user int64, options ...RequestOptionFunc) ([]*GPGKey, *Response, error) {
	return do[[]*GPGKey](s.client,
		withPath("users/%d/gpg_keys", user),
		withRequestOpts(options...),
	)
}

func (s *UsersService) GetGPGKeyForUser(user, key int64, options ...RequestOptionFunc) (*GPGKey, *Response, error) {
	return do[*GPGKey](s.client,
		withPath("users/%d/gpg_keys/%d", user, key),
		withRequestOpts(options...),
	)
}

func (s *UsersService) AddGPGKeyForUser(user int64, opt *AddGPGKeyOptions, options ...RequestOptionFunc) (*GPGKey, *Response, error) {
	return do[*GPGKey](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/gpg_keys", user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *UsersService) DeleteGPGKeyForUser(user, key int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("users/%d/gpg_keys/%d", user, key),
		withRequestOpts(options...),
	)
	return resp, err
}

// Email represents an Email.
//
// GitLab API docs:
// https://docs.gitlab.com/api/user_email_addresses/#list-all-email-addresses
type Email struct {
	ID          int64      `json:"id"`
	Email       string     `json:"email"`
	ConfirmedAt *time.Time `json:"confirmed_at"`
}

func (s *UsersService) ListEmails(options ...RequestOptionFunc) ([]*Email, *Response, error) {
	return do[[]*Email](s.client,
		withPath("user/emails"),
		withRequestOpts(options...),
	)
}

// ListEmailsForUserOptions represents the available ListEmailsForUser() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/user_email_addresses/#list-all-email-addresses-for-a-user
type ListEmailsForUserOptions struct {
	ListOptions
}

func (s *UsersService) ListEmailsForUser(user int64, opt *ListEmailsForUserOptions, options ...RequestOptionFunc) ([]*Email, *Response, error) {
	return do[[]*Email](s.client,
		withPath("users/%d/emails", user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *UsersService) GetEmail(email int64, options ...RequestOptionFunc) (*Email, *Response, error) {
	return do[*Email](s.client,
		withPath("user/emails/%d", email),
		withRequestOpts(options...),
	)
}

// AddEmailOptions represents the available AddEmail() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/user_email_addresses/#add-an-email-address
type AddEmailOptions struct {
	Email            *string `url:"email,omitempty" json:"email,omitempty"`
	SkipConfirmation *bool   `url:"skip_confirmation,omitempty" json:"skip_confirmation,omitempty"`
}

func (s *UsersService) AddEmail(opt *AddEmailOptions, options ...RequestOptionFunc) (*Email, *Response, error) {
	return do[*Email](s.client,
		withMethod(http.MethodPost),
		withPath("user/emails"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *UsersService) AddEmailForUser(user int64, opt *AddEmailOptions, options ...RequestOptionFunc) (*Email, *Response, error) {
	return do[*Email](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/emails", user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *UsersService) DeleteEmail(email int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("user/emails/%d", email),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *UsersService) DeleteEmailForUser(user, email int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("users/%d/emails/%d", user, email),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *UsersService) BlockUser(user int64, options ...RequestOptionFunc) error {
	_, _, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/block", user),
		withRequestOpts(options...),
	)
	return err
}

func (s *UsersService) UnblockUser(user int64, options ...RequestOptionFunc) error {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/unblock", user),
		withRequestOpts(options...),
	)
	if err != nil && resp == nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusForbidden:
		return ErrUserUnblockPrevented
	case http.StatusNotFound:
		return ErrUserNotFound
	default:
		return fmt.Errorf("%w: %d", errUnexpectedResultCode, resp.StatusCode)
	}
}

func (s *UsersService) BanUser(user int64, options ...RequestOptionFunc) error {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/ban", user),
		withRequestOpts(options...),
	)
	if err != nil && resp == nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusNotFound:
		return ErrUserNotFound
	default:
		return fmt.Errorf("%w: %d", errUnexpectedResultCode, resp.StatusCode)
	}
}

func (s *UsersService) UnbanUser(user int64, options ...RequestOptionFunc) error {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/unban", user),
		withRequestOpts(options...),
	)
	if err != nil && resp == nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusNotFound:
		return ErrUserNotFound
	default:
		return fmt.Errorf("%w: %d", errUnexpectedResultCode, resp.StatusCode)
	}
}

func (s *UsersService) DeactivateUser(user int64, options ...RequestOptionFunc) error {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/deactivate", user),
		withRequestOpts(options...),
	)
	if err != nil && resp == nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusForbidden:
		return ErrUserDeactivatePrevented
	case http.StatusNotFound:
		return ErrUserNotFound
	default:
		return fmt.Errorf("%w: %d", errUnexpectedResultCode, resp.StatusCode)
	}
}

func (s *UsersService) ActivateUser(user int64, options ...RequestOptionFunc) error {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/activate", user),
		withRequestOpts(options...),
	)
	if err != nil && resp == nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusForbidden:
		return ErrUserActivatePrevented
	case http.StatusNotFound:
		return ErrUserNotFound
	default:
		return fmt.Errorf("%w: %d", errUnexpectedResultCode, resp.StatusCode)
	}
}

func (s *UsersService) ApproveUser(user int64, options ...RequestOptionFunc) error {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/approve", user),
		withRequestOpts(options...),
	)
	if err != nil && resp == nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusForbidden:
		return ErrUserApprovePrevented
	case http.StatusNotFound:
		return ErrUserNotFound
	default:
		return fmt.Errorf("%w: %d", errUnexpectedResultCode, resp.StatusCode)
	}
}

func (s *UsersService) RejectUser(user int64, options ...RequestOptionFunc) error {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/reject", user),
		withRequestOpts(options...),
	)
	if err != nil && resp == nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return nil
	case http.StatusForbidden:
		return ErrUserRejectPrevented
	case http.StatusNotFound:
		return ErrUserNotFound
	case http.StatusConflict:
		return ErrUserConflict
	default:
		return fmt.Errorf("%w: %d", errUnexpectedResultCode, resp.StatusCode)
	}
}

// ImpersonationToken represents an impersonation token.
//
// GitLab API docs:
// https://docs.gitlab.com/api/user_tokens/#list-all-impersonation-tokens-for-a-user
type ImpersonationToken struct {
	ID         int64      `json:"id"`
	Name       string     `json:"name"`
	Active     bool       `json:"active"`
	Token      string     `json:"token"`
	Scopes     []string   `json:"scopes"`
	Revoked    bool       `json:"revoked"`
	CreatedAt  *time.Time `json:"created_at"`
	ExpiresAt  *ISOTime   `json:"expires_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
}

// GetAllImpersonationTokensOptions represents the available
// GetAllImpersonationTokens() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/user_tokens/#list-all-impersonation-tokens-for-a-user
type GetAllImpersonationTokensOptions struct {
	ListOptions
	State *string `url:"state,omitempty" json:"state,omitempty"`
}

func (s *UsersService) GetAllImpersonationTokens(user int64, opt *GetAllImpersonationTokensOptions, options ...RequestOptionFunc) ([]*ImpersonationToken, *Response, error) {
	return do[[]*ImpersonationToken](s.client,
		withPath("users/%d/impersonation_tokens", user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *UsersService) GetImpersonationToken(user, token int64, options ...RequestOptionFunc) (*ImpersonationToken, *Response, error) {
	return do[*ImpersonationToken](s.client,
		withPath("users/%d/impersonation_tokens/%d", user, token),
		withRequestOpts(options...),
	)
}

// CreateImpersonationTokenOptions represents the available
// CreateImpersonationToken() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/user_tokens/#create-an-impersonation-token
type CreateImpersonationTokenOptions struct {
	Name      *string    `url:"name,omitempty" json:"name,omitempty"`
	Scopes    *[]string  `url:"scopes,omitempty" json:"scopes,omitempty"`
	ExpiresAt *time.Time `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

func (s *UsersService) CreateImpersonationToken(user int64, opt *CreateImpersonationTokenOptions, options ...RequestOptionFunc) (*ImpersonationToken, *Response, error) {
	return do[*ImpersonationToken](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/impersonation_tokens", user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *UsersService) RevokeImpersonationToken(user, token int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("users/%d/impersonation_tokens/%d", user, token),
		withRequestOpts(options...),
	)
	return resp, err
}

// CreatePersonalAccessTokenOptions represents the available
// CreatePersonalAccessToken() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/user_tokens/#create-a-personal-access-token-for-a-user
type CreatePersonalAccessTokenOptions struct {
	Name        *string   `url:"name,omitempty" json:"name,omitempty"`
	Description *string   `url:"description,omitempty" json:"description,omitempty"`
	ExpiresAt   *ISOTime  `url:"expires_at,omitempty" json:"expires_at,omitempty"`
	Scopes      *[]string `url:"scopes,omitempty" json:"scopes,omitempty"`
}

func (s *UsersService) CreatePersonalAccessToken(user int64, opt *CreatePersonalAccessTokenOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error) {
	return do[*PersonalAccessToken](s.client,
		withMethod(http.MethodPost),
		withPath("users/%d/personal_access_tokens", user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// CreatePersonalAccessTokenForCurrentUserOptions represents the available
// CreatePersonalAccessTokenForCurrentUser() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/user_tokens/#create-a-personal-access-token
type CreatePersonalAccessTokenForCurrentUserOptions struct {
	Name        *string   `url:"name,omitempty" json:"name,omitempty"`
	Description *string   `url:"description,omitempty" json:"description,omitempty"`
	Scopes      *[]string `url:"scopes,omitempty" json:"scopes,omitempty"`
	ExpiresAt   *ISOTime  `url:"expires_at,omitempty" json:"expires_at,omitempty"`
}

func (s *UsersService) CreatePersonalAccessTokenForCurrentUser(opt *CreatePersonalAccessTokenForCurrentUserOptions, options ...RequestOptionFunc) (*PersonalAccessToken, *Response, error) {
	return do[*PersonalAccessToken](s.client,
		withMethod(http.MethodPost),
		withPath("user/personal_access_tokens"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UserActivity represents an entry in the user/activities response
//
// GitLab API docs:
// https://docs.gitlab.com/api/users/#list-a-users-activity
type UserActivity struct {
	Username       string   `json:"username"`
	LastActivityOn *ISOTime `json:"last_activity_on"`
}

// GetUserActivitiesOptions represents the options for GetUserActivities
//
// GitLab API docs:
// https://docs.gitlab.com/api/users/#list-a-users-activity
type GetUserActivitiesOptions struct {
	ListOptions
	From *ISOTime `url:"from,omitempty" json:"from,omitempty"`
}

func (s *UsersService) GetUserActivities(opt *GetUserActivitiesOptions, options ...RequestOptionFunc) ([]*UserActivity, *Response, error) {
	return do[[]*UserActivity](s.client,
		withPath("user/activities"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UserMembership represents a membership of the user in a namespace or project.
//
// GitLab API docs:
// https://docs.gitlab.com/api/users/#list-projects-and-groups-that-a-user-is-a-member-of
type UserMembership struct {
	SourceID    int64            `json:"source_id"`
	SourceName  string           `json:"source_name"`
	SourceType  string           `json:"source_type"`
	AccessLevel AccessLevelValue `json:"access_level"`
}

// GetUserMembershipOptions represents the options available to query user memberships.
//
// GitLab API docs:
// https://docs.gitlab.com/api/users/#list-projects-and-groups-that-a-user-is-a-member-of
type GetUserMembershipOptions struct {
	ListOptions
	Type *string `url:"type,omitempty" json:"type,omitempty"`
}

func (s *UsersService) GetUserMemberships(user int64, opt *GetUserMembershipOptions, options ...RequestOptionFunc) ([]*UserMembership, *Response, error) {
	return do[[]*UserMembership](s.client,
		withPath("users/%d/memberships", user),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *UsersService) DisableTwoFactor(user int64, options ...RequestOptionFunc) error {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPatch),
		withPath("users/%d/disable_two_factor", user),
		withRequestOpts(options...),
	)
	if err != nil && resp == nil {
		return err
	}

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusBadRequest:
		return ErrUserTwoFactorNotEnabled
	case http.StatusForbidden:
		return ErrUserDisableTwoFactorPrevented
	case http.StatusNotFound:
		return ErrUserNotFound
	default:
		return fmt.Errorf("%w: %d", errUnexpectedResultCode, resp.StatusCode)
	}
}

// UserRunner represents a GitLab runner linked to the current user.
//
// GitLab API docs:
// https://docs.gitlab.com/api/users/#create-a-runner-linked-to-a-user
type UserRunner struct {
	ID             int64      `json:"id"`
	Token          string     `json:"token"`
	TokenExpiresAt *time.Time `json:"token_expires_at"`
}

// CreateUserRunnerOptions represents the available CreateUserRunner() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/users/#create-a-runner-linked-to-a-user
type CreateUserRunnerOptions struct {
	RunnerType      *string   `url:"runner_type,omitempty" json:"runner_type,omitempty"`
	GroupID         *int64    `url:"group_id,omitempty" json:"group_id,omitempty"`
	ProjectID       *int64    `url:"project_id,omitempty" json:"project_id,omitempty"`
	Description     *string   `url:"description,omitempty" json:"description,omitempty"`
	Paused          *bool     `url:"paused,omitempty" json:"paused,omitempty"`
	Locked          *bool     `url:"locked,omitempty" json:"locked,omitempty"`
	RunUntagged     *bool     `url:"run_untagged,omitempty" json:"run_untagged,omitempty"`
	TagList         *[]string `url:"tag_list,omitempty" json:"tag_list,omitempty"`
	AccessLevel     *string   `url:"access_level,omitempty" json:"access_level,omitempty"`
	MaximumTimeout  *int64    `url:"maximum_timeout,omitempty" json:"maximum_timeout,omitempty"`
	MaintenanceNote *string   `url:"maintenance_note,omitempty" json:"maintenance_note,omitempty"`
}

func (s *UsersService) CreateUserRunner(opts *CreateUserRunnerOptions, options ...RequestOptionFunc) (*UserRunner, *Response, error) {
	return do[*UserRunner](s.client,
		withMethod(http.MethodPost),
		withPath("user/runners"),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

// CreateServiceAccountUserOptions represents the available CreateServiceAccountUser() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/user_service_accounts/#create-a-service-account-user
type CreateServiceAccountUserOptions struct {
	Name     *string `url:"name,omitempty" json:"name,omitempty"`
	Username *string `url:"username,omitempty" json:"username,omitempty"`
	Email    *string `url:"email,omitempty" json:"email,omitempty"`
}

func (s *UsersService) CreateServiceAccountUser(opts *CreateServiceAccountUserOptions, options ...RequestOptionFunc) (*User, *Response, error) {
	return do[*User](s.client,
		withMethod(http.MethodPost),
		withPath("service_accounts"),
		withAPIOpts(opts),
		withRequestOpts(options...),
	)
}

func (s *UsersService) ListServiceAccounts(opt *ListServiceAccountsOptions, options ...RequestOptionFunc) ([]*ServiceAccount, *Response, error) {
	return do[[]*ServiceAccount](s.client,
		withPath("service_accounts"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *UsersService) UploadAvatar(avatar io.Reader, filename string, options ...RequestOptionFunc) (*User, *Response, error) {
	return do[*User](s.client,
		withMethod(http.MethodPut),
		withPath("user/avatar"),
		withUpload(avatar, filename, UploadAvatar),
		withRequestOpts(options...),
	)
}

func (s *UsersService) DeleteUserIdentity(user int64, provider string, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("users/%d/identities/%s", user, provider),
		withRequestOpts(options...),
	)
	return resp, err
}

// userCoreBasicTemplate defines the common fields for a user in GraphQL queries.
var userCoreBasicTemplate = template.Must(template.New("UserCoreBasic").Parse(`
	id
	username
	name
	state
	createdAt
	avatarUrl
	webUrl
`))

// userCoreBasicGQL represents the UserCore GraphQL type. It unwraps to a *BasicUser type.
type userCoreBasicGQL struct {
	ID        gidGQL     `json:"id"`
	Username  string     `json:"username"`
	Name      string     `json:"name"`
	State     string     `json:"state"`
	CreatedAt *time.Time `json:"createdAt"`
	AvatarURL string     `json:"avatarUrl"`
	WebURL    string     `json:"webUrl"`
}

// unwrap converts the GraphQL data structure to a *BasicUser.
func (u userCoreBasicGQL) unwrap() *BasicUser {
	if u.Username == "" {
		return nil
	}

	return &BasicUser{
		ID:        u.ID.Int64,
		Username:  u.Username,
		Name:      u.Name,
		State:     u.State,
		Locked:    u.State != "active",
		CreatedAt: u.CreatedAt,
		AvatarURL: u.AvatarURL,
		WebURL:    u.WebURL,
	}
}
