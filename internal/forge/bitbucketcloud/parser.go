package bitbucketcloud

import (
	"errors"

	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/messages"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

func parsePullRequest(pullRequest map[string]any) (result forgedomain.BitbucketCloudProposalData, err error) {
	id1, has := pullRequest["id"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	id2, ok := id1.(float64)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	number := int(id2)
	title1, has := pullRequest["title"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	title2, ok := title1.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	body1, has := pullRequest["description"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	body2, ok := body1.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	state1, has := pullRequest["state"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	state2, ok := state1.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination1, has := pullRequest["destination"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination2, ok := destination1.(map[string]any)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination3, has := destination2["branch"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination4, ok := destination3.(map[string]any)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination5, has := destination4["name"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination6, ok := destination5.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source1, has := pullRequest["source"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source2, ok := source1.(map[string]any)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source3, has := source2["branch"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source4, ok := source3.(map[string]any)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source5, has := source4["name"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source6, ok := source5.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url1, has := pullRequest["links"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url2, ok := url1.(map[string]any)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url3, has := url2["html"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url4, ok := url3.(map[string]any)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url5, has := url4["href"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url6, ok := url5.(string)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	closeSourceBranch1, has := pullRequest["close_source_branch"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	closeSourceBranch2, ok := closeSourceBranch1.(bool)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	draft1, has := pullRequest["draft"]
	if !has {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	draft2, ok := draft1.(bool)
	if !ok {
		return result, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	return forgedomain.BitbucketCloudProposalData{
		ProposalData: forgedomain.ProposalData{
			Active:       state2 == "open",
			MergeWithAPI: false,
			Number:       number,
			Source:       gitdomain.NewLocalBranchName(source6),
			Target:       gitdomain.NewLocalBranchName(destination6),
			Title:        title2,
			Body:         NewOption(body2),
			URL:          url6,
		},
		CloseSourceBranch: closeSourceBranch2,
		Draft:             draft2,
	}, nil
}
