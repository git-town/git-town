package steps

import "github.com/cucumber/godog"

// OriginSteps defines Cucumber step implementations around Git origins.
func OriginSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my repo does not have a remote origin$`, func() error {
		err := fs.state.gitEnv.DevRepo.RemoveRemote("origin")
		if err != nil {
			return err
		}
		fs.state.gitEnv.OriginRepo = nil
		return nil
	})

	suite.Step(`^my repo's origin is "([^"]*)"$`, func(origin string) error {
		fs.state.gitEnv.DevShell.SetTestOrigin(origin)
		return nil
	})
}
