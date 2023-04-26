package test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"github.com/acarl005/stripansi"
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/eiannone/keyboard"
	"github.com/git-town/git-town/v8/src/cli"
	"github.com/git-town/git-town/v8/src/config"
	"github.com/git-town/git-town/v8/src/stringslice"
	"github.com/git-town/git-town/v8/test/gherkin"
	"github.com/git-town/git-town/v8/test/git"
	"github.com/git-town/git-town/v8/test/output"
)

// beforeSuiteMux ensures that we run BeforeSuite only once globally.
var beforeSuiteMux sync.Mutex //nolint:gochecknoglobals

// the global FixtureFactory instance.
var fixtureFactory *FixtureFactory //nolint:gochecknoglobals

// Steps defines Cucumber step implementations around Git workspace management.
func Steps(suite *godog.Suite, state *ScenarioState) {
	suite.BeforeScenario(func(scenario *messages.Pickle) {
		// create a Fixture for the scenario
		fixture, err := fixtureFactory.CreateFixture(scenario.GetName())
		if err != nil {
			log.Fatalf("cannot create environment for scenario %q: %s", scenario.GetName(), err)
		}
		// Godog only provides state for the entire feature.
		// We want state to be scenario-specific, hence we reset the shared state before each scenario.
		// This is a limitation of the current Godog implementation, which doesn't have a `ScenarioContext` method,
		// only a `FeatureContext` method.
		// See main_test.go for additional details.
		state.Reset(fixture)
		if hasTag(scenario, "@debug") {
			state.fixture.DevRepo.Debug = true
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
			gm, err := NewFixtureFactory(evalBaseDir)
			if err != nil {
				log.Fatalf("Cannot create memoized environment: %s", err)
			}
			fixtureFactory = &gm
		}
	})

	suite.AfterScenario(func(scenario *messages.Pickle, e error) {
		if e != nil {
			fmt.Printf("failed scenario, investigate state in %q\n", state.fixture.Dir)
		}
		if state.runErr != nil && !state.runErrChecked {
			cli.PrintError(fmt.Errorf("%s - scenario %q doesn't document error %w", scenario.GetUri(), scenario.GetName(), state.runErr))
			os.Exit(1)
		}
	})

	suite.Step(`^a branch "([^"]*)"$`, func(branch string) error {
		state.initialLocalBranches = append(state.initialLocalBranches, branch)
		return state.fixture.DevRepo.CreateBranch(branch, "main")
	})

	suite.Step(`^a coworker clones the repository$`, func() error {
		return state.fixture.AddCoworkerRepo()
	})

	suite.Step(`^a feature branch "([^"]+)" as a child of "([^"]+)"$`, func(branch, parentBranch string) error {
		err := state.fixture.DevRepo.CreateChildFeatureBranch(branch, parentBranch)
		if err != nil {
			return fmt.Errorf("cannot create feature branch %q: %w", branch, err)
		}
		state.initialLocalBranches = append(state.initialLocalBranches, branch)
		state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
		state.initialBranchHierarchy.AddRow(branch, parentBranch)
		return state.fixture.DevRepo.PushBranchToRemote(branch, config.OriginRemote)
	})

	suite.Step(`^a merge is now in progress$`, func() error {
		if !state.fixture.DevRepo.HasMergeInProgress() {
			return fmt.Errorf("expected merge in progress")
		}
		return nil
	})

	suite.Step(`^a (local )?feature branch "([^"]*)"$`, func(localStr, branch string) error {
		isLocal := localStr != ""
		err := state.fixture.DevRepo.CreateFeatureBranch(branch)
		if err != nil {
			return err
		}
		state.initialLocalBranches = append(state.initialLocalBranches, branch)
		state.initialBranchHierarchy.AddRow(branch, "main")
		if !isLocal {
			state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
			return state.fixture.DevRepo.PushBranchToRemote(branch, config.OriginRemote)
		}
		return nil
	})

	suite.Step(`^a perennial branch "([^"]+)"$`, func(branch string) error {
		err := state.fixture.DevRepo.CreatePerennialBranches(branch)
		if err != nil {
			return fmt.Errorf("cannot create perennial branch: %w", err)
		}
		state.initialLocalBranches = append(state.initialLocalBranches, branch)
		state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
		return state.fixture.DevRepo.PushBranchToRemote(branch, config.OriginRemote)
	})

	suite.Step(`^a rebase is now in progress$`, func() error {
		hasRebase, err := state.fixture.DevRepo.HasRebaseInProgress()
		if err != nil {
			return err
		}
		if !hasRebase {
			return fmt.Errorf("expected rebase in progress")
		}
		return nil
	})

	suite.Step(`^a remote feature branch "([^"]*)"$`, func(branch string) error {
		state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
		return state.fixture.OriginRepo.CreateBranch(branch, "main")
	})

	suite.Step(`^a remote tag "([^"]+)" not on a branch$`, func(name string) error {
		return state.fixture.OriginRepo.CreateStandaloneTag(name)
	})

	suite.Step(`^all branches are now synchronized$`, func() error {
		outOfSync, err := state.fixture.DevRepo.HasBranchesOutOfSync()
		if err != nil {
			return err
		}
		if outOfSync {
			return fmt.Errorf("expected no branches out of sync")
		}
		return nil
	})

	suite.Step(`^an uncommitted file$`, func() error {
		state.uncommittedFileName = "uncommitted file"
		state.uncommittedContent = "uncommitted content"
		return state.fixture.DevRepo.CreateFile(
			state.uncommittedFileName,
			state.uncommittedContent,
		)
	})

	suite.Step(`^an uncommitted file in folder "([^"]*)"$`, func(folder string) error {
		state.uncommittedFileName = fmt.Sprintf("%s/uncommitted file", folder)
		return state.fixture.DevRepo.CreateFile(
			state.uncommittedFileName,
			state.uncommittedContent,
		)
	})

	suite.Step(`^an uncommitted file with name "([^"]+)" and content "([^"]+)"$`, func(name, content string) error {
		state.uncommittedFileName = name
		state.uncommittedContent = content
		return state.fixture.DevRepo.CreateFile(name, content)
	})

	suite.Step(`^an upstream repo$`, func() error {
		return state.fixture.AddUpstream()
	})

	suite.Step(`^file "([^"]+)" still contains unresolved conflicts$`, func(name string) error {
		content, err := state.fixture.DevRepo.FileContent(name)
		if err != nil {
			return fmt.Errorf("cannot read file %q: %w", name, err)
		}
		if !strings.Contains(content, "<<<<<<<") {
			return fmt.Errorf("file %q does not contain unresolved conflicts", name)
		}
		return nil
	})

	suite.Step(`^file "([^"]*)" still has content "([^"]*)"$`, func(file, expectedContent string) error {
		actualContent, err := state.fixture.DevRepo.FileContent(file)
		if err != nil {
			return err
		}
		if expectedContent != actualContent {
			return fmt.Errorf("file content does not match\n\nEXPECTED: %q\n\nACTUAL:\n\n%q\n----------------------------", expectedContent, actualContent)
		}
		return nil
	})

	suite.Step(`^Git has version "([^"]*)"$`, func(version string) error {
		err := state.fixture.DevRepo.MockGit(version)
		return err
	})

	suite.Step(`^Git Town is no longer configured$`, func() error {
		if state.fixture.DevRepo.HasGitTownConfigNow() {
			return fmt.Errorf("unexpected Git Town configuration")
		}
		return nil
	})

	suite.Step(`^Git Town is not configured$`, func() error {
		err := state.fixture.DevRepo.Config.RemovePerennialBranchConfiguration()
		if err != nil {
			return err
		}
		return state.fixture.DevRepo.DeleteMainBranchConfiguration()
	})

	suite.Step(`^I add commit "([^"]*)" to the "([^"]*)" branch`, func(message, branch string) error {
		return state.fixture.DevRepo.CreateCommit(git.Commit{
			Branch:      branch,
			FileName:    "new_file",
			FileContent: "new content",
			Message:     message,
		})
	})

	suite.Step(`^I am not prompted for any parent branches$`, func() error {
		notExpected := "Please specify the parent branch of"
		if strings.Contains(state.runOutput, notExpected) {
			return fmt.Errorf("text found:\n\nDID NOT EXPECT: %q\n\nACTUAL\n\n%q\n----------------------------", notExpected, state.runOutput)
		}
		return nil
	})

	suite.Step(`^I am outside a Git repo$`, func() error {
		os.RemoveAll(filepath.Join(state.fixture.DevRepo.WorkingDir(), ".git"))
		return nil
	})

	suite.Step(`^I resolve the conflict in "([^"]*)"(?: with "([^"]*)")?$`, func(filename, content string) error {
		if content == "" {
			content = "resolved content"
		}
		err := state.fixture.DevRepo.CreateFile(filename, content)
		if err != nil {
			return err
		}
		err = state.fixture.DevRepo.StageFiles(filename)
		if err != nil {
			return err
		}
		return nil
	})

	suite.Step(`^I (?:run|ran) "(.+)"$`, func(command string) error {
		state.runOutput, state.runErr = state.fixture.DevRepo.RunString(command)
		return nil
	})

	suite.Step(`^I (?:run|ran) "([^"]+)" and answer(?:ed)? the prompts:$`, func(cmd string, input *messages.PickleStepArgument_PickleTable) error {
		state.runOutput, state.runErr = state.fixture.DevRepo.RunStringWith(cmd, &Options{Input: tableToInput(input)})
		return nil
	})

	suite.Step(`^I run "([^"]*)" and close the editor$`, func(cmd string) error {
		env := append(os.Environ(), "GIT_EDITOR=true")
		state.runOutput, state.runErr = state.fixture.DevRepo.RunStringWith(cmd, &Options{Env: env})
		return nil
	})

	suite.Step(`^I run "([^"]*)" and enter an empty commit message$`, func(cmd string) error {
		if err := state.fixture.DevRepo.MockCommitMessage(""); err != nil {
			return err
		}
		state.runOutput, state.runErr = state.fixture.DevRepo.RunString(cmd)
		return nil
	})

	suite.Step(`^I run "([^"]*)" and enter "([^"]*)" for the commit message$`, func(cmd, message string) error {
		if err := state.fixture.DevRepo.MockCommitMessage(message); err != nil {
			return err
		}
		state.runOutput, state.runErr = state.fixture.DevRepo.RunString(cmd)
		return nil
	})

	suite.Step(`^I run "([^"]*)", answer the prompts, and close the next editor:$`, func(cmd string, input *messages.PickleStepArgument_PickleTable) error {
		env := append(os.Environ(), "GIT_EDITOR=true")
		state.runOutput, state.runErr = state.fixture.DevRepo.RunStringWith(cmd, &Options{Env: env, Input: tableToInput(input)})
		return nil
	})

	suite.Step(`^I run "([^"]+)" in the "([^"]+)" folder$`, func(cmd, folderName string) error {
		state.runOutput, state.runErr = state.fixture.DevRepo.RunStringWith(cmd, &Options{Dir: folderName})
		return nil
	})

	suite.Step(`^inspect the repo$`, func() error {
		fmt.Printf("\nThe workspace is at %q\n", state.fixture.DevRepo.WorkingDir())
		_, _, err := keyboard.GetSingleKey()
		if err != nil {
			return fmt.Errorf("cannot read from os.Stdin: %w", err)
		}
		return nil
	})

	suite.Step(`^it does not print "(.+)"$`, func(text string) error {
		if strings.Contains(stripansi.Strip(state.runOutput), text) {
			return fmt.Errorf("text found: %q", text)
		}
		return nil
	})

	suite.Step(`^it prints:$`, func(expected *messages.PickleStepArgument_PickleDocString) error {
		if state.runErr != nil {
			return fmt.Errorf("unexpected error: %w", state.runErr)
		}
		if !strings.Contains(stripansi.Strip(state.runOutput), expected.Content) {
			return fmt.Errorf("text not found:\n\nEXPECTED:\n\n%q\n\nACTUAL:\n\n%q", expected.Content, state.runOutput)
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
		state.runErrChecked = true
		if !strings.Contains(stripansi.Strip(state.runOutput), expected.Content) {
			return fmt.Errorf("text not found:\n%s\n\nactual text:\n%s", expected.Content, state.runOutput)
		}
		if state.runErr == nil {
			return fmt.Errorf("expected error")
		}
		return nil
	})

	suite.Step(`^it runs no commands$`, func() error {
		commands := output.GitCommandsInGitTownOutput(state.runOutput)
		if len(commands) > 0 {
			for _, command := range commands {
				fmt.Println(command)
			}
			return fmt.Errorf("expected no commands but found %d commands", len(commands))
		}
		return nil
	})

	suite.Step(`^it runs the commands$`, func(input *messages.PickleStepArgument_PickleTable) error {
		commands := output.GitCommandsInGitTownOutput(state.runOutput)
		table := RenderExecutedGitCommands(commands, input)
		dataTable := gherkin.FromGherkin(input)
		expanded, err := dataTable.Expand(
			&state.fixture.DevRepo,
			state.fixture.OriginRepo,
		)
		if err != nil {
			return err
		}
		diff, errorCount := table.EqualDataTable(expanded)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the commands run\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching commands run, see diff above")
		}
		return nil
	})

	suite.Step(`^it runs without error$`, func() error {
		if state.runErr != nil {
			return fmt.Errorf("did not expect the Git Town command to produce an error: %w", state.runErr)
		}
		return nil
	})

	suite.Step(`^"([^"]*)" launches a new pull request with this url in my browser:$`, func(tool string, url *messages.PickleStepArgument_PickleDocString) error {
		want := fmt.Sprintf("%s called with: %s", tool, url.Content)
		want = strings.ReplaceAll(want, "?", `\?`)
		regex := regexp.MustCompile(want)
		have := state.runOutput
		if !regex.MatchString(have) {
			return fmt.Errorf("EXPECTED: a regex matching %q\nGOT: %q", want, have)
		}
		return nil
	})

	suite.Step(`^my repo does not have an origin$`, func() error {
		err := state.fixture.DevRepo.RemoveRemote(config.OriginRemote)
		if err != nil {
			return err
		}
		state.initialRemoteBranches = []string{}
		state.fixture.OriginRepo = nil
		return nil
	})

	suite.Step(`^my repo has a Git submodule$`, func() error {
		err := state.fixture.AddSubmoduleRepo()
		if err != nil {
			return err
		}
		return state.fixture.DevRepo.AddSubmodule(state.fixture.SubmoduleRepo.WorkingDir())
	})

	suite.Step(`^no branch hierarchy exists now$`, func() error {
		state.fixture.DevRepo.Config.Reload()
		if state.fixture.DevRepo.Config.HasBranchInformation() {
			branchInfo := state.fixture.DevRepo.Config.ParentBranchMap()
			return fmt.Errorf("unexpected Git Town branch hierarchy information: %+v", branchInfo)
		}
		return nil
	})

	suite.Step(`^no merge is in progress$`, func() error {
		if state.fixture.DevRepo.HasMergeInProgress() {
			return fmt.Errorf("expected no merge in progress")
		}
		return nil
	})

	suite.Step(`^no rebase is in progress$`, func() error {
		hasRebase, err := state.fixture.DevRepo.HasRebaseInProgress()
		if err != nil {
			return err
		}
		if hasRebase {
			return fmt.Errorf("expected no rebase in progress")
		}
		return nil
	})

	suite.Step(`^no tool to open browsers is installed$`, func() error {
		return state.fixture.DevRepo.MockNoCommandsInstalled()
	})

	suite.Step(`^no uncommitted files exist$`, func() error {
		files, err := state.fixture.DevRepo.UncommittedFiles()
		if err != nil {
			return fmt.Errorf("cannot determine uncommitted files: %w", err)
		}
		if len(files) > 0 {
			return fmt.Errorf("unexpected uncommitted files: %s", files)
		}
		return nil
	})

	suite.Step(`^now the initial commits exist$`, func() error {
		return compareExistingCommits(state, state.initialCommits)
	})

	suite.Step(`^now these commits exist$`, func(table *messages.PickleStepArgument_PickleTable) error {
		return compareExistingCommits(state, table)
	})

	suite.Step(`^offline mode is disabled$`, func() error {
		state.fixture.DevRepo.Config.Reload()
		isOffline, err := state.fixture.DevRepo.Config.IsOffline()
		if err != nil {
			return err
		}
		if isOffline {
			return fmt.Errorf("expected to not be offline but am")
		}
		return nil
	})

	suite.Step(`^offline mode is enabled$`, func() error {
		return state.fixture.DevRepo.Config.SetOffline(true)
	})

	suite.Step(`^origin deletes the "([^"]*)" branch$`, func(name string) error {
		state.initialRemoteBranches = stringslice.Remove(state.initialRemoteBranches, name)
		return state.fixture.OriginRepo.RemoveBranch(name)
	})

	suite.Step(`^Git setting "color.ui" is "([^"]*)"$`, func(value string) error {
		return state.fixture.DevRepo.Config.SetColorUI(value)
	})

	suite.Step(`^(?:local )?setting "([^"]*)" is "([^"]*)"$`, func(name, value string) error {
		_, err := state.fixture.DevRepo.Config.SetLocalConfigValue("git-town."+name, value)
		return err
	})

	suite.Step(`^global setting "([^"]*)" is "([^"]*)"$`, func(name, value string) error {
		_, err := state.fixture.DevRepo.Config.SetGlobalConfigValue("git-town."+name, value)
		return err
	})

	suite.Step(`^local setting "([^"]*)" no longer exists$`, func(name string) error {
		state.fixture.DevRepo.Config.Reload()
		newValue := state.fixture.DevRepo.Config.LocalConfigValue("git-town." + name)
		if newValue == "" {
			return nil
		}
		return fmt.Errorf("should not have local %q anymore but has value %q", name, newValue)
	})

	suite.Step(`^global setting "([^"]*)" no longer exists$`, func(name string) error {
		state.fixture.DevRepo.Config.Reload()
		newValue := state.fixture.DevRepo.Config.GlobalConfigValue("git-town." + name)
		if newValue == "" {
			return nil
		}
		return fmt.Errorf("should not have global %q anymore but has value %q", name, newValue)
	})

	suite.Step(`^setting "([^"]*)" is now "([^"]*)"$`, func(name, want string) error {
		state.fixture.DevRepo.Config.Reload()
		have := state.fixture.DevRepo.Config.LocalOrGlobalConfigValue("git-town." + name)
		if have != want {
			return fmt.Errorf("expected setting %q to be %q, but was %q", name, want, have)
		}
		return nil
	})

	suite.Step(`^local setting "([^"]*)" is now "([^"]*)"$`, func(name, want string) error {
		state.fixture.DevRepo.Config.Reload()
		have := state.fixture.DevRepo.Config.LocalConfigValue("git-town." + name)
		if have != want {
			return fmt.Errorf("expected local setting %q to be %q, but was %q", name, want, have)
		}
		return nil
	})

	suite.Step(`^global setting "([^"]*)" is (?:now|still) "([^"]*)"$`, func(name, want string) error {
		state.fixture.DevRepo.Config.Reload()
		have := state.fixture.DevRepo.Config.GlobalConfigValue("git-town." + name)
		if have != want {
			return fmt.Errorf("expected global setting %q to be %q, but was %q", name, want, have)
		}
		return nil
	})

	suite.Step(`^the branches "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		for _, branch := range []string{branch1, branch2} {
			err := state.fixture.DevRepo.CreateBranch(branch, "main")
			if err != nil {
				return err
			}
			state.initialLocalBranches = append(state.initialLocalBranches, branch)
		}
		return nil
	})

	suite.Step(`^the branches are now$`, func(table *messages.PickleStepArgument_PickleTable) error {
		existing, err := state.fixture.Branches()
		if err != nil {
			return err
		}
		diff, errCount := existing.EqualGherkin(table)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the branches\n\n", errCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching branches found, see the diff above")
		}
		return nil
	})

	suite.Step(`^the commits$`, func(table *messages.PickleStepArgument_PickleTable) error {
		state.initialCommits = table
		// create the commits
		commits, err := git.FromGherkinTable(table)
		if err != nil {
			return fmt.Errorf("cannot parse Gherkin table: %w", err)
		}
		err = state.fixture.CreateCommits(commits)
		if err != nil {
			return fmt.Errorf("cannot create commits: %w", err)
		}
		// restore the initial branch
		if state.initialCurrentBranch == "" {
			return state.fixture.DevRepo.CheckoutBranch("main")
		}
		if state.fixture.DevRepo.config.CurrentBranchCache.Value() != state.initialCurrentBranch {
			return state.fixture.DevRepo.CheckoutBranch(state.initialCurrentBranch)
		}
		return nil
	})

	suite.Step(`^the coworker fetches updates$`, func() error {
		return state.fixture.CoworkerRepo.Fetch()
	})

	suite.Step(`^the coworker is on the "([^"]*)" branch$`, func(branch string) error {
		return state.fixture.CoworkerRepo.CheckoutBranch(branch)
	})

	suite.Step(`^the coworker runs "([^"]+)"$`, func(command string) error {
		state.runOutput, state.runErr = state.fixture.CoworkerRepo.RunString(command)
		return nil
	})

	suite.Step(`^the coworker sets the parent branch of "([^"]*)" as "([^"]*)"$`, func(childBranch, parentBranch string) error {
		_ = state.fixture.CoworkerRepo.Config.SetParent(childBranch, parentBranch)
		return nil
	})

	suite.Step(`^the coworker sets the "sync-strategy" to "(merge|rebase)"$`, func(value string) error {
		syncStrategy, err := config.ToSyncStrategy(value)
		if err != nil {
			return err
		}
		_ = state.fixture.CoworkerRepo.Config.SetSyncStrategy(syncStrategy)
		return nil
	})

	suite.Step(`^the current branch is "([^"]*)"$`, func(name string) error {
		state.initialCurrentBranch = name
		if !stringslice.Contains(state.initialLocalBranches, name) {
			state.initialLocalBranches = append(state.initialLocalBranches, name)
			err := state.fixture.DevRepo.CreateBranch(name, "main")
			if err != nil {
				return err
			}
		}
		return state.fixture.DevRepo.CheckoutBranch(name)
	})

	suite.Step(`^the current branch is a (local )?(feature|perennial) branch "([^"]*)"$`, func(localStr, branchType, branch string) error {
		isLocal := localStr != ""
		var err error
		switch branchType {
		case "feature":
			err = state.fixture.DevRepo.CreateFeatureBranch(branch)
		case "perennial":
			err = state.fixture.DevRepo.CreatePerennialBranches(branch)
		default:
			panic(fmt.Sprintf("unknown branch type: %q", branchType))
		}
		if err != nil {
			return err
		}
		state.initialLocalBranches = append(state.initialLocalBranches, branch)
		if branchType == "feature" {
			state.initialBranchHierarchy.AddRow(branch, "main")
		}
		if !isLocal {
			state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
			err := state.fixture.DevRepo.PushBranchToRemote(branch, config.OriginRemote)
			if err != nil {
				return err
			}
		}
		state.initialCurrentBranch = branch
		if !state.fixture.DevRepo.config.CurrentBranchCache.Initialized() || state.fixture.DevRepo.config.CurrentBranchCache.Value() != branch {
			return state.fixture.DevRepo.CheckoutBranch(branch)
		}
		return nil
	})

	suite.Step(`^the current branch is "([^"]*)" and the previous branch is "([^"]*)"$`, func(current, previous string) error {
		state.initialCurrentBranch = current
		err := state.fixture.DevRepo.CheckoutBranch(previous)
		if err != nil {
			return err
		}
		return state.fixture.DevRepo.CheckoutBranch(current)
	})

	suite.Step(`^the current branch is (?:now|still) "([^"]*)"$`, func(expected string) error {
		state.fixture.DevRepo.config.CurrentBranchCache.Invalidate()
		actual, err := state.fixture.DevRepo.CurrentBranch()
		if err != nil {
			return fmt.Errorf("cannot determine current branch of developer repo: %w", err)
		}
		if actual != expected {
			return fmt.Errorf("expected active branch %q but is %q", expected, actual)
		}
		return nil
	})

	suite.Step(`^the initial branch hierarchy exists$`, func() error {
		have := state.fixture.DevRepo.BranchHierarchyTable()
		state.initialBranchHierarchy.Sort()
		diff, errCnt := have.EqualDataTable(state.initialBranchHierarchy)
		if errCnt > 0 {
			fmt.Printf("\nERROR! Found %d differences in the branch hierarchy\n\n", errCnt)
			fmt.Printf("INITIAL BRANCH HIERARCHY:\n%s\n", state.initialBranchHierarchy.String())
			fmt.Printf("CURRENT BRANCH HIERARCHY:\n%s\n", have.String())
			fmt.Println(diff)
			return fmt.Errorf("mismatching branches found, see the diff above")
		}
		return nil
	})

	suite.Step(`^the initial branches and hierarchy exist$`, func() error {
		// verify initial branches
		have, err := state.fixture.Branches()
		if err != nil {
			return err
		}
		want := state.InitialBranches()
		diff, errorCount := have.EqualDataTable(want)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing branches\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching branches found, see diff above")
		}
		// verify initial branch hierarchy
		state.initialBranchHierarchy.Sort()
		have = state.fixture.DevRepo.BranchHierarchyTable()
		diff, errCnt := have.EqualDataTable(state.initialBranchHierarchy)
		if errCnt > 0 {
			fmt.Printf("\nERROR! Found %d differences in the branch hierarchy\n\n", errCnt)
			fmt.Println(diff)
			return fmt.Errorf("mismatching branch hierarchy found, see the diff above")
		}
		return nil
	})

	suite.Step(`^the initial branches exist$`, func() error {
		have, err := state.fixture.Branches()
		if err != nil {
			return err
		}
		want := state.InitialBranches()
		// fmt.Printf("HAVE:\n%s\n", have.String())
		// fmt.Printf("WANT:\n%s\n", want.String())
		diff, errorCount := have.EqualDataTable(want)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing branches\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching branches found, see diff above")
		}
		return nil
	})

	suite.Step(`^the (local )?feature branches "([^"]+)" and "([^"]+)"$`, func(localStr, branch1, branch2 string) error {
		isLocal := localStr != ""
		for _, branch := range []string{branch1, branch2} {
			err := state.fixture.DevRepo.CreateFeatureBranch(branch)
			if err != nil {
				return err
			}
			state.initialLocalBranches = append(state.initialLocalBranches, branch)
			state.initialBranchHierarchy.AddRow(branch, "main")
			if !isLocal {
				err = state.fixture.DevRepo.PushBranchToRemote(branch, config.OriginRemote)
				if err != nil {
					return err
				}
				state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
			}
		}
		return nil
	})

	suite.Step(`^the (local )?feature branches "([^"]+)", "([^"]+)", and "([^"]+)"$`, func(localStr, branch1, branch2, branch3 string) error {
		isLocal := localStr != ""
		for _, branch := range []string{branch1, branch2, branch3} {
			err := state.fixture.DevRepo.CreateFeatureBranch(branch)
			if err != nil {
				return err
			}
			state.initialLocalBranches = append(state.initialLocalBranches, branch)
			state.initialBranchHierarchy.AddRow(branch, "main")
			if !isLocal {
				err = state.fixture.DevRepo.PushBranchToRemote(branch, config.OriginRemote)
				if err != nil {
					return err
				}
				state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
			}
		}
		return nil
	})

	suite.Step(`^the (local )?perennial branches "([^"]+)" and "([^"]+)"$`, func(localStr, branch1, branch2 string) error {
		isLocal := localStr != ""
		err := state.fixture.DevRepo.CreatePerennialBranches(branch1, branch2)
		if err != nil {
			return fmt.Errorf("cannot create perennial branches: %w", err)
		}
		state.initialLocalBranches = append(state.initialLocalBranches, branch1, branch2)
		if !isLocal {
			state.initialRemoteBranches = append(state.initialRemoteBranches, branch1, branch2)
			err = state.fixture.DevRepo.PushBranchToRemote(branch1, config.OriginRemote)
			if err != nil {
				return err
			}
			return state.fixture.DevRepo.PushBranchToRemote(branch2, config.OriginRemote)
		}
		return nil
	})

	suite.Step(`^the (local )?perennial branches "([^"]+)", "([^"]+)", and "([^"]+)"$`, func(localStr, branch1, branch2, branch3 string) error {
		isLocal := localStr != ""
		for _, branch := range []string{branch1, branch2, branch3} {
			err := state.fixture.DevRepo.CreatePerennialBranches(branch)
			if err != nil {
				return fmt.Errorf("cannot create perennial branches: %w", err)
			}
			state.initialLocalBranches = append(state.initialLocalBranches, branch)
			if !isLocal {
				err = state.fixture.DevRepo.PushBranchToRemote(branch, config.OriginRemote)
				if err != nil {
					return fmt.Errorf("cannot push perennial branch upstream: %w", err)
				}
				state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
			}
		}
		return nil
	})

	suite.Step(`^the main branch is "([^"]+)"$`, func(name string) error {
		return state.fixture.DevRepo.Config.SetMainBranch(name)
	})

	suite.Step(`^the main branch is now "([^"]+)"$`, func(name string) error {
		state.fixture.DevRepo.Config.Reload()
		actual := state.fixture.DevRepo.Config.MainBranch()
		if actual != name {
			return fmt.Errorf("expected %q, got %q", name, actual)
		}
		return nil
	})

	suite.Step(`^the origin is "([^"]*)"$`, func(origin string) error {
		state.fixture.DevRepo.SetTestOrigin(origin)
		return nil
	})

	suite.Step(`^the perennial branches are "([^"]+)"$`, func(name string) error {
		return state.fixture.DevRepo.Config.AddToPerennialBranches(name)
	})

	suite.Step(`^the perennial branches are "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		return state.fixture.DevRepo.Config.AddToPerennialBranches(branch1, branch2)
	})

	suite.Step(`^the perennial branches are not configured$`, func() error {
		return state.fixture.DevRepo.Config.RemovePerennialBranchConfiguration()
	})

	suite.Step(`^the perennial branches are now "([^"]+)"$`, func(name string) error {
		state.fixture.DevRepo.Config.Reload()
		actual := state.fixture.DevRepo.Config.PerennialBranches()
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 perennial branch, got %q", actual)
		}
		if actual[0] != name {
			return fmt.Errorf("expected %q, got %q", name, actual[0])
		}
		return nil
	})

	suite.Step(`^the perennial branches are now "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		state.fixture.DevRepo.Config.Reload()
		actual := state.fixture.DevRepo.Config.PerennialBranches()
		if len(actual) != 2 {
			return fmt.Errorf("expected 2 perennial branches, got %q", actual)
		}
		if actual[0] != branch1 || actual[1] != branch2 {
			return fmt.Errorf("expected %q, got %q", []string{branch1, branch2}, actual)
		}
		return nil
	})

	suite.Step(`^the previous Git branch is (?:now|still) "([^"]*)"$`, func(want string) error {
		err := state.fixture.DevRepo.CheckoutBranch("-")
		if err != nil {
			return err
		}
		have, err := state.fixture.DevRepo.CurrentBranch()
		if err != nil {
			return err
		}
		if have != want {
			return fmt.Errorf("expected previous branch %q but got %q", want, have)
		}
		return state.fixture.DevRepo.CheckoutBranch("-")
	})

	suite.Step(`^the tags$`, func(table *messages.PickleStepArgument_PickleTable) error {
		return state.fixture.CreateTags(table)
	})

	suite.Step(`^the uncommitted file is stashed$`, func() error {
		uncommittedFiles, err := state.fixture.DevRepo.UncommittedFiles()
		if err != nil {
			return err
		}
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
		hasFile, err := state.fixture.DevRepo.HasFile(
			state.uncommittedFileName,
			state.uncommittedContent,
		)
		if err != nil {
			return err
		}
		if !hasFile {
			return fmt.Errorf("expected file %q but didn't find it", state.uncommittedFileName)
		}
		return nil
	})

	suite.Step(`^there are still no perennial branches$`, func() error {
		state.fixture.DevRepo.Config.Reload()
		branches := state.fixture.DevRepo.Config.PerennialBranches()
		if len(branches) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", branches)
		}
		return nil
	})

	suite.Step(`^these committed files exist now$`, func(table *messages.PickleStepArgument_PickleTable) error {
		fileTable, err := state.fixture.DevRepo.FilesInBranches("main")
		if err != nil {
			return fmt.Errorf("cannot determine files in branches in the developer repo: %w", err)
		}
		diff, errorCount := fileTable.EqualGherkin(table)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing files\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching files found, see diff above")
		}
		return nil
	})

	suite.Step(`^these tags exist$`, func(table *messages.PickleStepArgument_PickleTable) error {
		tagTable, err := state.fixture.TagTable()
		if err != nil {
			return err
		}
		diff, errorCount := tagTable.EqualGherkin(table)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing tags\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching tags found, see diff above")
		}
		return nil
	})

	suite.Step(`^this branch hierarchy exists now$`, func(input *messages.PickleStepArgument_PickleTable) error {
		table := state.fixture.DevRepo.BranchHierarchyTable()
		diff, errCount := table.EqualGherkin(input)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the branch hierarchy\n\n", errCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching branches found, see the diff above")
		}
		return nil
	})

	suite.Step(`^tool "([^"]*)" is broken$`, func(name string) error {
		return state.fixture.DevRepo.MockBrokenCommand(name)
	})

	suite.Step(`^tool "([^"]*)" is installed$`, func(tool string) error {
		return state.fixture.DevRepo.MockCommand(tool)
	})
}
