package gitlab

import (
	"net/http"
	"time"
)

type (
	AuditEventsServiceInterface interface {
		// ListInstanceAuditEvents gets a list of audit events for instance.
		// Authentication as Administrator is required.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/audit_events/#retrieve-all-instance-audit-events
		// ListInstanceAuditEvents gets a list of audit events for instance.
		// Authentication as Administrator is required.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/audit_events/#retrieve-all-instance-audit-events
		ListInstanceAuditEvents(opt *ListAuditEventsOptions, options ...RequestOptionFunc) ([]*AuditEvent, *Response, error)

		// GetInstanceAuditEvent gets a specific instance audit event.
		// Authentication as Administrator is required.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/audit_events/#retrieve-single-instance-audit-event
		GetInstanceAuditEvent(event int64, options ...RequestOptionFunc) (*AuditEvent, *Response, error)

		// ListGroupAuditEvents gets a list of audit events for the specified group
		// viewable by the authenticated user.
		//
		// GitLab API docs:
		// https://docs.gitlab.com/api/audit_events/#retrieve-all-group-audit-events
		ListGroupAuditEvents(gid any, opt *ListAuditEventsOptions, options ...RequestOptionFunc) ([]*AuditEvent, *Response, error)

		// GetGroupAuditEvent gets a specific group audit event.
		//
		// GitLab API docs: https://docs.gitlab.com/api/audit_events/#retrieve-a-specific-group-audit-event
		GetGroupAuditEvent(gid any, event int64, options ...RequestOptionFunc) (*AuditEvent, *Response, error)

		// ListProjectAuditEvents gets a list of audit events for the specified project
		// viewable by the authenticated user.
		//
		// GitLab API docs: https://docs.gitlab.com/api/audit_events/#retrieve-all-project-audit-events
		ListProjectAuditEvents(pid any, opt *ListAuditEventsOptions, options ...RequestOptionFunc) ([]*AuditEvent, *Response, error)

		// GetProjectAuditEvent gets a specific project audit event.
		//
		// GitLab API docs: https://docs.gitlab.com/api/audit_events/#retrieve-a-specific-project-audit-event
		GetProjectAuditEvent(pid any, event int64, options ...RequestOptionFunc) (*AuditEvent, *Response, error)
	}

	// AuditEventsService handles communication with the project/group/instance
	// audit event related methods of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/audit_events/
	AuditEventsService struct {
		client *Client
	}
)

var _ AuditEventsServiceInterface = (*AuditEventsService)(nil)

// AuditEvent represents an audit event for a group, a project or the instance.
//
// GitLab API docs: https://docs.gitlab.com/api/audit_events/
type AuditEvent struct {
	ID         int64             `json:"id"`
	AuthorID   int64             `json:"author_id"`
	EntityID   int64             `json:"entity_id"`
	EntityType string            `json:"entity_type"`
	EventName  string            `json:"event_name"`
	Details    AuditEventDetails `json:"details"`
	CreatedAt  *time.Time        `json:"created_at"`
	EventType  string            `json:"event_type"`
}

// AuditEventDetails represents the details portion of an audit event for
// a group, a project or the instance. The exact fields that are returned
// for an audit event depend on the action being recorded.
//
// GitLab API docs: https://docs.gitlab.com/api/audit_events/
type AuditEventDetails struct {
	With          string `json:"with"`
	Add           string `json:"add"`
	As            string `json:"as"`
	Change        string `json:"change"`
	From          string `json:"from"`
	To            string `json:"to"`
	Remove        string `json:"remove"`
	CustomMessage string `json:"custom_message"`
	AuthorName    string `json:"author_name"`
	AuthorEmail   string `json:"author_email"`
	AuthorClass   string `json:"author_class"`
	TargetID      any    `json:"target_id"`
	TargetType    string `json:"target_type"`
	TargetDetails string `json:"target_details"`
	IPAddress     string `json:"ip_address"`
	EntityPath    string `json:"entity_path"`
	FailedLogin   string `json:"failed_login"`
	EventName     string `json:"event_name"`
}

// ListAuditEventsOptions represents the available ListProjectAuditEvents(),
// ListGroupAuditEvents() or ListInstanceAuditEvents() options.
//
// GitLab API docs: https://docs.gitlab.com/api/audit_events/
type ListAuditEventsOptions struct {
	ListOptions
	CreatedAfter  *time.Time `url:"created_after,omitempty" json:"created_after,omitempty"`
	CreatedBefore *time.Time `url:"created_before,omitempty" json:"created_before,omitempty"`
}

func (s *AuditEventsService) ListInstanceAuditEvents(opt *ListAuditEventsOptions, options ...RequestOptionFunc) ([]*AuditEvent, *Response, error) {
	return do[[]*AuditEvent](s.client,
		withMethod(http.MethodGet),
		withPath("audit_events"),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *AuditEventsService) GetInstanceAuditEvent(event int64, options ...RequestOptionFunc) (*AuditEvent, *Response, error) {
	return do[*AuditEvent](s.client,
		withMethod(http.MethodGet),
		withPath("audit_events/%d", event),
		withRequestOpts(options...),
	)
}

func (s *AuditEventsService) ListGroupAuditEvents(gid any, opt *ListAuditEventsOptions, options ...RequestOptionFunc) ([]*AuditEvent, *Response, error) {
	return do[[]*AuditEvent](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/audit_events", GroupID{gid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *AuditEventsService) GetGroupAuditEvent(gid any, event int64, options ...RequestOptionFunc) (*AuditEvent, *Response, error) {
	return do[*AuditEvent](s.client,
		withMethod(http.MethodGet),
		withPath("groups/%s/audit_events/%d", GroupID{gid}, event),
		withRequestOpts(options...),
	)
}

func (s *AuditEventsService) ListProjectAuditEvents(pid any, opt *ListAuditEventsOptions, options ...RequestOptionFunc) ([]*AuditEvent, *Response, error) {
	return do[[]*AuditEvent](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/audit_events", ProjectID{pid}),
		withAPIOpts(opt),
		withRequestOpts(options...),
	)
}

func (s *AuditEventsService) GetProjectAuditEvent(pid any, event int64, options ...RequestOptionFunc) (*AuditEvent, *Response, error) {
	return do[*AuditEvent](s.client,
		withMethod(http.MethodGet),
		withPath("projects/%s/audit_events/%d", ProjectID{pid}, event),
		withRequestOpts(options...),
	)
}
