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
	"net/http"
	"time"
)

type (
	IntegrationsServiceInterface interface {
		// ListActiveGroupIntegrations gets a list of all active group integrations.
		// The vulnerability_events field is only available for GitLab Enterprise Edition.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#list-all-active-integrations
		ListActiveGroupIntegrations(gid any, opt *ListActiveIntegrationsOptions, options ...RequestOptionFunc) ([]*Integration, *Response, error)

		// SetUpGroupHarbor sets up the Harbor integration for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#set-up-harbor
		SetUpGroupHarbor(gid any, opt *SetUpHarborOptions, options ...RequestOptionFunc) (*Integration, *Response, error)

		// DisableGroupHarbor disables the Harbor integration for a group.
		// Integration settings are reset.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#disable-harbor
		DisableGroupHarbor(gid any, options ...RequestOptionFunc) (*Response, error)

		// GetGroupHarborSettings gets the Harbor integration for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#get-harbor-settings
		GetGroupHarborSettings(gid any, options ...RequestOptionFunc) (*Integration, *Response, error)

		// SetGroupMicrosoftTeamsNotifications sets up Microsoft Teams notifications for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#set-up-microsoft-teams-notifications
		SetGroupMicrosoftTeamsNotifications(gid any, opt *SetMicrosoftTeamsNotificationsOptions, options ...RequestOptionFunc) (*Integration, *Response, error)

		// DisableGroupMicrosoftTeamsNotifications disables Microsoft Teams notifications
		// for a group. Integration settings are reset.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#disable-microsoft-teams-notifications
		DisableGroupMicrosoftTeamsNotifications(gid any, options ...RequestOptionFunc) (*Response, error)

		// GetGroupMicrosoftTeamsNotifications gets the Microsoft Teams notifications for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#get-microsoft-teams-notifications-settings
		GetGroupMicrosoftTeamsNotifications(gid any, options ...RequestOptionFunc) (*Integration, *Response, error)

		// SetUpGroupJira sets up the Jira integration for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#set-up-jira
		SetUpGroupJira(gid any, opt *SetUpJiraOptions, options ...RequestOptionFunc) (*Integration, *Response, error)

		// DisableGroupJira disables the Jira integration for a group.
		// Integration settings are reset.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#disable-jira
		DisableGroupJira(gid any, options ...RequestOptionFunc) (*Response, error)

		// GetGroupJiraSettings gets the Jira integration for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#get-jira-settings
		GetGroupJiraSettings(gid any, options ...RequestOptionFunc) (*Integration, *Response, error)

		// GetGroupSlackSettings gets the Slack integration for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#get-slack-settings
		GetGroupSlackSettings(gid any, options ...RequestOptionFunc) (*SlackIntegration, *Response, error)

		// SetGroupSlackSettings sets up the Slack integration for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#set-up-slack
		SetGroupSlackSettings(gid any, opt *SetGroupSlackOptions, options ...RequestOptionFunc) (*SlackIntegration, *Response, error)

		// DisableGroupSlack disables the Slack integration for a group.
		// Integration settings are reset.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#disable-slack
		DisableGroupSlack(gid any, options ...RequestOptionFunc) (*Response, error)

		// GetGroupDiscordSettings gets the Discord integration settings for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#discord
		GetGroupDiscordSettings(gid any, options ...RequestOptionFunc) (*DiscordIntegration, *Response, error)

		// GetGroupTelegramSettings gets the Telegram integration settings for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#telegram
		GetGroupTelegramSettings(gid any, options ...RequestOptionFunc) (*TelegramIntegration, *Response, error)

		// GetGroupMattermostSettings gets the Mattermost integration settings for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#mattermost-notifications
		GetGroupMattermostSettings(gid any, options ...RequestOptionFunc) (*MattermostIntegration, *Response, error)

		// GetGroupMatrixSettings gets the Matrix integration settings for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#matrix-notifications
		GetGroupMatrixSettings(gid any, options ...RequestOptionFunc) (*MatrixIntegration, *Response, error)

		// GetGroupGoogleChatSettings gets the Google Chat integration settings for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#google-chat
		GetGroupGoogleChatSettings(gid any, options ...RequestOptionFunc) (*GoogleChatIntegration, *Response, error)

		// SetProjectGoogleChatSettings sets up the Google Chat integration for a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/integrations.html#set-up-google-chat
		SetProjectGoogleChatSettings(pid any, opt *SetProjectGoogleChatOptions, options ...RequestOptionFunc) (*GoogleChatIntegration, *Response, error)

		// DisableProjectGoogleChat disables the Google Chat integration for a project.
		// Integration settings are reset.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/integrations.html#disable-google-chat
		DisableProjectGoogleChat(pid any, options ...RequestOptionFunc) (*Response, error)

		// GetProjectGoogleChatSettings gets the Google Chat integration settings for a project.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/integrations.html#get-google-chat-settings
		GetProjectGoogleChatSettings(pid any, options ...RequestOptionFunc) (*GoogleChatIntegration, *Response, error)

		// GetGroupMattermostIntegration retrieves the Mattermost integration for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#mattermost-notifications
		GetGroupMattermostIntegration(gid any, options ...RequestOptionFunc) (*GroupMattermostIntegration, *Response, error)

		// SetGroupMattermostIntegration creates or updates the Mattermost integration for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#mattermost-notifications
		SetGroupMattermostIntegration(gid any, opt *GroupMattermostIntegrationOptions, options ...RequestOptionFunc) (*GroupMattermostIntegration, *Response, error)

		// DeleteGroupMattermostIntegration removes the Mattermost integration from a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#mattermost-notifications
		DeleteGroupMattermostIntegration(gid any, options ...RequestOptionFunc) (*Response, error)

		// GetGroupMattermostSlashCommandsIntegration retrieves the Mattermost slash commands integration for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#mattermost-slash-commands
		GetGroupMattermostSlashCommandsIntegration(gid any, options ...RequestOptionFunc) (*GroupMattermostSlashCommandsIntegration, *Response, error)

		// SetGroupMattermostSlashCommandsIntegration creates or updates the Mattermost slash commands integration for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#mattermost-slash-commands
		SetGroupMattermostSlashCommandsIntegration(gid any, opt *GroupMattermostSlashCommandsIntegrationOptions, options ...RequestOptionFunc) (*GroupMattermostSlashCommandsIntegration, *Response, error)

		// DeleteGroupMattermostSlashCommandsIntegration removes the Mattermost slash commands integration from a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#mattermost-slash-commands
		DeleteGroupMattermostSlashCommandsIntegration(gid any, options ...RequestOptionFunc) (*Response, error)

		// GetGroupWebexTeamsSettings gets the Webex Teams integration for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#get-webex-teams-settings
		GetGroupWebexTeamsSettings(gid any, options ...RequestOptionFunc) (*WebexTeamsIntegration, *Response, error)

		// SetGroupWebexTeamsSettings sets up the Webex Teams integration for a group.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#set-up-webex-teams
		SetGroupWebexTeamsSettings(gid any, opt *SetGroupWebexTeamsOptions, options ...RequestOptionFunc) (*WebexTeamsIntegration, *Response, error)

		// DisableGroupWebexTeams disables the Webex Teams integration for a group.
		// Integration settings are reset.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/group_integrations/#disable-webex-teams
		DisableGroupWebexTeams(gid any, options ...RequestOptionFunc) (*Response, error)
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
	ID                             int64      `json:"id"`
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

// SlackIntegration represents the Slack integration settings.
// It embeds the generic Integration struct and adds Slack-specific properties.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#get-slack-settings
type SlackIntegration struct {
	Integration
	Properties SlackIntegrationProperties `json:"properties"`
}

// SlackIntegrationProperties represents Slack specific properties
// returned by the GitLab API.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#get-slack-settings
type SlackIntegrationProperties struct {
	Username                        string `json:"username"`
	Channel                         string `json:"channel"`
	NotifyOnlyBrokenPipelines       bool   `json:"notify_only_broken_pipelines"`
	BranchesToBeNotified            string `json:"branches_to_be_notified"`
	LabelsToBeNotified              string `json:"labels_to_be_notified"`
	LabelsToBeNotifiedBehavior      string `json:"labels_to_be_notified_behavior"`
	PushChannel                     string `json:"push_channel"`
	IssueChannel                    string `json:"issue_channel"`
	ConfidentialIssueChannel        string `json:"confidential_issue_channel"`
	MergeRequestChannel             string `json:"merge_request_channel"`
	NoteChannel                     string `json:"note_channel"`
	ConfidentialNoteChannel         string `json:"confidential_note_channel"`
	TagPushChannel                  string `json:"tag_push_channel"`
	PipelineChannel                 string `json:"pipeline_channel"`
	WikiPageChannel                 string `json:"wiki_page_channel"`
	DeploymentChannel               string `json:"deployment_channel"`
	IncidentChannel                 string `json:"incident_channel"`
	AlertChannel                    string `json:"alert_channel"`
	GroupMentionChannel             string `json:"group_mention_channel"`
	GroupConfidentialMentionChannel string `json:"group_confidential_mention_channel"`
}

// DiscordIntegration represents the Discord integration settings.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#discord
type DiscordIntegration struct {
	Integration
	Properties DiscordIntegrationProperties `json:"properties"`
}

// DiscordIntegrationProperties represents Discord specific properties.
type DiscordIntegrationProperties struct {
	NotifyOnlyBrokenPipelines bool   `json:"notify_only_broken_pipelines,omitempty"`
	BranchesToBeNotified      string `json:"branches_to_be_notified,omitempty"`
}

// TelegramIntegration represents the Telegram integration settings.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#telegram
type TelegramIntegration struct {
	Integration
	Properties TelegramIntegrationProperties `json:"properties"`
}

// TelegramIntegrationProperties represents Telegram specific properties.
type TelegramIntegrationProperties struct {
	Hostname                  string `json:"hostname,omitempty"`
	Room                      string `json:"room,omitempty"`
	Thread                    string `json:"thread,omitempty"`
	NotifyOnlyBrokenPipelines bool   `json:"notify_only_broken_pipelines,omitempty"`
	BranchesToBeNotified      string `json:"branches_to_be_notified,omitempty"`
}

// MattermostIntegration represents the Mattermost integration settings.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#mattermost-notifications
type MattermostIntegration struct {
	Integration
	Properties MattermostIntegrationProperties `json:"properties"`
}

// MattermostIntegrationProperties represents Mattermost specific properties.
type MattermostIntegrationProperties struct {
	Username                   string `json:"username,omitempty"`
	Channel                    string `json:"channel,omitempty"`
	NotifyOnlyBrokenPipelines  bool   `json:"notify_only_broken_pipelines,omitempty"`
	BranchesToBeNotified       string `json:"branches_to_be_notified,omitempty"`
	LabelsToBeNotified         string `json:"labels_to_be_notified,omitempty"`
	LabelsToBeNotifiedBehavior string `json:"labels_to_be_notified_behavior,omitempty"`
	PushChannel                string `json:"push_channel,omitempty"`
	IssueChannel               string `json:"issue_channel,omitempty"`
	ConfidentialIssueChannel   string `json:"confidential_issue_channel,omitempty"`
	MergeRequestChannel        string `json:"merge_request_channel,omitempty"`
	NoteChannel                string `json:"note_channel,omitempty"`
	ConfidentialNoteChannel    string `json:"confidential_note_channel,omitempty"`
	TagPushChannel             string `json:"tag_push_channel,omitempty"`
	PipelineChannel            string `json:"pipeline_channel,omitempty"`
	WikiPageChannel            string `json:"wiki_page_channel,omitempty"`
	DeploymentChannel          string `json:"deployment_channel,omitempty"`
	IncidentChannel            string `json:"incident_channel,omitempty"`
}

// MatrixIntegration represents the Matrix integration settings.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#matrix-notifications
type MatrixIntegration struct {
	Integration
	Properties MatrixIntegrationProperties `json:"properties"`
}

// MatrixIntegrationProperties represents Matrix specific properties.
type MatrixIntegrationProperties struct {
	Hostname                  string `json:"hostname,omitempty"`
	Room                      string `json:"room,omitempty"`
	NotifyOnlyBrokenPipelines bool   `json:"notify_only_broken_pipelines,omitempty"`
	BranchesToBeNotified      string `json:"branches_to_be_notified,omitempty"`
}

// GoogleChatIntegration represents the Google Chat integration settings.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#google-chat
type GoogleChatIntegration struct {
	Integration
	Properties GoogleChatIntegrationProperties `json:"properties"`
}

// GoogleChatIntegrationProperties represents Google Chat specific properties.
type GoogleChatIntegrationProperties struct {
	NotifyOnlyBrokenPipelines           bool   `json:"notify_only_broken_pipelines,omitempty"`
	NotifyOnlyWhenPipelineStatusChanges bool   `json:"notify_only_when_pipeline_status_changes,omitempty"`
	NotifyOnlyDefaultBranch             bool   `json:"notify_only_default_branch,omitempty"`
	BranchesToBeNotified                string `json:"branches_to_be_notified,omitempty"`
	PushEvents                          bool   `json:"push_events,omitempty"`
	IssuesEvents                        bool   `json:"issues_events,omitempty"`
	ConfidentialIssuesEvents            bool   `json:"confidential_issues_events,omitempty"`
	MergeRequestsEvents                 bool   `json:"merge_requests_events,omitempty"`
	TagPushEvents                       bool   `json:"tag_push_events,omitempty"`
	NoteEvents                          bool   `json:"note_events,omitempty"`
	ConfidentialNoteEvents              bool   `json:"confidential_note_events,omitempty"`
	PipelineEvents                      bool   `json:"pipeline_events,omitempty"`
	WikiPageEvents                      bool   `json:"wiki_page_events,omitempty"`
}

// WebexTeamsIntegration represents the WebexTeams integration settings.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#get-webex-teams-settings
type WebexTeamsIntegration struct {
	Integration
	Properties WebexTeamsIntegrationProperties `json:"properties"`
}

// WebexTeamsIntegrationProperties represents WebexTeams specific properties
type WebexTeamsIntegrationProperties struct {
	NotifyOnlyBrokenPipelines bool   `json:"notify_only_broken_pipelines,omitempty"`
	BranchesToBeNotified      string `json:"branches_to_be_notified,omitempty"`
}

// SetGroupWebexTeamsOptions represents the available SetGroupWebexTeamsSettings() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#set-up-webex-teams
type SetGroupWebexTeamsOptions struct {
	NotifyOnlyBrokenPipelines *bool   `url:"notify_only_broken_pipelines,omitempty" json:"notify_only_broken_pipelines,omitempty"`
	BranchesToBeNotified      *string `url:"branches_to_be_notified,omitempty" json:"branches_to_be_notified,omitempty"`
}

// SetProjectGoogleChatOptions represents the available SetProjectGoogleChatSettings() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/integrations.html#set-up-google-chat
type SetProjectGoogleChatOptions struct {
	Webhook                             *string `url:"webhook,omitempty" json:"webhook,omitempty"`
	NotifyOnlyBrokenPipelines           *bool   `url:"notify_only_broken_pipelines,omitempty" json:"notify_only_broken_pipelines,omitempty"`
	NotifyOnlyWhenPipelineStatusChanges *bool   `url:"notify_only_when_pipeline_status_changes,omitempty" json:"notify_only_when_pipeline_status_changes,omitempty"`
	NotifyOnlyDefaultBranch             *bool   `url:"notify_only_default_branch,omitempty" json:"notify_only_default_branch,omitempty"`
	BranchesToBeNotified                *string `url:"branches_to_be_notified,omitempty" json:"branches_to_be_notified,omitempty"`
	PushEvents                          *bool   `url:"push_events,omitempty" json:"push_events,omitempty"`
	IssuesEvents                        *bool   `url:"issues_events,omitempty" json:"issues_events,omitempty"`
	ConfidentialIssuesEvents            *bool   `url:"confidential_issues_events,omitempty" json:"confidential_issues_events,omitempty"`
	MergeRequestsEvents                 *bool   `url:"merge_requests_events,omitempty" json:"merge_requests_events,omitempty"`
	TagPushEvents                       *bool   `url:"tag_push_events,omitempty" json:"tag_push_events,omitempty"`
	NoteEvents                          *bool   `url:"note_events,omitempty" json:"note_events,omitempty"`
	ConfidentialNoteEvents              *bool   `url:"confidential_note_events,omitempty" json:"confidential_note_events,omitempty"`
	PipelineEvents                      *bool   `url:"pipeline_events,omitempty" json:"pipeline_events,omitempty"`
	WikiPageEvents                      *bool   `url:"wiki_page_events,omitempty" json:"wiki_page_events,omitempty"`
	UseInheritedSettings                *bool   `url:"use_inherited_settings,omitempty" json:"use_inherited_settings,omitempty"`
}

// ListActiveIntegrationsOptions represents the available
// ListActiveIntegrations() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#list-all-active-integrations
type ListActiveIntegrationsOptions struct {
	ListOptions
}

func (s *IntegrationsService) ListActiveGroupIntegrations(gid any, opt *ListActiveIntegrationsOptions, options ...RequestOptionFunc) ([]*Integration, *Response, error) {
	return do[[]*Integration](
		s.client,
		withPath("groups/%s/integrations", GroupID{gid}),
		withMethod(http.MethodGet),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
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

func (s *IntegrationsService) SetUpGroupHarbor(gid any, opt *SetUpHarborOptions, options ...RequestOptionFunc) (*Integration, *Response, error) {
	return do[*Integration](
		s.client,
		withPath("groups/%s/integrations/harbor", GroupID{gid}),
		withMethod(http.MethodPut),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) DisableGroupHarbor(gid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](
		s.client,
		withPath("groups/%s/integrations/harbor", GroupID{gid}),
		withMethod(http.MethodDelete),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *IntegrationsService) GetGroupHarborSettings(gid any, options ...RequestOptionFunc) (*Integration, *Response, error) {
	return do[*Integration](
		s.client,
		withPath("groups/%s/integrations/harbor", GroupID{gid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

// SetMicrosoftTeamsNotificationsOptions represents the available
// SetMicrosoftTeamsNotifications() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#set-up-microsoft-teams-notifications
type SetMicrosoftTeamsNotificationsOptions struct {
	Targets                   *string `url:"targets,omitempty" json:"targets,omitempty"`
	Webhook                   *string `url:"webhook,omitempty" json:"webhook,omitempty"`
	NotifyOnlyBrokenPipelines *bool   `url:"notify_only_broken_pipelines,omitempty" json:"notify_only_broken_pipelines,omitempty"`
	NotifyOnlyDefaultBranch   *bool   `url:"notify_only_default_branch,omitempty" json:"notify_only_default_branch,omitempty"`
	BranchesToBeNotified      *string `url:"branches_to_be_notified,omitempty" json:"branches_to_be_notified,omitempty"`
	PushEvents                *bool   `url:"push_events,omitempty" json:"push_events,omitempty"`
	IssuesEvents              *bool   `url:"issues_events,omitempty" json:"issues_events,omitempty"`
	ConfidentialIssuesEvents  *bool   `url:"confidential_issues_events,omitempty" json:"confidential_issues_events,omitempty"`
	MergeRequestsEvents       *bool   `url:"merge_requests_events,omitempty" json:"merge_requests_events,omitempty"`
	TagPushEvents             *bool   `url:"tag_push_events,omitempty" json:"tag_push_events,omitempty"`
	NoteEvents                *bool   `url:"note_events,omitempty" json:"note_events,omitempty"`
	ConfidentialNoteEvents    *bool   `url:"confidential_note_events,omitempty" json:"confidential_note_events,omitempty"`
	PipelineEvents            *bool   `url:"pipeline_events,omitempty" json:"pipeline_events,omitempty"`
	WikiPageEvents            *bool   `url:"wiki_page_events,omitempty" json:"wiki_page_events,omitempty"`
	UseInheritedSettings      *bool   `url:"use_inherited_settings,omitempty" json:"use_inherited_settings,omitempty"`
}

func (s *IntegrationsService) SetGroupMicrosoftTeamsNotifications(gid any, opt *SetMicrosoftTeamsNotificationsOptions, options ...RequestOptionFunc) (*Integration, *Response, error) {
	return do[*Integration](
		s.client,
		withPath("groups/%s/integrations/microsoft-teams", GroupID{gid}),
		withMethod(http.MethodPut),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) DisableGroupMicrosoftTeamsNotifications(gid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](
		s.client,
		withPath("groups/%s/integrations/microsoft-teams", GroupID{gid}),
		withMethod(http.MethodDelete),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *IntegrationsService) GetGroupMicrosoftTeamsNotifications(gid any, options ...RequestOptionFunc) (*Integration, *Response, error) {
	return do[*Integration](
		s.client,
		withPath("groups/%s/integrations/microsoft-teams", GroupID{gid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

// SetUpJiraOptions represents the available SetUpJira() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#set-up-jira
type SetUpJiraOptions struct {
	URL                          *string   `url:"url,omitempty" json:"url,omitempty"`
	APIURL                       *string   `url:"api_url,omitempty" json:"api_url,omitempty"`
	Username                     *string   `url:"username,omitempty" json:"username,omitempty"`
	Password                     *string   `url:"password,omitempty" json:"password,omitempty"`
	Active                       *bool     `url:"active,omitempty" json:"active,omitempty"`
	JiraAuthType                 *int64    `url:"jira_auth_type,omitempty" json:"jira_auth_type,omitempty"`
	JiraIssuePrefix              *string   `url:"jira_issue_prefix,omitempty" json:"jira_issue_prefix,omitempty"`
	JiraIssueRegex               *string   `url:"jira_issue_regex,omitempty" json:"jira_issue_regex,omitempty"`
	JiraIssueTransitionAutomatic *bool     `url:"jira_issue_transition_automatic,omitempty" json:"jira_issue_transition_automatic,omitempty"`
	JiraIssueTransitionID        *string   `url:"jira_issue_transition_id,omitempty" json:"jira_issue_transition_id,omitempty"`
	CommitEvents                 *bool     `url:"commit_events,omitempty" json:"commit_events,omitempty"`
	MergeRequestsEvents          *bool     `url:"merge_requests_events,omitempty" json:"merge_requests_events,omitempty"`
	CommentOnEventEnabled        *bool     `url:"comment_on_event_enabled,omitempty" json:"comment_on_event_enabled,omitempty"`
	IssuesEnabled                *bool     `url:"issues_enabled,omitempty" json:"issues_enabled,omitempty"`
	ProjectKeys                  *[]string `url:"project_keys,omitempty" json:"project_keys,omitempty"`
	UseInheritedSettings         *bool     `url:"use_inherited_settings,omitempty" json:"use_inherited_settings,omitempty"`
}

func (s *IntegrationsService) SetUpGroupJira(gid any, opt *SetUpJiraOptions, options ...RequestOptionFunc) (*Integration, *Response, error) {
	return do[*Integration](
		s.client,
		withPath("groups/%s/integrations/jira", GroupID{gid}),
		withMethod(http.MethodPut),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) DisableGroupJira(gid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](
		s.client,
		withPath("groups/%s/integrations/jira", GroupID{gid}),
		withMethod(http.MethodDelete),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *IntegrationsService) GetGroupJiraSettings(gid any, options ...RequestOptionFunc) (*Integration, *Response, error) {
	return do[*Integration](
		s.client,
		withPath("groups/%s/integrations/jira", GroupID{gid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

// SetGroupSlackOptions represents the available SetGroupSlackSettings() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/group_integrations/#set-up-slack-notifications
type SetGroupSlackOptions struct {
	Webhook                             *string `url:"webhook,omitempty" json:"webhook,omitempty"`
	Username                            *string `url:"username,omitempty" json:"username,omitempty"`
	Channel                             *string `url:"channel,omitempty" json:"channel,omitempty"`
	NotifyOnlyBrokenPipelines           *bool   `url:"notify_only_broken_pipelines,omitempty" json:"notify_only_broken_pipelines,omitempty"`
	NotifyOnlyWhenPipelineStatusChanges *bool   `url:"notify_only_when_pipeline_status_changes,omitempty" json:"notify_only_when_pipeline_status_changes,omitempty"`
	BranchesToBeNotified                *string `url:"branches_to_be_notified,omitempty" json:"branches_to_be_notified,omitempty"`
	LabelsToBeNotified                  *string `url:"labels_to_be_notified,omitempty" json:"labels_to_be_notified,omitempty"`
	LabelsToBeNotifiedBehavior          *string `url:"labels_to_be_notified_behavior,omitempty" json:"labels_to_be_notified_behavior,omitempty"`
	AlertEvents                         *bool   `url:"alert_events,omitempty" json:"alert_events,omitempty"`
	CommitEvents                        *bool   `url:"commit_events,omitempty" json:"commit_events,omitempty"`
	PushChannel                         *string `url:"push_channel,omitempty" json:"push_channel,omitempty"`
	PushEvents                          *bool   `url:"push_events,omitempty" json:"push_events,omitempty"`
	IssueChannel                        *string `url:"issue_channel,omitempty" json:"issue_channel,omitempty"`
	IssuesEvents                        *bool   `url:"issues_events,omitempty" json:"issues_events,omitempty"`
	JobEvents                           *bool   `url:"job_events,omitempty" json:"job_events,omitempty"`
	ConfidentialIssueChannel            *string `url:"confidential_issue_channel,omitempty" json:"confidential_issue_channel,omitempty"`
	ConfidentialIssuesEvents            *bool   `url:"confidential_issues_events,omitempty" json:"confidential_issues_events,omitempty"`
	MergeRequestChannel                 *string `url:"merge_request_channel,omitempty" json:"merge_request_channel,omitempty"`
	MergeRequestsEvents                 *bool   `url:"merge_requests_events,omitempty" json:"merge_requests_events,omitempty"`
	NoteChannel                         *string `url:"note_channel,omitempty" json:"note_channel,omitempty"`
	NoteEvents                          *bool   `url:"note_events,omitempty" json:"note_events,omitempty"`
	ConfidentialNoteChannel             *string `url:"confidential_note_channel,omitempty" json:"confidential_note_channel,omitempty"`
	ConfidentialNoteEvents              *bool   `url:"confidential_note_events,omitempty" json:"confidential_note_events,omitempty"`
	TagPushChannel                      *string `url:"tag_push_channel,omitempty" json:"tag_push_channel,omitempty"`
	TagPushEvents                       *bool   `url:"tag_push_events,omitempty" json:"tag_push_events,omitempty"`
	PipelineChannel                     *string `url:"pipeline_channel,omitempty" json:"pipeline_channel,omitempty"`
	PipelineEvents                      *bool   `url:"pipeline_events,omitempty" json:"pipeline_events,omitempty"`
	WikiPageChannel                     *string `url:"wiki_page_channel,omitempty" json:"wiki_page_channel,omitempty"`
	WikiPageEvents                      *bool   `url:"wiki_page_events,omitempty" json:"wiki_page_events,omitempty"`
	DeploymentChannel                   *string `url:"deployment_channel,omitempty" json:"deployment_channel,omitempty"`
	DeploymentEvents                    *bool   `url:"deployment_events,omitempty" json:"deployment_events,omitempty"`
	IncidentChannel                     *string `url:"incident_channel,omitempty" json:"incident_channel,omitempty"`
	IncidentEvents                      *bool   `url:"incident_events,omitempty" json:"incident_events,omitempty"`
	AlertChannel                        *string `url:"alert_channel,omitempty" json:"alert_channel,omitempty"`
	GroupMentionChannel                 *string `url:"group_mention_channel,omitempty" json:"group_mention_channel,omitempty"`
	GroupConfidentialMentionChannel     *string `url:"group_confidential_mention_channel,omitempty" json:"group_confidential_mention_channel,omitempty"`
	UseInheritedSettings                *bool   `url:"use_inherited_settings,omitempty" json:"use_inherited_settings,omitempty"`
}

func (s *IntegrationsService) SetGroupSlackSettings(gid any, opt *SetGroupSlackOptions, options ...RequestOptionFunc) (*SlackIntegration, *Response, error) {
	return do[*SlackIntegration](
		s.client,
		withPath("groups/%s/integrations/slack", GroupID{gid}),
		withMethod(http.MethodPut),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) DisableGroupSlack(gid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](
		s.client,
		withPath("groups/%s/integrations/slack", GroupID{gid}),
		withMethod(http.MethodDelete),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *IntegrationsService) GetGroupSlackSettings(gid any, options ...RequestOptionFunc) (*SlackIntegration, *Response, error) {
	return do[*SlackIntegration](
		s.client,
		withPath("groups/%s/integrations/slack", GroupID{gid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) GetGroupDiscordSettings(gid any, options ...RequestOptionFunc) (*DiscordIntegration, *Response, error) {
	return do[*DiscordIntegration](
		s.client,
		withPath("groups/%s/integrations/discord", GroupID{gid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) GetGroupTelegramSettings(gid any, options ...RequestOptionFunc) (*TelegramIntegration, *Response, error) {
	return do[*TelegramIntegration](
		s.client,
		withPath("groups/%s/integrations/telegram", GroupID{gid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) GetGroupMattermostSettings(gid any, options ...RequestOptionFunc) (*MattermostIntegration, *Response, error) {
	return do[*MattermostIntegration](
		s.client,
		withPath("groups/%s/integrations/mattermost", GroupID{gid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) GetGroupMatrixSettings(gid any, options ...RequestOptionFunc) (*MatrixIntegration, *Response, error) {
	return do[*MatrixIntegration](
		s.client,
		withPath("groups/%s/integrations/matrix", GroupID{gid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) GetGroupGoogleChatSettings(gid any, options ...RequestOptionFunc) (*GoogleChatIntegration, *Response, error) {
	return do[*GoogleChatIntegration](
		s.client,
		withPath("groups/%s/integrations/hangouts-chat", GroupID{gid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) SetProjectGoogleChatSettings(pid any, opt *SetProjectGoogleChatOptions, options ...RequestOptionFunc) (*GoogleChatIntegration, *Response, error) {
	return do[*GoogleChatIntegration](
		s.client,
		withPath("projects/%s/integrations/hangouts-chat", ProjectID{pid}),
		withMethod(http.MethodPut),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) DisableProjectGoogleChat(pid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](
		s.client,
		withPath("projects/%s/integrations/hangouts-chat", ProjectID{pid}),
		withMethod(http.MethodDelete),
		withRequestOpts(options...),
	)
	return resp, err
}

func (s *IntegrationsService) GetProjectGoogleChatSettings(pid any, options ...RequestOptionFunc) (*GoogleChatIntegration, *Response, error) {
	return do[*GoogleChatIntegration](
		s.client,
		withPath("projects/%s/integrations/hangouts-chat", ProjectID{pid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) GetGroupWebexTeamsSettings(gid any, options ...RequestOptionFunc) (*WebexTeamsIntegration, *Response, error) {
	return do[*WebexTeamsIntegration](
		s.client,
		withPath("groups/%s/integrations/webex-teams", GroupID{gid}),
		withMethod(http.MethodGet),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) SetGroupWebexTeamsSettings(gid any, opt *SetGroupWebexTeamsOptions, options ...RequestOptionFunc) (*WebexTeamsIntegration, *Response, error) {
	return do[*WebexTeamsIntegration](
		s.client,
		withPath("groups/%s/integrations/webex-teams", GroupID{gid}),
		withMethod(http.MethodPut),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *IntegrationsService) DisableGroupWebexTeams(gid any, options ...RequestOptionFunc) (*Response, error) {
	_, resp, err := do[none](
		s.client,
		withPath("groups/%s/integrations/webex-teams", GroupID{gid}),
		withMethod(http.MethodDelete),
		withRequestOpts(options...),
	)
	return resp, err
}
