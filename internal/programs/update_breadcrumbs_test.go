package programs_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/programs"
	"github.com/git-town/git-town/v22/internal/vm/opcodes"
	"github.com/git-town/git-town/v22/internal/vm/program"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestUpdateBreadcrumbsProgram(t *testing.T) {
	t.Parallel()

	main := gitdomain.NewLocalBranchName("main")
	branchOne := gitdomain.NewLocalBranchName("branch-1")
	branchTwo := gitdomain.NewLocalBranchName("branch-2")
	branchThree := gitdomain.NewLocalBranchName("branch-3")
	otherOne := gitdomain.NewLocalBranchName("other-1")
	otherTwo := gitdomain.NewLocalBranchName("other-2")

	t.Run("updates only the touched stack", func(t *testing.T) {
		t.Parallel()

		prog := NewMutable(&program.Program{})
		give := programs.UpdateBreadcrumbsArgs{
			Config: config.ValidatedConfig{
				NormalConfig: config.NormalConfig{
					Lineage: configdomain.NewLineageWith(configdomain.LineageData{
						branchOne:   main,
						branchThree: branchTwo,
						branchTwo:   branchOne,
						otherOne:    main,
						otherTwo:    otherOne,
					}),
					UnknownBranchType: configdomain.UnknownBranchType(configdomain.BranchTypeFeatureBranch),
				},
				ValidatedConfigData: configdomain.ValidatedConfigData{
					MainBranch: main,
				},
			},
			Program:         prog,
			TouchedBranches: gitdomain.LocalBranchNames{branchTwo},
		}

		programs.UpdateBreadcrumbsProgram(give)

		want := program.Program{
			&opcodes.ProposalUpdateBreadcrumb{Branch: branchOne},
			&opcodes.ProposalUpdateBreadcrumb{Branch: branchTwo},
			&opcodes.ProposalUpdateBreadcrumb{Branch: branchThree},
		}
		must.Eq(t, want, *prog.Value)
	})
}
