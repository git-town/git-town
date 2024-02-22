package undobranches

import (
	"github.com/git-town/git-town/v12/src/config/configdomain"
	"github.com/git-town/git-town/v12/src/undo/undodomain"
)

func CategorizeInconsistentChanges(changes undodomain.InconsistentChanges, config *configdomain.FullConfig) (perennials, features undodomain.InconsistentChanges) {
	for _, change := range changes {
		switch config.BranchType(change.Before.LocalName) {
		case configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
			perennials = append(perennials, change)
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeParkedBranch:
			features = append(features, change)
		}
	}
	return
}

func CategorizeLocalBranchChange(change LocalBranchChange, config *configdomain.FullConfig) (changedPerennials, changedFeatures LocalBranchChange) {
	for branch, change := range change {
		switch config.BranchType(branch) {
		case configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch:
			changedPerennials[branch] = change
		case configdomain.BranchTypeFeatureBranch, configdomain.BranchTypeObservedBranch, configdomain.BranchTypeParkedBranch:
			changedFeatures[branch] = change
		}
	}
	return
}

func CategorizeRemoteBranchChange(change RemoteBranchChange, config *configdomain.FullConfig) (perennialChanges, featureChanges RemoteBranchChange) {
	perennialChanges = RemoteBranchChange{}
	featureChanges = RemoteBranchChange{}
	for branch, change := range change {
		if config.IsFeatureBranch(branch.LocalBranchName()) {
			featureChanges[branch] = change
		} else {
			perennialChanges[branch] = change
		}
	}
	return
}

func CategorizeRemoteBranchesSHAs(shas RemoteBranchesSHAs, config *configdomain.FullConfig) (perennials, features RemoteBranchesSHAs) {
	perennials = RemoteBranchesSHAs{}
	features = RemoteBranchesSHAs{}
	for branch, sha := range shas {
		if config.IsFeatureBranch(branch.LocalBranchName()) {
			features[branch] = sha
		} else {
			perennials[branch] = sha
		}
	}
	return
}
