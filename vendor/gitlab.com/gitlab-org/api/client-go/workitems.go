package gitlab

import (
	"errors"
	"strconv"
	"strings"
	"text/template"
	"time"
)

type (
	WorkItemsServiceInterface interface {
		GetWorkItem(fullPath string, iid int64, options ...RequestOptionFunc) (*WorkItem, *Response, error)
		ListWorkItems(fullPath string, opt *ListWorkItemsOptions, options ...RequestOptionFunc) ([]*WorkItem, *Response, error)
	}

	// WorkItemsService handles communication with the work item related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/graphql/reference/#workitem
	WorkItemsService struct {
		client *Client
	}
)

var _ WorkItemsServiceInterface = (*WorkItemsService)(nil)

// WorkItem represents a GitLab work item.
//
// GitLab API docs: https://docs.gitlab.com/api/graphql/reference/#workitem
type WorkItem struct {
	ID          int64
	IID         int64
	Type        string
	State       string
	Status      string
	Title       string
	Description string
	CreatedAt   *time.Time
	UpdatedAt   *time.Time
	ClosedAt    *time.Time
	WebURL      string
	Author      *BasicUser
	Assignees   []*BasicUser
}

func (wi WorkItem) GID() string {
	return gidGQL{
		Type:  "WorkItem",
		Int64: wi.ID,
	}.String()
}

// workItemTemplate defines the common fields for a work item in GraphQL queries.
// It's chained from userCoreBasicTemplate so nested templates work.
var workItemTemplate = template.Must(template.Must(userCoreBasicTemplate.Clone()).New("WorkItem").Parse(`
	id
	iid
	workItemType {
		name
	}
	state
	title
	description
	author {
		{{ template "UserCoreBasic" }}
	}
	createdAt
	updatedAt
	closedAt
	webUrl
	features {
		assignees {
			assignees {
				nodes {
					{{ template "UserCoreBasic" }}
				}
			}
		}
		status {
			status {
				name
			}
		}
	}
`))

// getWorkItemTemplate is chained from workItemTemplate so it has access to both
// UserCoreBasic and WorkItem templates.
var getWorkItemTemplate = template.Must(template.Must(workItemTemplate.Clone()).New("GetWorkItem").Parse(`
	query GetWorkItem($fullPath: ID!, $iid: String!) {
		namespace(fullPath: $fullPath) {
			workItem(iid: $iid) {
				{{ template "WorkItem" }}
			}
		}
	}
`))

// GetWorkItem gets a single work item.
//
// fullPath is the full path to either a group or project.
// iid is the internal ID of the work item.
//
// GitLab API docs: https://docs.gitlab.com/api/graphql/reference/#namespaceworkitem
func (s *WorkItemsService) GetWorkItem(fullPath string, iid int64, options ...RequestOptionFunc) (*WorkItem, *Response, error) {
	var queryBuilder strings.Builder
	if err := getWorkItemTemplate.Execute(&queryBuilder, nil); err != nil {
		return nil, nil, err
	}

	q := GraphQLQuery{
		Query: queryBuilder.String(),
		Variables: map[string]any{
			"fullPath": fullPath,
			"iid":      strconv.FormatInt(iid, 10),
		},
	}

	var result struct {
		Data struct {
			Namespace struct {
				WorkItem *workItemGQL `json:"workItem"`
			} `json:"namespace"`
		}
		GenericGraphQLErrors
	}

	resp, err := s.client.GraphQL.Do(q, &result, options...)
	if err != nil {
		return nil, resp, err
	}

	if len(result.Errors) != 0 {
		return nil, resp, &GraphQLResponseError{
			Err:    errors.New("GraphQL query failed"),
			Errors: result.GenericGraphQLErrors,
		}
	}

	wiQL := result.Data.Namespace.WorkItem
	if wiQL == nil {
		return nil, resp, ErrNotFound
	}

	return wiQL.unwrap(), resp, nil
}

// ListWorkItemsOptions represents the available ListWorkItems() options.
//
// GitLab API docs: https://docs.gitlab.com/api/graphql/reference/#namespaceworkitems
type ListWorkItemsOptions struct {
	AssigneeUsernames    []string
	AssigneeWildcardID   *string
	AuthorUsername       *string
	Confidential         *bool
	CRMContactID         *string
	CRMOrganizationID    *string
	HealthStatusFilter   *string
	IDs                  []string
	IIDs                 []string
	IncludeAncestors     *bool
	IncludeDescendants   *bool
	IterationCadenceID   []string
	IterationID          []string
	IterationWildcardID  *string
	LabelName            []string
	MilestoneTitle       []string
	MilestoneWildcardID  *string
	MyReactionEmoji      *string
	ParentIDs            []string
	ReleaseTag           []string
	ReleaseTagWildcardID *string
	State                *string
	Subscribed           *string
	Types                []string
	Weight               *string
	WeightWildcardID     *string

	// Time filters
	ClosedAfter   *time.Time
	ClosedBefore  *time.Time
	CreatedAfter  *time.Time
	CreatedBefore *time.Time
	DueAfter      *time.Time
	DueBefore     *time.Time
	UpdatedAfter  *time.Time
	UpdatedBefore *time.Time

	// Sorting
	Sort *string

	// Search
	Search *string
	In     []string

	// Pagination
	After  *string
	Before *string
	First  *int64
	Last   *int64
}

