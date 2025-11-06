package gitlab

import (
	"fmt"
	"net/http"
)

type (
	ProjectStatisticsServiceInterface interface {
		// Last30DaysStatistics gets the project statistics for the last 30 days.
		//
		// GitLab API docs: https://docs.gitlab.com/api/project_statistics/#get-the-statistics-of-the-last-30-days
		Last30DaysStatistics(pid any, options ...RequestOptionFunc) (*ProjectStatistics, *Response, error)
	}

	// ProjectStatisticsService handles communication with the project statistics related methods
	// of the GitLab API.
	//
	// GitLab API docs: https://docs.gitlab.com/api/project_statistics
	ProjectStatisticsService struct {
		client *Client
	}
)

// ProjectStatistics represents the Project Statistics.
//
// GitLab API docs: https://docs.gitlab.com/api/project_statistics
type ProjectStatistics struct {
	Fetches FetchStats `json:"fetches"`
}

type FetchStats struct {
	Total int64      `json:"total"`
	Days  []DayStats `json:"days"`
}

type DayStats struct {
	Count int64  `json:"count"`
	Date  string `json:"date"`
}

var _ ProjectStatisticsServiceInterface = (*ProjectStatisticsService)(nil)

func (s *ProjectStatisticsService) Last30DaysStatistics(pid any, options ...RequestOptionFunc) (*ProjectStatistics, *Response, error) {
	project, err := parseID(pid)
	if err != nil {
		return nil, nil, err
	}

	u := fmt.Sprintf("projects/%s/statistics",
		PathEscape(project),
	)

	req, err := s.client.NewRequest(http.MethodGet, u, nil, options)
	if err != nil {
		return nil, nil, err
	}

	projectStats := new(ProjectStatistics)
	resp, err := s.client.Do(req, projectStats)
	if err != nil {
		return nil, resp, err
	}

	return projectStats, resp, err
}
