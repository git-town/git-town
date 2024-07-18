package cucumber

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/acarl005/stripansi"
	"github.com/cucumber/godog"
	cukemessages "github.com/cucumber/messages/go/v21"
	"github.com/git-town/git-town/v14/src/cli/dialog/components"
	"github.com/git-town/git-town/v14/src/cli/print"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/config/configfile"
	"github.com/git-town/git-town/v14/src/config/gitconfig"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/gohacks"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/git-town/git-town/v14/test/asserts"
	"github.com/git-town/git-town/v14/test/commands"
	"github.com/git-town/git-town/v14/test/datatable"
	"github.com/git-town/git-town/v14/test/filesystem"
	"github.com/git-town/git-town/v14/test/fixture"
	"github.com/git-town/git-town/v14/test/git"
	"github.com/git-town/git-town/v14/test/helpers"
	"github.com/git-town/git-town/v14/test/output"
	"github.com/git-town/git-town/v14/test/subshell"
	"github.com/git-town/git-town/v14/test/testruntime"
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
			panic("after-scenario hook has found no scenario state found to clean up")
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

	sc.Step("a Git repo clone", func(ctx context.Context) (context.Context, error) {
		scenarioName := ctx.Value(keyScenarioName).(string)
		scenarioTags := ctx.Value(keyScenarioTags).([]*cukemessages.PickleTag)
		fixture := fixtureFactory.CreateFixture(scenarioName)
		if helpers.HasTag(scenarioTags, "@debug") {
			fixture.DevRepo.GetOrPanic().Verbose = true
		}
		state := ScenarioState{
			fixture:              fixture,
			initialBranches:      None[datatable.DataTable](),
			initialCommits:       None[datatable.DataTable](),
			initialCurrentBranch: None[gitdomain.LocalBranchName](),
			initialDevSHAs:       None[map[string]gitdomain.SHA](),
			initialLineage:       None[datatable.DataTable](),
			initialOriginSHAs:    None[map[string]gitdomain.SHA](),
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

	sc.Step("a local Git repo", func(ctx context.Context) (context.Context, error) {
		scenarioName := ctx.Value(keyScenarioName).(string)
		scenarioTags := ctx.Value(keyScenarioTags).([]*cukemessages.PickleTag)
		fixture := fixtureFactory.CreateFixture(scenarioName)
		devRepo := fixture.DevRepo.GetOrPanic()
		if helpers.HasTag(scenarioTags, "@debug") {
			devRepo.Verbose = true
		}
		devRepo.RemoveRemote(gitdomain.RemoteOrigin)
		fixture.OriginRepo = NoneP[testruntime.TestRuntime]()
		state := ScenarioState{
			fixture:              fixture,
			initialBranches:      None[datatable.DataTable](),
			initialCommits:       None[datatable.DataTable](),
			initialCurrentBranch: None[gitdomain.LocalBranchName](),
			initialDevSHAs:       None[map[string]gitdomain.SHA](),
			initialLineage:       None[datatable.DataTable](),
			initialOriginSHAs:    None[map[string]gitdomain.SHA](),
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

	sc.Step(`^a merge is now in progress$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		if !devRepo.HasMergeInProgress(devRepo.TestRunner) {
			panic("expected merge in progress")
		}
	})

	sc.Step(`^a rebase is now in progress$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		repoStatus, err := devRepo.RepoStatus(devRepo.TestRunner)
		asserts.NoError(err)
		if !repoStatus.RebaseInProgress {
			panic("expected rebase in progress")
		}
	})

	sc.Step(`^a remote tag "([^"]+)" not on a branch$`, func(ctx context.Context, name string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.OriginRepo.GetOrPanic().CreateStandaloneTag(name)
	})

	sc.Step(`^all branches are now synchronized$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branchesOutOfSync, output := devRepo.HasBranchesOutOfSync()
		if branchesOutOfSync {
			panic("unexpected out of sync:\n" + output)
		}
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

	sc.Step(`^an uncommitted file in folder "([^"]*)"$`, func(ctx context.Context, folder string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		fileName := folder + "/uncommitted file"
		state.uncommittedFileName = Some(fileName)
		content := "uncommitted content"
		state.uncommittedContent = Some(content)
		devRepo.CreateFile(fileName, content)
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

	sc.Step(`^a remote "([^"]+)" pointing to "([^"]+)"`, func(ctx context.Context, name, url string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.AddRemote(gitdomain.Remote(name), url)
	})

	sc.Step(`^branch "([^"]+)" is active in another worktree`, func(ctx context.Context, branch string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.AddSecondWorktree(gitdomain.NewLocalBranchName(branch))
	})

	sc.Step(`^branch "([^"]+)" is (?:now|still) a contribution branch`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		if !devRepo.Config.Config.IsContributionBranch(branch) {
			return fmt.Errorf(
				"branch %q isn't contribution as expected.\nContribution branches: %s",
				branch,
				strings.Join(devRepo.Config.Config.ContributionBranches.Strings(), ", "),
			)
		}
		return nil
	})

	sc.Step(`^branch "([^"]+)" is (?:now|still) a feature branch`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		if devRepo.Config.Config.BranchType(branch) != configdomain.BranchTypeFeatureBranch {
			return fmt.Errorf("branch %q isn't a feature branch as expected", branch)
		}
		return nil
	})

	sc.Step(`^branch "([^"]+)" is (?:now|still) observed`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		if !devRepo.Config.Config.IsObservedBranch(branch) {
			return fmt.Errorf(
				"branch %q isn't observed as expected.\nObserved branches: %s",
				branch,
				strings.Join(devRepo.Config.Config.ObservedBranches.Strings(), ", "),
			)
		}
		return nil
	})

	sc.Step(`^branch "([^"]+)" is now parked`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		if !devRepo.Config.Config.IsParkedBranch(branch) {
			return fmt.Errorf(
				"branch %q isn't parked as expected.\nParked branches: %s",
				branch,
				strings.Join(devRepo.Config.Config.ParkedBranches.Strings(), ", "),
			)
		}
		return nil
	})

	sc.Step(`^branch "([^"]+)" is (?:now|still) perennial`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		if !devRepo.Config.Config.IsPerennialBranch(branch) {
			return fmt.Errorf(
				"branch %q isn't perennial as expected.\nPerennial branches: %s",
				branch,
				strings.Join(devRepo.Config.Config.PerennialBranches.Strings(), ", "),
			)
		}
		return nil
	})

	sc.Step(`^branch "([^"]+)" is now a feature branch`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		if devRepo.Config.Config.IsParkedBranch(branch) {
			return fmt.Errorf("branch %q is parked", branch)
		}
		if devRepo.Config.Config.IsObservedBranch(branch) {
			return fmt.Errorf("branch %q is observed", branch)
		}
		if devRepo.Config.Config.IsContributionBranch(branch) {
			return fmt.Errorf("branch %q is contribution", branch)
		}
		if devRepo.Config.Config.IsPerennialBranch(branch) {
			return fmt.Errorf("branch %q is perennial", branch)
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

	sc.Step(`^file "([^"]+)" with content$`, func(ctx context.Context, name string, content *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		filePath := filepath.Join(devRepo.WorkingDir, name)
		//nolint:gosec // need permission 700 here in order for tests to work
		return os.WriteFile(filePath, []byte(content.Content), 0o700)
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

	sc.Step(`^file "([^"]*)" (?:now|still) has content "([^"]*)"$`, func(ctx context.Context, file, expectedContent string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actualContent := devRepo.FileContent(file)
		if expectedContent != actualContent {
			return fmt.Errorf("file content does not match\n\nEXPECTED: %q\n\nACTUAL:\n\n%q\n----------------------------", expectedContent, actualContent)
		}
		return nil
	})

	sc.Step(`^Git has version "([^"]*)"$`, func(ctx context.Context, version string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.MockGit(version)
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

	sc.Step(`^Git Town setting "color.ui" is "([^"]*)"$`, func(ctx context.Context, value string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		return devRepo.SetColorUI(value)
	})

	// TODO: remove?
	sc.Step(`^Git Town parent setting for branch "([^"]*)" is "([^"]*)"$`, func(ctx context.Context, branch, value string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branchName := gitdomain.NewLocalBranchName(branch)
		configKey := gitconfig.NewParentKey(branchName)
		return devRepo.Config.GitConfig.SetLocalConfigValue(configKey, value)
	})

	sc.Step(`^local Git setting "init.defaultbranch" is "([^"]*)"$`, func(ctx context.Context, value string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.SetDefaultGitBranch(gitdomain.NewLocalBranchName(value))
	})

	sc.Step(`^global Git setting "alias\.(.*?)" is "([^"]*)"$`, func(ctx context.Context, name, value string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		key, hasKey := gitconfig.ParseKey("alias." + name).Get()
		if !hasKey {
			return fmt.Errorf("no key found for %q", name)
		}
		aliasableCommand, hasAliasableCommand := gitconfig.AliasableCommandForKey(key).Get()
		if !hasAliasableCommand {
			return fmt.Errorf("no aliasableCommand found for key %q", key)
		}
		return devRepo.SetGitAlias(aliasableCommand, value)
	})

	sc.Step(`^global Git setting "alias\.(.*?)" (?:now|still) doesn't exist$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		key, hasKey := gitconfig.ParseKey("alias." + name).Get()
		if !hasKey {
			return errors.New("key not found")
		}
		aliasableCommand, hasAliasableCommand := gitconfig.AliasableCommandForKey(key).Get()
		if !hasAliasableCommand {
			return errors.New("unknown alias: " + key.String())
		}
		command, has := devRepo.Config.Config.Aliases[aliasableCommand]
		if has {
			return fmt.Errorf("unexpected aliasableCommand %q: %q", key, command)
		}
		return nil
	})

	sc.Step(`^global Git setting "alias\.(.*?)" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, name, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		key, hasKey := gitconfig.ParseKey("alias." + name).Get()
		if !hasKey {
			return errors.New("key not found")
		}
		aliasableCommand, hasAliasableCommand := gitconfig.AliasableCommandForKey(key).Get()
		if !hasAliasableCommand {
			return fmt.Errorf("aliasableCommand not found for key %q", key)
		}
		have := devRepo.Config.Config.Aliases[aliasableCommand]
		if have != want {
			return fmt.Errorf("unexpected value for key %q: want %q have %q", name, want, have)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "([^"]*)" is "([^"]*)"$`, func(ctx context.Context, name, value string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		configKey, hasConfigKey := gitconfig.ParseKey("git-town." + name).Get()
		if !hasConfigKey {
			return fmt.Errorf("unknown configuration key: %q", name)
		}
		return devRepo.Config.GitConfig.SetGlobalConfigValue(configKey, value)
	})

	sc.Step(`^global Git Town setting "([^"]*)" (?:now|still) doesn't exist$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		configKey, hasConfigKey := gitconfig.ParseKey("git-town." + name).Get()
		if !hasConfigKey {
			return errors.New("unknown config key: " + name)
		}
		newValue, hasNewValue := devRepo.TestCommands.GlobalGitConfig(configKey).Get()
		if hasNewValue {
			return fmt.Errorf("should not have global %q anymore but has value %q", name, newValue)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "hosting-origin-hostname" is now "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.GlobalGitConfig.HostingOriginHostname.String()
		if have != want {
			return fmt.Errorf(`expected global setting "hosting-origin-hostname" to be %q, but was %q`, want, have)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "hosting-platform" is now "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.GlobalGitConfig.HostingPlatform
		if have.String() != want {
			return fmt.Errorf(`expected global setting "hosting-platform" to be %q, but was %q`, want, have)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "main-branch" is now "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.GlobalGitConfig.MainBranch.String()
		if have != want {
			return fmt.Errorf(`expected global setting "main-branch" to be %q, but was %q`, want, have)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "offline" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		wantBool, err := gohacks.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.Offline(wantBool)
		have, exists := devRepo.Config.GlobalGitConfig.Offline.Get()
		if !exists {
			return fmt.Errorf(`expected global setting "offline" to be %t, but doesn't exist`, want)
		}
		if have != want {
			return fmt.Errorf(`expected global setting "offline" to be %t, but was %t`, want, have)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "perennial-branches" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.GlobalGitConfig.PerennialBranches
		want := gitdomain.NewLocalBranchNames(strings.Split(wantStr, " ")...)
		if !cmp.Equal(have, want) {
			return fmt.Errorf(`expected global setting "perennial-branches" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "push-hook" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.GlobalGitConfig.PushHook.String()
		if !cmp.Equal(have, want) {
			return fmt.Errorf(`expected global setting "push-hook" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "push-new-branches" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have, has := devRepo.Config.GlobalGitConfig.PushNewBranches.Get()
		if !has {
			return errors.New(`expected global setting "push-new-branches" to exist but it doesn't`)
		}
		want, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		if have.Bool() != want {
			return fmt.Errorf(`expected global setting "push-new-branches" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "ship-delete-tracking-branch" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		want, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		have, has := devRepo.Config.GlobalGitConfig.ShipDeleteTrackingBranch.Get()
		if !has {
			return fmt.Errorf(`expected global setting "ship-delete-tracking-branch" to be %v, but doesn't exist`, want)
		}
		if have.Bool() != want {
			return fmt.Errorf(`expected global setting "ship-delete-tracking-branch" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "sync-before-ship" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		want, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		have, has := devRepo.Config.GlobalGitConfig.SyncBeforeShip.Get()
		if !has {
			return fmt.Errorf(`expected global setting "sync-before-ship" to be %v, but doesn't exist`, want)
		}
		if have.Bool() != want {
			return fmt.Errorf(`expected global setting "sync-before-ship" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "sync-feature-strategy" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		want, err := configdomain.NewSyncFeatureStrategy(wantStr)
		asserts.NoError(err)
		have, has := devRepo.Config.GlobalGitConfig.SyncFeatureStrategy.Get()
		if !has {
			return fmt.Errorf(`expected global setting "sync-feature-strategy" to be %v, but doesn't exist`, want)
		}
		if have != want {
			return fmt.Errorf(`expected global setting "sync-feature-strategy" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "sync-perennial-strategy" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		want, err := configdomain.NewSyncPerennialStrategy(wantStr)
		asserts.NoError(err)
		have, has := devRepo.Config.GlobalGitConfig.SyncPerennialStrategy.Get()
		if !has {
			return fmt.Errorf(`expected global setting "sync-perennial-strategy" to be %v, but it doesn't exist`, want)
		}
		if have != want {
			return fmt.Errorf(`expected global setting "sync-perennial-strategy" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^global Git Town setting "sync-upstream" is (?:now|still) "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		wantBool, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.SyncUpstream(wantBool)
		have, has := devRepo.Config.GlobalGitConfig.SyncUpstream.Get()
		if !has {
			return fmt.Errorf(`expected global setting "sync-upstream" to be %v, but doesn't exist`, want)
		}
		if have != want {
			return fmt.Errorf(`expected global setting "sync-upstream" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^I add commit "([^"]*)" to the "([^"]*)" branch`, func(ctx context.Context, message, branch string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.CreateCommit(git.Commit{
			Branch:   gitdomain.NewLocalBranchName(branch),
			FileName: "new_file",
			Message:  message,
		})
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
		envDirName := filesystem.FolderName(scenarioName) + "_" + fixtureFactory.Counter.ToString()
		envPath := filepath.Join(fixtureFactory.Dir, envDirName)
		asserts.NoError(os.Mkdir(envPath, 0o777))
		fixture := fixture.Fixture{
			CoworkerRepo:   NoneP[testruntime.TestRuntime](),
			DevRepo:        NoneP[testruntime.TestRuntime](),
			Dir:            envPath,
			OriginRepo:     NoneP[testruntime.TestRuntime](),
			SecondWorktree: NoneP[testruntime.TestRuntime](),
			SubmoduleRepo:  NoneP[testruntime.TestRuntime](),
			UpstreamRepo:   NoneP[testruntime.TestRuntime](),
		}
		state := ScenarioState{
			fixture:              fixture,
			initialBranches:      None[datatable.DataTable](),
			initialCommits:       None[datatable.DataTable](),
			initialCurrentBranch: None[gitdomain.LocalBranchName](),
			initialDevSHAs:       None[map[string]gitdomain.SHA](),
			initialLineage:       None[datatable.DataTable](),
			initialOriginSHAs:    None[map[string]gitdomain.SHA](),
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

	sc.Step(`^I resolve the conflict in "([^"]*)"(?: with "([^"]*)")?$`, func(ctx context.Context, filename, content string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		if content == "" {
			content = "resolved content"
		}
		devRepo.CreateFile(filename, content)
		devRepo.StageFiles(filename)
	})

	sc.Step(`^I resolve the conflict in "([^"]*)" in the other worktree$`, func(ctx context.Context, filename string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		content := "resolved content"
		secondWorkTree := state.fixture.SecondWorktree.GetOrPanic()
		secondWorkTree.CreateFile(filename, content)
		secondWorkTree.StageFiles(filename)
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

	sc.Step(`^I run "([^"]*)", enter into the dialog, and close the next editor:$`, func(ctx context.Context, cmd string, input *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.CaptureState()
		updateInitialSHAs(state)
		env := append(os.Environ(), "GIT_EDITOR=true")
		answers, err := helpers.TableToInputEnv(input)
		asserts.NoError(err)
		for dialogNumber, answer := range answers {
			env = append(env, fmt.Sprintf("%s%d=%s", components.TestInputKey, dialogNumber, answer))
		}
		var exitCode int
		var output string
		output, exitCode = devRepo.MustQueryStringCodeWith(cmd, &subshell.Options{Env: env})
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

	sc.Step(`^it does not print "(.+)"$`, func(ctx context.Context, text string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		if strings.Contains(stripansi.Strip(state.runOutput.GetOrPanic()), text) {
			return fmt.Errorf("text found: %q", text)
		}
		return nil
	})

	sc.Step(`^it prints:$`, func(ctx context.Context, expected *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		if exitCode := state.runExitCode.GetOrPanic(); exitCode != 0 {
			return fmt.Errorf("unexpected exit code %d", state.runExitCode)
		}
		if !strings.Contains(stripansi.Strip(state.runOutput.GetOrPanic()), expected.Content) {
			fmt.Println("ERROR: text not found:")
			fmt.Println("\nEXPECTED:", expected.Content)
			fmt.Println()
			fmt.Println("==================================================================")
			fmt.Println("ACTUAL OUTPUT START ==============================================")
			fmt.Println("==================================================================")
			fmt.Println()
			fmt.Println(state.runOutput.GetOrPanic())
			fmt.Println()
			fmt.Println("==================================================================")
			fmt.Println("ACTUAL OUTPUT END ================================================")
			fmt.Println("==================================================================")
			fmt.Println()
			return errors.New("expected text not found")
		}
		return nil
	})

	sc.Step(`^it prints no output$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		output := state.runOutput.GetOrPanic()
		if output != "" {
			return fmt.Errorf("expected no output but found %q", output)
		}
		return nil
	})

	sc.Step(`^it prints something like:$`, func(ctx context.Context, expected *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		regex := regexp.MustCompile(expected.Content)
		have := stripansi.Strip(state.runOutput.GetOrPanic())
		if !regex.MatchString(have) {
			return fmt.Errorf("EXPECTED: content matching %q\nGOT: %q", expected.Content, have)
		}
		return nil
	})

	sc.Step(`^it prints the error:$`, func(ctx context.Context, expected *godog.DocString) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.runExitCodeChecked = true
		if !strings.Contains(stripansi.Strip(state.runOutput.GetOrPanic()), expected.Content) {
			return fmt.Errorf("text not found:\n%s\n\nactual text:\n%s", expected.Content, state.runOutput)
		}
		if exitCode := state.runExitCode.GetOrPanic(); exitCode == 0 {
			return fmt.Errorf("unexpected exit code %d", state.runExitCode)
		}
		return nil
	})

	sc.Step(`^it runs no commands$`, func(ctx context.Context) error {
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

	sc.Step(`^it runs the commands$`, func(ctx context.Context, input *godog.Table) {
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
			panic("mismatching commands run, see diff above")
		}
	})

	sc.Step(`^it runs without error$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		if state.runExitCode.GetOrPanic() != 0 {
			return fmt.Errorf("did not expect the Git Town command to produce an exit code: %d", state.runExitCode)
		}
		return nil
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

	sc.Step(`^local Git Town setting "([^"]*)" (:?now|still) doesn't exist$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		configKey, hasConfigKey := gitconfig.ParseKey("git-town." + name).Get()
		if !hasConfigKey {
			return errors.New("unknown config key: " + name)
		}
		newValue, hasNewValue := devRepo.TestCommands.LocalGitConfig(configKey).Get()
		if hasNewValue {
			return fmt.Errorf("should not have local %q anymore but has value %q", name, newValue)
		}
		return nil
	})

	sc.Step(`^(?:local )?Git Town setting "([^"]*)" doesn't exist$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		configKey, hasConfigKey := gitconfig.ParseKey("git-town." + name).Get()
		if !hasConfigKey {
			return errors.New("unknown config key: " + name)
		}
		return devRepo.Config.GitConfig.RemoveLocalConfigValue(configKey)
	})

	sc.Step(`^(?:local )?Git Town setting "([^"]*)" is "([^"]*)"$`, func(ctx context.Context, name, value string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		configKey, hasConfigKey := gitconfig.ParseKey("git-town." + name).Get()
		if !hasConfigKey {
			return fmt.Errorf("unknown config key: %q", name)
		}
		return devRepo.Config.GitConfig.SetLocalConfigValue(configKey, value)
	})

	sc.Step(`^local Git Town setting "code-hosting-origin-hostname" now doesn't exist$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.HostingOriginHostname
		if have.IsSome() {
			return fmt.Errorf(`unexpected local setting "code-hosting-origin-hostname" with value %q`, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "hosting-platform" is now "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.HostingPlatform
		if have.String() != want {
			return fmt.Errorf(`expected local setting "hosting-platform" to be %q, but was %q`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "hosting-platform" (:?now|still) doesn't exist$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.HostingPlatform
		if value, has := have.Get(); has {
			return fmt.Errorf(`expected local setting "hosting-platform" to not exist but was %q`, value)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "gitea-token" is now "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.GiteaToken.String()
		if have != want {
			return fmt.Errorf(`expected local setting "gitea-token" to be %q, but was %q`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "github-token" is now "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.GitHubToken.String()
		if have != want {
			return fmt.Errorf(`expected local setting "github-token" to be %q, but was %q`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "github-token" now doesn't exist$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.GitHubToken
		if have.IsSome() {
			return fmt.Errorf(`unexpected local setting "github-token" with value %q`, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "gitlab-token" is now "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.GitLabToken.String()
		if have != want {
			return fmt.Errorf(`expected local setting "gitlab-token" to be %q, but was %q`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "hosting-origin-hostname" is now "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.HostingOriginHostname
		if have.String() != want {
			return fmt.Errorf(`expected local setting "hosting-origin-hostname" to be %q, but was %q`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "hosting-origin-hostname" now doesn't exist$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.HostingOriginHostname
		if have.IsSome() {
			return fmt.Errorf(`unexpected local setting "hosting-origin-hostname" with value %q`, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "main-branch" is now "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.MainBranch.String()
		if have != want {
			return fmt.Errorf(`expected local setting "main-branch" to be %q, but was %q`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "perennial-branches" is now "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.PerennialBranches
		want := gitdomain.NewLocalBranchNames(strings.Split(wantStr, " ")...)
		if !cmp.Equal(have, want) {
			return fmt.Errorf(`expected local setting "perennial-branches" to be %q, but was %q`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "perennial-regex" is now "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.PerennialRegex.String()
		if have != want {
			return fmt.Errorf(`expected local setting "perennial-regex" to be %q, but was %q`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "push-hook" is (:?now|still) not set$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.PushHook
		if have.IsSome() {
			return fmt.Errorf(`unexpected local setting "push-hook" %v`, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "push-hook" is now "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.PushHook.String()
		if !cmp.Equal(have, want) {
			return fmt.Errorf(`expected local setting "push-hook" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "push-new-branches" is (:?now|still) not set$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.PushNewBranches
		if value, has := have.Get(); has {
			return fmt.Errorf(`unexpected local setting "push-new-branches" %v`, value)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "push-new-branches" is now "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		want, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		pushNewBranches, has := devRepo.Config.LocalGitConfig.PushNewBranches.Get()
		if !has {
			return fmt.Errorf(`expected local setting "push-new-branches" to be %v, but it doesn't exist`, want)
		}
		have := pushNewBranches.Bool()
		if have != want {
			return fmt.Errorf(`expected local setting "push-new-branches" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "ship-delete-tracking-branch" is still not set$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have, has := devRepo.Config.LocalGitConfig.ShipDeleteTrackingBranch.Get()
		if has {
			return fmt.Errorf(`unexpected local setting "ship-delete-tracking-branch" %v`, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "ship-delete-tracking-branch" is now "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		want, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		have, has := devRepo.Config.LocalGitConfig.ShipDeleteTrackingBranch.Get()
		if !has {
			return fmt.Errorf(`expected local setting "ship-delete-tracking-branch" to be %v, but doesn't exist`, want)
		}
		if have.Bool() != want {
			return fmt.Errorf(`expected local setting "ship-delete-tracking-branch" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "sync-before-ship" is still not set$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have, has := devRepo.Config.LocalGitConfig.SyncBeforeShip.Get()
		if has {
			return fmt.Errorf(`unexpected local setting "sync-before-ship" %v`, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "sync-before-ship" is now "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		want, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		have, has := devRepo.Config.LocalGitConfig.SyncBeforeShip.Get()
		if !has {
			return fmt.Errorf(`expected local setting "sync-before-ship" to be %v, but doesn't exist`, want)
		}
		if have.Bool() != want {
			return fmt.Errorf(`expected local setting "sync-before-ship" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "sync-feature-strategy" is still not set$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have, has := devRepo.Config.LocalGitConfig.SyncFeatureStrategy.Get()
		if has {
			return fmt.Errorf(`unexpected local setting "sync-feature-strategy" %v`, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "sync-feature-strategy" is now "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		want, err := configdomain.NewSyncFeatureStrategy(wantStr)
		asserts.NoError(err)
		have, has := devRepo.Config.LocalGitConfig.SyncFeatureStrategy.Get()
		if !has {
			return fmt.Errorf(`expected local setting "sync-feature-strategy" to be %v, but doesn't exist`, want)
		}
		if have != want {
			return fmt.Errorf(`expected local setting "sync-feature-strategy" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "sync-perennial-strategy" is still not set$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have, has := devRepo.Config.LocalGitConfig.SyncPerennialStrategy.Get()
		if has {
			return fmt.Errorf(`unexpected local setting "sync-perennial-strategy" %v`, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "sync-perennial-strategy" is now "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		want, err := configdomain.NewSyncPerennialStrategy(wantStr)
		asserts.NoError(err)
		have, has := devRepo.Config.LocalGitConfig.SyncPerennialStrategy.Get()
		if !has {
			return fmt.Errorf(`expected local setting "sync-perennial-strategy" to be %v, but doesn't exist`, want)
		}
		if have != want {
			return fmt.Errorf(`expected local setting "sync-perennial-strategy" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "sync-upstream" is still not set$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have, has := devRepo.Config.LocalGitConfig.SyncUpstream.Get()
		if has {
			return fmt.Errorf(`unexpected local setting "sync-upstream" %v`, have)
		}
		return nil
	})

	sc.Step(`^local Git Town setting "sync-upstream" is now "([^"]*)"$`, func(ctx context.Context, wantStr string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		wantBool, err := strconv.ParseBool(wantStr)
		asserts.NoError(err)
		want := configdomain.SyncUpstream(wantBool)
		have, has := devRepo.Config.LocalGitConfig.SyncUpstream.Get()
		if !has {
			return fmt.Errorf(`expected local setting "sync-upstream" to be %v, but doesn't exist`, want)
		}
		if have != want {
			return fmt.Errorf(`expected local setting "sync-upstream" to be %v, but was %v`, want, have)
		}
		return nil
	})

	sc.Step(`^my repo does not have an origin$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		devRepo.RemoveRemote(gitdomain.RemoteOrigin)
		state.fixture.OriginRepo = NoneP[testruntime.TestRuntime]()
	})

	sc.Step(`^my repo has a Git submodule$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		state.fixture.AddSubmoduleRepo()
		devRepo.AddSubmodule(state.fixture.SubmoduleRepo.GetOrPanic().WorkingDir)
	})

	sc.Step(`^my repo's "([^"]*)" remote is "([^"]*)"$`, func(ctx context.Context, remoteName, remoteURL string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		remote := gitdomain.Remote(remoteName)
		devRepo.RemoveRemote(remote)
		devRepo.AddRemote(remote, remoteURL)
	})

	sc.Step(`^still no configuration file exists$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		_, err := devRepo.FileContentErr(configfile.FileName)
		if err == nil {
			return errors.New("expected no configuration file but found one")
		}
		return nil
	})

	sc.Step(`^no commits exist now$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		currentCommits := state.fixture.CommitTable(state.initialCommits.GetOrPanic().Cells[0])
		noCommits := datatable.DataTable{}
		noCommits.AddRow(state.initialCommits.GetOrPanic().Cells[0]...)
		errDiff, errCount := currentCommits.EqualDataTable(noCommits)
		if errCount == 0 {
			return
		}
		fmt.Println(errDiff)
		panic("found unexpected commits")
	})

	sc.Step(`^no lineage exists now$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		if devRepo.Config.Config.ContainsLineage() {
			lineage := devRepo.Config.Config.Lineage
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

	sc.Step(`^no uncommitted files exist$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		files := devRepo.UncommittedFiles()
		if len(files) > 0 {
			return fmt.Errorf("unexpected uncommitted files: %s", files)
		}
		return nil
	})

	sc.Step(`^offline mode is disabled$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		isOffline := devRepo.Config.Config.Offline
		if isOffline {
			return errors.New("expected to not be offline but am")
		}
		return nil
	})

	sc.Step(`^offline mode is enabled$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		return devRepo.Config.SetOffline(true)
	})

	sc.Step(`^origin deletes the "([^"]*)" branch$`, func(ctx context.Context, branch string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		state.fixture.OriginRepo.GetOrPanic().RemoveBranch(gitdomain.NewLocalBranchName(branch))
	})

	sc.Step(`^origin ships the "([^"]*)" branch$`, func(ctx context.Context, branch string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		originRepo := state.fixture.OriginRepo.GetOrPanic()
		originRepo.CheckoutBranch(gitdomain.NewLocalBranchName("main"))
		err := originRepo.MergeBranch(gitdomain.NewLocalBranchName(branch))
		asserts.NoError(err)
		originRepo.RemoveBranch(gitdomain.NewLocalBranchName(branch))
	})

	sc.Step("^the branches$", func(ctx context.Context, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		// TODO: uncomment this and make it work
		// if state.initialBranches.IsNone() {
		// 	initialTable := datatable.FromGherkin(table)
		// 	state.initialBranches = Some(initialTable)
		// } else {
		// 	state.initialBranches = None[datatable.DataTable]()
		// }
		for _, branchSetup := range datatable.ParseBranchSetupTable(table) {
			var repoToCreateBranchIn *testruntime.TestRuntime
			switch {
			case branchSetup.Locations.Is(git.LocationLocal), branchSetup.Locations.Is(git.LocationLocal, git.LocationOrigin):
				repoToCreateBranchIn = state.fixture.DevRepo.GetOrPanic()
			case branchSetup.Locations.Is(git.LocationOrigin):
				repoToCreateBranchIn = state.fixture.OriginRepo.GetOrPanic()
			case branchSetup.Locations.Is(git.LocationUpstream):
				repoToCreateBranchIn = state.fixture.UpstreamRepo.GetOrPanic()
			default:
				panic("unhandled location to create the new branch: " + branchSetup.Locations.String())
			}
			branchType, hasBranchType := branchSetup.BranchType.Get()
			if hasBranchType {
				switch branchType {
				case configdomain.BranchTypeMainBranch:
					panic("main branch exists already")
				case configdomain.BranchTypeFeatureBranch:
					repoToCreateBranchIn.CreateChildFeatureBranch(branchSetup.Name, branchSetup.Parent.GetOrElse("main"))
				case configdomain.BranchTypePerennialBranch:
					repoToCreateBranchIn.CreatePerennialBranches(branchSetup.Name)
				case configdomain.BranchTypeContributionBranch:
					repoToCreateBranchIn.CreateContributionBranches(branchSetup.Name)
				case configdomain.BranchTypeObservedBranch:
					repoToCreateBranchIn.CreateObservedBranches(branchSetup.Name)
				case configdomain.BranchTypeParkedBranch:
					repoToCreateBranchIn.CreateParkedBranches(branchSetup.Name)
				}
			} else {
				repoToCreateBranchIn.CreateBranch(branchSetup.Name, "main")
			}
			if len(branchSetup.Locations) > 1 {
				switch {
				case branchSetup.Locations.Is(git.LocationLocal, git.LocationOrigin):
					state.fixture.DevRepo.GetOrPanic().PushBranchToRemote(branchSetup.Name, gitdomain.RemoteOrigin)
				default:
					panic("unhandled location to push the new branch to: " + branchSetup.Locations.String())
				}
			}
		}
	})

	sc.Step(`^the branches are now$`, func(ctx context.Context, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		existing := state.fixture.Branches()
		diff, errCount := existing.EqualGherkin(table)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the branches\n\n", errCount)
			fmt.Println(diff)
			panic("mismatching branches found, see the diff above")
		}
	})

	sc.Step(`^the commits$`, func(ctx context.Context, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		// TODO: uncomment this and make it work
		// if state.initialCommits.IsNone() {
		// 	initialTable := datatable.FromGherkin(table)
		// 	state.initialCommits = Some(initialTable)
		// } else {
		// 	state.initialCommits = None[datatable.DataTable]()
		// }
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

	sc.Step(`^the configuration file is (?:now|still):$`, func(ctx context.Context, content *godog.DocString) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have, err := devRepo.FileContentErr(configfile.FileName)
		if err != nil {
			panic("no configuration file found")
		}
		have = strings.TrimSpace(have)
		want := strings.TrimSpace(content.Content)
		if have != want {
			fmt.Println(cmp.Diff(want, have))
			panic("mismatching config file content")
		}
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
		_ = state.fixture.CoworkerRepo.GetOrPanic().Config.SetParent(gitdomain.NewLocalBranchName(childBranch), gitdomain.NewLocalBranchName(parentBranch))
	})

	sc.Step(`^the coworker sets the "sync-feature-strategy" to "(merge|rebase)"$`, func(ctx context.Context, value string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		syncFeatureStrategy, err := configdomain.NewSyncFeatureStrategy(value)
		asserts.NoError(err)
		_ = state.fixture.CoworkerRepo.GetOrPanic().Config.SetSyncFeatureStrategy(syncFeatureStrategy)
	})

	sc.Step(`^the coworkers workspace now contains file "([^"]*)" with content "([^"]*)"$`, func(ctx context.Context, file, expectedContent string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		actualContent := state.fixture.CoworkerRepo.GetOrPanic().FileContent(file)
		if expectedContent != actualContent {
			return fmt.Errorf("file content does not match\n\nEXPECTED: %q\n\nACTUAL:\n\n%q\n----------------------------", expectedContent, actualContent)
		}
		return nil
	})

	sc.Step(`^the current branch is "([^"]*)"$`, func(ctx context.Context, name string) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branch := gitdomain.NewLocalBranchName(name)
		state.initialCurrentBranch = Some(branch)
		if !devRepo.BranchExists(devRepo.TestRunner, branch) {
			devRepo.CreateBranch(branch, gitdomain.NewLocalBranchName("main"))
		}
		devRepo.CheckoutBranch(branch)
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

	sc.Step(`^the current branch in the other worktree is (?:now|still) "([^"]*)"$`, func(ctx context.Context, expected string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		secondWorkTree := state.fixture.SecondWorktree.GetOrPanic()
		secondWorkTree.CurrentBranchCache.Invalidate()
		actual, err := secondWorkTree.CurrentBranch(secondWorkTree.TestCommands)
		if err != nil {
			return fmt.Errorf("cannot determine current branch of second worktree: %w", err)
		}
		if actual.String() != expected {
			return fmt.Errorf("expected active branch %q but is %q", expected, actual)
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

	sc.Step(`^the initial lineage exists$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.LineageTable()
		diff, errCnt := have.EqualDataTable(state.initialLineage.GetOrPanic())
		if errCnt > 0 {
			fmt.Printf("\nERROR! Found %d differences in the lineage\n\n", errCnt)
			fmt.Printf("INITIAL LINEAGE:\n%s\n", state.initialLineage.String())
			fmt.Printf("CURRENT LINEAGE:\n%s\n", have.String())
			fmt.Println(diff)
			panic("mismatching branches found, see the diff above")
		}
	})

	sc.Step(`^the initial branches and lineage exist$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		// verify initial branches
		currentBranches := state.fixture.Branches()
		// fmt.Printf("\nINITIAL:\n%s\n", initialBranches)
		// fmt.Printf("NOW:\n%s\n", currentBranches.String())
		diff, errorCount := currentBranches.EqualDataTable(state.initialBranches.GetOrPanic())
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing branches\n\n", errorCount)
			fmt.Println(diff)
			panic("mismatching branches found, see diff above")
		}
		// verify initial lineage
		currentLineage := devRepo.LineageTable()
		diff, errCnt := currentLineage.EqualDataTable(state.initialLineage.GetOrPanic())
		if errCnt > 0 {
			fmt.Printf("\nERROR! Found %d differences in the lineage\n\n", errCnt)
			fmt.Println(diff)
			panic("mismatching lineage found, see the diff above")
		}
	})

	sc.Step(`^the initial branches exist$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		have := state.fixture.Branches()
		want := state.initialBranches.GetOrPanic()
		// fmt.Printf("HAVE:\n%s\n", have.String())
		// fmt.Printf("WANT:\n%s\n", want.String())
		diff, errorCount := have.EqualDataTable(want)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing branches\n\n", errorCount)
			fmt.Println(diff)
			panic("mismatching branches found, see diff above")
		}
	})

	sc.Step(`^the initial commits exist$`, func(ctx context.Context) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		currentCommits := state.fixture.CommitTable(state.initialCommits.GetOrPanic().Cells[0])
		errDiff, errCount := state.initialCommits.GetOrPanic().EqualDataTable(currentCommits)
		if errCount == 0 {
			return
		}
		fmt.Println(errDiff)
		panic("current commits are not the same as the initial commits")
	})

	sc.Step(`^the main branch is "([^"]+)"$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		return devRepo.Config.SetMainBranch(gitdomain.NewLocalBranchName(name))
	})

	sc.Step(`^the main branch is (?:now|still) "([^"]*)"$`, func(ctx context.Context, want string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.Config.MainBranch
		if have.String() != want {
			return fmt.Errorf("expected %q, got %q", want, have)
		}
		return nil
	})

	sc.Step(`^the main branch is (?:now|still) not set$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		have := devRepo.Config.LocalGitConfig.MainBranch
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
		return devRepo.Config.SetPerennialBranches(gitdomain.NewLocalBranchNames(name))
	})

	sc.Step(`^the perennial branches are "([^"]+)" and "([^"]+)"$`, func(ctx context.Context, branch1, branch2 string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		return devRepo.Config.SetPerennialBranches(gitdomain.NewLocalBranchNames(branch1, branch2))
	})

	sc.Step(`^the perennial branches are not configured$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		return devRepo.RemovePerennialBranchConfiguration()
	})

	sc.Step(`^the perennial branches are (?:now|still) "([^"]+)"$`, func(ctx context.Context, name string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actual := devRepo.Config.LocalGitConfig.PerennialBranches
		if len(actual) != 1 {
			return fmt.Errorf("expected 1 perennial branch, got %q", actual)
		}
		if (actual)[0].String() != name {
			return fmt.Errorf("expected %q, got %q", name, (actual)[0])
		}
		return nil
	})

	sc.Step(`^the perennial branches are now "([^"]+)" and "([^"]+)"$`, func(ctx context.Context, branch1, branch2 string) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		actual := devRepo.Config.LocalGitConfig.PerennialBranches
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

	sc.Step(`^there are (?:now|still) no contribution branches$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branches := devRepo.Config.LocalGitConfig.ContributionBranches
		if len(branches) > 0 {
			return fmt.Errorf("expected no contribution branches, got %q", branches)
		}
		return nil
	})

	sc.Step(`^there are (?:now|still) no observed branches$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branches := devRepo.Config.LocalGitConfig.ObservedBranches
		if len(branches) > 0 {
			return fmt.Errorf("expected no observed branches, got %q", branches)
		}
		return nil
	})

	sc.Step(`^there are (?:now|still) no parked branches$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branches := devRepo.Config.LocalGitConfig.ParkedBranches
		if len(branches) > 0 {
			return fmt.Errorf("expected no parked branches, got %q", branches)
		}
		return nil
	})

	sc.Step(`^there are (?:now|still) no perennial branches$`, func(ctx context.Context) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		branches := devRepo.Config.LocalGitConfig.PerennialBranches
		if len(branches) > 0 {
			return fmt.Errorf("expected no perennial branches, got %q", branches)
		}
		return nil
	})

	sc.Step(`^these committed files exist now$`, func(ctx context.Context, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		fileTable := devRepo.FilesInBranches(gitdomain.NewLocalBranchName("main"))
		diff, errorCount := fileTable.EqualGherkin(table)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing files\n\n", errorCount)
			fmt.Println(diff)
			panic("mismatching files found, see diff above")
		}
	})

	sc.Step(`^these commits exist now$`, func(ctx context.Context, table *godog.Table) error {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		return state.compareGherkinTable(table)
	})

	sc.Step(`^these tags exist$`, func(ctx context.Context, table *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		tagTable := state.fixture.TagTable()
		diff, errorCount := tagTable.EqualGherkin(table)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing tags\n\n", errorCount)
			fmt.Println(diff)
			panic("mismatching tags found, see diff above")
		}
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
		if hasFile != "" {
			panic(hasFile)
		}
	})

	sc.Step(`^these branches exist now$`, func(ctx context.Context, input *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		currentBranches := state.fixture.Branches()
		// fmt.Printf("NOW:\n%s\n", currentBranches.String())
		diff, errorCount := currentBranches.EqualGherkin(input)
		if errorCount != 0 {
			fmt.Printf("\nERROR! Found %d differences in the existing branches\n\n", errorCount)
			fmt.Println(diff)
			panic("mismatching branches found, see diff above")
		}
	})

	sc.Step(`^this lineage exists now$`, func(ctx context.Context, input *godog.Table) {
		state := ctx.Value(keyScenarioState).(*ScenarioState)
		devRepo := state.fixture.DevRepo.GetOrPanic()
		table := devRepo.LineageTable()
		diff, errCount := table.EqualGherkin(input)
		if errCount > 0 {
			fmt.Printf("\nERROR! Found %d differences in the lineage\n\n", errCount)
			fmt.Println(diff)
			panic("mismatching branches found, see the diff above")
		}
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
}

func updateInitialSHAs(state *ScenarioState) {
	devRepo := state.fixture.DevRepo.GetOrPanic()
	if state.initialDevSHAs.IsNone() && state.insideGitRepo {
		state.initialDevSHAs = Some(devRepo.TestCommands.CommitSHAs())
	}
	if originRepo, hasOriginrepo := state.fixture.OriginRepo.Get(); state.initialOriginSHAs.IsNone() && state.insideGitRepo && hasOriginrepo {
		state.initialOriginSHAs = Some(originRepo.TestCommands.CommitSHAs())
	}
	if secondWorkTree, hasSecondWorkTree := state.fixture.SecondWorktree.Get(); state.initialWorktreeSHAs.IsNone() && state.insideGitRepo && hasSecondWorkTree {
		state.initialWorktreeSHAs = Some(secondWorkTree.TestCommands.CommitSHAs())
	}
}
