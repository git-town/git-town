package gh

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/forge/github"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

var (
	ghConnector Connector
	_           forgedomain.Connector    = ghConnector
	_           forgedomain.AuthVerifier = ghConnector
)

// Connector provides standardized connectivity for the given repository (github.com/owner/repo)
// via the GitHub API.
type Connector struct {
	Backend  subshelldomain.Querier
	Frontend subshelldomain.Runner
}

func (self Connector) CreateProposal(data forgedomain.CreateProposalArgs) error {
	args := []string{"pr", "create", "--head=" + data.Branch.String(), "--base=" + data.ParentBranch.String()}
	if title, hasTitle := data.ProposalTitle.Get(); hasTitle {
		args = append(args, "--title="+title.String())
	}
	if body, hasBody := data.ProposalBody.Get(); hasBody {
		args = append(args, "--body="+body.String())
	}
	if err := self.Frontend.Run("gh", args...); err != nil {
		return err
	}
	return self.Frontend.Run("gh", "pr", "view", "--web")
}

func (self Connector) DefaultProposalMessage(data forgedomain.ProposalData) string {
	return github.DefaultProposalMessage(data)
}

func (self Connector) OpenRepository(runner subshelldomain.Runner) error {
	return runner.Run("gh", "browse")
}

func (self Connector) VerifyConnection() forgedomain.VerifyConnectionResult {
	output, err := self.Backend.Query("gh", "auth", "status", "--active")
	if err != nil {
		return forgedomain.VerifyConnectionResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: err,
			AuthorizationError:  nil,
		}
	}
	return ParsePermissionsOutput(output)
}

func ParsePermissionsOutput(output string) forgedomain.VerifyConnectionResult {
	result := forgedomain.VerifyConnectionResult{
		AuthenticatedUser:   None[string](),
		AuthenticationError: nil,
		AuthorizationError:  nil,
	}
	lines := strings.Split(output, "\n")
	regex := regexp.MustCompile(`Logged in to github.com account (\w+)`)
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
	regex = regexp.MustCompile(`Token scopes: (.+)`)
	for _, line := range lines {
		matches := regex.FindStringSubmatch(line)
		if matches != nil {
			parts := strings.Split(matches[1], ", ")
			if slices.Contains(parts, "'repo'") {
				break
			}
			result.AuthorizationError = fmt.Errorf(messages.AuthorizationMissing, parts)
		}
	}
	return result
}
