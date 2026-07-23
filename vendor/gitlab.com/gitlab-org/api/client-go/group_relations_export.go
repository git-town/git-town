package gitlab

import (
	"bytes"
	"net/http"
	"time"
)

// GroupRelationsScheduleExportOptions represents the available ScheduleExport() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_relations_export/#schedule-new-export
type GroupRelationsScheduleExportOptions struct {
	Batched *bool `url:"batched,omitempty" json:"batched,omitempty"`
}

// ListGroupRelationsStatusOptions represents the available ListExportStatus() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_relations_export/#export-status
type ListGroupRelationsStatusOptions struct {
	ListOptions

	Relation *string `url:"relation,omitempty" json:"relation,omitempty"`
}

// GroupRelationsDownloadOptions represents the available ExportDownload() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_relations_export/#download-exported-relations
type GroupRelationsDownloadOptions struct {
	Relation    *string `url:"relation,omitempty" json:"relation,omitempty"`
	Batched     *bool   `url:"batched,omitempty" json:"batched,omitempty"`
	BatchNumber *int64  `url:"batch_number,omitempty" json:"batch_number,omitempty"`
}

type GroupRelationStatus struct {
	Relation     string    `json:"relation"`
	Status       int64     `json:"status"`
	Error        string    `json:"error"`
	UpdatedAt    time.Time `json:"updated_at"`
	Batched      bool      `json:"batched"`
	BatchesCount int64     `json:"batches_count"`
	Batches      []Batch   `json:"batches,omitempty"`
}

type Batch struct {
	Status       int64     `json:"status"`
	BatchNumber  int64     `json:"batch_number"`
	ObjectsCount int64     `json:"objects_count"`
	Error        string    `json:"error"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type (
	GroupRelationsExportServiceInterface interface {
		// ScheduleExport schedules a new export of group relations.
		//
		// GitLab API docs: https://docs.gitlab.com/api/group_relations_export/#schedule-new-export
		ScheduleExport(gid any, opt *GroupRelationsScheduleExportOptions, options ...RequestOptionFunc) (*Response, error)
		// ListExportStatus gets the status of group relations export.
		//
		// GitLab API docs: https://docs.gitlab.com/api/group_relations_export/#export-status
		ListExportStatus(gid any, opt *ListGroupRelationsStatusOptions, options ...RequestOptionFunc) ([]*GroupRelationStatus, *Response, error)
		// ExportDownload downloads the exported group relations.
		//
		// GitLab API docs: https://docs.gitlab.com/api/group_relations_export/#download-exported-relations
		ExportDownload(gid any, opt *GroupRelationsDownloadOptions, options ...RequestOptionFunc) (*bytes.Reader, *Response, error)
	}

	// GroupRelationsExportService handles communication with the group relations export related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/group_relations_export
	GroupRelationsExportService struct {
		client *Client
	}
)

var _ GroupRelationsExportServiceInterface = (*GroupRelationsExportService)(nil)

func (s *GroupRelationsExportService) ScheduleExport(gid any, opt *GroupRelationsScheduleExportOptions, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/export_relations", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *GroupRelationsExportService) ListExportStatus(gid any, opt *ListGroupRelationsStatusOptions, options ...RequestOptionFunc) ([]*GroupRelationStatus, *Response, error) {
	return do[[]*GroupRelationStatus](s.client,
		withPath("groups/%s/export_relations/status", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *GroupRelationsExportService) ExportDownload(gid any, opt *GroupRelationsDownloadOptions, options ...RequestOptionFunc) (*bytes.Reader, *Response, error) {
	buf, resp, err := do[bytes.Buffer](s.client,
		withPath("groups/%s/export_relations/download", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return bytes.NewReader(buf.Bytes()), resp, nil
}
