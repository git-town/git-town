package gh_test

import (
	"errors"
	"testing"

	"github.com/git-town/git-town/v23/internal/browser/browserdomain"
	"github.com/git-town/git-town/v23/internal/cli/print"
	"github.com/git-town/git-town/v23/internal/forge/forgedomain"
	"github.com/git-town/git-town/v23/internal/forge/gh"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestConnectorCreateProposal(t *testing.T) {
	t.Parallel()

	t.Run("autofills missing body", func(t *testing.T) {
		t.Parallel()
		frontend := &recordingRunner{}
		connector := newTestConnector(frontend)
		args := createProposalArgs()
		args.ProposalTitle = Some(gitdomain.ProposalTitle("my title"))

		err := connector.CreateProposal(args)

		must.NoError(t, err)
		must.Eq(t, [][]string{
			{"gh", "pr", "create", "--head=feature", "--base=main", "--title=my title", "--fill"},
			{"gh", "pr", "view"},
		}, frontend.calls)
	})

	t.Run("autofills missing title", func(t *testing.T) {
		t.Parallel()
		frontend := &recordingRunner{}
		connector := newTestConnector(frontend)
		args := createProposalArgs()
		args.ProposalBody = Some(gitdomain.ProposalBody("my body"))

		err := connector.CreateProposal(args)

		must.NoError(t, err)
		must.Eq(t, [][]string{
			{"gh", "pr", "create", "--head=feature", "--base=main", "--body=my body", "--fill"},
			{"gh", "pr", "view"},
		}, frontend.calls)
	})

	t.Run("autofills missing title and body", func(t *testing.T) {
		t.Parallel()
		frontend := &recordingRunner{}
		connector := newTestConnector(frontend)
		args := createProposalArgs()

		err := connector.CreateProposal(args)

		must.NoError(t, err)
		must.Eq(t, [][]string{
			{"gh", "pr", "create", "--head=feature", "--base=main", "--fill"},
			{"gh", "pr", "view"},
		}, frontend.calls)
	})

	t.Run("uses provided title and body", func(t *testing.T) {
		t.Parallel()
		frontend := &recordingRunner{}
		connector := newTestConnector(frontend)
		args := createProposalArgs()
		args.ProposalTitle = Some(gitdomain.ProposalTitle("my title"))
		args.ProposalBody = Some(gitdomain.ProposalBody("my body"))

		err := connector.CreateProposal(args)

		must.NoError(t, err)
		must.Eq(t, [][]string{
			{"gh", "pr", "create", "--head=feature", "--base=main", "--title=my title", "--body=my body"},
			{"gh", "pr", "view"},
		}, frontend.calls)
	})
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

func createProposalArgs() forgedomain.CreateProposalArgs {
	return forgedomain.CreateProposalArgs{
		Branch:         "feature",
		FrontendRunner: nil,
		MainBranch:     "main",
		ParentBranch:   "main",
		ProposalBody:   None[gitdomain.ProposalBody](),
		ProposalTitle:  None[gitdomain.ProposalTitle](),
	}
}

func newTestConnector(frontend *recordingRunner) gh.Connector {
	return gh.Connector{
		Backend:        proposalQuerier{},
		BrowserEnabled: browserdomain.BrowserEnabled(false),
		Frontend:       frontend,
		Log:            print.Logger{},
	}
}

type proposalQuerier struct{}

func (self proposalQuerier) Query(_ string, _ ...string) (string, error) {
	return `[{"baseRefName":"main","body":"body","headRefName":"feature","mergeable":"MERGEABLE","number":123,"title":"title","url":"https://github.com/git-town/git-town/pull/123"}]`, nil
}

func (self proposalQuerier) QueryTrim(_ string, _ ...string) (stringss.Trimmed, error) {
	return stringss.Trimmed(""), nil
}

func (self proposalQuerier) QueryZ(_ string, _ ...string) (stringss.ZeroDelineated, error) {
	return stringss.ZeroDelineated(""), nil
}

type recordingRunner struct {
	calls [][]string
}

func (self *recordingRunner) Run(executable string, args ...string) error {
	call := append([]string{executable}, args...)
	self.calls = append(self.calls, call)
	return nil
}

func (self *recordingRunner) RunWithEnv(_ []string, executable string, args ...string) error {
	return self.Run(executable, args...)
}
