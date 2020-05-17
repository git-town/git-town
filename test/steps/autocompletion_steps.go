package steps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// AutocompletionSteps defines Cucumber step implementations around Git branches.
// nolint:funlen
func AutocompletionSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^I have an empty fish autocompletion folder$`, func() error {
		return os.MkdirAll(fishFolderPath(state), 0744)
	})

	suite.Step(`^I have an existing Git autocompletion file$`, func() error {
		err := os.MkdirAll(fishFolderPath(state), 0744)
		if err != nil {
			return fmt.Errorf("cannot create fish folder: %w", err)
		}
		return ioutil.WriteFile(fishFilePath(state), []byte("existing content"), 0744)
	})

	suite.Step(`^I have no fish autocompletion file$`, func() error {
		// nothing to do here, the test directory has no data
		return nil
	})

	suite.Step(`^I now have a Git autocompletion file$`, func() error {
		fishPath := filepath.Join(state.gitEnv.Dir, ".config", "fish", "completions", "git.fish")
		_, err := os.Stat(fishPath)
		if os.IsNotExist(err) {
			return err
		}
		return nil
	})
	suite.Step(`^I still have my original Git autocompletion file$`, func() error {
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
}

func fishFolderPath(state *ScenarioState) string {
	return filepath.Join(state.gitEnv.Dir, ".config", "fish", "completions")
}

func fishFilePath(state *ScenarioState) string {
	return filepath.Join(fishFolderPath(state), "git.fish")
}
