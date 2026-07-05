package gh_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/forge/gh"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestConnectorSquashMergeProposal(t *testing.T) {
	t.Parallel()

	t.Run("with a commit message sends it as the body", func(t *testing.T) {
		t.Parallel()
		runner := &recordingRunner{}
		connector := gh.Connector{Frontend: runner}
		err := connector.SquashMergeProposal(forgedomain.ProposalNumber(1), Some(gitdomain.CommitMessage("my message")))
		must.NoError(t, err)
		must.Eq(t, [][]string{{"gh", "pr", "merge", "--squash", "--body=my message", "1"}}, runner.calls)
	})

	t.Run("without a commit message lets the forge choose it", func(t *testing.T) {
		t.Parallel()
		runner := &recordingRunner{}
		connector := gh.Connector{Frontend: runner}
		err := connector.SquashMergeProposal(forgedomain.ProposalNumber(1), None[gitdomain.CommitMessage]())
		must.NoError(t, err)
		must.Eq(t, [][]string{{"gh", "pr", "merge", "--squash", "1"}}, runner.calls)
	})
}

// recordingRunner is a subshelldomain.Runner that records the commands it is asked to run.
type recordingRunner struct {
	calls [][]string
}

func (self *recordingRunner) Run(executable string, args ...string) error {
	self.calls = append(self.calls, append([]string{executable}, args...))
	return nil
}

func (self *recordingRunner) RunWithEnv(_ []string, executable string, args ...string) error {
	return self.Run(executable, args...)
}

func TestParsePermissionsOutput(t *testing.T) {
	t.Parallel()

	t.Run("logged into github.com with correct scopes", func(t *testing.T) {
		t.Parallel()
		give := `
github.com
  ✓ Logged in to github.com account kevgo (keyring)
  - Active account: true
  - Git operations protocol: ssh
  - Token: gho_************************************
  - Token scopes: 'gist', 'read:org', 'repo'`[1:]
		have := gh.ParsePermissionsOutput(give)
		want := forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   Some("kevgo"),
			AuthenticationError: nil,
			AuthorizationError:  nil,
		}
		must.Eq(t, want, have)
	})

	t.Run("logged into github.com with missing scopes", func(t *testing.T) {
		t.Parallel()
		give := `
github.com
  ✓ Logged in to github.com account kevgo (keyring)
  - Active account: true
  - Git operations protocol: ssh
  - Token: gho_************************************
  - Token scopes: 'gist', 'read:org'`[1:]
		have := gh.ParsePermissionsOutput(give)
		want := forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   Some("kevgo"),
			AuthenticationError: nil,
			AuthorizationError:  errors.New(`cannot find "repo" scope: ['gist' 'read:org']`),
		}
		must.Eq(t, want, have)
	})

	t.Run("not logged in", func(t *testing.T) {
		t.Parallel()
		give := "You are not logged into any GitHub hosts. To log in, run: gh auth login"
		have := gh.ParsePermissionsOutput(give)
		want := forgedomain.VerifyCredentialsResult{
			AuthenticatedUser:   None[string](),
			AuthenticationError: errors.New("not logged in"),
			AuthorizationError:  nil,
		}
		must.Eq(t, want, have)
	})
}
