package test

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"sync"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"github.com/git-town/git-town/src/command"
)

// beforeSuiteMux ensures that we run BeforeSuite only once globally.
var beforeSuiteMux sync.Mutex

// the global GitManager instance
var gitManager *GitManager

// Steps defines Cucumber step implementations around Git workspace management.
// nolint: gocyclo,gocognit,funlen
func Steps(suite *godog.Suite, state *ScenarioState) {
	suite.BeforeScenario(func(scenario *messages.Pickle) {
		// create a GitEnvironment for the scenario
		gitEnvironment, err := gitManager.CreateScenarioEnvironment(scenario.GetName())
		if err != nil {
			log.Fatalf("cannot create environment for scenario %q: %s", scenario.GetName(), err)
		}
		// Godog only provides state for the entire feature.
		// We want state to be scenario-specific, hence we reset the shared state before each scenario.
		// This is a limitation of the current Godog implementation, which doesn't have a `ScenarioContext` method,
		// only a `FeatureContext` method.
		// See main_test.go for additional details.
		state.Reset(gitEnvironment)
		if hasTag(scenario, "@debug") {
			Debug = true
		}
	})

	suite.BeforeSuite(func() {
		// NOTE: we want to create only one global GitManager instance with one global memoized environment.
		beforeSuiteMux.Lock()
		defer beforeSuiteMux.Unlock()
		if gitManager == nil {
			baseDir, err := ioutil.TempDir("", "")
			if err != nil {
				log.Fatalf("cannot create base directory for feature specs: %s", err)
			}
			// Evaluate symlinks as Mac temp dir is symlinked
			evalBaseDir, err := filepath.EvalSymlinks(baseDir)
			if err != nil {
				log.Fatalf("cannot evaluate symlinks of base directory for feature specs: %s", err)
			}
			gitManager = NewGitManager(evalBaseDir)
			err = gitManager.CreateMemoizedEnvironment()
			if err != nil {
				log.Fatalf("Cannot create memoized environment: %s", err)
			}
		}
	})

	suite.AfterScenario(func(scenario *messages.Pickle, e error) {
		if e != nil {
			fmt.Printf("failed scenario, investigate state in %q\n", state.gitEnv.Dir)
		}
	})
	suite.Step(`^all branches are now synchronized$`, func() error {
		outOfSync, err := state.gitEnv.DevRepo.HasBranchesOutOfSync()
		if err != nil {
			return err
		}
		if outOfSync {
			return fmt.Errorf("expected no branches out of sync")
		}
		return nil
	})

	suite.Step(`^Git Town is in offline mode$`, func() error {
		state.gitEnv.DevRepo.SetOffline(true)
		return nil
	})

	suite.Step(`^Git Town is no longer configured for this repo$`, func() error {
		res, err := state.gitEnv.DevRepo.HasGitTownConfigNow()
		if err != nil {
			return err
		}
		if res {
			return fmt.Errorf("unexpected Git Town configuration")
		}
		return nil
	})

	suite.Step(`^Git Town is (?:now|still) aware of this branch hierarchy$`, func(input *messages.PickleStepArgument_PickleTable) error {
		table := DataTable{}
		table.AddRow("BRANCH", "PARENT")
		for _, row := range input.Rows[1:] {
			branch := row.Cells[0].Value
			state.gitEnv.DevRepo.Configuration.Reload()
			parentBranch := state.gitEnv.DevRepo.Configuration.GetParentBranch(branch)
			table.AddRow(branch, parentBranch)
		}
		diff, errCount := table.EqualGherkin(input)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the branch hierarchy\n\n", errCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching branches found, see the diff above")
		}
		return nil
	})

	suite.Step(`^Git Town now has no branch hierarchy information$`, func() error {
		state.gitEnv.DevRepo.Configuration.Reload()
		if state.gitEnv.DevRepo.Configuration.HasBranchInformation() {
			return fmt.Errorf("unexpected Git Town branch hierarchy information")
		}
		return nil
	})

	suite.Step(`^I am collaborating with a coworker$`, func() error {
		return state.gitEnv.AddCoworkerRepo()
	})

	suite.Step(`^I am in the project root folder$`, func() error {
		actual, err := state.gitEnv.DevRepo.LastActiveDir()
		if err != nil {
			return fmt.Errorf("cannot determine the current working directory: %w", err)
		}
		expected := state.gitEnv.DevRepo.WorkingDir()
		if actual != expected {
			return fmt.Errorf("expected to be in %q but am in %q", expected, actual)
		}
		return nil
	})

	suite.Step(`^I am not prompted for any parent branches$`, func() error {
		notExpected := "Please specify the parent branch of"
		if state.runRes.OutputContainsText(notExpected) {
			return fmt.Errorf("text found:\n\nDID NOT EXPECT: %q\n\nACTUAL\n\n%q\n----------------------------", notExpected, state.runRes.Output())
		}
		return nil
	})

	suite.Step(`^I am on the "([^"]*)" branch$`, func(branchName string) error {
		err := state.gitEnv.DevRepo.CheckoutBranch(branchName)
		if err != nil {
			return fmt.Errorf("cannot change to branch %q: %w", branchName, err)
		}
		return nil
	})

	suite.Step(`^I am on the "([^"]*)" branch with "([^"]*)" as the previous Git branch$`, func(current, previous string) error {
		err := state.gitEnv.DevRepo.CheckoutBranch(previous)
		if err != nil {
			return err
		}
		return state.gitEnv.DevRepo.CheckoutBranch(current)
	})

	suite.Step(`^I (?:end up|am still) on the "([^"]*)" branch$`, func(expected string) error {
		actual, err := state.gitEnv.DevRepo.CurrentBranch()
		if err != nil {
			return fmt.Errorf("cannot determine current branch of developer repo: %w", err)
		}
		if actual != expected {
			return fmt.Errorf("expected active branch %q but is %q", expected, actual)
		}
		return nil
	})

	suite.Step(`^I haven't configured Git Town yet$`, func() error {
		state.gitEnv.DevRepo.DeletePerennialBranchConfiguration()
		return state.gitEnv.DevRepo.DeleteMainBranchConfiguration()
	})

	suite.Step(`^I resolve the conflict in "([^"]*)"(?: with "([^"]*)")?$`, func(filename, content string) error {
		if content == "" {
			content = "resolved content"
		}
		err := state.gitEnv.DevRepo.CreateFile(filename, content)
		if err != nil {
			return err
		}
		err = state.gitEnv.DevRepo.StageFiles(filename)
		if err != nil {
			return err
		}
		return nil
	})

	suite.Step(`^I run "(.+)"$`, func(command string) error {
		state.runRes, state.runErr = state.gitEnv.DevShell.RunString(command)
		return nil
	})

	suite.Step(`^I run "([^"]+)" and answer the prompts:$`, func(cmd string, input *messages.PickleStepArgument_PickleTable) error {
		state.runRes, state.runErr = state.gitEnv.DevShell.RunStringWith(cmd, command.Options{Input: tableToInput(input)})
		return nil
	})

	suite.Step(`^I run "([^"]*)", answer the prompts, and close the next editor:$`, func(cmd string, input *messages.PickleStepArgument_PickleTable) error {
		env := append(os.Environ(), "GIT_EDITOR=true")
		state.runRes, state.runErr = state.gitEnv.DevShell.RunStringWith(cmd, command.Options{Env: env, Input: tableToInput(input)})
		return nil
	})

	suite.Step(`^I run "([^"]*)" and close the editor$`, func(cmd string) error {
		env := append(os.Environ(), "GIT_EDITOR=true")
		state.runRes, state.runErr = state.gitEnv.DevShell.RunStringWith(cmd, command.Options{Env: env})
		return nil
	})

	suite.Step(`^I run "([^"]*)" and enter an empty commit message$`, func(cmd string) error {
		state.runRes, state.runErr = state.gitEnv.DevShell.RunStringWith(cmd, command.Options{Input: []string{"dGZZ"}})
		return nil
	})

	suite.Step(`^I run "([^"]+)" in the "([^"]+)" folder$`, func(cmd, folderName string) error {
		state.runRes, state.runErr = state.gitEnv.DevShell.RunStringWith(cmd, command.Options{Dir: folderName})
		return nil
	})

	suite.Step(`^inspect the repo$`, func() error {
		fmt.Println(state.gitEnv.DevRepo.WorkingDir())
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		return nil
	})

	suite.Step(`^it does not print "([^\"]*)"$`, func(text string) error {
		if strings.Contains(state.runRes.OutputSanitized(), text) {
			return fmt.Errorf("text found: %q", text)
		}
		return nil
	})

	suite.Step(`^it prints:$`, func(expected *messages.PickleStepArgument_PickleDocString) error {
		if !strings.Contains(state.runRes.OutputSanitized(), expected.Content) {
			return fmt.Errorf("text not found:\n\nEXPECTED: %q\n\nACTUAL:\n\n%q", expected.Content, state.runRes.OutputSanitized())
		}
		return nil
	})

	suite.Step(`^it prints no output$`, func() error {
		output := state.runRes.OutputSanitized()
		if output != "" {
			return fmt.Errorf("expected no output but found %q", output)
		}
		return nil
	})

	suite.Step(`^it prints the error:$`, func(expected *messages.PickleStepArgument_PickleDocString) error {
		if !strings.Contains(state.runRes.OutputSanitized(), expected.Content) {
			return fmt.Errorf("text not found: %s\n\nactual text:\n%s", expected.Content, state.runRes.OutputSanitized())
		}
		if state.runErr == nil {
			return fmt.Errorf("expected error")
		}
		return nil
	})

	suite.Step(`^it prints the initial configuration prompt$`, func() error {
		expected := "Git Town needs to be configured"
		if !state.runRes.OutputContainsText(expected) {
			return fmt.Errorf("text not found:\n\nEXPECTED: %q\n\nACTUAL:\n\n%q\n----------------------------", expected, state.runRes.Output())
		}
		return nil
	})

	suite.Step(`^it runs no commands$`, func() error {
		commands := GitCommandsInGitTownOutput(state.runRes.Output())
		if len(commands) > 0 {
			for _, command := range commands {
				fmt.Println(command)
			}
			return fmt.Errorf("expected no commands but found %d commands", len(commands))
		}
		return nil
	})

	suite.Step(`^it runs the commands$`, func(input *messages.PickleStepArgument_PickleTable) error {
		commands := GitCommandsInGitTownOutput(state.runRes.Output())
		table := RenderExecutedGitCommands(commands, input)
		dataTable := FromGherkin(input)
		expanded := dataTable.Expand(
			state.gitEnv.DevRepo.WorkingDir(),
			&state.gitEnv.DevRepo,
			state.gitEnv.OriginRepo,
		)
		diff, errorCount := table.EqualDataTable(expanded)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the commands run\n\n", errorCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching commands run, see diff above")
		}
		return nil
	})

	suite.Step(`^"([^"]*)" launches a new pull request with this url in my browser:$`, func(tool string, url *messages.PickleStepArgument_PickleDocString) error {
		want := fmt.Sprintf("%s called with: %s", tool, url.Content)
		want = strings.ReplaceAll(want, "?", `\?`)
		regex := regexp.MustCompile(want)
		have := state.runRes.OutputSanitized()
		if !regex.MatchString(have) {
			return fmt.Errorf("EXPECTED: a regex matching %q\nGOT: %q", want, have)
		}
		return nil
	})

	suite.Step(`^my code base has a feature branch named "([^"]*)"$`, func(name string) error {
		err := state.gitEnv.DevRepo.CreateFeatureBranch(name)
		if err != nil {
			return err
		}
		return state.gitEnv.DevRepo.PushBranchSetUpstream(name)
	})

	suite.Step(`^my code base has a feature branch named "([^"]*)" as a child of "([^"]*)"$`, func(branch, parent string) error {
		err := state.gitEnv.DevRepo.CreateChildFeatureBranch(branch, parent)
		if err != nil {
			return err
		}
		return state.gitEnv.DevRepo.PushBranchSetUpstream(branch)
	})

	suite.Step(`^my computer has a broken "([^"]*)" tool installed$`, func(name string) error {
		return state.gitEnv.DevShell.MockBrokenCommand(name)
	})

	suite.Step(`^my computer has an empty fish autocompletion folder$`, func() error {
		return os.MkdirAll(fishFolderPath(state), 0744)
	})

	suite.Step(`^my computer has an existing Git autocompletion file$`, func() error {
		err := os.MkdirAll(fishFolderPath(state), 0744)
		if err != nil {
			return fmt.Errorf("cannot create fish folder: %w", err)
		}
		return ioutil.WriteFile(fishFilePath(state), []byte("existing content"), 0744)
	})

	suite.Step(`^my computer has Git "([^"]*)" installed$`, func(version string) error {
		err := state.gitEnv.DevShell.MockGit(version)
		return err
	})

	suite.Step(`^my computer has no fish autocompletion file$`, func() error {
		// nothing to do here, the test directory has no data
		return nil
	})

	suite.Step(`^my computer has no tool to open browsers installed$`, func() error {
		return state.gitEnv.DevShell.MockNoCommandsInstalled()
	})

	suite.Step(`^my computer has the "([^"]*)" tool installed$`, func(tool string) error {
		return state.gitEnv.DevShell.MockCommand(tool)
	})

	suite.Step(`^my computer now has a Git autocompletion file$`, func() error {
		fishPath := filepath.Join(state.gitEnv.Dir, ".config", "fish", "completions", "git.fish")
		_, err := os.Stat(fishPath)
		if os.IsNotExist(err) {
			return err
		}
		return nil
	})

	suite.Step(`^my computer still has the original Git autocompletion file$`, func() error {
		content, err := ioutil.ReadFile(fishFilePath(state))
		if err != nil {
			return err
		}
		contentStr := string(content)
		if contentStr != "existing content" {
			return fmt.Errorf("config file content was changed to %q", content)
		}
		return nil
	})

	suite.Step(`^my (?:coworker|origin) has a feature branch named "([^"]*)"$`, func(branch string) error {
		return state.gitEnv.OriginRepo.CreateBranch(branch, "main")
	})

	suite.Step(`^my coworker fetches updates$`, func() error {
		return state.gitEnv.CoworkerRepo.Fetch()
	})

	suite.Step(`^my coworker is on the "([^"]*)" branch$`, func(branchName string) error {
		return state.gitEnv.CoworkerRepo.CheckoutBranch(branchName)
	})

	suite.Step(`^my coworker runs "([^"]+)"$`, func(command string) error {
		state.runRes, state.runErr = state.gitEnv.CoworkerRepo.RunString(command)
		return nil
	})

	suite.Step(`^my coworker sets the parent branch of "([^"]*)" as "([^"]*)"$`, func(childBranch, parentBranch string) error {
		_ = state.gitEnv.CoworkerRepo.Configuration.SetParentBranch(childBranch, parentBranch)
		return nil
	})

	suite.Step(`^my repo does not have a remote origin$`, func() error {
		err := state.gitEnv.DevRepo.RemoveRemote("origin")
		if err != nil {
			return err
		}
		state.gitEnv.OriginRepo = nil
		return nil
	})

	suite.Step(`^my repo doesn't have a main branch configured$`, func() error {
		return state.gitEnv.DevRepo.DeleteMainBranchConfiguration()
	})

	suite.Step(`^my repo doesn't have any uncommitted files$`, func() error {
		files, err := state.gitEnv.DevRepo.UncommittedFiles()
		if err != nil {
			return fmt.Errorf("cannot determine uncommitted files: %w", err)
		}
		if len(files) > 0 {
			return fmt.Errorf("unexpected uncommitted files: %s", files)
		}
		return nil
	})

	suite.Step(`^my repo has a branch "([^"]*)"$`, func(branch string) error {
		return state.gitEnv.DevRepo.CreateBranch(branch, "main")
	})

	suite.Step(`^my repo has a feature branch named "([^"]*)" with no parent$`, func(branch string) error {
		return state.gitEnv.DevRepo.CreateFeatureBranchNoParent(branch)
	})

	suite.Step(`^my repo has a feature branch named "([^"]+)" as a child of "([^"]+)"$`, func(childBranch, parentBranch string) error {
		err := state.gitEnv.DevRepo.CreateChildFeatureBranch(childBranch, parentBranch)
		if err != nil {
			return fmt.Errorf("cannot create feature branch %q: %w", childBranch, err)
		}
		return state.gitEnv.DevRepo.PushBranchSetUpstream(childBranch)
	})

	suite.Step(`^my repo has a (local )?feature branch named "([^"]*)"$`, func(localStr, branch string) error {
		isLocal := localStr != ""
		err := state.gitEnv.DevRepo.CreateFeatureBranch(branch)
		if err != nil {
			return err
		}
		if !isLocal {
			return state.gitEnv.DevRepo.PushBranchSetUpstream(branch)
		}
		return nil
	})

	suite.Step(`^my repo has an upstream repo$`, func() error {
		return state.gitEnv.AddUpstream()
	})

	suite.Step(`^my repo has "color\.ui" set to "([^"]*)"$`, func(value string) error {
		_ = state.gitEnv.DevRepo.SetColorUI(value)
		return nil
	})

	suite.Step(`^my repo has "git-town.code-hosting-driver" set to "([^"]*)"$`, func(value string) error {
		_ = state.gitEnv.DevRepo.SetCodeHostingDriver(value)
		return nil
	})

	suite.Step(`^my repo has "git-town.code-hosting-origin-hostname" set to "([^"]*)"$`, func(value string) error {
		_ = state.gitEnv.DevRepo.SetCodeHostingOriginHostname(value)
		return nil
	})

	suite.Step(`^my repo has "git-town.ship-delete-remote-branch" set to "(true|false)"$`, func(value string) error {
		parsed, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		_ = state.gitEnv.DevRepo.SetShouldShipDeleteRemoteBranch(parsed)
		return nil
	})

	suite.Step(`^my repo has "git-town.sync-upstream" set to (true|false)$`, func(text string) error {
		value, err := strconv.ParseBool(text)
		if err != nil {
			return err
		}
		_ = state.gitEnv.DevRepo.SetShouldSyncUpstream(value)
		return nil
	})

	suite.Step(`^my repo has the branches "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		err := state.gitEnv.DevRepo.CreateBranch(branch1, "main")
		if err != nil {
			return err
		}
		return state.gitEnv.DevRepo.CreateBranch(branch2, "main")
	})

	suite.Step(`^my repo has the following tags$`, func(table *messages.PickleStepArgument_PickleTable) error {
		return state.gitEnv.CreateTags(table)
	})

	suite.Step(`^my repo has the (local )?feature branches "([^"]+)" and "([^"]+)"$`, func(localStr, branch1, branch2 string) error {
		isLocal := localStr != ""
		err := state.gitEnv.DevRepo.CreateFeatureBranch(branch1)
		if err != nil {
			return err
		}
		err = state.gitEnv.DevRepo.CreateFeatureBranch(branch2)
		if err != nil {
			return err
		}
		if !isLocal {
			err = state.gitEnv.DevRepo.PushBranchSetUpstream(branch1)
			if err != nil {
				return err
			}
			return state.gitEnv.DevRepo.PushBranchSetUpstream(branch2)
		}
		return nil
	})

	suite.Step(`^my repo has the (local )?perennial branches "([^"]+)" and "([^"]+)"$`, func(localStr, branch1, branch2 string) error {
		isLocal := localStr != ""
		err := state.gitEnv.DevRepo.CreatePerennialBranches(branch1, branch2)
		if err != nil {
			return fmt.Errorf("cannot create perennial branches: %w", err)
		}
		if !isLocal {
			err = state.gitEnv.DevRepo.PushBranchSetUpstream(branch1)
			if err != nil {
				return err
			}
			return state.gitEnv.DevRepo.PushBranchSetUpstream(branch2)
		}
		return nil
	})

	suite.Step(`^my repo has the perennial branch "([^"]+)"`, func(branch1 string) error {
		err := state.gitEnv.DevRepo.CreatePerennialBranches(branch1)
		if err != nil {
			return fmt.Errorf("cannot create perennial branches: %w", err)
		}
		return state.gitEnv.DevRepo.PushBranchSetUpstream(branch1)
	})

	suite.Step(`^my repo is left with my original commits$`, func() error {
		return compareExistingCommits(state, state.initialCommits)
	})

	suite.Step(`^my repo is now configured with no perennial branches$`, func() error {
		state.gitEnv.DevRepo.Configuration.Reload()
		branches := state.gitEnv.DevRepo.GetPerennialBranches()
		if len(branches) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", branches)
		}
		return nil
	})

	suite.Step(`^my repo knows about the remote branch$`, func() error {
		return state.gitEnv.DevRepo.Fetch()
	})

	suite.Step(`^my repo now has the following commits$`, func(table *messages.PickleStepArgument_PickleTable) error {
		return compareExistingCommits(state, table)
	})

	suite.Step(`^my repo now has the following tags$`, func(table *messages.PickleStepArgument_PickleTable) error {
		tagTable, err := state.gitEnv.TagTable()
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

	suite.Step(`^my repo (?:now|still) has a merge in progress$`, func() error {
		hasMerge, err := state.gitEnv.DevRepo.HasMergeInProgress()
		if err != nil {
			return err
		}
		if !hasMerge {
			return fmt.Errorf("expected merge in progress")
		}
		return nil
	})

	suite.Step(`^my repo (?:now|still) has a rebase in progress$`, func() error {
		hasRebase, err := state.gitEnv.DevRepo.HasRebaseInProgress()
		if err != nil {
			return err
		}
		if !hasRebase {
			return fmt.Errorf("expected rebase in progress")
		}
		return nil
	})

	suite.Step(`^my repo (?:now|still) has the following committed files$`, func(table *messages.PickleStepArgument_PickleTable) error {
		fileTable, err := state.gitEnv.DevRepo.FilesInBranches()
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

	suite.Step(`^my repo has a remote tag "([^"]+)" that is not on a branch$`, func(name string) error {
		return state.gitEnv.OriginRepo.CreateStandaloneTag(name)
	})

	suite.Step(`^my repo's origin is "([^"]*)"$`, func(origin string) error {
		state.gitEnv.DevShell.SetTestOrigin(origin)
		return nil
	})

	suite.Step(`^my uncommitted file is stashed$`, func() error {
		uncommittedFiles, err := state.gitEnv.DevRepo.UncommittedFiles()
		if err != nil {
			return err
		}
		for ucf := range uncommittedFiles {
			if uncommittedFiles[ucf] == state.uncommittedFileName {
				return fmt.Errorf("expected file %q to be stashed but it is still uncommitted", state.uncommittedFileName)
			}
		}
		stashSize, err := state.gitEnv.DevRepo.StashSize()
		if err != nil {
			return err
		}
		if stashSize != 1 {
			return fmt.Errorf("expected 1 stash but found %d", stashSize)
		}
		return nil
	})

	suite.Step(`^my workspace has an uncommitted file$`, func() error {
		state.uncommittedFileName = "uncommitted file"
		state.uncommittedContent = "uncommitted content"
		return state.gitEnv.DevRepo.CreateFile(
			state.uncommittedFileName,
			state.uncommittedContent,
		)
	})

	suite.Step(`^my workspace has an uncommitted file in folder "([^"]*)"$`, func(folder string) error {
		state.uncommittedFileName = fmt.Sprintf("%s/uncommitted file", folder)
		return state.gitEnv.DevRepo.CreateFile(
			state.uncommittedFileName,
			state.uncommittedContent,
		)
	})

	suite.Step(`^my workspace has an uncommitted file with name: "([^"]+)" and content: "([^"]+)"$`, func(name, content string) error {
		state.uncommittedFileName = name
		state.uncommittedContent = content
		return state.gitEnv.DevRepo.CreateFile(name, content)
	})

	suite.Step(`^my workspace has the uncommitted file again$`, func() error {
		hasFile, err := state.gitEnv.DevRepo.HasFile(
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

	suite.Step(`^my workspace is currently not a Git repo$`, func() error {
		os.RemoveAll(filepath.Join(state.gitEnv.DevRepo.WorkingDir(), ".git"))
		return nil
	})

	suite.Step(`^my workspace still contains my uncommitted file$`, func() error {
		hasFile, err := state.gitEnv.DevRepo.HasFile(
			state.uncommittedFileName,
			state.uncommittedContent,
		)
		if err != nil {
			return fmt.Errorf("cannot determine if workspace contains uncommitted file: %w", err)
		}
		if !hasFile {
			return fmt.Errorf("expected the uncommitted file but didn't find one")
		}
		return nil
	})

	suite.Step(`^my workspace still contains the file "([^"]*)" with content "([^"]*)"$`, func(file, expectedContent string) error {
		actualContent, err := state.gitEnv.DevRepo.FileContent(file)
		if err != nil {
			return err
		}
		if expectedContent != actualContent {
			return fmt.Errorf("file content does not match\n\nEXPECTED: %q\n\nACTUAL:\n\n%q\n----------------------------", expectedContent, actualContent)
		}
		return nil
	})

	suite.Step(`^offline mode is disabled$`, func() error {
		state.gitEnv.DevRepo.Configuration.Reload()
		if state.gitEnv.DevRepo.IsOffline() {
			return fmt.Errorf("expected to not be offline but am")
		}
		return nil
	})

	suite.Step(`^offline mode is enabled$`, func() error {
		state.gitEnv.DevRepo.Configuration.Reload()
		if !state.gitEnv.DevRepo.IsOffline() {
			return fmt.Errorf("expected to be offline but am not")
		}
		return nil
	})

	suite.Step(`^the "([^"]*)" branch gets deleted on the remote$`, func(name string) error {
		return state.gitEnv.OriginRepo.RemoveBranch(name)
	})

	suite.Step(`^the following commits exist in my repo$`, func(table *messages.PickleStepArgument_PickleTable) error {
		state.initialCommits = table
		commits, err := FromGherkinTable(table)
		if err != nil {
			return fmt.Errorf("cannot parse Gherkin table: %w", err)
		}
		return state.gitEnv.CreateCommits(commits)
	})

	suite.Step(`^the existing branches are$`, func(table *messages.PickleStepArgument_PickleTable) error {
		existing, err := state.gitEnv.Branches()
		if err != nil {
			return err
		}
		// remove the master branch from the remote since it exists only as a performance optimization
		existing.RemoveText("master, ")
		existing.RemoveText(", master")
		diff, errCount := existing.EqualGherkin(table)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the branches\n\n", errCount)
			fmt.Println(diff)
			return fmt.Errorf("mismatching branches found, see the diff above")
		}
		return nil
	})

	suite.Step(`^the global new-branch-push-flag configuration is (true|false)$`, func(text string) error {
		b, err := strconv.ParseBool(text)
		if err != nil {
			return err
		}
		_ = state.gitEnv.DevRepo.SetNewBranchPush(b, true)
		return nil
	})

	suite.Step(`^the main branch is configured as "([^"]+)"$`, func(name string) error {
		state.gitEnv.DevRepo.SetMainBranch(name)
		return nil
	})

	suite.Step(`^the main branch name is not configured$`, func() error {
		return state.gitEnv.DevRepo.DeleteMainBranchConfiguration()
	})

	suite.Step(`^the main branch is now configured as "([^"]+)"$`, func(name string) error {
		state.gitEnv.DevRepo.Configuration.Reload()
		actual := state.gitEnv.DevRepo.GetMainBranch()
		if actual != name {
			return fmt.Errorf("expected %q, got %q", name, actual)
		}
		return nil
	})

	suite.Step(`^the new-branch-push-flag configuration is (true|false)$`, func(value string) error {
		b, err := strconv.ParseBool(value)
		if err != nil {
			return err
		}
		state.gitEnv.DevRepo.SetNewBranchPush(b, false)
		return nil
	})

	suite.Step(`^the new-branch-push-flag configuration is now (true|false)$`, func(text string) error {
		want, err := strconv.ParseBool(text)
		if err != nil {
			return err
		}
		state.gitEnv.DevRepo.Configuration.Reload()
		have := state.gitEnv.DevRepo.ShouldNewBranchPush()
		if have != want {
			return fmt.Errorf("expected global new-branch-push-flag to be %t, but was %t", want, have)
		}
		return nil
	})

	suite.Step(`^the perennial branches are configured as "([^"]+)"$`, func(name string) error {
		state.gitEnv.DevRepo.AddToPerennialBranches(name)
		return nil
	})

	suite.Step(`^the perennial branches are configured as "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		state.gitEnv.DevRepo.AddToPerennialBranches(branch1, branch2)
		return nil
	})

	suite.Step(`^the perennial branches are now configured as "([^"]+)"$`, func(name string) error {
		state.gitEnv.DevRepo.Configuration.Reload()
		actual := state.gitEnv.DevRepo.GetPerennialBranches()
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 perennial branch, got %q", actual)
		}
		if actual[0] != name {
			return fmt.Errorf("expected %q, got %q", name, actual[0])
		}
		return nil
	})

	suite.Step(`^the perennial branches are now configured as "([^"]+)" and "([^"]+)"$`, func(branch1, branch2 string) error {
		state.gitEnv.DevRepo.Configuration.Reload()
		actual := state.gitEnv.DevRepo.GetPerennialBranches()
		if len(actual) != 2 {
			return fmt.Errorf("expected 2 perennial branches, got %q", actual)
		}
		if actual[0] != branch1 || actual[1] != branch2 {
			return fmt.Errorf("expected %q, got %q", []string{branch1, branch2}, actual)
		}
		return nil
	})

	suite.Step(`^the perennial branches are not configured$`, func() error {
		state.gitEnv.DevRepo.DeletePerennialBranchConfiguration()
		return nil
	})

	suite.Step(`^the previous Git branch is (?:now|still) "([^"]*)"$`, func(want string) error {
		err := state.gitEnv.DevRepo.CheckoutBranch("-")
		if err != nil {
			return err
		}
		have, err := state.gitEnv.DevRepo.CurrentBranch()
		if err != nil {
			return err
		}
		if have != want {
			return fmt.Errorf("expected previous branch %q but got %q", want, have)
		}
		return state.gitEnv.DevRepo.CheckoutBranch("-")
	})

	suite.Step(`^the pull-branch-strategy configuration is "(merge|rebase)"$`, func(value string) error {
		state.gitEnv.DevRepo.SetPullBranchStrategy(value)
		return nil
	})

	suite.Step(`^the pull-branch-strategy configuration is now "(merge|rebase)"$`, func(want string) error {
		state.gitEnv.DevRepo.Configuration.Reload()
		have := state.gitEnv.DevRepo.GetPullBranchStrategy()
		if have != want {
			return fmt.Errorf("expected pull-branch-strategy to be %q but was %q", want, have)
		}
		return nil
	})

	suite.Step(`^the remote deletes the "([^"]*)" branch$`, func(name string) error {
		return state.gitEnv.OriginRepo.RemoveBranch(name)
	})

	suite.Step(`^there is no merge in progress$`, func() error {
		hasMerge, err := state.gitEnv.DevRepo.HasMergeInProgress()
		if err != nil {
			return err
		}
		if hasMerge {
			return fmt.Errorf("expected no merge in progress")
		}
		return nil
	})

	suite.Step(`^there is no rebase in progress$`, func() error {
		hasRebase, err := state.gitEnv.DevRepo.HasRebaseInProgress()
		if err != nil {
			return err
		}
		if hasRebase {
			return fmt.Errorf("expected no rebase in progress")
		}
		return nil
	})
}
