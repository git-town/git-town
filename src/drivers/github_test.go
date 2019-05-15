package drivers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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

var _ = Describe("CodeHostingDriver - GitHub", func() {
	pullRequestBaseURL := "https://api.github.com/repos/Originate/git-town/pulls"
	currentPullRequestURL := pullRequestBaseURL + "?base=main&head=Originate%3Afeature&state=open"
	var driver CodeHostingDriver

	BeforeEach(func() {
		driver = GetDriver(DriverOptions{OriginURL: "git@github.com:Originate/git-town.git"})
		Expect(driver).NotTo(BeNil())
	})

	Describe("CanMergePullRequest", func() {
		It("returns false if the environment variable GITHUB_TOKEN is an empty string", func() {
			driver.SetAPIToken("")
			canMerge, _, err := driver.CanMergePullRequest("feature", "main")
			Expect(err).To(BeNil())
			Expect(canMerge).To(BeFalse())
		})

		Describe("environment variable GITHUB_TOKEN is a non-empty string", func() {
			BeforeEach(func() {
				driver.SetAPIToken("TOKEN")
			})

			It("returns request errors (getting the pull request number to merge)", func() {
				httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(404, ""))
				_, _, err := driver.CanMergePullRequest("feature", "main")
				Expect(err).To(HaveOccurred())
			})

			It("returns false if there is no pull request for the branch", func() {
				httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[]"))
				canMerge, _, err := driver.CanMergePullRequest("feature", "main")
				Expect(err).To(BeNil())
				Expect(canMerge).To(BeFalse())
			})

			It("returns false if there are multiple pull requests for the branch", func() {
				httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}, {"number": 2}]`))
				canMerge, _, err := driver.CanMergePullRequest("feature", "main")
				Expect(err).To(BeNil())
				Expect(canMerge).To(BeFalse())
			})

			It("returns true (and the default commit message) if there is one pull request for the branch", func() {
				httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1, "title": "my title" }]`))
				canMerge, defaultCommintMessage, err := driver.CanMergePullRequest("feature", "main")
				Expect(err).To(BeNil())
				Expect(canMerge).To(BeTrue())
				Expect(defaultCommintMessage).To(Equal("my title (#1)"))
			})
		})
	})

	Describe("MergePullRequest", func() {
		childPullRequestsURL := pullRequestBaseURL + "?base=feature&state=open"
		mergePullRequestURL := pullRequestBaseURL + "/1/merge"
		updatePullRequestBaseURL1 := pullRequestBaseURL + "/2"
		updatePullRequestBaseURL2 := pullRequestBaseURL + "/3"
		var options MergePullRequestOptions

		BeforeEach(func() {
			options = MergePullRequestOptions{
				Branch:        "feature",
				CommitMessage: "title\nextra detail1\nextra detail2",
				ParentBranch:  "main",
			}
			driver.SetAPIToken("TOKEN")
		})

		It("returns request errors (getting the pull request numbers against the shipped branch)", func() {
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(404, ""))
			_, err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
		})

		It("returns request errors (getting the pull request number to merge)", func() {
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(404, ""))
			_, err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
		})

		It("returns an error if pull request number not found", func() {
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[]"))
			_, err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("No pull request found"))
		})

		It("returns an error if multiple pull request numbers not found", func() {
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}, {"number": 2}]`))
			_, err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("Multiple pull requests found: 1, 2"))
		})

		It("returns request errors (merging the pull request)", func() {
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}]`))
			httpmock.RegisterResponder("PUT", mergePullRequestURL, httpmock.NewStringResponder(404, ""))
			_, err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
		})

		It("merges the pull request", func() {
			var mergeRequest *http.Request
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}]`))
			httpmock.RegisterResponder("PUT", mergePullRequestURL, func(req *http.Request) (*http.Response, error) {
				mergeRequest = req
				return httpmock.NewStringResponse(200, `{"sha": "abc123"}`), nil
			})
			sha, err := driver.MergePullRequest(options)
			Expect(err).To(BeNil())
			Expect(sha).To(Equal("abc123"))
			mergeParameters := GetRequestData(mergeRequest)
			Expect(mergeParameters["commit_title"]).To(Equal("title"))
			Expect(mergeParameters["commit_message"]).To(Equal("extra detail1\nextra detail2"))
			Expect(mergeParameters["merge_method"]).To(Equal("squash"))
		})

		It("updates the base of child pull requests", func() {
			var updateRequest1, updateRequest2 *http.Request
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, `[{"number": 2}, {"number": 3}]`))
			httpmock.RegisterResponder("PATCH", updatePullRequestBaseURL1, func(req *http.Request) (*http.Response, error) {
				updateRequest1 = req
				return httpmock.NewStringResponse(200, ""), nil
			})
			httpmock.RegisterResponder("PATCH", updatePullRequestBaseURL2, func(req *http.Request) (*http.Response, error) {
				updateRequest2 = req
				return httpmock.NewStringResponse(200, ""), nil
			})
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}]`))
			httpmock.RegisterResponder("PUT", mergePullRequestURL, httpmock.NewStringResponder(200, `{"sha": "abc123"}`))
			_, err := driver.MergePullRequest(options)
			Expect(err).To(BeNil())
			updateParameters1 := GetRequestData(updateRequest1)
			Expect(updateParameters1["base"]).To(Equal("main"))
			updateParameters2 := GetRequestData(updateRequest2)
			Expect(updateParameters2["base"]).To(Equal("main"))
		})
	})
})
