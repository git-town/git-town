package configdomain_test

import (
	"testing"

	"github.com/shoenig/test/must"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type mockConnectorProposalStackLineageBuilder struct{}

var _ forgedomain.Connector = (*mockConnectorProposalStackLineageBuilder)(nil)

func (self *mockConnectorProposalStackLineageBuilder) CreateProposal(_ forgedomain.CreateProposalArgs) error {
	return nil
}

func (self *mockConnectorProposalStackLineageBuilder) DefaultProposalMessage(_ forgedomain.ProposalData) string {
	return "mock"
}

func (self *mockConnectorProposalStackLineageBuilder) FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	return Some(func(branch, target gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error) {
		return Some(forgedomain.Proposal{
			Data: forgedomain.ProposalData{
				Body:         None[string](),
				MergeWithAPI: false,
				Number:       1,
				Source:       branch,
				Target:       target,
				Title:        "Test Mocker",
				URL:          "https://www.github.com/git-town/pull/1",
			},
			ForgeType: forgedomain.ForgeTypeCodeberg,
		}), nil
	})
}

func (self *mockConnectorProposalStackLineageBuilder) OpenRepository(_ subshelldomain.Runner) error {
	return nil
}

func (self *mockConnectorProposalStackLineageBuilder) SearchProposalFn() Option[func(_ gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)] {
	return None[func(_ gitdomain.LocalBranchName) (Option[forgedomain.Proposal], error)]()
}

func (self *mockConnectorProposalStackLineageBuilder) SquashMergeProposalFn() Option[func(_ int, _ gitdomain.CommitMessage) error] {
	return None[func(_ int, _ gitdomain.CommitMessage) error]()
}

func (self *mockConnectorProposalStackLineageBuilder) UpdateProposalBodyFn() Option[func(_ forgedomain.ProposalInterface, _ string) error] {
	return None[func(_ forgedomain.ProposalInterface, _ string) error]()
}

func (self *mockConnectorProposalStackLineageBuilder) UpdateProposalSourceFn() Option[func(_ forgedomain.ProposalInterface, _ gitdomain.LocalBranchName) error] {
	return None[func(_ forgedomain.ProposalInterface, _ gitdomain.LocalBranchName) error]()
}

func (self *mockConnectorProposalStackLineageBuilder) UpdateProposalTargetFn() Option[func(_ forgedomain.ProposalInterface, _ gitdomain.LocalBranchName) error] {
	return None[func(_ forgedomain.ProposalInterface, _ gitdomain.LocalBranchName) error]()
}

func (self *mockConnectorProposalStackLineageBuilder) VerifyConnection() forgedomain.VerifyConnectionResult {
	return forgedomain.VerifyConnectionResult{
		AuthenticatedUser:   None[string](),
		AuthenticationError: nil,
		AuthorizationError:  nil,
	}
}

func TestProposalStackLineageBuilder_CheckLineageAndProposals(t *testing.T) {
	t.Parallel()
	// arrange
	mainBranch := gitdomain.NewLocalBranchName("main")
	perennialBranch := gitdomain.NewLocalBranchName("development")
	featureBranch := gitdomain.NewLocalBranchName("git-town/proposal-stack-lineage")
	lineage := configdomain.NewLineage()
	lineage.Root(mainBranch)
	lineage.Set(perennialBranch, mainBranch)
	lineage.Set(featureBranch, perennialBranch)
	var connector forgedomain.Connector = &mockConnectorProposalStackLineageBuilder{}
	args := configdomain.ProposalStackLineageArgs{
		AfterStackDisplay:        []string{},
		BeforeStackDisplay:       []string{},
		Connector:                Some(connector),
		CurrentBranch:            featureBranch,
		CurrentBranchIndicator:   ":point_left:",
		IndentMarker:             "-",
		Lineage:                  lineage,
		MainAndPerennialBranches: Some(gitdomain.NewLocalBranchNames(mainBranch.String(), perennialBranch.String())),
	}
	expectedStackLineage := ` - main
   - development
     - PR https://www.github.com/git-town/pull/1 :point_left:
`

	// act
	actual := configdomain.NewProposalStackLineageBuilder(&args)

	// assert
	must.True(t, actual.IsSome())
	builder := actual.GetOrDefault()
	must.True(t, builder.GetProposal(mainBranch).IsNone())
	must.True(t, builder.GetProposal(perennialBranch).IsNone())
	must.True(t, builder.GetProposal(featureBranch).IsSome())

	stackLineageAsString := builder.Build(&args)
	must.EqOp(t, expectedStackLineage, stackLineageAsString.GetOrPanic())
}

func TestProposalStackLineageBuilder_ForgeConnectorNone(t *testing.T) {
	t.Parallel()
	// arrange
	args := configdomain.ProposalStackLineageArgs{
		AfterStackDisplay:        []string{},
		BeforeStackDisplay:       []string{},
		Connector:                None[forgedomain.Connector](),
		CurrentBranch:            gitdomain.LocalBranchName("main"),
		CurrentBranchIndicator:   ":point_left:",
		IndentMarker:             "-",
		Lineage:                  configdomain.NewLineage(),
		MainAndPerennialBranches: None[gitdomain.LocalBranchNames](),
	}
	expected := None[configdomain.ProposalStackLineageBuilder]()
	// act
	actual := configdomain.NewProposalStackLineageBuilder(&args)
	// assert
	must.EqOp(t, expected, actual)
}

func TestProposalStackLineageBuilder_NilArgs(t *testing.T) {
	t.Parallel()
	// arrange
	var args *configdomain.ProposalStackLineageArgs
	expected := None[configdomain.ProposalStackLineageBuilder]()
	// act
	actual := configdomain.NewProposalStackLineageBuilder(args)
	// assert
	must.EqOp(t, expected, actual)
}

func TestProposalStackLineageBuilder_NoLineageForMainOrPerennialLineage(t *testing.T) {
	t.Parallel()
	// arrange
	mainBranch := gitdomain.NewLocalBranchName("main")
	perennialBranch := gitdomain.NewLocalBranchName("development")
	lineage := configdomain.NewLineage()
	lineage.Root(mainBranch)
	lineage.Set(perennialBranch, mainBranch)
	var connector forgedomain.Connector = &mockConnectorProposalStackLineageBuilder{}
	args := configdomain.ProposalStackLineageArgs{
		AfterStackDisplay:        []string{},
		BeforeStackDisplay:       []string{},
		Connector:                Some(connector),
		CurrentBranch:            mainBranch,
		CurrentBranchIndicator:   ":point_left:",
		IndentMarker:             "-",
		Lineage:                  lineage,
		MainAndPerennialBranches: Some(gitdomain.NewLocalBranchNames(mainBranch.String(), perennialBranch.String())),
	}
	expected := None[configdomain.ProposalStackLineageBuilder]()

	// act
	actual := configdomain.NewProposalStackLineageBuilder(&args)

	// assert
	must.EqOp(t, expected, actual)
}
