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
		getPullRequestsUrl := "https://api.github.com/repos/Originate/git-town/pulls?access_token=TOKEN&base=main&head=Originate%3Afeature&state=open"
		mergePullRequestUrl := "https://api.github.com/repos/Originate/git-town/pulls/1/merge"
		var options MergePullRequestOptions

		BeforeEach(func() {
			options = MergePullRequestOptions{
				Branch:        "feature",
				CommitMessage: "message",
				CommitTitle:   "title",
				MergeMethod:   "squash",
				ParentBranch:  "main",
				Sha:           "sha",
				Repository:    "Originate/git-town",
			}
			os.Setenv("GITHUB_TOKEN", "TOKEN")
		})

		AfterEach(func() {
			os.Unsetenv("GITHUB_TOKEN")
		})

		It("returns request errors (getting the pull request number)", func() {
			httpmock.RegisterResponder("GET", getPullRequestsUrl, httpmock.NewStringResponder(404, ""))
			err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
		})

		It("returns an error if pull request number not found", func() {
			httpmock.RegisterResponder("GET", getPullRequestsUrl, httpmock.NewStringResponder(200, "[]"))
			err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("No pull request found"))
		})

		Describe("pull request found", func() {
			BeforeEach(func() {
				httpmock.RegisterResponder("GET", getPullRequestsUrl, httpmock.NewStringResponder(200, "[{\"number\": 1}]"))
			})

			It("returns request errors (merging the pull request)", func() {
				httpmock.RegisterResponder("PUT", mergePullRequestUrl, httpmock.NewStringResponder(404, ""))
				err := driver.MergePullRequest(options)
				Expect(err).ToNot(BeNil())
			})

			It("returns no error if the request succeeds", func() {
				httpmock.RegisterResponder("PUT", mergePullRequestUrl, func(req *http.Request) (*http.Response, error) {
					Expect(req.FormValue("commit_message")).To(Equal("message"))
					Expect(req.FormValue("commit_title")).To(Equal("title"))
					Expect(req.FormValue("merge_method")).To(Equal("squash"))
					Expect(req.FormValue("sha")).To(Equal("sha"))
					return httpmock.NewStringResponse(200, ""), nil
				})
				err := driver.MergePullRequest(options)
				Expect(err).To(BeNil())
			})
		})
	})
})
