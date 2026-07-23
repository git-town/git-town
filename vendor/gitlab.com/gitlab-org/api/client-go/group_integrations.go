package gitlab

import (
	"net/http"
	"time"
)

// GroupMattermostIntegration represents a Mattermost integration for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#mattermost-notifications
type GroupMattermostIntegration struct {
	Integration
	NotifyOnlyBrokenPipelines  bool                                  `json:"notify_only_broken_pipelines"`
	BranchesToBeNotified       string                                `json:"branches_to_be_notified"`
	LabelsToBeNotified         string                                `json:"labels_to_be_notified"`
	LabelsToBeNotifiedBehavior string                                `json:"labels_to_be_notified_behavior"`
	NotifyOnlyDefaultBranch    bool                                  `json:"notify_only_default_branch"`
	Properties                 *GroupMattermostIntegrationProperties `json:"properties"`
}

type GroupMattermostIntegrationProperties struct {
	WebHook                  string `json:"webhook"`
	Username                 string `json:"username"`
	Channel                  string `json:"channel"`
	PushChannel              string `json:"push_channel"`
	IssueChannel             string `json:"issue_channel"`
	ConfidentialIssueChannel string `json:"confidential_issue_channel"`
	MergeRequestChannel      string `json:"merge_request_channel"`
	NoteChannel              string `json:"note_channel"`
	ConfidentialNoteChannel  string `json:"confidential_note_channel"`
	TagPushChannel           string `json:"tag_push_channel"`
	PipelineChannel          string `json:"pipeline_channel"`
	WikiPageChannel          string `json:"wiki_page_channel"`
	DeploymentChannel        string `json:"deployment_channel"`
	AlertChannel             string `json:"alert_channel"`
	VulnerabilityChannel     string `json:"vulnerability_channel"`
}

// GroupMattermostIntegrationOptions represents the available options for
// creating or updating a Mattermost integration for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#mattermost-notifications
type GroupMattermostIntegrationOptions struct {
	WebHook                    *string `url:"webhook,omitempty" json:"webhook,omitempty"`
	Username                   *string `url:"username,omitempty" json:"username,omitempty"`
	Channel                    *string `url:"channel,omitempty" json:"channel,omitempty"`
	NotifyOnlyBrokenPipelines  *bool   `url:"notify_only_broken_pipelines,omitempty" json:"notify_only_broken_pipelines,omitempty"`
	BranchesToBeNotified       *string `url:"branches_to_be_notified,omitempty" json:"branches_to_be_notified,omitempty"`
	PushEvents                 *bool   `url:"push_events,omitempty" json:"push_events,omitempty"`
	IssuesEvents               *bool   `url:"issues_events,omitempty" json:"issues_events,omitempty"`
	ConfidentialIssuesEvents   *bool   `url:"confidential_issues_events,omitempty" json:"confidential_issues_events,omitempty"`
	MergeRequestsEvents        *bool   `url:"merge_requests_events,omitempty" json:"merge_requests_events,omitempty"`
	TagPushEvents              *bool   `url:"tag_push_events,omitempty" json:"tag_push_events,omitempty"`
	NoteEvents                 *bool   `url:"note_events,omitempty" json:"note_events,omitempty"`
	ConfidentialNoteEvents     *bool   `url:"confidential_note_events,omitempty" json:"confidential_note_events,omitempty"`
	PipelineEvents             *bool   `url:"pipeline_events,omitempty" json:"pipeline_events,omitempty"`
	WikiPageEvents             *bool   `url:"wiki_page_events,omitempty" json:"wiki_page_events,omitempty"`
	DeploymentEvents           *bool   `url:"deployment_events,omitempty" json:"deployment_events,omitempty"`
	AlertEvents                *bool   `url:"alert_events,omitempty" json:"alert_events,omitempty"`
	VulnerabilityEvents        *bool   `url:"vulnerability_events,omitempty" json:"vulnerability_events,omitempty"`
	PushChannel                *string `url:"push_channel,omitempty" json:"push_channel,omitempty"`
	IssueChannel               *string `url:"issue_channel,omitempty" json:"issue_channel,omitempty"`
	ConfidentialIssueChannel   *string `url:"confidential_issue_channel,omitempty" json:"confidential_issue_channel,omitempty"`
	MergeRequestChannel        *string `url:"merge_request_channel,omitempty" json:"merge_request_channel,omitempty"`
	NoteChannel                *string `url:"note_channel,omitempty" json:"note_channel,omitempty"`
	ConfidentialNoteChannel    *string `url:"confidential_note_channel,omitempty" json:"confidential_note_channel,omitempty"`
	TagPushChannel             *string `url:"tag_push_channel,omitempty" json:"tag_push_channel,omitempty"`
	PipelineChannel            *string `url:"pipeline_channel,omitempty" json:"pipeline_channel,omitempty"`
	WikiPageChannel            *string `url:"wiki_page_channel,omitempty" json:"wiki_page_channel,omitempty"`
	DeploymentChannel          *string `url:"deployment_channel,omitempty" json:"deployment_channel,omitempty"`
	AlertChannel               *string `url:"alert_channel,omitempty" json:"alert_channel,omitempty"`
	VulnerabilityChannel       *string `url:"vulnerability_channel,omitempty" json:"vulnerability_channel,omitempty"`
	LabelsToBeNotified         *string `url:"labels_to_be_notified,omitempty" json:"labels_to_be_notified,omitempty"`
	LabelsToBeNotifiedBehavior *string `url:"labels_to_be_notified_behavior,omitempty" json:"labels_to_be_notified_behavior,omitempty"`
	NotifyOnlyDefaultBranch    *bool   `url:"notify_only_default_branch,omitempty" json:"notify_only_default_branch,omitempty"`
	UseInheritedSettings       *bool   `url:"use_inherited_settings,omitempty" json:"use_inherited_settings,omitempty"`
}

