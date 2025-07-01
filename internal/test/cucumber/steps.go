package cucumber

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/cucumber/godog"
	messages "github.com/cucumber/messages/go/v21"
	"github.com/git-town/git-town/v21/internal/browser"
	"github.com/git-town/git-town/v21/internal/cli/dialog/components"
	"github.com/git-town/git-town/v21/internal/cli/print"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/configfile"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/test/commands"
	"github.com/git-town/git-town/v21/internal/test/datatable"
	"github.com/git-town/git-town/v21/internal/test/envvars"
	"github.com/git-town/git-town/v21/internal/test/filesystem"
	"github.com/git-town/git-town/v21/internal/test/fixture"
	"github.com/git-town/git-town/v21/internal/test/handlebars"
	"github.com/git-town/git-town/v21/internal/test/helpers"
	"github.com/git-town/git-town/v21/internal/test/output"
	"github.com/git-town/git-town/v21/internal/test/subshell"
	"github.com/git-town/git-town/v21/internal/test/testgit"
	"github.com/git-town/git-town/v21/pkg/asserts"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/google/go-cmp/cmp"
	"github.com/kballard/go-shellquote"
)

// the global FixtureFactory instance.
var fixtureFactory *fixture.Factory

// dedicated type for storing data in context.Context
type key int

// the key for storing the state in the context.Context
const (
	keyScenarioState key = iota
	keyScenarioName
	keyScenarioTags
)

func InitializeScenario(scenarioContext *godog.ScenarioContext) {
	scenarioContext.Before(func(ctx context.Context, scenario *godog.Scenario) (context.Context, error) {
		ctx = context.WithValue(ctx, keyScenarioName, scenario.Name)
		ctx = context.WithValue(ctx, keyScenarioTags, scenario.Tags)
		return ctx, nil
	})

	scenarioContext.After(func(ctx context.Context, scenario *godog.Scenario, err error) (context.Context, error) {
		ctxValue := ctx.Value(keyScenarioState)
		if ctxValue == nil {
			return ctx, errors.New("after-scenario hook has found no scenario state found to clean up")
		}
		state := ctxValue.(*ScenarioState)
		if err != nil {
			fmt.Printf("failed scenario %q in %s - investigate state in %s\n", scenario.Name, scenario.Uri, state.fixture.Dir)
			return ctx, nil //nolint:nilerr
		}
		exitCode := state.runExitCode.GetOrPanic()
		if exitCode != 0 && !state.runExitCodeChecked {
			print.Error(fmt.Errorf("%s - scenario %q doesn't document exit code %d", scenario.Uri, scenario.Name, exitCode))
			os.Exit(1)
		}
		if state != nil {
			state.fixture.Delete()
		}
		return ctx, nil
	})
}

func InitializeSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		factory := fixture.CreateFactory()
		fixtureFactory = &factory
	})
	ctx.AfterSuite(func() {
		fixtureFactory.Remove()
	})
	defineSteps(ctx.ScenarioContext())
}

