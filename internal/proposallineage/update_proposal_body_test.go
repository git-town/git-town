package proposallineage_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/proposallineage"
	"github.com/shoenig/test/must"
)

func TestUpdateProposalBody(t *testing.T) {
	t.Parallel()

	t.Run("append to end of body without marker", func(t *testing.T) {
		t.Parallel()
		body := gitdomain.ProposalBody("Proposal body text")
		lineageSection := "main\n  - feat-a\n    - feat-b"
		have := proposallineage.UpdateProposalBody(body, lineageSection)
		want := gitdomain.ProposalBody(`
Proposal body text

<!-- branch-stack-start -->
main
  - feat-a
    - feat-b
<!-- branch-stack-end -->
`[1:])
		must.EqOp(t, want, have)
	})

	t.Run("append to end of empty body", func(t *testing.T) {
		t.Parallel()
		body := gitdomain.ProposalBody("")
		lineageSection := "main\n  - feat-a\n    - feat-b"
		have := proposallineage.UpdateProposalBody(body, lineageSection)
		want := gitdomain.ProposalBody(`
<!-- branch-stack-start -->
main
  - feat-a
    - feat-b
<!-- branch-stack-end -->
`[1:])
		must.EqOp(t, want, have)
	})

	t.Run("insert after the marker", func(t *testing.T) {
		t.Parallel()
		body := gitdomain.ProposalBody(`
Proposal body text

<!-- branch-stack -->

### Next section

text
`[1:])
		lineageSection := "main\n  - feat-a\n    - feat-b"
		have := proposallineage.UpdateProposalBody(body, lineageSection)
		want := gitdomain.ProposalBody(`
Proposal body text

<!-- branch-stack -->
<!-- branch-stack-start -->
main
  - feat-a
    - feat-b
<!-- branch-stack-end -->

### Next section

text
`[1:])
		must.EqOp(t, want, have)
	})

	t.Run("replace existing lineage", func(t *testing.T) {
		t.Parallel()
		body := gitdomain.ProposalBody(`
Proposal body text

<!-- branch-stack -->
<!-- branch-stack-start -->
main
  - old-a
    - old-b
<!-- branch-stack-end -->
`[1:])
		lineageSection := "main\n  - feat-a\n    - feat-b"
		have := proposallineage.UpdateProposalBody(body, lineageSection)
		want := gitdomain.ProposalBody(`
Proposal body text

<!-- branch-stack -->
<!-- branch-stack-start -->
main
  - feat-a
    - feat-b
<!-- branch-stack-end -->
`[1:])
		must.EqOp(t, want, have)
	})
}
