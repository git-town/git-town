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
	"net/http"
	"time"
)

type (
	IntegrationsServiceInterface interface {
		ListActiveGroupIntegrations(gid any, opt *ListActiveIntegrationsOptions, options ...RequestOptionFunc) ([]*Integration, *Response, error)
		SetUpGroupHarbor(gid any, opt *SetUpHarborOptions, options ...RequestOptionFunc) (*Integration, *Response, error)
		DisableGroupHarbor(gid any, options ...RequestOptionFunc) (*Response, error)
		GetGroupHarborSettings(gid any, options ...RequestOptionFunc) (*Integration, *Response, error)
		SetGroupMicrosoftTeamsNotifications(gid any, opt *SetMicrosoftTeamsNotificationsOptions, options ...RequestOptionFunc) (*Integration, *Response, error)
		DisableGroupMicrosoftTeamsNotifications(gid any, options ...RequestOptionFunc) (*Response, error)
		GetGroupMicrosoftTeamsNotifications(gid any, options ...RequestOptionFunc) (*Integration, *Response, error)
	}

	// IntegrationsService handles communication with the group
	// integrations related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/ee/api/group_integrations.html
	IntegrationsService struct {
		client *Client
	}
)

var _ IntegrationsServiceInterface = (*IntegrationsService)(nil)

// Integration represents a GitLab group or project integration.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/
// https://docs.gitlab.com/api/project_integrations/
type Integration struct {
	ID                             int        `json:"id"`
	Title                          string     `json:"title"`
	Slug                           string     `json:"slug"`
	CreatedAt                      *time.Time `json:"created_at"`
	UpdatedAt                      *time.Time `json:"updated_at"`
	Active                         bool       `json:"active"`
	AlertEvents                    bool       `json:"alert_events"`
	CommitEvents                   bool       `json:"commit_events"`
	ConfidentialIssuesEvents       bool       `json:"confidential_issues_events"`
	ConfidentialNoteEvents         bool       `json:"confidential_note_events"`
	DeploymentEvents               bool       `json:"deployment_events"`
	GroupConfidentialMentionEvents bool       `json:"group_confidential_mention_events"`
	GroupMentionEvents             bool       `json:"group_mention_events"`
	IncidentEvents                 bool       `json:"incident_events"`
	IssuesEvents                   bool       `json:"issues_events"`
	JobEvents                      bool       `json:"job_events"`
	MergeRequestsEvents            bool       `json:"merge_requests_events"`
	NoteEvents                     bool       `json:"note_events"`
	PipelineEvents                 bool       `json:"pipeline_events"`
	PushEvents                     bool       `json:"push_events"`
	TagPushEvents                  bool       `json:"tag_push_events"`
	VulnerabilityEvents            bool       `json:"vulnerability_events"`
	WikiPageEvents                 bool       `json:"wiki_page_events"`
	CommentOnEventEnabled          bool       `json:"comment_on_event_enabled"`
	Inherited                      bool       `json:"inherited"`
}

// ListActiveIntegrationsOptions represents the available
// ListActiveIntegrations() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#list-all-active-integrations
type ListActiveIntegrationsOptions struct {
	ListOptions
}

// ListActiveGroupIntegrations gets a list of all active group integrations.
// The vulnerability_events field is only available for GitLab Enterprise Edition.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#list-all-active-integrations
func (s *IntegrationsService) ListActiveGroupIntegrations(gid any, opt *ListActiveIntegrationsOptions, options ...RequestOptionFunc) ([]*Integration, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/integrations", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	var integrations []*Integration
	resp, err := s.client.Do(req, &integrations)
	if err != nil {
		return nil, resp, err
	}

	return integrations, resp, nil
}

// SetUpHarborOptions represents the available SetUpGroupHarbor()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#set-up-harbor
type SetUpHarborOptions struct {
	URL                  *string `url:"url,omitempty" json:"url,omitempty"`
	ProjectName          *string `url:"project_name,omitempty" json:"project_name,omitempty"`
	Username             *string `url:"username,omitempty" json:"username,omitempty"`
	Password             *string `url:"password,omitempty" json:"password,omitempty"`
	UseInheritedSettings *bool   `url:"use_inherited_settings,omitempty" json:"use_inherited_settings,omitempty"`
}

