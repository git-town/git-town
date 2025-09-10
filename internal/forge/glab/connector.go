package glab

import (
	"errors"
	"regexp"
	"strings"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/gitlab"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

var (
	glabConnector Connector
	_             forgedomain.AuthVerifier = glabConnector
	_             forgedomain.Connector    = glabConnector
)

// Connector provides standardized connectivity for the given repository (github.com/owner/repo)
// via the GitHub API.
type Connector struct {
	Backend  subshelldomain.Querier
	Frontend subshelldomain.Runner
}

func (self Connector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	args := []string{"mr", "create", "--source-branch=" + data.Branch.String(), "--target-branch=" + data.ParentBranch.String()}
	title, hasTitle := data.ProposalTitle.Get()
	if hasTitle {
		args = append(args, "--title="+title.String())
	}
	body, hasBody := data.ProposalBody.Get()
	if hasBody {
		args = append(args, "--description="+body.String())
	}
	if !hasTitle || !hasBody {
		args = append(args, "--fill")
	}
	args = append(args, "--web")
	return self.Frontend.Run("glab", args...)
}

func (self Connector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return gitlab.DefaultProposalMessage(data)
}

func (self Connector) OpenRepository(runner subshelldomain.Runner) error {
	return runner.Run("glab", "repo", "view", "--web")
}

func (self Connector) VerifyCredentials() forgedomain.VerifyCredentialsResult {
	output, err := self.Backend.Query("glab", "auth", "status")
	if err != nil {
		return forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	return ParsePermissionsOutput(output)
}

func ParsePermissionsOutput(output string) forgedomain.VerifyCredentialsResult {
	result := forgedomain.VerifyCredentialsResult{
		AuthenticatedUser:   None[string](),
		AuthenticationError: nil,
		AuthorizationError:  nil,
	}
	lines := strings.Split(output, "\n")
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
