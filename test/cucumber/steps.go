package cucumber

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/v13/src/cli/dialog/components"
	"github.com/git-town/git-town/v13/src/cli/print"
	"github.com/git-town/git-town/v13/src/config/configdomain"
	"github.com/git-town/git-town/v13/src/config/configfile"
	"github.com/git-town/git-town/v13/src/config/gitconfig"
	"github.com/git-town/git-town/v13/src/git/gitdomain"
	"github.com/git-town/git-town/v13/src/gohacks"
	"github.com/git-town/git-town/v13/test/asserts"
	"github.com/git-town/git-town/v13/test/commands"
	"github.com/git-town/git-town/v13/test/datatable"
	"github.com/git-town/git-town/v13/test/fixture"
	"github.com/git-town/git-town/v13/test/git"
	"github.com/git-town/git-town/v13/test/helpers"
	"github.com/git-town/git-town/v13/test/output"
	"github.com/git-town/git-town/v13/test/subshell"
	"github.com/google/go-cmp/cmp"
)

// beforeSuiteMux ensures that we run BeforeSuite only once globally.
var beforeSuiteMux sync.Mutex //nolint:gochecknoglobals

// the global FixtureFactory instance.
var fixtureFactory *fixture.Factory //nolint:gochecknoglobals

// Steps defines Cucumber step implementations around Git workspace management.
func Steps(suite *godog.Suite, state *ScenarioState) {
	suite.BeforeScenario(func(scenario *messages.Pickle) {
		// create a Fixture for the scenario
		fixture := fixtureFactory.CreateFixture(scenario.GetName())
		// Godog only provides state for the entire feature.
		// We want state to be scenario-specific, hence we reset the shared state before each scenario.
		// This is a limitation of the current Godog implementation, which doesn't have a `ScenarioContext` method,
		// only a `FeatureContext` method.
		// See main_test.go for additional details.
		state.Reset(fixture)
		if helpers.HasTag(scenario, "@debug") {
			state.fixture.DevRepo.Verbose = true
		}
	})

	suite.BeforeSuite(func() {
		// NOTE: we want to create only one global FixtureFactory instance with one global memoized environment.
		beforeSuiteMux.Lock()
		defer beforeSuiteMux.Unlock()
		if fixtureFactory == nil {
			baseDir, err := os.MkdirTemp("", "")
			if err != nil {
				log.Fatalf("cannot create base directory for feature specs: %s", err)
			}
			// Evaluate symlinks as Mac temp dir is symlinked
			evalBaseDir, err := filepath.EvalSymlinks(baseDir)
			if err != nil {
				log.Fatalf("cannot evaluate symlinks of base directory for feature specs: %s", err)
			}
			gm := fixture.NewFactory(evalBaseDir)
			fixtureFactory = &gm
		}
	})

	suite.AfterScenario(func(scenario *messages.Pickle, e error) {
		if e != nil {
			fmt.Printf("failed scenario %q in %s - investigate state in %s\n", scenario.GetName(), scenario.GetUri(), state.fixture.Dir)
		}
		if state.runExitCode != 0 && !state.runExitCodeChecked {
			print.Error(fmt.Errorf("%s - scenario %q doesn't document exit code %d", scenario.GetUri(), scenario.GetName(), state.runExitCode))
			os.Exit(1)
		}
	})

	suite.Step(`^a branch "([^"]*)"$`, func(branch string) error {
		state.fixture.DevRepo.CreateBranch(gitdomain.NewLocalBranchName(branch), gitdomain.NewLocalBranchName("main"))
		return nil
	})

	suite.Step(`^a coworker clones the repository$`, func() error {
		state.fixture.AddCoworkerRepo()
		return nil
	})

	suite.Step(`^a feature branch "([^"]+)" as a child of "([^"]+)"$`, func(branchText, parentBranch string) error {
		branch := gitdomain.NewLocalBranchName(branchText)
		state.fixture.DevRepo.CreateChildFeatureBranch(branch, gitdomain.NewLocalBranchName(parentBranch))
		state.fixture.DevRepo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
		return nil
	})

	suite.Step(`^a folder "([^"]*)"$`, func(name string) error {
		state.fixture.DevRepo.CreateFolder(name)
		return nil
	})

	suite.Step(`^a known remote feature branch "([^"]*)"$`, func(branchText string) error {
		branch := gitdomain.NewLocalBranchName(branchText)
		// we are creating a remote branch in the remote repo --> it is a local branch there
		state.fixture.OriginRepo.CreateBranch(branch, gitdomain.NewLocalBranchName("main"))
		state.fixture.DevRepo.TestCommands.Fetch()
		return nil
	})

	suite.Step(`^a merge is now in progress$`, func() error {
		if !state.fixture.DevRepo.HasMergeInProgress() {
			return errors.New("expected merge in progress")
		}
		return nil
	})

	suite.Step(`^a (local )?feature branch "([^"]*)"$`, func(localStr, branchText string) error {
		branch := gitdomain.NewLocalBranchName(branchText)
		isLocal := localStr != ""
		state.fixture.DevRepo.CreateFeatureBranch(branch)
		if !isLocal {
			state.fixture.DevRepo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
			return nil
		}
		return nil
	})

	suite.Step(`^a parked branch "([^"]+)"$`, func(branchText string) error {
		branch := gitdomain.NewLocalBranchName(branchText)
		state.fixture.DevRepo.CreateParkedBranches(branch)
		state.fixture.DevRepo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
		return nil
	})

	suite.Step(`^a perennial branch "([^"]+)"$`, func(branchText string) error {
		branch := gitdomain.NewLocalBranchName(branchText)
		state.fixture.DevRepo.CreatePerennialBranches(branch)
		state.fixture.DevRepo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
		return nil
	})

	suite.Step(`^a rebase is now in progress$`, func() error {
		repoStatus, err := state.fixture.DevRepo.RepoStatus()
		asserts.NoError(err)
		if !repoStatus.RebaseInProgress {
			return errors.New("expected rebase in progress")
		}
		return nil
	})

	suite.Step(`^a remote feature branch "([^"]*)"$`, func(branchText string) error {
		branch := gitdomain.NewLocalBranchName(branchText)
		// we are creating a remote branch in the remote repo --> it is a local branch there
		state.fixture.OriginRepo.CreateBranch(branch, gitdomain.NewLocalBranchName("main"))
		return nil
	})

	suite.Step(`^a remote tag "([^"]+)" not on a branch$`, func(name string) error {
		state.fixture.OriginRepo.CreateStandaloneTag(name)
		return nil
	})

	suite.Step(`^all branches are now synchronized$`, func() error {
		branchesOutOfSync, output := state.fixture.DevRepo.HasBranchesOutOfSync()
		if branchesOutOfSync {
			return errors.New("unexpected out of sync:\n" + output)
		}
		return nil
	})

	suite.Step(`^an uncommitted file$`, func() error {
		state.uncommittedFileName = "uncommitted file"
		state.uncommittedContent = "uncommitted content"
		state.fixture.DevRepo.CreateFile(
			state.uncommittedFileName,
			state.uncommittedContent,
		)
		return nil
	})

	suite.Step(`^an uncommitted file in folder "([^"]*)"$`, func(folder string) error {
		state.uncommittedFileName = folder + "/uncommitted file"
		state.fixture.DevRepo.CreateFile(
			state.uncommittedFileName,
			state.uncommittedContent,
		)
		return nil
	})

	suite.Step(`^an uncommitted file with name "([^"]+)" and content "([^"]+)"$`, func(name, content string) error {
		state.uncommittedFileName = name
		state.uncommittedContent = content
		state.fixture.DevRepo.CreateFile(name, content)
		return nil
	})

	suite.Step(`^an upstream repo$`, func() error {
		state.fixture.AddUpstream()
		return nil
	})

	suite.Step(`^branch "([^"]+)" is active in another worktree`, func(branch string) error {
		state.fixture.AddSecondWorktree(gitdomain.NewLocalBranchName(branch))
		return nil
	})

	suite.Step(`^branch "([^"]+)" is (?:now|still) a contribution branch`, func(name string) error {
		branch := gitdomain.NewLocalBranchName(name)
		if !state.fixture.DevRepo.Config.FullConfig.IsContributionBranch(branch) {
			return fmt.Errorf(
				"branch %q isn't contribution as expected.\nContribution branches: %s",
				branch,
				strings.Join(state.fixture.DevRepo.Config.FullConfig.ContributionBranches.Strings(), ", "),
			)
		}
		return nil
	})

	suite.Step(`^branch "([^"]+)" is (?:now|still) a feature branch`, func(name string) error {
		branch := gitdomain.NewLocalBranchName(name)
		if state.fixture.DevRepo.Config.FullConfig.BranchType(branch) != configdomain.BranchTypeFeatureBranch {
			return fmt.Errorf("branch %q isn't a feature branch as expected", branch)
		}
		return nil
	})

	suite.Step(`^branch "([^"]+)" is (?:now|still) observed`, func(name string) error {
		branch := gitdomain.NewLocalBranchName(name)
		if !state.fixture.DevRepo.Config.FullConfig.IsObservedBranch(branch) {
			return fmt.Errorf(
				"branch %q isn't observed as expected.\nObserved branches: %s",
				branch,
				strings.Join(state.fixture.DevRepo.Config.FullConfig.ObservedBranches.Strings(), ", "),
			)
		}
		return nil
	})

	suite.Step(`^branch "([^"]+)" is now parked`, func(name string) error {
		branch := gitdomain.NewLocalBranchName(name)
		if !state.fixture.DevRepo.Config.FullConfig.IsParkedBranch(branch) {
			return fmt.Errorf(
				"branch %q isn't parked as expected.\nParked branches: %s",
				branch,
				strings.Join(state.fixture.DevRepo.Config.FullConfig.ParkedBranches.Strings(), ", "),
			)
		}
		return nil
	})

	suite.Step(`^branch "([^"]+)" is (?:now|still) perennial`, func(name string) error {
		branch := gitdomain.NewLocalBranchName(name)
		if !state.fixture.DevRepo.Config.FullConfig.IsPerennialBranch(branch) {
			return fmt.Errorf(
				"branch %q isn't perennial as expected.\nPerennial branches: %s",
				branch,
				strings.Join(state.fixture.DevRepo.Config.FullConfig.PerennialBranches.Strings(), ", "),
			)
		}
		return nil
	})

	suite.Step(`^branch "([^"]+)" is now a feature branch`, func(name string) error {
		branch := gitdomain.NewLocalBranchName(name)
		if state.fixture.DevRepo.Config.FullConfig.IsParkedBranch(branch) {
			return fmt.Errorf("branch %q is parked", branch)
		}
		if state.fixture.DevRepo.Config.FullConfig.IsObservedBranch(branch) {
			return fmt.Errorf("branch %q is observed", branch)
		}
		if state.fixture.DevRepo.Config.FullConfig.IsContributionBranch(branch) {
			return fmt.Errorf("branch %q is contribution", branch)
		}
		if state.fixture.DevRepo.Config.FullConfig.IsPerennialBranch(branch) {
			return fmt.Errorf("branch %q is perennial", branch)
		}
		return nil
	})

	suite.Step(`^display "([^"]+)"$`, func(command string) error {
		parts := strings.Split(command, " ")
		output, err := state.fixture.DevRepo.Runner.Query(parts[0], parts[1:]...)
		fmt.Println("XXXXXXXXXXXXXXXXX " + strings.ToUpper(command) + " START XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		fmt.Println(output)
		fmt.Println("XXXXXXXXXXXXXXXXX " + strings.ToUpper(command) + " END XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		return err
	})

	suite.Step(`^file "([^"]+)" still contains unresolved conflicts$`, func(name string) error {
		content := state.fixture.DevRepo.FileContent(name)
		if !strings.Contains(content, "<<<<<<<") {
			return fmt.Errorf("file %q does not contain unresolved conflicts", name)
		}
		return nil
	})

	suite.Step(`^file "([^"]*)" (?:now|still) has content "([^"]*)"$`, func(file, expectedContent string) error {
		actualContent := state.fixture.DevRepo.FileContent(file)
		if expectedContent != actualContent {
			return fmt.Errorf("file content does not match\n\nEXPECTED: %q\n\nACTUAL:\n\n%q\n----------------------------", expectedContent, actualContent)
		}
		return nil
	})

	suite.Step(`^Git has version "([^"]*)"$`, func(version string) error {
		state.fixture.DevRepo.MockGit(version)
		return nil
	})

	suite.Step(`^Git Town is no longer configured$`, func() error {
		return state.fixture.DevRepo.VerifyNoGitTownConfiguration()
	})

	suite.Step(`^Git Town is not configured$`, func() error {
		err := state.fixture.DevRepo.RemovePerennialBranchConfiguration()
		if err != nil {
			return err
		}
		state.fixture.DevRepo.RemoveMainBranchConfiguration()
		return nil
	})

	suite.Step(`^Git Town setting "color.ui" is "([^"]*)"$`, func(value string) error {
		return state.fixture.DevRepo.SetColorUI(value)
	})

	suite.Step(`^Git Town parent setting for branch "([^"]*)" is "([^"]*)"$`, func(branch, value string) error {
		branchName := gitdomain.NewLocalBranchName(branch)
		configKey := gitconfig.NewParentKey(branchName)
		return state.fixture.DevRepo.Config.GitConfig.SetLocalConfigValue(configKey, value)
	})

	suite.Step(`^local Git setting "init.defaultbranch" is "([^"]*)"$`, func(value string) error {
		state.fixture.DevRepo.SetDefaultGitBranch(gitdomain.NewLocalBranchName(value))
		return nil
	})

	suite.Step(`^global Git setting "alias\.(.*?)" is "([^"]*)"$`, func(name, value string) error {
		key := gitconfig.ParseKey("alias." + name)
		if key == nil {
			return fmt.Errorf("no key found for %q", name)
		}
		aliasableCommand := gitconfig.AliasableCommandForKey(*key)
		if aliasableCommand == nil {
			return fmt.Errorf("no aliasableCommand found for key %q", *key)
		}
		return state.fixture.DevRepo.SetGitAlias(*aliasableCommand, value)
	})

	suite.Step(`^global Git setting "alias\.(.*?)" (?:now|still) doesn't exist$`, func(name string) error {
		key := gitconfig.ParseKey("alias." + name)
		if key == nil {
			return errors.New("key not found")
		}
		aliasableCommand := gitconfig.AliasableCommandForKey(*key)
		command, has := state.fixture.DevRepo.Config.FullConfig.Aliases[*aliasableCommand]
		if !has {
			return nil
		}
		return fmt.Errorf("unexpected aliasableCommand %q: %q", *key, command)
	})

	suite.Step(`^global Git setting "alias\.(.*?)" is (?:now|still) "([^"]*)"$`, func(name, want string) error {
		key := gitconfig.ParseKey("alias." + name)
		if key == nil {
			return errors.New("key not found")
		}
		aliasableCommand := gitconfig.AliasableCommandForKey(*key)
		if aliasableCommand == nil {
			return fmt.Errorf("aliasableCommand not found for key %q", *key)
		}
		have := state.fixture.DevRepo.Config.FullConfig.Aliases[*aliasableCommand]
		if have != want {
			return fmt.Errorf("unexpected value for key %q: want %q have %q", name, want, have)
		}
		return nil
	})

	suite.Step(`^global Git Town setting "([^"]*)" is "([^"]*)"$`, func(name, value string) error {
		configKey := gitconfig.ParseKey("git-town." + name)
		if configKey == nil {
			return fmt.Errorf("unknown configuration key: %q", name)
		}
		return state.fixture.DevRepo.Config.GitConfig.SetGlobalConfigValue(*configKey, value)
	})

	suite.Step(`^global Git Town setting "([^"]*)" (?:now|still) doesn't exist$`, func(name string) error {
		configKey := gitconfig.ParseKey("git-town." + name)
		newValue := state.fixture.DevRepo.TestCommands.GlobalGitConfig(*configKey)
		if newValue != nil {
			return fmt.Errorf("should not have global %q anymore but has value %q", name, *newValue)
		}
		return nil
	})

	suite.Step(`^global Git Town setting "hosting-origin-hostname" is now "([^"]*)"$`, func(want string) error {
		have := state.fixture.DevRepo.Config.GlobalGitConfig.HostingOriginHostname
		if have.String() != want {
			return fmt.Errorf(`expected global setting "hosting-origin-hostname" to be %q, but was %q`, want, *have)
		}
		return nil
	})

	suite.Step(`^global Git Town setting "hosting-platform" is now "([^"]*)"$`, func(want string) error {
		have := state.fixture.DevRepo.Config.GlobalGitConfig.HostingPlatform
		if have.String() != want {
			return fmt.Errorf(`expected global setting "hosting-platform" to be %q, but was %q`, want, *have)
		}
		return nil
	})

	suite.Step(`^global Git Town setting "main-branch" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.GlobalGitConfig.MainBranch
		want := gitdomain.LocalBranchName(wantStr)
		if *have != want {
			return fmt.Errorf(`expected global setting "main-branch" to be %q, but was %q`, want, *have)
		}
		return nil
	})

	suite.Step(`^global Git Town setting "offline" is (?:now|still) "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.GlobalGitConfig.Offline
		wantBool, err := gohacks.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.Offline(wantBool)
		if *have != want {
			return fmt.Errorf(`expected global setting "offline" to be %t, but was %t`, want, *have)
		}
		return nil
	})

	suite.Step(`^global Git Town setting "perennial-branches" is (?:now|still) "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.GlobalGitConfig.PerennialBranches
		want := gitdomain.NewLocalBranchNames(strings.Split(wantStr, " ")...)
		if cmp.Equal(*have, want) {
			return nil
		}
		return fmt.Errorf(`expected global setting "perennial-branches" to be %v, but was %v`, want, *have)
	})

	suite.Step(`^global Git Town setting "push-hook" is (?:now|still) "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.GlobalGitConfig.PushHook
		wantBool, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.PushHook(wantBool)
		if cmp.Equal(*have, want) {
			return nil
		}
		return fmt.Errorf(`expected global setting "push-hook" to be %v, but was %v`, want, *have)
	})

	suite.Step(`^global Git Town setting "push-new-branches" is (?:now|still) "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.GlobalGitConfig.PushNewBranches
		wantBool, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.PushNewBranches(wantBool)
		if cmp.Equal(*have, want) {
			return nil
		}
		return fmt.Errorf(`expected global setting "push-new-branches" to be %v, but was %v`, want, *have)
	})

	suite.Step(`^global Git Town setting "ship-delete-tracking-branch" is (?:now|still) "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.GlobalGitConfig.ShipDeleteTrackingBranch
		wantBool, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.ShipDeleteTrackingBranch(wantBool)
		if cmp.Equal(*have, want) {
			return nil
		}
		return fmt.Errorf(`expected global setting "ship-delete-tracking-branch" to be %v, but was %v`, want, *have)
	})

	suite.Step(`^global Git Town setting "sync-before-ship" is (?:now|still) "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.GlobalGitConfig.SyncBeforeShip
		wantBool, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.SyncBeforeShip(wantBool)
		if cmp.Equal(*have, want) {
			return nil
		}
		return fmt.Errorf(`expected global setting "sync-before-ship" to be %v, but was %v`, want, *have)
	})

	suite.Step(`^global Git Town setting "sync-feature-strategy" is (?:now|still) "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.GlobalGitConfig.SyncFeatureStrategy
		want, err := configdomain.NewSyncFeatureStrategy(wantStr)
		asserts.NoError(err)
		if cmp.Equal(*have, want) {
			return nil
		}
		return fmt.Errorf(`expected global setting "sync-feature-strategy" to be %v, but was %v`, want, *have)
	})

	suite.Step(`^global Git Town setting "sync-perennial-strategy" is (?:now|still) "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.GlobalGitConfig.SyncPerennialStrategy
		want, err := configdomain.NewSyncPerennialStrategy(wantStr)
		asserts.NoError(err)
		if cmp.Equal(*have, want) {
			return nil
		}
		return fmt.Errorf(`expected global setting "sync-perennial-strategy" to be %v, but was %v`, want, *have)
	})

	suite.Step(`^global Git Town setting "sync-upstream" is (?:now|still) "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.GlobalGitConfig.SyncUpstream
		wantBool, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.SyncUpstream(wantBool)
		if cmp.Equal(*have, want) {
			return nil
		}
		return fmt.Errorf(`expected global setting "sync-upstream" to be %v, but was %v`, want, *have)
	})

	suite.Step(`^I add commit "([^"]*)" to the "([^"]*)" branch`, func(message, branch string) error {
		state.fixture.DevRepo.CreateCommit(git.Commit{
			Branch:   gitdomain.NewLocalBranchName(branch),
			FileName: "new_file",
			Message:  message,
		})
		return nil
	})

	suite.Step(`^I add (?:this commit|these commits) to the current branch:$`, func(table *messages.PickleStepArgument_PickleTable) error {
		for _, commit := range git.FromGherkinTable(table) {
			state.fixture.DevRepo.CreateFile(commit.FileName, commit.FileContent)
			state.fixture.DevRepo.StageFiles(commit.FileName)
			state.fixture.DevRepo.CommitStagedChanges(commit.Message)
		}
		return nil
	})

	suite.Step(`^I am not prompted for any parent branches$`, func() error {
		notExpected := "Please specify the parent branch of"
		if strings.Contains(state.runOutput, notExpected) {
			return fmt.Errorf("text found:\n\nDID NOT EXPECT: %q\n\nACTUAL\n\n%q\n----------------------------", notExpected, state.runOutput)
		}
		return nil
	})

	suite.Step(`^I am outside a Git repo$`, func() error {
		state.insideGitRepo = false
		os.RemoveAll(filepath.Join(state.fixture.DevRepo.WorkingDir, ".git"))
		return nil
	})

	suite.Step(`^I resolve the conflict in "([^"]*)"(?: with "([^"]*)")?$`, func(filename, content string) error {
		if content == "" {
			content = "resolved content"
		}
		state.fixture.DevRepo.CreateFile(filename, content)
		state.fixture.DevRepo.StageFiles(filename)
		return nil
	})

	suite.Step(`^I resolve the conflict in "([^"]*)" in the other worktree$`, func(filename string) error {
		content := "resolved content"
		state.fixture.SecondWorktree.CreateFile(filename, content)
		state.fixture.SecondWorktree.StageFiles(filename)
		return nil
	})

	suite.Step(`^I (?:run|ran) "(.+)"$`, func(command string) error {
		state.CaptureState()
		updateInitialSHAs(state)
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCode(command)
		state.fixture.DevRepo.Config.Reload()
		return nil
	})

	suite.Step(`^I run "([^"]*)" and close the editor$`, func(cmd string) error {
		state.CaptureState()
		updateInitialSHAs(state)
		env := append(os.Environ(), "GIT_EDITOR=true")
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Env: env})
		state.fixture.DevRepo.Config.Reload()
		return nil
	})

	suite.Step(`^I run "([^"]*)" and enter an empty commit message$`, func(cmd string) error {
		state.CaptureState()
		updateInitialSHAs(state)
		state.fixture.DevRepo.MockCommitMessage("")
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCode(cmd)
		state.fixture.DevRepo.Config.Reload()
		return nil
	})

	suite.Step(`^I run "([^"]*)" and enter "([^"]*)" for the commit message$`, func(cmd, message string) error {
		state.CaptureState()
		updateInitialSHAs(state)
		state.fixture.DevRepo.MockCommitMessage(message)
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCode(cmd)
		state.fixture.DevRepo.Config.Reload()
		return nil
	})

	suite.Step(`^I run "([^"]*)" in the other worktree and enter "([^"]*)" for the commit message$`, func(cmd, message string) error {
		state.CaptureState()
		updateInitialSHAs(state)
		state.fixture.SecondWorktree.MockCommitMessage(message)
		state.runOutput, state.runExitCode = state.fixture.SecondWorktree.MustQueryStringCode(cmd)
		state.fixture.SecondWorktree.Config.Reload()
		return nil
	})

	suite.Step(`^I (?:run|ran) "([^"]+)" and enter into the dialogs?:$`, func(cmd string, input *messages.PickleStepArgument_PickleTable) error {
		state.CaptureState()
		updateInitialSHAs(state)
		env := os.Environ()
		answers, err := helpers.TableToInputEnv(input)
		if err != nil {
			return err
		}
		for dialogNumber, answer := range answers {
			env = append(env, fmt.Sprintf("%s_%02d=%s", components.TestInputKey, dialogNumber, answer))
		}
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Env: env})
		state.fixture.DevRepo.Config.Reload()
		return nil
	})

	suite.Step(`^I run "([^"]*)", enter into the dialog, and close the next editor:$`, func(cmd string, input *messages.PickleStepArgument_PickleTable) error {
		state.CaptureState()
		updateInitialSHAs(state)
		env := append(os.Environ(), "GIT_EDITOR=true")
		answers, err := helpers.TableToInputEnv(input)
		if err != nil {
			return err
		}
		for dialogNumber, answer := range answers {
			env = append(env, fmt.Sprintf("%s%d=%s", components.TestInputKey, dialogNumber, answer))
		}
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Env: env})
		state.fixture.DevRepo.Config.Reload()
		return nil
	})

	suite.Step(`^I run "([^"]+)" in the "([^"]+)" folder$`, func(cmd, folderName string) error {
		state.CaptureState()
		updateInitialSHAs(state)
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Dir: folderName})
		state.fixture.DevRepo.Config.Reload()
		return nil
	})

	suite.Step(`^I run "([^"]+)" in the other worktree$`, func(cmd string) error {
		state.CaptureState()
		updateInitialSHAs(state)
		state.runOutput, state.runExitCode = state.fixture.SecondWorktree.MustQueryStringCode(cmd)
		state.fixture.SecondWorktree.Config.Reload()
		return nil
	})

	suite.Step(`^inspect the repo$`, func() error {
		fmt.Printf("\nThe workspace is at %s\n", state.fixture.DevRepo.WorkingDir)
		time.Sleep(1 * time.Hour)
		return nil
	})

	suite.Step(`^it does not print "(.+)"$`, func(text string) error {
		if strings.Contains(stripansi.Strip(state.runOutput), text) {
			return fmt.Errorf("text found: %q", text)
		}
		return nil
	})

	suite.Step(`^it prints:$`, func(expected *messages.PickleStepArgument_PickleDocString) error {
		if state.runExitCode != 0 {
			return fmt.Errorf("unexpected exit code %d", state.runExitCode)
		}
		if !strings.Contains(stripansi.Strip(state.runOutput), expected.Content) {
			fmt.Println("ERROR: text not found:")
			fmt.Println("\nEXPECTED:", expected.Content)
			fmt.Println()
			fmt.Println("==================================================================")
			fmt.Println("ACTUAL OUTPUT START ==============================================")
			fmt.Println("==================================================================")
			fmt.Println()
			fmt.Println(state.runOutput)
			fmt.Println()
			fmt.Println("==================================================================")
			fmt.Println("ACTUAL OUTPUT END ================================================")
			fmt.Println("==================================================================")
			fmt.Println()
			return errors.New("expected text not found")
		}
		return nil
	})

	suite.Step(`^it prints no output$`, func() error {
		output := state.runOutput
		if output != "" {
			return fmt.Errorf("expected no output but found %q", output)
		}
		return nil
	})

	suite.Step(`^it prints something like:$`, func(expected *messages.PickleStepArgument_PickleDocString) error {
		regex := regexp.MustCompile(expected.Content)
		have := stripansi.Strip(state.runOutput)
		if !regex.MatchString(have) {
			return fmt.Errorf("EXPECTED: content matching %q\nGOT: %q", expected.Content, have)
		}
		return nil
	})

	suite.Step(`^it prints the error:$`, func(expected *messages.PickleStepArgument_PickleDocString) error {
		state.runExitCodeChecked = true
		if !strings.Contains(stripansi.Strip(state.runOutput), expected.Content) {
			return fmt.Errorf("text not found:\n%s\n\nactual text:\n%s", expected.Content, state.runOutput)
		}
		if state.runExitCode == 0 {
			return fmt.Errorf("unexpected exit code %d", state.runExitCode)
		}
		return nil
	})

	suite.Step(`^it runs no commands$`, func() error {
		commands := output.GitCommandsInGitTownOutput(state.runOutput)
		if len(commands) > 0 {
			fmt.Println("\n\nERROR: Unexpected commands run!")
			for _, command := range commands {
				fmt.Printf("%s > %s\n", command.Branch, command.Command)
			}
			fmt.Println()
			fmt.Println()
			return fmt.Errorf("expected no commands but found %d commands", len(commands))
		}
		return nil
	})

	suite.Step(`^it runs the commands$`, func(input *messages.PickleStepArgument_PickleTable) error {
		commands := output.GitCommandsInGitTownOutput(state.runOutput)
		table := output.RenderExecutedGitCommands(commands, input)
		dataTable := datatable.FromGherkin(input)
		expanded := dataTable.Expand(
			&state.fixture.DevRepo,
			state.fixture.OriginRepo,
			state.fixture.SecondWorktree,
			state.initialDevSHAs,
			state.initialOriginSHAs,
			state.initialWorktreeSHAs,
		)
		diff, errorCount := table.EqualDataTable(expanded)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the commands run\n\n", errorCount)
			fmt.Println(diff)
			return errors.New("mismatching commands run, see diff above")
		}
		return nil
	})

	suite.Step(`^it runs without error$`, func() error {
		if state.runExitCode != 0 {
			return fmt.Errorf("did not expect the Git Town command to produce an exit code: %d", state.runExitCode)
		}
		return nil
	})

	suite.Step(`^"([^"]*)" launches a new proposal with this url in my browser:$`, func(tool string, url *messages.PickleStepArgument_PickleDocString) error {
		want := fmt.Sprintf("%s called with: %s", tool, url.Content)
		want = strings.ReplaceAll(want, "?", `\?`)
		regex := regexp.MustCompile(want)
		have := state.runOutput
		if !regex.MatchString(have) {
			return fmt.Errorf("EXPECTED: a regex matching %q\nGOT: %q", want, have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "([^"]*)" (:?now|still) doesn't exist$`, func(name string) error {
		configKey := gitconfig.ParseKey("git-town." + name)
		newValue := state.fixture.DevRepo.TestCommands.LocalGitConfig(*configKey)
		if newValue != nil {
			return fmt.Errorf("should not have local %q anymore but has value %q", name, *newValue)
		}
		return nil
	})

	suite.Step(`^(?:local )?Git Town setting "([^"]*)" doesn't exist$`, func(name string) error {
		configKey := gitconfig.ParseKey("git-town." + name)
		return state.fixture.DevRepo.Config.GitConfig.RemoveLocalConfigValue(*configKey)
	})

	suite.Step(`^(?:local )?Git Town setting "([^"]*)" is "([^"]*)"$`, func(name, value string) error {
		configKey := gitconfig.ParseKey("git-town." + name)
		if configKey == nil {
			return fmt.Errorf("unknown config key: %q", name)
		}
		return state.fixture.DevRepo.Config.GitConfig.SetLocalConfigValue(*configKey, value)
	})

	suite.Step(`^local Git Town setting "code-hosting-origin-hostname" now doesn't exist$`, func() error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.HostingOriginHostname
		if have == nil {
			return nil
		}
		return fmt.Errorf(`unexpected local setting "code-hosting-origin-hostname" with value %q`, *have)
	})

	suite.Step(`^local Git Town setting "hosting-platform" is now "([^"]*)"$`, func(want string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.HostingPlatform
		if have.String() != want {
			return fmt.Errorf(`expected local setting "hosting-platform" to be %q, but was %q`, want, *have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "hosting-platform" (:?now|still) doesn't exist$`, func() error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.HostingPlatform
		if have != nil {
			return fmt.Errorf(`expected local setting "hosting-platform" to not exist but was %q`, *have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "gitea-token" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.GiteaToken
		want := configdomain.GiteaToken(wantStr)
		if *have != want {
			return fmt.Errorf(`expected local setting "gitea-token" to be %q, but was %q`, want, have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "github-token" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.GitHubToken
		want := configdomain.GitHubToken(wantStr)
		if *have != want {
			return fmt.Errorf(`expected local setting "github-token" to be %q, but was %q`, want, have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "github-token" now doesn't exist$`, func() error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.GitHubToken
		if have == nil {
			return nil
		}
		return fmt.Errorf(`unexpected local setting "github-token" with value %q`, have)
	})

	suite.Step(`^local Git Town setting "gitlab-token" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.GitLabToken
		want := configdomain.GitLabToken(wantStr)
		if *have != want {
			return fmt.Errorf(`expected local setting "gitlab-token" to be %q, but was %q`, want, have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "hosting-origin-hostname" is now "([^"]*)"$`, func(want string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.HostingOriginHostname
		if have.String() != want {
			return fmt.Errorf(`expected local setting "hosting-origin-hostname" to be %q, but was %q`, want, *have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "hosting-origin-hostname" now doesn't exist$`, func() error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.HostingOriginHostname
		if have == nil {
			return nil
		}
		return fmt.Errorf(`unexpected local setting "hosting-origin-hostname" with value %q`, *have)
	})

	suite.Step(`^local Git Town setting "main-branch" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.MainBranch
		want := gitdomain.NewLocalBranchName(wantStr)
		if *have != want {
			return fmt.Errorf(`expected local setting "main-branch" to be %q, but was %q`, want, have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "perennial-branches" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.PerennialBranches
		want := gitdomain.NewLocalBranchNames(strings.Split(wantStr, " ")...)
		if cmp.Equal(*have, want) {
			return nil
		}
		return fmt.Errorf(`expected local setting "main-branch" to be %v, but was %v`, want, have)
	})

	suite.Step(`^local Git Town setting "perennial-regex" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.PerennialRegex
		want := configdomain.PerennialRegex(wantStr)
		if *have != want {
			return fmt.Errorf(`expected local setting "perennial-regex" to be %q, but was %q`, want, have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "push-hook" is (:?now|still) not set$`, func() error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.PushHook
		if have == nil {
			return nil
		}
		return fmt.Errorf(`unexpected local setting "push-hook" %v`, have)
	})

	suite.Step(`^local Git Town setting "push-hook" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.PushHook
		wantBool, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.PushHook(wantBool)
		if cmp.Equal(*have, want) {
			return nil
		}
		return fmt.Errorf(`expected local setting "push-hook" to be %v, but was %v`, want, have)
	})

	suite.Step(`^local Git Town setting "push-new-branches" is (:?now|still) not set$`, func() error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.PushNewBranches
		if have == nil {
			return nil
		}
		return fmt.Errorf(`unexpected local setting "push-new-branches" %v`, have)
	})

	suite.Step(`^local Git Town setting "push-new-branches" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.PushNewBranches
		wantBool, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.PushNewBranches(wantBool)
		if cmp.Equal(*have, want) {
			return nil
		}
		return fmt.Errorf(`expected local setting "push-new-branches" to be %v, but was %v`, want, have)
	})

	suite.Step(`^local Git Town setting "ship-delete-tracking-branch" is still not set$`, func() error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.ShipDeleteTrackingBranch
		if have == nil {
			return nil
		}
		return fmt.Errorf(`unexpected local setting "ship-delete-tracking-branch" %v`, have)
	})

	suite.Step(`^local Git Town setting "ship-delete-tracking-branch" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.ShipDeleteTrackingBranch
		wantBool, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.ShipDeleteTrackingBranch(wantBool)
		if *have != want {
			return fmt.Errorf(`expected local setting "ship-delete-tracking-branch" to be %v, but was %v`, want, have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "sync-before-ship" is still not set$`, func() error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.SyncBeforeShip
		if have == nil {
			return nil
		}
		return fmt.Errorf(`unexpected local setting "sync-before-ship" %v`, have)
	})

	suite.Step(`^local Git Town setting "sync-before-ship" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.SyncBeforeShip
		wantBool, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.SyncBeforeShip(wantBool)
		if *have != want {
			return fmt.Errorf(`expected local setting "sync-before-ship" to be %v, but was %v`, want, have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "sync-feature-strategy" is still not set$`, func() error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.SyncFeatureStrategy
		if have == nil {
			return nil
		}
		return fmt.Errorf(`expected local setting "sync-feature-strategy" %v`, have)
	})

	suite.Step(`^local Git Town setting "sync-feature-strategy" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.SyncFeatureStrategy
		want, err := configdomain.NewSyncFeatureStrategy(wantStr)
		asserts.NoError(err)
		if *have != want {
			return fmt.Errorf(`expected local setting "sync-feature-strategy" to be %v, but was %v`, want, have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "sync-perennial-strategy" is still not set$`, func() error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.SyncPerennialStrategy
		if have == nil {
			return nil
		}
		return fmt.Errorf(`unexpected local setting "sync-perennial-strategy" %v`, have)
	})

	suite.Step(`^local Git Town setting "sync-perennial-strategy" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.SyncPerennialStrategy
		want, err := configdomain.NewSyncPerennialStrategy(wantStr)
		asserts.NoError(err)
		if *have != want {
			return fmt.Errorf(`expected local setting "sync-perennial-strategy" to be %v, but was %v`, want, have)
		}
		return nil
	})

	suite.Step(`^local Git Town setting "sync-upstream" is still not set$`, func() error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.SyncUpstream
		if have == nil {
			return nil
		}
		return fmt.Errorf(`unexpected local setting "sync-upstream" %v`, have)
	})

	suite.Step(`^local Git Town setting "sync-upstream" is now "([^"]*)"$`, func(wantStr string) error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.SyncUpstream
		wantBool, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.SyncUpstream(wantBool)
		if *have != want {
			return fmt.Errorf(`expected local setting "sync-upstream" to be %v, but was %v`, want, have)
		}
		return nil
	})

	suite.Step(`^my repo does not have an origin$`, func() error {
		state.fixture.DevRepo.RemoveRemote(gitdomain.RemoteOrigin)
		state.fixture.OriginRepo = nil
		return nil
	})

	suite.Step(`^my repo has a Git submodule$`, func() error {
		state.fixture.AddSubmoduleRepo()
		state.fixture.DevRepo.AddSubmodule(state.fixture.SubmoduleRepo.WorkingDir)
		return nil
	})

	suite.Step(`^my repo's "([^"]*)" remote is "([^"]*)"$`, func(remoteName, remoteURL string) error {
		remote := gitdomain.Remote(remoteName)
		state.fixture.DevRepo.RemoveRemote(remote)
		state.fixture.DevRepo.AddRemote(remote, remoteURL)
		return nil
	})

	suite.Step(`^still no configuration file exists$`, func() error {
		_, err := state.fixture.DevRepo.FileContentErr(configfile.FileName)
		if err == nil {
			return errors.New("expected no configuration file but found one")
		}
		return nil
	})

	suite.Step(`^no commits exist now$`, func() error {
		currentCommits := state.fixture.CommitTable(state.initialCommits.Cells[0])
		noCommits := datatable.DataTable{}
		noCommits.AddRow(state.initialCommits.Cells[0]...)
		errDiff, errCount := currentCommits.EqualDataTable(noCommits)
		if errCount == 0 {
			return nil
		}
		fmt.Println(errDiff)
		return errors.New("found unexpected commits")
	})

	suite.Step(`^no lineage exists now$`, func() error {
		if state.fixture.DevRepo.Config.FullConfig.ContainsLineage() {
			lineage := state.fixture.DevRepo.Config.FullConfig.Lineage
			return fmt.Errorf("unexpected Git Town lineage information: %+v", lineage)
		}
		return nil
	})

	suite.Step(`^no merge is in progress$`, func() error {
		if state.fixture.DevRepo.HasMergeInProgress() {
			return errors.New("expected no merge in progress")
		}
		return nil
	})

	suite.Step(`^no rebase is in progress$`, func() error {
		repoStatus, err := state.fixture.DevRepo.RepoStatus()
		if err != nil {
			return err
		}
		if repoStatus.RebaseInProgress {
			return errors.New("expected no rebase in progress")
		}
		return nil
	})

	suite.Step(`^no tool to open browsers is installed$`, func() error {
		state.fixture.DevRepo.MockNoCommandsInstalled()
		return nil
	})

	suite.Step(`^no uncommitted files exist$`, func() error {
		files := state.fixture.DevRepo.UncommittedFiles()
		if len(files) > 0 {
			return fmt.Errorf("unexpected uncommitted files: %s", files)
		}
		return nil
	})

	suite.Step(`^offline mode is disabled$`, func() error {
		isOffline := state.fixture.DevRepo.Config.FullConfig.Offline
		if isOffline {
			return errors.New("expected to not be offline but am")
		}
		return nil
	})

	suite.Step(`^offline mode is enabled$`, func() error {
		return state.fixture.DevRepo.Config.SetOffline(true)
	})

	suite.Step(`^origin deletes the "([^"]*)" branch$`, func(branch string) error {
		state.fixture.OriginRepo.RemoveBranch(gitdomain.NewLocalBranchName(branch))
		return nil
	})

	suite.Step(`^origin ships the "([^"]*)" branch$`, func(branch string) error {
		state.fixture.OriginRepo.CheckoutBranch(gitdomain.NewLocalBranchName("main"))
		err := state.fixture.OriginRepo.MergeBranch(gitdomain.NewLocalBranchName(branch))
		if err != nil {
			return err
		}
		state.fixture.OriginRepo.RemoveBranch(gitdomain.NewLocalBranchName(branch))
		return nil
	})

	suite.Step(`^the branches "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		for _, branchName := range []string{branch1, branch2} {
			branch := gitdomain.NewLocalBranchName(branchName)
			state.fixture.DevRepo.CreateBranch(branch, gitdomain.NewLocalBranchName("main"))
		}
		return nil
	})

	suite.Step(`^the branches are now$`, func(table *messages.PickleStepArgument_PickleTable) error {
		existing := state.fixture.Branches()
		diff, errCount := existing.EqualGherkin(table)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the branches\n\n", errCount)
			fmt.Println(diff)
			return errors.New("mismatching branches found, see the diff above")
		}
		return nil
	})

	suite.Step(`^the commits$`, func(table *messages.PickleStepArgument_PickleTable) error {
		initialTable := datatable.FromGherkin(table)
		state.initialCommits = &initialTable
		// create the commits
		commits := git.FromGherkinTable(table)
		state.fixture.CreateCommits(commits)
		// restore the initial branch
		if state.initialCurrentBranch.IsEmpty() {
			state.fixture.DevRepo.CheckoutBranch(gitdomain.NewLocalBranchName("main"))
			return nil
		}
		// NOTE: reading the cached value here to keep the test suite fast by avoiding unnecessary disk access
		if state.fixture.DevRepo.CurrentBranchCache.Value() != state.initialCurrentBranch {
			state.fixture.DevRepo.CheckoutBranch(state.initialCurrentBranch)
			return nil
		}
		return nil
	})

	suite.Step(`^the committed configuration file:$`, func(content *messages.PickleStepArgument_PickleDocString) error {
		state.fixture.DevRepo.CreateFile(configfile.FileName, content.Content)
		state.fixture.DevRepo.StageFiles(configfile.FileName)
		state.fixture.DevRepo.CommitStagedChanges(commands.ConfigFileCommitMessage)
		state.fixture.DevRepo.PushBranch()
		return nil
	})

	suite.Step(`^the configuration file:$`, func(content *messages.PickleStepArgument_PickleDocString) error {
		state.fixture.DevRepo.CreateFile(configfile.FileName, content.Content)
		return nil
	})

	suite.Step(`^the configuration file is (?:now|still):$`, func(content *messages.PickleStepArgument_PickleDocString) error {
		have, err := state.fixture.DevRepo.FileContentErr(configfile.FileName)
		if err != nil {
			return errors.New("no configuration file found")
		}
		have = strings.TrimSpace(have)
		want := strings.TrimSpace(content.Content)
		if have != want {
			fmt.Println(cmp.Diff(want, have))
			return errors.New("mismatching config file content")
		}
		return nil
	})

	suite.Step(`^a contribution branch "([^"]+)"$`, func(branch string) error {
		state.fixture.DevRepo.CreateBranch(gitdomain.NewLocalBranchName(branch), "main")
		return state.fixture.DevRepo.Config.SetContributionBranches(gitdomain.NewLocalBranchNames(branch))
	})

	suite.Step(`^the contribution branches "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		state.fixture.DevRepo.CreateBranch(gitdomain.NewLocalBranchName(branch1), "main")
		state.fixture.DevRepo.CreateBranch(gitdomain.NewLocalBranchName(branch2), "main")
		return state.fixture.DevRepo.Config.SetContributionBranches(gitdomain.NewLocalBranchNames(branch1, branch2))
	})

	suite.Step(`^the coworker adds this commit to their current branch:$`, func(table *messages.PickleStepArgument_PickleTable) error {
		commits := git.FromGherkinTable(table)
		commit := commits[0]
		state.fixture.CoworkerRepo.CreateFile(commit.FileName, commit.FileContent)
		state.fixture.CoworkerRepo.StageFiles(commit.FileName)
		state.fixture.CoworkerRepo.CommitStagedChanges(commit.Message)
		return nil
	})

	suite.Step(`^the coworker fetches updates$`, func() error {
		state.fixture.CoworkerRepo.Fetch()
		return nil
	})

	suite.Step(`^the coworker is on the "([^"]*)" branch$`, func(branch string) error {
		state.fixture.CoworkerRepo.CheckoutBranch(gitdomain.NewLocalBranchName(branch))
		return nil
	})

	suite.Step(`^the coworker resolves the conflict in "([^"]*)"(?: with "([^"]*)")?$`, func(filename, content string) error {
		state.fixture.CoworkerRepo.CreateFile(filename, content)
		state.fixture.CoworkerRepo.StageFiles(filename)
		return nil
	})

	suite.Step(`^the coworker runs "([^"]+)"$`, func(command string) error {
		state.runOutput, state.runExitCode = state.fixture.CoworkerRepo.MustQueryStringCode(command)
		return nil
	})

	suite.Step(`^the coworker runs "([^"]*)" and closes the editor$`, func(cmd string) error {
		env := append(os.Environ(), "GIT_EDITOR=true")
		state.runOutput, state.runExitCode = state.fixture.CoworkerRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Env: env})
		return nil
	})

	suite.Step(`^the coworker sets the parent branch of "([^"]*)" as "([^"]*)"$`, func(childBranch, parentBranch string) error {
		_ = state.fixture.CoworkerRepo.Config.SetParent(gitdomain.NewLocalBranchName(childBranch), gitdomain.NewLocalBranchName(parentBranch))
		return nil
	})

	suite.Step(`^the coworker sets the "sync-feature-strategy" to "(merge|rebase)"$`, func(value string) error {
		syncFeatureStrategy, err := configdomain.NewSyncFeatureStrategy(value)
		if err != nil {
			return err
		}
		_ = state.fixture.CoworkerRepo.Config.SetSyncFeatureStrategy(syncFeatureStrategy)
		return nil
	})

	suite.Step(`^the coworkers workspace now contains file "([^"]*)" with content "([^"]*)"$`, func(file, expectedContent string) error {
		actualContent := state.fixture.CoworkerRepo.FileContent(file)
		if expectedContent != actualContent {
			return fmt.Errorf("file content does not match\n\nEXPECTED: %q\n\nACTUAL:\n\n%q\n----------------------------", expectedContent, actualContent)
		}
		return nil
	})

	suite.Step(`^the current branch is "([^"]*)"$`, func(name string) error {
		branch := gitdomain.NewLocalBranchName(name)
		state.initialCurrentBranch = branch
		if !state.fixture.DevRepo.BranchExists(branch) {
			state.fixture.DevRepo.CreateBranch(branch, gitdomain.NewLocalBranchName("main"))
		}
		state.fixture.DevRepo.CheckoutBranch(branch)
		return nil
	})

	suite.Step(`^the current branch is an? (local )?(feature|perennial|parked|contribution|observed) branch "([^"]*)"$`, func(localStr, branchType, branchName string) error {
		branch := gitdomain.NewLocalBranchName(branchName)
		isLocal := localStr != ""
		switch branchType {
		case "feature":
			state.fixture.DevRepo.CreateFeatureBranch(branch)
		case "perennial":
			state.fixture.DevRepo.CreatePerennialBranches(branch)
		case "parked":
			state.fixture.DevRepo.CreateParkedBranches(branch)
		case "contribution":
			state.fixture.DevRepo.CreateContributionBranches(branch)
		case "observed":
			state.fixture.DevRepo.CreateObservedBranches(branch)
		default:
			panic(fmt.Sprintf("unknown branch type: %q", branchType))
		}
		if !isLocal {
			state.fixture.DevRepo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
		}
		state.initialCurrentBranch = branch
		// NOTE: reading the cached value here to keep the test suite fast by avoiding unnecessary disk access
		if !state.fixture.DevRepo.CurrentBranchCache.Initialized() || state.fixture.DevRepo.CurrentBranchCache.Value() != branch {
			state.fixture.DevRepo.CheckoutBranch(branch)
		}
		return nil
	})

	suite.Step(`^the current branch is "([^"]*)" and the previous branch is "([^"]*)"$`, func(currentText, previousText string) error {
		current := gitdomain.NewLocalBranchName(currentText)
		previous := gitdomain.NewLocalBranchName(previousText)
		state.initialCurrentBranch = current
		state.fixture.DevRepo.CheckoutBranch(previous)
		state.fixture.DevRepo.CheckoutBranch(current)
		return nil
	})

	suite.Step(`^feature branch "([^"]*)" with these commits$`, func(name string, table *messages.PickleStepArgument_PickleTable) error {
		branch := gitdomain.NewLocalBranchName(name)
		state.fixture.DevRepo.CreateFeatureBranch(branch)
		state.fixture.DevRepo.CheckoutBranch(branch)
		state.fixture.DevRepo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
		for _, commit := range git.FromGherkinTable(table) {
			state.fixture.DevRepo.CreateFile(commit.FileName, commit.FileContent)
			state.fixture.DevRepo.StageFiles(commit.FileName)
			state.fixture.DevRepo.CommitStagedChanges(commit.Message)
			if commit.Locations.Contains(git.Location(gitdomain.RemoteOrigin)) {
				state.fixture.DevRepo.PushBranch()
			}
		}
		return nil
	})

	suite.Step(`^feature branch "([^"]*)" with these commits is a child of "([^"]*)"$`, func(name, parent string, table *messages.PickleStepArgument_PickleTable) error {
		branch := gitdomain.NewLocalBranchName(name)
		parentBranch := gitdomain.NewLocalBranchName(parent)
		state.fixture.DevRepo.CreateChildFeatureBranch(branch, parentBranch)
		state.fixture.DevRepo.CheckoutBranch(branch)
		state.fixture.DevRepo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
		for _, commit := range git.FromGherkinTable(table) {
			state.fixture.DevRepo.CreateFile(commit.FileName, commit.FileContent)
			state.fixture.DevRepo.StageFiles(commit.FileName)
			state.fixture.DevRepo.CommitStagedChanges(commit.Message)
			if commit.Locations.Contains(git.Location(gitdomain.RemoteOrigin)) {
				state.fixture.DevRepo.PushBranch()
			}
		}
		return nil
	})

	suite.Step(`^the current branch is (?:now|still) "([^"]*)"$`, func(expected string) error {
		state.fixture.DevRepo.CurrentBranchCache.Invalidate()
		actual, err := state.fixture.DevRepo.CurrentBranch()
		if err != nil {
			return fmt.Errorf("cannot determine current branch of developer repo: %w", err)
		}
		if actual.String() != expected {
			return fmt.Errorf("expected active branch %q but is %q", expected, actual)
		}
		return nil
	})

	suite.Step(`^the current branch in the other worktree is (?:now|still) "([^"]*)"$`, func(expected string) error {
		state.fixture.SecondWorktree.CurrentBranchCache.Invalidate()
		actual, err := state.fixture.SecondWorktree.CurrentBranch()
		if err != nil {
			return fmt.Errorf("cannot determine current branch of second worktree: %w", err)
		}
		if actual.String() != expected {
			return fmt.Errorf("expected active branch %q but is %q", expected, actual)
		}
		return nil
	})

	suite.Step(`^the initial lineage exists$`, func() error {
		have := state.fixture.DevRepo.LineageTable()
		diff, errCnt := have.EqualDataTable(*state.initialLineage)
		if errCnt > 0 {
			fmt.Printf("\nERROR! Found %d differences in the lineage\n\n", errCnt)
			fmt.Printf("INITIAL LINEAGE:\n%s\n", state.initialLineage.String())
			fmt.Printf("CURRENT LINEAGE:\n%s\n", have.String())
			fmt.Println(diff)
			return errors.New("mismatching branches found, see the diff above")
		}
		return nil
	})

	suite.Step(`^the initial branches and lineage exist$`, func() error {
		// verify initial branches
		currentBranches := state.fixture.Branches()
		initialBranches := state.initialBranches
		// fmt.Printf("\nINITIAL:\n%s\n", initialBranches)
		// fmt.Printf("NOW:\n%s\n", currentBranches.String())
		diff, errorCount := currentBranches.EqualDataTable(*initialBranches)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing branches\n\n", errorCount)
			fmt.Println(diff)
			return errors.New("mismatching branches found, see diff above")
		}
		// verify initial lineage
		currentLineage := state.fixture.DevRepo.LineageTable()
		diff, errCnt := currentLineage.EqualDataTable(*state.initialLineage)
		if errCnt > 0 {
			fmt.Printf("\nERROR! Found %d differences in the lineage\n\n", errCnt)
			fmt.Println(diff)
			return errors.New("mismatching lineage found, see the diff above")
		}
		return nil
	})

	suite.Step(`^the initial branches exist$`, func() error {
		have := state.fixture.Branches()
		want := state.initialBranches
		// fmt.Printf("HAVE:\n%s\n", have.String())
		// fmt.Printf("WANT:\n%s\n", want.String())
		diff, errorCount := have.EqualDataTable(*want)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing branches\n\n", errorCount)
			fmt.Println(diff)
			return errors.New("mismatching branches found, see diff above")
		}
		return nil
	})

	suite.Step(`^the initial commits exist$`, func() error {
		currentCommits := state.fixture.CommitTable(state.initialCommits.Cells[0])
		errDiff, errCount := state.initialCommits.EqualDataTable(currentCommits)
		if errCount == 0 {
			return nil
		}
		fmt.Println(errDiff)
		return errors.New("current commits are not the same as the initial commits")
	})

	suite.Step(`^the local feature branch "([^"]+)"$`, func(branch string) error {
		branchName := gitdomain.NewLocalBranchName(branch)
		state.fixture.DevRepo.CreateFeatureBranch(branchName)
		return nil
	})

	suite.Step(`^the (local )?feature branches "([^"]+)" and "([^"]+)"$`, func(localStr, branch1, branch2 string) error {
		isLocal := localStr != ""
		for _, branchText := range []string{branch1, branch2} {
			branch := gitdomain.NewLocalBranchName(branchText)
			state.fixture.DevRepo.CreateFeatureBranch(branch)
			if !isLocal {
				state.fixture.DevRepo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
			}
		}
		return nil
	})

	suite.Step(`^the (local )?feature branches "([^"]+)", "([^"]+)", and "([^"]+)"$`, func(localStr, branch1, branch2, branch3 string) error {
		isLocal := localStr != ""
		for _, branchText := range []string{branch1, branch2, branch3} {
			branch := gitdomain.NewLocalBranchName(branchText)
			state.fixture.DevRepo.CreateFeatureBranch(branch)
			if !isLocal {
				state.fixture.DevRepo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
			}
		}
		return nil
	})

	suite.Step(`^the local observed branch "([^"]+)"$`, func(name string) error {
		branch := gitdomain.NewLocalBranchName(name)
		state.fixture.DevRepo.CreateObservedBranches(branch)
		return nil
	})

	suite.Step(`^the (local )?perennial branches "([^"]+)" and "([^"]+)"$`, func(localStr, branch1Text, branch2Text string) error {
		branch1 := gitdomain.NewLocalBranchName(branch1Text)
		branch2 := gitdomain.NewLocalBranchName(branch2Text)
		isLocal := localStr != ""
		state.fixture.DevRepo.CreatePerennialBranches(branch1, branch2)
		if !isLocal {
			state.fixture.DevRepo.PushBranchToRemote(branch1, gitdomain.RemoteOrigin)
			state.fixture.DevRepo.PushBranchToRemote(branch2, gitdomain.RemoteOrigin)
		}
		return nil
	})

	suite.Step(`^the (local )?perennial branches "([^"]+)", "([^"]+)", and "([^"]+)"$`, func(localStr, branch1, branch2, branch3 string) error {
		isLocal := localStr != ""
		for _, branchText := range []string{branch1, branch2, branch3} {
			branch := gitdomain.NewLocalBranchName(branchText)
			state.fixture.DevRepo.CreatePerennialBranches(branch)
			if !isLocal {
				state.fixture.DevRepo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
			}
		}
		return nil
	})

	suite.Step(`^the main branch is "([^"]+)"$`, func(name string) error {
		return state.fixture.DevRepo.Config.SetMainBranch(gitdomain.NewLocalBranchName(name))
	})

	suite.Step(`^the main branch is (?:now|still) "([^"]*)"$`, func(want string) error {
		have := state.fixture.DevRepo.Config.FullConfig.MainBranch
		if have.String() != want {
			return fmt.Errorf("expected %q, got %q", want, have)
		}
		return nil
	})

	suite.Step(`^the main branch is (?:now|still) not set$`, func() error {
		have := state.fixture.DevRepo.Config.LocalGitConfig.MainBranch
		if have == nil {
			return nil
		}
		return fmt.Errorf("unexpected main branch setting %q", have)
	})

	suite.Step(`^an observed branch "([^"]+)"$`, func(name string) error {
		branch := gitdomain.NewLocalBranchName(name)
		state.fixture.DevRepo.CreateBranch(branch, "main")
		state.fixture.DevRepo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
		return state.fixture.DevRepo.Config.SetObservedBranches(gitdomain.NewLocalBranchNames(name))
	})

	suite.Step(`^the observed branches "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		return state.fixture.DevRepo.Config.SetObservedBranches(gitdomain.NewLocalBranchNames(branch1, branch2))
	})

	suite.Step(`^the origin is "([^"]*)"$`, func(origin string) error {
		state.fixture.DevRepo.SetTestOrigin(origin)
		return nil
	})

	suite.Step(`^the parked branches "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		return state.fixture.DevRepo.Config.SetParkedBranches(gitdomain.NewLocalBranchNames(branch1, branch2))
	})

	suite.Step(`^the perennial branches are "([^"]+)"$`, func(name string) error {
		return state.fixture.DevRepo.Config.SetPerennialBranches(gitdomain.NewLocalBranchNames(name))
	})

	suite.Step(`^the perennial branches are "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		return state.fixture.DevRepo.Config.SetPerennialBranches(gitdomain.NewLocalBranchNames(branch1, branch2))
	})

	suite.Step(`^the perennial branches are not configured$`, func() error {
		return state.fixture.DevRepo.RemovePerennialBranchConfiguration()
	})

	suite.Step(`^the perennial branches are (?:now|still) "([^"]+)"$`, func(name string) error {
		actual := state.fixture.DevRepo.Config.LocalGitConfig.PerennialBranches
		if len(*actual) != 1 {
			return fmt.Errorf("expected 1 perennial branch, got %q", actual)
		}
		if (*actual)[0].String() != name {
			return fmt.Errorf("expected %q, got %q", name, (*actual)[0])
		}
		return nil
	})

	suite.Step(`^the perennial branches are now "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		actual := state.fixture.DevRepo.Config.LocalGitConfig.PerennialBranches
		if len(*actual) != 2 {
			return fmt.Errorf("expected 2 perennial branches, got %q", actual)
		}
		if (*actual)[0].String() != branch1 || (*actual)[1].String() != branch2 {
			return fmt.Errorf("expected %q, got %q", []string{branch1, branch2}, actual)
		}
		return nil
	})

	suite.Step(`^the previous Git branch is (?:now|still) "([^"]*)"$`, func(want string) error {
		have := state.fixture.DevRepo.BackendCommands.PreviouslyCheckedOutBranch()
		if have.String() != want {
			return fmt.Errorf("expected previous branch %q but got %q", want, have)
		}
		return nil
	})

	suite.Step(`^there are (?:now|still) no contribution branches$`, func() error {
		branches := state.fixture.DevRepo.Config.LocalGitConfig.ContributionBranches
		if branches != nil && len(*branches) > 0 {
			return fmt.Errorf("expected no contribution branches, got %q", branches)
		}
		return nil
	})

	suite.Step(`^there are (?:now|still) no observed branches$`, func() error {
		branches := state.fixture.DevRepo.Config.LocalGitConfig.ObservedBranches
		if branches != nil && len(*branches) > 0 {
			return fmt.Errorf("expected no observed branches, got %q", branches)
		}
		return nil
	})

	suite.Step(`^there are (?:now|still) no parked branches$`, func() error {
		branches := state.fixture.DevRepo.Config.LocalGitConfig.ParkedBranches
		if branches != nil && len(*branches) > 0 {
			return fmt.Errorf("expected no parked branches, got %q", branches)
		}
		return nil
	})

	suite.Step(`^there are (?:now|still) no perennial branches$`, func() error {
		branches := state.fixture.DevRepo.Config.LocalGitConfig.PerennialBranches
		if branches != nil && len(*branches) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", branches)
		}
		return nil
	})

	suite.Step(`^these committed files exist now$`, func(table *messages.PickleStepArgument_PickleTable) error {
		fileTable := state.fixture.DevRepo.FilesInBranches(gitdomain.NewLocalBranchName("main"))
		diff, errorCount := fileTable.EqualGherkin(table)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing files\n\n", errorCount)
			fmt.Println(diff)
			return errors.New("mismatching files found, see diff above")
		}
		return nil
	})

	suite.Step(`^these commits exist now$`, func(table *messages.PickleStepArgument_PickleTable) error {
		return state.compareGherkinTable(table)
	})

	suite.Step(`^these tags exist$`, func(table *messages.PickleStepArgument_PickleTable) error {
		tagTable := state.fixture.TagTable()
		diff, errorCount := tagTable.EqualGherkin(table)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing tags\n\n", errorCount)
			fmt.Println(diff)
			return errors.New("mismatching tags found, see diff above")
		}
		return nil
	})

	suite.Step(`^the tags$`, func(table *messages.PickleStepArgument_PickleTable) error {
		state.fixture.CreateTags(table)
		return nil
	})

	suite.Step(`^the uncommitted file is stashed$`, func() error {
		uncommittedFiles := state.fixture.DevRepo.UncommittedFiles()
		for _, ucf := range uncommittedFiles {
			if ucf == state.uncommittedFileName {
				return fmt.Errorf("expected file %q to be stashed but it is still uncommitted", state.uncommittedFileName)
			}
		}
		stashSize, err := state.fixture.DevRepo.StashSize()
		if err != nil {
			return err
		}
		if stashSize != 1 {
			return fmt.Errorf("expected 1 stash but found %d", stashSize)
		}
		return nil
	})

	suite.Step(`^the uncommitted file still exists$`, func() error {
		hasFile := state.fixture.DevRepo.HasFile(
			state.uncommittedFileName,
			state.uncommittedContent,
		)
		if hasFile != "" {
			return errors.New(hasFile)
		}
		return nil
	})

	suite.Step(`^these branches exist now$`, func(input *messages.PickleStepArgument_PickleTable) error {
		currentBranches := state.fixture.Branches()
		// fmt.Printf("NOW:\n%s\n", currentBranches.String())
		diff, errorCount := currentBranches.EqualGherkin(input)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing branches\n\n", errorCount)
			fmt.Println(diff)
			return errors.New("mismatching branches found, see diff above")
		}
		return nil
	})

	suite.Step(`^this lineage exists now$`, func(input *messages.PickleStepArgument_PickleTable) error {
		table := state.fixture.DevRepo.LineageTable()
		diff, errCount := table.EqualGherkin(input)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the lineage\n\n", errCount)
			fmt.Println(diff)
			return errors.New("mismatching branches found, see the diff above")
		}
		return nil
	})

	suite.Step(`^tool "([^"]*)" is broken$`, func(name string) error {
		state.fixture.DevRepo.MockBrokenCommand(name)
		return nil
	})

	suite.Step(`^tool "([^"]*)" is installed$`, func(tool string) error {
		state.fixture.DevRepo.MockCommand(tool)
		return nil
	})
}

func updateInitialSHAs(state *ScenarioState) {
	if len(state.initialDevSHAs) == 0 && state.insideGitRepo {
		state.initialDevSHAs = state.fixture.DevRepo.TestCommands.CommitSHAs()
	}
	if len(state.initialOriginSHAs) == 0 && state.insideGitRepo && state.fixture.OriginRepo != nil {
		state.initialOriginSHAs = state.fixture.OriginRepo.TestCommands.CommitSHAs()
	}
	if len(state.initialWorktreeSHAs) == 0 && state.insideGitRepo && state.fixture.SecondWorktree != nil {
		state.initialWorktreeSHAs = state.fixture.SecondWorktree.TestCommands.CommitSHAs()
	}
}
