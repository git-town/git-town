package gitlab

import (
	"bytes"
	"io"
	"net/http"
	"time"
)

type (
	UsageDataServiceInterface interface {
		GetServicePing(options ...RequestOptionFunc) (*ServicePingData, *Response, error)
		GetMetricDefinitionsAsYAML(options ...RequestOptionFunc) (io.Reader, *Response, error)
		GetQueries(options ...RequestOptionFunc) (*ServicePingQueries, *Response, error)
		GetNonSQLMetrics(options ...RequestOptionFunc) (*ServicePingNonSqlMetrics, *Response, error)
		TrackEvent(opt *TrackEventOptions, options ...RequestOptionFunc) (*Response, error)
		TrackEvents(opt *TrackEventsOptions, options ...RequestOptionFunc) (*Response, error)
	}

	// UsageDataService handles communication with the service ping related
	// methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/usage_data.html
	UsageDataService struct {
		client *Client
	}
)

// ServicePingData represents a service ping data response.
type ServicePingData struct {
	RecordedAt *time.Time        `json:"recorded_at"`
	License    map[string]string `json:"license"`
	Counts     map[string]int    `json:"counts"`
}

// GetServicePing gets the current service ping data.
//
// GitLab API docs:
// https://docs.gitlab.com/api/usage_data.html#export-service-ping-data
func (s *UsageDataService) GetServicePing(options ...RequestOptionFunc) (*ServicePingData, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "usage_data/service_ping", nil, options)
	if err != nil {
		return nil, nil, err
	}

	sp := new(ServicePingData)
	resp, err := s.client.Do(req, sp)
	if err != nil {
		return nil, resp, err
	}

	return sp, resp, nil
}

// GetMetricDefinitionsAsYAML gets all metric definitions as a single YAML file.
//
// GitLab API docs:
// https://docs.gitlab.com/api/usage_data.html#export-metric-definitions-as-a-single-yaml-file
func (s *UsageDataService) GetMetricDefinitionsAsYAML(options ...RequestOptionFunc) (io.Reader, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "usage_data/metric_definitions", nil, options)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Accept", "text/yaml")

	var buf bytes.Buffer
	resp, err := s.client.Do(req, &buf)
	if err != nil {
		return nil, resp, err
	}

	return &buf, resp, nil
}

// ServicePingQueries represents the raw service ping SQL queries.
type ServicePingQueries struct {
	RecordedAt            *time.Time        `json:"recorded_at"`
	UUID                  string            `json:"uuid"`
	Hostname              string            `json:"hostname"`
	Version               string            `json:"version"`
	InstallationType      string            `json:"installation_type"`
	ActiveUserCount       string            `json:"active_user_count"`
	Edition               string            `json:"edition"`
	LicenseMD5            string            `json:"license_md5"`
	LicenseSHA256         string            `json:"license_sha256"`
	LicenseID             string            `json:"license_id"`
	HistoricalMaxUsers    int               `json:"historical_max_users"`
	Licensee              map[string]string `json:"licensee"`
	LicenseUserCount      int               `json:"license_user_count"`
	LicenseStartsAt       string            `json:"license_starts_at"`
	LicenseExpiresAt      string            `json:"license_expires_at"`
	LicensePlan           string            `json:"license_plan"`
	LicenseAddOns         map[string]int    `json:"license_add_ons"`
	LicenseTrial          string            `json:"license_trial"`
	LicenseSubscriptionID string            `json:"license_subscription_id"`
	License               map[string]string `json:"license"`
	Settings              map[string]string `json:"settings"`
	Counts                map[string]string `json:"counts"`
}

// GetQueries gets all raw SQL queries used to compute service ping.
//
// GitLab API docs:
// https://docs.gitlab.com/api/usage_data.html#export-service-ping-sql-queries
func (s *UsageDataService) GetQueries(options ...RequestOptionFunc) (*ServicePingQueries, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "usage_data/queries", nil, options)
	if err != nil {
		return nil, nil, err
	}

	sq := new(ServicePingQueries)
	resp, err := s.client.Do(req, sq)
	if err != nil {
		return nil, resp, err
	}

	return sq, resp, nil
}

// ServicePingNonSqlMetrics represents the non-SQL metrics used in service ping.
type ServicePingNonSqlMetrics struct {
	RecordedAt            string            `json:"recorded_at"`
	UUID                  string            `json:"uuid"`
	Hostname              string            `json:"hostname"`
	Version               string            `json:"version"`
	InstallationType      string            `json:"installation_type"`
	ActiveUserCount       int               `json:"active_user_count"`
	Edition               string            `json:"edition"`
	LicenseMD5            string            `json:"license_md5"`
	LicenseSHA256         string            `json:"license_sha256"`
	LicenseID             string            `json:"license_id"`
	HistoricalMaxUsers    int               `json:"historical_max_users"`
	Licensee              map[string]string `json:"licensee"`
	LicenseUserCount      int               `json:"license_user_count"`
	LicenseStartsAt       string            `json:"license_starts_at"`
	LicenseExpiresAt      string            `json:"license_expires_at"`
	LicensePlan           string            `json:"license_plan"`
	LicenseAddOns         map[string]int    `json:"license_add_ons"`
	LicenseTrial          string            `json:"license_trial"`
	LicenseSubscriptionID string            `json:"license_subscription_id"`
	License               map[string]string `json:"license"`
	Settings              map[string]string `json:"settings"`
}

// GetNonSQLMetrics gets all non-SQL metrics data used in the service ping.
//
// GitLab API docs:
// https://docs.gitlab.com/api/usage_data.html#usagedatanonsqlmetrics-api
func (s *UsageDataService) GetNonSQLMetrics(options ...RequestOptionFunc) (*ServicePingNonSqlMetrics, *Response, error) {
	req, err := s.client.NewRequest(http.MethodGet, "usage_data/non_sql_metrics", nil, options)
	if err != nil {
		return nil, nil, err
	}

	nsm := new(ServicePingNonSqlMetrics)
	resp, err := s.client.Do(req, nsm)
	if err != nil {
		return nil, resp, err
	}

	return nsm, resp, nil
}

// TrackEventOptions represents the available options for tracking events.
type TrackEventOptions struct {
	Event                string            `json:"event" url:"event"`
	SendToSnowplow       *bool             `json:"send_to_snowplow,omitempty" url:"send_to_snowplow,omitempty"`
	NamespaceID          *int              `json:"namespace_id,omitempty" url:"namespace_id,omitempty"`
	ProjectID            *int              `json:"project_id,omitempty" url:"project_id,omitempty"`
	AdditionalProperties map[string]string `json:"additional_properties,omitempty" url:"additional_properties,omitempty"`
}

// TrackEvent tracks an internal GitLab event.
//
// GitLab API docs:
// https://docs.gitlab.com/api/usage_data.html#events-tracking-api
func (s *UsageDataService) TrackEvent(opt *TrackEventOptions, options ...RequestOptionFunc) (*Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "usage_data/track_event", opt, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// TrackEventsOptions represents the available options for tracking multiple events.
type TrackEventsOptions struct {
	Events []TrackEventOptions `json:"events" url:"events"`
}

// TrackEvents tracks multiple internal GitLab events.
//
// GitLab API docs:
// https://docs.gitlab.com/api/usage_data.html#events-tracking-api
func (s *UsageDataService) TrackEvents(opt *TrackEventsOptions, options ...RequestOptionFunc) (*Response, error) {
	req, err := s.client.NewRequest(http.MethodPost, "usage_data/track_events", opt, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
