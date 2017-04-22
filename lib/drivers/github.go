package drivers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Originate/git-town/lib/git"
)

type GithubApiPullRequest struct {
	Number int `json:number`
}

// GithubCodeHostingDriver provides tools for working with repositories
// hosted on Github
type GithubCodeHostingDriver struct{}

// GetNewPullRequestURL returns the URL of the page
// to create a new pull request on Github
func (driver GithubCodeHostingDriver) GetNewPullRequestURL(repository string, branch string, parentBranch string) string {
	toCompare := branch
	if parentBranch != git.GetMainBranch() {
		toCompare = parentBranch + "..." + branch
	}
	return fmt.Sprintf("https://github.com/%s/compare/%s?expand=1", repository, toCompare)
}

// GetRepositoryURL returns the URL of the given repository on github.com
func (driver GithubCodeHostingDriver) GetRepositoryURL(repository string) string {
	return "https://github.com/" + repository
}

func (driver GithubCodeHostingDriver) GetPullRequestNumber(repository string, branch string, parentBranch string) (int, error) {
	query := url.Values{}
	query.Add("access_token", os.Getenv("GITHUB_TOKEN"))
	query.Add("base", parentBranch)
	query.Add("head", strings.Split(repository, "/")[0]+":"+branch)
	query.Add("state", "open")
	getPullRequestUrl := fmt.Sprintf("https://api.github.com/repos/%s/pulls?%s", repository, query.Encode())
	response, err := http.Get(getPullRequestUrl)
	if err != nil {
		return -1, err
	}
	if response.StatusCode >= 400 {
		return -1, fmt.Errorf("Request was not successful: %v", response)
	}
	if response.Body == nil {
		return 1, errors.New("Missing response body")
	}
	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()

	_, err = decoder.Token()
	if err != nil {
		return -1, fmt.Errorf("error decoding open bracket: %v", err)
	}

	if decoder.More() {
		pullRequest := GithubApiPullRequest{}
		err := decoder.Decode(&pullRequest)
		if err != nil {
			return -1, fmt.Errorf("error decoding pull request object: %v", err)
		}
		return pullRequest.Number, nil
	}

	return -1, errors.New("No pull request found")
}

func (driver GithubCodeHostingDriver) MergePullRequest(repository string, options MergePullRequestOptions) error {
	query := url.Values{}
	query.Add("access_token", os.Getenv("GITHUB_TOKEN"))
	query.Add("commit_title", options.CommitTitle)
	query.Add("commit_message", options.CommitMessage)
	query.Add("sha", options.Sha)
	query.Add("merge_method", options.MergeMethod)
	mergePullRequestUrl := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/merge?%s", repository, options.Number, query.Encode())
	client := &http.Client{}
	request, err := http.NewRequest("PUT", mergePullRequestUrl, nil)
	if err != nil {
		log.Fatal(err)
	}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	if response.StatusCode >= 400 {
		return fmt.Errorf("Request was not successful: %v", err)
	}
	return nil
}
