package fixture

import (
	"os"
	"path/filepath"

	"github.com/git-town/git-town/v14/test/filesystem"
	"github.com/git-town/git-town/v14/test/helpers"
)

// Factory manages the Git setup for the entire test suite.
// For each scenario, it provides a standardized, empty Fixture consisting of a local and remote Git repository.
//
// Setting up a Fixture is an expensive operation and has to be done for every scenario.
// As a performance optimization, Factory creates a fully set up Fixture (including the main branch and configuration)
// (the "memoized" environment) at the beginning of the test suite and makes copies of it for each scenario.
// Making copies of a fully set up Git repo is much faster than creating it from scratch.
// End-to-end tests run multi-threaded, all threads share a global Factory instance.
type Factory struct {
	counter helpers.AtomicCounter

	// path of the folder that this class operates in
	dir string

	// the memoized environment
	memoized Fixture
}

// NewFactory provides a new FixtureFactory instance operating in the given directory.
func NewFactory(dir string) Factory {
	return Factory{
		counter:  helpers.AtomicCounter{},
		dir:      dir,
		memoized: NewStandardFixture(filepath.Join(dir, "memoized")),
	}
}

// CreateFixture provides a new Fixture for the scenario with the given name.
func (self *Factory) CreateFixture(scenarioName string) Fixture {
	envDirName := filesystem.FolderName(scenarioName) + "_" + self.counter.ToString()
	envPath := filepath.Join(self.dir, envDirName)
	return CloneFixture(self.memoized, envPath)
}

func (self *Factory) Remove() {
	os.RemoveAll(self.dir)
}
