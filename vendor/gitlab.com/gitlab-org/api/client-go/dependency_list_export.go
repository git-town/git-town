package gitlab

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
)

type (
	DependencyListExportServiceInterface interface {
		CreateDependencyListExport(pipelineID int, opt *CreateDependencyListExportOptions, options ...RequestOptionFunc) (*DependencyListExport, *Response, error)
		GetDependencyListExport(id int, options ...RequestOptionFunc) (*DependencyListExport, *Response, error)
		DownloadDependencyListExport(id int, options ...RequestOptionFunc) (io.Reader, *Response, error)
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
	ID          int    `json:"id"`
	HasFinished bool   `json:"has_finished"`
	Self        string `json:"self"`
	Download    string `json:"download"`
}

const defaultExportType = "sbom"

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
func (s *DependencyListExportService) CreateDependencyListExport(pipelineID int, opt *CreateDependencyListExportOptions, options ...RequestOptionFunc) (*DependencyListExport, *Response, error) {
	// POST /pipelines/:id/dependency_list_exports
	createExportPath := fmt.Sprintf("pipelines/%d/dependency_list_exports", pipelineID)

	if opt == nil {
		opt = &CreateDependencyListExportOptions{}
	}
	if opt.ExportType == nil {
		opt.ExportType = Ptr(defaultExportType)
	}

	req, err := s.client.NewRequest(http.MethodPost, createExportPath, opt, options)
	if err != nil {
		return nil, nil, err
	}

	export := new(DependencyListExport)
	resp, err := s.client.Do(req, &export)
	if err != nil {
		return nil, resp, err
	}

	return export, resp, nil
}

// GetDependencyListExport gets metadata about a single dependency list export.
//
// GitLab docs:
// https://docs.gitlab.com/api/dependency_list_export/#get-single-dependency-list-export
func (s *DependencyListExportService) GetDependencyListExport(id int, options ...RequestOptionFunc) (*DependencyListExport, *Response, error) {
	// GET /dependency_list_exports/:id
	getExportPath := fmt.Sprintf("dependency_list_exports/%d", id)

	req, err := s.client.NewRequest(http.MethodGet, getExportPath, nil, options)
	if err != nil {
		return nil, nil, err
	}

	export := new(DependencyListExport)
	resp, err := s.client.Do(req, &export)
	if err != nil {
		return nil, resp, err
	}

	return export, resp, nil
}

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
func (s *DependencyListExportService) DownloadDependencyListExport(id int, options ...RequestOptionFunc) (io.Reader, *Response, error) {
	// GET /dependency_list_exports/:id/download
	downloadExportPath := fmt.Sprintf("dependency_list_exports/%d/download", id)

	req, err := s.client.NewRequest(http.MethodGet, downloadExportPath, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var sbomBuffer bytes.Buffer
	resp, err := s.client.Do(req, &sbomBuffer)
	if err != nil {
		return nil, resp, err
	}

	return &sbomBuffer, resp, nil
}
