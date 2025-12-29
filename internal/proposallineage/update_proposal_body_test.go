package proposallineage_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/proposallineage"
	"github.com/shoenig/test/must"
)

func TestUpdateProposalBody(t *testing.T) {
	t.Parallel()
	lineageSection := "main\n  - feat-a\n    - feat-b"

	t.Run("body with existing lineage", func(t *testing.T) {
		t.Parallel()
		body := gitdomain.ProposalBody(`
Proposal body text

<!-- branch-stack-start -->
main
  - old-a
    - old-b
<!-- branch-stack-end -->
`[1:])
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

	t.Run("body with marker and existing lineage", func(t *testing.T) {
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

	t.Run("body with marker", func(t *testing.T) {
		t.Parallel()
		body := gitdomain.ProposalBody(`
Proposal body text

<!-- branch-stack -->

### Next section

text
`[1:])
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

	t.Run("body without marker", func(t *testing.T) {
		t.Parallel()
		body := gitdomain.ProposalBody("Proposal body text")
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

	t.Run("empty body", func(t *testing.T) {
		t.Parallel()
		body := gitdomain.ProposalBody("")
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
}
