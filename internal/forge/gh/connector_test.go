package gh_test

import (
	"errors"
	"slices"
	"sync/atomic"
	"testing"

	"github.com/git-town/git-town/v24/internal/browser/browserdomain"
	"github.com/git-town/git-town/v24/internal/cli/print"
	"github.com/git-town/git-town/v24/internal/config/configdomain"
	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/forge/gh"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v24/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestConnectorSquashMergeProposal(t *testing.T) {
	t.Parallel()

	t.Run("multi-line message uses the first line as the subject and the rest as the body", func(t *testing.T) {
		t.Parallel()
		runner := &recordingRunner{}
		connector := gh.Connector{Frontend: runner}
		err := connector.SquashMergeProposal(forgedomain.ProposalNumber(1), Some(gitdomain.CommitMessage("my title\n\nmy body")))
		must.NoError(t, err)
		must.Eq(t, [][]string{{"gh", "pr", "merge", "--squash", "--subject=my title", "--body=my body", "1"}}, runner.calls)
	})

	t.Run("single-line message sends it as the commit subject", func(t *testing.T) {
		t.Parallel()
		runner := &recordingRunner{}
		connector := gh.Connector{Frontend: runner}
		err := connector.SquashMergeProposal(forgedomain.ProposalNumber(1), Some(gitdomain.CommitMessage("my title")))
		must.NoError(t, err)
		must.Eq(t, [][]string{{"gh", "pr", "merge", "--squash", "--subject=my title", "1"}}, runner.calls)
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

func TestConnectorCreateProposal(t *testing.T) {
	t.Parallel()

	t.Run("adds --fill when any prompt condition fails", func(t *testing.T) {
		t.Parallel()
		notInteractive := newTestConnector(true)
		notInteractive.Interactive = configdomain.Interactive("disabled in test")
		noTerminal := newTestConnector(true)
		noTerminal.StdioIsTerminal = false
		ghPromptDisabled := newTestConnector(true)
		ghPromptDisabled.GhPrompt = false
		promptSettingDisabled := newTestConnector(true)
		promptSettingDisabled.Backend = proposalQuerier{promptQueried: nil, promptSetting: Some(stringss.Trimmed("disabled"))}
		promptSettingUnreadable := newTestConnector(true)
		promptSettingUnreadable.Backend = proposalQuerier{promptQueried: nil, promptSetting: None[stringss.Trimmed]()}
		connectors := map[string]gh.Connector{
			"GH_PROMPT_DISABLED is set":         ghPromptDisabled,
			"gh's prompt setting can't be read": promptSettingUnreadable,
			"gh's prompt setting is disabled":   promptSettingDisabled,
			"no terminal":                       noTerminal,
			"not interactive":                   notInteractive,
		}
		for name, connector := range connectors {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				have := createProposalCalls(t, connector, None[gitdomain.ProposalTitle](), None[gitdomain.ProposalBody]())
				want := [][]string{
					{"gh", "pr", "create", "--head=feature", "--base=main", "--fill"},
					{"gh", "pr", "view"},
				}
				must.Eq(t, want, have)
			})
		}
	})

	t.Run("adds --fill when gh cannot prompt", func(t *testing.T) {
		t.Parallel()
		connector := newTestConnector(false)

		t.Run("missing body", func(t *testing.T) {
			t.Parallel()
			have := createProposalCalls(t, connector, Some(gitdomain.ProposalTitle("my title")), None[gitdomain.ProposalBody]())
			want := [][]string{
				{"gh", "pr", "create", "--head=feature", "--base=main", "--title=my title", "--fill"},
				{"gh", "pr", "view"},
			}
			must.Eq(t, want, have)
		})

		t.Run("missing title", func(t *testing.T) {
			t.Parallel()
			have := createProposalCalls(t, connector, None[gitdomain.ProposalTitle](), Some(gitdomain.ProposalBody("my body")))
			want := [][]string{
				{"gh", "pr", "create", "--head=feature", "--base=main", "--body=my body", "--fill"},
				{"gh", "pr", "view"},
			}
			must.Eq(t, want, have)
		})

		t.Run("missing title and body", func(t *testing.T) {
			t.Parallel()
			have := createProposalCalls(t, connector, None[gitdomain.ProposalTitle](), None[gitdomain.ProposalBody]())
			want := [][]string{
				{"gh", "pr", "create", "--head=feature", "--base=main", "--fill"},
				{"gh", "pr", "view"},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("doesn't read gh's prompt setting when title and body are provided", func(t *testing.T) {
		t.Parallel()
		querier := proposalQuerier{promptQueried: &atomic.Bool{}, promptSetting: None[stringss.Trimmed]()}
		connector := newTestConnector(true)
		connector.Backend = querier
		createProposalCalls(t, connector, Some(gitdomain.ProposalTitle("my title")), Some(gitdomain.ProposalBody("my body")))
		must.False(t, querier.promptQueried.Load())
	})

	t.Run("lets gh prompt when it can", func(t *testing.T) {
		t.Parallel()
		connector := newTestConnector(true)

		t.Run("missing body", func(t *testing.T) {
			t.Parallel()
			have := createProposalCalls(t, connector, Some(gitdomain.ProposalTitle("my title")), None[gitdomain.ProposalBody]())
			want := [][]string{
				{"gh", "pr", "create", "--head=feature", "--base=main", "--title=my title"},
				{"gh", "pr", "view"},
			}
			must.Eq(t, want, have)
		})

		t.Run("missing title", func(t *testing.T) {
			t.Parallel()
			have := createProposalCalls(t, connector, None[gitdomain.ProposalTitle](), Some(gitdomain.ProposalBody("my body")))
			want := [][]string{
				{"gh", "pr", "create", "--head=feature", "--base=main", "--body=my body"},
				{"gh", "pr", "view"},
			}
			must.Eq(t, want, have)
		})

		t.Run("missing title and body", func(t *testing.T) {
			t.Parallel()
			have := createProposalCalls(t, connector, None[gitdomain.ProposalTitle](), None[gitdomain.ProposalBody]())
			want := [][]string{
				{"gh", "pr", "create", "--head=feature", "--base=main"},
				{"gh", "pr", "view"},
			}
			must.Eq(t, want, have)
		})
	})

	t.Run("uses provided title and body", func(t *testing.T) {
		t.Parallel()
		connectors := map[string]gh.Connector{
			"gh can prompt":    newTestConnector(true),
			"gh cannot prompt": newTestConnector(false),
		}
		for name, connector := range connectors {
			t.Run(name, func(t *testing.T) {
				t.Parallel()
				have := createProposalCalls(t, connector, Some(gitdomain.ProposalTitle("my title")), Some(gitdomain.ProposalBody("my body")))
				want := [][]string{
					{"gh", "pr", "create", "--head=feature", "--base=main", "--title=my title", "--body=my body"},
					{"gh", "pr", "view"},
				}
				must.Eq(t, want, have)
			})
		}
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

// createProposalCalls creates a proposal with the given title and body and returns the gh commands that ran.
func createProposalCalls(t *testing.T, connector gh.Connector, title Option[gitdomain.ProposalTitle], body Option[gitdomain.ProposalBody]) [][]string {
	t.Helper()
	frontend := &recordingRunner{}
	connector.Frontend = frontend
	err := connector.CreateProposal(forgedomain.CreateProposalArgs{
		Branch:         "feature",
		FrontendRunner: nil,
		MainBranch:     "main",
		ParentBranch:   "main",
		ProposalBody:   body,
		ProposalTitle:  title,
	})
	must.NoError(t, err)
	return frontend.calls
}

func newTestConnector(canPrompt bool) gh.Connector {
	interactive := configdomain.InteractiveEnabled
	if !canPrompt {
		interactive = configdomain.Interactive("disabled in test")
	}
	return gh.Connector{
		Backend:         proposalQuerier{promptQueried: nil, promptSetting: Some(stringss.Trimmed("enabled"))},
		BrowserEnabled:  browserdomain.BrowserEnabled(false),
		Frontend:        nil,
		GhPrompt:        forgedomain.GhPrompt(canPrompt),
		Interactive:     interactive,
		Log:             print.Logger{},
		StdioIsTerminal: configdomain.StdioIsTerminal(canPrompt),
	}
}

type proposalQuerier struct {
	promptQueried *atomic.Bool             // records whether "gh config get prompt" ran
	promptSetting Option[stringss.Trimmed] // None makes "gh config get prompt" fail
}

func (self proposalQuerier) Query(_ string, _ ...string) (string, error) {
	return `[{"baseRefName":"main","body":"body","headRefName":"feature","mergeable":"MERGEABLE","number":123,"title":"title","url":"https://github.com/git-town/git-town/pull/123"}]`, nil
}

func (self proposalQuerier) QueryTrim(_ string, args ...string) (stringss.Trimmed, error) {
	if !slices.Equal(args, []string{"config", "get", "prompt"}) {
		return stringss.Trimmed(""), nil
	}
	if self.promptQueried != nil {
		self.promptQueried.Store(true)
	}
	setting, hasSetting := self.promptSetting.Get()
	if !hasSetting {
		return stringss.Trimmed(""), errors.New("cannot read gh config")
	}
	return setting, nil
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
