package config

import (
	"fmt"
	"slices"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/envconfig"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/git/giturl"
	"github.com/git-town/git-town/v21/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v21/internal/messages"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

type NormalConfig struct {
	configdomain.NormalConfigData
	DryRun     configdomain.DryRun                // whether to only print the Git commands but not execute them
	Env        configdomain.PartialConfig         // configuration data taken from environment variables
	File       Option[configdomain.PartialConfig] // content of git-town.toml, nil = no config file exists
	Git        configdomain.PartialConfig         // configuration data taken from Git metadata, in particular the unscoped Git metadata
	GitVersion git.Version                        // version of the installed Git executable
}

// removes the given branch from the lineage, and updates its children
func (self *NormalConfig) CleanupBranchFromLineage(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) {
	parent, hasParent := self.Git.Lineage.Parent(branch).Get()
	children := self.Lineage.Children(branch)
	for _, child := range children {
		if hasParent {
			self.Lineage = self.Lineage.Set(child, parent)
			_ = gitconfig.SetParent(runner, child, parent)
		} else {
			self.Lineage = self.Lineage.RemoveBranch(child)
			_ = gitconfig.RemoveParent(runner, parent)
		}
	}
	self.Lineage = self.Lineage.RemoveBranch(branch)
	_ = gitconfig.RemoveParent(runner, branch)
}

// DevURL provides the URL for the development remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *NormalConfig) DevURL(querier subshelldomain.Querier) Option[giturl.Parts] {
	return self.RemoteURL(querier, self.DevRemote)
}

// RemoteURL provides the URL for the given remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
// Caches its result so can be called repeatedly.
func (self *NormalConfig) RemoteURL(querier subshelldomain.Querier, remote gitdomain.Remote) Option[giturl.Parts] {
	urlStr, hasURLStr := self.RemoteURLString(querier, remote).Get()
	if !hasURLStr {
		return None[giturl.Parts]()
	}
	url, hasURL := giturl.Parse(urlStr).Get()
	if !hasURL {
		return None[giturl.Parts]()
	}
	if hostnameOverride, hasHostNameOverride := self.HostingOriginHostname.Get(); hasHostNameOverride {
		url.Host = hostnameOverride.String()
	}
	return Some(url)
}

// RemoteURLString provides the URL for the given remote.
// Tests can stub this through the GIT_TOWN_REMOTE environment variable.
func (self *NormalConfig) RemoteURLString(querier subshelldomain.Querier, remote gitdomain.Remote) Option[string] {
	remoteOverride := envconfig.RemoteURLOverride()
	if remoteOverride.IsSome() {
		return remoteOverride
	}
	return gitconfig.RemoteURL(querier, remote)
}

// RemoveParent removes the parent branch entry for the given branch from the Git configuration.
func (self *NormalConfig) RemoveParent(runner subshelldomain.Runner, branch gitdomain.LocalBranchName) {
	self.Git.Lineage = self.Git.Lineage.RemoveBranch(branch)
	_ = gitconfig.RemoveParent(runner, branch)
}

func (self *NormalConfig) RemovePerennialAncestors(runner subshelldomain.Runner, finalMessages stringslice.Collector) {
	for _, perennialBranch := range self.PerennialBranches {
		if self.Lineage.Parent(perennialBranch).IsSome() {
			_ = gitconfig.RemoveParent(runner, perennialBranch)
			self.Lineage = self.Lineage.RemoveBranch(perennialBranch)
			finalMessages.Add(fmt.Sprintf(messages.PerennialBranchRemovedParentEntry, perennialBranch))
		}
	}
}

// SetBranchTypeOverride registers the given branch names as contribution branches.
// The branches must exist.
func (self *NormalConfig) SetBranchTypeOverride(runner subshelldomain.Runner, branchType configdomain.BranchType, branches ...gitdomain.LocalBranchName) error {
	for _, branch := range branches {
		self.BranchTypeOverrides[branch] = branchType
		if err := gitconfig.SetBranchTypeOverride(runner, branch, branchType); err != nil {
			return err
		}
	}
	return nil
}

// SetParent marks the given branch as the direct parent of the other given branch
// in the Git Town configuration.
func (self *NormalConfig) SetParent(runner subshelldomain.Runner, branch, parentBranch gitdomain.LocalBranchName) error {
	if self.DryRun {
		return nil
	}
	self.Lineage = self.Lineage.Set(branch, parentBranch)
	return gitconfig.SetParent(runner, branch, parentBranch)
}

// SetPerennialBranches marks the given branches as perennial branches.
// TODO: inline into setup.go:savePerennialBranches
func (self *NormalConfig) SetPerennialBranches(runner subshelldomain.Runner, branches gitdomain.LocalBranchNames) error {
	self.PerennialBranches = branches
	if slices.Compare(self.Git.PerennialBranches, branches) == 0 {
		return nil
	}
	return gitconfig.SetPerennialBranches(runner, branches)
}

// SetPerennialRegex updates the locally configured perennial regex.
func (self *NormalConfig) SetPerennialRegex(runner subshelldomain.Runner, value configdomain.PerennialRegex) error {
	self.PerennialRegex = Some(value)
	existing, has := self.Git.PerennialRegex.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetPerennialRegex(runner, value)
}

// SetPushHook updates the locally configured push-hook strategy.
func (self *NormalConfig) SetPushHook(runner subshelldomain.Runner, value configdomain.PushHook) error {
	self.PushHook = value
	existing, has := self.Git.PushHook.Get()
	if has && existing == value {
		return nil
	}
	return gitconfig.SetPushHook(runner, value)
}
