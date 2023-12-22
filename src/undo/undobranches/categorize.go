package undobranches

import (
	"github.com/git-town/git-town/v11/src/sync/syncdomain"
	"github.com/git-town/git-town/v11/src/undo/undodomain"
)

func CategorizeInconsistentChanges(changes undodomain.InconsistentChanges, branchTypes syncdomain.BranchTypes) (perennials, features undodomain.InconsistentChanges) {
	for _, change := range changes {
		if branchTypes.IsFeatureBranch(change.Before.LocalName) {
			features = append(features, change)
		} else {
			perennials = append(perennials, change)
		}
	}
	return
}

func CategorizeLocalBranchChange(change undodomain.LocalBranchChange, branchTypes syncdomain.BranchTypes) (changedPerennials, changedFeatures undodomain.LocalBranchChange) {
	changedPerennials = undodomain.LocalBranchChange{}
	changedFeatures = undodomain.LocalBranchChange{}
	for branch, change := range change {
		if branchTypes.IsFeatureBranch(branch) {
			changedFeatures[branch] = change
		} else {
			changedPerennials[branch] = change
		}
	}
	return
}

func CategorizeRemoteBranchChange(change undodomain.RemoteBranchChange, branchTypes syncdomain.BranchTypes) (perennialChanges, featureChanges undodomain.RemoteBranchChange) {
	perennialChanges = undodomain.RemoteBranchChange{}
	featureChanges = undodomain.RemoteBranchChange{}
	for branch, change := range change {
		if branchTypes.IsFeatureBranch(branch.LocalBranchName()) {
			featureChanges[branch] = change
		} else {
			perennialChanges[branch] = change
		}
	}
	return
}

func CategorizeRemoteBranchesSHAs(shas undodomain.RemoteBranchesSHAs, branchTypes syncdomain.BranchTypes) (perennials, features undodomain.RemoteBranchesSHAs) {
	perennials = undodomain.RemoteBranchesSHAs{}
	features = undodomain.RemoteBranchesSHAs{}
	for branch, sha := range shas {
		if branchTypes.IsFeatureBranch(branch.LocalBranchName()) {
			features[branch] = sha
		} else {
			perennials[branch] = sha
		}
	}
	return
}
