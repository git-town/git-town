package forge_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/forge"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestUpdateProposalBodyUpdateWithStackLineage(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		description string
		currentBody string
		lineage     string
		want        string
	}{
		{
			description: "Empty body with empty lineage content",
			currentBody: "",
			lineage:     "",
			want:        "",
		},
		{
			description: "Non-Empty body with empty lineage content",
			currentBody: "Some description",
			lineage:     "",
			want:        "Some description",
		},
		{
			description: "Empty body with non-empty lineage content",
			currentBody: "",
			lineage: `
main
	feat-a
		feat-b
`,
			want: `<!-- branch-stack -->

main
	feat-a
		feat-b

<!-- branch-stack-end -->`,
		},
		{
			description: "Proposal Body with multiple line body",
			currentBody: `
Git-town is a town of Gitters.
These Gitters are the stackers of tomorrow`,
			lineage: `
main
	feat-a
		feat-b
`,
			want: `
Git-town is a town of Gitters.
These Gitters are the stackers of tomorrow

<!-- branch-stack -->

main
	feat-a
		feat-b

<!-- branch-stack-end -->`,
		},
		{
			description: "Proposal with template where branch-stack hidden comment is in the middle of the proposal body",
			currentBody: `
Git-town is a town of Gitters.
<!-- branch-stack -->
Please check the box that apply
- [ ] Add Tests
- [ ] Wrote Documentation
- [ ] Fixed Infra
`,
			lineage: `
main
	feat-a
		feat-b
`,
			want: `
Git-town is a town of Gitters.
<!-- branch-stack -->

main
	feat-a
		feat-b

<!-- branch-stack-end -->
Please check the box that apply
- [ ] Add Tests
- [ ] Wrote Documentation
- [ ] Fixed Infra
`,
		},
		{
			description: "Proposal existing proposal lineage in the middle of the body is updated",
			currentBody: `
Git-town is a town of Gitters.
<!-- branch-stack -->

main
	feat-a

<!-- branch-stack-end -->
Please check the box that apply
- [ ] Add Tests
- [ ] Wrote Documentation
- [ ] Fixed Infra
`,
			lineage: `
main
	feat-a
		feat-b
			feat-c
				feat-d
`,
			want: `
Git-town is a town of Gitters.
<!-- branch-stack -->

main
	feat-a
		feat-b
			feat-c
				feat-d

<!-- branch-stack-end -->
Please check the box that apply
- [ ] Add Tests
- [ ] Wrote Documentation
- [ ] Fixed Infra
`,
		},
		{
			description: "Proposal existing proposal lineage at the end of the body",
			currentBody: `
Git-town is a town of Gitters.
Please check the box that apply
- [ ] Add Tests
- [ ] Wrote Documentation
- [ ] Fixed Infra
<!-- branch-stack -->

main
	feat-a

<!-- branch-stack-end -->
`,
			lineage: `
main
	feat-a
		feat-b
			feat-c
				feat-d
`,
			want: `
Git-town is a town of Gitters.
Please check the box that apply
- [ ] Add Tests
- [ ] Wrote Documentation
- [ ] Fixed Infra
<!-- branch-stack -->

main
	feat-a
		feat-b
			feat-c
				feat-d

<!-- branch-stack-end -->
`,
		},
	} {
		t.Run(tc.description, func(t *testing.T) {
			t.Parallel()
			// act
			have := forge.ProposalBodyUpdateWithStackLineage(gitdomain.ProposalBody(tc.currentBody), tc.lineage)
			// assert
			must.EqOp(t, gitdomain.ProposalBody(tc.want), have)
		})
	}
}