func defineSteps(sc *godog.ScenarioContext) {
	sc.Step(`^a brand-new Git repo$`, func(ctx context.Context) (context.Context, error) {
		scenarioName := ctx.Value(keyScenarioName).(string)
		scenarioTags := ctx.Value(keyScenarioTags).([]*messages.PickleTag)
		fixture := fixtureFactory.CreateEmptyFixture(scenarioName)
		if helpers.HasTag(scenarioTags, "@debug") {
			fixture.DevRepo.GetOrPanic().Verbose = true
		}
		state := ScenarioState{
			beforeRunDevSHAs:     None[gitdomain.Commits](),
			beforeRunOriginSHAs:  None[gitdomain.Commits](),
			browserVariable:      None[string](),
			fixture:              fixture,
			initialBranches:      None[datatable.DataTable](),
			initialCommits:       None[datatable.DataTable](),
			initialCurrentBranch: None[gitdomain.LocalBranchName](),
			initialDevSHAs:       None[gitdomain.Commits](),
			initialLineage:       None[datatable.DataTable](),
			initialOriginSHAs:    None[gitdomain.Commits](),
			initialTags:          None[datatable.DataTable](),
			initialWorktreeSHAs:  None[gitdomain.Commits](),
			insideGitRepo:        true,
			runExitCode:          None[int](),
			runExitCodeChecked:   false,
			runOutput:            None[string](),
			uncommittedContent:   None[string](),
			uncommittedFileName:  None[string](),
		}
		return context.WithValue(ctx, keyScenarioState, &state), nil
	})

	sc.Step(`^a coworker clones the repository$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.AddCoworkerRepo()
	})

	sc.Step(`^a folder "([^"]*)"$`, func(ctx context.Context, name string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.CreateFolder(name)
	})

	sc.Step(`^a Git repo with origin$`, func(ctx context.Context) (context.Context, error) {
		scenarioName := ctx.Value(keyScenarioName).(string)
		scenarioTags := ctx.Value(keyScenarioTags).([]*messages.PickleTag)
		fixture := fixtureFactory.CreateFixture(scenarioName)
		if helpers.HasTag(scenarioTags, "@debug") {
			fixture.DevRepo.GetOrPanic().Verbose = true
			fixture.OriginRepo.GetOrPanic().Verbose = true
		}
		state := ScenarioState{
			beforeRunDevSHAs:     None[gitdomain.Commits](),
			beforeRunOriginSHAs:  None[gitdomain.Commits](),
			browserVariable:      None[string](),
			fixture:              fixture,
			initialBranches:      None[datatable.DataTable](),
			initialCommits:       None[datatable.DataTable](),
			initialCurrentBranch: None[gitdomain.LocalBranchName](),
			initialDevSHAs:       None[gitdomain.Commits](),
			initialLineage:       None[datatable.DataTable](),
			initialOriginSHAs:    None[gitdomain.Commits](),
			initialTags:          None[datatable.DataTable](),
			initialWorktreeSHAs:  None[gitdomain.Commits](),
			insideGitRepo:        true,
			runExitCode:          None[int](),
			runExitCodeChecked:   false,
			runOutput:            None[string](),
			uncommittedContent:   None[string](),
			uncommittedFileName:  None[string](),
		}
		return context.WithValue(ctx, keyScenarioState, &state), nil
	})

	sc.Step(`^all branches are now synchronized$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branchesOutOfSync, output := devRepo.HasBranchesOutOfSync()
		if branchesOutOfSync {
			return errors.New("unexpected out of sync:\n" + output)
		}
		return nil
	})

	sc.Step(`^a local Git repo$`, func(ctx context.Context) (context.Context, error) {
		scenarioName := ctx.Value(keyScenarioName).(string)
		scenarioTags := ctx.Value(keyScenarioTags).([]*messages.PickleTag)
		fixture := fixtureFactory.CreateFixture(scenarioName)
		devRepo := fixture.DevRepo.GetOrPanic()
		if helpers.HasTag(scenarioTags, "@debug") {
			devRepo.Verbose = true
		}
		devRepo.RemoveRemote(gitdomain.RemoteOrigin)
		fixture.OriginRepo = MutableNone[commands.TestCommands]()
		state := ScenarioState{
			beforeRunDevSHAs:     None[gitdomain.Commits](),
			beforeRunOriginSHAs:  None[gitdomain.Commits](),
			browserVariable:      None[string](),
			fixture:              fixture,
			initialBranches:      None[datatable.DataTable](),
			initialCommits:       None[datatable.DataTable](),
			initialCurrentBranch: None[gitdomain.LocalBranchName](),
			initialDevSHAs:       None[gitdomain.Commits](),
			initialLineage:       None[datatable.DataTable](),
			initialOriginSHAs:    None[gitdomain.Commits](),
			initialTags:          None[datatable.DataTable](),
			initialWorktreeSHAs:  None[gitdomain.Commits](),
			insideGitRepo:        true,
			runExitCode:          None[int](),
			runExitCodeChecked:   false,
			runOutput:            None[string](),
			uncommittedContent:   None[string](),
			uncommittedFileName:  None[string](),
		}
		return context.WithValue(ctx, keyScenarioState, &state), nil
	})

	sc.Step(`^a merge is now in progress$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		if !devRepo.Git.HasMergeInProgress(devRepo.TestRunner) {
			return errors.New("expected merge in progress")
		}
		return nil
	})

	sc.Step(`^an additional "([^"]+)" remote with URL "([^"]+)"$`, func(ctx context.Context, remote, url string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.AddRemote(gitdomain.Remote(remote), url)
	})

	sc.Step(`^an uncommitted file$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		filename := "uncommitted file"
		state.uncommittedFileName = Some(filename)
		content := "uncommitted content"
		state.uncommittedContent = Some(content)
		devRepo.CreateFile(
			filename,
			content,
		)
	})

	sc.Step(`^an uncommitted file with name "([^"]+)" and content "([^"]+)"$`, func(ctx context.Context, name, content string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.uncommittedFileName = Some(name)
		state.uncommittedContent = Some(content)
		devRepo.CreateFile(name, content)
	})

	sc.Step(`^an uncommitted file with name "([^"]+)" exists now$`, func(ctx context.Context, filename string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		files := devRepo.UncommittedFiles()
		want := []string{filename}
		if !reflect.DeepEqual(files, want) {
			return fmt.Errorf("expected %s but found %s", want, files)
		}
		return nil
	})

	sc.Step(`^an upstream repo$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.AddUpstream()
	})

	sc.Step(`^a proposal for this branch does not exist`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.TestRunner.ProposalOverride = Some(forgedomain.OverrideNoProposal)
	})

	sc.Step(`^a proposal for this branch exists at "([^"]+)"`, func(ctx context.Context, url string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.TestRunner.ProposalOverride = Some(url)
	})

	sc.Step(`^a rebase is (?:now|still) in progress$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		repoStatus := asserts.NoError1(devRepo.Git.RepoStatus(devRepo.TestRunner))
		if !repoStatus.RebaseInProgress {
			return errors.New("expected rebase in progress")
		}
		return nil
	})

	sc.Step(`^a remote "([^"]+)" pointing to "([^"]+)"`, func(ctx context.Context, name, url string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.AddRemote(gitdomain.Remote(name), url)
	})

	sc.Step(`^a remote tag "([^"]+)" not on a branch$`, func(ctx context.Context, name string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.OriginRepo.GetOrPanic().CreateStandaloneTag(name)
	})

	sc.Step(`^branch "([^"]+)" is active in another worktree`, func(ctx context.Context, branch string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.AddSecondWorktree(gitdomain.NewLocalBranchName(branch))
	})

	sc.Step(`^branch "([^"]+)" (?:now|still) has type "(\w+)"$`, func(ctx context.Context, branchName, branchTypeName string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(branchName)
		wantOpt := asserts.NoError1(configdomain.ParseBranchType(branchTypeName))
		want := wantOpt.GetOrPanic()
		have := devRepo.Config.BranchType(branch)
		if have != want {
			return fmt.Errorf("branch %q is %s", branch, have)
		}
		return nil
	})

	sc.Step(`^commit "([^"]+)" on branch "([^"]+)" now has this full commit message`, func(ctx context.Context, title, branchText string, expected *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(branchText)
		parent := devRepo.Config.NormalConfig.Lineage.Parent(branch).GetOrPanic()
		sha := devRepo.CommitSHA(devRepo, title, branch, parent.BranchName())
		have := asserts.NoError1(devRepo.Git.CommitMessage(devRepo, sha)).String()
		want := expected.Content
		if have != want {
			return fmt.Errorf("\nwant:\n%q\n\nhave:\n%q", want, have)
		}
		return nil
	})

	sc.Step(`^display "([^"]+)"$`, func(ctx context.Context, command string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		parts := strings.Split(command, " ")
		output, err := devRepo.TestRunner.Query(parts[0], parts[1:]...)
		fmt.Println("XXXXXXXXXXXXXXXXX " + strings.ToUpper(command) + " START XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		fmt.Println(output)
		fmt.Println("XXXXXXXXXXXXXXXXX " + strings.ToUpper(command) + " END XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX")
		return err
	})

	sc.Step(`^file "([^"]*)" (?:now|still) has content:$`, func(ctx context.Context, file string, expectedContent *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actualContent := strings.TrimSpace(devRepo.FileContent(file))
		expectedText := handlebars.Expand(expectedContent.Content, handlebars.ExpandArgs{
			BeforeRunDevSHAs:       state.beforeRunDevSHAs.GetOrPanic(),
			BeforeRunOriginSHAsOpt: state.beforeRunOriginSHAs,
			InitialDevCommits:      state.initialDevSHAs.GetOrPanic(),
			InitialOriginCommits:   state.initialOriginSHAs,
			InitialWorktreeCommits: state.initialWorktreeSHAs,
			LocalRepo:              devRepo,
			RemoteRepo:             state.fixture.OriginRepo.Value,
			WorktreeRepo:           state.fixture.SecondWorktree.Value,
		})
		if expectedText != actualContent {
			return fmt.Errorf("file content does not match\n\nEXPECTED:\n%q\n\nACTUAL:\n\n%q\n----------------------------", expectedText, actualContent)
		}
		return nil
	})

	sc.Step(`^file "([^"]+)" (?:now|still) has content "([^"]*)"$`, func(ctx context.Context, file, expectedContent string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actualContent := devRepo.FileContent(file)
		if expectedContent != actualContent {
			return fmt.Errorf("file content does not match\n\nEXPECTED:\n%q\n\nACTUAL:\n\n%q\n----------------------------", expectedContent, actualContent)
		}
		return nil
	})

	sc.Step(`^file "([^"]+)" with content$`, func(ctx context.Context, name string, content *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		filePath := filepath.Join(devRepo.WorkingDir, name)
		//nolint:gosec // need permission 700 here in order for tests to work
		return os.WriteFile(filePath, []byte(content.Content), 0o700)
	})

	sc.Step(`^Git has version "([^"]*)"$`, func(ctx context.Context, version string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.MockGit(version)
	})

	sc.Step(`^Git Town does not print "(.+)"$`, func(ctx context.Context, text string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		if strings.Contains(stripansi.Strip(state.runOutput.GetOrPanic()), text) {
			return fmt.Errorf("text found: %q", text)
		}
		return nil
	})

	sc.Step(`^Git Town is no longer configured$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		return devRepo.VerifyNoGitTownConfiguration()
	})

	sc.Step(`^Git Town is not configured$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		asserts.NoError(devRepo.RemovePerennialBranchConfiguration())
		devRepo.RemoveMainBranchConfiguration()
	})

	sc.Step(`^Git Town parent setting for branch "([^"]*)" is "([^"]*)"$`, func(ctx context.Context, branch, value string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branchName := gitdomain.NewLocalBranchName(branch)
		configKey := configdomain.NewParentKey(branchName)
		return devRepo.Config.NormalConfig.GitConfigAccess.SetConfigValue(configdomain.ConfigScopeLocal, configKey, value)
	})

	sc.Step(`^Git Town prints:$`, func(ctx context.Context, expected *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		if exitCode := state.runExitCode.GetOrPanic(); exitCode != 0 {
			return fmt.Errorf("unexpected exit code %d", exitCode)
		}
		output := stripansi.Strip(state.runOutput.GetOrPanic())
		if !strings.Contains(output, strings.TrimRight(expected.Content, "\n")) {
			fmt.Println("ERROR: text not found:")
			fmt.Println("==================================================================")
			fmt.Println("EXPECTED OUTPUT START ============================================")
			fmt.Println("==================================================================")
			fmt.Println()
			fmt.Println(expected.Content)
			fmt.Println()
			fmt.Println("==================================================================")
			fmt.Println("EXPECTED OUTPUT END ==============================================")
			fmt.Println("==================================================================")
			fmt.Println()
			fmt.Println("==================================================================")
			fmt.Println("ACTUAL OUTPUT START ==============================================")
			fmt.Println("==================================================================")
			fmt.Println()
			fmt.Println(output)
			fmt.Println()
			fmt.Println("==================================================================")
			fmt.Println("ACTUAL OUTPUT END ================================================")
			fmt.Println("==================================================================")
			fmt.Println()
			return errors.New("expected text not found")
		}
		return nil
	})

	sc.Step(`^Git Town prints no output$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		output := state.runOutput.GetOrPanic()
		if len(output) > 0 {
			return fmt.Errorf("expected no output but found %q", output)
		}
		return nil
	})

	sc.Step(`^Git Town prints something like:$`, func(ctx context.Context, expected *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		regex := regexp.MustCompile(expected.Content)
		have := stripansi.Strip(state.runOutput.GetOrPanic())
		if !regex.MatchString(have) {
			return fmt.Errorf("EXPECTED: content matching %q\nGOT: %q", expected.Content, have)
		}
		return nil
	})

	sc.Step(`^Git Town prints the error:$`, func(ctx context.Context, expected *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.runExitCodeChecked = true
		if !strings.Contains(stripansi.Strip(state.runOutput.GetOrPanic()), expected.Content) {
			return fmt.Errorf("text not found:\n%s\n\nactual text:\n%s", expected.Content, state.runOutput.GetOrDefault())
		}
		if exitCode := state.runExitCode.GetOrPanic(); exitCode == 0 {
			return fmt.Errorf("unexpected exit code %d", exitCode)
		}
		return nil
	})

	sc.Step(`^Git Town runs no commands$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		commands := output.GitCommandsInGitTownOutput(state.runOutput.GetOrPanic())
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

	sc.Step(`^Git Town runs the commands$`, func(ctx context.Context, input *godog.Table) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		commands := output.GitCommandsInGitTownOutput(state.runOutput.GetOrPanic())
		table := output.RenderExecutedGitCommands(commands, input)
		dataTable := datatable.FromGherkin(input)
		expanded := dataTable.Expand(handlebars.ExpandArgs{
			BeforeRunDevSHAs:       state.beforeRunDevSHAs.GetOrPanic(),
			BeforeRunOriginSHAsOpt: state.beforeRunOriginSHAs,
			InitialDevCommits:      state.initialDevSHAs.GetOrPanic(),
			InitialOriginCommits:   state.initialOriginSHAs,
			InitialWorktreeCommits: state.initialWorktreeSHAs,
			LocalRepo:              devRepo,
			RemoteRepo:             state.fixture.OriginRepo.Value,
			WorktreeRepo:           state.fixture.SecondWorktree.Value,
		})
		diff, errorCount := table.EqualDataTable(expanded)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the commands run\n\n", errorCount)
			fmt.Println(diff)
			return errors.New("mismatching commands run, see diff above")
		}
		return nil
	})

	sc.Step(`^Git Town runs without errors`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		exitCode := state.runExitCode.GetOrPanic()
		if exitCode != 0 {
			return errors.New("unexpected failure of scenario")
		}
		return nil
	})

	sc.Step(`^(global |local |)Git setting "([^"]+)" is "([^"]*)"$`, func(ctx context.Context, scope, key, value string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		parsedScope := configdomain.ParseConfigScope(scope)
		return devRepo.Config.NormalConfig.GitConfigAccess.SetConfigValue(parsedScope, configdomain.Key(key), value)
	})

	sc.Step(`^(global |local |)Git setting "([^"]+)" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, scope, name, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		parsedScope := configdomain.ParseConfigScope(scope)
		have := devRepo.SnapShots[parsedScope][configdomain.Key(name)]
		if have != want {
			return fmt.Errorf("unexpected value for key %q: want %q have %q", name, want, have)
		}
		return nil
	})

	sc.Step(`^(global |local |)Git setting "([^"]+)" (?:now|still) doesn't exist$`, func(ctx context.Context, scope, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		parsedScope := configdomain.ParseConfigScope(scope)
		have, has := devRepo.SnapShots[parsedScope][configdomain.Key(name)]
		if has {
			return fmt.Errorf("unexpected value for %q: %q", name, have)
		}
		return nil
	})

	sc.Step(`^I add an unrelated stash entry with file "([^"]+)"$`, func(ctx context.Context, filename string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.CreateFile(filename, "stash content")
		devRepo.StashOpenFiles()
	})

	sc.Step(`^I add commit "([^"]*)" to the "([^"]*)" branch$`, func(ctx context.Context, message, branch string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.CreateCommit(testgit.Commit{
			Branch:   gitdomain.NewLocalBranchName(branch),
			FileName: "new_file",
			Message:  gitdomain.CommitMessage(message),
		})
	})

	sc.Step(`^I add this commit to the "([^"]*)" branch$`, func(ctx context.Context, branch string, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		commit := testgit.FromGherkinTable(table)[0]
		commit.Branch = gitdomain.LocalBranchName(branch)
		devRepo.CreateCommit(commit)
	})

	sc.Step(`^I add this commit to the current branch:$`, func(ctx context.Context, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		commit := testgit.FromGherkinTable(table)[0]
		devRepo.CreateFile(commit.FileName, commit.FileContent)
		devRepo.StageFiles(commit.FileName)
		devRepo.CommitStagedChanges(commit.Message)
	})

	sc.Step(`^I amend this commit$`, func(ctx context.Context, table *godog.Table) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		commits := testgit.FromGherkinTable(table)
		if len(commits) != 1 {
			return errors.New("expected exactly one commit")
		}
		commit := commits[0]
		devRepo.CheckoutBranch(commit.Branch)
		devRepo.CreateFile(commit.FileName, commit.FileContent)
		asserts.NoError(devRepo.Run("git", "add", commit.FileName))
		return devRepo.Run("git", "commit", "--amend", "--message", commit.Message.String())
	})

	sc.Step(`^I am not prompted for any parent branches$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		notExpected := "Please specify the parent branch of"
		if strings.Contains(state.runOutput.GetOrPanic(), notExpected) {
			return fmt.Errorf("text found:\n\nDID NOT EXPECT: %q\n\nACTUAL\n\n%q\n----------------------------", notExpected, state.runOutput)
		}
		return nil
	})

	sc.Step(`^I am outside a Git repo$`, func(ctx context.Context) (context.Context, error) {
		scenarioName := ctx.Value(keyScenarioName).(string)
		// scenarioTags := ctx.Value(keyScenarioTags).([]*messages.PickleTag)
		envDirName := filesystem.FolderName(scenarioName) + "_" + fixtureFactory.Counter.NextAsString()
		envPath := filepath.Join(fixtureFactory.Dir, envDirName)
		asserts.NoError(os.Mkdir(envPath, 0o777))
		fixture := fixture.Fixture{
			CoworkerRepo:   MutableNone[commands.TestCommands](),
			DevRepo:        MutableNone[commands.TestCommands](),
			Dir:            envPath,
			OriginRepo:     MutableNone[commands.TestCommands](),
			SecondWorktree: MutableNone[commands.TestCommands](),
			SubmoduleRepo:  MutableNone[commands.TestCommands](),
			UpstreamRepo:   MutableNone[commands.TestCommands](),
		}
		state := ScenarioState{
			beforeRunDevSHAs:     None[gitdomain.Commits](),
			beforeRunOriginSHAs:  None[gitdomain.Commits](),
			browserVariable:      None[string](),
			fixture:              fixture,
			initialBranches:      None[datatable.DataTable](),
			initialCommits:       None[datatable.DataTable](),
			initialCurrentBranch: None[gitdomain.LocalBranchName](),
			initialDevSHAs:       None[gitdomain.Commits](),
			initialLineage:       None[datatable.DataTable](),
			initialOriginSHAs:    None[gitdomain.Commits](),
			initialTags:          None[datatable.DataTable](),
			initialWorktreeSHAs:  None[gitdomain.Commits](),
			insideGitRepo:        true,
			runExitCode:          None[int](),
			runExitCodeChecked:   false,
			runOutput:            None[string](),
			uncommittedContent:   None[string](),
			uncommittedFileName:  None[string](),
		}
		return context.WithValue(ctx, keyScenarioState, &state), nil
	})

	sc.Step(`^in a separate terminal I create branch "([^"]+)" with commits$`, func(ctx context.Context, branchName string, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		existingBranch := devRepo.Git.CurrentBranchCache.Value()
		newBranch := gitdomain.NewLocalBranchName(branchName)
		devRepo.CreateBranch(newBranch, "main")
		devRepo.CheckoutBranch(newBranch)
		for _, commit := range testgit.FromGherkinTable(table) {
			devRepo.CreateFile(commit.FileName, commit.FileContent)
			devRepo.StageFiles(commit.FileName)
			devRepo.CommitStagedChanges(commit.Message)
		}
		devRepo.CheckoutBranch(existingBranch)
	})

	sc.Step(`^inspect the commits$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		if devRepo, hasDevRepo := state.fixture.DevRepo.Get(); hasDevRepo {
			fmt.Println("\nsha")
			fmt.Println(asserts.NoError1(devRepo.Query("git", "branch", "-vva")))
		}
		if originRepo, hasOriginRepo := state.fixture.OriginRepo.Get(); hasOriginRepo {
			fmt.Println("\nsha-in-origin")
			fmt.Println(asserts.NoError1(originRepo.Query("git", "branch", "-vva")))
		}
		if initialSHAs, hasInitialSHAs := state.initialDevSHAs.Get(); hasInitialSHAs {
			fmt.Println("\nsha-initial")
			for _, commit := range initialSHAs {
				fmt.Printf("- %s (%s)\n", commit.SHA.Truncate(7), commit.Message)
			}
		}
		if initialOriginSHAs, hasInitialOriginSHAs := state.initialOriginSHAs.Get(); hasInitialOriginSHAs {
			fmt.Println("\nsha-in-origin-initial")
			for _, commit := range initialOriginSHAs {
				fmt.Printf("- %s (%s)\n", commit.SHA.Truncate(7), commit.Message)
			}
		}
		if worktreeRepo, hasWorktreeRepo := state.fixture.SecondWorktree.Get(); hasWorktreeRepo {
			fmt.Println("\nsha-in-worktree")
			fmt.Println(asserts.NoError1(worktreeRepo.Query("git", "branch", "-vva")))
		}
		if initialWorktreeSHAs, hasInitialWorktreeSHAs := state.initialWorktreeSHAs.Get(); hasInitialWorktreeSHAs {
			fmt.Println("\nsha-in-worktree-initial")
			for _, commit := range initialWorktreeSHAs {
				fmt.Printf("- %s (%s)\n", commit.SHA.Truncate(7), commit.Message)
			}
		}
		if devBeforeRunSHAs, hasDevBeforeRunSHAs := state.beforeRunDevSHAs.Get(); hasDevBeforeRunSHAs {
			fmt.Println("\nsha-before-run")
			for _, commit := range devBeforeRunSHAs {
				fmt.Printf("- %s (%s)\n", commit.SHA.Truncate(7), commit.Message)
			}
		}
		if originBeforeRunSHAs, hasOriginBeforeRunSHAs := state.beforeRunOriginSHAs.Get(); hasOriginBeforeRunSHAs {
			fmt.Println("\nsha-in-origin-before-run")
			for _, commit := range originBeforeRunSHAs {
				fmt.Printf("- %s (%s)\n", commit.SHA.Truncate(7), commit.Message)
			}
		}
		fmt.Println()
	})

	sc.Step(`^inspect the repo$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		fmt.Printf("\nThe workspace is at %s\n", devRepo.WorkingDir)
		time.Sleep(1 * time.Hour)
	})

	sc.Step(`^I pipe the following text into "([^"]+)":$`, func(ctx context.Context, cmd string, input *godog.DocString) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		env := os.Environ()
		if browserPath, has := state.browserVariable.Get(); has {
			env = envvars.Replace(env, browser.EnvVarName, browserPath)
		}
		output, exitCode := devRepo.MustQueryStringCodeWith(cmd, &subshell.Options{
			Env:   env,
			Input: Some(input.Content),
		})
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Reload()
	})

	sc.Step(`^I ran "(.+)"$`, func(ctx context.Context, command string) error {
		runCommand(ctx, command)
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		if exitCode, hasExitCode := state.runExitCode.Get(); hasExitCode {
			if exitCode != 0 {
				fmt.Println("Output from failed command:")
				fmt.Println(state.runOutput.GetOrDefault())
				return fmt.Errorf("unexpected exit code: %d", exitCode)
			}
		}
		return nil
	})

	sc.Step(`^I ran "(.+)" and ignore the error$`, func(ctx context.Context, command string) error {
		runCommand(ctx, command)
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		if exitCode, hasExitCode := state.runExitCode.Get(); hasExitCode {
			if exitCode == 0 {
				return errors.New("this command should fail")
			}
		}
		return nil
	})

	sc.Step(`^I rename the "([^"]+)" remote to "([^"]+)"$`, func(ctx context.Context, oldName, newName string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.RenameRemote(oldName, newName)
	})

	sc.Step(`^I resolve the conflict in "([^"]*)" in the other worktree$`, func(ctx context.Context, filename string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		content := "resolved content"
		secondWorkTree := state.fixture.SecondWorktree.GetOrPanic()
		secondWorkTree.CreateFile(filename, content)
		secondWorkTree.StageFiles(filename)
	})

	sc.Step(`^I resolve the conflict in "([^"]*)"(?: with "([^"]*)")?$`, func(ctx context.Context, filename, content string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		if content == "" {
			content = "resolved content"
		}
		content = strings.ReplaceAll(content, "\\n", "\n")
		devRepo.CreateFile(filename, content)
		devRepo.StageFiles(filename)
	})

	sc.Step(`^I run "(.+)"$`, func(ctx context.Context, command string) {
		runCommand(ctx, command)
	})

	sc.Step(`^I run "(.+)" with these environment variables$`, func(ctx context.Context, command string, envVars *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		env := os.Environ()
		for _, row := range envVars.Rows {
			env = append(env, fmt.Sprintf("%s=%s", row.Cells[0].Value, row.Cells[1].Value))
		}
		output, exitCode := devRepo.MustQueryStringCodeWith(command, &subshell.Options{Env: env})
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Reload()
	})

	sc.Step(`^I run "([^"]*)" and close the editor$`, func(ctx context.Context, cmd string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		env := append(os.Environ(), "GIT_EDITOR=true")
		output, exitCode := devRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Env: env})
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Reload()
	})

	sc.Step(`^I run "([^"]*)" and enter an empty commit message$`, func(ctx context.Context, cmd string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		devRepo.MockCommitMessage("")
		output, exitCode := devRepo.MustQueryStringCode(cmd)
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Reload()
	})

	sc.Step(`^I run "([^"]*)" and enter "([^"]*)" for the commit message$`, func(ctx context.Context, cmd, message string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		devRepo.MockCommitMessage(message)
		env := os.Environ()
		if browserPath, has := state.browserVariable.Get(); has {
			env = envvars.Replace(env, browser.EnvVarName, browserPath)
		}
		output, exitCode := devRepo.MustQueryStringCodeWith(cmd, &subshell.Options{
			Env: env,
		})
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Reload()
	})

	sc.Step(`^I run "([^"]+)" in the "([^"]+)" folder$`, func(ctx context.Context, cmd, folderName string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		output, exitCode := devRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Dir: folderName})
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Reload()
	})

	sc.Step(`^I run "([^"]+)" in the other worktree$`, func(ctx context.Context, cmd string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.CaptureState()
		updateInitialSHAs(state)
		secondWorkTree := state.fixture.SecondWorktree.GetOrPanic()
		output, exitCode := secondWorkTree.MustQueryStringCode(cmd)
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		secondWorkTree.Reload()
	})

	sc.Step(`^I run "([^"]*)" in the other worktree and enter "([^"]*)" for the commit message$`, func(ctx context.Context, cmd, message string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.CaptureState()
		updateInitialSHAs(state)
		secondWorkTree := state.fixture.SecondWorktree.GetOrPanic()
		secondWorkTree.MockCommitMessage(message)
		output, exitCode := secondWorkTree.MustQueryStringCode(cmd)
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		secondWorkTree.Reload()
	})

	sc.Step(`^I (?:run|ran) "([^"]+)" and enter into the dialogs?:$`, func(ctx context.Context, cmd string, input *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		env := os.Environ()
		if browserPath, has := state.browserVariable.Get(); has {
			env = envvars.Replace(env, browser.EnvVarName, browserPath)
		}
		answers := asserts.NoError1(helpers.TableToInputEnv(input))
		for dialogNumber, answer := range answers {
			env = append(env, fmt.Sprintf("%s_%02d=%s", components.TestInputKey, dialogNumber, answer))
		}
		output, exitCode := devRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Env: env})
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Reload()
	})

	sc.Step(`^my repo's "([^"]*)" remote is "([^"]*)"$`, func(ctx context.Context, remoteName, remoteURL string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		remote := gitdomain.Remote(remoteName)
		devRepo.RemoveRemote(remote)
		devRepo.AddRemote(remote, remoteURL)
	})

	sc.Step(`^my repo has a Git submodule$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.fixture.AddSubmoduleRepo()
		devRepo.AddSubmodule(state.fixture.SubmoduleRepo.GetOrPanic().WorkingDir)
	})

	sc.Step(`^no commits exist now$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		currentCommits := state.fixture.CommitTable(state.initialCommits.GetOrPanic().Cells[0])
		noCommits := datatable.DataTable{}
		noCommits.AddRow(state.initialCommits.GetOrPanic().Cells[0]...)
		errDiff, errCount := currentCommits.EqualDataTable(noCommits)
		if errCount == 0 {
			return nil
		}
		fmt.Println(errDiff)
		return errors.New("found unexpected commits")
	})

	sc.Step(`^no lineage exists now$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		if devRepo.Config.NormalConfig.Lineage.Len() > 0 {
			lineage := devRepo.Config.NormalConfig.Lineage
			return fmt.Errorf("unexpected Git Town lineage information: %+v", lineage)
		}
		return nil
	})

	sc.Step(`^no merge is in progress$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		if devRepo.Git.HasMergeInProgress(devRepo.TestRunner) {
			return errors.New("expected no merge in progress")
		}
		return nil
	})

	sc.Step(`^no rebase is now in progress$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		repoStatus := asserts.NoError1(devRepo.Git.RepoStatus(devRepo.TestRunner))
		if repoStatus.RebaseInProgress {
			return errors.New("expected no rebase in progress")
		}
		return nil
	})

	sc.Step(`^no tool to open browsers is installed$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.browserVariable = Some(string(browser.EnvVarNone))
	})

	sc.Step(`^no uncommitted files exist now$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		files := devRepo.UncommittedFiles()
		if len(files) > 0 {
			return fmt.Errorf("unexpected uncommitted files: %s", files)
		}
		return nil
	})

	sc.Step(`^offline mode is enabled$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		return devRepo.Config.NormalConfig.SetOffline(true)
	})

	sc.Step(`^origin deletes the "([^"]*)" branch$`, func(ctx context.Context, branch string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.OriginRepo.GetOrPanic().RemoveBranch(gitdomain.NewLocalBranchName(branch))
	})

	sc.Step(`^origin ships the "([^"]*)" branch using the "squash-merge" ship-strategy$`, func(ctx context.Context, branchName string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		branchToShip := gitdomain.NewLocalBranchName(branchName)
		originRepo := state.fixture.OriginRepo.GetOrPanic()
		commitMessage := asserts.NoError1(originRepo.Git.FirstCommitMessageInBranch(originRepo.TestRunner, branchToShip.BranchName(), "main"))
		message, hasCommitMessage := commitMessage.Get()
		if !hasCommitMessage {
			return errors.New("branch to ship contains no commits")
		}
		originRepo.CheckoutBranch("main")
		asserts.NoError(originRepo.Git.SquashMerge(originRepo.TestRunner, branchToShip))
		originRepo.StageFiles("-A")
		asserts.NoError(originRepo.Git.Commit(originRepo.TestRunner, configdomain.UseCustomMessage(message), gitdomain.NewAuthorOpt("CI <ci@acme.com>"), configdomain.CommitHookEnabled))
		originRepo.RemoveBranch(branchToShip)
		originRepo.CheckoutBranch("initial")
		return nil
	})

	sc.Step(`^origin ships the "([^"]*)" branch using the "squash-merge" ship-strategy and resolves the merge conflict in "([^"]+)" with "([^"]+)" and commits as "([^"]+)"$`, func(ctx context.Context, branchName, fileName, fileContent, commitMessage string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		branchToShip := gitdomain.NewLocalBranchName(branchName)
		originRepo := state.fixture.OriginRepo.GetOrPanic()
		originRepo.CheckoutBranch("main")
		err := originRepo.Git.SquashMerge(originRepo.TestRunner, branchToShip)
		if err == nil {
			panic("expected a merge conflict here")
		}
		originRepo.CreateFile(fileName, fileContent)
		originRepo.StageFiles("-A")
		asserts.NoError(originRepo.Git.Commit(originRepo.TestRunner, configdomain.UseCustomMessage(gitdomain.CommitMessage(commitMessage)), gitdomain.NewAuthorOpt("CI <ci@acme.com>"), configdomain.CommitHookEnabled))
		originRepo.RemoveBranch(branchToShip)
		originRepo.CheckoutBranch("initial")
	})

	sc.Step(`^origin ships the "([^"]*)" branch using the "squash-merge" ship-strategy as "([^"]+)"$`, func(ctx context.Context, branchName, commitMessage string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		branchToShip := gitdomain.NewLocalBranchName(branchName)
		originRepo := state.fixture.OriginRepo.GetOrPanic()
		originRepo.CheckoutBranch("main")
		asserts.NoError(originRepo.Git.SquashMerge(originRepo.TestRunner, branchToShip))
		originRepo.StageFiles("-A")
		asserts.NoError(originRepo.Git.Commit(originRepo.TestRunner, configdomain.UseCustomMessage(gitdomain.CommitMessage(commitMessage)), gitdomain.NewAuthorOpt("CI <ci@acme.com>"), configdomain.CommitHookEnabled))
		originRepo.RemoveBranch(branchToShip)
		originRepo.CheckoutBranch("initial")
		return nil
	})

	sc.Step(`^the branches$`, func(ctx context.Context, table *godog.Table) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		for _, branchSetup := range datatable.ParseBranchSetupTable(table) {
			var repoToCreateBranchIn *commands.TestCommands
			switch {
			case
				branchSetup.Locations.Is(testgit.LocationLocal),
				branchSetup.Locations.Is(testgit.LocationLocal, testgit.LocationOrigin):
				repoToCreateBranchIn = state.fixture.DevRepo.GetOrPanic()
			case branchSetup.Locations.Is(testgit.LocationOrigin):
				repoToCreateBranchIn = state.fixture.OriginRepo.GetOrPanic()
			case branchSetup.Locations.Is(testgit.LocationUpstream):
				repoToCreateBranchIn = state.fixture.UpstreamRepo.GetOrPanic()
			default:
				return errors.New("unhandled location to create the new branch: " + branchSetup.Locations.String())
			}
			branchType, hasBranchType := branchSetup.BranchType.Get()
			if hasBranchType {
				switch branchType {
				case configdomain.BranchTypeMainBranch:
					return errors.New("main branch exists already")
				case configdomain.BranchTypeFeatureBranch:
					repoToCreateBranchIn.CreateFeatureBranch(branchSetup.Name, branchSetup.Parent.GetOrPanic().BranchName())
				case
					configdomain.BranchTypePerennialBranch,
					configdomain.BranchTypeContributionBranch,
					configdomain.BranchTypeObservedBranch,
					configdomain.BranchTypeParkedBranch,
					configdomain.BranchTypePrototypeBranch:
					repoToCreateBranchIn.CreateBranchOfType(branchSetup.Name, branchSetup.Parent, branchType)
				}
			} else {
				repoToCreateBranchIn.CreateBranch(branchSetup.Name, "main")
				if parent, hasParent := branchSetup.Parent.Get(); hasParent {
					asserts.NoError(repoToCreateBranchIn.Config.NormalConfig.SetParent(branchSetup.Name, parent))
				}
			}
			if len(branchSetup.Locations) > 1 {
				switch {
				case branchSetup.Locations.Is(testgit.LocationLocal, testgit.LocationOrigin):
					state.fixture.DevRepo.GetOrPanic().PushBranchToRemote(branchSetup.Name, gitdomain.RemoteOrigin)
				default:
					return errors.New("unhandled location to push the new branch to: " + branchSetup.Locations.String())
				}
			}
		}
		return nil
	})

	sc.Step(`^the branches are now$`, func(ctx context.Context, table *godog.Table) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		existing := state.fixture.Branches()
		diff, errCount := existing.EqualGherkin(table)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the branches\n\n", errCount)
			fmt.Println(diff)
			return errors.New("mismatching branches found, see the diff above")
		}
		return nil
	})

	sc.Step(`^the branches contain these files:$`, func(ctx context.Context, godogTable *godog.Table) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		repo := state.fixture.DevRepo.GetOrPanic()
		branches, _ := asserts.NoError2(repo.LocalBranches())
		haveTable := datatable.DataTable{}
		haveTable.AddRow("BRANCH", "NAME")
		for _, branch := range branches {
			repo.CheckoutBranch(branch)
			firstFileInBranch := true
			for _, file := range repo.FilesInWorkspace() {
				if firstFileInBranch {
					haveTable.AddRow(branch.String(), file)
					firstFileInBranch = false
				} else {
					haveTable.AddRow("", file)
				}
			}
		}
		wantTable := datatable.FromGherkin(godogTable)
		diff, errCnt := haveTable.EqualDataTable(wantTable)
		if errCnt > 0 {
			fmt.Println(diff)
			return fmt.Errorf("found %d differences", errCnt)
		}
		return nil
	})

	sc.Step(`^the commits$`, func(ctx context.Context, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		commits := testgit.FromGherkinTable(table)
		state.fixture.CreateCommits(commits)
	})

	sc.Step(`^the committed configuration file:$`, func(ctx context.Context, content *godog.DocString) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.CreateFile(configfile.FileName, content.Content)
		devRepo.StageFiles(configfile.FileName)
		devRepo.CommitStagedChanges(commands.ConfigFileCommitMessage)
		devRepo.PushBranch()
	})

	sc.Step(`^the configuration file:$`, func(ctx context.Context, content *godog.DocString) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.CreateFile(configfile.FileName, content.Content)
	})

	sc.Step(`^the configuration file is (?:now|still):$`, func(ctx context.Context, content *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have, err := devRepo.FileContentErr(configfile.FileName)
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

	sc.Step(`^the coworker adds this commit to their current branch:$`, func(ctx context.Context, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		commits := testgit.FromGherkinTable(table)
		commit := commits[0]
		coworkerRepo := state.fixture.CoworkerRepo.GetOrPanic()
		coworkerRepo.CreateFile(commit.FileName, commit.FileContent)
		coworkerRepo.StageFiles(commit.FileName)
		coworkerRepo.CommitStagedChanges(commit.Message)
	})

	sc.Step(`^the coworker fetches updates$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.CoworkerRepo.GetOrPanic().Fetch()
	})

	sc.Step(`^the coworker is on the "([^"]*)" branch$`, func(ctx context.Context, branch string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.CoworkerRepo.GetOrPanic().CheckoutBranch(gitdomain.NewLocalBranchName(branch))
	})

	sc.Step(`^the coworker pushes a new "([^"]+)" branch with these commits$`, func(ctx context.Context, branchName string, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		branch := gitdomain.NewLocalBranchName(branchName)
		coworkerRepo := state.fixture.CoworkerRepo.GetOrPanic()
		coworkerRepo.CreateBranch(branch, "main")
		coworkerRepo.CheckoutBranch(branch)
		for _, commit := range testgit.FromGherkinTable(table) {
			coworkerRepo.CreateFile(commit.FileName, commit.FileContent)
			coworkerRepo.StageFiles(commit.FileName)
			coworkerRepo.CommitStagedChanges(commit.Message)
		}
		coworkerRepo.PushBranchToRemote(branch, gitdomain.RemoteOrigin)
	})

	sc.Step(`^the coworker pushes these commits to the "([^"]+)" branch$`, func(ctx context.Context, branchName string, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		branch := gitdomain.NewLocalBranchName(branchName)
		coworkerRepo := state.fixture.CoworkerRepo.GetOrPanic()
		coworkerRepo.CheckoutBranch(branch)
		for _, commit := range testgit.FromGherkinTable(table) {
			coworkerRepo.CreateFile(commit.FileName, commit.FileContent)
			coworkerRepo.StageFiles(commit.FileName)
			coworkerRepo.CommitStagedChanges(commit.Message)
		}
		coworkerRepo.PushBranch()
	})

	sc.Step(`^the coworker resolves the conflict in "([^"]*)"(?: with "([^"]*)")?$`, func(ctx context.Context, filename, content string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		coworkerRepo := state.fixture.CoworkerRepo.GetOrPanic()
		coworkerRepo.CreateFile(filename, content)
		coworkerRepo.StageFiles(filename)
	})

	sc.Step(`^the coworker runs "([^"]+)"$`, func(ctx context.Context, command string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		output, exitCode := state.fixture.CoworkerRepo.GetOrPanic().MustQueryStringCode(command)
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
	})

	sc.Step(`^the coworker runs "([^"]*)" and closes the editor$`, func(ctx context.Context, cmd string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		env := append(os.Environ(), "GIT_EDITOR=true")
		output, exitCode := state.fixture.CoworkerRepo.GetOrPanic().MustQueryStringCodeWith(cmd, &subshell.Options{Env: env})
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
	})

	sc.Step(`^the coworker sets the parent branch of "([^"]*)" as "([^"]*)"$`, func(ctx context.Context, childBranch, parentBranch string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		_ = state.fixture.CoworkerRepo.GetOrPanic().Config.NormalConfig.SetParent(gitdomain.NewLocalBranchName(childBranch), gitdomain.NewLocalBranchName(parentBranch))
	})

	sc.Step(`^the coworker sets the "sync-feature-strategy" to "(merge|rebase)"$`, func(ctx context.Context, value string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		syncFeatureStrategy := asserts.NoError1(configdomain.ParseSyncFeatureStrategy(value))
		_ = state.fixture.CoworkerRepo.GetOrPanic().Config.NormalConfig.SetSyncFeatureStrategy(syncFeatureStrategy.GetOrPanic())
	})

	sc.Step(`^the coworkers workspace now contains file "([^"]*)" with content "([^"]*)"$`, func(ctx context.Context, file, expectedContent string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		actualContent := state.fixture.CoworkerRepo.GetOrPanic().FileContent(file)
		if expectedContent != actualContent {
			return fmt.Errorf("file content does not match\n\nEXPECTED: %q\n\nACTUAL:\n\n%q\n----------------------------", expectedContent, actualContent)
		}
		return nil
	})

	sc.Step(`^the current branch in the other worktree is (?:now|still) "([^"]*)"$`, func(ctx context.Context, expected string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		secondWorkTree := state.fixture.SecondWorktree.GetOrPanic()
		secondWorkTree.Git.CurrentBranchCache.Invalidate()
		actual, err := secondWorkTree.Git.CurrentBranch(secondWorkTree)
		if err != nil {
			return fmt.Errorf("cannot determine current branch of second worktree: %w", err)
		}
		if actual.String() != expected {
			return fmt.Errorf("expected active branch %q but is %q", expected, actual)
		}
		return nil
	})

	sc.Step(`^the current branch is "([^"]*)"$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		state.initialCurrentBranch = Some(branch)
		if !devRepo.Git.BranchExists(devRepo.TestRunner, branch) {
			return fmt.Errorf("cannot check out non-existing branch: %q", branch)
		}
		devRepo.CheckoutBranch(branch)
		return nil
	})

	sc.Step(`^the current branch is "([^"]*)" and the previous branch is "([^"]*)"$`, func(ctx context.Context, currentText, previousText string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		current := gitdomain.NewLocalBranchName(currentText)
		previous := gitdomain.NewLocalBranchName(previousText)
		state.initialCurrentBranch = Some(current)
		devRepo.CheckoutBranch(previous)
		devRepo.CheckoutBranch(current)
	})

	sc.Step(`^the current branch is (?:now|still) "([^"]*)"$`, func(ctx context.Context, expected string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.Git.CurrentBranchCache.Invalidate()
		actual, err := devRepo.Git.CurrentBranch(devRepo.TestRunner)
		if err != nil {
			return fmt.Errorf("cannot determine current branch of developer repo: %w", err)
		}
		if actual.String() != expected {
			return fmt.Errorf("expected active branch %q but is %q", expected, actual)
		}
		return nil
	})

	sc.Step(`^the currently checked out commit is "([^"]+)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have, err := devRepo.Git.CommitMessage(devRepo, "HEAD")
		if err != nil {
			return fmt.Errorf("cannot determine current commit: %w", err)
		}
		if have.String() != want {
			return fmt.Errorf("expected commit %q but got %q", want, have)
		}
		return nil
	})

	sc.Step(`^the home directory contains file "([^"]+)" with content$`, func(ctx context.Context, filename string, docString *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		filePath := filepath.Join(devRepo.HomeDir, filename)
		//nolint:gosec // need permission 700 here in order for tests to work
		return os.WriteFile(filePath, []byte(docString.Content), 0o700)
	})

	sc.Step(`^the initial branches and lineage exist now$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		// verify initial branches
		currentBranches := state.fixture.Branches()
		initialBranches := state.initialBranches.GetOrPanic()
		// fmt.Printf("\nINITIAL:\n%s\n", initialBranches.String())
		// fmt.Printf("NOW:\n%s\n", currentBranches.String())
		diff, errorCount := currentBranches.EqualDataTable(initialBranches)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing branches\n\n", errorCount)
			fmt.Println(diff)
			return errors.New("mismatching branches found, see diff above")
		}
		// verify initial lineage
		currentLineage := devRepo.LineageTable()
		diff, errCnt := currentLineage.EqualDataTable(state.initialLineage.GetOrPanic())
		if errCnt > 0 {
			fmt.Printf("\nERROR! Found %d differences in the lineage\n\n", errCnt)
			fmt.Println(diff)
			return errors.New("mismatching lineage found, see the diff above")
		}
		return nil
	})

	sc.Step(`^the initial branches exist now$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		have := state.fixture.Branches()
		want := state.initialBranches.GetOrPanic()
		// fmt.Printf("HAVE:\n%s\n", have.String())
		// fmt.Printf("WANT:\n%s\n", want.String())
		diff, errorCount := have.EqualDataTable(want)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing branches\n\n", errorCount)
			fmt.Println(diff)
			return errors.New("mismatching branches found, see diff above")
		}
		return nil
	})

	sc.Step(`^the initial commits exist now$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		currentCommits := state.fixture.CommitTable(state.initialCommits.GetOrPanic().Cells[0])
		initialCommits := state.initialCommits.GetOrPanic()
		errDiff, errCount := initialCommits.EqualDataTable(currentCommits)
		if errCount == 0 {
			return nil
		}
		fmt.Println(errDiff)
		return errors.New("current commits are not the same as the initial commits")
	})

	sc.Step(`^the initial lineage exists now$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.LineageTable()
		diff, errCnt := have.EqualDataTable(state.initialLineage.GetOrPanic())
		if errCnt > 0 {
			fmt.Printf("\nERROR! Found %d differences in the lineage\n\n", errCnt)
			fmt.Printf("INITIAL LINEAGE:\n%s\n", state.initialLineage.String())
			fmt.Printf("CURRENT LINEAGE:\n%s\n", have.String())
			fmt.Println(diff)
			return errors.New("mismatching branches found, see the diff above")
		}
		return nil
	})

	sc.Step(`^the initial tags exist now$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		currentTags := state.fixture.TagTable()
		initialTags := state.initialTags.GetOrPanic()
		errDiff, errCount := initialTags.EqualDataTable(currentTags)
		if errCount == 0 {
			return nil
		}
		fmt.Println(errDiff)
		return errors.New("current tags are not the same as the initial commits")
	})

	sc.Step(`^the main branch is "([^"]+)"$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		return devRepo.Config.SetMainBranch(gitdomain.NewLocalBranchName(name))
	})

	sc.Step(`^the main branch is (?:now|still) "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.UnvalidatedConfig.MainBranch
		if have.String() != want {
			return fmt.Errorf("expected %q, got %q", want, have)
		}
		return nil
	})

	sc.Step(`^the main branch is (?:now|still) not set$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.NormalConfig.GitConfig.MainBranch
		if branch, has := have.Get(); has {
			return fmt.Errorf("unexpected main branch setting %q", branch)
		}
		return nil
	})

	sc.Step(`^the origin is "([^"]*)"$`, func(ctx context.Context, origin string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.SetTestOrigin(origin)
	})

	sc.Step(`^the perennial branches are "([^"]+)"$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		return devRepo.Config.NormalConfig.SetPerennialBranches(gitdomain.NewLocalBranchNames(name))
	})

	sc.Step(`^the perennial branches are (?:now|still) "([^"]+)"$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actual := devRepo.Config.NormalConfig.PartialBranchesOfType(configdomain.BranchTypePerennialBranch)
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 perennial branch, got %q", actual)
		}
		if (actual)[0].String() != name {
			return fmt.Errorf("expected %q, got %q", name, (actual)[0])
		}
		return nil
	})

	sc.Step(`^the previous Git branch is (?:now|still) "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Git.PreviouslyCheckedOutBranch(devRepo.TestRunner)
		if have.String() != want {
			return fmt.Errorf("expected previous branch %q but got %q", want, have)
		}
		return nil
	})

	sc.Step(`^there are (?:now|still) no perennial branches$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branches := devRepo.Config.NormalConfig.GitConfig.PerennialBranches
		if len(branches) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", branches)
		}
		return nil
	})

	sc.Step(`^these branches exist now$`, func(ctx context.Context, input *godog.Table) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
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

	sc.Step(`^these commits exist now$`, func(ctx context.Context, table *godog.Table) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		return state.compareGherkinTable(table)
	})

	sc.Step(`^these committed files exist now$`, func(ctx context.Context, table *godog.Table) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		fileTable := devRepo.FilesInBranches("main")
		diff, errorCount := fileTable.EqualGherkin(table)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing files\n\n", errorCount)
			fmt.Println(diff)
			return errors.New("mismatching files found, see diff above")
		}
		return nil
	})

	sc.Step(`^these tags exist now$`, func(ctx context.Context, table *godog.Table) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		tagTable := state.fixture.TagTable()
		diff, errorCount := tagTable.EqualGherkin(table)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing tags\n\n", errorCount)
			fmt.Println(diff)
			return errors.New("mismatching tags found, see diff above")
		}
		return nil
	})

	sc.Step(`^the tags$`, func(ctx context.Context, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.CreateTags(table)
	})

	sc.Step(`^the uncommitted file has content:$`, func(ctx context.Context, content *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		msg := devRepo.HasFile(
			state.uncommittedFileName.GetOrPanic(),
			content.Content,
		)
		if len(msg) > 0 {
			return errors.New(msg)
		}
		return nil
	})

	sc.Step(`^the uncommitted file still exists$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		msg := devRepo.HasFile(
			state.uncommittedFileName.GetOrPanic(),
			state.uncommittedContent.GetOrPanic(),
		)
		if len(msg) > 0 {
			return errors.New(msg)
		}
		return nil
	})

	sc.Step(`^this lineage exists now$`, func(ctx context.Context, input *godog.Table) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		table := devRepo.LineageTable()
		diff, errCount := table.EqualGherkin(input)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the lineage\n\n", errCount)
			fmt.Println(diff)
			return errors.New("mismatching branches found, see the diff above")
		}
		return nil
	})

	sc.Step(`^tool "([^"]*)" is broken$`, func(ctx context.Context, name string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.MockBrokenCommand(name)
		state.browserVariable = Some(name)
	})

	sc.Step(`^tool "([^"]*)" is installed$`, func(ctx context.Context, tool string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.MockCommand(tool)
		state.browserVariable = Some(tool)
	})

	// This step exists to avoid re-creating commits with the same SHA as existing commits
	// because that can cause flaky tests.
	sc.Step(`wait 1 second to ensure new Git timestamps`, func() {
		time.Sleep(1 * time.Second)
	})
}

