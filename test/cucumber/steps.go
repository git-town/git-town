package cucumber

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/cucumber/godog"
	cukemessages "github.com/cucumber/messages/go/v21"
	"github.com/git-town/git-town/v16/internal/cli/dialog/components"
	"github.com/git-town/git-town/v16/internal/cli/print"
	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/config/configfile"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/internal/hosting/hostingdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/git-town/git-town/v16/test/asserts"
	"github.com/git-town/git-town/v16/test/commands"
	"github.com/git-town/git-town/v16/test/datatable"
	"github.com/git-town/git-town/v16/test/filesystem"
	"github.com/git-town/git-town/v16/test/fixture"
	"github.com/git-town/git-town/v16/test/git"
	"github.com/git-town/git-town/v16/test/helpers"
	"github.com/git-town/git-town/v16/test/output"
	"github.com/git-town/git-town/v16/test/subshell"
	"github.com/google/go-cmp/cmp"
	"github.com/kballard/go-shellquote"
)

// the global FixtureFactory instance.
var fixtureFactory *fixture.Factory //nolint:gochecknoglobals

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
	sc.Step(`^a coworker clones the repository$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.AddCoworkerRepo()
	})

	sc.Step(`^a folder "([^"]*)"$`, func(ctx context.Context, name string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.CreateFolder(name)
	})

	sc.Step(`a Git repo with origin`, func(ctx context.Context) (context.Context, error) {
		scenarioName := ctx.Value(keyScenarioName).(string)
		scenarioTags := ctx.Value(keyScenarioTags).([]*cukemessages.PickleTag)
		fixture := fixtureFactory.CreateFixture(scenarioName)
		if helpers.HasTag(scenarioTags, "@debug") {
			fixture.DevRepo.GetOrPanic().Verbose = true
			fixture.OriginRepo.GetOrPanic().Verbose = true
		}
		state := ScenarioState{
			fixture:              fixture,
			initialBranches:      None[datatable.DataTable](),
			initialCommits:       None[datatable.DataTable](),
			initialCurrentBranch: None[gitdomain.LocalBranchName](),
			initialDevSHAs:       None[map[string]gitdomain.SHA](),
			initialLineage:       None[datatable.DataTable](),
			initialOriginSHAs:    None[map[string]gitdomain.SHA](),
			initialTags:          None[datatable.DataTable](),
			initialWorktreeSHAs:  None[map[string]gitdomain.SHA](),
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
		scenarioTags := ctx.Value(keyScenarioTags).([]*cukemessages.PickleTag)
		fixture := fixtureFactory.CreateFixture(scenarioName)
		devRepo := fixture.DevRepo.GetOrPanic()
		if helpers.HasTag(scenarioTags, "@debug") {
			devRepo.Verbose = true
		}
		devRepo.RemoveRemote(gitdomain.RemoteOrigin)
		fixture.OriginRepo = MutableNone[commands.TestCommands]()
		state := ScenarioState{
			fixture:              fixture,
			initialBranches:      None[datatable.DataTable](),
			initialCommits:       None[datatable.DataTable](),
			initialCurrentBranch: None[gitdomain.LocalBranchName](),
			initialDevSHAs:       None[map[string]gitdomain.SHA](),
			initialLineage:       None[datatable.DataTable](),
			initialOriginSHAs:    None[map[string]gitdomain.SHA](),
			initialTags:          None[datatable.DataTable](),
			initialWorktreeSHAs:  None[map[string]gitdomain.SHA](),
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
		if !devRepo.HasMergeInProgress(devRepo.TestRunner) {
			return errors.New("expected merge in progress")
		}
		return nil
	})

	sc.Step(`^an additional "([^"]+)" remote with URL "([^"]+)"$`, func(ctx context.Context, remote, url string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.AddRemote(gitdomain.NewRemote(remote), url)
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

	sc.Step(`^an upstream repo$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.AddUpstream()
	})

	sc.Step(`^a proposal for this branch does not exist`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.TestRunner.ProposalOverride = Some(hostingdomain.OverrideNoProposal)
	})

	sc.Step(`^a proposal for this branch exists at "([^"]+)"`, func(ctx context.Context, url string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.TestRunner.ProposalOverride = Some(url)
	})

	sc.Step(`^a rebase is (?:now|still) in progress$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		repoStatus, err := devRepo.RepoStatus(devRepo.TestRunner)
		asserts.NoError(err)
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

	sc.Step(`^branch "([^"]+)" is (?:now|still) a contribution branch`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		if !slices.Contains(devRepo.Config.NormalConfig.ContributionBranches, branch) {
			return fmt.Errorf(
				"branch %q isn't contribution as expected.\nContribution branches: %s",
				branch,
				strings.Join(devRepo.Config.NormalConfig.ContributionBranches.Strings(), ", "),
			)
		}
		return nil
	})

	sc.Step(`^branch "([^"]+)" is (?:now|still) a feature branch`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		branchType := devRepo.Config.BranchType(branch)
		if branchType != configdomain.BranchTypeFeatureBranch {
			return fmt.Errorf("branch %q is a %s", branch, branchType)
		}
		return nil
	})

	sc.Step(`^branch "([^"]+)" is (?:now|still) observed`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		if !slices.Contains(devRepo.Config.NormalConfig.ObservedBranches, branch) {
			return fmt.Errorf(
				"branch %q isn't observed as expected.\nObserved branches: %s",
				branch,
				strings.Join(devRepo.Config.NormalConfig.ObservedBranches.Strings(), ", "),
			)
		}
		return nil
	})

	sc.Step(`^branch "([^"]+)" is (?:now|still) parked`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		if !slices.Contains(devRepo.Config.NormalConfig.ParkedBranches, branch) {
			return fmt.Errorf(
				"branch %q isn't parked as expected.\nParked branches: %s",
				branch,
				strings.Join(devRepo.Config.NormalConfig.ParkedBranches.Strings(), ", "),
			)
		}
		return nil
	})

	sc.Step(`^branch "([^"]+)" is (?:now|still) perennial`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		if !devRepo.Config.NormalConfig.IsPerennialBranch(branch) {
			return fmt.Errorf(
				"branch %q isn't perennial as expected.\nPerennial branches: %s",
				branch,
				strings.Join(devRepo.Config.NormalConfig.PerennialBranches.Strings(), ", "),
			)
		}
		return nil
	})

	sc.Step(`^branch "([^"]+)" is (?:now|still) prototype`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		if !slices.Contains(devRepo.Config.NormalConfig.PrototypeBranches, branch) {
			return fmt.Errorf(
				"branch %q isn't prototype as expected.\nPrototype branches: %s",
				branch,
				devRepo.Config.NormalConfig.PrototypeBranches.Join(", "),
			)
		}
		return nil
	})

	sc.Step(`^custom global Git setting "alias\.(.*?)" is "([^"]*)"$`, func(ctx context.Context, name, value string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		aliasableCommand := configdomain.AliasableCommand(name)
		return devRepo.SetGitAlias(aliasableCommand, value)
	})

	sc.Step(`^custom global Git setting "alias\.(.*?)" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, name, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		aliasableCommand := configdomain.AliasableCommand(name)
		have, err := devRepo.LoadGitAlias(aliasableCommand)
		asserts.NoError(err)
		if have != want {
			return fmt.Errorf("unexpected value for key %q: want %q have %q", name, want, have)
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

	sc.Step(`^file "([^"]*)" (?:now|still) has content "([^"]*)"$`, func(ctx context.Context, file, expectedContent string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actualContent := devRepo.FileContent(file)
		if expectedContent != actualContent {
			return fmt.Errorf("file content does not match\n\nEXPECTED: %q\n\nACTUAL:\n\n%q\n----------------------------", expectedContent, actualContent)
		}
		return nil
	})

	sc.Step(`^file "([^"]+)" still contains unresolved conflicts$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		content := devRepo.FileContent(name)
		if !strings.Contains(content, "<<<<<<<") {
			return fmt.Errorf("file %q does not contain unresolved conflicts", name)
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
		err := devRepo.RemovePerennialBranchConfiguration()
		asserts.NoError(err)
		devRepo.RemoveMainBranchConfiguration()
	})

	sc.Step(`^Git Town parent setting for branch "([^"]*)" is "([^"]*)"$`, func(ctx context.Context, branch, value string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branchName := gitdomain.NewLocalBranchName(branch)
		configKey := configdomain.NewParentKey(branchName)
		return devRepo.Config.NormalConfig.GitConfig.SetConfigValue(configdomain.ConfigScopeLocal, configKey, value)
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
			fmt.Println("EXPECTED OUTPUT START ==============================================")
			fmt.Println("==================================================================")
			fmt.Println()
			fmt.Println(expected.Content)
			fmt.Println()
			fmt.Println("==================================================================")
			fmt.Println("EXPECTED OUTPUT END ================================================")
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

	sc.Step(`^Git Town prints an error like:$`, func(ctx context.Context, expected *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.runExitCodeChecked = true
		regex := regexp.MustCompile(expected.Content)
		have := stripansi.Strip(state.runOutput.GetOrPanic())
		if !regex.MatchString(have) {
			return fmt.Errorf("text not found:\n%s\n\nactual text:\n%s", expected.Content, state.runOutput.GetOrDefault())
		}
		if exitCode := state.runExitCode.GetOrPanic(); exitCode == 0 {
			return fmt.Errorf("unexpected exit code %d", exitCode)
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
		expanded := dataTable.Expand(
			devRepo,
			state.fixture.OriginRepo.Value,
			state.fixture.SecondWorktree.Value,
			state.initialDevSHAs.GetOrPanic(),
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

	sc.Step(`^Git Town runs without errors`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		exitCode := state.runExitCode.GetOrPanic()
		if exitCode != 0 {
			return errors.New("unexpected failure of scenario")
		}
		return nil
	})

	sc.Step(`^global Git setting "alias\.(.*?)" is "([^"]*)"$`, func(ctx context.Context, name, value string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		key, hasKey := configdomain.ParseKey(configdomain.AliasKeyPrefix + name).Get()
		if !hasKey {
			return fmt.Errorf("no key found for %q", name)
		}
		return devRepo.SetGitAlias(configdomain.NewAliasKey(key).GetOrPanic().AliasableCommand(), value)
	})

	sc.Step(`^global Git setting "alias\.(.*?)" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, name, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		key, hasKey := configdomain.ParseKey(configdomain.AliasKeyPrefix + name).Get()
		if !hasKey {
			return errors.New("key not found")
		}
		aliasableCommand := configdomain.NewAliasKey(key).GetOrPanic().AliasableCommand()
		have := devRepo.Config.NormalConfig.Aliases[aliasableCommand]
		if have != want {
			return fmt.Errorf("unexpected value for key %q: want %q have %q", name, want, have)
		}
		return nil
	})

	sc.Step(`^global Git setting "alias\.(.*?)" (?:now|still) doesn't exist$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		key, hasKey := configdomain.ParseKey(configdomain.AliasKeyPrefix + name).Get()
		if !hasKey {
			return errors.New("key not found")
		}
		aliasableCommand := configdomain.NewAliasKey(key).GetOrPanic().AliasableCommand()
		command, has := devRepo.Config.NormalConfig.Aliases[aliasableCommand]
		if has {
			return fmt.Errorf("unexpected aliasableCommand %q: %q", key, command)
		}
		return nil
	})

	sc.Step(`^global Git setting "([^"]+)" is "([^"]+)"$`, func(ctx context.Context, name, value string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		return devRepo.Config.NormalConfig.GitConfig.SetConfigValue(configdomain.ConfigScopeGlobal, configdomain.Key(name), value)
	})

	sc.Step(`^(global |local |)Git Town setting "([^"]+)" is "([^"]+)"$`, func(ctx context.Context, locality, name, value string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		key, hasKey := configdomain.ParseKey("git-town." + name).Get()
		if !hasKey {
			return fmt.Errorf("unknown config key: %q", name)
		}
		scope := configdomain.ParseConfigScope(strings.TrimSpace(locality))
		return devRepo.Config.NormalConfig.GitConfig.SetConfigValue(scope, key, value)
	})

	sc.Step(`^(global |local |)Git Town setting "([^"]+)" is (?:now|still) "([^"]+)"$`, func(ctx context.Context, locality, name, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		key, hasKey := configdomain.ParseKey("git-town." + name).Get()
		if !hasKey {
			return fmt.Errorf("unknown config key: %q", name)
		}
		scope := configdomain.ParseConfigScope(strings.TrimSpace(locality))
		have, has := devRepo.GitConfig(scope, key).Get()
		if !has {
			return fmt.Errorf(`expected %s setting %q to be %q but doesn't exist`, scope, name, want)
		}
		if have != want {
			return fmt.Errorf(`expected %s setting %q to be %q but is %q`, scope, name, want, have)
		}
		return nil
	})

	sc.Step(`^(global |local |)Git Town setting "([^"]+)" (:?now|still) doesn't exist$`, func(ctx context.Context, locality, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		key, hasKey := configdomain.ParseKey("git-town." + name).Get()
		if !hasKey {
			return fmt.Errorf("unknown config key: %q", name)
		}
		scope := configdomain.ParseConfigScope(strings.TrimSpace(locality))
		if value, hasValue := devRepo.GitConfig(scope, key).Get(); hasValue {
			return fmt.Errorf("should not have %s setting %q anymore but it exists and has value %q", scope, name, value)
		}
		return nil
	})

	sc.Step(`^I add commit "([^"]*)" to the "([^"]*)" branch$`, func(ctx context.Context, message, branch string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.CreateCommit(git.Commit{
			Branch:   gitdomain.NewLocalBranchName(branch),
			FileName: "new_file",
			Message:  gitdomain.CommitMessage(message),
		})
	})

	sc.Step(`^I add this commit to the "([^"]*)" branch$`, func(ctx context.Context, branch string, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		commit := git.FromGherkinTable(table)[0]
		commit.Branch = gitdomain.LocalBranchName(branch)
		devRepo.CreateCommit(commit)
	})

	sc.Step(`^I add this commit to the current branch:$`, func(ctx context.Context, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		commit := git.FromGherkinTable(table)[0]
		devRepo.CreateFile(commit.FileName, commit.FileContent)
		devRepo.StageFiles(commit.FileName)
		devRepo.CommitStagedChanges(commit.Message)
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
		// scenarioTags := ctx.Value(keyScenarioTags).([]*cukemessages.PickleTag)
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
			fixture:              fixture,
			initialBranches:      None[datatable.DataTable](),
			initialCommits:       None[datatable.DataTable](),
			initialCurrentBranch: None[gitdomain.LocalBranchName](),
			initialDevSHAs:       None[map[string]gitdomain.SHA](),
			initialLineage:       None[datatable.DataTable](),
			initialOriginSHAs:    None[map[string]gitdomain.SHA](),
			initialTags:          None[datatable.DataTable](),
			initialWorktreeSHAs:  None[map[string]gitdomain.SHA](),
			insideGitRepo:        true,
			runExitCode:          None[int](),
			runExitCodeChecked:   false,
			runOutput:            None[string](),
			uncommittedContent:   None[string](),
			uncommittedFileName:  None[string](),
		}
		return context.WithValue(ctx, keyScenarioState, &state), nil
	})

	sc.Step(`^I delete the local "([^"]+)" branch manually$`, func(ctx context.Context, branchName string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.DeleteLocalBranch(devRepo, gitdomain.LocalBranchName(branchName))
	})

	sc.Step(`^in a separate terminal I create branch "([^"]+)" with commits$`, func(ctx context.Context, branchName string, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		existingBranch := devRepo.CurrentBranchCache.Value()
		newBranch := gitdomain.NewLocalBranchName(branchName)
		devRepo.CreateBranch(newBranch, "main")
		devRepo.CheckoutBranch(newBranch)
		for _, commit := range git.FromGherkinTable(table) {
			devRepo.CreateFile(commit.FileName, commit.FileContent)
			devRepo.StageFiles(commit.FileName)
			devRepo.CommitStagedChanges(commit.Message)
		}
		devRepo.CheckoutBranch(existingBranch)
	})

	sc.Step(`^inspect the commits$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		fmt.Println("DEV")
		output, err := devRepo.Query("git", "branch", "-vva")
		fmt.Println(output)
		return err
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
		output, exitCode := devRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Env: env, Input: Some(input.Content)})
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Config.Reload()
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

	sc.Step(`^I run "([^"]*)" and close the editor$`, func(ctx context.Context, cmd string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		env := append(os.Environ(), "GIT_EDITOR=true")
		var exitCode int
		var output string
		output, exitCode = devRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Env: env})
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Config.Reload()
	})

	sc.Step(`^I run "([^"]*)" and enter an empty commit message$`, func(ctx context.Context, cmd string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		devRepo.MockCommitMessage("")
		var exitCode int
		var output string
		output, exitCode = devRepo.MustQueryStringCode(cmd)
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Config.Reload()
	})

	sc.Step(`^I run "([^"]*)" and enter "([^"]*)" for the commit message$`, func(ctx context.Context, cmd, message string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		devRepo.MockCommitMessage(message)
		var exitCode int
		var output string
		output, exitCode = devRepo.MustQueryStringCode(cmd)
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Config.Reload()
	})

	sc.Step(`^I run "([^"]+)" in the "([^"]+)" folder$`, func(ctx context.Context, cmd, folderName string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		var exitCode int
		var output string
		output, exitCode = devRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Dir: folderName})
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Config.Reload()
	})

	sc.Step(`^I run "([^"]+)" in the other worktree$`, func(ctx context.Context, cmd string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.CaptureState()
		updateInitialSHAs(state)
		secondWorkTree := state.fixture.SecondWorktree.GetOrPanic()
		var exitCode int
		var output string
		output, exitCode = secondWorkTree.MustQueryStringCode(cmd)
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		secondWorkTree.Config.Reload()
	})

	sc.Step(`^I run "([^"]*)" in the other worktree and enter "([^"]*)" for the commit message$`, func(ctx context.Context, cmd, message string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.CaptureState()
		updateInitialSHAs(state)
		secondWorkTree := state.fixture.SecondWorktree.GetOrPanic()
		secondWorkTree.MockCommitMessage(message)
		var exitCode int
		var output string
		output, exitCode = secondWorkTree.MustQueryStringCode(cmd)
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		secondWorkTree.Config.Reload()
	})

	sc.Step(`^I (?:run|ran) "(.+)"$`, func(ctx context.Context, command string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo, hasDevRepo := state.fixture.DevRepo.Get()
		if hasDevRepo {
			state.CaptureState()
			updateInitialSHAs(state)
		}
		var exitCode int
		var runOutput string
		if hasDevRepo {
			runOutput, exitCode = devRepo.MustQueryStringCode(command)
			devRepo.Config.Reload()
		} else {
			parts, err := shellquote.Split(command)
			asserts.NoError(err)
			cmd, args := parts[0], parts[1:]
			subProcess := exec.Command(cmd, args...) // #nosec
			subProcess.Dir = state.fixture.Dir
			subProcess.Env = append(subProcess.Environ(), "LC_ALL=C")
			outputBytes, _ := subProcess.CombinedOutput()
			runOutput = string(outputBytes)
			exitCode = subProcess.ProcessState.ExitCode()
		}
		state.runOutput = Some(runOutput)
		state.runExitCode = Some(exitCode)
	})

	sc.Step(`^I (?:run|ran) "([^"]+)" and enter into the dialogs?:$`, func(ctx context.Context, cmd string, input *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		env := os.Environ()
		answers, err := helpers.TableToInputEnv(input)
		asserts.NoError(err)
		for dialogNumber, answer := range answers {
			env = append(env, fmt.Sprintf("%s_%02d=%s", components.TestInputKey, dialogNumber, answer))
		}
		var exitCode int
		var output string
		output, exitCode = devRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Env: env})
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
		devRepo.Config.Reload()
	})

	sc.Step(`^"([^"]*)" launches a new proposal with this url in my browser:$`, func(ctx context.Context, tool string, url *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		want := fmt.Sprintf("%s called with: %s", tool, url.Content)
		want = strings.ReplaceAll(want, "?", `\?`)
		regex := regexp.MustCompile(want)
		have := state.runOutput.GetOrPanic()
		if !regex.MatchString(have) {
			return fmt.Errorf("EXPECTED: a regex matching %q\nGOT: %q", want, have)
		}
		return nil
	})

	sc.Step(`^local Git setting "init.defaultbranch" is "([^"]*)"$`, func(ctx context.Context, value string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.SetDefaultGitBranch(gitdomain.NewLocalBranchName(value))
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
		if devRepo.Config.NormalConfig.ContainsLineage() {
			lineage := devRepo.Config.NormalConfig.Lineage
			return fmt.Errorf("unexpected Git Town lineage information: %+v", lineage)
		}
		return nil
	})

	sc.Step(`^no merge is in progress$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		if devRepo.HasMergeInProgress(devRepo.TestRunner) {
			return errors.New("expected no merge in progress")
		}
		return nil
	})

	sc.Step(`^no rebase is in progress$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		repoStatus, err := devRepo.RepoStatus(devRepo.TestRunner)
		asserts.NoError(err)
		if repoStatus.RebaseInProgress {
			return errors.New("expected no rebase in progress")
		}
		return nil
	})

	sc.Step(`^no tool to open browsers is installed$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.MockNoCommandsInstalled()
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
		commitMessage, err := originRepo.FirstCommitMessageInBranch(originRepo.TestRunner, branchToShip.BranchName(), "main")
		asserts.NoError(err)
		if commitMessage.IsNone() {
			return errors.New("branch to ship contains no commits")
		}
		originRepo.CheckoutBranch("main")
		err = originRepo.SquashMerge(originRepo.TestRunner, branchToShip)
		asserts.NoError(err)
		originRepo.StageFiles("-A")
		err = originRepo.Commit(originRepo.TestRunner, commitMessage, false, gitdomain.NewAuthorOpt("CI <ci@acme.com>"))
		asserts.NoError(err)
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
				branchSetup.Locations.Is(git.LocationLocal),
				branchSetup.Locations.Is(git.LocationLocal, git.LocationOrigin):
				repoToCreateBranchIn = state.fixture.DevRepo.GetOrPanic()
			case branchSetup.Locations.Is(git.LocationOrigin):
				repoToCreateBranchIn = state.fixture.OriginRepo.GetOrPanic()
			case branchSetup.Locations.Is(git.LocationUpstream):
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
					repoToCreateBranchIn.CreateChildFeatureBranch(branchSetup.Name, branchSetup.Parent.GetOrPanic())
				case configdomain.BranchTypePerennialBranch:
					repoToCreateBranchIn.CreatePerennialBranch(branchSetup.Name)
				case configdomain.BranchTypeContributionBranch:
					repoToCreateBranchIn.CreateContributionBranch(branchSetup.Name)
				case configdomain.BranchTypeObservedBranch:
					repoToCreateBranchIn.CreateObservedBranch(branchSetup.Name)
				case configdomain.BranchTypeParkedBranch:
					repoToCreateBranchIn.CreateParkedBranch(branchSetup.Name, branchSetup.Parent.GetOrPanic())
				case configdomain.BranchTypePrototypeBranch:
					repoToCreateBranchIn.CreatePrototypeBranch(branchSetup.Name, branchSetup.Parent.GetOrPanic())
				default:
					return errors.New("unhandled branch type: " + branchType.String())
				}
			} else {
				repoToCreateBranchIn.CreateBranch(branchSetup.Name, "main")
				if parent, hasParent := branchSetup.Parent.Get(); hasParent {
					err := repoToCreateBranchIn.Config.NormalConfig.SetParent(branchSetup.Name, parent)
					asserts.NoError(err)
				}
			}
			if len(branchSetup.Locations) > 1 {
				switch {
				case branchSetup.Locations.Is(git.LocationLocal, git.LocationOrigin):
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

	sc.Step(`^the commits$`, func(ctx context.Context, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		// create the commits
		commits := git.FromGherkinTable(table)
		state.fixture.CreateCommits(commits)
		// restore the initial branch
		initialBranch, hasInitialBranch := state.initialCurrentBranch.Get()
		if !hasInitialBranch {
			devRepo.CheckoutBranch(gitdomain.NewLocalBranchName("main"))
			return
		}
		// NOTE: reading the cached value here to keep the test suite fast by avoiding unnecessary disk access
		if devRepo.CurrentBranchCache.Value() != initialBranch {
			devRepo.CheckoutBranch(initialBranch)
			return
		}
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

	sc.Step(`^the contribution branches are (?:now|still) "([^"]+)"$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actual := devRepo.Config.NormalConfig.LocalGitConfig.ContributionBranches
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 contribution branch, got %q", actual)
		}
		if (actual)[0].String() != name {
			return fmt.Errorf("expected %q, got %q", name, (actual)[0])
		}
		return nil
	})

	sc.Step(`^the coworker adds this commit to their current branch:$`, func(ctx context.Context, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		commits := git.FromGherkinTable(table)
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
		for _, commit := range git.FromGherkinTable(table) {
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
		for _, commit := range git.FromGherkinTable(table) {
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
		var exitCode int
		var output string
		output, exitCode = state.fixture.CoworkerRepo.GetOrPanic().MustQueryStringCode(command)
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
	})

	sc.Step(`^the coworker runs "([^"]*)" and closes the editor$`, func(ctx context.Context, cmd string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		env := append(os.Environ(), "GIT_EDITOR=true")
		var exitCode int
		var output string
		output, exitCode = state.fixture.CoworkerRepo.GetOrPanic().MustQueryStringCodeWith(cmd, &subshell.Options{Env: env})
		state.runOutput = Some(output)
		state.runExitCode = Some(exitCode)
	})

	sc.Step(`^the coworker sets the parent branch of "([^"]*)" as "([^"]*)"$`, func(ctx context.Context, childBranch, parentBranch string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		_ = state.fixture.CoworkerRepo.GetOrPanic().Config.NormalConfig.SetParent(gitdomain.NewLocalBranchName(childBranch), gitdomain.NewLocalBranchName(parentBranch))
	})

	sc.Step(`^the coworker sets the "sync-feature-strategy" to "(merge|rebase)"$`, func(ctx context.Context, value string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		syncFeatureStrategy, err := configdomain.ParseSyncFeatureStrategy(value)
		asserts.NoError(err)
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
		secondWorkTree.CurrentBranchCache.Invalidate()
		actual, err := secondWorkTree.CurrentBranch(secondWorkTree)
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
		if !devRepo.BranchExists(devRepo.TestRunner, branch) {
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
		devRepo.CurrentBranchCache.Invalidate()
		actual, err := devRepo.CurrentBranch(devRepo.TestRunner)
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
		have := devRepo.CurrentCommitMessage()
		if have != want {
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
		errDiff, errCount := state.initialCommits.GetOrPanic().EqualDataTable(currentCommits)
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
		errDiff, errCount := state.initialTags.GetOrPanic().EqualDataTable(currentTags)
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
		have := devRepo.Config.NormalConfig.LocalGitConfig.MainBranch
		if branch, has := have.Get(); has {
			return fmt.Errorf("unexpected main branch setting %q", branch)
		}
		return nil
	})

	sc.Step(`^the observed branches are (?:now|still) "([^"]+)"$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actual := devRepo.Config.NormalConfig.LocalGitConfig.ObservedBranches
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 observed branch, got %q", actual)
		}
		if (actual)[0].String() != name {
			return fmt.Errorf("expected %q, got %q", name, (actual)[0])
		}
		return nil
	})

	sc.Step(`^the origin is "([^"]*)"$`, func(ctx context.Context, origin string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.SetTestOrigin(origin)
	})

	sc.Step(`^the parked branches are (?:now|still) "([^"]+)"$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actual := devRepo.Config.NormalConfig.LocalGitConfig.ParkedBranches
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 parked branch, got %q", actual)
		}
		if (actual)[0].String() != name {
			return fmt.Errorf("expected %q, got %q", name, (actual)[0])
		}
		return nil
	})

	sc.Step(`^the perennial branches are "([^"]+)"$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		return devRepo.Config.NormalConfig.SetPerennialBranches(gitdomain.NewLocalBranchNames(name))
	})

	sc.Step(`^the perennial branches are (?:now|still) "([^"]+)"$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actual := devRepo.Config.NormalConfig.LocalGitConfig.PerennialBranches
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 perennial branch, got %q", actual)
		}
		if (actual)[0].String() != name {
			return fmt.Errorf("expected %q, got %q", name, (actual)[0])
		}
		return nil
	})

	sc.Step(`^the perennial branches are (?:now|still) "([^"]+)" and "([^"]+)"$`, func(ctx context.Context, branch1, branch2 string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actual := devRepo.Config.NormalConfig.LocalGitConfig.PerennialBranches
		if len(actual) != 2 {
			return fmt.Errorf("expected 2 perennial branches, got %q", actual)
		}
		if (actual)[0].String() != branch1 || (actual)[1].String() != branch2 {
			return fmt.Errorf("expected %q, got %q", []string{branch1, branch2}, actual)
		}
		return nil
	})

	sc.Step(`^the previous Git branch is (?:now|still) "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Commands.PreviouslyCheckedOutBranch(devRepo.TestRunner)
		if have.String() != want {
			return fmt.Errorf("expected previous branch %q but got %q", want, have)
		}
		return nil
	})

	sc.Step(`^the prototype branches are (?:now|still) "([^"]+)"$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actual := devRepo.Config.NormalConfig.LocalGitConfig.PrototypeBranches
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 prototype branch, got %q", actual)
		}
		if (actual)[0].String() != name {
			return fmt.Errorf("expected %q, got %q", name, (actual)[0])
		}
		return nil
	})

	sc.Step(`^there are (?:now|still) no contribution branches$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branches := devRepo.Config.NormalConfig.LocalGitConfig.ContributionBranches
		if len(branches) > 0 {
			return fmt.Errorf("expected no contribution branches, got %q", branches)
		}
		return nil
	})

	sc.Step(`^there are (?:now|still) no observed branches$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branches := devRepo.Config.NormalConfig.LocalGitConfig.ObservedBranches
		if len(branches) > 0 {
			return fmt.Errorf("expected no observed branches, got %q", branches)
		}
		return nil
	})

	sc.Step(`^there are (?:now|still) no parked branches$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branches := devRepo.Config.NormalConfig.LocalGitConfig.ParkedBranches
		if len(branches) > 0 {
			return fmt.Errorf("expected no parked branches, got %q", branches)
		}
		return nil
	})

	sc.Step(`^there are (?:now|still) no perennial branches$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branches := devRepo.Config.NormalConfig.LocalGitConfig.PerennialBranches
		if len(branches) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", branches)
		}
		return nil
	})

	sc.Step(`^there are (?:now|still) no prototype branches$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branches := devRepo.Config.NormalConfig.LocalGitConfig.PrototypeBranches
		if len(branches) > 0 {
			return fmt.Errorf("expected no prototype branches, got %q", branches)
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
		fileTable := devRepo.FilesInBranches(gitdomain.NewLocalBranchName("main"))
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

	sc.Step(`^the uncommitted file is stashed$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		uncommittedFiles := devRepo.UncommittedFiles()
		for _, ucf := range uncommittedFiles {
			if ucf == state.uncommittedFileName.GetOrPanic() {
				return fmt.Errorf("expected file %q to be stashed but it is still uncommitted", state.uncommittedFileName)
			}
		}
		stashSize, err := devRepo.StashSize(devRepo.TestRunner)
		asserts.NoError(err)
		if stashSize != 1 {
			return fmt.Errorf("expected 1 stash but found %d", stashSize)
		}
		return nil
	})

	sc.Step(`^the uncommitted file still exists$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		hasFile := devRepo.HasFile(
			state.uncommittedFileName.GetOrPanic(),
			state.uncommittedContent.GetOrPanic(),
		)
		if len(hasFile) > 0 {
			panic(hasFile)
		}
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
	})

	sc.Step(`^tool "([^"]*)" is installed$`, func(ctx context.Context, tool string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.MockCommand(tool)
	})

	// This step exists to avoid re-creating commits with the same SHA as existing commits
	// because that can cause flaky tests.
	sc.Step(`wait 1 second to ensure new Git timestamps`, func() {
		time.Sleep(1 * time.Second)
	})
}

func updateInitialSHAs(state *ScenarioState) {
	devRepo := state.fixture.DevRepo.GetOrPanic()
	if state.initialDevSHAs.IsNone() && state.insideGitRepo {
		state.initialDevSHAs = Some(devRepo.CommitSHAs())
	}
	if originRepo, hasOriginrepo := state.fixture.OriginRepo.Get(); state.initialOriginSHAs.IsNone() && state.insideGitRepo && hasOriginrepo {
		state.initialOriginSHAs = Some(originRepo.CommitSHAs())
	}
	if secondWorkTree, hasSecondWorkTree := state.fixture.SecondWorktree.Get(); state.initialWorktreeSHAs.IsNone() && state.insideGitRepo && hasSecondWorkTree {
		state.initialWorktreeSHAs = Some(secondWorkTree.CommitSHAs())
	}
}
