package undobranches

import (
	"github.com/git-town/git-town/v14/internal/config/configdomain"
	"github.com/git-town/git-town/v14/internal/undo/undodomain"
)

func CategorizeInconsistentChanges(changes undodomain.InconsistentChanges, config configdomain.ValidatedConfig) (perennials, features undodomain.InconsistentChanges) {
	for _, change := range changes {
		if localBefore, hasLocalBefore := change.Before.LocalName.Get(); hasLocalBefore {
			if config.IsMainOrPerennialBranch(localBefore) {
				perennials = append(perennials, change)
			} else {
				features = append(features, change)
			}
		}
	}
	return
}

func CategorizeLocalBranchChange(change LocalBranchChange, config configdomain.ValidatedConfig) (changedPerennials, changedFeatures LocalBranchChange) {
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

func CategorizeRemoteBranchChange(change RemoteBranchChange, config configdomain.ValidatedConfig) (perennialChanges, featureChanges RemoteBranchChange) {
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

func CategorizeRemoteBranchesSHAs(shas RemoteBranchesSHAs, config configdomain.ValidatedConfig) (perennials, features RemoteBranchesSHAs) {
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
