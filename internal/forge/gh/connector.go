package gh

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/github"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/messages"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Connector provides standardized connectivity for the given repository (github.com/owner/repo)
// via the GitHub API.
type Connector struct {
	forgedomain.Data
	runner Runner
	ghPath string // full path of the gh executable
	log    print.Logger
}

// NewConnector provides a fully configured gh.Connector instance
// if the current repo is hosted on GitHub, otherwise nil.
func NewConnector(args NewConnectorArgs) (Connector, error) {
	tokenSource := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: args.APIToken.String()})
	httpClient := oauth2.NewClient(context.Background(), tokenSource)
	githubClient := github.NewClient(httpClient)
	if args.RemoteURL.Host != "github.com" {
		url := "https://" + args.RemoteURL.Host
		var err error
		githubClient, err = githubClient.WithEnterpriseURLs(url, url)
		if err != nil {
			return Connector{}, fmt.Errorf(messages.GitHubEnterpriseInitializeError, err)
		}
	}
	return Connector{
		APIToken: args.APIToken,
		Data: forgedomain.Data{
			Hostname:     args.RemoteURL.Host,
			Organization: args.RemoteURL.Org,
			Repository:   args.RemoteURL.Repo,
		},
		client: githubClient,
		log:    args.Log,
	}, nil
}

type NewConnectorArgs struct {
	APIToken  Option[configdomain.GitHubToken]
	Log       print.Logger
	RemoteURL giturl.Parts
}

type Runner interface {
	Query(executable string, args ...string) (string, error)
}

func (self Connector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return github.DefaultProposalMessage(data)
}

func (self Connector) FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	return Some(self.findProposal)
}

func (self Connector) findProposal(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
	out, err := self.runner.Query(self.ghPath, "pr", "list", "--head="+branch.String(), "--base="+target.String(), "--json=number,title,body,mergeable,headRefName,baseRefName,url")
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	var parsed []ghData
	err = json.Unmarshal([]byte(out), &parsed)
	if err != nil {
		return None[forgedomain.Proposal](), err
	}
	if len(parsed) == 0 {
		return None[forgedomain.Proposal](), nil
	}
	if len(parsed) > 1 {
		return None[forgedomain.Proposal](), fmt.Errorf("found more than one pull request: %d", len(parsed))
	}
	pr := parsed[0]
	proposal := forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Body:         NewOption(pr.Body),
			MergeWithAPI: pr.Mergeable == "MERGEABLE",
			Number:       pr.Number,
			Source:       gitdomain.NewLocalBranchName(pr.HeadRefName),
			Target:       gitdomain.NewLocalBranchName(pr.BaseRefName),
			Title:        pr.Title,
			URL:          pr.URL,
		},
		ForgeType: forgedomain.ForgeTypeGitHub,
	}
	return Some(proposal), nil
}

type ghData struct {
	BaseRefName string `json:"baseRefName"`
	Body        string `json:"body"`
	HeadRefName string `json:"headRefName"`
	Mergeable   string `json:"mergeable"`
	Number      int    `json:"number"`
	Title       string `json:"title"`
	URL         string `json:"url"`
}

func (self Connector) NewProposalURL(data forgedomain.NewProposalURLData) (string, error) {
	return github.NewProposalURL(data, self.RepositoryURL())
}

func (self Connector) RepositoryURL() string {
	return github.RepositoryURL(self.HostnameWithStandardPort(), self.Organization, self.Repository)
}
