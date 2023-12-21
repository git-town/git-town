package undo

import (
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

func CategorizeInconsistentChanges(changes domain.InconsistentChanges, branchTypes domain.BranchTypes) (perennials, features domain.InconsistentChanges) {
	for _, change := range changes {
		if branchTypes.IsFeatureBranch(change.Before.LocalName) {
			features = append(features, change)
		} else {
			perennials = append(perennials, change)
		}
	}
	return
}

func CategorizeLocalBranchChange(change domain.LocalBranchChange, branchTypes domain.BranchTypes) (changedPerennials, changedFeatures domain.LocalBranchChange) {
	changedPerennials = domain.LocalBranchChange{}
	changedFeatures = domain.LocalBranchChange{}
	for branch, change := range change {
		if branchTypes.IsFeatureBranch(branch) {
			changedFeatures[branch] = change
		} else {
			changedPerennials[branch] = change
		}
	}
	return
}

func CategorizeLocalBranchNames(branchNames gitdomain.LocalBranchNames, branchTypes domain.BranchTypes) (perennials, features gitdomain.LocalBranchNames) {
	for _, branch := range branchNames {
		if branchTypes.IsFeatureBranch(branch) {
			features = append(features, branch)
		} else {
			perennials = append(perennials, branch)
		}
	}
	return
}

func CategorizeRemoteBranchChange(change domain.RemoteBranchChange, branchTypes domain.BranchTypes) (perennialChanges, featureChanges domain.RemoteBranchChange) {
	perennialChanges = domain.RemoteBranchChange{}
	featureChanges = domain.RemoteBranchChange{}
	for branch, change := range change {
		if branchTypes.IsFeatureBranch(branch.LocalBranchName()) {
			featureChanges[branch] = change
		} else {
			perennialChanges[branch] = change
		}
	}
	return
}

func CategorizeRemoteBranchesSHAs(shas domain.RemoteBranchesSHAs, branchTypes domain.BranchTypes) (perennials, features domain.RemoteBranchesSHAs) {
	perennials = domain.RemoteBranchesSHAs{}
	features = domain.RemoteBranchesSHAs{}
	for branch, sha := range shas {
		if branchTypes.IsFeatureBranch(branch.LocalBranchName()) {
			features[branch] = sha
		} else {
			perennials[branch] = sha
		}
	}
	return
}
