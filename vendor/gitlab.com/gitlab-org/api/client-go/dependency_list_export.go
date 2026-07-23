package gitlab

import (
	"bytes"
	"io"
	"net/http"
)

type (
	DependencyListExportServiceInterface interface {
		// CreateDependencyListExport creates a new CycloneDX JSON export for all the project dependencies
		// detected in a pipeline.
		//
		// If an authenticated user does not have permission to read_dependency, this request returns a 403
		// Forbidden status code.
		//
		// SBOM exports can be only accessed by the export’s author.
		//
		// GitLab docs:
		// https://docs.gitlab.com/api/dependency_list_export/#create-a-dependency-list-export
		CreateDependencyListExport(pipelineID int64, opt *CreateDependencyListExportOptions, options ...RequestOptionFunc) (*DependencyListExport, *Response, error)

		// GetDependencyListExport gets metadata about a single dependency list export.
		//
		// GitLab docs:
		// https://docs.gitlab.com/api/dependency_list_export/#get-single-dependency-list-export
		GetDependencyListExport(id int64, options ...RequestOptionFunc) (*DependencyListExport, *Response, error)

		// DownloadDependencyListExport downloads a single dependency list export.
		//
		// The github.com/CycloneDX/cyclonedx-go package can be used to parse the data from the returned io.Reader.
		//
		//	sbom := new(cdx.BOM)
		//	decoder := cdx.NewBOMDecoder(reader, cdx.BOMFileFormatJSON)
		//
		//	if err = decoder.Decode(sbom); err != nil {
		//		panic(err)
		//	}
		//
		// GitLab docs:
		// https://docs.gitlab.com/api/dependency_list_export/#download-dependency-list-export
		DownloadDependencyListExport(id int64, options ...RequestOptionFunc) (io.Reader, *Response, error)
	}

	// DependencyListExportService handles communication with the dependency list export
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/dependency_list_export/
	DependencyListExportService struct {
		client *Client
	}
)

var _ DependencyListExportServiceInterface = (*DependencyListExportService)(nil)

// CreateDependencyListExportOptions represents the available CreateDependencyListExport()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/dependency_list_export/#create-a-dependency-list-export
type CreateDependencyListExportOptions struct {
	ExportType *string `url:"export_type" json:"export_type"`
}

// DependencyListExport represents a request for a GitLab project's dependency list.
//
// GitLab API docs:
// https://docs.gitlab.com/api/dependency_list_export/#create-a-dependency-list-export
type DependencyListExport struct {
	ID          int64  `json:"id"`
	HasFinished bool   `json:"has_finished"`
	Self        string `json:"self"`
	Download    string `json:"download"`
}

const defaultExportType = "sbom"

func (s *DependencyListExportService) CreateDependencyListExport(pipelineID int64, opt *CreateDependencyListExportOptions, options ...RequestOptionFunc) (*DependencyListExport, *Response, error) {
	if opt == nil {
		opt = &CreateDependencyListExportOptions{}
	}
	if opt.ExportType == nil {
		opt.ExportType = Ptr(defaultExportType)
	}

	return do[*DependencyListExport](s.client,
		withMethod(http.MethodPost),
		withPath("pipelines/%d/dependency_list_exports", pipelineID),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *DependencyListExportService) GetDependencyListExport(id int64, options ...RequestOptionFunc) (*DependencyListExport, *Response, error) {
	return do[*DependencyListExport](s.client,
		withPath("dependency_list_exports/%d", id),
		withRequestOpts(options...),
	)
}

func (s *DependencyListExportService) DownloadDependencyListExport(id int64, options ...RequestOptionFunc) (io.Reader, *Response, error) {
	buf, resp, err := do[bytes.Buffer](s.client,
		withPath("dependency_list_exports/%d/download", id),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return &buf, resp, nil
}
