package drivers_test

import (
	"net/http"
	"os"

	. "github.com/Originate/git-town/lib/drivers"
	httpmock "gopkg.in/jarcoal/httpmock.v1"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Github", func() {
	var driver GithubCodeHostingDriver

	BeforeEach(func() {
		driver = GithubCodeHostingDriver{}
	})

	Describe("GetPullRequestNumber", func() {
		pullRequestUrl := "https://api.github.com/repos/Originate/git-town/pulls?access_token=TOKEN&base=main&head=Originate%3Afeature&state=open"

		BeforeEach(func() {
			os.Setenv("GITHUB_TOKEN", "TOKEN")
		})

		AfterEach(func() {
			os.Unsetenv("GITHUB_TOKEN")
		})

		It("returns request errors", func() {
			httpmock.RegisterResponder("GET", pullRequestUrl, httpmock.NewStringResponder(404, ""))
			_, err := driver.GetPullRequestNumber("Originate/git-town", "feature", "main")
			Expect(err).ToNot(BeNil())
		})

		It("returns an error if pull request not found", func() {
			httpmock.RegisterResponder("GET", pullRequestUrl, httpmock.NewStringResponder(200, "[]"))
			_, err := driver.GetPullRequestNumber("Originate/git-town", "feature", "main")
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("No pull request found"))
		})

		It("returns the number if pull request found", func() {
			httpmock.RegisterResponder("GET", pullRequestUrl, httpmock.NewStringResponder(200, "[{\"number\": 1}]"))
			number, err := driver.GetPullRequestNumber("Originate/git-town", "feature", "main")
			Expect(err).To(BeNil())
			Expect(number).To(Equal(1))
		})
	})

	Describe("MergePullRequest", func() {
		mergePullRequestUrl := "https://api.github.com/repos/Originate/git-town/pulls/1/merge"
		var options MergePullRequestOptions

		BeforeEach(func() {
			options = MergePullRequestOptions{
				CommitMessage: "message",
				CommitTitle:   "title",
				MergeMethod:   "squash",
				Number:        1,
				Sha:           "sha",
			}
			os.Setenv("GITHUB_TOKEN", "TOKEN")
		})

		AfterEach(func() {
			os.Unsetenv("GITHUB_TOKEN")
		})

		It("returns request errors", func() {
			httpmock.RegisterResponder("PUT", mergePullRequestUrl, httpmock.NewStringResponder(404, ""))
			err := driver.MergePullRequest("Originate/git-town", options)
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
			err := driver.MergePullRequest("Originate/git-town", options)
			Expect(err).To(BeNil())
		})
	})
})
