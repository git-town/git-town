package gitlab

import (
	"net/http"
	"time"
)

type (
	// GroupSSHCertificatesServiceInterface defines methods for the GroupSSHCertificatesService.
	GroupSSHCertificatesServiceInterface interface {
		ListGroupSSHCertificates(gid any, options ...RequestOptionFunc) ([]*GroupSSHCertificate, *Response, error)
		CreateGroupSSHCertificate(gid any, opt *CreateGroupSSHCertificateOptions, options ...RequestOptionFunc) (*GroupSSHCertificate, *Response, error)
		DeleteGroupSSHCertificate(gid any, cert int64, options ...RequestOptionFunc) (*Response, error)
	}

	// GroupSSHCertificatesService handles communication with the group
	// SSH certificate related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/group_ssh_certificates/
	GroupSSHCertificatesService struct {
		client *Client
	}
)

var _ GroupSSHCertificatesServiceInterface = (*GroupSSHCertificatesService)(nil)

// GroupSSHCertificate represents a GitLab Group SSH certificate.
//
// GitLab API docs: https://docs.gitlab.com/api/group_ssh_certificates/
type GroupSSHCertificate struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Key       string     `json:"key"`
	CreatedAt *time.Time `json:"created_at"`
}

// ListGroupSSHCertificates gets a list of SSH certificates for a specified
// group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_ssh_certificates/#get-all-ssh-certificates-for-a-particular-group
func (s *GroupSSHCertificatesService) ListGroupSSHCertificates(gid any, options ...RequestOptionFunc) ([]*GroupSSHCertificate, *Response, error) {
	return do[[]*GroupSSHCertificate](s.client,
		withPath("groups/%s/ssh_certificates", GroupID{gid}),
		withRequestOpts(options...),
	)
}

// CreateGroupSSHCertificateOptions represents the available
// CreateGroupSSHCertificate() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_ssh_certificates/#create-ssh-certificate
type CreateGroupSSHCertificateOptions struct {
	Key   *string `url:"key,omitempty" json:"key,omitempty"`
	Title *string `url:"title,omitempty" json:"title,omitempty"`
}

// CreateGroupSSHCertificate creates a new SSH certificate in the group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_ssh_certificates/#create-ssh-certificate
func (s *GroupSSHCertificatesService) CreateGroupSSHCertificate(gid any, opt *CreateGroupSSHCertificateOptions, options ...RequestOptionFunc) (*GroupSSHCertificate, *Response, error) {
	return do[*GroupSSHCertificate](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/ssh_certificates", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteGroupSSHCertificate deletes a SSH certificate from a specified group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_ssh_certificates/#delete-group-ssh-certificate
func (s *GroupSSHCertificatesService) DeleteGroupSSHCertificate(gid any, cert int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("groups/%s/ssh_certificates/%d", GroupID{gid}, cert),
		withRequestOpts(options...),
	)
	return resp, err
}
