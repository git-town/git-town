//nolint:tagliatelle // we integrate with remote APIs only
package bitbucketdatacenter

type PullRequestResponse struct {
	IsLastPage    bool          `json:"isLastPage"`
	Limit         int           `json:"limit"`
	NextPageStart int           `json:"nextPageStart"`
	Size          int           `json:"size"`
	Start         int           `json:"start"`
	Values        []PullRequest `json:"values"`
}

type User struct {
	Active       bool   `json:"active"`
	DisplayName  string `json:"displayName"`
	EmailAddress string `json:"emailAddress"`
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

type Participant struct {
	Approved           bool   `json:"approved"`
	LastReviewedCommit string `json:"lastReviewedCommit"`
	Role               string `json:"role"`
	Status             string `json:"status"`
	User               User   `json:"user"`
}

type PullRequest struct {
	Closed       bool          `json:"closed"`
	ClosedDate   int64         `json:"closedDate"`
	CreatedDate  int64         `json:"createdDate"`
	Description  string        `json:"description"`
	Draft        bool          `json:"draft"`
	FromRef      Ref           `json:"fromRef"`
	ID           int           `json:"id"`
	Locked       bool          `json:"locked"`
	Open         bool          `json:"open"`
	Participants []Participant `json:"participants"`
	Reviewers    []Participant `json:"reviewers"`
	State        string        `json:"state"`
	Title        string        `json:"title"`
	ToRef        Ref           `json:"toRef"`
	UpdatedDate  int64         `json:"updatedDate"`
	Version      int           `json:"version"`
}

type Project struct {
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	ID          int    `json:"id"`
	Key         string `json:"key"`
	Name        string `json:"name"`
	Namespace   string `json:"namespace"`
	Public      bool   `json:"public"`
	Scope       string `json:"scope"`
	Type        string `json:"type"`
}

type Ref struct {
	DisplayID    string `json:"displayId"`
	ID           string `json:"id"`
	LatestCommit string `json:"latestCommit"`
	Repository   struct {
		Repository
		Origin Repository `json:"origin"`
	} `json:"repository"`
	Type string `json:"type"`
}

type Repository struct {
	Archived      bool     `json:"archived"`
	DefaultBranch string   `json:"defaultBranch"`
	Description   string   `json:"description"`
	Forkable      bool     `json:"forkable"`
	HierarchyID   string   `json:"hierarchyId"`
	ID            int      `json:"id"`
	Name          string   `json:"name"`
	Partition     int      `json:"partition"`
	Project       Project  `json:"project"`
	Public        bool     `json:"public"`
	RelatedLinks  struct{} `json:"relatedLinks"`
	ScmID         string   `json:"scmId"`
	Scope         string   `json:"scope"`
	Slug          string   `json:"slug"`
	State         string   `json:"state"`
	StatusMessage string   `json:"statusMessage"`
}
