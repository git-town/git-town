package drivers_test

import (
	"net/http"
	"os"

	. "github.com/Originate/git-town/src/drivers"
	httpmock "gopkg.in/jarcoal/httpmock.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Github", func() {
	var driver GithubCodeHostingDriver

	BeforeEach(func() {
		driver = GithubCodeHostingDriver{}
	})

	Describe("MergePullRequest", func() {
		childPullRequestsURL := "https://api.github.com/repos/Originate/git-town/pulls?access_token=TOKEN&base=feature&state=open"
		currentPullRequestURL := "https://api.github.com/repos/Originate/git-town/pulls?access_token=TOKEN&base=main&head=Originate%3Afeature&state=open"
		mergePullRequestURL := "https://api.github.com/repos/Originate/git-town/pulls/1/merge"
		updatePullRequestBaseURL1 := "https://api.github.com/repos/Originate/git-town/pulls/2/update?access_token=TOKEN&base=main"
		updatePullRequestBaseURL2 := "https://api.github.com/repos/Originate/git-town/pulls/3/update?access_token=TOKEN&base=main"
		var options MergePullRequestOptions

		BeforeEach(func() {
			options = MergePullRequestOptions{
				Branch:        "feature",
				CommitMessage: "message",
				CommitTitle:   "title",
				MergeMethod:   "squash",
				ParentBranch:  "main",
				Repository:    "Originate/git-town",
				Sha:           "sha",
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
			var mergeReq *http.Request
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[{\"number\": 1}]"))
			httpmock.RegisterResponder("PUT", mergePullRequestURL, func(req *http.Request) (*http.Response, error) {
				mergeReq = req
				return httpmock.NewStringResponse(200, ""), nil
			})
			err := driver.MergePullRequest(options)
			Expect(err).To(BeNil())
			Expect(mergeReq.FormValue("commit_message")).To(Equal("message"))
			Expect(mergeReq.FormValue("commit_title")).To(Equal("title"))
			Expect(mergeReq.FormValue("merge_method")).To(Equal("squash"))
			Expect(mergeReq.FormValue("sha")).To(Equal("sha"))
		})

		It("updates the base of child pull requests", func() {
			var updateBaseReq1, updateBaseReq2 *http.Request
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[{\"number\": 2}, {\"number\": 3}]"))
			httpmock.RegisterResponder("PATCH", updatePullRequestBaseURL1, func(req *http.Request) (*http.Response, error) {
				updateBaseReq1 = req
				return httpmock.NewStringResponse(200, ""), nil
			})
			httpmock.RegisterResponder("PATCH", updatePullRequestBaseURL2, func(req *http.Request) (*http.Response, error) {
				updateBaseReq2 = req
				return httpmock.NewStringResponse(200, ""), nil
			})
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[{\"number\": 1}]"))
			httpmock.RegisterResponder("PUT", mergePullRequestURL, httpmock.NewStringResponder(200, ""))
			err := driver.MergePullRequest(options)
			Expect(err).To(BeNil())
			Expect(updateBaseReq1.FormValue("base")).To(Equal("main"))
			Expect(updateBaseReq2.FormValue("base")).To(Equal("main"))
		})
	})
})
