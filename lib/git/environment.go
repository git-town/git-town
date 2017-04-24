package git

import "github.com/Originate/git-town/lib/util"

func isRepository() bool {
	_, err := util.GetFullCommandOutput("git", "rev-parse")
	return err == nil
}

func EnsureIsRepository() {
	util.Ensure(isRepository(), "This is not a git repository.")
}
