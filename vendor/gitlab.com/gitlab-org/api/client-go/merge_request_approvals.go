//
// Copyright 2021, Sander van Harmelen
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
	"net/http"
	"time"
)

type (
	MergeRequestApprovalsServiceInterface interface {
		ApproveMergeRequest(pid any, mr int, opt *ApproveMergeRequestOptions, options ...RequestOptionFunc) (*MergeRequestApprovals, *Response, error)
		UnapproveMergeRequest(pid any, mr int, options ...RequestOptionFunc) (*Response, error)
		ResetApprovalsOfMergeRequest(pid any, mr int, options ...RequestOptionFunc) (*Response, error)
		GetConfiguration(pid any, mr int, options ...RequestOptionFunc) (*MergeRequestApprovals, *Response, error)
		ChangeApprovalConfiguration(pid any, mergeRequest int, opt *ChangeMergeRequestApprovalConfigurationOptions, options ...RequestOptionFunc) (*MergeRequest, *Response, error)
		GetApprovalRules(pid any, mergeRequest int, options ...RequestOptionFunc) ([]*MergeRequestApprovalRule, *Response, error)
		GetApprovalState(pid any, mergeRequest int, options ...RequestOptionFunc) (*MergeRequestApprovalState, *Response, error)
		CreateApprovalRule(pid any, mergeRequest int, opt *CreateMergeRequestApprovalRuleOptions, options ...RequestOptionFunc) (*MergeRequestApprovalRule, *Response, error)
		UpdateApprovalRule(pid any, mergeRequest int, approvalRule int, opt *UpdateMergeRequestApprovalRuleOptions, options ...RequestOptionFunc) (*MergeRequestApprovalRule, *Response, error)
		DeleteApprovalRule(pid any, mergeRequest int, approvalRule int, options ...RequestOptionFunc) (*Response, error)
	}

	// MergeRequestApprovalsService handles communication with the merge request
	// approvals related methods of the GitLab API. This includes reading/updating
	// approval settings and approve/unapproving merge requests
	//
	// GitLab API docs: https://docs.gitlab.com/api/merge_request_approvals/
	MergeRequestApprovalsService struct {
		client *Client
	}
)

var _ MergeRequestApprovalsServiceInterface = (*MergeRequestApprovalsService)(nil)

// MergeRequestApprovals represents GitLab merge request approvals.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#single-merge-request-approval
type MergeRequestApprovals struct {
	ID                             int                          `json:"id"`
	IID                            int                          `json:"iid"`
	ProjectID                      int                          `json:"project_id"`
	Title                          string                       `json:"title"`
	Description                    string                       `json:"description"`
	State                          string                       `json:"state"`
	CreatedAt                      *time.Time                   `json:"created_at"`
	UpdatedAt                      *time.Time                   `json:"updated_at"`
	MergeStatus                    string                       `json:"merge_status"`
	Approved                       bool                         `json:"approved"`
	ApprovalsBeforeMerge           int                          `json:"approvals_before_merge"`
	ApprovalsRequired              int                          `json:"approvals_required"`
	ApprovalsLeft                  int                          `json:"approvals_left"`
	RequirePasswordToApprove       bool                         `json:"require_password_to_approve"`
	ApprovedBy                     []*MergeRequestApproverUser  `json:"approved_by"`
	SuggestedApprovers             []*BasicUser                 `json:"suggested_approvers"`
	Approvers                      []*MergeRequestApproverUser  `json:"approvers"`
	ApproverGroups                 []*MergeRequestApproverGroup `json:"approver_groups"`
	UserHasApproved                bool                         `json:"user_has_approved"`
	UserCanApprove                 bool                         `json:"user_can_approve"`
	ApprovalRulesLeft              []*MergeRequestApprovalRule  `json:"approval_rules_left"`
	HasApprovalRules               bool                         `json:"has_approval_rules"`
	MergeRequestApproversAvailable bool                         `json:"merge_request_approvers_available"`
	MultipleApprovalRulesAvailable bool                         `json:"multiple_approval_rules_available"`
}

func (m MergeRequestApprovals) String() string {
	return Stringify(m)
}

