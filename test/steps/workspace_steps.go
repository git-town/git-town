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
		fs.uncommittedFileName = "uncommitted file"
		fs.uncommittedContent = "uncommitted content"
		return fs.gitEnv.DevRepo.CreateFile(
			fs.uncommittedFileName,
			fs.uncommittedContent,
		)
	})

	suite.Step(`^my workspace has an uncommitted file in folder "([^"]*)"$`, func(folder string) error {
		fs.uncommittedFileName = fmt.Sprintf("%s/uncommitted file", folder)
		return fs.gitEnv.DevRepo.CreateFile(
			fs.uncommittedFileName,
			fs.uncommittedContent,
		)
	})

	suite.Step(`^my workspace has an uncommitted file with name: "([^"]+)" and content: "([^"]+)"$`, func(name, content string) error {
		fs.uncommittedFileName = name
		fs.uncommittedContent = content
		return fs.gitEnv.DevRepo.CreateFile(name, content)
	})

	suite.Step(`^my workspace has the uncommitted file again$`, func() error {
		hasFile, err := fs.gitEnv.DevRepo.HasFile(
			fs.uncommittedFileName,
			fs.uncommittedContent,
		)
		if err != nil {
			return err
		}
		if !hasFile {
			return fmt.Errorf("expected file %q but didn't find it", fs.uncommittedFileName)
		}
		return nil
	})

	suite.Step(`^my workspace is currently not a Git repository$`, func() error {
		os.RemoveAll(filepath.Join(fs.gitEnv.DevRepo.Dir, ".git"))
		return nil
	})

	suite.Step(`^my workspace still contains my uncommitted file$`, func() error {
		hasFile, err := fs.gitEnv.DevRepo.HasFile(
			fs.uncommittedFileName,
			fs.uncommittedContent,
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