// listWorkItemsTemplate is chained from workItemTemplate so it has access to both
// UserCoreBasic and WorkItem templates.
var listWorkItemsTemplate = template.Must(template.Must(workItemTemplate.Clone()).New("ListWorkItems").Parse(`
	query ListWorkItems(
		$fullPath: ID!
		$assigneeUsernames: [String!]
		$assigneeWildcardId: AssigneeWildcardId
		$authorUsername: String
		$confidential: Boolean
		$crmContactId: String
		$crmOrganizationId: String
		$healthStatusFilter: HealthStatusFilter
		$ids: [WorkItemID!]
		$iids: [String!]
		$includeAncestors: Boolean
		$includeDescendants: Boolean
		$iterationCadenceId: [IterationsCadenceID!]
		$iterationId: [ID]
		$iterationWildcardId: IterationWildcardId
		$labelName: [String!]
		$milestoneTitle: [String!]
		$milestoneWildcardId: MilestoneWildcardId
		$myReactionEmoji: String
		$parentIds: [WorkItemID!]
		$releaseTag: [String!]
		$releaseTagWildcardId: ReleaseTagWildcardId
		$state: IssuableState
		$subscribed: SubscriptionStatus
		$types: [IssueType!]
		$weight: String
		$weightWildcardId: WeightWildcardId
		$closedAfter: Time
		$closedBefore: Time
		$createdAfter: Time
		$createdBefore: Time
		$dueAfter: Time
		$dueBefore: Time
		$updatedAfter: Time
		$updatedBefore: Time
		$sort: WorkItemSort
		$search: String
		$in: [IssuableSearchableField!]
		$after: String
		$before: String
		$first: Int
		$last: Int
	) {
		namespace(fullPath: $fullPath) {
			workItems(
				assigneeUsernames: $assigneeUsernames
				assigneeWildcardId: $assigneeWildcardId
				authorUsername: $authorUsername
				confidential: $confidential
				crmContactId: $crmContactId
				crmOrganizationId: $crmOrganizationId
				healthStatusFilter: $healthStatusFilter
				ids: $ids
				iids: $iids
				includeAncestors: $includeAncestors
				includeDescendants: $includeDescendants
				iterationCadenceId: $iterationCadenceId
				iterationId: $iterationId
				iterationWildcardId: $iterationWildcardId
				labelName: $labelName
				milestoneTitle: $milestoneTitle
				milestoneWildcardId: $milestoneWildcardId
				myReactionEmoji: $myReactionEmoji
				parentIds: $parentIds
				releaseTag: $releaseTag
				releaseTagWildcardId: $releaseTagWildcardId
				state: $state
				subscribed: $subscribed
				types: $types
				weight: $weight
				weightWildcardId: $weightWildcardId
				closedAfter: $closedAfter
				closedBefore: $closedBefore
				createdAfter: $createdAfter
				createdBefore: $createdBefore
				dueAfter: $dueAfter
				dueBefore: $dueBefore
				updatedAfter: $updatedAfter
				updatedBefore: $updatedBefore
				sort: $sort
				search: $search
				in: $in
				after: $after
				before: $before
				first: $first
				last: $last
			) {
				nodes {
					{{ template "WorkItem" }}
				}
				pageInfo {
					endCursor
					hasNextPage
					startCursor
					hasPreviousPage
				}
			}
		}
	}
`))

