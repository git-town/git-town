package steps

import (
	"bufio"
	"fmt"
	"os"

	"github.com/cucumber/godog"
)

// DebugSteps implements Gherkin steps that help debug the test setup.
func DebugSteps(suite *godog.Suite, state *ScenarioState) {
	suite.Step(`^inspect the repo$`, func() error {
		fmt.Println(state.gitEnv.DevRepo.Dir)
		reader := bufio.NewReader(os.Stdin)
		_, _ = reader.ReadString('\n')
		return nil
	})
}
