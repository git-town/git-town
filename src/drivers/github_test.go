package drivers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	. "github.com/Originate/git-town/src/drivers"
	httpmock "gopkg.in/jarcoal/httpmock.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func GetRequestData(request *http.Request) map[string]interface{} {
	dataStr, err := ioutil.ReadAll(request.Body)
	if err != nil {
		panic(err)
	}
	data := map[string]interface{}{}
	err = json.Unmarshal(dataStr, &data)
	if err != nil {
		panic(err)
	}
	return data
}

var _ = Describe("Github", func() {
	var driver GithubCodeHostingDriver

	BeforeEach(func() {
		driver = GithubCodeHostingDriver{}
	})

	Describe("MergePullRequest", func() {
		pullRequestBaseURL := "https://api.github.com/repos/Originate/git-town/pulls"
		childPullRequestsURL := pullRequestBaseURL + "?base=feature&state=open"
		currentPullRequestURL := pullRequestBaseURL + "base=main&head=Originate%3Afeature&state=open"
		mergePullRequestURL := pullRequestBaseURL + "/1/merge"
		updatePullRequestBaseURL1 := pullRequestBaseURL + "/2"
		updatePullRequestBaseURL2 := pullRequestBaseURL + "/3"
		var options MergePullRequestOptions

		BeforeEach(func() {
			options = MergePullRequestOptions{
				Branch:        "feature",
				CommitMessage: "message",
				ParentBranch:  "main",
				Owner:         "Originate",
				Repository:    "git-town",
			}
			os.Setenv("GITHUB_TOKEN", "TOKEN")
		})

		AfterEach(func() {
			os.Unsetenv("GITHUB_TOKEN")
		})

		It("returns request errors (getting the pull request numbers against the shipped branch)", func() {
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(404, ""))
			err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
		})

		It("returns request errors (getting the pull request number to merge)", func() {
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(404, ""))
			err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
		})

		It("returns an error if pull request number not found", func() {
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[]"))
			err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("No pull request found"))
		})

		It("returns an error if multiple pull request numbers not found", func() {
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[{\"number\": 1}, {\"number\": 2}]"))
			err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("Multiple pull requests found: 1, 2"))
		})

		It("returns request errors (merging the pull request)", func() {
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[{\"number\": 1}]"))
			httpmock.RegisterResponder("PUT", mergePullRequestURL, httpmock.NewStringResponder(404, ""))
			err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
		})

		It("merges the pull request", func() {
			var mergeRequest *http.Request
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[{\"number\": 1}]"))
			httpmock.RegisterResponder("PUT", mergePullRequestURL, func(req *http.Request) (*http.Response, error) {
				mergeRequest = req
				return httpmock.NewStringResponse(200, ""), nil
			})
			err := driver.MergePullRequest(options)
			Expect(err).To(BeNil())
			mergeParameters := GetRequestData(mergeRequest)
			Expect(mergeParameters["commit_message"]).To(Equal("message"))
			Expect(mergeParameters["merge_method"]).To(Equal("squash"))
		})

		It("updates the base of child pull requests", func() {
			var updateRequest1, updateRequest2 *http.Request
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[{\"number\": 2}, {\"number\": 3}]"))
			httpmock.RegisterResponder("PATCH", updatePullRequestBaseURL1, func(req *http.Request) (*http.Response, error) {
				updateRequest1 = req
				return httpmock.NewStringResponse(200, ""), nil
			})
			httpmock.RegisterResponder("PATCH", updatePullRequestBaseURL2, func(req *http.Request) (*http.Response, error) {
				updateRequest2 = req
				return httpmock.NewStringResponse(200, ""), nil
			})
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[{\"number\": 1}]"))
			httpmock.RegisterResponder("PUT", mergePullRequestURL, httpmock.NewStringResponder(200, ""))
			err := driver.MergePullRequest(options)
			Expect(err).To(BeNil())
			updateParameters1 := GetRequestData(updateRequest1)
			Expect(updateParameters1["base"]).To(Equal("main"))
			updateParameters2 := GetRequestData(updateRequest2)
			Expect(updateParameters2["base"]).To(Equal("main"))
		})
	})
})
