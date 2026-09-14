package bitbucketcloud

import (
	"errors"
	"strings"

	"github.com/git-town/git-town/v24/internal/forge/forgedomain"
	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/git-town/git-town/v24/internal/gohacks/stringss"
	"github.com/git-town/git-town/v24/internal/messages"
)

func parsePullRequest(pullRequest map[string]any) (forgedomain.BitbucketCloudProposalData, error) {
	var emptyResult forgedomain.BitbucketCloudProposalData
	id1, has := pullRequest["id"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	id2, ok := id1.(float64)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	number := forgedomain.NewProposalNumberFromFloat64(id2)
	title1, has := pullRequest["title"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	title2, ok := title1.(string)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	body1, has := pullRequest["description"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	body2, ok := body1.(string)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	state1, has := pullRequest["state"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	state2, ok := state1.(string)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	state3 := strings.ToLower(state2)
	isActive := state3 == "open" || state3 == "new"
	destination1, has := pullRequest["destination"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination2, ok := destination1.(map[string]any)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination3, has := destination2["branch"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination4, ok := destination3.(map[string]any)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination5, has := destination4["name"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	destination6, ok := destination5.(string)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source1, has := pullRequest["source"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source2, ok := source1.(map[string]any)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source3, has := source2["branch"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source4, ok := source3.(map[string]any)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source5, has := source4["name"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	source6, ok := source5.(string)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url1, has := pullRequest["links"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url2, ok := url1.(map[string]any)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url3, has := url2["html"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url4, ok := url3.(map[string]any)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url5, has := url4["href"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	url6, ok := url5.(string)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	closeSourceBranch1, has := pullRequest["close_source_branch"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	closeSourceBranch2, ok := closeSourceBranch1.(bool)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	draft1, has := pullRequest["draft"]
	if !has {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	draft2, ok := draft1.(bool)
	if !ok {
		return emptyResult, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	reviewers, reviewerAccountIDs, err := parseReviewers(pullRequest)
	if err != nil {
		return emptyResult, err
	}
	return forgedomain.BitbucketCloudProposalData{
		ProposalData: forgedomain.ProposalData{
			Active:       isActive,
			MergeWithAPI: false,
			Number:       number,
			Source:       gitdomain.LocalBranchNameOrPanic(stringss.Trim(source6)),
			Target:       gitdomain.LocalBranchNameOrPanic(stringss.Trim(destination6)),
			Title:        gitdomain.ProposalTitle(title2),
			Body:         gitdomain.NewProposalBodyOpt(body2),
			URL:          url6,
		},
		CloseSourceBranch:  closeSourceBranch2,
		Draft:              draft2,
		Reviewers:          reviewers,
		ReviewerAccountIDs: reviewerAccountIDs,
	}, nil
}

func parseReviewers(pullRequest map[string]any) ([]string, []string, error) {
	reviewers1, has := pullRequest["reviewers"]
	if !has {
		return []string{}, []string{}, nil
	}
	reviewers2, ok := reviewers1.([]any)
	if !ok {
		return nil, nil, errors.New(messages.APIUnexpectedResultDataStructure)
	}
	reviewerUUIDs := make([]string, 0, len(reviewers2))
	reviewerAccountIDs := make([]string, 0, len(reviewers2))
	for _, reviewer1 := range reviewers2 {
		reviewer2, ok := reviewer1.(map[string]any)
		if !ok {
			return nil, nil, errors.New(messages.APIUnexpectedResultDataStructure)
		}
		uuid1, has := reviewer2["uuid"]
		if !has {
			return nil, nil, errors.New(messages.APIUnexpectedResultDataStructure)
		}
		uuid2, ok := uuid1.(string)
		if !ok {
			return nil, nil, errors.New(messages.APIUnexpectedResultDataStructure)
		}
		reviewerUUIDs = append(reviewerUUIDs, uuid2)
		accountID1, has := reviewer2["account_id"]
		if has {
			accountID2, ok := accountID1.(string)
			if !ok {
				return nil, nil, errors.New(messages.APIUnexpectedResultDataStructure)
			}
			reviewerAccountIDs = append(reviewerAccountIDs, accountID2)
		}
	}
	return reviewerUUIDs, reviewerAccountIDs, nil
}