// MergeRequestApproverGroup  represents GitLab project level merge request approver group.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#project-approval-rules
type MergeRequestApproverGroup struct {
	Group struct {
		ID                   int    `json:"id"`
		Name                 string `json:"name"`
		Path                 string `json:"path"`
		Description          string `json:"description"`
		Visibility           string `json:"visibility"`
		AvatarURL            string `json:"avatar_url"`
		WebURL               string `json:"web_url"`
		FullName             string `json:"full_name"`
		FullPath             string `json:"full_path"`
		LFSEnabled           bool   `json:"lfs_enabled"`
		RequestAccessEnabled bool   `json:"request_access_enabled"`
	}
}

// MergeRequestApprovalRule represents a GitLab merge request approval rule.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#get-merge-request-approval-rules
type MergeRequestApprovalRule struct {
	ID                   int                  `json:"id"`
	Name                 string               `json:"name"`
	RuleType             string               `json:"rule_type"`
	ReportType           string               `json:"report_type"`
	EligibleApprovers    []*BasicUser         `json:"eligible_approvers"`
	ApprovalsRequired    int                  `json:"approvals_required"`
	SourceRule           *ProjectApprovalRule `json:"source_rule"`
	Users                []*BasicUser         `json:"users"`
	Groups               []*Group             `json:"groups"`
	ContainsHiddenGroups bool                 `json:"contains_hidden_groups"`
	Section              string               `json:"section"`
	ApprovedBy           []*BasicUser         `json:"approved_by"`
	Approved             bool                 `json:"approved"`
}

// MergeRequestApprovalState represents a GitLab merge request approval state.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#get-the-approval-state-of-merge-requests
type MergeRequestApprovalState struct {
	ApprovalRulesOverwritten bool                        `json:"approval_rules_overwritten"`
	Rules                    []*MergeRequestApprovalRule `json:"rules"`
}

// String is a stringify for MergeRequestApprovalRule
func (s MergeRequestApprovalRule) String() string {
	return Stringify(s)
}

// MergeRequestApproverUser  represents GitLab project level merge request approver user.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#project-approval-rules
type MergeRequestApproverUser struct {
	User *BasicUser
}

// ApproveMergeRequestOptions represents the available ApproveMergeRequest() options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#approve-merge-request
type ApproveMergeRequestOptions struct {
	SHA *string `url:"sha,omitempty" json:"sha,omitempty"`
}

