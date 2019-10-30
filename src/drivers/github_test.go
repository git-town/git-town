package drivers_test

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"testing"

	. "github.com/Originate/git-town/src/drivers"
	"github.com/stretchr/testify/assert"
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

var pullRequestBaseURL = "https://api.github.com/repos/Originate/git-town/pulls"
var currentPullRequestURL = pullRequestBaseURL + "?base=main&head=Originate%3Afeature&state=open"

func TestCodeHostingDriver_CanMergePullRequest_ReturnsFalseIfGithubTokenIsEmpty(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	driver := GetDriver(DriverOptions{OriginURL: "git@github.com:Originate/git-town.git"})
	assert.NotNil(t, driver)

	driver.SetAPIToken("")
	canMerge, _, err := driver.CanMergePullRequest("feature", "main")

	assert.Nil(t, err)
	assert.False(t, canMerge)
}

func TestCodeHostingDriver_CanMergePullRequest_ReturnsErrorGettingPullRequestNumber(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	driver := GetDriver(DriverOptions{OriginURL: "git@github.com:Originate/git-town.git"})
	assert.NotNil(t, driver)
	driver.SetAPIToken("TOKEN")
	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(404, ""))
	_, _, err := driver.CanMergePullRequest("feature", "main")
	assert.Error(t, err)
}

func TestCodeHostingDriver_CanMergePullRequest_ReturnsFalseIfNoPullRequestForBranch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	driver := GetDriver(DriverOptions{OriginURL: "git@github.com:Originate/git-town.git"})
	assert.NotNil(t, driver)
	driver.SetAPIToken("TOKEN")

	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[]"))
	canMerge, _, err := driver.CanMergePullRequest("feature", "main")

	assert.Nil(t, err)
	assert.False(t, canMerge)
}

func TestCodeHostingDriver_CanMergePullRequest_ReturnsFalseIfMultiplePullRequestsForBranch(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	driver := GetDriver(DriverOptions{OriginURL: "git@github.com:Originate/git-town.git"})
	assert.NotNil(t, driver)
	driver.SetAPIToken("TOKEN")

	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}, {"number": 2}]`))
	canMerge, _, err := driver.CanMergePullRequest("feature", "main")

	assert.Nil(t, err)
	assert.False(t, canMerge)
}

func TestCodeHostingDriver_CanMergePullRequest_OnePullRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	driver := GetDriver(DriverOptions{OriginURL: "git@github.com:Originate/git-town.git"})
	assert.NotNil(t, driver)
	driver.SetAPIToken("TOKEN")

	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1, "title": "my title" }]`))
	canMerge, defaultCommintMessage, err := driver.CanMergePullRequest("feature", "main")

	assert.Nil(t, err)
	assert.True(t, canMerge)
	assert.Equal(t, "my title (#1)", defaultCommintMessage)
}

var childPullRequestsURL = pullRequestBaseURL + "?base=feature&state=open"
var mergePullRequestURL = pullRequestBaseURL + "/1/merge"
var updatePullRequestBaseURL1 = pullRequestBaseURL + "/2"
var updatePullRequestBaseURL2 = pullRequestBaseURL + "/3"

func TestCodeHostingDriver_MergePullRequest_ReturnsRequestErrorForGetPullRequestIds(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	driver := GetDriver(DriverOptions{OriginURL: "git@github.com:Originate/git-town.git"})
	assert.NotNil(t, driver)
	options := MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	driver.SetAPIToken("TOKEN")

	httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(404, ""))
	_, err := driver.MergePullRequest(options)

	Expect(err).ToNot(BeNil())
}

func TestCodeHostingDriver_MergePullRequest_ReturnsRequestErrorForGetPullRequestToMerge(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	driver := GetDriver(DriverOptions{OriginURL: "git@github.com:Originate/git-town.git"})
	assert.NotNil(t, driver)
	options := MergePullRequestOptions{
		Branch:        "feature",
		CommitMessage: "title\nextra detail1\nextra detail2",
		ParentBranch:  "main",
	}
	driver.SetAPIToken("TOKEN")

	httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
	httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(404, ""))

	_, err := driver.MergePullRequest(options)
	assert.Error(t, err)
}

var _ = Describe("CodeHostingDriver - GitHub", func() {
	var driver CodeHostingDriver
	BeforeEach(func() {
		driver = GetDriver(DriverOptions{OriginURL: "git@github.com:Originate/git-town.git"})
		Expect(driver).NotTo(BeNil())
	})
	Describe("MergePullRequest", func() {
		var options MergePullRequestOptions
		BeforeEach(func() {
			options = MergePullRequestOptions{
				Branch:        "feature",
				CommitMessage: "title\nextra detail1\nextra detail2",
				ParentBranch:  "main",
			}
			driver.SetAPIToken("TOKEN")
		})

		It("returns an error if pull request number not found", func() {
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, "[]"))
			_, err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("no pull request found"))
		})

		It("returns an error if multiple pull request numbers not found", func() {
			httpmock.RegisterResponder("GET", childPullRequestsURL, httpmock.NewStringResponder(200, "[]"))
			httpmock.RegisterResponder("GET", currentPullRequestURL, httpmock.NewStringResponder(200, `[{"number": 1}, {"number": 2}]`))
			_, err := driver.MergePullRequest(options)
			Expect(err).ToNot(BeNil())
			Expect(err.Error()).To(Equal("multiple pull requests found: 1, 2"))
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
