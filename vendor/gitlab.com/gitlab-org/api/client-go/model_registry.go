package gitlab

import (
	"bytes"
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
	// The modelVersionID can be an int or a string like "candidate:5",
	// so we convert it to a string for the URL.
	mvid, err := parseID(modelVersionID)
	if err != nil {
		return nil, nil, err
	}

	buf, resp, err := do[bytes.Buffer](s.client,
		withPath("projects/%s/packages/ml_models/%s/files/%s/%s",
			ProjectID{pid},
			// the following URI components must not escape `.` which is what withPath does by default
			// without NoEscape.
			NoEscape{url.PathEscape(mvid)},
			NoEscape{url.PathEscape(path)},
			NoEscape{url.PathEscape(filename)},
		),
		withRequestOpts(options...),
	)
	if err != nil {
		return nil, resp, err
	}
	return bytes.NewReader(buf.Bytes()), resp, nil
}
