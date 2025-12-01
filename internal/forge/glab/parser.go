package glab

import (
	"encoding/json"
	"errors"
	"regexp"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func ParseJSONOutput(output string) ([]forgedomain.Proposal, error) {
	var parsed []jsonData
	err := json.Unmarshal([]byte(output), &parsed)
	if err != nil || len(parsed) == 0 {
		return []forgedomain.Proposal{}, err
	}
	result := make([]forgedomain.Proposal, len(parsed))
	for d, data := range parsed {
		result[d] = createProposal(data)
	}
	return result, nil
}

func ParsePermissionsOutput(output string) forgedomain.VerifyCredentialsResult {
	result := forgedomain.VerifyCredentialsResult{
		AuthenticatedUser:   None[string](),
		AuthenticationError: nil,
		AuthorizationError:  nil,
	}
	lines := stringslice.NonEmptyLines(output)
	regex := regexp.MustCompile(`Logged in to \S+ as (\S+) `)
	for _, line := range lines {
		matches := regex.FindStringSubmatch(line)
		if matches != nil {
			result.AuthenticatedUser = NewOption(matches[1])
			break
		}
	}
	if result.AuthenticatedUser.IsNone() {
		result.AuthenticationError = errors.New(messages.AuthenticationMissing)
	}
	return result
}

type jsonData struct {
	Description  gitdomain.ProposalBody  `json:"description"`
	Mergeable    string                  `json:"detailed_merge_status"` //nolint:tagliatelle
	Number       int                     `json:"iid"`                   //nolint:tagliatelle
	SourceBranch string                  `json:"source_branch"`         //nolint:tagliatelle
	State        string                  `json:"state"`                 //nolint:tagliatelle
	TargetBranch string                  `json:"target_branch"`         //nolint:tagliatelle
	Title        gitdomain.ProposalTitle `json:"title"`
	URL          string                  `json:"web_url"` //nolint:tagliatelle
}

func createProposal(data jsonData) forgedomain.Proposal {
	return forgedomain.Proposal{
		Data: forgedomain.ProposalData{
			Active:       data.State == "open",
			Body:         NewOption(data.Description),
			MergeWithAPI: data.Mergeable == "mergeable",
			Number:       data.Number,
			Source:       data.SourceBranch,
			Target:       data.TargetBranch,
			Title:        data.Title,
			URL:          data.URL,
		},
		ForgeType: forgedomain.ForgeTypeGitLab,
	}
}
