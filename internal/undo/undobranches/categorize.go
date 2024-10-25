package undobranches

import (
	"github.com/git-town/git-town/v16/internal/config"
	"github.com/git-town/git-town/v16/internal/undo/undodomain"
)

func CategorizeInconsistentChanges(changes undodomain.InconsistentChanges, config config.ValidatedConfig) (perennials, features undodomain.InconsistentChanges) {
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

func CategorizeLocalBranchChange(change LocalBranchChange, config config.ValidatedConfig) (changedPerennials, changedFeatures LocalBranchChange) {
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

func CategorizeRemoteBranchChange(change RemoteBranchChange, config config.ValidatedConfig) (perennialChanges, featureChanges RemoteBranchChange) {
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

func CategorizeRemoteBranchesSHAs(shas RemoteBranchesSHAs, config config.ValidatedConfig) (perennials, features RemoteBranchesSHAs) {
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
