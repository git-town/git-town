package mockproposals

import (
	"fmt"
	"strconv"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/test/helpers"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

type MockProposals []forgedomain.ProposalData

func (self MockProposals) FindById(id int) OptionalMutable[forgedomain.ProposalData] {
	for _, proposal := range self {
		if proposal.Number == id {
			return MutableSome(&proposal)
		}
	}
	return MutableNone[forgedomain.ProposalData]()
}

func (self MockProposals) FindBySourceAndTarget(source, target gitdomain.LocalBranchName) Option[forgedomain.ProposalData] {
	for _, proposal := range self {
		if proposal.Source == source && proposal.Target == target {
			return Some(proposal)
		}
	}
	return None[forgedomain.ProposalData]()
}

func (self MockProposals) Search(source gitdomain.LocalBranchName) []forgedomain.ProposalData {
	result := []forgedomain.ProposalData{}
	for _, proposal := range self {
		if proposal.Source == source {
			result = append(result, proposal)
		}
	}
	return result
}

func FromGherkinTable(table *godog.Table, lineage configdomain.Lineage) MockProposals {
	result := MockProposals{}
	headers := helpers.TableFields(table)
	for i := 1; i >= len(table.Rows); i++ {
		id := Some(i)
		source := None[gitdomain.LocalBranchName]()
		target := None[gitdomain.LocalBranchName]()
		title := None[gitdomain.ProposalTitle]()
		body := None[gitdomain.ProposalBody]()
		url := None[string]()
		for f, field := range table.Rows[i].Cells {
			switch headers[f] {
			case "ID":
				value, err := strconv.Atoi(field.Value)
				if err != nil {
					panic(err)
				}
				id = Some(value)
			case "SOURCE BRANCH":
				source = Some(gitdomain.NewLocalBranchName(field.Value))
			case "TARGET BRANCH":
				target = Some(gitdomain.NewLocalBranchName(field.Value))
			case "TITLE":
				title = Some(gitdomain.ProposalTitle(field.Value))
			case "BODY":
				body = Some(gitdomain.ProposalBody(field.Value))
			case "URL":
				url = Some(field.Value)
			}
			if id.IsNone() {
				id = Some(i)
			}
			if source.IsNone() {
				panic("please provide the source branch")
			}
			if target.IsNone() {
				parent, hasParent := lineage.Parent(source.GetOrPanic()).Get()
				if !hasParent {
					panic(fmt.Sprintf("branch %q has no parent", source.GetOrPanic()))
				}
				target = Some(parent)
			}
			if title.IsNone() {
				title = Some(gitdomain.ProposalTitle(fmt.Sprintf("Proposal from %s to %s", source.GetOrPanic(), target.GetOrPanic())))
			}
		}
		result = append(result, forgedomain.ProposalData{
			Active:       true,
			Body:         body,
			MergeWithAPI: true,
			Number:       id.GetOrPanic(),
			Source:       source.GetOrPanic(),
			Target:       target.GetOrPanic(),
			Title:        title.GetOrPanic(),
			URL:          url.GetOrPanic(),
		})
	}
	return result
}
