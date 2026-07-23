package gitlab

import (
	"bytes"
	"net/http"
	"time"
)

type (
	AttestationsServiceInterface interface {
		// ListAttestations gets a list of all attestations
		//
		// GitLab API docs: https://docs.gitlab.com/api/attestations/#list-all-attestations
		ListAttestations(pid any, subjectDigest string, options ...RequestOptionFunc) ([]*Attestation, *Response, error)

		// DownloadAttestation
		//
		// GitLab API docs: https://docs.gitlab.com/api/attestations/#download-an-attestation
		DownloadAttestation(pid any, attestationIID int64, options ...RequestOptionFunc) ([]byte, *Response, error)
	}

	// AttestationsService handles communication with the keys related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/attestations
	AttestationsService struct {
		client *Client
	}
)

var _ AttestationsServiceInterface = (*AttestationsService)(nil)

type Attestation struct {
	ID            int64      `json:"id"`
	IID           int64      `json:"iid"`
	ProjectID     int64      `json:"project_id"`
	BuildID       int64      `json:"build_id"`
	Status        string     `json:"status"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	ExpireAt      *time.Time `json:"expire_at"`
	PredicateKind string     `json:"predicate_kind"`
	PredicateType string     `json:"predicate_type"`
	SubjectDigest string     `json:"subject_digest"`
	DownloadURL   string     `json:"download_url"`
}

func (s *AttestationsService) ListAttestations(pid any, subjectDigest string, options ...RequestOptionFunc) ([]*Attestation, *Response, error) {
	return do[[]*Attestation](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/attestations/%s", ProjectID{pid}, subjectDigest),
		withRequestOpts(options...),
	)
}

func (s *AttestationsService) DownloadAttestation(pid any, attestationIID int64, options ...RequestOptionFunc) ([]byte, *Response, error) {
	b, resp, err := do[bytes.Buffer](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/attestations/%d/download", ProjectID{pid}, attestationIID),
		withRequestOpts(options...),
	)

	return b.Bytes(), resp, err
}
