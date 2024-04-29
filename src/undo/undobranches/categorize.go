package undobranches

import (
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/undo/undodomain"
)

func CategorizeInconsistentChanges(changes undodomain.InconsistentChanges, config configdomain.FullConfig) (perennials, features undodomain.InconsistentChanges) {
	for _, change := range changes {
		if config.IsMainOrPerennialBranch(change.Before.LocalName) {
			perennials = append(perennials, change)
		} else {
			features = append(features, change)
		}
	}
	return
}

func CategorizeLocalBranchChange(change LocalBranchChange, config configdomain.FullConfig) (changedPerennials, changedFeatures LocalBranchChange) {
	changedPerennials = LocalBranchChange{}
	changedFeatures = LocalBranchChange{}
	for branch, change := range change {
		if config.IsMainOrPerennialBranch(branch) {
			changedPerennials[branch] = change
		} else {
			changedFeatures[branch] = change
		}
	}
	return
}

func CategorizeRemoteBranchChange(change RemoteBranchChange, config configdomain.FullConfig) (perennialChanges, featureChanges RemoteBranchChange) {
	perennialChanges = RemoteBranchChange{}
	featureChanges = RemoteBranchChange{}
	for branch, change := range change {
		if config.IsMainOrPerennialBranch(branch.LocalBranchName()) {
			perennialChanges[branch] = change
		} else {
			featureChanges[branch] = change
		}
	}
	return
}

func CategorizeRemoteBranchesSHAs(shas RemoteBranchesSHAs, config configdomain.FullConfig) (perennials, features RemoteBranchesSHAs) {
	perennials = RemoteBranchesSHAs{}
	features = RemoteBranchesSHAs{}
	for branch, sha := range shas {
		if config.IsMainOrPerennialBranch(branch.LocalBranchName()) {
			perennials[branch] = sha
		} else {
			features[branch] = sha
		}
	}
	return
}
