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
	"io"
	"net/http"
	"time"
)

type (
	AlertManagementServiceInterface interface {
		// UploadMetricImage uploads a metric image to a project alert.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/alert_management_alerts/#upload-metric-image
		UploadMetricImage(pid any, alertIID int64, content io.Reader, filename string, opt *UploadMetricImageOptions, options ...RequestOptionFunc) (*MetricImage, *Response, error)

		// ListMetricImages lists all the metric images for a project alert.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/alert_management_alerts/#list-all-metric-images
		ListMetricImages(pid any, alertIID int64, opt *ListMetricImagesOptions, options ...RequestOptionFunc) ([]*MetricImage, *Response, error)

		// UpdateMetricImage updates a metric image for a project alert.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/alert_management_alerts/#update-a-metric-image
		UpdateMetricImage(pid any, alertIID int64, id int64, opt *UpdateMetricImageOptions, options ...RequestOptionFunc) (*MetricImage, *Response, error)

		// DeleteMetricImage deletes a metric image for a project alert.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/alert_management_alerts/#delete-a-metric-image
		DeleteMetricImage(pid any, alertIID int64, id int64, options ...RequestOptionFunc) (*Response, error)
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
	ID        int64      `json:"id"`
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

func (s *AlertManagementService) UploadMetricImage(pid any, alertIID int64, content io.Reader, filename string, opt *UploadMetricImageOptions, options ...RequestOptionFunc) (*MetricImage, *Response, error) {
	return do[*MetricImage](s.client,
		withMethod(http.MethodPost),
		withPath("projects/%s/alert_management_alerts/%d/metric_images", ProjectID{pid}, alertIID),
		withUpload(content, filename, UploadFile),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// ListMetricImagesOptions represents the available ListMetricImages() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/alert_management_alerts/#list-metric-images
type ListMetricImagesOptions struct {
	ListOptions
}

func (s *AlertManagementService) ListMetricImages(pid any, alertIID int64, opt *ListMetricImagesOptions, options ...RequestOptionFunc) ([]*MetricImage, *Response, error) {
	return do[[]*MetricImage](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/alert_management_alerts/%d/metric_images", ProjectID{pid}, alertIID),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// UpdateMetricImageOptions represents the available UpdateMetricImage() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/alert_management_alerts/#update-metric-image
type UpdateMetricImageOptions struct {
	URL     *string `url:"url,omitempty" json:"url,omitempty"`
	URLText *string `url:"url_text,omitempty" json:"url_text,omitempty"`
}

func (s *AlertManagementService) UpdateMetricImage(pid any, alertIID int64, id int64, opt *UpdateMetricImageOptions, options ...RequestOptionFunc) (*MetricImage, *Response, error) {
	return do[*MetricImage](s.client,
		withMethod(http.MethodPut),
		withPath("projects/%s/alert_management_alerts/%d/metric_images/%d", ProjectID{pid}, alertIID, id),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *AlertManagementService) DeleteMetricImage(pid any, alertIID int64, id int64, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](s.client,
		withMethod(http.MethodDelete),
		withPath("projects/%s/alert_management_alerts/%d/metric_images/%d", ProjectID{pid}, alertIID, id),
		withRequestOpts(options...),
	)
	return resp, err
}