func runCommand(ctx context.Context, command string) {
	state := ctx.Value(keyScenarioState).(*ScenarioState)
	devRepo, hasDevRepo := state.fixture.DevRepo.Get()
	if hasDevRepo {
		state.CaptureState()
		updateInitialSHAs(state)
	}
	var exitCode int
	var runOutput string
	env := os.Environ()
	if browserVariable, hasBrowserOverride := state.browserVariable.Get(); hasBrowserOverride {
		env = envvars.Replace(env, browser.EnvVarName, browserVariable)
	}
	if hasDevRepo {
		runOutput, exitCode = devRepo.MustQueryStringCodeWith(command, &subshell.Options{
			Env: env,
		})
		devRepo.Reload()
	} else {
		parts := asserts.NoError1(shellquote.Split(command))
		cmd, args := parts[0], parts[1:]
		subProcess := exec.Command(cmd, args...) // #nosec
		subProcess.Dir = state.fixture.Dir
		outputBytes, _ := subProcess.CombinedOutput()
		runOutput = string(outputBytes)
		exitCode = subProcess.ProcessState.ExitCode()
	}
	state.runOutput = Some(runOutput)
	state.runExitCode = Some(exitCode)
}

func updateInitialSHAs(state *ScenarioState) {
	devRepo := state.fixture.DevRepo.GetOrPanic()
	devSHAs := devRepo.CommitSHAs()
	if state.initialDevSHAs.IsNone() && state.insideGitRepo {
		state.initialDevSHAs = Some(devSHAs)
	}
	state.beforeRunDevSHAs = Some(devSHAs)
	if originRepo, hasOriginrepo := state.fixture.OriginRepo.Get(); hasOriginrepo && state.insideGitRepo {
		originSHAs := originRepo.CommitSHAs()
		if state.initialOriginSHAs.IsNone() {
			state.initialOriginSHAs = Some(originSHAs)
		}
		state.beforeRunOriginSHAs = Some(originSHAs)
	}
	if secondWorkTree, hasSecondWorkTree := state.fixture.SecondWorktree.Get(); hasSecondWorkTree && state.insideGitRepo {
		workTreeSHAs := secondWorkTree.CommitSHAs()
		if state.initialWorktreeSHAs.IsNone() {
			state.initialWorktreeSHAs = Some(workTreeSHAs)
		}
	}
}
