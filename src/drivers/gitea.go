package drivers

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"code.gitea.io/sdk/gitea"
	"github.com/git-town/git-town/v7/src/drivers/helpers"
	"golang.org/x/oauth2"
)

// GiteaCodeHostingDriver provides access to the API of Gitea installations.
type GiteaCodeHostingDriver struct {
	originURL  string
	hostname   string
	apiToken   string
	client     *gitea.Client
	owner      string
	repository string
	log        logFn
}

// LoadGitea provides a Gitea driver instance if the given repo configuration is for a Gitea repo,
// otherwise nil.
func LoadGitea(config config, log logFn) *GiteaCodeHostingDriver {
	driverType := config.GetCodeHostingDriverName()
	originURL := config.GetRemoteOriginURL()
	hostname := helpers.GetURLHostname(originURL)
	manualHostName := config.GetCodeHostingOriginHostname()
	if manualHostName != "" {
		hostname = manualHostName
	}
	if driverType != "gitea" && hostname != "gitea.com" {
		return nil
	}
	repositoryParts := strings.SplitN(helpers.GetURLRepositoryName(originURL), "/", 2)
	if len(repositoryParts) != 2 {
		return nil
	}
	owner := repositoryParts[0]
	repository := repositoryParts[1]
	return &GiteaCodeHostingDriver{
		originURL:  originURL,
		hostname:   hostname,
		apiToken:   config.GetGiteaToken(),
		log:        log,
		owner:      owner,
		repository: repository,
	}
}

func (d *GiteaCodeHostingDriver) LoadPullRequestInfo(branch, parentBranch string) (result PullRequestInfo, err error) {
	if d.apiToken == "" {
		return result, nil
	}
	d.connect()
	openPullRequests, err := d.client.ListRepoPullRequests(d.owner, d.repository, gitea.ListPullRequestsOptions{
		ListOptions: gitea.ListOptions{
			PageSize: 50,
		},
		State: gitea.StateOpen,
	})
	if err != nil {
		return result, err
	}
	baseName := parentBranch
	headName := d.owner + "/" + branch
	pullRequests := filterPullRequests(openPullRequests, baseName, headName)
	if len(pullRequests) != 1 {
		return result, nil
	}
	pullRequest := pullRequests[0]
	if !pullRequest.Mergeable {
		return result, nil
	}
	result.CanMergeWithAPI = true
	result.DefaultCommitMessage = getDefaultCommitMessage(pullRequest)
	result.PullRequestNumber = pullRequest.Index
	return result, nil
}

func (d *GiteaCodeHostingDriver) NewPullRequestURL(branch string, parentBranch string) (string, error) {
	toCompare := parentBranch + "..." + branch
	return fmt.Sprintf("%s/compare/%s", d.RepositoryURL(), url.PathEscape(toCompare)), nil
}

func (d *GiteaCodeHostingDriver) RepositoryURL() string {
	return fmt.Sprintf("https://%s/%s/%s", d.hostname, d.owner, d.repository)
}

func (d *GiteaCodeHostingDriver) HostingServiceName() string {
	return "Gitea"
}

func (d *GiteaCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) (mergeSha string, err error) {
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

// Helper

func (d *GiteaCodeHostingDriver) connect() {
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

func (d *GiteaCodeHostingDriver) apiMergePullRequest(pullRequestNumber int64, commitTitle, commitMessage string) (mergeSha string, err error) {
	_, err = d.client.MergePullRequest(d.owner, d.repository, pullRequestNumber, gitea.MergePullRequestOption{
		Style:   gitea.MergeStyleSquash,
		Title:   commitTitle,
		Message: commitMessage,
	})
	if err != nil {
		return "", err
	}
	pullRequest, err := d.client.GetPullRequest(d.owner, d.repository, pullRequestNumber)
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
//nolint:unparam
func (d *GiteaCodeHostingDriver) apiRetargetPullRequests(pullRequests []*gitea.PullRequest, newBaseName string) error {
	for _, pullRequest := range pullRequests {
		// RE-ENABLE AFTER https://github.com/go-gitea/gitea/issues/11552 and remove the nolint above
		// if options.LogRequests {
		// 	helpers.PrintLog(fmt.Sprintf("Gitea API: Updating base branch for PR #%d to #%s", *pullRequest.Index, newBaseName))
		// }
		d.log("Gitea API: Updating base branch for PR #%d to #%s", 1, newBaseName)
		d.log("The Gitea API currently does not support retargeting, please restarget #%d manually, see https://github.com/go-gitea/gitea/issues/11552", pullRequest.Index)
		// _, err = d.client.EditPullRequest(d.owner, d.repository, *pullRequest.Index, &gitea.EditPullRequestOption{
		// 	Base: newBaseName
		// })
		// if err != nil {
		//	return err
		// }
	}
	return nil
}
