package gitlab

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
)

type (
	ModelRegistryServiceInterface interface {
		DownloadMachineLearningModelPackage(pid, modelVersionID any, path string, filename string, options ...RequestOptionFunc) (*bytes.Reader, *Response, error)
	}

	// ModelRegistryService handles communication with the model registry related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/model_registry/
	ModelRegistryService struct {
		client *Client
	}
)

var _ ModelRegistryServiceInterface = (*ModelRegistryService)(nil)

// DownloadMachineLearningModelPackage downloads a machine learning model package file.
//
// GitLab API docs: https://docs.gitlab.com/api/model_registry/#download-a-model-package-file
func (s *ModelRegistryService) DownloadMachineLearningModelPackage(pid, modelVersionID any, path string, filename string, options ...RequestOptionFunc) (*bytes.Reader, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}

	// The modelVersionID can be an int or a string like "candidate:5",
	// so we convert it to a string for the URL.
	mvid, err := parseID(modelVersionID)
	if err != nil {
		return nil, nil, err
	}

	u := fmt.Sprintf("projects/%s/packages/ml_models/%s/files/%s/%s",
		PathEscape(project),
		url.PathEscape(mvid),
		url.PathEscape(path),
		url.PathEscape(filename),
	)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	registryDownload := new(bytes.Buffer)
	resp, err := s.client.Do(req, registryDownload)
	if err != nil {
		return nil, resp, err
	}

	return bytes.NewReader(registryDownload.Bytes()), resp, err
}
