package mockproposals

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/cucumber/godog"
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks"
	"github.com/git-town/git-town/v22/internal/test/datatable"
	"github.com/git-town/git-town/v22/internal/test/helpers"
	"github.com/git-town/git-town/v22/pkg/asserts"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func FromGherkinTable(table *godog.Table, lineage configdomain.Lineage) []forgedomain.ProposalData {
	result := make([]forgedomain.ProposalData, 0, len(table.Rows)-1)
	headers := helpers.TableFields(table)
	for r := 1; r < len(table.Rows); r++ {
		row := table.Rows[r]
		id := Some(r)
		source := None[gitdomain.LocalBranchName]()
		target := None[gitdomain.LocalBranchName]()
		title := None[gitdomain.ProposalTitle]()
		body := None[gitdomain.ProposalBody]()
		url := None[string]()
		for f, field := range row.Cells {
			switch headers[f] {
			case "ID":
				id = Some(asserts.NoError1(strconv.Atoi(field.Value)))
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
		result = append(result, forgedomain.ProposalData{
			Active:       true,
			Body:         body,
			MergeWithAPI: true,
			Number:       forgedomain.ProposalNumber(id.GetOrPanic()),
			Source:       source.GetOrPanic(),
			Target:       target.GetOrPanic(),
			Title:        title.GetOrZero(),
			URL:          url.GetOrZero(),
		})
	}
	return result
}

func ToDataTable(proposals []forgedomain.ProposalData, fields []string) datatable.DataTable {
	result := datatable.DataTable{}
	result.AddRow(fields...)
	for _, proposal := range proposals {
		row := make([]string, len(fields))
		for f, field := range fields {
			switch field {
			case "ID":
				row[f] = proposal.Number.String()
			case "SOURCE BRANCH":
				row[f] = proposal.Source.String()
			case "TARGET BRANCH":
				row[f] = proposal.Target.String()
			case "TITLE":
				row[f] = proposal.Title.String()
			case "BODY":
				row[f] = proposal.Body.GetOrZero().String()
			case "URL":
				row[f] = proposal.URL
			default:
				panic("unknown field: " + field)
			}
		}
		result.AddRow(row...)
	}
	return result
}

func ToDocString(proposals []forgedomain.ProposalData) string {
	result := strings.Builder{}
	for _, proposal := range proposals {
		result.WriteString("url:")
		if proposal.URL != "" {
			result.WriteString(" ")
			result.WriteString(proposal.URL)
		}
		result.WriteString("\nnumber: ")
		result.WriteString(proposal.Number.String())
		result.WriteString("\nsource: ")
		result.WriteString(proposal.Source.String())
		result.WriteString("\ntarget: ")
		result.WriteString(proposal.Target.String())
		result.WriteString("\nbody:\n")
		if body, hasBody := proposal.Body.Get(); hasBody {
			result.WriteString(gohacks.IndentLines(body.String(), 2))
		}
		result.WriteString("\n")
		result.WriteString("\n")
	}
	return strings.TrimSpace(result.String())
}
