package bitbucketdatacenter

type PullRequestResponse struct {
	Values        []PullRequest `json:"values"`
	Size          int           `json:"size"`
	IsLastPage    bool          `json:"isLastPage"`
	NextPageStart int           `json:"nextPageStart"`
	Start         int           `json:"start"`
	Limit         int           `json:"limit"`
}

type User struct {
	Slug         string `json:"slug"`
	EmailAddress string `json:"emailAddress"`
	Active       bool   `json:"active"`
	Name         string `json:"name"`
	Id           int    `json:"id"`
	Type         string `json:"type"`
	DisplayName  string `json:"displayName"`
}

type Participant struct {
	User               User   `json:"user"`
	LastReviewedCommit string `json:"lastReviewedCommit"`
	Approved           bool   `json:"approved"`
	Status             string `json:"status"`
	Role               string `json:"role"`
}

type PullRequest struct {
	ClosedDate   int64         `json:"closedDate"`
	FromRef      Ref           `json:"fromRef"`
	Participants []Participant `json:"participants"`
	Reviewers    []Participant `json:"reviewers"`
	CreatedDate  int64         `json:"createdDate"`
	ToRef        Ref           `json:"toRef"`
	Draft        bool          `json:"draft"`
	UpdatedDate  int64         `json:"updatedDate"`
	Version      int           `json:"version"`
	Locked       bool          `json:"locked"`
	Description  string        `json:"description"`
	Closed       bool          `json:"closed"`
	Title        string        `json:"title"`
	Id           int           `json:"id"`
	State        string        `json:"state"`
	Open         bool          `json:"open"`
}

type Project struct {
	Avatar      string `json:"avatar"`
	Description string `json:"description"`
	Namespace   string `json:"namespace"`
	Scope       string `json:"scope"`
	Name        string `json:"name"`
	Key         string `json:"key"`
	Id          int    `json:"id"`
	Type        string `json:"type"`
	Public      bool   `json:"public"`
}

type Ref struct {
	DisplayId    string `json:"displayId"`
	LatestCommit string `json:"latestCommit"`
	Repository   struct {
		HierarchyId   string `json:"hierarchyId"`
		ScmId         string `json:"scmId"`
		Slug          string `json:"slug"`
		StatusMessage string `json:"statusMessage"`
		Archived      bool   `json:"archived"`
		Forkable      bool   `json:"forkable"`
		DefaultBranch string `json:"defaultBranch"`
		Partition     int    `json:"partition"`
		RelatedLinks  struct {
		} `json:"relatedLinks"`
		Project     Project `json:"project"`
		Description string  `json:"description"`
		Scope       string  `json:"scope"`
		Origin      struct {
			HierarchyId   string `json:"hierarchyId"`
			ScmId         string `json:"scmId"`
			Slug          string `json:"slug"`
			StatusMessage string `json:"statusMessage"`
			Archived      bool   `json:"archived"`
			Forkable      bool   `json:"forkable"`
			DefaultBranch string `json:"defaultBranch"`
			Partition     int    `json:"partition"`
			RelatedLinks  struct {
			} `json:"relatedLinks"`
			Project     Project `json:"project"`
			Description string  `json:"description"`
			Scope       string  `json:"scope"`
			Name        string  `json:"name"`
			Id          int     `json:"id"`
			State       string  `json:"state"`
			Public      bool    `json:"public"`
		} `json:"origin"`
		Name   string `json:"name"`
		Id     int    `json:"id"`
		State  string `json:"state"`
		Public bool   `json:"public"`
	} `json:"repository"`
	Id   string `json:"id"`
	Type string `json:"type"`
}
