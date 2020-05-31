package drivers

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/oauth2"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/src/git"
)

type giteaCodeHostingDriver struct {
	originURL  string
	hostname   string
	apiToken   string
	client     *gitea.Client
	owner      string
	repository string
}

func (d *giteaCodeHostingDriver) CanBeUsed(driverType string) bool {
	return driverType == "gitea" || d.hostname == "gitea.com"
}

func (d *giteaCodeHostingDriver) CanMergePullRequest(branch, parentBranch string) (canMerge bool, defaultCommitMessage string, pullRequestNumber int64, err error) {
	if d.apiToken == "" {
		return false, "", 0, nil
	}
	d.connect()
	openPullRequests, err := d.client.ListRepoPullRequests(d.owner, d.repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		return false, "", 0, err
	}
	baseName := parentBranch
	headName := d.owner + "/" + branch
	pullRequest, err := identifyPullRequest(filterPullRequests(openPullRequests, baseName, headName))
	if err != nil {
		return false, "", 0, nil
	}
	return true, getDefaultCommitMessage(pullRequest), int(pullRequest.Index), nil
}

func (d *giteaCodeHostingDriver) GetAPIToken() string {
	return git.Config().GetGiteaToken()
}

func (d *giteaCodeHostingDriver) GetNewPullRequestURL(branch string, parentBranch string) string {
	toCompare := parentBranch + "..." + branch
	return fmt.Sprintf("%s/compare/%s", d.GetRepositoryURL(), url.PathEscape(toCompare))
}

func (d *giteaCodeHostingDriver) GetRepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", d.hostname, d.owner, d.repository)
}

func (d *giteaCodeHostingDriver) HostingServiceName() string {
	return "Gitea"
}

func (d *giteaCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
	d.connect()
	openPullRequests, err := d.client.ListRepoPullRequests(d.owner, d.repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		return "", err
	}
	baseName := options.Branch
	newBaseName := options.ParentBranch
	err = d.apiRetargetPullRequests(filterPullRequests(openPullRequests, baseName, ""), newBaseName)
	if err != nil {
		return "", err
	}

	commitMessageParts := strings.SplitN(options.CommitMessage, "\n", 2)
	commitTitle := commitMessageParts[0]
	commitMessage := ""
	if len(commitMessageParts) == 2 {
		commitMessage = commitMessageParts[1]
	}
	return d.apiMergePullRequest(options.PullRequestNumber, commitTitle, commitMessage)
}

func (d *giteaCodeHostingDriver) SetAPIToken(apiToken string) {
	d.apiToken = apiToken
}

func (d *giteaCodeHostingDriver) SetOriginHostname(originHostname string) {
	d.hostname = originHostname
}

func (d *giteaCodeHostingDriver) SetOriginURL(originURL string) {
	d.originURL = originURL
	d.hostname = git.Config().GetURLHostname(originURL)
	d.client = nil
	repositoryParts := strings.SplitN(git.Config().GetURLRepositoryName(originURL), "/", 2)
	if len(repositoryParts) == 2 {
		d.owner = repositoryParts[0]
		d.repository = repositoryParts[1]
	}
}

func init() {
	registry.RegisterDriver(&giteaCodeHostingDriver{})
}

// Helper

func (d *giteaCodeHostingDriver) connect() {
	if d.client == nil {
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: d.apiToken},
		)
		tc := oauth2.NewClient(context.Background(), ts)
		d.client = gitea.NewClientWithHTTP(fmt.Sprintf("https://%s", d.hostname), tc)
	}
}

func getDefaultCommitMessage(pullRequest *gitea.PullRequest) string {
	return fmt.Sprintf("%s (#%d)", pullRequest.Title, pullRequest.Index)
}

func identifyPullRequest(filteredPullRequests []*gitea.PullRequest) (*gitea.PullRequest, error) {
	if len(filteredPullRequests) == 0 {
		return nil, errors.New("no pull request found")
	}
	if len(filteredPullRequests) > 1 {
		pullRequestNumbersAsStrings := make([]string, len(filteredPullRequests))
		for i, filteredPullRequest := range filteredPullRequests {
			pullRequestNumbersAsStrings[i] = strconv.FormatInt(filteredPullRequest.Index, 10)

		}
		return nil, fmt.Errorf("multiple pull requests found: %s", strings.Join(pullRequestNumbersAsStrings, ", "))
	}

	return filteredPullRequests[0], nil
}

func filterPullRequests(pullRequests []*gitea.PullRequest, baseName, headName string) []*gitea.PullRequest {
	pullRequestsFiltered := []*gitea.PullRequest{}
	for _, pullRequest := range pullRequests {
		if pullRequest.Base.Name != baseName {
			break
		}
		if headName != "" && pullRequest.Head.Name != headName {
			break
		}
		pullRequestsFiltered = append(pullRequestsFiltered, pullRequest)
	}
	return pullRequestsFiltered
}

func (d *giteaCodeHostingDriver) apiMergePullRequest(pullRequestNumber int64, commitTitle, commitMessage string) (mergeSha string, err error) {
	printLog(fmt.Sprintf("Gitea API: Merging PR #%d", pullRequestNumber))
	_, err = d.client.MergePullRequest(d.owner, d.repository, int64(pullRequestNumber), gitea.MergePullRequestOption{
		Style:   gitea.MergeStyleSquash,
		Title:   commitTitle,
		Message: commitMessage,
	})
	if err != nil {
		return "", err
	}
	pullRequest, err := d.client.GetPullRequest(d.owner, d.repository, int64(pullRequestNumber))
	if err != nil {
		return "", err
	}
	return *pullRequest.MergedCommitID, nil
}

// retargetPullRequests retargets pullrequests onto a new base branch
// this comes in handy when an ancestor got merged, so that children can be retargeted to the ancestor's own target branch
// example:
//   ancerstor -> master
//   children1 -> ancestor  --> master (retargeted to master after merge)
//   children2 -> ancestor  --> master (retargeted to master after merge)
func (d *giteaCodeHostingDriver) apiRetargetPullRequests(pullRequests []*gitea.PullRequest, newBaseName string) error {
	for _, pullRequest := range pullRequests {
		// if options.LogRequests {
		// 	printLog(fmt.Sprintf("Gitea API: Updating base branch for PR #%d", *pullRequest.Index))
		// }
		printLog(fmt.Sprintf("The Gitea API currently does not support retargeting, please restarget #%d manually, see https://github.com/go-gitea/gitea/issues/11552", pullRequest.Index))
		// _, err = d.client.EditPullRequest(d.owner, d.repository, *pullRequest.Index, &gitea.EditPullRequestOption{
		// 	Base: newBaseName
		// })
		// if err != nil {
		//	return err
		// }
	}
	return nil
}
