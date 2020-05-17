package steps

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cucumber/godog"
)

// WorkspaceSteps defines Cucumber step implementations around Git workspace management.
func WorkspaceSteps(suite *godog.Suite, fs *FeatureState) {
	suite.Step(`^my workspace has an uncommitted file$`, func() error {
		fs.state.uncommittedFileName = "uncommitted file"
		fs.state.uncommittedContent = "uncommitted content"
		return fs.state.gitEnv.DevRepo.CreateFile(
			fs.state.uncommittedFileName,
			fs.state.uncommittedContent,
		)
	})

	suite.Step(`^my workspace has an uncommitted file in folder "([^"]*)"$`, func(folder string) error {
		fs.state.uncommittedFileName = fmt.Sprintf("%s/uncommitted file", folder)
		return fs.state.gitEnv.DevRepo.CreateFile(
			fs.state.uncommittedFileName,
			fs.state.uncommittedContent,
		)
	})

	suite.Step(`^my workspace has an uncommitted file with name: "([^"]+)" and content: "([^"]+)"$`, func(name, content string) error {
		fs.state.uncommittedFileName = name
		fs.state.uncommittedContent = content
		return fs.state.gitEnv.DevRepo.CreateFile(name, content)
	})

	suite.Step(`^my workspace has the uncommitted file again$`, func() error {
		hasFile, err := fs.state.gitEnv.DevRepo.HasFile(
			fs.state.uncommittedFileName,
			fs.state.uncommittedContent,
		)
		if err != nil {
			return err
		}
		if !hasFile {
			return fmt.Errorf("expected file %q but didn't find it", fs.state.uncommittedFileName)
		}
		return nil
	})

	suite.Step(`^my workspace is currently not a Git repository$`, func() error {
		os.RemoveAll(filepath.Join(fs.state.gitEnv.DevRepo.Dir, ".git"))
		return nil
	})

	suite.Step(`^my workspace still contains my uncommitted file$`, func() error {
		hasFile, err := fs.state.gitEnv.DevRepo.HasFile(
			fs.state.uncommittedFileName,
			fs.state.uncommittedContent,
		)
		if err != nil {
			return fmt.Errorf("cannot determine if workspace contains uncommitted file: %w", err)
		}
		if !hasFile {
			return fmt.Errorf("expected the uncommitted file but didn't find one")
		}
		return nil
	})
}
