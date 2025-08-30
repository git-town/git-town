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
	"fmt"
	"io"
	"net/http"
	"time"
)

type (
	AlertManagementServiceInterface interface {
		UploadMetricImage(pid any, alertIID int, content io.Reader, filename string, opt *UploadMetricImageOptions, options ...RequestOptionFunc) (*MetricImage, *Response, error)
		ListMetricImages(pid any, alertIID int, opt *ListMetricImagesOptions, options ...RequestOptionFunc) ([]*MetricImage, *Response, error)
		UpdateMetricImage(pid any, alertIID int, id int, opt *UpdateMetricImageOptions, options ...RequestOptionFunc) (*MetricImage, *Response, error)
		DeleteMetricImage(pid any, alertIID int, id int, options ...RequestOptionFunc) (*Response, error)
	}

	// AlertManagementService handles communication with the alert management
	// related methods of the GitLab API.
	//
	// GitLab API docs:
	// https://docs.gitlab.com/api/alert_management_alerts/
	AlertManagementService struct {
		client *Client
	}
)

var _ AlertManagementServiceInterface = (*AlertManagementService)(nil)

// MetricImage represents a single metric image file.
//
// GitLab API docs:
// https://docs.gitlab.com/api/alert_management_alerts/
type MetricImage struct {
	ID        int        `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	Filename  string     `json:"filename"`
	FilePath  string     `json:"file_path"`
	URL       string     `json:"url"`
	URLText   string     `json:"url_text"`
}

// UploadMetricImageOptions represents the available UploadMetricImage() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/alert_management_alerts/#upload-metric-image
type UploadMetricImageOptions struct {
	URL     *string `url:"url,omitempty" json:"url,omitempty"`
	URLText *string `url:"url_text,omitempty" json:"url_text,omitempty"`
}

// UploadMetricImageOptions uploads a metric image to a project alert.
//
// GitLab API docs:
// https://docs.gitlab.com/api/alert_management_alerts/#upload-metric-image
func (s *AlertManagementService) UploadMetricImage(pid any, alertIID int, content io.Reader, filename string, opt *UploadMetricImageOptions, options ...RequestOptionFunc) (*MetricImage, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/alert_management_alerts/%d/metric_images", PathEscape(project), alertIID)

	req, err := s.client.UploadRequest(http.MethodPost, u, content, filename, UploadFile, opt, options)
	if err != nil {
		return nil, nil, err
	}

	mi := new(MetricImage)
	resp, err := s.client.Do(req, mi)
	if err != nil {
		return nil, resp, err
	}

	return mi, resp, nil
}

// ListMetricImagesOptions represents the available ListMetricImages() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/alert_management_alerts/#list-metric-images
type ListMetricImagesOptions struct {
	ListOptions
}

// ListMetricImages lists all the metric images for a project alert.
//
// GitLab API docs:
// https://docs.gitlab.com/api/alert_management_alerts/#list-metric-images
func (s *AlertManagementService) ListMetricImages(pid any, alertIID int, opt *ListMetricImagesOptions, options ...RequestOptionFunc) ([]*MetricImage, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/alert_management_alerts/%d/metric_images", PathEscape(project), alertIID)

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var mis []*MetricImage
	resp, err := s.client.Do(req, &mis)
	if err != nil {
		return nil, resp, err
	}

	return mis, resp, nil
}

// UpdateMetricImageOptions represents the available UpdateMetricImage() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/alert_management_alerts/#update-metric-image
type UpdateMetricImageOptions struct {
	URL     *string `url:"url,omitempty" json:"url,omitempty"`
	URLText *string `url:"url_text,omitempty" json:"url_text,omitempty"`
}

// UpdateMetricImage updates a metric image for a project alert.
//
// GitLab API docs:
// https://docs.gitlab.com/api/alert_management_alerts/#update-metric-image
func (s *AlertManagementService) UpdateMetricImage(pid any, alertIID int, id int, opt *UpdateMetricImageOptions, options ...RequestOptionFunc) (*MetricImage, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/alert_management_alerts/%d/metric_images/%d", PathEscape(project), alertIID, id)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	mi := new(MetricImage)
	resp, err := s.client.Do(req, mi)
	if err != nil {
		return nil, resp, err
	}

	return mi, resp, nil
}

// DeleteMetricImage deletes a metric image for a project alert.
//
// GitLab API docs:
// https://docs.gitlab.com/api/alert_management_alerts/#delete-metric-image
func (s *AlertManagementService) DeleteMetricImage(pid any, alertIID int, id int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/alert_management_alerts/%d/metric_images/%d", PathEscape(project), alertIID, id)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
