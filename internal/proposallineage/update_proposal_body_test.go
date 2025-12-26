package proposallineage_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/proposallineage"
	"github.com/shoenig/test/must"
)

func TestUpdateProposalBody(t *testing.T) {
	t.Parallel()

	t.Run("append to end of empty body", func(t *testing.T) {
		t.Parallel()
		body := gitdomain.ProposalBody("")
		lineageSection := "main\n-  feat-a\n-    feat-b"
		have := proposallineage.UpdateProposalBody(body, lineageSection)
		want := gitdomain.ProposalBody("main\n- feat-a\n-   feat-b")
		must.EqOp(t, want, have)
	})

	t.Run("append to end of body without marker", func(t *testing.T) {
		t.Parallel()
		body := gitdomain.ProposalBody("Proposal body text")
		lineageSection := "main\n-  feat-a\n-    feat-b"
		have := proposallineage.UpdateProposalBody(body, lineageSection)
		want := gitdomain.ProposalBody("Proposal body text\n\nmain\n- feat-a\n-   feat-b")
		must.EqOp(t, want, have)
	})

	t.Run("insert after the marker", func(t *testing.T) {
		t.Parallel()
		body := gitdomain.ProposalBody("Proposal body text\n\n<!-- branch-stack -->\n\n### Other section\n\ntext\n")
		lineageSection := "main\n-  feat-a\n-    feat-b"
		have := proposallineage.UpdateProposalBody(body, lineageSection)
		want := gitdomain.ProposalBody("Proposal body text\n\n<!-- branch-stack -->\nmain\n- feat-a\n-   feat-b\n\n### Other section\n\ntext\n")
		must.EqOp(t, want, have)
	})
}
