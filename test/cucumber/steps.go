package cucumber

import (
	"errors"
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
	"github.com/git-town/git-town/v11/src/cli/print"
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/gohacks/slice"
	"github.com/git-town/git-town/v11/test/asserts"
	"github.com/git-town/git-town/v11/test/datatable"
	"github.com/git-town/git-town/v11/test/fixture"
	"github.com/git-town/git-town/v11/test/git"
	"github.com/git-town/git-town/v11/test/helpers"
	"github.com/git-town/git-town/v11/test/output"
	"github.com/git-town/git-town/v11/test/subshell"
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
		state.initialLocalBranches = append(state.initialLocalBranches, domain.NewLocalBranchName(branch))
		state.fixture.DevRepo.CreateBranch(domain.NewLocalBranchName(branch), domain.NewLocalBranchName("main"))
		return nil
	})

	suite.Step(`^a coworker clones the repository$`, func() error {
		state.fixture.AddCoworkerRepo()
		return nil
	})

	suite.Step(`^a feature branch "([^"]+)" as a child of "([^"]+)"$`, func(branchText, parentBranch string) error {
		branch := domain.NewLocalBranchName(branchText)
		state.fixture.DevRepo.CreateChildFeatureBranch(branch, domain.NewLocalBranchName(parentBranch))
		state.initialLocalBranches = append(state.initialLocalBranches, branch)
		state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
		state.initialLineage.AddRow(branchText, parentBranch)
		state.fixture.DevRepo.PushBranchToRemote(branch, domain.OriginRemote)
		return nil
	})

	suite.Step(`^a merge is now in progress$`, func() error {
		if !state.fixture.DevRepo.HasMergeInProgress() {
			return fmt.Errorf("expected merge in progress")
		}
		return nil
	})

	suite.Step(`^a (local )?feature branch "([^"]*)"$`, func(localStr, branchText string) error {
		branch := domain.NewLocalBranchName(branchText)
		isLocal := localStr != ""
		asserts.NoError(state.fixture.DevRepo.CreateFeatureBranch(branch))
		state.initialLocalBranches = append(state.initialLocalBranches, branch)
		state.initialLineage.AddRow(branchText, "main")
		if !isLocal {
			state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
			state.fixture.DevRepo.PushBranchToRemote(branch, domain.OriginRemote)
			return nil
		}
		return nil
	})

	suite.Step(`^a perennial branch "([^"]+)"$`, func(branchText string) error {
		branch := domain.NewLocalBranchName(branchText)
		state.fixture.DevRepo.CreatePerennialBranches(branch)
		state.initialLocalBranches = append(state.initialLocalBranches, branch)
		state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
		state.fixture.DevRepo.PushBranchToRemote(branch, domain.OriginRemote)
		return nil
	})

	suite.Step(`^a rebase is now in progress$`, func() error {
		repoStatus, err := state.fixture.DevRepo.RepoStatus()
		asserts.NoError(err)
		if !repoStatus.RebaseInProgress {
			return fmt.Errorf("expected rebase in progress")
		}
		return nil
	})

	suite.Step(`^a remote feature branch "([^"]*)"$`, func(branchText string) error {
		branch := domain.NewLocalBranchName(branchText)
		state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
		// we are creating a remote branch in the remote repo --> it is a local branch there
		state.fixture.OriginRepo.CreateBranch(branch, domain.NewLocalBranchName("main"))
		return nil
	})

	suite.Step(`^a remote tag "([^"]+)" not on a branch$`, func(name string) error {
		state.fixture.OriginRepo.CreateStandaloneTag(name)
		return nil
	})

	suite.Step(`^all branches are now synchronized$`, func() error {
		if state.fixture.DevRepo.HasBranchesOutOfSync() {
			return fmt.Errorf("expected no branches out of sync")
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
		state.uncommittedFileName = fmt.Sprintf("%s/uncommitted file", folder)
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
		state.fixture.AddSecondWorktree(domain.NewLocalBranchName(branch))
		return nil
	})

	suite.Step(`^display "([^"]+)"$`, func(command string) error {
		parts := strings.Split(command, " ")
		output, err := state.fixture.DevRepo.BackendRunner.Query(parts[0], parts[1:]...)
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

	suite.Step(`^file "([^"]*)" still has content "([^"]*)"$`, func(file, expectedContent string) error {
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
		if state.fixture.DevRepo.HasGitTownConfigNow() {
			return fmt.Errorf("unexpected Git Town configuration")
		}
		return nil
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
		branchName := domain.NewLocalBranchName(branch)
		configKey := configdomain.NewParentKey(branchName)
		return state.fixture.DevRepo.GitTown.SetLocalConfigValue(configKey, value)
	})

	suite.Step(`^Git Town setting "([^"]*)" is now "([^"]*)"$`, func(name, want string) error {
		configKey := configdomain.ParseKey("git-town." + name)
		have := state.fixture.DevRepo.GitTown.LocalOrGlobalConfigValue(*configKey)
		if have != want {
			return fmt.Errorf("expected setting %q to be %q, but was %q", name, want, have)
		}
		return nil
	})

	suite.Step(`^global Git Town setting "([^"]*)" is "([^"]*)"$`, func(name, value string) error {
		configKey := configdomain.ParseKey("git-town." + name)
		return state.fixture.DevRepo.GitTown.SetGlobalConfigValue(*configKey, value)
	})

	suite.Step(`^global Git Town setting "([^"]*)" no longer exists$`, func(name string) error {
		configKey := configdomain.ParseKey("git-town." + name)
		newValue := state.fixture.DevRepo.GitTown.GlobalConfigValue(*configKey)
		if newValue == "" {
			return nil
		}
		return fmt.Errorf("should not have global %q anymore but has value %q", name, newValue)
	})

	suite.Step(`^global Git Town setting "([^"]*)" is (?:now|still) "([^"]*)"$`, func(name, want string) error {
		configKey := configdomain.ParseKey("git-town." + name)
		have := state.fixture.DevRepo.GitTown.GlobalConfigValue(*configKey)
		if have != want {
			return fmt.Errorf("expected global setting %q to be %q, but was %q", name, want, have)
		}
		return nil
	})

	suite.Step(`^I add commit "([^"]*)" to the "([^"]*)" branch`, func(message, branch string) error {
		state.fixture.DevRepo.CreateCommit(git.Commit{
			Branch:      domain.NewLocalBranchName(branch),
			FileName:    "new_file",
			FileContent: "new content",
			Message:     message,
		})
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

	suite.Step(`^I (?:run|ran) "(.+)"$`, func(command string) error {
		updateInitialSHAs(state)
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCode(command)
		state.fixture.DevRepo.GitTown.Reload()
		return nil
	})

	suite.Step(`^I (?:run|ran) "([^"]+)" and answer(?:ed)? the prompts:$`, func(cmd string, input *messages.PickleStepArgument_PickleTable) error {
		updateInitialSHAs(state)
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Input: helpers.TableToInput(input)})
		state.fixture.DevRepo.GitTown.Reload()
		return nil
	})

	suite.Step(`^I run "([^"]*)" and close the editor$`, func(cmd string) error {
		updateInitialSHAs(state)
		env := append(os.Environ(), "GIT_EDITOR=true")
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Env: env})
		state.fixture.DevRepo.GitTown.Reload()
		return nil
	})

	suite.Step(`^I run "([^"]*)" and enter an empty commit message$`, func(cmd string) error {
		updateInitialSHAs(state)
		state.fixture.DevRepo.MockCommitMessage("")
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCode(cmd)
		state.fixture.DevRepo.GitTown.Reload()
		return nil
	})

	suite.Step(`^I run "([^"]*)" and enter "([^"]*)" for the commit message$`, func(cmd, message string) error {
		updateInitialSHAs(state)
		state.fixture.DevRepo.MockCommitMessage(message)
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCode(cmd)
		state.fixture.DevRepo.GitTown.Reload()
		return nil
	})

	suite.Step(`^I run "([^"]*)", answer the prompts, and close the next editor:$`, func(cmd string, input *messages.PickleStepArgument_PickleTable) error {
		updateInitialSHAs(state)
		env := append(os.Environ(), "GIT_EDITOR=true")
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Env: env, Input: helpers.TableToInput(input)})
		state.fixture.DevRepo.GitTown.Reload()
		return nil
	})

	suite.Step(`^I run "([^"]+)" in the "([^"]+)" folder$`, func(cmd, folderName string) error {
		updateInitialSHAs(state)
		state.runOutput, state.runExitCode = state.fixture.DevRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Dir: folderName})
		state.fixture.DevRepo.GitTown.Reload()
		return nil
	})

	suite.Step(`^inspect the repo$`, func() error {
		fmt.Printf("\nThe workspace is at %s\n", state.fixture.DevRepo.WorkingDir)
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
			return fmt.Errorf("expected text not found")
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
			return fmt.Errorf("expected exit code %d", state.runExitCode)
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
			state.initialDevSHAs,
			state.initialOriginSHAs,
		)
		diff, errorCount := table.EqualDataTable(expanded)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the commands run\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching commands run, see diff above")
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

	suite.Step(`^local Git Town setting "([^"]*)" no longer exists$`, func(name string) error {
		configKey := configdomain.ParseKey("git-town." + name)
		newValue := state.fixture.DevRepo.GitTown.LocalConfigValue(*configKey)
		if newValue == "" {
			return nil
		}
		return fmt.Errorf("should not have local %q anymore but has value %q", name, newValue)
	})

	suite.Step(`^(?:local )?Git Town setting "([^"]*)" is "([^"]*)"$`, func(name, value string) error {
		configKey := configdomain.ParseKey("git-town." + name)
		return state.fixture.DevRepo.GitTown.SetLocalConfigValue(*configKey, value)
	})

	suite.Step(`^local Git Town setting "([^"]*)" is now "([^"]*)"$`, func(name, want string) error {
		configKey := configdomain.ParseKey("git-town." + name)
		have := state.fixture.DevRepo.GitTown.LocalConfigValue(*configKey)
		if have != want {
			return fmt.Errorf("expected local setting %q to be %q, but was %q", name, want, have)
		}
		return nil
	})

	suite.Step(`^my repo does not have an origin$`, func() error {
		state.fixture.DevRepo.RemoveRemote(domain.OriginRemote)
		state.initialRemoteBranches = domain.LocalBranchNames{}
		state.fixture.OriginRepo = nil
		return nil
	})

	suite.Step(`^my repo has a Git submodule$`, func() error {
		state.fixture.AddSubmoduleRepo()
		state.fixture.DevRepo.AddSubmodule(state.fixture.SubmoduleRepo.WorkingDir)
		return nil
	})

	suite.Step(`^no lineage exists now$`, func() error {
		if state.fixture.DevRepo.GitTown.ContainsLineage() {
			branchInfo := state.fixture.DevRepo.GitTown.Lineage(state.fixture.DevRepo.GitTown.RemoveLocalConfigValue)
			return fmt.Errorf("unexpected Git Town lineage information: %+v", branchInfo)
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
		repoStatus, err := state.fixture.DevRepo.RepoStatus()
		if err != nil {
			return err
		}
		if repoStatus.RebaseInProgress {
			return fmt.Errorf("expected no rebase in progress")
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
		isOffline, err := state.fixture.DevRepo.GitTown.IsOffline()
		if err != nil {
			return err
		}
		if isOffline {
			return fmt.Errorf("expected to not be offline but am")
		}
		return nil
	})

	suite.Step(`^offline mode is enabled$`, func() error {
		return state.fixture.DevRepo.GitTown.SetOffline(true)
	})

	suite.Step(`^origin deletes the "([^"]*)" branch$`, func(branch string) error {
		slice.Remove(&state.initialRemoteBranches, domain.NewLocalBranchName(branch))
		state.fixture.OriginRepo.RemoveBranch(domain.NewLocalBranchName(branch))
		return nil
	})

	suite.Step(`^origin ships the "([^"]*)" branch$`, func(branch string) error {
		state.fixture.OriginRepo.CheckoutBranch(domain.NewLocalBranchName("main"))
		err := state.fixture.OriginRepo.MergeBranch(domain.NewLocalBranchName(branch))
		if err != nil {
			return err
		}
		state.fixture.OriginRepo.RemoveBranch(domain.NewLocalBranchName(branch))
		slice.Remove(&state.initialRemoteBranches, domain.NewLocalBranchName(branch))
		return nil
	})

	suite.Step(`^the branches "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		for _, branchName := range []string{branch1, branch2} {
			branch := domain.NewLocalBranchName(branchName)
			state.fixture.DevRepo.CreateBranch(branch, domain.NewLocalBranchName("main"))
			state.initialLocalBranches = append(state.initialLocalBranches, branch)
		}
		return nil
	})

	suite.Step(`^the branches are now$`, func(table *messages.PickleStepArgument_PickleTable) error {
		existing := state.fixture.Branches()
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
		commits := git.FromGherkinTable(table)
		state.fixture.CreateCommits(commits)
		// restore the initial branch
		if state.initialCurrentBranch.IsEmpty() {
			state.fixture.DevRepo.CheckoutBranch(domain.NewLocalBranchName("main"))
			return nil
		}
		// NOTE: reading the cached value here to keep the test suite fast by avoiding unnecessary disk access
		if state.fixture.DevRepo.CurrentBranchCache.Value() != state.initialCurrentBranch {
			state.fixture.DevRepo.CheckoutBranch(state.initialCurrentBranch)
			return nil
		}
		return nil
	})

	suite.Step(`^the coworker fetches updates$`, func() error {
		state.fixture.CoworkerRepo.Fetch()
		return nil
	})

	suite.Step(`^the coworker is on the "([^"]*)" branch$`, func(branch string) error {
		state.fixture.CoworkerRepo.CheckoutBranch(domain.NewLocalBranchName(branch))
		return nil
	})

	suite.Step(`^the coworker runs "([^"]+)"$`, func(command string) error {
		state.runOutput, state.runExitCode = state.fixture.CoworkerRepo.MustQueryStringCode(command)
		return nil
	})

	suite.Step(`^the coworker sets the parent branch of "([^"]*)" as "([^"]*)"$`, func(childBranch, parentBranch string) error {
		_ = state.fixture.CoworkerRepo.GitTown.SetParent(domain.NewLocalBranchName(childBranch), domain.NewLocalBranchName(parentBranch))
		return nil
	})

	suite.Step(`^the coworker sets the "sync-feature-strategy" to "(merge|rebase)"$`, func(value string) error {
		syncFeatureStrategy, err := configdomain.NewSyncFeatureStrategy(value)
		if err != nil {
			return err
		}
		_ = state.fixture.CoworkerRepo.GitTown.SetSyncFeatureStrategy(syncFeatureStrategy)
		return nil
	})

	suite.Step(`^the current branch is "([^"]*)"$`, func(name string) error {
		branch := domain.NewLocalBranchName(name)
		state.initialCurrentBranch = branch
		if !slice.Contains(state.initialLocalBranches, branch) {
			state.initialLocalBranches = append(state.initialLocalBranches, branch)
			state.fixture.DevRepo.CreateBranch(branch, domain.NewLocalBranchName("main"))
		}
		state.fixture.DevRepo.CheckoutBranch(branch)
		return nil
	})

	suite.Step(`^the current branch is a (local )?(feature|perennial) branch "([^"]*)"$`, func(localStr, branchType, branchName string) error {
		branch := domain.NewLocalBranchName(branchName)
		isLocal := localStr != ""
		var err error
		switch branchType {
		case "feature":
			err = state.fixture.DevRepo.CreateFeatureBranch(branch)
		case "perennial":
			state.fixture.DevRepo.CreatePerennialBranches(branch)
		default:
			panic(fmt.Sprintf("unknown branch type: %q", branchType))
		}
		if err != nil {
			return err
		}
		state.initialLocalBranches = append(state.initialLocalBranches, branch)
		if branchType == "feature" {
			state.initialLineage.AddRow(branchName, "main")
		}
		if !isLocal {
			state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
			state.fixture.DevRepo.PushBranchToRemote(branch, domain.OriginRemote)
		}
		state.initialCurrentBranch = branch
		// NOTE: reading the cached value here to keep the test suite fast by avoiding unnecessary disk access
		if !state.fixture.DevRepo.CurrentBranchCache.Initialized() || state.fixture.DevRepo.CurrentBranchCache.Value() != branch {
			state.fixture.DevRepo.CheckoutBranch(branch)
		}
		return nil
	})

	suite.Step(`^the current branch is "([^"]*)" and the previous branch is "([^"]*)"$`, func(currentText, previousText string) error {
		current := domain.NewLocalBranchName(currentText)
		previous := domain.NewLocalBranchName(previousText)
		state.initialCurrentBranch = current
		state.fixture.DevRepo.CheckoutBranch(previous)
		state.fixture.DevRepo.CheckoutBranch(current)
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

	suite.Step(`^the initial lineage exists$`, func() error {
		have := state.fixture.DevRepo.LineageTable()
		state.initialLineage.Sort()
		diff, errCnt := have.EqualDataTable(state.initialLineage)
		if errCnt > 0 {
			fmt.Printf("\nERROR! Found %d differences in the lineage\n\n", errCnt)
			fmt.Printf("INITIAL LINEAGE:\n%s\n", state.initialLineage.String())
			fmt.Printf("CURRENT LINEAGE:\n%s\n", have.String())
			fmt.Println(diff)
			return fmt.Errorf("mismatching branches found, see the diff above")
		}
		return nil
	})

	suite.Step(`^the initial branches and lineage exist$`, func() error {
		// verify initial branches
		have := state.fixture.Branches()
		want := state.InitialBranches()
		diff, errorCount := have.EqualDataTable(want)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing branches\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching branches found, see diff above")
		}
		// verify initial lineage
		state.initialLineage.Sort()
		have = state.fixture.DevRepo.LineageTable()
		diff, errCnt := have.EqualDataTable(state.initialLineage)
		if errCnt > 0 {
			fmt.Printf("\nERROR! Found %d differences in the lineage\n\n", errCnt)
			fmt.Println(diff)
			return fmt.Errorf("mismatching lineage found, see the diff above")
		}
		return nil
	})

	suite.Step(`^the initial branches exist$`, func() error {
		have := state.fixture.Branches()
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

	suite.Step(`^the initial commits exist$`, func() error {
		return state.compareTable(state.initialCommits)
	})

	suite.Step(`^the (local )?feature branches "([^"]+)" and "([^"]+)"$`, func(localStr, branch1, branch2 string) error {
		isLocal := localStr != ""
		for _, branchText := range []string{branch1, branch2} {
			branch := domain.NewLocalBranchName(branchText)
			err := state.fixture.DevRepo.CreateFeatureBranch(branch)
			if err != nil {
				return err
			}
			state.initialLocalBranches = append(state.initialLocalBranches, branch)
			state.initialLineage.AddRow(branchText, "main")
			if !isLocal {
				state.fixture.DevRepo.PushBranchToRemote(branch, domain.OriginRemote)
				state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
			}
		}
		return nil
	})

	suite.Step(`^the (local )?feature branches "([^"]+)", "([^"]+)", and "([^"]+)"$`, func(localStr, branch1, branch2, branch3 string) error {
		isLocal := localStr != ""
		for _, branchText := range []string{branch1, branch2, branch3} {
			branch := domain.NewLocalBranchName(branchText)
			err := state.fixture.DevRepo.CreateFeatureBranch(branch)
			if err != nil {
				return err
			}
			state.initialLocalBranches = append(state.initialLocalBranches, branch)
			state.initialLineage.AddRow(branchText, "main")
			if !isLocal {
				state.fixture.DevRepo.PushBranchToRemote(branch, domain.OriginRemote)
				state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
			}
		}
		return nil
	})

	suite.Step(`^the (local )?perennial branches "([^"]+)" and "([^"]+)"$`, func(localStr, branch1Text, branch2Text string) error {
		branch1 := domain.NewLocalBranchName(branch1Text)
		branch2 := domain.NewLocalBranchName(branch2Text)
		isLocal := localStr != ""
		state.fixture.DevRepo.CreatePerennialBranches(branch1, branch2)
		state.initialLocalBranches = append(state.initialLocalBranches, branch1, branch2)
		if !isLocal {
			state.initialRemoteBranches = append(state.initialRemoteBranches, branch1, branch2)
			state.fixture.DevRepo.PushBranchToRemote(branch1, domain.OriginRemote)
			state.fixture.DevRepo.PushBranchToRemote(branch2, domain.OriginRemote)
		}
		return nil
	})

	suite.Step(`^the (local )?perennial branches "([^"]+)", "([^"]+)", and "([^"]+)"$`, func(localStr, branch1, branch2, branch3 string) error {
		isLocal := localStr != ""
		for _, branchText := range []string{branch1, branch2, branch3} {
			branch := domain.NewLocalBranchName(branchText)
			state.fixture.DevRepo.CreatePerennialBranches(branch)
			state.initialLocalBranches = append(state.initialLocalBranches, branch)
			if !isLocal {
				state.fixture.DevRepo.PushBranchToRemote(branch, domain.OriginRemote)
				state.initialRemoteBranches = append(state.initialRemoteBranches, branch)
			}
		}
		return nil
	})

	suite.Step(`^the main branch is "([^"]+)"$`, func(name string) error {
		return state.fixture.DevRepo.GitTown.SetMainBranch(domain.NewLocalBranchName(name))
	})

	suite.Step(`^the main branch is now "([^"]+)"$`, func(want string) error {
		have := state.fixture.DevRepo.GitTown.MainBranch()
		if have.String() != want {
			return fmt.Errorf("expected %q, got %q", want, have)
		}
		return nil
	})

	suite.Step(`^the origin is "([^"]*)"$`, func(origin string) error {
		state.fixture.DevRepo.SetTestOrigin(origin)
		return nil
	})

	suite.Step(`^the perennial branches are "([^"]+)"$`, func(name string) error {
		return state.fixture.DevRepo.GitTown.AddToPerennialBranches(domain.NewLocalBranchName(name))
	})

	suite.Step(`^the perennial branches are "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		return state.fixture.DevRepo.GitTown.AddToPerennialBranches(domain.NewLocalBranchName(branch1), domain.NewLocalBranchName(branch2))
	})

	suite.Step(`^the perennial branches are not configured$`, func() error {
		return state.fixture.DevRepo.RemovePerennialBranchConfiguration()
	})

	suite.Step(`^the perennial branches are now "([^"]+)"$`, func(name string) error {
		actual := state.fixture.DevRepo.GitTown.PerennialBranches()
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 perennial branch, got %q", actual)
		}
		if actual[0].String() != name {
			return fmt.Errorf("expected %q, got %q", name, actual[0])
		}
		return nil
	})

	suite.Step(`^the perennial branches are now "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		actual := state.fixture.DevRepo.GitTown.PerennialBranches()
		if len(actual) != 2 {
			return fmt.Errorf("expected 2 perennial branches, got %q", actual)
		}
		if actual[0].String() != branch1 || actual[1].String() != branch2 {
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

	suite.Step(`^there are still no perennial branches$`, func() error {
		branches := state.fixture.DevRepo.GitTown.PerennialBranches()
		if len(branches) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", branches)
		}
		return nil
	})

	suite.Step(`^these committed files exist now$`, func(table *messages.PickleStepArgument_PickleTable) error {
		fileTable := state.fixture.DevRepo.FilesInBranches(domain.NewLocalBranchName("main"))
		diff, errorCount := fileTable.EqualGherkin(table)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing files\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching files found, see diff above")
		}
		return nil
	})

	suite.Step(`^these commits exist now$`, func(table *messages.PickleStepArgument_PickleTable) error {
		return state.compareTable(table)
	})

	suite.Step(`^these tags exist$`, func(table *messages.PickleStepArgument_PickleTable) error {
		tagTable := state.fixture.TagTable()
		diff, errorCount := tagTable.EqualGherkin(table)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing tags\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching tags found, see diff above")
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
		stashSnapshot, err := state.fixture.DevRepo.StashSnapshot()
		if err != nil {
			return err
		}
		if stashSnapshot != 1 {
			return fmt.Errorf("expected 1 stash but found %d", stashSnapshot)
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

	suite.Step(`^this branch lineage exists now$`, func(input *messages.PickleStepArgument_PickleTable) error {
		table := state.fixture.DevRepo.LineageTable()
		diff, errCount := table.EqualGherkin(input)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the lineage\n\n", errCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching branches found, see the diff above")
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
}
