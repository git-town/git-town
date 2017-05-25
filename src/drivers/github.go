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

	"github.com/Originate/git-town/src/git"
)

type githubAPIPullRequest struct {
	Number int `json:"number"`
}

// GithubCodeHostingDriver provides tools for working with repositories
// hosted on Github
type GithubCodeHostingDriver struct{}

// GetNewPullRequestURL returns the URL of the page
// to create a new pull request on Github
func (driver GithubCodeHostingDriver) GetNewPullRequestURL(repository string, branch, parentBranch string) string {
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

// MergePullRequest merges the pull request through the Github API
func (driver GithubCodeHostingDriver) MergePullRequest(options MergePullRequestOptions) error {
	childPullRequestNumbers, err := driver.getPullRequestNumbersAgainst(options.Repository, options.Branch)
	if err != nil {
		return err
	}
	for _, childPullRequestNumber := range childPullRequestNumbers {
		err = driver.updatePullRequestBase(options.Repository, childPullRequestNumber, options.ParentBranch)
		if err != nil {
			return err
		}
	}
	pullRequestNumber, err := driver.getPullRequestNumberFor(options.Repository, options.Branch, options.ParentBranch)
	if err != nil {
		return err
	}
	query := url.Values{}
	query.Add("access_token", os.Getenv("GITHUB_TOKEN"))
	query.Add("commit_title", options.CommitTitle)
	query.Add("commit_message", options.CommitMessage)
	query.Add("sha", options.Sha)
	query.Add("merge_method", options.MergeMethod)
	mergePullRequestURL := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/merge?%s", options.Repository, pullRequestNumber, query.Encode())
	client := &http.Client{}
	request, err := http.NewRequest("PUT", mergePullRequestURL, nil)
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

// Helpers

func (driver GithubCodeHostingDriver) getPullRequestNumberFor(repository string, branch string, parentBranch string) (int, error) {
	query := url.Values{}
	query.Add("access_token", os.Getenv("GITHUB_TOKEN"))
	query.Add("base", parentBranch)
	query.Add("head", strings.Split(repository, "/")[0]+":"+branch)
	query.Add("state", "open")
	getPullRequestURL := fmt.Sprintf("https://api.github.com/repos/%s/pulls?%s", repository, query.Encode())
	response, err := http.Get(getPullRequestURL)
	if err != nil {
		return -1, err
	}
	numbers, err := driver.parsePullRequestsResponse(response)
	if err != nil {
		return -1, err
	}
	if len(numbers) == 0 {
		return -1, errors.New("No pull request found")
	}
	if len(numbers) > 1 {
		return -1, fmt.Errorf("Multiple pull requests found: %s", strings.Trim(strings.Replace(fmt.Sprint(numbers), " ", ", ", -1), "[]"))
	}
	return numbers[0], nil
}

func (driver GithubCodeHostingDriver) getPullRequestNumbersAgainst(repository string, branch string) ([]int, error) {
	query := url.Values{}
	query.Add("access_token", os.Getenv("GITHUB_TOKEN"))
	query.Add("base", branch)
	query.Add("state", "open")
	getPullRequestURL := fmt.Sprintf("https://api.github.com/repos/%s/pulls?%s", repository, query.Encode())
	response, err := http.Get(getPullRequestURL)
	if err != nil {
		return []int{}, err
	}
	return driver.parsePullRequestsResponse(response)
}

func (driver GithubCodeHostingDriver) parsePullRequestsResponse(response *http.Response) ([]int, error) {
	if response.StatusCode >= 400 {
		return []int{}, fmt.Errorf("Request was not successful: %v", response)
	}
	if response.Body == nil {
		return []int{}, errors.New("Missing response body")
	}
	decoder := json.NewDecoder(response.Body)
	defer response.Body.Close()
	_, err := decoder.Token()
	if err != nil {
		return []int{}, fmt.Errorf("error decoding open bracket: %v", err)
	}
	result := []int{}
	for decoder.More() {
		pullRequest := githubAPIPullRequest{}
		err := decoder.Decode(&pullRequest)
		if err != nil {
			return []int{}, fmt.Errorf("error decoding pull request object: %v", err)
		}
		result = append(result, pullRequest.Number)
	}
	return result, nil
}

func (driver GithubCodeHostingDriver) updatePullRequestBase(repository string, pullRequestNumber int, base string) error {
	query := url.Values{}
	query.Add("access_token", os.Getenv("GITHUB_TOKEN"))
	query.Add("base", base)
	updatePullRequestURL := fmt.Sprintf("https://api.github.com/repos/%s/pulls/%d/update?%s", repository, pullRequestNumber, query.Encode())
	client := &http.Client{}
	request, err := http.NewRequest("PATCH", updatePullRequestURL, nil)
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