// ListWorkItems lists workitems in a given namespace (group or project).
//
// GitLab API docs: https://docs.gitlab.com/api/graphql/reference/#namespaceworkitems
func (s *WorkItemsService) ListWorkItems(fullPath string, opt *ListWorkItemsOptions, options ...RequestOptionFunc) ([]*WorkItem, *Response, error) {
	var queryBuilder strings.Builder

	if err := listWorkItemsTemplate.Execute(&queryBuilder, nil); err != nil {
		return nil, nil, err
	}

	vars := map[string]any{
		"fullPath":             fullPath,
		"assigneeUsernames":    opt.AssigneeUsernames,
		"assigneeWildcardId":   opt.AssigneeWildcardID,
		"authorUsername":       opt.AuthorUsername,
		"confidential":         opt.Confidential,
		"crmContactId":         opt.CRMContactID,
		"crmOrganizationId":    opt.CRMOrganizationID,
		"healthStatusFilter":   opt.HealthStatusFilter,
		"ids":                  opt.IDs,
		"iids":                 opt.IIDs,
		"includeAncestors":     opt.IncludeAncestors,
		"includeDescendants":   opt.IncludeDescendants,
		"iterationCadenceId":   opt.IterationCadenceID,
		"iterationId":          opt.IterationID,
		"iterationWildcardId":  opt.IterationWildcardID,
		"labelName":            opt.LabelName,
		"milestoneTitle":       opt.MilestoneTitle,
		"milestoneWildcardId":  opt.MilestoneWildcardID,
		"myReactionEmoji":      opt.MyReactionEmoji,
		"parentIds":            opt.ParentIDs,
		"releaseTag":           opt.ReleaseTag,
		"releaseTagWildcardId": opt.ReleaseTagWildcardID,
		"state":                opt.State,
		"subscribed":           opt.Subscribed,
		"types":                opt.Types,
		"weight":               opt.Weight,
		"weightWildcardId":     opt.WeightWildcardID,
		"closedAfter":          opt.ClosedAfter,
		"closedBefore":         opt.ClosedBefore,
		"createdAfter":         opt.CreatedAfter,
		"createdBefore":        opt.CreatedBefore,
		"dueAfter":             opt.DueAfter,
		"dueBefore":            opt.DueBefore,
		"updatedAfter":         opt.UpdatedAfter,
		"updatedBefore":        opt.UpdatedBefore,
		"sort":                 opt.Sort,
		"search":               opt.Search,
		"in":                   opt.In,
		"after":                opt.After,
		"before":               opt.Before,
		"first":                opt.First,
		"last":                 opt.Last,
	}

	query := GraphQLQuery{
		Query:     queryBuilder.String(),
		Variables: vars,
	}

	var result struct {
		Data struct {
			Namespace struct {
				WorkItems connectionGQL[workItemGQL] `json:"workItems"`
			} `json:"namespace"`
		}
		GenericGraphQLErrors
	}

	resp, err := s.client.GraphQL.Do(query, &result, options...)
	if err != nil {
		return nil, resp, err
	}

	if len(result.Errors) != 0 {
		return nil, resp, &GraphQLResponseError{
			Err:    errors.New("GraphQL query failed"),
			Errors: result.GenericGraphQLErrors,
		}
	}

	var ret []*WorkItem

	for _, wi := range result.Data.Namespace.WorkItems.Nodes {
		ret = append(ret, wi.unwrap())
	}

	resp.PageInfo = &result.Data.Namespace.WorkItems.PageInfo

	return ret, resp, nil
}

// workItemGQL represents the JSON structure returned by the GraphQL query.
// It is used to parse the response and convert it to the more user-friendly WorkItem type.
type workItemGQL struct {
	ID           gidGQL `json:"id"`
	IID          iidGQL `json:"iid"`
	WorkItemType struct {
		Name string `json:"name"`
	} `json:"workItemType"`
	State       string              `json:"state"`
	Title       string              `json:"title"`
	Description string              `json:"description"`
	CreatedAt   *time.Time          `json:"createdAt"`
	UpdatedAt   *time.Time          `json:"updatedAt"`
	ClosedAt    *time.Time          `json:"closedAt"`
	Author      userCoreBasicGQL    `json:"author"`
	Features    workItemFeaturesGQL `json:"features"`
	WebURL      string              `json:"webUrl"`
}

func (w workItemGQL) unwrap() *WorkItem {
	var assignees []*BasicUser

	for _, a := range w.Features.Assignees.Assignees.Nodes {
		assignees = append(assignees, a.unwrap())
	}

	return &WorkItem{
		ID:          w.ID.Int64,
		IID:         int64(w.IID),
		Type:        w.WorkItemType.Name,
		State:       w.State,
		Status:      w.Features.Status.Status.Name,
		Title:       w.Title,
		Description: w.Description,
		CreatedAt:   w.CreatedAt,
		UpdatedAt:   w.UpdatedAt,
		ClosedAt:    w.ClosedAt,
		WebURL:      w.WebURL,
		Author:      w.Author.unwrap(),
		Assignees:   assignees,
	}
}

type workItemFeaturesGQL struct {
	Assignees struct {
		Assignees struct {
			Nodes []userCoreBasicGQL `json:"nodes"`
		} `json:"assignees"`
	} `json:"assignees"`
	Status struct {
		Status struct {
			Name string
		}
	}
}
