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
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

type (
	GroupImportExportServiceInterface interface {
		ScheduleExport(gid any, options ...RequestOptionFunc) (*Response, error)
		ExportDownload(gid any, options ...RequestOptionFunc) (*bytes.Reader, *Response, error)
		ImportFile(opt *GroupImportFileOptions, options ...RequestOptionFunc) (*Response, error)
	}

	// GroupImportExportService handles communication with the group import export
	// related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/group_import_export/
	GroupImportExportService struct {
		client *Client
	}
)

var _ GroupImportExportServiceInterface = (*GroupImportExportService)(nil)

// ScheduleExport starts a new group export.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_import_export/#schedule-new-export
func (s *GroupImportExportService) ScheduleExport(gid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodPost),
		withPath("groups/%s/export", GroupID{gid}),
		withAPIOpts(nil),
		withRequestOpts(options...),
	)
	return resp, err
}

// ExportDownload downloads the finished export.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_import_export/#export-download
func (s *GroupImportExportService) ExportDownload(gid any, options ...RequestOptionFunc) (*bytes.Reader, *Response, error) {
	buf, resp, err := do[bytes.Buffer](s.client,
		withPath("groups/%s/export/download", GroupID{gid}),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return bytes.NewReader(buf.Bytes()), resp, nil
}

// GroupImportFileOptions represents the available ImportFile() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_import_export/#import-a-file
type GroupImportFileOptions struct {
	Name     *string `url:"name,omitempty" json:"name,omitempty"`
	Path     *string `url:"path,omitempty" json:"path,omitempty"`
	File     *string `url:"file,omitempty" json:"file,omitempty"`
	ParentID *int64  `url:"parent_id,omitempty" json:"parent_id,omitempty"`
}

// ImportFile imports a file.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_import_export/#import-a-file
func (s *GroupImportExportService) ImportFile(opt *GroupImportFileOptions, options ...RequestOptionFunc) (*Response, error) {
	// First check if we got all required options.
	if opt.Name == nil || *opt.Name == "" {
		return nil, errors.New("missing required option: Name")
	}
	if opt.Path == nil || *opt.Path == "" {
		return nil, errors.New("missing required option: Path")
	}
	if opt.File == nil || *opt.File == "" {
		return nil, errors.New("missing required option: File")
	}

	f, err := os.Open(*opt.File)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b := &bytes.Buffer{}
	w := multipart.NewWriter(b)

	_, filename := filepath.Split(*opt.File)
	fw, err := w.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return nil, err
	}

	// Populate the additional fields.
	fw, err = w.CreateFormField("name")
	if err != nil {
		return nil, err
	}

	_, err = fw.Write([]byte(*opt.Name))
	if err != nil {
		return nil, err
	}

	fw, err = w.CreateFormField("path")
	if err != nil {
		return nil, err
	}

	_, err = fw.Write([]byte(*opt.Path))
	if err != nil {
		return nil, err
	}

	if opt.ParentID != nil {
		fw, err = w.CreateFormField("parent_id")
		if err != nil {
			return nil, err
		}

		_, err = fw.Write([]byte(strconv.FormatInt(*opt.ParentID, 10)))
		if err != nil {
			return nil, err
		}
	}

	if err = w.Close(); err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest(http.MethodPost, "groups/import", nil, options)
	if err != nil {
		return nil, err
	}

	// Set the buffer as the request body.
	if err = req.SetBody(b); err != nil {
		return nil, err
	}

	// Overwrite the default content type.
	req.Header.Set("Content-Type", w.FormDataContentType())

	return s.client.Do(req, nil)
}