// SetUpGroupHarbor sets up the Harbor integration for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#set-up-harbor
func (s *IntegrationsService) SetUpGroupHarbor(gid any, opt *SetUpHarborOptions, options ...RequestOptionFunc) (*Integration, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/integrations/harbor", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	integration := new(Integration)
	resp, err := s.client.Do(req, integration)
	if err != nil {
		return nil, resp, err
	}
	return integration, resp, nil
}

// DisableGroupHarbor disables the Harbor integration for a group.
// Integration settings are reset.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#disable-harbor
func (s *IntegrationsService) DisableGroupHarbor(gid any, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/integrations/harbor", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetGroupHarborSettings gets the Harbor integration for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#get-harbor-settings
func (s *IntegrationsService) GetGroupHarborSettings(gid any, options ...RequestOptionFunc) (*Integration, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/integrations/harbor", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	integration := new(Integration)
	resp, err := s.client.Do(req, integration)
	if err != nil {
		return nil, nil, err
	}
	return integration, resp, nil
}

// SetMicrosoftTeamsNotificationsOptions represents the available
// SetMicrosoftTeamsNotifications() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#set-up-microsoft-teams-notifications
type SetMicrosoftTeamsNotificationsOptions struct {
	Targets                   *string `url:"targets,omitempty"`
	Webhook                   *string `url:"webhook,omitempty"`
	NotifyOnlyBrokenPipelines *bool   `url:"notify_only_broken_pipelines,omitempty"`
	NotifyOnlyDefaultBranch   *bool   `url:"notify_only_default_branch,omitempty"`
	BranchesToBeNotified      *string `url:"branches_to_be_notified,omitempty"`
	PushEvents                *bool   `url:"push_events,omitempty"`
	IssuesEvents              *bool   `url:"issues_events,omitempty"`
	ConfidentialIssuesEvents  *bool   `url:"confidential_issues_events,omitempty"`
	MergeRequestsEvents       *bool   `url:"merge_requests_events,omitempty"`
	TagPushEvents             *bool   `url:"tag_push_events,omitempty"`
	NoteEvents                *bool   `url:"note_events,omitempty"`
	ConfidentialNoteEvents    *bool   `url:"confidential_note_events,omitempty"`
	PipelineEvents            *bool   `url:"pipeline_events,omitempty"`
	WikiPageEvents            *bool   `url:"wiki_page_events,omitempty"`
	UseInheritedSettings      *bool   `url:"use_inherited_settings,omitempty"`
}

// SetGroupMicrosoftTeamsNotifications sets up Microsoft Teams notifications for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#set-up-microsoft-teams-notifications
func (s *IntegrationsService) SetGroupMicrosoftTeamsNotifications(gid any, opt *SetMicrosoftTeamsNotificationsOptions, options ...RequestOptionFunc) (*Integration, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/integrations/microsoft_teams", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	integration := new(Integration)
	resp, err := s.client.Do(req, integration)
	if err != nil {
		return nil, resp, err
	}
	return integration, resp, nil
}

// DisableGroupMicrosoftTeamsNotifications disables Microsoft Teams notifications
// for a group. Integration settings are reset.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#disable-microsoft-teams-notifications
func (s *IntegrationsService) DisableGroupMicrosoftTeamsNotifications(gid any, options ...RequestOptionFunc) (*Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("groups/%s/integrations/microsoft_teams", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	resp, err := s.client.Do(req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetGroupMicrosoftTeamsNotifications gets the Microsoft Teams notifications for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#get-microsoft-teams-notifications-settings
func (s *IntegrationsService) GetGroupMicrosoftTeamsNotifications(gid any, options ...RequestOptionFunc) (*Integration, *Response, error) {
	group, err := parseID(gid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("groups/%s/integrations/microsoft_teams", PathEscape(group))

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	integration := new(Integration)
	resp, err := s.client.Do(req, integration)
	if err != nil {
		return nil, nil, err
	}
	return integration, resp, nil
}
