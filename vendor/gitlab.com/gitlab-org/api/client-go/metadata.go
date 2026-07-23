//
// Copyright 2022, Timo Furrer <tuxtimo@gmail.com>
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

type (
	MetadataServiceInterface interface {
		GetMetadata(options ...RequestOptionFunc) (*Metadata, *Response, error)
	}

	// MetadataService handles communication with the GitLab server instance to
	// retrieve its metadata information via the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/metadata/
	MetadataService struct {
		client *Client
	}
)

var _ MetadataServiceInterface = (*MetadataService)(nil)

// Metadata represents a GitLab instance version.
//
// GitLab API docs: https://docs.gitlab.com/api/metadata/
type Metadata struct {
	Version    string      `json:"version"`
	Revision   string      `json:"revision"`
	KAS        MetadataKAS `json:"kas"`
	Enterprise bool        `json:"enterprise"`
}

func (s Metadata) String() string {
	return Stringify(s)
}

// MetadataKAS represents a GitLab instance version metadata KAS.
//
// GitLab API docs: https://docs.gitlab.com/api/metadata/
type MetadataKAS struct {
	Enabled             bool   `json:"enabled"`
	ExternalURL         string `json:"externalUrl"`
	ExternalK8SProxyURL string `json:"externalK8sProxyUrl"`
	Version             string `json:"version"`
}

func (k MetadataKAS) String() string {
	return Stringify(k)
}

// GetMetadata gets a GitLab server instance meteadata.
//
// GitLab API docs: https://docs.gitlab.com/api/metadata/
func (s *MetadataService) GetMetadata(options ...RequestOptionFunc) (*Metadata, *Response, error) {
	return do[*Metadata](s.client,
		withPath("metadata"),
		withRequestOpts(options...),
	)
}
