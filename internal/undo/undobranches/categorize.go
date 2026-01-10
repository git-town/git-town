package undobranches

import (
	"github.com/git-town/git-town/v22/internal/config"
	"github.com/git-town/git-town/v22/internal/undo/undodomain"
)

func CategorizeInconsistentChanges(changes undodomain.InconsistentChanges, config config.ValidatedConfig) Categorized[undodomain.InconsistentChanges] {
	perennials := undodomain.InconsistentChanges{}
	features := undodomain.InconsistentChanges{}
	for _, change := range changes {
		if localBefore, hasLocalBefore := change.Before.LocalName.Get(); hasLocalBefore {
			if config.IsMainOrPerennialBranch(localBefore) {
				perennials = append(perennials, change)
			} else {
				features = append(features, change)
			}
		}
	}
	return Categorized[undodomain.InconsistentChanges]{
		Features:   features,
		Perennials: perennials,
	}
}

func CategorizeLocalBranchChange(change LocalBranchChange, config config.ValidatedConfig) Categorized[LocalBranchChange] {
	changedPerennials := LocalBranchChange{}
	changedFeatures := LocalBranchChange{}
	for branch, change := range change { // okay to iterate the map in random order because we assign to new maps
		if config.IsMainOrPerennialBranch(branch) {
			changedPerennials[branch] = change
		} else {
			changedFeatures[branch] = change
		}
	}
	return Categorized[LocalBranchChange]{
		Features:   changedFeatures,
		Perennials: changedPerennials,
	}
}

func CategorizeRemoteBranchChange(change RemoteBranchChange, config config.ValidatedConfig) Categorized[RemoteBranchChange] {
	perennialChanges := RemoteBranchChange{}
	featureChanges := RemoteBranchChange{}
	for branch, change := range change { // okay to iterate the map in random order because we assign to new maps
		if config.IsMainOrPerennialBranch(branch.LocalBranchName()) {
			perennialChanges[branch] = change
		} else {
			featureChanges[branch] = change
		}
	}
	return Categorized[RemoteBranchChange]{
		Features:   featureChanges,
		Perennials: perennialChanges,
	}
}

func CategorizeRemoteBranchesSHAs(shas RemoteBranchesSHAs, config config.ValidatedConfig) Categorized[RemoteBranchesSHAs] {
	perennials := RemoteBranchesSHAs{}
	features := RemoteBranchesSHAs{}
	for branch, sha := range shas { // okay to iterate the map in random order because we assign to new maps
		if config.IsMainOrPerennialBranch(branch.LocalBranchName()) {
			perennials[branch] = sha
		} else {
			features[branch] = sha
		}
	}
	return Categorized[RemoteBranchesSHAs]{
		Features:   features,
		Perennials: perennials,
	}
}

type Categorized[T any] struct {
	Features   T
	Perennials T
}