// ApproveMergeRequest approves a merge request on GitLab. If a non-empty sha
// is provided then it must match the sha at the HEAD of the MR.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#approve-merge-request
func (s *MergeRequestApprovalsService) ApproveMergeRequest(pid any, mr int, opt *ApproveMergeRequestOptions, options ...RequestOptionFunc) (*MergeRequestApprovals, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/approve", PathEscape(project), mr)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(MergeRequestApprovals)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// UnapproveMergeRequest unapproves a previously approved merge request on GitLab.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#unapprove-merge-request
func (s *MergeRequestApprovalsService) UnapproveMergeRequest(pid any, mr int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/unapprove", PathEscape(project), mr)

	req, err := s.client.NewRequest(http.MethodPost, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// ResetApprovalsOfMergeRequest clear all approvals of merge request on GitLab.
// Available only for bot users based on project or group tokens.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#reset-approvals-of-a-merge-request
func (s *MergeRequestApprovalsService) ResetApprovalsOfMergeRequest(pid any, mr int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/reset_approvals", PathEscape(project), mr)

	req, err := s.client.NewRequest(http.MethodPut, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}

// GetConfiguration shows information about single merge request approvals
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#single-merge-request-approval
func (s *MergeRequestApprovalsService) GetConfiguration(pid any, mr int, options ...RequestOptionFunc) (*MergeRequestApprovals, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/approvals", PathEscape(project), mr)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(MergeRequestApprovals)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// ChangeMergeRequestApprovalConfigurationOptions represents the available
// ChangeMergeRequestApprovalConfiguration() options.
//
// Deprecated: in GitLab 16.0
type ChangeMergeRequestApprovalConfigurationOptions struct {
	ApprovalsRequired *int `url:"approvals_required,omitempty" json:"approvals_required,omitempty"`
}

// ChangeApprovalConfiguration updates the approval configuration of a merge request.
//
// Deprecated: in GitLab 16.0
func (s *MergeRequestApprovalsService) ChangeApprovalConfiguration(pid any, mergeRequest int, opt *ChangeMergeRequestApprovalConfigurationOptions, options ...RequestOptionFunc) (*MergeRequest, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/approvals", PathEscape(project), mergeRequest)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	m := new(MergeRequest)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}

	return m, resp, nil
}

// GetApprovalRules requests information about a merge request’s approval rules
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#get-merge-request-approval-rules
func (s *MergeRequestApprovalsService) GetApprovalRules(pid any, mergeRequest int, options ...RequestOptionFunc) ([]*MergeRequestApprovalRule, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/approval_rules", PathEscape(project), mergeRequest)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var par []*MergeRequestApprovalRule
	resp, err := s.client.Do(req, &par)
	if err != nil {
		return nil, resp, err
	}

	return par, resp, nil
}

// GetApprovalState requests information about a merge request’s approval state
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#get-the-approval-state-of-merge-requests
func (s *MergeRequestApprovalsService) GetApprovalState(pid any, mergeRequest int, options ...RequestOptionFunc) (*MergeRequestApprovalState, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/approval_state", PathEscape(project), mergeRequest)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	var pas *MergeRequestApprovalState
	resp, err := s.client.Do(req, &pas)
	if err != nil {
		return nil, resp, err
	}

	return pas, resp, nil
}

// CreateMergeRequestApprovalRuleOptions represents the available CreateApprovalRule()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#create-merge-request-rule
type CreateMergeRequestApprovalRuleOptions struct {
	Name                  *string `url:"name,omitempty" json:"name,omitempty"`
	ApprovalsRequired     *int    `url:"approvals_required,omitempty" json:"approvals_required,omitempty"`
	ApprovalProjectRuleID *int    `url:"approval_project_rule_id,omitempty" json:"approval_project_rule_id,omitempty"`
	UserIDs               *[]int  `url:"user_ids,omitempty" json:"user_ids,omitempty"`
	GroupIDs              *[]int  `url:"group_ids,omitempty" json:"group_ids,omitempty"`
}

// CreateApprovalRule creates a new MR level approval rule.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#create-merge-request-rule
func (s *MergeRequestApprovalsService) CreateApprovalRule(pid any, mergeRequest int, opt *CreateMergeRequestApprovalRuleOptions, options ...RequestOptionFunc) (*MergeRequestApprovalRule, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/approval_rules", PathEscape(project), mergeRequest)

	req, err := s.client.NewRequest(http.MethodPost, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	par := new(MergeRequestApprovalRule)
	resp, err := s.client.Do(req, &par)
	if err != nil {
		return nil, resp, err
	}

	return par, resp, nil
}

// UpdateMergeRequestApprovalRuleOptions represents the available UpdateApprovalRule()
// options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#update-merge-request-rule
type UpdateMergeRequestApprovalRuleOptions struct {
	Name              *string `url:"name,omitempty" json:"name,omitempty"`
	ApprovalsRequired *int    `url:"approvals_required,omitempty" json:"approvals_required,omitempty"`
	UserIDs           *[]int  `url:"user_ids,omitempty" json:"user_ids,omitempty"`
	GroupIDs          *[]int  `url:"group_ids,omitempty" json:"group_ids,omitempty"`
}

// UpdateApprovalRule updates an existing approval rule with new options.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#update-merge-request-rule
func (s *MergeRequestApprovalsService) UpdateApprovalRule(pid any, mergeRequest int, approvalRule int, opt *UpdateMergeRequestApprovalRuleOptions, options ...RequestOptionFunc) (*MergeRequestApprovalRule, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/approval_rules/%d", PathEscape(project), mergeRequest, approvalRule)

	req, err := s.client.NewRequest(http.MethodPut, u, opt, options)
	if err != nil {
		return nil, nil, err
	}

	par := new(MergeRequestApprovalRule)
	resp, err := s.client.Do(req, &par)
	if err != nil {
		return nil, resp, err
	}

	return par, resp, nil
}

// DeleteApprovalRule deletes a mr level approval rule.
//
// GitLab API docs:
// https://docs.gitlab.com/api/merge_request_approvals/#delete-merge-request-rule
func (s *MergeRequestApprovalsService) DeleteApprovalRule(pid any, mergeRequest int, approvalRule int, options ...RequestOptionFunc) (*Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, err
	}
	u := fmt.Sprintf("projects/%s/merge_requests/%d/approval_rules/%d", PathEscape(project), mergeRequest, approvalRule)

	req, err := s.client.NewRequest(http.MethodDelete, u, nil, options)
	if err != nil {
		return nil, err
	}

	return s.client.Do(req, nil)
}
