package mockproposals_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/test/mockproposals"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestToDocString(t *testing.T) {
	t.Parallel()

	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		proposals := []forgedomain.ProposalData{}
		have := mockproposals.ToDocString(proposals)
		want := ""
		must.Eq(t, want, have)
	})

	t.Run("single proposal with multi-line body", func(t *testing.T) {
		t.Parallel()
		proposals := []forgedomain.ProposalData{
			{
				Number: forgedomain.ProposalNumber(789),
				Source: gitdomain.NewLocalBranchName("feature"),
				Target: gitdomain.NewLocalBranchName("main"),
				Body:   gitdomain.NewProposalBodyOpt("Line 1\nLine 2\nLine 3"),
				URL:    "https://example.com/pr/789",
			},
		}
		have := mockproposals.ToDocString(proposals)
		want := `number: 789
url: https://example.com/pr/789
source: feature
target: main
body:
  Line 1
  Line 2
  Line 3`
		must.Eq(t, want, have)
	})

	t.Run("single proposal with empty body", func(t *testing.T) {
		t.Parallel()
		proposals := []forgedomain.ProposalData{
			{
				Number: forgedomain.ProposalNumber(456),
				Source: gitdomain.NewLocalBranchName("bugfix"),
				Target: gitdomain.NewLocalBranchName("develop"),
				Body:   None[gitdomain.ProposalBody](),
				URL:    "https://example.com/pr/456",
			},
		}
		have := mockproposals.ToDocString(proposals)
		want := `number: 456
url: https://example.com/pr/456
source: bugfix
target: develop
body:`[1:]
		must.Eq(t, want, have)
	})

	t.Run("multiple proposals", func(t *testing.T) {
		t.Parallel()
		proposals := []forgedomain.ProposalData{
			{
				Number: forgedomain.ProposalNumber(1),
				Source: gitdomain.NewLocalBranchName("branch-1"),
				Target: gitdomain.NewLocalBranchName("main"),
				Body:   gitdomain.NewProposalBodyOpt("Body 1"),
				URL:    "https://example.com/pr/1",
			},
			{
				Number: forgedomain.ProposalNumber(2),
				Source: gitdomain.NewLocalBranchName("branch-2"),
				Target: gitdomain.NewLocalBranchName("main"),
				Body:   gitdomain.NewProposalBodyOpt("Body 2"),
				URL:    "https://example.com/pr/2",
			},
		}
		have := mockproposals.ToDocString(proposals)
		want := `
url: https://example.com/pr/1
number: 1
source: branch-1
target: main
body:
  Body 1

url: https://example.com/pr/2
number: 2
source: branch-2
target: main
body:
  Body 2`[1:]
		must.Eq(t, want, have)
	})

	t.Run("proposal with empty URL and body", func(t *testing.T) {
		t.Parallel()
		proposals := []forgedomain.ProposalData{
			{
				Number: forgedomain.ProposalNumber(999),
				Source: gitdomain.NewLocalBranchName("test"),
				Target: gitdomain.NewLocalBranchName("main"),
				Body:   gitdomain.NewProposalBodyOpt(""),
				URL:    "",
			},
		}
		have := mockproposals.ToDocString(proposals)
		want := `
url:
number: 999
source: test
target: main
body:`[1:]
		must.Eq(t, want, have)
	})
}
