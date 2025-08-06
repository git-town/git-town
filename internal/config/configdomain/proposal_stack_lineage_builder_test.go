package configdomain_test

import (
	"fmt"
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
		var prNumber int
		for _, char := range branch {
			prNumber += int(char)
		}
		return Some(forgedomain.Proposal{
			Data: forgedomain.ProposalData{
				Body:         None[string](),
				MergeWithAPI: false,
				Number:       1,
				Source:       branch,
				Target:       target,
				Title:        "Test Mocker",
				URL:          fmt.Sprintf("https://www.github.com/git-town/pull/%d", prNumber),
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
	featureBranchA := gitdomain.NewLocalBranchName("a")
	featureBranchB := gitdomain.NewLocalBranchName("b")
	lineage := configdomain.NewLineageWith(configdomain.LineageData{
		featureBranchA: mainBranch,
		featureBranchB: featureBranchA,
	})
	var connector forgedomain.Connector = &mockConnectorProposalStackLineageBuilder{}
	args := configdomain.ProposalStackLineageArgs{
		AfterStackDisplay:        []string{},
		BeforeStackDisplay:       []string{},
		Connector:                Some(connector),
		CurrentBranch:            featureBranchA,
		CurrentBranchIndicator:   ":point_left:",
		IndentMarker:             "-",
		Lineage:                  lineage,
		MainAndPerennialBranches: gitdomain.NewLocalBranchNames(mainBranch.String()),
	}
	expectedStackLineage := ` - main
   - PR https://www.github.com/git-town/pull/97 :point_left:
     - PR https://www.github.com/git-town/pull/98
`

	// act
	actual := configdomain.NewProposalStackLineageBuilder(&args)

	// assert
	builder, hasBuilder := actual.Get()
	must.True(t, hasBuilder)
	must.True(t, builder.GetProposal(mainBranch).IsNone())
	must.True(t, builder.GetProposal(featureBranchB).IsSome())
	must.True(t, builder.GetProposal(featureBranchA).IsSome())

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
		MainAndPerennialBranches: gitdomain.NewLocalBranchNames("main"),
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

func TestProposalStackLineageBuilder_NoLineageForMainAndPerennialBranches(t *testing.T) {
	t.Parallel()
	// arrange
	mainBranch := gitdomain.NewLocalBranchName("main")
	featureBranchA := gitdomain.NewLocalBranchName("a")
	lineage := configdomain.NewLineageWith(configdomain.LineageData{
		featureBranchA: mainBranch,
	})
	var connector forgedomain.Connector = &mockConnectorProposalStackLineageBuilder{}
	args := configdomain.ProposalStackLineageArgs{
		AfterStackDisplay:        []string{},
		BeforeStackDisplay:       []string{},
		Connector:                Some(connector),
		CurrentBranch:            mainBranch,
		CurrentBranchIndicator:   ":point_left:",
		IndentMarker:             "-",
		Lineage:                  lineage,
		MainAndPerennialBranches: lineage.Roots(),
	}
	expected := None[configdomain.ProposalStackLineageBuilder]()

	// act
	actual := configdomain.NewProposalStackLineageBuilder(&args)

	// assert
	must.EqOp(t, expected, actual)
}
