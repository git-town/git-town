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
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	ProjectImportExportServiceInterface interface {
		ScheduleExport(pid any, opt *ScheduleExportOptions, options ...RequestOptionFunc) (*Response, error)
		ExportStatus(pid any, options ...RequestOptionFunc) (*ExportStatus, *Response, error)
		ExportDownload(pid any, options ...RequestOptionFunc) ([]byte, *Response, error)
		ImportFromFile(archive io.Reader, opt *ImportFileOptions, options ...RequestOptionFunc) (*ImportStatus, *Response, error)
		ImportStatus(pid any, options ...RequestOptionFunc) (*ImportStatus, *Response, error)
	}

	// ProjectImportExportService handles communication with the project
	// import/export related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/project_import_export/
	ProjectImportExportService struct {
		client *Client
	}
)

var _ ProjectImportExportServiceInterface = (*ProjectImportExportService)(nil)

// ImportStatus represents a project import status.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_import_export/#import-status
type ImportStatus struct {
	ID                int        `json:"id"`
	Description       string     `json:"description"`
	Name              string     `json:"name"`
	NameWithNamespace string     `json:"name_with_namespace"`
	Path              string     `json:"path"`
	PathWithNamespace string     `json:"path_with_namespace"`
	CreateAt          *time.Time `json:"create_at"`
	ImportStatus      string     `json:"import_status"`
	ImportType        string     `json:"import_type"`
	CorrelationID     string     `json:"correlation_id"`
	ImportError       string     `json:"import_error"`
}

func (s ImportStatus) String() string {
	return Stringify(s)
}

// ExportStatus represents a project export status.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_import_export/#export-status
type ExportStatus struct {
	ID                int        `json:"id"`
	Description       string     `json:"description"`
	Name              string     `json:"name"`
	NameWithNamespace string     `json:"name_with_namespace"`
	Path              string     `json:"path"`
	PathWithNamespace string     `json:"path_with_namespace"`
	CreatedAt         *time.Time `json:"created_at"`
	ExportStatus      string     `json:"export_status"`
	Message           string     `json:"message"`
	Links             struct {
		APIURL string `json:"api_url"`
		WebURL string `json:"web_url"`
	} `json:"_links"`
}

func (s ExportStatus) String() string {
	return Stringify(s)
}

// ScheduleExportOptions represents the available ScheduleExport() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_import_export/#schedule-an-export
type ScheduleExportOptions struct {
	Description *string `url:"description,omitempty" json:"description,omitempty"`
	Upload      struct {
		URL        *string `url:"url,omitempty" json:"url,omitempty"`
		HTTPMethod *string `url:"http_method,omitempty" json:"http_method,omitempty"`
	} `url:"upload,omitempty" json:"upload,omitempty"`
}

// ScheduleExport schedules a project export.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_import_export/#schedule-an-export
func (s *ProjectImportExportService) ScheduleExport(pid any, opt *ScheduleExportOptions, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/export", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ExportStatus get the status of export.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_import_export/#export-status
func (s *ProjectImportExportService) ExportStatus(pid any, options ...RequestOptionFunc) (*ExportStatus, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/export", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	es := new(ExportStatus)
	resp, err := s.client.Do(req, es)
	if err != nil {
		return nil, resp, err
	}

	return es, resp, nil
}

// ExportDownload download the finished export.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_import_export/#export-download
func (s *ProjectImportExportService) ExportDownload(pid any, options ...RequestOptionFunc) ([]byte, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/export/download", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var b bytes.Buffer
	resp, err := s.client.Do(req, &b)
	if err != nil {
		return nil, resp, err
	}

	return b.Bytes(), resp, err
}

// ImportFileOptions represents the available ImportFile() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_import_export/#import-a-file
type ImportFileOptions struct {
	Namespace      *string               `url:"namespace,omitempty" json:"namespace,omitempty"`
	Name           *string               `url:"name,omitempty" json:"name,omitempty"`
	Path           *string               `url:"path,omitempty" json:"path,omitempty"`
	Overwrite      *bool                 `url:"overwrite,omitempty" json:"overwrite,omitempty"`
	OverrideParams *CreateProjectOptions `url:"override_params,omitempty" json:"override_params,omitempty"`
}

// Import a project from an archive file.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_import_export/#import-a-file
func (s *ProjectImportExportService) ImportFromFile(archive io.Reader, opt *ImportFileOptions, options ...RequestOptionFunc) (*ImportStatus, *Response, error) {
	req, err := s.client.UploadRequest(
		http.MethodPost,
		"projects/import",
		archive,
		"archive.tar.gz",
		UploadFile,
		opt,
		options,
	)
	if err != nil {
		return nil, nil, err
	}

	is := new(ImportStatus)
	resp, err := s.client.Do(req, is)
	if err != nil {
		return nil, resp, err
	}

	return is, resp, nil
}

// ImportStatus get the status of an import.
//
// GitLab API docs:
// https://docs.gitlab.com/api/project_import_export/#import-status
func (s *ProjectImportExportService) ImportStatus(pid any, options ...RequestOptionFunc) (*ImportStatus, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/import", PathEscape(project))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	is := new(ImportStatus)
	resp, err := s.client.Do(req, is)
	if err != nil {
		return nil, resp, err
	}

	return is, resp, nil
}
