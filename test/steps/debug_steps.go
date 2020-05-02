package steps

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cucumber/godog"
)

// DebugSteps defines Gherkin step implementations around merge conflicts.
func DebugSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^inspect the repo$`, func() error {
		fmt.Println(fs.activeScenarioState.gitEnvironment.DeveloperRepo.Dir)
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		return nil
	})
}