// GetGroupMattermostIntegration retrieves the Mattermost integration for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#mattermost-notifications
func (s *IntegrationsService) GetGroupMattermostIntegration(gid any, options ...RequestOptionFunc) (*GroupMattermostIntegration, *Response, error) {
	return do[*GroupMattermostIntegration](
		s.client,
		withPath("groups/%s/integrations/mattermost", GroupID{gid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

// SetGroupMattermostIntegration creates or updates the Mattermost integration for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#mattermost-notifications
func (s *IntegrationsService) SetGroupMattermostIntegration(gid any, opt *GroupMattermostIntegrationOptions, options ...RequestOptionFunc) (*GroupMattermostIntegration, *Response, error) {
	return do[*GroupMattermostIntegration](
		s.client,
		withPath("groups/%s/integrations/mattermost", GroupID{gid}),
		withMethod(http.MethodPut),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteGroupMattermostIntegration removes the Mattermost integration from a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#mattermost-notifications
func (s *IntegrationsService) DeleteGroupMattermostIntegration(gid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](
		s.client,
		withPath("groups/%s/integrations/mattermost", GroupID{gid}),
		withMethod(http.MethodDelete),
		withRequestOpts(options...),
	)
	return resp, err
}

// GroupMattermostSlashCommandsIntegration represents a Mattermost slash commands integration for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#mattermost-slash-commands
type GroupMattermostSlashCommandsIntegration struct {
	ID        int        `json:"id"`
	Title     string     `json:"title"`
	Slug      string     `json:"slug"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	Token     string     `json:"token"`
}

// GroupMattermostSlashCommandsIntegrationOptions represents the available options for
// creating or updating a Mattermost slash commands integration for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#mattermost-slash-commands
type GroupMattermostSlashCommandsIntegrationOptions struct {
	Token                *string `url:"token,omitempty" json:"token,omitempty"`
	UseInheritedSettings *bool   `url:"use_inherited_settings,omitempty" json:"use_inherited_settings,omitempty"`
}

// GetGroupMattermostSlashCommandsIntegration retrieves the Mattermost slash commands integration for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#mattermost-slash-commands
func (s *IntegrationsService) GetGroupMattermostSlashCommandsIntegration(gid any, options ...RequestOptionFunc) (*GroupMattermostSlashCommandsIntegration, *Response, error) {
	return do[*GroupMattermostSlashCommandsIntegration](
		s.client,
		withPath("groups/%s/integrations/mattermost-slash-commands", GroupID{gid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

// SetGroupMattermostSlashCommandsIntegration creates or updates the Mattermost slash commands integration for a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#mattermost-slash-commands
func (s *IntegrationsService) SetGroupMattermostSlashCommandsIntegration(gid any, opt *GroupMattermostSlashCommandsIntegrationOptions, options ...RequestOptionFunc) (*GroupMattermostSlashCommandsIntegration, *Response, error) {
	return do[*GroupMattermostSlashCommandsIntegration](
		s.client,
		withPath("groups/%s/integrations/mattermost-slash-commands", GroupID{gid}),
		withMethod(http.MethodPut),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

// DeleteGroupMattermostSlashCommandsIntegration removes the Mattermost slash commands integration from a group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#mattermost-slash-commands
func (s *IntegrationsService) DeleteGroupMattermostSlashCommandsIntegration(gid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](
		s.client,
		withPath("groups/%s/integrations/mattermost-slash-commands", GroupID{gid}),
		withMethod(http.MethodDelete),
		withRequestOpts(options...),
	)
	return resp, err
}
